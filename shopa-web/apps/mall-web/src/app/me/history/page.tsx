"use client";

import { useI18n } from "@shopa/ui";

const rows = [
  { id: "O20260301011", title: "家用体脂秤", amount: 99, date: "2026-03-01" },
  { id: "O20260218007", title: "防蓝光护眼台灯", amount: 129, date: "2026-02-18" },
  { id: "O20260202003", title: "纯棉四件套", amount: 239, date: "2026-02-02" }
];

export default function HistoryPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <section className="tb-floor">
      <header className="tb-floor-header">
        <h2>{isZh ? "历史订单" : "Order History"}</h2>
      </header>
      <div className="tb-checkout-items">
        <table>
          <thead>
            <tr>
              <th>{isZh ? "订单号" : "Order No"}</th>
              <th>{isZh ? "商品" : "Item"}</th>
              <th>{isZh ? "金额" : "Amount"}</th>
              <th>{isZh ? "下单时间" : "Created At"}</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.id}>
                <td>{row.id}</td>
                <td>{row.title}</td>
                <td>¥{row.amount}</td>
                <td>{row.date}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
