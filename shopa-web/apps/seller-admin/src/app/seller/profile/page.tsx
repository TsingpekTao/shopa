"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Button,
  Card,
  Col,
  Form,
  Input,
  Row,
  Skeleton,
  Space,
  Tag,
  Timeline,
  Typography,
  notification
} from "antd";
import { useI18n } from "@shopa/ui";
import {
  fetchSellerWorkbench,
  getMyApplication,
  listMyApplications,
  updateApplicationDraft
} from "@/features/seller-shop/api";
import { SellerApplicationDetail, SellerApplicationListItem } from "@/features/seller-shop/types";

const { Title, Text } = Typography;

type ProfileForm = {
  entityName: string;
  contactName: string;
  contactPhone: string;
  contactEmail: string;
  shopName: string;
  shopDisplayName: string;
  servicePhone: string;
  description: string;
};

const draftUpdateMask = [
  "entity.entity_name",
  "entity.contact_name",
  "entity.contact_phone",
  "entity.contact_email",
  "shop.shop_display_name",
  "shop.service_phone",
  "shop.description"
];

function pickPreferredApplication(applications: SellerApplicationListItem[]) {
  if (!applications.length) {
    return undefined;
  }

  const toTime = (input?: string) => {
    if (!input) {
      return 0;
    }
    const ms = Date.parse(input);
    return Number.isFinite(ms) ? ms : 0;
  };

  const latestOf = (codes: string[]) => {
    const matched = applications.filter((item) => codes.includes((item.applicationStatusCode || "").toUpperCase()));
    if (!matched.length) {
      return undefined;
    }
    return matched.sort((a, b) => {
      const bTime = toTime(b.updatedAt) || toTime(b.submittedAt);
      const aTime = toTime(a.updatedAt) || toTime(a.submittedAt);
      if (bTime !== aTime) {
        return bTime - aTime;
      }
      return (b.version || 0) - (a.version || 0);
    })[0];
  };

  return (
    latestOf(["REVIEWING", "SUBMITTED"]) ||
    latestOf(["APPROVED"]) ||
    latestOf(["DRAFT"]) ||
    latestOf(["REJECTED", "SYSTEM_REJECTED"]) ||
    applications[0]
  );
}

function resolveStatusTag(statusCode: string, isZh: boolean) {
  const code = (statusCode || "").toUpperCase();
  if (code === "APPROVED") return <Tag color="success">{isZh ? "已通过" : "APPROVED"}</Tag>;
  if (code === "SUBMITTED" || code === "REVIEWING") return <Tag color="processing">{isZh ? "审核中" : code}</Tag>;
  if (code === "DRAFT") return <Tag>{isZh ? "草稿" : "DRAFT"}</Tag>;
  if (code === "REJECTED" || code === "SYSTEM_REJECTED") return <Tag color="error">{isZh ? "已驳回" : code}</Tag>;
  return <Tag>{code || "UNSPECIFIED"}</Tag>;
}

function resolveShopStatusTag(statusCode: string, isZh: boolean) {
  const code = (statusCode || "").toUpperCase();
  if (code === "ACTIVE") return <Tag color="success">{isZh ? "营业中" : "ACTIVE"}</Tag>;
  if (code === "PENDING") return <Tag color="processing">{isZh ? "待开通" : "PENDING"}</Tag>;
  if (code === "DISABLED") return <Tag color="error">{isZh ? "已停用" : "DISABLED"}</Tag>;
  return <Tag>{code || "UNKNOWN"}</Tag>;
}

function calcCertCompleteness(detail: SellerApplicationDetail | null): string {
  const ext = detail?.entity?.ext || {};
  const required = ["id_card_front_asset_id", "id_card_back_asset_id", "business_license_asset_id"];
  const filled = required.filter((key) => Boolean((ext as Record<string, string>)[key])).length;
  return `${Math.round((filled / required.length) * 100)}%`;
}

