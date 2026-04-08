package bff

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	adminv1 "github.com/TsingpekTao/shopa/edge-gateway/api/admin/v1"
	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/consts"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	pointsv1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sBff 负责聚合多个下游微服务的 gRPC 能力，对外输出 BFF 视角的数据组合结果。
type sBff struct {
	// once 保证所有 gRPC 连接和 client 只初始化一次，避免每次请求重复建连。
	once sync.Once

	// 下面这些连接分别对应各个下游服务，初始化后会在整个进程内复用。
	userProfileConn *grpc.ClientConn
	pointsConn      *grpc.ClientConn
	sellerShopConn  *grpc.ClientConn
	catalogConn     *grpc.ClientConn
	inventoryConn   *grpc.ClientConn
	orderConn       *grpc.ClientConn
	aftersaleConn   *grpc.ClientConn

	// 下面这些 client 对应具体的下游服务接口，供各个聚合函数直接发起 RPC 调用。
	userProfileClient userprofilev1.UserProfileServiceClient
	pointsClient      pointsv1.PointsServiceClient
	sellerAppClient   sellershopv1.SellerApplicationServiceClient
	sellerAdmin       sellershopv1.AdminSellerServiceClient
	sellerInternal    sellershopv1.InternalShopServiceClient
	catalogSeller     catalogv1.SellerProductServiceClient
	catalogAdmin      catalogv1.AdminProductReviewServiceClient
	inventoryInternal inventoryv1.InternalInventoryServiceClient
	orderBuyer        orderv1.BuyerOrderServiceClient
	orderInternal     orderv1.InternalOrderServiceClient
	orderSeller       orderv1.SellerOrderServiceClient
	aftersaleBuyer    aftersalev1.BuyerAfterSaleServiceClient
	aftersaleSeller   aftersalev1.SellerAfterSaleServiceClient

	// initErr 记录首次初始化下游连接时出现的错误，后续请求会直接复用这次初始化结果。
	initErr error
}

// New 构造 BFF 逻辑实例，供 service 层注册和注入使用。
func New() *sBff {
	// 这里仅返回空实例，真正的下游连接会在首次请求时按需初始化。
	return &sBff{}
}

func init() {
	// 在包初始化时注册 BFF 实现，方便 controller 通过 service 层统一调用。
	service.RegisterBff(New())
}

