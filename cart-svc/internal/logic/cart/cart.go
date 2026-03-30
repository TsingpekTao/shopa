package cart

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	"github.com/TsingpekTao/shopa/cart-svc/internal/consts"
	"github.com/TsingpekTao/shopa/cart-svc/internal/dao"
	"github.com/TsingpekTao/shopa/cart-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/cart-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/cart-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// checkoutTokenConsumeLua 使用 Redis Lua 保证 checkout token 的消费原子性。
// 返回值约定：
// 0: token 不存在或已过期；
// 1: 校验通过并成功消费；
// 2: token 与 user_id 不匹配（不消费 token）。
const checkoutTokenConsumeLua = `
local owner = redis.call('GET', KEYS[2])
if not owner then
  return {0}
end
if owner ~= ARGV[1] then
  return {2}
end
local payload = redis.call('GET', KEYS[1])
if not payload then
  return {0}
end
redis.call('DEL', KEYS[1])
redis.call('DEL', KEYS[2])
return {1, payload}
`

type sCart struct{}

type redisCartItem struct {
	UserID            uint64 `json:"user_id"`
	SkuNo             string `json:"sku_no"`
	SpuNo             string `json:"spu_no"`
	ShopNo            string `json:"shop_no"`
	Qty               uint32 `json:"qty"`
	Checked           bool   `json:"checked"`
	Status            int32  `json:"status"`
	InvalidReasonCode string `json:"invalid_reason_code"`
	SpuTitle          string `json:"spu_title"`
	SkuName           string `json:"sku_name"`
	SkuImageAssetID   uint64 `json:"sku_image_asset_id"`
	SalePrice         uint64 `json:"sale_price"`
	MarketPrice       uint64 `json:"market_price"`
	SaleAttrsJSON     string `json:"sale_attrs_json"`
	CreatedAtUnix     int64  `json:"created_at_unix"`
	UpdatedAtUnix     int64  `json:"updated_at_unix"`
}

// New 创建购物车服务实现。
func New() *sCart {
	return &sCart{}
}

func init() {
	svc := New()
	service.RegisterCart(svc)
	svc.startBackupSyncWorker()
}

// AddItem 买家加购，购物车以 Redis Hash 为主存，并刷新 30 天滑动 TTL。
func (s *sCart) AddItem(ctx context.Context, req *v1.AddItemReq) (*v1.AddItemRes, error) {
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	if req.GetQty() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "qty must be greater than 0")
	}

	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Unix()
	skuNo := strings.TrimSpace(req.GetSkuNo())
	item, exists := items[skuNo]
	if !exists {
		item = &redisCartItem{
			UserID:            userID,
			SkuNo:             skuNo,
			SpuNo:             strings.TrimSpace(req.GetSpuNo()),
			ShopNo:            strings.TrimSpace(req.GetShopNo()),
			Qty:               0,
			Checked:           req.GetChecked(),
			Status:            int32(v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE),
			CreatedAtUnix:     now,
			UpdatedAtUnix:     now,
			SaleAttrsJSON:     "[]",
			InvalidReasonCode: "",
		}
	}
	item.Qty += req.GetQty()
	item.UpdatedAtUnix = now
	if !exists {
		item.Checked = req.GetChecked()
	}
	if item.Status == int32(v1.CartItemStatus_CART_ITEM_STATUS_UNSPECIFIED) {
		item.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE)
	}
	items[skuNo] = item

	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	_ = s.markDirtyUser(ctx, userID)

	return &v1.AddItemRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// UpdateItemQty 更新购物车中某个 SKU 的购买数量。
func (s *sCart) UpdateItemQty(ctx context.Context, req *v1.UpdateItemQtyReq) (*v1.UpdateItemQtyRes, error) {
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	if req.GetQty() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "qty must be greater than 0")
	}

	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	item, ok := items[strings.TrimSpace(req.GetSkuNo())]
	if !ok {
		return nil, gerror.NewCode(gcode.CodeNotFound, "cart item not found")
	}
	item.Qty = req.GetQty()
	item.UpdatedAtUnix = time.Now().UTC().Unix()
	items[item.SkuNo] = item

	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	_ = s.markDirtyUser(ctx, userID)

	return &v1.UpdateItemQtyRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// ToggleItemChecked 设置单个商品勾选状态。
