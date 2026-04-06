"use client";

import { useMemo } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button, Card, Form, Input, InputNumber, Space, notification, Typography } from "antd";
import { MinusCircleOutlined, PlusOutlined } from "@ant-design/icons";
import { submitDraft } from "./api";
import { CatalogDraftPayload, CatalogSkuDraft } from "./types";

const { Paragraph, Text } = Typography;

type Props = {
  onSaved?: () => void;
};

const defaultDraft: CatalogDraftPayload = {
  draft: {
    title: "",
    summary: "",
    categoryId: 0,
    brandNo: "",
    detailImageAssetIds: [],
    attributeValues: {
      brand: "",
      color: ""
    }
  },
  skus: [
    {
      name: "Default SKU",
      salePrice: 0,
      marketPrice: 0,
      saleSpecs: {
        color: "",
        size: ""
      }
    }
  ]
};

export function CatalogSection({ onSaved }: Props) {
  const [form] = Form.useForm<CatalogDraftPayload>();

  const initial = useMemo(() => defaultDraft, []);

  const mutation = useMutation({
    mutationFn: submitDraft,
    onSuccess: () => {
      notification.success({ message: "Draft saved", description: "Your product draft is now stored in the gateway." });
      onSaved?.();
    },
    onError: () => {
      notification.error({ message: "Unable to save draft", description: "Please try again or contact support." });
    }
  });

  const onFinish = (values: CatalogDraftPayload) => {
    mutation.mutate(values);
  };

  return (
    <Card aria-label="Catalog builder" title="Publish product" bordered>
      <Paragraph>
        We stitch media uploads, catalog data, and inventory together. The draft state persists locally so you can safely refresh the page.
      </Paragraph>
      <Form form={form} layout="vertical" onFinish={onFinish} initialValues={initial}>
        <Form.Item
          name={["draft", "title"]}
          label="Product Title"
          rules={[{ required: true, message: "A product title is required for catalog visibility." }]}
        >
          <Input aria-label="Product title" placeholder="Product title" />
        </Form.Item>
        <Form.Item
          name={["draft", "summary"]}
          label="Summary"
          rules={[{ required: true, message: "Describe this product to help buyers." }]}
        >
          <Input.TextArea aria-label="Product summary" rows={3} placeholder="Short product description" />
        </Form.Item>
        <Space direction="vertical" size="small" style={{ width: "100%" }}>
          <Form.Item label="Key attributes">
            <Space style={{ width: "100%" }}>
              <Form.Item name={["draft", "attributeValues", "brand"]} rules={[{ required: true }]} style={{ flex: 1 }}>
                <Input placeholder="Brand" />
              </Form.Item>
              <Form.Item name={["draft", "attributeValues", "color"]} style={{ flex: 1 }}>
                <Input placeholder="Color" />
              </Form.Item>
            </Space>
          </Form.Item>
          <Form.List name="skus">
            {(fields, { add, remove }) => (
              <>
                {fields.map((field) => (
                  <Space key={field.key} align="baseline" style={{ width: "100%", gap: "1rem" }}>
                    <Form.Item
                      {...field}
                      name={[field.name, "name"]}
                      label="SKU Name"
                      rules={[{ required: true }]}
                      style={{ flex: 1 }}
                    >
                      <Input placeholder="SKU name" />
                    </Form.Item>
                    <Form.Item
                      {...field}
                      name={[field.name, "salePrice"]}
                      label="Sale Price"
                      rules={[{ required: true }]}
                      style={{ width: 160 }}
                    >
                      <InputNumber
                        aria-label="Sale price"
                        min={0}
                        precision={2}
                        style={{ width: "100%" }}
                        controls={false}
                      />
                    </Form.Item>
                    <Form.Item
                      {...field}
                      name={[field.name, "marketPrice"]}
                      label="Market Price"
                      rules={[{ required: true }]}
                      style={{ width: 160 }}
                    >
                      <InputNumber
                        min={0}
                        precision={2}
                        style={{ width: "100%" }}
                        aria-label="Market price"
                        controls={false}
                      />
                    </Form.Item>
                    <Button type="text" icon={<MinusCircleOutlined />} onClick={() => remove(field.name)} aria-label="Remove SKU" />
                  </Space>
                ))}
                <Button type="dashed" onClick={() => add(defaultDraft.skus[0])} icon={<PlusOutlined />} style={{ width: "100%" }}>
                  Add SKU
                </Button>
              </>
            )}
          </Form.List>
        </Space>
        <Button type="primary" htmlType="submit" loading={mutation.isLoading} style={{ marginTop: "1rem" }}>
          Save draft
        </Button>
      </Form>
    </Card>
  );
}

export default CatalogSection;
