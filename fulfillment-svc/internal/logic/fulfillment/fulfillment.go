package fulfillment

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/fulfillment-svc/api/v1"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/dao"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// defaultPageSize 默认分页大小，避免单次请求返回量过大
	defaultPageSize = 20
	// maxPageSize 最大分页限制，避免一次拉取过多记录
	maxPageSize = 100
)

// sFulfillment 实现 fulfillment 服务的核心逻辑
type sFulfillment struct{}

// New 创建 fulfillment 逻辑 handler
func New() *sFulfillment {
	return &sFulfillment{}
}

func init() {
	// 在组件初始化阶段注册 handler，供 controller 调度
	service.RegisterFulfillment(New())
}

// CreateShipment 生成发货单并写入数据库
func (s *sFulfillment) CreateShipment(ctx context.Context, req *v1.CreateShipmentReq) (*v1.CreateShipmentRes, error) {
	// 校验请求参数
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/shop_no are required")
	}
	// 看看该子订单是否已有发货单
	var existed entity.FulfillmentShipment
	_ = dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().SubOrderNo, req.GetSubOrderNo()).Scan(&existed)
	if existed.Id > 0 {
		return &v1.CreateShipmentRes{Shipment: toProtoShipment(ctx, &existed, nil)}, nil
	}
	// 读取上下文中的用户 ID 并生成业务单号
	userID, _ := userIDFromContextOptional(ctx)
	shipmentNo := generateBizNo("SP")
	addr := req.GetReceiverAddressStruct()
	// 构建并写入发货单记录
	_, err := dao.FulfillmentShipment.Ctx(ctx).Data(do.FulfillmentShipment{
		ShipmentNo:           shipmentNo,
		OrderNo:              strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:           strings.TrimSpace(req.GetSubOrderNo()),
		ShopNo:               strings.TrimSpace(req.GetShopNo()),
		UserId:               userID,
		ShipmentStatus:       int(v1.ShipmentStatus_SHIPMENT_STATUS_WAIT_SHIP),
		ReceiverName:         strings.TrimSpace(req.GetReceiverName()),
		ReceiverPhone:        strings.TrimSpace(req.GetReceiverPhone()),
		ReceiverAddress:      strings.TrimSpace(req.GetReceiverAddress()),
		ReceiverCountryCode:  addr.GetCountryCode(),
		ReceiverProvinceCode: addr.GetProvinceCode(),
		ReceiverProvinceName: addr.GetProvinceName(),
		ReceiverCityCode:     addr.GetCityCode(),
		ReceiverCityName:     addr.GetCityName(),
		ReceiverDistrictCode: addr.GetDistrictCode(),
		ReceiverDistrictName: addr.GetDistrictName(),
		ReceiverStreet:       addr.GetStreet(),
		ReceiverDetail:       addr.GetDetail(),
		ReceiverPostalCode:   addr.GetPostalCode(),
		Version:              1,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "create shipment failed")
	}
	// 写入完成后立即回读这条记录
	row, err := s.getShipmentByNo(ctx, shipmentNo)
	if err != nil {
		return nil, err
	}
	return &v1.CreateShipmentRes{Shipment: toProtoShipment(ctx, row, nil)}, nil
}

// MarkShipmentShipped 标记发货单为已发货并更新物流信息
func (s *sFulfillment) MarkShipmentShipped(ctx context.Context, req *v1.MarkShipmentShippedReq) (*v1.MarkShipmentShippedRes, error) {
	// 校验请求携带发货单号
	// 校验参数
	// 校验发货单号
	// 校验必填字段
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	// 查出这笔发货单
	// 读取发货单
	// 复用已有方法读取订单状态
	// 读取发货单信息
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	// 验证版本号防止并发写覆盖
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	// 尝试使用请求中提供的发货时间，不存在则使用当前时间
	shippedAt := protoTsToGTime(req.GetShippedAt())
	if shippedAt == nil {
		shippedAt = gtime.Now()
	}
	// 写入物流公司、快递单号并推进状态
	_, err = dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().ShipmentNo, row.ShipmentNo).Data(do.FulfillmentShipment{
		LogisticsCompanyCode: strings.TrimSpace(req.GetLogisticsCompanyCode()),
		LogisticsCompanyName: strings.TrimSpace(req.GetLogisticsCompanyName()),
		LogisticsNo:          strings.TrimSpace(req.GetLogisticsNo()),
		ShipmentStatus:       int(v1.ShipmentStatus_SHIPMENT_STATUS_SHIPPED),
		ShippedAt:            shippedAt,
		Version:              row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "mark shipment shipped failed")
	}
	return &v1.MarkShipmentShippedRes{
		ShipmentNo:     row.ShipmentNo,
		ShipmentStatus: v1.ShipmentStatus_SHIPMENT_STATUS_SHIPPED,
	}, nil
}

