"use client";

import { ReactNode, useEffect } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { notification } from "antd";
import { apiClient } from "@/lib/api-client";
import { I18nProvider, useI18n } from "@shopa/ui";
import { sellerMessages } from "@/lib/i18n-messages";

const queryClient = new QueryClient();

type Props = {
  children: ReactNode;
};

function DegradedNotifier() {
  const { t } = useI18n();

  useEffect(() => {
    notification.config({
      placement: "topRight",
      maxCount: 3
    });

    apiClient.configure({
      onDegraded: (fields) => {
        if (!fields.length) {
          return;
        }
        const degradedFields = fields.join(", ");
        notification.warning({
          message: t("notify_partial"),
          description: t("notify_degraded_fields").replace("{fields}", degradedFields),
          key: "shopa-degraded"
        });
      }
    });
  }, [t]);

  return null;
}

export function Providers({ children }: Props) {
  return (
    <I18nProvider messages={sellerMessages}>
      <QueryClientProvider client={queryClient}>
        <DegradedNotifier />
        {children}
      </QueryClientProvider>
    </I18nProvider>
  );
}
