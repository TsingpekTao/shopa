"use client";

import Link from "next/link";
import { Alert, Empty } from "antd";
import { useQuery } from "@tanstack/react-query";
import { fetchDashboardOverview } from "@/features/admin/api";
import { AdminPage } from "@/components/admin-page";

type DashboardMetricMeta = {
  label: string;
  tone: string;
  hint: string;
  href: string;
};

const METRIC_META: Record<string, DashboardMetricMeta> = {
  totalShops: {
    label: "已开通店铺",
    tone: "#1677ff",
    hint: "平台当前已审核通过并开通的店铺数量",
    href: "/admin/shops"
  },
  pendingMerchantReviews: {
    label: "待审核商家申请",
    tone: "#d48806",
    hint: "等待平台审核处理的商家入驻申请单",
    href: "/admin/review/merchant?statuses=2&statuses=3"
  },
  pendingProductReviews: {
    label: "待审核商品",
    tone: "#fa8c16",
    hint: "等待平台审核的商品提审单",
    href: "/admin/review/product?status=2"
  }
};

export default function AdminDashboardPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["admin", "dashboard-overview"],
    queryFn: fetchDashboardOverview
  });

  const metrics = (data?.metrics ?? []).map((metric) => {
    const meta = METRIC_META[metric.key] ?? {
      label: metric.label || metric.key,
      tone: "#5b3f2a",
      hint: "平台实时指标",
      href: "#"
    };
    return {
      ...metric,
      label: meta.label,
      tone: meta.tone,
      hint: meta.hint,
      href: meta.href
    };
  });

  return (
    <AdminPage title="商城运营看板" subtitle="查看平台实时经营数据与当前待处理事项。">
      {error ? <Alert type="error" showIcon message="看板数据加载失败，请稍后重试。" /> : null}

      {!error && isLoading ? (
        <div className="admin-stat-grid">
          {Array.from({ length: 3 }).map((_, idx) => (
            <div className="admin-stat-card" key={idx}>
              <div className="admin-stat-label">加载中</div>
              <div className="admin-stat-value">--</div>
            </div>
          ))}
        </div>
      ) : null}

      {!error && !isLoading ? (
        metrics.length > 0 ? (
          <div className="admin-stat-grid">
            {metrics.map((metric) => (
              <Link key={metric.key} href={metric.href} style={{ color: "inherit", textDecoration: "none" }}>
                <div className="admin-stat-card" style={{ cursor: "pointer" }}>
                  <div className="admin-stat-label" style={{ color: metric.tone }}>
                    {metric.label}
                  </div>
                  <div className="admin-stat-value">{metric.value}</div>
                  <div style={{ marginTop: 10, color: "#8c6f55", fontSize: 13, lineHeight: 1.7 }}>{metric.hint}</div>
                  <div style={{ marginTop: 14, color: metric.tone, fontSize: 13, fontWeight: 600 }}>点击查看列表</div>
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <Empty description="暂无可展示指标" />
        )
      ) : null}
    </AdminPage>
  );
}
