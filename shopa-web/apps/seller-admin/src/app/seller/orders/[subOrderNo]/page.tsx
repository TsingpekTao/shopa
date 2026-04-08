"use client";

import Link from "next/link";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { Alert, Button, Card, Col, Descriptions, Empty, Row, Skeleton, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { ArrowLeftOutlined } from "@ant-design/icons";
import { getSellerOrderDetail } from "@/features/order/api";
import type {
  SellerOrderDetail,
  SellerOrderItem,
  SellerOrderMain,
  SellerOrderSub,
  SellerPaymentStatusCode,
  SellerSubOrderStatusCode
} from "@/features/order/types";

const { Title, Text, Paragraph } = Typography;

type PageProps = {
  params: {
    subOrderNo: string;
  };
};

function safeDecode(value: string): string {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function formatDateTime(value?: string): string {
  if (!value) {
    return "-";
  }
  const parsed = Date.parse(value);
  if (Number.isNaN(parsed)) {
    return value;
  }
  const date = new Date(parsed);
  const yyyy = date.getFullYear();
  const mm = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const hh = String(date.getHours()).padStart(2, "0");
  const min = String(date.getMinutes()).padStart(2, "0");
  return `${yyyy}-${mm}-${dd} ${hh}:${min}`;
}

function formatMoney(amount: number): string {
  return `CNY ${(amount / 100).toFixed(2)}`;
}

function statusColor(code: string): string {
  switch ((code || "").toUpperCase()) {
    case "PAID":
    case "SHIPPED":
    case "COMPLETED":
      return "success";
    case "PAYING":
    case "WAIT_SHIP":
    case "FULFILLING":
      return "processing";
    case "PENDING_PAY":
    case "UNPAID":
      return "warning";
    case "PAY_FAILED":
    case "CANCELED":
    case "CLOSED":
    case "REFUNDING":
    case "REFUNDED":
      return "error";
    default:
      return "default";
  }
}

function paymentStatusLabel(status: SellerPaymentStatusCode): string {
  switch (status) {
    case "UNPAID":
      return "未支付";
    case "PAYING":
      return "支付中";
    case "PAID":
      return "已支付";
    case "PAY_FAILED":
      return "支付失败";
    case "REFUNDED":
      return "已退款";
    default:
      return "未知";
  }
}

function subOrderStatusLabel(status: SellerSubOrderStatusCode): string {
  switch (status) {
    case "PENDING_PAY":
      return "待付款";
    case "PAID":
      return "已支付";
    case "WAIT_SHIP":
      return "待发货";
    case "SHIPPED":
      return "已发货";
    case "COMPLETED":
      return "已完成";
    case "CANCELED":
      return "已取消";
    case "CLOSED":
      return "已关闭";
    case "REFUNDING":
      return "退款中";
    case "REFUNDED":
      return "已退款";
    default:
      return "未知";
  }
}

function orderStatusLabel(status: string): string {
  switch ((status || "").toUpperCase()) {
    case "PENDING_PAY":
      return "待付款";
    case "PAID":
      return "已支付";
    case "FULFILLING":
      return "履约中";
    case "COMPLETED":
      return "已完成";
    case "CANCELED":
      return "已取消";
    case "CLOSED":
      return "已关闭";
    case "REFUNDING":
      return "退款中";
    case "REFUNDED":
      return "已退款";
    default:
      return "未知";
  }
}

function formatAddress(order: SellerOrderMain | null): string {
  const address = order?.address;
  if (!address) {
    return "暂无地址快照";
  }
  return [
    [address.receiverName, address.receiverPhone].filter(Boolean).join(" / "),
    [address.provinceName, address.cityName, address.districtName, address.street, address.detail].filter(Boolean).join(" "),
    address.postalCode
  ]
    .filter(Boolean)
    .join(" | ");
}

function formatSaleAttrs(value: string): string {
  const text = value.trim();
  if (!text) {
    return "-";
  }
  try {
    const parsed = JSON.parse(text) as unknown;
    if (Array.isArray(parsed)) {
      return parsed
        .map((item) => {
          if (typeof item === "string") {
            return item;
          }
          if (item && typeof item === "object") {
            const record = item as Record<string, unknown>;
            return String(record.value ?? record.attr_value ?? record.label ?? record.name ?? "").trim();
          }
          return "";
        })
        .filter(Boolean)
        .join(" / ");
    }
    if (parsed && typeof parsed === "object") {
      return Object.values(parsed as Record<string, unknown>)
        .map((item) => String(item ?? "").trim())
        .filter(Boolean)
        .join(" / ");
    }
  } catch {
    return text;
  }
  return text;
}

function pickCurrentSubOrder(detail: SellerOrderDetail, subOrderNo: string): SellerOrderSub | null {
  if (detail.subOrder?.subOrderNo === subOrderNo) {
    return detail.subOrder;
  }
  return detail.order?.subOrders.find((item) => item.subOrderNo === subOrderNo) ?? detail.subOrder ?? null;
}

export default function SellerOrderDetailPage({ params }: PageProps) {
  const subOrderNo = useMemo(() => safeDecode(params.subOrderNo || ""), [params.subOrderNo]);

  const detailQuery = useQuery({
    queryKey: ["seller-order-detail", subOrderNo],
    queryFn: () => getSellerOrderDetail(subOrderNo),
    enabled: Boolean(subOrderNo)
  });

  const detail = detailQuery.data;
  const order = detail?.order ?? null;
  const currentSubOrder = detail ? pickCurrentSubOrder(detail, subOrderNo) : null;

  const itemColumns = useMemo<ColumnsType<SellerOrderItem>>(
    () => [
      {
        title: "商品信息",
        key: "product",
        render: (_, row) => (
          <Space direction="vertical" size={0}>
            <Text strong>{row.spuTitle || row.skuName || row.skuNo || "-"}</Text>
            <Text type="secondary">{[row.skuName, row.skuNo, row.spuNo].filter(Boolean).join(" / ") || "-"}</Text>
          </Space>
        )
      },
      {
        title: "销售属性",
        dataIndex: "saleAttrsJson",
        key: "saleAttrsJson",
        render: (value: string) => <Text type="secondary">{formatSaleAttrs(value)}</Text>
      },
      {
        title: "数量",
        dataIndex: "qty",
        key: "qty",
        width: 100
      },
      {
        title: "成交单价",
        dataIndex: "salePrice",
        key: "salePrice",
        width: 140,
        render: (value: number) => formatMoney(value)
      },
      {
        title: "小计",
        key: "subtotal",
        width: 140,
        render: (_, row) => formatMoney(row.salePrice * row.qty)
      }
    ],
    []
  );

  const relatedSubOrderColumns = useMemo<ColumnsType<SellerOrderSub>>(
    () => [
      {
        title: "子订单号",
        dataIndex: "subOrderNo",
        key: "subOrderNo",
        render: (value: string) => <Text copyable>{value || "-"}</Text>
      },
      {
        title: "店铺编号",
        dataIndex: "shopNo",
        key: "shopNo",
        width: 160,
        render: (value: string) => value || "-"
      },
      {
        title: "状态",
        dataIndex: "subStatus",
        key: "subStatus",
        width: 140,
        render: (value: SellerSubOrderStatusCode) => <Tag color={statusColor(value)}>{subOrderStatusLabel(value)}</Tag>
      },
      {
        title: "应付金额",
        key: "payable",
        width: 140,
        render: (_, row) => formatMoney(row.amount.payableAmount)
      },
      {
        title: "操作",
        key: "action",
        width: 140,
        render: (_, row) =>
          row.subOrderNo ? <Link href={`/seller/orders/${encodeURIComponent(row.subOrderNo)}`}>查看详情</Link> : "-"
      }
    ],
    []
  );

  if (detailQuery.isLoading) {
    return <Skeleton active paragraph={{ rows: 12 }} />;
  }

  if (detailQuery.isError) {
    const message = detailQuery.error instanceof Error ? detailQuery.error.message : "订单加载失败";
    return (
      <section className="seller-page">
        <Alert type="error" showIcon message="订单详情加载失败" description={message} />
      </section>
    );
  }

  if (!order || !currentSubOrder) {
    return (
      <section className="seller-page">
        <Alert type="warning" showIcon message="未找到订单详情" description={subOrderNo || "-"} />
      </section>
    );
  }

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Space size={8} wrap>
            <Tag color={statusColor(currentSubOrder.subStatus)}>{subOrderStatusLabel(currentSubOrder.subStatus)}</Tag>
            <Tag color={statusColor(order.paymentStatus)}>{paymentStatusLabel(order.paymentStatus)}</Tag>
          </Space>
          <Title level={3} style={{ marginTop: 8 }}>
            订单详情
          </Title>
          <Text type="secondary">{`子订单 ${currentSubOrder.subOrderNo} / 主订单 ${order.orderNo}`}</Text>
        </div>
        <Space wrap>
          <Link href="/seller/conversations">
            <Button icon={<ArrowLeftOutlined />}>返回会话中心</Button>
          </Link>
        </Space>
      </header>

      <div className="seller-kpi-grid">
        <Card>
          <Statistic title="子订单应付金额" value={formatMoney(currentSubOrder.amount.payableAmount)} />
        </Card>
        <Card>
          <Statistic title="订单实付金额" value={formatMoney(order.amount.paidAmount)} />
        </Card>
        <Card>
          <Statistic title="积分抵扣" value={formatMoney(order.pointsDiscountAmount)} />
        </Card>
        <Card>
          <Statistic title="支付截止时间" value={formatDateTime(order.payDeadlineAt)} />
        </Card>
      </div>

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={14}>
          <Card title="订单信息">
            <Descriptions column={1} size="small" labelStyle={{ width: 160 }}>
              <Descriptions.Item label="主订单号">
                <Text copyable>{order.orderNo}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="子订单号">
                <Text copyable>{currentSubOrder.subOrderNo}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="店铺编号">{currentSubOrder.shopNo || "-"}</Descriptions.Item>
              <Descriptions.Item label="订单状态">
                <Tag color={statusColor(order.orderStatus)}>{orderStatusLabel(order.orderStatus)}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="支付状态">
                <Tag color={statusColor(order.paymentStatus)}>{paymentStatusLabel(order.paymentStatus)}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="下单时间">{formatDateTime(currentSubOrder.createdAt || order.createdAt)}</Descriptions.Item>
              <Descriptions.Item label="更新时间">{formatDateTime(currentSubOrder.updatedAt || order.updatedAt)}</Descriptions.Item>
              <Descriptions.Item label="支付时间">{formatDateTime(order.paidAt)}</Descriptions.Item>
              <Descriptions.Item label="买家备注">{currentSubOrder.buyerRemark || order.buyerRemark || "-"}</Descriptions.Item>
              <Descriptions.Item label="商家备注">{currentSubOrder.sellerRemark || "-"}</Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>

        <Col xs={24} xl={10}>
          <Card title="收货地址">
            <Paragraph style={{ marginBottom: 0 }}>{formatAddress(order)}</Paragraph>
          </Card>
          <Card title="金额明细" style={{ marginTop: 14 }}>
            <Descriptions column={1} size="small" labelStyle={{ width: 160 }}>
              <Descriptions.Item label="商品金额">{formatMoney(currentSubOrder.amount.goodsAmount)}</Descriptions.Item>
              <Descriptions.Item label="运费">{formatMoney(currentSubOrder.amount.freightAmount)}</Descriptions.Item>
              <Descriptions.Item label="优惠金额">{formatMoney(currentSubOrder.amount.discountAmount)}</Descriptions.Item>
              <Descriptions.Item label="积分抵扣">{formatMoney(currentSubOrder.amount.pointsDiscountAmount)}</Descriptions.Item>
              <Descriptions.Item label="应付金额">{formatMoney(currentSubOrder.amount.payableAmount)}</Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
      </Row>

      <Card title="商品列表">
        {currentSubOrder.items.length > 0 ? (
          <Table<SellerOrderItem>
            rowKey={(row) => row.itemNo || `${row.subOrderNo}-${row.skuNo}`}
            columns={itemColumns}
            dataSource={currentSubOrder.items}
            pagination={false}
          />
        ) : (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无商品明细" />
        )}
      </Card>

      <Card title="关联子订单">
        <Table<SellerOrderSub> rowKey={(row) => row.subOrderNo} columns={relatedSubOrderColumns} dataSource={order.subOrders} pagination={false} />
      </Card>
    </section>
  );
}
