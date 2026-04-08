package aftersale

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	paymentv1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/dao"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// defaultPageSize 定义售后列表默认分页大小，避免调用方不传时一次性扫太多数据。
	defaultPageSize = 20
	// maxPageSize 限制单次列表查询上限，避免超大分页拖垮数据库和网关。
	maxPageSize = 100
)

// sAfterSale 是售后领域的核心逻辑实现，承接 controller 转进来的所有售后业务请求。
type sAfterSale struct{}

func New() *sAfterSale {
	// 这里没有额外依赖初始化，直接返回空实现即可让 service 层注册使用。
	return &sAfterSale{}
}

func init() {
	// 在包初始化阶段注册售后逻辑实现，确保 controller 能通过 service 门面找到它。
	service.RegisterAfterSale(New())
}

func (s *sAfterSale) CreateAfterSale(ctx context.Context, req *v1.CreateAfterSaleReq) (*v1.CreateAfterSaleRes, error) {
	// 创建售后单时必须拿到订单、子单和条目编号，否则无法精确定位售后对象。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetItemNo()) == "" {
		// 关键入参缺失时直接返回参数错误，避免落下脏售后单。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/item_no are required")
	}
	// 从上下文里提取当前买家 user_id，确保售后单绑定到真实发起人。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		// 无法识别用户身份时不能继续创建售后单，直接把鉴权错误返回上层。
		return nil, err
	}
	// 把证据附件列表序列化成 JSON，便于直接存入售后主表。
	evidenceJSON, _ := json.Marshal(req.GetEvidenceAssetIds())
	// 生成售后单号，后续所有查询、取消、审批都会围绕这个业务号展开。
	afterSaleNo := generateBizNo("AS")
	// 把售后申请主记录写入数据库，初始状态固定为“待卖家审核”。
	_, err = dao.AfterSaleCase.Ctx(ctx).Data(do.AfterSaleCase{
		AfterSaleNo:          afterSaleNo,
		OrderNo:              strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:           strings.TrimSpace(req.GetSubOrderNo()),
		ItemNo:               strings.TrimSpace(req.GetItemNo()),
		UserId:               userID,
		Qty:                  req.GetQty(),
		AfterSaleType:        int(req.GetAfterSaleType()),
		AfterSaleStatus:      int(v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW),
		ApplyRefundAmount:    req.GetApplyRefundAmount(),
		ApprovedRefundAmount: 0,
		ReasonCode:           strings.TrimSpace(req.GetReasonCode()),
		ReasonDesc:           strings.TrimSpace(req.GetReasonDesc()),
		EvidenceAssetIdsJson: string(evidenceJSON),
		BuyerRemark:          strings.TrimSpace(req.GetBuyerRemark()),
		Version:              1,
	}).Insert()
	if err != nil {
		// 主记录插入失败时直接返回，避免接口层误判为售后创建成功。
		return nil, gerror.Wrap(err, "create after_sale_case failed")
	}
	// 重新按业务号回查一次最新记录，统一走同一套 proto 映射逻辑返回给前端。
	row, err := s.getAfterSaleByNo(ctx, afterSaleNo)
	if err != nil {
		// 理论上刚插入成功就应该能读到，读不到说明数据库状态异常，需要显式返回错误。
		return nil, err
	}
	// 把数据库实体转换成对外协议对象，作为创建接口的最终响应。
	return &v1.CreateAfterSaleRes{AfterSale: toProtoAfterSaleCase(row)}, nil
}

