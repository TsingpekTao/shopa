"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Button, Card, Form, Image, Input, Space, Tag, Typography, Upload, notification } from "antd";
import type { UploadProps } from "antd";
import { LoadingOutlined, PlusOutlined } from "@ant-design/icons";
import { useSearchParams } from "next/navigation";
import { useI18n } from "@shopa/ui";
import {
  createApplicationDraft,
  listMyApplications,
  getMyApplication,
  submitApplication,
  updateApplicationDraft
} from "@/features/seller-shop/api";
import {
  issueAssetReadUrl,
  OSS_CORS_BLOCKED_ERROR,
  syncOnboardingCertificateBindings,
  uploadCertificateAsset
} from "@/features/media/api";
import { useOnboardingStore } from "@/features/onboarding/store";

const { Title, Text } = Typography;

type UploadFieldKey = "idCardFrontAssetId" | "idCardBackAssetId" | "businessLicenseAssetId";

type UploadMeta = {
  loading: boolean;
  previewUrl: string;
};

type ApiErrorLike = {
  code?: number;
  message?: string;
};

const uploadFields: UploadFieldKey[] = ["idCardFrontAssetId", "idCardBackAssetId", "businessLicenseAssetId"];

const uploadTitleMap: Record<UploadFieldKey, { zh: string; en: string }> = {
  idCardFrontAssetId: { zh: "身份证人像面", en: "ID Card Front" },
  idCardBackAssetId: { zh: "身份证国徽面", en: "ID Card Back" },
  businessLicenseAssetId: { zh: "营业执照", en: "Business License" }
};

const maxUploadMb = Number(process.env.NEXT_PUBLIC_MEDIA_CERT_MAX_MB ?? 10);
const draftUpdateMask = [
  "entity.entity_name",
  "entity.contact_name",
  "entity.contact_phone",
  "entity.contact_email",
  "entity.ext",
  "shop.shop_name",
  "shop.shop_display_name",
  "shop.service_phone",
  "shop.description"
];

function isShopNameReservedError(message?: string): boolean {
  return /shop_name already reserved or used/i.test(message || "");
}

function buildPayload(state: ReturnType<typeof useOnboardingStore.getState>) {
  return {
    entity: {
      merchantTypeCode: "INDIVIDUAL",
      entityName: state.entityName,
      contactName: state.contactName,
      contactPhone: state.contactPhone,
      contactEmail: state.contactEmail,
      ext: {
        id_card_front_asset_id: state.idCardFrontAssetId,
        id_card_back_asset_id: state.idCardBackAssetId,
        business_license_asset_id: state.businessLicenseAssetId
      }
    },
    shop: {
      shopName: state.shopName,
      shopDisplayName: state.shopDisplayName,
      shopTypeCode: "STANDARD",
      servicePhone: state.servicePhone,
      description: state.description
    }
  };
}

