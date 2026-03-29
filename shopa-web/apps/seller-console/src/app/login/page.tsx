"use client";

import Link from "next/link";
import { useState } from "react";
import { Alert, Button, Card, Form, Input, Tabs, Typography, notification } from "antd";
import { AlipayCircleOutlined, LockOutlined, QrcodeOutlined, UserOutlined } from "@ant-design/icons";
import { useAuth } from "@/features/iam/useAuth";

const { Title, Text } = Typography;

export default function LoginPage() {
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);
  const [liveMessage, setLiveMessage] = useState("");

  const handlePasswordLogin = async (values: { identifier: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.identifier.trim(), values.password);
      notification.success({ message: "登录成功，正在跳转卖家中心" });
      setLiveMessage("登录成功，正在跳转。");
      window.location.href = "/seller/workbench";
    } catch (_err) {
      notification.error({ message: "登录失败，请检查账号或密码" });
      setLiveMessage("登录失败，请检查账号或密码。");
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="tb-login-wrap">
      <div className="tb-login-banner">
        <h2>SHOPA 商家中心</h2>
        <p>品牌运营、发品管理、订单履约与营销投放，一站式完成。</p>
      </div>

      <Card className="tb-login-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          欢迎登录
        </Title>
        <Text type="secondary">登录后可进入卖家中心。</Text>

        <Tabs
          className="tb-login-tabs"
          defaultActiveKey="password"
          items={[
            {
              key: "password",
              label: "账号密码登录",
              children: (
                <Form layout="vertical" onFinish={handlePasswordLogin} autoComplete="on">
                  <Form.Item
                    name="identifier"
                    label="账号"
                    rules={[{ required: true, message: "请输入手机号/邮箱/卖家ID" }]}
                  >
                    <Input
                      autoComplete="username"
                      placeholder="手机号 / 邮箱 / 卖家ID"
                      prefix={<UserOutlined />}
                    />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    label="密码"
                    rules={[{ required: true, message: "请输入密码" }]}
                  >
                    <Input.Password
                      autoComplete="current-password"
                      placeholder="请输入密码"
                      prefix={<LockOutlined />}
                    />
                  </Form.Item>
                  <div className="tb-sr-only" aria-live="polite" role="status">
                    {liveMessage}
                  </div>
                  <div className="tb-login-options">
                    <Link href="/forgot-password">忘记密码？</Link>
                  </div>
                  <Button type="primary" htmlType="submit" block loading={loading}>
                    登录
                  </Button>
                </Form>
              )
            },
            {
              key: "qrcode",
              label: "扫码登录",
              children: (
                <div className="tb-qr-login">
                  <div className="tb-qr-box" aria-hidden="true">
                    <QrcodeOutlined />
                  </div>
                  <p>请使用手机淘宝或支付宝扫码登录</p>
                  <Text type="secondary">扫码后请在手机端确认登录</Text>
                </div>
              )
            }
          ]}
        />

        <div className="tb-third-party">
          <Text type="secondary">第三方账号登录</Text>
          <div className="tb-third-party-actions">
            <Button icon={<AlipayCircleOutlined />} className="tb-third-party-btn">
              支付宝登录
            </Button>
          </div>
        </div>

        <Alert
          style={{ marginTop: 14 }}
          type="info"
          showIcon
          message="提示"
          description="第三方登录能力接入中，当前可先使用账号密码登录。"
        />

        <div className="tb-register-hint">
          <Text type="secondary">还没有账号？</Text>
          <Link href="/register">去注册</Link>
        </div>
      </Card>
    </section>
  );
}
