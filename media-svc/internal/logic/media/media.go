package media

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/media-svc/api/v1"
	"github.com/TsingpekTao/shopa/media-svc/internal/consts"
	"github.com/TsingpekTao/shopa/media-svc/internal/dao"
	"github.com/TsingpekTao/shopa/media-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/media-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/media-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 领域模型速览（仅注释说明，不参与运行）：
// 1) 资产树：
//   - ORIGINAL 资产：初始化上传创建，parent_asset_id=0。
//   - DERIVED 资产：由父资产派生，parent_asset_id 指向父节点。
//   - root_asset_id：整棵派生树统一根，用于状态聚合、链路展示与 GC 候选归并。
//
// 2) 状态机（核心在 process_status / risk_status）：
//   - InitUpload：先落库为 pending（或按策略直接 success）。
//   - CompleteUpload：上传完成后从 pending 推进到 processing/success。
//   - UpdateAssetProcessStatus：异步任务回写终态与错误信息。
//
// 3) 绑定关系：
//   - media_asset_binding 采用“新增记录 + 失活旧记录”模式，保留审计历史。
//   - is_active=0 + unbound_at!=nil 代表解绑完成，不做物理删除。
//
// 4) URL 发放：
//   - 公有读：返回公共 URL。
//   - 私有读：返回带过期时间与签名参数的时效 URL。
//
// 5) GC 语义：
//   - 本文件不直接删除对象存储实体，只提供 storage_object_hash/root_asset_id 等归并线索，
//     供离线/异步 GC 任务判断“对象是否仍被引用”。
//
// sMedia 承载媒体域逻辑，实现外部与内部 RPC 的统一服务入口。
type sMedia struct{}

// New 创建媒体逻辑服务实例。
func New() *sMedia {
	return &sMedia{}
}

func init() {
	// 进程启动时注册实现到 service 层，供 RPC 框架按接口分发调用。
	service.RegisterMedia(New())
}

func (s *sMedia) InitUpload(ctx context.Context, req *v1.InitUploadReq) (*v1.InitUploadRes, error) {
	// 1) 基础入参校验：场景、文件名、MIME 是上传初始化最小必填集合。
	// sceneCode 将参与对象路径生成、场景策略匹配、后续绑定校验。
	sceneCode := strings.TrimSpace(req.GetSceneCode())
	if sceneCode == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code is required")
	}
	// file_name 允许业务原名透传，后续会在 buildObjectKey 中做安全字符替换。
	if strings.TrimSpace(req.GetFileName()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "file_name is required")
	}
	// mime_type 用于策略白名单校验和后续回传 UploadTicket 的 Content-Type。
	if strings.TrimSpace(req.GetMimeType()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "mime_type is required")
	}

	// 2) 读取场景策略并做前置拦截：大小、MIME 任一不满足都直接拒绝。
	policy, err := s.getScenePolicy(ctx, sceneCode)
	if err != nil {
		return nil, err
	}
	// max_size_bytes=0 视为不限制；>0 时执行硬限制。
	if policy.MaxSizeBytes > 0 && req.GetSizeBytes() > policy.MaxSizeBytes {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "size exceeds scene max_size_bytes")
	}
	// allow_mime_json 非空且可解析时采用白名单精确匹配。
	if !mimeAllowed(policy.AllowMimeJson, req.GetMimeType()) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "mime_type is not allowed by scene policy")
	}

	// 3) 组装资源初始化状态（逐字段解释变量用途）：
	// - now：本次初始化统一时间基准，避免多次取 now 带来的毫秒级漂移。
	// - storageProvider/bucket：对象存储定位维度，参与 hashStorageObject 计算。
	// - objectKey：对象键（资产树叶子节点的物理定位）。
	// - publicURL：仅公开读资产预计算，私有资产留空等待签发。
	// - riskStatus/processStatus/processProgress：上传后状态机起点。
	// - processedAt：仅在“处理同步完成”分支写入。
	var (
		now             = gtime.Now()                                  // 统一时间戳基准（含时区/纳秒）。
		storageProvider = cfgString(ctx, "storage.provider", "minio")  // 存储提供方标识（如 minio/s3）。
		bucket          = cfgString(ctx, "storage.bucket", "shopa")    // 存储桶名，后续参与 URL 与对象 hash。
		objectKey       = buildObjectKey(sceneCode, req.GetFileName()) // 生成安全对象键。
		publicURL       = ""                                           // 公开读场景填充；私有读保持空。
		riskStatus      = uint(consts.RiskStatusPending)               // 风控默认 pending。
		processStatus   = uint(consts.ProcessStatusPending)            // 处理默认 pending。
		processProgress = uint(0)                                      // 初始进度 0。
		processedAt     *gtime.Time                                    // 默认空，表示尚未完成处理。
	)

	// ACL 为公开读时可直接生成公网地址，后续读 URL 可直接复用。
	if uint(policy.AclType) == consts.ACLTypePublicRead {
		publicURL = buildPublicURL(ctx, bucket, objectKey)
	}
	// 风控状态机：risk_async=0 时，不需要异步审核，直接标记通过。
	if policy.RiskAsyncEnabled == 0 {
		riskStatus = uint(consts.RiskStatusPassed)
	}
	// 处理状态机：process_async=0 时，初始化即进入 success(100%)。
	if policy.ProcessAsyncEnabled == 0 {
		processStatus = uint(consts.ProcessStatusSuccess)
		processProgress = 100
		processedAt = now
	}

	// 透传业务追踪字段（非绑定关系）：用于链路审计、问题回溯、离线统计归因。
	traceBizType := req.GetBizType()
	traceBizNo := req.GetBizNo()
	// 从 metadata 提取上传人 ID；解析失败不报错，落 0 代表系统/匿名上下文。
	uploaderID, _ := userIDFromMetadata(ctx, "x-user-id", "user_id", "uid")

	// 4) 写入 media_asset 主记录。
	// 说明：StorageObjectHash = SHA256(provider|bucket|objectKey)：
	// - 语义上是“对象物理定位指纹”，不是内容哈希；
	// - 可用于同对象多引用归并、GC 候选聚合、异常排查比对。
	insertRes, err := dao.MediaAsset.Ctx(ctx).Data(do.MediaAsset{
		SceneCode:         sceneCode,                                             // 场景码：策略/绑定都以此隔离。
		AclType:           uint(policy.AclType),                                  // 访问控制：公开读或私有读。
		AssetRole:         consts.AssetRoleOriginal,                              // 初始化上传恒为 ORIGINAL。
		DerivedKind:       consts.DerivedKindUnspecified,                         // 原始资产无派生类型。
		FileName:          req.GetFileName(),                                     // 业务展示名。
		MimeType:          req.GetMimeType(),                                     // 媒体 MIME。
		SizeBytes:         req.GetSizeBytes(),                                    // 媒体字节数。
		ChecksumSha256:    req.GetChecksumSha256(),                               // 可选内容校验摘要（调用方提供）。
		StorageProvider:   storageProvider,                                       // 存储提供方。
		Bucket:            bucket,                                                // 存储桶。
		ObjectKey:         objectKey,                                             // 对象键。
		StorageObjectHash: hashStorageObject(storageProvider, bucket, objectKey), // 对象定位指纹。
		PublicUrl:         publicURL,                                             // 公开读 URL（私有读为空）。
		RiskStatus:        riskStatus,                                            // 风控状态起点。
		ProcessStatus:     processStatus,                                         // 处理状态起点。
		ProcessProgress:   processProgress,                                       // 处理进度起点。
		ProcessedAt:       processedAt,                                           // 处理完成时间（仅同步处理场景）。
		UploaderUserId:    uploaderID,                                            // 上传人。
		TraceBizType:      traceBizType,                                          // 追踪业务类型。
		TraceBizNo:        traceBizNo,                                            // 追踪业务编号。
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "insert media_asset failed")
	}
	// 自增主键作为 asset_id，对外作为全局资产标识。
	lastID, err := insertRes.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "read inserted asset id failed")
	}
	assetID := uint64(lastID)

	// 5) 资产树根节点回写：
	// ORIGINAL 自己就是根，因此 root_asset_id=asset_id。
	// 后续 DERIVED 资产统一挂到这个 root，下游查询/GC 以 root 为聚合维度。
	_, err = dao.MediaAsset.Ctx(ctx).
		Where(dao.MediaAsset.Columns().AssetId, assetID).
		Data(do.MediaAsset{RootAssetId: assetID}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "set root_asset_id failed")
	}

	// 回读最新记录，确保响应体包含数据库最终落盘值。
	asset, err := s.getAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	// 6) 返回上传凭证：客户端直接 PUT 到对象存储，服务端不接收大文件流。
	uploadURL := buildUploadURL(ctx, bucket, objectKey)
	return &v1.InitUploadRes{
		Asset: toProtoAsset(asset),
		UploadTicket: &v1.UploadTicket{
			Method:    "PUT",     // 约定直接上传方式。
			UploadUrl: uploadURL, // 上传目标地址（可能是网关或对象存储 endpoint）。
			Headers: map[string]string{
				"Content-Type": req.GetMimeType(), // 上传时显式声明 MIME，便于存储层记录。
			},
			ExpiredAt: timestamppb.New(time.Now().Add(15 * time.Minute)), // 凭证有效期。
		},
	}, nil
}

