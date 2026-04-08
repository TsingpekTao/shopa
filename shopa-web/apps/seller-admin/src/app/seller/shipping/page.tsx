"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Alert, Button, Card, Empty, Input, Select, Skeleton, Space, Statistic, Table, Tag, Typography, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { ReloadOutlined, SearchOutlined } from "@ant-design/icons";
import { listSellerProducts } from "@/features/catalog/api";
import { buildSellerChatShopOptions } from "@/features/chat/shop-selection";
import { getSellerCurrentShopNo, setSellerCurrentShopNo } from "@/features/navigation/nav-badges";
import { listSellerSubOrders, markSubOrderShipped } from "@/features/order/api";
import type { SellerOrderItem, SellerOrderSub } from "@/features/order/types";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const { Text, Title } = Typography;

const SHIPPING_QUERY_KEY = "seller-shipping";
const SELLER_SHIPPING_LAST_SHOP_KEY = "seller-shipping-last-shop-no";
const SHIPPING_STATUSES = ["SUB_ORDER_STATUS_PAID", "SUB_ORDER_STATUS_WAIT_SHIP"];

function formatMoney(amount: number): string {
  return `CNY ${(amount / 100).toFixed(2)}`;
}

function formatDateTime(raw?: string): string {
  if (!raw) {
    return "-";
  }
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }
  return date.toLocaleString("zh-CN", { hour12: false });
}

function subOrderStatusMeta(status: string) {
  switch ((status || "").toUpperCase()) {
    case "SHIPPED":
      return { text: "已发货", color: "success" as const };
    case "WAIT_SHIP":
      return { text: "待发货", color: "processing" as const };
    case "PAID":
      return { text: "已付款", color: "gold" as const };
    case "PENDING_PAY":
      return { text: "待付款", color: "warning" as const };
    case "CANCELED":
      return { text: "已取消", color: "default" as const };
    case "CLOSED":
      return { text: "已关闭", color: "default" as const };
    default:
      return { text: status || "未知状态", color: "default" as const };
  }
}

function summarizeItems(items: SellerOrderItem[]): string {
  return items
    .map((item) => [item.spuTitle, item.skuName, item.skuNo].find(Boolean) ?? "")
    .filter(Boolean)
    .join(" / ");
}

