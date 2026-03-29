"use client";

import { useEffect, useState } from "react";
import { Button, Card, Form, Input, Select, Space, Typography, message } from "antd";
import { useRouter } from "next/navigation";
import { hasPermission, useI18n } from "@shopa/ui";
import { Permissions } from "@/features/admin/permissions";
import { useAuth } from "@/features/iam/useAuth";
import { useAuthStore } from "@/features/iam/store";

const { Title, Text } = Typography;

function resolveAdminHome(permissions: string[]) {
  if (hasPermission(permissions, Permissions.DashboardMallView)) {
    return "/admin/dashboard";
  }
  if (hasPermission(permissions, Permissions.MerchantReviewView)) {
    return "/admin/review/merchant";
  }
  if (hasPermission(permissions, Permissions.ProductReviewView)) {
    return "/admin/review/product";
  }
  if (hasPermission(permissions, Permissions.CSConversationView)) {
    return "/admin/cs/conversations";
  }
  return "/403";
}

export default function AdminLoginPage() {
  const router = useRouter();
  const { locale, setLocale } = useI18n();
  const { login, bootstrap } = useAuth();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let active = true;
    (async () => {
      const ok = await bootstrap();
      if (!active || !ok) {
        return;
      }
      const permissions = useAuthStore.getState().permissions;
      router.replace(resolveAdminHome(permissions));
    })();
    return () => {
      active = false;
    };
  }, [bootstrap, router]);

  const isZh = locale === "zh-CN";

  return (
    <div className="admin-login-wrap">
      <Card className="admin-login-card">
        <Space direction="vertical" size={10} style={{ width: "100%" }}>
          <Space style={{ width: "100%", justifyContent: "space-between" }}>
            <Title level={4} style={{ margin: 0 }}>
              {isZh ? "Shopa 管理后台" : "Shopa Admin Console"}
            </Title>
            <Select
              value={locale}
              style={{ width: 120 }}
              onChange={(value) => setLocale(value === "zh-CN" ? "zh-CN" : "en-US")}
              options={[
                { label: "中文", value: "zh-CN" },
                { label: "English", value: "en-US" }
              ]}
            />
          </Space>

          <Text type="secondary">{isZh ? "使用管理员账号登录并按权限访问模块。" : "Sign in with an admin account to access authorized modules."}</Text>

          <Form
            layout="vertical"
            onFinish={async (values: { identifier: string; password: string }) => {
              setLoading(true);
              try {
                await login(values.identifier.trim(), values.password);
                const permissions = useAuthStore.getState().permissions;
                router.replace(resolveAdminHome(permissions));
              } catch (err) {
                const text = String((err as Error)?.message || "");
                if (text.includes("MFA_REQUIRED")) {
                  message.error(
                    isZh
                      ? "账号触发了二次验证（MFA），当前后台登录页暂未接入该流程。请先使用 127.0.0.1 访问并关闭本地 mfaOnIpChange。"
                      : "MFA challenge is required. Current admin login page does not support MFA flow yet."
                  );
                } else {
                  message.error(isZh ? "登录失败，请检查账号或密码。" : "Login failed. Please check your credentials.");
                }
              } finally {
                setLoading(false);
              }
            }}
          >
            <Form.Item
              name="identifier"
              label={isZh ? "账号" : "Account"}
              rules={[{ required: true, message: isZh ? "请输入账号" : "Please input account" }]}
            >
              <Input autoComplete="username" placeholder={isZh ? "手机号 / 邮箱 / 用户名" : "Phone / Email / Username"} />
            </Form.Item>
            <Form.Item
              name="password"
              label={isZh ? "密码" : "Password"}
              rules={[{ required: true, message: isZh ? "请输入密码" : "Please input password" }]}
            >
              <Input.Password autoComplete="current-password" placeholder={isZh ? "请输入密码" : "Password"} />
            </Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              {isZh ? "登录管理后台" : "Sign In"}
            </Button>
          </Form>

          <Text type="secondary">
            {isZh ? "买家商城入口：" : "Buyer mall:"}{" "}
            <a href="http://127.0.0.1:3000" target="_blank" rel="noreferrer">
              http://127.0.0.1:3000
            </a>
          </Text>
        </Space>
      </Card>
    </div>
  );
}
