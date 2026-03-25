package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/dao"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sCatalog struct{}

// New 创建 Catalog 领域服务实例。
func New() *sCatalog {
	return &sCatalog{}
}

// init 在包加载时注册 Catalog 服务实现到 service 层。
func init() {
	service.RegisterCatalog(New())
}

// CreateProductDraft 创建 SPU 草稿，并在事务中写入 SPU 主记录与 SPU 属性。
func (s *sCatalog) CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (*v1.CreateProductDraftRes, error) {
	// 1) 参数前置校验：确保必要业务字段完整。
	if strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	if strings.TrimSpace(req.GetTitle()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "title is required")
	}
	if req.GetCategoryId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "category_id is required")
	}

	// 2) 生成业务号并序列化图片数组，便于落库到 JSON 字段。
	spuNo := generateBizNo("SPU")
	mainImagesJSON, err := marshalJSON(req.GetMainImageAssetIds())
	if err != nil {
		return nil, err
	}
	detailImagesJSON, err := marshalJSON(req.GetDetailImageAssetIds())
	if err != nil {
		return nil, err
	}

	// 3) 使用本地事务写入 SPU 主记录 + 属性明细，保证草稿一致性。
	err = dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = tx.Model(dao.CatalogSpu.Table()).Data(do.CatalogSpu{
			SpuNo:                   spuNo,
			ShopNo:                  req.GetShopNo(),
			Title:                   req.GetTitle(),
			SubTitle:                req.GetSubTitle(),
			CategoryId:              req.GetCategoryId(),
			BrandNo:                 req.GetBrandNo(),
			MainImageAssetIdsJson:   mainImagesJSON,
			DetailImageAssetIdsJson: detailImagesJSON,
			SpuStatus:               uint(v1.SpuStatus_SPU_STATUS_DRAFT),
			SpuStockStatus:          uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK),
			MinSalePrice:            uint64(0),
			MaxSalePrice:            uint64(0),
			MinMarketPrice:          uint64(0),
			MaxMarketPrice:          uint64(0),
			PublishTime:             protoTsToGTime(req.GetPublishTime()),
			Version:                 uint(1),
			ReviewStatus:            uint(v1.ReviewStatus_REVIEW_STATUS_PENDING),
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert catalog_spu failed")
		}
		return s.replaceSpuAttrsTx(ctx, tx, spuNo, req.GetSpuAttrs())
	})
	if err != nil {
		return nil, err
	}

	// 4) 回读聚合结构，保证返回值与数据库最终状态一致。
	aggregate, err := s.getProductAggregate(ctx, spuNo, false)
	if err != nil {
		return nil, err
	}
	return &v1.CreateProductDraftRes{Product: aggregate}, nil
}

// UpdateProductDraft 按 FieldMask 更新草稿字段，使用版本号避免并发覆盖。
func (s *sCatalog) UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (*v1.UpdateProductDraftRes, error) {
	// 1) 校验业务号、期望版本和补丁对象。
	if strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	if req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}
	if req.GetPatch() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "patch is required")
	}

	var (
		patch          = req.GetPatch()
		updateMask     = req.GetUpdateMask()
		paths          []string
		updateData     do.CatalogSpu
		replaceSpuAttr bool
	)

	// 2) 若未传 update_mask，则按默认可编辑字段全量处理。
	if updateMask == nil || len(updateMask.GetPaths()) == 0 {
		paths = []string{
			"title", "sub_title", "category_id", "brand_no",
			"main_image_asset_ids", "detail_image_asset_ids", "spu_attrs", "publish_time",
		}
	} else {
		paths = updateMask.GetPaths()
	}

	// 3) 将 patch 按路径映射到 DAO 对象，仅更新被允许的字段。
	for _, p := range paths {
		path := normalizePatchPath(p)
		switch path {
		case "title":
			updateData.Title = patch.GetTitle()
		case "sub_title":
			updateData.SubTitle = patch.GetSubTitle()
		case "category_id":
			updateData.CategoryId = patch.GetCategoryId()
		case "brand_no":
			updateData.BrandNo = patch.GetBrandNo()
		case "main_image_asset_ids":
			j, err := marshalJSON(patch.GetMainImageAssetIds())
			if err != nil {
				return nil, err
			}
			updateData.MainImageAssetIdsJson = j
		case "detail_image_asset_ids":
			j, err := marshalJSON(patch.GetDetailImageAssetIds())
			if err != nil {
				return nil, err
			}
			updateData.DetailImageAssetIdsJson = j
		case "spu_attrs":
			replaceSpuAttr = true
		case "publish_time":
			updateData.PublishTime = protoTsToGTime(patch.GetPublishTime())
		default:
			return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "unsupported update_mask path: %s", p)
		}
	}

	// 4) 事务更新：版本号乐观锁 + 可编辑状态校验 + 属性重建（可选）。
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		updateData.Version = uint(req.GetExpectedVersion()) + 1
		result, err := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
			Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
			WhereNull(dao.CatalogSpu.Columns().DeletedAt).
			WhereIn(dao.CatalogSpu.Columns().SpuStatus, []uint{
				uint(v1.SpuStatus_SPU_STATUS_DRAFT),
				uint(v1.SpuStatus_SPU_STATUS_REJECTED),
				uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			}).
			Data(updateData).
			Update()
		if err != nil {
			return gerror.Wrap(err, "update catalog_spu failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not editable")
		}
		if replaceSpuAttr {
			return s.replaceSpuAttrsTx(ctx, tx, req.GetSpuNo(), patch.GetSpuAttrs())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 5) 更新后重新回读，返回最新聚合结果。
	aggregate, err := s.getProductAggregate(ctx, req.GetSpuNo(), false)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateProductDraftRes{Product: aggregate}, nil
}