func (s *sMedia) CompleteUpload(ctx context.Context, req *v1.CompleteUploadReq) (*v1.CompleteUploadRes, error) {
	// 上传完成回调：
	// - 补齐对象存储返回的元信息（etag/size/checksum 等）；
	// - 将 process_status 从初始化态推进到 processing/success。
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	// 先读取资产，确保 scene_code 与软删除状态合法。
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}
	// 场景策略决定“上传完成后”状态机走向（同步完成 or 异步处理中）。
	policy, err := s.getScenePolicy(ctx, asset.SceneCode)
	if err != nil {
		return nil, err
	}

	// data 只更新“完成上传后可确认”的字段，避免覆盖初始化已写入字段。
	data := do.MediaAsset{
		Etag:           req.GetEtag(),           // 对象存储返回的实体标签。
		MimeType:       req.GetMimeType(),       // 以完成上传时的值为准。
		SizeBytes:      req.GetSizeBytes(),      // 实际对象大小。
		ChecksumSha256: req.GetChecksumSha256(), // 调用方确认的摘要。
	}
	// process 状态机推进：
	// - process_async=0：直接进入 success(100%) 并写 processed_at；
	// - process_async=1：进入 processing(0%)，等待异步 worker 回写。
	if policy.ProcessAsyncEnabled == 0 {
		data.ProcessStatus = consts.ProcessStatusSuccess
		data.ProcessProgress = 100
		data.ProcessedAt = gtime.Now()
	} else {
		data.ProcessStatus = consts.ProcessStatusProcessing
		data.ProcessProgress = 0
	}
	// risk 状态机推进：risk_async=0 时，上传完成即可判定通过。
	if policy.RiskAsyncEnabled == 0 {
		data.RiskStatus = consts.RiskStatusPassed
	}

	// 仅更新未软删记录，避免误改历史资产（软删是 GC 判断的重要维度之一）。
	_, err = dao.MediaAsset.Ctx(ctx).
		Where(dao.MediaAsset.Columns().AssetId, req.GetAssetId()).
		WhereNull(dao.MediaAsset.Columns().DeletedAt).
		Data(data).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update asset complete status failed")
	}
	updated, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}
	// 返回更新后的完整资产，便于调用方立即观察状态机当前位置。
	return &v1.CompleteUploadRes{Asset: toProtoAsset(updated)}, nil
}

