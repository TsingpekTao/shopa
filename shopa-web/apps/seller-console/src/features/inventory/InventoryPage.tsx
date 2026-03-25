"use client";

import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Modal,
  Table,
  Tag,
  Typography,
  notification
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { adjustInventory, fetchInventoryList } from "./api";
import { FailedItem, InventoryAdjustPayload, InventoryRecord } from "./types";

const { Paragraph } = Typography;

export function InventoryPage() {
  const [selected, setSelected] = useState<InventoryRecord | null>(null);
  const [failedItems, setFailedItems] = useState<FailedItem[]>([]);
  const [form] = Form.useForm<InventoryAdjustPayload>();

  const { data = [], isLoading, refetch } = useQuery({
    queryKey: ["inventory-list"],
    queryFn: fetchInventoryList
  });

  const mutation = useMutation({
    mutationFn: adjustInventory,
    onSuccess: (result) => {
      setFailedItems(result.failedItems);
      notification.success({ message: "Inventory adjusted" });
      form.resetFields();
      setSelected(null);
      void refetch();
    },
    onError: () => {
      notification.error({ message: "Adjustment failed" });
    }
  });

  const onConfirmAdjust = () => {
    if (!selected) {
      return;
    }

    const payload: InventoryAdjustPayload = {
      skuNo: selected.skuNo,
      delta: form.getFieldValue("delta"),
      reason: form.getFieldValue("reason")
    };

    mutation.mutate(payload);
  };

  const columns: ColumnsType<InventoryRecord> = [
    { title: "SKU", dataIndex: "skuNo", key: "skuNo" },
    { title: "Name", dataIndex: "skuName", key: "skuName" },
    { title: "Available", dataIndex: "availableQty", key: "availableQty", render: (v: number) => <Tag color="green">{v}</Tag> },
    { title: "Locked", dataIndex: "lockedQty", key: "lockedQty" },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (value: string) => <Tag>{value}</Tag>
    },
    {
      title: "Action",
      key: "action",
      render: (_: unknown, record: InventoryRecord) => (
        <Button type="link" onClick={() => setSelected(record)}>
          Adjust
        </Button>
      )
    }
  ];

  return (
    <Card title="Inventory Control">
      <Paragraph>Adjust inventory and inspect failed items from batch operations.</Paragraph>
      <Table<InventoryRecord>
        rowKey="skuNo"
        loading={isLoading || mutation.isLoading}
        dataSource={data}
        pagination={false}
        columns={columns}
      />

      {failedItems.length > 0 && (
        <Alert
          style={{ marginTop: 16 }}
          type="warning"
          showIcon
          message="Partial failures"
          description={failedItems.map((item) => `${item.skuNo}: ${item.reason}`).join(" | ")}
        />
      )}

      <Modal
        open={Boolean(selected)}
        title={`Adjust ${selected?.skuName ?? ""}`}
        onCancel={() => setSelected(null)}
        onOk={onConfirmAdjust}
        confirmLoading={mutation.isLoading}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="Delta" name="delta" rules={[{ required: true, message: "Please input delta" }]}>
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item label="Reason" name="reason">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}
