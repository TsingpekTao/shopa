"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Form,
  Input,
  Modal,
  Segmented,
  Space,
  Steps,
  Table,
  Tag,
  Typography,
  notification
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { EyeOutlined, PlusOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { createProductDraft, listSellerProducts } from "@/features/catalog/api";
import { SellerProductSpu, SpuStatusCode } from "@/features/catalog/types";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const { Title, Text } = Typography;

type DraftFilter = "ALL" | "DRAFT" | "READY" | "REVIEWING";

type DraftRow = {
  id: string;
  name: string;
  category: string;
  status: DraftFilter | "REJECTED";
  updatedAt: string;
  mediaCount: number;
};

type CreateDraftForm = {
  title: string;
  subTitle?: string;
};

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

function mapDraftStatus(status: SpuStatusCode): DraftRow["status"] {
  if (status === "SPU_STATUS_DRAFT") {
    return "DRAFT";
  }
  if (status === "SPU_STATUS_REVIEWING") {
    return "REVIEWING";
  }
  if (status === "SPU_STATUS_REJECTED") {
    return "REJECTED";
  }
  return "READY";
}

function statusColor(status: DraftRow["status"]) {
  if (status === "READY") return "processing";
  if (status === "REVIEWING") return "success";
  if (status === "REJECTED") return "error";
  return "default";
}

function statusesForFilter(filter: DraftFilter): SpuStatusCode[] | undefined {
  if (filter === "DRAFT") {
    return ["SPU_STATUS_DRAFT"];
  }
  if (filter === "REVIEWING") {
    return ["SPU_STATUS_REVIEWING"];
  }
  if (filter === "READY") {
    return ["SPU_STATUS_APPROVED", "SPU_STATUS_ON_SHELF", "SPU_STATUS_OFF_SHELF"];
  }
  return undefined;
}

export default function PublishPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [filter, setFilter] = useState<DraftFilter>("ALL");
  const [createOpen, setCreateOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form] = Form.useForm<CreateDraftForm>();

  const workbenchQuery = useQuery({
    queryKey: ["seller-workbench-for-publish"],
    queryFn: async () => fetchSellerWorkbench()
  });

  const draftsQuery = useQuery({
    queryKey: ["seller-publish-drafts", filter],
    queryFn: async () =>
      listSellerProducts({
        page: 1,
        pageSize: 50,
        statuses: statusesForFilter(filter)
      })
  });

  const shopNo = workbenchQuery.data?.shops?.[0]?.shopNo ? String(workbenchQuery.data.shops[0].shopNo) : "";

  const rows = useMemo<DraftRow[]>(() => {
    return (draftsQuery.data?.products ?? []).map((item: SellerProductSpu) => ({
      id: item.spuNo,
      name: item.title || "-",
      category: item.categoryId > 0 ? String(item.categoryId) : "-",
      status: mapDraftStatus(item.spuStatus),
      updatedAt: formatDateTime(item.updatedAt),
      mediaCount: item.mainImageAssetIds.length + item.detailImageAssetIds.length
    }));
  }, [draftsQuery.data?.products]);

  const submitCreate = async () => {
    try {
      const values = await form.validateFields();
      if (!shopNo) {
        notification.warning({
          message: isZh ? "暂无可用店铺" : "No shop available",
          description: isZh ? "请先完成店铺入驻审核并创建店铺。" : "Please finish onboarding and create a shop first."
        });
        return;
      }

      setCreating(true);
      const created = await createProductDraft({
        shopNo,
        title: values.title,
        subTitle: values.subTitle ?? "",
        categoryId: 0
      });

      notification.success({
        message: isZh ? "草稿创建成功" : "Draft created",
        description: `${isZh ? "SPU 编号" : "SPU"}: ${created.spuNo || "-"}`
      });
      setCreateOpen(false);
      form.resetFields();
      await draftsQuery.refetch();
    } catch (error) {
      if (error instanceof Error && error.message) {
        notification.error({ message: isZh ? "创建失败" : "Create failed", description: error.message });
      }
    } finally {
      setCreating(false);
    }
  };

  const columns = useMemo<ColumnsType<DraftRow>>(
    () => [
      { title: isZh ? "草稿编号" : "Draft ID", dataIndex: "id", key: "id", width: 180 },
      { title: isZh ? "商品名称" : "Product", dataIndex: "name", key: "name" },
      { title: isZh ? "类目ID" : "Category ID", dataIndex: "category", key: "category", width: 110 },
      {
        title: isZh ? "状态" : "Status",
        dataIndex: "status",
        key: "status",
        width: 130,
        render: (value: DraftRow["status"]) => <Tag color={statusColor(value)}>{value}</Tag>
      },
      { title: isZh ? "素材数" : "Media", dataIndex: "mediaCount", key: "mediaCount", width: 100 },
      { title: isZh ? "更新时间" : "Updated", dataIndex: "updatedAt", key: "updatedAt", width: 180 },
      {
        title: isZh ? "操作" : "Action",
        key: "action",
        width: 140,
        render: (_, row) => (
          <Link href="/seller/products">
            <Button size="small" icon={<EyeOutlined />}>
              {isZh ? `查看 ${row.id}` : "View"}
            </Button>
          </Link>
        )
      }
    ],
    [isZh]
  );

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "发布商品" : "Publish Product"}</Title>
          <Text type="secondary">{isZh ? "按流程完成素材、信息、库存并提交审核。" : "Complete media, details and stock, then submit for review."}</Text>
        </div>
        <Space wrap>
          <Link href="/seller/media/library">
            <Button icon={<EyeOutlined />}>{isZh ? "查看素材库" : "Media Library"}</Button>
          </Link>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
            {isZh ? "新建商品草稿" : "New Draft"}
          </Button>
        </Space>
      </header>

      {!shopNo ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "当前账号暂无店铺" : "No shop found"}
          description={isZh ? "请先完成店铺入驻审核，之后才能创建商品草稿。" : "Complete onboarding before creating product drafts."}
        />
      ) : null}

      <Card>
        <Steps
          responsive
          current={1}
          items={[
            { title: isZh ? "选择类目" : "Category", description: isZh ? "匹配经营类目" : "Select category" },
            { title: isZh ? "完善信息" : "Details", description: isZh ? "标题、属性、主图" : "Title, specs, images" },
            { title: isZh ? "设置库存" : "Inventory", description: isZh ? "SKU 与库存" : "SKU and stock" },
            { title: isZh ? "提交审核" : "Review", description: isZh ? "平台审核" : "Platform review" }
          ]}
        />
      </Card>

      <Card
        title={isZh ? "商品草稿" : "Drafts"}
        extra={
          <Segmented
            value={filter}
            onChange={(val) => setFilter(val as DraftFilter)}
            options={[
              { label: isZh ? "全部" : "All", value: "ALL" },
              { label: "DRAFT", value: "DRAFT" },
              { label: "READY", value: "READY" },
              { label: "REVIEWING", value: "REVIEWING" }
            ]}
          />
        }
      >
        <Table<DraftRow>
          loading={draftsQuery.isLoading}
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={rows}
          pagination={false}
          locale={{ emptyText: draftsQuery.isError ? (isZh ? "加载失败，请刷新重试" : "Load failed, please retry") : undefined }}
        />
      </Card>

      <Modal
        title={isZh ? "新建商品草稿" : "Create Product Draft"}
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={() => void submitCreate()}
        okButtonProps={{ loading: creating }}
        okText={isZh ? "创建" : "Create"}
        cancelText={isZh ? "取消" : "Cancel"}
      >
        <Form<CreateDraftForm> form={form} layout="vertical">
          <Form.Item label={isZh ? "商品标题" : "Title"} name="title" rules={[{ required: true, message: isZh ? "请输入商品标题" : "Please input title" }]}> 
            <Input maxLength={80} />
          </Form.Item>
          <Form.Item label={isZh ? "副标题" : "Subtitle"} name="subTitle">
            <Input maxLength={120} />
          </Form.Item>
        </Form>
      </Modal>
    </section>
  );
}
