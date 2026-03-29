"use client";

import Link from "next/link";
import { useState } from "react";
import { Alert, Button, Card, Form, Input, Typography, notification } from "antd";
import { LockOutlined, UserOutlined } from "@ant-design/icons";
import { useAuth } from "@/features/iam/useAuth";
import { fetchSellerWorkbench, listMyApplications } from "@/features/seller-shop/api";
import { SellerApplicationListItem } from "@/features/seller-shop/types";
import { useI18n } from "@shopa/ui";

const { Title, Text } = Typography;

type ApiErrorLike = {
  code?: number;
  message?: string;
};

function resolveApplicationRouteByStatus(item: SellerApplicationListItem | undefined) {
  if (!item?.applicationNo) {
    return "/onboarding/start";
  }
  const code = (item.applicationStatusCode || "").toUpperCase();
  if (code === "APPROVED") {
    return "/seller/workbench";
  }
  if (code === "SUBMITTED" || code === "REVIEWING") {
    return `/onboarding/reviewing?applicationNo=${item.applicationNo}`;
  }
  if (code === "REJECTED" || code === "SYSTEM_REJECTED") {
    return `/onboarding/rejected?applicationNo=${item.applicationNo}`;
  }
  return `/onboarding/start?applicationNo=${item.applicationNo}`;
}

function pickPreferredApplication(applications: SellerApplicationListItem[]) {
  if (!applications.length) {
    return undefined;
  }
  const approved = applications.find((item) => (item.applicationStatusCode || "").toUpperCase() === "APPROVED");
  if (approved) {
    return approved;
  }
  const reviewing = applications.find((item) => {
    const code = (item.applicationStatusCode || "").toUpperCase();
    return code === "SUBMITTED" || code === "REVIEWING";
  });
  if (reviewing) {
    return reviewing;
  }
  const rejected = applications.find((item) => {
    const code = (item.applicationStatusCode || "").toUpperCase();
    return code === "REJECTED" || code === "SYSTEM_REJECTED";
  });
  if (rejected) {
    return rejected;
  }
  return applications[0];
}

async function resolveSellerLanding() {
  const workbench = await fetchSellerWorkbench();
  if (workbench.shopsTotal > 0) {
    return "/seller/workbench";
  }

  const fromWorkbench = pickPreferredApplication(workbench.latestApplications ?? []);
  if (fromWorkbench && (fromWorkbench.applicationStatusCode || "").toUpperCase() === "APPROVED") {
    return "/seller/workbench";
  }

  try {
    const mine = await listMyApplications({ page: 1, pageSize: 100 });
    const fromList = pickPreferredApplication(mine.applications ?? []);
    if (fromList) {
      return resolveApplicationRouteByStatus(fromList);
    }
  } catch (error) {
    console.warn("fallback list applications failed", error);
  }

  if (fromWorkbench) {
    return resolveApplicationRouteByStatus(fromWorkbench);
  }

  return "/onboarding/start";
}

function parseRetryAfterSeconds(message?: string): number {
  const matched = (message || "").match(/retry after\s+(\d+)\s+seconds?/i);
  if (!matched) {
    return 30 * 60;
  }
  const value = Number(matched[1]);
  return Number.isFinite(value) && value > 0 ? value : 30 * 60;
}

function formatWaitHint(seconds: number, isZh: boolean): string {
  const mins = Math.max(1, Math.ceil(seconds / 60));
  return isZh ? `请约 ${mins} 分钟后再登录。` : `Please retry in about ${mins} minutes.`;
}

export default function LoginPage() {
  const { login } = useAuth();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loading, setLoading] = useState(false);

  const handlePasswordLogin = async (values: { identifier: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.identifier.trim(), values.password);
      const target = await resolveSellerLanding();
      notification.success({ message: isZh ? "登录成功" : "Login success" });
      window.location.href = target;
    } catch (error) {
      const apiErr = error as ApiErrorLike;
      if (apiErr?.code === 101004) {
        const retryAfterSeconds = parseRetryAfterSeconds(apiErr?.message);
        notification.error({
          message: isZh ? "当前账号登录次数过多，请稍后再试" : "Too many login attempts, please try later",
          description: formatWaitHint(retryAfterSeconds, isZh)
        });
      } else {
        notification.error({
          message: isZh ? "登录失败，请检查账号或密码" : "Login failed, please check account or password"
        });
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="tb-login-wrap">
      <div className="tb-login-banner">
        <h2>{isZh ? "SHOPA 商家后台" : "SHOPA Seller Admin"}</h2>
        <p>{isZh ? "统一管理入驻申请、商品、库存与店铺运营。" : "Manage your applications, products, inventory, and shop operations."}</p>
      </div>

      <Card className="tb-login-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          {isZh ? "商家登录" : "Seller Login"}
        </Title>
        <Text type="secondary">{isZh ? "新用户请先注册，再提交入驻资料。" : "If you are new, register first and submit onboarding materials."}</Text>

        <Form layout="vertical" onFinish={handlePasswordLogin} autoComplete="on" style={{ marginTop: 16 }}>
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
            {isZh ? "登录商家后台" : "Login Seller Admin"}
          </Button>
        </Form>

        <Alert
          style={{ marginTop: 14 }}
          type="info"
          showIcon
          message={isZh ? "入驻提醒" : "Onboarding Gate"}
          description={isZh ? "如果还没有通过审核的店铺，登录后会自动引导到入驻流程。" : "No approved shop yet? You will be redirected to onboarding flow after login."}
        />

        <div className="tb-register-hint">
          <Text type="secondary">{isZh ? "还没有账号？" : "No account yet?"}</Text>
          <Link href="/register">{isZh ? "去注册" : "Register"}</Link>
        </div>
      </Card>
    </section>
  );
}

