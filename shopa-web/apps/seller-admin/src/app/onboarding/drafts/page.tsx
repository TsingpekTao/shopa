"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Alert, Button, Card, Empty, Space, Table, Tag, Typography, notification } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useI18n } from "@shopa/ui";
import { listMyApplications } from "@/features/seller-shop/api";
import { SellerApplicationListItem } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

function getStatusColor(status: string): string {
  switch (status) {
    case "DRAFT":
      return "default";
    case "SUBMITTED":
    case "REVIEWING":
      return "processing";
    case "APPROVED":
      return "success";
    case "REJECTED":
    case "SYSTEM_REJECTED":
      return "error";
    case "CANCELLED":
      return "warning";
    default:
      return "default";
  }
}

function resolveApplicationRoute(item: SellerApplicationListItem): string {
  const code = (item.applicationStatusCode || "").toUpperCase();
  if (code === "SUBMITTED" || code === "REVIEWING") {
    return `/onboarding/reviewing?applicationNo=${item.applicationNo}`;
  }
  if (code === "REJECTED" || code === "SYSTEM_REJECTED") {
    return `/onboarding/rejected?applicationNo=${item.applicationNo}`;
  }
  if (code === "APPROVED") {
    return "/seller/workbench";
  }
  return `/onboarding/start?applicationNo=${item.applicationNo}`;
}

function formatDate(value?: string): string {
  if (!value) {
    return "-";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")} ${String(date.getHours()).padStart(2, "0")}:${String(date.getMinutes()).padStart(2, "0")}`;
}

export default function OnboardingDraftsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(false);
  const [rows, setRows] = useState<SellerApplicationListItem[]>([]);

  useEffect(() => {
    setLoading(true);
    listMyApplications({ page: 1, pageSize: 30 })
      .then((res) => {
        setRows(res.applications ?? []);
      })
      .catch((err) => {
        console.error("load drafts failed", err);
        notification.error({
          message: isZh ? "草稿箱加载失败" : "Failed to load drafts"
        });
      })
      .finally(() => setLoading(false));
  }, [isZh]);

  const columns: ColumnsType<SellerApplicationListItem> = useMemo(
    () => [
      {
        title: isZh ? "申请单号" : "Application No",
        dataIndex: "applicationNo",
        key: "applicationNo",
        width: 280,
        render: (value: string) => <Text copyable>{value}</Text>
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "applicationStatusCode",
        key: "applicationStatusCode",
        width: 160,
        render: (value: string) => <Tag color={getStatusColor((value || "").toUpperCase())}>{value || "UNSPECIFIED"}</Tag>
      },
      {
        title: isZh ? "最后更新时间" : "Updated At",
        dataIndex: "updatedAt",
        key: "updatedAt",
        width: 200,
        render: (value?: string) => formatDate(value)
      },
      {
        title: isZh ? "操作" : "Action",
        key: "action",
        render: (_, record) => {
          const code = (record.applicationStatusCode || "").toUpperCase();
          const actionText =
            code === "SUBMITTED" || code === "REVIEWING"
              ? isZh
                ? "查看进度"
                : "View Progress"
              : code === "APPROVED"
                ? isZh
                  ? "进入工作台"
                  : "Go Workbench"
                : isZh
                  ? "继续填写"
                  : "Continue";

          return (
            <Button type="link" onClick={() => (window.location.href = resolveApplicationRoute(record))}>
              {actionText}
            </Button>
          );
        }
      }
    ],
    [isZh]
  );

  return (
    <section style={{ display: "grid", gap: 16 }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12, flexWrap: "wrap" }}>
        <div>
          <Title level={3} style={{ marginBottom: 6 }}>
            {isZh ? "入驻草稿箱" : "Onboarding Drafts"}
          </Title>
          <Text type="secondary">
            {isZh ? "你保存的入驻申请都在这里，可继续填写或查看审核进度。" : "All your onboarding applications are here."}
          </Text>
        </div>
        <Space>
          <Link href="/onboarding/start">
            <Button type="primary">{isZh ? "新建申请" : "New Application"}</Button>
          </Link>
        </Space>
      </header>

      <Card>
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          message={isZh ? "草稿保存说明" : "Draft tips"}
          description={
            isZh
              ? "点击“保存草稿”后，申请会出现在本页。若已提交审核，会在这里显示审核状态。"
              : "Saved drafts and submitted applications will appear here."
          }
        />

        <Table<SellerApplicationListItem>
          loading={loading}
          rowKey={(record) => record.applicationNo}
          columns={columns}
          dataSource={rows}
          pagination={false}
          locale={{
            emptyText: (
              <Empty
                description={isZh ? "暂无草稿" : "No drafts yet"}
                image={Empty.PRESENTED_IMAGE_SIMPLE}
              />
            )
          }}
        />
      </Card>
    </section>
  );
}
