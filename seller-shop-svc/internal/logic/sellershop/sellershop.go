package sellershop

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/dao"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100

	shopNameReserveLimitPerOwner = 2
	shopNameReserveHours         = 48

	registryStatusReserved = 1
	registryStatusBound    = 2
	registryStatusReleased = 3
	outboxStatusNew        = 1
)

// SellerShop 是 seller-shop 领域服务实现，负责入驻申请、审核流与店铺状态流转。
type sSellerShop struct{}

// outboxEventPayload 是 seller-shop 发到 MQ 的统一事件包裹结构。
// Data 字段承载不同事件类型的业务载荷（如赋权/撤权请求）。
type outboxEventPayload struct {
	EventID       string      `json:"event_id"`
	RequestID     string      `json:"request_id"`
	OccurredAt    string      `json:"occurred_at"`
	Producer      string      `json:"producer"`
	AggregateType string      `json:"aggregate_type"`
	AggregateNo   string      `json:"aggregate_no"`
	Data          interface{} `json:"data"`
}

// New 创建 seller-shop 逻辑实例。
func New() *sSellerShop {
	return &sSellerShop{}
}

func init() {
	service.RegisterSellerShop(New())
}

// CreateApplicationDraft 创建入驻草稿并预占店铺名。
func (s *sSellerShop) CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (*v1.CreateApplicationDraftRes, error) {
	if req.GetEntity() == nil || req.GetShop() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "entity and shop are required")
	}
	if strings.TrimSpace(req.GetShop().GetShopName()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop.shop_name is required")
	}

	ownerUserID, ok := userIDFromMetadata(ctx, "x-user-id")
	if !ok {
		ownerUserID = req.GetEntity().GetOwnerUserId()
	}
	if ownerUserID == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "owner_user_id is required")
	}

	// applicationNo 是对外业务号，不暴露数据库自增 ID。
	applicationNo := generateBizNo("APP")
	entityProfile := proto.Clone(req.GetEntity()).(*v1.EntityProfile)
	shopProfile := proto.Clone(req.GetShop()).(*v1.ShopProfile)
	entityProfile.OwnerUserId = ownerUserID
	if strings.TrimSpace(entityProfile.GetEntityNo()) == "" {
		entityProfile.EntityNo = generateBizNo("ENT")
	}
	if strings.TrimSpace(shopProfile.GetShopNo()) == "" {
		shopProfile.ShopNo = generateBizNo("SHP")
	}

	// shopNameNorm 用于唯一约束比较，避免大小写/空格差异导致“伪不同名”。
	shopNameNorm := normalizeShopName(shopProfile.GetShopName())
	if shopNameNorm == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop.shop_name is invalid")
	}

	entityDraftJSON, err := marshalProtoJSON(entityProfile)
	if err != nil {
		return nil, err
	}
	shopDraftJSON, err := marshalProtoJSON(shopProfile)
	if err != nil {
		return nil, err
	}

	err = dao.SellerApplication.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if reserveErr := reserveShopNameTx(ctx, tx, ownerUserID, applicationNo, shopProfile.GetShopName(), shopNameNorm); reserveErr != nil {
			return reserveErr
		}
		_, insertErr := tx.Model(dao.SellerApplication.Table()).
			Data(do.SellerApplication{
				ApplicationNo:   applicationNo,
				Version:         1,
				OwnerUserId:     ownerUserID,
				Status:          uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT),
				EntityNo:        entityProfile.GetEntityNo(),
				ShopNo:          shopProfile.GetShopNo(),
				EntityName:      entityProfile.GetEntityName(),
				ShopName:        shopProfile.GetShopName(),
				ShopNameNorm:    shopNameNorm,
				EntityDraftJson: entityDraftJSON,
				ShopDraftJson:   shopDraftJSON,
			}).
			Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "insert draft failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	app, err := queryApplicationByNo(ctx, applicationNo)
	if err != nil {
		return nil, err
	}
	return &v1.CreateApplicationDraftRes{Application: toProtoApplication(app)}, nil
}

// UpdateApplicationDraft 更新草稿内容并刷新店铺名预占记录。
func (s *sSellerShop) UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (*v1.UpdateApplicationDraftRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	if req.GetExpectedVersion() <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}

	app, err := queryApplicationByNo(ctx, req.GetApplicationNo())
	if err != nil {
		return nil, err
	}
	if app.Status != uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "only draft can be updated")
	}
	// expected_version + 当前 version 做乐观锁，防止并发覆盖。
	if app.Version != uint(req.GetExpectedVersion()) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}

	entityProfile, _ := unmarshalEntityProfile(app.EntityDraftJson)
	shopProfile, _ := unmarshalShopProfile(app.ShopDraftJson)
	if entityProfile == nil {
		entityProfile = &v1.EntityProfile{}
	}
	if shopProfile == nil {
		shopProfile = &v1.ShopProfile{}
	}

	if req.GetEntity() != nil {
		entityProfile = proto.Clone(req.GetEntity()).(*v1.EntityProfile)
	}
	if req.GetShop() != nil {
		shopProfile = proto.Clone(req.GetShop()).(*v1.ShopProfile)
	}

	entityProfile.OwnerUserId = app.OwnerUserId
	if strings.TrimSpace(entityProfile.GetEntityNo()) == "" {
		entityProfile.EntityNo = app.EntityNo
	}
	if strings.TrimSpace(shopProfile.GetShopNo()) == "" {
		shopProfile.ShopNo = app.ShopNo
	}

	shopNameNorm := normalizeShopName(shopProfile.GetShopName())
	if shopNameNorm == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop.shop_name is invalid")
	}

	entityDraftJSON, err := marshalProtoJSON(entityProfile)
	if err != nil {
		return nil, err
	}
	shopDraftJSON, err := marshalProtoJSON(shopProfile)
	if err != nil {
		return nil, err
	}

	err = dao.SellerApplication.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.SellerApplication.Columns()
		result, updateErr := tx.Model(dao.SellerApplication.Table()).
			Where(cols.ApplicationNo, app.ApplicationNo).
			Where(cols.Version, app.Version).
			Where(cols.Status, uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT)).
			Data(do.SellerApplication{
				Version:         app.Version + 1,
				EntityNo:        entityProfile.GetEntityNo(),
				ShopNo:          shopProfile.GetShopNo(),
				EntityName:      entityProfile.GetEntityName(),
				ShopName:        shopProfile.GetShopName(),
				ShopNameNorm:    shopNameNorm,
				EntityDraftJson: entityDraftJSON,
				ShopDraftJson:   shopDraftJSON,
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "update draft failed")
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return gerror.Wrap(rowsErr, "get affected rows failed")
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "draft update conflict")
		}
		if reserveErr := updateReservedShopNameTx(ctx, tx, app.OwnerUserId, app.ApplicationNo, shopProfile.GetShopName(), shopNameNorm); reserveErr != nil {
			return reserveErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := queryApplicationByNo(ctx, app.ApplicationNo)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateApplicationDraftRes{Application: toProtoApplication(updated)}, nil
}