func (s *sCart) ToggleItemChecked(ctx context.Context, req *v1.ToggleItemCheckedReq) (*v1.ToggleItemCheckedRes, error) {
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	item, ok := items[strings.TrimSpace(req.GetSkuNo())]
	if !ok {
		return nil, gerror.NewCode(gcode.CodeNotFound, "cart item not found")
	}
	item.Checked = req.GetChecked()
	item.UpdatedAtUnix = time.Now().UTC().Unix()
	items[item.SkuNo] = item

	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	_ = s.markDirtyUser(ctx, userID)

	return &v1.ToggleItemCheckedRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// BatchToggleItems 批量设置勾选状态，常用于“全选/取消全选”。
func (s *sCart) BatchToggleItems(ctx context.Context, req *v1.BatchToggleItemsReq) (*v1.BatchToggleItemsRes, error) {
	if len(req.GetSkuNos()) == 0 {
		return &v1.BatchToggleItemsRes{Affected: 0, Summary: &v1.CartSummary{}}, nil
	}
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	affected := uint32(0)
	fields := make(map[string]any, len(req.GetSkuNos()))
	now := time.Now().UTC().Unix()
	for _, skuNo := range req.GetSkuNos() {
		skuNo = strings.TrimSpace(skuNo)
		item, ok := items[skuNo]
		if !ok {
			continue
		}
		item.Checked = req.GetChecked()
		item.UpdatedAtUnix = now
		items[skuNo] = item
		encoded, marshalErr := json.Marshal(item)
		if marshalErr != nil {
			return nil, gerror.Wrapf(marshalErr, "marshal cart item failed, sku_no=%s", skuNo)
		}
		fields[skuNo] = string(encoded)
		affected++
	}
	if len(fields) > 0 {
		if _, err = g.Redis().HSet(ctx, s.cartKey(userID), fields); err != nil {
			return nil, gerror.Wrap(err, "batch toggle cart items failed")
		}
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		_ = s.markDirtyUser(ctx, userID)
	}

	return &v1.BatchToggleItemsRes{
		Affected: affected,
		Summary:  summarize(items),
	}, nil
}

// RemoveItems 从购物车中删除指定商品。
func (s *sCart) RemoveItems(ctx context.Context, req *v1.RemoveItemsReq) (*v1.RemoveItemsRes, error) {
	if len(req.GetSkuNos()) == 0 {
		return &v1.RemoveItemsRes{Removed: 0, Summary: &v1.CartSummary{}}, nil
	}
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	keysToDelete := make([]string, 0, len(req.GetSkuNos()))
	removed := uint32(0)
	for _, skuNo := range req.GetSkuNos() {
		skuNo = strings.TrimSpace(skuNo)
		if skuNo == "" {
			continue
		}
		if _, ok := items[skuNo]; ok {
			delete(items, skuNo)
			keysToDelete = append(keysToDelete, skuNo)
			removed++
		}
	}
	if len(keysToDelete) > 0 {
		if _, err = g.Redis().HDel(ctx, s.cartKey(userID), keysToDelete...); err != nil {
			return nil, gerror.Wrap(err, "remove cart items failed")
		}
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		_ = s.markDirtyUser(ctx, userID)
	}

	return &v1.RemoveItemsRes{
		Removed: removed,
		Summary: summarize(items),
	}, nil
}

// ClearInvalidItems 清空失效商品。
func (s *sCart) ClearInvalidItems(ctx context.Context, _ *v1.ClearInvalidItemsReq) (*v1.ClearInvalidItemsRes, error) {
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	keysToDelete := make([]string, 0)
	removed := uint32(0)
	for skuNo, item := range items {
		if v1.CartItemStatus(item.Status) == v1.CartItemStatus_CART_ITEM_STATUS_INVALID {
			keysToDelete = append(keysToDelete, skuNo)
			delete(items, skuNo)
			removed++
		}
	}
	if len(keysToDelete) > 0 {
		if _, err = g.Redis().HDel(ctx, s.cartKey(userID), keysToDelete...); err != nil {
			return nil, gerror.Wrap(err, "clear invalid cart items failed")
		}
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		_ = s.markDirtyUser(ctx, userID)
	}

	return &v1.ClearInvalidItemsRes{
		Removed: removed,
		Summary: summarize(items),
	}, nil
}

// GetMyCart 获取当前用户购物车；Redis miss 时支持从 MySQL 备份冷恢复。
func (s *sCart) GetMyCart(ctx context.Context, req *v1.GetMyCartReq) (*v1.GetMyCartRes, error) {
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	items, loadedFromBackup, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	out := make([]*v1.CartItem, 0, len(items))
	for _, item := range items {
		if req.GetOnlyChecked() && !item.Checked {
			continue
		}
		out = append(out, toProtoCartItem(item))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].UpdatedAt == nil {
			return false
		}
		if out[j].UpdatedAt == nil {
			return true
		}
		return out[i].UpdatedAt.AsTime().After(out[j].UpdatedAt.AsTime())
	})

	return &v1.GetMyCartRes{
		Items:            out,
		Summary:          summarize(items),
		LoadedFromBackup: loadedFromBackup,
	}, nil
}