// BuildMyOverview 聚合当前登录用户的个人概览数据。
func (s *sBff) BuildMyOverview(ctx context.Context, accessToken string) (*mev1.GetOverviewRes, error) {
	// 先校验 accessToken，拿到当前请求对应的用户身份、角色和账号状态。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 个人概览依赖多个下游服务，先确保 gRPC client 已经初始化完成。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先构造基础返回体，把鉴权阶段已经拿到的身份信息直接写入响应。
	out := &mev1.GetOverviewRes{
		UserID:            verified.UserID,
		AccountStatusCode: verified.AccountStatusCode,
	}
	for _, role := range verified.Roles {
		// 把鉴权结果里的角色列表原样转成前端需要的 RoleItem 结构。
		out.Roles = append(out.Roles, mev1.RoleItem{
			RoleCode:      role.RoleCode,
			ScopeTypeCode: role.ScopeTypeCode,
			ScopeNo:       role.ScopeNo,
		})
	}

	var (
		// mu 保护 out 和 degradedFields，避免并发 goroutine 同时写同一份内存导致数据竞争。
		mu sync.Mutex
		// degradedFields 用于记录哪些下游能力降级了，前端可以据此展示部分失败状态。
		degradedFields []string
	)
	// 用 errgroup 并发拉取用户资料和积分，减少概览页的整体等待时间。
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 从用户资料服务拉取昵称和头像；地址信息在概览页不需要，所以显式关闭。
		res, err := s.userProfileClient.GetProfileByUserId(egCtx, &userprofilev1.GetProfileByUserIdReq{
			UserId:           verified.UserID,
			IncludeAddresses: false,
		})
		if err != nil {
			// 用户资料查询失败时不让整个页面报错，只标记该字段降级。
			mu.Lock()
			degradedFields = append(degradedFields, "user_profile")
			mu.Unlock()
			return nil
		}
		if res.GetProfile() != nil {
			// 只有 profile 存在时才回填展示信息，避免空指针并保留默认空值语义。
			mu.Lock()
			out.DisplayName = res.GetProfile().GetDisplayName()
			if res.GetProfile().GetAvatar() != nil {
				// 头像是可选字段，需要二次判空后再取 URL。
				out.AvatarURL = res.GetProfile().GetAvatar().GetUrl()
			}
			mu.Unlock()
		}
		return nil
	})

	eg.Go(func() error {
		// 并发调用积分服务，补齐用户当前可展示的积分余额。
		res, err := s.pointsClient.GetPointsByUserId(egCtx, &pointsv1.GetPointsByUserIdReq{UserId: verified.UserID})
		if err != nil {
			// 积分失败同样走降级路径，避免单个下游异常拖垮整个概览页。
			mu.Lock()
			degradedFields = append(degradedFields, "points")
			mu.Unlock()
			return nil
		}
		// 成功时回填积分余额，供用户中心首页直接展示。
		mu.Lock()
		out.Points = res.GetBalance()
		mu.Unlock()
		return nil
	})

	// 等待所有并发查询结束，再统一计算是否为部分返回。
	_ = eg.Wait()
	// 只要存在任意降级字段，就把 Partial 标成 true，提醒前端结果不完整。
	out.Partial = len(degradedFields) > 0
	// 对降级字段去重，避免多个分支重复追加造成前端展示噪音。
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// BuildAdminOverview 构造管理员个人概览视图。
func (s *sBff) BuildAdminOverview(ctx context.Context, accessToken string) (*adminv1.GetOverviewRes, error) {
	// 先校验 accessToken，确认当前请求确实来自已登录主体。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "admin:me:overview") {
		// 管理员概览要求显式权限，避免普通用户读取后台身份信息。
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}

	// 只收集非空角色编码，避免把空角色透出到前端。
	roles := make([]string, 0, len(verified.Roles))
	for _, role := range verified.Roles {
		if role.RoleCode == "" {
			continue
		}
		roles = append(roles, role.RoleCode)
	}

	// 这里直接复用鉴权结果返回，不需要额外访问下游服务。
	return &adminv1.GetOverviewRes{
		UserID:            verified.UserID,
		AccountStatusCode: verified.AccountStatusCode,
		Roles:             uniqueStrings(roles),
		Permissions:       uniqueStrings(verified.Permissions),
	}, nil
}

