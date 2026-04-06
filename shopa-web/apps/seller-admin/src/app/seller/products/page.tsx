"use client";

import { useEffect, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Drawer,
  Empty,
  Input,
  Modal,
  Select,
  Skeleton,
  Space,
  Table,
  Tag,
  Typography,
  notification
} from "antd";
import type { ColumnsType } from "antd/es/table";
import type { TablePaginationConfig } from "antd/es/table";
import { DeleteOutlined, DownloadOutlined, EyeOutlined, ReloadOutlined, TagsOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { deleteProductDraft, getSellerProductDetail, listSellerProducts } from "@/features/catalog/api";
import { SellerProductSku, SellerProductSpu, SpuStatusCode, StockStatusCode } from "@/features/catalog/types";
import {
  batchGetSellerProductStoreCategoryBindings,
  fetchSellerWorkbench,
  listSellerStoreCategories,
  updateSellerProductStoreCategoryBinding
} from "@/features/seller-shop/api";
import {
  SellerProductStoreCategoryBinding,
  SellerStoreCategory,
  SellerWorkbenchResponse
} from "@/features/seller-shop/types";

const { Title, Text } = Typography;

type StatusFilter = "ALL" | SpuStatusCode;
type StoreCategoryFilterValue = "ALL" | "UNCLASSIFIED" | `l1:${number}` | `leaf:${number}`;

const statusOptions: Array<{ value: StatusFilter; zh: string; en: string }> = [
  { value: "ALL", zh: "全部状态", en: "All status" },
  { value: "SPU_STATUS_DRAFT", zh: "草稿", en: "Draft" },
  { value: "SPU_STATUS_REVIEWING", zh: "审核中", en: "Reviewing" },
  { value: "SPU_STATUS_APPROVED", zh: "已通过", en: "Approved" },
  { value: "SPU_STATUS_ON_SHELF", zh: "在架", en: "On shelf" },
  { value: "SPU_STATUS_OFF_SHELF", zh: "下架", en: "Off shelf" },
  { value: "SPU_STATUS_REJECTED", zh: "驳回", en: "Rejected" }
];

function statusColor(status: SpuStatusCode) {
  if (status === "SPU_STATUS_ON_SHELF") return "success";
  if (status === "SPU_STATUS_REVIEWING") return "processing";
  if (status === "SPU_STATUS_REJECTED") return "error";
  if (status === "SPU_STATUS_APPROVED") return "blue";
  return "default";
}

function stockStatusColor(status: StockStatusCode) {
  if (status === "STOCK_STATUS_OUT_OF_STOCK") return "error";
  if (status === "STOCK_STATUS_IN_STOCK") return "success";
  return "default";
}

function statusText(status: SpuStatusCode, isZh: boolean): string {
  switch (status) {
    case "SPU_STATUS_DRAFT":
      return isZh ? "草稿" : "Draft";
    case "SPU_STATUS_REVIEWING":
      return isZh ? "审核中" : "Reviewing";
    case "SPU_STATUS_APPROVED":
      return isZh ? "已通过" : "Approved";
    case "SPU_STATUS_ON_SHELF":
      return isZh ? "在架" : "On shelf";
    case "SPU_STATUS_OFF_SHELF":
      return isZh ? "下架" : "Off shelf";
    case "SPU_STATUS_REJECTED":
      return isZh ? "驳回" : "Rejected";
    case "SPU_STATUS_FROZEN":
      return isZh ? "冻结" : "Frozen";
    case "SPU_STATUS_DELETED":
      return isZh ? "已删除" : "Deleted";
    default:
      return "UNSPECIFIED";
  }
}

function stockStatusText(status: StockStatusCode, isZh: boolean): string {
  switch (status) {
    case "STOCK_STATUS_IN_STOCK":
      return isZh ? "有货" : "In stock";
    case "STOCK_STATUS_OUT_OF_STOCK":
      return isZh ? "缺货" : "Out of stock";
    default:
      return isZh ? "未知" : "Unknown";
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

function toPriceRange(minSalePrice: number, maxSalePrice: number): string {
  if (minSalePrice <= 0 && maxSalePrice <= 0) {
    return "-";
  }
  if (minSalePrice === maxSalePrice || maxSalePrice <= 0) {
    return `¥${minSalePrice}`;
  }
  return `¥${minSalePrice} ~ ¥${maxSalePrice}`;
}

function canDeleteProduct(status: SpuStatusCode): boolean {
  return status === "SPU_STATUS_DRAFT" || status === "SPU_STATUS_REJECTED" || status === "SPU_STATUS_OFF_SHELF";
}

function flattenLeafCategoryOptions(categories: SellerStoreCategory[]) {
  const options: Array<{ value: number; label: string }> = [];
  categories.forEach((category) => {
    if (category.children.length === 0) {
      options.push({ value: category.id, label: category.name });
      return;
    }
    category.children.forEach((child) => {
      options.push({ value: child.id, label: `${category.name} / ${child.name}` });
    });
  });
  return options;
}

function buildCategoryFilterOptions(categories: SellerStoreCategory[], isZh: boolean) {
  const options: Array<{ label: string; value: StoreCategoryFilterValue }> = [
    { label: isZh ? "全部店内分类" : "All store categories", value: "ALL" },
    { label: isZh ? "未分类" : "Unclassified", value: "UNCLASSIFIED" }
  ];

  categories.forEach((category) => {
    options.push({
      label: `${isZh ? "一级" : "L1"} · ${category.name}`,
      value: `l1:${category.id}`
    });

    if (category.children.length === 0) {
      options.push({
        label: category.name,
        value: `leaf:${category.id}`
      });
      return;
    }

    category.children.forEach((child) => {
      options.push({
        label: `${category.name} / ${child.name}`,
        value: `leaf:${child.id}`
      });
    });
  });

  return options;
}

function resolveBindingLabel(
  binding: SellerProductStoreCategoryBinding | undefined,
  categories: SellerStoreCategory[],
  isZh: boolean
) {
  if (!binding?.storeCategoryId) {
    return isZh ? "未分类" : "Unclassified";
  }
  if (binding.storeCategoryName) {
    if (binding.storeCategoryPath.length >= 2) {
      const primary = categories.find((item) => item.id === binding.storeCategoryPath[0]);
      if (primary) {
        return `${primary.name} / ${binding.storeCategoryName}`;
      }
    }
    return binding.storeCategoryName;
  }
  return isZh ? "未分类" : "Unclassified";
}

function matchesCategoryFilter(
  row: SellerProductSpu,
  filterValue: StoreCategoryFilterValue,
  bindingMap: Map<string, SellerProductStoreCategoryBinding>
) {
  if (filterValue === "ALL") {
    return true;
  }
  const binding = bindingMap.get(row.spuNo);
  if (filterValue === "UNCLASSIFIED") {
    return !binding?.storeCategoryId;
  }
  if (filterValue.startsWith("l1:")) {
    const target = Number(filterValue.slice(3));
    return binding?.storeCategoryPath.includes(target) ?? false;
  }
  if (filterValue.startsWith("leaf:")) {
    const target = Number(filterValue.slice(5));
    return binding?.storeCategoryId === target;
  }
  return true;
}

async function listProductsAcrossPages(params: {
  keyword: string;
  statusFilter: StatusFilter;
}): Promise<SellerProductSpu[]> {
  const pageSize = 100;
  const items: SellerProductSpu[] = [];
  let page = 1;
  let total = 0;

  while (page <= 20) {
    const response = await listSellerProducts({
      page,
      pageSize,
      keyword: params.keyword,
      statuses: params.statusFilter === "ALL" ? undefined : [params.statusFilter]
    });
    items.push(...response.products);
    total = response.total;
    if (items.length >= total || response.products.length < pageSize) {
      break;
    }
    page += 1;
  }

  return items;
}

export default function ProductsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [keyword, setKeyword] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");
  const [storeCategoryFilter, setStoreCategoryFilter] = useState<StoreCategoryFilterValue>("ALL");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [detailSpuNo, setDetailSpuNo] = useState<string | null>(null);
  const [selectedSpuNos, setSelectedSpuNos] = useState<string[]>([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const [batchCategoryId, setBatchCategoryId] = useState<number | null>(null);
  const [isUpdatingBatchCategory, setIsUpdatingBatchCategory] = useState(false);

  const filterActive = storeCategoryFilter !== "ALL";

  const workbenchQuery = useQuery<SellerWorkbenchResponse>({
    queryKey: ["seller-workbench", "products-v2"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30_000
  });
  const shop = workbenchQuery.data?.shops?.[0];

  const categoriesQuery = useQuery({
    queryKey: ["seller-store-categories", shop?.shopNo, "products"],
    queryFn: () => listSellerStoreCategories(shop?.shopNo ?? ""),
    enabled: Boolean(shop?.shopNo),
    staleTime: 30_000
  });

  const productsQuery = useQuery({
    queryKey: ["seller-products", page, pageSize, keyword, statusFilter],
    queryFn: async () =>
      listSellerProducts({
        page,
        pageSize,
        keyword,
        statuses: statusFilter === "ALL" ? undefined : [statusFilter]
      }),
    enabled: !filterActive
  });

  const filteredUniverseQuery = useQuery({
    queryKey: ["seller-products-universe", keyword, statusFilter, storeCategoryFilter],
    queryFn: () => listProductsAcrossPages({ keyword, statusFilter }),
    enabled: filterActive
  });

  const detailQuery = useQuery({
    queryKey: ["seller-product-detail", detailSpuNo],
    queryFn: async () => {
      if (!detailSpuNo) {
        return null;
      }
      return getSellerProductDetail(detailSpuNo);
    },
    enabled: Boolean(detailSpuNo)
  });

  const bindingSourceRows = filterActive ? filteredUniverseQuery.data ?? [] : productsQuery.data?.products ?? [];
  const bindingSourceSpuNos = useMemo(
    () => Array.from(new Set(bindingSourceRows.map((item) => item.spuNo).filter(Boolean))),
    [bindingSourceRows]
  );

  const bindingsQuery = useQuery({
    queryKey: ["seller-product-store-category-bindings", shop?.shopNo, bindingSourceSpuNos.join("|")],
    queryFn: () => batchGetSellerProductStoreCategoryBindings(shop?.shopNo ?? "", bindingSourceSpuNos),
    enabled: Boolean(shop?.shopNo) && bindingSourceSpuNos.length > 0,
    staleTime: 10_000
  });

  const categories = categoriesQuery.data ?? [];
  const bindingMap = useMemo(() => {
    const map = new Map<string, SellerProductStoreCategoryBinding>();
    for (const item of bindingsQuery.data ?? []) {
      map.set(item.spuNo, item);
    }
    return map;
  }, [bindingsQuery.data]);

  const filteredRows = useMemo(() => {
    if (!filterActive) {
      return [];
    }
    return (filteredUniverseQuery.data ?? []).filter((row) => matchesCategoryFilter(row, storeCategoryFilter, bindingMap));
  }, [bindingMap, filterActive, filteredUniverseQuery.data, storeCategoryFilter]);

  const allRowsForSelection = filterActive ? filteredRows : productsQuery.data?.products ?? [];
  const rows = useMemo(() => {
    if (!filterActive) {
      return productsQuery.data?.products ?? [];
    }
    const start = (page - 1) * pageSize;
    return filteredRows.slice(start, start + pageSize);
  }, [filterActive, filteredRows, page, pageSize, productsQuery.data?.products]);
  const total = filterActive ? filteredRows.length : productsQuery.data?.total ?? 0;

  const selectedProducts = useMemo(
    () => allRowsForSelection.filter((item) => selectedSpuNos.includes(item.spuNo)),
    [allRowsForSelection, selectedSpuNos]
  );
  const deletableSelectedProducts = useMemo(
    () => selectedProducts.filter((item) => canDeleteProduct(item.spuStatus)),
    [selectedProducts]
  );
  const leafCategoryOptions = useMemo(() => flattenLeafCategoryOptions(categories), [categories]);
  const categoryFilterOptions = useMemo(() => buildCategoryFilterOptions(categories, isZh), [categories, isZh]);

  useEffect(() => {
    setSelectedSpuNos((prev) => prev.filter((key) => allRowsForSelection.some((row) => row.spuNo === key)));
  }, [allRowsForSelection]);

  useEffect(() => {
    setPage(1);
  }, [keyword, statusFilter, storeCategoryFilter]);

  const exportCurrentPage = () => {
    if (!rows.length) {
      notification.info({ message: isZh ? "当前页没有可导出的数据" : "No data to export on this page" });
      return;
    }
    const head = ["spu_no", "title", "store_category", "status", "stock_status", "price_range", "updated_at"];
    const body = rows.map((item) =>
      [
        item.spuNo,
        item.title,
        resolveBindingLabel(bindingMap.get(item.spuNo), categories, false),
        statusText(item.spuStatus, false),
        stockStatusText(item.spuStockStatus, false),
        toPriceRange(item.minSalePrice, item.maxSalePrice),
        formatDateTime(item.updatedAt)
      ]
        .map((field) => `"${String(field).replaceAll('"', '""')}"`)
        .join(",")
    );
    const csv = `\uFEFF${head.join(",")}\n${body.join("\n")}`;
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `seller-products-${Date.now()}.csv`;
    anchor.click();
    URL.revokeObjectURL(url);
  };

  const refreshAll = async () => {
    if (filterActive) {
      await filteredUniverseQuery.refetch();
    } else {
      await productsQuery.refetch();
    }
    await bindingsQuery.refetch();
    await categoriesQuery.refetch();
  };

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? 10);
  };

  const handleDeleteProducts = async (targets: SellerProductSpu[]) => {
    if (!targets.length) {
      notification.warning({ message: isZh ? "请选择可删除的商品" : "Please select deletable products" });
      return;
    }
    setIsDeleting(true);
    try {
      const failed: Array<{ spuNo: string; reason: string }> = [];
      for (const target of targets) {
        try {
          await deleteProductDraft({
            shopNo: target.shopNo,
            spuNo: target.spuNo,
            expectedVersion: target.version,
            reasonCode: "SELLER_MANUAL_DELETE"
          });
        } catch (error) {
          failed.push({
            spuNo: target.spuNo,
            reason: error instanceof Error ? error.message : String(error)
          });
        }
      }
      if (failed.length) {
        notification.error({
          message: isZh ? "部分商品删除失败" : "Some products failed to delete",
          description: failed.map((item) => `${item.spuNo}: ${item.reason}`).join("；")
        });
      } else {
        notification.success({ message: isZh ? "删除成功" : "Delete success" });
      }
      setSelectedSpuNos((prev) => prev.filter((spuNo) => !targets.some((item) => item.spuNo === spuNo)));
      await refreshAll();
    } finally {
      setIsDeleting(false);
    }
  };

  const confirmDeleteTargets = (targets: SellerProductSpu[]) => {
    if (!targets.length) {
      return;
    }
    Modal.confirm({
      title: isZh ? "确认删除草稿" : "Confirm delete draft",
      content: isZh
        ? `将删除 ${targets.length} 个商品。仅草稿、驳回、下架状态允许删除。`
        : `This will delete ${targets.length} products. Only draft, rejected, or off-shelf items can be deleted.`,
      onOk: () => handleDeleteProducts(targets)
    });
  };

  const handleBatchDelete = () => {
    confirmDeleteTargets(deletableSelectedProducts);
  };

  const handleBatchUpdateStoreCategory = async () => {
    if (!shop?.shopNo) {
      notification.warning({ message: isZh ? "当前账号没有可操作店铺" : "No available shop for this account" });
      return;
    }
    if (!selectedProducts.length) {
      notification.warning({ message: isZh ? "请先选择商品" : "Select products first" });
      return;
    }

    setIsUpdatingBatchCategory(true);
    try {
      const failed: string[] = [];
      for (const product of selectedProducts) {
        try {
          await updateSellerProductStoreCategoryBinding(shop.shopNo, product.spuNo, Number(batchCategoryId ?? 0));
        } catch (error) {
          failed.push(`${product.spuNo}: ${error instanceof Error ? error.message : "更新失败"}`);
        }
      }

      if (failed.length) {
        notification.error({
          message: isZh ? "部分商品分类更新失败" : "Some product categories failed to update",
          description: failed.join("；")
        });
      } else {
        notification.success({
          message: isZh ? "店内分类已批量更新" : "Store categories updated",
          description:
            batchCategoryId && batchCategoryId > 0
              ? isZh
                ? `已将 ${selectedProducts.length} 个商品归入指定店内分类。`
                : `${selectedProducts.length} products moved into the selected store category.`
              : isZh
                ? `已将 ${selectedProducts.length} 个商品改为未分类。`
                : `${selectedProducts.length} products were marked unclassified.`
        });
      }

      await bindingsQuery.refetch();
      if (filterActive) {
        await filteredUniverseQuery.refetch();
      } else {
        await productsQuery.refetch();
      }
    } finally {
      setIsUpdatingBatchCategory(false);
    }
  };

  const columns = useMemo<ColumnsType<SellerProductSpu>>(
    () => [
      { title: "SPU", dataIndex: "spuNo", key: "spuNo", width: 150 },
      { title: isZh ? "商品名称" : "Product", dataIndex: "title", key: "title" },
      {
        title: isZh ? "店内分类" : "Store Category",
        key: "storeCategory",
        width: 220,
        render: (_, row) => {
          const binding = bindingMap.get(row.spuNo);
          const label = resolveBindingLabel(binding, categories, isZh);
          return (
            <Tag color={binding?.storeCategoryId ? "gold" : "default"} style={{ maxWidth: 190, overflow: "hidden" }}>
              {label}
            </Tag>
          );
        }
      },
      {
        title: isZh ? "类目ID" : "Category ID",
        dataIndex: "categoryId",
        key: "categoryId",
        width: 120,
        render: (value: number) => (value > 0 ? value : "-")
      },
      {
        title: isZh ? "售价区间" : "Price",
        key: "price",
        width: 170,
        render: (_, row) => toPriceRange(row.minSalePrice, row.maxSalePrice)
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "spuStatus",
        key: "status",
        width: 130,
        render: (value: SpuStatusCode) => <Tag color={statusColor(value)}>{statusText(value, isZh)}</Tag>
      },
      {
        title: isZh ? "库存状态" : "Stock",
        dataIndex: "spuStockStatus",
        key: "spuStockStatus",
        width: 130,
        render: (value: StockStatusCode) => <Tag color={stockStatusColor(value)}>{stockStatusText(value, isZh)}</Tag>
      },
      {
        title: isZh ? "更新时间" : "Updated",
        dataIndex: "updatedAt",
        key: "updatedAt",
        width: 180,
        render: (value?: string) => formatDateTime(value)
      },
      {
        title: isZh ? "操作" : "Action",
        key: "action",
        width: 180,
        render: (_, row) => (
          <Space size={6}>
            <Button size="small" icon={<EyeOutlined />} onClick={() => setDetailSpuNo(row.spuNo)}>
              {isZh ? "详情" : "Detail"}
            </Button>
            <Button
              size="small"
              danger
              icon={<DeleteOutlined />}
              disabled={isDeleting || !canDeleteProduct(row.spuStatus)}
              onClick={() => confirmDeleteTargets([row])}
            >
              {isZh ? "删除" : "Delete"}
            </Button>
          </Space>
        )
      }
    ],
    [bindingMap, categories, isDeleting, isZh]
  );

  const tableLoading = workbenchQuery.isLoading || categoriesQuery.isLoading || (filterActive ? filteredUniverseQuery.isLoading : productsQuery.isLoading) || bindingsQuery.isLoading;
  const tableError = workbenchQuery.isError || categoriesQuery.isError || (filterActive ? filteredUniverseQuery.isError : productsQuery.isError) || bindingsQuery.isError;

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "商品管理" : "Products"}</Title>
          <Text type="secondary">
            {isZh ? "查看商品状态、价格与库存，并按店内分类做运营管理。" : "Manage product status, pricing, stock, and in-shop category placement."}
          </Text>
        </div>
        <Space wrap>
          <Button icon={<ReloadOutlined />} onClick={() => void refreshAll()}>
            {isZh ? "刷新" : "Refresh"}
          </Button>
          <Button icon={<DownloadOutlined />} onClick={exportCurrentPage}>
            {isZh ? "导出当前页" : "Export Page"}
          </Button>
          <Button type="primary" href="/seller/publish">
            {isZh ? "新建商品草稿" : "Create Draft"}
          </Button>
          <Button danger icon={<DeleteOutlined />} disabled={!deletableSelectedProducts.length || isDeleting} onClick={handleBatchDelete}>
            {isZh
              ? `删除选中 (${deletableSelectedProducts.length})`
              : `Delete selected (${deletableSelectedProducts.length})`}
          </Button>
        </Space>
        {selectedSpuNos.length ? (
          <Text type="secondary">
            {isZh
              ? `已选择 ${selectedSpuNos.length} 个商品，可批量改店内分类；删除仅支持草稿/驳回/下架商品。`
              : `${selectedSpuNos.length} item(s) selected. Bulk store-category update is available; delete works only for draft/rejected/off-shelf items.`}
          </Text>
        ) : null}
      </header>

      {filterActive ? (
        <Alert
          type="info"
          showIcon
          message={isZh ? "当前正在按店内分类筛选" : "Store category filter is active"}
          description={
            isZh
              ? "筛选结果会基于当前关键词和状态条件拉取匹配商品全集，再按店内分类分页展示。"
              : "The page loads all products matching the current keyword/status filters, then paginates them by store category."
          }
        />
      ) : null}

      <Card>
        <div className="seller-toolbar" style={{ display: "grid", gap: 12, gridTemplateColumns: "minmax(220px, 1.6fr) repeat(3, minmax(180px, 1fr))" }}>
          <Input
            allowClear
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder={isZh ? "搜索商品名 / SPU" : "Search product / SPU"}
          />
          <Select
            value={statusFilter}
            onChange={(value) => setStatusFilter(value)}
            options={statusOptions.map((item) => ({ label: isZh ? item.zh : item.en, value: item.value }))}
          />
          <Select
            value={storeCategoryFilter}
            onChange={(value) => setStoreCategoryFilter(value)}
            options={categoryFilterOptions}
            suffixIcon={<TagsOutlined />}
          />
          <Space.Compact style={{ width: "100%" }}>
            <Select
              value={batchCategoryId}
              allowClear
              placeholder={isZh ? "批量设置店内分类" : "Bulk set store category"}
              options={[
                { value: 0, label: isZh ? "改为未分类" : "Set unclassified" },
                ...leafCategoryOptions
              ]}
              onChange={(value) => setBatchCategoryId(typeof value === "number" ? value : null)}
              style={{ width: "100%" }}
            />
            <Button
              type="primary"
              loading={isUpdatingBatchCategory}
              disabled={!selectedSpuNos.length || batchCategoryId === null}
              onClick={() => void handleBatchUpdateStoreCategory()}
            >
              {isZh ? "应用" : "Apply"}
            </Button>
          </Space.Compact>
        </div>

        {!shop && !workbenchQuery.isLoading ? (
          <Empty description={isZh ? "当前账号没有可操作的店铺" : "This account has no available shop"} />
        ) : null}

        {tableLoading && !rows.length ? <Skeleton active paragraph={{ rows: 8 }} /> : null}

        <Table<SellerProductSpu>
          rowKey={(record) => record.spuNo}
          rowSelection={{
            selectedRowKeys: selectedSpuNos,
            onChange: (keys) => setSelectedSpuNos(keys as string[]),
            getCheckboxProps: () => ({
              disabled: isDeleting || isUpdatingBatchCategory
            })
          }}
          loading={tableLoading}
          columns={columns}
          dataSource={rows}
          onChange={onTableChange}
          locale={{
            emptyText: tableError
              ? isZh
                ? "加载失败，请刷新重试"
                : "Load failed, please retry"
              : filterActive && !rows.length
                ? isZh
                  ? "当前分类下没有商品"
                  : "No products in this store category"
                : undefined
          }}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: false
          }}
        />
      </Card>

      <Drawer
        width={680}
        open={Boolean(detailSpuNo)}
        onClose={() => setDetailSpuNo(null)}
        title={isZh ? "商品详情" : "Product Detail"}
      >
        {detailQuery.isLoading ? <Skeleton active paragraph={{ rows: 8 }} /> : null}
        {!detailQuery.isLoading && !detailQuery.data?.spu ? <Empty description={isZh ? "暂无详情数据" : "No detail data"} /> : null}

        {!detailQuery.isLoading && detailQuery.data?.spu ? (
          <section className="seller-page" style={{ gap: 10 }}>
            <Card size="small">
              <Space direction="vertical" size={4}>
                <Text strong>{detailQuery.data.spu.title || "-"}</Text>
                <Text type="secondary">SPU: {detailQuery.data.spu.spuNo}</Text>
                <Text type="secondary">
                  {isZh ? "状态：" : "Status: "}
                  {statusText(detailQuery.data.spu.spuStatus, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "店内分类：" : "Store category: "}
                  {resolveBindingLabel(bindingMap.get(detailQuery.data.spu.spuNo), categories, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "更新时间：" : "Updated: "}
                  {formatDateTime(detailQuery.data.spu.updatedAt)}
                </Text>
              </Space>
            </Card>

            <Card size="small" title={isZh ? "SKU 列表" : "SKU List"}>
              <Table<SellerProductSku>
                rowKey={(record) => record.skuNo}
                pagination={false}
                dataSource={detailQuery.data.skus}
                locale={{ emptyText: isZh ? "暂无 SKU" : "No SKUs" }}
                columns={[
                  { title: "SKU", dataIndex: "skuNo", key: "skuNo", width: 160 },
                  { title: isZh ? "名称" : "Name", dataIndex: "skuName", key: "skuName" },
                  {
                    title: isZh ? "售价" : "Sale Price",
                    dataIndex: "salePrice",
                    key: "salePrice",
                    width: 120,
                    render: (value: number) => `¥${value}`
                  },
                  {
                    title: isZh ? "市场价" : "Market Price",
                    dataIndex: "marketPrice",
                    key: "marketPrice",
                    width: 120,
                    render: (value: number) => `¥${value}`
                  },
                  {
                    title: isZh ? "库存状态" : "Stock",
                    dataIndex: "stockStatus",
                    key: "stockStatus",
                    width: 120,
                    render: (value: StockStatusCode) => <Tag color={stockStatusColor(value)}>{stockStatusText(value, isZh)}</Tag>
                  }
                ]}
              />
            </Card>
          </section>
        ) : null}
      </Drawer>
    </section>
  );
}
