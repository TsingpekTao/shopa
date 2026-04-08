"use client";

import { ReactNode, useCallback, useEffect, useMemo, useState, type FormEvent } from "react";
import { Avatar, Button, ConfigProvider, Dropdown, Layout, Space } from "antd";
import type { MenuProps } from "antd";
import "./styles.css";
import { useI18n } from "./i18n";
import { buildGlobalSearchHref, normalizeGlobalSearchTab, type GlobalSearchTab } from "./search-form";

type Props = {
  children: ReactNode;
  mode?: "default" | "mall";
  shellState?: {
    cartTypeCount?: number;
    unreadMessages?: number;
  };
};

type MallAuthSnapshot = {
  loggedIn: boolean;
  displayName: string;
  avatarText: string;
  avatarUrl: string;
};

type ProfileSummaryResponse = {
  code?: number;
  data?: {
    profile?: {
      display_name?: string;
      displayName?: string;
      avatar?: {
        asset_id?: number | string;
        assetId?: number | string;
        url?: string;
      };
      ext?: Record<string, string>;
    };
  };
  profile?: {
    display_name?: string;
    displayName?: string;
    avatar?: {
      asset_id?: number | string;
      assetId?: number | string;
      url?: string;
    };
    ext?: Record<string, string>;
  };
};

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";
const ACCESS_TOKEN_KEY = "shopa_mall_access_token";
const LAST_IDENTIFIER_KEY = "shopa_mall_last_identifier";
const DISPLAY_NAME_KEY = "shopa_mall_display_name";
const AVATAR_URL_KEY = "shopa_mall_avatar_url";
const PERSISTED_AUTH_KEY = "shopa-mall-auth";

function go(path: string) {
  if (typeof window === "undefined") {
    return;
  }
  window.location.href = path;
}

function resolveMallApiBaseURL(): string {
  if (typeof window === "undefined") {
    return "http://127.0.0.1:8000";
  }
  const envBaseURL = process.env.NEXT_PUBLIC_API_BASE_URL?.trim();
  if (envBaseURL) {
    return envBaseURL;
  }
  const { protocol, hostname } = window.location;
  return `${protocol}//${hostname}:8000`;
}

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

function decodeBase64Url(value: string): string {
  const normalized = value.replaceAll("-", "+").replaceAll("_", "/");
  const padding = normalized.length % 4 === 0 ? "" : "=".repeat(4 - (normalized.length % 4));
  return atob(`${normalized}${padding}`);
}

