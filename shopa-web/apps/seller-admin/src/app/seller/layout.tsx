"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  AppstoreOutlined,
  ContainerOutlined,
  DeploymentUnitOutlined,
  LineChartOutlined,
  PictureOutlined,
  ShopOutlined,
  TeamOutlined,
  UserOutlined
} from "@ant-design/icons";
import { Layout, Menu, Skeleton, Typography } from "antd";
import { useI18n } from "@shopa/ui";
import { useAuthStore } from "@/features/iam/store";

const { Sider, Content } = Layout;
const { Text } = Typography;

type Props = {
  children: ReactNode;
};

export default function SellerLayout({ children }: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const [ready, setReady] = useState(false);

  const menuItems = useMemo(
    () => [
      {
        key: "/seller/workbench",
        icon: <AppstoreOutlined />,
        label: <Link href="/seller/workbench">{isZh ? "工作台" : "Workbench"}</Link>
      },
      {
        key: "/seller/sales",
        icon: <LineChartOutlined />,
        label: <Link href="/seller/sales">{isZh ? "销售分析" : "Sales Analytics"}</Link>
      },
      {
        key: "/seller/publish",
        icon: <DeploymentUnitOutlined />,
        label: <Link href="/seller/publish">{isZh ? "发布商品" : "Publish Product"}</Link>
      },
      {
        key: "/seller/products",
        icon: <ContainerOutlined />,
        label: <Link href="/seller/products">{isZh ? "商品管理" : "Products"}</Link>
      },
      {
        key: "/seller/media/library",
        icon: <PictureOutlined />,
        label: <Link href="/seller/media/library">{isZh ? "素材库" : "Media Library"}</Link>
      },
      {
        key: "/seller/inventory",
        icon: <ShopOutlined />,
        label: <Link href="/seller/inventory">{isZh ? "库存管理" : "Inventory"}</Link>
      },
      {
        key: "/seller/profile",
        icon: <UserOutlined />,
        label: <Link href="/seller/profile">{isZh ? "店主资料" : "Profile"}</Link>
      },
      {
        key: "/seller/me",
        icon: <TeamOutlined />,
        label: <Link href="/seller/me">{isZh ? "我的" : "My Account"}</Link>
      }
    ],
    [isZh]
  );

  useEffect(() => {
    const accessToken = tokenPair?.accessToken ?? localStorage.getItem("shopa_seller_access_token");
    if (!accessToken) {
      router.replace("/login");
      return;
    }
    setReady(true);
  }, [router, tokenPair?.accessToken]);

  const selectedKey = useMemo(() => {
    const found = menuItems.find((item) => pathname.startsWith(item.key));
    return found ? [found.key] : ["/seller/workbench"];
  }, [menuItems, pathname]);

  if (!ready) {
    return <Skeleton active />;
  }

  return (
    <Layout className="seller-workspace-layout">
      <Sider width={252} className="seller-workspace-sider" breakpoint="lg" collapsedWidth="0">
        <div className="seller-workspace-sider-head">
          <Text strong>{isZh ? "商家工作区" : "Seller Workspace"}</Text>
          <Text type="secondary">{isZh ? "经营导航" : "Workspace Navigation"}</Text>
        </div>
        <Menu mode="inline" selectedKeys={selectedKey} items={menuItems} className="seller-workspace-menu" />
      </Sider>
      <Layout>
        <Content className="seller-workspace-content">{children}</Content>
      </Layout>
    </Layout>
  );
}
