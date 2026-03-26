"use client";

import { Card, Col, Row, Statistic, Typography } from "antd";
import { useQuery } from "@tanstack/react-query";
import { fetchInventorySummary } from "./api";

const { Text } = Typography;

type Props = {
  compact?: boolean;
};

export function InventorySection({ compact = false }: Props) {
  const { data, isLoading } = useQuery({
    queryKey: ["inventory-summary"],
    queryFn: fetchInventorySummary,
    staleTime: 30 * 1000
  });

  return (
    <Card title={compact ? "Inventory Snapshot" : "Inventory Summary"} loading={isLoading}>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={8}>
          <Statistic title="Total SKUs" value={data?.totalSkus ?? 0} />
        </Col>
        <Col xs={24} sm={8}>
          <Statistic title="Low Stock" value={data?.lowStockCount ?? 0} />
        </Col>
        <Col xs={24} sm={8}>
          <Statistic title="Out Of Stock" value={data?.outOfStockCount ?? 0} />
        </Col>
      </Row>
      <Text type="secondary">Last updated: {data?.lastUpdated ?? "N/A"}</Text>
    </Card>
  );
}
