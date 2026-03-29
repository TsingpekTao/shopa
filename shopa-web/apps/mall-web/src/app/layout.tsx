import "./globals.css";
import { ReactNode } from "react";
import { AppShell } from "@shopa/ui";
import { Providers } from "./providers";

export const metadata = {
  title: "Shopa Mall",
  description: "Shopa buyer storefront"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>
          <AppShell mode="mall">{children}</AppShell>
        </Providers>
      </body>
    </html>
  );
}
