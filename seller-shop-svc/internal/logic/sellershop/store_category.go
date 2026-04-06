package sellershop

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/dao"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	storeCategoryLevelL1                   = 1
	storeCategoryLevelL2                   = 2
	outboxEventTypeStoreCategoryAssignment = "ShopStoreCategoryAssignmentChanged"
)

type storeCategoryAssignmentPayload struct {
	ShopNo            string   `json:"shop_no"`
	SpuNo             string   `json:"spu_no"`
	StoreCategoryId   uint64   `json:"store_category_id"`
	StoreCategoryL1   uint64   `json:"store_category_l1"`
	StoreCategoryL2   uint64   `json:"store_category_l2"`
	StoreCategoryPath []uint64 `json:"store_category_path"`
	UpdatedAt         string   `json:"updated_at"`
}

func (s *sSellerShop) ListStoreCategories(ctx context.Context, req *v1.ListStoreCategoriesReq) (*v1.ListStoreCategoriesRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	if shopNo == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	shop, err := queryShopByNo(ctx, shopNo)
	if err != nil {
		return nil, err
	}
	if req.GetBuyerSide() && shop.BuyerVisible != 1 {
		return &v1.ListStoreCategoriesRes{Categories: []*v1.StoreCategory{}}, nil
	}

	rows, err := listStoreCategoryRows(ctx, shopNo, req.GetBuyerSide())
	if err != nil {
		return nil, err
	}
	return &v1.ListStoreCategoriesRes{Categories: buildStoreCategoryTree(rows)}, nil
}

func (s *sSellerShop) ListBuyerStoreCategories(ctx context.Context, req *v1.ListBuyerStoreCategoriesReq) (*v1.ListBuyerStoreCategoriesRes, error) {
	out, err := s.ListStoreCategories(ctx, &v1.ListStoreCategoriesReq{ShopNo: req.GetShopNo(), BuyerSide: true})
	if err != nil {
		return nil, err
	}
	return &v1.ListBuyerStoreCategoriesRes{Categories: out.Categories}, nil
}

func (s *sSellerShop) CreateStoreCategory(ctx context.Context, req *v1.CreateStoreCategoryReq) (*v1.CreateStoreCategoryRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	name := strings.TrimSpace(req.GetName())
	if shopNo == "" || name == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no and name are required")
	}
	if _, err := queryShopByNo(ctx, shopNo); err != nil {
		return nil, err
	}

	level := storeCategoryLevelL1
	if req.GetParentId() > 0 {
		parent, err := queryStoreCategoryByID(ctx, shopNo, req.GetParentId())
		if err != nil {
			return nil, err
		}
		if parent.Level >= storeCategoryLevelL2 {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "only two levels are supported")
		}
		level = int(parent.Level) + 1
	}

	result, err := dao.ShopStoreCategory.Ctx(ctx).Data(do.ShopStoreCategory{
		ShopNo:       shopNo,
		ParentId:     req.GetParentId(),
		Name:         name,
		Level:        level,
		SortOrder:    req.GetSortOrder(),
		IsVisible:    boolToInt(req.GetIsVisible()),
		IsDeleted:    0,
		ProductCount: 0,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "create store category failed")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "read store category id failed")
	}
	row, err := queryStoreCategoryByID(ctx, shopNo, uint64(id))
	if err != nil {
		return nil, err
	}
	return &v1.CreateStoreCategoryRes{Category: toStoreCategoryProto(row)}, nil
}

func (s *sSellerShop) UpdateStoreCategory(ctx context.Context, req *v1.UpdateStoreCategoryReq) (*v1.UpdateStoreCategoryRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	if shopNo == "" || req.GetCategoryId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no and category_id are required")
	}
	row, err := queryStoreCategoryByID(ctx, shopNo, req.GetCategoryId())
	if err != nil {
		return nil, err
	}

	data := do.ShopStoreCategory{}
	if name := strings.TrimSpace(req.GetName()); name != "" {
		data.Name = name
	}
	if req.SetSortOrder {
		data.SortOrder = req.GetSortOrder()
	}
	if req.SetIsVisible {
		data.IsVisible = boolToInt(req.GetIsVisible())
	}
	if _, err = dao.ShopStoreCategory.Ctx(ctx).
		Where(dao.ShopStoreCategory.Columns().Id, row.Id).
		Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).
		Data(data).
		Update(); err != nil {
		return nil, gerror.Wrap(err, "update store category failed")
	}
	updated, err := queryStoreCategoryByID(ctx, shopNo, row.Id)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateStoreCategoryRes{Category: toStoreCategoryProto(updated)}, nil
}

