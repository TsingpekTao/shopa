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
	"github.com/google/uuid"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sMedia implements media domain logic for both external and internal RPC services.
type sMedia struct{}

// New creates media logic instance.
func New() *sMedia {
	return &sMedia{}
}

func init() {
	service.RegisterMedia(New())
}

func (s *sMedia) InitUpload(ctx context.Context, req *v1.InitUploadReq) (*v1.InitUploadRes, error) {
	sceneCode := strings.TrimSpace(req.GetSceneCode())
	if sceneCode == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code is required")
	}
	if strings.TrimSpace(req.GetFileName()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "file_name is required")
	}
	if strings.TrimSpace(req.GetMimeType()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "mime_type is required")
	}

	policy, err := s.getScenePolicy(ctx, sceneCode)
	if err != nil {
		return nil, err
	}
	if policy.MaxSizeBytes > 0 && req.GetSizeBytes() > policy.MaxSizeBytes {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "size exceeds scene max_size_bytes")
	}
	if !mimeAllowed(policy.AllowMimeJson, req.GetMimeType()) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "mime_type is not allowed by scene policy")
	}

	var (
		now             = gtime.Now()
		storageProvider = cfgString(ctx, "storage.provider", "minio")
		bucket          = cfgString(ctx, "storage.bucket", "shopa")
		objectKey       = buildObjectKey(sceneCode, req.GetFileName())
		publicURL       = ""
		riskStatus      = uint(consts.RiskStatusPending)
		processStatus   = uint(consts.ProcessStatusPending)
		processProgress = uint(0)
		processedAt     *gtime.Time
	)

	if uint(policy.AclType) == consts.ACLTypePublicRead {
		publicURL = buildPublicURL(ctx, bucket, objectKey)
	}
	if policy.RiskAsyncEnabled == 0 {
		riskStatus = uint(consts.RiskStatusPassed)
	}
	if policy.ProcessAsyncEnabled == 0 {
		processStatus = uint(consts.ProcessStatusSuccess)
		processProgress = 100
		processedAt = now
	}

	traceBizType := req.GetBizType()
	traceBizNo := req.GetBizNo()
	uploaderID, _ := userIDFromMetadata(ctx, "x-user-id", "user_id", "uid")

	insertRes, err := dao.MediaAsset.Ctx(ctx).Data(do.MediaAsset{
		SceneCode:         sceneCode,
		AclType:           uint(policy.AclType),
		AssetRole:         consts.AssetRoleOriginal,
		DerivedKind:       consts.DerivedKindUnspecified,
		FileName:          req.GetFileName(),
		MimeType:          req.GetMimeType(),
		SizeBytes:         req.GetSizeBytes(),
		ChecksumSha256:    req.GetChecksumSha256(),
		StorageProvider:   storageProvider,
		Bucket:            bucket,
		ObjectKey:         objectKey,
		StorageObjectHash: hashStorageObject(storageProvider, bucket, objectKey),
		PublicUrl:         publicURL,
		RiskStatus:        riskStatus,
		ProcessStatus:     processStatus,
		ProcessProgress:   processProgress,
		ProcessedAt:       processedAt,
		UploaderUserId:    uploaderID,
		TraceBizType:      traceBizType,
		TraceBizNo:        traceBizNo,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "insert media_asset failed")
	}
	lastID, err := insertRes.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "read inserted asset id failed")
	}
	assetID := uint64(lastID)

	_, err = dao.MediaAsset.Ctx(ctx).
		Where(dao.MediaAsset.Columns().AssetId, assetID).
		Data(do.MediaAsset{RootAssetId: assetID}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "set root_asset_id failed")
	}

	asset, err := s.getAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	uploadURL := buildUploadURL(ctx, bucket, objectKey)
	return &v1.InitUploadRes{
		Asset: toProtoAsset(asset),
		UploadTicket: &v1.UploadTicket{
			Method:    "PUT",
			UploadUrl: uploadURL,
			Headers: map[string]string{
				"Content-Type": req.GetMimeType(),
			},
			ExpiredAt: timestamppb.New(time.Now().Add(15 * time.Minute)),
		},
	}, nil
}

func (s *sMedia) CompleteUpload(ctx context.Context, req *v1.CompleteUploadReq) (*v1.CompleteUploadRes, error) {
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}
	policy, err := s.getScenePolicy(ctx, asset.SceneCode)
	if err != nil {
		return nil, err
	}

	data := do.MediaAsset{
		Etag:           req.GetEtag(),
		MimeType:       req.GetMimeType(),
		SizeBytes:      req.GetSizeBytes(),
		ChecksumSha256: req.GetChecksumSha256(),
	}
	if policy.ProcessAsyncEnabled == 0 {
		data.ProcessStatus = consts.ProcessStatusSuccess
		data.ProcessProgress = 100
		data.ProcessedAt = gtime.Now()
	} else {
		data.ProcessStatus = consts.ProcessStatusProcessing
		data.ProcessProgress = 0
	}
	if policy.RiskAsyncEnabled == 0 {
		data.RiskStatus = consts.RiskStatusPassed
	}

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
	return &v1.CompleteUploadRes{Asset: toProtoAsset(updated)}, nil
}