function tryDecodeIdentityFromToken(token: string): string {
  const parts = token.split(".");
  if (parts.length < 2) {
    return "";
  }
  try {
    const payload = JSON.parse(decodeBase64Url(parts[1])) as Record<string, unknown>;
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

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

async function issueMallAssetReadUrl(assetId: string, ttlSeconds = 900): Promise<string> {
  if (!assetId) {
    return "";
  }

  try {
    const params = new URLSearchParams({
      assetIds: assetId,
      ttlSeconds: String(ttlSeconds)
    });
    const response = await fetch(`/api/media/read-urls?${params.toString()}`, {
      method: "GET",
      cache: "no-store"
    });
    if (!response.ok) {
      return "";
    }

    const payload = (await response.json()) as {
      items?: Array<{ assetId?: string | number; asset_id?: string | number; url?: string }>;
    };
    const item = (payload.items ?? []).find((candidate) => toString(candidate.assetId ?? candidate.asset_id) === assetId);
    return toString(item?.url);
  } catch {
    return "";
  }
}

async function fetchMallProfileSummary(token: string): Promise<{ displayName: string; avatarUrl: string }> {
  if (!token) {
    return { displayName: "", avatarUrl: "" };
  }

  const response = await fetch(`${resolveMallApiBaseURL()}/v1/me/profile?include_addresses=false`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`
    },
    cache: "no-store"
  });
  if (!response.ok) {
    return { displayName: "", avatarUrl: "" };
  }

  const payload = (await response.json()) as ProfileSummaryResponse;
  if (typeof payload.code === "number" && payload.code !== 0) {
    return { displayName: "", avatarUrl: "" };
  }

  const profile = payload.data?.profile ?? payload.profile;
  const fromProfile = profile?.display_name ?? profile?.displayName ?? "";
  const avatarAssetId = toString(profile?.avatar?.asset_id ?? profile?.avatar?.assetId);
  const avatarUrl = (await issueMallAssetReadUrl(avatarAssetId)) || profile?.avatar?.url?.trim() || "";

  const normalized = normalizeDisplayName(fromProfile);
  if (normalized) {
    return { displayName: normalized, avatarUrl };
  }

  const fromExt = profile?.ext;
  if (!fromExt) {
    return { displayName: "", avatarUrl: "" };
  }
  const fallback = [fromExt.nickname, fromExt.username, fromExt.email, fromExt.phone].find((item) => typeof item === "string");
  return {
    displayName: typeof fallback === "string" ? normalizeDisplayName(fallback) : "",
    avatarUrl
  };
}

function readMallAuthSnapshot(): MallAuthSnapshot {
  if (typeof window === "undefined") {
    return { loggedIn: false, displayName: "", avatarText: "S", avatarUrl: "" };
  }

  const token = window.localStorage.getItem(ACCESS_TOKEN_KEY)?.trim() || tryReadPersistedAccessToken();
  const loggedIn = token.length > 0;
  const displayNameCandidates = [
    window.localStorage.getItem(DISPLAY_NAME_KEY)?.trim() ?? "",
    window.localStorage.getItem(LAST_IDENTIFIER_KEY)?.trim() ?? "",
    tryDecodeIdentityFromToken(token)
  ];
  const displayName = displayNameCandidates.map((item) => normalizeDisplayName(item)).find((item) => item.length > 0) ?? "";
  const avatarText = (displayName || "S").slice(0, 1).toUpperCase();
  const avatarUrl = window.localStorage.getItem(AVATAR_URL_KEY)?.trim() ?? "";

  return {
    loggedIn,
    displayName,
    avatarText,
    avatarUrl
  };
}

export function AppShell({ children, mode = "default", shellState }: Props) {
  const { locale, setLocale, t } = useI18n();
  const isMallMode = mode === "mall";
  const sellerEntryUrl = "http://127.0.0.1:3100/";
  const [hasHydrated, setHasHydrated] = useState(false);
  const [searchTab, setSearchTab] = useState<GlobalSearchTab>("item");
  const [mallAuth, setMallAuth] = useState<MallAuthSnapshot>({
    loggedIn: false,
    displayName: "",
    avatarText: "S",
    avatarUrl: ""
  });

  useEffect(() => {
    setHasHydrated(true);
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const params = new URLSearchParams(window.location.search);
    setSearchTab(normalizeGlobalSearchTab(params.get("tab")));
  }, []);

  const refreshMallAuth = useCallback(() => {
    if (!isMallMode) {
      return;
    }
    setMallAuth(readMallAuthSnapshot());
  }, [isMallMode]);

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
    if (!hasHydrated || !isMallMode || !mallAuth.loggedIn) {
      return;
    }

    const token = window.localStorage.getItem(ACCESS_TOKEN_KEY)?.trim() || tryReadPersistedAccessToken();
    if (!token) {
      return;
    }

    let mounted = true;
    fetchMallProfileSummary(token)
      .then(({ displayName, avatarUrl }) => {
        if (!mounted) {
          return;
        }
        if (displayName) {
          window.localStorage.setItem(DISPLAY_NAME_KEY, displayName);
        }
        if (avatarUrl) {
          window.localStorage.setItem(AVATAR_URL_KEY, avatarUrl);
        }
        if (displayName || avatarUrl) {
          window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
        }
      })
      .catch(() => {
        // Ignore background hydration failures and keep cached summary.
      });

    return () => {
      mounted = false;
    };
  }, [hasHydrated, isMallMode, mallAuth.loggedIn]);

  const isMallLoggedIn = isMallMode && hasHydrated && mallAuth.loggedIn;
  const mallDisplayName = useMemo(() => {
    if (!isMallLoggedIn) {
      return "";
    }
    return mallAuth.displayName || (locale === "zh-CN" ? "\u5df2\u767b\u5f55\u7528\u6237" : "Signed-in user");
  }, [isMallLoggedIn, locale, mallAuth.displayName]);

  const cartTypeCount = shellState?.cartTypeCount ?? 0;
  const unreadMessages = shellState?.unreadMessages ?? 0;
  const cartBadgeText = cartTypeCount > 99 ? "99+" : String(cartTypeCount);
  const unreadBadgeText = unreadMessages > 99 ? "99+" : String(unreadMessages);

  const settingMenu: MenuProps = {
    items: [
      { key: "/settings", label: t("settings_center") },
      { key: "/me/address", label: t("settings_address") },
      { key: "/me/orders", label: t("settings_orders") },
      { key: "/me/points", label: locale === "zh-CN" ? "我的积分" : "My Points" },
      { key: "/me/history", label: t("settings_history") }
    ],
    onClick: ({ key }) => go(String(key))
  };

  const userMenu: MenuProps = {
    items: [
      { key: "/me/profile", label: t("profile_view") },
      { key: "/me/orders", label: t("profile_orders") },
      { key: "/me/points", label: locale === "zh-CN" ? "我的积分" : "My Points" },
      ...(isMallLoggedIn ? [] : [{ key: "/login", label: t("profile_switch") }])
    ],
    onClick: ({ key }) => go(String(key))
  };

  const helpMenu: MenuProps = {
    items: [
      { key: "/help?tab=official", label: t("help_official") },
      { key: "/help?tab=merchant", label: t("help_merchant") },
      { key: "/help?tab=history", label: t("help_history") }
    ],
    onClick: ({ key }) => go(String(key))
  };

  const handleSearchSubmit = useCallback(
    (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      const formData = new FormData(event.currentTarget);
      go(buildGlobalSearchHref(toString(formData.get("q")), searchTab));
    },
    [searchTab]
  );

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
                <a href="/me/profile">{locale === "zh-CN" ? `\u4f60\u597d\uff0c${mallDisplayName}` : `Hi, ${mallDisplayName}`}</a>
              ) : (
                <a href="/login">{t("topbar_login_tip")}</a>
              )}

              <a href="/cart" className="tb-topbar-action tb-topbar-action--with-badge" aria-label={t("topbar_cart")}>
                <span>{t("topbar_cart")}</span>
                {cartTypeCount > 0 ? <span className="tb-topbar-badge tb-topbar-badge--alert">{cartBadgeText}</span> : null}
              </a>

              <a href="/me/orders" className="tb-topbar-action">
                {t("topbar_orders")}
              </a>

              <a href="/me/messages" className="tb-topbar-action tb-topbar-action--with-badge" aria-label={t("topbar_messages")}>
                <span>{t("topbar_messages")}</span>
                {unreadMessages > 0 ? <span className="tb-topbar-badge tb-topbar-badge--alert">{unreadBadgeText}</span> : null}
              </a>

              <Dropdown menu={helpMenu} trigger={["hover", "click"]}>
                <button type="button" className="tb-topbar-drop" aria-label={t("topbar_help")}>
                  {t("topbar_help")}
                </button>
              </Dropdown>

              <a href="/search?tab=favorites">{t("topbar_favorites")}</a>
              <a href={sellerEntryUrl} target="_blank" rel="noreferrer">
                {t("topbar_open_store")}
              </a>
            </div>

            <div className="tb-topbar-links">
              <a href={isMallLoggedIn ? "/me/profile" : "/register"}>
                {isMallLoggedIn ? (locale === "zh-CN" ? "\u4e2a\u4eba\u4e2d\u5fc3" : "Profile Center") : t("topbar_register")}
              </a>
              <a href="/contact">{t("topbar_contact")}</a>
              <select
                className="tb-lang-select"
                value={locale}
                onChange={(event) => setLocale(event.target.value === "zh-CN" ? "zh-CN" : "en-US")}
                aria-label={t("language_switch")}
              >
                <option value="zh-CN">{"\u4e2d\u6587"}</option>
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

            <form className="tb-search" action="/search" method="get" onSubmit={handleSearchSubmit}>
              <div className="tb-search-mode" role="tablist" aria-label={locale === "zh-CN" ? "搜索类型" : "Search type"}>
                <button
                  type="button"
                  className={searchTab === "item" ? "tb-search-mode-btn is-active" : "tb-search-mode-btn"}
                  aria-pressed={searchTab === "item"}
                  onClick={() => setSearchTab("item")}
                >
                  {locale === "zh-CN" ? "商品" : "Items"}
                </button>
                <button
                  type="button"
                  className={searchTab === "shop" ? "tb-search-mode-btn is-active" : "tb-search-mode-btn"}
                  aria-pressed={searchTab === "shop"}
                  onClick={() => setSearchTab("shop")}
                >
                  {locale === "zh-CN" ? "店铺" : "Shops"}
                </button>
              </div>
              <input type="hidden" name="tab" value={searchTab} />
              <label htmlFor="global-search" className="tb-sr-only">
                {t("search_label")}
              </label>
              <input id="global-search" name="q" className="tb-search-input" placeholder={t("search_placeholder")} autoComplete="off" />
              <button type="submit">{t("search_button")}</button>
            </form>

            <div className="tb-header-actions">
              {!isMallLoggedIn ? (
                <a href="/login" className="tb-login-btn">
                  {t("header_login")}
                </a>
              ) : null}

              <Dropdown menu={settingMenu} trigger={["click"]} placement="bottomRight">
                <Button className="tb-setting-btn" type="default">
                  {t("header_settings")}
                </Button>
              </Dropdown>

              <Dropdown menu={userMenu} trigger={["click"]} placement="bottomRight">
                <Space className="tb-user-btn" size={8}>
                  <Avatar size={28} src={isMallLoggedIn ? mallAuth.avatarUrl || undefined : undefined}>
                    {isMallLoggedIn ? mallAuth.avatarText : "T"}
                  </Avatar>
                  <span className="tb-user-name">{isMallLoggedIn ? mallDisplayName : t("header_profile")}</span>
                </Space>
              </Dropdown>
            </div>
          </div>

          <nav className="tb-nav" aria-label={t("mall_nav")}>
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
            <a href={sellerEntryUrl} target="_blank" rel="noreferrer">
              {t("footer_open_store")}
            </a>
          </div>
          <p className="tb-footer-text">{t("footer_copy")}</p>
        </Layout.Footer>
      </Layout>
    </ConfigProvider>
  );
}
