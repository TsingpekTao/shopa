"use client";

import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Button, Descriptions, Empty, Image, Input, Modal, Select, Space, Spin, Table, Tag, Typography, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useQuery } from "@tanstack/react-query";
import { AccessControl } from "@shopa/ui";
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

const STATUS_OPTIONS = [
  { value: 1, label: "草稿" },
  { value: 2, label: "已提交" },
  { value: 3, label: "审核中" },
  { value: 4, label: "已通过" },
  { value: 5, label: "已驳回" },
  { value: 6, label: "已取消" },
  { value: 7, label: "系统驳回" }
];

const REJECT_REASON_OPTIONS = [
  { value: "MATERIAL_INVALID", label: "资料无效" },
  { value: "MATERIAL_INCOMPLETE", label: "资料不完整" },
  { value: "IDENTITY_MISMATCH", label: "身份信息不匹配" },
  { value: "RISK_CONTROL", label: "风控拦截" }
];

type CredentialAsset = {
  key: string;
  label: string;
  assetId: string;
};

type RejectTarget = {
  applicationNo: string;
  expectedVersion: number;
};

function statusMeta(status?: number) {
  switch (status) {
    case 1:
      return { color: "default", label: "草稿" };
    case 2:
      return { color: "processing", label: "已提交" };
    case 3:
      return { color: "gold", label: "审核中" };
    case 4:
      return { color: "success", label: "已通过" };
    case 5:
      return { color: "error", label: "已驳回" };
    case 6:
      return { color: "default", label: "已取消" };
    case 7:
      return { color: "volcano", label: "系统驳回" };
    default:
      return { color: "default", label: "未知" };
  }
}

function renderStatus(status?: number) {
  const meta = statusMeta(status);
  return <Tag color={meta.color}>{meta.label}</Tag>;
}

function displayValue(value: string | number | undefined | null): string {
  if (value === null || value === undefined) {
    return "-";
  }
  const text = String(value).trim();
  return text.length > 0 ? text : "-";
}

