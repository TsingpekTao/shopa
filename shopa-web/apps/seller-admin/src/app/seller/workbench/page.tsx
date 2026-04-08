"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Button, Card, Col, Progress, Row, Skeleton, Space, Statistic, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { ArrowRightOutlined, PlusOutlined, ReloadOutlined, ShopOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";
import { SellerApplicationListItem, SellerWorkbenchResponse } from "@/features/seller-shop/types";
import { getSellerStatDisplayValue, hasSellerDegradedField } from "@/features/seller-shop/sales";

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

function getWorkbenchDegradedMessages(fields: string[] | undefined, isZh: boolean) {
  const messages: string[] = [];
  if (hasSellerDegradedField(fields, "seller_shops")) {
    messages.push(isZh ? "店铺信息稍后更新" : "Shop data will update soon");
  }
  if (hasSellerDegradedField(fields, "seller_applications")) {
    messages.push(isZh ? "申请记录稍后更新" : "Application records will update soon");
  }
  if (hasSellerDegradedField(fields, "catalog_products")) {
    messages.push(isZh ? "商品统计稍后更新" : "Product totals will update soon");
  }
  return messages;
}

export default function WorkbenchPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<SellerWorkbenchResponse | null>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    let active = true;

    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const res = await fetchSellerWorkbench();
        if (!active) {
          return;
        }
        setData(res);
      } catch (err) {
        if (!active) {
          return;
        }
        const nextError = err instanceof Error && err.message ? err.message : isZh ? "工作台暂时无法加载" : "Workbench is unavailable";
        setError(nextError);
        setData(null);
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    load();
    return () => {
      active = false;
    };
  }, [isZh]);

  const shopRows = useMemo<ShopRow[]>(() => {
    return (data?.shops || []).map((item: any) => ({
      shopNo: String(item.shopNo || ""),
      shopDisplayName: String(item.shopDisplayName || item.shopName || "-"),
      shopStatusCode: String(item.shopStatusCode || "UNKNOWN"),
      buyerVisible: Boolean(item.buyerVisible)
    }));
  }, [data?.shops]);

  const applicationRows = useMemo<SellerApplicationListItem[]>(() => data?.latestApplications || [], [data?.latestApplications]);
  const degradedFields = data?.degradedFields;
  const degradedMessages = useMemo(() => getWorkbenchDegradedMessages(degradedFields, isZh), [degradedFields, isZh]);

  const shopsDegraded = hasSellerDegradedField(degradedFields, "seller_shops");
  const appsDegraded = hasSellerDegradedField(degradedFields, "seller_applications");
  const productsDegraded = hasSellerDegradedField(degradedFields, "catalog_products");

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

  const productTotal = productsDegraded ? undefined : data?.productSummary?.total;
  const onShelf = productsDegraded ? undefined : data?.productSummary?.onShelf;
  const draft = productsDegraded ? undefined : data?.productSummary?.draft;
  const review = productsDegraded ? undefined : data?.productSummary?.reviewing;
  const rejected = productsDegraded ? undefined : data?.productSummary?.rejected;

  const qualityPercent =
    typeof productTotal === "number" && productTotal > 0 && typeof onShelf === "number"
      ? Math.round((onShelf / productTotal) * 100)
      : 0;

  if (loading) {
    return <Skeleton active paragraph={{ rows: 8 }} />;
  }

  if (!data) {
    return (
      <section className="seller-page">
        <header className="seller-page-head">
          <div>
            <Title level={3}>{isZh ? "商家工作台" : "Seller Workbench"}</Title>
            <Text type="secondary">{error || (isZh ? "工作台暂时无法加载" : "Workbench is unavailable")}</Text>
          </div>
          <Button icon={<ReloadOutlined />} onClick={() => window.location.reload()}>
            {isZh ? "重新加载" : "Reload"}
          </Button>
        </header>
      </section>
    );
  }

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "商家工作台" : "Seller Workbench"}</Title>
          <Text type="secondary">
            {isZh ? "查看店铺、商品和申请进度。" : "Review shops, products and application progress."}
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

      {degradedMessages.length > 0 ? (
        <Card size="small" style={{ marginBottom: 14 }}>
          <Space wrap>
            {degradedMessages.map((message) => (
              <Tag key={message}>{message}</Tag>
            ))}
          </Space>
        </Card>
      ) : null}

      <div className="seller-kpi-grid">
        <Card>
          <Statistic
            title={isZh ? "店铺数" : "Shops"}
            value={getSellerStatDisplayValue(data.shopsTotal, shopsDegraded)}
            prefix={<ShopOutlined />}
          />
        </Card>
        <Card>
          <Statistic title={isZh ? "申请单" : "Applications"} value={getSellerStatDisplayValue(data.applicationsTotal, appsDegraded)} />
        </Card>
        <Card>
          <Statistic title={isZh ? "商品总数" : "Total Products"} value={getSellerStatDisplayValue(productTotal, productsDegraded)} />
        </Card>
        <Card>
          <Statistic title={isZh ? "在架商品" : "On Shelf"} value={getSellerStatDisplayValue(onShelf, productsDegraded)} />
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
              locale={{ emptyText: shopsDegraded ? (isZh ? "店铺信息暂不可用" : "Shop data unavailable") : isZh ? "暂无店铺" : "No shops yet" }}
            />
          </Card>
        </Col>

        <Col xs={24} xl={8}>
          <Card title={isZh ? "商品概况" : "Product Overview"} className="seller-side-card">
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "在架率" : "On-shelf ratio"}</span>
                <strong>{productsDegraded ? "--" : `${qualityPercent}%`}</strong>
              </div>
              <Progress percent={productsDegraded ? 0 : qualityPercent} strokeColor="#ff6a00" showInfo={false} />
              <div>
                <span>{isZh ? "草稿" : "Draft"}</span>
                <strong>{getSellerStatDisplayValue(draft, productsDegraded)}</strong>
              </div>
              <div>
                <span>{isZh ? "审核中" : "Reviewing"}</span>
                <strong>{getSellerStatDisplayValue(review, productsDegraded)}</strong>
              </div>
              <div>
                <span>{isZh ? "驳回" : "Rejected"}</span>
                <strong>{getSellerStatDisplayValue(rejected, productsDegraded)}</strong>
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
          locale={{ emptyText: appsDegraded ? (isZh ? "申请记录暂不可用" : "Application data unavailable") : isZh ? "暂无申请记录" : "No recent applications" }}
        />
      </Card>
    </section>
  );
}
