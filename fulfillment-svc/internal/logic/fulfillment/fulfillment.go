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
	defaultPageSize = 20
	maxPageSize     = 100
)

type sFulfillment struct{}

func New() *sFulfillment {
	return &sFulfillment{}
}

func init() {
	service.RegisterFulfillment(New())
}

func (s *sFulfillment) CreateShipment(ctx context.Context, req *v1.CreateShipmentReq) (*v1.CreateShipmentRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/shop_no are required")
	}
	var existed entity.FulfillmentShipment
	_ = dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().SubOrderNo, req.GetSubOrderNo()).Scan(&existed)
	if existed.Id > 0 {
		return &v1.CreateShipmentRes{Shipment: toProtoShipment(ctx, &existed, nil)}, nil
	}
	userID, _ := userIDFromContextOptional(ctx)
	shipmentNo := generateBizNo("SP")
	addr := req.GetReceiverAddressStruct()
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
	row, err := s.getShipmentByNo(ctx, shipmentNo)
	if err != nil {
		return nil, err
	}
	return &v1.CreateShipmentRes{Shipment: toProtoShipment(ctx, row, nil)}, nil
}

func (s *sFulfillment) MarkShipmentShipped(ctx context.Context, req *v1.MarkShipmentShippedReq) (*v1.MarkShipmentShippedRes, error) {
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	shippedAt := protoTsToGTime(req.GetShippedAt())
	if shippedAt == nil {
		shippedAt = gtime.Now()
	}
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

func (s *sFulfillment) ListShopShipments(ctx context.Context, req *v1.ListShopShipmentsReq) (*v1.ListShopShipmentsRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().ShopNo, req.GetShopNo()).OrderDesc(dao.FulfillmentShipment.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.FulfillmentShipment.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.FulfillmentShipment.Columns().ShipmentStatus, statuses)
	}
	var rows []entity.FulfillmentShipment
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list shop shipments failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	list := make([]*v1.Shipment, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoShipment(ctx, &row, nil))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopShipmentsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sFulfillment) GetShipmentDetail(ctx context.Context, req *v1.GetShipmentDetailReq) (*v1.GetShipmentDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetShipmentNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shipment_no is required")
	}
	row, err := s.getShipmentByNo(ctx, req.GetShipmentNo())
	if err != nil {
		return nil, err
	}
	nodes, _ := s.listTrackingNodes(ctx, row.ShipmentNo)
	return &v1.GetShipmentDetailRes{Shipment: toProtoShipment(ctx, row, nodes)}, nil
}

func (s *sFulfillment) GetMyOrderLogistics(ctx context.Context, req *v1.GetMyOrderLogisticsReq) (*v1.GetMyOrderLogisticsRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	userID, _ := userIDFromContextOptional(ctx)
	model := dao.FulfillmentShipment.Ctx(ctx).Where(dao.FulfillmentShipment.Columns().OrderNo, req.GetOrderNo())
	if userID > 0 {
		model = model.Where(dao.FulfillmentShipment.Columns().UserId, userID)
	}
	var rows []entity.FulfillmentShipment
	if err := model.OrderDesc(dao.FulfillmentShipment.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query order logistics failed")
	}
	list := make([]*v1.Shipment, 0, len(rows))
	for _, row := range rows {
		nodes, _ := s.listTrackingNodes(ctx, row.ShipmentNo)
		list = append(list, toProtoShipment(ctx, &row, nodes))
	}
	return &v1.GetMyOrderLogisticsRes{Shipments: list}, nil
}

func (s *sFulfillment) IngestTrackingCallback(ctx context.Context, req *v1.IngestTrackingCallbackReq) (*v1.IngestTrackingCallbackRes, error) {
	if req == nil || strings.TrimSpace(req.GetLogisticsCompanyCode()) == "" || strings.TrimSpace(req.GetLogisticsNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "logistics_company_code/logistics_no are required")
	}
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
	changed := false
	for _, node := range req.GetNodes() {
		if strings.TrimSpace(node.GetNodeNo()) == "" {
			continue
		}
		var existed entity.FulfillmentTrackingNode
		_ = dao.FulfillmentTrackingNode.Ctx(ctx).Where(dao.FulfillmentTrackingNode.Columns().NodeNo, node.GetNodeNo()).Scan(&existed)
		if existed.Id > 0 {
			continue
		}
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
	targetStatus := shipment.ShipmentStatus
	deliveredAt := shipment.DeliveredAt
	for _, node := range req.GetNodes() {
		code := strings.ToUpper(strings.TrimSpace(node.GetStatusCode()))
		content := strings.ToUpper(strings.TrimSpace(node.GetContent()))
		if strings.Contains(code, "DELIVER") || strings.Contains(code, "SIGN") || strings.Contains(content, "签收") {
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_DELIVERED)
			deliveredAt = protoTsToGTime(node.GetEventTime())
			break
		}
		if strings.Contains(code, "EXCEPTION") || strings.Contains(code, "FAIL") {
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_EXCEPTION)
		} else if targetStatus < uint(v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT) {
			targetStatus = uint(v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT)
		}
	}
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

func (s *sFulfillment) listTrackingNodes(ctx context.Context, shipmentNo string) ([]*v1.TrackingNode, error) {
	var rows []entity.FulfillmentTrackingNode
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

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
