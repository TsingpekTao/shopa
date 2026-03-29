import "./globals.css";
import { ReactNode } from "react";
import { Providers } from "./providers";
import { SellerShell } from "./seller-shell";

export const metadata = {
  title: "Shopa Seller Admin",
  description: "Shopa seller backend"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>
          <SellerShell>{children}</SellerShell>
        </Providers>
      </body>
    </html>
  );
}
