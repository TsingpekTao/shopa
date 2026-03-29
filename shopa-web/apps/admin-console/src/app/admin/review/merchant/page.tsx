"use client";

import { useEffect, useMemo, useState } from "react";
import { Button, Descriptions, Empty, Image, Input, Modal, Select, Space, Spin, Table, Tag, Typography, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useQuery } from "@tanstack/react-query";
import { AccessControl, useI18n } from "@shopa/ui";
import {
  approveMerchantApplication,
  fetchAssetReadUrl,
  fetchMerchantApplicationDetail,
  fetchMerchantApplications,
  rejectMerchantApplication
} from "@/features/admin/api";
import { Permissions } from "@/features/admin/permissions";
import { MerchantApplication, MerchantEntityProfile, QualificationDoc } from "@/features/admin/types";
import { AdminPage } from "@/components/admin-page";

const { Text, Paragraph } = Typography;
const { TextArea } = Input;

type CredentialAsset = {
  key: string;
  label: string;
  assetId: string;
};

type RejectTarget = {
  applicationNo: string;
  expectedVersion: number;
};

const REJECT_REASON_OPTIONS = [
  { value: "MATERIAL_INVALID", label: "MATERIAL_INVALID" },
  { value: "MATERIAL_INCOMPLETE", label: "MATERIAL_INCOMPLETE" },
  { value: "IDENTITY_MISMATCH", label: "IDENTITY_MISMATCH" },
  { value: "RISK_CONTROL", label: "RISK_CONTROL" }
];

const STATUS_FILTER_OPTIONS = [
  { value: 1, label: "DRAFT" },
  { value: 2, label: "SUBMITTED" },
  { value: 3, label: "REVIEWING" },
  { value: 4, label: "APPROVED" },
  { value: 5, label: "REJECTED" },
  { value: 6, label: "CANCELLED" },
  { value: 7, label: "SYSTEM_REJECTED" }
];

function statusTag(status: number | undefined, isZh: boolean) {
  switch (status) {
    case 1:
      return <Tag color="default">{isZh ? "草稿" : "DRAFT"}</Tag>;
    case 2:
      return <Tag color="processing">{isZh ? "已提交" : "SUBMITTED"}</Tag>;
    case 3:
      return <Tag color="gold">{isZh ? "审核中" : "REVIEWING"}</Tag>;
    case 4:
      return <Tag color="success">{isZh ? "已通过" : "APPROVED"}</Tag>;
    case 5:
      return <Tag color="error">{isZh ? "已驳回" : "REJECTED"}</Tag>;
    case 6:
      return <Tag>{isZh ? "已取消" : "CANCELLED"}</Tag>;
    case 7:
      return <Tag color="volcano">{isZh ? "系统驳回" : "SYSTEM_REJECTED"}</Tag>;
    default:
      return <Tag>{isZh ? "未知" : "UNKNOWN"}</Tag>;
  }
}

function displayValue(value: string | number | undefined | null): string {
  if (value === null || value === undefined) {
    return "-";
  }
  const text = String(value).trim();
  return text.length > 0 ? text : "-";
}

function formatDateTime(value: string | undefined, locale: string): string {
  if (!value) {
    return "-";
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return value;
  }
  return new Intl.DateTimeFormat(locale, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false
  }).format(parsed);
}

function toAssetId(value: unknown): string | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  const normalized = String(value).trim();
  if (!/^\d+$/.test(normalized)) {
    return undefined;
  }
  return normalized;
}

function docTypeLabel(docTypeCode: string | undefined, isZh: boolean): string {
  const normalized = (docTypeCode ?? "").toUpperCase();
  if (normalized.includes("ID_CARD_FRONT")) {
    return isZh ? "身份证人像面" : "ID Card Front";
  }
  if (normalized.includes("ID_CARD_BACK")) {
    return isZh ? "身份证国徽面" : "ID Card Back";
  }
  if (normalized.includes("BUSINESS_LICENSE")) {
    return isZh ? "营业执照" : "Business License";
  }
  if (normalized.length > 0) {
    return normalized;
  }
  return isZh ? "资质文件" : "Qualification File";
}

function findExtAssetId(ext: Record<string, string> | undefined, keys: string[]): string | undefined {
  if (!ext) {
    return undefined;
  }
  for (const key of keys) {
    const value = toAssetId(ext[key]);
    if (value) {
      return value;
    }
  }
  return undefined;
}

