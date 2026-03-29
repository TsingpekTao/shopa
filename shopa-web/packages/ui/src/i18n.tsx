"use client";

import { createContext, ReactNode, useContext, useEffect, useMemo, useState } from "react";

export type LocaleCode = "zh-CN" | "en-US";

export type MessageDict = Record<string, Partial<Record<LocaleCode, string>>>;

type I18nContextValue = {
  locale: LocaleCode;
  setLocale: (locale: LocaleCode) => void;
  t: (key: string) => string;
};

const I18nContext = createContext<I18nContextValue | null>(null);

const STORAGE_KEY = "shopa_locale";

function detectInitialLocale(defaultLocale: LocaleCode): LocaleCode {
  if (typeof window === "undefined") {
    return defaultLocale;
  }
  const stored = window.localStorage.getItem(STORAGE_KEY);
  if (stored === "zh-CN" || stored === "en-US") {
    return stored;
  }
  const language = (window.navigator.language || "").toLowerCase();
  if (language.startsWith("zh")) {
    return "zh-CN";
  }
  return "en-US";
}

export function I18nProvider({
  children,
  messages = {},
  defaultLocale = "zh-CN"
}: {
  children: ReactNode;
  messages?: MessageDict;
  defaultLocale?: LocaleCode;
}) {
  const [locale, setLocale] = useState<LocaleCode>(defaultLocale);
  const [hasResolvedLocale, setHasResolvedLocale] = useState(false);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const detected = detectInitialLocale(defaultLocale);
    setLocale(detected);
    setHasResolvedLocale(true);
  }, [defaultLocale]);

  useEffect(() => {
    if (typeof window === "undefined" || !hasResolvedLocale) {
      return;
    }
    window.localStorage.setItem(STORAGE_KEY, locale);
    if (window.document?.documentElement) {
      window.document.documentElement.lang = locale;
    }
  }, [hasResolvedLocale, locale]);

  const value = useMemo<I18nContextValue>(() => {
    return {
      locale,
      setLocale,
      t: (key: string) => {
        const row = messages[key];
        if (!row) {
          return key;
        }
        return row[locale] ?? row[defaultLocale] ?? key;
      }
    };
  }, [defaultLocale, locale, messages]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  const context = useContext(I18nContext);
  if (!context) {
    return {
      locale: "en-US" as LocaleCode,
      setLocale: () => {},
      t: (key: string) => key
    };
  }
  return context;
}
