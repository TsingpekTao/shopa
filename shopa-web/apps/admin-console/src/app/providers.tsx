"use client";

import { ReactNode, useEffect } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { notification } from "antd";
import { I18nProvider, useI18n } from "@shopa/ui";
import { apiClient } from "@/lib/api-client";
import { adminMessages } from "@/lib/i18n-messages";

const queryClient = new QueryClient();

function DegradedNotifier() {
  const { t } = useI18n();

  useEffect(() => {
    notification.config({ placement: "topRight", maxCount: 3 });
    apiClient.configure({
      onDegraded: (fields) => {
        if (!fields.length) {
          return;
        }
        notification.warning({
          message: t("notify_partial"),
          description: t("notify_degraded_fields").replace("{fields}", fields.join(", ")),
          key: "shopa-admin-degraded"
        });
      }
    });
  }, [t]);

  return null;
}

export function Providers({ children }: { children: ReactNode }) {
  return (
    <I18nProvider messages={adminMessages}>
      <QueryClientProvider client={queryClient}>
        <DegradedNotifier />
        {children}
      </QueryClientProvider>
    </I18nProvider>
  );
}
