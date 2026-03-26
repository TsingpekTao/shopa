"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { AppstoreOutlined, ContainerOutlined, DeploymentUnitOutlined, PictureOutlined, ShopOutlined, UserOutlined } from "@ant-design/icons";
import { Layout, Menu, Skeleton, Typography } from "antd";
import { useAuthStore } from "@/features/iam/store";

const { Sider, Content } = Layout;
const { Text } = Typography;

type Props = {
  children: ReactNode;
};

const menuItems = [
  { key: "/seller/workbench", icon: <AppstoreOutlined />, label: <Link href="/seller/workbench">Workbench</Link> },
  { key: "/seller/publish", icon: <DeploymentUnitOutlined />, label: <Link href="/seller/publish">Publish Product</Link> },
  { key: "/seller/products", icon: <ContainerOutlined />, label: <Link href="/seller/products">Products</Link> },
  { key: "/seller/media/library", icon: <PictureOutlined />, label: <Link href="/seller/media/library">Media</Link> },
  { key: "/seller/inventory", icon: <ShopOutlined />, label: <Link href="/seller/inventory">Inventory</Link> },
  { key: "/seller/profile", icon: <UserOutlined />, label: <Link href="/seller/profile">Profile</Link> }
];

export default function SellerLayout({ children }: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const accessToken = tokenPair?.accessToken ?? localStorage.getItem("shopa_access_token");
    if (!accessToken) {
      router.replace("/login");
      return;
    }
    setReady(true);
  }, [router, tokenPair?.accessToken]);

  const selectedKey = useMemo(() => {
    const found = menuItems.find((item) => pathname.startsWith(item.key));
    return found ? [found.key] : ["/seller/workbench"];
  }, [pathname]);

  if (!ready) {
    return <Skeleton active />;
  }

  return (
    <Layout style={{ minHeight: "calc(100vh - 118px)", marginTop: 24 }}>
      <Sider width={250} style={{ background: "#fff", borderRadius: 12, border: "1px solid #ffe0c9" }}>
        <div style={{ padding: "16px 16px 8px" }}>
          <Text strong style={{ color: "#8a4a16" }}>
            Seller Navigation
          </Text>
        </div>
        <Menu mode="inline" selectedKeys={selectedKey} items={menuItems} aria-label="Seller workspace navigation" />
      </Sider>
      <Layout>
        <Content
          style={{
            marginLeft: 16,
            padding: 24,
            background: "#fff",
            borderRadius: 12,
            border: "1px solid #ffe7d3"
          }}
        >
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}