// SubmitApplication 把草稿状态切换为已提交，并固化审核快照字段。
// 这里使用 version 条件更新，防止多端重复提交覆盖。
func (s *sSellerShop) SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (*v1.SubmitApplicationRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	if req.GetExpectedVersion() <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}

	app, err := queryApplicationByNo(ctx, req.GetApplicationNo())
	if err != nil {
		return nil, err
	}
	if app.Status != uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "only draft can be submitted")
	}
	if app.Version != uint(req.GetExpectedVersion()) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}

	cols := dao.SellerApplication.Columns()
	result, err := dao.SellerApplication.Ctx(ctx).
		Where(cols.ApplicationNo, app.ApplicationNo).
		Where(cols.Version, app.Version).
		Where(cols.Status, uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT)).
		Data(do.SellerApplication{
			Version:             app.Version + 1,
			Status:              uint(v1.ApplicationStatus_APPLICATION_STATUS_SUBMITTED),
			EntitySubmittedJson: app.EntityDraftJson,
			ShopSubmittedJson:   app.ShopDraftJson,
			SubmittedAt:         gtime.Now(),
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "submit application failed")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, gerror.Wrap(err, "get affected rows failed")
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "application submit conflict")
	}

	updated, err := queryApplicationByNo(ctx, app.ApplicationNo)
	if err != nil {
		return nil, err
	}
	return &v1.SubmitApplicationRes{Application: toProtoApplication(updated)}, nil
}

// ResubmitApplication 针对被驳回申请创建“下一张”申请单，保持链式追踪。
// 老申请通过 next_application_no 指向新单，避免并发重复重提。
func (s *sSellerShop) ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (*v1.ResubmitApplicationRes, error) {
	if strings.TrimSpace(req.GetRejectedApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "rejected_application_no is required")
	}
	if req.GetExpectedVersion() <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}

	oldApp, err := queryApplicationByNo(ctx, req.GetRejectedApplicationNo())
	if err != nil {
		return nil, err
	}
	if oldApp.Status != uint(v1.ApplicationStatus_APPLICATION_STATUS_REJECTED) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "only rejected application can be resubmitted")
	}
	if oldApp.Version != uint(req.GetExpectedVersion()) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}
	if strings.TrimSpace(oldApp.NextApplicationNo) != "" {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "resubmit already created")
	}

	entityProfile, _ := unmarshalEntityProfile(oldApp.EntityDraftJson)
	shopProfile, _ := unmarshalShopProfile(oldApp.ShopDraftJson)
	if req.GetEntity() != nil {
		entityProfile = proto.Clone(req.GetEntity()).(*v1.EntityProfile)
	}
	if req.GetShop() != nil {
		shopProfile = proto.Clone(req.GetShop()).(*v1.ShopProfile)
	}
	if entityProfile == nil || shopProfile == nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "application draft data is invalid")
	}

	entityProfile.OwnerUserId = oldApp.OwnerUserId
	if strings.TrimSpace(entityProfile.GetEntityNo()) == "" {
		entityProfile.EntityNo = oldApp.EntityNo
	}
	if strings.TrimSpace(shopProfile.GetShopNo()) == "" {
		shopProfile.ShopNo = generateBizNo("SHP")
	}

	shopNameNorm := normalizeShopName(shopProfile.GetShopName())
	if shopNameNorm == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop.shop_name is invalid")
	}

	newApplicationNo := generateBizNo("APP")
	entityDraftJSON, err := marshalProtoJSON(entityProfile)
	if err != nil {
		return nil, err
	}
	shopDraftJSON, err := marshalProtoJSON(shopProfile)
	if err != nil {
		return nil, err
	}

	newStatus := uint(v1.ApplicationStatus_APPLICATION_STATUS_DRAFT)
	entitySubmittedJSON := ""
	shopSubmittedJSON := ""
	var submittedAt *gtime.Time
	if req.GetSubmitImmediately() {
		newStatus = uint(v1.ApplicationStatus_APPLICATION_STATUS_SUBMITTED)
		entitySubmittedJSON = entityDraftJSON
		shopSubmittedJSON = shopDraftJSON
		submittedAt = gtime.Now()
	}

	err = dao.SellerApplication.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.SellerApplication.Columns()
		result, updateErr := tx.Model(dao.SellerApplication.Table()).
			Where(cols.ApplicationNo, oldApp.ApplicationNo).
			Where(cols.Version, oldApp.Version).
			Where(cols.Status, uint(v1.ApplicationStatus_APPLICATION_STATUS_REJECTED)).
			Where(cols.NextApplicationNo, "").
			Data(do.SellerApplication{
				Version:           oldApp.Version + 1,
				NextApplicationNo: newApplicationNo,
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "update previous application failed")
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return gerror.Wrap(rowsErr, "get affected rows failed")
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "resubmit conflict")
		}

		if reserveErr := reserveShopNameTx(ctx, tx, oldApp.OwnerUserId, newApplicationNo, shopProfile.GetShopName(), shopNameNorm); reserveErr != nil {
			return reserveErr
		}

		_, insertErr := tx.Model(dao.SellerApplication.Table()).
			Data(do.SellerApplication{
				ApplicationNo:         newApplicationNo,
				Version:               1,
				PreviousApplicationNo: oldApp.ApplicationNo,
				OwnerUserId:           oldApp.OwnerUserId,
				Status:                newStatus,
				EntityNo:              entityProfile.GetEntityNo(),
				ShopNo:                shopProfile.GetShopNo(),
				EntityName:            entityProfile.GetEntityName(),
				ShopName:              shopProfile.GetShopName(),
				ShopNameNorm:          shopNameNorm,
				EntityDraftJson:       entityDraftJSON,
				ShopDraftJson:         shopDraftJSON,
				EntitySubmittedJson:   entitySubmittedJSON,
				ShopSubmittedJson:     shopSubmittedJSON,
				SubmittedAt:           submittedAt,
			}).
			Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "insert resubmitted application failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	created, err := queryApplicationByNo(ctx, newApplicationNo)
	if err != nil {
		return nil, err
	}
	return &v1.ResubmitApplicationRes{Application: toProtoApplication(created)}, nil
}

