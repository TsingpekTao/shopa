"use client";

import { useI18n } from "@shopa/ui";

function InfoPage({ title, content }: { title: { zh: string; en: string }; content: { zh: string; en: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  return (
    <section className="tb-floor">
      <header className="tb-floor-header">
        <h2>{isZh ? title.zh : title.en}</h2>
      </header>
      <p>{isZh ? content.zh : content.en}</p>
    </section>
  );
}

export function AboutContent() {
  return (
    <InfoPage
      title={{ zh: "关于我们", en: "About Us" }}
      content={{ zh: "Shopa 致力于打造高质量、可信赖的电商交易平台。", en: "Shopa is dedicated to building a reliable marketplace for quality shopping." }}
    />
  );
}

export function ContactContent() {
  return (
    <InfoPage
      title={{ zh: "联系我们", en: "Contact" }}
      content={{ zh: "客服热线：400-800-2026，邮箱：support@shopa.com", en: "Hotline: 400-800-2026, Email: support@shopa.com" }}
    />
  );
}

export function ServiceContent() {
  return (
    <InfoPage
      title={{ zh: "商家服务", en: "Merchant Service" }}
      content={{ zh: "为商家提供运营、成长与数据支持，帮助店铺稳定增长。", en: "We provide operation, growth and data support for merchants." }}
    />
  );
}

export function HelpContent() {
  return (
    <InfoPage
      title={{ zh: "帮助中心", en: "Help Center" }}
      content={{ zh: "可在这里查看支付、配送、售后与账号相关常见问题。", en: "Find FAQs about payment, shipping, after-sales and account security." }}
    />
  );
}