// BuildAdminDashboardOverview 汇总平台后台看板的核心指标。
func (s *sBff) BuildAdminDashboardOverview(ctx context.Context, accessToken string) (*adminv1.GetDashboardOverviewRes, error) {
	// 后台看板是管理视角数据，必须先校验 token 和看板权限。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "dashboard:mall:view") {
		// 没有看板权限时直接拒绝，避免敏感经营数据泄露。
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}
	// 看板依赖商家和商品审核两个下游服务，先确保 client 可用。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// metrics 先初始化为 0，保证下游降级时仍能返回结构完整的默认值。
	metrics := map[string]int64{
		"totalShops":             0,
		"pendingMerchantReviews": 0,
		"pendingProductReviews":  0,
	}
	// degradedFields 记录哪些指标因为下游失败而退化为默认值。
	degradedFields := make([]string, 0, 3)
	// 多个 goroutine 会同时写 metrics 和 degradedFields，需要互斥保护。
	var mu sync.Mutex
	// 并发查询多个统计项，减少管理看板的总体耗时。
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 通过只取 1 条数据的分页请求复用 total 字段，统计已审核通过的店铺数。
		res, err := s.sellerAdmin.ListApplications(egCtx, &sellershopv1.ListApplicationsReq{
			Page:     1,
			PageSize: 1,
			Statuses: []sellershopv1.ApplicationStatus{
				sellershopv1.ApplicationStatus_APPLICATION_STATUS_APPROVED,
			},
		})
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			// 单个指标失败时只标记降级，不影响其他指标继续返回。
			degradedFields = append(degradedFields, "total_shops")
			return nil
		}
		metrics["totalShops"] = res.GetTotal()
		return nil
	})

	eg.Go(func() error {
		// 同样复用申请列表 total，统计待审核中的商家入驻申请数。
		res, err := s.sellerAdmin.ListApplications(egCtx, &sellershopv1.ListApplicationsReq{
			Page:     1,
			PageSize: 1,
			Statuses: []sellershopv1.ApplicationStatus{
				sellershopv1.ApplicationStatus_APPLICATION_STATUS_SUBMITTED,
				sellershopv1.ApplicationStatus_APPLICATION_STATUS_REVIEWING,
			},
		})
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			degradedFields = append(degradedFields, "pending_merchant_reviews")
			return nil
		}
		metrics["pendingMerchantReviews"] = res.GetTotal()
		return nil
	})

	eg.Go(func() error {
		// 商品审核任务也只拉 1 条记录，通过 total 汇总当前待审商品数。
		res, err := s.catalogAdmin.ListReviewTasks(egCtx, &catalogv1.ListReviewTasksReq{
			Page:     1,
			PageSize: 1,
			Statuses: []catalogv1.SpuStatus{
				catalogv1.SpuStatus_SPU_STATUS_REVIEWING,
			},
		})
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			degradedFields = append(degradedFields, "pending_product_reviews")
			return nil
		}
		metrics["pendingProductReviews"] = res.GetTotal()
		return nil
	})

	// 看板允许部分降级，所以这里忽略 goroutine 内部已经吞掉的单项错误。
	_ = eg.Wait()

	// 返回时统一附带降级标记，前端可据此提示某些指标是兜底值。
	return &adminv1.GetDashboardOverviewRes{
		Metrics: []adminv1.DashboardMetric{
			{Key: "totalShops", Label: "已开通店铺", Value: metrics["totalShops"]},
			{Key: "pendingMerchantReviews", Label: "待审商家入驻", Value: metrics["pendingMerchantReviews"]},
			{Key: "pendingProductReviews", Label: "待审商品", Value: metrics["pendingProductReviews"]},
		},
		Partial:        len(degradedFields) > 0,
		DegradedFields: uniqueStrings(degradedFields),
	}, nil
}

// BuildAdminShopInsights 聚合单个店铺的基础洞察信息。
func (s *sBff) BuildAdminShopInsights(ctx context.Context, accessToken string, shopNo string) (*adminv1.GetShopInsightsRes, error) {
	// 先校验管理员身份，避免未授权主体读取店铺洞察。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "dashboard:shop:view") {
		// 店铺洞察属于后台能力，必须具备对应权限才能访问。
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}
	if strings.TrimSpace(shopNo) == "" {
		// shopNo 是查询目标，没有它就无法定位店铺。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	// 洞察页需要访问店铺服务，先确保下游 client 已初始化。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先构造默认返回体，即使下游降级也能返回最小可用结构。
	out := &adminv1.GetShopInsightsRes{
		ShopNo:         shopNo,
		ShopName:       "",
		ShopStatusCode: "",
	}
	// 通过店铺内部服务查询基础店铺信息，补齐名称和状态。
	shopRes, shopErr := s.sellerInternal.GetShopByNo(ctx, &sellershopv1.GetShopByNoReq{ShopNo: shopNo})
	if shopErr != nil || shopRes.GetShop() == nil {
		// 如果店铺基础信息不可用，则降级返回并显式标记缺失字段。
		out.Partial = true
		out.DegradedFields = []string{"shop_basic_info_unavailable"}
		return out, nil
	}
	// 查询成功时把店铺名称和状态码透出给前端。
	out.ShopName = shopRes.GetShop().GetShopName()
	out.ShopStatusCode = enumCode(shopRes.GetShop().GetStatus().String(), "SHOP_STATUS_")
	return out, nil
}

// BuildAdminConversationList 当前只返回降级占位结果，因为会话服务尚未接入。
func (s *sBff) BuildAdminConversationList(ctx context.Context, accessToken string) (*adminv1.ListConversationsRes, error) {
	// 先校验管理员 accessToken，确认请求主体合法。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "cs:conversation:view") {
		// 会话列表属于客服后台能力，没有权限时必须拒绝访问。
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}
	// 当前实现还没接入真正的会话服务，因此返回空列表并显式声明为降级结果。
	return &adminv1.ListConversationsRes{
		Items:          []adminv1.ConversationItem{},
		Partial:        true,
		DegradedFields: []string{"conversation_service_not_integrated"},
	}, nil
}