// GetMyApplication 查询当前用户的申请单详情。
func (s *sSellerShop) GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (*v1.GetMyApplicationRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	cols := dao.SellerApplication.Columns()
	model := dao.SellerApplication.Ctx(ctx).Where(cols.ApplicationNo, req.GetApplicationNo())
	if ownerUserID, ok := userIDFromMetadata(ctx, "x-user-id"); ok {
		model = model.Where(cols.OwnerUserId, ownerUserID)
	}

	var app entity.SellerApplication
	if err := model.Scan(&app); err != nil {
		return nil, gerror.Wrap(err, "query application failed")
	}
	if app.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "application not found")
	}
	return &v1.GetMyApplicationRes{Application: toProtoApplication(&app)}, nil
}

// ListMyApplications 分页列出当前用户申请单。
func (s *sSellerShop) ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (*v1.ListMyApplicationsRes, error) {
	ownerUserID, ok := userIDFromMetadata(ctx, "x-user-id")
	if !ok || ownerUserID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "x-user-id metadata is required")
	}

	page, pageSize := normalizePaging(req.GetPage(), req.GetPageSize())
	cols := dao.SellerApplication.Columns()
	model := dao.SellerApplication.Ctx(ctx).Where(cols.OwnerUserId, ownerUserID)
	if statuses := enumApplicationStatuses(req.GetStatuses()); len(statuses) > 0 {
		model = model.WhereIn(cols.Status, statuses)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count applications failed")
	}

	var apps []*entity.SellerApplication
	if err = model.Page(page, pageSize).OrderDesc(cols.Id).Scan(&apps); err != nil {
		return nil, gerror.Wrap(err, "list applications failed")
	}

	items := make([]*v1.SellerApplication, 0, len(apps))
	for _, app := range apps {
		items = append(items, toProtoApplication(app))
	}
	return &v1.ListMyApplicationsRes{
		Applications: items,
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Total:        int64(total),
	}, nil
}

// ListApplications 管理端分页列出申请单。
func (s *sSellerShop) ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (*v1.ListApplicationsRes, error) {
	page, pageSize := normalizePaging(req.GetPage(), req.GetPageSize())
	cols := dao.SellerApplication.Columns()
	model := dao.SellerApplication.Ctx(ctx)
	if statuses := enumApplicationStatuses(req.GetStatuses()); len(statuses) > 0 {
		model = model.WhereIn(cols.Status, statuses)
	}

	keyword := strings.TrimSpace(req.GetKeyword())
	if keyword != "" {
		like := "%" + keyword + "%"
		model = model.Wheref("( %s LIKE ? OR %s LIKE ? )", cols.ShopName, cols.EntityName, like, like)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count applications failed")
	}

	var apps []*entity.SellerApplication
	if err = model.Page(page, pageSize).OrderDesc(cols.Id).Scan(&apps); err != nil {
		return nil, gerror.Wrap(err, "list applications failed")
	}

	items := make([]*v1.SellerApplication, 0, len(apps))
	for _, app := range apps {
		items = append(items, toProtoApplication(app))
	}
	return &v1.ListApplicationsRes{
		Applications: items,
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Total:        int64(total),
	}, nil
}