func (s *sAfterSale) CancelAfterSale(ctx context.Context, req *v1.CancelAfterSaleReq) (*v1.CancelAfterSaleRes, error) {
	// 取消售后时必须传售后单号，否则无法命中要取消的申请。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 缺主键参数时直接拒绝请求，避免误取消其他单据。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 提取当前操作人的 user_id，用于校验售后单归属。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		// 用户身份不明时不能操作售后单，直接返回错误。
		return nil, err
	}
	// 先读取当前售后单快照，后续归属校验和版本校验都依赖它。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 售后单不存在或查询失败时不能继续取消。
		return nil, err
	}
	// 只有售后单所属买家自己才能取消，防止越权操作他人售后。
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "after_sale case does not belong to user")
	}
	// 如果前端带了乐观锁版本，就必须和数据库当前版本一致，避免覆盖并发更新。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	// 取消售后会推动版本号递增，防止旧快照继续写入。
	nextVersion := row.Version + 1
	// 把售后状态改成“已取消”，同时记录取消原因和关闭时间。
	_, err = dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).Data(do.AfterSaleCase{
		AfterSaleStatus:  int(v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED),
		CancelReasonCode: strings.TrimSpace(req.GetReasonCode()),
		Version:          nextVersion,
		ClosedAt:         gtime.Now(),
	}).Update()
	if err != nil {
		// 更新失败时不返回成功，确保取消语义与数据库状态保持一致。
		return nil, gerror.Wrap(err, "cancel after_sale_case failed")
	}
	// 返回最小必要字段给调用方，明确告知本次取消已经生效。
	return &v1.CancelAfterSaleRes{
		AfterSaleNo:     row.AfterSaleNo,
		AfterSaleStatus: v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED,
	}, nil
}