// UpsertSkuDrafts 批量新增/更新 SKU 草稿，并在 replace_all 时软删未提交的旧 SKU。
func (s *sCatalog) UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (*v1.UpsertSkuDraftsRes, error) {
	// 1) 基础参数校验。
	if strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	if req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}
	if len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items is required")
	}

	// 2) 在单事务内处理：SPU 版本递增、SKU upsert、属性明细重建、聚合重算。
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
			Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
			WhereNull(dao.CatalogSpu.Columns().DeletedAt).
			Data(do.CatalogSpu{
				Version: uint(req.GetExpectedVersion()) + 1,
			}).
			Update()
		if err != nil {
			return gerror.Wrap(err, "bump spu version failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict")
		}

		// 查询当前有效 SKU，后续用于更新判定与 replace_all 对比。
		var existingSkus []*entity.CatalogSku
		err = tx.Model(dao.CatalogSku.Table()).
			Where(dao.CatalogSku.Columns().SpuNo, req.GetSpuNo()).
			WhereNull(dao.CatalogSku.Columns().DeletedAt).
			Scan(&existingSkus)
		if err != nil {
			return gerror.Wrap(err, "query existing sku failed")
		}
		existingByNo := make(map[string]*entity.CatalogSku, len(existingSkus))
		for _, row := range existingSkus {
			existingByNo[row.SkuNo] = row
		}

		// 逐条处理本次传入 SKU：生成规格哈希、执行新增/更新。
		touched := make(map[string]struct{}, len(req.GetItems()))
		for _, item := range req.GetItems() {
			skuNo := strings.TrimSpace(item.GetSkuNo())
			if skuNo == "" {
				skuNo = generateBizNo("SKU")
			}
			saleSpecsJSON, saleSpecsHash, err := normalizeSaleAttrs(item.GetSaleAttrs())
			if err != nil {
				return err
			}

			if old, ok := existingByNo[skuNo]; ok {
				// SKU 脱离草稿后不允许改销售规格，避免与库存、购物车语义冲突。
				if old.SkuStatus != uint(v1.SkuStatus_SKU_STATUS_DRAFT) && old.SaleSpecsHash != saleSpecsHash {
					return gerror.NewCodef(gcode.CodeInvalidParameter, "sku %s sale attrs are immutable after draft", skuNo)
				}
				_, err = tx.Model(dao.CatalogSku.Table()).
					Where(dao.CatalogSku.Columns().SkuNo, skuNo).
					Data(do.CatalogSku{
						SkuName:         item.GetSkuName(),
						SkuImageAssetId: item.GetSkuImageAssetId(),
						SkuStatus:       uint(item.GetSkuStatus()),
						SalePrice:       item.GetSalePrice(),
						MarketPrice:     item.GetMarketPrice(),
						SaleSpecsJson:   saleSpecsJSON,
						SaleSpecsHash:   saleSpecsHash,
						SortOrder:       item.GetSortOrder(),
						Version:         old.Version + 1,
					}).
					Update()
				if err != nil {
					return gerror.Wrapf(err, "update sku %s failed", skuNo)
				}
			} else {
				_, err = tx.Model(dao.CatalogSku.Table()).
					Data(do.CatalogSku{
						SkuNo:           skuNo,
						SpuNo:           req.GetSpuNo(),
						SkuName:         item.GetSkuName(),
						SkuImageAssetId: item.GetSkuImageAssetId(),
						SkuStatus:       uint(item.GetSkuStatus()),
						StockStatus:     uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK),
						StockVersion:    uint64(0),
						SalePrice:       item.GetSalePrice(),
						MarketPrice:     item.GetMarketPrice(),
						SaleSpecsJson:   saleSpecsJSON,
						SaleSpecsHash:   saleSpecsHash,
						SortOrder:       item.GetSortOrder(),
						Version:         uint(1),
					}).
					Insert()
				if err != nil {
					return gerror.Wrapf(err, "insert sku %s failed", skuNo)
				}
			}

			// 重建 SKU 销售属性明细，确保明细表与主表 JSON 一致。
			if err = s.replaceSkuSaleAttrsTx(ctx, tx, req.GetSpuNo(), skuNo, item.GetSaleAttrs()); err != nil {
				return err
			}
			touched[skuNo] = struct{}{}
		}

		// replace_all=true 时，未被本次触达的旧 SKU 会被软删除。
		if req.GetReplaceAll() {
			for _, old := range existingSkus {
				if _, ok := touched[old.SkuNo]; ok {
					continue
				}
				_, err = tx.Model(dao.CatalogSku.Table()).
					Where(dao.CatalogSku.Columns().SkuNo, old.SkuNo).
					WhereNull(dao.CatalogSku.Columns().DeletedAt).
					Data(do.CatalogSku{
						SkuStatus: uint(v1.SkuStatus_SKU_STATUS_DELETED),
						DeletedAt: gtime.Now(),
						Version:   old.Version + 1,
					}).
					Update()
				if err != nil {
					return gerror.Wrapf(err, "soft delete sku %s failed", old.SkuNo)
				}
			}
		}
		// SKU 变更后重算 SPU 聚合字段（价格区间/库存状态）。
		_, err = s.recomputeSpuAggregationTx(ctx, tx, req.GetSpuNo())
		return err
	})
	if err != nil {
		return nil, err
	}

	aggregate, err := s.getProductAggregate(ctx, req.GetSpuNo(), false)
	if err != nil {
		return nil, err
	}
	return &v1.UpsertSkuDraftsRes{
		Skus:       aggregate.GetSkus(),
		SpuVersion: aggregate.GetSpu().GetVersion(),
	}, nil
}

// SubmitProductReview 将商品从草稿态提交到审核流程。
func (s *sCatalog) SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (*v1.SubmitProductReviewRes, error) {
	return s.submitReview(ctx, req.GetSpuNo(), req.GetExpectedVersion(), req.GetSubmitNote(), false)
}

// ResubmitProductReview 将驳回商品重新提交审核。
func (s *sCatalog) ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (*v1.ResubmitProductReviewRes, error) {
	res, err := s.submitReview(ctx, req.GetSpuNo(), req.GetExpectedVersion(), req.GetSubmitNote(), true)
	if err != nil {
		return nil, err
	}
	return &v1.ResubmitProductReviewRes{
		SpuNo:     res.GetSpuNo(),
		SpuStatus: res.GetSpuStatus(),
	}, nil
}

// SetProductOnShelf 执行立即上架或定时上架状态切换。
func (s *sCatalog) SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (*v1.SetProductOnShelfRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	// 默认按“立即上架”处理；若传了未来时间，则仅更新发布计划，状态保持 APPROVED。
	nextStatus := v1.SpuStatus_SPU_STATUS_ON_SHELF
	publishTime := gtime.Now()
	if ts := req.GetEffectivePublishTime(); ts != nil {
		publishTime = gtime.NewFromTime(ts.AsTime())
		if ts.AsTime().After(time.Now()) {
			nextStatus = v1.SpuStatus_SPU_STATUS_APPROVED
		}
	}

	// 乐观锁更新：只有版本匹配且状态允许（APPROVED/OFF_SHELF）才允许切换上架状态。
	result, err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		WhereIn(dao.CatalogSpu.Columns().SpuStatus, []uint{
			uint(v1.SpuStatus_SPU_STATUS_APPROVED),
			uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
		}).
		Data(do.CatalogSpu{
			SpuStatus:   uint(nextStatus),
			PublishTime: publishTime,
			Version:     uint(req.GetExpectedVersion()) + 1,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "set product on shelf failed")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not allowed")
	}
	// 返回 nextStatus，调用方可直接据此判断是“已上架”还是“待定时生效”。
	return &v1.SetProductOnShelfRes{SpuNo: req.GetSpuNo(), SpuStatus: nextStatus}, nil
}

// SetProductOffShelf 执行主动下架。
func (s *sCatalog) SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (*v1.SetProductOffShelfRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	// 下架同样走版本条件，避免多端操作互相覆盖。
	result, err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Data(do.CatalogSpu{
			SpuStatus: uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			Version:   uint(req.GetExpectedVersion()) + 1,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "set product off shelf failed")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict")
	}
	return &v1.SetProductOffShelfRes{
		SpuNo:     req.GetSpuNo(),
		SpuStatus: v1.SpuStatus_SPU_STATUS_OFF_SHELF,
	}, nil
}

