package bff

import (
	"context"
	"strings"
	"sync"
	"time"

	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	pointsv1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sBff 负责网关 BFF 聚合：
// 1) 聚合多个下游服务，输出前端友好结构。
// 2) 对旁路失败执行“字段级降级”，而不是整页失败。
// 3) 对列表做截断控制，避免响应体过重。
type sBff struct {
	// once 保证所有下游 gRPC 客户端仅初始化一次。
	once sync.Once

	// 以下是到下游服务的长连接。
	userProfileConn *grpc.ClientConn
	pointsConn      *grpc.ClientConn
	sellerShopConn  *grpc.ClientConn
	catalogConn     *grpc.ClientConn
	inventoryConn   *grpc.ClientConn

	// 以下是对应下游 RPC 客户端。
	userProfileClient userprofilev1.UserProfileServiceClient
	pointsClient      pointsv1.PointsServiceClient
	sellerAppClient   sellershopv1.SellerApplicationServiceClient
	sellerInternal    sellershopv1.InternalShopServiceClient
	catalogSeller     catalogv1.SellerProductServiceClient
	inventoryInternal inventoryv1.InternalInventoryServiceClient

	// initErr 记录初始化失败原因，避免每次请求都重复建连重试。
	initErr error
}

// New 创建 BFF 服务实例。
func New() *sBff {
	return &sBff{}
}

// init 启动时注册 BFF 实现到 service 门面。
func init() {
	service.RegisterBff(New())
}

// BuildMyOverview 构建“我的概览”聚合数据：
// - 主身份信息来自 IAM 鉴权结果。
// - 资料来自 user-profile。
// - 积分来自 points。
// - 任一路失败时返回 partial=true，并标注 degraded_fields。
func (s *sBff) BuildMyOverview(ctx context.Context, accessToken string) (*mev1.GetOverviewRes, error) {
	// 第一步串行鉴权：身份无效时直接失败，不继续下游调用。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 确保下游客户端已初始化。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先填充可立即得到的核心身份字段。
	out := &mev1.GetOverviewRes{
		UserID:            verified.UserID,
		AccountStatusCode: verified.AccountStatusCode,
	}
	// 透传角色给前端，用于入口显隐和能力控制。
	for _, role := range verified.Roles {
		out.Roles = append(out.Roles, mev1.RoleItem{
			RoleCode:      role.RoleCode,
			ScopeTypeCode: role.ScopeTypeCode,
			ScopeNo:       role.ScopeNo,
		})
	}

	var (
		// mu 保护并发写 out 与 degradedFields。
		mu sync.Mutex
		// degradedFields 记录本次聚合中发生降级的字段。
		degradedFields []string
	)
	// errgroup 并发请求多个下游，降低整体 RT。
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 拉取基础资料（昵称/头像）；地址不需要，减少下游负载。
		res, err := s.userProfileClient.GetProfileByUserId(egCtx, &userprofilev1.GetProfileByUserIdReq{
			UserId:           verified.UserID,
			IncludeAddresses: false,
		})
		if err != nil {
			// profile 失败仅记降级，不中断整页。
			mu.Lock()
			degradedFields = append(degradedFields, "user_profile")
			mu.Unlock()
			return nil
		}
		if res.GetProfile() != nil {
			mu.Lock()
			out.DisplayName = res.GetProfile().GetDisplayName()
			if res.GetProfile().GetAvatar() != nil {
				out.AvatarURL = res.GetProfile().GetAvatar().GetUrl()
			}
			mu.Unlock()
		}
		return nil
	})

	eg.Go(func() error {
		// 拉取积分；失败时降级为默认 0（保持字段存在）。
		res, err := s.pointsClient.GetPointsByUserId(egCtx, &pointsv1.GetPointsByUserIdReq{UserId: verified.UserID})
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "points")
			mu.Unlock()
			return nil
		}
		mu.Lock()
		out.Points = res.GetBalance()
		mu.Unlock()
		return nil
	})

	// 每个 goroutine 已将错误“转译”为降级信号，这里不再硬失败。
	_ = eg.Wait()
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// BuildSellerWorkbench 构建卖家工作台：
// - 店铺摘要（最多 10 条）。
// - 最近申请（最多 5 条）。
// - 商品状态聚合摘要。
func (s *sBff) BuildSellerWorkbench(ctx context.Context, accessToken string) (*sellerv1.GetWorkbenchRes, error) {
	// 串行鉴权，确保 user_id 可信。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 确保 gRPC 客户端可用。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 预置默认响应，防止降级时字段缺失。
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
		// mu 保护并发写响应对象。
		mu sync.Mutex
		// degradedFields 记录本次工作台聚合中降级项。
		degradedFields []string
	)
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 拉取“我名下店铺”。
		res, err := s.sellerInternal.ListShopsByOwnerUserId(egCtx, &sellershopv1.ListShopsByOwnerUserIdReq{
			OwnerUserId: verified.UserID,
		})
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "seller_shops")
			mu.Unlock()
			return nil
		}

		shops := res.GetShops()
		// total 保留真实总数，前端可据此展示“更多”。
		total := uint64(len(shops))
		// 工作台仅展示前 10 条，避免首屏过载。
		limit := 10
		if len(shops) < limit {
			limit = len(shops)
		}
		// 预分配容量，减少 append 扩容开销。
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

		mu.Lock()
		out.Shops = list
		out.ShopsTotal = total
		// 截断标记告诉前端：这里只是摘要，不是全量。
		out.ShopsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 拉取“我的申请”，用于工作台最近申请卡片。
		res, err := s.sellerAppClient.ListMyApplications(egCtx, &sellershopv1.ListMyApplicationsReq{
			Page:     1,
			PageSize: 20,
		})
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "seller_applications")
			mu.Unlock()
			return nil
		}

		apps := res.GetApplications()
		// total 使用分页接口返回值，而非当前页长度。
		total := uint64(res.GetTotal())
		// 工作台固定展示最近 5 条。
		limit := 5
		if len(apps) < limit {
			limit = len(apps)
		}
		list := make([]sellerv1.SellerApplicationSummary, 0, limit)
		for _, row := range apps[:limit] {
			list = append(list, sellerv1.SellerApplicationSummary{
				ApplicationNo:         row.GetApplicationNo(),
				ApplicationStatusCode: enumCode(row.GetStatus().String(), "APPLICATION_STATUS_"),
				// version 给前端后续做 expected_version 并发控制。
				Version:     int32(row.GetVersion()),
				SubmittedAt: tsToString(row.GetSubmittedAt()),
				UpdatedAt:   tsToString(row.GetUpdatedAt()),
			})
		}

		mu.Lock()
		out.LatestApplications = list
		out.ApplicationsTotal = total
		// 截断标记用于触发“查看全部申请”入口。
		out.ApplicationsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 并发拉取商品状态摘要，避免串行拖慢工作台。
		summary, err := s.loadSellerProductSummary(egCtx)
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "catalog_products")
			mu.Unlock()
			return nil
		}
		mu.Lock()
		out.ProductSummary = summary
		mu.Unlock()
		return nil
	})

	// 旁路错误在 goroutine 内已转为降级，此处无需硬失败。
	_ = eg.Wait()
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// BuildSellerShopDashboard 构建店铺仪表盘：
// 1) 先鉴权并校验 shop_no 归属。
// 2) 并发拉取店铺信息、商品摘要、库存风险摘要。
// 3) 暂缺能力用 degraded_fields 明确标注。
func (s *sBff) BuildSellerShopDashboard(ctx context.Context, accessToken, shopNo string) (*sellerv1.GetShopDashboardRes, error) {
	// shopNo 是店铺维度主键，不允许为空。
	if strings.TrimSpace(shopNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	// 统一鉴权。
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 初始化下游客户端。
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 归属校验：仅店铺 owner 可访问该 dashboard。
	ownerRes, err := s.sellerInternal.IsUserShopOwner(ctx, &sellershopv1.IsUserShopOwnerReq{
		UserId: verified.UserID,
		ShopNo: shopNo,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "check shop owner failed")
	}
	if !ownerRes.GetIsOwner() {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "user is not owner of this shop")
	}

	// 先给默认结构，确保降级时响应结构稳定。
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
		// mu 保护并发写共享结果。
		mu sync.Mutex
		// degradedFields 汇总本次请求的降级项。
		degradedFields []string
	)
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 查询店铺基本信息（名称、状态）。
		shopRes, err := s.sellerInternal.GetShopByNo(egCtx, &sellershopv1.GetShopByNoReq{ShopNo: shopNo})
		if err != nil {
			return gerror.Wrap(err, "query shop failed")
		}
		if shopRes.GetShop() == nil {
			return gerror.NewCode(gcode.CodeNotFound, "shop not found")
		}
		mu.Lock()
		out.ShopName = shopRes.GetShop().GetShopName()
		out.ShopStatusCode = enumCode(shopRes.GetShop().GetStatus().String(), "SHOP_STATUS_")
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 当前 catalog-svc 暂不支持按 shop_no 精确过滤。
		// 这里回退到卖家级摘要，并显式标注为降级字段。
		summary, err := s.loadSellerProductSummary(egCtx)
		if err != nil {
			mu.Lock()
			degradedFields = append(degradedFields, "catalog_products")
			mu.Unlock()
			return nil
		}
		mu.Lock()
		out.ProductSummary = summary
		degradedFields = append(degradedFields, "product_summary_not_filtered_by_shop")
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 当前 inventory-svc 暂无 shop 维度库存风险聚合接口。
		// 因此保留默认值并标注降级，而不是让整页失败。
		_ = s.inventoryInternal
		mu.Lock()
		degradedFields = append(degradedFields, "inventory_risk_unavailable")
		mu.Unlock()
		return nil
	})

	if err = eg.Wait(); err != nil {
		return nil, err
	}
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// loadSellerProductSummary 聚合卖家商品状态统计：
// 并发调用 ListMyProducts(total-only) 拉取各状态计数，再回填汇总结构。
func (s *sBff) loadSellerProductSummary(ctx context.Context) (sellerv1.SellerProductSummary, error) {
	// statusCase 用于抽象“统计键”与“下游状态枚举”的映射关系。
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
		// mu 保护 counter 并发写。
		mu sync.Mutex
		// counter 存放各状态商品数。
		counter = map[string]uint64{
			"on":        0,
			"off":       0,
			"reviewing": 0,
			"draft":     0,
			"rejected":  0,
		}
	)

	eg, egCtx := errgroup.WithContext(ctx)
	for _, item := range cases {
		// 重新绑定循环变量，避免 goroutine 闭包捕获同一地址。
		item := item
		eg.Go(func() error {
			// 只查总数，PageSize=1 即可，减少下游返回负载。
			res, err := s.catalogSeller.ListMyProducts(egCtx, &catalogv1.ListMyProductsReq{
				Page:     1,
				PageSize: 1,
				Statuses: []catalogv1.SpuStatus{item.status},
			})
			if err != nil {
				return err
			}
			mu.Lock()
			counter[item.key] = uint64(res.GetTotal())
			mu.Unlock()
			return nil
		})
	}

	// 任一状态统计失败时，交由上层决定是否降级。
	if err := eg.Wait(); err != nil {
		return sellerv1.SellerProductSummary{}, err
	}

	// 再拉一次不带状态过滤的总数。
	totalRes, err := s.catalogSeller.ListMyProducts(ctx, &catalogv1.ListMyProductsReq{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return sellerv1.SellerProductSummary{}, err
	}

	return sellerv1.SellerProductSummary{
		Total:     uint64(totalRes.GetTotal()),
		OnShelf:   counter["on"],
		OffShelf:  counter["off"],
		Reviewing: counter["reviewing"],
		Draft:     counter["draft"],
		Rejected:  counter["rejected"],
	}, nil
}

