"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Alert, Button, Card, Typography } from "antd";
import { useI18n } from "@shopa/ui";

const { Title, Text } = Typography;

export default function OnboardingRejectedPage() {
  const searchParams = useSearchParams();
  const applicationNo = searchParams.get("applicationNo") ?? "";
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <section style={{ display: "grid", gap: 16 }}>
      <Title level={3} style={{ marginBottom: 0 }}>
        {isZh ? "申请被驳回" : "Application Rejected"}
      </Title>
      <Card>
        <Alert
          type="warning"
          showIcon
          message={isZh ? "你的商家申请未通过" : "Your seller application was rejected"}
          description={isZh ? "请补充或修改入驻资料后再次提交。" : "Please update your onboarding materials and submit again."}
        />
        <div style={{ marginTop: 12 }}>
          <Text>
            {isZh ? "申请单号：" : "Application No: "}
            {applicationNo || "N/A"}
          </Text>
        </div>
        <div style={{ marginTop: 12 }}>
          <Link href={`/onboarding/start${applicationNo ? `?applicationNo=${applicationNo}` : ""}`}>
            <Button type="primary">{isZh ? "修改并重新提交" : "Edit and Resubmit"}</Button>
          </Link>
        </div>
      </Card>
    </section>
  );
}
