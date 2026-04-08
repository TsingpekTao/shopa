---
title: 预发货退款流程 V1 / 售后中心
date: 2026-04-08
team: seller-admin
---

# 预发货退款流程 V1（Seller After-Sale Center）

## Summary

“预发货退款流程 V1”要求在 seller-admin 内新增一个「售后中心」导航入口，并提供可按当前店铺查看的退款批次列表，辅以详情/审核页（同意/驳回）。后端接口已经准备好，前端需要补充 API/type 层、状态/原因的显示逻辑，并保证中英文文案没有乱码。当前任务也顺便修复我接触到的发货中心页面中的乱码文本。

## Goals

1. 在 `SellerLayout` 的导航菜单里新增“售后中心”入口（按店铺 scope，优先显示待审核），并复用已有的 shop picker + badge 工具。
2. 实现 `/seller/after-sales`（列表）和 `/seller/after-sales/[refundBatchNo]`（详情）两个路由，列表提供按以下状态筛选的 chips/tab：`PENDING_SELLER_REVIEW`（默认）、`WAIT_REFUND_TASK`、`REFUND_PROCESSING`、`REFUNDED`、`SELLER_REJECTED`、`CANCELED`、`CLOSED`，并允许显示所有状态。
3. 详情页展示批次摘要、买家+订单信息、每个 AfterSaleCase 的进度/金额，以及 refund task timeline；提供「同意退款（可写备注）」和「驳回退款（选择原因/备注）」两个操作。
4. 在 `seller-admin` 侧补充 aftersale API 方法、类型/状态映射、金额/时间格式等 helper，并添加至少 1-2 个映射/状态测试（例如状态 -> 颜色、拒绝原因 label）。
5. 修复发货中心页面 (`seller/shipping/page.tsx`) 中可见的中文乱码（例如状态标签、表头、按钮文案等）。

## Constraints

- 所有变更仅限 `apps/seller-admin/**`（包括 `globals.css` 等样式文件）。
- 页面风格要贴近现有的 seller-admin/会话中心，遵循现有的 CSS 命名空间。
- 先写测试，再实现（符合 TDD 要求），测试命令仍然可以用 `npx tsx` 运行。
- 不添加与 mall-web 或后端无关的修改；不回滚别人改动。
- 在 detail 页中保留 case 级别的数据（如果你希望只要批次摘要我可以调整，但当前设想是 case+task 信息都要可见）。

## Proposed Solution

### Navigation & shop context

1. 在 `seller/layout.tsx` 的 menu items 中插入新的 key `/seller/after-sales`（图标可选 `ExclamationCircleOutlined` 或类似）并确保 `selectedKey` 逻辑包含该路径，同时保持店铺 badge 机制。
2. 列表/详情页各自维护和持久化一个本地 shop 选择（参考 `SELLER_SHIPPING_LAST_SHOP_KEY` 逻辑），并同步 `SELLER_CURRENT_SHOP_NO`（保持与 badge/会话中心一致的当前店铺）。

### List page (`/seller/after-sales`)

- 页面结构类似 `shipping`/`conversations`：顶部 shop select + 工具栏（状态 chips、刷新、搜索），中间展示统计卡片（待审核数量、金额、当前 shop 等）。
- 状态 chips/tab 采用 Ant Design `Segmented` + `Space`，默认选中 `PENDING_SELLER_REVIEW`；用户可以点击显示任意状态，API `statuses` 参数同步 chips。
- 列表使用 `Table`，列包括：批次号、订单号、状态（带颜色）、申请金额、审核金额、案例数、子订单、创建/截止时间、操作按钮（“查看详情”）。
- 通过 `useQuery` 调用 `/v1/aftersale/seller/shops/{shopNo}/refund-batches`（带 `statuses`、`page_size`、`next_cursor`），分页/加载更多暂时不做，先取 50 条即可；默认把 `PENDING_SELLER_REVIEW` 排在最前，可通过 `orderBy` (e.g., `sort` in UI after fetching).
- 支持 `message` 提示错误、`Skeleton/Empty` 友好占位。

### Detail page (`/seller/after-sales/[refundBatchNo]`)

- 使用 `useParams` 或 `useRouter` 读取 `refundBatchNo`，再调用 `/v1/aftersale/seller/refund-batches/{refundBatchNo}?shopNo=...` 获取 `RefundBatchDetail`。
- 上方展示 `Card` 的批次信息：批次号、状态（颜色）、申请/审核金额、总案例数、创建/更新时间、审核截止、自动触发时间等。
- 下面分两栏：左侧 `AfterSaleCase` 列表，展示每个 case 的类型（退货/仅退款）、订单号、子订单号、申请/同意金额、原因文本、状态；右侧 `RefundTask` 时间线（或列表）展示任务号、状态、退款金额、错误/重试信息。
- 页面底部提供 action buttons（“同意退款”、“驳回退款”），点击弹出 `Modal` 或 `Drawer` 让卖家填写可选备注/选择拒绝原因（参考 `RejectReasonCode_*` enum）。提交后调用对应 POST 接口刷新数据并显示 `message.success`。
- Approve payload: `{ refund_batch_no, shop_no, seller_reply?, idempotency_key? }`. Reject payload also requires `reject_reason_code`.
- Actions 之后 invalidate detail query (and list query via `queryClient.invalidateQueries`).

### API/Type layer

- 新建 `features/aftersale` 目录，包含 `api.ts`, `types.ts`, `helpers.ts`。
- `types.ts` 提供 `AfterSaleStatusCode`, `RejectReasonCode`, `RefundBatch`, `RefundBatchDetail`, `AfterSaleCase`, `RefundTask`, `RefundBatchListParams/Response` 等类型 with proper normalization helpers (similar to `order/api.ts`).
- `helpers.ts` 提供：
  - `normalizeAfterSaleStatus(value): AfterSaleStatusCode` (support numeric/string/upper/lower/different shapes).
  - `getAfterSaleStatusMeta(status): { label: string; color: string; badge?: string }`.
  - `getRejectReasonOptions()` returning label/value for select.
  - `formatCny(amountInCents)` for amounts.
- Tests in `features/aftersale/helpers.test.ts` or similar cover normalization (string vs number, unknown value) and reason options label mapping per requirements.
- API functions use `apiClient.get/post` (list/detail, approve/reject) and rely on `normalize` functions before returning typed data.

### Styles & Copy

- 新页面采用 `.seller-after-sale-page` 等命名空间，充分复用 `globals.css` 中的卡片/Toolbar 样式。
- `globals.css` 需要补充 `.seller-after-sale-chips`, `.seller-after-sale-table-card` 等，尽量重用现有变量。
- 同时把 `seller/shipping/page.tsx` 里的中文字符串改为可读中文（例如“订单编号”、“商品”）并修复状态标签。

## Testing Strategy

1. 新增 `features/aftersale/helpers.test.ts`，测试 `normalizeAfterSaleStatus` 正确处理数字/字符串、`getRejectReasonOptions` 中的 label/payload。这样至少覆盖两种映射 helper，满足 TDD 要求。
2. 运行 `npx tsx src/features/aftersale/helpers.test.ts` 作为验证命令，确保 TypeScript 直接执行。

## Next Steps

1. 请确认上述设计（尤其是 detail 还包含 case 级信息）是否符合预期。
2. 若确认，我会把此内容存为 spec 文档并请你再审核一遍；只有在 spec 被批准后再转入 implementation planning。