func (s *sMedia) IssueReadUrl(ctx context.Context, req *v1.IssueReadUrlReq) (*v1.IssueReadUrlRes, error) {
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}

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

	ttl := int(req.GetTtlSeconds())
	if ttl <= 0 {
		ttl = consts.DefaultIssueReadTTLSeconds
	}
	if ttl > consts.MaxIssueReadTTLSeconds {
		ttl = consts.MaxIssueReadTTLSeconds
	}
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
	signedURL := buildSignedReadURL(ctx, asset, expireAt.Unix())
	return &v1.IssueReadUrlRes{
		AssetId:   asset.AssetId,
		Url:       signedURL,
		IsPublic:  false,
		ExpiredAt: timestamppb.New(expireAt),
	}, nil
}

func (s *sMedia) GetAssetProcessStatus(ctx context.Context, req *v1.GetAssetProcessStatusReq) (*v1.GetAssetProcessStatusRes, error) {
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	asset, err := s.getAssetByID(ctx, req.GetAssetId())
	if err != nil {
		return nil, err
	}
	rootID := asset.RootAssetId
	if rootID == 0 {
		rootID = asset.AssetId
	}

	var derivedRows []*entity.MediaAsset
	cols := dao.MediaAsset.Columns()
	if err = dao.MediaAsset.Ctx(ctx).
		Where(cols.RootAssetId, rootID).
		Where(cols.AssetRole, consts.AssetRoleDerived).
		WhereNull(cols.DeletedAt).
		OrderAsc(cols.AssetId).
		Scan(&derivedRows); err != nil {
		return nil, gerror.Wrap(err, "query derived assets failed")
	}

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

	operatorID, _ := userIDFromMetadata(ctx, "x-user-id", "user_id", "uid")
	requestID := metadataValue(ctx, "x-request-id")
	now := gtime.Now()

	err = dao.MediaAssetBinding.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.MediaAssetBinding.Columns()
		for _, item := range req.GetItems() {
			asset := assetMap[item.GetAssetId()]
			if asset == nil {
				return gerror.NewCode(gcode.CodeNotFound, fmt.Sprintf("asset %d not found", item.GetAssetId()))
			}
			field := strings.TrimSpace(item.GetBindingField())
			_, _ = tx.Model(dao.MediaAssetBinding.Table()).
				Where(cols.SceneCode, req.GetSceneCode()).
				Where(cols.BizType, req.GetBizType()).
				Where(cols.BizNo, req.GetBizNo()).
				Where(cols.BindingField, field).
				Where(cols.AssetId, item.GetAssetId()).
				Where(cols.IsActive, consts.BindingActive).
				Data(do.MediaAssetBinding{IsActive: consts.BindingInactive, UnboundAt: now}).
				Update()

			_, insertErr := tx.Model(dao.MediaAssetBinding.Table()).Data(do.MediaAssetBinding{
				SceneCode:      req.GetSceneCode(),
				BizType:        req.GetBizType(),
				BizNo:          req.GetBizNo(),
				BindingField:   field,
				AssetId:        item.GetAssetId(),
				SortOrder:      item.GetSortOrder(),
				IsActive:       consts.BindingActive,
				OperatorUserId: operatorID,
				RequestId:      requestID,
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

	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), "", false)
	if err != nil {
		return nil, err
	}
	return &v1.BatchBindAssetsToBizRes{Bindings: toProtoBindingList(bindings)}, nil
}

func (s *sMedia) ReplaceBizAssetBindings(ctx context.Context, req *v1.ReplaceBizAssetBindingsReq) (*v1.ReplaceBizAssetBindingsRes, error) {
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

	err = dao.MediaAssetBinding.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.MediaAssetBinding.Columns()
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

		for _, item := range req.GetItems() {
			_, insertErr := tx.Model(dao.MediaAssetBinding.Table()).Data(do.MediaAssetBinding{
				SceneCode:      req.GetSceneCode(),
				BizType:        req.GetBizType(),
				BizNo:          req.GetBizNo(),
				BindingField:   req.GetBindingField(),
				AssetId:        item.GetAssetId(),
				SortOrder:      item.GetSortOrder(),
				IsActive:       consts.BindingActive,
				OperatorUserId: operatorID,
				RequestId:      requestID,
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

	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), req.GetBindingField(), false)
	if err != nil {
		return nil, err
	}
	return &v1.ReplaceBizAssetBindingsRes{Bindings: toProtoBindingList(bindings)}, nil
}

func (s *sMedia) BatchUnbindAssetsFromBiz(ctx context.Context, req *v1.BatchUnbindAssetsFromBizReq) (*v1.BatchUnbindAssetsFromBizRes, error) {
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" || strings.TrimSpace(req.GetBindingField()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no/binding_field are required")
	}
	cols := dao.MediaAssetBinding.Columns()
	model := dao.MediaAssetBinding.Ctx(ctx).
		Where(cols.SceneCode, req.GetSceneCode()).
		Where(cols.BizType, req.GetBizType()).
		Where(cols.BizNo, req.GetBizNo()).
		Where(cols.BindingField, req.GetBindingField()).
		Where(cols.IsActive, consts.BindingActive)
	if len(req.GetAssetIds()) > 0 {
		model = model.WhereIn(cols.AssetId, toInterfaceSliceUint64(req.GetAssetIds()))
	}
	ret, err := model.Data(do.MediaAssetBinding{
		IsActive:  consts.BindingInactive,
		UnboundAt: gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "unbind assets failed")
	}
	rows, _ := ret.RowsAffected()
	return &v1.BatchUnbindAssetsFromBizRes{AffectedRows: uint64(rows)}, nil
}