func (s *sMedia) IssueReadUrl(ctx context.Context, req *v1.IssueReadUrlReq) (*v1.IssueReadUrlRes, error) {
	// 读链接签发：
	// - 公有 ACL：直接返回公共访问地址；
	// - 私有 ACL：签发短期有效 URL，超时后自动失效。
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}

	// 公有读路径：
	// 1) 优先复用库里已保存的 public_url（避免配置变更引起抖动）；
	// 2) 若为空，则按当前配置实时拼装兜底。
	if asset.AclType == consts.ACLTypePublicRead {
		url := asset.PublicUrl
		if strings.TrimSpace(url) == "" {
			url = buildPublicURL(ctx, asset.Bucket, asset.ObjectKey)
		}
		return &v1.IssueReadUrlRes{
			AssetId:  asset.AssetId,
			Url:      url,
			IsPublic: true,
		}, nil
	}

	// 私有读路径：TTL 上下界保护，防止调用方传入极端值。
	ttl := int(req.GetTtlSeconds())
	if ttl <= 0 {
		// 未传或非法值时，使用系统默认有效期。
		ttl = consts.DefaultIssueReadTTLSeconds
	}
	if ttl > consts.MaxIssueReadTTLSeconds {
		// 防止签发过长有效期 URL 导致越权窗口扩大。
		ttl = consts.MaxIssueReadTTLSeconds
	}
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
	// 用资产身份 + 到期时间生成签名，URL 天然带时效语义。
	signedURL := buildSignedReadURL(ctx, asset, expireAt.Unix())
	return &v1.IssueReadUrlRes{
		AssetId:   asset.AssetId,             // 对应资产 ID。
		Url:       signedURL,                 // 私有签名 URL。
		IsPublic:  false,                     // 明确该链接非永久公开。
		ExpiredAt: timestamppb.New(expireAt), // 明确回传过期时间，便于调用方缓存控制。
	}, nil
}

func (s *sMedia) GetAssetProcessStatus(ctx context.Context, req *v1.GetAssetProcessStatusReq) (*v1.GetAssetProcessStatusRes, error) {
	// 查询单资产处理状态，并按资产树根节点聚合同链路 DERIVED 资产摘要。
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}
	rootID := asset.RootAssetId
	// 兼容历史数据：root_asset_id 为空时退化为自身 ID，保证“树根可计算”。
	if rootID == 0 {
		rootID = asset.AssetId
	}

	var derivedRows []*entity.MediaAsset
	cols := dao.MediaAsset.Columns()
	// 只取同根 + 派生角色 + 未软删记录：
	// - 同根：确保聚合范围属于同一资产树；
	// - 派生角色：不重复返回原始节点；
	// - 未软删：避免把已逻辑删除节点展示给上层。
	if err = dao.MediaAsset.Ctx(ctx).
		Where(cols.RootAssetId, rootID).
		Where(cols.AssetRole, consts.AssetRoleDerived).
		WhereNull(cols.DeletedAt).
		OrderAsc(cols.AssetId).
		Scan(&derivedRows); err != nil {
		return nil, gerror.Wrap(err, "query derived assets failed")
	}

	// 转为轻量摘要结构，既表达链路状态，又避免泄露无关底层字段。
	derived := make([]*v1.DerivedAssetSummary, 0, len(derivedRows))
	for _, row := range derivedRows {
		derived = append(derived, &v1.DerivedAssetSummary{
			AssetId:       row.AssetId,
			ParentAssetId: row.ParentAssetId,
			DerivedKind:   v1.DerivedKind(row.DerivedKind),
			MimeType:      row.MimeType,
			PublicUrl:     row.PublicUrl,
			ProcessStatus: v1.ProcessStatus(row.ProcessStatus),
		})
	}

	return &v1.GetAssetProcessStatusRes{
		AssetId:             asset.AssetId,
		ProcessStatus:       v1.ProcessStatus(asset.ProcessStatus),
		ProcessProgress:     uint32(asset.ProcessProgress),
		ProcessErrorCode:    asset.ProcessErrorCode,
		ProcessErrorMessage: asset.ProcessErrorMessage,
		DerivedAssets:       derived,
	}, nil
}