func (s *sAfterSale) GetMyAfterSaleDetail(ctx context.Context, req *v1.GetMyAfterSaleDetailReq) (*v1.GetMyAfterSaleDetailRes, error) {
	// 买家查询售后详情时必须指定售后单号。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 没有售后单号就无法定位详情对象，直接返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 提取当前买家身份，用于校验数据归属。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		// 无法识别买家时不允许继续查询售后详情。
		return nil, err
	}
	// 先读取售后主单详情。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 售后单查不到时直接把错误返回调用方。
		return nil, err
	}
	// 买家只能看自己的售后单，避免用户之间信息泄露。
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "after_sale case does not belong to user")
	}
	// 再按售后单号补充查询关联退款任务，方便详情页一次性拿全信息。
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, row.AfterSaleNo)
	// 把售后主单和退款任务一起返回，减少前端二次请求。
	return &v1.GetMyAfterSaleDetailRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) ListMyAfterSales(ctx context.Context, req *v1.ListMyAfterSalesReq) (*v1.ListMyAfterSalesRes, error) {
	// 售后列表只允许基于当前登录买家维度查询，先拿到 user_id。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		// 没有买家身份就不能返回个人售后列表。
		return nil, err
	}
	// 统一规整分页大小，避免非法 page_size 直接打到数据库。
	pageSize := normalizePageSize(req.GetPageSize())
	// 解析游标，支持买家售后列表的稳定翻页。
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		// 游标不合法时直接返回参数错误，防止分页语义错乱。
		return nil, err
	}
	// 先构造“当前买家 + 倒序分页”的基础查询模型。
	model := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().UserId, userID).OrderDesc(dao.AfterSaleCase.Columns().Id).Limit(pageSize + 1)
	// 带了游标时只查询游标之前的数据，保持游标翻页稳定性。
	if cursorID > 0 {
		model = model.WhereLT(dao.AfterSaleCase.Columns().Id, cursorID)
	}
	// 调用方带了状态筛选时，把枚举值统一转成数据库里的 int 状态码。
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			// 逐条把协议层状态枚举压平成数据库查询可用的整型状态。
			statuses = append(statuses, int(status))
		}
		// 把状态筛选条件挂进查询模型，只返回目标状态集合内的售后单。
		model = model.WhereIn(dao.AfterSaleCase.Columns().AfterSaleStatus, statuses)
	}
	// 承接数据库扫描出来的售后记录列表。
	var rows []entity.AfterSaleCase
	if err = model.Scan(&rows); err != nil {
		// 列表查询失败时直接返回数据库错误，避免前端看到空列表误判。
		return nil, gerror.Wrap(err, "list my after sales failed")
	}
	// 默认先认为当前页没有更多数据。
	hasMore := false
	// 多取一条命中时说明后面还有下一页。
	if len(rows) > pageSize {
		// 标记还有下一页，供前端决定是否继续翻页。
		hasMore = true
		// 裁掉额外探测出来的一条，只把当前页大小的数据返回给前端。
		rows = rows[:pageSize]
	}
	// 预分配响应列表，避免循环 append 时反复扩容。
	list := make([]*v1.AfterSaleCase, 0, len(rows))
	for _, row := range rows {
		// 逐条把数据库实体转换成协议对象，保持列表返回结构统一。
		list = append(list, toProtoAfterSaleCase(&row))
	}
	// 默认没有下一页时就返回空游标。
	next := ""
	// 有下一页时用最后一条记录的 ID 作为下次查询游标。
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	// 返回当前页售后列表、下一页游标和是否还有更多数据。
	return &v1.ListMyAfterSalesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sAfterSale) ListShopAfterSales(ctx context.Context, req *v1.ListShopAfterSalesReq) (*v1.ListShopAfterSalesRes, error) {
	// 店铺维度查询必须明确 shop_no，否则无法限定卖家数据边界。
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		// 缺店铺编号时直接拒绝，避免把全库售后单暴露给卖家端。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	// 统一规整分页大小，避免非正常 page_size 打穿数据库。
	pageSize := normalizePageSize(req.GetPageSize())
	// 解析游标，支持卖家后台稳定翻页。
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		// 游标非法时不能继续查询，否则会导致分页结果错乱。
		return nil, err
	}
	// 构造“当前店铺 + 倒序翻页”的基础查询模型。
	model := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().ShopNo, req.GetShopNo()).OrderDesc(dao.AfterSaleCase.Columns().Id).Limit(pageSize + 1)
	// 带游标时只取游标之前的数据，保证连续翻页不重不漏。
	if cursorID > 0 {
		model = model.WhereLT(dao.AfterSaleCase.Columns().Id, cursorID)
	}
	// 卖家端带了状态筛选时，把协议枚举统一转换成数据库状态码。
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			// 逐条把协议层状态压平成 int，便于 WhereIn 查询直接使用。
			statuses = append(statuses, int(status))
		}
		// 把状态集合挂到查询模型，只返回命中的售后单。
		model = model.WhereIn(dao.AfterSaleCase.Columns().AfterSaleStatus, statuses)
	}
	// 承接数据库扫描出来的售后单列表。
	var rows []entity.AfterSaleCase
	if err = model.Scan(&rows); err != nil {
		// 数据库查询失败时直接返回，避免卖家后台把故障误判成“暂无数据”。
		return nil, gerror.Wrap(err, "list shop after sales failed")
	}
	// 默认当前页没有更多数据。
	hasMore := false
	// 多查出一条时说明还有下一页。
	if len(rows) > pageSize {
		// 标记后续还有更多结果。
		hasMore = true
		// 裁掉额外探测的一条，保证返回数量与 page_size 一致。
		rows = rows[:pageSize]
	}
	// 预分配响应切片，减少循环追加时的内存扩容。
	list := make([]*v1.AfterSaleCase, 0, len(rows))
	for _, row := range rows {
		// 逐条把数据库实体转成协议对象，供卖家端列表直接渲染。
		list = append(list, toProtoAfterSaleCase(&row))
	}
	// 默认没有下一页时保持空游标。
	next := ""
	// 有下一页时用最后一条记录 ID 作为新的翻页游标。
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	// 返回卖家售后列表、翻页游标和 hasMore 标识。
	return &v1.ListShopAfterSalesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sAfterSale) GetShopAfterSaleDetail(ctx context.Context, req *v1.GetShopAfterSaleDetailReq) (*v1.GetShopAfterSaleDetailRes, error) {
	// 卖家查详情时必须指定售后单号。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 没有售后单号无法定位记录，直接返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 先加载售后主记录，作为详情页主体数据。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 查不到售后单时直接返回错误，避免后续继续查退款任务。
		return nil, err
	}
	// 再补充查询当前售后单关联的退款任务，供卖家侧一起展示。
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, row.AfterSaleNo)
	// 把售后主单和退款任务组合成详情响应。
	return &v1.GetShopAfterSaleDetailRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) ApproveAfterSale(ctx context.Context, req *v1.ApproveAfterSaleReq) (*v1.ApproveAfterSaleRes, error) {
	// 审核通过时必须指定要审批的售后单号。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 没有售后单号时直接拒绝审批请求。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 先读取当前售后单，后续要基于当前版本和订单信息创建退款任务。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 售后单不存在时不能继续审批。
		return nil, err
	}
	// 前端带了 expected_version 时必须做乐观锁校验，避免并发审核覆盖。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	// 预留退款任务指针，事务提交后会回查最新任务并返回给调用方。
	var task *entity.RefundTask
	// 审核通过要同时更新售后状态并创建退款任务，这里必须放在一个事务里。
	err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 审核通过后售后主单版本递增，防止旧快照继续写回。
		nextVersion := row.Version + 1
		if _, e := tx.Model(dao.AfterSaleCase.Table()).
			Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
			Data(do.AfterSaleCase{
				AfterSaleStatus:      int(v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK),
				ApprovedRefundAmount: req.GetApprovedRefundAmount(),
				SellerReply:          strings.TrimSpace(req.GetSellerReply()),
				Version:              nextVersion,
			}).Update(); e != nil {
			// 主单状态更新失败时整个事务直接回滚，避免只创建任务不改状态。
			return gerror.Wrap(e, "approve after_sale_case failed")
		}
		// 生成退款任务单号，后续支付/退款补偿都以这个业务号为锚点。
		taskNo := generateBizNo("RF")
		if _, e := tx.Model(dao.RefundTask.Table()).Data(do.RefundTask{
			RefundTaskNo: taskNo,
			AfterSaleNo:  row.AfterSaleNo,
			OrderNo:      row.OrderNo,
			SubOrderNo:   row.SubOrderNo,
			RefundAmount: req.GetApprovedRefundAmount(),
			Status:       int(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			Version:      1,
		}).Insert(); e != nil {
			// 退款任务创建失败时也要回滚主单状态，确保审批和任务创建强一致。
			return gerror.Wrap(e, "create refund_task failed")
		}
		// 两步都成功时提交事务，正式生效“待退款任务执行”的状态切换。
		return nil
	})
	if err != nil {
		// 事务失败时直接把错误返回给调用方，避免误判审核通过。
		return nil, err
	}
	// 事务提交后重新回查售后主单，确保返回的是数据库最新状态。
	row, _ = s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	// 再回查当前售后单下关联的退款任务，供卖家端直接展示。
	task, _ = s.getRefundTaskByAfterSaleNo(ctx, req.GetAfterSaleNo())
	// 把审批后的主单和任务快照一起返回。
	return &v1.ApproveAfterSaleRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) RejectAfterSale(ctx context.Context, req *v1.RejectAfterSaleReq) (*v1.RejectAfterSaleRes, error) {
	// 审核驳回同样必须指定售后单号。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 缺主键时不能继续驳回流程。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 先读取售后单快照，后续要基于当前版本做状态推进。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 售后单不存在时直接返回错误。
		return nil, err
	}
	// 卖家如果带了预期版本，就必须和数据库当前版本一致。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	// 驳回会把版本递增，防止前端旧数据继续覆盖新状态。
	nextVersion := row.Version + 1
	// 把售后状态切到“卖家已驳回”，并记录驳回原因和卖家回复。
	_, err = dao.AfterSaleCase.Ctx(ctx).
		Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
		Data(do.AfterSaleCase{
			AfterSaleStatus:  int(v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED),
			RejectReasonCode: int(req.GetRejectReasonCode()),
			SellerReply:      strings.TrimSpace(req.GetSellerReply()),
			Version:          nextVersion,
		}).Update()
	if err != nil {
		// 状态更新失败时不能返回成功，避免状态机走岔。
		return nil, gerror.Wrap(err, "reject after_sale_case failed")
	}
	// 返回最小必要字段，告诉前端当前售后单已经进入驳回态。
	return &v1.RejectAfterSaleRes{
		AfterSaleNo:     row.AfterSaleNo,
		AfterSaleStatus: v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED,
	}, nil
}

