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
// checkoutTokenConsumeLua 使用 Redis Lua 保证 checkout token 的消费具备原子性。
// 返回值约定：
// 0: token 不存在或已经过期；
// 1: token 校验通过并成功消费；
// 2: token 和 user_id 不匹配，此时不消费 token。
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

// sCart 是购物车领域的核心逻辑实现，负责 Redis 主存、结算快照和备份同步。
type sCart struct{}

// redisCartItem 是 Redis Hash 中单个购物车 SKU 的缓存快照结构。
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
// New 创建购物车逻辑实现。
func New() *sCart {
	// 当前购物车逻辑没有额外依赖初始化，直接返回空结构体即可。
	return &sCart{}
}

func init() {
	// 先创建购物车逻辑实例，后续注册和后台 worker 启动都复用它。
	svc := New()
	// 把购物车逻辑实现注册到 service 门面，供 controller 和其他逻辑调用。
	service.RegisterCart(svc)
	// 启动后台备份同步任务，把 Redis 主存数据周期性刷到 MySQL。
	svc.startBackupSyncWorker()
}

// AddItem 买家加购，购物车以 Redis Hash 为主存，并刷新 30 天滑动 TTL。
func (s *sCart) AddItem(ctx context.Context, req *v1.AddItemReq) (*v1.AddItemRes, error) {
	// sku_no 是购物车项的最小主键，没有它就无法定位商品。
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	// 加购数量必须大于 0，避免产生没有业务意义的空购物车项。
	if req.GetQty() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "qty must be greater than 0")
	}

	// 先提取当前买家 user_id，确保购物车严格按用户隔离。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 加载当前用户购物车快照，优先走 Redis，必要时支持从备份表冷恢复。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	// 统一记录本次修改时间，避免单个请求里出现多个不同时间戳。
	now := time.Now().UTC().Unix()
	// 规整 sku_no，避免不同空格形式造成重复 key。
	skuNo := strings.TrimSpace(req.GetSkuNo())
	// 先看购物车里是否已经有这个 SKU，决定是新增还是累加。
	item, exists := items[skuNo]
	if !exists {
		// 首次加购该 SKU 时，先构造一条新的购物车缓存项。
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
	// 加购的核心语义就是数量累加到现有值上。
	item.Qty += req.GetQty()
	// 刷新最后修改时间，供前端排序和异步同步使用。
	item.UpdatedAtUnix = now
	// 只有首次新增时才采用请求里的勾选状态作为初始值。
	if !exists {
		item.Checked = req.GetChecked()
	}
	// 遇到历史脏数据状态未指定时，主动纠正为可用态。
	if item.Status == int32(v1.CartItemStatus_CART_ITEM_STATUS_UNSPECIFIED) {
		item.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE)
	}
	// 把最新购物车项放回内存快照，便于后续汇总直接复用。
	items[skuNo] = item

	// 单条回写 Redis 主存，保证加购结果立刻对查询接口可见。
	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	// 刷新购物车滑动 TTL，保持活跃用户购物车数据存活。
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	// 把当前用户标记为 dirty，交给后台 worker 异步刷 MySQL 备份。
	_ = s.markDirtyUser(ctx, userID)

	// 返回最新购物车项和汇总结果，方便前端同步刷新角标和列表。
	return &v1.AddItemRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// UpdateItemQty 更新购物车中某个 SKU 的购买数量。
func (s *sCart) UpdateItemQty(ctx context.Context, req *v1.UpdateItemQtyReq) (*v1.UpdateItemQtyRes, error) {
	// 改数量时必须指定 SKU，才能精确命中购物车项。
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	// 数量必须大于 0，0 数量统一由删除接口处理。
	if req.GetQty() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "qty must be greater than 0")
	}

	// 提取当前买家身份，确保只能改自己的购物车。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 读取当前购物车快照，后续会在内存 map 上直接更新目标项。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	// 从购物车 map 中取出目标 SKU 项。
	item, ok := items[strings.TrimSpace(req.GetSkuNo())]
	if !ok {
		// 购物车里没有该商品时，返回 not found 更符合语义。
		return nil, gerror.NewCode(gcode.CodeNotFound, "cart item not found")
	}
	// 用调用方指定的新数量覆盖原值。
	item.Qty = req.GetQty()
	// 刷新更新时间，方便列表排序和后续同步。
	item.UpdatedAtUnix = time.Now().UTC().Unix()
	// 把修改后的对象写回内存快照。
	items[item.SkuNo] = item

	// 单条回写 Redis 主存。
	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	// 刷新购物车 TTL。
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	// 标记当前用户购物车有变更，等待后台备份同步。
	_ = s.markDirtyUser(ctx, userID)

	// 返回最新购物车项和汇总结果。
	return &v1.UpdateItemQtyRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// ToggleItemChecked 设置单个商品勾选状态。