function extractCredentialAssets(entity: MerchantEntityProfile | undefined, isZh: boolean): CredentialAsset[] {
  if (!entity) {
    return [];
  }
  const items: CredentialAsset[] = [];
  const ext = entity.ext;
  const extDefinitions: Array<{ key: string; label: string; keys: string[] }> = [
    {
      key: "id_card_front",
      label: isZh ? "身份证人像面" : "ID Card Front",
      keys: ["id_card_front_asset_id", "idCardFrontAssetId"]
    },
    {
      key: "id_card_back",
      label: isZh ? "身份证国徽面" : "ID Card Back",
      keys: ["id_card_back_asset_id", "idCardBackAssetId"]
    },
    {
      key: "business_license",
      label: isZh ? "营业执照" : "Business License",
      keys: ["business_license_asset_id", "businessLicenseAssetId"]
    }
  ];

  for (const definition of extDefinitions) {
    const assetId = findExtAssetId(ext, definition.keys);
    if (!assetId) {
      continue;
    }
    items.push({
      key: definition.key,
      label: definition.label,
      assetId
    });
  }

  const legalDocs: QualificationDoc[] = entity.legalSubject?.docs ?? [];
  legalDocs.forEach((doc, index) => {
    const assetId = toAssetId(doc.assetId);
    if (!assetId) {
      return;
    }
    items.push({
      key: `doc_${index}`,
      label: docTypeLabel(doc.docTypeCode, isZh),
      assetId
    });
  });

  const deduped = new Map<string, CredentialAsset>();
  items.forEach((item) => {
    const uniqueKey = `${item.label}_${item.assetId}`;
    if (!deduped.has(uniqueKey)) {
      deduped.set(uniqueKey, item);
    }
  });
  return Array.from(deduped.values());
}

function renderExtList(ext: Record<string, string> | undefined) {
  if (!ext || Object.keys(ext).length === 0) {
    return <Text type="secondary">-</Text>;
  }
  return (
    <div className="admin-ext-list">
      {Object.entries(ext).map(([key, value]) => (
        <div key={key} className="admin-ext-item">
          <Text code>{key}</Text>
          <span>{value}</span>
        </div>
      ))}
    </div>
  );
}

