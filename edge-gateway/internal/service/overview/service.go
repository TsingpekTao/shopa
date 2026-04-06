package overview

import (
	"context"
	"strings"
	"sync"
	"time"

	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	pointsv1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service 是一个较早期的概览聚合实现：
// - 串行调用 IAM、资料、积分服务。
// - 与 logic/bff 不同，这里不做字段级降级，任一步失败即返回错误。
type Service struct {
	// once 用于惰性初始化全部 gRPC 连接。
	once sync.Once

	// 下游连接。
	iamConn         *grpc.ClientConn
	userProfileConn *grpc.ClientConn
	pointsConn      *grpc.ClientConn

	// 下游客户端。
	iamClient         iamv1.InternalServiceClient
	userProfileClient userprofilev1.UserProfileServiceClient
	pointsClient      pointsv1.PointsServiceClient

	// initErr 记录首次初始化失败原因。
	initErr error
}

var (
	// serviceOnce + serviceInst 实现进程内单例。
	serviceOnce sync.Once
	serviceInst *Service
)

// New 返回 Service 单例。
func New() *Service {
	serviceOnce.Do(func() {
		serviceInst = &Service{}
	})
	return serviceInst
}

// OverviewResult 是概览聚合输出结构。
type OverviewResult struct {
	UserID        uint64
	AccountStatus int32
	Roles         []RoleItem
	DisplayName   string
	AvatarURL     string
	Points        uint64
}

// RoleItem 描述用户角色及其作用域。
type RoleItem struct {
	RoleCode  int32
	ScopeType int32
	ScopeID   uint64
}

// BuildMyOverview 构建“我的概览”：
// 1) 校验 access token。
// 2) 拉 IAM 用户信息。
// 3) 拉资料和积分。
// 4) 组装统一结果。
func (s *Service) BuildMyOverview(ctx context.Context, accessToken string) (*OverviewResult, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "missing access token")
	}
	if err := s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 先做 token 合法性校验。
	verifyRes, err := s.iamClient.VerifyAccessToken(ctx, &iamv1.VerifyAccessTokenReq{AccessToken: accessToken})
	if err != nil {
		return nil, err
	}
	if !verifyRes.GetValid() || verifyRes.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid access token")
	}

	userID := verifyRes.GetUserId()

	// 拉 IAM 账户详情（状态、角色）。
	authRes, err := s.iamClient.GetAuthUserById(ctx, &iamv1.GetAuthUserByIdReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	// 拉用户资料。
	profileRes, err := s.userProfileClient.GetProfileByUserId(ctx, &userprofilev1.GetProfileByUserIdReq{
		UserId:           userID,
		IncludeAddresses: false,
	})
	if err != nil {
		return nil, err
	}
	// 拉积分。
	pointsRes, err := s.pointsClient.GetPointsByUserId(ctx, &pointsv1.GetPointsByUserIdReq{UserId: userID})
	if err != nil {
		return nil, err
	}

	// 先填充稳定字段。
	out := &OverviewResult{
		UserID: userID,
		Points: pointsRes.GetBalance(),
	}
	if authRes.GetUser() != nil {
		out.AccountStatus = int32(authRes.GetUser().GetAccountStatus())
		for _, role := range authRes.GetUser().GetRoles() {
			out.Roles = append(out.Roles, RoleItem{
				RoleCode:  int32(role.GetRoleCode()),
				ScopeType: int32(role.GetScopeType()),
				ScopeID:   role.GetScopeId(),
			})
		}
	}
	if profileRes.GetProfile() != nil {
		out.DisplayName = profileRes.GetProfile().GetDisplayName()
		if profileRes.GetProfile().GetAvatar() != nil {
			out.AvatarURL = profileRes.GetProfile().GetAvatar().GetUrl()
		}
	}
	return out, nil
}

// ensureClients 惰性初始化 IAM/资料/积分 gRPC 客户端。
func (s *Service) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		iamAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.iamGrpc", "127.0.0.1:9001").String())
		userProfileAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002").String())
		pointsAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.pointsGrpc", "127.0.0.1:9012").String())

		// 连接初始化统一使用 5 秒超时。
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.iamConn, s.initErr = grpc.DialContext(timeoutCtx, iamAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			return
		}
		s.userProfileConn, s.initErr = grpc.DialContext(timeoutCtx, userProfileAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			return
		}
		s.pointsConn, s.initErr = grpc.DialContext(timeoutCtx, pointsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			return
		}

		s.iamClient = iamv1.NewInternalServiceClient(s.iamConn)
		s.userProfileClient = userprofilev1.NewUserProfileServiceClient(s.userProfileConn)
		s.pointsClient = pointsv1.NewPointsServiceClient(s.pointsConn)
	})
	return s.initErr
}
