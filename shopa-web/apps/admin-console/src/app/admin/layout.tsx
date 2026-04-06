"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Button, Layout, Menu, Space, Spin, Typography } from "antd";
import { PermissionProvider, hasPermission } from "@shopa/ui";
import { Permissions } from "@/features/admin/permissions";
import { useAuth } from "@/features/iam/useAuth";

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

type NavItem = {
  key: string;
  permission: string;
  label: string;
  hint: string;
};

const NAV_ITEMS: NavItem[] = [
  {
    key: "/admin/dashboard",
    permission: Permissions.DashboardMallView,
    label: "商城看板",
    hint: "查看整体运营指标与异常预警"
  },
  {
    key: "/admin/shops",
    permission: Permissions.DashboardMallView,
    label: "店铺列表",
    hint: "查看平台已开通店铺与基础信息"
  },
  {
    key: "/admin/review/merchant",
    permission: Permissions.MerchantReviewView,
    label: "商家审核",
    hint: "处理商家入驻申请与资质审核"
  },
  {
    key: "/admin/review/product",
    permission: Permissions.ProductReviewView,
    label: "商品审核",
    hint: "审核商品发布、冻结与下架操作"
  },
  {
    key: "/admin/cs/conversations",
    label: "客服工作台",
    permission: Permissions.CSConversationView,
    hint: "查看客服会话与待办处理任务"
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
  if (pathname === "/admin/shops") {
    return Permissions.DashboardMallView;
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
  const { permissions, bootstrap, logout } = useAuth();
  const [ready, setReady] = useState(false);
  const [collapsed, setCollapsed] = useState(false);

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
      label: item.label
    }));
  }, [permissions]);

  const selectedKey = useMemo(() => {
    const hit = findActiveNav(pathname);
    return hit?.key ? [hit.key] : [];
  }, [pathname]);

  const currentTitle = useMemo(() => {
    const hit = findActiveNav(pathname);
    if (!hit) {
      return {
        title: "后台管理",
        hint: "统一权限控制与审计中心"
      };
    }
    return {
      title: hit.label,
      hint: hit.hint
    };
  }, [pathname]);

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
          <div className="admin-brand" aria-label="Shopa 管理台">
            <span className="admin-brand-dot" />
            {!collapsed ? "SHOPA 管理台" : "管理"}
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
                {collapsed ? "展开" : "收起"}
              </Button>
              <div>
                <div className="admin-header-title">{currentTitle.title}</div>
                <div className="admin-header-subtitle">{currentTitle.hint}</div>
              </div>
            </div>

            <Space size={10}>
              <a href="http://127.0.0.1:3000" target="_blank" rel="noreferrer">
                <Button>打开商城前台</Button>
              </a>
              <Button
                danger
                onClick={async () => {
                  await logout();
                  router.replace("/login");
                }}
              >
                退出登录
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
