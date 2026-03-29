"use client";

import { useI18n } from "@shopa/ui";

export default function AddressPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <section className="tb-floor">
      <header className="tb-floor-header">
        <h2>{isZh ? "收货地址" : "Address Book"}</h2>
      </header>
      <div className="tb-checkout-address-card">
        <p>{isZh ? "张三 13300001111" : "Alex 13300001111"}<span>{isZh ? "默认" : "Default"}</span></p>
        <small>{isZh ? "浙江省杭州市西湖区古墩路 88 号" : "No.88 Gudun Road, Hangzhou"}</small>
      </div>
    </section>
  );
}
