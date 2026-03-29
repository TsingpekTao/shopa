package bff

import (
	"context"
	"strings"
	"sync"
	"time"

	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	adminv1 "github.com/TsingpekTao/shopa/edge-gateway/api/admin/v1"
	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/consts"
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

// sBff 璐熻矗缃戝叧 BFF 鑱氬悎锛?
// 1) 鑱氬悎澶氫釜涓嬫父鏈嶅姟锛岃緭鍑哄墠绔弸濂界粨鏋勩€?
// 2) 瀵规梺璺け璐ユ墽琛屸€滃瓧娈电骇闄嶇骇鈥濓紝鑰屼笉鏄暣椤靛け璐ャ€?
// 3) 瀵瑰垪琛ㄥ仛鎴柇鎺у埗锛岄伩鍏嶅搷搴斾綋杩囬噸銆?
type sBff struct {
	// once 淇濊瘉鎵€鏈変笅娓?gRPC 瀹㈡埛绔粎鍒濆鍖栦竴娆°€?
	once sync.Once

	// 浠ヤ笅鏄埌涓嬫父鏈嶅姟鐨勯暱杩炴帴銆?
	userProfileConn *grpc.ClientConn
	pointsConn      *grpc.ClientConn
	sellerShopConn  *grpc.ClientConn
	catalogConn     *grpc.ClientConn
	inventoryConn   *grpc.ClientConn

	// 浠ヤ笅鏄搴斾笅娓?RPC 瀹㈡埛绔€?
	userProfileClient userprofilev1.UserProfileServiceClient
	pointsClient      pointsv1.PointsServiceClient
	sellerAppClient   sellershopv1.SellerApplicationServiceClient
	sellerInternal    sellershopv1.InternalShopServiceClient
	catalogSeller     catalogv1.SellerProductServiceClient
	inventoryInternal inventoryv1.InternalInventoryServiceClient

	// initErr 璁板綍鍒濆鍖栧け璐ュ師鍥狅紝閬垮厤姣忔璇锋眰閮介噸澶嶅缓杩為噸璇曘€?
	initErr error
}

// New 鍒涘缓 BFF 鏈嶅姟瀹炰緥銆?
func New() *sBff {
	return &sBff{}
}

// init 鍚姩鏃舵敞鍐?BFF 瀹炵幇鍒?service 闂ㄩ潰銆?
func init() {
	service.RegisterBff(New())
}

// BuildMyOverview 鏋勫缓鈥滄垜鐨勬瑙堚€濊仛鍚堟暟鎹細
// - 涓昏韩浠戒俊鎭潵鑷?IAM 閴存潈缁撴灉銆?
// - 璧勬枡鏉ヨ嚜 user-profile銆?
// - 绉垎鏉ヨ嚜 points銆?
// - 浠讳竴璺け璐ユ椂杩斿洖 partial=true锛屽苟鏍囨敞 degraded_fields銆?
func (s *sBff) BuildMyOverview(ctx context.Context, accessToken string) (*mev1.GetOverviewRes, error) {
	// 绗竴姝ヤ覆琛岄壌鏉冿細韬唤鏃犳晥鏃剁洿鎺ュけ璐ワ紝涓嶇户缁笅娓歌皟鐢ㄣ€?
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 纭繚涓嬫父瀹㈡埛绔凡鍒濆鍖栥€?
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 鍏堝～鍏呭彲绔嬪嵆寰楀埌鐨勬牳蹇冭韩浠藉瓧娈点€?
	out := &mev1.GetOverviewRes{
		UserID:            verified.UserID,
		AccountStatusCode: verified.AccountStatusCode,
	}
	// 閫忎紶瑙掕壊缁欏墠绔紝鐢ㄤ簬鍏ュ彛鏄鹃殣鍜岃兘鍔涙帶鍒躲€?
	for _, role := range verified.Roles {
		out.Roles = append(out.Roles, mev1.RoleItem{
			RoleCode:      role.RoleCode,
			ScopeTypeCode: role.ScopeTypeCode,
			ScopeNo:       role.ScopeNo,
		})
	}

	var (
		// mu 淇濇姢骞跺彂鍐?out 涓?degradedFields銆?
		mu sync.Mutex
		// degradedFields 璁板綍鏈鑱氬悎涓彂鐢熼檷绾х殑瀛楁銆?
		degradedFields []string
	)
	// errgroup 骞跺彂璇锋眰澶氫釜涓嬫父锛岄檷浣庢暣浣?RT銆?
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 鎷夊彇鍩虹璧勬枡锛堟樀绉?澶村儚锛夛紱鍦板潃涓嶉渶瑕侊紝鍑忓皯涓嬫父璐熻浇銆?
		res, err := s.userProfileClient.GetProfileByUserId(egCtx, &userprofilev1.GetProfileByUserIdReq{
			UserId:           verified.UserID,
			IncludeAddresses: false,
		})
		if err != nil {
			// profile 澶辫触浠呰闄嶇骇锛屼笉涓柇鏁撮〉銆?
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
		// 鎷夊彇绉垎锛涘け璐ユ椂闄嶇骇涓洪粯璁?0锛堜繚鎸佸瓧娈靛瓨鍦級銆?
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

	// 姣忎釜 goroutine 宸插皢閿欒鈥滆浆璇戔€濅负闄嶇骇淇″彿锛岃繖閲屼笉鍐嶇‖澶辫触銆?
	_ = eg.Wait()
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

func (s *sBff) BuildAdminOverview(ctx context.Context, accessToken string) (*adminv1.GetOverviewRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "admin:me:overview") {
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}

	roles := make([]string, 0, len(verified.Roles))
	for _, role := range verified.Roles {
		if role.RoleCode == "" {
			continue
		}
		roles = append(roles, role.RoleCode)
	}

	return &adminv1.GetOverviewRes{
		UserID:            verified.UserID,
		AccountStatusCode: verified.AccountStatusCode,
		Roles:             uniqueStrings(roles),
		Permissions:       uniqueStrings(verified.Permissions),
	}, nil
}

func (s *sBff) BuildAdminDashboardOverview(ctx context.Context, accessToken string) (*adminv1.GetDashboardOverviewRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "dashboard:mall:view") {
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}

	metrics := []adminv1.DashboardMetric{
		{Key: "totalShops", Label: "Total Shops", Value: 0},
		{Key: "pendingMerchantReviews", Label: "Pending Merchant Reviews", Value: 0},
		{Key: "pendingProductReviews", Label: "Pending Product Reviews", Value: 0},
	}
	degradedFields := []string{"mall_metrics_unavailable"}

	return &adminv1.GetDashboardOverviewRes{
		Metrics:        metrics,
		Partial:        true,
		DegradedFields: degradedFields,
	}, nil
}

func (s *sBff) BuildAdminShopInsights(ctx context.Context, accessToken string, shopNo string) (*adminv1.GetShopInsightsRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "dashboard:shop:view") {
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}
	if strings.TrimSpace(shopNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	out := &adminv1.GetShopInsightsRes{
		ShopNo:         shopNo,
		ShopName:       "",
		ShopStatusCode: "",
	}
	shopRes, shopErr := s.sellerInternal.GetShopByNo(ctx, &sellershopv1.GetShopByNoReq{ShopNo: shopNo})
	if shopErr != nil || shopRes.GetShop() == nil {
		out.Partial = true
		out.DegradedFields = []string{"shop_basic_info_unavailable"}
		return out, nil
	}
	out.ShopName = shopRes.GetShop().GetShopName()
	out.ShopStatusCode = enumCode(shopRes.GetShop().GetStatus().String(), "SHOP_STATUS_")
	return out, nil
}

func (s *sBff) BuildAdminConversationList(ctx context.Context, accessToken string) (*adminv1.ListConversationsRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if !hasPermission(verified.Permissions, "cs:conversation:view") {
		return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
	}
	return &adminv1.ListConversationsRes{
		Items:          []adminv1.ConversationItem{},
		Partial:        true,
		DegradedFields: []string{"conversation_service_not_integrated"},
	}, nil
}

// BuildSellerWorkbench 鏋勫缓鍗栧宸ヤ綔鍙帮細
// - 搴楅摵鎽樿锛堟渶澶?10 鏉★級銆?
// - 鏈€杩戠敵璇凤紙鏈€澶?5 鏉★級銆?
// - 鍟嗗搧鐘舵€佽仛鍚堟憳瑕併€?
func (s *sBff) BuildSellerWorkbench(ctx context.Context, accessToken string) (*sellerv1.GetWorkbenchRes, error) {
	// 涓茶閴存潈锛岀‘淇?user_id 鍙俊銆?
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 纭繚 gRPC 瀹㈡埛绔彲鐢ㄣ€?
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 棰勭疆榛樿鍝嶅簲锛岄槻姝㈤檷绾ф椂瀛楁缂哄け銆?
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
		// mu 淇濇姢骞跺彂鍐欏搷搴斿璞°€?
		mu sync.Mutex
		// degradedFields 璁板綍鏈宸ヤ綔鍙拌仛鍚堜腑闄嶇骇椤广€?
		degradedFields []string
	)
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 鎷夊彇鈥滄垜鍚嶄笅搴楅摵鈥濄€?
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
		// total 淇濈暀鐪熷疄鎬绘暟锛屽墠绔彲鎹灞曠ず鈥滄洿澶氣€濄€?
		total := uint64(len(shops))
		// 宸ヤ綔鍙颁粎灞曠ず鍓?10 鏉★紝閬垮厤棣栧睆杩囪浇銆?
		limit := 10
		if len(shops) < limit {
			limit = len(shops)
		}
		// 棰勫垎閰嶅閲忥紝鍑忓皯 append 鎵╁寮€閿€銆?
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
		// 鎴柇鏍囪鍛婅瘔鍓嶇锛氳繖閲屽彧鏄憳瑕侊紝涓嶆槸鍏ㄩ噺銆?
		out.ShopsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 鎷夊彇鈥滄垜鐨勭敵璇封€濓紝鐢ㄤ簬宸ヤ綔鍙版渶杩戠敵璇峰崱鐗囥€?
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
		// total 浣跨敤鍒嗛〉鎺ュ彛杩斿洖鍊硷紝鑰岄潪褰撳墠椤甸暱搴︺€?
		total := uint64(res.GetTotal())
		// 宸ヤ綔鍙板浐瀹氬睍绀烘渶杩?5 鏉°€?
		limit := 5
		if len(apps) < limit {
			limit = len(apps)
		}
		list := make([]sellerv1.SellerApplicationSummary, 0, limit)
		for _, row := range apps[:limit] {
			list = append(list, sellerv1.SellerApplicationSummary{
				ApplicationNo:         row.GetApplicationNo(),
				ApplicationStatusCode: enumCode(row.GetStatus().String(), "APPLICATION_STATUS_"),
				// version 缁欏墠绔悗缁仛 expected_version 骞跺彂鎺у埗銆?
				Version:     int32(row.GetVersion()),
				SubmittedAt: tsToString(row.GetSubmittedAt()),
				UpdatedAt:   tsToString(row.GetUpdatedAt()),
			})
		}

		mu.Lock()
		out.LatestApplications = list
		out.ApplicationsTotal = total
		// 鎴柇鏍囪鐢ㄤ簬瑙﹀彂鈥滄煡鐪嬪叏閮ㄧ敵璇封€濆叆鍙ｃ€?
		out.ApplicationsTruncated = int(total) > limit
		mu.Unlock()
		return nil
	})

	eg.Go(func() error {
		// 骞跺彂鎷夊彇鍟嗗搧鐘舵€佹憳瑕侊紝閬垮厤涓茶鎷栨參宸ヤ綔鍙般€?
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

	// 鏃佽矾閿欒鍦?goroutine 鍐呭凡杞负闄嶇骇锛屾澶勬棤闇€纭け璐ャ€?
	_ = eg.Wait()
	out.Partial = len(degradedFields) > 0
	out.DegradedFields = uniqueStrings(degradedFields)
	return out, nil
}

// BuildSellerShopDashboard 鏋勫缓搴楅摵浠〃鐩橈細
// 1) 鍏堥壌鏉冨苟鏍￠獙 shop_no 褰掑睘銆?
// 2) 骞跺彂鎷夊彇搴楅摵淇℃伅銆佸晢鍝佹憳瑕併€佸簱瀛橀闄╂憳瑕併€?
// 3) 鏆傜己鑳藉姏鐢?degraded_fields 鏄庣‘鏍囨敞銆?
func (s *sBff) BuildSellerShopDashboard(ctx context.Context, accessToken, shopNo string) (*sellerv1.GetShopDashboardRes, error) {
	// shopNo 鏄簵閾虹淮搴︿富閿紝涓嶅厑璁镐负绌恒€?
	if strings.TrimSpace(shopNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	// 缁熶竴閴存潈銆?
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	// 鍒濆鍖栦笅娓稿鎴风銆?
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 褰掑睘鏍￠獙锛氫粎搴楅摵 owner 鍙闂 dashboard銆?
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

	// 鍏堢粰榛樿缁撴瀯锛岀‘淇濋檷绾ф椂鍝嶅簲缁撴瀯绋冲畾銆?
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
		// mu 淇濇姢骞跺彂鍐欏叡浜粨鏋溿€?
		mu sync.Mutex
		// degradedFields 姹囨€绘湰娆¤姹傜殑闄嶇骇椤广€?
		degradedFields []string
	)
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// 鏌ヨ搴楅摵鍩烘湰淇℃伅锛堝悕绉般€佺姸鎬侊級銆?
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
		// 褰撳墠 catalog-svc 鏆備笉鏀寔鎸?shop_no 绮剧‘杩囨护銆?
		// 杩欓噷鍥為€€鍒板崠瀹剁骇鎽樿锛屽苟鏄惧紡鏍囨敞涓洪檷绾у瓧娈点€?
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
		// 褰撳墠 inventory-svc 鏆傛棤 shop 缁村害搴撳瓨椋庨櫓鑱氬悎鎺ュ彛銆?
		// 鍥犳淇濈暀榛樿鍊煎苟鏍囨敞闄嶇骇锛岃€屼笉鏄鏁撮〉澶辫触銆?
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

// loadSellerProductSummary 鑱氬悎鍗栧鍟嗗搧鐘舵€佺粺璁★細
// 骞跺彂璋冪敤 ListMyProducts(total-only) 鎷夊彇鍚勭姸鎬佽鏁帮紝鍐嶅洖濉眹鎬荤粨鏋勩€?
func (s *sBff) loadSellerProductSummary(ctx context.Context) (sellerv1.SellerProductSummary, error) {
	// statusCase 鐢ㄤ簬鎶借薄鈥滅粺璁￠敭鈥濅笌鈥滀笅娓哥姸鎬佹灇涓锯€濈殑鏄犲皠鍏崇郴銆?
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
		// mu 淇濇姢 counter 骞跺彂鍐欍€?
		mu sync.Mutex
		// counter 瀛樻斁鍚勭姸鎬佸晢鍝佹暟銆?
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
		// 閲嶆柊缁戝畾寰幆鍙橀噺锛岄伩鍏?goroutine 闂寘鎹曡幏鍚屼竴鍦板潃銆?
		item := item
		eg.Go(func() error {
			// 鍙煡鎬绘暟锛孭ageSize=1 鍗冲彲锛屽噺灏戜笅娓歌繑鍥炶礋杞姐€?
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

	// 浠讳竴鐘舵€佺粺璁″け璐ユ椂锛屼氦鐢变笂灞傚喅瀹氭槸鍚﹂檷绾с€?
	if err := eg.Wait(); err != nil {
		return sellerv1.SellerProductSummary{}, err
	}

	// 鍐嶆媺涓€娆′笉甯︾姸鎬佽繃婊ょ殑鎬绘暟銆?
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

// ensureClients 鎯版€у垵濮嬪寲鎵€鏈変笅娓?gRPC 瀹㈡埛绔€?
func (s *sBff) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		var err error

		// 鐢ㄦ埛璧勬枡鏈嶅姟銆?
		s.userProfileConn, err = dialUpstream(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002")
		if err != nil {
			s.initErr = err
			return
		}
		// 绉垎鏈嶅姟銆?
		s.pointsConn, err = dialUpstream(ctx, "upstream.pointsGrpc", "127.0.0.1:8012")
		if err != nil {
			s.initErr = err
			return
		}
		// 鍗栧搴楅摵鏈嶅姟銆?
		s.sellerShopConn, err = dialUpstream(ctx, "upstream.sellerShopGrpc", "127.0.0.1:50051")
		if err != nil {
			s.initErr = err
			return
		}
		// 鍟嗗搧鏈嶅姟銆?
		s.catalogConn, err = dialUpstream(ctx, "upstream.catalogGrpc", "127.0.0.1:9004")
		if err != nil {
			s.initErr = err
			return
		}
		// 搴撳瓨鏈嶅姟銆?
		s.inventoryConn, err = dialUpstream(ctx, "upstream.inventoryGrpc", "127.0.0.1:9005")
		if err != nil {
			s.initErr = err
			return
		}

		// 杩炴帴鎴愬姛鍚庡疄渚嬪寲瀹㈡埛绔€?
		s.userProfileClient = userprofilev1.NewUserProfileServiceClient(s.userProfileConn)
		s.pointsClient = pointsv1.NewPointsServiceClient(s.pointsConn)
		s.sellerAppClient = sellershopv1.NewSellerApplicationServiceClient(s.sellerShopConn)
		s.sellerInternal = sellershopv1.NewInternalShopServiceClient(s.sellerShopConn)
		s.catalogSeller = catalogv1.NewSellerProductServiceClient(s.catalogConn)
		s.inventoryInternal = inventoryv1.NewInternalInventoryServiceClient(s.inventoryConn)
	})
	return s.initErr
}

// dialUpstream 鎸夐厤缃敭鍒涘缓涓嬫父 gRPC 杩炴帴銆?
func dialUpstream(ctx context.Context, cfgKey, defaultAddr string) (*grpc.ClientConn, error) {
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, cfgKey, defaultAddr).String())
	// 缁熶竴 5 绉掕秴鏃讹紝閬垮厤棣栬繛闃诲璇锋眰绾跨▼銆?
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// tsToString 灏?protobuf 鏃堕棿杞崲涓?RFC3339 瀛楃涓诧紝渚夸簬鍓嶇灞曠ず銆?
func tsToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format(time.RFC3339)
}

// uniqueStrings 瀵瑰瓧绗︿覆鍒囩墖鍘婚噸骞惰繃婊ょ┖鍊硷紝纭繚 degraded_fields 绋冲畾銆?
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

func hasPermission(granted []string, required string) bool {
	required = strings.TrimSpace(required)
	if required == "" {
		return true
	}
	for _, item := range granted {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == required || item == "*:*:*" {
			return true
		}
	}
	return false
}

// enumCode converts enum names with prefix into plain business codes.
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