// BuildSellerWorkbench 构建卖家工作台首页的聚合数据。
func (s *sBff) BuildSellerWorkbench(ctx context.Context, accessToken string) (*sellerv1.GetWorkbenchRes, error) {
	// 先校验 accessToken，拿到当前卖家主体身份。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 工作台依赖多个下游服务，先统一初始化 gRPC client。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先构造默认返回体，确保任一模块降级时仍有稳定的兜底结构。
	out := &sellerv1.GetWorkbenchRes{
		UserID: verified.UserID,
		ProductSummary: sellerv1.SellerProductSummary{
			Total:     0,
			OnShelf:   0,
			OffShelf:  0,
			Reviewing: 0,
			Draft:     0,
			Rejected:  0,
		},
	}

	var (
		// mu 保护 out 和 degradedFields，避免多个异步查询同时写响应对象。
		mu sync.Mutex
		// degradedFields 用于记录店铺、申请单、商品汇总等模块的降级原因。
		degradedFields []string
	)
	// 各个卡片数据彼此独立，使用 errgroup 并发拉取可以缩短工作台首屏耗时。
	eg, egCtx := errgroup.WithContext(ctx)
	// 某些卖家侧服务依赖 x-user-id metadata 做主体识别，这里提前注入到下游上下文。
	sellerUserCtx := withOutgoingUserMetadata(egCtx, verified.UserID)

	eg.Go(func() error {
		// 查询当前用户名下的店铺列表，用于工作台展示最近的店铺概览。
		res, err := s.sellerInternal.ListShopsByOwnerUserId(egCtx, &sellershopv1.ListShopsByOwnerUserIdReq{
			OwnerUserId: verified.UserID,
		})
		if err != nil {
			// 店铺列表失败时只标记 seller_shops 降级，不影响其他模块继续返回。
			mu.Lock()
			degradedFields = append(degradedFields, "seller_shops")
			mu.Unlock()
			return nil
		}

		// shops 保存当前卖家名下所有店铺，后续会做截断展示。
		shops := res.GetShops()
		// total 记录真实店铺总数，供前端显示“共多少家店”。
		total := uint64(len(shops))
		// 工作台只展示前 10 家店，避免首页列表过长影响可读性。
		limit := 10
		if len(shops) < limit {
			limit = len(shops)
		}
		// list 只承载截断后的展示数据，和 total 一起表达“部分展示”的语义。
		list := make([]sellerv1.SellerShopSummary, 0, limit)
		for _, row := range shops[:limit] {
			list = append(list, sellerv1.SellerShopSummary{
				ShopNo:          row.GetShopNo(),
				ShopName:        row.GetShopName(),
				ShopDisplayName: row.GetShopDisplayName(),
				ShopStatusCode:  enumCode(row.GetStatus().String(), "SHOP_STATUS_"),
				BuyerVisible:    row.GetBuyerVisible(),
				UpdatedAt:       "",
			})
		}

		// 成功后一次性回写店铺列表、总数和是否截断标记，保持响应字段一致性。
		mu.Lock()
		out.Shops = list
		out.ShopsTotal = total
		out.ShopsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 拉取当前卖家的入驻申请列表，工作台只展示最近一小部分申请记录。
		res, err := s.sellerAppClient.ListMyApplications(sellerUserCtx, &sellershopv1.ListMyApplicationsReq{
			Page:     1,
			PageSize: 20,
		})
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "seller_applications")
			mu.Unlock()
			return nil
		}

		// apps 保存服务返回的申请列表，后续会截断为首页所需的最近 5 条。
		apps := res.GetApplications()
		// total 记录申请单总数，便于前端给出完整数量感知。
		total := uint64(res.GetTotal())
		// 工作台只展示最近 5 条申请，避免把列表页能力搬到首页。
		limit := 5
		if len(apps) < limit {
			limit = len(apps)
		}
		// list 只保留首页卡片需要的关键信息，减少无关字段透出。
		list := make([]sellerv1.SellerApplicationSummary, 0, limit)
		for _, row := range apps[:limit] {
			list = append(list, sellerv1.SellerApplicationSummary{
				ApplicationNo:         row.GetApplicationNo(),
				ApplicationStatusCode: enumCode(row.GetStatus().String(), "APPLICATION_STATUS_"),
				Version:               int32(row.GetVersion()),
				SubmittedAt:           tsToString(row.GetSubmittedAt()),
				UpdatedAt:             tsToString(row.GetUpdatedAt()),
			})
		}

		// 成功后回填申请摘要和截断信息，供工作台卡片直接渲染。
		mu.Lock()
		out.LatestApplications = list
		out.ApplicationsTotal = total
		out.ApplicationsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 商品汇总单独查询，避免某个商品统计失败时拖累店铺和申请模块。
		summary, err := s.loadSellerProductSummary(egCtx)
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "catalog_products")
			mu.Unlock()
			return nil
		}
		// 成功时用统计结果覆盖默认值。
		mu.Lock()
		out.ProductSummary = summary
		mu.Unlock()
		return nil
	})

	// 工作台允许部分模块降级，所以这里等待所有任务结束后统一收口。
	_ = eg.Wait()
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// BuildSellerShopDashboard 提供单个店铺维度的卖家看板。
func (s *sBff) BuildSellerShopDashboard(ctx context.Context, accessToken, shopNo string) (*sellerv1.GetShopDashboardRes, error) {
	// shopNo 是店铺看板的定位参数，缺失时直接判为非法请求。
	if strings.TrimSpace(shopNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	// 先校验 accessToken，确认当前请求主体合法。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 后续会访问多个下游服务，先确保 client 都已经初始化完成。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先校验当前用户是否真的是该店铺 owner，避免越权读取店铺经营数据。
	ownerRes, err := s.sellerInternal.IsUserShopOwner(ctx, &sellershopv1.IsUserShopOwnerReq{
		UserId: verified.UserID,
		ShopNo: shopNo,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "check shop owner failed")
	}
	if !ownerRes.GetIsOwner() {
		// 非店主没有资格查看店铺看板，必须直接拒绝。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "user is not owner of this shop")
	}

	// 先构造默认返回体，保证任一子模块降级时也能稳定返回结构。
	out := &sellerv1.GetShopDashboardRes{
		ShopNo: shopNo,
		ProductSummary: sellerv1.SellerProductSummary{
			Total:     0,
			OnShelf:   0,
			OffShelf:  0,
			Reviewing: 0,
			Draft:     0,
			Rejected:  0,
		},
		InventoryRisk: sellerv1.InventoryRiskSummary{
			LowStockSkuCount:   0,
			OutOfStockSkuCount: 0,
		},
	}

	var (
		// mu 保护 out 和 degradedFields，避免并发任务写入冲突。
		mu sync.Mutex
		// degradedFields 记录哪些看板模块因为下游能力不足而走了降级路径。
		degradedFields []string
	)
	// 店铺信息、商品汇总、库存风险彼此独立，适合并发拉取。
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 查询店铺基础信息，补齐看板顶部所需的店铺名和状态码。
		shopRes, err := s.sellerInternal.GetShopByNo(egCtx, &sellershopv1.GetShopByNoReq{ShopNo: shopNo})
		if err != nil {
			return gerror.Wrap(err, "query shop failed")
		}
		if shopRes.GetShop() == nil {
			return gerror.NewCode(gcode.CodeNotFound, "shop not found")
		}
		// 店铺基础信息是看板主干数据，成功后立即写回响应对象。
		mu.Lock()
		out.ShopName = shopRes.GetShop().GetShopName()
		out.ShopStatusCode = enumCode(shopRes.GetShop().GetStatus().String(), "SHOP_STATUS_")
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 商品汇总直接复用公共方法，减少重复统计代码。
		summary, err := s.loadSellerProductSummary(egCtx)
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "catalog_products")
			mu.Unlock()
			return nil
		}
		mu.Lock()
		out.ProductSummary = summary
		// 当前汇总没有按 shopNo 过滤，所以要明确标记这是一个降级能力，提醒前端不要误解。
		degradedFields = append(degradedFields, "product_summary_not_filtered_by_shop")
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 库存风险能力目前还没真正接入，只保留字段占位并显式标记降级。
		_ = s.inventoryInternal
		mu.Lock()
		degradedFields = append(degradedFields, "inventory_risk_unavailable")
		mu.Unlock()
		return nil
	})

	if err = eg.Wait(); err != nil {
		// 只要主干 goroutine 返回硬错误，就直接终止整个店铺看板请求。
		return nil, err
	}
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// loadSellerProductSummary 按商品状态汇总卖家侧商品数量。
func (s *sBff) loadSellerProductSummary(ctx context.Context) (sellerv1.SellerProductSummary, error) {
	// statusCase 把返回字段 key 和下游枚举状态绑定起来，便于后续循环统计。
	type statusCase struct {
		key    string
		status catalogv1.SpuStatus
	}
	cases := []statusCase{
		{key: "on", status: catalogv1.SpuStatus_SPU_STATUS_ON_SHELF},
		{key: "off", status: catalogv1.SpuStatus_SPU_STATUS_OFF_SHELF},
		{key: "reviewing", status: catalogv1.SpuStatus_SPU_STATUS_REVIEWING},
		{key: "draft", status: catalogv1.SpuStatus_SPU_STATUS_DRAFT},
		{key: "rejected", status: catalogv1.SpuStatus_SPU_STATUS_REJECTED},
	}

	var (
		// mu 保护 counter，避免多个状态统计并发写 map 触发数据竞争。
		mu sync.Mutex
		// counter 暂存各个商品状态对应的总数，最后统一组装成返回结构。
		counter = map[string]uint64{
			"on":        0,
			"off":       0,
			"reviewing": 0,
			"draft":     0,
			"rejected":  0,
		}
	)

	// 每个状态都可以独立统计，所以用并发请求缩短总体耗时。
	eg, egCtx := errgroup.WithContext(ctx)
	for _, item := range cases {
		// 复制循环变量，避免 goroutine 捕获同一个 item 导致统计结果串值。
		item := item
		eg.Go(func() error {
			// 这里只请求 1 条记录，真正需要的是 total 字段返回的状态总数。
			res, err := s.catalogSeller.ListMyProducts(egCtx, &catalogv1.ListMyProductsReq{
				Page:     1,
				PageSize: 1,
				Statuses: []catalogv1.SpuStatus{item.status},
			})
			if err != nil {
				return err
			}
			// 每个状态统计完成后写入对应 key，供最终汇总使用。
			mu.Lock()
			counter[item.key] = uint64(res.GetTotal())
			mu.Unlock()
			return nil
		})
	}

	// 只要任一状态统计失败，就让调用方决定是否整块降级。
	if err := eg.Wait(); err != nil {
		return sellerv1.SellerProductSummary{}, err
	}

	// 再补发一次不带状态过滤的请求，用 total 字段拿到商品总量。
	totalRes, err := s.catalogSeller.ListMyProducts(ctx, &catalogv1.ListMyProductsReq{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return sellerv1.SellerProductSummary{}, err
	}

	// 最终把分状态数量和总量一起组装成卖家工作台需要的商品摘要。
	return sellerv1.SellerProductSummary{
		Total:     uint64(totalRes.GetTotal()),
		OnShelf:   counter["on"],
		OffShelf:  counter["off"],
		Reviewing: counter["reviewing"],
		Draft:     counter["draft"],
		Rejected:  counter["rejected"],
	}, nil
}

// ensureClients 以单例方式初始化所有下游 gRPC 连接和 client。
func (s *sBff) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		var err error

		// 先初始化用户资料服务连接，后续个人概览等接口都会依赖它。
		s.userProfileConn, err = dialUpstream(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化积分服务连接，供个人概览等场景读取积分余额。
		s.pointsConn, err = dialUpstream(ctx, "upstream.pointsGrpc", "127.0.0.1:9021")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化 seller-shop 服务连接，供店铺、申请单、所有者校验等能力复用。
		s.sellerShopConn, err = dialUpstream(ctx, "upstream.sellerShopGrpc", "127.0.0.1:50051")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化 catalog 服务连接，供商品审核和商品汇总统计使用。
		s.catalogConn, err = dialUpstream(ctx, "upstream.catalogGrpc", "127.0.0.1:9004")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化 inventory 服务连接，虽然当前库存风险仍是占位，但预留 client 方便后续接入。
		s.inventoryConn, err = dialUpstream(ctx, "upstream.inventoryGrpc", "127.0.0.1:9005")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化 order 服务连接，供卖家工作台与销售报表读取订单支付数据。
		s.orderConn, err = dialUpstream(ctx, "upstream.orderGrpc", "127.0.0.1:9008")
		if err != nil {
			s.initErr = err
			return
		}
		// 初始化 aftersale 服务连接，供销售分析聚合退款申请数据。
		s.aftersaleConn, err = dialUpstream(ctx, "upstream.aftersaleGrpc", "127.0.0.1:9010")
		if err != nil {
			s.initErr = err
			return
		}

		// 连接建立完成后，统一构造各个服务的 typed client，供业务逻辑直接使用。
		s.userProfileClient = userprofilev1.NewUserProfileServiceClient(s.userProfileConn)
		s.pointsClient = pointsv1.NewPointsServiceClient(s.pointsConn)
		s.sellerAppClient = sellershopv1.NewSellerApplicationServiceClient(s.sellerShopConn)
		s.sellerAdmin = sellershopv1.NewAdminSellerServiceClient(s.sellerShopConn)
		s.sellerInternal = sellershopv1.NewInternalShopServiceClient(s.sellerShopConn)
		s.catalogSeller = catalogv1.NewSellerProductServiceClient(s.catalogConn)
		s.catalogAdmin = catalogv1.NewAdminProductReviewServiceClient(s.catalogConn)
		s.inventoryInternal = inventoryv1.NewInternalInventoryServiceClient(s.inventoryConn)
		s.orderBuyer = orderv1.NewBuyerOrderServiceClient(s.orderConn)
		s.orderInternal = orderv1.NewInternalOrderServiceClient(s.orderConn)
		s.orderSeller = orderv1.NewSellerOrderServiceClient(s.orderConn)
		s.aftersaleBuyer = aftersalev1.NewBuyerAfterSaleServiceClient(s.aftersaleConn)
		s.aftersaleSeller = aftersalev1.NewSellerAfterSaleServiceClient(s.aftersaleConn)
	})
	return s.initErr
}

