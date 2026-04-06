"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Checkbox,
  Col,
  Divider,
  Empty,
  Form,
  Input,
  InputNumber,
  Radio,
  Row,
  Space,
  Steps,
  Tag,
  Typography,
  notification
} from "antd";
import {
  AppstoreAddOutlined,
  MinusCircleOutlined,
  PictureOutlined,
  PlusOutlined,
  ReloadOutlined,
  RocketOutlined,
  SafetyCertificateOutlined,
  ShopOutlined,
  ShoppingOutlined,
  TagOutlined
} from "@ant-design/icons";
import { useRouter, useSearchParams } from "next/navigation";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import {
  CatalogDraftPayload,
  CatalogSkuDraft,
  ProductDetailResult,
  SaveDraftResult
} from "@/features/catalog/types";
import {
  fetchMyProduct,
  mapProductDetailToDraftPayload,
  saveProductDraft,
  submitProductReview
} from "@/features/catalog/api";
import { adjustInventory } from "@/features/inventory/api";
import { replaceMediaBindings } from "@/features/media/api";
import { MediaLibrary } from "@/features/media/MediaLibrary";
import { MediaAsset } from "@/features/media/types";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";
import { SellerWorkbenchResponse } from "@/features/seller-shop/types";

const { Title, Paragraph, Text } = Typography;

type PublishWorkspaceStore = {
  sessionBizNo?: string;
  lastSavedAt?: string;
  syncedStockBySku: Record<string, number>;
  ensureSession: () => void;
  markSavedAt: () => void;
  resetSyncedStock: () => void;
  markSyncedStock: (skuNo: string, qty: number) => void;
};

const usePublishWorkspaceStore = create<PublishWorkspaceStore>()(
  persist(
    (set, get) => ({
      syncedStockBySku: {},
      ensureSession: () => {
        if (get().sessionBizNo) {
          return;
        }
        const random =
          typeof crypto !== "undefined" && typeof crypto.randomUUID === "function"
            ? crypto.randomUUID()
            : `${Date.now()}-${Math.random().toString(36).slice(2)}`;
        set({ sessionBizNo: `seller-publish-${random}` });
      },
      markSavedAt: () => set({ lastSavedAt: new Date().toISOString() }),
      resetSyncedStock: () => set({ syncedStockBySku: {} }),
      markSyncedStock: (skuNo, qty) =>
        set((state) => ({
          syncedStockBySku: {
            ...state.syncedStockBySku,
            [skuNo]: qty
          }
        }))
    }),
    { name: "shopa-seller-publish-workspace" }
  )
);

const defaultDraftValues: CatalogDraftPayload = {
  draft: {
    title: "",
    summary: "",
    categoryId: 0,
    brandNo: "",
    mainImageAssetId: undefined,
    detailImageAssetIds: [],
    submitNote: "",
    attributeValues: {
      brand: "",
      material: "",
      origin: "",
      craft: "",
      fitScene: "",
      servicePromise: ""
    }
  },
  skus: [
    {
      name: "Default SKU",
      salePrice: 0,
      marketPrice: 0,
      saleSpecs: { color: "", size: "" },
      initialStock: 0
    }
  ]
};

function formatDateTime(value?: string) {
  if (!value) {
    return "Not saved yet";
  }
  return new Date(value).toLocaleString();
}

function formatCurrency(value?: number) {
  return `¥${Number(value ?? 0).toFixed(2)}`;
}

function collectUsedAssetIds(values: CatalogDraftPayload): string[] {
  const ids = new Set<string>();
  if (values.draft.mainImageAssetId) {
    ids.add(String(values.draft.mainImageAssetId));
  }
  for (const item of values.draft.detailImageAssetIds ?? []) {
    if (item) {
      ids.add(String(item));
    }
  }
  for (const sku of values.skus ?? []) {
    if (sku.skuImageAssetId) {
      ids.add(String(sku.skuImageAssetId));
    }
  }
  return Array.from(ids);
}

function getPriceRangeText(skus: CatalogSkuDraft[]) {
  const prices = skus.map((sku) => Number(sku.salePrice ?? 0)).filter((price) => price > 0);
  if (!prices.length) {
    return "Waiting for price";
  }
  const min = Math.min(...prices);
  const max = Math.max(...prices);
  return min === max ? formatCurrency(min) : `${formatCurrency(min)} - ${formatCurrency(max)}`;
}

