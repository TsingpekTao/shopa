package auth

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
)

// HasShopRole 判断用户是否已经拥有指定店铺作用域下的卖家角色。
func (s *Service) HasShopRole(ctx context.Context, req *v1.HasShopRoleReq) (*v1.HasShopRoleRes, error) {
	if req.GetUserId() == 0 || req.GetScopeId() == 0 {
		return nil, errs.New(errs.CodeInvalidParam, "user_id and scope_id are required")
	}

	roleCols := dao.IamUserRole.Columns()
	count, err := dao.IamUserRole.Ctx(ctx).
		Where(roleCols.UserId, req.GetUserId()).
		Where(roleCols.RoleCode, consts.RoleCodeSeller).
		Where(roleCols.ScopeType, consts.ScopeTypeShop).
		Where(roleCols.ScopeId, req.GetScopeId()).
		Where(roleCols.Status, consts.StatusActive).
		Count()
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.HasShopRoleRes{HasRole: count > 0}, nil
}

// AssignShopSellerRole 为用户授予店铺卖家角色；已存在时按幂等成功返回。
func (s *Service) AssignShopSellerRole(ctx context.Context, req *v1.AssignShopSellerRoleReq) (*v1.AssignShopSellerRoleRes, error) {
	if req.GetUserId() == 0 || req.GetScopeId() == 0 {
		return nil, errs.New(errs.CodeInvalidParam, "user_id and scope_id are required")
	}

	auth, err := s.findAuthByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if auth == nil {
		return nil, errs.New(errs.CodeInvalidParam, "user not found")
	}

	err = dao.IamUserRole.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		role, err := s.findShopSellerRoleTx(ctx, tx, req.GetUserId(), req.GetScopeId())
		if err != nil {
			return err
		}

		if role != nil {
			if role.Status == consts.StatusActive {
				return nil
			}

			_, err = tx.Model(dao.IamUserRole.Table()).
				Where(dao.IamUserRole.Columns().Id, role.Id).
				Data(do.IamUserRole{Status: consts.StatusActive}).
				Update()
			return err
		}

		_, err = tx.Model(dao.IamUserRole.Table()).
			Data(do.IamUserRole{
				UserId:    req.GetUserId(),
				RoleCode:  consts.RoleCodeSeller,
				ScopeType: consts.ScopeTypeShop,
				ScopeId:   req.GetScopeId(),
				Status:    consts.StatusActive,
			}).
			Insert()
		return err
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.AssignShopSellerRoleRes{Assigned: true}, nil
}

// RevokeShopSellerRoleAndBumpToken 回收指定店铺卖家角色，并提升 token_version 使旧令牌失效。
func (s *Service) RevokeShopSellerRoleAndBumpToken(ctx context.Context, req *v1.RevokeShopSellerRoleAndBumpTokenReq) (*v1.RevokeShopSellerRoleAndBumpTokenRes, error) {
	if req.GetUserId() == 0 || req.GetScopeId() == 0 {
		return nil, errs.New(errs.CodeInvalidParam, "user_id and scope_id are required")
	}

	var (
		revoked      bool
		tokenVersion uint32
	)

	err := dao.IamUserRole.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		auth, err := s.findAuthByUserIDTx(ctx, tx, req.GetUserId())
		if err != nil {
			return err
		}
		if auth == nil {
			return errs.New(errs.CodeInvalidParam, "user not found")
		}

		role, err := s.findShopSellerRoleTx(ctx, tx, req.GetUserId(), req.GetScopeId())
		if err != nil {
			return err
		}
		if role == nil || role.Status != consts.StatusActive {
			tokenVersion = uint32(auth.TokenVersion)
			return nil
		}

		_, err = tx.Model(dao.IamUserRole.Table()).
			Where(dao.IamUserRole.Columns().Id, role.Id).
			Data(do.IamUserRole{Status: consts.StatusDisabled}).
			Update()
		if err != nil {
			return err
		}

		nextTokenVersion := auth.TokenVersion + 1
		_, err = tx.Model(dao.IamUserAuth.Table()).
			Where(dao.IamUserAuth.Columns().UserId, req.GetUserId()).
			Data(do.IamUserAuth{TokenVersion: nextTokenVersion}).
			Update()
		if err != nil {
			return err
		}

		revoked = true
		tokenVersion = uint32(nextTokenVersion)
		return nil
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.RevokeShopSellerRoleAndBumpTokenRes{
		Revoked:      revoked,
		TokenVersion: tokenVersion,
	}, nil
}

// findShopSellerRoleTx 在事务内读取指定店铺作用域下的卖家角色。
func (s *Service) findShopSellerRoleTx(ctx context.Context, tx gdb.TX, userID, scopeID uint64) (*entity.IamUserRole, error) {
	roleCols := dao.IamUserRole.Columns()
	var role entity.IamUserRole
	err := tx.Model(dao.IamUserRole.Table()).
		Where(roleCols.UserId, userID).
		Where(roleCols.RoleCode, consts.RoleCodeSeller).
		Where(roleCols.ScopeType, consts.ScopeTypeShop).
		Where(roleCols.ScopeId, scopeID).
		Scan(&role)
	if err != nil {
		return nil, err
	}
	if role.Id == 0 {
		return nil, nil
	}
	return &role, nil
}

// findAuthByUserIDTx 在事务内读取账户主记录，保证角色回收与 token_version 更新原子一致。
func (s *Service) findAuthByUserIDTx(ctx context.Context, tx gdb.TX, userID uint64) (*entity.IamUserAuth, error) {
	authCols := dao.IamUserAuth.Columns()
	var auth entity.IamUserAuth
	err := tx.Model(dao.IamUserAuth.Table()).
		Where(authCols.UserId, userID).
		Scan(&auth)
	if err != nil {
		return nil, err
	}
	if auth.UserId == 0 {
		return nil, nil
	}
	return &auth, nil
}