// DeleteProductDraft 软删草稿商品及其 SKU。
func (s *sCatalog) DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (*emptypb.Empty, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	now := gtime.Now()
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先软删 SPU 主记录；只允许在可删除状态删除（草稿/驳回/下架）。
		result, err := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
			Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
			WhereNull(dao.CatalogSpu.Columns().DeletedAt).
			WhereIn(dao.CatalogSpu.Columns().SpuStatus, []uint{
				uint(v1.SpuStatus_SPU_STATUS_DRAFT),
				uint(v1.SpuStatus_SPU_STATUS_REJECTED),
				uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			}).
			Data(do.CatalogSpu{
				SpuStatus: uint(v1.SpuStatus_SPU_STATUS_DELETED),
				DeletedAt: now,
				Version:   uint(req.GetExpectedVersion()) + 1,
			}).Update()
		if err != nil {
			return gerror.Wrap(err, "delete spu failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not allowed")
		}

		// 再批量软删该 SPU 下全部 SKU，避免出现“主商品删了但 SKU 还在”的脏读。
		_, err = tx.Model(dao.CatalogSku.Table()).
			Where(dao.CatalogSku.Columns().SpuNo, req.GetSpuNo()).
			WhereNull(dao.CatalogSku.Columns().DeletedAt).
			Data(do.CatalogSku{
				SkuStatus: uint(v1.SkuStatus_SKU_STATUS_DELETED),
				DeletedAt: now,
			}).Update()
		if err != nil {
			return gerror.Wrap(err, "delete skus failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetMyProduct 返回卖家视角的商品聚合详情与审核信息。
func (s *sCatalog) GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (*v1.GetMyProductRes, error) {
	// 商家视角读取完整聚合（SPU+SKU+属性），不做买家可见性裁剪。
	aggregate, err := s.getProductAggregate(ctx, req.GetSpuNo(), false)
	if err != nil {
		return nil, err
	}
	spu, err := s.getSpuEntity(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetMyProductRes{
		Product: aggregate,
		Review:  spuReviewInfo(spu),
	}, nil
}

// ListMyProducts 分页查询卖家商品列表。
func (s *sCatalog) ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (*v1.ListMyProductsRes, error) {
	page, pageSize := normalizePage(req.GetPage(), req.GetPageSize())
	// 仅过滤已软删数据，其他状态由调用方自行传 statuses 筛选。
	model := dao.CatalogSpu.Ctx(ctx).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt)

	if len(req.GetStatuses()) > 0 {
		statuses := make([]uint, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			statuses = append(statuses, uint(st))
		}
		model = model.WhereIn(dao.CatalogSpu.Columns().SpuStatus, statuses)
	}
	if kw := strings.TrimSpace(req.GetKeyword()); kw != "" {
		// 目前按标题做模糊匹配，后续可切到检索索引服务。
		model = model.WhereLike(dao.CatalogSpu.Columns().Title, "%"+kw+"%")
	}
	total, err := model.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count my products failed")
	}

	var rows []*entity.CatalogSpu
	err = model.Page(page, pageSize).OrderDesc(dao.CatalogSpu.Columns().UpdatedAt).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "list my products failed")
	}

	items := make([]*v1.ProductSpu, 0, len(rows))
	for _, row := range rows {
		items = append(items, toProtoSpu(row, nil))
	}
	return &v1.ListMyProductsRes{
		Products: items,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    int64(total),
	}, nil
}

// ListReviewTasks 分页查询审核任务列表。
func (s *sCatalog) ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (*v1.ListReviewTasksRes, error) {
	page, pageSize := normalizePage(req.GetPage(), req.GetPageSize())
	model := dao.CatalogReviewTask.Ctx(ctx)
	if kw := strings.TrimSpace(req.GetKeyword()); kw != "" {
		model = model.WhereLike(dao.CatalogReviewTask.Columns().SpuNo, "%"+kw+"%")
	}
	if len(req.GetStatuses()) > 0 {
		// 入参是 SPU 业务状态，这里先映射为 review_task 的审核状态字段再筛选。
		reviewStatuses := spuStatusToReviewStatus(req.GetStatuses())
		if len(reviewStatuses) > 0 {
			model = model.WhereIn(dao.CatalogReviewTask.Columns().ReviewStatus, reviewStatuses)
		}
	}
	total, err := model.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count review tasks failed")
	}

	var rows []*entity.CatalogReviewTask
	err = model.Page(page, pageSize).OrderDesc(dao.CatalogReviewTask.Columns().SubmittedAt).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "list review tasks failed")
	}
	tasks := make([]*v1.ReviewTask, 0, len(rows))
	for _, row := range rows {
		task := &v1.ReviewTask{
			TaskNo:      row.TaskNo,
			SpuNo:       row.SpuNo,
			ShopNo:      row.ShopNo,
			SpuStatus:   reviewStatusToSpuStatus(row.ReviewStatus),
			SubmittedAt: toProtoTs(row.SubmittedAt),
		}
		// 任务表是审核快照，标题以 SPU 实时数据补全，失败不阻断主流程。
		spu, _ := s.getSpuEntity(ctx, row.SpuNo)
		if spu != nil {
			task.Title = spu.Title
			task.SpuStatus = v1.SpuStatus(spu.SpuStatus)
		}
		tasks = append(tasks, task)
	}
	return &v1.ListReviewTasksRes{
		Tasks:    tasks,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    int64(total),
	}, nil
}

// GetReviewDetail 查询审核详情（商品聚合 + 审核信息）。
func (s *sCatalog) GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (*v1.GetReviewDetailRes, error) {
	aggregate, err := s.getProductAggregate(ctx, req.GetSpuNo(), false)
	if err != nil {
		return nil, err
	}
	spu, err := s.getSpuEntity(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetReviewDetailRes{
		Product: aggregate,
		Review:  spuReviewInfo(spu),
	}, nil
}

// ApproveProduct 通过审核并更新审核任务状态。
func (s *sCatalog) ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (*v1.ApproveProductRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	now := gtime.Now()
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 审核通过：REVIEWING -> APPROVED（待上架），并写入审核意见。
		result, err := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
			Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
			Where(dao.CatalogSpu.Columns().SpuStatus, uint(v1.SpuStatus_SPU_STATUS_REVIEWING)).
			Data(do.CatalogSpu{
				SpuStatus:     uint(v1.SpuStatus_SPU_STATUS_APPROVED),
				ReviewStatus:  uint(v1.ReviewStatus_REVIEW_STATUS_APPROVED),
				ReviewedAt:    now,
				ReviewComment: req.GetReviewComment(),
				Version:       uint(req.GetExpectedVersion()) + 1,
			}).Update()
		if err != nil {
			return gerror.Wrap(err, "approve spu failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not reviewing")
		}
		// 同事务收口审核任务，避免“商品已通过但任务仍 pending”。
		return s.finishLatestReviewTaskTx(ctx, tx, req.GetSpuNo(), uint(v1.ReviewStatus_REVIEW_STATUS_APPROVED), "", "", req.GetReviewComment())
	})
	if err != nil {
		return nil, err
	}
	return &v1.ApproveProductRes{
		SpuNo:     req.GetSpuNo(),
		SpuStatus: v1.SpuStatus_SPU_STATUS_APPROVED,
	}, nil
}