func (s *sAfterSale) ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error) {
	// 执行退款任务时必须提供退款任务单号。
	if req == nil || strings.TrimSpace(req.GetRefundTaskNo()) == "" {
		// 缺少任务主键时无法定位要执行的退款任务。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_task_no is required")
	}
	// 先读取当前退款任务快照，判断任务是否已执行完成。
	task, err := s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	if err != nil {
		// 任务不存在时直接返回，避免继续推进售后状态。
		return nil, err
	}
	// 已经成功的任务再次执行时直接幂等返回成功，避免重复退款。
	if task.Status == uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED) {
		finalCashRefundAmount := task.FinalCashRefundAmount
		if finalCashRefundAmount == 0 && task.RefundAmount >= task.PointsCashOffsetAmount {
			finalCashRefundAmount = task.RefundAmount - task.PointsCashOffsetAmount
		}
		return &v1.ExecuteRefundTaskRes{
			Task:                   toProtoRefundTask(task),
			RefundSucceeded:        true,
			ShouldRetry:            false,
			FinalCashRefundAmount:  finalCashRefundAmount,
			PointsCashOffsetAmount: task.PointsCashOffsetAmount,
		}, nil
	}

	// 先把本地任务标记为 processing，避免并发执行同一退款任务时重复打外部服务。
	_, err = dao.RefundTask.Ctx(ctx).
		Where(dao.RefundTask.Columns().RefundTaskNo, task.RefundTaskNo).
		WhereIn(dao.RefundTask.Columns().Status, []uint{
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_FAILED),
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PROCESSING),
		}).
		Data(do.RefundTask{
			Status:           uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PROCESSING),
			RetryCount:       task.RetryCount + 1,
			NextRetryAt:      nil,
			LastErrorCode:    "",
			LastErrorMessage: "",
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "mark refund task processing failed")
	}

	afterSaleRow, err := s.getAfterSaleByNo(ctx, task.AfterSaleNo)
	if err != nil {
		return nil, err
	}
	conn, paymentClient, err := s.newPaymentClient(ctx)
	if err != nil {
		return nil, err
	}
	execRes, execErr := paymentClient.ExecuteRefundTask(ctx, &paymentv1.ExecuteRefundTaskReq{
		RefundTaskNo:   task.RefundTaskNo,
		IdempotencyKey: fmt.Sprintf("execute_%s", task.RefundTaskNo),
	})
	conn.Close()
	if execErr != nil {
		if markErr := s.markLocalRefundTaskFailed(ctx, task, "PAYMENT_EXECUTE_FAILED", execErr.Error()); markErr != nil {
			return nil, markErr
		}
		task, _ = s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
		return &v1.ExecuteRefundTaskRes{
			Task:            toProtoRefundTask(task),
			RefundSucceeded: false,
			ShouldRetry:     true,
		}, nil
	}

	orderConn, orderClient, err := s.newOrderClient(ctx)
	if err != nil {
		return nil, err
	}
	finalizeRes, finalizeErr := orderClient.FinalizeSubOrderRefund(ctx, &orderv1.FinalizeSubOrderRefundReq{
		OrderNo:              task.OrderNo,
		SubOrderNo:           task.SubOrderNo,
		RefundNo:             task.AfterSaleNo,
		ApprovedRefundAmount: task.RefundAmount,
	})
	orderConn.Close()
	if finalizeErr != nil {
		if markErr := s.markLocalRefundTaskFailed(ctx, task, "ORDER_FINALIZE_FAILED", finalizeErr.Error()); markErr != nil {
			return nil, markErr
		}
		task, _ = s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
		return &v1.ExecuteRefundTaskRes{
			Task:            toProtoRefundTask(task),
			RefundSucceeded: false,
			ShouldRetry:     true,
		}, nil
	}

	finalCashRefundAmount := task.RefundAmount
	if finalizeRes.GetPointsCashOffsetAmount() >= finalCashRefundAmount {
		finalCashRefundAmount = 0
	} else {
		finalCashRefundAmount -= finalizeRes.GetPointsCashOffsetAmount()
	}
	err = dao.RefundTask.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model(dao.RefundTask.Table()).
			Where(dao.RefundTask.Columns().RefundTaskNo, task.RefundTaskNo).
			Data(do.RefundTask{
				Status:                 uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED),
				NextRetryAt:            nil,
				LastErrorCode:          "",
				LastErrorMessage:       "",
				PointsReturnAmount:     finalizeRes.GetPointsReturnAmount(),
				PointsReverseAmount:    finalizeRes.GetPointsReverseAmount(),
				PointsCashOffsetAmount: finalizeRes.GetPointsCashOffsetAmount(),
				FinalCashRefundAmount:  finalCashRefundAmount,
				AccountDebtAfter:       finalizeRes.GetAccountDebtAfter(),
			}).Update(); e != nil {
			return gerror.Wrap(e, "mark refund_task succeeded failed")
		}
		if _, e := tx.Model(dao.AfterSaleCase.Table()).
			Where(dao.AfterSaleCase.Columns().AfterSaleNo, task.AfterSaleNo).
			Data(do.AfterSaleCase{
				AfterSaleStatus:      uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED),
				ApprovedRefundAmount: task.RefundAmount,
				ClosedAt:             gtime.Now(),
				Version:              afterSaleRow.Version + 1,
			}).Update(); e != nil {
			return gerror.Wrap(e, "mark after_sale_case refunded failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	task, _ = s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	return &v1.ExecuteRefundTaskRes{
		Task:                   toProtoRefundTask(task),
		RefundSucceeded:        execRes.GetStatus() == paymentv1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED,
		ShouldRetry:            false,
		FinalCashRefundAmount:  finalCashRefundAmount,
		PointsCashOffsetAmount: finalizeRes.GetPointsCashOffsetAmount(),
	}, nil
}