func (s *sMedia) BatchGetAssetProcessStatus(ctx context.Context, req *v1.BatchGetAssetProcessStatusReq) (*v1.BatchGetAssetProcessStatusRes, error) {
	// 顺序复用单查逻辑，保证单查与批量查行为一致。
	out := make([]*v1.GetAssetProcessStatusRes, 0, len(req.GetAssetIds()))
	for _, assetID := range req.GetAssetIds() {
		item, err := s.GetAssetProcessStatus(ctx, &v1.GetAssetProcessStatusReq{AssetId: assetID})
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return &v1.BatchGetAssetProcessStatusRes{Statuses: out}, nil
}

func (s *sMedia) BatchBindAssetsToBiz(ctx context.Context, req *v1.BatchBindAssetsToBizReq) (*v1.BatchBindAssetsToBizRes, error) {
	// 批量绑定入口（绑定语义）：
	// - 一个业务对象可有多个 binding_field（如 cover/gallery/detail）；
	// - 每条绑定记录是“业务对象 <-> 资产”的一条历史事实；
	// - 当前有效绑定由 is_active=1 表示，解绑通过失活实现。
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no are required")
	}
	if len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items is required")
	}
	policy, err := s.getScenePolicy(ctx, req.GetSceneCode())
	if err != nil {
		return nil, err
	}
	if policy.MaxCount > 0 {
		// max_count 对每个 binding_field 独立生效，不是整个请求全局总数。
		fieldCnt := map[string]int{}
		for _, item := range req.GetItems() {
			field := strings.TrimSpace(item.GetBindingField())
			if field == "" {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter, "binding_field is required")
			}
			fieldCnt[field]++
		}
		for _, cnt := range fieldCnt {
			if cnt > int(policy.MaxCount) {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items exceed scene max_count")
			}
		}
	}

	// 收集并校验资产集合：必须存在、同场景、且仅 ORIGINAL 可被业务直接绑定。
	// （DERIVED 资产用于加工链路展示，不直接作为业务主资产。）
	assetIDs := make([]uint64, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item.GetAssetId() == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "item.asset_id is required")
		}
		assetIDs = append(assetIDs, item.GetAssetId())
	}
	assetMap, err := s.validateBindableAssets(ctx, req.GetSceneCode(), assetIDs)
	if err != nil {
		return nil, err
	}

	// 审计字段：操作人 + 请求号，便于追踪“是谁在何次请求里调整了绑定”。
	operatorID, _ := userIDFromMetadata(ctx, "x-user-id", "user_id", "uid")
	requestID := metadataValue(ctx, "x-request-id")
	now := gtime.Now() // 统一解绑时间戳，保证同事务内时间一致。

	// 事务语义：对每个 item 执行“幂等失活 -> 插入新激活记录”。
	// 这样做的好处：
	// - 保留历史（用于审计和回放）；
	// - 对重复请求天然幂等；
	// - 一旦失败整体回滚，不会出现半成功绑定。
	err = dao.MediaAssetBinding.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.MediaAssetBinding.Columns()
		for _, item := range req.GetItems() {
			asset := assetMap[item.GetAssetId()]
			if asset == nil {
				return gerror.NewCode(gcode.CodeNotFound, fmt.Sprintf("asset %d not found", item.GetAssetId()))
			}
			field := strings.TrimSpace(item.GetBindingField())
			// 对同 key（scene/biz/field/asset）的活跃记录先失活：
			// - 若原本不存在活跃记录，更新 0 行也视为正常（幂等）；
			// - 若存在活跃记录，先打上 unbound_at 再插入新版本。
			_, _ = tx.Model(dao.MediaAssetBinding.Table()).
				Where(cols.SceneCode, req.GetSceneCode()).
				Where(cols.BizType, req.GetBizType()).
				Where(cols.BizNo, req.GetBizNo()).
				Where(cols.BindingField, field).
				Where(cols.AssetId, item.GetAssetId()).
				Where(cols.IsActive, consts.BindingActive).
				Data(do.MediaAssetBinding{IsActive: consts.BindingInactive, UnboundAt: now}).
				Update()

			// 插入新活跃绑定版本；唯一约束冲突统一转业务错误码，便于前端提示。
			_, insertErr := tx.Model(dao.MediaAssetBinding.Table()).Data(do.MediaAssetBinding{
				SceneCode:      req.GetSceneCode(),   // 场景隔离维度。
				BizType:        req.GetBizType(),     // 业务类型。
				BizNo:          req.GetBizNo(),       // 业务主键。
				BindingField:   field,                // 业务字段（如 cover）。
				AssetId:        item.GetAssetId(),    // 绑定资产 ID。
				SortOrder:      item.GetSortOrder(),  // 字段内展示顺序。
				IsActive:       consts.BindingActive, // 新版本置为活跃。
				OperatorUserId: operatorID,           // 审计：操作人。
				RequestId:      requestID,            // 审计：请求号。
			}).Insert()
			if insertErr != nil {
				if isDuplicateErr(insertErr) {
					return gerror.NewCode(gcode.CodeBusinessValidationFailed, "duplicate active binding or sort_order conflict")
				}
				return gerror.Wrap(insertErr, "insert media_asset_binding failed")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 绑定完成后返回该业务对象当前活跃绑定快照。
	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), "", false)
	if err != nil {
		return nil, err
	}
	return &v1.BatchBindAssetsToBizRes{Bindings: toProtoBindingList(bindings)}, nil
}

func (s *sMedia) ReplaceBizAssetBindings(ctx context.Context, req *v1.ReplaceBizAssetBindingsReq) (*v1.ReplaceBizAssetBindingsRes, error) {
	// 字段级整体替换：
	// - 目标是“这个 binding_field 现在应该精确等于 req.items”；
	// - 实现策略是先批量失活旧记录，再插入新集合。
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" || strings.TrimSpace(req.GetBindingField()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no/binding_field are required")
	}

	policy, err := s.getScenePolicy(ctx, req.GetSceneCode())
	if err != nil {
		return nil, err
	}
	if policy.MaxCount > 0 && len(req.GetItems()) > int(policy.MaxCount) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items exceed scene max_count")
	}

	// 先做资产合法性校验，避免事务做到一半才因脏数据失败。
	assetIDs := make([]uint64, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item.GetAssetId() == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "item.asset_id is required")
		}
		assetIDs = append(assetIDs, item.GetAssetId())
	}
	_, err = s.validateBindableAssets(ctx, req.GetSceneCode(), assetIDs)
	if err != nil {
		return nil, err
	}

	operatorID, _ := userIDFromMetadata(ctx, "x-user-id", "user_id", "uid")
	requestID := metadataValue(ctx, "x-request-id")
	now := gtime.Now()

	// replace 同样在事务中完成，确保“清空旧集合 + 写入新集合”原子可见。
	err = dao.MediaAssetBinding.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.MediaAssetBinding.Columns()
		// 第一步：批量失活目标字段下现有活跃绑定（逻辑解绑，不删历史）。
		_, updErr := tx.Model(dao.MediaAssetBinding.Table()).
			Where(cols.SceneCode, req.GetSceneCode()).
			Where(cols.BizType, req.GetBizType()).
			Where(cols.BizNo, req.GetBizNo()).
			Where(cols.BindingField, req.GetBindingField()).
			Where(cols.IsActive, consts.BindingActive).
			Data(do.MediaAssetBinding{IsActive: consts.BindingInactive, UnboundAt: now}).
			Update()
		if updErr != nil {
			return gerror.Wrap(updErr, "deactivate old bindings failed")
		}

		// 第二步：插入新集合，形成新的活跃版本快照。
		for _, item := range req.GetItems() {
			_, insertErr := tx.Model(dao.MediaAssetBinding.Table()).Data(do.MediaAssetBinding{
				SceneCode:      req.GetSceneCode(),    // 场景维度。
				BizType:        req.GetBizType(),      // 业务类型。
				BizNo:          req.GetBizNo(),        // 业务编号。
				BindingField:   req.GetBindingField(), // 被替换字段。
				AssetId:        item.GetAssetId(),     // 新绑定资产。
				SortOrder:      item.GetSortOrder(),   // 新顺序。
				IsActive:       consts.BindingActive,  // 新记录激活。
				OperatorUserId: operatorID,            // 审计：操作人。
				RequestId:      requestID,             // 审计：请求 ID。
			}).Insert()
			if insertErr != nil {
				if isDuplicateErr(insertErr) {
					return gerror.NewCode(gcode.CodeBusinessValidationFailed, "duplicate active binding or sort_order conflict")
				}
				return gerror.Wrap(insertErr, "insert replaced binding failed")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 返回替换后该字段的最新活跃绑定列表。
	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), req.GetBindingField(), false)
	if err != nil {
		return nil, err
	}
	return &v1.ReplaceBizAssetBindingsRes{Bindings: toProtoBindingList(bindings)}, nil
}

