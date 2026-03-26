import "./globals.css";
import { ReactNode } from "react";
import { AppShell } from "@shopa/ui";
import { Providers } from "./providers";

export const metadata = {
  title: "Shopa Seller Console",
  description: "Seller console for managing shops, products and inventory"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <AppShell>{children}</AppShell>
        </Providers>
      </body>
    </html>
  );
}