export default function OnboardingStartPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const searchParams = useSearchParams();
  const [loadingServer, setLoadingServer] = useState(false);
  const [saving, setSaving] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [uploadMeta, setUploadMeta] = useState<Record<UploadFieldKey, UploadMeta>>({
    idCardFrontAssetId: { loading: false, previewUrl: "" },
    idCardBackAssetId: { loading: false, previewUrl: "" },
    businessLicenseAssetId: { loading: false, previewUrl: "" }
  });

  const state = useOnboardingStore();
  const applicationNoFromQuery = searchParams.get("applicationNo") ?? "";

  const setUploadLoading = (field: UploadFieldKey, loading: boolean) => {
    setUploadMeta((prev) => ({
      ...prev,
      [field]: { ...prev[field], loading }
    }));
  };

  const setUploadPreview = (field: UploadFieldKey, previewUrl: string) => {
    setUploadMeta((prev) => ({
      ...prev,
      [field]: { ...prev[field], previewUrl }
    }));
  };

  const hydratePreviewByAssetId = async (field: UploadFieldKey, assetId?: string) => {
    if (!assetId || uploadMeta[field].previewUrl) {
      return;
    }
    try {
      const readUrl = await issueAssetReadUrl(assetId);
      if (readUrl) {
        setUploadPreview(field, readUrl);
      }
    } catch (err) {
      console.warn("load asset read url failed", field, assetId, err);
    }
  };

  useEffect(() => {
    const targetNo = applicationNoFromQuery || state.applicationNo;
    if (!targetNo) {
      return;
    }

    setLoadingServer(true);
    getMyApplication(targetNo)
      .then(async (draft) => {
        state.applyServerDraft(draft);
        const ext = draft.entity?.ext || {};
        await Promise.all([
          hydratePreviewByAssetId("idCardFrontAssetId", String(ext.id_card_front_asset_id || "")),
          hydratePreviewByAssetId("idCardBackAssetId", String(ext.id_card_back_asset_id || "")),
          hydratePreviewByAssetId("businessLicenseAssetId", String(ext.business_license_asset_id || ""))
        ]);
      })
      .catch(() => {
        notification.warning({ message: isZh ? "加载服务端草稿失败，已使用本地草稿" : "Load server draft failed, using local draft" });
      })
      .finally(() => setLoadingServer(false));
  }, [applicationNoFromQuery]);

  useEffect(() => {
    uploadFields.forEach((field) => {
      void hydratePreviewByAssetId(field, state[field]);
    });
  }, [state.idCardFrontAssetId, state.idCardBackAssetId, state.businessLicenseAssetId]);

  const canSubmit = useMemo(() => {
    return Boolean(
      state.entityName &&
        state.contactName &&
        state.contactPhone &&
        state.shopName &&
        state.shopDisplayName &&
        state.idCardFrontAssetId &&
        state.idCardBackAssetId &&
        state.businessLicenseAssetId
    );
  }, [
    state.businessLicenseAssetId,
    state.contactName,
    state.contactPhone,
    state.entityName,
    state.idCardBackAssetId,
    state.idCardFrontAssetId,
    state.shopDisplayName,
    state.shopName
  ]);

  const validateUploadFile = (file: File): boolean => {
    const isImage = file.type.startsWith("image/");
    if (!isImage) {
      notification.warning({
        message: isZh ? "仅支持图片文件" : "Only image files are supported"
      });
      return false;
    }
    const isInLimit = file.size <= maxUploadMb * 1024 * 1024;
    if (!isInLimit) {
      notification.warning({
        message: isZh ? `图片大小不能超过 ${maxUploadMb}MB` : `Image size must be <= ${maxUploadMb}MB`
      });
      return false;
    }
    return true;
  };

  const onUploadCertificate = async (field: UploadFieldKey, file: File) => {
    if (!validateUploadFile(file)) {
      return;
    }

    setUploadLoading(field, true);
    try {
      const bizNo = (state.applicationNo || applicationNoFromQuery || "").trim();
      const result = await uploadCertificateAsset(file, {
        bizType: "seller_application",
        bizNo: bizNo || undefined
      });
      state.setField(field, result.assetId);
      setUploadPreview(field, result.previewUrl);
      notification.success({
        message: isZh ? "上传成功" : "Upload success",
        description: `${isZh ? "素材ID" : "Asset ID"}: ${result.assetId}`
      });
    } catch (err) {
      console.error("upload certificate failed", err);
      const message = err instanceof Error ? err.message : "";
      if (message === OSS_CORS_BLOCKED_ERROR) {
        notification.error({
          message: isZh ? "上传失败：对象存储未开启跨域" : "Upload failed: object storage CORS is disabled",
          description: isZh
            ? "请在 OSS Bucket 配置 CORS，允许当前前端域名发起 OPTIONS/PUT。"
            : "Configure OSS bucket CORS to allow OPTIONS/PUT from current frontend origin."
        });
      } else {
        notification.error({ message: isZh ? "上传失败，请稍后重试" : "Upload failed, please retry later" });
      }
    } finally {
      setUploadLoading(field, false);
    }
  };

  const buildUploadProps = (field: UploadFieldKey): UploadProps => ({
    accept: "image/jpeg,image/png,image/webp,image/gif",
    maxCount: 1,
    showUploadList: false,
    beforeUpload: async (file) => {
      await onUploadCertificate(field, file as File);
      return false;
    }
  });

  const syncBindingsIfPossible = async (applicationNo: string) => {
    if (!applicationNo) {
      return;
    }
    try {
      await syncOnboardingCertificateBindings({
        applicationNo,
        idCardFrontAssetId: state.idCardFrontAssetId,
        idCardBackAssetId: state.idCardBackAssetId,
        businessLicenseAssetId: state.businessLicenseAssetId
      });
    } catch (err) {
      console.error("sync media bindings failed", err);
      notification.warning({
        message: isZh ? "证件绑定同步失败" : "Certificate binding sync failed",
        description: isZh ? "草稿已保存，请稍后重新保存重试。" : "Draft saved, please save again later to retry."
      });
    }
  };

  const saveDraft = async (): Promise<{ applicationNo: string; version: number }> => {
    setSaving(true);
    try {
      const payload = buildPayload(state);
      let appNo = state.applicationNo;
      let version = state.expectedVersion || 0;
      if (state.applicationNo) {
        const updated = await updateApplicationDraft(state.applicationNo, {
          expectedVersion: state.expectedVersion || 1,
          ...payload,
          updateMask: draftUpdateMask
        });
        state.applyServerDraft(updated);
        appNo = updated.applicationNo || appNo;
        version = updated.version || version;
      } else {
        try {
          const created = await createApplicationDraft(payload);
          state.applyServerDraft(created);
          appNo = created.applicationNo || appNo;
          version = created.version || version;
        } catch (err) {
          const apiErr = err as ApiErrorLike;
          if (!isShopNameReservedError(apiErr?.message)) {
            throw err;
          }

          // 店铺名已被当前用户某个草稿占用时，自动复用最近草稿继续编辑。
          const mine = await listMyApplications({ page: 1, pageSize: 20, statuses: [1] });
          const fallbackDraft = (mine.applications || []).find((item) => item.applicationNo);
          if (!fallbackDraft?.applicationNo) {
            throw err;
          }
          const latest = await getMyApplication(fallbackDraft.applicationNo);
          const updated = await updateApplicationDraft(fallbackDraft.applicationNo, {
            expectedVersion: latest.version || 1,
            ...payload,
            updateMask: draftUpdateMask
          });
          state.applyServerDraft(updated);
          appNo = updated.applicationNo || fallbackDraft.applicationNo;
          version = updated.version || latest.version || version;
          notification.info({
            message: isZh ? "已复用历史草稿" : "Reused existing draft",
            description: isZh
              ? "检测到同名店铺草稿，已自动切换为更新该草稿。"
              : "An existing draft with the same shop name was found and reused."
          });
        }
      }

      await syncBindingsIfPossible(appNo);
      notification.success({
        message: isZh ? "草稿已保存" : "Draft saved",
        description: isZh ? "可前往草稿箱继续编辑。": "You can continue editing in Draft Box."
      });
      return { applicationNo: appNo, version };
    } catch (_err) {
      notification.error({ message: isZh ? "草稿保存失败" : "Failed to save draft" });
      throw _err;
    } finally {
      setSaving(false);
    }
  };

  const handleSubmit = async () => {
    if (!canSubmit) {
      notification.warning({
        message: isZh ? "请先填写必填项并上传三项证件素材" : "Please fill required fields and upload all 3 certificates"
      });
      return;
    }

    setSubmitting(true);
    try {
      let draftNo = state.applicationNo;
      let version = state.expectedVersion || 0;

      if (!draftNo) {
        notification.warning({
          message: isZh ? "请先保存草稿再提交审核" : "Please save draft before submitting",
          description: isZh
            ? "当前提交流程不会自动保存草稿，请先点击“保存草稿”。"
            : "Current submit flow does not auto-save draft, please click Save Draft first."
        });
        return;
      }

      // 提交前拉一次最新版本，避免 expectedVersion 过期导致冲突。
      const latest = await getMyApplication(draftNo);
      if (latest.version && latest.version > 0) {
        version = latest.version;
      }
      if (!version || version <= 0) {
        notification.error({
          message: isZh ? "提交失败：草稿版本无效" : "Submit failed: invalid draft version",
          description: isZh ? "请刷新页面后重试。" : "Please refresh and retry."
        });
        return;
      }

      await syncBindingsIfPossible(draftNo);
      const submitted = await submitApplication(draftNo, { expectedVersion: version });
      state.applyServerDraft(submitted);
      notification.success({ message: isZh ? "申请已提交" : "Application submitted" });
      window.location.href = `/onboarding/reviewing?applicationNo=${submitted.applicationNo}`;
    } catch (err) {
      const apiErr = err as ApiErrorLike;
      notification.error({
        message: isZh ? "提交失败" : "Submit failed",
        description: apiErr?.message ? (isZh ? `原因：${apiErr.message}` : `Reason: ${apiErr.message}`) : undefined
      });
    } finally {
      setSubmitting(false);
    }
  };

  const renderUploadCard = (field: UploadFieldKey) => {
    const label = isZh ? uploadTitleMap[field].zh : uploadTitleMap[field].en;
    const value = state[field];
    const meta = uploadMeta[field];
    return (
      <Card key={field} size="small" style={{ width: "100%" }}>
        <Space direction="vertical" size={10} style={{ width: "100%" }}>
          <Text strong>{label}</Text>
          <Upload {...buildUploadProps(field)}>
            <div
              style={{
                border: "1px dashed #d9d9d9",
                borderRadius: 8,
                minHeight: 180,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                cursor: "pointer",
                background: "#fafafa"
              }}
            >
              {meta.previewUrl ? (
                <Image src={meta.previewUrl} alt={label} width={220} height={140} style={{ objectFit: "cover" }} preview />
              ) : (
                <Space direction="vertical" size={4} align="center">
                  {meta.loading ? <LoadingOutlined /> : <PlusOutlined />}
                  <Text type="secondary">{isZh ? "点击上传图片" : "Click to upload image"}</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {isZh ? `支持 JPG/PNG/WEBP，最大 ${maxUploadMb}MB` : `JPG/PNG/WEBP up to ${maxUploadMb}MB`}
                  </Text>
                </Space>
              )}
            </div>
          </Upload>
          {value ? <Tag color="success">{`${isZh ? "素材ID" : "Asset ID"}: ${value}`}</Tag> : <Tag>{isZh ? "未上传" : "Not uploaded"}</Tag>}
        </Space>
      </Card>
    );
  };

  return (
    <section style={{ display: "grid", gap: 16 }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: 12, flexWrap: "wrap" }}>
        <div>
          <Title level={3} style={{ marginBottom: 6 }}>
            {isZh ? "商家入驻申请" : "Seller Onboarding"}
          </Title>
          <Text type="secondary">
            {isZh ? "请完善资料并提交审核，通过后可进入商家工作台。" : "Complete your profile and submit materials for review."}
          </Text>
        </div>
        <Link href="/onboarding/drafts">
          <Button>{isZh ? "进入草稿箱" : "Draft Box"}</Button>
        </Link>
      </header>

      <Card loading={loadingServer}>
        <Form layout="vertical">
          <Title level={5}>{isZh ? "主体信息" : "Entity Info"}</Title>
          <Form.Item label={isZh ? "主体名称" : "Entity Name"} required>
            <Input value={state.entityName} onChange={(e) => state.setField("entityName", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "联系人姓名" : "Contact Name"} required>
            <Input value={state.contactName} onChange={(e) => state.setField("contactName", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "联系人手机" : "Contact Phone"} required>
            <Input value={state.contactPhone} onChange={(e) => state.setField("contactPhone", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "联系人邮箱" : "Contact Email"}>
            <Input value={state.contactEmail} onChange={(e) => state.setField("contactEmail", e.target.value)} />
          </Form.Item>

          <Title level={5}>{isZh ? "店铺信息" : "Shop Info"}</Title>
          <Form.Item label={isZh ? "店铺名称" : "Shop Name"} required>
            <Input value={state.shopName} onChange={(e) => state.setField("shopName", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "店铺展示名" : "Shop Display Name"} required>
            <Input value={state.shopDisplayName} onChange={(e) => state.setField("shopDisplayName", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "客服电话" : "Service Phone"}>
            <Input value={state.servicePhone} onChange={(e) => state.setField("servicePhone", e.target.value)} />
          </Form.Item>
          <Form.Item label={isZh ? "店铺简介" : "Description"}>
            <Input.TextArea rows={3} value={state.description} onChange={(e) => state.setField("description", e.target.value)} />
          </Form.Item>

          <Title level={5}>{isZh ? "证件素材上传" : "Certificate Uploads"}</Title>
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit,minmax(240px,1fr))", gap: 12, marginBottom: 16 }}>
            {uploadFields.map((field) => renderUploadCard(field))}
          </div>

          <Space>
            <Button onClick={saveDraft} loading={saving}>
              {isZh ? "保存草稿" : "Save Draft"}
            </Button>
            <Link href="/onboarding/drafts">
              <Button>{isZh ? "查看草稿箱" : "View Drafts"}</Button>
            </Link>
            <Button type="primary" onClick={handleSubmit} loading={submitting}>
              {isZh ? "提交审核" : "Submit For Review"}
            </Button>
            <Button danger onClick={state.clearDraft}>
              {isZh ? "清空本地草稿" : "Clear Local Draft"}
            </Button>
          </Space>
        </Form>
      </Card>

      <Text type="secondary">
        {isZh ? "若审核驳回，请前往 " : "If rejected, go to "}
        <Link href="/onboarding/rejected">{isZh ? "驳回页" : "Rejected Page"}</Link>
        {isZh ? " 修改后重新提交。" : " to modify and resubmit."}
      </Text>
    </section>
  );
}