// PrepareCheckout 生成结算快照 token（默认 5 分钟有效），订单侧必须使用该快照下单。
func (s *sCart) PrepareCheckout(ctx context.Context, req *v1.PrepareCheckoutReq) (*v1.PrepareCheckoutRes, error) {
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address_id is required")
	}

	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	selected := s.selectCheckoutItems(items, req)
	if len(selected) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "no cart item selected for checkout")
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].SkuNo < selected[j].SkuNo
	})

	snapshotItems := make([]*v1.CheckoutSnapshotItem, 0, len(selected))
	goodsAmount := uint64(0)
	for _, item := range selected {
		snapshotItems = append(snapshotItems, &v1.CheckoutSnapshotItem{
			SkuNo:           item.SkuNo,
			SpuNo:           item.SpuNo,
			ShopNo:          item.ShopNo,
			Qty:             item.Qty,
			SettlePrice:     item.SalePrice,
			MarketPrice:     item.MarketPrice,
			SpuTitle:        item.SpuTitle,
			SkuName:         item.SkuName,
			SkuImageAssetId: item.SkuImageAssetID,
			SaleAttrsJson:   item.SaleAttrsJSON,
		})
		goodsAmount += item.SalePrice * uint64(item.Qty)
	}

	freightAmount := uint64(0)
	payableAmount := goodsAmount + freightAmount
	now := time.Now().UTC()
	expireAt := now.Add(time.Duration(consts.CheckoutTokenTTLSeconds) * time.Second)
	token := s.newCheckoutToken()

	snapshot := &v1.CheckoutSnapshot{
		CheckoutToken: token,
		UserId:        userID,
		Items:         snapshotItems,
		GoodsAmount:   goodsAmount,
		FreightAmount: freightAmount,
		PayableAmount: payableAmount,
		CreatedAt:     timestamppb.New(now),
		ExpireAt:      timestamppb.New(expireAt),
	}
	snapshot.SnapshotDigest = digestSnapshot(snapshot)

	payload, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(snapshot)
	if err != nil {
		return nil, gerror.Wrap(err, "marshal checkout snapshot failed")
	}
	if err = g.Redis().SetEX(ctx, s.checkoutTokenPayloadKey(token), string(payload), consts.CheckoutTokenTTLSeconds); err != nil {
		return nil, gerror.Wrap(err, "save checkout token payload failed")
	}
	if err = g.Redis().SetEX(ctx, s.checkoutTokenOwnerKey(token), strconv.FormatUint(userID, 10), consts.CheckoutTokenTTLSeconds); err != nil {
		return nil, gerror.Wrap(err, "save checkout token owner failed")
	}

	return &v1.PrepareCheckoutRes{Snapshot: snapshot}, nil
}