// dialUpstream 从配置里读取下游地址并建立 gRPC 连接。
func dialUpstream(ctx context.Context, cfgKey, defaultAddr string) (*grpc.ClientConn, error) {
	// 先读取配置里的地址，没有配置时退回到代码内置默认地址。
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, cfgKey, defaultAddr).String())
	// 建连使用独立超时上下文，避免上游请求 ctx 被取消后影响初始化控制。
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// withOutgoingUserMetadata 把当前登录用户 ID 注入下游 gRPC metadata，供用户态服务按主体做权限过滤。
func withOutgoingUserMetadata(ctx context.Context, userID uint64) context.Context {
	if userID == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(userID, 10))
}

// tsToString 把 protobuf Timestamp 转成统一的 RFC3339 字符串格式。
func tsToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		// 时间为空时返回空字符串，保持前端已有的空值语义。
		return ""
	}
	return ts.AsTime().Format(time.RFC3339)
}

// uniqueStrings 对字符串列表做去重并过滤空白项。
func uniqueStrings(in []string) []string {
	// m 用来记录已经出现过的值，避免重复写入结果切片。
	m := make(map[string]struct{}, len(in))
	// out 保留原始遍历顺序，方便前端按追加顺序展示降级字段。
	out := make([]string, 0, len(in))
	for _, item := range in {
		// 统一先去掉首尾空白，避免空串和空白串污染结果。
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := m[item]; ok {
			continue
		}
		m[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

// hasPermission 判断权限列表里是否包含目标权限。
func hasPermission(granted []string, required string) bool {
	required = strings.TrimSpace(required)
	if required == "" {
		// 调用方没有声明必需权限时，默认视为无需拦截。
		return true
	}
	for _, item := range granted {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == required || item == "*:*:*" {
			// 命中精确权限或全量通配权限时即可放行。
			return true
		}
	}
	return false
}

// enumCode 去掉枚举常量前缀，把下游枚举名转成前端更易消费的 code。
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
