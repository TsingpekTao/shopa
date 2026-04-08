"use client";

import { useEffect, useMemo, useState } from "react";
import { Alert, Button, Card, Col, Empty, Row, Segmented, Select, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { BarChartOutlined, LineChartOutlined, PieChartOutlined, ReloadOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";
import {
  fetchSellerSalesAnalytics,
  formatSellerSalesRefundRate,
  getSellerSalesMetricValue,
  getSellerStatDisplayValue,
  hasSellerDegradedField,
  SellerSalesMetricType
} from "@/features/seller-shop/sales";
import {
  buildSellerSalesChartPoints,
  buildSellerSalesLineAreaPath,
  buildSellerSalesLinePath,
  buildSellerSalesLinePoints,
  buildSellerSalesPieSlices,
  SellerSalesChartMode
} from "@/features/seller-shop/sales-chart";
import { SellerSalesAnalyticsResponse, SellerSalesBucket, SellerSalesRange } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

type ShopOption = {
  label: string;
  value: string;
};

const CHART_WIDTH = 760;
const CHART_HEIGHT = 280;
const PIE_COLORS = ["#ff6a00", "#ff934d", "#ffb266", "#ffcc85", "#ffdfad", "#f6b37e"];

function toCurrency(value: number) {
  return `¥${(Number(value || 0) / 100).toLocaleString("zh-CN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`;
}

function toCompactMetricValue(metric: SellerSalesMetricType, value: number) {
  if (metric === "ORDERS") {
    return `${Math.round(value)}`;
  }
  const currency = Number(value || 0) / 100;
  if (currency >= 10000) {
    return `¥${(currency / 10000).toFixed(1)}w`;
  }
  if (currency >= 1000) {
    return `¥${currency.toFixed(0)}`;
  }
  return `¥${currency.toFixed(currency >= 100 ? 0 : 2)}`;
}

function getMetricDisplay(metric: SellerSalesMetricType, value: number) {
  return metric === "ORDERS" ? Number(value || 0).toLocaleString("zh-CN") : toCurrency(value);
}

function getMetricLabel(metric: SellerSalesMetricType, isZh: boolean) {
  return metric === "ORDERS" ? (isZh ? "支付订单数" : "Paid Orders") : "GMV";
}

function getChartModeLabel(mode: SellerSalesChartMode, isZh: boolean) {
  switch (mode) {
    case "LINE":
      return isZh ? "折线图" : "Line";
    case "PIE":
      return isZh ? "环形图" : "Donut";
    default:
      return isZh ? "柱状图" : "Bar";
  }
}

function getSalesDegradedMessages(fields: string[] | undefined, isZh: boolean) {
  const messages: string[] = [];
  if (hasSellerDegradedField(fields, "sales_orders")) {
    messages.push(isZh ? "支付订单数据暂时不可用" : "Paid-order data is temporarily unavailable");
  }
  if (hasSellerDegradedField(fields, "sales_refunds")) {
    messages.push(isZh ? "退款数据暂时不可用" : "Refund data is temporarily unavailable");
  }
  return messages;
}

function polarToCartesian(cx: number, cy: number, radius: number, angle: number) {
  return {
    x: cx + Math.cos(angle) * radius,
    y: cy + Math.sin(angle) * radius
  };
}

function describeDonutSlice(
  cx: number,
  cy: number,
  outerRadius: number,
  innerRadius: number,
  startAngle: number,
  endAngle: number
) {
  const outerStart = polarToCartesian(cx, cy, outerRadius, startAngle);
  const outerEnd = polarToCartesian(cx, cy, outerRadius, endAngle);
  const innerEnd = polarToCartesian(cx, cy, innerRadius, endAngle);
  const innerStart = polarToCartesian(cx, cy, innerRadius, startAngle);
  const largeArcFlag = endAngle - startAngle > Math.PI ? 1 : 0;

  return [
    `M ${outerStart.x} ${outerStart.y}`,
    `A ${outerRadius} ${outerRadius} 0 ${largeArcFlag} 1 ${outerEnd.x} ${outerEnd.y}`,
    `L ${innerEnd.x} ${innerEnd.y}`,
    `A ${innerRadius} ${innerRadius} 0 ${largeArcFlag} 0 ${innerStart.x} ${innerStart.y}`,
    "Z"
  ].join(" ");
}

export default function SalesAnalyticsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [range, setRange] = useState<SellerSalesRange>("7D");
  const [metric, setMetric] = useState<SellerSalesMetricType>("GMV");
  const [chartMode, setChartMode] = useState<SellerSalesChartMode>("BAR");
  const [shopOptions, setShopOptions] = useState<ShopOption[]>([{ label: isZh ? "全部店铺" : "All Shops", value: "ALL" }]);
  const [selectedShop, setSelectedShop] = useState<string>("ALL");
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<SellerSalesAnalyticsResponse | null>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    let active = true;

    fetchSellerWorkbench()
      .then((res) => {
        if (!active) {
          return;
        }
        const options: ShopOption[] = [
          { label: isZh ? "全部店铺" : "All Shops", value: "ALL" },
          ...(res.shops || []).map((item: any) => ({
            label: String(item.shopDisplayName || item.shopName || item.shopNo || (isZh ? "店铺" : "Shop")),
            value: String(item.shopNo || "").trim()
          }))
        ].filter((item) => item.value || item.value === "ALL");
        setShopOptions(options);
      })
      .catch(() => {
        if (active) {
          setShopOptions([{ label: isZh ? "全部店铺" : "All Shops", value: "ALL" }]);
        }
      });

    return () => {
      active = false;
    };
  }, [isZh]);

  useEffect(() => {
    let active = true;

    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const res = await fetchSellerSalesAnalytics({
          range,
          shopNo: selectedShop === "ALL" ? "" : selectedShop
        });
        if (!active) {
          return;
        }
        setData(res);
      } catch (err) {
        if (!active) {
          return;
        }
        setData(null);
        setError(err instanceof Error && err.message ? err.message : isZh ? "销售分析暂时无法加载" : "Sales analytics is unavailable");
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    void load();
    return () => {
      active = false;
    };
  }, [isZh, range, selectedShop]);

  const degradedFields = data?.degradedFields;
  const degradedMessages = useMemo(() => getSalesDegradedMessages(degradedFields, isZh), [degradedFields, isZh]);
  const ordersDegraded = hasSellerDegradedField(degradedFields, "sales_orders");
  const refundsDegraded = hasSellerDegradedField(degradedFields, "sales_refunds");
  const series = data?.series ?? [];

  const points = useMemo(() => buildSellerSalesChartPoints(series, metric), [metric, series]);
  const linePoints = useMemo(() => buildSellerSalesLinePoints(points, CHART_WIDTH, CHART_HEIGHT), [points]);
  const linePath = useMemo(() => buildSellerSalesLinePath(points, CHART_WIDTH, CHART_HEIGHT), [points]);
  const lineAreaPath = useMemo(() => buildSellerSalesLineAreaPath(points, CHART_WIDTH, CHART_HEIGHT), [points]);
  const pieSlices = useMemo(() => buildSellerSalesPieSlices(series, metric), [metric, series]);

  const peakPoint = useMemo(() => {
    if (!points.length) {
      return null;
    }
    return points.reduce((best, point) => (point.value > best.value ? point : best), points[0]);
  }, [points]);

  const averageOrdersPerBucket = useMemo(() => {
    if (ordersDegraded || !data?.summary || series.length === 0) {
      return "--";
    }
    return (data.summary.paidOrderCount / series.length).toFixed(1);
  }, [data?.summary, ordersDegraded, series.length]);

  const chartSummaryValue = useMemo(() => {
    const totalValue = points.reduce((sum, point) => sum + point.value, 0);
    return getMetricDisplay(metric, totalValue);
  }, [metric, points]);

  const columns: ColumnsType<SellerSalesBucket> = useMemo(
    () => [
      { title: isZh ? "周期" : "Period", dataIndex: "bucketLabel", key: "bucketLabel", width: 140 },
      {
        title: "GMV",
        dataIndex: "gmv",
        key: "gmv",
        width: 160,
        render: (value: number) => (ordersDegraded ? "--" : toCurrency(value))
      },
      {
        title: isZh ? "支付订单数" : "Paid Orders",
        dataIndex: "paidOrderCount",
        key: "paidOrderCount",
        width: 120,
        render: (value: number) => (ordersDegraded ? "--" : value)
      },
      {
        title: isZh ? "付款买家数" : "Paid Buyers",
        dataIndex: "paidBuyerCount",
        key: "paidBuyerCount",
        width: 120,
        render: (value: number) => (ordersDegraded ? "--" : value)
      },
      {
        title: isZh ? "退款金额" : "Refund Amount",
        dataIndex: "refundAmount",
        key: "refundAmount",
        width: 150,
        render: (value: number) => (refundsDegraded ? "--" : toCurrency(value))
      },
      {
        title: isZh ? "退款率" : "Refund Rate",
        key: "refundRate",
        width: 120,
        render: (_, row) => (refundsDegraded ? "--" : formatSellerSalesRefundRate(row.refundRate))
      }
    ],
    [isZh, ordersDegraded, refundsDegraded]
  );

  const barChart = (
    <div className="seller-sales-bar-chart">
      {points.map((point) => (
        <div key={point.key} className="seller-sales-bar-col">
          <strong>{toCompactMetricValue(metric, point.value)}</strong>
          <div className="seller-sales-bar-track">
            <div className="seller-sales-bar-fill" style={{ height: `${Math.max(point.ratio * 100, point.value > 0 ? 8 : 0)}%` }} />
          </div>
          <span>{point.label}</span>
        </div>
      ))}
    </div>
  );

  const lineChart = (
    <div className="seller-sales-line-wrap">
      <svg viewBox={`0 0 ${CHART_WIDTH} ${CHART_HEIGHT}`} className="seller-sales-line-svg" role="img" aria-label="sales trend chart">
        {[0, 1, 2, 3].map((index) => {
          const y = 18 + ((CHART_HEIGHT - 36) / 3) * index;
          return <line key={index} x1="22" x2={String(CHART_WIDTH - 22)} y1={String(y)} y2={String(y)} className="seller-sales-line-grid" />;
        })}
        {lineAreaPath ? <path d={lineAreaPath} className="seller-sales-line-area" /> : null}
        {linePath ? <path d={linePath} className="seller-sales-line-path" /> : null}
        {linePoints.map((point) => (
          <g key={point.key}>
            <circle cx={point.x} cy={point.y} r="5" className="seller-sales-line-dot" />
            {point.showLabel ? (
              <text x={point.x} y={point.y - 12} textAnchor="middle" className="seller-sales-line-value">
                {toCompactMetricValue(metric, point.value)}
              </text>
            ) : null}
          </g>
        ))}
      </svg>
      <div className="seller-sales-line-axis">
        {points.map((point) => (
          <span key={point.key}>{point.label}</span>
        ))}
      </div>
    </div>
  );

  const pieChart = pieSlices.length ? (
    <div className="seller-sales-pie-layout">
      <svg viewBox="0 0 280 280" className="seller-sales-pie-svg" role="img" aria-label="sales composition chart">
        {pieSlices.map((slice, index) => (
          <path
            key={slice.key}
            d={describeDonutSlice(140, 140, 106, 60, slice.startAngle, slice.endAngle)}
            fill={PIE_COLORS[index % PIE_COLORS.length]}
          />
        ))}
        <circle cx="140" cy="140" r="54" fill="#fffaf6" />
        <text x="140" y="128" textAnchor="middle" className="seller-sales-pie-center-label">
          {getMetricLabel(metric, isZh)}
        </text>
        <text x="140" y="154" textAnchor="middle" className="seller-sales-pie-center-value">
          {chartSummaryValue}
        </text>
      </svg>
      <div className="seller-sales-pie-legend">
        {pieSlices.map((slice, index) => (
          <div key={slice.key} className="seller-sales-pie-legend-item">
            <span className="seller-sales-pie-dot" style={{ background: PIE_COLORS[index % PIE_COLORS.length] }} />
            <div>
              <strong>{slice.label}</strong>
              <small>{getMetricDisplay(metric, slice.value)}</small>
            </div>
            <em>{(slice.ratio * 100).toFixed(1)}%</em>
          </div>
        ))}
      </div>
    </div>
  ) : (
    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "当前筛选下没有可展示的数据" : "No data for the current filters"} />
  );

  const chartContent = ordersDegraded ? (
    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "支付订单数据暂时不可用" : "Paid-order data unavailable"} />
  ) : !series.length ? (
    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "当前筛选下没有数据" : "No data for the current filters"} />
  ) : chartMode === "LINE" ? (
    lineChart
  ) : chartMode === "PIE" ? (
    pieChart
  ) : (
    barChart
  );

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "销售分析" : "Sales Analytics"}</Title>
          <Text type="secondary">
            {isZh ? "查看支付成交、买家转化和退款趋势。" : "Track paid sales, buyer conversion, and refund trends."}
          </Text>
        </div>

        <Space wrap>
          <Select value={selectedShop} onChange={setSelectedShop} options={shopOptions} style={{ minWidth: 220 }} />
          <Segmented
            value={range}
            onChange={(value) => setRange(value as SellerSalesRange)}
            options={[
              { label: "7D", value: "7D" },
              { label: "14D", value: "14D" },
              { label: "30D", value: "30D" },
              { label: "8W", value: "8W" }
            ]}
          />
          <Segmented
            value={metric}
            onChange={(value) => setMetric(value as SellerSalesMetricType)}
            options={[
              { label: "GMV", value: "GMV" },
              { label: isZh ? "订单" : "Orders", value: "ORDERS" }
            ]}
          />
        </Space>
      </header>

      {degradedMessages.length > 0 ? (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 14 }}
          message={isZh ? "部分报表字段暂不可用" : "Some report fields are temporarily unavailable"}
          description={
            <Space wrap>
              {degradedMessages.map((message) => (
                <Tag key={message}>{message}</Tag>
              ))}
            </Space>
          }
        />
      ) : null}

      {error && !loading && !data ? (
        <Card style={{ marginBottom: 14 }}>
          <Space direction="vertical" size={8}>
            <Text>{error}</Text>
            <Button icon={<ReloadOutlined />} onClick={() => window.location.reload()}>
              {isZh ? "重新加载" : "Reload"}
            </Button>
          </Space>
        </Card>
      ) : null}

      <div className="seller-kpi-grid">
        <Card loading={loading}>
          <Statistic title="GMV" value={ordersDegraded ? "--" : toCurrency(data?.summary?.gmv ?? 0)} />
        </Card>
        <Card loading={loading}>
          <Statistic title={isZh ? "支付订单数" : "Paid Orders"} value={getSellerStatDisplayValue(data?.summary?.paidOrderCount, ordersDegraded)} />
        </Card>
        <Card loading={loading}>
          <Statistic title={isZh ? "付款买家数" : "Paid Buyers"} value={getSellerStatDisplayValue(data?.summary?.paidBuyerCount, ordersDegraded)} />
        </Card>
        <Card loading={loading}>
          <Statistic title={isZh ? "客单价" : "AOV"} value={ordersDegraded ? "--" : toCurrency(data?.summary?.avgOrderValue ?? 0)} />
        </Card>
      </div>

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={16}>
          <Card
            title={isZh ? "趋势图" : "Trend"}
            extra={
              <Segmented
                value={chartMode}
                onChange={(value) => setChartMode(value as SellerSalesChartMode)}
                options={[
                  { label: <span><BarChartOutlined /> {isZh ? "柱状" : "Bar"}</span>, value: "BAR" },
                  { label: <span><LineChartOutlined /> {isZh ? "折线" : "Line"}</span>, value: "LINE" },
                  { label: <span><PieChartOutlined /> {isZh ? "环形" : "Donut"}</span>, value: "PIE" }
                ]}
              />
            }
            loading={loading}
          >
            <div className="seller-sales-panel-head">
              <div className="seller-sales-panel-copy">
                <strong>{getMetricLabel(metric, isZh)}</strong>
                <span>{isZh ? `当前使用${getChartModeLabel(chartMode, isZh)}展示` : `Displaying ${getChartModeLabel(chartMode, isZh)} view`}</span>
              </div>
              <div className="seller-sales-panel-stats">
                <div>
                  <span>{isZh ? "总览" : "Total"}</span>
                  <strong>{ordersDegraded ? "--" : chartSummaryValue}</strong>
                </div>
                <div>
                  <span>{isZh ? "峰值周期" : "Peak Period"}</span>
                  <strong>{ordersDegraded || !peakPoint ? "--" : `${peakPoint.label} · ${getMetricDisplay(metric, peakPoint.value)}`}</strong>
                </div>
              </div>
            </div>
            <div className="seller-sales-chart-shell">{chartContent}</div>
          </Card>
        </Col>

        <Col xs={24} xl={8}>
          <Card title={isZh ? "关键指标" : "Key Metrics"} className="seller-side-card" loading={loading}>
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "退款金额" : "Refund Amount"}</span>
                <strong>{refundsDegraded ? "--" : toCurrency(data?.summary?.refundAmount ?? 0)}</strong>
              </div>
              <div>
                <span>{isZh ? "退款率" : "Refund Rate"}</span>
                <strong>{refundsDegraded ? "--" : formatSellerSalesRefundRate(data?.summary?.refundRate)}</strong>
              </div>
              <div>
                <span>{isZh ? "平均每周期订单数" : "Avg Orders / Bucket"}</span>
                <strong>{averageOrdersPerBucket}</strong>
              </div>
              <div>
                <span>{isZh ? "当前图表口径" : "Current Metric"}</span>
                <strong>{getMetricLabel(metric, isZh)}</strong>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Card title={isZh ? "周期明细" : "Details"} loading={loading}>
        <Table<SellerSalesBucket>
          rowKey={(record) => record.bucketKey}
          columns={columns}
          dataSource={series}
          pagination={{ pageSize: 12, showSizeChanger: false }}
          locale={{ emptyText: isZh ? "暂无数据" : "No data" }}
        />
      </Card>
    </section>
  );
}
