"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Button, Layout, Menu, Select, Space, Spin, Typography } from "antd";
import { PermissionProvider, hasPermission, useI18n } from "@shopa/ui";
import { Permissions } from "@/features/admin/permissions";
import { useAuth } from "@/features/iam/useAuth";

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

type NavItem = {
  key: string;
  labelZh: string;
  labelEn: string;
  permission: string;
  hintZh: string;
  hintEn: string;
};

const NAV_ITEMS: NavItem[] = [
  {
    key: "/admin/dashboard",
    labelZh: "商城看板",
    labelEn: "Dashboard",
    permission: Permissions.DashboardMallView,
    hintZh: "查看整体运营指标与异常预警",
    hintEn: "Overview of platform health and operations"
  },
  {
    key: "/admin/review/merchant",
    labelZh: "商家审核",
    labelEn: "Merchant Review",
    permission: Permissions.MerchantReviewView,
    hintZh: "处理商家入驻申请与资质审核",
    hintEn: "Process merchant onboarding applications"
  },
  {
    key: "/admin/review/product",
    labelZh: "商品审核",
    labelEn: "Product Review",
    permission: Permissions.ProductReviewView,
    hintZh: "审核商品发布、冻结与下架操作",
    hintEn: "Review product approvals and controls"
  },
  {
    key: "/admin/cs/conversations",
    labelZh: "客服工作台",
    labelEn: "Customer Service",
    permission: Permissions.CSConversationView,
    hintZh: "查看客服会话与待办处理任务",
    hintEn: "Review customer support conversations"
  }
];

function routePermission(pathname: string): string {
  if (pathname.startsWith("/admin/review/merchant")) {
    return Permissions.MerchantReviewView;
  }
  if (pathname.startsWith("/admin/review/product")) {
    return Permissions.ProductReviewView;
  }
  if (pathname.startsWith("/admin/cs/conversations")) {
    return Permissions.CSConversationView;
  }
  if (pathname.startsWith("/admin/shops/")) {
    return Permissions.DashboardShopView;
  }
  if (pathname.startsWith("/admin/dashboard")) {
    return Permissions.DashboardMallView;
  }
  return "";
}

function findActiveNav(pathname: string) {
  const sorted = [...NAV_ITEMS].sort((a, b) => b.key.length - a.key.length);
  return sorted.find((item) => pathname.startsWith(item.key));
}

export default function AdminLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { locale, setLocale } = useI18n();
  const { permissions, bootstrap, logout } = useAuth();
  const [ready, setReady] = useState(false);
  const [collapsed, setCollapsed] = useState(false);
  const isZh = locale === "zh-CN";

  useEffect(() => {
    let active = true;
    (async () => {
      const ok = await bootstrap();
      if (!active) {
        return;
      }
      if (!ok) {
        router.replace("/login");
        return;
      }
      setReady(true);
    })();
    return () => {
      active = false;
    };
  }, [bootstrap, router]);

  useEffect(() => {
    if (!ready) {
      return;
    }
    const required = routePermission(pathname);
    if (required && !hasPermission(permissions, required)) {
      router.replace("/403");
    }
  }, [pathname, permissions, ready, router]);

  const visibleNav = useMemo(() => {
    return NAV_ITEMS.filter((item) => hasPermission(permissions, item.permission)).map((item) => ({
      key: item.key,
      label: isZh ? item.labelZh : item.labelEn
    }));
  }, [isZh, permissions]);

  const selectedKey = useMemo(() => {
    const hit = findActiveNav(pathname);
    return hit?.key ? [hit.key] : [];
  }, [pathname]);

  const currentTitle = useMemo(() => {
    const hit = findActiveNav(pathname);
    if (!hit) {
      return {
        title: isZh ? "后台管理" : "Admin Console",
        hint: isZh ? "统一权限控制与审计中心" : "Unified permission and audit center"
      };
    }
    return {
      title: isZh ? hit.labelZh : hit.labelEn,
      hint: isZh ? hit.hintZh : hit.hintEn
    };
  }, [isZh, pathname]);

  if (!ready) {
    return (
      <div className="admin-login-wrap">
        <Spin size="large" />
      </div>
    );
  }

  return (
    <PermissionProvider permissions={permissions}>
      <Layout className="admin-shell">
        <Sider
          className="admin-sider"
          width={232}
          collapsible
          collapsed={collapsed}
          onCollapse={setCollapsed}
          breakpoint="lg"
          collapsedWidth={72}
          trigger={null}
        >
          <div className="admin-brand" aria-label="Shopa Admin">
            <span className="admin-brand-dot" />
            {!collapsed ? "SHOPA ADMIN" : "SA"}
          </div>
          <Menu
            className="admin-menu"
            mode="inline"
            selectedKeys={selectedKey}
            items={visibleNav}
            onClick={(e) => router.push(e.key)}
          />
        </Sider>

        <Layout>
          <Header className="admin-header">
            <div className="admin-header-left">
              <Button type="text" onClick={() => setCollapsed((v) => !v)}>
                {collapsed ? (isZh ? "展开" : "Expand") : isZh ? "收起" : "Collapse"}
              </Button>
              <div>
                <div className="admin-header-title">{currentTitle.title}</div>
                <div className="admin-header-subtitle">{currentTitle.hint}</div>
              </div>
            </div>

            <Space size={10}>
              <a href="http://127.0.0.1:3000" target="_blank" rel="noreferrer">
                <Button>{isZh ? "打开商城前台" : "Open Mall"}</Button>
              </a>
              <Select
                value={locale}
                style={{ width: 120 }}
                onChange={(value) => setLocale(value === "zh-CN" ? "zh-CN" : "en-US")}
                options={[
                  { label: "中文", value: "zh-CN" },
                  { label: "English", value: "en-US" }
                ]}
              />
              <Button
                danger
                onClick={async () => {
                  await logout();
                  router.replace("/login");
                }}
              >
                {isZh ? "退出登录" : "Logout"}
              </Button>
            </Space>
          </Header>

          <Content className="admin-content-wrap">
            <main className="admin-page">{children}</main>
          </Content>
        </Layout>
      </Layout>
    </PermissionProvider>
  );
}
