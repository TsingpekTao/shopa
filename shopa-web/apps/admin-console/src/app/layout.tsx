import "./globals.css";
import "antd/dist/reset.css";
import { ReactNode } from "react";
import { Providers } from "./providers";

export const metadata = {
  title: "Shopa Admin Console",
  description: "Shopa admin backend"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