// RejectProduct 驳回审核并记录驳回原因。
func (s *sCatalog) RejectProduct(ctx context.Context, req *v1.RejectProductReq) (*v1.RejectProductRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	now := gtime.Now()
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 审核驳回：记录驳回码和驳回说明，便于商家定位整改点。
		result, err := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
			Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
			Where(dao.CatalogSpu.Columns().SpuStatus, uint(v1.SpuStatus_SPU_STATUS_REVIEWING)).
			Data(do.CatalogSpu{
				SpuStatus:        uint(v1.SpuStatus_SPU_STATUS_REJECTED),
				ReviewStatus:     uint(v1.ReviewStatus_REVIEW_STATUS_REJECTED),
				ReviewedAt:       now,
				RejectReasonCode: req.GetRejectReasonCode(),
				RejectComment:    req.GetRejectComment(),
				Version:          uint(req.GetExpectedVersion()) + 1,
			}).Update()
		if err != nil {
			return gerror.Wrap(err, "reject spu failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not reviewing")
		}
		return s.finishLatestReviewTaskTx(ctx, tx, req.GetSpuNo(), uint(v1.ReviewStatus_REVIEW_STATUS_REJECTED), req.GetRejectReasonCode(), req.GetRejectComment(), "")
	})
	if err != nil {
		return nil, err
	}
	return &v1.RejectProductRes{
		SpuNo:     req.GetSpuNo(),
		SpuStatus: v1.SpuStatus_SPU_STATUS_REJECTED,
	}, nil
}

// FreezeProduct 冻结商品，阻止继续售卖。
func (s *sCatalog) FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (*v1.FreezeProductRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	// 冻结用于风控/运营强制处置，不要求当前必须是上架态。
	result, err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Data(do.CatalogSpu{
			SpuStatus: uint(v1.SpuStatus_SPU_STATUS_FROZEN),
			Version:   uint(req.GetExpectedVersion()) + 1,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "freeze product failed")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict")
	}
	return &v1.FreezeProductRes{SpuNo: req.GetSpuNo(), SpuStatus: v1.SpuStatus_SPU_STATUS_FROZEN}, nil
}

// UnfreezeProduct 解除冻结并回到下架态。
func (s *sCatalog) UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (*v1.UnfreezeProductRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	// 仅允许从 FROZEN 解冻回 OFF_SHELF，防止越权直接恢复售卖。
	result, err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
		Where(dao.CatalogSpu.Columns().SpuStatus, uint(v1.SpuStatus_SPU_STATUS_FROZEN)).
		Data(do.CatalogSpu{
			SpuStatus: uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			Version:   uint(req.GetExpectedVersion()) + 1,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "unfreeze product failed")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not frozen")
	}
	return &v1.UnfreezeProductRes{SpuNo: req.GetSpuNo(), SpuStatus: v1.SpuStatus_SPU_STATUS_OFF_SHELF}, nil
}

// ForceOffShelf 运营强制下架商品。
func (s *sCatalog) ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (*v1.ForceOffShelfRes, error) {
	if strings.TrimSpace(req.GetSpuNo()) == "" || req.GetExpectedVersion() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	// 平台强制下架入口，状态直接切 OFF_SHELF。
	result, err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.CatalogSpu.Columns().Version, req.GetExpectedVersion()).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Data(do.CatalogSpu{
			SpuStatus: uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			Version:   uint(req.GetExpectedVersion()) + 1,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "force off shelf failed")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict")
	}
	return &v1.ForceOffShelfRes{SpuNo: req.GetSpuNo(), SpuStatus: v1.SpuStatus_SPU_STATUS_OFF_SHELF}, nil
}

// GetProductDetail 返回买家可见的商品详情。
func (s *sCatalog) GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (*v1.GetProductDetailRes, error) {
	aggregate, err := s.getProductAggregate(ctx, req.GetSpuNo(), true)
	if err != nil {
		return nil, err
	}
	if aggregate.GetSpu().GetSpuStatus() != v1.SpuStatus_SPU_STATUS_ON_SHELF {
		// 买家侧强约束：非上架商品统一按不存在处理，避免暴露后台状态。
		return nil, gerror.NewCode(gcode.CodeNotFound, "product not on shelf")
	}
	return &v1.GetProductDetailRes{Product: aggregate}, nil
}

// ListProducts 返回买家商品列表（按分类/排序）。
func (s *sCatalog) ListProducts(ctx context.Context, req *v1.ListProductsReq) (*v1.ListProductsRes, error) {
	return s.listBuyerProducts(ctx, req.GetCategoryId(), "", req.GetSortBy(), req.GetPage(), req.GetPageSize())
}

// SearchProducts 返回买家商品搜索结果。
func (s *sCatalog) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error) {
	res, err := s.listBuyerProducts(ctx, req.GetCategoryId(), req.GetKeyword(), req.GetSortBy(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return &v1.SearchProductsRes{
		Items:    res.GetItems(),
		Page:     res.GetPage(),
		PageSize: res.GetPageSize(),
		Total:    res.GetTotal(),
	}, nil
}

// BatchGetSpuByNo 供内部服务按 spu_no 批量查询 SPU。
func (s *sCatalog) BatchGetSpuByNo(ctx context.Context, req *v1.BatchGetSpuByNoReq) (*v1.BatchGetSpuByNoRes, error) {
	if len(req.GetSpuNos()) == 0 {
		return &v1.BatchGetSpuByNoRes{Spus: []*v1.ProductSpu{}}, nil
	}
	var rows []*entity.CatalogSpu
	err := dao.CatalogSpu.Ctx(ctx).
		WhereIn(dao.CatalogSpu.Columns().SpuNo, req.GetSpuNos()).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "batch get spu failed")
	}
	out := make([]*v1.ProductSpu, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProtoSpu(row, nil))
	}
	return &v1.BatchGetSpuByNoRes{Spus: out}, nil
}

// BatchGetSkuByNo 供内部服务按 sku_no 批量查询 SKU。
func (s *sCatalog) BatchGetSkuByNo(ctx context.Context, req *v1.BatchGetSkuByNoReq) (*v1.BatchGetSkuByNoRes, error) {
	if len(req.GetSkuNos()) == 0 {
		return &v1.BatchGetSkuByNoRes{Skus: []*v1.ProductSku{}}, nil
	}
	var rows []*entity.CatalogSku
	err := dao.CatalogSku.Ctx(ctx).
		WhereIn(dao.CatalogSku.Columns().SkuNo, req.GetSkuNos()).
		WhereNull(dao.CatalogSku.Columns().DeletedAt).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "batch get sku failed")
	}
	saleAttrsMap, err := s.querySkuSaleAttrs(ctx, skuNoList(rows))
	if err != nil {
		return nil, err
	}
	out := make([]*v1.ProductSku, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProtoSku(row, saleAttrsMap[row.SkuNo]))
	}
	return &v1.BatchGetSkuByNoRes{Skus: out}, nil
}