func (s *sSellerShop) SortStoreCategories(ctx context.Context, req *v1.SortStoreCategoriesReq) (*v1.SortStoreCategoriesRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	if shopNo == "" || len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no and items are required")
	}
	for _, item := range req.GetItems() {
		if item == nil || item.GetCategoryId() == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "category_id is required")
		}
	}
	if err := dao.ShopStoreCategory.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range req.GetItems() {
			if _, err := tx.Model(dao.ShopStoreCategory.Table()).Where(dao.ShopStoreCategory.Columns().Id, item.GetCategoryId()).Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).Data(do.ShopStoreCategory{SortOrder: item.GetSortOrder()}).Update(); err != nil {
				return gerror.Wrap(err, "sort store categories failed")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	rows, err := listStoreCategoryRows(ctx, shopNo, false)
	if err != nil {
		return nil, err
	}
	return &v1.SortStoreCategoriesRes{Categories: buildStoreCategoryTree(rows)}, nil
}

func (s *sSellerShop) DeleteStoreCategory(ctx context.Context, req *v1.DeleteStoreCategoryReq) (*v1.DeleteStoreCategoryRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	if shopNo == "" || req.GetCategoryId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no and category_id are required")
	}
	row, err := queryStoreCategoryByID(ctx, shopNo, req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	if err = dao.ShopStoreCategory.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		childCount, err := tx.Model(dao.ShopStoreCategory.Table()).Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).Where(dao.ShopStoreCategory.Columns().ParentId, row.Id).Where(dao.ShopStoreCategory.Columns().IsDeleted, 0).Count()
		if err != nil {
			return gerror.Wrap(err, "count child categories failed")
		}
		if childCount > 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "category has child categories")
		}
		bindingCount, err := tx.Model(dao.ShopStoreCategoryProduct.Table()).Where(dao.ShopStoreCategoryProduct.Columns().ShopNo, shopNo).Where(dao.ShopStoreCategoryProduct.Columns().StoreCategoryId, row.Id).Count()
		if err != nil {
			return gerror.Wrap(err, "count category bindings failed")
		}
		if bindingCount > 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "category still has bound products")
		}
		_, err = tx.Model(dao.ShopStoreCategory.Table()).Where(dao.ShopStoreCategory.Columns().Id, row.Id).Data(do.ShopStoreCategory{IsDeleted: 1, IsVisible: 0}).Update()
		if err != nil {
			return gerror.Wrap(err, "soft delete store category failed")
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &v1.DeleteStoreCategoryRes{Ok: true}, nil
}

func (s *sSellerShop) GetProductStoreCategoryBinding(ctx context.Context, req *v1.GetProductStoreCategoryBindingReq) (*v1.GetProductStoreCategoryBindingRes, error) {
	binding, err := queryStoreCategoryBinding(ctx, strings.TrimSpace(req.GetShopNo()), strings.TrimSpace(req.GetSpuNo()))
	if err != nil {
		return nil, err
	}
	return &v1.GetProductStoreCategoryBindingRes{Binding: binding}, nil
}