// ConsumeCheckoutToken 通过 Lua 原子脚本消费 token，避免并发重放下单。
func (s *sCart) ConsumeCheckoutToken(ctx context.Context, req *v1.ConsumeCheckoutTokenReq) (*v1.ConsumeCheckoutTokenRes, error) {
	token := strings.TrimSpace(req.GetCheckoutToken())
	if token == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout_token is required")
	}
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}

	resultVar, err := g.Redis().Eval(
		ctx,
		checkoutTokenConsumeLua,
		2,
		[]string{s.checkoutTokenPayloadKey(token), s.checkoutTokenOwnerKey(token)},
		[]any{strconv.FormatUint(req.GetUserId(), 10)},
	)
	if err != nil {
		return nil, gerror.Wrap(err, "consume checkout token by lua failed")
	}
	result := gconv.SliceAny(resultVar.Val())
	if len(result) == 0 {
		return &v1.ConsumeCheckoutTokenRes{
			TokenStatus: v1.CheckoutTokenStatus_CHECKOUT_TOKEN_STATUS_NOT_FOUND_OR_EXPIRED,
			ErrorCode:   "TOKEN_NOT_FOUND",
		}, nil
	}

	code := gconv.Int(result[0])
	switch code {
	case 0:
		return &v1.ConsumeCheckoutTokenRes{
			TokenStatus: v1.CheckoutTokenStatus_CHECKOUT_TOKEN_STATUS_NOT_FOUND_OR_EXPIRED,
			ErrorCode:   "TOKEN_NOT_FOUND_OR_EXPIRED",
		}, nil
	case 2:
		return &v1.ConsumeCheckoutTokenRes{
			TokenStatus: v1.CheckoutTokenStatus_CHECKOUT_TOKEN_STATUS_USER_MISMATCH,
			ErrorCode:   "TOKEN_USER_MISMATCH",
		}, nil
	case 1:
		if len(result) < 2 {
			return nil, gerror.NewCode(gcode.CodeInternalError, "token payload missing")
		}
		payload := gconv.String(result[1])
		var snapshot v1.CheckoutSnapshot
		if err = protojson.Unmarshal([]byte(payload), &snapshot); err != nil {
			return nil, gerror.Wrap(err, "unmarshal checkout snapshot failed")
		}
		return &v1.ConsumeCheckoutTokenRes{
			TokenStatus: v1.CheckoutTokenStatus_CHECKOUT_TOKEN_STATUS_OK,
			Snapshot:    &snapshot,
		}, nil
	default:
		return nil, gerror.NewCodef(gcode.CodeInternalError, "unknown token consume code: %d", code)
	}
}

