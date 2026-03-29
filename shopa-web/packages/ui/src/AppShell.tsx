"use client";

import { ReactNode, useCallback, useEffect, useMemo, useState } from "react";
import { Avatar, Button, ConfigProvider, Dropdown, Layout, Space } from "antd";
import type { MenuProps } from "antd";
import "./styles.css";
import { useI18n } from "./i18n";

type Props = {
  children: ReactNode;
  mode?: "default" | "mall";
};

function go(path: string) {
  if (typeof window === "undefined") {
    return;
  }
  window.location.href = path;
}

type MallAuthSnapshot = {
  loggedIn: boolean;
  displayName: string;
  avatarText: string;
};

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";
const ACCESS_TOKEN_KEY = "shopa_mall_access_token";
const LAST_IDENTIFIER_KEY = "shopa_mall_last_identifier";
const DISPLAY_NAME_KEY = "shopa_mall_display_name";
const PERSISTED_AUTH_KEY = "shopa-mall-auth";

function tryReadPersistedAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const raw = window.localStorage.getItem(PERSISTED_AUTH_KEY);
  if (!raw) {
    return "";
  }
  try {
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function tryDecodeIdentityFromToken(token: string): string {
  const parts = token.split(".");
  if (parts.length < 2) {
    return "";
  }
  try {
    const payload = JSON.parse(atob(parts[1])) as Record<string, unknown>;
    const fields = [
      payload.display_name,
      payload.displayName,
      payload.nickname,
      payload.name,
      payload.username,
      payload.account,
      payload.phone,
      payload.email
    ];
    for (const field of fields) {
      if (typeof field === "string" && field.trim()) {
        return field.trim();
      }
    }
    return "";
  } catch {
    return "";
  }
}

function normalizeDisplayName(raw: string): string {
  const value = raw.trim();
  if (!value) {
    return "";
  }
  if (/^\d{11}$/.test(value)) {
    return `${value.slice(0, 3)}****${value.slice(-4)}`;
  }
  if (/^\d+$/.test(value)) {
    return "";
  }
  return value;
}

async function fetchMallProfileDisplayName(token: string): Promise<string> {
  if (!token) {
    return "";
  }
  const response = await fetch("/v1/me/profile?include_addresses=false", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`
    }
  });
  if (!response.ok) {
    return "";
  }
  const payload = (await response.json()) as {
    code?: number;
    data?: { profile?: { display_name?: string; displayName?: string; ext?: Record<string, string> } };
    profile?: { display_name?: string; displayName?: string; ext?: Record<string, string> };
  };
  if (typeof payload.code === "number" && payload.code !== 0) {
    return "";
  }
  const profile = payload.data?.profile ?? payload.profile;
  const fromProfile = profile?.display_name ?? profile?.displayName ?? "";
  if (typeof fromProfile === "string") {
    const normalized = normalizeDisplayName(fromProfile);
    if (normalized) {
      return normalized;
    }
  }
  const fromExt = profile?.ext;
  if (fromExt) {
    const fallback = [fromExt.nickname, fromExt.username, fromExt.email, fromExt.phone].find((item) => typeof item === "string");
    if (typeof fallback === "string") {
      return normalizeDisplayName(fallback);
    }
  }
  return "";
}

function readMallAuthSnapshot(locale: string): MallAuthSnapshot {
  if (typeof window === "undefined") {
    return { loggedIn: false, displayName: "", avatarText: "S" };
  }
  const token = window.localStorage.getItem(ACCESS_TOKEN_KEY)?.trim() || tryReadPersistedAccessToken();
  const loggedIn = token.length > 0;
  const candidates = [
    window.localStorage.getItem(DISPLAY_NAME_KEY)?.trim() ?? "",
    window.localStorage.getItem(LAST_IDENTIFIER_KEY)?.trim() ?? "",
    tryDecodeIdentityFromToken(token)
  ];
  const displayName = candidates.map((item) => normalizeDisplayName(item)).find((item) => item.length > 0) ?? "";
  const avatarText = (displayName || "S").slice(0, 1).toUpperCase();
  return { loggedIn, displayName, avatarText };
}

export function AppShell({ children, mode = "default" }: Props) {
  const { locale, setLocale, t } = useI18n();
  const isMallMode = mode === "mall";
  const sellerEntryUrl = "http://127.0.0.1:3100/";
  const [hasHydrated, setHasHydrated] = useState(false);
  const [mallAuth, setMallAuth] = useState<MallAuthSnapshot>({
    loggedIn: false,
    displayName: "",
    avatarText: "S"
  });

  useEffect(() => {
    setHasHydrated(true);
  }, []);

  const refreshMallAuth = useCallback(() => {
    if (!isMallMode) {
      return;
    }
    setMallAuth(readMallAuthSnapshot(locale));
  }, [isMallMode, locale]);

  useEffect(() => {
    if (!isMallMode || !hasHydrated) {
      return;
    }
    refreshMallAuth();
    const onStorage = () => refreshMallAuth();
    const onAuthChanged = () => refreshMallAuth();
    window.addEventListener("storage", onStorage);
    window.addEventListener(AUTH_CHANGED_EVENT, onAuthChanged);
    return () => {
      window.removeEventListener("storage", onStorage);
      window.removeEventListener(AUTH_CHANGED_EVENT, onAuthChanged);
    };
  }, [hasHydrated, isMallMode, refreshMallAuth]);

  useEffect(() => {
    if (!hasHydrated || !isMallMode || !mallAuth.loggedIn || mallAuth.displayName) {
      return;
    }
    const token = window.localStorage.getItem(ACCESS_TOKEN_KEY)?.trim() || tryReadPersistedAccessToken();
    if (!token) {
      return;
    }
    let mounted = true;
    fetchMallProfileDisplayName(token)
      .then((name) => {
        if (!mounted || !name) {
          return;
        }
        window.localStorage.setItem(DISPLAY_NAME_KEY, name);
        window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
      })
      .catch(() => {
        // ignore profile hydrate errors
      });
    return () => {
      mounted = false;
    };
  }, [hasHydrated, isMallMode, mallAuth.displayName, mallAuth.loggedIn]);

  const isMallLoggedIn = isMallMode && hasHydrated && mallAuth.loggedIn;
  const mallDisplayName = useMemo(() => {
    if (!isMallLoggedIn) {
      return "";
    }
    return mallAuth.displayName || (locale === "zh-CN" ? "已登录用户" : "Signed-in user");
  }, [isMallLoggedIn, mallAuth.displayName, locale]);

  const settingMenu: MenuProps = {
    items: [
      { key: "/settings", label: t("settings_center") },
      { key: "/me/address", label: t("settings_address") },
      { key: "/me/orders", label: t("settings_orders") },
      { key: "/me/history", label: t("settings_history") }
    ],
    onClick: ({ key }) => go(String(key))
  };

  const userMenu: MenuProps = {
    items: [
      { key: "/me/profile", label: t("profile_view") },
      { key: "/me/orders", label: t("profile_orders") },
      ...(isMallLoggedIn ? [] : [{ key: "/login", label: t("profile_switch") }])
    ],
    onClick: ({ key }) => {
      go(String(key));
    }
  };

  const helpMenu: MenuProps = {
    items: [
      { key: "/help?tab=official", label: t("help_official") },
      { key: "/help?tab=merchant", label: t("help_merchant") },
      { key: "/help?tab=history", label: t("help_history") }
    ],
    onClick: ({ key }) => go(String(key))
  };

  return (
    <ConfigProvider
      componentSize="small"
      theme={{
        token: {
          colorPrimary: "#ff5000",
          borderRadius: 8,
          colorText: "#3c3c3c",
          colorBgContainer: "#ffffff"
        }
      }}
    >
      <a className="skip-link" href="#main-content">
        {t("skip_main")}
      </a>
      <Layout className="app-shell" style={{ fontSize: 13, lineHeight: 1.45 }}>
        <div className="tb-topbar" role="banner">
          <div className="tb-topbar-inner">
            <div className="tb-topbar-links">
              {isMallLoggedIn ? (
                <a href="/me/profile">{locale === "zh-CN" ? `你好，${mallDisplayName}` : `Hi, ${mallDisplayName}`}</a>
              ) : (
                <a href="/login">{t("topbar_login_tip")}</a>
              )}
              <a href="/cart">{t("topbar_cart")}</a>
              <a href="/me/orders">{t("topbar_orders")}</a>
              <Dropdown menu={helpMenu} trigger={["hover", "click"]}>
                <button type="button" className="tb-topbar-drop" aria-label={t("topbar_help")}>{t("topbar_help")}</button>
              </Dropdown>
              <a href="/search?tab=favorites">{t("topbar_favorites")}</a>
              <a href={sellerEntryUrl} target="_blank" rel="noreferrer">
                {t("topbar_open_store")}
              </a>
            </div>
            <div className="tb-topbar-links">
              <a href={isMallLoggedIn ? "/me/profile" : "/register"}>
                {isMallLoggedIn ? (locale === "zh-CN" ? "个人中心" : "Profile Center") : t("topbar_register")}
              </a>
              <a href="/contact">{t("topbar_contact")}</a>
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
        </div>

        <Layout.Header className="tb-header" role="banner">
          <div className="tb-header-inner">
            <a href="/" className="tb-logo" aria-label={t("logo_aria")}>
              SHOPA
            </a>
            <form className="tb-search" action="/search" method="get">
              <label htmlFor="global-search" className="tb-sr-only">
                {t("search_label")}
              </label>
              <input id="global-search" name="q" placeholder={t("search_placeholder")} autoComplete="off" />
              <button type="submit">{t("search_button")}</button>
            </form>

            <div className="tb-header-actions">
              {!isMallLoggedIn && (
                <a href="/login" className="tb-login-btn">
                  {t("header_login")}
                </a>
              )}
              <Dropdown menu={settingMenu} trigger={["click"]} placement="bottomRight">
                <Button className="tb-setting-btn" type="default">
                  {t("header_settings")}
                </Button>
              </Dropdown>
              <Dropdown menu={userMenu} trigger={["click"]} placement="bottomRight">
                <Space className="tb-user-btn" size={8}>
                  <Avatar size={28}>{isMallLoggedIn ? mallAuth.avatarText : "T"}</Avatar>
                  <span className="tb-user-name">{isMallLoggedIn ? mallDisplayName : t("header_profile")}</span>
                </Space>
              </Dropdown>
            </div>
          </div>

          <nav className="tb-nav" aria-label={t("mall_nav")}
          >
            <a href="/">{t("nav_home")}</a>
            <a href="/search?q=electronics">{t("nav_electronics")}</a>
            <a href="/search?q=fashion">{t("nav_fashion")}</a>
            <a href="/search?q=home">{t("nav_home_living")}</a>
            <a href="/search?q=beauty">{t("nav_beauty")}</a>
            <a href="/search?q=fresh">{t("nav_fresh")}</a>
          </nav>
        </Layout.Header>

        <Layout.Content id="main-content" className="app-content" role="main" tabIndex={-1}>
          {children}
        </Layout.Content>

        <Layout.Footer className="tb-footer">
          <div className="tb-footer-links">
            <a href="/about">{t("footer_about")}</a>
            <a href="/contact">{t("footer_contact")}</a>
            <a href="/settings">{t("footer_settings")}</a>
            <a href="/service">{t("footer_service")}</a>
            <a href="/help">{t("footer_help")}</a>
            <a href={sellerEntryUrl} target="_blank" rel="noreferrer">{t("footer_open_store")}</a>
          </div>
          <p className="tb-footer-text">{t("footer_copy")}</p>
        </Layout.Footer>
      </Layout>
    </ConfigProvider>
  );
}

