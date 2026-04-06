"use client";

import Link from "next/link";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button, Card, Space, Table, Tag, Typography } from "antd";
import { fetchProductList } from "./api";
import { CatalogProduct } from "./types";

const statusColor: Record<string, string> = {
  draft: "default",
  reviewing: "orange",
  approved: "green",
  onShelf: "blue",
  offShelf: "volcano"
};

const { Paragraph, Text } = Typography;

export function ProductList() {
  const { data = [], isLoading } = useQuery({
    queryKey: ["catalog-products"],
    queryFn: fetchProductList,
    staleTime: 60 * 1000
  });

  const summary = useMemo(
    () => ({
      total: data.length,
      onShelf: data.filter((item) => item.status === "onShelf").length,
      drafts: data.filter((item) => item.status === "draft").length
    }),
    [data]
  );

  const columns = [
    {
      title: "SPU",
      dataIndex: "spuNo",
      key: "spuNo",
      render: (value: string) => <Text strong>{value}</Text>
    },
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
      render: (value: string) => <Paragraph ellipsis={{ rows: 2, expandable: false }}>{value}</Paragraph>
    },
    {
      title: "Price",
      dataIndex: "salePrice",
      key: "salePrice",
      align: "right" as const,
      render: (value: number) => <Text aria-label={`sale price ${value}`}>¥{value.toFixed(2)}</Text>
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (value: CatalogProduct["status"]) => <Tag color={statusColor[value] ?? "default"}>{value}</Tag>
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: unknown, record: CatalogProduct) => (
        <Link href={`/seller/publish?spuNo=${record.spuNo}`}>
          <Button size="small">View</Button>
        </Link>
      )
    }
  ];

  return (
    <Card
      title="Product List"
      extra={
        <Link href="/seller/publish">
          <Button type="primary" aria-label="Create new product">New Product</Button>
        </Link>
      }
    >
      <Space direction="horizontal" size="large" style={{ marginBottom: 16 }} aria-live="polite">
        <Typography.Text strong>Total: {summary.total}</Typography.Text>
        <Typography.Text type="success">On Shelf: {summary.onShelf}</Typography.Text>
        <Typography.Text type="secondary">Drafts: {summary.drafts}</Typography.Text>
      </Space>
      <Table<CatalogProduct>
        rowKey="spuNo"
        loading={isLoading}
        pagination={{ pageSize: 10, showSizeChanger: false, responsive: true }}
        dataSource={data}
        columns={columns}
        aria-label="My product catalog"
      />
    </Card>
  );
}
