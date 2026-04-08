"use client";

import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Alert, Button, Card, Empty, Input, InputNumber, Modal, Select, Space, Spin, Steps, Tag, Typography, notification } from "antd";
import { CheckCircleFilled, DeleteOutlined, ExclamationCircleFilled, FolderOpenOutlined, LoadingOutlined, PlusOutlined, UploadOutlined } from "@ant-design/icons";
import { useRouter, useSearchParams } from "next/navigation";
import { ProductDetailResult, SaveDraftResult } from "@/features/catalog/types";
import { fetchMyProduct, saveProductDraft, submitProductReview } from "@/features/catalog/api";
import { adjustInventory, batchGetSkuInventory } from "@/features/inventory/api";
import { fetchMediaLibrary, replaceMediaBindings, uploadSellerMediaAsset } from "@/features/media/api";
import { MediaAsset } from "@/features/media/types";
import {
  fetchSellerWorkbench,
  getSellerProductStoreCategoryBinding,
  listSellerStoreCategories,
  updateSellerProductStoreCategoryBinding
} from "@/features/seller-shop/api";
import { SellerStoreCategory, SellerWorkbenchResponse } from "@/features/seller-shop/types";

const { Title, Paragraph, Text } = Typography;

type Row = {
  key: string;
  size: string;
  skuNo?: string;
  stockQty: number;
  syncedQty: number;
  syncStatus: "idle" | "syncing" | "success" | "error";
  syncMessage?: string;
};

function rid() {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") return crypto.randomUUID();
  return `${Date.now()}-${Math.random().toString(36).slice(2)}`;
}

function emptyRow(size = ""): Row {
  return { key: rid(), size, stockQty: 0, syncedQty: 0, syncStatus: "idle" };
}

function trimSize(v: string) {
  return v.trim();
}

function uniqRows(rows: Row[]) {
  const used = new Set<string>();
  return rows.filter((row) => {
    const size = trimSize(row.size);
    if (!size) return true;
    if (used.has(size)) return false;
    used.add(size);
    return true;
  });
}

function money(v?: number | null) {
  return `¥${Number(v ?? 0).toFixed(2)}`;
}

function statusText(v: string) {
  if (v === "reviewing") return "审核中";
  if (v === "approved") return "已通过";
  if (v === "onShelf") return "已上架";
  if (v === "offShelf") return "已下架";
  return "草稿中";
}

function assetPreviewSrc(asset?: MediaAsset) {
  return asset?.publicUrl || asset?.thumbnail || "";
}

function flattenLeafCategories(categories: SellerStoreCategory[]) {
  const leafNodes: Array<{ value: number; label: string }> = [];
  categories.forEach((category) => {
    if (category.children.length === 0) {
      leafNodes.push({ value: category.id, label: category.name });
      return;
    }
    category.children.forEach((child) => {
      leafNodes.push({ value: child.id, label: `${category.name} / ${child.name}` });
    });
  });
  return leafNodes;
}

function resolveStoreCategoryLabel(categories: SellerStoreCategory[], storeCategoryId: number) {
  if (!storeCategoryId) {
    return "未分类";
  }
  for (const category of categories) {
    if (category.id === storeCategoryId) {
      return category.name;
    }
    for (const child of category.children) {
      if (child.id === storeCategoryId) {
        return `${category.name} / ${child.name}`;
      }
    }
  }
  return "未分类";
}