func (s *sMedia) BatchUnbindAssetsFromBiz(ctx context.Context, req *v1.BatchUnbindAssetsFromBizReq) (*v1.BatchUnbindAssetsFromBizRes, error) {
	// 批量解绑：
	// - 不删记录，只把 is_active 改为 0 并写 unbound_at；
	// - 支持解绑整个字段，或仅解绑字段下的指定资产子集。
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" || strings.TrimSpace(req.GetBindingField()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no/binding_field are required")
	}
	cols := dao.MediaAssetBinding.Columns()
	// 先圈定“目标业务对象 + 字段 + 当前活跃记录”的解绑范围。
	model := dao.MediaAssetBinding.Ctx(ctx).
		Where(cols.SceneCode, req.GetSceneCode()).
		Where(cols.BizType, req.GetBizType()).
		Where(cols.BizNo, req.GetBizNo()).
		Where(cols.BindingField, req.GetBindingField()).
		Where(cols.IsActive, consts.BindingActive)
	if len(req.GetAssetIds()) > 0 {
		// 指定 asset_ids 时，只解绑这些资产，不影响同字段其他资产。
		model = model.WhereIn(cols.AssetId, toInterfaceSliceUint64(req.GetAssetIds()))
	}
	// 逻辑解绑：写失活标记与解绑时间，供审计与回溯使用。
	ret, err := model.Data(do.MediaAssetBinding{
		IsActive:  consts.BindingInactive, // 置为失活。
		UnboundAt: gtime.Now(),            // 记录解绑发生时间。
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "unbind assets failed")
	}
	rows, _ := ret.RowsAffected()
	return &v1.BatchUnbindAssetsFromBizRes{AffectedRows: uint64(rows)}, nil
}

func (s *sMedia) GetBizAssets(ctx context.Context, req *v1.GetBizAssetsReq) (*v1.GetBizAssetsRes, error) {
	// 查询业务资产视图：
	// - 先查绑定关系；
	// - 再按绑定中的 asset_id 批量补齐资产详情；
	// - 最终按绑定排序输出，保证前端稳定渲染。
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no are required")
	}
	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), req.GetBindingField(), req.GetIncludeInactive())
	if err != nil {
		return nil, err
	}
	// 从绑定中抽 asset_id，避免逐条查库（N+1）。
	assetIDs := make([]uint64, 0, len(bindings))
	for _, row := range bindings {
		assetIDs = append(assetIDs, row.AssetId)
	}
	assetMap, err := s.queryAssetsByIDs(ctx, assetIDs)
	if err != nil {
		return nil, err
	}

	out := make([]*v1.BizAsset, 0, len(bindings))
	for _, row := range bindings {
		item := &v1.BizAsset{
			Binding: toProtoBinding(row), // 绑定侧字段始终返回。
		}
		if asset := assetMap[row.AssetId]; asset != nil {
			// 资产可能因软删/权限等原因查不到，查到时再补详情。
			item.Asset = toProtoAsset(asset)
		}
		out = append(out, item)
	}
	return &v1.GetBizAssetsRes{Assets: out}, nil
}

func (s *sMedia) CreateDerivedAsset(ctx context.Context, req *v1.CreateDerivedAssetReq) (*v1.CreateDerivedAssetRes, error) {
	// 创建派生资产（DERIVED）：
	// - 典型场景：缩略图、转码产物、水印图；
	// - 资产树语义：parent_asset_id 指向直接父节点，root_asset_id 保持整树一致；
	// - GC 语义：派生节点与原始节点共享同一 root，便于整链路引用判断。
	if req.GetParentAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "parent_asset_id is required")
	}
	if req.GetDerivedKind() == v1.DerivedKind_DERIVED_KIND_UNSPECIFIED {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "derived_kind is required")
	}
	parent, err := s.getAssetByID(ctx, req.GetParentAssetId())
	if err != nil {
		return nil, err
	}
	// 场景码可显式覆盖；未传则继承父资产场景，保持策略一致性。
	sceneCode := strings.TrimSpace(req.GetSceneCode())
	if sceneCode == "" {
		sceneCode = parent.SceneCode
	}

	// 存储参数支持覆盖：
	// - 传值：按调用方指定落盘；
	// - 不传：优先继承父资产，再回退系统默认值。
	storageProvider := strings.TrimSpace(req.GetStorageProvider())
	if storageProvider == "" {
		storageProvider = cfgString(ctx, "storage.provider", "minio")
	}
	bucket := strings.TrimSpace(req.GetBucket())
	if bucket == "" {
		if strings.TrimSpace(parent.Bucket) != "" {
			bucket = parent.Bucket
		} else {
			bucket = cfgString(ctx, "storage.bucket", "shopa")
		}
	}
	objectKey := strings.TrimSpace(req.GetObjectKey())
	if objectKey == "" {
		objectKey = buildObjectKey(sceneCode, req.GetFileName())
	}
	publicURL := strings.TrimSpace(req.GetPublicUrl())
	if publicURL == "" && parent.AclType == consts.ACLTypePublicRead {
		// 父资产公开读时，若未显式传 public_url，则按当前配置生成。
		publicURL = buildPublicURL(ctx, bucket, objectKey)
	}

	// root_asset_id 统一归并到根资源：
	// - 父节点已有 root 时直接复用；
	// - 父节点无 root（兼容历史数据）则父节点自己为 root。
	rootID := parent.RootAssetId
	if rootID == 0 {
		rootID = parent.AssetId
	}
	// 派生资源默认写入 success(100%)：
	// 约定由上游处理任务在“产物可用”后再调用本接口。
	insertRes, err := dao.MediaAsset.Ctx(ctx).Data(do.MediaAsset{
		ParentAssetId:     parent.AssetId,                                        // 直接父节点。
		RootAssetId:       rootID,                                                // 树根节点（整链路统一）。
		SceneCode:         sceneCode,                                             // 场景码。
		AclType:           parent.AclType,                                        // ACL 继承父资产。
		AssetRole:         consts.AssetRoleDerived,                               // 角色明确为 DERIVED。
		DerivedKind:       uint(req.GetDerivedKind()),                            // 派生类型（缩略图/转码等）。
		FileName:          req.GetFileName(),                                     // 产物文件名。
		MimeType:          req.GetMimeType(),                                     // 产物 MIME。
		SizeBytes:         req.GetSizeBytes(),                                    // 产物大小。
		ChecksumSha256:    req.GetChecksumSha256(),                               // 产物摘要。
		Etag:              req.GetEtag(),                                         // 存储 etag。
		StorageProvider:   storageProvider,                                       // 存储提供方。
		Bucket:            bucket,                                                // 存储桶。
		ObjectKey:         objectKey,                                             // 对象键。
		StorageObjectHash: hashStorageObject(storageProvider, bucket, objectKey), // 对象定位哈希。
		PublicUrl:         publicURL,                                             // 公有 URL（若有）。
		RiskStatus:        parent.RiskStatus,                                     // 风控状态继承父资产。
		ProcessStatus:     consts.ProcessStatusSuccess,                           // 派生结果默认已处理成功。
		ProcessProgress:   100,                                                   // 处理进度 100。
		ProcessedAt:       gtime.Now(),                                           // 派生完成时间。
		UploaderUserId:    parent.UploaderUserId,                                 // 上传人继承。
		TraceBizType:      parent.TraceBizType,                                   // 业务追踪继承。
		TraceBizNo:        parent.TraceBizNo,                                     // 业务追踪继承。
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "insert derived asset failed")
	}
	lastID, err := insertRes.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "read derived asset id failed")
	}
	// 回读保证响应内容与 DB 一致（含默认值/触发器可能写入字段）。
	derived, err := s.getAssetByID(ctx, uint64(lastID))
	if err != nil {
		return nil, err
	}
	return &v1.CreateDerivedAssetRes{Asset: toProtoAsset(derived)}, nil
}

