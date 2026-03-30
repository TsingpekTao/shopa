"use client";

import Link from "next/link";
import { useI18n } from "@shopa/ui";

const rows = [
  {
    id: "O20260327001",
    shopNo: "SHOP1003",
    spuNo: "P2004",
    titleZh: "主动降噪蓝牙耳机",
    titleEn: "Active Noise Cancelling Earbuds",
    amount: 269,
    statusZh: "待发货",
    statusEn: "Pending Shipment"
  },
  {
    id: "O20260327002",
    shopNo: "SHOP1005",
    spuNo: "P2002",
    titleZh: "轻弹缓震跑鞋",
    titleEn: "Responsive Running Shoes",
    amount: 199,
    statusZh: "配送中",
    statusEn: "In Delivery"
  }
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
              <th>{isZh ? "服务" : "Service"}</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.id}>
                <td>{row.id}</td>
                <td>{isZh ? row.titleZh : row.titleEn}</td>
                <td>¥{row.amount}</td>
                <td>{isZh ? row.statusZh : row.statusEn}</td>
                <td>
                  <Link
                    href={`/me/messages?shop_no=${encodeURIComponent(row.shopNo)}&spu_no=${encodeURIComponent(row.spuNo)}&order_no=${encodeURIComponent(row.id)}&mode=after-sale`}
                    className="tb-orders-chat-link"
                  >
                    {isZh ? "售后咨询" : "After-sale Chat"}
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