func (s *sAfterSale) RetryRefundTask(ctx context.Context, req *v1.RetryRefundTaskReq) (*v1.RetryRefundTaskRes, error) {
	// 人工或调度重试退款任务时必须带退款任务号。
	if req == nil || strings.TrimSpace(req.GetRefundTaskNo()) == "" {
		// 缺主键时无法把指定任务重置为待执行状态。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_task_no is required")
	}
	// 把退款任务状态重置为待执行，并把 next_retry_at 置为当前时间让 worker 立即可见。
	_, err := dao.RefundTask.Ctx(ctx).
		Where(dao.RefundTask.Columns().RefundTaskNo, req.GetRefundTaskNo()).
		Data(do.RefundTask{
			Status:           int(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			NextRetryAt:      gtime.Now(),
			LastErrorCode:    strings.TrimSpace(req.GetReasonCode()),
			LastErrorMessage: "retry requested",
		}).Update()
	if err != nil {
		// 更新失败时不能返回成功，避免任务调度丢失。
		return nil, gerror.Wrap(err, "retry refund task failed")
	}
	// 回查最新任务快照，返回给调用方确认重试状态已经生效。
	task, err := s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	if err != nil {
		// 理论上刚更新成功就应该能读到任务，读不到时直接返回错误。
		return nil, err
	}
	// 返回最新退款任务对象，供后台或 worker 面板直接展示。
	return &v1.RetryRefundTaskRes{Task: toProtoRefundTask(task)}, nil
}

func (s *sAfterSale) GetAfterSaleSnapshotByNo(ctx context.Context, req *v1.GetAfterSaleSnapshotByNoReq) (*v1.GetAfterSaleSnapshotByNoRes, error) {
	// 内部快照查询接口必须指定售后单号。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		// 没有售后单号就无法返回内部快照。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	// 先读取售后主单快照。
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		// 主单读取失败时直接返回错误。
		return nil, err
	}
	// 再补充关联退款任务快照，供内部联调用同一份事实来源。
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, req.GetAfterSaleNo())
	// 把售后主单和退款任务组合成内部快照响应。
	return &v1.GetAfterSaleSnapshotByNoRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) getAfterSaleByNo(ctx context.Context, afterSaleNo string) (*entity.AfterSaleCase, error) {
	// 预留实体对象承接数据库扫描结果。
	var row entity.AfterSaleCase
	// 按业务号精确查询售后主单，这是所有售后读写流程的基础入口。
	if err := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().AfterSaleNo, afterSaleNo).Scan(&row); err != nil {
		// 查询失败时直接包装数据库错误返回，便于上层统一处理。
		return nil, gerror.Wrap(err, "query after_sale_case failed")
	}
	// 主键仍然为 0 说明数据库没有命中对应售后单。
	if row.Id == 0 {
		// 未命中时返回 not found，避免上层误把空对象当成有效记录。
		return nil, gerror.NewCode(gcode.CodeNotFound, "after_sale case not found")
	}
	// 返回命中的售后实体，供上层继续做状态判断或协议映射。
	return &row, nil
}

