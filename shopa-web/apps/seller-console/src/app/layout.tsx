import "./globals.css";
import { ReactNode } from "react";
import { AppShell } from "@shopa/ui";
import { Providers } from "./providers";

export const metadata = {
  title: "Shopa Mall",
  description: "Shopa 商城前端，覆盖买家浏览与卖家入口"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>
          <AppShell>{children}</AppShell>
        </Providers>
      </body>
    </html>
  );
}
