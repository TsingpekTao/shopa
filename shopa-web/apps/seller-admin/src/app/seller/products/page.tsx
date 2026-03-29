"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Button,
  Card,
  Drawer,
  Empty,
  Input,
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
import { DownloadOutlined, EyeOutlined, ReloadOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { getSellerProductDetail, listSellerProducts } from "@/features/catalog/api";
import { SellerProductSku, SellerProductSpu, SpuStatusCode, StockStatusCode } from "@/features/catalog/types";

const { Title, Text } = Typography;

type StatusFilter = "ALL" | SpuStatusCode;

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

export default function ProductsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [keyword, setKeyword] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [detailSpuNo, setDetailSpuNo] = useState<string | null>(null);

  const productsQuery = useQuery({
    queryKey: ["seller-products", page, pageSize, keyword, statusFilter],
    queryFn: async () =>
      listSellerProducts({
        page,
        pageSize,
        keyword,
        statuses: statusFilter === "ALL" ? undefined : [statusFilter]
      })
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

  const rows = productsQuery.data?.products ?? [];
  const total = productsQuery.data?.total ?? 0;

  const exportCurrentPage = () => {
    if (!rows.length) {
      notification.info({ message: isZh ? "当前页没有可导出的数据" : "No data to export on this page" });
      return;
    }
    const head = ["spu_no", "title", "status", "stock_status", "price_range", "updated_at"];
    const body = rows.map((item) =>
      [
        item.spuNo,
        item.title,
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

  const onTableChange = (pagination: TablePaginationConfig) => {
    setPage(pagination.current ?? 1);
    setPageSize(pagination.pageSize ?? 10);
  };

  const columns = useMemo<ColumnsType<SellerProductSpu>>(
    () => [
      { title: "SPU", dataIndex: "spuNo", key: "spuNo", width: 150 },
      { title: isZh ? "商品名称" : "Product", dataIndex: "title", key: "title" },
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
        width: 120,
        render: (_, row) => (
          <Button size="small" icon={<EyeOutlined />} onClick={() => setDetailSpuNo(row.spuNo)}>
            {isZh ? "详情" : "Detail"}
          </Button>
        )
      }
    ],
    [isZh]
  );

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "商品管理" : "Products"}</Title>
          <Text type="secondary">{isZh ? "查看商品状态、价格与库存，进行批量运营管理。" : "Manage product status, pricing and stock in one place."}</Text>
        </div>
        <Space wrap>
          <Button icon={<ReloadOutlined />} onClick={() => void productsQuery.refetch()}>
            {isZh ? "刷新" : "Refresh"}
          </Button>
          <Button icon={<DownloadOutlined />} onClick={exportCurrentPage}>
            {isZh ? "导出当前页" : "Export Page"}
          </Button>
          <Button type="primary" href="/seller/publish">
            {isZh ? "新建商品草稿" : "Create Draft"}
          </Button>
        </Space>
      </header>

      <Card>
        <div className="seller-toolbar">
          <Input
            allowClear
            value={keyword}
            onChange={(e) => {
              setKeyword(e.target.value);
              setPage(1);
            }}
            placeholder={isZh ? "搜索商品名 / SPU" : "Search product / SPU"}
          />
          <Select
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter(value);
              setPage(1);
            }}
            options={statusOptions.map((item) => ({ label: isZh ? item.zh : item.en, value: item.value }))}
          />
        </div>

        <Table<SellerProductSpu>
          rowKey={(record) => record.spuNo}
          loading={productsQuery.isLoading}
          columns={columns}
          dataSource={rows}
          onChange={onTableChange}
          locale={{ emptyText: productsQuery.isError ? (isZh ? "加载失败，请刷新重试" : "Load failed, please retry") : undefined }}
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