// ListShopShipments 提供店铺维度的发货单分页查询
func (s *sFulfillment) ListShopShipments(ctx context.Context, req *v1.ListShopShipmentsReq) (*v1.ListShopShipmentsRes, error) {
	// 校验 shop_no 参数
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	// 规范分页大小
	pageSize := normalizePageSize(req.GetPageSize())
	// 解析 next_cursor
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	// 构造查询模型，按 ID 逆序分页
	model := dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().ShopNo, req.GetShopNo()).OrderDesc(dao.FulfillmentShipment.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		// cursor 表示 read 持续状态，从该 ID 之后继续读取
		model = model.WhereLT(dao.FulfillmentShipment.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		// 把 proto 定义的状态枚举转换成数据库可用的整型切片
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.FulfillmentShipment.Columns().ShipmentStatus, statuses)
	}
	// rows 用于收集查询结果
	var rows []entity.FulfillmentShipment
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list shop shipments failed")
	}
	// 分页候选：如果拉取的数量超过 pageSize，即存在下一页
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	// 转换实体为 proto 输出
	list := make([]*v1.Shipment, 0, len(rows))
	// 将每条节点实体转换成 proto
	for _, row := range rows {
		list = append(list, toProtoShipment(ctx, &row, nil))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		// 用最后一条的 ID 作为 next_cursor，供客户端继续分页
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopShipmentsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

// GetShipmentDetail 返回指定发货单及其轨迹节点
func (s *sFulfillment) GetShipmentDetail(ctx context.Context, req *v1.GetShipmentDetailReq) (*v1.GetShipmentDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	// 查询其轨迹节点
	nodes, _ := s.listTrackingNodes(ctx, row.ShipmentNo)
	return &v1.GetShipmentDetailRes{Shipment: toProtoShipment(ctx, row, nodes)}, nil
}

// GetMyOrderLogistics 查询当前用户在订单维度的所有发货记录
func (s *sFulfillment) GetMyOrderLogistics(ctx context.Context, req *v1.GetMyOrderLogisticsReq) (*v1.GetMyOrderLogisticsRes, error) {
	// 确保传入 order_no
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	// 尝试读取当前用户 ID 以限制自己的订单
	userID, _ := userIDFromContextOptional(ctx)
	model := dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().OrderNo, req.GetOrderNo())
	if userID > 0 {
		model = model.Where(dao.FulfillmentShipment.Columns().UserId, userID)
	}
	// rows 存储查询结果
	var rows []entity.FulfillmentShipment
	if err := model.OrderDesc(dao.FulfillmentShipment.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query order logistics failed")
	}
	list := make([]*v1.Shipment, 0, len(rows))
	// 每条记录都附带轨迹节点返回
	for _, row := range rows {
		nodes, _ := s.listTrackingNodes(ctx, row.ShipmentNo)
		list = append(list, toProtoShipment(ctx, &row, nodes))
	}
	return &v1.GetMyOrderLogisticsRes{Shipments: list}, nil
}

