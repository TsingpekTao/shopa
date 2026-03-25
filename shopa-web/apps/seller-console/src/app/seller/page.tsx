"use client";

import Link from "next/link";
import { Col, Row } from "antd";
import { FeatureCard } from "@shopa/ui";

const cards = [
  {
    title: "Workbench",
    description: "View seller KPI, applications and shop status.",
    href: "/seller/workbench"
  },
  {
    title: "Publish Product",
    description: "One-page flow: media, draft, inventory and review submit.",
    href: "/seller/publish"
  },
  {
    title: "Inventory",
    description: "Adjust SKU stock with failure details.",
    href: "/seller/inventory"
  }
];

export default function SellerHomePage() {
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