// GetSkuSnapshotForOrder 为订单创建提供不可变商品快照。
func (s *sCatalog) GetSkuSnapshotForOrder(ctx context.Context, req *v1.GetSkuSnapshotForOrderReq) (*v1.GetSkuSnapshotForOrderRes, error) {
	if len(req.GetSkuNos()) == 0 {
		return &v1.GetSkuSnapshotForOrderRes{Snapshots: []*v1.SkuOrderSnapshot{}}, nil
	}
	var skus []*entity.CatalogSku
	err := dao.CatalogSku.Ctx(ctx).
		WhereIn(dao.CatalogSku.Columns().SkuNo, req.GetSkuNos()).
		WhereNull(dao.CatalogSku.Columns().DeletedAt).
		Scan(&skus)
	if err != nil {
		return nil, gerror.Wrap(err, "query skus for snapshot failed")
	}
	if len(skus) == 0 {
		return &v1.GetSkuSnapshotForOrderRes{Snapshots: []*v1.SkuOrderSnapshot{}}, nil
	}
	spuNos := uniqueSpuNosFromSkus(skus)
	var spus []*entity.CatalogSpu
	err = dao.CatalogSpu.Ctx(ctx).
		WhereIn(dao.CatalogSpu.Columns().SpuNo, spuNos).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Scan(&spus)
	if err != nil {
		return nil, gerror.Wrap(err, "query spus for snapshot failed")
	}
	spuMap := make(map[string]*entity.CatalogSpu, len(spus))
	for _, row := range spus {
		spuMap[row.SpuNo] = row
	}
	saleAttrsMap, err := s.querySkuSaleAttrs(ctx, skuNoList(skus))
	if err != nil {
		return nil, err
	}

	snapshots := make([]*v1.SkuOrderSnapshot, 0, len(skus))
	for _, sku := range skus {
		spu := spuMap[sku.SpuNo]
		if spu == nil {
			// 正常不会发生，保险起见跳过异常脏数据，避免整体请求失败。
			continue
		}
		snapshots = append(snapshots, &v1.SkuOrderSnapshot{
			SkuNo:           sku.SkuNo,
			SpuNo:           sku.SpuNo,
			ShopNo:          spu.ShopNo,
			SpuTitle:        spu.Title,
			SkuName:         sku.SkuName,
			SkuImageAssetId: sku.SkuImageAssetId,
			SalePrice:       sku.SalePrice,
			MarketPrice:     sku.MarketPrice,
			SaleAttrs:       saleAttrsMap[sku.SkuNo],
			SnapshotVersion: uint32(sku.Version),
		})
	}
	return &v1.GetSkuSnapshotForOrderRes{Snapshots: snapshots}, nil
}

// UpsertSkuStockProjection 写入库存投影，仅接受更高 stock_version 的事件。
func (s *sCatalog) UpsertSkuStockProjection(ctx context.Context, req *v1.UpsertSkuStockProjectionReq) (*v1.UpsertSkuStockProjectionRes, error) {
	var (
		updatedRows uint64
		skippedRows uint64
		spuSet      = make(map[string]struct{})
	)
	for _, item := range req.GetItems() {
		if strings.TrimSpace(item.GetSkuNo()) == "" {
			skippedRows++
			continue
		}
		// 关键点：只接受更高版本库存事件，天然抵御 MQ 乱序与重复投递。
		result, err := dao.CatalogSku.Ctx(ctx).
			Where(dao.CatalogSku.Columns().SkuNo, item.GetSkuNo()).
			Where(fmt.Sprintf("%s < ?", dao.CatalogSku.Columns().StockVersion), item.GetStockVersion()).
			Data(do.CatalogSku{
				StockStatus:      uint(item.GetStockStatus()),
				StockVersion:     item.GetStockVersion(),
				LastStockEventId: item.GetSourceEventId(),
			}).
			Update()
		if err != nil {
			return nil, gerror.Wrap(err, "upsert stock projection failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			// 版本未前进说明是旧消息或重复消息，计入 skipped 即可。
			skippedRows++
			continue
		}
		updatedRows++
		if item.GetSpuNo() != "" {
			spuSet[item.GetSpuNo()] = struct{}{}
		}
	}
	for spuNo := range spuSet {
		// 只重算受影响 SPU，避免全表重算带来的写放大。
		if _, err := s.recomputeSpuAggregation(ctx, spuNo); err != nil {
			return nil, err
		}
	}
	return &v1.UpsertSkuStockProjectionRes{UpdatedRows: updatedRows, SkippedRows: skippedRows}, nil
}

// RecomputeSpuAggregation 主动触发 SPU 聚合字段重算。
func (s *sCatalog) RecomputeSpuAggregation(ctx context.Context, req *v1.RecomputeSpuAggregationReq) (*v1.RecomputeSpuAggregationRes, error) {
	updated, err := s.recomputeSpuAggregation(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	return &v1.RecomputeSpuAggregationRes{
		SpuNo:   req.GetSpuNo(),
		Updated: updated,
	}, nil
}

// submitReview 封装提审事务：改状态、写提交时间、落审核任务。
func (s *sCatalog) submitReview(ctx context.Context, spuNo string, expectedVersion uint32, submitNote string, onlyRejected bool) (*v1.SubmitProductReviewRes, error) {
	if strings.TrimSpace(spuNo) == "" || expectedVersion == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no and expected_version are required")
	}
	now := gtime.Now()
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先定位目标 SPU，并带 expectedVersion 乐观锁，避免并发提交覆盖。
		model := tx.Model(dao.CatalogSpu.Table()).
			Where(dao.CatalogSpu.Columns().SpuNo, spuNo).
			Where(dao.CatalogSpu.Columns().Version, expectedVersion).
			WhereNull(dao.CatalogSpu.Columns().DeletedAt)
		if onlyRejected {
			// 重提入口只允许 REJECTED 状态进入，防止绕过常规审核流程。
			model = model.Where(dao.CatalogSpu.Columns().SpuStatus, uint(v1.SpuStatus_SPU_STATUS_REJECTED))
		} else {
			model = model.WhereIn(dao.CatalogSpu.Columns().SpuStatus, []uint{
				uint(v1.SpuStatus_SPU_STATUS_DRAFT),
				uint(v1.SpuStatus_SPU_STATUS_REJECTED),
				uint(v1.SpuStatus_SPU_STATUS_OFF_SHELF),
			})
		}
		result, err := model.Data(do.CatalogSpu{
			SpuStatus:        uint(v1.SpuStatus_SPU_STATUS_REVIEWING),
			ReviewStatus:     uint(v1.ReviewStatus_REVIEW_STATUS_PENDING),
			SubmittedAt:      now,
			RejectReasonCode: "",
			RejectComment:    "",
			Version:          uint(expectedVersion) + 1,
		}).Update()
		if err != nil {
			return gerror.Wrap(err, "submit review failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "spu version conflict or status not submittable")
		}

		// 写一条审核任务快照，供审核后台按任务流处理。
		spu, err := s.getSpuEntityTx(ctx, tx, spuNo)
		if err != nil {
			return err
		}
		_, err = tx.Model(dao.CatalogReviewTask.Table()).Data(do.CatalogReviewTask{
			TaskNo:             generateBizNo("RVW"),
			SpuNo:              spuNo,
			ShopNo:             spu.ShopNo,
			SpuVersionAtSubmit: uint(expectedVersion),
			SubmitNote:         submitNote,
			ReviewStatus:       uint(v1.ReviewStatus_REVIEW_STATUS_PENDING),
			SubmittedAt:        now,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert review task failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.SubmitProductReviewRes{
		SpuNo:     spuNo,
		SpuStatus: v1.SpuStatus_SPU_STATUS_REVIEWING,
	}, nil
}

// listBuyerProducts 封装买家列表查询和排序逻辑。
func (s *sCatalog) listBuyerProducts(ctx context.Context, categoryID uint64, keyword string, sortBy v1.SortBy, pageReq, pageSizeReq int32) (*v1.ListProductsRes, error) {
	page, pageSize := normalizePage(pageReq, pageSizeReq)
	model := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuStatus, uint(v1.SpuStatus_SPU_STATUS_ON_SHELF)).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt)
	if categoryID > 0 {
		model = model.Where(dao.CatalogSpu.Columns().CategoryId, categoryID)
	}
	if kw := strings.TrimSpace(keyword); kw != "" {
		model = model.WhereLike(dao.CatalogSpu.Columns().Title, "%"+kw+"%")
	}
	total, err := model.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count products failed")
	}
	var rows []*entity.CatalogSpu
	err = model.Page(page, pageSize).Order(buyerSortExpr(sortBy)).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "list products failed")
	}
	items := make([]*v1.BuyerProductCard, 0, len(rows))
	for _, row := range rows {
		mainImages := parseUint64Slice(row.MainImageAssetIdsJson)
		var cover uint64
		if len(mainImages) > 0 {
			cover = mainImages[0]
		}
		items = append(items, &v1.BuyerProductCard{
			SpuNo:             row.SpuNo,
			Title:             row.Title,
			CategoryId:        row.CategoryId,
			BrandNo:           row.BrandNo,
			CoverImageAssetId: cover,
			MinSalePrice:      row.MinSalePrice,
			MaxSalePrice:      row.MaxSalePrice,
			SpuStockStatus:    v1.StockStatus(row.SpuStockStatus),
			SpuStatus:         v1.SpuStatus(row.SpuStatus),
			SoldCount:         int64(row.SoldCount),
		})
	}
	return &v1.ListProductsRes{
		Items:    items,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Total:    int64(total),
	}, nil
}

