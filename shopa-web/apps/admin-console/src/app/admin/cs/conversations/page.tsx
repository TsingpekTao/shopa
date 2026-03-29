"use client";

import { Alert, Table, Tag } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { fetchConversations } from "@/features/admin/api";
import { ConversationItem } from "@/features/admin/types";
import { AdminPage } from "@/components/admin-page";

export default function ConversationsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const query = useQuery({
    queryKey: ["admin", "cs-conversations"],
    queryFn: fetchConversations
  });

  const columns: ColumnsType<ConversationItem> = [
    { title: "Conversation No", dataIndex: "conversationNo", key: "conversationNo", width: 220 },
    { title: isZh ? "会话标题" : "Title", dataIndex: "title", key: "title", width: 280 },
    { title: isZh ? "最后消息时间" : "Last Message At", dataIndex: "lastMessageAt", key: "lastMessageAt", width: 220 },
    {
      title: isZh ? "状态" : "Status",
      dataIndex: "statusCode",
      key: "statusCode",
      render: (value) => <Tag>{value || "UNKNOWN"}</Tag>
    }
  ];

  return (
    <AdminPage
      title={isZh ? "客服会话工作台" : "Customer Service Workspace"}
      subtitle={isZh ? "查看用户会话列表与处理状态（V1 占位数据）。" : "View conversation list and statuses (V1 placeholder)."}
    >
      <div className="admin-table">
        <Table rowKey="conversationNo" loading={query.isLoading} dataSource={query.data?.items ?? []} columns={columns} pagination={false} />
      </div>

      {query.data?.partial ? (
        <Alert
          style={{ marginTop: 14 }}
          type="warning"
          showIcon
          message={isZh ? "当前模块返回占位数据" : "Placeholder data returned"}
          description={(query.data.degradedFields ?? []).join(", ")}
        />
      ) : null}
    </AdminPage>
  );
}