// GetApplicationDetail 管理端查询申请单详情。
func (s *sSellerShop) GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (*v1.GetApplicationDetailRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	app, err := queryApplicationByNo(ctx, req.GetApplicationNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetApplicationDetailRes{Application: toProtoApplication(app)}, nil
}

// ApproveApplication 管理审核通过：回写主体主数据、建店、更新申请状态并写 outbox 赋权事件。
// 店铺初始状态为 PROVISIONING，待 IAM 赋权确认后再转 ACTIVE。
func (s *sSellerShop) ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (*v1.ApproveApplicationRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	if req.GetExpectedVersion() <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}

	app, err := queryApplicationByNo(ctx, req.GetApplicationNo())
	if err != nil {
		return nil, err
	}
	if app.Version != uint(req.GetExpectedVersion()) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}
	if app.Status != uint(v1.ApplicationStatus_APPLICATION_STATUS_SUBMITTED) &&
		app.Status != uint(v1.ApplicationStatus_APPLICATION_STATUS_REVIEWING) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "application status is not approvable")
	}

	entityProfile, _ := unmarshalEntityProfile(nonEmptyString(app.EntitySubmittedJson, app.EntityDraftJson))
	shopProfile, _ := unmarshalShopProfile(nonEmptyString(app.ShopSubmittedJson, app.ShopDraftJson))
	if entityProfile == nil || shopProfile == nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "submitted snapshot is invalid")
	}

	// reviewerID 用于审计“谁审批了申请单”；requestID 用于跨服务链路追踪。
	reviewerID, _ := userIDFromMetadata(ctx, "x-operator-user-id", "x-user-id")
	requestID := metadataValue(ctx, "x-request-id")

	err = dao.SellerApplication.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		entityID, upsertErr := upsertEntityTx(ctx, tx, app.OwnerUserId, entityProfile)
		if upsertErr != nil {
			return upsertErr
		}
		if syncErr := replaceEntityDocsTx(ctx, tx, entityID, entityProfile.GetEntityNo(), entityProfile.GetLegalSubject().GetDocs()); syncErr != nil {
			return syncErr
		}

		shopNameNorm := normalizeShopName(shopProfile.GetShopName())
		if shopNameNorm == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, "shop.shop_name is invalid")
		}

		_, insertErr := tx.Model(dao.SellerShop.Table()).
			Data(do.SellerShop{
				ShopNo:              shopProfile.GetShopNo(),
				EntityId:            entityID,
				EntityNo:            entityProfile.GetEntityNo(),
				ApplicationNo:       app.ApplicationNo,
				OwnerUserId:         app.OwnerUserId,
				ShopName:            shopProfile.GetShopName(),
				ShopNameNorm:        shopNameNorm,
				ShopDisplayName:     shopProfile.GetShopDisplayName(),
				ShopTypeCode:        shopProfile.GetShopTypeCode(),
				MainCategoryIdsJson: mustMarshalJSON(shopProfile.GetMainCategoryIds()),
				LogoAssetId:         shopProfile.GetLogoAssetId(),
				BannerAssetId:       shopProfile.GetBannerAssetId(),
				ServicePhone:        shopProfile.GetServicePhone(),
				ServiceEmail:        shopProfile.GetServiceEmail(),
				Description:         shopProfile.GetDescription(),
				ExtJson:             mustMarshalJSON(shopProfile.GetExt()),
				Status:              uint(v1.ShopStatus_SHOP_STATUS_PROVISIONING),
				BuyerVisible:        0,
				Version:             1,
				ProvisionStartedAt:  gtime.Now(),
			}).
			Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "create shop failed")
		}

		cols := dao.SellerApplication.Columns()
		result, updateErr := tx.Model(dao.SellerApplication.Table()).
			Where(cols.ApplicationNo, app.ApplicationNo).
			Where(cols.Version, app.Version).
			WhereIn(cols.Status, []uint{
				uint(v1.ApplicationStatus_APPLICATION_STATUS_SUBMITTED),
				uint(v1.ApplicationStatus_APPLICATION_STATUS_REVIEWING),
			}).
			Data(do.SellerApplication{
				Version:       app.Version + 1,
				Status:        uint(v1.ApplicationStatus_APPLICATION_STATUS_APPROVED),
				ReviewedAt:    gtime.Now(),
				ReviewerId:    reviewerID,
				ReviewComment: req.GetReviewComment(),
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "update approved status failed")
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return gerror.Wrap(rowsErr, "get affected rows failed")
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "approve conflict")
		}

		_, regErr := tx.Model(dao.SellerShopNameRegistry.Table()).
			Where(dao.SellerShopNameRegistry.Columns().ApplicationNo, app.ApplicationNo).
			Data(do.SellerShopNameRegistry{
				Status:      registryStatusBound,
				BoundShopNo: shopProfile.GetShopNo(),
				ExpireAt:    gtime.NewFromStr("2099-12-31 23:59:59"),
			}).
			Update()
		if regErr != nil {
			return gerror.Wrap(regErr, "update shop name registry failed")
		}

		event := buildOutboxEvent(requestID, "SHOP", shopProfile.GetShopNo(), map[string]interface{}{
			"user_id":        app.OwnerUserId,
			"application_no": app.ApplicationNo,
			"entity_no":      entityProfile.GetEntityNo(),
			"shop_no":        shopProfile.GetShopNo(),
		})
		_, outboxErr := tx.Model(dao.SellerOutboxEvent.Table()).
			Data(do.SellerOutboxEvent{
				EventId:       event.EventID,
				EventType:     "SellerRoleAssignRequested",
				AggregateType: "SHOP",
				AggregateNo:   shopProfile.GetShopNo(),
				RequestId:     requestID,
				PayloadJson:   mustMarshalJSON(event),
				Status:        outboxStatusNew,
				AvailableAt:   gtime.Now(),
			}).
			Insert()
		if outboxErr != nil {
			return gerror.Wrap(outboxErr, "insert outbox event failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.ApproveApplicationRes{
		ApplicationNo: app.ApplicationNo,
		ShopNo:        shopProfile.GetShopNo(),
		ShopStatus:    v1.ShopStatus_SHOP_STATUS_PROVISIONING,
	}, nil
}