// getProductAggregate 组装 SPU + SKU + 属性的聚合返回结构。
func (s *sCatalog) getProductAggregate(ctx context.Context, spuNo string, buyer bool) (*v1.ProductAggregate, error) {
	spu, err := s.getSpuEntity(ctx, spuNo)
	if err != nil {
		return nil, err
	}
	spuAttrs, err := s.querySpuAttrs(ctx, spuNo)
	if err != nil {
		return nil, err
	}

	skuModel := dao.CatalogSku.Ctx(ctx).
		Where(dao.CatalogSku.Columns().SpuNo, spuNo).
		WhereNull(dao.CatalogSku.Columns().DeletedAt)
	if buyer {
		skuModel = skuModel.Where(dao.CatalogSku.Columns().SkuStatus, uint(v1.SkuStatus_SKU_STATUS_ENABLED))
	}
	var skus []*entity.CatalogSku
	err = skuModel.OrderAsc(dao.CatalogSku.Columns().SortOrder).Scan(&skus)
	if err != nil {
		return nil, gerror.Wrap(err, "query product skus failed")
	}

	saleAttrsMap, err := s.querySkuSaleAttrs(ctx, skuNoList(skus))
	if err != nil {
		return nil, err
	}
	outSkus := make([]*v1.ProductSku, 0, len(skus))
	for _, row := range skus {
		outSkus = append(outSkus, toProtoSku(row, saleAttrsMap[row.SkuNo]))
	}
	return &v1.ProductAggregate{
		Spu:  toProtoSpu(spu, spuAttrs),
		Skus: outSkus,
	}, nil
}

// getSpuEntity 查询单个 SPU 实体（非事务版本）。
func (s *sCatalog) getSpuEntity(ctx context.Context, spuNo string) (*entity.CatalogSpu, error) {
	var spu entity.CatalogSpu
	err := dao.CatalogSpu.Ctx(ctx).
		Where(dao.CatalogSpu.Columns().SpuNo, spuNo).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Scan(&spu)
	if err != nil {
		return nil, gerror.Wrap(err, "query spu failed")
	}
	if spu.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "spu not found")
	}
	return &spu, nil
}

// getSpuEntityTx 在事务内查询单个 SPU 实体。
func (s *sCatalog) getSpuEntityTx(ctx context.Context, tx gdb.TX, spuNo string) (*entity.CatalogSpu, error) {
	var spu entity.CatalogSpu
	err := tx.Model(dao.CatalogSpu.Table()).
		Where(dao.CatalogSpu.Columns().SpuNo, spuNo).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Scan(&spu)
	if err != nil {
		return nil, gerror.Wrap(err, "query spu failed")
	}
	if spu.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "spu not found")
	}
	return &spu, nil
}

// querySpuAttrs 查询 SPU 属性并转换为 API 结构。
func (s *sCatalog) querySpuAttrs(ctx context.Context, spuNo string) ([]*v1.AttributeValue, error) {
	var rows []*entity.CatalogSpuAttrValue
	err := dao.CatalogSpuAttrValue.Ctx(ctx).
		Where(dao.CatalogSpuAttrValue.Columns().SpuNo, spuNo).
		OrderAsc(dao.CatalogSpuAttrValue.Columns().SortOrder).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "query spu attrs failed")
	}
	out := make([]*v1.AttributeValue, 0, len(rows))
	for _, row := range rows {
		out = append(out, &v1.AttributeValue{
			AttrCode: row.AttrCode,
			AttrName: row.AttrName,
			Scope:    v1.AttributeScope(row.AttrScope),
			Value:    row.AttrValue,
		})
	}
	return out, nil
}

// querySkuSaleAttrs 批量查询 SKU 销售属性并按 sku_no 分组。
func (s *sCatalog) querySkuSaleAttrs(ctx context.Context, skuNos []string) (map[string][]*v1.SkuSaleAttr, error) {
	out := make(map[string][]*v1.SkuSaleAttr)
	if len(skuNos) == 0 {
		return out, nil
	}
	var rows []*entity.CatalogSkuSaleAttrValue
	err := dao.CatalogSkuSaleAttrValue.Ctx(ctx).
		WhereIn(dao.CatalogSkuSaleAttrValue.Columns().SkuNo, skuNos).
		OrderAsc(dao.CatalogSkuSaleAttrValue.Columns().SortOrder).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "query sku sale attrs failed")
	}
	for _, row := range rows {
		out[row.SkuNo] = append(out[row.SkuNo], &v1.SkuSaleAttr{
			AttrCode: row.AttrCode,
			AttrName: row.AttrName,
			Value:    row.AttrValue,
		})
	}
	return out, nil
}

