package overview

import (
	"context"
	"strings"
	"sync"
	"time"

	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	pointsv1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	once sync.Once

	iamConn         *grpc.ClientConn
	userProfileConn *grpc.ClientConn
	pointsConn      *grpc.ClientConn

	iamClient         iamv1.InternalServiceClient
	userProfileClient userprofilev1.UserProfileServiceClient
	pointsClient      pointsv1.PointsServiceClient

	initErr error
}

var (
	serviceOnce sync.Once
	serviceInst *Service
)

func New() *Service {
	serviceOnce.Do(func() {
		serviceInst = &Service{}
	})
	return serviceInst
}

type OverviewResult struct {
	UserID        uint64
	AccountStatus int32
	Roles         []RoleItem
	DisplayName   string
	AvatarURL     string
	Points        uint64
}

type RoleItem struct {
	RoleCode  int32
	ScopeType int32
	ScopeID   uint64
}

func (s *Service) BuildMyOverview(ctx context.Context, accessToken string) (*OverviewResult, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "missing access token")
	}
	if err := s.ensureClients(ctx); err != nil {
		return nil, err
	}

	verifyRes, err := s.iamClient.VerifyAccessToken(ctx, &iamv1.VerifyAccessTokenReq{AccessToken: accessToken})
	if err != nil {
		return nil, err
	}
	if !verifyRes.GetValid() || verifyRes.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid access token")
	}

	userID := verifyRes.GetUserId()
	authRes, err := s.iamClient.GetAuthUserById(ctx, &iamv1.GetAuthUserByIdReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	profileRes, err := s.userProfileClient.GetProfileByUserId(ctx, &userprofilev1.GetProfileByUserIdReq{
		UserId:           userID,
		IncludeAddresses: false,
	})
	if err != nil {
		return nil, err
	}
	pointsRes, err := s.pointsClient.GetPointsByUserId(ctx, &pointsv1.GetPointsByUserIdReq{UserId: userID})
	if err != nil {
		return nil, err
	}

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

func (s *Service) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		iamAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.iamGrpc", "127.0.0.1:9001").String())
		userProfileAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002").String())
		pointsAddr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.pointsGrpc", "127.0.0.1:8012").String())

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
