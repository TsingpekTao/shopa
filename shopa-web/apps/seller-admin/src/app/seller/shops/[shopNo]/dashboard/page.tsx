"use client";

import { useEffect, useMemo, useState } from "react";
import { Alert, Card, Col, Progress, Row, Skeleton, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useI18n } from "@shopa/ui";
import { fetchShopDashboard } from "@/features/seller-shop/api";
import { ShopDashboardResponse } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

type TaskRow = {
  key: string;
  task: string;
  owner: string;
  status: "DONE" | "RUNNING" | "TODO";
  updatedAt: string;
};

const tasks: TaskRow[] = [
  { key: "1", task: "主图优化", owner: "运营A", status: "RUNNING", updatedAt: "2026-03-28 11:12" },
  { key: "2", task: "低库存补货", owner: "仓储B", status: "TODO", updatedAt: "2026-03-28 10:03" },
  { key: "3", task: "活动价格检查", owner: "运营C", status: "DONE", updatedAt: "2026-03-27 22:45" }
];

function statusColor(status: TaskRow["status"]) {
  if (status === "DONE") return "success";
  if (status === "RUNNING") return "processing";
  return "default";
}

export default function ShopDashboardPage({ params }: { params: { shopNo: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<ShopDashboardResponse | null>(null);

  useEffect(() => {
    fetchShopDashboard(params.shopNo)
      .then((res) => setData(res))
      .finally(() => setLoading(false));
  }, [params.shopNo]);

  const statusTag = useMemo(() => {
    const code = (data?.shopStatusCode || "UNKNOWN").toUpperCase();
    const color = code === "ACTIVE" ? "success" : code === "PENDING" ? "processing" : "default";
    return <Tag color={color}>{code}</Tag>;
  }, [data?.shopStatusCode]);

  const taskColumns: ColumnsType<TaskRow> = [
    { title: isZh ? "任务" : "Task", dataIndex: "task", key: "task" },
    { title: isZh ? "负责人" : "Owner", dataIndex: "owner", key: "owner", width: 120 },
    { title: isZh ? "状态" : "Status", dataIndex: "status", key: "status", width: 120, render: (v) => <Tag color={statusColor(v)}>{v}</Tag> },
    { title: isZh ? "更新时间" : "Updated", dataIndex: "updatedAt", key: "updatedAt", width: 180 }
  ];

  if (loading) {
    return <Skeleton active paragraph={{ rows: 8 }} />;
  }

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "店铺看板" : "Shop Dashboard"}</Title>
          <Text type="secondary">
            {isZh ? "店铺编号：" : "Shop No: "}
            {data?.shopNo}
          </Text>
        </div>
        <div>{statusTag}</div>
      </header>

      {data?.partial ? (
        <Alert type="warning" showIcon message={isZh ? "部分数据不可用" : "Partial data"} description={(data.degradedFields || []).join(", ")} />
      ) : null}

      <div className="seller-kpi-grid">
        <Card>
          <Text type="secondary">{isZh ? "店铺名称" : "Shop Name"}</Text>
          <Title level={4} style={{ marginTop: 8 }}>
            {data?.shopName || "-"}
          </Title>
        </Card>
        <Card>
          <Text type="secondary">{isZh ? "商品总数" : "Total Products"}</Text>
          <Title level={4} style={{ marginTop: 8 }}>
            {data?.productSummary?.total ?? 0}
          </Title>
        </Card>
        <Card>
          <Text type="secondary">{isZh ? "低库存SKU" : "Low Stock SKU"}</Text>
          <Title level={4} style={{ marginTop: 8 }}>
            {data?.inventoryRisk?.lowStockSkuCount ?? 0}
          </Title>
        </Card>
        <Card>
          <Text type="secondary">{isZh ? "缺货SKU" : "Out of Stock SKU"}</Text>
          <Title level={4} style={{ marginTop: 8 }}>
            {data?.inventoryRisk?.outOfStockSkuCount ?? 0}
          </Title>
        </Card>
      </div>

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={12}>
          <Card title={isZh ? "商品结构" : "Product Structure"}>
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "在架" : "On shelf"}</span>
                <strong>{data?.productSummary?.onShelf ?? 0}</strong>
              </div>
              <div>
                <span>{isZh ? "下架" : "Off shelf"}</span>
                <strong>{data?.productSummary?.offShelf ?? 0}</strong>
              </div>
              <div>
                <span>{isZh ? "草稿" : "Draft"}</span>
                <strong>{data?.productSummary?.draft ?? 0}</strong>
              </div>
              <div>
                <span>{isZh ? "审核中" : "Reviewing"}</span>
                <strong>{data?.productSummary?.reviewing ?? 0}</strong>
              </div>
              <div>
                <span>{isZh ? "驳回" : "Rejected"}</span>
                <strong>{data?.productSummary?.rejected ?? 0}</strong>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={24} xl={12}>
          <Card title={isZh ? "库存风险指数" : "Inventory Risk Index"}>
            <Progress
              percent={Math.min(100, (data?.inventoryRisk?.lowStockSkuCount ?? 0) * 8 + (data?.inventoryRisk?.outOfStockSkuCount ?? 0) * 12)}
              strokeColor="#ff6a00"
            />
            <Text type="secondary">{isZh ? "指数越高说明补货优先级越高。" : "Higher score means higher replenishment priority."}</Text>
          </Card>
        </Col>
      </Row>

      <Card title={isZh ? "运营任务" : "Operations Tasks"}>
        <Table<TaskRow> rowKey="key" columns={taskColumns} dataSource={tasks} pagination={false} />
      </Card>
    </section>
  );
}
