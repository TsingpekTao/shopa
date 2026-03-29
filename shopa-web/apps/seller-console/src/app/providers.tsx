"use client";

import { ReactNode, useEffect } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { notification } from "antd";
import { apiClient } from "@shopa/api-client";

const queryClient = new QueryClient();

type Props = {
  children: ReactNode;
};

export function Providers({ children }: Props) {
  useEffect(() => {
    notification.config({
      placement: "topRight",
      maxCount: 3
    });

    apiClient.configure({
      getToken: () => {
        const token = localStorage.getItem("shopa_access_token");
        if (token) {
          return token;
        }
        const persisted = localStorage.getItem("shopa-auth");
        if (!persisted) {
          return null;
        }
        try {
          const parsed = JSON.parse(persisted) as {
            state?: { tokenPair?: { accessToken?: string; refreshToken?: string } };
          };
          const accessToken = parsed?.state?.tokenPair?.accessToken ?? null;
          const refreshToken = parsed?.state?.tokenPair?.refreshToken ?? null;
          if (accessToken) {
            localStorage.setItem("shopa_access_token", accessToken);
          }
          if (refreshToken) {
            localStorage.setItem("shopa_refresh_token", refreshToken);
          }
          return accessToken;
        } catch (_err) {
          return null;
        }
      },
      onUnauthorized: () => {
        localStorage.removeItem("shopa_access_token");
        localStorage.removeItem("shopa_refresh_token");
        localStorage.removeItem("shopa-auth");
        window.location.href = "/login";
      },
      onDegraded: (fields) => {
        if (!fields.length) {
          return;
        }
        notification.warning({
          message: "Partial data returned",
          description: `Degraded fields: ${fields.join(", ")}`,
          key: "shopa-degraded"
        });
      }
    });
  }, []);

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}