function formatDateTime(value?: string): string {
  if (!value) {
    return "-";
  }
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return value;
  }
  return new Intl.DateTimeFormat("zh-CN", {
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
  return /^\d+$/.test(normalized) ? normalized : undefined;
}

function docTypeLabel(docTypeCode?: string) {
  const normalized = (docTypeCode ?? "").toUpperCase();
  if (normalized.includes("ID_CARD_FRONT")) {
    return "身份证人像面";
  }
  if (normalized.includes("ID_CARD_BACK")) {
    return "身份证国徽面";
  }
  if (normalized.includes("BUSINESS_LICENSE")) {
    return "营业执照";
  }
  return normalized || "资质文件";
}

function renderExtList(ext?: Record<string, string>) {
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

function findExtAssetId(ext: Record<string, string> | undefined, keys: string[]) {
  if (!ext) {
    return undefined;
  }
  for (const key of keys) {
    const assetId = toAssetId(ext[key]);
    if (assetId) {
      return assetId;
    }
  }
  return undefined;
}

function extractCredentialAssets(entity?: MerchantEntityProfile): CredentialAsset[] {
  if (!entity) {
    return [];
  }
  const items: CredentialAsset[] = [];
  const ext = entity.ext;
  const extDefinitions = [
    { key: "id_card_front", label: "身份证人像面", keys: ["id_card_front_asset_id", "idCardFrontAssetId"] },
    { key: "id_card_back", label: "身份证国徽面", keys: ["id_card_back_asset_id", "idCardBackAssetId"] },
    { key: "business_license", label: "营业执照", keys: ["business_license_asset_id", "businessLicenseAssetId"] }
  ];

  extDefinitions.forEach((definition) => {
    const assetId = findExtAssetId(ext, definition.keys);
    if (assetId) {
      items.push({ key: definition.key, label: definition.label, assetId });
    }
  });

  const legalDocs: QualificationDoc[] = entity.legalSubject?.docs ?? [];
  legalDocs.forEach((doc, index) => {
    const assetId = toAssetId(doc.assetId);
    if (!assetId) {
      return;
    }
    items.push({
      key: `doc_${index}`,
      label: docTypeLabel(doc.docTypeCode),
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

export default function MerchantReviewPage() {
  const searchParams = useSearchParams();
  const initialStatuses = useMemo(
    () =>
      searchParams
        .getAll("statuses")
        .map((value) => Number(value))
        .filter((value) => Number.isFinite(value) && value > 0),
    [searchParams]
  );
  const initialKeyword = searchParams.get("keyword")?.trim() ?? "";
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [keywordInput, setKeywordInput] = useState(initialKeyword);
  const [keyword, setKeyword] = useState(initialKeyword);
  const [statusFilters, setStatusFilters] = useState<number[]>(initialStatuses);
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
  const credentialAssets = useMemo(() => extractCredentialAssets(detailEntity), [detailEntity]);

  useEffect(() => {
    const nextStatuses = searchParams
      .getAll("statuses")
      .map((value) => Number(value))
      .filter((value) => Number.isFinite(value) && value > 0);
    const nextKeyword = searchParams.get("keyword")?.trim() ?? "";
    setStatusFilters(nextStatuses);
    setKeywordInput(nextKeyword);
    setKeyword(nextKeyword);
    setPage(1);
  }, [searchParams]);

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

  const submitReject = async () => {
    if (!rejectTarget?.applicationNo) {
      message.error("申请单号为空，无法驳回");
      return;
    }
    if (!rejectComment.trim()) {
      message.warning("请填写驳回说明");
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
      message.success("已完成驳回");
      setRejectOpen(false);
      setRejectTarget(null);
      setRejectComment("");
      await query.refetch();
    } finally {
      setRejectSubmitting(false);
      setLoadingNo("");
    }
  };

  const columns: ColumnsType<MerchantApplication> = [
    {
      title: "申请单号",
      dataIndex: "applicationNo",
      key: "applicationNo",
      width: 230
    },
    {
      title: "主体名称",
      dataIndex: ["entity", "entityName"],
      key: "entityName",
      width: 180
    },
    {
      title: "联系人",
      dataIndex: ["entity", "contactName"],
      key: "contactName",
      width: 140
    },
    {
      title: "联系电话",
      dataIndex: ["entity", "contactPhone"],
      key: "contactPhone",
      width: 160
    },
    {
      title: "状态",
      key: "status",
      width: 120,
      render: (_, row) => renderStatus(row.status)
    },
    {
      title: "操作",
      key: "actions",
      width: 360,
      render: (_, row) => (
        <Space>
          <Button
            size="small"
            onClick={async () => {
              if (!row.applicationNo) {
                message.error("申请单号为空，无法查看详情");
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
                message.error(`获取详情失败：${err?.message || "未知错误"}`);
              } finally {
                setDetailLoading(false);
              }
            }}
          >
            查看详情
          </Button>

          <AccessControl require={Permissions.MerchantReviewApprove}>
            <Button
              type="primary"
              size="small"
              loading={loadingNo === `${row.applicationNo}-approve`}
              onClick={async () => {
                try {
                  setLoadingNo(`${row.applicationNo}-approve`);
                  await approveMerchantApplication({
                    applicationNo: row.applicationNo,
                    expectedVersion: row.version ?? 1
                  });
                  message.success("已提交通过操作");
                  await query.refetch();
                } finally {
                  setLoadingNo("");
                }
              }}
            >
              通过
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
              驳回
            </Button>
          </AccessControl>
        </Space>
      )
    }
  ];

  return (
    <AdminPage title="商家入驻审核" subtitle="审核商家资质材料，执行通过或驳回操作。">
      <Space wrap size={12} style={{ marginBottom: 12 }}>
        <Input
          allowClear
          value={keywordInput}
          onChange={(event) => setKeywordInput(event.target.value)}
          onPressEnter={applyFilters}
          style={{ width: 320 }}
          placeholder="搜索主体名称或店铺名称"
        />
        <Select
          mode="multiple"
          allowClear
          maxTagCount="responsive"
          style={{ width: 320 }}
          value={statusFilters}
          options={STATUS_OPTIONS}
          onChange={(values) => {
            setPage(1);
            setStatusFilters(values);
          }}
          placeholder="按状态筛选"
        />
        <Button type="primary" onClick={applyFilters}>
          搜索
        </Button>
        <Button onClick={resetFilters}>重置</Button>
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
        title="填写驳回原因"
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
            <Text type="secondary">申请单号</Text>
            <div>{rejectTarget?.applicationNo ?? "-"}</div>
          </div>
          <div>
            <Text type="secondary">驳回原因</Text>
            <Select
              value={rejectReasonCode}
              style={{ width: "100%", marginTop: 6 }}
              options={REJECT_REASON_OPTIONS}
              onChange={(value) => setRejectReasonCode(value)}
            />
          </div>
          <div>
            <Text type="secondary">驳回说明</Text>
            <TextArea
              value={rejectComment}
              onChange={(event) => setRejectComment(event.target.value)}
              showCount
              maxLength={300}
              rows={5}
              style={{ marginTop: 6 }}
              placeholder="请填写驳回原因，商家会看到这段说明"
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
        title="商家申请详情"
        width={1040}
        footer={
          <Button
            onClick={() => {
              setDetailOpen(false);
              setDetailApplication(null);
            }}
          >
            关闭
          </Button>
        }
        destroyOnClose
      >
        {detailLoading ? (
          <div className="admin-detail-loading">
            <Spin />
          </div>
        ) : !detailApplication ? (
          <Empty description="暂无详情数据" />
        ) : (
          <div className="admin-merchant-detail">
            <section className="admin-detail-block">
              <h4 className="admin-detail-title">申请信息</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label="申请单号">{displayValue(detailApplication.applicationNo)}</Descriptions.Item>
                <Descriptions.Item label="状态">{renderStatus(detailApplication.status)}</Descriptions.Item>
                <Descriptions.Item label="版本号">{displayValue(detailApplication.version)}</Descriptions.Item>
                <Descriptions.Item label="申请人用户 ID">{displayValue(detailApplication.ownerUserId)}</Descriptions.Item>
                <Descriptions.Item label="审核人用户 ID">{displayValue(detailApplication.reviewerId)}</Descriptions.Item>
                <Descriptions.Item label="创建时间">{formatDateTime(detailApplication.createdAt)}</Descriptions.Item>
                <Descriptions.Item label="提交时间">{formatDateTime(detailApplication.submittedAt)}</Descriptions.Item>
                <Descriptions.Item label="开始审核时间">{formatDateTime(detailApplication.reviewStartedAt)}</Descriptions.Item>
                <Descriptions.Item label="审核完成时间">{formatDateTime(detailApplication.reviewedAt)}</Descriptions.Item>
                <Descriptions.Item label="更新时间">{formatDateTime(detailApplication.updatedAt)}</Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">主体信息</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label="主体编号">{displayValue(detailEntity?.entityNo)}</Descriptions.Item>
                <Descriptions.Item label="商户类型">{displayValue(detailEntity?.merchantTypeCode)}</Descriptions.Item>
                <Descriptions.Item label="主体名称">{displayValue(detailEntity?.entityName)}</Descriptions.Item>
                <Descriptions.Item label="联系人">{displayValue(detailEntity?.contactName)}</Descriptions.Item>
                <Descriptions.Item label="联系电话">{displayValue(detailEntity?.contactPhone)}</Descriptions.Item>
                <Descriptions.Item label="联系邮箱">{displayValue(detailEntity?.contactEmail)}</Descriptions.Item>
                <Descriptions.Item label="主体扩展信息" span={2}>
                  {renderExtList(detailEntity?.ext)}
                </Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">店铺信息</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label="店铺编号">{displayValue(detailShop?.shopNo)}</Descriptions.Item>
                <Descriptions.Item label="店铺类型">{displayValue(detailShop?.shopTypeCode)}</Descriptions.Item>
                <Descriptions.Item label="店铺名称">{displayValue(detailShop?.shopName)}</Descriptions.Item>
                <Descriptions.Item label="展示名称">{displayValue(detailShop?.shopDisplayName)}</Descriptions.Item>
                <Descriptions.Item label="客服电话">{displayValue(detailShop?.servicePhone)}</Descriptions.Item>
                <Descriptions.Item label="客服邮箱">{displayValue(detailShop?.serviceEmail)}</Descriptions.Item>
                <Descriptions.Item label="主营类目 ID" span={2}>
                  {displayValue(detailShop?.mainCategoryIds?.join(", "))}
                </Descriptions.Item>
                <Descriptions.Item label="店铺描述" span={2}>
                  <Paragraph className="admin-detail-paragraph">{displayValue(detailShop?.description)}</Paragraph>
                </Descriptions.Item>
                <Descriptions.Item label="店铺扩展信息" span={2}>
                  {renderExtList(detailShop?.ext)}
                </Descriptions.Item>
              </Descriptions>
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">资质材料</h4>
              {credentialAssets.length === 0 ? (
                <Text type="secondary">未提交可识别的资质文件</Text>
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
                            {assetUrlLoading ? "加载预览中..." : "暂时无法加载预览，请检查 media 服务读链路"}
                          </div>
                        )}
                        <div className="admin-cert-actions">
                          {previewUrl ? (
                            <a href={previewUrl} target="_blank" rel="noreferrer">
                              新窗口查看
                            </a>
                          ) : (
                            <Text type="secondary">保留 assetId 以便继续排查</Text>
                          )}
                        </div>
                      </article>
                    );
                  })}
                </div>
              )}
            </section>

            <section className="admin-detail-block">
              <h4 className="admin-detail-title">审核信息</h4>
              <Descriptions bordered size="small" column={2}>
                <Descriptions.Item label="审核备注" span={2}>
                  <Paragraph className="admin-detail-paragraph">{displayValue(detailApplication.reviewComment)}</Paragraph>
                </Descriptions.Item>
                <Descriptions.Item label="驳回原因编码">
                  {displayValue(detailApplication.latestReject?.rejectReasonCode)}
                </Descriptions.Item>
                <Descriptions.Item label="驳回时间">
                  {formatDateTime(detailApplication.latestReject?.rejectedAt)}
                </Descriptions.Item>
                <Descriptions.Item label="驳回说明" span={2}>
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