func (s *sCart) ToggleItemChecked(ctx context.Context, req *v1.ToggleItemCheckedReq) (*v1.ToggleItemCheckedRes, error) {
	// 切换勾选状态时必须指定 SKU。
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	// 提取当前买家身份，避免操作到其他人的购物车。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 加载购物车快照，定位目标 SKU 项。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	// 从快照里取出目标购物车项。
	item, ok := items[strings.TrimSpace(req.GetSkuNo())]
	if !ok {
		// 目标商品不存在时直接返回 not found。
		return nil, gerror.NewCode(gcode.CodeNotFound, "cart item not found")
	}
	// 按请求值更新勾选状态。
	item.Checked = req.GetChecked()
	// 同步刷新更新时间，保证列表排序稳定。
	item.UpdatedAtUnix = time.Now().UTC().Unix()
	// 把更新后的对象回填到内存快照。
	items[item.SkuNo] = item

	// 单条回写 Redis，保证勾选状态即时生效。
	if err = s.writeCartItem(ctx, userID, item); err != nil {
		return nil, err
	}
	// 刷新购物车 TTL。
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, err
	}
	// 标记脏用户，等待后台刷盘。
	_ = s.markDirtyUser(ctx, userID)

	// 返回最新项和汇总结果给前端。
	return &v1.ToggleItemCheckedRes{
		Item:    toProtoCartItem(item),
		Summary: summarize(items),
	}, nil
}

// BatchToggleItems 批量设置勾选状态，常用于“全选/取消全选”。
func (s *sCart) BatchToggleItems(ctx context.Context, req *v1.BatchToggleItemsReq) (*v1.BatchToggleItemsRes, error) {
	// 没有任何 SKU 需要处理时，直接返回空结果。
	if len(req.GetSkuNos()) == 0 {
		return &v1.BatchToggleItemsRes{Affected: 0, Summary: &v1.CartSummary{}}, nil
	}
	// 提取当前买家身份，确保批量勾选只作用于本人购物车。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	// 加载购物车快照，后续在内存 map 上批量更新勾选状态。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	// affected 统计这次真正更新了多少个 SKU。
	affected := uint32(0)
	// fields 用来收集 HSet 需要的 field-value 参数，方便一次性批量写 Redis。
	fields := make(map[string]any, len(req.GetSkuNos()))
	// 统一本次批量操作的更新时间，避免同一批次出现不同时间戳。
	now := time.Now().UTC().Unix()
	for _, skuNo := range req.GetSkuNos() {
		// 规整 sku_no，避免空白字符影响 key 命中。
		skuNo = strings.TrimSpace(skuNo)
		// 先从购物车快照里取出目标项。
		item, ok := items[skuNo]
		if !ok {
			// 购物车里没有该 SKU 时直接跳过，不算受影响数量。
			continue
		}
		// 按请求统一设置勾选状态。
		item.Checked = req.GetChecked()
		// 把更新时间改成当前批次统一时间。
		item.UpdatedAtUnix = now
		// 更新内存快照，方便后续汇总直接复用。
		items[skuNo] = item
		// 先把购物车项编码成 JSON 字符串，作为 Redis Hash 的值。
		encoded, marshalErr := json.Marshal(item)
		if marshalErr != nil {
			// 任意一条编码失败都直接中断，避免半批成功半批失败。
			return nil, gerror.Wrapf(marshalErr, "marshal cart item failed, sku_no=%s", skuNo)
		}
		// 把编码后的结果收进批量 HSet 参数。
		fields[skuNo] = string(encoded)
		// 记录这条购物车项已被成功纳入本次更新。
		affected++
	}
	// 只有真的有字段需要更新时，才发起 Redis 批量写。
	if len(fields) > 0 {
		// 使用 HSet 批量更新 Redis 主存，减少网络往返。
		if _, err = g.Redis().HSet(ctx, s.cartKey(userID), fields); err != nil {
			return nil, gerror.Wrap(err, "batch toggle cart items failed")
		}
		// 刷新购物车 TTL，延续活跃购物车生命周期。
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		// 标记当前用户购物车已变更，等待后台异步刷盘。
		_ = s.markDirtyUser(ctx, userID)
	}

	// 返回受影响数量和最新购物车汇总。
	return &v1.BatchToggleItemsRes{
		Affected: affected,
		Summary:  summarize(items),
	}, nil
}