// RejectApplication 记录驳回原因并释放店名预占，允许后续重提。
func (s *sSellerShop) RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (*v1.RejectApplicationRes, error) {
	if strings.TrimSpace(req.GetApplicationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "application_no is required")
	}
	if req.GetExpectedVersion() <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}
	if strings.TrimSpace(req.GetRejectReasonCode()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reject_reason_code is required")
	}

	app, err := queryApplicationByNo(ctx, req.GetApplicationNo())
	if err != nil {
		return nil, err
	}
	if app.Version != uint(req.GetExpectedVersion()) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}

	rejectInfo := &v1.RejectInfo{
		RejectReasonCode: req.GetRejectReasonCode(),
		RejectComment:    req.GetRejectComment(),
		RejectedAt:       timestamppb.Now(),
	}
	rejectJSON, err := marshalProtoJSON(rejectInfo)
	if err != nil {
		return nil, err
	}

	reviewerID, _ := userIDFromMetadata(ctx, "x-operator-user-id", "x-user-id")
	err = dao.SellerApplication.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.SellerApplication.Columns()
		result, updateErr := tx.Model(dao.SellerApplication.Table()).
			Where(cols.ApplicationNo, app.ApplicationNo).
			Where(cols.Version, app.Version).
			Data(do.SellerApplication{
				Version:          app.Version + 1,
				Status:           uint(v1.ApplicationStatus_APPLICATION_STATUS_REJECTED),
				ReviewedAt:       gtime.Now(),
				ReviewerId:       reviewerID,
				ReviewComment:    req.GetRejectComment(),
				LatestRejectJson: rejectJSON,
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "reject application failed")
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return gerror.Wrap(rowsErr, "get affected rows failed")
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "reject conflict")
		}

		_, releaseErr := tx.Model(dao.SellerShopNameRegistry.Table()).
			Where(dao.SellerShopNameRegistry.Columns().ApplicationNo, app.ApplicationNo).
			Data(do.SellerShopNameRegistry{Status: registryStatusReleased, ExpireAt: gtime.Now()}).
			Update()
		if releaseErr != nil {
			return gerror.Wrap(releaseErr, "release reserved name failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RejectApplicationRes{
		ApplicationNo: app.ApplicationNo,
		Status:        v1.ApplicationStatus_APPLICATION_STATUS_REJECTED,
	}, nil
}

// FreezeShop 将店铺置为冻结流程，并发出撤权事件让 IAM 回收卖家权限。
func (s *sSellerShop) FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (*v1.FreezeShopRes, error) {
	return s.changeShopStatusWithRevoke(ctx, req.GetShopNo(), req.GetExpectedVersion(), req.GetReasonCode(), req.GetReason(), v1.ShopStatus_SHOP_STATUS_FREEZING)
}

// CloseShop 将店铺置为关店流程，并发出撤权事件让 IAM 回收卖家权限。
func (s *sSellerShop) CloseShop(ctx context.Context, req *v1.CloseShopReq) (*v1.CloseShopRes, error) {
	result, err := s.changeShopStatusWithRevoke(ctx, req.GetShopNo(), req.GetExpectedVersion(), req.GetReasonCode(), req.GetReason(), v1.ShopStatus_SHOP_STATUS_CLOSING)
	if err != nil {
		return nil, err
	}
	return &v1.CloseShopRes{
		ShopNo:     result.GetShopNo(),
		ShopStatus: result.GetShopStatus(),
	}, nil
}

// GetShopByNo 内部按 shop_no 查询店铺。
func (s *sSellerShop) GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (*v1.GetShopByNoRes, error) {
	if strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	row, err := queryShopByNo(ctx, req.GetShopNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetShopByNoRes{Shop: toProtoShop(row)}, nil
}

// BatchGetShopsByNo 内部批量查询店铺。
func (s *sSellerShop) BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (*v1.BatchGetShopsByNoRes, error) {
	if len(req.GetShopNos()) == 0 {
		return &v1.BatchGetShopsByNoRes{}, nil
	}
	var rows []*entity.SellerShop
	if err := dao.SellerShop.Ctx(ctx).WhereIn(dao.SellerShop.Columns().ShopNo, req.GetShopNos()).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "batch query shops failed")
	}
	out := make([]*v1.Shop, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProtoShop(row))
	}
	return &v1.BatchGetShopsByNoRes{Shops: out}, nil
}

// ListShopsByOwnerUserId 内部按 owner_user_id 列出店铺摘要。
func (s *sSellerShop) ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (*v1.ListShopsByOwnerUserIdRes, error) {
	if req.GetOwnerUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "owner_user_id is required")
	}
	cols := dao.SellerShop.Columns()
	var rows []*entity.SellerShop
	if err := dao.SellerShop.Ctx(ctx).Where(cols.OwnerUserId, req.GetOwnerUserId()).OrderDesc(cols.Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list shops failed")
	}
	out := make([]*v1.ShopSummary, 0, len(rows))
	for _, row := range rows {
		out = append(out, &v1.ShopSummary{
			ShopNo:          row.ShopNo,
			ShopName:        row.ShopName,
			ShopDisplayName: row.ShopDisplayName,
			Status:          v1.ShopStatus(row.Status),
			BuyerVisible:    row.BuyerVisible == 1,
		})
	}
	return &v1.ListShopsByOwnerUserIdRes{Shops: out}, nil
}

// IsUserShopOwner 校验某 user_id 是否是指定店铺拥有者，供内部鉴权调用。
func (s *sSellerShop) IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (*v1.IsUserShopOwnerRes, error) {
	if req.GetUserId() == 0 || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id and shop_no are required")
	}
	cols := dao.SellerShop.Columns()
	count, err := dao.SellerShop.Ctx(ctx).Where(cols.OwnerUserId, req.GetUserId()).Where(cols.ShopNo, req.GetShopNo()).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "check shop owner failed")
	}
	return &v1.IsUserShopOwnerRes{IsOwner: count > 0}, nil
}

