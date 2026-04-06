import "./globals.css";
import { ReactNode } from "react";
import { Providers } from "./providers";
import { ShellAwareAppShell } from "@/features/shell/ShellAwareAppShell";

export const metadata = {
  title: "Shopa Mall",
  description: "Shopa buyer storefront"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>
          <ShellAwareAppShell>{children}</ShellAwareAppShell>
        </Providers>
      </body>
    </html>
  );
}