func (s *sMedia) GetBizAssets(ctx context.Context, req *v1.GetBizAssetsReq) (*v1.GetBizAssetsRes, error) {
	if strings.TrimSpace(req.GetSceneCode()) == "" || strings.TrimSpace(req.GetBizType()) == "" || strings.TrimSpace(req.GetBizNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "scene_code/biz_type/biz_no are required")
	}
	bindings, err := s.queryBindings(ctx, req.GetSceneCode(), req.GetBizType(), req.GetBizNo(), req.GetBindingField(), req.GetIncludeInactive())
	if err != nil {
		return nil, err
	}
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
			Binding: toProtoBinding(row),
		}
		if asset := assetMap[row.AssetId]; asset != nil {
			item.Asset = toProtoAsset(asset)
		}
		out = append(out, item)
	}
	return &v1.GetBizAssetsRes{Assets: out}, nil
}

func (s *sMedia) CreateDerivedAsset(ctx context.Context, req *v1.CreateDerivedAssetReq) (*v1.CreateDerivedAssetRes, error) {
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
	sceneCode := strings.TrimSpace(req.GetSceneCode())
	if sceneCode == "" {
		sceneCode = parent.SceneCode
	}

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
		publicURL = buildPublicURL(ctx, bucket, objectKey)
	}

	rootID := parent.RootAssetId
	if rootID == 0 {
		rootID = parent.AssetId
	}
	insertRes, err := dao.MediaAsset.Ctx(ctx).Data(do.MediaAsset{
		ParentAssetId:     parent.AssetId,
		RootAssetId:       rootID,
		SceneCode:         sceneCode,
		AclType:           parent.AclType,
		AssetRole:         consts.AssetRoleDerived,
		DerivedKind:       uint(req.GetDerivedKind()),
		FileName:          req.GetFileName(),
		MimeType:          req.GetMimeType(),
		SizeBytes:         req.GetSizeBytes(),
		ChecksumSha256:    req.GetChecksumSha256(),
		Etag:              req.GetEtag(),
		StorageProvider:   storageProvider,
		Bucket:            bucket,
		ObjectKey:         objectKey,
		StorageObjectHash: hashStorageObject(storageProvider, bucket, objectKey),
		PublicUrl:         publicURL,
		RiskStatus:        parent.RiskStatus,
		ProcessStatus:     consts.ProcessStatusSuccess,
		ProcessProgress:   100,
		ProcessedAt:       gtime.Now(),
		UploaderUserId:    parent.UploaderUserId,
		TraceBizType:      parent.TraceBizType,
		TraceBizNo:        parent.TraceBizNo,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "insert derived asset failed")
	}
	lastID, err := insertRes.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "read derived asset id failed")
	}
	derived, err := s.getAssetByID(ctx, uint64(lastID))
	if err != nil {
		return nil, err
	}
	return &v1.CreateDerivedAssetRes{Asset: toProtoAsset(derived)}, nil
}

func (s *sMedia) UpdateAssetProcessStatus(ctx context.Context, req *v1.UpdateAssetProcessStatusReq) (*v1.UpdateAssetProcessStatusRes, error) {
	if req.GetAssetId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "asset_id is required")
	}
	data := do.MediaAsset{
		ProcessStatus:       uint(req.GetProcessStatus()),
		ProcessProgress:     req.GetProcessProgress(),
		ProcessErrorCode:    req.GetProcessErrorCode(),
		ProcessErrorMessage: req.GetProcessErrorMessage(),
	}
	if req.GetProcessedAt() != nil {
		data.ProcessedAt = gtime.NewFromTime(req.GetProcessedAt().AsTime())
	} else if req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_SUCCESS ||
		req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_FAILED ||
		req.GetProcessStatus() == v1.ProcessStatus_PROCESS_STATUS_PARTIAL_SUCCESS {
		data.ProcessedAt = gtime.Now()
	}

	ret, err := dao.MediaAsset.Ctx(ctx).
		Where(dao.MediaAsset.Columns().AssetId, req.GetAssetId()).
		WhereNull(dao.MediaAsset.Columns().DeletedAt).
		Data(data).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update process status failed")
	}
	rows, _ := ret.RowsAffected()
	return &v1.UpdateAssetProcessStatusRes{Updated: rows > 0}, nil
}

