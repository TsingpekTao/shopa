"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Alert, Button, Card, Col, Progress, Row, Skeleton, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { ArrowRightOutlined, PlusOutlined, ShopOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";
import { SellerApplicationListItem, SellerWorkbenchResponse } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

type ShopRow = {
  shopNo: string;
  shopDisplayName: string;
  shopStatusCode: string;
  buyerVisible: boolean;
};

function statusColor(status: string) {
  const code = (status || "").toUpperCase();
  if (code === "ACTIVE") return "success";
  if (code === "PENDING") return "processing";
  if (code === "DISABLED") return "error";
  return "default";
}

function appStatusColor(status: string) {
  const code = (status || "").toUpperCase();
  if (code === "APPROVED") return "success";
  if (code === "SUBMITTED" || code === "REVIEWING") return "processing";
  if (code === "REJECTED" || code === "SYSTEM_REJECTED") return "error";
  return "default";
}

export default function WorkbenchPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<SellerWorkbenchResponse | null>(null);

  useEffect(() => {
    fetchSellerWorkbench()
      .then((res) => setData(res))
      .finally(() => setLoading(false));
  }, []);

  const shopRows = useMemo<ShopRow[]>(() => {
    return (data?.shops || []).map((item: any) => ({
      shopNo: String(item.shopNo || ""),
      shopDisplayName: String(item.shopDisplayName || item.shopName || "-"),
      shopStatusCode: String(item.shopStatusCode || "UNKNOWN"),
      buyerVisible: Boolean(item.buyerVisible)
    }));
  }, [data?.shops]);

  const applicationRows = useMemo<SellerApplicationListItem[]>(() => data?.latestApplications || [], [data?.latestApplications]);

  const shopColumns = useMemo<ColumnsType<ShopRow>>(
    () => [
      {
        title: isZh ? "店铺" : "Shop",
        dataIndex: "shopDisplayName",
        key: "shopDisplayName",
        render: (_, row) => (
          <Space direction="vertical" size={0}>
            <Text strong>{row.shopDisplayName}</Text>
            <Text type="secondary">{row.shopNo}</Text>
          </Space>
        )
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "shopStatusCode",
        key: "shopStatusCode",
        width: 120,
        render: (value: string) => <Tag color={statusColor(value)}>{value}</Tag>
      },
      {
        title: isZh ? "买家可见" : "Buyer Visible",
        dataIndex: "buyerVisible",
        key: "buyerVisible",
        width: 120,
        render: (value: boolean) => <Tag color={value ? "success" : "default"}>{value ? (isZh ? "是" : "Yes") : isZh ? "否" : "No"}</Tag>
      },
      {
        title: isZh ? "操作" : "Action",
        key: "action",
        width: 140,
        render: (_, row) => (
          <Link href={`/seller/shops/${row.shopNo}/dashboard`}>
            {isZh ? "进入看板" : "Open Dashboard"}
          </Link>
        )
      }
    ],
    [isZh]
  );

  const appColumns = useMemo<ColumnsType<SellerApplicationListItem>>(
    () => [
      {
        title: isZh ? "申请单号" : "Application No",
        dataIndex: "applicationNo",
        key: "applicationNo",
        render: (value: string) => <Text copyable>{value}</Text>
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "applicationStatusCode",
        key: "applicationStatusCode",
        width: 140,
        render: (value: string) => <Tag color={appStatusColor(value)}>{value || "UNSPECIFIED"}</Tag>
      },
      {
        title: isZh ? "更新时间" : "Updated",
        dataIndex: "updatedAt",
        key: "updatedAt",
        width: 200,
        render: (value?: string) => (value ? value.replace("T", " ").slice(0, 19) : "-")
      }
    ],
    [isZh]
  );

  const productTotal = data?.productSummary?.total ?? 0;
  const onShelf = data?.productSummary?.onShelf ?? 0;
  const draft = data?.productSummary?.draft ?? 0;
  const review = data?.productSummary?.reviewing ?? 0;
  const rejected = data?.productSummary?.rejected ?? 0;

  const qualityPercent = productTotal > 0 ? Math.round((onShelf / productTotal) * 100) : 0;

  if (loading) {
    return <Skeleton active paragraph={{ rows: 8 }} />;
  }

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "商家工作台" : "Seller Workbench"}</Title>
          <Text type="secondary">
            {isZh ? "查看店铺与商品经营状态，并快速进入核心操作。" : "Track stores and products, then jump to key actions."}
          </Text>
        </div>
        <Space wrap>
          <Link href="/seller/publish">
            <Button type="primary" icon={<PlusOutlined />}>
              {isZh ? "发布商品" : "Publish Product"}
            </Button>
          </Link>
          <Link href="/onboarding/drafts">
            <Button>{isZh ? "入驻草稿箱" : "Onboarding Drafts"}</Button>
          </Link>
        </Space>
      </header>

      {data?.partial ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "部分服务降级" : "Partial data"}
          description={(data.degradedFields || []).join(", ")}
        />
      ) : null}

      <div className="seller-kpi-grid">
        <Card>
          <Statistic title={isZh ? "店铺数" : "Shops"} value={data?.shopsTotal ?? 0} prefix={<ShopOutlined />} />
        </Card>
        <Card>
          <Statistic title={isZh ? "申请单" : "Applications"} value={data?.applicationsTotal ?? 0} />
        </Card>
        <Card>
          <Statistic title={isZh ? "商品总数" : "Total Products"} value={productTotal} />
        </Card>
        <Card>
          <Statistic title={isZh ? "在架商品" : "On Shelf"} value={onShelf} />
        </Card>
      </div>

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={16}>
          <Card
            title={isZh ? "我的店铺" : "My Shops"}
            extra={
              <Link href="/seller/profile">
                {isZh ? "店铺资料" : "Shop Profile"} <ArrowRightOutlined />
              </Link>
            }
          >
            <Table<ShopRow>
              rowKey={(record) => record.shopNo}
              columns={shopColumns}
              dataSource={shopRows}
              pagination={false}
              locale={{ emptyText: isZh ? "暂无店铺" : "No shops yet" }}
            />
          </Card>
        </Col>

        <Col xs={24} xl={8}>
          <Card title={isZh ? "商品健康度" : "Product Health"} className="seller-side-card">
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "在架率" : "On-shelf ratio"}</span>
                <strong>{qualityPercent}%</strong>
              </div>
              <Progress percent={qualityPercent} strokeColor="#ff6a00" />
              <div>
                <span>{isZh ? "草稿" : "Draft"}</span>
                <strong>{draft}</strong>
              </div>
              <div>
                <span>{isZh ? "审核中" : "Reviewing"}</span>
                <strong>{review}</strong>
              </div>
              <div>
                <span>{isZh ? "驳回" : "Rejected"}</span>
                <strong>{rejected}</strong>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Card title={isZh ? "最近申请" : "Recent Applications"}>
        <Table<SellerApplicationListItem>
          rowKey={(record) => record.applicationNo}
          columns={appColumns}
          dataSource={applicationRows}
          pagination={false}
          locale={{ emptyText: isZh ? "暂无申请记录" : "No recent applications" }}
        />
      </Card>
    </section>
  );
}
