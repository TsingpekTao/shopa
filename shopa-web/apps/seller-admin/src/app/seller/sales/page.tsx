"use client";

import { useEffect, useMemo, useState } from "react";
import { Card, Col, Row, Segmented, Select, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { LineChartOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const { Title, Text } = Typography;

type RangeType = "7D" | "14D" | "30D" | "8W";
type MetricType = "GMV" | "ORDERS";

type ShopOption = {
  label: string;
  value: string;
};

type Point = {
  key: string;
  label: string;
  gmv: number;
  orders: number;
  buyers: number;
  refunds: number;
};

function pad2(value: number): string {
  return String(value).padStart(2, "0");
}

function formatDateLabel(date: Date) {
  return `${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`;
}

function generateDailySeries(days: number): Point[] {
  const now = new Date();
  const rows: Point[] = [];
  for (let i = days - 1; i >= 0; i -= 1) {
    const date = new Date(now);
    date.setDate(now.getDate() - i);
    const seed = days - i;
    const gmv = Math.max(800, Math.round(1800 + seed * 120 + Math.sin(seed * 1.3) * 260 + (seed % 3) * 90));
    const orders = Math.max(10, Math.round(gmv / 88 + (seed % 5)));
    const buyers = Math.max(8, orders - 2 - (seed % 3));
    const refunds = Math.max(0, Math.round(orders * 0.025 + (seed % 2 === 0 ? 1 : 0)));

    rows.push({
      key: `D-${i}`,
      label: formatDateLabel(date),
      gmv,
      orders,
      buyers,
      refunds
    });
  }
  return rows;
}

function generateWeeklySeries(weeks: number): Point[] {
  const rows: Point[] = [];
  for (let i = weeks - 1; i >= 0; i -= 1) {
    const seed = weeks - i;
    const gmv = Math.max(5000, Math.round(9300 + seed * 460 + Math.sin(seed * 0.8) * 900 + (seed % 2) * 520));
    const orders = Math.max(60, Math.round(gmv / 92 + (seed % 8)));
    const buyers = Math.max(50, orders - 8 - (seed % 5));
    const refunds = Math.max(1, Math.round(orders * 0.02 + (seed % 2)));

    rows.push({
      key: `W-${i}`,
      label: `W${seed}`,
      gmv,
      orders,
      buyers,
      refunds
    });
  }
  return rows;
}

function toCurrency(value: number) {
  return `¥${value.toLocaleString("zh-CN")}`;
}

export default function SalesAnalyticsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [range, setRange] = useState<RangeType>("7D");
  const [metric, setMetric] = useState<MetricType>("GMV");
  const [shopOptions, setShopOptions] = useState<ShopOption[]>([{ label: isZh ? "全部店铺" : "All Shops", value: "ALL" }]);
  const [selectedShop, setSelectedShop] = useState<string>("ALL");

  useEffect(() => {
    fetchSellerWorkbench()
      .then((res) => {
        const options: ShopOption[] = [
          { label: isZh ? "全部店铺" : "All Shops", value: "ALL" },
          ...(res.shops || []).map((item: any) => ({
            label: String(item.shopDisplayName || item.shopName || item.shopNo || "Shop"),
            value: String(item.shopNo || item.shopDisplayName || "UNKNOWN")
          }))
        ];
        setShopOptions(options);
      })
      .catch(() => {
        setShopOptions([{ label: isZh ? "全部店铺" : "All Shops", value: "ALL" }]);
      });
  }, [isZh]);

  const points = useMemo<Point[]>(() => {
    if (range === "8W") {
      return generateWeeklySeries(8);
    }
    if (range === "14D") {
      return generateDailySeries(14);
    }
    if (range === "30D") {
      return generateDailySeries(30);
    }
    return generateDailySeries(7);
  }, [range, selectedShop]);

  const summary = useMemo(() => {
    const totalGmv = points.reduce((sum, item) => sum + item.gmv, 0);
    const totalOrders = points.reduce((sum, item) => sum + item.orders, 0);
    const totalBuyers = points.reduce((sum, item) => sum + item.buyers, 0);
    const totalRefunds = points.reduce((sum, item) => sum + item.refunds, 0);
    const avgOrderValue = totalOrders > 0 ? totalGmv / totalOrders : 0;
    const refundRate = totalOrders > 0 ? (totalRefunds / totalOrders) * 100 : 0;

    return { totalGmv, totalOrders, totalBuyers, totalRefunds, avgOrderValue, refundRate };
  }, [points]);

  const maxMetric = useMemo(() => {
    return Math.max(...points.map((item) => (metric === "GMV" ? item.gmv : item.orders)), 1);
  }, [metric, points]);

  const columns: ColumnsType<Point> = useMemo(
    () => [
      { title: isZh ? "周期" : "Period", dataIndex: "label", key: "label", width: 120 },
      {
        title: "GMV",
        dataIndex: "gmv",
        key: "gmv",
        width: 160,
        render: (value: number) => toCurrency(value)
      },
      { title: isZh ? "订单数" : "Orders", dataIndex: "orders", key: "orders", width: 120 },
      { title: isZh ? "付款买家" : "Buyers", dataIndex: "buyers", key: "buyers", width: 120 },
      { title: isZh ? "退款单" : "Refunds", dataIndex: "refunds", key: "refunds", width: 120 },
      {
        title: isZh ? "退款率" : "Refund Rate",
        key: "refundRate",
        width: 120,
        render: (_, row) => {
          const rate = row.orders > 0 ? ((row.refunds / row.orders) * 100).toFixed(2) : "0.00";
          return `${rate}%`;
        }
      }
    ],
    [isZh]
  );

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "销售分析" : "Sales Analytics"}</Title>
          <Text type="secondary">{isZh ? "查看最近几天或几周的销量、GMV、退款趋势。" : "Inspect recent sales, GMV and refund trends by day or week."}</Text>
        </div>

        <Space wrap>
          <Select value={selectedShop} onChange={setSelectedShop} options={shopOptions} style={{ minWidth: 200 }} />
          <Segmented
            value={range}
            onChange={(value) => setRange(value as RangeType)}
            options={[
              { label: "7D", value: "7D" },
              { label: "14D", value: "14D" },
              { label: "30D", value: "30D" },
              { label: "8W", value: "8W" }
            ]}
          />
          <Segmented
            value={metric}
            onChange={(value) => setMetric(value as MetricType)}
            options={[
              { label: "GMV", value: "GMV" },
              { label: isZh ? "订单" : "Orders", value: "ORDERS" }
            ]}
          />
        </Space>
      </header>

      <div className="seller-kpi-grid">
        <Card>
          <Statistic title="GMV" value={summary.totalGmv} prefix="¥" precision={0} />
        </Card>
        <Card>
          <Statistic title={isZh ? "订单数" : "Orders"} value={summary.totalOrders} />
        </Card>
        <Card>
          <Statistic title={isZh ? "付款买家" : "Buyers"} value={summary.totalBuyers} />
        </Card>
        <Card>
          <Statistic title={isZh ? "客单价" : "AOV"} value={summary.avgOrderValue} prefix="¥" precision={2} />
        </Card>
      </div>

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={16}>
          <Card title={isZh ? "趋势图" : "Trend"} extra={<Tag icon={<LineChartOutlined />}>{metric}</Tag>}>
            <div className="seller-sales-chart">
              {points.map((item) => {
                const value = metric === "GMV" ? item.gmv : item.orders;
                const height = Math.max(8, Math.round((value / maxMetric) * 100));
                return (
                  <div key={item.key} className="seller-sales-col">
                    <div className="seller-sales-bar" style={{ height: `${height}%` }} />
                    <span>{item.label}</span>
                  </div>
                );
              })}
            </div>
          </Card>
        </Col>

        <Col xs={24} xl={8}>
          <Card title={isZh ? "关键指标" : "Key Metrics"} className="seller-side-card">
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "总退款单" : "Total Refunds"}</span>
                <strong>{summary.totalRefunds}</strong>
              </div>
              <div>
                <span>{isZh ? "退款率" : "Refund Rate"}</span>
                <strong>{summary.refundRate.toFixed(2)}%</strong>
              </div>
              <div>
                <span>{isZh ? "平均每日订单" : "Avg Orders / Period"}</span>
                <strong>{points.length > 0 ? (summary.totalOrders / points.length).toFixed(1) : "0"}</strong>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Card title={isZh ? "明细数据" : "Detailed Data"}>
        <Table<Point> rowKey={(record) => record.key} columns={columns} dataSource={points} pagination={{ pageSize: 12, showSizeChanger: false }} />
      </Card>
    </section>
  );
}