// replaceSpuAttrsTx 在事务内重建 SPU 属性明细。
func (s *sCatalog) replaceSpuAttrsTx(ctx context.Context, tx gdb.TX, spuNo string, attrs []*v1.AttributeValue) error {
	_, err := tx.Model(dao.CatalogSpuAttrValue.Table()).
		Where(dao.CatalogSpuAttrValue.Columns().SpuNo, spuNo).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "delete spu attrs failed")
	}
	if len(attrs) == 0 {
		return nil
	}
	items := make([]do.CatalogSpuAttrValue, 0, len(attrs))
	for idx, attr := range attrs {
		items = append(items, do.CatalogSpuAttrValue{
			SpuNo:     spuNo,
			AttrCode:  attr.GetAttrCode(),
			AttrName:  attr.GetAttrName(),
			AttrScope: uint(attr.GetScope()),
			AttrValue: attr.GetValue(),
			SortOrder: idx + 1,
		})
	}
	_, err = tx.Model(dao.CatalogSpuAttrValue.Table()).Data(items).Insert()
	if err != nil {
		return gerror.Wrap(err, "insert spu attrs failed")
	}
	return nil
}

// replaceSkuSaleAttrsTx 在事务内重建 SKU 销售属性明细。
func (s *sCatalog) replaceSkuSaleAttrsTx(ctx context.Context, tx gdb.TX, spuNo, skuNo string, attrs []*v1.SkuSaleAttr) error {
	_, err := tx.Model(dao.CatalogSkuSaleAttrValue.Table()).
		Where(dao.CatalogSkuSaleAttrValue.Columns().SkuNo, skuNo).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "delete sku sale attrs failed")
	}
	if len(attrs) == 0 {
		return nil
	}
	items := make([]do.CatalogSkuSaleAttrValue, 0, len(attrs))
	for idx, attr := range attrs {
		items = append(items, do.CatalogSkuSaleAttrValue{
			SkuNo:     skuNo,
			SpuNo:     spuNo,
			AttrCode:  attr.GetAttrCode(),
			AttrName:  attr.GetAttrName(),
			AttrValue: attr.GetValue(),
			SortOrder: idx + 1,
		})
	}
	_, err = tx.Model(dao.CatalogSkuSaleAttrValue.Table()).Data(items).Insert()
	if err != nil {
		return gerror.Wrap(err, "insert sku sale attrs failed")
	}
	return nil
}

// finishLatestReviewTaskTx 结束最新待处理审核任务并写入审核结论。
func (s *sCatalog) finishLatestReviewTaskTx(ctx context.Context, tx gdb.TX, spuNo string, reviewStatus uint, rejectCode, rejectComment, reviewComment string) error {
	var task entity.CatalogReviewTask
	err := tx.Model(dao.CatalogReviewTask.Table()).
		Where(dao.CatalogReviewTask.Columns().SpuNo, spuNo).
		Where(dao.CatalogReviewTask.Columns().ReviewStatus, uint(v1.ReviewStatus_REVIEW_STATUS_PENDING)).
		OrderDesc(dao.CatalogReviewTask.Columns().Id).
		Scan(&task)
	if err != nil {
		return gerror.Wrap(err, "query latest review task failed")
	}
	if task.Id == 0 {
		return nil
	}
	_, err = tx.Model(dao.CatalogReviewTask.Table()).
		Where(dao.CatalogReviewTask.Columns().TaskNo, task.TaskNo).
		Data(do.CatalogReviewTask{
			ReviewStatus:     reviewStatus,
			RejectReasonCode: rejectCode,
			RejectComment:    rejectComment,
			ReviewComment:    reviewComment,
			ReviewedAt:       gtime.Now(),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "update latest review task failed")
	}
	return nil
}

// recomputeSpuAggregation 事务封装的 SPU 聚合重算入口。
func (s *sCatalog) recomputeSpuAggregation(ctx context.Context, spuNo string) (bool, error) {
	var updated bool
	err := dao.CatalogSpu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var err error
		updated, err = s.recomputeSpuAggregationTx(ctx, tx, spuNo)
		return err
	})
	if err != nil {
		return false, err
	}
	return updated, nil
}

// recomputeSpuAggregationTx 在事务内按 SKU 重算 SPU 聚合价格和有货状态。
func (s *sCatalog) recomputeSpuAggregationTx(ctx context.Context, tx gdb.TX, spuNo string) (bool, error) {
	var skus []*entity.CatalogSku
	err := tx.Model(dao.CatalogSku.Table()).
		Where(dao.CatalogSku.Columns().SpuNo, spuNo).
		WhereNull(dao.CatalogSku.Columns().DeletedAt).
		Where(dao.CatalogSku.Columns().SkuStatus+" <> ?", uint(v1.SkuStatus_SKU_STATUS_DELETED)).
		Scan(&skus)
	if err != nil {
		return false, gerror.Wrap(err, "query sku for aggregation failed")
	}

	var (
		minSale    uint64
		maxSale    uint64
		minMarket  uint64
		maxMarket  uint64
		stockState = uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
	)
	for i, row := range skus {
		if i == 0 {
			minSale = row.SalePrice
			maxSale = row.SalePrice
			minMarket = row.MarketPrice
			maxMarket = row.MarketPrice
		} else {
			if row.SalePrice < minSale {
				minSale = row.SalePrice
			}
			if row.SalePrice > maxSale {
				maxSale = row.SalePrice
			}
			if row.MarketPrice < minMarket {
				minMarket = row.MarketPrice
			}
			if row.MarketPrice > maxMarket {
				maxMarket = row.MarketPrice
			}
		}
		if row.SkuStatus == uint(v1.SkuStatus_SKU_STATUS_ENABLED) && row.StockStatus == uint(v1.StockStatus_STOCK_STATUS_IN_STOCK) {
			stockState = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
		}
	}
	result, err := tx.Model(dao.CatalogSpu.Table()).
		Where(dao.CatalogSpu.Columns().SpuNo, spuNo).
		WhereNull(dao.CatalogSpu.Columns().DeletedAt).
		Data(do.CatalogSpu{
			MinSalePrice:   minSale,
			MaxSalePrice:   maxSale,
			MinMarketPrice: minMarket,
			MaxMarketPrice: maxMarket,
			SpuStockStatus: stockState,
		}).
		Update()
	if err != nil {
		return false, gerror.Wrap(err, "update spu aggregation failed")
	}
	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

// toProtoSpu 将 SPU 实体转换为 Proto 返回结构。
func toProtoSpu(row *entity.CatalogSpu, attrs []*v1.AttributeValue) *v1.ProductSpu {
	return &v1.ProductSpu{
		SpuNo:               row.SpuNo,
		ShopNo:              row.ShopNo,
		Title:               row.Title,
		SubTitle:            row.SubTitle,
		CategoryId:          row.CategoryId,
		BrandNo:             row.BrandNo,
		MainImageAssetIds:   parseUint64Slice(row.MainImageAssetIdsJson),
		DetailImageAssetIds: parseUint64Slice(row.DetailImageAssetIdsJson),
		SpuAttrs:            attrs,
		SpuStatus:           v1.SpuStatus(row.SpuStatus),
		SpuStockStatus:      v1.StockStatus(row.SpuStockStatus),
		MinSalePrice:        row.MinSalePrice,
		MaxSalePrice:        row.MaxSalePrice,
		MinMarketPrice:      row.MinMarketPrice,
		MaxMarketPrice:      row.MaxMarketPrice,
		PublishTime:         toProtoTs(row.PublishTime),
		Version:             uint32(row.Version),
		CreatedAt:           toProtoTs(row.CreatedAt),
		UpdatedAt:           toProtoTs(row.UpdatedAt),
	}
}

// toProtoSku 将 SKU 实体转换为 Proto 返回结构。
func toProtoSku(row *entity.CatalogSku, saleAttrs []*v1.SkuSaleAttr) *v1.ProductSku {
	return &v1.ProductSku{
		SkuNo:           row.SkuNo,
		SpuNo:           row.SpuNo,
		SkuName:         row.SkuName,
		SkuImageAssetId: row.SkuImageAssetId,
		SkuStatus:       v1.SkuStatus(row.SkuStatus),
		StockStatus:     v1.StockStatus(row.StockStatus),
		StockVersion:    row.StockVersion,
		SalePrice:       row.SalePrice,
		MarketPrice:     row.MarketPrice,
		SaleAttrs:       saleAttrs,
		SortOrder:       int32(row.SortOrder),
		Version:         uint32(row.Version),
		CreatedAt:       toProtoTs(row.CreatedAt),
		UpdatedAt:       toProtoTs(row.UpdatedAt),
	}
}

// spuReviewInfo 从 SPU 实体提取审核信息。
func spuReviewInfo(spu *entity.CatalogSpu) *v1.ReviewInfo {
	if spu == nil {
		return &v1.ReviewInfo{}
	}
	return &v1.ReviewInfo{
		ReviewStatus:     v1.ReviewStatus(spu.ReviewStatus),
		RejectReasonCode: spu.RejectReasonCode,
		RejectComment:    spu.RejectComment,
		ReviewerId:       spu.ReviewerId,
		ReviewedAt:       toProtoTs(spu.ReviewedAt),
	}
}

// parseUint64Slice 解析 JSON 数组字符串为 uint64 列表。
func parseUint64Slice(v string) []uint64 {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	var out []uint64
	_ = json.Unmarshal([]byte(v), &out)
	return out
}

// marshalJSON 将对象序列化为 JSON 字符串。
func marshalJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", gerror.Wrap(err, "marshal json failed")
	}
	return string(b), nil
}