export default function PublishPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const querySpuNo = (searchParams.get("spuNo") || "").trim() || undefined;
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const routeSyncSpuNoRef = useRef<string | undefined>();
  const [step, setStep] = useState(0);
  const [editorOpen, setEditorOpen] = useState(Boolean(querySpuNo));
  const [sessionBizNo] = useState(() => `seller-publish-${rid()}`);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [brandNo, setBrandNo] = useState("");
  const [categoryId, setCategoryId] = useState<number | null>(null);
  const [price, setPrice] = useState<number | null>(null);
  const [storeCategoryId, setStoreCategoryId] = useState(0);
  const [rows, setRows] = useState<Row[]>([emptyRow()]);
  const [assets, setAssets] = useState<MediaAsset[]>([]);
  const [libraryAssets, setLibraryAssets] = useState<MediaAsset[]>([]);
  const [mainImageId, setMainImageId] = useState<string | undefined>();
  const [detailImageIds, setDetailImageIds] = useState<string[]>([]);
  const [libraryModalOpen, setLibraryModalOpen] = useState(false);
  const [librarySelection, setLibrarySelection] = useState<string[]>([]);
  const [spuNo, setSpuNo] = useState<string | undefined>(querySpuNo);
  const [spuVersion, setSpuVersion] = useState(0);
  const [reviewStatus, setReviewStatus] = useState("draft");
  const [hydrated, setHydrated] = useState<string | undefined>();
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [loadingInventory, setLoadingInventory] = useState(false);
  const [libraryLoading, setLibraryLoading] = useState(false);

  const workbenchQuery = useQuery<SellerWorkbenchResponse>({ queryKey: ["seller-workbench", "publish-v2"], queryFn: fetchSellerWorkbench, staleTime: 30000 });
  const detailQuery = useQuery<ProductDetailResult>({ queryKey: ["seller-product", querySpuNo], queryFn: () => fetchMyProduct(querySpuNo as string), enabled: Boolean(querySpuNo), staleTime: 0 });

  const shop = useMemo(() => workbenchQuery.data?.shops?.[0], [workbenchQuery.data]);
  const storeCategoriesQuery = useQuery({
    queryKey: ["seller-store-categories", shop?.shopNo, "publish"],
    queryFn: () => listSellerStoreCategories(shop?.shopNo ?? ""),
    enabled: Boolean(shop?.shopNo),
    staleTime: 30_000
  });
  const bindingQuery = useQuery({
    queryKey: ["seller-product-store-category", shop?.shopNo, querySpuNo],
    queryFn: () => getSellerProductStoreCategoryBinding(shop?.shopNo ?? "", querySpuNo ?? ""),
    enabled: Boolean(shop?.shopNo && querySpuNo),
    staleTime: 0
  });
  const bizNo = spuNo || sessionBizNo;
  const inventoryKey = useMemo(() => rows.map((row) => row.skuNo || "").filter(Boolean).join("|"), [rows]);
  const sizeRows = useMemo(() => uniqRows(rows).map((row) => ({ ...row, size: trimSize(row.size) })).filter((row) => row.size), [rows]);
  const mainAsset = useMemo(() => assets.find((asset) => String(asset.assetId) === mainImageId), [assets, mainImageId]);
  const leafCategoryOptions = useMemo(() => flattenLeafCategories(storeCategoriesQuery.data ?? []), [storeCategoriesQuery.data]);
  const storeCategoryLabel = useMemo(
    () => resolveStoreCategoryLabel(storeCategoriesQuery.data ?? [], storeCategoryId),
    [storeCategoriesQuery.data, storeCategoryId]
  );
  const issues = useMemo(() => {
    const list: string[] = [];
    const sizes = rows.map((row) => trimSize(row.size)).filter(Boolean);
    if (!title.trim()) list.push("请填写商品名称");
    if (!description.trim()) list.push("请填写商品卖点摘要");
    if (!Number(categoryId ?? 0)) list.push("请填写商品类别");
    if (Number(price ?? 0) <= 0) list.push("请填写商品价格");
    if (!mainImageId) list.push("请选择商品主图");
    if (!sizes.length) list.push("请至少填写一个尺码");
    if (new Set(sizes).size !== sizes.length) list.push("尺码存在重复");
    return list;
  }, [categoryId, description, mainImageId, price, rows, title]);

  useEffect(() => {
    if (querySpuNo && routeSyncSpuNoRef.current === querySpuNo) {
      routeSyncSpuNoRef.current = undefined;
      setEditorOpen(true);
      return;
    }
    setSpuNo(querySpuNo);
    setSpuVersion(0);
    setReviewStatus("draft");
    setStep(0);
    setEditorOpen(Boolean(querySpuNo));
    setHydrated(undefined);
    if (!querySpuNo) {
      setTitle("");
      setDescription("");
      setBrandNo("");
      setCategoryId(null);
      setStoreCategoryId(0);
      setPrice(null);
      setRows([emptyRow()]);
      setAssets([]);
      setMainImageId(undefined);
      setDetailImageIds([]);
    }
  }, [querySpuNo]);

  useEffect(() => {
    let cancelled = false;
    void (async () => {
      try {
        const list = await fetchMediaLibrary({ bizNo });
        if (!cancelled) setAssets(list);
      } catch {
        if (!cancelled) setAssets([]);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [bizNo]);

  useEffect(() => {
    let cancelled = false;
    setLibraryLoading(true);
    void (async () => {
      try {
        const list = await fetchMediaLibrary();
        if (!cancelled) setLibraryAssets(list);
      } catch {
        if (!cancelled) setLibraryAssets([]);
      } finally {
        if (!cancelled) setLibraryLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (!detailQuery.data || !querySpuNo || hydrated === querySpuNo) return;
    const spu = detailQuery.data.product.spu;
    const skus = detailQuery.data.product.skus || [];
    setTitle(spu.title || "");
    setDescription(spu.summary || "");
    setBrandNo(spu.brandNo || "");
    setCategoryId(spu.categoryId || null);
    setPrice(skus[0]?.salePrice ?? null);
    setMainImageId(spu.mainImageAssetIds?.[0]);
    setDetailImageIds((spu.detailImageAssetIds || []).filter((id) => id !== spu.mainImageAssetIds?.[0]));
    setSpuNo(spu.spuNo || querySpuNo);
    setSpuVersion(spu.version || 0);
    setReviewStatus(detailQuery.data.review?.reviewStatus || spu.status || "draft");
    setRows(skus.length ? skus.map((sku) => ({ ...emptyRow(sku.saleAttrs?.size || sku.skuName || ""), skuNo: sku.skuNo })) : [emptyRow()]);
    setHydrated(querySpuNo);
  }, [detailQuery.data, hydrated, querySpuNo]);

  useEffect(() => {
    if (!querySpuNo) {
      return;
    }
    setStoreCategoryId(bindingQuery.data?.storeCategoryId ?? 0);
  }, [bindingQuery.data?.storeCategoryId, querySpuNo]);
  useEffect(() => {
    let cancelled = false;
    const skuNos = inventoryKey ? inventoryKey.split("|").filter(Boolean) : [];
    if (!skuNos.length) return;
    setLoadingInventory(true);
    void (async () => {
      try {
        const stocks = await batchGetSkuInventory(skuNos);
        if (cancelled) return;
        const map = new Map(stocks.map((item) => [item.skuNo, item.totalQty]));
        setRows((prev) => prev.map((row) => {
          if (!row.skuNo) return row;
          const qty = Number(map.get(row.skuNo) ?? 0);
          return { ...row, stockQty: qty, syncedQty: qty, syncStatus: qty > 0 ? "success" : row.syncStatus, syncMessage: qty > 0 ? `已同步库存 ${qty}` : row.syncMessage };
        }));
      } catch {
        if (!cancelled) setRows((prev) => prev.map((row) => ({ ...row, syncMessage: row.syncMessage || "未获取到库存，默认按 0 展示" })));
      } finally {
        if (!cancelled) setLoadingInventory(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [inventoryKey]);

  useEffect(() => {
    if (!assets.length) return;
    if (!mainImageId || !assets.some((asset) => String(asset.assetId) === mainImageId)) setMainImageId(String(assets[0].assetId));
  }, [assets, mainImageId]);

  useEffect(() => {
    setDetailImageIds((prev) => prev.filter((id) => id !== mainImageId));
  }, [mainImageId]);

  const resetComposer = () => {
    setStep(0);
    setTitle("");
    setDescription("");
    setBrandNo("");
    setCategoryId(null);
    setStoreCategoryId(0);
    setPrice(null);
    setRows([emptyRow()]);
    setAssets([]);
    setMainImageId(undefined);
    setDetailImageIds([]);
    setSpuNo(undefined);
    setSpuVersion(0);
    setReviewStatus("draft");
    setHydrated(undefined);
  };

  const openCreateModal = () => {
    if (!querySpuNo) {
      resetComposer();
      router.replace("/seller/publish");
    }
    setEditorOpen(true);
  };

  const closeEditor = () => {
    setEditorOpen(false);
    if (!querySpuNo) {
      resetComposer();
      router.replace("/seller/publish");
    }
  };

  const patchRow = (key: string, patch: Partial<Row>) => setRows((prev) => prev.map((row) => (row.key === key ? { ...row, ...patch } : row)));
  const addRow = () => setRows((prev) => [...prev, emptyRow()]);
  const removeRow = (key: string) => setRows((prev) => (prev.length === 1 ? [emptyRow()] : prev.filter((row) => row.key !== key)));

  const syncMedia = async (targetSpuNo: string) => {
    const ids = [mainImageId, ...detailImageIds].filter(Boolean) as string[];
    if (!ids.length) return;
    await replaceMediaBindings({ bizNo: targetSpuNo, assetIds: ids });
  };

  const mergeSavedRows = (result: SaveDraftResult, sourceRows: Row[]) =>
    result.product.skus.length
      ? result.product.skus.map((sku, index) => {
          const matched = sourceRows.find(
            (row) => trimSize(row.size) === trimSize(sku.saleAttrs.size || "")
          );
          return {
            key: matched?.key || rid(),
            size: sku.saleAttrs.size || matched?.size || `尺码${index + 1}`,
            skuNo: sku.skuNo,
            stockQty: matched?.stockQty || 0,
            syncedQty: matched?.syncedQty || 0,
            syncStatus: matched?.syncStatus || "idle",
            syncMessage: matched?.syncMessage
          };
        })
      : [emptyRow()];

  const saveDraft = async () => {
    if (!shop) {
      notification.warning({ message: "当前账号没有可用店铺", description: "请先完成店铺入驻审核。" });
      return null;
    }
    if (issues.length) {
      notification.warning({ message: issues[0], description: "请先把当前步骤的必填信息补齐。" });
      return null;
    }
    const payload = {
      draft: { title: title.trim(), summary: description.trim(), categoryId: Number(categoryId ?? 0), brandNo: brandNo.trim(), mainImageAssetId: mainImageId, detailImageAssetIds: detailImageIds, submitNote: "", attributeValues: {} },
      skus: sizeRows.map((row) => ({ skuNo: row.skuNo, name: `${title.trim() || "商品"}-${row.size}`, salePrice: Number(price ?? 0), marketPrice: Number(price ?? 0), saleSpecs: { color: "", size: row.size }, initialStock: row.stockQty }))
    };
    setSaving(true);
    try {
      const result = await saveProductDraft({ shopNo: shop.shopNo, draft: payload.draft, skus: payload.skus, spuNo, expectedVersion: spuVersion });
      try {
        await syncMedia(result.product.spu.spuNo);
      } catch (error) {
        notification.warning({ message: "草稿已保存，图片绑定需要重试", description: error instanceof Error ? error.message : "图片绑定失败" });
      }
      try {
        await updateSellerProductStoreCategoryBinding(shop.shopNo, result.product.spu.spuNo, storeCategoryId);
      } catch (error) {
        notification.warning({
          message: "草稿已保存，店内分类绑定需要重试",
          description: error instanceof Error ? error.message : "店内分类绑定失败"
        });
      }
      setSpuNo(result.product.spu.spuNo);
      setSpuVersion(result.product.spu.version);
      setReviewStatus(result.product.spu.status);
      setRows(mergeSavedRows(result, sizeRows));
      routeSyncSpuNoRef.current = result.product.spu.spuNo;
      router.replace(`/seller/publish?spuNo=${encodeURIComponent(result.product.spu.spuNo)}`);
      notification.success({ message: "草稿已保存", description: `商品草稿 ${result.product.spu.spuNo} 已成功写入。` });
      return result;
    } catch (error) {
      notification.error({ message: "保存草稿失败", description: error instanceof Error ? error.message : "请稍后重试" });
      return null;
    } finally {
      setSaving(false);
    }
  };

  const uploadFiles = async (files: FileList | null) => {
    if (!files?.length) return;
    setUploading(true);
    try {
      let existing = assets.map((asset) => asset.assetId);
      const next: MediaAsset[] = [];
      for (const file of Array.from(files)) {
        const asset = await uploadSellerMediaAsset({ file, bizNo, existingAssetIds: existing });
        existing = [...existing, asset.assetId];
        next.push(asset);
      }
      setAssets((prev) => [...prev, ...next]);
      if (!mainImageId && next[0]) setMainImageId(String(next[0].assetId));
      notification.success({ message: "图片上传完成", description: `成功上传 ${next.length} 张图片。` });
    } catch (error) {
      notification.error({ message: "图片上传失败", description: error instanceof Error ? error.message : "请检查网络或 OSS 配置后重试" });
    } finally {
      setUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const toggleDetail = (assetId: string) => {
    if (assetId === mainImageId) return;
    setDetailImageIds((prev) => (prev.includes(assetId) ? prev.filter((item) => item !== assetId) : [...prev, assetId]));
  };

  const toggleLibrarySelection = (assetId: string) => {
    setLibrarySelection((prev) => {
      const next = prev.includes(assetId) ? prev.filter((item) => item !== assetId) : [...prev, assetId];
      return next;
    });
  };

  const mergeLibraryAssetsIntoDraft = (assetIds: string[]) => {
    if (!assetIds.length) return;
    setAssets((prev) => {
      const existingIds = new Set(prev.map((asset) => String(asset.assetId)));
      const appended = libraryAssets.filter(
        (asset) => assetIds.includes(String(asset.assetId)) && !existingIds.has(String(asset.assetId))
      );
      return appended.length ? [...prev, ...appended] : prev;
    });
  };

  const handleLibrarySetMain = () => {
    if (!librarySelection.length) {
      notification.warning({ message: "请先在素材库中选择一张图片" });
      return;
    }
    const assetId = librarySelection[0];
    mergeLibraryAssetsIntoDraft([assetId]);
    setMainImageId(assetId);
    setDetailImageIds((prev) => prev.filter((id) => id !== assetId));
    setLibrarySelection([]);
    setLibraryModalOpen(false);
  };

  const handleLibraryAddDetails = () => {
    if (!librarySelection.length) {
      notification.warning({ message: "请先选择要加入详情图的素材" });
      return;
    }
    mergeLibraryAssetsIntoDraft(librarySelection);
    setDetailImageIds((prev) => {
      const next = [...prev];
      librarySelection.forEach((assetId) => {
        if (!next.includes(assetId) && assetId !== mainImageId) {
          next.push(assetId);
        }
      });
      return next;
    });
    setLibrarySelection([]);
    setLibraryModalOpen(false);
  };

  const saveInventory = async (prefetchedResult?: SaveDraftResult | null) => {
    if (!shop) return false;
    const needDraft = rows.some((row) => trimSize(row.size) && !row.skuNo);
    let result: SaveDraftResult | null = prefetchedResult ?? null;
    let workingRows = result ? mergeSavedRows(result, rows) : rows;
    if ((!spuNo || needDraft) && !result) {
      result = await saveDraft();
      if (!result) return false;
      workingRows = mergeSavedRows(result, rows);
    }
    const targetSpuNo = spuNo || result?.product.spu.spuNo;
    if (!targetSpuNo) return false;
    setSyncing(true);
    let hasError = false;
    for (const row of workingRows) {
      if (!trimSize(row.size) || !row.skuNo) continue;
      const desired = Number(row.stockQty ?? 0);
      const delta = desired - Number(row.syncedQty ?? 0);
      if (delta === 0) {
        patchRow(row.key, { syncStatus: "success", syncMessage: desired > 0 ? `已同步库存 ${desired}` : "库存为 0，无需同步" });
        continue;
      }
      patchRow(row.key, { syncStatus: "syncing", syncMessage: "正在同步库存..." });
      const resp = await adjustInventory({ shopNo: shop.shopNo, skuNo: row.skuNo, spuNo: targetSpuNo, delta, bizNo: `seller-publish:${targetSpuNo}:${row.skuNo}`, reasonCode: "ADJUST_REASON_CODE_SELLER_REPLENISH", remark: "卖家发布商品时同步初始库存" });
      if (resp.success) {
        patchRow(row.key, { syncedQty: desired, syncStatus: "success", syncMessage: `已同步库存 ${desired}` });
      } else {
        hasError = true;
        patchRow(row.key, { syncStatus: "error", syncMessage: resp.failedItems[0]?.reason || "库存同步失败，请稍后再试" });
      }
    }
    setSyncing(false);
    if (hasError) {
      notification.warning({ message: "部分尺码库存同步失败", description: "失败的尺码已经保留在列表中，可以修改后再次保存库存。" });
      return false;
    }
    notification.success({ message: "库存已保存", description: "各尺码库存已经同步到库存服务。" });
    return true;
  };
  const next = async () => {
    if (step === 0) {
      if (issues.length) {
        notification.warning({ message: issues[0], description: "请先把当前步骤的必填信息补齐。" });
        return;
      }
      setStep(1);
      return;
    }
    if (step === 1) {
      const needSync = rows.some((row) => Number(row.stockQty ?? 0) > 0 && Number(row.syncedQty ?? 0) !== Number(row.stockQty ?? 0));
      if (needSync) {
        const ok = await saveInventory();
        if (!ok) return;
      }
      setStep(2);
    }
  };

  const submit = async () => {
    if (!shop) return;
    const result = await saveDraft();
    if (!result) return;
    const needSync = rows.some((row) => Number(row.stockQty ?? 0) > 0 && Number(row.syncedQty ?? 0) !== Number(row.stockQty ?? 0));
    if (needSync) {
      const ok = await saveInventory(result);
      if (!ok) return;
    }
    setSubmitting(true);
    try {
      await submitProductReview({ shopNo: shop.shopNo, spuNo: result.product.spu.spuNo, expectedVersion: result.product.spu.version, submitNote: "" });
      setReviewStatus("reviewing");
      notification.success({ message: "提交审核成功", description: `商品 ${result.product.spu.spuNo} 已进入审核队列。` });
    } catch (error) {
      notification.error({ message: "提交审核失败", description: error instanceof Error ? error.message : "请稍后重试" });
    } finally {
      setSubmitting(false);
    }
  };

  if (workbenchQuery.isLoading) return <Card loading style={{ minHeight: 560, borderRadius: 24 }} />;

  if (!shop) {
    return (
      <Card bordered={false} style={{ borderRadius: 24, minHeight: 420 }}>
        <Empty description="当前账号还没有可发布商品的店铺。" image={Empty.PRESENTED_IMAGE_SIMPLE}>
          <Space>
            <Link href="/seller/workbench"><Button type="primary">前往商家工作台</Button></Link>
            <Link href="/seller/profile"><Button>查看商家资料</Button></Link>
          </Space>
        </Empty>
      </Card>
    );
  }

  return (
    <div style={{ minHeight: "100%", background: "radial-gradient(circle at top, rgba(255,221,185,0.42), transparent 28%), linear-gradient(180deg, #fffaf5 0%, #fff7ef 36%, #fffdf9 100%)", padding: "24px 0 48px" }}>
      <div style={{ width: "min(840px, calc(100vw - 32px))", margin: "0 auto", display: "grid", gap: 20 }}>
        <Card bordered={false} style={{ borderRadius: 28, border: "1px solid #f1dfcf", boxShadow: "0 20px 56px rgba(82, 44, 10, 0.08)", background: "linear-gradient(180deg, rgba(255,255,255,0.96), rgba(255,248,241,0.98))", padding: 28 }}>
          <Space direction="vertical" size={18} style={{ width: "100%" }}>
            <div style={{ display: "grid", gap: 8 }}>
              <Text style={{ fontSize: 12, letterSpacing: 2, color: "#b56a2b", textTransform: "uppercase" }}>Seller Studio</Text>
              <Title level={3} style={{ margin: 0, color: "#3d2006" }}>商品发布</Title>
              <Paragraph style={{ margin: 0, color: "#7f6244", lineHeight: 1.8 }}>
                从这里进入商品发布弹窗。新建商品时按三步完成信息填写、库存设置和提交审核，也可以直接去草稿箱继续编辑已有商品。
              </Paragraph>
            </div>
            <Space wrap size={12}>
              <Button type="primary" size="large" onClick={openCreateModal} style={{ borderRadius: 14, background: "#ff7a1a", minWidth: 136 }}>
                发布商品
              </Button>
              <Link href="/seller/products?scope=drafts">
                <Button size="large" icon={<FolderOpenOutlined />} style={{ borderRadius: 14, minWidth: 136 }}>
                  查看草稿箱
                </Button>
              </Link>
            </Space>
            <div style={{ display: "grid", gap: 10, gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))" }}>
              <div style={{ borderRadius: 18, border: "1px solid #f2e4d6", background: "#fffaf6", padding: 16 }}>
                <Text type="secondary">当前店铺</Text>
                <div style={{ marginTop: 6 }}><Text strong>{shop.shopDisplayName || shop.shopName}</Text></div>
              </div>
              <div style={{ borderRadius: 18, border: "1px solid #f2e4d6", background: "#fffaf6", padding: 16 }}>
                <Text type="secondary">发布状态</Text>
                <div style={{ marginTop: 6 }}><Text strong>{statusText(reviewStatus)}</Text></div>
              </div>
              <div style={{ borderRadius: 18, border: "1px solid #f2e4d6", background: "#fffaf6", padding: 16 }}>
                <Text type="secondary">当前草稿</Text>
                <div style={{ marginTop: 6 }}><Text strong>{spuNo || "还未创建"}</Text></div>
              </div>
            </div>
          </Space>
        </Card>

        {detailQuery.isLoading ? <Card bordered={false} style={{ borderRadius: 24, minHeight: 220 }}><Space direction="vertical" align="center" style={{ width: "100%", padding: "48px 0" }}><Spin /><Text type="secondary">正在加载商品草稿...</Text></Space></Card> : null}
        {detailQuery.isError ? <Alert type="warning" showIcon message="商品草稿加载失败" description={detailQuery.error instanceof Error ? detailQuery.error.message : "请稍后重试"} /> : null}
        <Modal
          open={editorOpen}
          onCancel={closeEditor}
          footer={null}
          width={960}
          destroyOnClose={false}
          styles={{ body: { padding: 0, background: "#fffaf5" } }}
          title={
            <div style={{ display: "flex", flexWrap: "wrap", gap: 8, alignItems: "center" }}>
              <Text strong style={{ fontSize: 18, color: "#3d2006" }}>商品发布</Text>
              <Tag color="gold">店铺：{shop.shopDisplayName || shop.shopName}</Tag>
              <Tag color={reviewStatus === "reviewing" ? "processing" : "default"}>状态：{statusText(reviewStatus)}</Tag>
              {spuNo ? <Tag color="blue">草稿：{spuNo}</Tag> : <Tag>新建商品</Tag>}
            </div>
          }
        >
          <div style={{ padding: 24, display: "grid", gap: 18 }}>
            <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #f2dfcf", boxShadow: "0 14px 36px rgba(98, 60, 24, 0.06)", padding: 12 }}>
              <Steps current={step} responsive items={[{ title: "填写商品信息", description: "名称、描述、图片、价格、尺码" }, { title: "填写商品库存", description: "按尺码逐项填写库存" }, { title: "提交审核", description: "确认摘要并提交审核" }]} />
            </Card>

            {step === 0 ? <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #f2dfcf", boxShadow: "0 14px 36px rgba(98, 60, 24, 0.06)", padding: 28 }}>
              <Space direction="vertical" size={24} style={{ width: "100%" }}>
                <div><Title level={4} style={{ marginBottom: 8 }}>填写商品信息</Title><Paragraph style={{ margin: 0, color: "#7a6a58" }}>先把商品的基础资料补齐，保存后会生成商品草稿和对应的尺码规格。</Paragraph></div>
                <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
                  <label style={{ display: "grid", gap: 8 }}><Text strong>商品名称</Text><Input size="large" placeholder="例如：冬季直筒牛仔裤" value={title} onChange={(e) => setTitle(e.target.value)} /></label>
                  <label style={{ display: "grid", gap: 8 }}><Text strong>商品类别（类目 ID）</Text><InputNumber size="large" min={1} precision={0} style={{ width: "100%" }} placeholder="例如：10001" value={categoryId} onChange={(value) => setCategoryId(typeof value === "number" ? value : null)} /></label>
                </div>
                <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
                  <label style={{ display: "grid", gap: 8 }}><Text strong>商品卖点摘要</Text><Input.TextArea rows={5} placeholder="建议写清楚卖点、面料、版型和适用场景，买家会先看到这里。" value={description} onChange={(e) => setDescription(e.target.value)} /></label>
                  <label style={{ display: "grid", gap: 8 }}><Text strong>品牌编码（可选）</Text><Input size="large" placeholder="例如：BRAND_DENIM_LAB" value={brandNo} onChange={(e) => setBrandNo(e.target.value)} /></label>
                </div>
                <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
                  <label style={{ display: "grid", gap: 8 }}>
                    <Text strong>店内分类（可选）</Text>
                    <Select
                      size="large"
                      value={storeCategoryId || undefined}
                      allowClear
                      placeholder={leafCategoryOptions.length ? "选择店铺自己的分类" : "请先去“店内分类”页创建分类"}
                      options={leafCategoryOptions}
                      loading={storeCategoriesQuery.isLoading}
                      onChange={(value) => setStoreCategoryId(typeof value === "number" ? value : 0)}
                      onClear={() => setStoreCategoryId(0)}
                    />
                  </label>
                  <div style={{ display: "grid", gap: 8 }}>
                    <Text strong>当前店内分类</Text>
                    <div style={{ minHeight: 48, borderRadius: 14, border: "1px solid #f0dfcf", background: "#fffaf6", display: "flex", alignItems: "center", padding: "0 16px", color: "#6b4b2d", fontWeight: 600 }}>
                      {storeCategoryLabel}
                    </div>
                  </div>
                </div>
                {!leafCategoryOptions.length ? (
                  <Alert
                    type="info"
                    showIcon
                    message="店内分类还未配置"
                    description={
                      <span>
                        当前商品仍可先保存草稿；如果希望买家在店铺页按“上衣 / 下装 / 外套 / 鞋子”这类店内分类浏览，请先到{" "}
                        <Link href="/seller/store-categories">店内分类</Link> 页面创建分类。
                      </span>
                    }
                  />
                ) : null}
                <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
                  <label style={{ display: "grid", gap: 8 }}><Text strong>统一售价</Text><InputNumber size="large" min={0.01} precision={2} style={{ width: "100%" }} placeholder="请输入售价" value={price} onChange={(value) => setPrice(typeof value === "number" ? value : null)} addonBefore="¥" /></label>
                  <div style={{ display: "grid", gap: 8 }}><Text strong>当前价格预览</Text><div style={{ minHeight: 48, borderRadius: 14, border: "1px solid #f0dfcf", background: "#fffaf6", display: "flex", alignItems: "center", padding: "0 16px", color: "#ff5a1f", fontWeight: 700, fontSize: 20 }}>{price ? money(price) : "待填写"}</div></div>
                </div>
                  <div style={{ display: "grid", gap: 14 }}>
                    <Space align="center" style={{ justifyContent: "space-between", width: "100%" }}>
                      <div style={{ display: "grid", gap: 4 }}><Text strong>商品图片</Text><Text type="secondary">支持上传多张图片，选 1 张主图，其余可勾选为详情图。</Text></div>
                      <div>
                        <input ref={fileInputRef} hidden type="file" accept="image/*" multiple onChange={(e) => void uploadFiles(e.target.files)} />
                        <Button type="primary" icon={<UploadOutlined />} loading={uploading} onClick={() => fileInputRef.current?.click()} style={{ borderRadius: 12, background: "#ff7a1a" }}>上传图片</Button>
                      </div>
                    </Space>
                    <Space wrap style={{ justifyContent: "flex-end" }}>
                      <Button type="default" icon={<FolderOpenOutlined />} onClick={() => setLibraryModalOpen(true)}>从素材库选择</Button>
                    </Space>
                  {mainAsset ? (
                    <div style={{ display: "grid", gap: 12, gridTemplateColumns: "minmax(0, 1.2fr) minmax(240px, 320px)" }}>
                      <div style={{ borderRadius: 20, overflow: "hidden", border: "1px solid #f0dfcf", background: "#fff6ef", minHeight: 280 }}>
                        {assetPreviewSrc(mainAsset) ? (
                          <img src={assetPreviewSrc(mainAsset)} alt={mainAsset.name} style={{ width: "100%", height: "100%", objectFit: "cover", display: "block" }} />
                        ) : (
                          <div style={{ minHeight: 280, display: "flex", alignItems: "center", justifyContent: "center", color: "#9b7b5a" }}>主图预览生成中</div>
                        )}
                      </div>
                      <div style={{ borderRadius: 20, border: "1px solid #f0dfcf", background: "#fffaf6", padding: 18, display: "grid", gap: 10, alignContent: "start" }}>
                        <Text strong>当前主图</Text>
                        <Text>{mainAsset.name}</Text>
                        <Text type="secondary">已上传 {assets.length} 张图片</Text>
                        <Text type="secondary">详情图 {detailImageIds.length} 张</Text>
                        <Tag color="orange" style={{ width: "fit-content" }}>主图已选</Tag>
                      </div>
                    </div>
                  ) : null}
                  {!assets.length ? <div style={{ border: "1px dashed #f0d7c2", borderRadius: 18, padding: "40px 24px", textAlign: "center", background: "#fffaf5" }}><Text type="secondary">还没有上传图片，先上传主图和详情图。</Text></div> : <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))" }}>{assets.map((asset) => {
                    const id = String(asset.assetId);
                    const isMain = mainImageId === id;
                    const isDetail = detailImageIds.includes(id);
                    const previewSrc = assetPreviewSrc(asset);
                    return <div key={asset.id} style={{ borderRadius: 18, border: isMain ? "2px solid #ff7a1a" : "1px solid #f0dfcf", background: isMain ? "#fff4ea" : "#fff", overflow: "hidden" }}><div style={{ aspectRatio: "1 / 1", background: "#fff6ef" }}>{previewSrc ? <img src={previewSrc} alt={asset.name} style={{ width: "100%", height: "100%", objectFit: "cover" }} /> : <div style={{ width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center", color: "#9b7b5a" }}>预览生成中</div>}</div><div style={{ padding: 12, display: "grid", gap: 8 }}><Text strong ellipsis>{asset.name}</Text><Space wrap><Button type={isMain ? "primary" : "default"} size="small" onClick={() => setMainImageId(id)} style={isMain ? { background: "#ff7a1a" } : undefined}>{isMain ? "当前主图" : "设为主图"}</Button><Button size="small" disabled={isMain} onClick={() => toggleDetail(id)}>{isDetail ? "移出详情图" : "加入详情图"}</Button></Space></div></div>;
                  })}</div>}
                </div>
                <div style={{ display: "grid", gap: 14 }}>
                  <Space align="center" style={{ justifyContent: "space-between", width: "100%" }}><div style={{ display: "grid", gap: 4 }}><Text strong>尺码列表</Text><Text type="secondary">先定义商品有哪些尺码，库存下一步单独填写。</Text></div><Button icon={<PlusOutlined />} onClick={addRow} style={{ borderRadius: 12 }}>新增尺码</Button></Space>
                  <div style={{ display: "grid", gap: 10 }}>{rows.map((row, index) => <div key={row.key} style={{ display: "grid", gap: 12, gridTemplateColumns: "minmax(0, 1fr) auto", padding: 14, borderRadius: 16, border: "1px solid #f1e1d2", background: "#fffaf6" }}><div style={{ display: "grid", gap: 8 }}><Text strong>尺码 {index + 1}</Text><Input size="large" placeholder="例如：S / M / L / XL" value={row.size} onChange={(e) => patchRow(row.key, { size: e.target.value, syncStatus: "idle", syncMessage: undefined, skuNo: undefined, syncedQty: 0 })} /></div><Button danger icon={<DeleteOutlined />} onClick={() => removeRow(row.key)} style={{ alignSelf: "end", borderRadius: 12 }}>删除</Button></div>)}</div>
                </div>
              </Space>
            </Card> : null}

            {step === 1 ? <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #f2dfcf", boxShadow: "0 14px 36px rgba(98, 60, 24, 0.06)", padding: 28 }}>
              <Space direction="vertical" size={20} style={{ width: "100%" }}>
                <div><Title level={4} style={{ marginBottom: 8 }}>填写商品库存</Title><Paragraph style={{ margin: 0, color: "#7a6a58" }}>按尺码分别填写库存。库存可以是 0；如果填写了非 0 数量，进入下一步前会先自动同步到库存服务。</Paragraph></div>
                {loadingInventory ? <Alert type="info" showIcon message="正在读取现有库存" description="已有商品会尽量回填真实库存，读取失败时默认按 0 展示。" /> : null}
                <div style={{ display: "grid", gap: 12 }}>{sizeRows.map((row) => <div key={row.key} style={{ display: "grid", gap: 14, gridTemplateColumns: "minmax(120px, 160px) minmax(0, 1fr) minmax(120px, 160px)", alignItems: "center", padding: 16, borderRadius: 18, border: "1px solid #f1e1d2", background: "#fffaf6" }}><div style={{ display: "grid", gap: 4 }}><Text type="secondary">尺码</Text><Text strong>{row.size}</Text></div><div style={{ display: "grid", gap: 4 }}><Text type="secondary">SKU 编号</Text><Text>{row.skuNo || "保存草稿后生成"}</Text></div><div style={{ display: "grid", gap: 8 }}><Text type="secondary">库存数量</Text><InputNumber size="large" min={0} precision={0} style={{ width: "100%" }} value={row.stockQty} onChange={(value) => patchRow(row.key, { stockQty: typeof value === "number" ? value : 0, syncStatus: "idle", syncMessage: undefined })} /></div><div style={{ gridColumn: "1 / -1", display: "flex", alignItems: "center", gap: 8, color: "#8a6b4a" }}>{row.syncStatus === "syncing" ? <LoadingOutlined style={{ color: "#fa8c16" }} /> : row.syncStatus === "success" ? <CheckCircleFilled style={{ color: "#52c41a" }} /> : row.syncStatus === "error" ? <ExclamationCircleFilled style={{ color: "#ff4d4f" }} /> : <span style={{ width: 14, height: 14, borderRadius: 999, background: "#e8e8e8", display: "inline-block" }} /> }<Text style={{ color: "inherit" }}>{row.syncMessage || (row.skuNo ? "等待同步" : "保存草稿后生成 SKU")}</Text></div></div>)}</div>
              </Space>
            </Card> : null}

            {step === 2 ? <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #f2dfcf", boxShadow: "0 14px 36px rgba(98, 60, 24, 0.06)", padding: 28 }}>
              <Space direction="vertical" size={20} style={{ width: "100%" }}>
                <div><Title level={4} style={{ marginBottom: 8 }}>提交审核</Title><Paragraph style={{ margin: 0, color: "#7a6a58" }}>这里展示提交前的核心摘要。满足条件后，页面会自动补保存草稿，并在需要时先同步库存，再提交审核。</Paragraph></div>
                <div style={{ display: "grid", gap: 12 }}>
                  <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">商品名称</Text><div style={{ marginTop: 6 }}><Text strong>{title || "未填写"}</Text></div></div>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">商品类别</Text><div style={{ marginTop: 6 }}><Text strong>{categoryId || "未填写"}</Text></div></div>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">店内分类</Text><div style={{ marginTop: 6 }}><Text strong>{storeCategoryLabel}</Text></div></div>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">品牌编码</Text><div style={{ marginTop: 6 }}><Text strong>{brandNo || "未填写"}</Text></div></div>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">统一售价</Text><div style={{ marginTop: 6 }}><Text strong>{price ? money(price) : "未填写"}</Text></div></div>
                    <div style={{ padding: 16, borderRadius: 16, background: "#fffaf6", border: "1px solid #f1e1d2" }}><Text type="secondary">当前状态</Text><div style={{ marginTop: 6 }}><Text strong>{statusText(reviewStatus)}</Text></div></div>
                  </div>
                  <div style={{ padding: 18, borderRadius: 18, border: "1px solid #f1e1d2", background: "#fffaf6", display: "grid", gap: 10 }}><Text strong>图片摘要</Text><Text>主图：{mainImageId ? "已选择" : "未选择"}</Text><Text>图片总数：{assets.length}</Text><Text>详情图数量：{detailImageIds.length}</Text></div>
                  <div style={{ padding: 18, borderRadius: 18, border: "1px solid #f1e1d2", background: "#fffaf6", display: "grid", gap: 10 }}><Text strong>尺码与库存摘要</Text>{sizeRows.length ? <div style={{ display: "grid", gap: 8 }}>{sizeRows.map((row) => <div key={row.key} style={{ display: "flex", justifyContent: "space-between", gap: 12 }}><Text>{row.size}</Text><Text type="secondary">库存 {row.stockQty}</Text></div>)}</div> : <Text type="secondary">暂无尺码</Text>}</div>
                  {reviewStatus === "reviewing" ? <Alert type="success" showIcon message="商品已经在审核中" description="你可以返回商品列表查看审核进度。" /> : null}
                </div>
              </Space>
            </Card> : null}

            <Card bordered={false} style={{ borderRadius: 24, border: "1px solid #f2dfcf", boxShadow: "0 14px 36px rgba(98, 60, 24, 0.06)", padding: 20 }}>
              <div style={{ display: "flex", justifyContent: "space-between", gap: 12, flexWrap: "wrap" }}>
                <Space wrap>
                  {step > 0 ? <Button size="large" onClick={() => setStep((prev) => Math.max(0, prev - 1))} style={{ borderRadius: 12 }}>上一步</Button> : null}
                  <Link href="/seller/products?scope=drafts"><Button size="large" style={{ borderRadius: 12 }}>查看草稿箱</Button></Link>
                </Space>
                <Space wrap>
                  {step === 0 ? <><Button size="large" loading={saving} onClick={() => void saveDraft()} style={{ borderRadius: 12 }}>保存草稿</Button><Button type="primary" size="large" loading={saving} disabled={saving || uploading} onClick={() => void next()} style={{ borderRadius: 12, background: "#ff7a1a" }}>下一步</Button></> : null}
                  {step === 1 ? <><Button size="large" loading={syncing || saving} onClick={() => void saveInventory()} style={{ borderRadius: 12 }}>保存库存</Button><Button type="primary" size="large" loading={syncing || saving} disabled={syncing || saving} onClick={() => void next()} style={{ borderRadius: 12, background: "#ff7a1a" }}>下一步</Button></> : null}
                  {step === 2 ? <Button type="primary" size="large" loading={saving || syncing || submitting} disabled={reviewStatus === "reviewing"} onClick={() => void submit()} style={{ borderRadius: 12, background: "#ff7a1a" }}>提交审核</Button> : null}
                </Space>
              </div>
            </Card>
          </div>
        </Modal>
        <Modal
          open={libraryModalOpen}
          onCancel={() => {
            setLibraryModalOpen(false);
            setLibrarySelection([]);
          }}
          title="从素材库选择图片"
          width={860}
          footer={
            <Space style={{ width: "100%", justifyContent: "flex-end" }}>
              <Button onClick={() => {
                setLibraryModalOpen(false);
                setLibrarySelection([]);
              }}>取消</Button>
              <Button type="default" onClick={handleLibraryAddDetails} disabled={!librarySelection.length}>
                加入详情图
              </Button>
              <Button type="primary" onClick={handleLibrarySetMain} disabled={!librarySelection.length}>
                设为主图
              </Button>
            </Space>
          }
        >
          {libraryLoading ? (
            <div style={{ display: "flex", justifyContent: "center", padding: "48px 0" }}>
              <Spin indicator={<LoadingOutlined style={{ fontSize: 28 }} spin />} />
            </div>
          ) : libraryAssets.length ? (
            <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))" }}>
              {libraryAssets.map((asset) => {
                const id = String(asset.assetId);
                const selected = librarySelection.includes(id);
                const previewSrc = assetPreviewSrc(asset);
                return (
                  <div
                    key={`${asset.id}-${id}`}
                    onClick={() => toggleLibrarySelection(id)}
                    style={{
                      borderRadius: 20,
                      border: selected ? "2px solid #ff7a1a" : "1px solid #f0dfcf",
                      background: selected ? "#fff4ea" : "#fff",
                      overflow: "hidden",
                      cursor: "pointer",
                      display: "grid"
                    }}
                  >
                    <div style={{ aspectRatio: "1 / 1", background: "#fff6ef", position: "relative" }}>
                      {previewSrc ? (
                        <img src={previewSrc} alt={asset.name} style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                      ) : (
                        <div style={{ width: "100%", height: "100%", display: "flex", alignItems: "center", justifyContent: "center", color: "#9b7b5a" }}>
                          预览生成中
                        </div>
                      )}
                      {selected && (
                        <span
                          style={{
                            position: "absolute",
                            top: 8,
                            right: 8,
                            background: "#ff7a1a",
                            color: "#fff",
                            borderRadius: 999,
                            padding: "0 10px",
                            fontSize: 12,
                            fontWeight: 600
                          }}
                        >
                          已选
                        </span>
                      )}
                    </div>
                    <div style={{ padding: 12, display: "grid", gap: 6 }}>
                      <Text strong ellipsis>{asset.name}</Text>
                      <Text type="secondary" ellipsis>{`ID: ${asset.assetId}`}</Text>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <Empty description="素材库当前暂无图片" />
          )}
        </Modal>
      </div>
    </div>
  );
}
