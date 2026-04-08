"use client";

import { ReactNode, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import {
  AppstoreOutlined,
  CarOutlined,
  ContainerOutlined,
  DeploymentUnitOutlined,
  LineChartOutlined,
  MessageOutlined,
  PictureOutlined,
  SafetyCertificateOutlined,
  ShopOutlined,
  TagsOutlined,
  TeamOutlined,
  UserOutlined
} from "@ant-design/icons";
import { Layout, Menu, Skeleton, Typography } from "antd";
import { useI18n } from "@shopa/ui";
import { listSellerProducts } from "@/features/catalog/api";
import { listSellerConversations } from "@/features/chat/api";
import { buildSellerChatShopOptions } from "@/features/chat/shop-selection";
import { useAuth } from "@/features/iam/useAuth";
import { useAuthStore } from "@/features/iam/store";
import { listSellerSubOrders } from "@/features/order/api";
import { listSellerRefundBatches } from "@/features/refund/api";
import {
  countSellerPendingShipments,
  formatSellerNavBadgeCount,
  getSellerCurrentShopNo,
  resolveSellerCurrentShopNo,
  SELLER_CURRENT_SHOP_KEY,
  SELLER_CURRENT_SHOP_EVENT,
  sumSellerConversationUnread
} from "@/features/navigation/nav-badges";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const { Sider, Content } = Layout;
const { Text } = Typography;
const SHIPPING_BADGE_STATUSES = ["SUB_ORDER_STATUS_PAID", "SUB_ORDER_STATUS_WAIT_SHIP"];
const ACTIVE_REFUND_BADGE_STATUSES = new Set(["PENDING_SELLER_REVIEW", "WAIT_REFUND_TASK", "REFUND_PROCESSING"]);

type Props = {
  children: ReactNode;
};

function buildSellerMenuLabel(href: string, text: string, badgeCount = "") {
  return (
    <div className="seller-workspace-menu-label">
      <Link href={href}>{text}</Link>
      {badgeCount ? <span className="seller-workspace-menu-badge">{badgeCount}</span> : null}
    </div>
  );
}

export default function SellerLayout({ children }: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const { bootstrap } = useAuth();
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const [ready, setReady] = useState(false);
  const [persistedShopNo, setPersistedShopNo] = useState("");

  const workbenchQuery = useQuery({
    queryKey: ["seller-layout-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 60_000,
    enabled: ready
  });

  useEffect(() => {
    if (!ready || typeof window === "undefined") {
      return;
    }
    const syncCurrentShopNo = () => {
      setPersistedShopNo(getSellerCurrentShopNo());
    };
    syncCurrentShopNo();

    const handleStorage = (event: StorageEvent) => {
      if (event.key && event.key !== SELLER_CURRENT_SHOP_KEY) {
        return;
      }
      syncCurrentShopNo();
    };

    window.addEventListener("storage", handleStorage);
    window.addEventListener(SELLER_CURRENT_SHOP_EVENT, syncCurrentShopNo as EventListener);
    return () => {
      window.removeEventListener("storage", handleStorage);
      window.removeEventListener(SELLER_CURRENT_SHOP_EVENT, syncCurrentShopNo as EventListener);
    };
  }, [ready]);

  const workbenchShops = useMemo(() => {
    return (workbenchQuery.data?.shops ?? [])
      .map((shop) => ({
        shopNo: String(shop.shopNo ?? "").trim(),
        shopName: String(shop.shopDisplayName || shop.shopName || shop.shopNo || "").trim(),
        shopStatusCode: String(shop.shopStatusCode ?? "").trim()
      }))
      .filter((shop) => shop.shopNo);
  }, [workbenchQuery.data?.shops]);

  const fallbackCatalogShopsQuery = useQuery({
    queryKey: ["seller-layout-fallback-shops"],
    queryFn: async () => {
      const result = await listSellerProducts({ page: 1, pageSize: 100 });
      return Array.from(new Set(result.products.map((product) => product.shopNo.trim()).filter(Boolean)));
    },
    enabled: ready && (workbenchShops.length === 0 || Boolean(workbenchQuery.data?.partial)),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const shops = useMemo(() => {
    return buildSellerChatShopOptions({
      workbenchShops,
      catalogShopNos: fallbackCatalogShopsQuery.data ?? [],
      persistedShopNo
    });
  }, [fallbackCatalogShopsQuery.data, persistedShopNo, workbenchShops]);

  const currentShopNo = useMemo(() => {
    return resolveSellerCurrentShopNo(
      shops.map((shop) => shop.shopNo),
      persistedShopNo
    );
  }, [persistedShopNo, shops]);

  const conversationsBadgeQuery = useQuery({
    queryKey: ["seller-layout-conversation-badge", currentShopNo],
    queryFn: async () => {
      const result = await listSellerConversations(currentShopNo, 60, "");
      return sumSellerConversationUnread(result.list);
    },
    enabled: ready && Boolean(currentShopNo),
    staleTime: 10_000,
    refetchInterval: 12_000,
    refetchOnWindowFocus: false
  });

  const shippingBadgeQuery = useQuery({
    queryKey: ["seller-layout-shipping-badge", currentShopNo],
    queryFn: async () => {
      const result = await listSellerSubOrders({
        shopNo: currentShopNo,
        statuses: SHIPPING_BADGE_STATUSES,
        pageSize: 200
      });
      return countSellerPendingShipments(result.subOrders);
    },
    enabled: ready && Boolean(currentShopNo),
    staleTime: 10_000,
    refetchInterval: 12_000,
    refetchOnWindowFocus: false
  });

  const refundBadgeQuery = useQuery({
    queryKey: ["seller-layout-refund-badge", currentShopNo],
    queryFn: async () => {
      const result = await listSellerRefundBatches(currentShopNo);
      return result.list.filter((batch) => ACTIVE_REFUND_BADGE_STATUSES.has(batch.batchStatus)).length;
    },
    enabled: ready && Boolean(currentShopNo),
    staleTime: 10_000,
    refetchInterval: 12_000,
    refetchOnWindowFocus: false
  });

  const conversationBadge = formatSellerNavBadgeCount(conversationsBadgeQuery.data ?? 0);
  const shippingBadge = formatSellerNavBadgeCount(shippingBadgeQuery.data ?? 0);
  const refundBadge = formatSellerNavBadgeCount(refundBadgeQuery.data ?? 0);

  const menuItems = useMemo(
    () => [
      {
        key: "/seller/workbench",
        icon: <AppstoreOutlined />,
        label: buildSellerMenuLabel("/seller/workbench", isZh ? "工作台" : "Workbench")
      },
      {
        key: "/seller/sales",
        icon: <LineChartOutlined />,
        label: buildSellerMenuLabel("/seller/sales", isZh ? "销售分析" : "Sales Analytics")
      },
      {
        key: "/seller/publish",
        icon: <DeploymentUnitOutlined />,
        label: buildSellerMenuLabel("/seller/publish", isZh ? "发布商品" : "Publish Product")
      },
      {
        key: "/seller/products",
        icon: <ContainerOutlined />,
        label: buildSellerMenuLabel("/seller/products", isZh ? "商品管理" : "Products")
      },
      {
        key: "/seller/store-categories",
        icon: <TagsOutlined />,
        label: buildSellerMenuLabel("/seller/store-categories", isZh ? "店内分类" : "Store Categories")
      },
      {
        key: "/seller/media/library",
        icon: <PictureOutlined />,
        label: buildSellerMenuLabel("/seller/media/library", isZh ? "素材库" : "Media Library")
      },
      {
        key: "/seller/inventory",
        icon: <ShopOutlined />,
        label: buildSellerMenuLabel("/seller/inventory", isZh ? "库存管理" : "Inventory")
      },
      {
        key: "/seller/shipping",
        icon: <CarOutlined />,
        label: buildSellerMenuLabel("/seller/shipping", isZh ? "发货中心" : "Shipping", shippingBadge)
      },
      {
        key: "/seller/after-sales",
        icon: <SafetyCertificateOutlined />,
        label: buildSellerMenuLabel("/seller/after-sales", isZh ? "售后中心" : "After-Sales", refundBadge)
      },
      {
        key: "/seller/conversations",
        icon: <MessageOutlined />,
        label: buildSellerMenuLabel("/seller/conversations", isZh ? "会话中心" : "Conversations", conversationBadge)
      },
      {
        key: "/seller/profile",
        icon: <UserOutlined />,
        label: buildSellerMenuLabel("/seller/profile", isZh ? "店主资料" : "Profile")
      },
      {
        key: "/seller/me",
        icon: <TeamOutlined />,
        label: buildSellerMenuLabel("/seller/me", isZh ? "我的" : "My Account")
      }
    ],
    [conversationBadge, isZh, refundBadge, shippingBadge]
  );

  useEffect(() => {
    let active = true;
    (async () => {
      const accessToken = tokenPair?.accessToken ?? localStorage.getItem("shopa_seller_access_token");
      if (accessToken) {
        if (active) {
          setReady(true);
        }
        return;
      }

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
  }, [bootstrap, router, tokenPair?.accessToken]);

  const selectedKey = useMemo(() => {
    const effectivePath = pathname.startsWith("/seller/orders/") ? "/seller/shipping" : pathname;
    const found = menuItems.find((item) => effectivePath.startsWith(item.key));
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