func (s *sSellerShop) BatchGetProductStoreCategoryBindings(ctx context.Context, req *v1.BatchGetProductStoreCategoryBindingsReq) (*v1.BatchGetProductStoreCategoryBindingsRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	if shopNo == "" || len(req.GetSpuNos()) == 0 {
		return &v1.BatchGetProductStoreCategoryBindingsRes{Items: []*v1.ProductStoreCategoryBinding{}}, nil
	}
	spuNos := make([]string, 0, len(req.GetSpuNos()))
	for _, spuNo := range req.GetSpuNos() {
		if normalized := strings.TrimSpace(spuNo); normalized != "" {
			spuNos = append(spuNos, normalized)
		}
	}
	if len(spuNos) == 0 {
		return &v1.BatchGetProductStoreCategoryBindingsRes{Items: []*v1.ProductStoreCategoryBinding{}}, nil
	}
	var rows []*entity.ShopStoreCategoryProduct
	if err := dao.ShopStoreCategoryProduct.Ctx(ctx).Where(dao.ShopStoreCategoryProduct.Columns().ShopNo, shopNo).WhereIn(dao.ShopStoreCategoryProduct.Columns().SpuNo, spuNos).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list product category bindings failed")
	}
	items := make([]*v1.ProductStoreCategoryBinding, 0, len(rows))
	for _, row := range rows {
		item, err := buildBindingProto(ctx, shopNo, row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &v1.BatchGetProductStoreCategoryBindingsRes{Items: items}, nil
}

func (s *sSellerShop) UpdateProductStoreCategoryBinding(ctx context.Context, req *v1.UpdateProductStoreCategoryBindingReq) (*v1.UpdateProductStoreCategoryBindingRes, error) {
	shopNo := strings.TrimSpace(req.GetShopNo())
	spuNo := strings.TrimSpace(req.GetSpuNo())
	var err error
	if shopNo == "" || spuNo == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no and spu_no are required")
	}
	if _, err := queryShopByNo(ctx, shopNo); err != nil {
		return nil, err
	}

	var targetCategory *entity.ShopStoreCategory
	if req.GetStoreCategoryId() > 0 {
		targetCategory, err = queryStoreCategoryByID(ctx, shopNo, req.GetStoreCategoryId())
		if err != nil {
			return nil, err
		}
		childCount, err := dao.ShopStoreCategory.Ctx(ctx).Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).Where(dao.ShopStoreCategory.Columns().ParentId, targetCategory.Id).Where(dao.ShopStoreCategory.Columns().IsDeleted, 0).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "count child categories failed")
		}
		if childCount > 0 {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "only leaf categories can bind products")
		}
	}

	currentBinding, err := queryBindingRow(ctx, shopNo, spuNo)
	if err != nil {
		return nil, err
	}
	oldCategoryID := uint64(0)
	if currentBinding != nil {
		oldCategoryID = currentBinding.StoreCategoryId
	}
	newCategoryID := req.GetStoreCategoryId()
	if oldCategoryID == newCategoryID {
		binding, err := queryStoreCategoryBinding(ctx, shopNo, spuNo)
		if err != nil {
			return nil, err
		}
		return &v1.UpdateProductStoreCategoryBindingRes{Binding: binding}, nil
	}

	assignmentPayload, err := buildAssignmentPayload(ctx, shopNo, spuNo, targetCategory)
	if err != nil {
		return nil, err
	}
	requestID := metadataValue(ctx, "x-request-id")
	updatedAt := gtime.Now()

	if err = dao.ShopStoreCategoryProduct.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := gtime.Now()
		if currentBinding == nil && newCategoryID > 0 {
			if _, err := tx.Model(dao.ShopStoreCategoryProduct.Table()).Data(do.ShopStoreCategoryProduct{ShopNo: shopNo, SpuNo: spuNo, StoreCategoryId: newCategoryID, CreatedAt: now, UpdatedAt: now}).Insert(); err != nil {
				return gerror.Wrap(err, "create product category binding failed")
			}
		} else if currentBinding != nil && newCategoryID == 0 {
			if _, err := tx.Model(dao.ShopStoreCategoryProduct.Table()).Where(dao.ShopStoreCategoryProduct.Columns().Id, currentBinding.Id).Delete(); err != nil {
				return gerror.Wrap(err, "delete product category binding failed")
			}
		} else if currentBinding != nil {
			if _, err := tx.Model(dao.ShopStoreCategoryProduct.Table()).Where(dao.ShopStoreCategoryProduct.Columns().Id, currentBinding.Id).Data(do.ShopStoreCategoryProduct{StoreCategoryId: newCategoryID, UpdatedAt: now}).Update(); err != nil {
				return gerror.Wrap(err, "update product category binding failed")
			}
		}
		if oldCategoryID > 0 {
			if err := adjustStoreCategoryCountTx(ctx, tx, oldCategoryID, -1); err != nil {
				return err
			}
		}
		if newCategoryID > 0 {
			if err := adjustStoreCategoryCountTx(ctx, tx, newCategoryID, 1); err != nil {
				return err
			}
		}
		event := buildOutboxEvent(requestID, "SHOP_STORE_CATEGORY_PRODUCT", spuNo, assignmentPayload)
		if _, err := tx.Model(dao.SellerOutboxEvent.Table()).Data(do.SellerOutboxEvent{EventId: event.EventID, EventType: outboxEventTypeStoreCategoryAssignment, AggregateType: "SHOP_STORE_CATEGORY_PRODUCT", AggregateNo: spuNo, RequestId: requestID, PayloadJson: mustMarshalJSON(event), Status: outboxStatusNew, AvailableAt: updatedAt}).Insert(); err != nil {
			return gerror.Wrap(err, "insert category assignment outbox failed")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	binding, err := queryStoreCategoryBinding(ctx, shopNo, spuNo)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateProductStoreCategoryBindingRes{Binding: binding}, nil
}