// RemoveItems 从购物车中删除指定商品。
func (s *sCart) RemoveItems(ctx context.Context, req *v1.RemoveItemsReq) (*v1.RemoveItemsRes, error) {
	// 没传待删除 SKU 列表时，直接按无操作返回。
	if len(req.GetSkuNos()) == 0 {
		return &v1.RemoveItemsRes{Removed: 0, Summary: &v1.CartSummary{}}, nil
	}
	// 提取当前买家 user_id，确保删除动作只命中本人购物车。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 先加载购物车快照，方便同步更新内存视图和 Redis 主存。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	// 收集真正需要从 Redis Hash 删除的 field 列表。
	keysToDelete := make([]string, 0, len(req.GetSkuNos()))
	// 统计实际删除成功的商品数量。
	removed := uint32(0)
	for _, skuNo := range req.GetSkuNos() {
		// 先规整 sku_no，去掉前后空白。
		skuNo = strings.TrimSpace(skuNo)
		if skuNo == "" {
			// 空 SKU 没有业务意义，直接跳过。
			continue
		}
		if _, ok := items[skuNo]; ok {
			// 购物车里确实存在该 SKU 时，先从内存快照删除。
			delete(items, skuNo)
			// 把它追加到 Redis 待删除列表里。
			keysToDelete = append(keysToDelete, skuNo)
			// 仅命中存在项时才累计删除数。
			removed++
		}
	}
	// 只有存在待删除字段时，才真正对 Redis 发起删除。
	if len(keysToDelete) > 0 {
		// 从 Redis Hash 中删除对应购物车项。
		if _, err = g.Redis().HDel(ctx, s.cartKey(userID), keysToDelete...); err != nil {
			return nil, gerror.Wrap(err, "remove cart items failed")
		}
		// 刷新购物车 TTL，保持活跃购物车数据继续存活。
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		// 标记脏用户，等待后台同步到备份表。
		_ = s.markDirtyUser(ctx, userID)
	}

	// 返回删除数量和删除后的汇总信息。
	return &v1.RemoveItemsRes{
		Removed: removed,
		Summary: summarize(items),
	}, nil
}

// ClearInvalidItems 清空失效商品。
func (s *sCart) ClearInvalidItems(ctx context.Context, _ *v1.ClearInvalidItemsReq) (*v1.ClearInvalidItemsRes, error) {
	// 先获取当前买家身份，确保清理范围限定在本人购物车。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	// 加载当前购物车快照，后续筛出状态无效的商品。
	items, _, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	// keysToDelete 用于收集失效商品的 SKU，稍后统一从 Redis 删除。
	keysToDelete := make([]string, 0)
	// removed 统计本次清理掉多少个失效商品。
	removed := uint32(0)
	for skuNo, item := range items {
		// 只有状态明确为 INVALID 的购物车项才属于可清理对象。
		if v1.CartItemStatus(item.Status) == v1.CartItemStatus_CART_ITEM_STATUS_INVALID {
			// 记录待删 SKU，方便后续批量 HDel。
			keysToDelete = append(keysToDelete, skuNo)
			// 同步从内存快照删除，保证汇总信息立刻反映最新状态。
			delete(items, skuNo)
			// 统计清理成功的失效项数量。
			removed++
		}
	}
	// 只有确实存在失效项时才发起 Redis 删除请求。
	if len(keysToDelete) > 0 {
		// 批量删除 Redis 主存中的失效 SKU。
		if _, err = g.Redis().HDel(ctx, s.cartKey(userID), keysToDelete...); err != nil {
			return nil, gerror.Wrap(err, "clear invalid cart items failed")
		}
		// 刷新 TTL，保持购物车滑动过期语义不变。
		if err = s.refreshCartTTL(ctx, userID); err != nil {
			return nil, err
		}
		// 标记用户购物车已更新，交给异步 worker 做备份同步。
		_ = s.markDirtyUser(ctx, userID)
	}

	// 返回清理数量和清理后的汇总结果。
	return &v1.ClearInvalidItemsRes{
		Removed: removed,
		Summary: summarize(items),
	}, nil
}