func (s *sMedia) UpdateAssetProcessStatus(ctx context.Context, req *v1.UpdateAssetProcessStatusReq) (*v1.UpdateAssetProcessStatusRes, error) {
	// 异步处理回写入口：
	// - worker 可更新 process_status/process_progress/error 信息；
	// - 终态时可写 processed_at，未传则由服务端补当前时间。
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	data := do.MediaAsset{
		ProcessStatus:       uint(req.GetProcessStatus()), // 新状态（processing/success/failed...）。
		ProcessProgress:     req.GetProcessProgress(),     // 新进度（调用方保证区间合法）。
		ProcessErrorCode:    req.GetProcessErrorCode(),    // 业务错误码。
		ProcessErrorMessage: req.GetProcessErrorMessage(), // 错误文本。
	}
	// 优先使用调用方明确传入的 processed_at。
	if req.GetProcessedAt() != nil {
		data.ProcessedAt = gtime.NewFromTime(req.GetProcessedAt().AsTime())
	} else if req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_SUCCESS ||
		req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_FAILED ||
		req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_PARTIAL_SUCCESS {
		// 终态且未显式提供时间时，使用服务端当前时间作为完成时间。
		data.ProcessedAt = gtime.Now()
	}

	// 仅允许更新未软删资产，避免历史资产被异步任务误覆盖。
	ret, err := dao.MediaAsset.Ctx(ctx).
		Where(dao.MediaAsset.Columns().AssetId, req.GetAssetId()).
		WhereNull(dao.MediaAsset.Columns().DeletedAt).
		Data(data).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update process status failed")
	}
	// Updated=true 代表至少命中一行（资产存在且未删除）。
	rows, _ := ret.RowsAffected()
	return &v1.UpdateAssetProcessStatusRes{Updated: rows > 0}, nil
}

func (s *sMedia) getScenePolicy(ctx context.Context, sceneCode string) (*entity.MediaScenePolicy, error) {
	// 场景策略读取统一入口：只暴露启用状态，停用/不存在都按 not found 处理。
	var row entity.MediaScenePolicy
	cols := dao.MediaScenePolicy.Columns()
	if err := dao.MediaScenePolicy.Ctx(ctx).
		Where(cols.SceneCode, sceneCode).
		Where(cols.Status, consts.ScenePolicyStatusEnabled).
		Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query media_scene_policy failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "scene policy not found or disabled")
	}
	return &row, nil
}

func (s *sMedia) getAssetByID(ctx context.Context, assetID uint64) (*entity.MediaAsset, error) {
	// 统一资产读取入口：默认过滤软删除记录（与 GC 逻辑保持一致）。
	var row entity.MediaAsset
	cols := dao.MediaAsset.Columns()
	if err := dao.MediaAsset.Ctx(ctx).
		Where(cols.AssetId, assetID).
		WhereNull(cols.DeletedAt).
		Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query media_asset failed")
	}
	if row.AssetId == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "asset not found")
	}
	return &row, nil
}

func (s *sMedia) validateBindableAssets(ctx context.Context, sceneCode string, assetIDs []uint64) (map[uint64]*entity.MediaAsset, error) {
	// 绑定前校验：
	// 1) 资产存在；
	// 2) 场景匹配；
	// 3) 仅允许 ORIGINAL 资产直接绑定业务对象（DERIVED 由链路消费）。
	assetMap, err := s.queryAssetsByIDs(ctx, assetIDs)
	if err != nil {
		return nil, err
	}
	for _, id := range assetIDs {
		row := assetMap[id]
		if row == nil {
			return nil, gerror.NewCode(gcode.CodeNotFound, fmt.Sprintf("asset %d not found", id))
		}
		// 绑定关系不允许跨场景，避免不同场景策略互相污染。
		if strings.TrimSpace(row.SceneCode) != sceneCode {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, fmt.Sprintf("asset %d scene mismatch", id))
		}
		// 仅 ORIGINAL 可绑定业务，防止把派生结果当主资产误绑定。
		if row.AssetRole != consts.AssetRoleOriginal {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, fmt.Sprintf("asset %d is not ORIGINAL", id))
		}
	}
	return assetMap, nil
}