export default function ProfilePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [form] = Form.useForm<ProfileForm>();

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [shopStatusCode, setShopStatusCode] = useState("");
  const [applicationNo, setApplicationNo] = useState("");
  const [applicationStatusCode, setApplicationStatusCode] = useState("");
  const [applicationVersion, setApplicationVersion] = useState(0);
  const [applicationDetail, setApplicationDetail] = useState<SellerApplicationDetail | null>(null);
  const [shopNameRaw, setShopNameRaw] = useState("");
  const [shopTypeCode, setShopTypeCode] = useState("");
  const [merchantTypeCode, setMerchantTypeCode] = useState("");
  const [serviceEmailRaw, setServiceEmailRaw] = useState("");

  const isDraftEditable = useMemo(() => (applicationStatusCode || "").toUpperCase() === "DRAFT", [applicationStatusCode]);

  useEffect(() => {
    setLoading(true);

    Promise.all([fetchSellerWorkbench(), listMyApplications({ page: 1, pageSize: 50 })])
      .then(async ([workbench, mine]) => {
        const firstShop = (workbench.shops || [])[0];
        setShopStatusCode(String((firstShop as any)?.shopStatusCode || ""));

        const target = pickPreferredApplication(mine.applications || []);
        if (!target?.applicationNo) {
          form.setFieldsValue({
            entityName: "",
            contactName: "",
            contactPhone: "",
            contactEmail: "",
            shopName: String((firstShop as any)?.shopName || ""),
            shopDisplayName: String((firstShop as any)?.shopDisplayName || (firstShop as any)?.shopName || ""),
            servicePhone: "",
            description: ""
          });
          return;
        }

        const detail = await getMyApplication(target.applicationNo);
        setApplicationDetail(detail);
        setApplicationNo(detail.applicationNo);
        setApplicationStatusCode(String(detail.applicationStatusCode || ""));
        setApplicationVersion(detail.version || 0);
        setShopNameRaw(String(detail.shop?.shopName || ""));
        setShopTypeCode(String(detail.shop?.shopTypeCode || ""));
        setMerchantTypeCode(String(detail.entity?.merchantTypeCode || ""));
        setServiceEmailRaw(String(detail.shop?.serviceEmail || ""));

        form.setFieldsValue({
          entityName: String(detail.entity?.entityName || ""),
          contactName: String(detail.entity?.contactName || ""),
          contactPhone: String(detail.entity?.contactPhone || ""),
          contactEmail: String(detail.entity?.contactEmail || ""),
          shopName: String(detail.shop?.shopName || (firstShop as any)?.shopName || ""),
          shopDisplayName: String(detail.shop?.shopDisplayName || (firstShop as any)?.shopDisplayName || ""),
          servicePhone: String(detail.shop?.servicePhone || ""),
          description: String(detail.shop?.description || "")
        });
      })
      .catch((error) => {
        console.error("load profile data failed", error);
        notification.error({ message: isZh ? "加载店主资料失败" : "Failed to load seller profile" });
      })
      .finally(() => setLoading(false));
  }, [form, isZh]);

  const saveProfile = async (values: ProfileForm) => {
    if (!applicationNo || !applicationVersion) {
      notification.warning({
        message: isZh ? "未找到可编辑申请单" : "No editable application found",
        description: isZh ? "请先在入驻流程中创建草稿。" : "Please create onboarding draft first."
      });
      return;
    }

    if (!isDraftEditable) {
      notification.info({
        message: isZh ? "当前状态不可直接修改" : "Current status is read-only",
        description: isZh ? "审核中/已通过状态请前往入驻流程页处理。" : "Please use onboarding flow to modify this state."
      });
      return;
    }

    setSaving(true);
    try {
      const updated = await updateApplicationDraft(applicationNo, {
        expectedVersion: applicationVersion,
        entity: {
          merchantTypeCode: merchantTypeCode || "INDIVIDUAL",
          entityName: values.entityName,
          contactName: values.contactName,
          contactPhone: values.contactPhone,
          contactEmail: values.contactEmail,
          ext: applicationDetail?.entity?.ext || {}
        },
        shop: {
          shopName: shopNameRaw || values.shopDisplayName,
          shopDisplayName: values.shopDisplayName,
          shopTypeCode: shopTypeCode || "STANDARD",
          servicePhone: values.servicePhone,
          serviceEmail: serviceEmailRaw || undefined,
          description: values.description,
          ext: applicationDetail?.shop?.ext || {}
        },
        updateMask: draftUpdateMask
      });

      setApplicationDetail(updated);
      setApplicationVersion(updated.version || applicationVersion);
      notification.success({ message: isZh ? "资料已同步到后端" : "Profile synced to backend" });
    } catch (error) {
      console.error("save profile failed", error);
      notification.error({ message: isZh ? "保存失败，请稍后重试" : "Save failed" });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 8 }} />;
  }

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "店主资料" : "Seller Profile"}</Title>
          <Text type="secondary">
            {isZh ? "数据来自已提交入驻申请与店铺信息。" : "Data comes from your onboarding application and shop profile."}
          </Text>
        </div>
        <Space>
          <Link href="/onboarding/start">
            <Button>{isZh ? "前往入驻页" : "Open Onboarding"}</Button>
          </Link>
          <Button type="primary" onClick={() => form.submit()} loading={saving}>
            {isZh ? "保存资料" : "Save"}
          </Button>
        </Space>
      </header>

      {!applicationNo ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "暂无可用入驻申请" : "No onboarding application found"}
          description={
            isZh ? "请先创建入驻申请并提交资料，资料页才会显示后端数据。" : "Create onboarding application first to load backend profile."
          }
        />
      ) : null}

      {!isDraftEditable && applicationNo ? (
        <Alert
          type="info"
          showIcon
          message={isZh ? "当前资料来自已提交/已审核数据" : "Viewing submitted/approved data"}
          description={
            isZh
              ? "当前申请状态不支持直接编辑，若需修改请走入驻流程。"
              : "Current application status is read-only. Use onboarding flow to modify."
          }
        />
      ) : null}

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={16}>
          <Card>
            <Form form={form} layout="vertical" onFinish={saveProfile}>
              <Row gutter={12}>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "主体名称" : "Entity Name"} name="entityName" rules={[{ required: true }]}> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "联系人" : "Contact"} name="contactName" rules={[{ required: true }]}> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "联系电话" : "Phone"} name="contactPhone" rules={[{ required: true }]}> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="Email" name="contactEmail"> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "店铺名" : "Shop Name"} name="shopName">
                    <Input disabled />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "店铺展示名" : "Shop Display Name"} name="shopDisplayName" rules={[{ required: true }]}> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label={isZh ? "客服热线" : "Service Phone"} name="servicePhone"> 
                    <Input disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
                <Col xs={24}>
                  <Form.Item label={isZh ? "店铺简介" : "Description"} name="description"> 
                    <Input.TextArea rows={4} disabled={!isDraftEditable} />
                  </Form.Item>
                </Col>
              </Row>
            </Form>
          </Card>
        </Col>

        <Col xs={24} xl={8}>
          <Card title={isZh ? "资质状态" : "Qualification"} className="seller-side-card">
            <Space direction="vertical" size={10} style={{ width: "100%" }}>
              <div className="seller-kv-item">
                <span>{isZh ? "申请单号" : "Application No"}</span>
                <Text code>{applicationNo || "-"}</Text>
              </div>
              <div className="seller-kv-item">
                <span>{isZh ? "申请状态" : "Application Status"}</span>
                {resolveStatusTag(applicationStatusCode, isZh)}
              </div>
              <div className="seller-kv-item">
                <span>{isZh ? "店铺状态" : "Store Status"}</span>
                {resolveShopStatusTag(shopStatusCode, isZh)}
              </div>
              <div className="seller-kv-item">
                <span>{isZh ? "证件完整度" : "Certificate Completeness"}</span>
                <Tag color="processing">{calcCertCompleteness(applicationDetail)}</Tag>
              </div>
            </Space>
          </Card>

          <Card title={isZh ? "最近动态" : "Recent Timeline"} className="seller-side-card">
            <Timeline
              items={[
                {
                  children:
                    (isZh ? "提交时间：" : "Submitted At: ") +
                    (applicationDetail?.submittedAt ? applicationDetail.submittedAt.replace("T", " ").slice(0, 19) : "-")
                },
                {
                  children:
                    (isZh ? "最近更新时间：" : "Last Updated: ") +
                    (applicationDetail?.updatedAt ? applicationDetail.updatedAt.replace("T", " ").slice(0, 19) : "-")
                },
                {
                  children: isZh ? "当前资料来自后端申请单数据" : "Profile is synced from backend application"
                }
              ]}
            />
          </Card>
        </Col>
      </Row>
    </section>
  );
}