// MarkItemsOrdered 订单创建后，清理购物车中已下单的 SKU。
func (s *sCart) MarkItemsOrdered(ctx context.Context, req *v1.MarkItemsOrderedReq) (*v1.MarkItemsOrderedRes, error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	if len(req.GetSkuNos()) == 0 {
		return &v1.MarkItemsOrderedRes{Affected: 0}, nil
	}

	key := s.cartKey(req.GetUserId())
	validSkuNos := make([]string, 0, len(req.GetSkuNos()))
	for _, skuNo := range req.GetSkuNos() {
		skuNo = strings.TrimSpace(skuNo)
		if skuNo != "" {
			validSkuNos = append(validSkuNos, skuNo)
		}
	}
	if len(validSkuNos) == 0 {
		return &v1.MarkItemsOrderedRes{Affected: 0}, nil
	}

	affected, err := g.Redis().HDel(ctx, key, validSkuNos...)
	if err != nil {
		return nil, gerror.Wrap(err, "mark ordered items failed")
	}
	if err = s.refreshCartTTL(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	_ = s.markDirtyUser(ctx, req.GetUserId())

	_, _ = dao.CartItemBackup.Ctx(ctx).
		Where(dao.CartItemBackup.Columns().UserId, req.GetUserId()).
		WhereIn(dao.CartItemBackup.Columns().SkuNo, validSkuNos).
		Delete()

	return &v1.MarkItemsOrderedRes{
		Affected: uint32(affected),
	}, nil
}

// BatchUpsertSkuProjection 批量回写 SKU 展示投影（标题/价格/是否可售等）。
func (s *sCart) BatchUpsertSkuProjection(ctx context.Context, req *v1.BatchUpsertSkuProjectionReq) (*v1.BatchUpsertSkuProjectionRes, error) {
	if len(req.GetItems()) == 0 {
		return &v1.BatchUpsertSkuProjectionRes{AffectedRows: 0}, nil
	}

	var affectedRows uint64
	err := dao.CartItemBackup.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.CartItemBackup.Columns()
		for _, item := range req.GetItems() {
			skuNo := strings.TrimSpace(item.GetSkuNo())
			if skuNo == "" {
				continue
			}
			data := do.CartItemBackup{
				SpuNo:           item.GetSpuNo(),
				ShopNo:          item.GetShopNo(),
				SpuTitle:        item.GetSpuTitle(),
				SkuName:         item.GetSkuName(),
				SkuImageAssetId: item.GetSkuImageAssetId(),
				SalePrice:       item.GetSalePrice(),
				MarketPrice:     item.GetMarketPrice(),
				SaleAttrsJson:   nonEmptyJSON(item.GetSaleAttrsJson()),
			}
			if item.GetSellable() {
				data.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE)
				data.InvalidReasonCode = ""
			} else {
				data.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_INVALID)
				data.InvalidReasonCode = nonEmpty(item.GetInvalidReasonCode(), "NOT_SELLABLE")
			}

			result, updateErr := tx.Model(dao.CartItemBackup.Table()).
				Where(cols.SkuNo, skuNo).
				Data(data).
				Update()
			if updateErr != nil {
				return gerror.Wrapf(updateErr, "update cart projection failed, sku_no=%s", skuNo)
			}
			if result != nil {
				rows, rowsErr := result.RowsAffected()
				if rowsErr == nil {
					affectedRows += uint64(rows)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.BatchUpsertSkuProjectionRes{
		AffectedRows: affectedRows,
	}, nil
}

// startBackupSyncWorker 启动异步刷盘任务，把 Redis 主存数据批量同步到 MySQL 备份表。
func (s *sCart) startBackupSyncWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			if err := s.flushDirtyUsers(ctx, 200); err != nil {
				g.Log().Warningf(ctx, "[cart-svc] flush dirty users failed: %+v", err)
			}
		}
	}()
}

func (s *sCart) flushDirtyUsers(ctx context.Context, limit int) error {
	if limit <= 0 {
		limit = 200
	}
	values, err := g.Redis().ZRange(
		ctx,
		consts.DirtyUsersZSetKey,
		0,
		int64(limit-1),
	)
	if err != nil {
		return gerror.Wrap(err, "query dirty users failed")
	}
	for _, v := range values {
		userID := gconv.Uint64(v.Val())
		if userID == 0 {
			continue
		}
		if syncErr := s.syncUserBackup(ctx, userID); syncErr != nil {
			g.Log().Warningf(ctx, "[cart-svc] sync user backup failed, user_id=%d, err=%+v", userID, syncErr)
			continue
		}
		_, _ = g.Redis().ZRem(ctx, consts.DirtyUsersZSetKey, strconv.FormatUint(userID, 10))
	}
	return nil
}