export default function MerchantReviewPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [loadingNo, setLoadingNo] = useState("");
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailApplication, setDetailApplication] = useState<MerchantApplication | null>(null);
  const [assetUrlMap, setAssetUrlMap] = useState<Record<string, string>>({});
  const [assetUrlLoading, setAssetUrlLoading] = useState(false);
  const [rejectOpen, setRejectOpen] = useState(false);
  const [rejectSubmitting, setRejectSubmitting] = useState(false);
  const [rejectTarget, setRejectTarget] = useState<RejectTarget | null>(null);
  const [rejectReasonCode, setRejectReasonCode] = useState("MATERIAL_INVALID");
  const [rejectComment, setRejectComment] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [keywordInput, setKeywordInput] = useState("");
  const [keyword, setKeyword] = useState("");
  const [statusFilters, setStatusFilters] = useState<number[]>([]);

  const query = useQuery({
    queryKey: ["admin", "merchant-applications", page, pageSize, keyword, statusFilters.join(",")],
    queryFn: () =>
      fetchMerchantApplications({
        page,
        pageSize,
        keyword: keyword || undefined,
        statuses: statusFilters.length ? statusFilters : undefined
      })
  });

  const detailEntity = detailApplication?.entitySubmitted ?? detailApplication?.entityDraft ?? detailApplication?.entity;
  const detailShop = detailApplication?.shopSubmitted ?? detailApplication?.shopDraft ?? detailApplication?.shop;
  const credentialAssets = useMemo(() => extractCredentialAssets(detailEntity, isZh), [detailEntity, isZh]);

  useEffect(() => {
    if (!detailOpen || credentialAssets.length === 0) {
      setAssetUrlMap({});
      setAssetUrlLoading(false);
      return;
    }
    let cancelled = false;
    setAssetUrlLoading(true);
    Promise.all(
      credentialAssets.map(async (item) => {
        try {
          const response = await fetchAssetReadUrl(item.assetId);
          return [item.assetId, response.url] as const;
        } catch {
          return [item.assetId, ""] as const;
        }
      })
    )
      .then((entries) => {
        if (cancelled) {
          return;
        }
        const nextMap: Record<string, string> = {};
        entries.forEach(([assetId, url]) => {
          if (url) {
            nextMap[assetId] = url;
          }
        });
        setAssetUrlMap(nextMap);
      })
      .finally(() => {
        if (!cancelled) {
          setAssetUrlLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [credentialAssets, detailOpen]);

  const submitReject = async () => {
    if (!rejectTarget?.applicationNo) {
      message.error(isZh ? "申请单号为空，无法驳回" : "Missing application number");
      return;
    }
    if (!rejectComment.trim()) {
      message.warning(isZh ? "请填写驳回原因说明" : "Please provide reject comment");
      return;
    }
    try {
      setRejectSubmitting(true);
      setLoadingNo(`${rejectTarget.applicationNo}-reject`);
      await rejectMerchantApplication({
        applicationNo: rejectTarget.applicationNo,
        expectedVersion: rejectTarget.expectedVersion,
        rejectReasonCode,
        rejectComment: rejectComment.trim()
      });
      message.success(isZh ? "驳回成功，已附带驳回原因" : "Rejected with reason");
      setRejectOpen(false);
      setRejectTarget(null);
      setRejectComment("");
      await query.refetch();
    } finally {
      setRejectSubmitting(false);
      setLoadingNo("");
    }
  };

  const applyFilters = () => {
    setPage(1);
    setKeyword(keywordInput.trim());
  };

  const resetFilters = () => {
    setPage(1);
    setPageSize(20);
    setKeywordInput("");
    setKeyword("");
    setStatusFilters([]);
  };

  const columns: ColumnsType<MerchantApplication> = useMemo(
    () => [
      { title: "Application No", dataIndex: "applicationNo", key: "applicationNo", width: 220 },
      { title: isZh ? "主体名称" : "Entity", dataIndex: ["entity", "entityName"], key: "entityName", width: 180 },
      { title: isZh ? "联系人" : "Contact", dataIndex: ["entity", "contactName"], key: "contactName", width: 120 },
      { title: isZh ? "联系电话" : "Phone", dataIndex: ["entity", "contactPhone"], key: "contactPhone", width: 150 },
      {
        title: isZh ? "状态" : "Status",
        key: "status",
        width: 140,
        render: (_, row) => statusTag(row.status, isZh)
      },
      {
        title: isZh ? "操作" : "Actions",
        key: "actions",
        width: 400,
        render: (_, row) => (
          <Space>
            <Button
              size="small"
              onClick={async () => {
                if (!row.applicationNo) {
                  message.error(isZh ? "申请单号为空，无法查看详情" : "Missing application number");
                  return;
                }
                setDetailOpen(true);
                setDetailLoading(true);
                setDetailApplication(null);
                setAssetUrlMap({});
                try {
                  const detail = await fetchMerchantApplicationDetail(row.applicationNo);
                  setDetailApplication(detail.application);
                } catch (err: any) {
                  message.error(isZh ? `获取详情失败：${err?.message || "unknown"}` : `Failed to load detail: ${err?.message || "unknown"}`);
                } finally {
                  setDetailLoading(false);
                }
              }}
            >
              {isZh ? "查看详情" : "Detail"}
            </Button>

            <AccessControl require={Permissions.MerchantReviewApprove}>
              <Button
                type="primary"
                size="small"
                loading={loadingNo === `${row.applicationNo}-approve`}
                onClick={async () => {
                  try {
                    setLoadingNo(`${row.applicationNo}-approve`);
                    await approveMerchantApplication({ applicationNo: row.applicationNo, expectedVersion: row.version ?? 1 });
                    message.success(isZh ? "已提交通过操作" : "Approve submitted");
                    await query.refetch();
                  } finally {
                    setLoadingNo("");
                  }
                }}
              >
                {isZh ? "通过" : "Approve"}
              </Button>
            </AccessControl>

            <AccessControl require={Permissions.MerchantReviewReject}>
              <Button
                danger
                size="small"
                loading={loadingNo === `${row.applicationNo}-reject`}
                onClick={() => {
                  setRejectTarget({
                    applicationNo: row.applicationNo,
                    expectedVersion: row.version ?? 1
                  });
                  setRejectReasonCode("MATERIAL_INVALID");
                  setRejectComment("");
                  setRejectOpen(true);
                }}
              >
                {isZh ? "驳回" : "Reject"}
              </Button>
            </AccessControl>
          </Space>
        )
      }
    ],
    [isZh, loadingNo, query]
  );

  return (
    <AdminPage
      title={isZh ? "商家入驻审核" : "Merchant Onboarding Review"}
      subtitle={isZh ? "审核商家资质材料，执行通过或驳回操作。" : "Review merchant qualification materials and process decisions."}
    >
      <Space wrap size={12} style={{ marginBottom: 12 }}>
        <Input
          allowClear
          value={keywordInput}
          onChange={(event) => setKeywordInput(event.target.value)}
          onPressEnter={applyFilters}
          style={{ width: 320 }}
          placeholder={isZh ? "搜索主体名称/店铺名称" : "Search entity/shop name"}
        />
        <Select
          mode="multiple"
          allowClear
          maxTagCount="responsive"
          style={{ width: 320 }}
          value={statusFilters}
          options={STATUS_FILTER_OPTIONS}
          onChange={(values) => {
            setPage(1);
            setStatusFilters(values);
          }}
          placeholder={isZh ? "状态筛选" : "Filter by status"}
        />
        <Button type="primary" onClick={applyFilters}>
          {isZh ? "搜索" : "Search"}
        </Button>
        <Button onClick={resetFilters}>{isZh ? "重置" : "Reset"}</Button>
      </Space>

      <div className="admin-table">
        <Table
          rowKey="applicationNo"
          loading={query.isLoading}
          dataSource={query.data?.applications ?? []}
          columns={columns}
          pagination={{
            current: query.data?.page ?? 1,
            pageSize: query.data?.pageSize ?? 20,
            total: query.data?.total ?? 0,
            showSizeChanger: true,
            pageSizeOptions: [10, 20, 50, 100]
          }}
          onChange={(pagination) => {
            setPage(pagination.current ?? 1);
            setPageSize(pagination.pageSize ?? pageSize);
          }}
          scroll={{ x: 1280 }}
        />
      </div>

      <Modal
        open={rejectOpen}
        title={isZh ? "填写驳回原因" : "Reject Application"}
        onCancel={() => {
          if (rejectSubmitting) {
            return;
          }
          setRejectOpen(false);
          setRejectTarget(null);
          setRejectComment("");
        }}
        onOk={submitReject}
        confirmLoading={rejectSubmitting}
        destroyOnClose
      >
        <Space direction="vertical" size={12} style={{ width: "100%" }}>
          <div>
            <Text type="secondary">{isZh ? "申请单号" : "Application No"}</Text>
            <div>{rejectTarget?.applicationNo ?? "-"}</div>
          </div>
          <div>
            <Text type="secondary">{isZh ? "驳回原因码" : "Reject Reason Code"}</Text>
            <Select
              value={rejectReasonCode}
              style={{ width: "100%", marginTop: 6 }}
              options={REJECT_REASON_OPTIONS}
              onChange={(value) => setRejectReasonCode(value)}
            />
          </div>
          <div>
            <Text type="secondary">{isZh ? "驳回说明（必填）" : "Reject Comment (Required)"}</Text>
            <TextArea
              value={rejectComment}
              onChange={(event) => setRejectComment(event.target.value)}
              showCount
              maxLength={300}
              rows={5}
              style={{ marginTop: 6 }}
              placeholder={isZh ? "请输入驳回原因，商家可看到该说明" : "Explain the rejection reason for merchant."}
            />
          </div>
        </Space>
      </Modal>

      <Modal
        open={detailOpen}
        onCancel={() => {
          setDetailOpen(false);
          setDetailApplication(null);
        }}
        title={isZh ? "商家申请详情" : "Merchant Application Detail"}
        width={1040}
        footer={
          <Button
            onClick={() => {
              setDetailOpen(false);
              setDetailApplication(null);
            }}
          >
            {isZh ? "关闭" : "Close"}
          </Button>
        }
        destroyOnClose
      >
        {detailLoading ? (
          <div className="admin-detail-loading">
            <Spin />
          </div>
        ) : !detailApplication ? (
          <Empty description={isZh ? "暂无详情" : "No detail data"} />
        ) : (
          <div className="admin-merchant-detail">
            <section className="admin-detail-block">
              <h4 className="admin-detail-title">{isZh ? "申请信息" : "Application Info"}</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label={isZh ? "申请单号" : "Application No"}>{displayValue(detailApplication.applicationNo)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "状态" : "Status"}>{statusTag(detailApplication.status, isZh)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "版本" : "Version"}>{displayValue(detailApplication.version)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "申请人用户ID" : "Owner User ID"}>{displayValue(detailApplication.ownerUserId)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "审核人用户ID" : "Reviewer User ID"}>{displayValue(detailApplication.reviewerId)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "创建时间" : "Created At"}>{formatDateTime(detailApplication.createdAt, locale)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "提交时间" : "Submitted At"}>{formatDateTime(detailApplication.submittedAt, locale)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "开始审核时间" : "Review Started At"}>
                  {formatDateTime(detailApplication.reviewStartedAt, locale)}
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "审核完成时间" : "Reviewed At"}>{formatDateTime(detailApplication.reviewedAt, locale)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "更新时间" : "Updated At"}>{formatDateTime(detailApplication.updatedAt, locale)}</Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">{isZh ? "主体信息" : "Entity Info"}</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label={isZh ? "主体编号" : "Entity No"}>{displayValue(detailEntity?.entityNo)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "商户类型" : "Merchant Type"}>{displayValue(detailEntity?.merchantTypeCode)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "主体名称" : "Entity Name"}>{displayValue(detailEntity?.entityName)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "联系人" : "Contact Name"}>{displayValue(detailEntity?.contactName)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "联系电话" : "Contact Phone"}>{displayValue(detailEntity?.contactPhone)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "联系邮箱" : "Contact Email"}>{displayValue(detailEntity?.contactEmail)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "主体扩展信息" : "Entity Ext"} span={2}>
                  {renderExtList(detailEntity?.ext)}
                </Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">{isZh ? "店铺信息" : "Shop Info"}</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label={isZh ? "店铺编号" : "Shop No"}>{displayValue(detailShop?.shopNo)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "店铺类型" : "Shop Type"}>{displayValue(detailShop?.shopTypeCode)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "店铺名称" : "Shop Name"}>{displayValue(detailShop?.shopName)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "展示名称" : "Display Name"}>{displayValue(detailShop?.shopDisplayName)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "客服电话" : "Service Phone"}>{displayValue(detailShop?.servicePhone)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "客服邮箱" : "Service Email"}>{displayValue(detailShop?.serviceEmail)}</Descriptions.Item>
                <Descriptions.Item label={isZh ? "主营类目ID" : "Main Category IDs"} span={2}>
                  {displayValue(detailShop?.mainCategoryIds?.join(", "))}
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "店铺描述" : "Description"} span={2}>
                  <Paragraph className="admin-detail-paragraph">{displayValue(detailShop?.description)}</Paragraph>
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "店铺扩展信息" : "Shop Ext"} span={2}>
                  {renderExtList(detailShop?.ext)}
                </Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">{isZh ? "资质材料" : "Qualification Materials"}</h4>
              {credentialAssets.length === 0 ? (
                <Text type="secondary">{isZh ? "未提交可识别的证件素材" : "No recognizable credential assets submitted."}</Text>
              ) : (
                <div className="admin-cert-grid">
                  {credentialAssets.map((asset) => {
                    const previewUrl = assetUrlMap[asset.assetId];
                    return (
                      <article className="admin-cert-card" key={`${asset.key}_${asset.assetId}`}>
                        <div className="admin-cert-head">
                          <strong>{asset.label}</strong>
                          <Tag>{`asset#${asset.assetId}`}</Tag>
                        </div>
                        {previewUrl ? (
                          <Image className="admin-cert-image" src={previewUrl} alt={asset.label} />
                        ) : (
                          <div className="admin-cert-placeholder">
                            {assetUrlLoading
                              ? isZh
                                ? "加载预览中..."
                                : "Loading preview..."
                              : isZh
                                ? "无法加载预览，保留 assetId 供追踪"
                                : "Preview unavailable, keep asset id for tracing"}
                          </div>
                        )}
                        <div className="admin-cert-actions">
                          {previewUrl ? (
                            <a href={previewUrl} target="_blank" rel="noreferrer">
                              {isZh ? "新窗口查看" : "Open in new tab"}
                            </a>
                          ) : (
                            <Text type="secondary">{isZh ? "请检查 media-svc 读链接权限" : "Check media read-url permission"}</Text>
                          )}
                        </div>
                      </article>
                    );
                  })}
                </div>
              )}
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">{isZh ? "审核信息" : "Review Notes"}</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label={isZh ? "审核备注" : "Review Comment"} span={2}>
                  <Paragraph className="admin-detail-paragraph">{displayValue(detailApplication.reviewComment)}</Paragraph>
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "驳回原因编码" : "Reject Reason Code"}>
                  {displayValue(detailApplication.latestReject?.rejectReasonCode)}
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "驳回时间" : "Rejected At"}>
                  {formatDateTime(detailApplication.latestReject?.rejectedAt, locale)}
                </Descriptions.Item>
                <Descriptions.Item label={isZh ? "驳回备注" : "Reject Comment"} span={2}>
                  <Paragraph className="admin-detail-paragraph">{displayValue(detailApplication.latestReject?.rejectComment)}</Paragraph>
                </Descriptions.Item>
              </Descriptions>
            </section>
          </div>
        )}
      </Modal>
    </AdminPage>
  );
}