// IngestTrackingCallback 接收第三方物流的节点回调并更新状态
func (s *sFulfillment) IngestTrackingCallback(ctx context.Context, req *v1.IngestTrackingCallbackReq) (*v1.IngestTrackingCallbackRes, error) {
	// 验证 payload 包含物流公司代码与运单号
	if req == nil || strings.TrimSpace(req.GetLogisticsCompanyCode()) == "" || strings.TrimSpace(req.GetLogisticsNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "logistics_company_code/logistics_no are required")
	}
	// 查询当前记录确定属于哪笔发货单
	var shipment entity.FulfillmentShipment
	if err := dao.FulfillmentShipment.Ctx(ctx).
		Where(dao.FulfillmentShipment.Columns().LogisticsCompanyCode, req.GetLogisticsCompanyCode()).
		Where(dao.FulfillmentShipment.Columns().LogisticsNo, req.GetLogisticsNo()).
		Scan(&shipment); err != nil {
		return nil, gerror.Wrap(err, "query shipment by logistics failed")
	}
	if shipment.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "shipment not found")
	}
	// 记录是否插入了新的轨迹节点或刷新状态
	changed := false
	for _, node := range req.GetNodes() {
		// 跳过缺少 node_no 的无效节点
		if strings.TrimSpace(node.GetNodeNo()) == "" {
			continue
		}
		// 去重：已有 node_no 的节点跳过
		var existed entity.FulfillmentTrackingNode
		_ = dao.FulfillmentTrackingNode.Ctx(ctx).Where(dao.FulfillmentTrackingNode.Columns().NodeNo, node.GetNodeNo()).Scan(&existed)
		if existed.Id > 0 {
			continue
		}
		// 插入新节点
		_, err := dao.FulfillmentTrackingNode.Ctx(ctx).Data(do.FulfillmentTrackingNode{
			NodeNo:     strings.TrimSpace(node.GetNodeNo()),
			ShipmentNo: shipment.ShipmentNo,
			StatusCode: strings.TrimSpace(node.GetStatusCode()),
			Content:    strings.TrimSpace(node.GetContent()),
			Location:   strings.TrimSpace(node.GetLocation()),
			EventTime:  protoTsToGTime(node.GetEventTime()),
		}).Insert()
		if err == nil {
			changed = true
		}
	}
	// 以当前状态为基础，遍历节点计算最晚状态
	targetStatus := shipment.ShipmentStatus
	deliveredAt := shipment.DeliveredAt
	for _, node := range req.GetNodes() {
		// 统一把状态码与描述转换成大写，便于关键字匹配
		code := strings.ToUpper(strings.TrimSpace(node.GetStatusCode()))
		content := strings.ToUpper(strings.TrimSpace(node.GetContent()))
		// DELIVER/SIGN 或包含签收描述则视作已送达
		if strings.Contains(code, "DELIVER") || strings.Contains(code, "SIGN") || strings.Contains(content, "签收") {
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_DELIVERED)
			deliveredAt = protoTsToGTime(node.GetEventTime())
			break
		}
		// 异常/失败关键词命中则视为异常状态
		if strings.Contains(code, "EXCEPTION") || strings.Contains(code, "FAIL") {
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_EXCEPTION)
		} else if targetStatus < uint(v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT) {
			// 否则状态最低保证为在途
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT)
		}
	}
	// 状态变更时刷新数据库记录
	if targetStatus != shipment.ShipmentStatus {
		_, _ = dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().ShipmentNo, shipment.ShipmentNo).Data(do.FulfillmentShipment{
			ShipmentStatus: targetStatus,
			DeliveredAt:    deliveredAt,
		}).Update()
		changed = true
	}
	return &v1.IngestTrackingCallbackRes{
		ShipmentNo:     shipment.ShipmentNo,
		ShipmentStatus: v1.ShipmentStatus(targetStatus),
		Changed:        changed,
	}, nil
}

// SyncTrackingByShipment 查询单条发货单的当前物流状态（供同步接口使用）
func (s *sFulfillment) SyncTrackingByShipment(ctx context.Context, req *v1.SyncTrackingByShipmentReq) (*v1.SyncTrackingByShipmentRes, error) {
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	return &v1.SyncTrackingByShipmentRes{
		ShipmentNo:     row.ShipmentNo,
		ShipmentStatus: v1.ShipmentStatus(row.ShipmentStatus),
	}, nil
}

// GetShipmentSnapshot 获取发货单及其所有轨迹节点的快照
func (s *sFulfillment) GetShipmentSnapshot(ctx context.Context, req *v1.GetShipmentSnapshotReq) (*v1.GetShipmentSnapshotRes, error) {
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	nodes, _ := s.listTrackingNodes(ctx, row.ShipmentNo)
	return &v1.GetShipmentSnapshotRes{Shipment: toProtoShipment(ctx, row, nodes)}, nil
}