func (s *sCart) syncUserBackup(ctx context.Context, userID uint64) error {
	items, _, err := s.loadCartItems(ctx, userID, false)
	if err != nil {
		return err
	}

	return dao.CartItemBackup.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.CartItemBackup.Columns()
		if _, delErr := tx.Model(dao.CartItemBackup.Table()).
			Where(cols.UserId, userID).
			Delete(); delErr != nil {
			return gerror.Wrap(delErr, "clear old cart backup failed")
		}

		if len(items) > 0 {
			batch := make([]do.CartItemBackup, 0, len(items))
			for _, item := range items {
				batch = append(batch, do.CartItemBackup{
					UserId:            userID,
					SkuNo:             item.SkuNo,
					SpuNo:             item.SpuNo,
					ShopNo:            item.ShopNo,
					Qty:               item.Qty,
					Checked:           boolToInt(item.Checked),
					Status:            item.Status,
					InvalidReasonCode: item.InvalidReasonCode,
					SpuTitle:          item.SpuTitle,
					SkuName:           item.SkuName,
					SkuImageAssetId:   item.SkuImageAssetID,
					SalePrice:         item.SalePrice,
					MarketPrice:       item.MarketPrice,
					SaleAttrsJson:     nonEmptyJSON(item.SaleAttrsJSON),
				})
			}
			if _, insErr := tx.Model(dao.CartItemBackup.Table()).Data(batch).Insert(); insErr != nil {
				return gerror.Wrap(insErr, "insert cart backup batch failed")
			}
		}

		checkpointCols := dao.CartSyncCheckpoint.Columns()
		var checkpoint entity.CartSyncCheckpoint
		if scanErr := tx.Model(dao.CartSyncCheckpoint.Table()).
			Where(checkpointCols.UserId, userID).
			Scan(&checkpoint); scanErr != nil {
			return gerror.Wrap(scanErr, "query cart sync checkpoint failed")
		}
		now := gtime.Now()
		if checkpoint.Id == 0 {
			if _, insErr := tx.Model(dao.CartSyncCheckpoint.Table()).Data(do.CartSyncCheckpoint{
				UserId:          userID,
				LastSyncedAt:    now,
				LastSyncVersion: uint64(now.Unix()),
			}).Insert(); insErr != nil {
				return gerror.Wrap(insErr, "insert cart sync checkpoint failed")
			}
		} else {
			if _, updErr := tx.Model(dao.CartSyncCheckpoint.Table()).
				Where(checkpointCols.UserId, userID).
				Data(do.CartSyncCheckpoint{
					LastSyncedAt:    now,
					LastSyncVersion: uint64(now.Unix()),
				}).
				Update(); updErr != nil {
				return gerror.Wrap(updErr, "update cart sync checkpoint failed")
			}
		}
		return nil
	})
}

func (s *sCart) selectCheckoutItems(items map[string]*redisCartItem, req *v1.PrepareCheckoutReq) []*redisCartItem {
	result := make([]*redisCartItem, 0)
	switch req.GetScope() {
	case v1.CheckoutScope_CHECKOUT_SCOPE_PARTIAL:
		allow := make(map[string]struct{}, len(req.GetSkuNos()))
		for _, skuNo := range req.GetSkuNos() {
			skuNo = strings.TrimSpace(skuNo)
			if skuNo != "" {
				allow[skuNo] = struct{}{}
			}
		}
		for skuNo, item := range items {
			if _, ok := allow[skuNo]; !ok {
				continue
			}
			if item.Qty == 0 || v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
				continue
			}
			result = append(result, item)
		}
	default:
		for _, item := range items {
			if !item.Checked {
				continue
			}
			if item.Qty == 0 || v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
				continue
			}
			result = append(result, item)
		}
	}
	return result
}

func (s *sCart) loadCartItems(ctx context.Context, userID uint64, enableColdRestore bool) (map[string]*redisCartItem, bool, error) {
	key := s.cartKey(userID)
	records, err := g.Redis().HGetAll(ctx, key)
	if err != nil {
		return nil, false, gerror.Wrap(err, "read redis cart hash failed")
	}
	if records != nil && !records.IsNil() && len(records.Map()) > 0 {
		items, parseErr := parseRedisCartHash(records.Map())
		if parseErr != nil {
			return nil, false, parseErr
		}
		return items, false, nil
	}
	if !enableColdRestore {
		return map[string]*redisCartItem{}, false, nil
	}

	backupItems, err := s.loadCartItemsFromBackup(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	if len(backupItems) == 0 {
		return map[string]*redisCartItem{}, false, nil
	}

	redisFields := make(map[string]any, len(backupItems))
	for skuNo, item := range backupItems {
		payload, marshalErr := json.Marshal(item)
		if marshalErr != nil {
			return nil, false, gerror.Wrapf(marshalErr, "marshal cold restore item failed, sku_no=%s", skuNo)
		}
		redisFields[skuNo] = string(payload)
	}
	if err = g.Redis().HMSet(ctx, key, redisFields); err != nil {
		return nil, false, gerror.Wrap(err, "cold restore cart data to redis failed")
	}
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, false, err
	}
	return backupItems, true, nil
}

