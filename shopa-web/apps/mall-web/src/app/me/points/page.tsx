"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { getMyOverview } from "@/features/overview/api";

function readAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const direct = window.localStorage.getItem("shopa_mall_access_token")?.trim();
  if (direct) {
    return direct;
  }
  const persisted = window.localStorage.getItem("shopa-mall-auth");
  if (!persisted) {
    return "";
  }
  try {
    const parsed = JSON.parse(persisted) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

export default function PointsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [accessToken, setAccessToken] = useState("");

  useEffect(() => {
    setAccessToken(readAccessToken());
  }, []);

  const overviewQuery = useQuery({
    queryKey: ["mall-me-overview", "points-page"],
    queryFn: getMyOverview,
    enabled: accessToken.length > 0,
    staleTime: 30_000,
    retry: 1
  });

  if (!accessToken) {
    return (
      <section className="tb-points-page">
        <div className="tb-profile-empty" role="status" aria-live="polite">
          <h2>{isZh ? "你还没有登录" : "You're not signed in"}</h2>
          <p>{isZh ? "登录后可以查看当前积分余额。" : "Sign in to view your points balance."}</p>
          <div className="tb-profile-empty-actions">
            <Link href="/login">{isZh ? "去登录" : "Go to Login"}</Link>
            <Link href="/me/profile">{isZh ? "返回个人中心" : "Back to Profile"}</Link>
          </div>
        </div>
      </section>
    );
  }

  const overview = overviewQuery.data;
  const pointsUnavailable = Boolean(overview?.partial && overview.degradedFields?.includes("points"));

  return (
    <section className="tb-points-page">
      <header className="tb-points-hero">
        <div>
          <p className="tb-profile-hero-tag">{isZh ? "积分中心" : "Points Center"}</p>
          <h1>{isZh ? "我的积分" : "My Points"}</h1>
        </div>
        <div className="tb-points-balance-card">
          <span>{isZh ? "当前可用积分" : "Available Points"}</span>
          <strong>{pointsUnavailable ? "--" : (overview?.points ?? 0)}</strong>
        </div>
      </header>

      <div className="tb-points-grid">
        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "积分余额" : "Points Balance"}</h3>
          </header>
          {overviewQuery.isLoading ? <p>{isZh ? "正在加载积分..." : "Loading points..."}</p> : null}
          {overviewQuery.isError ? <p className="tb-profile-tip is-error">{isZh ? "积分加载失败，请稍后重试。" : "Failed to load points."}</p> : null}
          {!overviewQuery.isLoading && !overviewQuery.isError ? (
            <dl className="tb-profile-list">
              <div>
                <dt>{isZh ? "当前积分" : "Current Points"}</dt>
                <dd>{pointsUnavailable ? "--" : (overview?.points ?? 0)}</dd>
              </div>
              <div>
                <dt>{isZh ? "账户状态" : "Account Status"}</dt>
                <dd>{overview?.accountStatusCode || "--"}</dd>
              </div>
              <div>
                <dt>{isZh ? "当前账号" : "Current Account"}</dt>
                <dd>{overview?.displayName || "--"}</dd>
              </div>
            </dl>
          ) : null}
        </article>

        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "快捷入口" : "Quick Access"}</h3>
          </header>
          <ul className="tb-profile-security">
            <li>
              <span>{isZh ? "个人中心" : "Profile Center"}</span>
              <Link href="/me/profile">{isZh ? "前往" : "Open"}</Link>
            </li>
            <li>
              <span>{isZh ? "我的订单" : "My Orders"}</span>
              <Link href="/me/orders">{isZh ? "查看" : "View"}</Link>
            </li>
            <li>
              <span>{isZh ? "地址管理" : "Addresses"}</span>
              <Link href="/me/address">{isZh ? "管理" : "Manage"}</Link>
            </li>
          </ul>
          {overview?.partial ? (
            <p className="tb-profile-tip">{isZh ? "部分数据正在更新中。" : "Some data is still updating."}</p>
          ) : null}
        </article>
      </div>
    </section>
  );
}