func (s *sMedia) queryAssetsByIDs(ctx context.Context, assetIDs []uint64) (map[uint64]*entity.MediaAsset, error) {
	// 批量查询资产并返回 map：
	// - 查询阶段按 IN 一次取回；
	// - 输出阶段以 asset_id 为 key，方便 O(1) 访问。
	result := make(map[uint64]*entity.MediaAsset, len(assetIDs))
	if len(assetIDs) == 0 {
		return result, nil
	}
	// 先做去重和零值过滤，避免无效 IN 条件。
	uniq := make([]uint64, 0, len(assetIDs))
	seen := make(map[uint64]struct{}, len(assetIDs))
	for _, id := range assetIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}

	var rows []*entity.MediaAsset
	cols := dao.MediaAsset.Columns()
	if err := dao.MediaAsset.Ctx(ctx).
		WhereIn(cols.AssetId, toInterfaceSliceUint64(uniq)).
		WhereNull(cols.DeletedAt).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "batch query media_asset failed")
	}
	for _, row := range rows {
		result[row.AssetId] = row
	}
	return result, nil
}

func (s *sMedia) queryBindings(
	ctx context.Context,
	sceneCode string,
	bizType string,
	bizNo string,
	bindingField string,
	includeInactive bool,
) ([]*entity.MediaAssetBinding, error) {
	// 通用绑定查询：
	// - 固定按 scene/bizType/bizNo 定位业务对象；
	// - 可选按 binding_field 缩小范围；
	// - includeInactive 控制是否读取解绑历史。
	cols := dao.MediaAssetBinding.Columns()
	model := dao.MediaAssetBinding.Ctx(ctx).
		Where(cols.SceneCode, sceneCode).
		Where(cols.BizType, bizType).
		Where(cols.BizNo, bizNo)
	if strings.TrimSpace(bindingField) != "" {
		// 只看指定字段（如只看 cover，不看 gallery）。
		model = model.Where(cols.BindingField, bindingField)
	}
	if !includeInactive {
		// 默认仅返回活跃绑定，避免把历史解绑记录暴露给业务层。
		model = model.Where(cols.IsActive, consts.BindingActive)
	}

	var rows []*entity.MediaAssetBinding
	// 排序规则：
	// 1) 先按字段分组；
	// 2) 活跃优先；
	// 3) 再按业务指定顺序和主键稳定排序。
	if err := model.
		OrderAsc(cols.BindingField).
		OrderDesc(cols.IsActive).
		OrderAsc(cols.SortOrder).
		OrderAsc(cols.Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query media_asset_binding failed")
	}
	return rows, nil
}

func toProtoAsset(row *entity.MediaAsset) *v1.Asset {
	// DB 实体 -> RPC 结构，字段基本一一映射，枚举值显式做类型转换。
	if row == nil {
		return nil
	}
	return &v1.Asset{
		AssetId:             row.AssetId,
		ParentAssetId:       row.ParentAssetId,
		RootAssetId:         row.RootAssetId,
		SceneCode:           row.SceneCode,
		AclType:             v1.AclType(row.AclType),
		AssetRole:           v1.AssetRole(row.AssetRole),
		DerivedKind:         v1.DerivedKind(row.DerivedKind),
		FileName:            row.FileName,
		MimeType:            row.MimeType,
		SizeBytes:           row.SizeBytes,
		ChecksumSha256:      row.ChecksumSha256,
		Etag:                row.Etag,
		PublicUrl:           row.PublicUrl,
		RiskStatus:          v1.RiskStatus(row.RiskStatus),
		ProcessStatus:       v1.ProcessStatus(row.ProcessStatus),
		ProcessProgress:     uint32(row.ProcessProgress),
		ProcessErrorCode:    row.ProcessErrorCode,
		ProcessErrorMessage: row.ProcessErrorMessage,
		ProcessedAt:         toProtoTimestamp(row.ProcessedAt),
		CreatedAt:           toProtoTimestamp(row.CreatedAt),
		UpdatedAt:           toProtoTimestamp(row.UpdatedAt),
	}
}

func toProtoBinding(row *entity.MediaAssetBinding) *v1.BizAssetBinding {
	// 绑定实体 -> RPC 结构，is_active 数值转 bool。
	if row == nil {
		return nil
	}
	return &v1.BizAssetBinding{
		BindingId:    row.Id,
		SceneCode:    row.SceneCode,
		BizType:      row.BizType,
		BizNo:        row.BizNo,
		BindingField: row.BindingField,
		AssetId:      row.AssetId,
		SortOrder:    int32(row.SortOrder),
		IsActive:     row.IsActive == consts.BindingActive,
		CreatedAt:    toProtoTimestamp(row.CreatedAt),
		UpdatedAt:    toProtoTimestamp(row.UpdatedAt),
	}
}

func toProtoBindingList(rows []*entity.MediaAssetBinding) []*v1.BizAssetBinding {
	// 批量转换，空集合返回 nil 以保持历史接口语义。
	if len(rows) == 0 {
		return nil
	}
	out := make([]*v1.BizAssetBinding, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProtoBinding(row))
	}
	return out
}

func toProtoTimestamp(t *gtime.Time) *timestamppb.Timestamp {
	// 统一时间转换：空值/零值返回 nil，避免把“未知时间”误表述成 Unix epoch。
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(time.Unix(0, t.TimestampNano()))
}

func cfgString(ctx context.Context, key string, defaultVal string) string {
	// 从配置中心读取字符串：
	// - 配置不存在/空白时回退默认值；
	// - 读取失败由 MustGet 抛 panic（保持与现有项目约定一致）。
	v := strings.TrimSpace(g.Cfg().MustGet(ctx, key).String())
	if v == "" {
		return defaultVal
	}
	return v
}

