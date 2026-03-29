"use client";

import Link from "next/link";
import { Col, Row } from "antd";
import { FeatureCard, useI18n } from "@shopa/ui";

export default function SellerHomePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const cards = [
    {
      title: isZh ? "工作台" : "Workbench",
      description: isZh ? "查看商家指标、申请状态和店铺总览。" : "View seller KPI, applications and shop status.",
      href: "/seller/workbench"
    },
    {
      title: isZh ? "销售分析" : "Sales Analytics",
      description: isZh ? "按天或按周查看最近销售趋势与核心指标。" : "Track recent sales trends by day or week.",
      href: "/seller/sales"
    },
    {
      title: isZh ? "发布商品" : "Publish Product",
      description: isZh ? "在一个流程中完成素材、草稿、库存和提审。" : "One-page flow: media, draft, inventory and review submit.",
      href: "/seller/publish"
    },
    {
      title: isZh ? "库存管理" : "Inventory",
      description: isZh ? "按 SKU 调整库存并查看失败明细。" : "Adjust SKU stock with failure details.",
      href: "/seller/inventory"
    },
    {
      title: isZh ? "我的" : "My Account",
      description: isZh ? "管理账号安全，执行退出登录与账号注销操作。" : "Manage account security and sign-out actions.",
      href: "/seller/me"
    }
  ];

  return (
    <Row gutter={[24, 24]}>
      {cards.map((item) => (
        <Col xs={24} md={12} lg={8} key={item.href}>
          <Link href={item.href}>
            <FeatureCard title={item.title} description={item.description} />
          </Link>
        </Col>
      ))}
    </Row>
  );
}
