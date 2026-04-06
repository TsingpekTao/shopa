"use client";

import { Alert, Descriptions, Skeleton, Tag } from "antd";
import { useQuery } from "@tanstack/react-query";
import { fetchShopInsights } from "@/features/admin/api";
import { AdminPage } from "@/components/admin-page";

export default function ShopInsightsPage({ params }: { params: { shopNo: string } }) {
  const shopNo = params.shopNo;

  const query = useQuery({
    queryKey: ["admin", "shop-insights", shopNo],
    queryFn: () => fetchShopInsights(shopNo),
    enabled: Boolean(shopNo)
  });

  return (
    <AdminPage title="店铺详情" subtitle="查看指定店铺的基础状态与运营信息。">
      {query.isLoading ? <Skeleton active /> : null}

      {!query.isLoading && query.data ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="店铺编号">{query.data.shopNo}</Descriptions.Item>
          <Descriptions.Item label="店铺名称">{query.data.shopName || "-"}</Descriptions.Item>
          <Descriptions.Item label="店铺状态">
            {query.data.shopStatusCode ? <Tag color="blue">{query.data.shopStatusCode}</Tag> : "-"}
          </Descriptions.Item>
        </Descriptions>
      ) : null}

      {!query.isLoading && !query.data ? <Alert type="error" showIcon message="加载店铺详情失败" /> : null}
    </AdminPage>
  );
}