func buildObjectKey(sceneCode string, fileName string) string {
	// 生成对象路径：
	// 说明：scene/yyyyMMdd/unixNano_uuid_filename
	// 通过替换路径分隔符防止文件名注入目录结构（目录穿越防护）。
	safeName := strings.ReplaceAll(strings.TrimSpace(fileName), " ", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	if safeName == "" {
		safeName = "file.bin"
	}
	now := time.Now()
	// 使用 unixNano + uuid 双保险降低冲突概率，scene 小写便于前缀归类。
	return path.Join(
		strings.ToLower(sceneCode),
		now.Format("20060102"),
		fmt.Sprintf("%d_%s_%s", now.UnixNano(), uuid.NewString(), safeName),
	)
}

func buildPublicURL(ctx context.Context, bucket string, objectKey string) string {
	// 公开访问 URL 拼装策略：
	// 1) 优先 cdnBaseUrl（外网加速）；
	// 2) 回退 storage.endpoint（直连存储网关）。
	base := cfgString(ctx, "storage.cdnBaseUrl", "")
	if base == "" {
		base = cfgString(ctx, "storage.endpoint", "")
	}
	if base == "" {
		return ""
	}
	base = strings.TrimRight(base, "/")
	key := strings.TrimLeft(objectKey, "/")
	if strings.TrimSpace(bucket) == "" {
		// bucket 为空时兼容 endpoint 已内置 bucket 的部署模式。
		return fmt.Sprintf("%s/%s", base, key)
	}
	return fmt.Sprintf("%s/%s/%s", base, bucket, key)
}

func buildUploadURL(ctx context.Context, bucket string, objectKey string) string {
	// 上传地址拼装：
	// - 优先 uploadEndpoint（读写分离时常见）；
	// - 其次 endpoint；
	// - 都没有时返回本地占位域名，避免空串。
	base := cfgString(ctx, "storage.uploadEndpoint", "")
	if base == "" {
		base = cfgString(ctx, "storage.endpoint", "")
	}
	base = strings.TrimRight(base, "/")
	key := strings.TrimLeft(objectKey, "/")
	if base == "" {
		return fmt.Sprintf("https://upload.shopa.local/%s/%s", bucket, key)
	}
	return fmt.Sprintf("%s/%s/%s", base, bucket, key)
}

func buildSignedReadURL(ctx context.Context, asset *entity.MediaAsset, expUnix int64) string {
	// 轻量签名 URL：
	// - 输入：asset_id + exp + secret；
	// - 输出：base?exp=...&sig=sha256(...)；
	// - 用途：在对象存储原生签名之外提供统一签发层。
	base := strings.TrimSpace(asset.PublicUrl)
	if base == "" {
		// 无缓存 public_url 时按当前配置重新拼装。
		base = buildPublicURL(ctx, asset.Bucket, asset.ObjectKey)
	}
	if base == "" {
		// 极端配置缺失时也返回可识别的占位 URL，避免响应为空。
		base = fmt.Sprintf("https://read.shopa.local/%d", asset.AssetId)
	}
	secret := cfgString(ctx, "storage.secretKey", "shopa-dev-secret")
	// 签名原文保持稳定格式，便于多语言侧复算校验。
	sigRaw := fmt.Sprintf("%d|%d|%s", asset.AssetId, expUnix, secret)
	sum := sha256.Sum256([]byte(sigRaw))
	sig := hex.EncodeToString(sum[:])
	separator := "?"
	if strings.Contains(base, "?") {
		// base 已带 query 时改用 & 追加参数。
		separator = "&"
	}
	return fmt.Sprintf("%s%sexp=%d&sig=%s", base, separator, expUnix, sig)
}

func hashStorageObject(provider string, bucket string, objectKey string) string {
	// 对存储定位三元组做稳定哈希：
	// - 这是“对象位置指纹”，不是文件内容指纹；
	// - GC 任务可据此归并“同物理对象被多资产引用”的关系。
	raw := provider + "|" + bucket + "|" + objectKey
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func mimeAllowed(allowMIMEJSON string, mimeType string) bool {
	// MIME 白名单策略：
	// - 配置为空/非法时默认放行（保持兼容）；
	// - 配置有效时要求精确匹配（忽略大小写）。
	raw := strings.TrimSpace(allowMIMEJSON)
	if raw == "" || raw == "null" {
		// 未配置白名单时默认放行，保持向后兼容。
		return true
	}
	var allowList []string
	if err := json.Unmarshal([]byte(raw), &allowList); err != nil {
		// 配置格式错误时也放行，避免因配置脏数据导致全量上传失败。
		return true
	}
	if len(allowList) == 0 {
		return true
	}
	target := strings.ToLower(strings.TrimSpace(mimeType))
	for _, item := range allowList {
		if strings.ToLower(strings.TrimSpace(item)) == target {
			return true
		}
	}
	return false
}

func metadataValue(ctx context.Context, keys ...string) string {
	// 从 gRPC metadata 中按候选 key 顺序取值，兼容大小写 key。
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	for _, key := range keys {
		// 先取小写 key，兼容代理层统一 lowercase 的情况。
		values := md.Get(strings.ToLower(key))
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
		// 再取原始 key，兼容直传大小写。
		values = md.Get(key)
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}

func userIDFromMetadata(ctx context.Context, keys ...string) (uint64, bool) {
	// 从 metadata 解析正整数用户 ID，返回 (id, 是否成功)。
	for _, key := range keys {
		raw := metadataValue(ctx, key)
		if raw == "" {
			continue
		}
		id, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && id > 0 {
			return id, true
		}
	}
	return 0, false
}

func toInterfaceSliceUint64(ids []uint64) []interface{} {
	// 适配 GoFrame WhereIn：[]uint64 -> []interface{}。
	if len(ids) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		out = append(out, id)
	}
	return out
}

func isDuplicateErr(err error) bool {
	// 识别常见唯一约束冲突（MySQL duplicate entry）并转业务可读错误。
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate entry")
}

// 对 map 的 key 排序，保证调用方拼装响应时顺序稳定（便于测试与排查）。
func sortedUint64Keys(m map[uint64]*entity.MediaAsset) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
