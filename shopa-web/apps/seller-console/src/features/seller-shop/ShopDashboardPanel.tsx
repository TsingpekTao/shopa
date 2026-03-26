"use client";

import {Avatar, Badge, Card, Col, Progress, Row, Space, Typography} from "antd";
import {useQuery} from "@tanstack/react-query";
import {DegradedBanner} from "@shopa/ui";
import {fetchShopDashboard} from "./api";

const {Text} = Typography;

export function ShopDashboardPanel({shopNo}: {shopNo: string}) {
  const {data, isLoading} = useQuery({
    queryKey: ["shop-dashboard", shopNo],
    queryFn: () => fetchShopDashboard(shopNo),
    staleTime: 30 * 1000
  });

  return (
    <>
      {data?.partial && <DegradedBanner fields={data.degradedFields ?? []} />}
      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Card loading={isLoading} bordered={false}>
            <Space direction="vertical">
              <Text type="secondary">Shop Identity</Text>
              <Space align="center">
                <Avatar>{data?.shopName?.slice(0, 1)}</Avatar>
                <div>
                  <Text strong>{data?.shopName}</Text>
                  <br />
                  <Text type="secondary">{data?.shopNo}</Text>
                </div>
              </Space>
              <Badge color="cyan" text={`Status: ${data?.shopStatusCode ?? "N/A"}`} />
            </Space>
          </Card>
        </Col>

        <Col span={8}>
          <Card title="Inventory Risk" loading={isLoading}>
            <Space direction="vertical" size="middle">
              <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                <Text>Low stock SKUs</Text>
                <Text strong>{data?.inventoryRisk.lowStockSkuCount ?? 0}</Text>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                <Text>Out of stock</Text>
                <Text strong>{data?.inventoryRisk.outOfStockSkuCount ?? 0}</Text>
              </div>
              <Progress percent={Math.min(100, (data?.inventoryRisk.outOfStockSkuCount ?? 0) * 10)} status="active" />
            </Space>
          </Card>
        </Col>

        <Col span={8}>
          <Card title="Product Momentum" loading={isLoading}>
            <Space direction="vertical">
              <Text>Total SKUs</Text>
              <Progress type="dashboard" percent={Math.min(100, (data?.productSummary.total ?? 0) / 100)} format={() => `${data?.productSummary.total ?? 0}`} />
              <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                <Text>On shelf</Text>
                <Text>{data?.productSummary.onShelf ?? 0}</Text>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                <Text>Drafts</Text>
                <Text>{data?.productSummary.draft ?? 0}</Text>
              </div>
            </Space>
          </Card>
        </Col>
      </Row>
    </>
  );
}