// normalizeSaleAttrs 规范化销售属性并生成稳定哈希，支持唯一性判断。
func normalizeSaleAttrs(attrs []*v1.SkuSaleAttr) (string, string, error) {
	normalized := make([]*v1.SkuSaleAttr, 0, len(attrs))
	for _, attr := range attrs {
		normalized = append(normalized, &v1.SkuSaleAttr{
			AttrCode: strings.TrimSpace(attr.GetAttrCode()),
			AttrName: strings.TrimSpace(attr.GetAttrName()),
			Value:    strings.TrimSpace(attr.GetValue()),
		})
	}
	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].GetAttrCode() == normalized[j].GetAttrCode() {
			return normalized[i].GetValue() < normalized[j].GetValue()
		}
		return normalized[i].GetAttrCode() < normalized[j].GetAttrCode()
	})
	b, err := json.Marshal(normalized)
	if err != nil {
		return "", "", gerror.Wrap(err, "marshal sale attrs failed")
	}
	sum := sha256.Sum256(b)
	return string(b), hex.EncodeToString(sum[:]), nil
}

// skuNoList 从 SKU 实体列表提取 sku_no 数组。
func skuNoList(skus []*entity.CatalogSku) []string {
	out := make([]string, 0, len(skus))
	for _, row := range skus {
		out = append(out, row.SkuNo)
	}
	return out
}

// uniqueSpuNosFromSkus 从 SKU 列表提取去重后的 spu_no。
func uniqueSpuNosFromSkus(skus []*entity.CatalogSku) []string {
	m := make(map[string]struct{})
	for _, row := range skus {
		if row.SpuNo == "" {
			continue
		}
		m[row.SpuNo] = struct{}{}
	}
	out := make([]string, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	return out
}

// normalizePage 统一分页参数并限制最大页大小。
func normalizePage(pageReq, pageSizeReq int32) (int, int) {
	page := int(pageReq)
	pageSize := int(pageSizeReq)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// normalizePatchPath 标准化 FieldMask 路径。
func normalizePatchPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "patch.")
	path = strings.ReplaceAll(path, ".", "_")
	return path
}

// buyerSortExpr 根据排序枚举生成 SQL 排序表达式。
func buyerSortExpr(sortBy v1.SortBy) string {
	switch sortBy {
	case v1.SortBy_SORT_BY_PRICE_ASC:
		return dao.CatalogSpu.Columns().MinSalePrice + " ASC, " + dao.CatalogSpu.Columns().Id + " DESC"
	case v1.SortBy_SORT_BY_PRICE_DESC:
		return dao.CatalogSpu.Columns().MinSalePrice + " DESC, " + dao.CatalogSpu.Columns().Id + " DESC"
	case v1.SortBy_SORT_BY_SALES_DESC:
		return dao.CatalogSpu.Columns().SoldCount + " DESC, " + dao.CatalogSpu.Columns().Id + " DESC"
	case v1.SortBy_SORT_BY_NEWEST:
		return dao.CatalogSpu.Columns().CreatedAt + " DESC, " + dao.CatalogSpu.Columns().Id + " DESC"
	default:
		return dao.CatalogSpu.Columns().Id + " DESC"
	}
}

// spuStatusToReviewStatus 将 SPU 状态筛选映射为审核状态筛选。
func spuStatusToReviewStatus(statuses []v1.SpuStatus) []uint {
	m := make(map[uint]struct{})
	for _, st := range statuses {
		switch st {
		case v1.SpuStatus_SPU_STATUS_REVIEWING:
			m[uint(v1.ReviewStatus_REVIEW_STATUS_PENDING)] = struct{}{}
		case v1.SpuStatus_SPU_STATUS_APPROVED:
			m[uint(v1.ReviewStatus_REVIEW_STATUS_APPROVED)] = struct{}{}
		case v1.SpuStatus_SPU_STATUS_REJECTED:
			m[uint(v1.ReviewStatus_REVIEW_STATUS_REJECTED)] = struct{}{}
		}
	}
	out := make([]uint, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	return out
}

// reviewStatusToSpuStatus 将审核状态映射回 SPU 状态。
func reviewStatusToSpuStatus(reviewStatus uint) v1.SpuStatus {
	switch v1.ReviewStatus(reviewStatus) {
	case v1.ReviewStatus_REVIEW_STATUS_PENDING:
		return v1.SpuStatus_SPU_STATUS_REVIEWING
	case v1.ReviewStatus_REVIEW_STATUS_APPROVED:
		return v1.SpuStatus_SPU_STATUS_APPROVED
	case v1.ReviewStatus_REVIEW_STATUS_REJECTED:
		return v1.SpuStatus_SPU_STATUS_REJECTED
	default:
		return v1.SpuStatus_SPU_STATUS_UNSPECIFIED
	}
}

// toProtoTs 将 gtime 时间转换为 protobuf Timestamp。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// protoTsToGTime 将 protobuf Timestamp 转换为 gtime 时间。
func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

// generateBizNo 生成业务号（前缀+时间戳+随机尾号）。
func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