func (s *sCart) loadCartItemsFromBackup(ctx context.Context, userID uint64) (map[string]*redisCartItem, error) {
	cols := dao.CartItemBackup.Columns()
	var rows []*entity.CartItemBackup
	if err := dao.CartItemBackup.Ctx(ctx).
		Where(cols.UserId, userID).
		OrderAsc(cols.Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query cart backup failed")
	}

	items := make(map[string]*redisCartItem, len(rows))
	for _, row := range rows {
		items[row.SkuNo] = &redisCartItem{
			UserID:            row.UserId,
			SkuNo:             row.SkuNo,
			SpuNo:             row.SpuNo,
			ShopNo:            row.ShopNo,
			Qty:               uint32(row.Qty),
			Checked:           row.Checked == 1,
			Status:            int32(row.Status),
			InvalidReasonCode: row.InvalidReasonCode,
			SpuTitle:          row.SpuTitle,
			SkuName:           row.SkuName,
			SkuImageAssetID:   row.SkuImageAssetId,
			SalePrice:         row.SalePrice,
			MarketPrice:       row.MarketPrice,
			SaleAttrsJSON:     nonEmptyJSON(row.SaleAttrsJson),
			CreatedAtUnix:     toUnix(row.CreatedAt),
			UpdatedAtUnix:     toUnix(row.UpdatedAt),
		}
	}
	return items, nil
}

func (s *sCart) writeCartItem(ctx context.Context, userID uint64, item *redisCartItem) error {
	encoded, err := json.Marshal(item)
	if err != nil {
		return gerror.Wrapf(err, "marshal cart item failed, sku_no=%s", item.SkuNo)
	}
	_, err = g.Redis().HSet(ctx, s.cartKey(userID), map[string]any{
		item.SkuNo: string(encoded),
	})
	if err != nil {
		return gerror.Wrap(err, "write cart item to redis failed")
	}
	return nil
}

func (s *sCart) markDirtyUser(ctx context.Context, userID uint64) error {
	_, err := g.Redis().ZAdd(ctx, consts.DirtyUsersZSetKey, nil, gredis.ZAddMember{
		Score:  float64(time.Now().UTC().Unix()),
		Member: strconv.FormatUint(userID, 10),
	})
	if err != nil {
		return gerror.Wrap(err, "mark dirty user failed")
	}
	return nil
}

func (s *sCart) refreshCartTTL(ctx context.Context, userID uint64) error {
	_, err := g.Redis().Expire(ctx, s.cartKey(userID), consts.CartKeyTTLSeconds)
	if err != nil {
		return gerror.Wrap(err, "refresh cart ttl failed")
	}
	return nil
}

func (s *sCart) cartKey(userID uint64) string {
	return consts.CartKeyPrefix + strconv.FormatUint(userID, 10)
}

func (s *sCart) checkoutTokenPayloadKey(token string) string {
	return consts.CheckoutTokenKeyPrefix + token
}

func (s *sCart) checkoutTokenOwnerKey(token string) string {
	return consts.CheckoutTokenKeyPrefix + "owner:" + token
}

func (s *sCart) newCheckoutToken() string {
	return "chk_" + strings.ReplaceAll(guid.S(), "-", "")
}

func (s *sCart) mustUserID(ctx context.Context) (uint64, error) {
	// 优先读取 gRPC metadata，适配服务间 RPC。
	if userID, ok := userIDFromMetadata(ctx, "x-user-id"); ok && userID > 0 {
		return userID, nil
	}
	// 再读取 HTTP 头，适配 REST API 调用。
	if req := ghttp.RequestFromCtx(ctx); req != nil {
		header := strings.TrimSpace(req.Header.Get("X-User-Id"))
		if header == "" {
			header = strings.TrimSpace(req.Header.Get("x-user-id"))
		}
		if header != "" {
			if parsed, err := strconv.ParseUint(header, 10, 64); err == nil && parsed > 0 {
				return parsed, nil
			}
		}
	}
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "x-user-id is required")
}

func parseRedisCartHash(raw map[string]any) (map[string]*redisCartItem, error) {
	items := make(map[string]*redisCartItem, len(raw))
	for skuNo, value := range raw {
		payload := gconv.String(value)
		if strings.TrimSpace(payload) == "" {
			continue
		}
		var item redisCartItem
		if err := json.Unmarshal([]byte(payload), &item); err != nil {
			return nil, gerror.Wrapf(err, "unmarshal redis cart item failed, sku_no=%s", skuNo)
		}
		if strings.TrimSpace(item.SkuNo) == "" {
			item.SkuNo = skuNo
		}
		items[item.SkuNo] = &item
	}
	return items, nil
}