func (s *sSellerShop) ReconcileStoreCategoryCounts(ctx context.Context, req *v1.ReconcileStoreCategoryCountsReq) (*v1.ReconcileStoreCategoryCountsRes, error) {
	counts := make(map[uint64]int)
	rows, err := dao.ShopStoreCategoryProduct.Ctx(ctx).
		Fields(dao.ShopStoreCategoryProduct.Columns().StoreCategoryId+" AS store_category_id", "COUNT(1) AS cnt").
		Group(dao.ShopStoreCategoryProduct.Columns().StoreCategoryId).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "aggregate category binding counts failed")
	}
	for _, row := range rows {
		counts[row["store_category_id"].Uint64()] = row["cnt"].Int()
	}
	var categories []*entity.ShopStoreCategory
	if err = dao.ShopStoreCategory.Ctx(ctx).Scan(&categories); err != nil {
		return nil, gerror.Wrap(err, "list store categories failed")
	}
	updated := 0
	for _, category := range categories {
		realCount := counts[category.Id]
		if int(category.ProductCount) == realCount {
			continue
		}
		if _, err = dao.ShopStoreCategory.Ctx(ctx).Where(dao.ShopStoreCategory.Columns().Id, category.Id).Data(do.ShopStoreCategory{ProductCount: realCount}).Update(); err != nil {
			return nil, gerror.Wrap(err, "reconcile store category count failed")
		}
		updated++
	}
	return &v1.ReconcileStoreCategoryCountsRes{UpdatedCategories: int32(updated)}, nil
}

func listStoreCategoryRows(ctx context.Context, shopNo string, buyerSide bool) ([]*entity.ShopStoreCategory, error) {
	model := dao.ShopStoreCategory.Ctx(ctx).
		Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).
		Where(dao.ShopStoreCategory.Columns().IsDeleted, 0)
	if buyerSide {
		model = model.Where(dao.ShopStoreCategory.Columns().IsVisible, 1)
	}
	var rows []*entity.ShopStoreCategory
	if err := model.OrderAsc(dao.ShopStoreCategory.Columns().ParentId).OrderAsc(dao.ShopStoreCategory.Columns().SortOrder).OrderAsc(dao.ShopStoreCategory.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list store categories failed")
	}
	return rows, nil
}

func buildStoreCategoryTree(rows []*entity.ShopStoreCategory) []*v1.StoreCategory {
	if len(rows) == 0 {
		return []*v1.StoreCategory{}
	}
	index := make(map[uint64]*v1.StoreCategory, len(rows))
	roots := make([]*v1.StoreCategory, 0)
	for _, row := range rows {
		node := toStoreCategoryProto(row)
		index[row.Id] = node
	}
	for _, row := range rows {
		node := index[row.Id]
		if row.ParentId == 0 {
			roots = append(roots, node)
			continue
		}
		parent := index[row.ParentId]
		if parent == nil {
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	for _, root := range roots {
		sortStoreCategoryNodes(root.Children)
		root.ProductCount = aggregateStoreCategoryCount(root)
	}
	sortStoreCategoryNodes(roots)
	return roots
}

func sortStoreCategoryNodes(nodes []*v1.StoreCategory) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].SortOrder == nodes[j].SortOrder {
			return nodes[i].Id < nodes[j].Id
		}
		return nodes[i].SortOrder < nodes[j].SortOrder
	})
}

func aggregateStoreCategoryCount(node *v1.StoreCategory) int32 {
	if node == nil {
		return 0
	}
	total := node.ProductCount
	for _, child := range node.Children {
		total += aggregateStoreCategoryCount(child)
	}
	return total
}

func queryStoreCategoryByID(ctx context.Context, shopNo string, categoryID uint64) (*entity.ShopStoreCategory, error) {
	var row entity.ShopStoreCategory
	if err := dao.ShopStoreCategory.Ctx(ctx).
		Where(dao.ShopStoreCategory.Columns().Id, categoryID).
		Where(dao.ShopStoreCategory.Columns().ShopNo, shopNo).
		Where(dao.ShopStoreCategory.Columns().IsDeleted, 0).
		Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query store category failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "store category not found")
	}
	return &row, nil
}

