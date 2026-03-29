"use client";

import { Alert, Empty, Tag } from "antd";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { fetchDashboardOverview } from "@/features/admin/api";
import { AdminPage } from "@/components/admin-page";

const COLOR_BY_KEY: Record<string, string> = {
  totalShops: "blue",
  pendingMerchantReviews: "gold",
  pendingProductReviews: "orange"
};

export default function AdminDashboardPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const { data, isLoading, error } = useQuery({
    queryKey: ["admin", "dashboard-overview"],
    queryFn: fetchDashboardOverview
  });

  return (
    <AdminPage
      title={isZh ? "商城运营看板" : "Mall Dashboard"}
      subtitle={isZh ? "用于查看平台整体运行状态与关键待办。" : "Monitor platform health and key pending tasks."}
    >
      {error ? <Alert type="error" showIcon message={isZh ? "看板数据加载失败，请稍后重试。" : "Failed to load dashboard data."} /> : null}

      {!error && isLoading ? (
        <div className="admin-stat-grid">
          {Array.from({ length: 3 }).map((_, idx) => (
            <div className="admin-stat-card" key={idx}>
              <div className="admin-stat-label">{isZh ? "加载中" : "Loading"}</div>
              <div className="admin-stat-value">--</div>
            </div>
          ))}
        </div>
      ) : null}

      {!error && !isLoading ? (
        (data?.metrics?.length ?? 0) > 0 ? (
          <div className="admin-stat-grid">
            {data?.metrics.map((metric) => (
              <div className="admin-stat-card" key={metric.key}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <div className="admin-stat-label">{metric.label}</div>
                  <Tag color={COLOR_BY_KEY[metric.key] || "default"}>{metric.key}</Tag>
                </div>
                <div className="admin-stat-value">{metric.value}</div>
              </div>
            ))}
          </div>
        ) : (
          <Empty description={isZh ? "暂无可展示指标" : "No metrics available"} />
        )
      ) : null}

      {data?.partial ? (
        <Alert
          style={{ marginTop: 14 }}
          type="warning"
          showIcon
          message={isZh ? "当前看板处于降级模式" : "Dashboard is in degraded mode"}
          description={(data.degradedFields ?? []).join(", ")}
        />
      ) : null}
    </AdminPage>
  );
}