export default function SellerShippingPage() {
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();
  const [selectedShopNo, setSelectedShopNo] = useState("");
  const [persistedShopNo, setPersistedShopNo] = useState("");
  const [keywordInput, setKeywordInput] = useState("");
  const [keyword, setKeyword] = useState("");
  const [processingSubOrderNo, setProcessingSubOrderNo] = useState("");

  const workbenchQuery = useQuery({
    queryKey: ["seller-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 60_000
  });

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    setPersistedShopNo(getSellerCurrentShopNo() || window.localStorage.getItem(SELLER_SHIPPING_LAST_SHOP_KEY)?.trim() || "");
  }, []);

  const workbenchShops = useMemo(() => {
    return (workbenchQuery.data?.shops ?? [])
      .map((shop) => ({
        shopNo: String(shop.shopNo ?? "").trim(),
        shopName: String(shop.shopDisplayName || shop.shopName || shop.shopNo || "").trim(),
        shopStatusCode: String(shop.shopStatusCode ?? "").trim()
      }))
      .filter((shop) => shop.shopNo);
  }, [workbenchQuery.data?.shops]);

  const fallbackCatalogShopsQuery = useQuery({
    queryKey: ["seller-shipping-fallback-shops"],
    queryFn: async () => {
      const result = await listSellerProducts({ page: 1, pageSize: 100 });
      return Array.from(new Set(result.products.map((product) => product.shopNo.trim()).filter(Boolean)));
    },
    enabled: workbenchShops.length === 0 || Boolean(workbenchQuery.data?.partial),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const shops = useMemo(() => {
    return buildSellerChatShopOptions({
      workbenchShops,
      catalogShopNos: fallbackCatalogShopsQuery.data ?? [],
      persistedShopNo
    });
  }, [fallbackCatalogShopsQuery.data, persistedShopNo, workbenchShops]);

  const normalizedSelectedShopNo = selectedShopNo.trim();

  useEffect(() => {
    if (!shops.length) {
      setSelectedShopNo("");
      return;
    }
    setSelectedShopNo((current) => {
      const normalizedCurrent = current.trim();
      if (normalizedCurrent && shops.some((shop) => shop.shopNo === normalizedCurrent)) {
        return normalizedCurrent;
      }
      return shops[0].shopNo;
    });
  }, [shops]);

  useEffect(() => {
    if (!normalizedSelectedShopNo || typeof window === "undefined") {
      return;
    }
    setSellerCurrentShopNo(normalizedSelectedShopNo);
    window.localStorage.setItem(SELLER_SHIPPING_LAST_SHOP_KEY, normalizedSelectedShopNo);
    setPersistedShopNo(normalizedSelectedShopNo);
  }, [normalizedSelectedShopNo]);

  const shippingQuery = useQuery({
    queryKey: [SHIPPING_QUERY_KEY, normalizedSelectedShopNo, keyword],
    queryFn: () =>
      listSellerSubOrders({
        shopNo: normalizedSelectedShopNo,
        statuses: SHIPPING_STATUSES,
        keyword,
        pageSize: 200
      }),
    enabled: Boolean(normalizedSelectedShopNo),
    staleTime: 15_000,
    refetchOnWindowFocus: false
  });

  const shipMutation = useMutation({
    mutationFn: async (subOrderNo: string) => {
      setProcessingSubOrderNo(subOrderNo);
      return markSubOrderShipped(subOrderNo);
    },
    onSuccess: async () => {
      messageApi.success("订单发货状态已更新");
      await queryClient.invalidateQueries({ queryKey: [SHIPPING_QUERY_KEY, normalizedSelectedShopNo] });
      await shippingQuery.refetch();
    },
    onError: (error) => {
      messageApi.error(error instanceof Error ? error.message : "发货失败");
    },
    onSettled: () => {
      setProcessingSubOrderNo("");
    }
  });

  const rows = shippingQuery.data?.subOrders ?? [];
  const totalPayable = rows.reduce((sum, row) => sum + row.amount.payableAmount, 0);
  const totalItems = rows.reduce((sum, row) => sum + row.items.reduce((itemSum, item) => itemSum + item.qty, 0), 0);

  const columns = useMemo<ColumnsType<SellerOrderSub>>(
    () => [
      {
        title: "订单信息",
        key: "order",
        width: 220,
        render: (_, row) => (
          <div className="seller-shipping-order-cell">
            <Text strong>{row.orderNo || "-"}</Text>
            <Text type="secondary">子单号：{row.subOrderNo || "-"}</Text>
            <Text type="secondary">店铺：{row.shopNo || "-"}</Text>
          </div>
        )
      },
      {
        title: "商品",
        key: "items",
        render: (_, row) => (
          <div className="seller-shipping-items-cell">
            <Text>{summarizeItems(row.items) || "-"}</Text>
            <Text type="secondary">共 {row.items.reduce((sum, item) => sum + item.qty, 0)} 件</Text>
          </div>
        )
      },
      {
        title: "应付金额",
        key: "amount",
        width: 120,
        render: (_, row) => <span className="seller-shipping-amount">{formatMoney(row.amount.payableAmount)}</span>
      },
      {
        title: "状态",
        key: "status",
        width: 104,
        render: (_, row) => {
          const meta = subOrderStatusMeta(row.subStatus);
          return <Tag color={meta.color}>{meta.text}</Tag>;
        }
      },
      {
        title: "买家备注",
        dataIndex: "buyerRemark",
        key: "buyerRemark",
        width: 180,
        render: (value: string) => value || "-"
      },
      {
        title: "下单时间",
        dataIndex: "createdAt",
        key: "createdAt",
        width: 156,
        render: (value: string) => formatDateTime(value)
      },
      {
        title: "操作",
        key: "action",
        width: 164,
        render: (_, row) => (
          <Space size={6} wrap>
            <Button
              type="primary"
              size="small"
              loading={shipMutation.isLoading && processingSubOrderNo === row.subOrderNo}
              onClick={() => shipMutation.mutate(row.subOrderNo)}
            >
              立即发货
            </Button>
            <Link href={`/seller/orders/${encodeURIComponent(row.subOrderNo)}`}>
              <Button size="small">订单详情</Button>
            </Link>
          </Space>
        )
      }
    ],
    [processingSubOrderNo, shipMutation.isLoading, shipMutation]
  );

  const emptyNode = workbenchQuery.isLoading ? (
    <Skeleton active />
  ) : shops.length === 0 ? (
    <Empty description="当前账号下暂无可发货店铺" />
  ) : (
    <Empty description="当前没有待发货订单" />
  );

  return (
    <section className="seller-page seller-shipping-page">
      {contextHolder}
      <header className="seller-page-head seller-shipping-head">
        <div className="seller-shipping-head-copy">
          <Title level={3}>发货中心</Title>
          <Text type="secondary">集中处理已付款订单，确认发货后买家侧会自动进入待收货。</Text>
        </div>
        <Space wrap className="seller-shipping-toolbar">
          <Select
            value={normalizedSelectedShopNo || undefined}
            placeholder="选择店铺"
            style={{ minWidth: 220 }}
            loading={workbenchQuery.isLoading}
            options={shops.map((shop) => ({
              label: shop.shopName,
              value: shop.shopNo
            }))}
            onChange={(value) => setSelectedShopNo(value.trim())}
          />
          <Input
            allowClear
            value={keywordInput}
            placeholder="搜索订单号或商品"
            style={{ width: 220 }}
            onChange={(event) => setKeywordInput(event.target.value)}
            onPressEnter={() => setKeyword(keywordInput.trim())}
          />
          <Button icon={<SearchOutlined />} onClick={() => setKeyword(keywordInput.trim())}>
            搜索
          </Button>
          <Button
            icon={<ReloadOutlined />}
            onClick={() => {
              void Promise.all([workbenchQuery.refetch(), fallbackCatalogShopsQuery.refetch(), shippingQuery.refetch()]);
            }}
          >
            刷新
          </Button>
        </Space>
      </header>

      {shippingQuery.isError ? (
        <Alert
          type="error"
          showIcon
          message="待发货订单加载失败"
          description={shippingQuery.error instanceof Error ? shippingQuery.error.message : "请稍后重试。"}
        />
      ) : null}

      <div className="seller-kpi-grid seller-shipping-kpi-grid">
        <Card className="seller-shipping-kpi-card" bordered={false}>
          <Statistic title="待发货子单" value={rows.length} />
        </Card>
        <Card className="seller-shipping-kpi-card" bordered={false}>
          <Statistic title="商品件数" value={totalItems} />
        </Card>
        <Card className="seller-shipping-kpi-card" bordered={false}>
          <Statistic title="待发货金额" value={formatMoney(totalPayable)} />
        </Card>
        <Card className="seller-shipping-kpi-card" bordered={false}>
          <Statistic title="当前店铺" value={normalizedSelectedShopNo || "-"} />
        </Card>
      </div>

      <Card bordered={false} className="seller-shipping-table-card">
        {normalizedSelectedShopNo ? (
          <Table
            className="seller-shipping-table"
            rowKey={(row) => row.subOrderNo || `${row.orderNo}-${row.shopNo}`}
            loading={shippingQuery.isLoading}
            columns={columns}
            dataSource={rows}
            pagination={false}
            size="small"
            locale={{ emptyText: emptyNode }}
          />
        ) : (
          emptyNode
        )}
      </Card>
    </section>
  );
}