// ensureClients 惰性初始化所有下游 gRPC 客户端。
func (s *sBff) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		var err error

		// 用户资料服务。
		s.userProfileConn, err = dialUpstream(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002")
		if err != nil {
			s.initErr = err
			return
		}
		// 积分服务。
		s.pointsConn, err = dialUpstream(ctx, "upstream.pointsGrpc", "127.0.0.1:8012")
		if err != nil {
			s.initErr = err
			return
		}
		// 卖家店铺服务。
		s.sellerShopConn, err = dialUpstream(ctx, "upstream.sellerShopGrpc", "127.0.0.1:50051")
		if err != nil {
			s.initErr = err
			return
		}
		// 商品服务。
		s.catalogConn, err = dialUpstream(ctx, "upstream.catalogGrpc", "127.0.0.1:9004")
		if err != nil {
			s.initErr = err
			return
		}
		// 库存服务。
		s.inventoryConn, err = dialUpstream(ctx, "upstream.inventoryGrpc", "127.0.0.1:9005")
		if err != nil {
			s.initErr = err
			return
		}

		// 连接成功后实例化客户端。
		s.userProfileClient = userprofilev1.NewUserProfileServiceClient(s.userProfileConn)
		s.pointsClient = pointsv1.NewPointsServiceClient(s.pointsConn)
		s.sellerAppClient = sellershopv1.NewSellerApplicationServiceClient(s.sellerShopConn)
		s.sellerInternal = sellershopv1.NewInternalShopServiceClient(s.sellerShopConn)
		s.catalogSeller = catalogv1.NewSellerProductServiceClient(s.catalogConn)
		s.inventoryInternal = inventoryv1.NewInternalInventoryServiceClient(s.inventoryConn)
	})
	return s.initErr
}

// dialUpstream 按配置键创建下游 gRPC 连接。
func dialUpstream(ctx context.Context, cfgKey, defaultAddr string) (*grpc.ClientConn, error) {
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, cfgKey, defaultAddr).String())
	// 统一 5 秒超时，避免首连阻塞请求线程。
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// tsToString 将 protobuf 时间转换为 RFC3339 字符串，便于前端展示。
func tsToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format(time.RFC3339)
}

// uniqueStrings 对字符串切片去重并过滤空值，确保 degraded_fields 稳定。
func uniqueStrings(in []string) []string {
	m := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, item := range in {
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

// enumCode 将带前缀枚举值转换为前端可读业务码。
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