func (s *sAfterSale) getRefundTaskByNo(ctx context.Context, refundTaskNo string) (*entity.RefundTask, error) {
	// 预留退款任务实体承接数据库查询结果。
	var row entity.RefundTask
	// 按退款任务号精确查询退款任务，用于执行、重试和详情展示。
	if err := dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().RefundTaskNo, refundTaskNo).Scan(&row); err != nil {
		// 查询失败时包装底层错误返回，便于任务调度方感知真实原因。
		return nil, gerror.Wrap(err, "query refund_task failed")
	}
	// 主键为 0 说明当前退款任务号不存在。
	if row.Id == 0 {
		// 不存在时返回 not found，避免后续状态机空转。
		return nil, gerror.NewCode(gcode.CodeNotFound, "refund task not found")
	}
	// 返回命中的退款任务实体，供上层继续推进状态。
	return &row, nil
}

func (s *sAfterSale) getRefundTaskByAfterSaleNo(ctx context.Context, afterSaleNo string) (*entity.RefundTask, error) {
	// 预留退款任务实体承接“按售后单号反查任务”的结果。
	var row entity.RefundTask
	// 按售后单号查询退款任务，适用于详情页和审批后回显。
	if err := dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().AfterSaleNo, afterSaleNo).Scan(&row); err != nil {
		// 查询失败时包装错误返回，避免上层把数据库异常误判成“没有任务”。
		return nil, gerror.Wrap(err, "query refund task by after_sale_no failed")
	}
	// 主键为 0 说明当前售后单还没有创建任何退款任务。
	if row.Id == 0 {
		// 这里返回 nil,nil，表示“暂无关联任务”而不是查询错误。
		return nil, nil
	}
	// 返回命中的退款任务实体，供详情和内部快照复用。
	return &row, nil
}

