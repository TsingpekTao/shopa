"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Button, Empty, Input, Space, Table, Tag } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useQuery } from "@tanstack/react-query";
import { fetchMerchantApplications } from "@/features/admin/api";
import type { MerchantApplication } from "@/features/admin/types";
import { AdminPage } from "@/components/admin-page";

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

export default function AdminShopsPage() {
  const searchParams = useSearchParams();
  const initialKeyword = searchParams.get("keyword")?.trim() ?? "";
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [keywordInput, setKeywordInput] = useState(initialKeyword);
  const [keyword, setKeyword] = useState(initialKeyword);

  useEffect(() => {
    const nextKeyword = searchParams.get("keyword")?.trim() ?? "";
    setKeywordInput(nextKeyword);
    setKeyword(nextKeyword);
    setPage(1);
  }, [searchParams]);

  const query = useQuery({
    queryKey: ["admin", "shops", page, pageSize, keyword],
    queryFn: () =>
      fetchMerchantApplications({
        page,
        pageSize,
        keyword: keyword || undefined,
        statuses: [4]
      })
  });

  const columns: ColumnsType<MerchantApplication> = [
    {
      title: "店铺编号",
      key: "shopNo",
      width: 220,
      render: (_, row) => row.shop?.shopNo || "-"
    },
    {
      title: "店铺名称",
      key: "shopName",
      width: 220,
      render: (_, row) => row.shop?.shopName || row.shop?.shopDisplayName || "-"
    },
    {
      title: "主体名称",
      key: "entityName",
      width: 220,
      render: (_, row) => row.entity?.entityName || "-"
    },
    {
      title: "联系人",
      key: "contactName",
      width: 140,
      render: (_, row) => row.entity?.contactName || "-"
    },
    {
      title: "联系电话",
      key: "contactPhone",
      width: 160,
      render: (_, row) => row.entity?.contactPhone || "-"
    },
    {
      title: "状态",
      key: "status",
      width: 120,
      render: () => <Tag color="success">已开通</Tag>
    },
    {
      title: "通过时间",
      key: "reviewedAt",
      width: 180,
      render: (_, row) => formatDateTime(row.reviewedAt)
    },
    {
      title: "操作",
      key: "actions",
      width: 160,
      render: (_, row) =>
        row.shop?.shopNo ? (
          <Link href={`/admin/shops/${row.shop.shopNo}/insights`} style={{ color: "#1677ff" }}>
            查看详情
          </Link>
        ) : (
          "-"
        )
    }
  ];

  return (
    <AdminPage title="已开通店铺" subtitle="查看平台已经审核通过并开通的店铺列表。">
      <Space wrap size={12} style={{ marginBottom: 12 }}>
        <Input
          allowClear
          value={keywordInput}
          onChange={(event) => setKeywordInput(event.target.value)}
          onPressEnter={() => {
            setPage(1);
            setKeyword(keywordInput.trim());
          }}
          style={{ width: 320 }}
          placeholder="搜索店铺名称或主体名称"
        />
        <Button
          type="primary"
          onClick={() => {
            setPage(1);
            setKeyword(keywordInput.trim());
          }}
        >
          搜索
        </Button>
        <Button
          onClick={() => {
            setPage(1);
            setPageSize(20);
            setKeywordInput("");
            setKeyword("");
          }}
        >
          重置
        </Button>
      </Space>

      {query.data?.applications?.length ? (
        <div className="admin-table">
          <Table
            rowKey={(row) => row.shop?.shopNo || row.applicationNo}
            loading={query.isLoading}
            dataSource={query.data?.applications ?? []}
            columns={columns}
            pagination={{
              current: query.data?.page ?? 1,
              pageSize: query.data?.pageSize ?? 20,
              total: query.data?.total ?? 0,
              showSizeChanger: true,
              pageSizeOptions: [10, 20, 50, 100]
            }}
            onChange={(pagination) => {
              setPage(pagination.current ?? 1);
              setPageSize(pagination.pageSize ?? pageSize);
            }}
            scroll={{ x: 1280 }}
          />
        </div>
      ) : (
        <Empty description={query.isLoading ? "加载中" : "暂无店铺数据"} />
      )}
    </AdminPage>
  );
}