func queryBindingRow(ctx context.Context, shopNo, spuNo string) (*entity.ShopStoreCategoryProduct, error) {
	var row entity.ShopStoreCategoryProduct
	if err := dao.ShopStoreCategoryProduct.Ctx(ctx).Where(dao.ShopStoreCategoryProduct.Columns().ShopNo, shopNo).Where(dao.ShopStoreCategoryProduct.Columns().SpuNo, spuNo).Scan(&row); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no rows in result set") {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query product store category binding failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

func queryStoreCategoryBinding(ctx context.Context, shopNo, spuNo string) (*v1.ProductStoreCategoryBinding, error) {
	row, err := queryBindingRow(ctx, shopNo, spuNo)
	if err != nil || row == nil {
		return nil, err
	}
	return buildBindingProto(ctx, shopNo, row)
}

func buildBindingProto(ctx context.Context, shopNo string, row *entity.ShopStoreCategoryProduct) (*v1.ProductStoreCategoryBinding, error) {
	if row == nil {
		return nil, nil
	}
	category, err := queryStoreCategoryByID(ctx, shopNo, row.StoreCategoryId)
	if err != nil {
		return nil, err
	}
	payload, err := buildAssignmentPayload(ctx, shopNo, row.SpuNo, category)
	if err != nil {
		return nil, err
	}
	return &v1.ProductStoreCategoryBinding{ShopNo: shopNo, SpuNo: row.SpuNo, StoreCategoryId: payload.StoreCategoryId, StoreCategoryL1: payload.StoreCategoryL1, StoreCategoryL2: payload.StoreCategoryL2, StoreCategoryPath: payload.StoreCategoryPath, StoreCategoryName: category.Name, UpdatedAt: toProtoTimestamp(row.UpdatedAt)}, nil
}

func buildAssignmentPayload(ctx context.Context, shopNo, spuNo string, category *entity.ShopStoreCategory) (*storeCategoryAssignmentPayload, error) {
	payload := &storeCategoryAssignmentPayload{ShopNo: shopNo, SpuNo: spuNo, StoreCategoryPath: []uint64{}, UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if category == nil {
		return payload, nil
	}
	payload.StoreCategoryId = category.Id
	payload.StoreCategoryPath = append(payload.StoreCategoryPath, category.Id)
	if category.Level == storeCategoryLevelL1 {
		payload.StoreCategoryL1 = category.Id
		return payload, nil
	}
	payload.StoreCategoryL2 = category.Id
	parent, err := queryStoreCategoryByID(ctx, shopNo, category.ParentId)
	if err != nil {
		return nil, err
	}
	payload.StoreCategoryL1 = parent.Id
	payload.StoreCategoryPath = []uint64{parent.Id, category.Id}
	return payload, nil
}

func adjustStoreCategoryCountTx(ctx context.Context, tx gdb.TX, categoryID uint64, delta int) error {
	if categoryID == 0 || delta == 0 {
		return nil
	}
	cols := dao.ShopStoreCategory.Columns()
	update := do.ShopStoreCategory{}
	if delta > 0 {
		update.ProductCount = gdb.Raw(fmt.Sprintf("%s + %d", cols.ProductCount, delta))
	} else {
		update.ProductCount = gdb.Raw(fmt.Sprintf("GREATEST(%s - %d, 0)", cols.ProductCount, -delta))
	}
	_, err := tx.Model(dao.ShopStoreCategory.Table()).Where(cols.Id, categoryID).Data(update).Update()
	if err != nil {
		return gerror.Wrap(err, "adjust store category count failed")
	}
	return nil
}

func toStoreCategoryProto(row *entity.ShopStoreCategory) *v1.StoreCategory {
	if row == nil {
		return nil
	}
	return &v1.StoreCategory{Id: row.Id, ShopNo: row.ShopNo, ParentId: row.ParentId, Name: row.Name, Level: int32(row.Level), SortOrder: int32(row.SortOrder), IsVisible: row.IsVisible == 1, IsDeleted: row.IsDeleted == 1, ProductCount: int32(row.ProductCount), CreatedAt: toProtoTimestamp(row.CreatedAt), UpdatedAt: toProtoTimestamp(row.UpdatedAt), Children: []*v1.StoreCategory{}}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
