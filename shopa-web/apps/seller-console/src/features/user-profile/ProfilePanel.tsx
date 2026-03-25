"use client";

import { useEffect, useState } from "react";
import { Avatar, Button, Card, Form, Input, Space, Skeleton, Alert, Typography } from "antd";
import { fetchProfile, saveProfile } from "./api";
import { Profile } from "./types";

const { Paragraph } = Typography;

export function ProfilePanel() {
  const [form] = Form.useForm<Profile>();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        const profile = await fetchProfile();
        form.setFieldsValue(profile);
      } catch {
        setError("Unable to load profile, please retry.");
      } finally {
        setLoading(false);
      }
    };
    void load();
  }, [form]);

  const onFinish = async (values: Profile) => {
    setSaving(true);
    try {
      await saveProfile(values);
      setError(null);
    } catch {
      setError("Failed to save profile. Try again in a moment.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card
      title={
        <Space align="center">
          <Avatar style={{ backgroundColor: "#f2545b" }} size={46}>
            商
          </Avatar>
          <div>
            <Paragraph style={{ margin: 0, fontWeight: 600 }}>商家资料</Paragraph>
            <Paragraph style={{ margin: 0, color: "#7a7a7a" }} type="secondary">
              更新后后台自动同步给客服与审核系统。
            </Paragraph>
          </div>
        </Space>
      }
      bordered
      style={{ minHeight: 380 }}
    >
      {error && (
        <Alert
          type="warning"
          message={error}
          showIcon
          style={{ marginBottom: 16 }}
          aria-live="polite"
        />
      )}
      <Skeleton active loading={loading} paragraph={{ rows: 5 }} title={false}>
        <Form form={form} layout="vertical" onFinish={onFinish} autoComplete="off">
          <Form.Item
            label="显示名称"
            name="displayName"
            rules={[{ required: true, message: "请输入显示名称" }]}
            tooltip="会用于商品详情与店铺名展示"
          >
            <Input placeholder="请输入店铺对外昵称" />
          </Form.Item>
          <Form.Item
            label="邮箱"
            name="email"
            rules={[{ type: "email", message: "请输入合法邮箱" }]}
          >
            <Input placeholder="contact@example.com" />
          </Form.Item>
          <Form.Item
            label="手机号"
            name="phone"
            rules={[{ required: true, message: "请输入联系手机号" }]}
          >
            <Input placeholder="13800138000" />
          </Form.Item>
          <Form.Item label="店铺介绍" name="bio">
            <Input.TextArea
              rows={3}
              placeholder="简单描述营业范围或主营产品，利于审核与推荐"
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving} aria-label="保存资料">
              保存资料
            </Button>
          </Form.Item>
        </Form>
      </Skeleton>
    </Card>
  );
}
