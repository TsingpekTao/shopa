"use client";

import { useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Drawer,
  Empty,
  Image,
  Input,
  Select,
  Skeleton,
  Space,
  Table,
  Tag,
  Typography,
  message
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { AccessControl, useI18n } from "@shopa/ui";
import {
  approveProduct,
  fetchAssetReadUrl,
  fetchInventorySnapshots,
  fetchProductReviewDetail,
  fetchProductReviewTasks,
  fetchShopInsights,
  forceOffShelf,
  freezeProduct,
  rejectProduct,
  unfreezeProduct
} from "@/features/admin/api";
import { Permissions } from "@/features/admin/permissions";
import { ProductInventorySnapshot, ProductReviewDetail, ProductReviewTask } from "@/features/admin/types";
import { AdminPage } from "@/components/admin-page";

const { Title, Text } = Typography;

type StatusFilter = "ALL" | number;

type GroupedTask = {
  shopNo: string;
  shopName: string;
  items: ProductReviewTask[];
};

const STATUS_OPTIONS: Array<{ value: StatusFilter; zh: string; en: string }> = [
  { value: "ALL", zh: "全部状态", en: "All status" },
  { value: 2, zh: "审核中", en: "Reviewing" },
  { value: 3, zh: "已通过", en: "Approved" },
  { value: 4, zh: "在架", en: "On shelf" },
  { value: 5, zh: "已下架", en: "Off shelf" },
  { value: 6, zh: "已驳回", en: "Rejected" },
  { value: 7, zh: "已冻结", en: "Frozen" }
];

function statusText(status: number | undefined, isZh: boolean) {
  switch (status) {
    case 1:
      return isZh ? "草稿" : "Draft";
    case 2:
      return isZh ? "审核中" : "Reviewing";
    case 3:
      return isZh ? "已通过" : "Approved";
    case 4:
      return isZh ? "在架" : "On shelf";
    case 5:
      return isZh ? "已下架" : "Off shelf";
    case 6:
      return isZh ? "已驳回" : "Rejected";
    case 7:
      return isZh ? "已冻结" : "Frozen";
    case 8:
      return isZh ? "已删除" : "Deleted";
    default:
      return isZh ? "未知" : "Unknown";
  }
}

function statusColor(status: number | undefined) {
  switch (status) {
    case 2:
      return "processing";
    case 3:
      return "blue";
    case 4:
      return "success";
    case 5:
      return "default";
    case 6:
      return "error";
    case 7:
      return "volcano";
    default:
      return "default";
  }
}

function stockStatusText(status: number | undefined, isZh: boolean) {
  switch (status) {
    case 1:
      return isZh ? "有货" : "In stock";
    case 2:
      return isZh ? "缺货" : "Out of stock";
    default:
      return isZh ? "未知" : "Unknown";
  }
}

function stockStatusColor(status: number | undefined) {
  switch (status) {
    case 1:
      return "success";
    case 2:
      return "error";
    default:
      return "default";
  }
}

function formatDateTime(value?: string) {
  if (!value) {
    return "-";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  const yyyy = date.getFullYear();
  const mm = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const hh = String(date.getHours()).padStart(2, "0");
  const min = String(date.getMinutes()).padStart(2, "0");
  return `${yyyy}-${mm}-${dd} ${hh}:${min}`;
}

function toPrice(value?: number) {
  if (!value || value <= 0) {
    return "-";
  }
  return `¥${value}`;
}

function toPriceRange(min?: number, max?: number) {
  if (!min && !max) {
    return "-";
  }
  if ((min ?? 0) <= 0 && (max ?? 0) <= 0) {
    return "-";
  }
  if (min === max || !max) {
    return toPrice(min);
  }
  return `¥${min} ~ ¥${max}`;
}

function renderAttrText(values?: Record<string, string>) {
  const entries = Object.entries(values ?? {}).filter(([, value]) => String(value || "").trim());
  if (!entries.length) {
    return "-";
  }
  return entries.map(([key, value]) => `${key}: ${value}`).join(" / ");
}

function normalizeKeyword(value: string) {
  return value.trim().toLowerCase();
}

export default function ProductReviewPage() {
  const searchParams = useSearchParams();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const initialStatus = (() => {
    const raw = searchParams.get("status");
    const parsed = Number(raw);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : "ALL";
  })() as StatusFilter;
  const initialShopKeyword = searchParams.get("keyword")?.trim() ?? "";
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(40);
  const [shopKeywordInput, setShopKeywordInput] = useState(initialShopKeyword);
  const [shopKeyword, setShopKeyword] = useState(initialShopKeyword);
  const [statusFilter, setStatusFilter] = useState<StatusFilter>(initialStatus);
  const [loadingKey, setLoadingKey] = useState("");
  const [detailTask, setDetailTask] = useState<ProductReviewTask | null>(null);

  const tasksQuery = useQuery({
    queryKey: ["admin", "product-review-tasks", page, pageSize, statusFilter],
    queryFn: () =>
      fetchProductReviewTasks({
        page,
        pageSize,
        statuses: statusFilter === "ALL" ? undefined : [statusFilter]
      })
  });

  const shopNos = useMemo(() => {
    return Array.from(new Set((tasksQuery.data?.tasks ?? []).map((item) => item.shopNo).filter(Boolean))) as string[];
  }, [tasksQuery.data?.tasks]);

  const shopMapQuery = useQuery({
    queryKey: ["admin", "review-shop-map", shopNos],
    enabled: shopNos.length > 0,
    queryFn: async () => {
      const pairs = await Promise.all(
        shopNos.map(async (shopNo) => {
          try {
            const shop = await fetchShopInsights(shopNo);
            return [shopNo, shop.shopName || shopNo] as const;
          } catch {
            return [shopNo, shopNo] as const;
          }
        })
      );
      return Object.fromEntries(pairs);
    }
  });

  const detailQuery = useQuery({
    queryKey: ["admin", "product-review-detail", detailTask?.spuNo],
    enabled: Boolean(detailTask?.spuNo),
    retry: false,
    queryFn: async () => {
      if (!detailTask?.spuNo) {
        return null;
      }
      return fetchProductReviewDetail(detailTask.spuNo);
    }
  });

  const detailAssetIds = useMemo(() => {
    const detail = detailQuery.data?.product;
    const skuAssetIds = (detail?.skus ?? []).map((sku) => sku.skuImageAssetId).filter(Boolean) as string[];
    return Array.from(
      new Set([...(detail?.spu?.mainImageAssetIds ?? []), ...(detail?.spu?.detailImageAssetIds ?? []), ...skuAssetIds].filter(Boolean))
    );
  }, [detailQuery.data]);

  const detailSkuNos = useMemo(() => {
    return Array.from(
      new Set((detailQuery.data?.product.skus ?? []).map((item) => item.skuNo).filter(Boolean))
    );
  }, [detailQuery.data]);

  const detailInventoryQuery = useQuery({
    queryKey: ["admin", "product-review-inventory", detailSkuNos],
    enabled: detailSkuNos.length > 0,
    queryFn: async () => fetchInventorySnapshots(detailSkuNos)
  });

  const detailImageQuery = useQuery({
    queryKey: ["admin", "product-review-images", detailAssetIds],
    enabled: detailAssetIds.length > 0,
    queryFn: async () => {
      const pairs = await Promise.all(
        detailAssetIds.map(async (assetId) => {
          try {
            const result = await fetchAssetReadUrl(assetId);
            return [assetId, result.url] as const;
          } catch {
            return [assetId, ""] as const;
          }
        })
      );
      return Object.fromEntries(pairs);
    }
  });

  const groupedTasks = useMemo(() => {
    const keyword = normalizeKeyword(shopKeyword);
    const shopMap = shopMapQuery.data ?? {};
    const bucket = new Map<string, GroupedTask>();

    for (const task of tasksQuery.data?.tasks ?? []) {
      if (statusFilter !== "ALL" && task.spuStatus !== statusFilter) {
        continue;
      }
      const shopNo = task.shopNo || "UNKNOWN_SHOP";
      const shopName = shopMap[shopNo] || shopNo;
      const combinedSearchText = `${shopName} ${shopNo} ${task.title ?? ""}`.toLowerCase();
      if (keyword && !combinedSearchText.includes(keyword)) {
        continue;
      }
      if (!bucket.has(shopNo)) {
        bucket.set(shopNo, {
          shopNo,
          shopName,
          items: []
        });
      }
      bucket.get(shopNo)?.items.push(task);
    }

    return Array.from(bucket.values()).sort((a, b) => a.shopName.localeCompare(b.shopName, "zh-CN"));
  }, [shopKeyword, shopMapQuery.data, statusFilter, tasksQuery.data?.tasks]);

  const detailImages = useMemo(() => {
    const imageMap = detailImageQuery.data ?? {};
    const spu = detailQuery.data?.product.spu;
    const mainImage = (spu?.mainImageAssetIds ?? [])
      .map((assetId) => ({ assetId, url: imageMap[assetId] ?? "" }))
      .find((item) => item.url);
    const gallery = (spu?.detailImageAssetIds ?? [])
      .map((assetId) => ({ assetId, url: imageMap[assetId] ?? "" }))
      .filter((item) => item.url);
    return {
      imageMap,
      mainImage,
      gallery
    };
  }, [detailImageQuery.data, detailQuery.data?.product.spu]);

  const detailInventoryMap = useMemo(() => {
    return new Map<string, ProductInventorySnapshot>(
      (detailInventoryQuery.data ?? []).map((item) => [item.skuNo, item])
    );
  }, [detailInventoryQuery.data]);

  const detailSpuInventorySummary = useMemo(() => {
    const rows = Array.from(detailInventoryMap.values());
    if (!rows.length) {
      return null;
    }
    return rows.reduce(
      (acc, item) => {
        acc.totalQty += Number(item.totalQty ?? 0);
        acc.availableQty += Number(item.availableQty ?? 0);
        if (Number(item.stockStatus ?? 0) === 1) {
          acc.hasInStockSku = true;
        }
        return acc;
      },
      { totalQty: 0, availableQty: 0, hasInStockSku: false }
    );
  }, [detailInventoryMap]);

  const refreshAll = async () => {
    await Promise.all([tasksQuery.refetch(), shopMapQuery.refetch()]);
  };

  const loadVersion = async (spuNo: string) => {
    const detail = await fetchProductReviewDetail(spuNo);
    const version = detail.product.spu?.version ?? 0;
    if (!version) {
      throw new Error(isZh ? "未获取到商品版本号" : "Missing product version");
    }
    return version;
  };

  const runAction = async (key: string, action: () => Promise<void>, successText: string) => {
    try {
      setLoadingKey(key);
      await action();
      message.success(successText);
      await refreshAll();
      if (detailTask?.spuNo && key.startsWith(detailTask.spuNo)) {
        await detailQuery.refetch();
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : isZh ? "操作失败" : "Action failed";
      message.error(errorMessage);
    } finally {
      setLoadingKey("");
    }
  };

  const columns = useMemo<ColumnsType<ProductReviewTask>>(
    () => [
      {
        title: "SPU",
        dataIndex: "spuNo",
        key: "spuNo",
        width: 180
      },
      {
        title: isZh ? "商品标题" : "Product",
        dataIndex: "title",
        key: "title",
        render: (value?: string) => value || "-"
      },
      {
        title: isZh ? "审核状态" : "Status",
        key: "spuStatus",
        width: 120,
        render: (_, row) => <Tag color={statusColor(row.spuStatus)}>{statusText(row.spuStatus, isZh)}</Tag>
      },
      {
        title: isZh ? "提审时间" : "Submitted At",
        dataIndex: "submittedAt",
        key: "submittedAt",
        width: 170,
        render: (value?: string) => formatDateTime(value)
      },
      {
        title: isZh ? "操作" : "Actions",
        key: "actions",
        width: 520,
        render: (_, row) => (
          <Space wrap>
            <Button size="small" onClick={() => setDetailTask(row)}>
              {isZh ? "查看详情" : "Detail"}
            </Button>

            <AccessControl require={Permissions.ProductReviewApprove}>
              <Button
                type="primary"
                size="small"
                loading={loadingKey === `${row.spuNo}-approve`}
                onClick={() =>
                  void runAction(
                    `${row.spuNo}-approve`,
                    async () => {
                      const version = await loadVersion(row.spuNo);
                      await approveProduct({ spuNo: row.spuNo, expectedVersion: version });
                    },
                    isZh ? "审核通过成功" : "Approved"
                  )
                }
              >
                {isZh ? "通过" : "Approve"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewReject}>
              <Button
                danger
                size="small"
                loading={loadingKey === `${row.spuNo}-reject`}
                onClick={() =>
                  void runAction(
                    `${row.spuNo}-reject`,
                    async () => {
                      const version = await loadVersion(row.spuNo);
                      await rejectProduct({
                        spuNo: row.spuNo,
                        expectedVersion: version,
                        rejectReasonCode: "CONTENT_INVALID",
                        rejectComment: isZh ? "请补充商品信息后重新提交" : "Please complete product information and resubmit"
                      });
                    },
                    isZh ? "驳回成功" : "Rejected"
                  )
                }
              >
                {isZh ? "驳回" : "Reject"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewFreeze}>
              <Button
                size="small"
                loading={loadingKey === `${row.spuNo}-freeze`}
                onClick={() =>
                  void runAction(
                    `${row.spuNo}-freeze`,
                    async () => {
                      const version = await loadVersion(row.spuNo);
                      await freezeProduct({ spuNo: row.spuNo, expectedVersion: version, reasonCode: "ADMIN_CONTROL" });
                    },
                    isZh ? "冻结成功" : "Frozen"
                  )
                }
              >
                {isZh ? "冻结" : "Freeze"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewUnfreeze}>
              <Button
                size="small"
                loading={loadingKey === `${row.spuNo}-unfreeze`}
                onClick={() =>
                  void runAction(
                    `${row.spuNo}-unfreeze`,
                    async () => {
                      const version = await loadVersion(row.spuNo);
                      await unfreezeProduct({ spuNo: row.spuNo, expectedVersion: version, reasonCode: "ADMIN_CONTROL" });
                    },
                    isZh ? "解冻成功" : "Unfrozen"
                  )
                }
              >
                {isZh ? "解冻" : "Unfreeze"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewForceOffShelf}>
              <Button
                size="small"
                danger
                loading={loadingKey === `${row.spuNo}-off`}
                onClick={() =>
                  void runAction(
                    `${row.spuNo}-off`,
                    async () => {
                      const version = await loadVersion(row.spuNo);
                      await forceOffShelf({ spuNo: row.spuNo, expectedVersion: version, reasonCode: "ADMIN_FORCE" });
                    },
                    isZh ? "强制下架成功" : "Forced off shelf"
                  )
                }
              >
                {isZh ? "强制下架" : "Force Off Shelf"}
              </Button>
            </AccessControl>
          </Space>
        )
      }
    ],
    [detailQuery, detailTask?.spuNo, isZh, loadingKey, shopMapQuery, tasksQuery]
  );

  const detailData: ProductReviewDetail | null = detailQuery.data ?? null;
  const detailSpu = detailData?.product.spu;
  const detailShopName = (detailSpu?.shopNo && shopMapQuery.data?.[detailSpu.shopNo]) || detailTask?.shopNo || "-";

  return (
    <AdminPage
      title={isZh ? "商品审核工作台" : "Product Review Workspace"}
      subtitle={isZh ? "先按店铺筛，再看该店铺下所有待处理商品。审核动作统一从详情拿真实版本号，避免误操作。" : "Filter by shop, then review all products under that shop with real versioned actions."}
      extra={
        <Space wrap>
          <Input
            allowClear
            value={shopKeywordInput}
            onChange={(event) => setShopKeywordInput(event.target.value)}
            onPressEnter={() => {
              setPage(1);
              setShopKeyword(shopKeywordInput);
            }}
            placeholder={isZh ? "按店铺名 / 店铺编号搜索" : "Search by shop name / shop no"}
            style={{ width: 240 }}
          />
          <Select
            value={statusFilter}
            onChange={(value) => {
              setPage(1);
              setStatusFilter(value);
            }}
            style={{ width: 160 }}
            options={STATUS_OPTIONS.map((item) => ({ label: isZh ? item.zh : item.en, value: item.value }))}
          />
          <Button
            type="primary"
            onClick={() => {
              setPage(1);
              setShopKeyword(shopKeywordInput);
            }}
          >
            {isZh ? "筛选" : "Filter"}
          </Button>
          <Button
            onClick={() => {
              setShopKeywordInput("");
              setShopKeyword("");
              setStatusFilter("ALL");
              setPage(1);
              void refreshAll();
            }}
          >
            {isZh ? "重置" : "Reset"}
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size={16} style={{ width: "100%" }}>
        <Card size="small">
          <Space split={<span style={{ color: "#d9d9d9" }}>|</span>} wrap>
            <Text>{isZh ? `当前页任务 ${tasksQuery.data?.tasks.length ?? 0} 条` : `${tasksQuery.data?.tasks.length ?? 0} tasks`}</Text>
            <Text>{isZh ? `匹配店铺 ${groupedTasks.length} 家` : `${groupedTasks.length} shops matched`}</Text>
            <Text>{isZh ? `总任务数 ${tasksQuery.data?.total ?? 0}` : `${tasksQuery.data?.total ?? 0} total`}</Text>
          </Space>
        </Card>

        {tasksQuery.isLoading ? <Skeleton active paragraph={{ rows: 10 }} /> : null}

        {!tasksQuery.isLoading && !groupedTasks.length ? (
          <Card>
            <Empty description={isZh ? "当前筛选条件下没有商品审核任务" : "No product review tasks found"} />
          </Card>
        ) : null}

        {!tasksQuery.isLoading
          ? groupedTasks.map((group) => (
              <Card
                key={group.shopNo}
                title={
                  <Space direction="vertical" size={2}>
                    <Title level={5} style={{ margin: 0 }}>
                      {group.shopName || group.shopNo}
                    </Title>
                    <Text type="secondary">{group.shopNo}</Text>
                  </Space>
                }
                extra={<Tag color="blue">{isZh ? `${group.items.length} 个商品` : `${group.items.length} products`}</Tag>}
              >
                <Table<ProductReviewTask>
                  rowKey={(record) => record.taskNo || record.spuNo}
                  dataSource={group.items}
                  columns={columns}
                  pagination={false}
                  scroll={{ x: 1380 }}
                />
              </Card>
            ))
          : null}

        <Card size="small">
          <Space align="center" wrap>
            <Button disabled={page <= 1} onClick={() => setPage((current) => Math.max(1, current - 1))}>
              {isZh ? "上一页" : "Previous"}
            </Button>
            <Text>{isZh ? `第 ${page} 页` : `Page ${page}`}</Text>
            <Button
              disabled={(tasksQuery.data?.tasks.length ?? 0) < pageSize}
              onClick={() => setPage((current) => current + 1)}
            >
              {isZh ? "下一页" : "Next"}
            </Button>
            <Select
              value={pageSize}
              style={{ width: 120 }}
              options={[
                { value: 20, label: isZh ? "每页 20 条" : "20 / page" },
                { value: 40, label: isZh ? "每页 40 条" : "40 / page" },
                { value: 60, label: isZh ? "每页 60 条" : "60 / page" }
              ]}
              onChange={(value) => {
                setPage(1);
                setPageSize(value);
              }}
            />
          </Space>
        </Card>
      </Space>

      <Drawer width={760} open={Boolean(detailTask)} onClose={() => setDetailTask(null)} title={isZh ? "商品审核详情" : "Review Detail"}>
        {detailQuery.isLoading ? <Skeleton active paragraph={{ rows: 10 }} /> : null}

        {detailQuery.isError ? (
          <Space direction="vertical" size={16} style={{ width: "100%" }}>
            <Alert
              type="warning"
              showIcon
              message={isZh ? "该审核任务对应的商品主记录不存在" : "The product record for this task no longer exists"}
              description={
                isZh
                  ? "后端返回了 `query spu failed: sql: no rows in result set`。页面已改为兜底展示任务信息，不再直接崩溃。"
                  : "The backend returned no product row. The page now falls back to task-level information instead of crashing."
              }
            />
            <Card size="small">
              <Descriptions
                column={1}
                size="small"
                items={[
                  { key: "shop", label: isZh ? "店铺" : "Shop", children: detailTask ? `${detailTask.shopNo || "-"}` : "-" },
                  { key: "spu", label: "SPU", children: detailTask?.spuNo || "-" },
                  { key: "title", label: isZh ? "商品标题" : "Title", children: detailTask?.title || "-" },
                  {
                    key: "status",
                    label: isZh ? "状态" : "Status",
                    children: <Tag color={statusColor(detailTask?.spuStatus)}>{statusText(detailTask?.spuStatus, isZh)}</Tag>
                  },
                  {
                    key: "submittedAt",
                    label: isZh ? "提审时间" : "Submitted At",
                    children: formatDateTime(detailTask?.submittedAt)
                  }
                ]}
              />
            </Card>
          </Space>
        ) : null}

        {!detailQuery.isLoading && !detailQuery.isError && !detailSpu ? (
          <Empty description={isZh ? "暂无详情数据" : "No detail data"} />
        ) : null}

        {!detailQuery.isLoading && !detailQuery.isError && detailSpu ? (
          <Space direction="vertical" size={16} style={{ width: "100%" }}>
            <Card size="small">
              <Space direction="vertical" size={12} style={{ width: "100%" }}>
                <Space direction="vertical" size={4}>
                  <Title level={5} style={{ margin: 0 }}>
                    {detailSpu.title || "-"}
                  </Title>
                  <Text type="secondary">{detailSpu.subTitle || (isZh ? "暂无副标题" : "No subtitle")}</Text>
                </Space>

                <Space wrap>
                  <Tag color={statusColor(detailSpu.spuStatus)}>{statusText(detailSpu.spuStatus, isZh)}</Tag>
                  <Tag>{`${isZh ? "售价" : "Price"} ${toPriceRange(detailSpu.minSalePrice, detailSpu.maxSalePrice)}`}</Tag>
                </Space>

                <Descriptions
                  column={1}
                  size="small"
                  items={[
                    { key: "shop", label: isZh ? "店铺" : "Shop", children: `${detailShopName} (${detailSpu.shopNo || "-"})` },
                    { key: "spu", label: "SPU", children: detailSpu.spuNo },
                    { key: "category", label: isZh ? "类目 ID" : "Category ID", children: detailSpu.categoryId || "-" },
                    { key: "brand", label: isZh ? "品牌编号" : "Brand", children: detailSpu.brandNo || "-" },
                    { key: "attrs", label: isZh ? "商品属性" : "Attributes", children: renderAttrText(detailSpu.attributeValues) },
                    {
                      key: "realStock",
                      label: isZh ? "实时库存" : "Live Stock",
                      children: detailSpuInventorySummary ? (
                        <Space wrap size={[8, 8]}>
                          <Tag color={detailSpuInventorySummary.hasInStockSku ? "success" : "error"}>
                            {detailSpuInventorySummary.hasInStockSku
                              ? isZh ? "有货" : "In stock"
                              : isZh ? "缺货" : "Out of stock"}
                          </Tag>
                          <Text>
                            {isZh
                              ? `可售 ${detailSpuInventorySummary.availableQty} / 总库存 ${detailSpuInventorySummary.totalQty}`
                              : `Available ${detailSpuInventorySummary.availableQty} / Total ${detailSpuInventorySummary.totalQty}`}
                          </Text>
                        </Space>
                      ) : (
                        <Text type="secondary">{isZh ? "未读取到实时库存" : "Live inventory unavailable"}</Text>
                      )
                    },
                    { key: "version", label: isZh ? "当前版本" : "Version", children: detailSpu.version || "-" },
                    {
                      key: "review",
                      label: isZh ? "审核备注" : "Review",
                      children: detailData?.review?.rejectComment || detailData?.review?.rejectReasonCode || (isZh ? "暂无" : "N/A")
                    }
                  ]}
                />
              </Space>
            </Card>

            <Card size="small" title={isZh ? "商品图片" : "Images"}>
              <Space direction="vertical" size={16} style={{ width: "100%" }}>
                <div>
                  <Text strong>{isZh ? "主图" : "Main Image"}</Text>
                  <div style={{ marginTop: 10 }}>
                    {detailImages.mainImage?.url ? (
                      <Image src={detailImages.mainImage.url} alt={detailSpu.title} width={240} height={240} style={{ borderRadius: 16, objectFit: "cover" }} />
                    ) : (
                      <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "暂无主图" : "No main image"} />
                    )}
                  </div>
                </div>
                <div>
                  <Text strong>{isZh ? "详情图" : "Gallery"}</Text>
                  <div style={{ marginTop: 10 }}>
                    {detailImages.gallery.length ? (
                      <Image.PreviewGroup>
                        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(120px, 1fr))", gap: 12 }}>
                          {detailImages.gallery.map((item) => (
                            <Image
                              key={item.assetId}
                              src={item.url}
                              alt={`${detailSpu.title || "product"}-${item.assetId}`}
                              height={120}
                              style={{ borderRadius: 14, objectFit: "cover" }}
                            />
                          ))}
                        </div>
                      </Image.PreviewGroup>
                    ) : (
                      <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "暂无详情图" : "No gallery images"} />
                    )}
                  </div>
                </div>
              </Space>
            </Card>

            <Card size="small" title={isZh ? "SKU 列表" : "SKU List"}>
              {detailInventoryQuery.isFetching ? (
                <Alert
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                  message={isZh ? "正在读取实时库存" : "Loading live inventory"}
                />
              ) : null}
              <Table
                rowKey={(record) => record.skuNo}
                dataSource={detailData?.product.skus ?? []}
                pagination={false}
                locale={{ emptyText: isZh ? "暂无 SKU" : "No SKUs" }}
                columns={[
                  { title: "SKU", dataIndex: "skuNo", key: "skuNo", width: 160 },
                  {
                    title: isZh ? "图片" : "Image",
                    key: "image",
                    width: 92,
                    render: (_, row: ProductReviewDetail["product"]["skus"][number]) =>
                      row.skuImageAssetId && detailImages.imageMap[row.skuImageAssetId] ? (
                        <Image src={detailImages.imageMap[row.skuImageAssetId]} alt={row.skuName} width={56} height={56} style={{ borderRadius: 12, objectFit: "cover" }} />
                      ) : (
                        <Text type="secondary">-</Text>
                      )
                  },
                  { title: isZh ? "名称" : "Name", dataIndex: "skuName", key: "skuName" },
                  {
                    title: isZh ? "规格" : "Specs",
                    key: "saleAttrs",
                    render: (_, row: ProductReviewDetail["product"]["skus"][number]) => renderAttrText(row.saleAttrs)
                  },
                  {
                    title: isZh ? "售价" : "Sale Price",
                    dataIndex: "salePrice",
                    key: "salePrice",
                    width: 120,
                    render: (value?: number) => toPrice(value)
                  },
                  {
                    title: isZh ? "市场价" : "Market Price",
                    dataIndex: "marketPrice",
                    key: "marketPrice",
                    width: 120,
                    render: (value?: number) => toPrice(value)
                  },
                  {
                    title: isZh ? "库存数量" : "Qty",
                    key: "inventoryQty",
                    width: 150,
                    render: (_, row: ProductReviewDetail["product"]["skus"][number]) => {
                      const snapshot = detailInventoryMap.get(row.skuNo);
                      if (!snapshot) {
                        return <Text type="secondary">-</Text>;
                      }
                      return `${snapshot.availableQty ?? 0} / ${snapshot.totalQty ?? 0}`;
                    }
                  },
                  {
                    title: isZh ? "库存状态" : "Stock",
                    dataIndex: "stockStatus",
                    key: "stockStatus",
                    width: 120,
                    render: (value: number | undefined, row: ProductReviewDetail["product"]["skus"][number]) => {
                      const liveStatus = detailInventoryMap.get(row.skuNo)?.stockStatus;
                      const finalStatus = liveStatus ?? value;
                      return <Tag color={stockStatusColor(finalStatus)}>{stockStatusText(finalStatus, isZh)}</Tag>;
                    }
                  }
                ]}
              />
            </Card>
          </Space>
        ) : null}
      </Drawer>
    </AdminPage>
  );
}