func toProtoAfterSaleCase(row *entity.AfterSaleCase) *v1.AfterSaleCase {
	// 空实体没有可转换内容，直接返回 nil 让上层保持空值语义。
	if row == nil {
		return nil
	}
	// 预留附件证据列表，后续从 JSON 字段里反序列化出来。
	var evidence []uint64
	var selectedItemNos []string
	// 证据字段解码失败时这里忽略错误，保持主流程尽量可用。
	_ = json.Unmarshal([]byte(row.EvidenceAssetIdsJson), &evidence)
	_ = json.Unmarshal([]byte(row.SelectedItemNosJson), &selectedItemNos)
	// 把数据库实体逐字段映射成协议对象，作为对外统一返回结构。
	return &v1.AfterSaleCase{
		AfterSaleNo:          row.AfterSaleNo,
		OrderNo:              row.OrderNo,
		SubOrderNo:           row.SubOrderNo,
		ItemNo:               row.ItemNo,
		UserId:               row.UserId,
		ShopNo:               row.ShopNo,
		SpuNo:                row.SpuNo,
		SkuNo:                row.SkuNo,
		Qty:                  uint32(row.Qty),
		AfterSaleType:        v1.AfterSaleType(row.AfterSaleType),
		AfterSaleStatus:      v1.AfterSaleStatus(row.AfterSaleStatus),
		ApplyRefundAmount:    row.ApplyRefundAmount,
		ApprovedRefundAmount: row.ApprovedRefundAmount,
		ReasonCode:           row.ReasonCode,
		ReasonDesc:           row.ReasonDesc,
		EvidenceAssetIds:     evidence,
		BuyerRemark:          row.BuyerRemark,
		SellerReply:          row.SellerReply,
		RejectReasonCode:     v1.RejectReasonCode(row.RejectReasonCode),
		Version:              row.Version,
		CreatedAt:            toProtoTs(row.CreatedAt),
		UpdatedAt:            toProtoTs(row.UpdatedAt),
		ClosedAt:             toProtoTs(row.ClosedAt),
		CancelReasonCode:     row.CancelReasonCode,
		RefundBatchNo:        row.RefundBatchNo,
		ScopeCode:            row.ScopeCode,
		ReviewDeadlineAt:     toProtoTs(row.ReviewDeadlineAt),
		AutoApprovedAt:       toProtoTs(row.AutoApprovedAt),
		SelectedItemNos:      selectedItemNos,
		PaymentNo:            row.PaymentNo,
	}
}

