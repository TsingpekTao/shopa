"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Alert, Button, Card, InputNumber, Space, Table, Tag, Typography, notification } from "antd";
import type { ColumnsType } from "antd/es/table";
import { WarningOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { getSellerProductDetail, listSellerProducts } from "@/features/catalog/api";
import { SellerProductSku, StockStatusCode } from "@/features/catalog/types";
import { batchAdjustSellerStock, batchGetSkuInventory } from "@/features/inventory/api";
import { SellerInventoryStock } from "@/features/inventory/types";

const { Title, Text } = Typography;

type InventoryRow = {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  product: string;
  skuName: string;
  stockStatus: StockStatusCode;
  availableQty: number | null;
  totalQty: number | null;
  lockedQty: number | null;
  safeStock: number;
  updatedAt?: string;
};

type InventoryQueryResult = {
  rows: InventoryRow[];
  inventoryQueryable: boolean;
};

const SAFE_STOCK = 10;

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

function toStatusColor(status: StockStatusCode, availableQty: number | null, safeStock: number) {
  if (status === "STOCK_STATUS_OUT_OF_STOCK") {
    return "error";
  }
  if (availableQty !== null && availableQty <= safeStock) {
    return "warning";
  }
  if (status === "STOCK_STATUS_IN_STOCK") {
    return "success";
  }
  return "default";
}

function toStatusText(status: StockStatusCode, isZh: boolean) {
  if (status === "STOCK_STATUS_OUT_OF_STOCK") {
    return isZh ? "缺货" : "Out of stock";
  }
  if (status === "STOCK_STATUS_IN_STOCK") {
    return isZh ? "有货" : "In stock";
  }
  return isZh ? "未知" : "Unknown";
}

function chunkArray<T>(input: T[], size: number): T[][] {
  if (size <= 0) {
    return [input];
  }
  const out: T[][] = [];
  for (let i = 0; i < input.length; i += size) {
    out.push(input.slice(i, i + size));
  }
  return out;
}

async function loadInventoryRows(): Promise<InventoryQueryResult> {
  const listRes = await listSellerProducts({ page: 1, pageSize: 80 });
  const spus = listRes.products;
  if (!spus.length) {
    return { rows: [], inventoryQueryable: true };
  }

  const detailSettled = await Promise.allSettled(spus.map((item) => getSellerProductDetail(item.spuNo)));

  const skuMeta = new Map<
    string,
    {
      spuNo: string;
      shopNo: string;
      product: string;
      sku: SellerProductSku;
    }
  >();

  detailSettled.forEach((result) => {
    if (result.status !== "fulfilled" || !result.value.spu) {
      return;
    }
    const spu = result.value.spu;
    result.value.skus.forEach((sku) => {
      if (!sku.skuNo) {
        return;
      }
      skuMeta.set(sku.skuNo, {
        spuNo: spu.spuNo,
        shopNo: spu.shopNo,
        product: spu.title,
        sku
      });
    });
  });

  const skuNos = Array.from(skuMeta.keys());
  if (!skuNos.length) {
    return { rows: [], inventoryQueryable: true };
  }

  let inventoryQueryable = true;
  const stockMap = new Map<string, SellerInventoryStock>();

  try {
    const chunks = chunkArray(skuNos, 100);
    const chunkResults = await Promise.all(chunks.map((group) => batchGetSkuInventory(group)));
    chunkResults.flat().forEach((stock) => {
      stockMap.set(stock.skuNo, stock);
    });
  } catch {
    inventoryQueryable = false;
  }

  const rows: InventoryRow[] = skuNos.map((skuNo) => {
    const meta = skuMeta.get(skuNo);
    const stock = stockMap.get(skuNo);

    return {
      skuNo,
      spuNo: meta?.spuNo ?? "",
      shopNo: stock?.shopNo || meta?.shopNo || "",
      product: meta?.product || "-",
      skuName: meta?.sku.skuName || skuNo,
      stockStatus: (stock?.stockStatus as StockStatusCode) || meta?.sku.stockStatus || "STOCK_STATUS_UNSPECIFIED",
      availableQty: stock ? stock.availableQty : null,
      totalQty: stock ? stock.totalQty : null,
      lockedQty: stock ? stock.lockedQty : null,
      safeStock: SAFE_STOCK,
      updatedAt: stock?.updatedAt ?? meta?.sku.updatedAt
    };
  });

  rows.sort((a, b) => {
    const score = (row: InventoryRow) => {
      if (row.stockStatus === "STOCK_STATUS_OUT_OF_STOCK") return 3;
      if (row.availableQty !== null && row.availableQty <= row.safeStock) return 2;
      if (row.stockStatus === "STOCK_STATUS_IN_STOCK") return 1;
      return 0;
    };
    return score(b) - score(a);
  });

  return { rows, inventoryQueryable };
}

export default function InventoryPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [adjustValues, setAdjustValues] = useState<Record<string, number>>({});

  const inventoryQuery = useQuery({
    queryKey: ["seller-inventory-rows"],
    queryFn: loadInventoryRows
  });

  const adjustMutation = useMutation({
    mutationFn: async ({ row, delta }: { row: InventoryRow; delta: number }) => {
      if (!row.shopNo) {
        throw new Error(isZh ? "当前 SKU 缺少 shop_no，无法调整库存" : "shop_no missing for this sku");
      }
      return batchAdjustSellerStock({
        shopNo: row.shopNo,
        items: [
          {
            skuNo: row.skuNo,
            spuNo: row.spuNo,
            deltaTotalQty: delta,
            reasonCode: "ADJUST_REASON_CODE_SELLER_CORRECTION",
            remark: "seller-admin quick adjust"
          }
        ],
        idempotencyKey: `seller-admin:${row.skuNo}:${Date.now()}`
      });
    },
    onSuccess: async () => {
      notification.success({ message: isZh ? "库存调整成功" : "Inventory adjusted" });
      setAdjustValues({});
      await inventoryQuery.refetch();
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : isZh ? "库存调整失败" : "Adjust failed";
      notification.error({ message });
    }
  });

  const rows = inventoryQuery.data?.rows ?? [];

  const riskCount = useMemo(() => {
    return rows.filter((item) => item.stockStatus === "STOCK_STATUS_OUT_OF_STOCK" || (item.availableQty !== null && item.availableQty <= item.safeStock)).length;
  }, [rows]);

  const runAdjust = async (row: InventoryRow, delta?: number) => {
    const finalDelta = delta ?? adjustValues[row.skuNo] ?? 0;
    if (finalDelta === 0) {
      notification.info({ message: isZh ? "请输入非 0 的调整数量" : "Please input non-zero delta" });
      return;
    }
    await adjustMutation.mutateAsync({ row, delta: finalDelta });
  };

  const columns = useMemo<ColumnsType<InventoryRow>>(
    () => [
      { title: "SKU", dataIndex: "skuNo", key: "skuNo", width: 170 },
      { title: isZh ? "商品" : "Product", dataIndex: "product", key: "product", width: 180 },
      { title: isZh ? "SKU 名称" : "SKU Name", dataIndex: "skuName", key: "skuName" },
      {
        title: isZh ? "当前库存" : "Stock",
        key: "stock",
        width: 130,
        render: (_, row) => {
          if (row.availableQty === null) {
            return <Tag>{isZh ? "不可用" : "N/A"}</Tag>;
          }
          return <Tag color={toStatusColor(row.stockStatus, row.availableQty, row.safeStock)}>{row.availableQty}</Tag>;
        }
      },
      {
        title: isZh ? "库存状态" : "Status",
        dataIndex: "stockStatus",
        key: "stockStatus",
        width: 130,
        render: (value: StockStatusCode, row) => (
          <Tag color={toStatusColor(value, row.availableQty, row.safeStock)}>{toStatusText(value, isZh)}</Tag>
        )
      },
      {
        title: isZh ? "安全库存" : "Safe Stock",
        dataIndex: "safeStock",
        key: "safeStock",
        width: 110
      },
      {
        title: isZh ? "更新时间" : "Updated",
        dataIndex: "updatedAt",
        key: "updatedAt",
        width: 180,
        render: (value?: string) => formatDateTime(value)
      },
      {
        title: isZh ? "快速调整" : "Quick Adjust",
        key: "adjust",
        width: 260,
        render: (_, row) => (
          <Space>
            <InputNumber
              value={adjustValues[row.skuNo] ?? 0}
              onChange={(value) => {
                setAdjustValues((prev) => ({ ...prev, [row.skuNo]: Number(value ?? 0) }));
              }}
              min={-999999}
              max={999999}
              step={1}
            />
            <Button size="small" onClick={() => void runAdjust(row)} loading={adjustMutation.isLoading}>
              {isZh ? "提交" : "Apply"}
            </Button>
            <Button size="small" onClick={() => void runAdjust(row, 10)} loading={adjustMutation.isLoading}>
              +10
            </Button>
          </Space>
        )
      }
    ],
    [adjustMutation.isLoading, adjustValues, isZh]
  );

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "库存管理" : "Inventory"}</Title>
          <Text type="secondary">{isZh ? "查看预警 SKU，并通过后端接口快速调整库存。" : "Monitor risky SKUs and adjust stock via backend API."}</Text>
        </div>
        <Space>
          <Tag color={riskCount > 0 ? "error" : "success"} icon={<WarningOutlined />}>
            {isZh ? `预警 ${riskCount} 个SKU` : `${riskCount} risky SKUs`}
          </Tag>
          <Button type="primary" onClick={() => void inventoryQuery.refetch()}>
            {isZh ? "同步库存" : "Sync Inventory"}
          </Button>
        </Space>
      </header>

      {inventoryQuery.data && !inventoryQuery.data.inventoryQueryable ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "库存查询接口暂不可用" : "Inventory query unavailable"}
          description={
            isZh
              ? "当前页面已降级为商品侧库存状态视图；库存调整接口仍可使用。"
              : "The page is downgraded to catalog stock status view; adjust API is still available."
          }
        />
      ) : null}

      <Card>
        <Table<InventoryRow>
          rowKey={(record) => record.skuNo}
          loading={inventoryQuery.isLoading}
          columns={columns}
          dataSource={rows}
          pagination={{ pageSize: 10, showSizeChanger: false }}
          locale={{ emptyText: inventoryQuery.isError ? (isZh ? "加载失败，请刷新重试" : "Load failed, please retry") : undefined }}
        />
      </Card>
    </section>
  );
}
