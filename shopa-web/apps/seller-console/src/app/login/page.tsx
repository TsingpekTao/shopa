"use client";

import { useState } from "react";
import { Button, Card, Form, Input, notification, Space, Typography } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useAuth } from "@/features/iam/useAuth";

const { Title, Text } = Typography;

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [liveMessage, setLiveMessage] = useState<string | null>(null);
  const { login } = useAuth();

  const onFinish = async (values: { identifier: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.identifier.trim(), values.password);
      notification.success({
        message: "Login successful",
        description: "Redirecting you to the workbench.",
        placement: "topRight"
      });
      window.location.href = "/seller/workbench";
    } catch (error) {
      const description = "Login failed. Please verify your credentials.";
      notification.error({ message: "Authentication error", description, placement: "topRight" });
      setLiveMessage(description);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-shell">
      <Card className="login-card">
        <Space direction="vertical" size="middle" style={{ width: "100%" }}>
          <Title level={3} style={{ margin: 0 }}>
            Seller Sign In
          </Title>
          <Text type="secondary">Manage your shop, publish new goods, and monitor performance.</Text>
        </Space>
        <Form
          layout="vertical"
          onFinish={onFinish}
          style={{ marginTop: 32 }}
          autoComplete="on"
          requiredMark="optional"
        >
          <Form.Item
            name="identifier"
            label="Account"
            rules={[{ required: true, message: "Please enter your phone, email or seller ID" }]}
            normalize={(value) => value?.trim()}
          >
            <Input
              autoComplete="username"
              placeholder="phone · email · seller ID"
              prefix={<UserOutlined />}
              disabled={loading}
            />
          </Form.Item>
          <Form.Item
            name="password"
            label="Password"
            rules={[{ required: true, message: "Password is required" }]}
          >
            <Input.Password
              autoComplete="current-password"
              prefix={<LockOutlined />}
              placeholder="Enter password"
              disabled={loading}
            />
          </Form.Item>
          <div aria-live="polite" role="status" className="sr-only">
            {liveMessage}
          </div>
          <Button type="primary" htmlType="submit" loading={loading} block className="login-submit">
            Sign In
          </Button>
        </Form>
      </Card>
    </div>
  );
}