func toProtoRefundTask(row *entity.RefundTask) *v1.RefundTask {
	// 空退款任务直接保持 nil，避免伪造一条空任务响应。
	if row == nil {
		return nil
	}
	// 把数据库退款任务实体转换成对外协议对象。
	return &v1.RefundTask{
		RefundTaskNo:           row.RefundTaskNo,
		AfterSaleNo:            row.AfterSaleNo,
		OrderNo:                row.OrderNo,
		SubOrderNo:             row.SubOrderNo,
		PayNo:                  row.PayNo,
		RefundAmount:           row.RefundAmount,
		Status:                 v1.RefundTaskStatus(row.Status),
		RetryCount:             uint32(row.RetryCount),
		NextRetryAt:            toProtoTs(row.NextRetryAt),
		LastErrorCode:          row.LastErrorCode,
		LastErrorMessage:       row.LastErrorMessage,
		CreatedAt:              toProtoTs(row.CreatedAt),
		UpdatedAt:              toProtoTs(row.UpdatedAt),
		PointsReturnAmount:     row.PointsReturnAmount,
		PointsReverseAmount:    row.PointsReverseAmount,
		PointsCashOffsetAmount: row.PointsCashOffsetAmount,
		FinalCashRefundAmount:  row.FinalCashRefundAmount,
		AccountDebtAfter:       row.AccountDebtAfter,
	}
}

func userIDFromContext(ctx context.Context) (uint64, error) {
	// 优先从 HTTP 请求上下文提取 user_id，因为买家前台和网关透传通常都走这里。
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			// 逐个兼容不同大小写和历史头字段命名，避免不同入口接入不一致。
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				// 把头部字符串解析成 uint64，后续所有归属校验都基于这个 user_id。
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil && uid > 0 {
					// 解析成功且 user_id 合法时直接返回，避免继续走后续降级分支。
					return uid, nil
				}
			}
		}
	}
	// HTTP 头里取不到时，再从 gRPC metadata 里兼容提取，支持跨服务内部调用。
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		// 逐个尝试不同 metadata key，兼容历史调用方命名差异。
		if val := md.Get(key); val != nil {
			// 用 gconv 统一把 metadata 值转成 uint64，兼容 string/int 等多种载荷。
			uid := gconv.Uint64(val)
			if uid > 0 {
				// 只要取到合法 user_id 就立即返回，后续逻辑按买家身份继续执行。
				return uid, nil
			}
		}
	}
	// 两条链路都拿不到身份时，按未授权处理，阻止匿名访问售后能力。
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

func normalizePageSize(reqSize int32) int {
	// 先把协议层 int32 分页值转成当前逻辑层更常用的 int。
	size := int(reqSize)
	// 调用方不传或传入非法值时，回退到默认分页大小。
	if size <= 0 {
		size = defaultPageSize
	}
	// 大于上限时强制截断，避免超大分页直接压垮数据库。
	if size > maxPageSize {
		size = maxPageSize
	}
	// 返回归一化后的分页大小，供列表查询复用。
	return size
}

func parseCursor(cursor string) (uint64, error) {
	// 先裁掉游标首尾空白，兼容前端拼接参数时的空格噪声。
	cursor = strings.TrimSpace(cursor)
	// 空游标表示从第一页开始查询，这里直接返回 0 作为“无游标”语义。
	if cursor == "" {
		return 0, nil
	}
	// 非空游标必须能被解析成 uint64，才能作为主键翻页游标使用。
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		// 解析失败时明确返回参数错误，避免调用方误用非法游标继续翻页。
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	// 返回解析成功的游标值，供列表接口继续做 where id < cursor 查询。
	return id, nil
}

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	// 空时间字段保持 nil，避免对外返回伪造时间戳。
	if t == nil {
		return nil
	}
	// 把 GoFrame 时间对象转换成 protobuf 时间戳，供 HTTP/gRPC 协议统一输出。
	return timestamppb.New(t.Time)
}

func generateBizNo(prefix string) string {
	// 先抓当前时间，保证业务号中的时间段和随机段来自同一时刻。
	now := time.Now()
	// 用“业务前缀 + 时间戳 + 纳秒尾部”的方式生成简单唯一业务号。
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