// GetMyCart 获取当前用户购物车；Redis miss 时支持从 MySQL 备份冷恢复。
func (s *sCart) GetMyCart(ctx context.Context, req *v1.GetMyCartReq) (*v1.GetMyCartRes, error) {
	// 先拿到当前买家身份，确定购物车读取范围。
	userID, err := s.mustUserID(ctx)
	if err != nil {
		return nil, err
	}
	// 加载购物车快照，并记录是否发生了从 MySQL 备份表冷恢复。
	items, loadedFromBackup, err := s.loadCartItems(ctx, userID, true)
	if err != nil {
		return nil, err
	}

	// 预分配响应切片，承接转换后的购物车项。
	out := make([]*v1.CartItem, 0, len(items))
	for _, item := range items {
		// 只看已勾选商品时，跳过未勾选项。
		if req.GetOnlyChecked() && !item.Checked {
			continue
		}
		// 把 Redis 快照项转换成对外协议对象。
		out = append(out, toProtoCartItem(item))
	}
	// 按最近更新时间倒序排列，保持最新改动的购物车项优先展示。
	sort.Slice(out, func(i, j int) bool {
		if out[i].UpdatedAt == nil {
			return false
		}
		if out[j].UpdatedAt == nil {
			return true
		}
		return out[i].UpdatedAt.AsTime().After(out[j].UpdatedAt.AsTime())
	})

	// 返回购物车列表、汇总数据以及是否命中冷恢复标记。
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

// MarkItemsOrdered ??????????????????? SKU?
func (s *sCart) MarkItemsOrdered(ctx context.Context, req *v1.MarkItemsOrderedReq) (*v1.MarkItemsOrderedRes, error) {
	// ?????????? user_id?????????????????
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	// ??????? SKU ?????? 0?????????
	if len(req.GetSkuNos()) == 0 {
		return &v1.MarkItemsOrderedRes{Affected: 0}, nil
	}

	// ?????????? Redis ????
	key := s.cartKey(req.GetUserId())
	// validSkuNos ????????? SKU????????? Redis HDel?
	validSkuNos := make([]string, 0, len(req.GetSkuNos()))
	for _, skuNo := range req.GetSkuNos() {
		// ?????????????????????????
		skuNo = strings.TrimSpace(skuNo)
		if skuNo != "" {
			// ???? SKU ??????????
			validSkuNos = append(validSkuNos, skuNo)
		}
	}
	// ????????? SKU?????????
	if len(validSkuNos) == 0 {
		return &v1.MarkItemsOrderedRes{Affected: 0}, nil
	}

	// ? Redis ???????????????????????? SKU?
	affected, err := g.Redis().HDel(ctx, key, validSkuNos...)
	if err != nil {
		return nil, gerror.Wrap(err, "mark ordered items failed")
	}
	// ?????????? TTL?????????????????
	if err = s.refreshCartTTL(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	// ??????? dirty????? worker ????????????
	_ = s.markDirtyUser(ctx, req.GetUserId())

	// ??????????? SKU????????????????????
	_, _ = dao.CartItemBackup.Ctx(ctx).
		Where(dao.CartItemBackup.Columns().UserId, req.GetUserId()).
		WhereIn(dao.CartItemBackup.Columns().SkuNo, validSkuNos).
		Delete()

	// ????? Redis ??? field ???????????
	return &v1.MarkItemsOrderedRes{
		Affected: uint32(affected),
	}, nil
}

// BatchUpsertSkuProjection ???? SKU ???????????????????
func (s *sCart) BatchUpsertSkuProjection(ctx context.Context, req *v1.BatchUpsertSkuProjectionReq) (*v1.BatchUpsertSkuProjectionRes, error) {
	// ????????????????? 0?
	if len(req.GetItems()) == 0 {
		return &v1.BatchUpsertSkuProjectionRes{AffectedRows: 0}, nil
	}

	// affectedRows ???????????????????
	var affectedRows uint64
	// ?????????????? SKU ?????????????????????
	err := dao.CartItemBackup.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// ????????????? Columns ???
		cols := dao.CartItemBackup.Columns()
		for _, item := range req.GetItems() {
			// sku_no ????????????????????
			skuNo := strings.TrimSpace(item.GetSkuNo())
			if skuNo == "" {
				// ? SKU ?????????????????
				continue
			}
			// ?????????????????? DO ????
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
			// ????????? ACTIVE?????????
			if item.GetSellable() {
				data.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE)
				data.InvalidReasonCode = ""
			} else {
				// ????????????? INVALID????????????
				data.Status = int32(v1.CartItemStatus_CART_ITEM_STATUS_INVALID)
				// ????????????? NOT_SELLABLE ????????
				data.InvalidReasonCode = nonEmpty(item.GetInvalidReasonCode(), "NOT_SELLABLE")
			}

			// ? sku_no ????????????? Redis ??????????????????
			result, updateErr := tx.Model(dao.CartItemBackup.Table()).
				Where(cols.SkuNo, skuNo).
				Data(data).
				Update()
			if updateErr != nil {
				return gerror.Wrapf(updateErr, "update cart projection failed, sku_no=%s", skuNo)
			}
			if result != nil {
				// ???? update ????????????????????
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

	// ???????????????????????
	return &v1.BatchUpsertSkuProjectionRes{
		AffectedRows: affectedRows,
	}, nil
}

// startBackupSyncWorker ?????????? Redis ???????? MySQL ????
func (s *sCart) startBackupSyncWorker() {
	// ???? goroutine ???????????????
	go func() {
		// ??????? dirty user ????????????
		ticker := time.NewTicker(1 * time.Minute)
		// worker ????? ticker ????????
		defer ticker.Stop()
		for range ticker.C {
			// ???????????????????????????
			ctx := context.Background()
			if err := s.flushDirtyUsers(ctx, 200); err != nil {
				// ????????????? worker??????????????
				g.Log().Warningf(ctx, "[cart-svc] flush dirty users failed: %+v", err)
			}
		}
	}()
}

func (s *sCart) flushDirtyUsers(ctx context.Context, limit int) error {
	// limit ??????????????? worker ??????? 0?
	if limit <= 0 {
		limit = 200
	}
	// ??????????????????????????
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
		// ZSet member ??? user_id ???????? uint64?
		userID := gconv.Uint64(v.Val())
		if userID == 0 {
			// ??????????????????????
			continue
		}
		if syncErr := s.syncUserBackup(ctx, userID); syncErr != nil {
			// ???????????????????????????
			g.Log().Warningf(ctx, "[cart-svc] sync user backup failed, user_id=%d, err=%+v", userID, syncErr)
			continue
		}
		// ???????????????????????
		_, _ = g.Redis().ZRem(ctx, consts.DirtyUsersZSetKey, strconv.FormatUint(userID, 10))
	}
	return nil
}

func (s *sCart) syncUserBackup(ctx context.Context, userID uint64) error {
	// ?? Redis ????????????????? MySQL ??????? Redis ????
	items, _, err := s.loadCartItems(ctx, userID, false)
	if err != nil {
		return err
	}

	// ?????????????????????? checkpoint?????????
	return dao.CartItemBackup.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// ???????? where/update ?????
		cols := dao.CartItemBackup.Columns()
		// ???????????????????????????????????
		if _, delErr := tx.Model(dao.CartItemBackup.Table()).
			Where(cols.UserId, userID).
			Delete(); delErr != nil {
			return gerror.Wrap(delErr, "clear old cart backup failed")
		}

		if len(items) > 0 {
			// batch ???????????????????????
			batch := make([]do.CartItemBackup, 0, len(items))
			for _, item := range items {
				// ? Redis ????????? MySQL ?????????????????
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
			// ??????????????????????
			if _, insErr := tx.Model(dao.CartItemBackup.Table()).Data(batch).Insert(); insErr != nil {
				return gerror.Wrap(insErr, "insert cart backup batch failed")
			}
		}

		// checkpoint ?????????????????????
		checkpointCols := dao.CartSyncCheckpoint.Columns()
		var checkpoint entity.CartSyncCheckpoint
		if scanErr := tx.Model(dao.CartSyncCheckpoint.Table()).
			Where(checkpointCols.UserId, userID).
			Scan(&checkpoint); scanErr != nil {
			return gerror.Wrap(scanErr, "query cart sync checkpoint failed")
		}
		// ????????????????????
		now := gtime.Now()
		if checkpoint.Id == 0 {
			// ??? checkpoint ???????????????????????
			if _, insErr := tx.Model(dao.CartSyncCheckpoint.Table()).Data(do.CartSyncCheckpoint{
				UserId:          userID,
				LastSyncedAt:    now,
				LastSyncVersion: uint64(now.Unix()),
			}).Insert(); insErr != nil {
				return gerror.Wrap(insErr, "insert cart sync checkpoint failed")
			}
		} else {
			// ??? checkpoint ???????????????????
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
	// result ????????????????
	result := make([]*redisCartItem, 0)
	switch req.GetScope() {
	case v1.CheckoutScope_CHECKOUT_SCOPE_PARTIAL:
		// ????????????? SKU ????? allow set?
		allow := make(map[string]struct{}, len(req.GetSkuNos()))
		for _, skuNo := range req.GetSkuNos() {
			// ??????????? SKU ?????????????
			skuNo = strings.TrimSpace(skuNo)
			if skuNo != "" {
				// ???? SKU ????????
				allow[skuNo] = struct{}{}
			}
		}
		for skuNo, item := range items {
			// ??????????????????
			if _, ok := allow[skuNo]; !ok {
				continue
			}
			// ??? 0 ???? ACTIVE ???????????????
			if item.Qty == 0 || v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
				continue
			}
			// ????????????????
			result = append(result, item)
		}
	default:
		for _, item := range items {
			// ?????????????????
			if !item.Checked {
				continue
			}
			// ??????? 0 ???????????????????
			if item.Qty == 0 || v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
				continue
			}
			// ??????????????
			result = append(result, item)
		}
	}
	// ??????????????????? checkout snapshot?
	return result
}

func (s *sCart) loadCartItems(ctx context.Context, userID uint64, enableColdRestore bool) (map[string]*redisCartItem, bool, error) {
	// ??? Redis ?????????? Redis ?????
	key := s.cartKey(userID)
	// ?? Redis Hash ?????????????????
	records, err := g.Redis().HGetAll(ctx, key)
	if err != nil {
		return nil, false, gerror.Wrap(err, "read redis cart hash failed")
	}
	if records != nil && !records.IsNil() && len(records.Map()) > 0 {
		// Redis ???????????????????????
		items, parseErr := parseRedisCartHash(records.Map())
		if parseErr != nil {
			return nil, false, parseErr
		}
		return items, false, nil
	}
	// ???????Redis miss ???????????
	if !enableColdRestore {
		return map[string]*redisCartItem{}, false, nil
	}

	// Redis miss ?????????? MySQL ?????????
	backupItems, err := s.loadCartItemsFromBackup(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	// ??????????????????????
	if len(backupItems) == 0 {
		return map[string]*redisCartItem{}, false, nil
	}

	// ??????????? Redis Hash ??????????
	redisFields := make(map[string]any, len(backupItems))
	for skuNo, item := range backupItems {
		// ????????????? JSON ????? Redis?
		payload, marshalErr := json.Marshal(item)
		if marshalErr != nil {
			return nil, false, gerror.Wrapf(marshalErr, "marshal cold restore item failed, sku_no=%s", skuNo)
		}
		redisFields[skuNo] = string(payload)
	}
	// ????????? Redis?????????
	if err = g.Redis().HMSet(ctx, key, redisFields); err != nil {
		return nil, false, gerror.Wrap(err, "cold restore cart data to redis failed")
	}
	// ???????? TTL????????????
	if err = s.refreshCartTTL(ctx, userID); err != nil {
		return nil, false, err
	}
	// ??????????? loadedFromBackup ??? true?
	return backupItems, true, nil
}

func (s *sCart) loadCartItemsFromBackup(ctx context.Context, userID uint64) (map[string]*redisCartItem, error) {
	// ?????????????????
	cols := dao.CartItemBackup.Columns()
	// rows ??????????????????
	var rows []*entity.CartItemBackup
	if err := dao.CartItemBackup.Ctx(ctx).
		Where(cols.UserId, userID).
		OrderAsc(cols.Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query cart backup failed")
	}

	// items ????????? Redis ???key ??? sku_no ???
	items := make(map[string]*redisCartItem, len(rows))
	for _, row := range rows {
		// ??????????? Redis ???????????????
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
	// ????????????? map?
	return items, nil
}

func (s *sCart) writeCartItem(ctx context.Context, userID uint64, item *redisCartItem) error {
	// ??????? Redis ?? JSON ??????? Hash value ??
	encoded, err := json.Marshal(item)
	if err != nil {
		return gerror.Wrapf(err, "marshal cart item failed, sku_no=%s", item.SkuNo)
	}
	// ? sku_no ?? field ?? Redis Hash?????????????
	_, err = g.Redis().HSet(ctx, s.cartKey(userID), map[string]any{
		item.SkuNo: string(encoded),
	})
	if err != nil {
		return gerror.Wrap(err, "write cart item to redis failed")
	}
	return nil
}

func (s *sCart) markDirtyUser(ctx context.Context, userID uint64) error {
	// ?????????score ?????????????? worker ???
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
	// ??????????? TTL??????????? 30 ??????
	_, err := g.Redis().Expire(ctx, s.cartKey(userID), consts.CartKeyTTLSeconds)
	if err != nil {
		return gerror.Wrap(err, "refresh cart ttl failed")
	}
	return nil
}

func (s *sCart) cartKey(userID uint64) string {
	// ??????????? user_id?????????
	return consts.CartKeyPrefix + strconv.FormatUint(userID, 10)
}

func (s *sCart) checkoutTokenPayloadKey(token string) string {
	// payload key ???????????????????
	return consts.CheckoutTokenKeyPrefix + token
}

func (s *sCart) checkoutTokenOwnerKey(token string) string {
	// owner key ???? token ?? user_id??? token ????
	return consts.CheckoutTokenKeyPrefix + "owner:" + token
}

func (s *sCart) newCheckoutToken() string {
	// ?????????? token?????????????
	return "chk_" + strings.ReplaceAll(guid.S(), "-", "")
}

func (s *sCart) mustUserID(ctx context.Context) (uint64, error) {
	// ??? gRPC metadata ?? user_id?????? RPC ???
	if userID, ok := userIDFromMetadata(ctx, "x-user-id"); ok && userID > 0 {
		return userID, nil
	}
	// ?? HTTP Header ?? user_id??? REST API ?????
	if req := ghttp.RequestFromCtx(ctx); req != nil {
		header := strings.TrimSpace(req.Header.Get("X-User-Id"))
		if header == "" {
			// ????????????????????????
			header = strings.TrimSpace(req.Header.Get("x-user-id"))
		}
		if header != "" {
			// ???????? user_id ??????
			if parsed, err := strconv.ParseUint(header, 10, 64); err == nil && parsed > 0 {
				return parsed, nil
			}
		}
	}
	// ?????????????????????
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "x-user-id is required")
}

func parseRedisCartHash(raw map[string]any) (map[string]*redisCartItem, error) {
	// items ?? Redis Hash ??????????
	items := make(map[string]*redisCartItem, len(raw))
	for skuNo, value := range raw {
		// ?? Redis value ??????? JSON?
		payload := gconv.String(value)
		if strings.TrimSpace(payload) == "" {
			// ? payload ????????????
			continue
		}
		// ?????? Redis ???????
		var item redisCartItem
		if err := json.Unmarshal([]byte(payload), &item); err != nil {
			return nil, gerror.Wrapf(err, "unmarshal redis cart item failed, sku_no=%s", skuNo)
		}
		if strings.TrimSpace(item.SkuNo) == "" {
			// ???????? sku_no ??????? hash field ?????
			item.SkuNo = skuNo
		}
		// ?? map?????? sku_no ?????
		items[item.SkuNo] = &item
	}
	// ????????????
	return items, nil
}

func summarize(items map[string]*redisCartItem) *v1.CartSummary {
	// summary ????????????????????
	summary := &v1.CartSummary{}
	for _, item := range items {
		if item == nil {
			// ?????????????? panic?
			continue
		}
		// ???????????????
		summary.TotalItemCount += item.Qty
		if item.Checked {
			// ??????????????????
			summary.CheckedItemCount += item.Qty
		}
		if !item.Checked {
			// ???????????????
			continue
		}
		if v1.CartItemStatus(item.Status) != v1.CartItemStatus_CART_ITEM_STATUS_ACTIVE {
			// ?????????????????????
			continue
		}
		// ??????????????????????
		summary.CheckedGoodsAmount += item.SalePrice * uint64(item.Qty)
	}
	// ???????????????????????????
	summary.CheckedPayableAmount = summary.CheckedGoodsAmount
	return summary
}

func toProtoCartItem(item *redisCartItem) *v1.CartItem {
	// ??????? nil???????????
	if item == nil {
		return nil
	}
	// ??? Redis ????????? proto ???
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
	// ?????????????
	if snapshot == nil {
		return ""
	}
	// ?????????????????????????????????
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
		// ?????????????? digest ????????????
		payload.Items = append(payload.Items, digestItem{
			SkuNo:       item.GetSkuNo(),
			Qty:         item.GetQty(),
			SettlePrice: item.GetSettlePrice(),
		})
	}
	// ?????? sku_no ?????????????????????????
	sort.Slice(payload.Items, func(i, j int) bool { return payload.Items[i].SkuNo < payload.Items[j].SkuNo })
	// ???????? JSON??? sha256 ????
	raw, _ := json.Marshal(payload)
	// ?? sha256??? checkout snapshot ???????
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func unixToProtoTs(unixSec int64) *timestamppb.Timestamp {
	// ????????????????? nil?
	if unixSec <= 0 {
		return nil
	}
	// ? Unix ????????? protobuf Timestamp?
	return timestamppb.New(time.Unix(unixSec, 0).UTC())
}

func toUnix(t *gtime.Time) int64 {
	// ????????????? 0???????????
	if t == nil || t.IsZero() {
		return 0
	}
	// ???? Unix ?????????????????
	return t.Time.Unix()
}

func boolToInt(v bool) int {
	// true ??? 1?????????????????
	if v {
		return 1
	}
	// false ??? 0?
	return 0
}

func nonEmpty(values ...string) string {
	for _, value := range values {
		// ????????????????????
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	// ????????????????
	return ""
}

func nonEmptyJSON(raw string) string {
	// ? JSON ?????????????????????????
	if strings.TrimSpace(raw) == "" {
		return "[]"
	}
	// ???????????????????
	return raw
}

func metadataValue(ctx context.Context, key string) string {
	// ? incoming metadata ???? key??? gRPC ???????
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	// metadata key ???????????????????????
	values := md.Get(strings.ToLower(strings.TrimSpace(key)))
	if len(values) > 0 {
		// ?????????????????
		return strings.TrimSpace(values[0])
	}
	// ????????????
	return ""
}

func userIDFromMetadata(ctx context.Context, keys ...string) (uint64, bool) {
	for _, key := range keys {
		// ?????? key????????? user_id ??????
		raw := metadataValue(ctx, key)
		if raw == "" {
			continue
		}
		// ???????? user_id ????????
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && parsed > 0 {
			return parsed, true
		}
	}
	// ?? key ???????? user_id ??? false?
	return 0, false
}