func (s *sMedia) getScenePolicy(ctx context.Context, sceneCode string) (*entity.MediaScenePolicy, error) {
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
	assetMap, err := s.queryAssetsByIDs(ctx, assetIDs)
	if err != nil {
		return nil, err
	}
	for _, id := range assetIDs {
		row := assetMap[id]
		if row == nil {
			return nil, gerror.NewCode(gcode.CodeNotFound, fmt.Sprintf("asset %d not found", id))
		}
		if strings.TrimSpace(row.SceneCode) != sceneCode {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, fmt.Sprintf("asset %d scene mismatch", id))
		}
		if row.AssetRole != consts.AssetRoleOriginal {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, fmt.Sprintf("asset %d is not ORIGINAL", id))
		}
	}
	return assetMap, nil
}

func (s *sMedia) queryAssetsByIDs(ctx context.Context, assetIDs []uint64) (map[uint64]*entity.MediaAsset, error) {
	result := make(map[uint64]*entity.MediaAsset, len(assetIDs))
	if len(assetIDs) == 0 {
		return result, nil
	}
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
	cols := dao.MediaAssetBinding.Columns()
	model := dao.MediaAssetBinding.Ctx(ctx).
		Where(cols.SceneCode, sceneCode).
		Where(cols.BizType, bizType).
		Where(cols.BizNo, bizNo)
	if strings.TrimSpace(bindingField) != "" {
		model = model.Where(cols.BindingField, bindingField)
	}
	if !includeInactive {
		model = model.Where(cols.IsActive, consts.BindingActive)
	}

	var rows []*entity.MediaAssetBinding
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
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(time.Unix(0, t.TimestampNano()))
}

func cfgString(ctx context.Context, key string, defaultVal string) string {
	v := strings.TrimSpace(g.Cfg().MustGet(ctx, key).String())
	if v == "" {
		return defaultVal
	}
	return v
}

func buildObjectKey(sceneCode string, fileName string) string {
	safeName := strings.ReplaceAll(strings.TrimSpace(fileName), " ", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	if safeName == "" {
		safeName = "file.bin"
	}
	now := time.Now()
	return path.Join(
		strings.ToLower(sceneCode),
		now.Format("20060102"),
		fmt.Sprintf("%d_%s_%s", now.UnixNano(), uuid.NewString(), safeName),
	)
}

func buildPublicURL(ctx context.Context, bucket string, objectKey string) string {
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
		return fmt.Sprintf("%s/%s", base, key)
	}
	return fmt.Sprintf("%s/%s/%s", base, bucket, key)
}

func buildUploadURL(ctx context.Context, bucket string, objectKey string) string {
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
	base := strings.TrimSpace(asset.PublicUrl)
	if base == "" {
		base = buildPublicURL(ctx, asset.Bucket, asset.ObjectKey)
	}
	if base == "" {
		base = fmt.Sprintf("https://read.shopa.local/%d", asset.AssetId)
	}
	secret := cfgString(ctx, "storage.secretKey", "shopa-dev-secret")
	sigRaw := fmt.Sprintf("%d|%d|%s", asset.AssetId, expUnix, secret)
	sum := sha256.Sum256([]byte(sigRaw))
	sig := hex.EncodeToString(sum[:])
	separator := "?"
	if strings.Contains(base, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%sexp=%d&sig=%s", base, separator, expUnix, sig)
}

func hashStorageObject(provider string, bucket string, objectKey string) string {
	raw := provider + "|" + bucket + "|" + objectKey
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func mimeAllowed(allowMIMEJSON string, mimeType string) bool {
	raw := strings.TrimSpace(allowMIMEJSON)
	if raw == "" || raw == "null" {
		return true
	}
	var allowList []string
	if err := json.Unmarshal([]byte(raw), &allowList); err != nil {
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	for _, key := range keys {
		values := md.Get(strings.ToLower(key))
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
		values = md.Get(key)
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}

func userIDFromMetadata(ctx context.Context, keys ...string) (uint64, bool) {
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
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate entry")
}

// Keep deterministic ordering in callers that assemble response from maps.
func sortedUint64Keys(m map[uint64]*entity.MediaAsset) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