// changeShopStatusWithRevoke 冻结/关店时更新状态并写入撤权 outbox 事件。
func (s *sSellerShop) changeShopStatusWithRevoke(
	ctx context.Context,
	shopNo string,
	expectedVersion int32,
	reasonCode string,
	reason string,
	toStatus v1.ShopStatus,
) (*v1.FreezeShopRes, error) {
	if strings.TrimSpace(shopNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	if expectedVersion <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version is required")
	}

	shop, err := queryShopByNo(ctx, shopNo)
	if err != nil {
		return nil, err
	}
	if shop.Version != uint(expectedVersion) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "version mismatch")
	}

	// requestID 透传到 outbox，便于后续 MQ/消费者日志对账。
	requestID := metadataValue(ctx, "x-request-id")
	err = dao.SellerShop.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.SellerShop.Columns()
		updateData := do.SellerShop{
			Version:      shop.Version + 1,
			Status:       uint(toStatus),
			BuyerVisible: 0,
		}
		if toStatus == v1.ShopStatus_SHOP_STATUS_FREEZING {
			updateData.FreezeReasonCode = reasonCode
			updateData.FreezeReason = reason
		}
		if toStatus == v1.ShopStatus_SHOP_STATUS_CLOSING {
			updateData.CloseReasonCode = reasonCode
			updateData.CloseReason = reason
		}

		result, updateErr := tx.Model(dao.SellerShop.Table()).
			Where(cols.ShopNo, shop.ShopNo).
			Where(cols.Version, shop.Version).
			Data(updateData).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "update shop status failed")
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return gerror.Wrap(rowsErr, "get affected rows failed")
		}
		if rows == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "shop status conflict")
		}

		event := buildOutboxEvent(requestID, "SHOP", shop.ShopNo, map[string]interface{}{
			"user_id":     shop.OwnerUserId,
			"shop_no":     shop.ShopNo,
			"scope_id":    shop.Id,
			"reason_code": reasonCode,
		})
		_, outboxErr := tx.Model(dao.SellerOutboxEvent.Table()).
			Data(do.SellerOutboxEvent{
				EventId:       event.EventID,
				EventType:     "SellerRoleRevokeRequested",
				AggregateType: "SHOP",
				AggregateNo:   shop.ShopNo,
				RequestId:     requestID,
				PayloadJson:   mustMarshalJSON(event),
				Status:        outboxStatusNew,
				AvailableAt:   gtime.Now(),
			}).
			Insert()
		if outboxErr != nil {
			return gerror.Wrap(outboxErr, "insert outbox event failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.FreezeShopRes{ShopNo: shop.ShopNo, ShopStatus: toStatus}, nil
}

// queryApplicationByNo 按业务号读取申请单，不存在时返回 CodeNotFound。
func queryApplicationByNo(ctx context.Context, applicationNo string) (*entity.SellerApplication, error) {
	var app entity.SellerApplication
	if err := dao.SellerApplication.Ctx(ctx).Where(dao.SellerApplication.Columns().ApplicationNo, applicationNo).Scan(&app); err != nil {
		return nil, gerror.Wrap(err, "query application by no failed")
	}
	if app.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "application not found")
	}
	return &app, nil
}

// queryShopByNo 按业务号读取店铺，不存在时返回 CodeNotFound。
func queryShopByNo(ctx context.Context, shopNo string) (*entity.SellerShop, error) {
	var row entity.SellerShop
	if err := dao.SellerShop.Ctx(ctx).Where(dao.SellerShop.Columns().ShopNo, shopNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query shop by no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "shop not found")
	}
	return &row, nil
}

// reserveShopNameTx 在事务内预占店名。
// 关键约束：
// 1) 每个 owner 同时最多保留 shopNameReserveLimitPerOwner 个 RESERVED；
// 2) 依赖唯一约束阻止不同申请抢同名。
func reserveShopNameTx(ctx context.Context, tx gdb.TX, ownerUserID uint64, applicationNo, shopName, shopNameNorm string) error {
	cols := dao.SellerShopNameRegistry.Columns()
	count, err := tx.Model(dao.SellerShopNameRegistry.Table()).
		Where(cols.OwnerUserId, ownerUserID).
		Where(cols.Status, registryStatusReserved).
		Count()
	if err != nil {
		return gerror.Wrap(err, "count reserved names failed")
	}
	if count >= shopNameReserveLimitPerOwner {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "shop name reserve quota exceeded")
	}

	_, err = tx.Model(dao.SellerShopNameRegistry.Table()).
		Data(do.SellerShopNameRegistry{
			ShopNameNorm:  shopNameNorm,
			ShopName:      shopName,
			OwnerUserId:   ownerUserID,
			ApplicationNo: applicationNo,
			Status:        registryStatusReserved,
			ExpireAt:      gtime.NewFromTime(time.Now().Add(shopNameReserveHours * time.Hour)),
		}).
		Insert()
	if err != nil {
		if isDuplicateErr(err) {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "shop_name already reserved or used")
		}
		return gerror.Wrap(err, "reserve shop name failed")
	}
	return nil
}

// updateReservedShopNameTx 更新申请草稿时同步更新店铺名预占记录。
func updateReservedShopNameTx(ctx context.Context, tx gdb.TX, ownerUserID uint64, applicationNo, shopName, shopNameNorm string) error {
	cols := dao.SellerShopNameRegistry.Columns()
	result, err := tx.Model(dao.SellerShopNameRegistry.Table()).
		Where(cols.ApplicationNo, applicationNo).
		Where(cols.OwnerUserId, ownerUserID).
		Data(do.SellerShopNameRegistry{
			ShopNameNorm: shopNameNorm,
			ShopName:     shopName,
			Status:       registryStatusReserved,
			ExpireAt:     gtime.NewFromTime(time.Now().Add(shopNameReserveHours * time.Hour)),
		}).
		Update()
	if err != nil {
		if isDuplicateErr(err) {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "shop_name already reserved or used")
		}
		return gerror.Wrap(err, "update reserved shop name failed")
	}
	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return gerror.Wrap(rowsErr, "get affected rows failed")
	}
	if rows == 0 {
		return reserveShopNameTx(ctx, tx, ownerUserID, applicationNo, shopName, shopNameNorm)
	}
	return nil
}

