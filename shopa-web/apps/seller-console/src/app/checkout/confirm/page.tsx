export default function CheckoutConfirmPage() {
  return (
    <main className="tb-checkout-page">
      <h1>确认订单</h1>
      <section className="tb-checkout-address">
        <h2>收货地址</h2>
        <div className="tb-checkout-address-card">
          <p>
            张三 13364027679
            <span>默认</span>
          </p>
          <p>北京市 朝阳区 酒仙桥街道 88 号 A 座 1201</p>
          <small>订单提交后将固化地址快照，后续资料变更不会影响本订单。</small>
        </div>
      </section>

      <section className="tb-checkout-items">
        <h2>商品清单</h2>
        <table>
          <thead>
            <tr>
              <th>商品</th>
              <th>单价</th>
              <th>数量</th>
              <th>小计</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>轻便跑鞋（黑色 42）</td>
              <td>￥189</td>
              <td>1</td>
              <td>￥189</td>
            </tr>
            <tr>
              <td>无线耳机（白色 标准版）</td>
              <td>￥139</td>
              <td>2</td>
              <td>￥278</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section className="tb-checkout-submit">
        <p>
          应付总额：<strong>￥467</strong>
        </p>
        <button type="button">提交订单</button>
      </section>
    </main>
  );
}
