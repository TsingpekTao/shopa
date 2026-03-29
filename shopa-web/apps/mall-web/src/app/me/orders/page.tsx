"use client";

import { useI18n } from "@shopa/ui";

const rows = [
  { id: "O20260327001", title: "主动降噪蓝牙耳机", amount: 269, statusZh: "待发货", statusEn: "Pending Shipment" },
  { id: "O20260327002", title: "轻弹缓震跑鞋", amount: 199, statusZh: "配送中", statusEn: "In Delivery" }
];

export default function OrdersPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <section className="tb-floor">
      <header className="tb-floor-header">
        <h2>{isZh ? "当前订单" : "Current Orders"}</h2>
      </header>
      <div className="tb-checkout-items">
        <table>
          <thead>
            <tr>
              <th>{isZh ? "订单号" : "Order No"}</th>
              <th>{isZh ? "商品" : "Item"}</th>
              <th>{isZh ? "金额" : "Amount"}</th>
              <th>{isZh ? "状态" : "Status"}</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.id}>
                <td>{row.id}</td>
                <td>{row.title}</td>
                <td>¥{row.amount}</td>
                <td>{isZh ? row.statusZh : row.statusEn}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
