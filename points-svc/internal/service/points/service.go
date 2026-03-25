package points

import (
	"context"
	"fmt"
	"sync"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/dao"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Service struct{}

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

func (s *Service) InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (*v1.InitPointsAccountIfAbsentRes, error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	created, balance, _, err := s.initAccountIfAbsent(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &v1.InitPointsAccountIfAbsentRes{
		Created: created,
		Balance: balance,
	}, nil
}

func (s *Service) GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (*v1.GetPointsByUserIdRes, error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	account, err := s.getAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	if account == nil {
		return &v1.GetPointsByUserIdRes{
			UserId:  req.GetUserId(),
			Balance: 0,
			Status:  0,
		}, nil
	}
	return &v1.GetPointsByUserIdRes{
		UserId:  account.UserId,
		Balance: toUint64(account.Balance),
		Status:  uint32(account.Status),
	}, nil
}

// InitAccountFromRegisterEvent 给 MQ 消费者调用，避免写重复 RPC 适配代码。
func (s *Service) InitAccountFromRegisterEvent(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	_, _, _, err := s.initAccountIfAbsent(ctx, userID)
	return err
}

func (s *Service) initAccountIfAbsent(ctx context.Context, userID uint64) (created bool, balance uint64, status uint32, err error) {
	r, err := dao.PointsAccount.Ctx(ctx).Data(do.PointsAccount{
		UserId:  userID,
		Balance: 0,
		Status:  1,
	}).InsertIgnore()
	if err != nil {
		return false, 0, 0, gerror.Wrap(err, "insert points_account failed")
	}
	if r != nil {
		if rows, rowsErr := r.RowsAffected(); rowsErr == nil && rows > 0 {
			created = true
		}
	}

	if created {
		_, err = dao.PointsLedger.Ctx(ctx).Data(do.PointsLedger{
			UserId:       userID,
			BizType:      "REGISTER_INIT",
			BizId:        fmt.Sprintf("register:%d", userID),
			Delta:        0,
			BalanceAfter: 0,
		}).InsertIgnore()
		if err != nil {
			return false, 0, 0, gerror.Wrap(err, "insert points_ledger failed")
		}
	}

	account, err := s.getAccount(ctx, userID)
	if err != nil {
		return false, 0, 0, err
	}
	if account == nil {
		return created, 0, 0, nil
	}
	return created, toUint64(account.Balance), uint32(account.Status), nil
}

func (s *Service) getAccount(ctx context.Context, userID uint64) (*entity.PointsAccount, error) {
	var account entity.PointsAccount
	cols := dao.PointsAccount.Columns()
	if err := dao.PointsAccount.Ctx(ctx).Where(cols.UserId, userID).Scan(&account); err != nil {
		return nil, gerror.Wrap(err, "query points_account failed")
	}
	if account.UserId == 0 {
		return nil, nil
	}
	return &account, nil
}

func toUint64(v int64) uint64 {
	if v <= 0 {
		return 0
	}
	return uint64(v)
}
