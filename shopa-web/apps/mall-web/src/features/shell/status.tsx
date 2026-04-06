"use client";

import { ReactNode, createContext, useContext, useEffect, useMemo, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getMyCart } from "@/features/cart/api";
import { getUnreadSummary } from "@/features/chat/api";
import type { MyCart } from "@/features/cart/types";

const SHELL_CART_QUERY_KEY = ["shell", "cart"] as const;
const SHELL_UNREAD_QUERY_KEY = ["shell", "unread"] as const;
const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";

const emptyCart: MyCart = {
  items: [],
  summary: {
    totalItemCount: 0,
    checkedItemCount: 0,
    checkedGoodsAmount: 0,
    checkedPayableAmount: 0
  },
  loadedFromBackup: false
};

export type ShellStatus = {
  cartTypeCount: number;
  unreadMessages: number;
  isLoading: boolean;
};

const ShellStatusContext = createContext<ShellStatus>({
  cartTypeCount: 0,
  unreadMessages: 0,
  isLoading: true
});

function readMallAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const direct = window.localStorage.getItem("shopa_mall_access_token")?.trim() ?? "";
  if (direct) {
    return direct;
  }
  try {
    const raw = window.localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

export function ShellStatusProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    const refresh = () => {
      const nextHasToken = readMallAccessToken().length > 0;
      setHasToken(nextHasToken);
      if (nextHasToken) {
        void queryClient.invalidateQueries({ queryKey: SHELL_CART_QUERY_KEY });
        void queryClient.invalidateQueries({ queryKey: SHELL_UNREAD_QUERY_KEY });
      } else {
        queryClient.setQueryData(SHELL_CART_QUERY_KEY, emptyCart);
        queryClient.setQueryData(SHELL_UNREAD_QUERY_KEY, {
          totalUnreadConversations: 0,
          totalUnreadMessages: 0
        });
      }
    };

    refresh();
    window.addEventListener("storage", refresh);
    window.addEventListener(AUTH_CHANGED_EVENT, refresh);
    return () => {
      window.removeEventListener("storage", refresh);
      window.removeEventListener(AUTH_CHANGED_EVENT, refresh);
    };
  }, [queryClient]);

  const cartQuery = useQuery({
    queryKey: SHELL_CART_QUERY_KEY,
    queryFn: async () => {
      try {
        return await getMyCart();
      } catch (error) {
        console.error("shell cart failed", error);
        return emptyCart;
      }
    },
    enabled: hasToken,
    staleTime: 60_000,
    refetchInterval: hasToken ? 60_000 : false,
    refetchOnWindowFocus: false,
    retry: false
  });

  const unreadQuery = useQuery({
    queryKey: SHELL_UNREAD_QUERY_KEY,
    queryFn: async () => {
      try {
        return await getUnreadSummary();
      } catch (error) {
        console.error("shell unread summary failed", error);
        return { totalUnreadConversations: 0, totalUnreadMessages: 0 };
      }
    },
    enabled: hasToken,
    staleTime: 30_000,
    refetchInterval: hasToken ? 20_000 : false,
    refetchOnWindowFocus: false,
    retry: false
  });

  const cartTypeCount = useMemo(() => {
    const items = cartQuery.data?.items ?? [];
    const types = new Set<string>();
    items.forEach((item) => {
      const candidate = item.skuNo || item.spuNo || `${item.shopNo}-${item.skuNo}`;
      if (candidate) {
        types.add(candidate);
      }
    });
    return types.size;
  }, [cartQuery.data]);

  const unreadMessages = useMemo(() => unreadQuery.data?.totalUnreadMessages ?? 0, [unreadQuery.data]);

  const state = useMemo(
    () => ({
      cartTypeCount,
      unreadMessages,
      isLoading: hasToken && (cartQuery.isLoading || unreadQuery.isLoading)
    }),
    [cartTypeCount, unreadMessages, cartQuery.isLoading, hasToken, unreadQuery.isLoading]
  );

  return <ShellStatusContext.Provider value={state}>{children}</ShellStatusContext.Provider>;
}

export function useShellStatus() {
  return useContext(ShellStatusContext);
}