func summarize(items map[string]*redisCartItem) *v1.CartSummary {
	summary := &v1.CartSummary{}
	for _, item := range items {
		if item == nil {
			continue
		}
		summary.TotalItemCount += item.Qty
		if item.Checked {
			summary.CheckedItemCount += item.Qty
		}
		if !item.Checked {
			continue
		}
		if v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
			continue
		}
		summary.CheckedGoodsAmount += item.SalePrice * uint64(item.Qty)
	}
	// 当前阶段不计算券与运费拆分，结算金额=商品金额。
	summary.CheckedPayableAmount = summary.CheckedGoodsAmount
	return summary
}

func toProtoCartItem(item *redisCartItem) *v1.CartItem {
	if item == nil {
		return nil
	}
	return &v1.CartItem{
		UserId:            item.UserID,
		SkuNo:             item.SkuNo,
		SpuNo:             item.SpuNo,
		ShopNo:            item.ShopNo,
		Qty:               item.Qty,
		Checked:           item.Checked,
		Status:            v1.CartItemStatus(item.Status),
		InvalidReasonCode: item.InvalidReasonCode,
		SpuTitle:          item.SpuTitle,
		SkuName:           item.SkuName,
		SkuImageAssetId:   item.SkuImageAssetID,
		SalePrice:         item.SalePrice,
		MarketPrice:       item.MarketPrice,
		SaleAttrsJson:     item.SaleAttrsJSON,
		CreatedAt:         unixToProtoTs(item.CreatedAtUnix),
		UpdatedAt:         unixToProtoTs(item.UpdatedAtUnix),
	}
}

func digestSnapshot(snapshot *v1.CheckoutSnapshot) string {
	if snapshot == nil {
		return ""
	}
	// 使用固定顺序字段生成摘要，避免前后端串联时出现“同内容不同摘要”。
	type digestItem struct {
		SkuNo       string `json:"sku_no"`
		Qty         uint32 `json:"qty"`
		SettlePrice uint64 `json:"settle_price"`
	}
	payload := struct {
		UserID        uint64       `json:"user_id"`
		GoodsAmount   uint64       `json:"goods_amount"`
		FreightAmount uint64       `json:"freight_amount"`
		PayableAmount uint64       `json:"payable_amount"`
		Items         []digestItem `json:"items"`
	}{
		UserID:        snapshot.GetUserId(),
		GoodsAmount:   snapshot.GetGoodsAmount(),
		FreightAmount: snapshot.GetFreightAmount(),
		PayableAmount: snapshot.GetPayableAmount(),
		Items:         make([]digestItem, 0, len(snapshot.GetItems())),
	}
	for _, item := range snapshot.GetItems() {
		payload.Items = append(payload.Items, digestItem{
			SkuNo:       item.GetSkuNo(),
			Qty:         item.GetQty(),
			SettlePrice: item.GetSettlePrice(),
		})
	}
	sort.Slice(payload.Items, func(i, j int) bool { return payload.Items[i].SkuNo < payload.Items[j].SkuNo })
	raw, _ := json.Marshal(payload)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func unixToProtoTs(unixSec int64) *timestamppb.Timestamp {
	if unixSec <= 0 {
		return nil
	}
	return timestamppb.New(time.Unix(unixSec, 0).UTC())
}

func toUnix(t *gtime.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.Time.Unix()
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func nonEmptyJSON(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "[]"
	}
	return raw
}

func metadataValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(strings.ToLower(strings.TrimSpace(key)))
	if len(values) > 0 {
		return strings.TrimSpace(values[0])
	}
	return ""
}

func userIDFromMetadata(ctx context.Context, keys ...string) (uint64, bool) {
	for _, key := range keys {
		raw := metadataValue(ctx, key)
		if raw == "" {
			continue
		}
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && parsed > 0 {
			return parsed, true
		}
	}
	return 0, false
}
