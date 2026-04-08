# 预发货退款流程 V1 设计

**日期**: 2026-04-08

## 背景
- 买家当前仅能在售后标签中临时查看退款状态，没有专属的退款入口，也无法直接从待发货订单发起退款。
- 目标是在 mall-web 端打造预发货退款的完整流程：点击订单进入申请，先预览再正式提交，同时用一个列表/详情页维护当前所有退款批次，保持 TAOBAO 风格。

## 数据 & API
- 新建 eatures/refund 模块，定义和平滑 RefundBatch、RefundBatchDetail、AfterSaleCase、RefundTask、PreviewRefundRes 等类型，统一把 snake_case 字段映射成 camelCase，并导出 AfterSaleStatus 相关枚举/label 工具。
- 新增几个 helper：ormatAfterSaleStatusLabel(status, isZh)、isRefundBatchCancelable(status)、hasRefundableSubOrder(order) 等，对外抛给页面复用，配合状态文本和按钮可用性。
- 使用 piClient 编写新的接口函数，将请求数据序列化成后端期望的字段：
  * previewRefund({ orderNo, subOrderNo, itemNo, selectedItemNos })
  * pplyRefundBatch({ targets, reasonCode, reasonDesc, buyerRemark })
  * listRefundBatches({ pageSize?, statuses? })
  * getRefundBatchDetail(refundBatchNo)
  * cancelRefundBatch(refundBatchNo)
- preview 接口的 snapshot 里会带 item 列表和可退金额，apply 页面会根据 	argets 里传递的 orderNo/subOrderNo/itemNo 自动预选对应行。

## 页面 & 路由
1. /me/refunds/page.tsx
   - 使用 useQuery 加载 listRefundBatches，展示每条批次号、状态 Chip、申请/批准金额、创建时间、案例数。
   - 状态 Chip 通过 ormatAfterSaleStatusLabel 区分色彩；可执行 Cancel 的状态（如 PENDING_SELLER_REVIEW、SELLER_REJECTED、WAIT_REFUND_TASK）会在卡片右侧出现按钮。
   - 列表项可跳转 /me/refunds/{refundBatchNo} 查看详情。
   - 在 page.tsx 上方或侧边提供  申请退款 主要入口，链接到 /me/refunds/apply，并在订单页面的 	b-order-card-actions 中加入 申请退款 字样的按钮，仅在 hasRefundableSubOrder(order) 的场景下显示，直接带上 ?orderNo=...&subOrderNo=...&itemNo=...。

2. /me/refunds/[refundBatchNo]/page.tsx
   - 加载单个批次详情（getRefundBatchDetail），展示 batch 头部（状态、总申请/批准金额、创建时间、review deadline/auto approve）、每个 AfterSaleCase 的订单/商品摘要、申请原因、当前 fterSaleStatus、关联 refund tasks（按照时间排，显示状态、金额、错误信息）。
   - 状态高亮、条目之间用 cards 或 	b-panel 风格视觉分隔。
   - 如果入口仍然在可取消状态，页脚也提供取消按钮。

3. /me/refunds/apply/page.tsx
   - 读取 orderNo/subOrderNo/itemNo query，调用 previewRefund，展示 snapshot 里 efundableAmount、包含商品、付款状态等内容，并允许用户选择 reason code（写几个固定选项，例如 changed_mind、ound_a_better_price、shipping_delay）和填写备注。
   - 需要验证 reason 描述非空，点击 提交退款 时调 pplyRefundBatch，成功后 outer.push(/me/refunds/) 并弹出 message.success。
   - 预览阶段也会提醒用户退款的预计金额，必要时展示 selectedItemNos 说明用户正在退货的项（若 API 返回）。

4. 订单页修改：
   - 依赖 hasRefundableSubOrder 判断是否有待发货子单，并在 	b-order-card-actions 追加 Link 到 /me/refunds/apply?...。
   - 仅在兑现条件满足时显示按钮，保持 OrderCard 现有按钮堆栈。

## 交互 & 流程
- 进入 /me/refunds 直接看到自己的退款记录，状态按 AfterSaleStatus 排布。卡片支持一键跳转、取消（状态允许）、提醒联系售后。
- 点击某个批次进入详情页，看到所有 AfterSaleCase（商品、原因、申请 / 审批金额、状态）和最新的 RefundTask 日志（状态/错误）。
- 点击订单列表中的 申请退款 或 /me/refunds 的 申请退款 入口，会跳到 /me/refunds/apply 先调用 preview，确认自动带入金额/商品清单，再选择 reason 及备注后点击提交，成功后自动进入新批次的详情页。
- 取消操作会调用 cancelRefundBatch，成功后刷新列表并提示用户。

## 测试计划
- 在 pps/mall-web/tests（或靠近 eatures/refund）新建测试文件覆盖两个关键 helper：
  1. hasRefundableSubOrder（确认 PAID/WAIT_SHIP 等子单会被认为可申请，这与订单页面动作位置直接相关）。
  2. ormatAfterSaleStatusLabel（验证中英文分别对应 待审核/退款中/已退款 等状态段）。
- 利用 itest（新增 devDependency）或 Node 自带的 
ode --test 把这两个测试先写好，跑一遍确保先失败再修复，满足 TDD 要求。
- 未来可扩展到 snapshot 测试 RefundBatch 卡片的 status chips，但这次先聚焦 helper。

## 假设 & 待确认
1. 退款列表页需要展示全部历史状态而不是仅限当前可申请项（如有不同请告知）。
2. Preview 接口不会再额外要求 scope 之外的数据，如需特殊 selectedItemNos 还需补充参数来源。
3. 理由选项采用字符串 code（例如 changed_mind），可在后台校验，但前端只需要展示友好中文文案和对应 code。
