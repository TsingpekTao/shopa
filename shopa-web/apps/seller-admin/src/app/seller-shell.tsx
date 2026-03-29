"use client";

import { ReactNode } from "react";
import { usePathname } from "next/navigation";
import { useI18n } from "@shopa/ui";

export function SellerShell({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const { locale, setLocale, t } = useI18n();

  const isPortalHome = pathname === "/";
  const isManagement = pathname.startsWith("/seller") || pathname.startsWith("/onboarding");

  if (isPortalHome) {
    return (
      <div className="app-shell">
        <main id="main-content" className="app-content" role="main">
          {children}
        </main>
      </div>
    );
  }

  return (
    <div className="app-shell">
      <header className="tb-header" style={{ height: 72 }}>
        <div className="tb-header-inner">
          <a href="/" className="tb-logo" aria-label={t("seller_admin_title")}>
            SHOPA SELLER
          </a>
          <div className="tb-header-actions">
            <a href="http://127.0.0.1:3000" className="tb-login-btn">
              {t("seller_admin_back_mall")}
            </a>
            {isManagement ? (
              <a href="/onboarding/drafts" className="tb-login-btn">
                Drafts
              </a>
            ) : null}
            <select
              className="tb-lang-select"
              value={locale}
              onChange={(e) => setLocale(e.target.value === "zh-CN" ? "zh-CN" : "en-US")}
              aria-label={t("language_switch")}
            >
              <option value="zh-CN">中文</option>
              <option value="en-US">English</option>
            </select>
          </div>
        </div>
      </header>
      <main id="main-content" className="app-content" role="main">
        {children}
      </main>
    </div>
  );
}