// upsertEntityTx 把审核通过快照回写到 seller_entity 主数据。
// 这一步保证 entity 表始终是“最新且已审核通过”的权威数据。
func upsertEntityTx(ctx context.Context, tx gdb.TX, ownerUserID uint64, p *v1.EntityProfile) (uint64, error) {
	cols := dao.SellerEntity.Columns()
	var row entity.SellerEntity
	if err := tx.Model(dao.SellerEntity.Table()).Where(cols.EntityNo, p.GetEntityNo()).Scan(&row); err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "no rows in result set") {
			return 0, gerror.Wrap(err, "query entity failed")
		}
	}

	legal := p.GetLegalSubject()
	if row.Id == 0 {
		result, err := tx.Model(dao.SellerEntity.Table()).
			Data(do.SellerEntity{
				EntityNo:                  p.GetEntityNo(),
				OwnerUserId:               ownerUserID,
				MerchantTypeCode:          p.GetMerchantTypeCode(),
				EntityName:                p.GetEntityName(),
				ContactName:               p.GetContactName(),
				ContactPhone:              p.GetContactPhone(),
				ContactEmail:              p.GetContactEmail(),
				SubjectTypeCode:           legal.GetSubjectTypeCode(),
				SubjectName:               legal.GetSubjectName(),
				SubjectCertNo:             legal.GetSubjectCertNo(),
				LegalRepresentativeName:   legal.GetLegalRepresentativeName(),
				LegalRepresentativeCertNo: legal.GetLegalRepresentativeCertNo(),
				CertValidFrom:             toGTime(legal.GetCertValidFrom()),
				CertValidUntil:            toGTime(legal.GetCertValidUntil()),
				LegalSubjectExtJson:       mustMarshalJSON(legal.GetExt()),
				ExtJson:                   mustMarshalJSON(p.GetExt()),
				Version:                   1,
			}).
			Insert()
		if err != nil {
			return 0, gerror.Wrap(err, "insert entity failed")
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, gerror.Wrap(err, "get entity id failed")
		}
		return uint64(id), nil
	}

	_, err := tx.Model(dao.SellerEntity.Table()).
		Where(cols.Id, row.Id).
		Data(do.SellerEntity{
			Version:                   row.Version + 1,
			OwnerUserId:               ownerUserID,
			MerchantTypeCode:          p.GetMerchantTypeCode(),
			EntityName:                p.GetEntityName(),
			ContactName:               p.GetContactName(),
			ContactPhone:              p.GetContactPhone(),
			ContactEmail:              p.GetContactEmail(),
			SubjectTypeCode:           legal.GetSubjectTypeCode(),
			SubjectName:               legal.GetSubjectName(),
			SubjectCertNo:             legal.GetSubjectCertNo(),
			LegalRepresentativeName:   legal.GetLegalRepresentativeName(),
			LegalRepresentativeCertNo: legal.GetLegalRepresentativeCertNo(),
			CertValidFrom:             toGTime(legal.GetCertValidFrom()),
			CertValidUntil:            toGTime(legal.GetCertValidUntil()),
			LegalSubjectExtJson:       mustMarshalJSON(legal.GetExt()),
			ExtJson:                   mustMarshalJSON(p.GetExt()),
		}).
		Update()
	if err != nil {
		return 0, gerror.Wrap(err, "update entity failed")
	}
	return row.Id, nil
}

// replaceEntityDocsTx 全量替换主体资质文档。
// 使用“先删后插”策略，保持 docs 与审核快照一致。
func replaceEntityDocsTx(ctx context.Context, tx gdb.TX, entityID uint64, entityNo string, docs []*v1.QualificationDoc) error {
	_, err := tx.Model(dao.SellerEntityDoc.Table()).Where(dao.SellerEntityDoc.Columns().EntityId, entityID).Delete()
	if err != nil {
		return gerror.Wrap(err, "delete entity docs failed")
	}
	for _, docItem := range docs {
		if docItem == nil {
			continue
		}
		_, insertErr := tx.Model(dao.SellerEntityDoc.Table()).
			Data(do.SellerEntityDoc{
				EntityId:    entityID,
				EntityNo:    entityNo,
				DocTypeCode: docItem.GetDocTypeCode(),
				AssetId:     docItem.GetAssetId(),
				ValidFrom:   toGTime(docItem.GetValidFrom()),
				ValidUntil:  toGTime(docItem.GetValidUntil()),
				Issuer:      docItem.GetIssuer(),
				ExtJson:     mustMarshalJSON(docItem.GetExt()),
				Status:      1,
			}).
			Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "insert entity doc failed")
		}
	}
	return nil
}

// toProtoApplication 把申请实体转换为对外 proto，包含 draft/submitted 双快照。
func toProtoApplication(app *entity.SellerApplication) *v1.SellerApplication {
	if app == nil {
		return nil
	}
	result := &v1.SellerApplication{
		ApplicationNo:         app.ApplicationNo,
		Version:               int32(app.Version),
		PreviousApplicationNo: app.PreviousApplicationNo,
		NextApplicationNo:     app.NextApplicationNo,
		OwnerUserId:           app.OwnerUserId,
		Status:                v1.ApplicationStatus(app.Status),
		LatestReject:          unmarshalRejectInfo(app.LatestRejectJson),
		SubmittedAt:           toProtoTimestamp(app.SubmittedAt),
		ReviewStartedAt:       toProtoTimestamp(app.ReviewStartedAt),
		ReviewedAt:            toProtoTimestamp(app.ReviewedAt),
		ReviewerId:            app.ReviewerId,
		ReviewComment:         app.ReviewComment,
		CreatedAt:             toProtoTimestamp(app.CreatedAt),
		UpdatedAt:             toProtoTimestamp(app.UpdatedAt),
	}
	result.EntityDraft, _ = unmarshalEntityProfile(app.EntityDraftJson)
	result.ShopDraft, _ = unmarshalShopProfile(app.ShopDraftJson)
	result.EntitySubmitted, _ = unmarshalEntityProfile(app.EntitySubmittedJson)
	result.ShopSubmitted, _ = unmarshalShopProfile(app.ShopSubmittedJson)
	return result
}

