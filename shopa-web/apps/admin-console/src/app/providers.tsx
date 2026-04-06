"use client";

import { ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { I18nProvider } from "@shopa/ui";
import { adminMessages } from "@/lib/i18n-messages";

const queryClient = new QueryClient();

export function Providers({ children }: { children: ReactNode }) {
  return (
    <I18nProvider messages={adminMessages}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </I18nProvider>
  );
}
