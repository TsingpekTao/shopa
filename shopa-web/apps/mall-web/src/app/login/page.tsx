"use client";

import Link from "next/link";
import { useState } from "react";
import { Alert, Button, Card, Form, Input, Tabs, Typography, notification } from "antd";
import { AlipayCircleOutlined, LockOutlined, QrcodeOutlined, UserOutlined } from "@ant-design/icons";
import { useAuth } from "@/features/iam/useAuth";
import { useI18n } from "@shopa/ui";

const { Title, Text } = Typography;

export default function LoginPage() {
  const { login } = useAuth();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(false);

  const handlePasswordLogin = async (values: { identifier: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.identifier.trim(), values.password);
      const searchParams = new URLSearchParams(window.location.search);
      const redirect = searchParams.get("redirect");
      notification.success({ message: isZh ? "登录成功" : "Login success" });
      window.location.href = redirect && redirect.startsWith("/") ? redirect : "/";
    } catch (_err) {
      notification.error({ message: isZh ? "登录失败，请检查账号或密码" : "Login failed, please check account or password" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="tb-login-wrap">
      <div className="tb-login-banner">
        <h2>{isZh ? "SHOPA 商城登录" : "SHOPA Mall Login"}</h2>
        <p>{isZh ? "作为买家登录后，可浏览商品、下单并跟踪物流。" : "Login as a customer to browse products, place orders, and track deliveries."}</p>
      </div>

      <Card className="tb-login-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          {isZh ? "欢迎回来" : "Welcome back"}
        </Title>
        <Text type="secondary">{isZh ? "此登录入口仅用于买家端。" : "This login is for buyer side only."}</Text>

        <Tabs
          className="tb-login-tabs"
          defaultActiveKey="password"
          items={[
            {
              key: "password",
              label: isZh ? "账号登录" : "Account Login",
              children: (
                <Form layout="vertical" onFinish={handlePasswordLogin} autoComplete="on">
                  <Form.Item
                    name="identifier"
                    label={isZh ? "账号" : "Account"}
                    rules={[{ required: true, message: isZh ? "请输入手机号/邮箱/账号" : "Please input phone/email/account" }]}
                  >
                    <Input autoComplete="username" placeholder={isZh ? "手机号 / 邮箱 / 账号" : "Phone / Email / Account"} prefix={<UserOutlined />} />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    label={isZh ? "密码" : "Password"}
                    rules={[{ required: true, message: isZh ? "请输入密码" : "Please input password" }]}
                  >
                    <Input.Password autoComplete="current-password" placeholder={isZh ? "密码" : "Password"} prefix={<LockOutlined />} />
                  </Form.Item>
                  <div className="tb-login-options">
                    <Link href="/forgot-password">{isZh ? "忘记密码？" : "Forgot password?"}</Link>
                  </div>
                  <Button type="primary" htmlType="submit" block loading={loading}>
                    {isZh ? "登录" : "Login"}
                  </Button>
                </Form>
              )
            },
            {
              key: "qrcode",
              label: isZh ? "扫码登录" : "QR Login",
              children: (
                <div className="tb-qr-login">
                  <div className="tb-qr-box" aria-hidden="true">
                    <QrcodeOutlined />
                  </div>
                  <p>{isZh ? "使用手机 App 扫码快速登录。" : "Scan by mobile app to login quickly."}</p>
                  <Text type="secondary">{isZh ? "当前本地环境为演示模式。" : "QR login is a demo mode in current local environment."}</Text>
                </div>
              )
            }
          ]}
        />

        <div className="tb-third-party">
          <Text type="secondary">{isZh ? "第三方账号登录" : "Third-party login"}</Text>
          <div className="tb-third-party-actions">
            <Button icon={<AlipayCircleOutlined />} className="tb-third-party-btn">
              {isZh ? "支付宝登录" : "Alipay Login"}
            </Button>
          </div>
        </div>

        <Alert
          style={{ marginTop: 14 }}
          type="info"
          showIcon
          message={isZh ? "提示" : "Tip"}
          description={
            isZh
              ? "卖家后台已独立，请访问 http://127.0.0.1:3100/ 进行商家操作。"
              : "Seller backend is separated now. Open http://127.0.0.1:3100/ for seller operations."
          }
        />

        <div className="tb-register-hint">
          <Text type="secondary">{isZh ? "还没有账号？" : "No account yet?"}</Text>
          <Link href="/register">{isZh ? "去注册" : "Register"}</Link>
        </div>
      </Card>
    </section>
  );
}

