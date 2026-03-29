"use client";

import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Alert, Button, Card, Result, Skeleton, Space, Tag, Typography, notification } from "antd";
import { CheckCircleFilled, ClockCircleOutlined, ReloadOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { getMyApplication } from "@/features/seller-shop/api";
import { SellerApplicationDetail } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

const POLL_INTERVAL_MS = 8000;

function normalizeStatusCode(value?: string): string {
  return (value || "").toUpperCase();
}

function statusTag(statusCode: string, isZh: boolean) {
  switch (statusCode) {
    case "SUBMITTED":
      return <Tag color="processing">{isZh ? "已提交" : "SUBMITTED"}</Tag>;
    case "REVIEWING":
      return <Tag color="gold">{isZh ? "审核中" : "REVIEWING"}</Tag>;
    case "APPROVED":
      return <Tag color="success">{isZh ? "已通过" : "APPROVED"}</Tag>;
    case "REJECTED":
      return <Tag color="error">{isZh ? "已驳回" : "REJECTED"}</Tag>;
    case "SYSTEM_REJECTED":
      return <Tag color="volcano">{isZh ? "系统驳回" : "SYSTEM_REJECTED"}</Tag>;
    default:
      return <Tag>{isZh ? "未知" : "UNKNOWN"}</Tag>;
  }
}

export default function OnboardingReviewingPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const applicationNo = (searchParams.get("applicationNo") ?? "").trim();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [application, setApplication] = useState<SellerApplicationDetail | null>(null);
  const [nowTs, setNowTs] = useState(Date.now());
  const redirectedRef = useRef(false);

  const statusCode = normalizeStatusCode(application?.applicationStatusCode);

  const pollApplication = async (silent = false) => {
    if (!applicationNo) {
      setLoading(false);
      return;
    }
    if (!silent) {
      setRefreshing(true);
    }
    try {
      const result = await getMyApplication(applicationNo);
      setApplication(result);
    } catch (error) {
      if (!silent) {
        notification.error({
          message: isZh ? "获取审核状态失败" : "Failed to fetch review status"
        });
      }
      console.error("load application reviewing status failed", error);
    } finally {
      setLoading(false);
      setRefreshing(false);
      setNowTs(Date.now());
    }
  };

  useEffect(() => {
    void pollApplication(true);
    const timer = window.setInterval(() => {
      void pollApplication(true);
    }, POLL_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, [applicationNo]);

  useEffect(() => {
    if (redirectedRef.current) {
      return;
    }
    if (statusCode !== "APPROVED") {
      return;
    }
    redirectedRef.current = true;
    notification.success({
      message: isZh ? "审核通过，欢迎进入商家工作台" : "Approved, welcome to seller workbench",
      description: isZh ? "正在为你跳转..." : "Redirecting..."
    });
    const timer = window.setTimeout(() => router.replace("/seller/workbench"), 1200);
    return () => window.clearTimeout(timer);
  }, [statusCode, isZh, router]);

  const updatedAtText = useMemo(() => {
    if (!application?.updatedAt) {
      return "-";
    }
    const date = new Date(application.updatedAt);
    if (Number.isNaN(date.getTime())) {
      return application.updatedAt;
    }
    return new Intl.DateTimeFormat(locale, {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false
    }).format(date);
  }, [application?.updatedAt, locale]);

  if (!applicationNo) {
    return (
      <section className="tb-onboard-status-wrap">
        <Result
          status="warning"
          title={isZh ? "缺少申请单号" : "Missing Application Number"}
          subTitle={isZh ? "请从入驻申请页重新进入审核页。" : "Please return from onboarding application page."}
          extra={
            <Link href="/onboarding/start">
              <Button type="primary">{isZh ? "去填写申请" : "Go to onboarding"}</Button>
            </Link>
          }
        />
      </section>
    );
  }

  if (loading) {
    return (
      <section className="tb-onboard-status-wrap">
        <Card>
          <Skeleton active paragraph={{ rows: 7 }} />
        </Card>
      </section>
    );
  }

  if (statusCode === "REJECTED" || statusCode === "SYSTEM_REJECTED") {
    return (
      <section className="tb-onboard-status-wrap">
        <Alert
          type="warning"
          showIcon
          message={isZh ? "申请已被驳回" : "Application was rejected"}
          description={isZh ? "请根据驳回原因补充资料后重新提交。" : "Please revise your materials and resubmit."}
        />
        <Card>
          <Space direction="vertical" size={12}>
            <Text>
              {isZh ? "申请单号：" : "Application No: "}
              <Text code>{applicationNo}</Text>
            </Text>
            <Link href={`/onboarding/rejected?applicationNo=${applicationNo}`}>
              <Button type="primary">{isZh ? "去修改并重提" : "Edit and resubmit"}</Button>
            </Link>
          </Space>
        </Card>
      </section>
    );
  }

  if (statusCode === "APPROVED") {
    return (
      <section className="tb-onboard-status-wrap">
        <Card className="tb-onboard-success-hero">
          <Result
            icon={<CheckCircleFilled style={{ color: "#16a34a" }} />}
            title={isZh ? "恭喜，已正式成为商家" : "Congratulations, you are now a seller"}
            subTitle={
              isZh
                ? "你的店铺资质审核已通过，现在可以进入商家工作台管理商品和店铺。"
                : "Your onboarding is approved. You can now manage products and shop in workbench."
            }
            extra={
              <Space wrap>
                <Button type="primary" size="large" onClick={() => router.replace("/seller/workbench")}>
                  {isZh ? "立即进入工作台" : "Open Workbench"}
                </Button>
                <Link href="/seller/publish">
                  <Button size="large">{isZh ? "开始上架商品" : "Publish Product"}</Button>
                </Link>
              </Space>
            }
          />
        </Card>
      </section>
    );
  }

  return (
    <section className="tb-onboard-status-wrap">
      <header className="tb-onboard-status-head">
        <div>
          <Title level={3} style={{ marginBottom: 6 }}>
            {isZh ? "申请审核中" : "Application Under Review"}
          </Title>
          <Text type="secondary">
            {isZh
              ? "审核通过后会自动跳转到商家工作台。你也可以手动刷新状态。"
              : "You will be redirected to seller workbench once approved. You can also refresh manually."}
          </Text>
        </div>
        <Button icon={<ReloadOutlined />} loading={refreshing} onClick={() => void pollApplication(false)}>
          {isZh ? "刷新状态" : "Refresh Status"}
        </Button>
      </header>

      <div className="tb-onboard-status-grid">
        <Card className="tb-onboard-main-card">
          <div className="tb-onboard-step-title">
            <ClockCircleOutlined />
            <span>{isZh ? "审核进度" : "Review Progress"}</span>
          </div>
          <ul className="tb-onboard-step-list">
            <li className="is-done">
              <strong>{isZh ? "资料已提交" : "Materials Submitted"}</strong>
              <span>{isZh ? "平台已接收入驻资料。" : "Platform has received your onboarding materials."}</span>
            </li>
            <li className="is-active">
              <strong>{isZh ? "人工审核中" : "Manual Review in Progress"}</strong>
              <span>{isZh ? "审核员正在核验资质文件与店铺信息。" : "Reviewer is checking qualification files and shop information."}</span>
            </li>
            <li>
              <strong>{isZh ? "审核通过后开通店铺" : "Shop Activation After Approval"}</strong>
              <span>{isZh ? "通过后会自动跳转至商家工作台。" : "You will be redirected to seller workbench automatically."}</span>
            </li>
          </ul>
        </Card>

        <Card className="tb-onboard-side-card">
          <Space direction="vertical" size={10} style={{ width: "100%" }}>
            <div>
              <Text type="secondary">{isZh ? "申请单号" : "Application No"}</Text>
              <div>
                <Text code>{applicationNo}</Text>
              </div>
            </div>
            <div>
              <Text type="secondary">{isZh ? "当前状态" : "Current Status"}</Text>
              <div>{statusTag(statusCode, isZh)}</div>
            </div>
            <div>
              <Text type="secondary">{isZh ? "最近更新时间" : "Last Updated At"}</Text>
              <div>{updatedAtText}</div>
            </div>
            <div>
              <Text type="secondary">{isZh ? "本次轮询时间" : "Current Poll Time"}</Text>
              <div>{new Date(nowTs).toLocaleTimeString(locale === "zh-CN" ? "zh-CN" : "en-US", { hour12: false })}</div>
            </div>
          </Space>
        </Card>
      </div>
    </section>
  );
}
