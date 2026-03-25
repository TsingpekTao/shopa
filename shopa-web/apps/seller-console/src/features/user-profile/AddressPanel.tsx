"use client";

import { useEffect, useState } from "react";
import {
  Card,
  Form,
  Input,
  Button,
  List,
  Divider,
  Tag,
  Space,
  Skeleton,
  Alert,
  message
} from "antd";
import { Address } from "./types";
import { fetchAddresses, saveAddress, setDefaultAddress } from "./api";

export function AddressPanel() {
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [form] = Form.useForm<Partial<Address>>();

  useEffect(() => {
    void load();
  }, []);

  async function load() {
    setLoading(true);
    setError(null);
    try {
      const list = await fetchAddresses();
      setAddresses(list);
    } catch {
      setError("地址加载失败，请稍后再试");
    } finally {
      setLoading(false);
    }
  }

  const onFinish = async (values: Partial<Address>) => {
    setSaving(true);
    try {
      await saveAddress(values);
      message.success("地址已保存");
      form.resetFields();
      await load();
    } catch (err) {
      message.error("保存地址失败");
    } finally {
      setSaving(false);
    }
  };

  const onSetDefault = async (id: string) => {
    try {
      await setDefaultAddress(id);
      message.success("默认地址已更新");
      await load();
    } catch {
      message.error("更新默认地址失败");
    }
  };

  return (
    <Card title="收货地址" style={{ minHeight: 420 }} bordered>
      {error && (
        <Alert
          type="warning"
          showIcon
          message={error}
          style={{ marginBottom: 16 }}
          aria-live="polite"
        />
      )}
      <Skeleton active loading={loading} paragraph={{ rows: 6 }} title={{ width: "40%" }}>
        <List
          dataSource={addresses}
          locale={{ emptyText: "目前还没有收货地址" }}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button
                  key="default"
                  type={item.isDefault ? "primary" : "link"}
                  onClick={() => onSetDefault(item.id)}
                  aria-label={`设置 ${item.label} 为默认地址`}
                >
                  {item.isDefault ? "默认地址" : "设为默认"}
                </Button>
              ]}
            >
              <List.Item.Meta
                title={
                  <Space>
                    <span>{item.label}</span>
                    {item.isDefault && <Tag color="success">默认</Tag>}
                  </Space>
                }
                description={
                  <>
                    <div>{item.fullAddress}</div>
                    <div style={{ color: "#6e6e6e" }}>
                      {item.contact} · {item.phone}
                    </div>
                  </>
                }
              />
            </List.Item>
          )}
        />
      </Skeleton>
      <Divider>添加 / 编辑地址</Divider>
      <Form form={form} layout="vertical" onFinish={onFinish} autoComplete="off">
        <Form.Item label="收件标签" name="label" rules={[{ required: true }]}>
          <Input placeholder="家 / 仓库 / 门店" />
        </Form.Item>
        <Form.Item label="联系人" name="contact" rules={[{ required: true }]}>
          <Input placeholder="张三" />
        </Form.Item>
        <Form.Item label="联系电话" name="phone" rules={[{ required: true }]}>
          <Input placeholder="13800138000" />
        </Form.Item>
        <Form.Item label="详细地址" name="fullAddress" rules={[{ required: true }]}>
          <Input.TextArea rows={2} placeholder="填写完整地址信息" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} aria-label="保存地址">
            保存地址
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
