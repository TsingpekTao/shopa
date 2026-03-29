"use client";

import { useMemo, useState } from "react";
import { Button, Modal, Space, Table, Tag, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useQuery } from "@tanstack/react-query";
import { AccessControl, useI18n } from "@shopa/ui";
import {
  approveProduct,
  fetchProductReviewDetail,
  fetchProductReviewTasks,
  forceOffShelf,
  freezeProduct,
  rejectProduct,
  unfreezeProduct
} from "@/features/admin/api";
import { Permissions } from "@/features/admin/permissions";
import { ProductReviewTask } from "@/features/admin/types";
import { AdminPage } from "@/components/admin-page";

function reviewStatusTag(status?: number) {
  switch (status) {
    case 1:
      return <Tag color="gold">REVIEWING</Tag>;
    case 2:
      return <Tag color="green">APPROVED</Tag>;
    case 3:
      return <Tag color="red">REJECTED</Tag>;
    default:
      return <Tag>UNKNOWN</Tag>;
  }
}

export default function ProductReviewPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loadingKey, setLoadingKey] = useState("");

  const query = useQuery({
    queryKey: ["admin", "product-review-tasks"],
    queryFn: () => fetchProductReviewTasks({ page: 1, pageSize: 20 })
  });

  const columns: ColumnsType<ProductReviewTask> = useMemo(
    () => [
      { title: "SPU No", dataIndex: "spuNo", key: "spuNo", width: 220 },
      { title: isZh ? "商品标题" : "Title", dataIndex: "title", key: "title", width: 300 },
      {
        title: isZh ? "审核状态" : "Review Status",
        key: "status",
        width: 150,
        render: (_, row) => reviewStatusTag(row.review?.reviewStatus)
      },
      {
        title: isZh ? "操作" : "Actions",
        key: "actions",
        width: 760,
        render: (_, row) => (
          <Space wrap>
            <Button
              size="small"
              onClick={async () => {
                const detail = await fetchProductReviewDetail(row.spuNo);
                Modal.info({
                  title: `${isZh ? "商品详情" : "Product Detail"} - ${row.spuNo}`,
                  width: 760,
                  content: <pre className="admin-inline-pre">{JSON.stringify(detail, null, 2)}</pre>
                });
              }}
            >
              {isZh ? "查看详情" : "Detail"}
            </Button>

            <AccessControl require={Permissions.ProductReviewApprove}>
              <Button
                type="primary"
                size="small"
                loading={loadingKey === `${row.spuNo}-approve`}
                onClick={async () => {
                  try {
                    setLoadingKey(`${row.spuNo}-approve`);
                    await approveProduct({ spuNo: row.spuNo, expectedVersion: 1 });
                    message.success(isZh ? "通过成功" : "Approved");
                    await query.refetch();
                  } finally {
                    setLoadingKey("");
                  }
                }}
              >
                {isZh ? "通过" : "Approve"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewReject}>
              <Button
                danger
                size="small"
                loading={loadingKey === `${row.spuNo}-reject`}
                onClick={async () => {
                  try {
                    setLoadingKey(`${row.spuNo}-reject`);
                    await rejectProduct({
                      spuNo: row.spuNo,
                      expectedVersion: 1,
                      rejectReasonCode: "CONTENT_INVALID"
                    });
                    message.success(isZh ? "驳回成功" : "Rejected");
                    await query.refetch();
                  } finally {
                    setLoadingKey("");
                  }
                }}
              >
                {isZh ? "驳回" : "Reject"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewFreeze}>
              <Button
                size="small"
                loading={loadingKey === `${row.spuNo}-freeze`}
                onClick={async () => {
                  try {
                    setLoadingKey(`${row.spuNo}-freeze`);
                    await freezeProduct({ spuNo: row.spuNo, expectedVersion: 1, reasonCode: "ADMIN_CONTROL" });
                    message.success(isZh ? "冻结成功" : "Frozen");
                  } finally {
                    setLoadingKey("");
                  }
                }}
              >
                {isZh ? "冻结" : "Freeze"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewUnfreeze}>
              <Button
                size="small"
                loading={loadingKey === `${row.spuNo}-unfreeze`}
                onClick={async () => {
                  try {
                    setLoadingKey(`${row.spuNo}-unfreeze`);
                    await unfreezeProduct({ spuNo: row.spuNo, expectedVersion: 1, reasonCode: "ADMIN_CONTROL" });
                    message.success(isZh ? "解冻成功" : "Unfrozen");
                  } finally {
                    setLoadingKey("");
                  }
                }}
              >
                {isZh ? "解冻" : "Unfreeze"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.ProductReviewForceOffShelf}>
              <Button
                size="small"
                danger
                loading={loadingKey === `${row.spuNo}-off`}
                onClick={async () => {
                  try {
                    setLoadingKey(`${row.spuNo}-off`);
                    await forceOffShelf({ spuNo: row.spuNo, expectedVersion: 1, reasonCode: "ADMIN_FORCE" });
                    message.success(isZh ? "强制下架成功" : "Forced off shelf");
                  } finally {
                    setLoadingKey("");
                  }
                }}
              >
                {isZh ? "强制下架" : "Force Off Shelf"}
              </Button>
            </AccessControl>
          </Space>
        )
      }
    ],
    [isZh, loadingKey, query]
  );

  return (
    <AdminPage
      title={isZh ? "商品审核中心" : "Product Review Center"}
      subtitle={isZh ? "统一处理商品审批、驳回、冻结与强制下架。" : "Handle approve/reject/freeze/off-shelf actions in one place."}
    >
      <div className="admin-table">
        <Table
          rowKey="spuNo"
          loading={query.isLoading}
          dataSource={query.data?.tasks ?? []}
          columns={columns}
          pagination={{
            current: query.data?.page ?? 1,
            pageSize: query.data?.pageSize ?? 20,
            total: query.data?.total ?? 0,
            showSizeChanger: false
          }}
          scroll={{ x: 1680 }}
        />
      </div>
    </AdminPage>
  );
}