// toProtoShop 把店铺实体转换为对外 proto。
func toProtoShop(row *entity.SellerShop) *v1.Shop {
	if row == nil {
		return nil
	}
	return &v1.Shop{
		ShopNo:          row.ShopNo,
		EntityNo:        row.EntityNo,
		ApplicationNo:   row.ApplicationNo,
		OwnerUserId:     row.OwnerUserId,
		ShopName:        row.ShopName,
		ShopDisplayName: row.ShopDisplayName,
		Status:          v1.ShopStatus(row.Status),
		BuyerVisible:    row.BuyerVisible == 1,
		Version:         int32(row.Version),
		CreatedAt:       toProtoTimestamp(row.CreatedAt),
		UpdatedAt:       toProtoTimestamp(row.UpdatedAt),
	}
}

// enumApplicationStatuses 过滤掉 UNSPECIFIED，转换成数据库枚举值集合。
func enumApplicationStatuses(statuses []v1.ApplicationStatus) []uint {
	out := make([]uint, 0, len(statuses))
	for _, status := range statuses {
		if status == v1.ApplicationStatus_APPLICATION_STATUS_UNSPECIFIED {
			continue
		}
		out = append(out, uint(status))
	}
	return out
}

// normalizePaging 归一化分页参数，防止过大 page_size 压垮查询。
func normalizePaging(pageIn, sizeIn int32) (page int, size int) {
	page = int(pageIn)
	size = int(sizeIn)
	if page < defaultPage {
		page = defaultPage
	}
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return
}

// buildOutboxEvent 生成 outbox 统一事件壳，便于后续重试与对账。
func buildOutboxEvent(requestID, aggregateType, aggregateNo string, data interface{}) *outboxEventPayload {
	return &outboxEventPayload{
		EventID:       generateEventID(),
		RequestID:     requestID,
		OccurredAt:    time.Now().UTC().Format(time.RFC3339Nano),
		Producer:      "seller-shop-svc",
		AggregateType: aggregateType,
		AggregateNo:   aggregateNo,
		Data:          data,
	}
}

// metadataValue 从 gRPC metadata 读取指定键的首个值。
func metadataValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(strings.ToLower(strings.TrimSpace(key)))
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

// userIDFromMetadata 依次尝试多个 metadata key，解析出 user_id。
func userIDFromMetadata(ctx context.Context, keys ...string) (uint64, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, false
	}
	for _, key := range keys {
		values := md.Get(strings.ToLower(strings.TrimSpace(key)))
		if len(values) == 0 {
			continue
		}
		parsed, err := strconv.ParseUint(strings.TrimSpace(values[0]), 10, 64)
		if err == nil && parsed > 0 {
			return parsed, true
		}
	}
	return 0, false
}

// generateBizNo 生成对外业务号（带前缀+时间+随机后缀）。
func generateBizNo(prefix string) string {
	now := time.Now().UTC()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.Nanosecond()%1000000)
}

// generateEventID 生成事件唯一标识，用于跨服务幂等和追踪。
func generateEventID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("evt_%d_%d", now.UnixNano(), now.Nanosecond()%1000)
}

// normalizeShopName 统一店铺名比较口径（去空格+小写）。
func normalizeShopName(name string) string {
	normalized := strings.TrimSpace(strings.ToLower(name))
	normalized = strings.ReplaceAll(normalized, " ", "")
	return normalized
}

// marshalProtoJSON 把 proto message 序列化为 JSON 字符串存库。
func marshalProtoJSON(msg proto.Message) (string, error) {
	if msg == nil {
		return "{}", nil
	}
	bytes, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {
		return "", gerror.Wrap(err, "marshal proto json failed")
	}
	return string(bytes), nil
}

// unmarshalEntityProfile 反序列化数据载荷。
func unmarshalEntityProfile(raw string) (*v1.EntityProfile, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var out v1.EntityProfile
	if err := protojson.Unmarshal([]byte(raw), &out); err != nil {
		return nil, gerror.Wrap(err, "unmarshal entity profile failed")
	}
	return &out, nil
}

// unmarshalShopProfile 反序列化数据载荷。
func unmarshalShopProfile(raw string) (*v1.ShopProfile, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var out v1.ShopProfile
	if err := protojson.Unmarshal([]byte(raw), &out); err != nil {
		return nil, gerror.Wrap(err, "unmarshal shop profile failed")
	}
	return &out, nil
}

// unmarshalRejectInfo 反序列化数据载荷。
func unmarshalRejectInfo(raw string) *v1.RejectInfo {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out v1.RejectInfo
	if err := protojson.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return &out
}

// toProtoTimestamp 把 gtime 转成 protobuf timestamp。
func toProtoTimestamp(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(time.Unix(t.Timestamp(), 0).UTC())
}

// toGTime 把 protobuf timestamp 转成 gtime。
func toGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

// mustMarshalJSON 用于“尽量不因日志/扩展字段序列化失败而中断流程”的场景。
// 失败时返回 "{}"，避免主链路报错。
func mustMarshalJSON(v interface{}) string {
	if v == nil {
		return "{}"
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// nonEmptyString 返回第一个非空字符串，常用于 submitted/draft 回退读取。
func nonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

// isDuplicateErr 判断是否为唯一键冲突（兼容不同驱动报错文案）。
func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	low := strings.ToLower(err.Error())
	return strings.Contains(low, "duplicate") || strings.Contains(low, "1062")
}