// getShipmentByNo 通过发货单号查发货单，找不到时返回 404
func (s *sFulfillment) getShipmentByNo(ctx context.Context, shipmentNo string) (*entity.FulfillmentShipment, error) {
	var row entity.FulfillmentShipment
	if err := dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().ShipmentNo, shipmentNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query shipment failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "shipment not found")
	}
	return &row, nil
}

// listTrackingNodes 查询指定发货单的所有轨迹节点并转为 proto
func (s *sFulfillment) listTrackingNodes(ctx context.Context, shipmentNo string) ([]*v1.TrackingNode, error) {
	var rows []entity.FulfillmentTrackingNode
	// 按事件时间顺序读取轨迹节点
	if err := dao.FulfillmentTrackingNode.Ctx(ctx).
		Where(dao.FulfillmentTrackingNode.Columns().ShipmentNo, shipmentNo).
		OrderAsc(dao.FulfillmentTrackingNode.Columns().EventTime).
		OrderAsc(dao.FulfillmentTrackingNode.Columns().Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query tracking nodes failed")
	}
	list := make([]*v1.TrackingNode, 0, len(rows))
	for _, row := range rows {
		list = append(list, &v1.TrackingNode{
			NodeNo:     row.NodeNo,
			ShipmentNo: row.ShipmentNo,
			StatusCode: row.StatusCode,
			Content:    row.Content,
			Location:   row.Location,
			EventTime:  toProtoTs(row.EventTime),
			CreatedAt:  toProtoTs(row.CreatedAt),
		})
	}
	return list, nil
}

// toProtoShipment 把实体转为 proto 包括地址和轨迹节点
func toProtoShipment(ctx context.Context, row *entity.FulfillmentShipment, nodes []*v1.TrackingNode) *v1.Shipment {
	if row == nil {
		return nil
	}
	if nodes == nil {
		nodes = []*v1.TrackingNode{}
	}
	return &v1.Shipment{
		ShipmentNo:           row.ShipmentNo,
		OrderNo:              row.OrderNo,
		SubOrderNo:           row.SubOrderNo,
		ShopNo:               row.ShopNo,
		UserId:               row.UserId,
		LogisticsCompanyCode: row.LogisticsCompanyCode,
		LogisticsCompanyName: row.LogisticsCompanyName,
		LogisticsNo:          row.LogisticsNo,
		ShipmentStatus:       v1.ShipmentStatus(row.ShipmentStatus),
		ShippedAt:            toProtoTs(row.ShippedAt),
		DeliveredAt:          toProtoTs(row.DeliveredAt),
		ReceiverName:         row.ReceiverName,
		ReceiverPhone:        row.ReceiverPhone,
		ReceiverAddress:      row.ReceiverAddress,
		Version:              row.Version,
		CreatedAt:            toProtoTs(row.CreatedAt),
		UpdatedAt:            toProtoTs(row.UpdatedAt),
		TrackingNodes:        nodes,
		ReceiverAddressStruct: &v1.ReceiverAddress{
			CountryCode:  row.ReceiverCountryCode,
			ProvinceCode: row.ReceiverProvinceCode,
			ProvinceName: row.ReceiverProvinceName,
			CityCode:     row.ReceiverCityCode,
			CityName:     row.ReceiverCityName,
			DistrictCode: row.ReceiverDistrictCode,
			DistrictName: row.ReceiverDistrictName,
			Street:       row.ReceiverStreet,
			Detail:       row.ReceiverDetail,
			PostalCode:   row.ReceiverPostalCode,
		},
	}
}

// userIDFromContextOptional 先从 HTTP header，再从 gRPC metadata 读取 user_id
func userIDFromContextOptional(ctx context.Context) (uint64, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil {
					return uid, nil
				}
				return 0, gerror.Wrap(err, "parse x-user-id failed")
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			return gconv.Uint64(val), nil
		}
	}
	return 0, nil
}

// normalizePageSize 规范分页大小，限制在合法范围内
func normalizePageSize(reqSize int32) int {
	size := int(reqSize)
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return size
}

// parseCursor 将 next_cursor 解析成 uint64 作为分页起点
func parseCursor(cursor string) (uint64, error) {
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	return id, nil
}

// toProtoTs 将 GoFrame gtime 转成 protobuf Timestamp
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// protoTsToGTime 将 protobuf Timestamp 转回 GoFrame gtime
func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

// generateBizNo 生成带前缀的唯一业务流水号
func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
