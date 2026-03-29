"use client";

import { Alert, Descriptions, Skeleton, Tag } from "antd";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { fetchShopInsights } from "@/features/admin/api";
import { AdminPage } from "@/components/admin-page";

export default function ShopInsightsPage({ params }: { params: { shopNo: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const shopNo = params.shopNo;

  const query = useQuery({
    queryKey: ["admin", "shop-insights", shopNo],
    queryFn: () => fetchShopInsights(shopNo),
    enabled: Boolean(shopNo)
  });

  return (
    <AdminPage
      title={isZh ? "店铺运营洞察" : "Shop Insights"}
      subtitle={isZh ? "查看指定店铺基础状态与可用运营信息。" : "Inspect base status and operational details for a specific shop."}
    >
      {query.isLoading ? <Skeleton active /> : null}

      {!query.isLoading && query.data ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Shop No">{query.data.shopNo}</Descriptions.Item>
          <Descriptions.Item label={isZh ? "店铺名称" : "Shop Name"}>{query.data.shopName || "-"}</Descriptions.Item>
          <Descriptions.Item label={isZh ? "店铺状态" : "Shop Status"}>
            {query.data.shopStatusCode ? <Tag color="blue">{query.data.shopStatusCode}</Tag> : "-"}
          </Descriptions.Item>
        </Descriptions>
      ) : null}

      {!query.isLoading && !query.data ? <Alert type="error" showIcon message={isZh ? "加载店铺洞察失败" : "Failed to load shop insights"} /> : null}

      {query.data?.partial ? (
        <Alert
          style={{ marginTop: 14 }}
          type="warning"
          showIcon
          message={isZh ? "部分字段降级" : "Partial degraded fields"}
          description={(query.data.degradedFields ?? []).join(", ")}
        />
      ) : null}
    </AdminPage>
  );
}
