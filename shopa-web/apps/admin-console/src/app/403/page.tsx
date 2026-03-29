"use client";

import Link from "next/link";
import { Button, Result, Space } from "antd";
import { useI18n } from "@shopa/ui";

export default function ForbiddenPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <div className="admin-login-wrap">
      <Result
        status="403"
        title="403"
        subTitle={isZh ? "你没有访问此页面的权限。" : "You do not have permission to access this page."}
        extra={
          <Space>
            <Link href="/admin/dashboard">
              <Button type="primary">{isZh ? "返回工作台" : "Back to dashboard"}</Button>
            </Link>
            <Link href="/login">
              <Button>{isZh ? "返回登录" : "Back to login"}</Button>
            </Link>
          </Space>
        }
      />
    </div>
  );
}
