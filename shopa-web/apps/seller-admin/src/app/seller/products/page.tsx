"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
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
import { buildSellerProductEditPath, canResumeSellerProductEditing } from "@/features/catalog/draft-edit";
import { applySellerProductLiveStock, resolveSellerProductLiveStockStatus } from "@/features/catalog/product-stock-view";
import { formatSellerProductPriceRange, resolveSellerProductsScope } from "@/features/catalog/products-view";
import { SellerProductSku, SellerProductSpu, SpuStatusCode, StockStatusCode } from "@/features/catalog/types";
import { batchGetSkuInventory } from "@/features/inventory/api";
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

const productStatusOptions: Array<{ value: StatusFilter; zh: string; en: string }> = [
  { value: "ALL", zh: "全部状态", en: "All status" },
  { value: "SPU_STATUS_DRAFT", zh: "草稿", en: "Draft" },
  { value: "SPU_STATUS_REVIEWING", zh: "审核中", en: "Reviewing" },
  { value: "SPU_STATUS_APPROVED", zh: "已通过", en: "Approved" },
  { value: "SPU_STATUS_ON_SHELF", zh: "在架", en: "On shelf" },
  { value: "SPU_STATUS_OFF_SHELF", zh: "下架", en: "Off shelf" },
  { value: "SPU_STATUS_REJECTED", zh: "驳回", en: "Rejected" }
];

function getStatusText(status: SpuStatusCode, isZh: boolean): string {
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
      return isZh ? "未指定" : "Unspecified";
  }
}

function getStockStatusText(status: StockStatusCode, isZh: boolean): string {
  switch (status) {
    case "STOCK_STATUS_IN_STOCK":
      return isZh ? "有货" : "In stock";
    case "STOCK_STATUS_OUT_OF_STOCK":
      return isZh ? "缺货" : "Out of stock";
    default:
      return isZh ? "未知" : "Unknown";
  }
}

