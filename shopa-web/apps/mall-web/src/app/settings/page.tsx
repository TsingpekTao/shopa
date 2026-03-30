"use client";

import Link from "next/link";
import { useI18n } from "@shopa/ui";
import { AddressBookManager } from "@/components/address-book-manager";

export default function SettingsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  const items = isZh
    ? [
        { href: "/me/profile", title: "个人信息", desc: "修改头像、昵称、绑定信息" },
        { href: "/me/address", title: "地址管理", desc: "新增或编辑收货地址" },
        { href: "/me/orders", title: "订单管理", desc: "查看当前订单状态" },
        { href: "/me/history", title: "历史订单", desc: "查看过往购买记录" }
      ]
    : [
        { href: "/me/profile", title: "Profile", desc: "Update avatar and account info" },
        { href: "/me/address", title: "Address", desc: "Create or edit delivery addresses" },
        { href: "/me/orders", title: "Orders", desc: "Check active order status" },
        { href: "/me/history", title: "History", desc: "Review previous purchases" }
      ];

  return (
    <section className="tb-floor">
      <header className="tb-floor-header">
        <h2>{isZh ? "设置中心" : "Settings Center"}</h2>
      </header>

      <div className="tb-sec-grid">
        {items.map((item) => (
          <article className="tb-sec-card" key={item.href}>
            <h3>{item.title}</h3>
            <p className="tb-sale-progress">{item.desc}</p>
            <Link href={item.href}>{isZh ? "立即进入" : "Open"}</Link>
          </article>
        ))}
      </div>

      <div style={{ marginTop: 16 }}>
        <AddressBookManager embedded />
      </div>
    </section>
  );
}
