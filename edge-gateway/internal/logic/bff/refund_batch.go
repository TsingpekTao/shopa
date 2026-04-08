package bff

import (
	"context"
	"strings"

	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
)

const buyerRefundScopeLabel = "当前按该商品所属整笔子单退款处理"

func (s *sBff) BuildBuyerRefundPreview(ctx context.Context, accessToken string, req *mev1.PreviewRefundReq) (*mev1.PreviewRefundRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if req == nil || strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.SubOrderNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "orderNo/subOrderNo are required")
	}
	buyerCtx := withOutgoingUserMetadata(ctx, verified.UserID)
	if err = s.ensureBuyerOwnsOrder(buyerCtx, req.OrderNo); err != nil {
		return nil, err
	}
	res, err := s.orderInternal.PreviewSubOrderRefund(buyerCtx, &orderv1.PreviewSubOrderRefundReq{
		OrderNo:    req.OrderNo,
		SubOrderNo: req.SubOrderNo,
		ItemNo:     req.ItemNo,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "preview refund failed")
	}
	return &mev1.PreviewRefundRes{
		Snapshot:   res.GetSnapshot(),
		ScopeCode:  "SUB_ORDER",
		ScopeLabel: buyerRefundScopeLabel,
	}, nil
}

func (s *sBff) ApplyBuyerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.ApplyRefundBatchReq) (*aftersalev1.ApplyRefundBatchRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if req == nil || len(req.GetTargets()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "targets are required")
	}
	buyerCtx := withOutgoingUserMetadata(ctx, verified.UserID)
	if err = s.ensureBuyerOwnsTargets(buyerCtx, req.GetTargets()); err != nil {
		return nil, err
	}
	return s.aftersaleBuyer.ApplyRefundBatch(buyerCtx, req)
}

func (s *sBff) ListBuyerRefundBatches(ctx context.Context, accessToken string, req *mev1.ListRefundBatchesReq) (*aftersalev1.ListMyRefundBatchesRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	return s.aftersaleBuyer.ListMyRefundBatches(withOutgoingUserMetadata(ctx, verified.UserID), &aftersalev1.ListMyRefundBatchesReq{
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (s *sBff) GetBuyerRefundBatchDetail(ctx context.Context, accessToken string, req *mev1.GetRefundBatchDetailReq) (*aftersalev1.GetMyRefundBatchDetailRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	return s.aftersaleBuyer.GetMyRefundBatchDetail(withOutgoingUserMetadata(ctx, verified.UserID), &aftersalev1.GetMyRefundBatchDetailReq{
		RefundBatchNo: req.RefundBatchNo,
	})
}

func (s *sBff) CancelBuyerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.CancelRefundBatchReq) (*aftersalev1.CancelRefundBatchRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	return s.aftersaleBuyer.CancelRefundBatch(withOutgoingUserMetadata(ctx, verified.UserID), req)
}

func (s *sBff) ListSellerRefundBatches(ctx context.Context, accessToken string, req *sellerv1.ListRefundBatchesReq) (*aftersalev1.ListShopRefundBatchesRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if err = s.ensureSellerOwnsShop(ctx, verified.UserID, req.ShopNo); err != nil {
		return nil, err
	}
	return s.aftersaleSeller.ListShopRefundBatches(ctx, &aftersalev1.ListShopRefundBatchesReq{
		ShopNo:     req.ShopNo,
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (s *sBff) GetSellerRefundBatchDetail(ctx context.Context, accessToken string, req *sellerv1.GetRefundBatchDetailReq) (*aftersalev1.GetShopRefundBatchDetailRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if err = s.ensureSellerOwnsShop(ctx, verified.UserID, req.ShopNo); err != nil {
		return nil, err
	}
	return s.aftersaleSeller.GetShopRefundBatchDetail(ctx, &aftersalev1.GetShopRefundBatchDetailReq{
		RefundBatchNo: req.RefundBatchNo,
		ShopNo:        req.ShopNo,
	})
}

func (s *sBff) ApproveSellerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.ApproveRefundBatchReq) (*aftersalev1.ApproveRefundBatchRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if err = s.ensureSellerOwnsShop(ctx, verified.UserID, req.GetShopNo()); err != nil {
		return nil, err
	}
	return s.aftersaleSeller.ApproveRefundBatch(ctx, req)
}

func (s *sBff) RejectSellerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.RejectRefundBatchReq) (*aftersalev1.RejectRefundBatchRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}
	if err = s.ensureSellerOwnsShop(ctx, verified.UserID, req.GetShopNo()); err != nil {
		return nil, err
	}
	return s.aftersaleSeller.RejectRefundBatch(ctx, req)
}

func (s *sBff) ensureBuyerOwnsTargets(ctx context.Context, targets []*aftersalev1.ApplyRefundTarget) error {
	seenOrders := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		if target == nil || strings.TrimSpace(target.GetOrderNo()) == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, "target orderNo is required")
		}
		if _, ok := seenOrders[target.GetOrderNo()]; ok {
			continue
		}
		if err := s.ensureBuyerOwnsOrder(ctx, target.GetOrderNo()); err != nil {
			return err
		}
		seenOrders[target.GetOrderNo()] = struct{}{}
	}
	return nil
}

func (s *sBff) ensureBuyerOwnsOrder(ctx context.Context, orderNo string) error {
	res, err := s.orderBuyer.GetMyOrderDetail(ctx, &orderv1.GetMyOrderDetailReq{OrderNo: orderNo})
	if err != nil {
		return gerror.Wrap(err, "verify buyer order ownership failed")
	}
	if res.GetOrder() == nil || strings.TrimSpace(res.GetOrder().GetOrderNo()) == "" {
		return gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	return nil
}

func (s *sBff) ensureSellerOwnsShop(ctx context.Context, userID uint64, shopNo string) error {
	if strings.TrimSpace(shopNo) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "shopNo is required")
	}
	res, err := s.sellerInternal.IsUserShopOwner(ctx, &sellershopv1.IsUserShopOwnerReq{
		UserId: userID,
		ShopNo: shopNo,
	})
	if err != nil {
		return gerror.Wrap(err, "check shop owner failed")
	}
	if !res.GetIsOwner() {
		return gerror.NewCode(gcode.CodeNotAuthorized, "user is not owner of this shop")
	}
	return nil
}