function getMallHighlights(values: Record<string, string> | undefined) {
  const items = [
    { label: "Brand", value: values?.brand ?? "" },
    { label: "Material", value: values?.material ?? "" },
    { label: "Origin", value: values?.origin ?? "" },
    { label: "Craft", value: values?.craft ?? "" },
    { label: "Scene", value: values?.fitScene ?? "" },
    { label: "Promise", value: values?.servicePromise ?? "" }
  ];
  return items.filter((item) => item.value.trim());
}

function InventoryStatusHint({ messages }: { messages: string[] }) {
  if (!messages.length) {
    return (
      <Alert
        type="info"
        showIcon
        message="Inventory sync notes"
        description="Inventory sync is a separate action. The page writes to inventory-svc on a best-effort basis and tells you exactly what failed."
      />
    );
  }
  return (
    <Alert
      type={messages.some((item) => item.includes("failed") || item.includes("later")) ? "warning" : "success"}
      showIcon
      message="Inventory sync result"
      description={
        <ul style={{ margin: 0, paddingLeft: 18 }}>
          {messages.map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      }
    />
  );
}

export default function PublishPage() {
  const [form] = Form.useForm<CatalogDraftPayload>();
  const router = useRouter();
  const searchParams = useSearchParams();
  const querySpuNo = (searchParams.get("spuNo") || "").trim() || undefined;
  const {
    sessionBizNo,
    lastSavedAt,
    syncedStockBySku,
    ensureSession,
    markSavedAt,
    resetSyncedStock,
    markSyncedStock
  } = usePublishWorkspaceStore();
  const [currentSpuNo, setCurrentSpuNo] = useState<string | undefined>(querySpuNo);
  const [currentSpuVersion, setCurrentSpuVersion] = useState(0);
  const [currentReviewStatus, setCurrentReviewStatus] = useState("draft");
  const [mediaAssets, setMediaAssets] = useState<MediaAsset[]>([]);
  const [savingDraft, setSavingDraft] = useState(false);
  const [syncingInventory, setSyncingInventory] = useState(false);
  const [submittingReview, setSubmittingReview] = useState(false);
  const [inventoryMessages, setInventoryMessages] = useState<string[]>([]);
  const [hydratedSpuNo, setHydratedSpuNo] = useState<string | undefined>();

  useEffect(() => {
    ensureSession();
  }, [ensureSession]);

  useEffect(() => {
    setCurrentSpuNo(querySpuNo);
    setCurrentSpuVersion(0);
    setCurrentReviewStatus("draft");
    setInventoryMessages([]);
    resetSyncedStock();
    if (!querySpuNo) {
      form.setFieldsValue(defaultDraftValues);
      setHydratedSpuNo(undefined);
    }
  }, [form, querySpuNo, resetSyncedStock]);

  const workbenchQuery = useQuery<SellerWorkbenchResponse>({
    queryKey: ["seller-workbench", "publish"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30_000
  });

  const currentShop = useMemo(() => workbenchQuery.data?.shops?.[0], [workbenchQuery.data]);

  const detailQuery = useQuery<ProductDetailResult>({
    queryKey: ["seller-product", querySpuNo],
    queryFn: () => fetchMyProduct(querySpuNo as string),
    enabled: Boolean(querySpuNo),
    staleTime: 0
  });

  useEffect(() => {
    if (!detailQuery.data || !querySpuNo || hydratedSpuNo === querySpuNo) {
      return;
    }
    const mapped = mapProductDetailToDraftPayload(detailQuery.data);
    form.setFieldsValue(mapped);
    setCurrentSpuNo(detailQuery.data.product.spu.spuNo || querySpuNo);
    setCurrentSpuVersion(detailQuery.data.product.spu.version || 0);
    setCurrentReviewStatus(detailQuery.data.review?.reviewStatus || detailQuery.data.product.spu.status || "draft");
    setHydratedSpuNo(querySpuNo);
  }, [detailQuery.data, form, hydratedSpuNo, querySpuNo]);

  const watchedDraft =
    (Form.useWatch("draft", form) as CatalogDraftPayload["draft"] | undefined) ?? defaultDraftValues.draft;
  const watchedSkus = (Form.useWatch("skus", form) as CatalogSkuDraft[] | undefined) ?? defaultDraftValues.skus;
  const watchedMainImageAssetId = Form.useWatch(["draft", "mainImageAssetId"], form);
  const watchedDetailImageAssetIds = Form.useWatch(["draft", "detailImageAssetIds"], form) ?? [];

  const assetOptions = useMemo(
    () =>
      mediaAssets.map((asset) => ({
        label: asset.name,
        value: String(asset.assetId),
        preview: asset.publicUrl || asset.thumbnail
      })),
    [mediaAssets]
  );

  useEffect(() => {
    if (!watchedMainImageAssetId && mediaAssets.length === 1) {
      form.setFieldValue(["draft", "mainImageAssetId"], String(mediaAssets[0].assetId));
    }
  }, [form, mediaAssets, watchedMainImageAssetId]);

  const mediaBizNo = useMemo(() => {
    if (!sessionBizNo) {
      return "seller-publish-session";
    }
    if (mediaAssets.length > 0 || !currentSpuNo) {
      return sessionBizNo;
    }
    return currentSpuNo;
  }, [currentSpuNo, mediaAssets.length, sessionBizNo]);

  const mallHighlights = useMemo(() => getMallHighlights(watchedDraft.attributeValues), [watchedDraft.attributeValues]);

  const stepItems = useMemo(() => {
    const hasPrice = watchedSkus.some((sku) => Number(sku.salePrice ?? 0) > 0);
    const inventoryPending = watchedSkus.some((sku) => Number(sku.initialStock ?? 0) > 0);
    const inventorySynced = inventoryPending
      ? watchedSkus.every((sku) => {
          const skuNo = sku.skuNo ?? "";
          if (!skuNo && Number(sku.initialStock ?? 0) > 0) {
            return false;
          }
          if (!skuNo) {
            return true;
          }
          return Number(syncedStockBySku[skuNo] ?? -1) === Number(sku.initialStock ?? 0);
        })
      : true;
    return [
      {
        title: "Assets",
        description: assetOptions.length || watchedMainImageAssetId ? "Media ready for display" : "Upload hero and detail images",
        status: assetOptions.length || watchedMainImageAssetId ? "finish" : "process"
      },
      {
        title: "Draft",
        description: currentSpuNo ? `SPU ${currentSpuNo} saved` : hasPrice ? "Form is ready to save" : "Fill core product info",
        status: currentSpuNo ? "finish" : hasPrice ? "process" : "wait"
      },
      {
        title: "Inventory",
        description: inventoryPending ? (inventorySynced ? "Initial stock synced" : "Inventory still waiting") : "No stock sync in this round",
        status: inventoryPending ? (inventorySynced ? "finish" : "process") : "wait"
      },
      {
        title: "Review",
        description: currentReviewStatus === "reviewing" ? "Already in review queue" : "Submit only after final check",
        status: currentReviewStatus === "reviewing" ? "finish" : currentSpuNo ? "process" : "wait"
      }
    ] as const;
  }, [assetOptions.length, currentReviewStatus, currentSpuNo, syncedStockBySku, watchedMainImageAssetId, watchedSkus]);

  async function syncMediaBindingsForSpu(targetSpuNo: string, values: CatalogDraftPayload) {
    const assetIds = collectUsedAssetIds(values);
    if (!assetIds.length) {
      return;
    }
    await replaceMediaBindings({ bizNo: targetSpuNo, assetIds });
  }

  async function saveCurrentDraft() {
    if (!currentShop) {
      notification.warning({
        message: "No seller shop available",
        description: "Please finish seller onboarding first."
      });
      return null;
    }

    const values = (await form.validateFields()) as CatalogDraftPayload;
    if (!values.draft.mainImageAssetId) {
      notification.warning({
        message: "Main image is required",
        description: "Pick one main image from the media library."
      });
      return null;
    }

    setSavingDraft(true);
    try {
      const result = await saveProductDraft({
        shopNo: currentShop.shopNo,
        draft: values.draft,
        skus: values.skus,
        spuNo: currentSpuNo,
        expectedVersion: currentSpuVersion
      });

      try {
        await syncMediaBindingsForSpu(result.product.spu.spuNo, values);
      } catch (bindingError) {
        notification.warning({
          message: "Draft saved, media binding needs retry",
          description: bindingError instanceof Error ? bindingError.message : "Media binding failed"
        });
      }

      setCurrentSpuNo(result.product.spu.spuNo);
      setCurrentSpuVersion(result.product.spu.version);
      setCurrentReviewStatus(result.product.spu.status);
      markSavedAt();
      router.replace(`/seller/publish?spuNo=${encodeURIComponent(result.product.spu.spuNo)}`);

      form.setFieldsValue({
        draft: values.draft,
        skus: result.product.skus.map((sku, index) => ({
          skuNo: sku.skuNo,
          name: values.skus[index]?.name ?? sku.skuName,
          skuImageAssetId: values.skus[index]?.skuImageAssetId ?? sku.skuImageAssetId,
          salePrice: values.skus[index]?.salePrice ?? sku.salePrice,
          marketPrice: values.skus[index]?.marketPrice ?? sku.marketPrice,
          saleSpecs: values.skus[index]?.saleSpecs ?? sku.saleAttrs,
          initialStock: values.skus[index]?.initialStock ?? 0
        }))
      });

      notification.success({
        message: "Draft saved",
        description: `SPU ${result.product.spu.spuNo} has been written to catalog-svc.`
      });
      return result;
    } catch (error) {
      notification.error({
        message: "Save failed",
        description: error instanceof Error ? error.message : "Draft save failed"
      });
      return null;
    } finally {
      setSavingDraft(false);
    }
  }

  async function handleSyncInventory() {
    if (!currentShop) {
      return;
    }
    let saveResult: SaveDraftResult | null = null;
    const currentValues = (form.getFieldsValue(true) || defaultDraftValues) as CatalogDraftPayload;
    const hasPendingSkuNo = currentValues.skus.some((sku) => Number(sku.initialStock ?? 0) > 0 && !sku.skuNo);

    if (!currentSpuNo || hasPendingSkuNo) {
      saveResult = await saveCurrentDraft();
      if (!saveResult) {
        return;
      }
    }

    const latestValues = (form.getFieldsValue(true) || defaultDraftValues) as CatalogDraftPayload;
    const stockTargets = latestValues.skus.filter((sku) => Number(sku.initialStock ?? 0) > 0 && sku.skuNo);
    if (!stockTargets.length) {
      notification.info({
        message: "No stock to sync",
        description: "You can submit for review directly if initial stock is not needed."
      });
      return;
    }

    setSyncingInventory(true);
    const messages: string[] = [];
    const effectiveSpuNo = currentSpuNo ?? saveResult?.product.spu.spuNo ?? "draft";

    for (const sku of stockTargets) {
      const desiredQty = Number(sku.initialStock ?? 0);
      const alreadySyncedQty = Number(syncedStockBySku[sku.skuNo as string] ?? 0);
      const delta = desiredQty - alreadySyncedQty;
      if (delta === 0) {
        messages.push(`${sku.name}: already synced to ${desiredQty}`);
        continue;
      }

      const response = await adjustInventory({
        shopNo: currentShop.shopNo,
        skuNo: sku.skuNo as string,
        spuNo: effectiveSpuNo,
        delta,
        bizNo: `publish:${effectiveSpuNo}:${sku.skuNo}`,
        reasonCode: "SELLER_REPLENISH",
        remark: "Initial stock sync from seller publish page"
      });

      if (response.success) {
        markSyncedStock(sku.skuNo as string, desiredQty);
        messages.push(`${sku.name}: synced ${desiredQty}`);
      } else {
        const reason = response.failedItems[0]?.reason ?? "inventory service unavailable";
        messages.push(`${sku.name}: failed, retry later (${reason})`);
      }
    }

    setInventoryMessages(messages);
    setSyncingInventory(false);
  }

  async function handleSubmitReview() {
    if (!currentShop) {
      return;
    }
    const saveResult = await saveCurrentDraft();
    if (!saveResult) {
      return;
    }
    const values = (form.getFieldsValue(true) || defaultDraftValues) as CatalogDraftPayload;
    setSubmittingReview(true);
    try {
      await submitProductReview({
        shopNo: currentShop.shopNo,
        spuNo: saveResult.product.spu.spuNo,
        expectedVersion: saveResult.product.spu.version,
        submitNote: values.draft.submitNote
      });
      setCurrentReviewStatus("reviewing");
      notification.success({
        message: "Review submitted",
        description: `SPU ${saveResult.product.spu.spuNo} is now in review queue.`
      });
    } catch (error) {
      notification.error({
        message: "Review submit failed",
        description: error instanceof Error ? error.message : "Submit review failed"
      });
    } finally {
      setSubmittingReview(false);
    }
  }

  if (workbenchQuery.isLoading) {
    return <Card loading style={{ minHeight: 540 }} />;
  }

  if (!currentShop) {
    return (
      <Card bordered={false} style={{ borderRadius: 28, border: "1px solid #ffe5cf", background: "linear-gradient(180deg, #fff9f2, #ffffff)" }}>
        <Empty
          description="This account does not have a seller shop that can publish products yet."
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        >
          <Space>
            <Link href="/seller/workbench">
              <Button type="primary">Open Seller Workbench</Button>
            </Link>
            <Link href="/seller/profile">
              <Button>Open Seller Profile</Button>
            </Link>
          </Space>
        </Empty>
      </Card>
    );
  }

  return (
    <div style={{ display: "grid", gap: 24 }}>
      <section
        style={{
          borderRadius: 28,
          padding: 28,
          color: "#4b2300",
          background: "radial-gradient(circle at right top, rgba(255,255,255,0.74), transparent 32%), linear-gradient(135deg, #fff0de 0%, #ffd5ae 45%, #fff6ef 100%)",
          border: "1px solid rgba(255, 146, 73, 0.24)"
        }}
      >
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} xl={16}>
            <Space direction="vertical" size={10}>
              <Tag color="orange" style={{ width: "fit-content", borderRadius: 999, padding: "4px 12px" }}>
                Seller Product Studio
              </Tag>
              <Title level={2} style={{ margin: 0 }}>Ship a real seller publishing workflow</Title>
              <Paragraph style={{ margin: 0, maxWidth: 780, fontSize: 15, lineHeight: 1.85, color: "#7f4413" }}>
                Assets go to media-svc, drafts go to catalog-svc, and initial stock syncs to inventory-svc. This page is wired to real services and shows honest failures instead of fake success.
              </Paragraph>
              <Space wrap>
                <Tag icon={<ShopOutlined />} color="gold">Shop: {currentShop.shopDisplayName || currentShop.shopName}</Tag>
                <Tag color="blue">Session: {sessionBizNo}</Tag>
                {currentSpuNo ? <Tag color="green">Editing SPU: {currentSpuNo}</Tag> : <Tag>New product draft</Tag>}
                <Tag color={currentReviewStatus === "reviewing" ? "processing" : "default"}>Status: {currentReviewStatus}</Tag>
              </Space>
            </Space>
          </Col>
          <Col xs={24} xl={8}>
            <Card bordered={false} style={{ borderRadius: 24, background: "rgba(255,255,255,0.82)", boxShadow: "0 22px 60px rgba(114, 54, 8, 0.12)" }}>
              <Space direction="vertical" size={8} style={{ width: "100%" }}>
                <Text strong>Current workspace</Text>
                <Text type="secondary">Last save: {formatDateTime(lastSavedAt)}</Text>
                <Text type="secondary">Price band: {getPriceRangeText(watchedSkus)}</Text>
                <Text type="secondary">Draft mode: {currentSpuNo ? "Edit existing draft" : "Create new draft"}</Text>
              </Space>
            </Card>
          </Col>
        </Row>
      </section>

      {detailQuery.isError ? (
        <Alert type="warning" showIcon message="Draft load failed" description={detailQuery.error instanceof Error ? detailQuery.error.message : "Unable to load current draft"} />
      ) : null}

      <Form form={form} layout="vertical" initialValues={defaultDraftValues}>
        <Row gutter={[24, 24]} align="top">
          <Col xs={24} xl={17}>
            <div style={{ display: "grid", gap: 24 }}>
              <MediaLibrary
                bizNo={mediaBizNo}
                title="1. Product media vault"
                description="Upload hero images, detail images, and SKU images first. The media area shows real uploaded assets for this publishing session or SPU."
                onAssetsChange={setMediaAssets}
              />

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3", boxShadow: "0 20px 45px rgba(33, 22, 14, 0.05)" }} title="2. Product identity">
                <Row gutter={16}>
                  <Col xs={24} md={16}>
                    <Form.Item name={["draft", "title"]} label="Buyer-facing title" rules={[{ required: true, message: "Title is required" }]}>
                      <Input placeholder="Example: Camellia silk dress 2026 spring new arrival" size="large" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={8}>
                    <Form.Item name={["draft", "categoryId"]} label="Category ID" rules={[{ required: true, message: "Category ID is required" }]}>
                      <InputNumber min={1} precision={0} style={{ width: "100%" }} size="large" placeholder="Real category ID" />
                    </Form.Item>
                  </Col>
                </Row>
                <Form.Item name={["draft", "summary"]} label="Sub title / selling summary" rules={[{ required: true, message: "Summary is required" }]}>
                  <Input.TextArea rows={4} placeholder="Write the material, fit, scene, and the first three selling points buyers should notice." />
                </Form.Item>
                <Row gutter={16}>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "brandNo"]} label="Brand code" rules={[{ required: true, message: "Brand code is required" }]}>
                      <Input placeholder="Example: BRAND_CAMELLIA" size="large" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "attributeValues", "brand"]} label="Brand display name">
                      <Input placeholder="Example: Camellia Atelier" size="large" />
                    </Form.Item>
                  </Col>
                </Row>
              </Card>

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3", boxShadow: "0 20px 45px rgba(33, 22, 14, 0.05)" }} title="3. Mall-facing imagery">
                {!assetOptions.length ? (
                  <Alert type="info" showIcon message="No uploaded assets yet" description="Upload images above first. The hero image, gallery, and SKU image selectors will show them here." />
                ) : (
                  <Space direction="vertical" size={20} style={{ width: "100%" }}>
                    <Form.Item name={["draft", "mainImageAssetId"]} label="Hero image" rules={[{ required: true, message: "Hero image is required" }]}>
                      <Radio.Group style={{ width: "100%" }}>
                        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(170px, 1fr))", gap: 12 }}>
                          {assetOptions.map((asset) => (
                            <label key={asset.value} style={{ display: "grid", gap: 10, padding: 12, borderRadius: 18, cursor: "pointer", border: watchedMainImageAssetId === asset.value ? "2px solid #ff7a1a" : "1px solid #f0d9c3", background: watchedMainImageAssetId === asset.value ? "#fff4ea" : "#fff" }}>
                              <Radio value={asset.value}>{asset.label}</Radio>
                              <img src={asset.preview} alt={asset.label} style={{ width: "100%", aspectRatio: "1 / 1", objectFit: "cover", borderRadius: 14 }} />
                            </label>
                          ))}
                        </div>
                      </Radio.Group>
                    </Form.Item>

                    <Form.Item name={["draft", "detailImageAssetIds"]} label="Gallery / detail images">
                      <Checkbox.Group style={{ width: "100%" }}>
                        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(170px, 1fr))", gap: 12 }}>
                          {assetOptions.map((asset) => (
                            <label key={`detail-${asset.value}`} style={{ display: "grid", gap: 10, padding: 12, borderRadius: 18, cursor: "pointer", border: watchedDetailImageAssetIds.includes(asset.value) ? "2px solid #ff9a43" : "1px solid #f0d9c3", background: watchedDetailImageAssetIds.includes(asset.value) ? "#fff7ef" : "#fff" }}>
                              <Checkbox value={asset.value}>{asset.label}</Checkbox>
                              <img src={asset.preview} alt={asset.label} style={{ width: "100%", aspectRatio: "1 / 1", objectFit: "cover", borderRadius: 14 }} />
                            </label>
                          ))}
                        </div>
                      </Checkbox.Group>
                    </Form.Item>
                  </Space>
                )}
              </Card>

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3", boxShadow: "0 20px 45px rgba(33, 22, 14, 0.05)" }} title="4. Buyer-facing highlights">
                <Row gutter={16}>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "attributeValues", "material"]} label="Material">
                      <Input placeholder="Example: 93% silk / 7% spandex" size="large" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "attributeValues", "origin"]} label="Origin">
                      <Input placeholder="Example: Hangzhou" size="large" />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={16}>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "attributeValues", "craft"]} label="Craft">
                      <Input placeholder="Example: silk twill, pleated finish" size="large" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name={["draft", "attributeValues", "fitScene"]} label="Scene">
                      <Input placeholder="Example: commute / date / travel" size="large" />
                    </Form.Item>
                  </Col>
                </Row>
                <Form.Item name={["draft", "attributeValues", "servicePromise"]} label="Service promise">
                  <Input.TextArea rows={3} placeholder="Example: ship within 72 hours, 7-day return, damage protection." />
                </Form.Item>
                <Form.Item name={["draft", "submitNote"]} label="Review note">
                  <Input.TextArea rows={3} placeholder="Optional internal note for review, such as material proof or copyright explanation." />
                </Form.Item>
              </Card>

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3", boxShadow: "0 20px 45px rgba(33, 22, 14, 0.05)" }} title="5. SKU pricing and opening stock">
                <Form.List name="skus">
                  {(fields, { add, remove }) => (
                    <Space direction="vertical" size={18} style={{ width: "100%" }}>
                      {fields.map((field, index) => (
                        <Card key={field.key} size="small" title={`SKU #${index + 1}`} extra={fields.length > 1 ? <Button danger type="text" icon={<MinusCircleOutlined />} onClick={() => remove(field.name)}>Remove</Button> : null} style={{ borderRadius: 20, background: "#fffaf6" }}>
                          <Row gutter={16}>
                            <Col xs={24} md={10}>
                              <Form.Item {...field} name={[field.name, "name"]} label="SKU label" rules={[{ required: true, message: "SKU label is required" }]}>
                                <Input placeholder="Example: Cream / S" size="large" />
                              </Form.Item>
                            </Col>
                            <Col xs={12} md={7}>
                              <Form.Item {...field} name={[field.name, "saleSpecs", "color"]} label="Color">
                                <Input placeholder="Color" size="large" />
                              </Form.Item>
                            </Col>
                            <Col xs={12} md={7}>
                              <Form.Item {...field} name={[field.name, "saleSpecs", "size"]} label="Size">
                                <Input placeholder="Size" size="large" />
                              </Form.Item>
                            </Col>
                          </Row>
                          <Row gutter={16}>
                            <Col xs={24} md={6}>
                              <Form.Item {...field} name={[field.name, "salePrice"]} label="Sale price" rules={[{ required: true, message: "Sale price is required" }]}>
                                <InputNumber min={0.01} precision={2} style={{ width: "100%" }} size="large" />
                              </Form.Item>
                            </Col>
                            <Col xs={24} md={6}>
                              <Form.Item {...field} name={[field.name, "marketPrice"]} label="Market price" rules={[{ required: true, message: "Market price is required" }]}>
                                <InputNumber min={0.01} precision={2} style={{ width: "100%" }} size="large" />
                              </Form.Item>
                            </Col>
                            <Col xs={24} md={6}>
                              <Form.Item {...field} name={[field.name, "initialStock"]} label="Initial stock">
                                <InputNumber min={0} precision={0} style={{ width: "100%" }} size="large" />
                              </Form.Item>
                            </Col>
                            <Col xs={24} md={6}>
                              <Form.Item {...field} name={[field.name, "skuImageAssetId"]} label="SKU image">
                                <Radio.Group style={{ width: "100%" }}>
                                  <Space direction="vertical" style={{ width: "100%" }}>
                                    <Radio value={undefined}>Use hero image</Radio>
                                    {assetOptions.slice(0, 6).map((asset) => (
                                      <Radio key={`sku-image-${field.key}-${asset.value}`} value={asset.value}>
                                        {asset.label}
                                      </Radio>
                                    ))}
                                  </Space>
                                </Radio.Group>
                              </Form.Item>
                            </Col>
                          </Row>
                          <Form.Item {...field} name={[field.name, "skuNo"]} hidden>
                            <Input />
                          </Form.Item>
                        </Card>
                      ))}
                      <Button
                        type="dashed"
                        icon={<PlusOutlined />}
                        onClick={() => add({ name: "New SKU", salePrice: 0, marketPrice: 0, saleSpecs: { color: "", size: "" }, initialStock: 0 })}
                        style={{ height: 48, borderRadius: 16 }}
                      >
                        Add SKU
                      </Button>
                    </Space>
                  )}
                </Form.List>
                <Divider />
                <InventoryStatusHint messages={inventoryMessages} />
              </Card>
            </div>
          </Col>

          <Col xs={24} xl={7}>
            <div style={{ position: "sticky", top: 24, display: "grid", gap: 20 }}>
              <Card bordered={false} style={{ borderRadius: 24, overflow: "hidden", background: "linear-gradient(180deg, rgba(255,118,41,0.12) 0%, rgba(255,255,255,1) 34%)", border: "1px solid #ffd9bf", boxShadow: "0 24px 60px rgba(82, 42, 12, 0.08)" }}>
                <Space direction="vertical" size={18} style={{ width: "100%" }}>
                  <div>
                    <Text strong style={{ fontSize: 16 }}>Launch cockpit</Text>
                    <Paragraph type="secondary" style={{ margin: "8px 0 0" }}>
                      This side panel reflects the real state of the publish flow and no longer hides integration failures.
                    </Paragraph>
                  </div>
                  <Steps direction="vertical" size="small" current={Math.max(stepItems.findIndex((item) => item.status === "process"), 0)} items={stepItems.map((item) => ({ title: item.title, description: item.description, status: item.status as never }))} />
                  <Alert
                    type={currentReviewStatus === "reviewing" ? "success" : "info"}
                    showIcon
                    message={currentReviewStatus === "reviewing" ? "Already in review" : "Not submitted yet"}
                    description={currentReviewStatus === "reviewing" ? "You can go back to product list and wait for the result." : "Save the draft first, then confirm inventory and image mapping."}
                  />
                  <Button type="primary" size="large" icon={<SafetyCertificateOutlined />} loading={savingDraft} onClick={() => void saveCurrentDraft()}>
                    Save draft
                  </Button>
                  <Button size="large" icon={<ReloadOutlined />} loading={syncingInventory} onClick={() => void handleSyncInventory()}>
                    Sync initial stock
                  </Button>
                  <Button size="large" icon={<RocketOutlined />} loading={submittingReview} onClick={() => void handleSubmitReview()} disabled={currentReviewStatus === "reviewing"}>
                    Submit review
                  </Button>
                </Space>
              </Card>

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3" }}>
                <Space direction="vertical" size={12} style={{ width: "100%" }}>
                  <Space align="center">
                    <ShoppingOutlined style={{ color: "#ff6a00" }} />
                    <Text strong>Mall card preview</Text>
                  </Space>
                  <div style={{ borderRadius: 20, overflow: "hidden", border: "1px solid #f0dac5", background: "#fffaf5" }}>
                    <div style={{ aspectRatio: "1 / 1", background: "linear-gradient(135deg, #ffe9d7, #fff6ed)", display: "grid", placeItems: "center" }}>
                      {assetOptions.find((item) => item.value === watchedMainImageAssetId)?.preview ? (
                        <img src={assetOptions.find((item) => item.value === watchedMainImageAssetId)?.preview} alt="product hero preview" style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                      ) : (
                        <Space direction="vertical" align="center">
                          <PictureOutlined style={{ fontSize: 28, color: "#ff8f3d" }} />
                          <Text type="secondary">Waiting for hero image</Text>
                        </Space>
                      )}
                    </div>
                    <div style={{ padding: 16, display: "grid", gap: 10 }}>
                      <Text strong style={{ fontSize: 15, lineHeight: 1.6 }}>{watchedDraft.title || "Buyer-facing title will appear here"}</Text>
                      <Text type="secondary" style={{ lineHeight: 1.7 }}>{watchedDraft.summary || "The summary shown here maps directly to what buyers feel first on search and PDP."}</Text>
                      <Text style={{ color: "#ff5000", fontWeight: 700, fontSize: 22 }}>{getPriceRangeText(watchedSkus)}</Text>
                      <Space wrap>
                        {mallHighlights.length ? mallHighlights.slice(0, 4).map((item) => (
                          <Tag key={item.label} color="orange">
                            {item.label}: {item.value}
                          </Tag>
                        )) : <Tag>Add material, origin, and promise to make this look like a real mall card</Tag>}
                      </Space>
                    </div>
                  </div>
                </Space>
              </Card>

              <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #ffe7d3" }}>
                <Space direction="vertical" size={10} style={{ width: "100%" }}>
                  <Space align="center">
                    <AppstoreAddOutlined style={{ color: "#1677ff" }} />
                    <Text strong>Final checklist</Text>
                  </Space>
                  <Text type="secondary">Hero image: {watchedMainImageAssetId ? "ready" : "missing"}</Text>
                  <Text type="secondary">Detail images: {(watchedDetailImageAssetIds ?? []).length}</Text>
                  <Text type="secondary">SKU count: {watchedSkus.length}</Text>
                  <Text type="secondary">Last save: {formatDateTime(lastSavedAt)}</Text>
                  <Text type="secondary">Current version: {currentSpuVersion}</Text>
                  <Text type="secondary">Price band: {getPriceRangeText(watchedSkus)}</Text>
                  <Divider style={{ margin: "6px 0" }} />
                  <Link href="/seller/products">
                    <Button block icon={<TagOutlined />}>View my products</Button>
                  </Link>
                </Space>
              </Card>
            </div>
          </Col>
        </Row>
      </Form>
    </div>
  );
}
