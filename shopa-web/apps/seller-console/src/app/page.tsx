"use client";

import Link from "next/link";
import { Button, Card, Col, Row, Space, Tag, Typography } from "antd";

const { Title, Paragraph } = Typography;

export default function HomePage() {
  const modules = [
    "edge-gateway",
    "iam-svc",
    "user-profile-svc",
    "seller-shop-svc",
    "media-svc",
    "catalog-svc",
    "inventory-svc"
  ];

  return (
    <div style={{ maxWidth: 1100, margin: "0 auto" }}>
      <Card
        style={{
          borderRadius: 16,
          border: "1px solid #ffd9bc",
          boxShadow: "0 16px 28px rgba(247, 127, 37, 0.15)"
        }}
      >
        <Tag color="orange">Seller Console V1</Tag>
        <Title level={2} style={{ marginTop: 12 }}>
          Multi-Service Seller Workspace
        </Title>
        <Paragraph>
          Unified seller workspace powered by edge-gateway with coordinated IAM, profile, shop,
          media, catalog and inventory modules.
        </Paragraph>
        <Space>
          <Link href="/login">
            <Button type="primary" size="large">
              Sign In
            </Button>
          </Link>
          <Link href="/seller/workbench">
            <Button size="large">Open Workbench</Button>
          </Link>
        </Space>
      </Card>

      <div style={{ marginTop: 20 }}>
        <Row gutter={[16, 16]}>
          {modules.map((name) => (
            <Col xs={24} sm={12} lg={8} key={name}>
              <Card size="small" title={name} style={{ borderRadius: 14, borderColor: "#ffe3cc" }}>
                <Paragraph type="secondary" style={{ marginBottom: 0 }}>
                  Connected through edge-gateway contracts with UI fallback and degraded handling.
                </Paragraph>
              </Card>
            </Col>
          ))}
        </Row>
      </div>
    </div>
  );
}