function getStoreCategoryFilterOptions(categories: SellerStoreCategory[], isZh: boolean) {
  const options: Array<{ label: string; value: StoreCategoryFilterValue }> = [
    { label: isZh ? "全部店内分类" : "All store categories", value: "ALL" },
    { label: isZh ? "未分类" : "Unclassified", value: "UNCLASSIFIED" }
  ];

  categories.forEach((category) => {
    options.push({
      label: `${isZh ? "一级" : "L1"} / ${category.name}`,
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

function getStoreCategoryBindingLabel(
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

async function loadVisibleProductLiveStock(products: SellerProductSpu[]) {
  const uniqueProducts = Array.from(new Map(products.map((product) => [product.spuNo, product])).values()).filter((product) => product.spuNo);
  if (!uniqueProducts.length) {
    return {
      snapshots: new Map(),
      inventoryQueryable: true
    };
  }

  const detailPairs = await Promise.all(
    uniqueProducts.map(async (product) => ({
      spuNo: product.spuNo,
      detail: await getSellerProductDetail(product.spuNo)
    }))
  );
  const detailMap = new Map(detailPairs.map((item) => [item.spuNo, item.detail]));
  const skuNos = Array.from(
    new Set(detailPairs.flatMap((item) => item.detail.skus.map((sku) => sku.skuNo).filter(Boolean)))
  );

  let inventoryQueryable = true;
  let inventoryStocks: Awaited<ReturnType<typeof batchGetSkuInventory>> = [];
  if (skuNos.length > 0) {
    try {
      inventoryStocks = await batchGetSkuInventory(skuNos);
    } catch {
      inventoryQueryable = false;
    }
  }

  const snapshots = new Map(
    uniqueProducts.map((product) => [
      product.spuNo,
      resolveSellerProductLiveStockStatus({
        spu: product,
        skus: detailMap.get(product.spuNo)?.skus,
        inventoryStocks
      })
    ])
  );

  return {
    snapshots,
    inventoryQueryable
  };
}

export default function ProductsPage() {
  const searchParams = useSearchParams();
  const scope = resolveSellerProductsScope(searchParams.get("scope"));
  const isDraftScope = scope === "drafts";
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [keyword, setKeyword] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>(isDraftScope ? "SPU_STATUS_DRAFT" : "ALL");
  const [storeCategoryFilter, setStoreCategoryFilter] = useState<StoreCategoryFilterValue>("ALL");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [detailSpuNo, setDetailSpuNo] = useState<string | null>(null);
  const [selectedSpuNos, setSelectedSpuNos] = useState<string[]>([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const [batchCategoryId, setBatchCategoryId] = useState<number | null>(null);
  const [isUpdatingBatchCategory, setIsUpdatingBatchCategory] = useState(false);

  const filterActive = storeCategoryFilter !== "ALL";
  const effectiveStatusFilter: StatusFilter = isDraftScope ? "SPU_STATUS_DRAFT" : statusFilter;

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
    queryKey: ["seller-products", page, pageSize, keyword, effectiveStatusFilter, scope],
    queryFn: async () =>
      listSellerProducts({
        page,
        pageSize,
        keyword,
        statuses: effectiveStatusFilter === "ALL" ? undefined : [effectiveStatusFilter]
      }),
    enabled: !filterActive
  });

  const filteredUniverseQuery = useQuery({
    queryKey: ["seller-products-universe", keyword, effectiveStatusFilter, storeCategoryFilter, scope],
    queryFn: () => listProductsAcrossPages({ keyword, statusFilter: effectiveStatusFilter }),
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
  const visibleSpuKey = useMemo(() => rows.map((item) => item.spuNo).join("|"), [rows]);
  const visibleStockQuery = useQuery({
    queryKey: ["seller-products-live-stock", visibleSpuKey],
    queryFn: () => loadVisibleProductLiveStock(rows),
    enabled: rows.length > 0,
    staleTime: 10_000
  });
  const displayRows = useMemo(
    () => applySellerProductLiveStock(rows, visibleStockQuery.data?.snapshots ?? new Map()),
    [rows, visibleStockQuery.data?.snapshots]
  );

  const selectedProducts = useMemo(
    () => allRowsForSelection.filter((item) => selectedSpuNos.includes(item.spuNo)),
    [allRowsForSelection, selectedSpuNos]
  );
  const deletableSelectedProducts = useMemo(
    () => selectedProducts.filter((item) => canDeleteProduct(item.spuStatus)),
    [selectedProducts]
  );
  const leafCategoryOptions = useMemo(() => flattenLeafCategoryOptions(categories), [categories]);
  const categoryFilterOptions = useMemo(() => getStoreCategoryFilterOptions(categories, isZh), [categories, isZh]);

  useEffect(() => {
    setSelectedSpuNos((prev) => prev.filter((key) => allRowsForSelection.some((row) => row.spuNo === key)));
  }, [allRowsForSelection]);

  useEffect(() => {
    setPage(1);
  }, [keyword, effectiveStatusFilter, storeCategoryFilter, scope]);

  useEffect(() => {
    if (isDraftScope) {
      setStatusFilter("SPU_STATUS_DRAFT");
    }
  }, [isDraftScope]);

  const exportCurrentPage = () => {
    if (!displayRows.length) {
      notification.info({ message: isZh ? "当前页没有可导出的数据" : "No data to export on this page" });
      return;
    }
    const head = ["spu_no", "title", "store_category", "status", "stock_status", "price_range", "updated_at"];
    const body = displayRows.map((item) =>
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

  const exportCurrentPageClean = () => {
    if (!displayRows.length) {
      notification.info({ message: isZh ? "当前页面没有可导出的数据" : "No data to export on this page" });
      return;
    }
    const head = ["spu_no", "title", "store_category", "status", "stock_status", "price_range", "updated_at"];
    const body = displayRows.map((item) =>
      [
        item.spuNo,
        item.title,
        getStoreCategoryBindingLabel(bindingMap.get(item.spuNo), categories, false),
        getStatusText(item.spuStatus, false),
        getStockStatusText(item.spuStockStatus, false),
        formatSellerProductPriceRange(item.minSalePrice, item.maxSalePrice),
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
          const label = getStoreCategoryBindingLabel(binding, categories, isZh);
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
        render: (_, row) => formatSellerProductPriceRange(row.minSalePrice, row.maxSalePrice)
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "spuStatus",
        key: "status",
        width: 130,
        render: (value: SpuStatusCode) => <Tag color={statusColor(value)}>{getStatusText(value, isZh)}</Tag>
      },
      {
        title: isZh ? "库存状态" : "Stock",
        dataIndex: "spuStockStatus",
        key: "spuStockStatus",
        width: 130,
        render: (value: StockStatusCode) => <Tag color={stockStatusColor(value)}>{getStockStatusText(value, isZh)}</Tag>
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

  const tableColumns = useMemo<ColumnsType<SellerProductSpu>>(
    () =>
      columns.map((column) => {
        if (column.key === "title") {
          return {
            ...column,
            render: (value: string, row: SellerProductSpu) => {
              const editable = canResumeSellerProductEditing(row.spuStatus);
              const href = buildSellerProductEditPath(row.spuNo);
              return (
                <Space direction="vertical" size={2}>
                  <Text strong>{value || row.spuNo || "-"}</Text>
                  {editable ? <Link href={href}>{isZh ? "继续编辑草稿" : "Continue editing"}</Link> : null}
                </Space>
              );
            }
          };
        }

        if (column.key === "action") {
          return {
            ...column,
            width: 260,
            render: (_: unknown, row: SellerProductSpu) => {
              const editable = canResumeSellerProductEditing(row.spuStatus);
              const href = buildSellerProductEditPath(row.spuNo);
              return (
                <Space size={6} wrap>
                  <Button size="small" icon={<EyeOutlined />} onClick={() => setDetailSpuNo(row.spuNo)}>
                    {isZh ? "璇︽儏" : "Detail"}
                  </Button>
                  {editable ? (
                    <Button size="small" type="primary" href={href}>
                      {isZh ? "继续编辑" : "Edit"}
                    </Button>
                  ) : null}
                  <Button
                    size="small"
                    danger
                    icon={<DeleteOutlined />}
                    disabled={isDeleting || !canDeleteProduct(row.spuStatus)}
                    onClick={() => confirmDeleteTargets([row])}
                  >
                    {isZh ? "鍒犻櫎" : "Delete"}
                  </Button>
                </Space>
              );
            }
          };
        }

        return column;
      }),
    [columns, isDeleting, isZh, confirmDeleteTargets]
  );

  const tableLoading =
    workbenchQuery.isLoading ||
    categoriesQuery.isLoading ||
    (filterActive ? filteredUniverseQuery.isLoading : productsQuery.isLoading) ||
    bindingsQuery.isLoading ||
    (rows.length > 0 && visibleStockQuery.isLoading);
  const tableError = workbenchQuery.isError || categoriesQuery.isError || (filterActive ? filteredUniverseQuery.isError : productsQuery.isError) || bindingsQuery.isError;

  const handleDeleteProductsClean = async (targets: SellerProductSpu[]) => {
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

  const confirmDeleteTargetsClean = (targets: SellerProductSpu[]) => {
    if (!targets.length) {
      return;
    }
    Modal.confirm({
      title: isZh ? "确认删除草稿" : "Confirm delete draft",
      content: isZh
        ? `将删除 ${targets.length} 个商品。仅草稿、驳回、下架状态允许删除。`
        : `This will delete ${targets.length} products. Only draft, rejected, or off-shelf items can be deleted.`,
      onOk: () => handleDeleteProductsClean(targets)
    });
  };

  const handleBatchDeleteClean = () => {
    confirmDeleteTargetsClean(deletableSelectedProducts);
  };

  const handleBatchUpdateStoreCategoryClean = async () => {
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
          failed.push(`${product.spuNo}: ${error instanceof Error ? error.message : isZh ? "更新失败" : "Update failed"}`);
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

  const displayColumns = useMemo<ColumnsType<SellerProductSpu>>(
    () => [
      { title: "SPU", dataIndex: "spuNo", key: "spuNo", width: 150 },
      {
        title: isZh ? "商品名称" : "Product",
        dataIndex: "title",
        key: "title",
        render: (value: string, row) => {
          const editable = canResumeSellerProductEditing(row.spuStatus);
          const href = buildSellerProductEditPath(row.spuNo);
          return (
            <Space direction="vertical" size={2}>
              <Text strong>{value || row.spuNo || "-"}</Text>
              {editable ? <Link href={href}>{isZh ? "继续编辑草稿" : "Continue editing"}</Link> : null}
            </Space>
          );
        }
      },
      {
        title: isZh ? "店内分类" : "Store Category",
        key: "storeCategory",
        width: 220,
        render: (_, row) => {
          const binding = bindingMap.get(row.spuNo);
          return (
            <Tag color={binding?.storeCategoryId ? "gold" : "default"} style={{ maxWidth: 190, overflow: "hidden" }}>
              {getStoreCategoryBindingLabel(binding, categories, isZh)}
            </Tag>
          );
        }
      },
      {
        title: isZh ? "类目 ID" : "Category ID",
        dataIndex: "categoryId",
        key: "categoryId",
        width: 120,
        render: (value: number) => (value > 0 ? value : "-")
      },
      {
        title: isZh ? "售价区间（元）" : "Price (CNY)",
        key: "price",
        width: 180,
        render: (_, row) => formatSellerProductPriceRange(row.minSalePrice, row.maxSalePrice)
      },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "spuStatus",
        key: "status",
        width: 130,
        render: (value: SpuStatusCode) => <Tag color={statusColor(value)}>{getStatusText(value, isZh)}</Tag>
      },
      {
        title: isZh ? "库存状态" : "Stock",
        dataIndex: "spuStockStatus",
        key: "spuStockStatus",
        width: 130,
        render: (value: StockStatusCode) => <Tag color={stockStatusColor(value)}>{getStockStatusText(value, isZh)}</Tag>
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
        width: 260,
        render: (_, row) => {
          const editable = canResumeSellerProductEditing(row.spuStatus);
          const href = buildSellerProductEditPath(row.spuNo);
          return (
            <Space size={6} wrap>
              <Button size="small" icon={<EyeOutlined />} onClick={() => setDetailSpuNo(row.spuNo)}>
                {isZh ? "详情" : "Detail"}
              </Button>
              {editable ? (
                <Button size="small" type="primary" href={href}>
                  {isZh ? "继续编辑" : "Edit"}
                </Button>
              ) : null}
              <Button
                size="small"
                danger
                icon={<DeleteOutlined />}
                disabled={isDeleting || !canDeleteProduct(row.spuStatus)}
                onClick={() => confirmDeleteTargetsClean([row])}
              >
                {isZh ? "删除" : "Delete"}
              </Button>
            </Space>
          );
        }
      }
    ],
    [bindingMap, categories, isDeleting, isZh]
  );

  const pageTitle = isDraftScope ? (isZh ? "商品草稿箱" : "Draft Box") : isZh ? "商品管理" : "Products";
  const pageDescription = isDraftScope
    ? isZh
      ? "这里只展示尚未提交审核的商品草稿，价格按元展示。"
      : "Only products not yet submitted for review are shown here, with prices displayed in CNY."
    : isZh
      ? "查看商品状态、售价和库存，并按店内分类做运营管理。"
      : "Manage product status, pricing, stock, and in-shop category placement.";
  const detailSpu = detailQuery.data?.spu ?? null;
  const detailSkus = detailQuery.data?.skus ?? [];

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{pageTitle}</Title>
          <Text type="secondary">{pageDescription}</Text>
        </div>
        <Space wrap>
          <Button icon={<ReloadOutlined />} onClick={() => void refreshAll()}>
            {isZh ? "刷新" : "Refresh"}
          </Button>
          <Button icon={<DownloadOutlined />} onClick={exportCurrentPageClean}>
            {isZh ? "导出当前页" : "Export Page"}
          </Button>
          <Button type="primary" href="/seller/publish">
            {isZh ? "新建商品草稿" : "Create Draft"}
          </Button>
          <Button danger icon={<DeleteOutlined />} disabled={!deletableSelectedProducts.length || isDeleting} onClick={handleBatchDeleteClean}>
            {isZh ? `删除选中（${deletableSelectedProducts.length}）` : `Delete selected (${deletableSelectedProducts.length})`}
          </Button>
        </Space>
        {selectedSpuNos.length ? (
          <Text type="secondary">
            {isZh
              ? `已选择 ${selectedSpuNos.length} 个商品，可批量修改店内分类；删除仅支持草稿、驳回、下架商品。`
              : `${selectedSpuNos.length} item(s) selected. Bulk category updates are available; delete works only for draft, rejected, or off-shelf products.`}
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
              ? "筛选结果会先按当前关键字和状态拉取商品，再按店内分类分页展示。"
              : "The page first loads products by keyword and status, then paginates them by store category."
          }
        />
      ) : null}

      {visibleStockQuery.data && !visibleStockQuery.data.inventoryQueryable ? (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 12 }}
          message={isZh ? "当前正在使用商品侧库存兜底" : "Using catalog stock fallback"}
          description={
            isZh
              ? "库存服务暂时不可用，商品管理页已回退到商品聚合库存状态；恢复后会自动切回实时 SKU 库存。"
              : "The inventory service is temporarily unavailable, so this page is falling back to catalog stock status until live SKU inventory recovers."
          }
        />
      ) : null}

      <Card>
        <div
          className="seller-toolbar"
          style={{
            display: "grid",
            gap: 12,
            gridTemplateColumns: isDraftScope
              ? "minmax(220px, 1.8fr) minmax(180px, 1fr) minmax(180px, 1fr)"
              : "minmax(220px, 1.6fr) repeat(3, minmax(180px, 1fr))"
          }}
        >
          <Input
            allowClear
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder={isZh ? "搜索商品名称 / SPU" : "Search product / SPU"}
          />
          {!isDraftScope ? (
            <Select
              value={statusFilter}
              onChange={(value) => setStatusFilter(value)}
              options={productStatusOptions.map((item) => ({ label: isZh ? item.zh : item.en, value: item.value }))}
            />
          ) : null}
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
                { value: 0, label: isZh ? "设为未分类" : "Set unclassified" },
                ...leafCategoryOptions
              ]}
              onChange={(value) => setBatchCategoryId(typeof value === "number" ? value : null)}
              style={{ width: "100%" }}
            />
            <Button
              type="primary"
              loading={isUpdatingBatchCategory}
              disabled={!selectedSpuNos.length || batchCategoryId === null}
              onClick={() => void handleBatchUpdateStoreCategoryClean()}
            >
              {isZh ? "应用" : "Apply"}
            </Button>
          </Space.Compact>
        </div>

        {!shop && !workbenchQuery.isLoading ? (
          <Empty description={isZh ? "当前账号没有可操作的店铺" : "This account has no available shop"} />
        ) : null}

        {tableLoading && !displayRows.length ? <Skeleton active paragraph={{ rows: 8 }} /> : null}

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
          columns={displayColumns}
          dataSource={displayRows}
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
                : isDraftScope
                  ? isZh
                    ? "当前没有未提交审核的草稿"
                    : "No unsubmitted drafts"
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
        extra={
          detailSpu && canResumeSellerProductEditing(detailSpu!.spuStatus) ? (
            <Button
              type="primary"
              href={buildSellerProductEditPath(detailSpu!.spuNo)}
              onClick={() => setDetailSpuNo(null)}
            >
              {isZh ? "继续编辑" : "Edit"}
            </Button>
          ) : null
        }
        title={isZh ? "商品详情" : "Product Detail"}
      >
        {detailQuery.isLoading ? <Skeleton active paragraph={{ rows: 8 }} /> : null}
        {!detailQuery.isLoading && !detailSpu ? (
          <Empty description={isZh ? "暂无详情数据" : "No detail data"} />
        ) : null}

        {!detailQuery.isLoading && detailSpu ? (
          <section className="seller-page" style={{ gap: 10 }}>
            <Card size="small">
              <Space direction="vertical" size={4}>
                <Text strong>{detailSpu!.title || "-"}</Text>
                <Text type="secondary">SPU: {detailSpu!.spuNo}</Text>
                <Text type="secondary">
                  {isZh ? "状态：" : "Status: "}
                  {getStatusText(detailSpu.spuStatus, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "店内分类：" : "Store category: "}
                  {getStoreCategoryBindingLabel(bindingMap.get(detailSpu.spuNo), categories, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "更新时间：" : "Updated: "}
                  {formatDateTime(detailSpu.updatedAt)}
                </Text>
                <Text type="secondary">
                  {isZh ? "售价区间：" : "Price range: "}
                  {formatSellerProductPriceRange(detailSpu.minSalePrice, detailSpu.maxSalePrice)}
                </Text>
              </Space>
            </Card>

            <Card size="small" title={isZh ? "SKU 列表" : "SKU List"}>
              <Table<SellerProductSku>
                rowKey={(record) => record.skuNo}
                pagination={false}
                dataSource={detailSkus}
                locale={{ emptyText: isZh ? "暂无 SKU" : "No SKUs" }}
                columns={[
                  { title: "SKU", dataIndex: "skuNo", key: "skuNo", width: 160 },
                  { title: isZh ? "名称" : "Name", dataIndex: "skuName", key: "skuName" },
                  {
                    title: isZh ? "售价（元）" : "Sale Price",
                    dataIndex: "salePrice",
                    key: "salePrice",
                    width: 120,
                    render: (value: number) => `¥${value.toFixed(2)}`
                  },
                  {
                    title: isZh ? "市场价（元）" : "Market Price",
                    dataIndex: "marketPrice",
                    key: "marketPrice",
                    width: 120,
                    render: (value: number) => `¥${value.toFixed(2)}`
                  },
                  {
                    title: isZh ? "库存状态" : "Stock",
                    dataIndex: "stockStatus",
                    key: "stockStatus",
                    width: 120,
                    render: (value: StockStatusCode) => <Tag color={stockStatusColor(value)}>{getStockStatusText(value, isZh)}</Tag>
                  }
                ]}
              />
            </Card>
          </section>
        ) : null}
      </Drawer>
    </section>
  );

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
          <Button icon={<DownloadOutlined />} onClick={exportCurrentPageClean}>
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
            options={productStatusOptions.map((item) => ({ label: isZh ? item.zh : item.en, value: item.value }))}
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
          columns={tableColumns}
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
        extra={
          detailSpu && canResumeSellerProductEditing(detailSpu!.spuStatus) ? (
            <Button
              type="primary"
              href={buildSellerProductEditPath(detailSpu!.spuNo)}
              onClick={() => setDetailSpuNo(null)}
            >
              {isZh ? "继续编辑" : "Edit"}
            </Button>
          ) : null
        }
        title={isZh ? "商品详情" : "Product Detail"}
      >
        {detailQuery.isLoading ? <Skeleton active paragraph={{ rows: 8 }} /> : null}
        {!detailQuery.isLoading && !detailQuery.data?.spu ? <Empty description={isZh ? "暂无详情数据" : "No detail data"} /> : null}

        {!detailQuery.isLoading && detailSpu ? (
          <section className="seller-page" style={{ gap: 10 }}>
            <Card size="small">
              <Space direction="vertical" size={4}>
                <Text strong>{detailSpu!.title || "-"}</Text>
                <Text type="secondary">SPU: {detailSpu!.spuNo}</Text>
                <Text type="secondary">
                  {isZh ? "状态：" : "Status: "}
                  {statusText(detailSpu!.spuStatus, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "店内分类：" : "Store category: "}
                  {resolveBindingLabel(bindingMap.get(detailSpu!.spuNo), categories, isZh)}
                </Text>
                <Text type="secondary">
                  {isZh ? "更新时间：" : "Updated: "}
                  {formatDateTime(detailSpu!.updatedAt)}
                </Text>
              </Space>
            </Card>

            <Card size="small" title={isZh ? "SKU 列表" : "SKU List"}>
              <Table<SellerProductSku>
                rowKey={(record) => record.skuNo}
                pagination={false}
                dataSource={detailSkus}
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
