"use client";

import { useMemo, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Alert,
  Button,
  Card,
  Checkbox,
  Drawer,
  Empty,
  List,
  Space,
  Spin,
  Tag,
  Typography,
  notification
} from "antd";
import {
  DeleteOutlined,
  EyeOutlined,
  InboxOutlined,
  PlusOutlined,
  ReloadOutlined
} from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import {
  batchUnbindSellerMediaAssets,
  fetchMediaLibrary,
  SELLER_MEDIA_LIBRARY_BIZ_NO,
  uploadSellerMediaAsset
} from "@/features/media/api";
import { MediaAsset } from "@/features/media/types";

const { Title, Paragraph, Text } = Typography;

type UploadDraftStatus = "pending" | "uploading" | "done" | "failed";

type UploadDraft = {
  id: string;
  file: File;
  status: UploadDraftStatus;
  error?: string;
};

function previewSrc(asset?: MediaAsset) {
  return asset?.publicUrl || asset?.thumbnail || "";
}

function assetStatusLabel(status: MediaAsset["status"], isZh: boolean) {
  if (!isZh) {
    switch (status) {
      case "done":
        return "Ready";
      case "failed":
        return "Failed";
      case "processing":
        return "Processing";
      case "uploading":
        return "Uploading";
      default:
        return "Idle";
    }
  }
  switch (status) {
    case "done":
      return "可用";
    case "failed":
      return "失败";
    case "processing":
      return "处理中";
    case "uploading":
      return "上传中";
    default:
      return "待处理";
  }
}

function draftStatusLabel(status: UploadDraftStatus, isZh: boolean) {
  if (!isZh) {
    switch (status) {
      case "uploading":
        return "Uploading";
      case "done":
        return "Done";
      case "failed":
        return "Failed";
      default:
        return "Pending";
    }
  }
  switch (status) {
    case "uploading":
      return "上传中";
    case "done":
      return "已完成";
    case "failed":
      return "失败";
    default:
      return "待上传";
  }
}

function draftStatusColor(status: UploadDraftStatus) {
  switch (status) {
    case "uploading":
      return "processing";
    case "done":
      return "success";
    case "failed":
      return "error";
    default:
      return "default";
  }
}

function formatFileSize(size: number, isZh: boolean) {
  if (size < 1024) {
    return `${size} B`;
  }
  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KB`;
  }
  const text = `${(size / (1024 * 1024)).toFixed(2)} MB`;
  return isZh ? text : text;
}

function buildDraftId(file: File) {
  return `${file.name}-${file.size}-${file.lastModified}`;
}

function dedupeDrafts(previous: UploadDraft[], files: File[]) {
  const draftMap = new Map(previous.map((item) => [item.id, item]));
  for (const file of files) {
    const id = buildDraftId(file);
    if (!draftMap.has(id)) {
      draftMap.set(id, {
        id,
        file,
        status: "pending"
      });
    }
  }
  return Array.from(draftMap.values());
}

export default function MediaLibraryPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [selectedAsset, setSelectedAsset] = useState<MediaAsset | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [uploadDrafts, setUploadDrafts] = useState<UploadDraft[]>([]);
  const [uploading, setUploading] = useState(false);

  const libraryQuery = useQuery({
    queryKey: ["seller-media-library"],
    queryFn: () => fetchMediaLibrary(),
    staleTime: 15000
  });

  const assets = libraryQuery.data ?? [];

  const pendingDraftCount = useMemo(
    () => uploadDrafts.filter((item) => item.status === "pending" || item.status === "failed").length,
    [uploadDrafts]
  );

  const toggleSelected = (assetId: string) => {
    setSelectedIds((prev) =>
      prev.includes(assetId) ? prev.filter((item) => item !== assetId) : [...prev, assetId]
    );
  };

  const refreshLibrary = async () => {
    await libraryQuery.refetch();
  };

  const addFilesToQueue = (files: FileList | null) => {
    if (!files?.length) {
      return;
    }
    const imageFiles = Array.from(files).filter((file) => file.type.startsWith("image/"));
    setUploadDrafts((prev) => dedupeDrafts(prev, imageFiles));
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const removeDraft = (draftId: string) => {
    setUploadDrafts((prev) => prev.filter((item) => item.id !== draftId));
  };

  const clearDrafts = () => {
    if (uploading) {
      return;
    }
    setUploadDrafts([]);
  };

  const startBatchUpload = async () => {
    if (!uploadDrafts.length) {
      notification.warning({
        message: isZh ? "请先选择图片" : "Please pick images first"
      });
      return;
    }

    setUploading(true);
    let successCount = 0;
    let failedCount = 0;
    let existingAssetIds = assets.map((asset) => asset.assetId);

    for (const draft of uploadDrafts) {
      if (draft.status === "done") {
        continue;
      }

      setUploadDrafts((prev) =>
        prev.map((item) =>
          item.id === draft.id
            ? {
                ...item,
                status: "uploading",
                error: undefined
              }
            : item
        )
      );

      try {
        const asset = await uploadSellerMediaAsset({
          file: draft.file,
          bizNo: SELLER_MEDIA_LIBRARY_BIZ_NO,
          existingAssetIds
        });
        existingAssetIds = [...existingAssetIds, asset.assetId];
        successCount += 1;
        setUploadDrafts((prev) =>
          prev.map((item) =>
            item.id === draft.id
              ? {
                  ...item,
                  status: "done",
                  error: undefined
                }
              : item
          )
        );
      } catch (error) {
        failedCount += 1;
        setUploadDrafts((prev) =>
          prev.map((item) =>
            item.id === draft.id
              ? {
                  ...item,
                  status: "failed",
                  error: error instanceof Error ? error.message : isZh ? "上传失败，请稍后重试" : "Upload failed"
                }
              : item
          )
        );
      }
    }

    await refreshLibrary();
    setUploading(false);

    if (successCount > 0 && failedCount === 0) {
      notification.success({
        message: isZh ? "素材上传完成" : "Upload completed",
        description: isZh ? `已成功上传 ${successCount} 张图片。` : `${successCount} image(s) uploaded successfully.`
      });
    } else if (successCount > 0 && failedCount > 0) {
      notification.warning({
        message: isZh ? "部分图片上传完成" : "Partial upload completed",
        description: isZh
          ? `成功 ${successCount} 张，失败 ${failedCount} 张。失败项已保留在队列里，修正后可以继续上传。`
          : `${successCount} succeeded and ${failedCount} failed.`
      });
    } else {
      notification.error({
        message: isZh ? "素材上传失败" : "Upload failed",
        description: isZh ? "这批图片都没有上传成功，请检查错误信息后重试。" : "No image was uploaded successfully."
      });
    }
  };

  const handleBatchDelete = async (targetIds: string[]) => {
    if (!targetIds.length) {
      notification.warning({ message: isZh ? "请先选择素材" : "Please select assets first" });
      return;
    }
    setDeleting(true);
    try {
      await batchUnbindSellerMediaAssets({ assetIds: targetIds });
      notification.success({
        message: isZh ? "素材已移出素材库" : "Assets removed",
        description: isZh ? `已移除 ${targetIds.length} 个素材。` : `${targetIds.length} asset(s) removed.`
      });
      setSelectedIds((prev) => prev.filter((id) => !targetIds.includes(id)));
      if (selectedAsset && targetIds.includes(String(selectedAsset.assetId))) {
        setSelectedAsset(null);
      }
      await refreshLibrary();
    } catch (error) {
      notification.error({
        message: isZh ? "删除素材失败" : "Remove failed",
        description: error instanceof Error ? error.message : isZh ? "请稍后重试" : "Please retry later"
      });
    } finally {
      setDeleting(false);
    }
  };

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "素材库" : "Media Library"}</Title>
          <Text type="secondary">
            {isZh
              ? "支持一次选择多张商品图片，先进入上传队列，再统一开始上传。上传完成后可直接在发布商品时复用。"
              : "Pick multiple product images at once, queue them, and upload them together."}
          </Text>
        </div>
        <Space wrap>
          <Button icon={<ReloadOutlined />} onClick={() => void refreshLibrary()}>
            {isZh ? "刷新" : "Refresh"}
          </Button>
          <Button icon={<PlusOutlined />} onClick={() => fileInputRef.current?.click()}>
            {isZh ? "选择图片" : "Choose images"}
          </Button>
          <Button type="primary" icon={<InboxOutlined />} loading={uploading} onClick={() => void startBatchUpload()}>
            {isZh ? "开始批量上传" : "Start upload"}
          </Button>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            hidden
            onChange={(event) => addFilesToQueue(event.target.files)}
          />
        </Space>
      </header>

      <Card
        style={{ marginBottom: 20 }}
        bodyStyle={{ display: "grid", gap: 18 }}
      >
        <div
          onClick={() => fileInputRef.current?.click()}
          onDragOver={(event) => event.preventDefault()}
          onDrop={(event) => {
            event.preventDefault();
            addFilesToQueue(event.dataTransfer.files);
          }}
          style={{
            border: "1px dashed #f2c9a5",
            background: "linear-gradient(180deg, #fffaf4 0%, #fff4ea 100%)",
            borderRadius: 24,
            padding: "28px 24px",
            cursor: "pointer"
          }}
        >
          <Space direction="vertical" size={10} style={{ width: "100%", textAlign: "center" }}>
            <div style={{ fontSize: 32, color: "#d97706" }}>
              <InboxOutlined />
            </div>
            <Text strong style={{ fontSize: 18, color: "#7c3f10" }}>
              {isZh ? "把多张商品图一次性拖到这里，或点击选择图片" : "Drop multiple images here or click to browse"}
            </Text>
            <Text type="secondary">
              {isZh ? "支持 JPG / PNG / WebP。选择后会先进入上传队列，不会立刻上传。" : "JPG / PNG / WebP supported."}
            </Text>
          </Space>
        </div>

        {uploadDrafts.length ? (
          <div style={{ display: "grid", gap: 12 }}>
            <div className="seller-toolbar" style={{ justifyContent: "space-between", gap: 12 }}>
              <Text type="secondary">
                {isZh
                  ? `当前队列 ${uploadDrafts.length} 张，待上传 ${pendingDraftCount} 张`
                  : `${uploadDrafts.length} file(s) in queue, ${pendingDraftCount} pending`}
              </Text>
              <Button onClick={clearDrafts} disabled={uploading}>
                {isZh ? "清空队列" : "Clear queue"}
              </Button>
            </div>

            <List
              dataSource={uploadDrafts}
              renderItem={(draft) => (
                <List.Item
                  actions={[
                    <Button
                      key="remove"
                      type="text"
                      danger
                      disabled={uploading && draft.status === "uploading"}
                      onClick={() => removeDraft(draft.id)}
                    >
                      {isZh ? "移除" : "Remove"}
                    </Button>
                  ]}
                  style={{
                    background: "#fffaf5",
                    border: "1px solid #f4dfcf",
                    borderRadius: 16,
                    padding: "14px 16px"
                  }}
                >
                  <List.Item.Meta
                    title={
                      <Space wrap size={[8, 8]}>
                        <Text strong>{draft.file.name}</Text>
                        <Tag color={draftStatusColor(draft.status)}>{draftStatusLabel(draft.status, isZh)}</Tag>
                      </Space>
                    }
                    description={
                      <Space direction="vertical" size={4}>
                        <Text type="secondary">{formatFileSize(draft.file.size, isZh)}</Text>
                        {draft.error ? <Text type="danger">{draft.error}</Text> : null}
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          </div>
        ) : (
          <Alert
            type="info"
            showIcon
            message={isZh ? "还没有加入上传队列" : "No files in queue"}
            description={
              isZh
                ? "先选择多张图片，它们会先出现在这里。确认无误后再点击“开始批量上传”。"
                : "Pick multiple images first, then start the batch upload."
            }
          />
        )}
      </Card>

      <Card>
        <div className="seller-toolbar" style={{ justifyContent: "space-between", gap: 12 }}>
          <Text type="secondary">
            {isZh ? `当前共有 ${assets.length} 个素材` : `${assets.length} assets in library`}
          </Text>
          <Button
            danger
            icon={<DeleteOutlined />}
            disabled={!selectedIds.length || deleting}
            loading={deleting}
            onClick={() => void handleBatchDelete(selectedIds)}
          >
            {isZh
              ? `批量删除${selectedIds.length ? `（${selectedIds.length}）` : ""}`
              : `Delete selected${selectedIds.length ? ` (${selectedIds.length})` : ""}`}
          </Button>
        </div>

        {libraryQuery.isLoading ? (
          <div style={{ display: "flex", justifyContent: "center", padding: "48px 0" }}>
            <Spin />
          </div>
        ) : assets.length ? (
          <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
            {assets.map((asset) => {
              const assetId = String(asset.assetId);
              const checked = selectedIds.includes(assetId);
              const cover = previewSrc(asset);
              return (
                <article
                  key={asset.id}
                  style={{
                    borderRadius: 20,
                    border: checked ? "1px solid #f08a24" : "1px solid #f0dfcf",
                    background: "#fffaf5",
                    boxShadow: checked ? "0 12px 30px rgba(219, 126, 39, 0.16)" : "0 10px 24px rgba(98, 60, 24, 0.06)",
                    overflow: "hidden"
                  }}
                >
                  <div style={{ position: "relative", aspectRatio: "4 / 3", background: "#f7e7da" }}>
                    {cover ? (
                      <img
                        src={cover}
                        alt={asset.name}
                        style={{ width: "100%", height: "100%", objectFit: "cover", display: "block" }}
                      />
                    ) : (
                      <div
                        style={{
                          width: "100%",
                          height: "100%",
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          color: "#9b7b5a"
                        }}
                      >
                        {isZh ? "预览生成中" : "Preview pending"}
                      </div>
                    )}
                    <div style={{ position: "absolute", top: 12, left: 12 }}>
                      <Checkbox checked={checked} onChange={() => toggleSelected(assetId)} />
                    </div>
                  </div>
                  <div style={{ padding: 16, display: "grid", gap: 10 }}>
                    <div style={{ display: "grid", gap: 4 }}>
                      <Text strong ellipsis={{ tooltip: asset.name }}>
                        {asset.name}
                      </Text>
                      <Text type="secondary">Asset ID: {asset.assetId}</Text>
                    </div>
                    <Space wrap size={[8, 8]}>
                      <Tag color={asset.status === "done" ? "success" : asset.status === "failed" ? "error" : "processing"}>
                        {assetStatusLabel(asset.status, isZh)}
                      </Tag>
                      {asset.uploadedAt ? <Tag>{asset.uploadedAt}</Tag> : null}
                    </Space>
                    <Space>
                      <Button icon={<EyeOutlined />} onClick={() => setSelectedAsset(asset)}>
                        {isZh ? "查看" : "Preview"}
                      </Button>
                      <Button danger icon={<DeleteOutlined />} loading={deleting} onClick={() => void handleBatchDelete([assetId])}>
                        {isZh ? "删除" : "Delete"}
                      </Button>
                    </Space>
                  </div>
                </article>
              );
            })}
          </div>
        ) : (
          <Empty description={isZh ? "素材库还是空的，先批量上传几张商品图吧。" : "No assets yet."} />
        )}
      </Card>

      <Drawer
        open={Boolean(selectedAsset)}
        width={460}
        onClose={() => setSelectedAsset(null)}
        title={isZh ? "素材详情" : "Asset detail"}
      >
        {selectedAsset ? (
          <Space direction="vertical" size={14} style={{ width: "100%" }}>
            <div style={{ borderRadius: 20, overflow: "hidden", border: "1px solid #f0dfcf", background: "#fff7f0" }}>
              {previewSrc(selectedAsset) ? (
                <img
                  src={previewSrc(selectedAsset)}
                  alt={selectedAsset.name}
                  style={{ width: "100%", display: "block", aspectRatio: "4 / 3", objectFit: "cover" }}
                />
              ) : (
                <div style={{ aspectRatio: "4 / 3", display: "flex", alignItems: "center", justifyContent: "center" }}>
                  {isZh ? "预览生成中" : "Preview pending"}
                </div>
              )}
            </div>
            <Paragraph style={{ marginBottom: 0 }}>
              <Text strong>{isZh ? "文件名：" : "File: "}</Text>
              {selectedAsset.name}
            </Paragraph>
            <Paragraph style={{ marginBottom: 0 }}>
              <Text strong>Asset ID: </Text>
              {selectedAsset.assetId}
            </Paragraph>
            <Paragraph style={{ marginBottom: 0 }}>
              <Text strong>{isZh ? "状态：" : "Status: "}</Text>
              {assetStatusLabel(selectedAsset.status, isZh)}
            </Paragraph>
            <Paragraph style={{ marginBottom: 0 }}>
              <Text strong>{isZh ? "说明：" : "Note: "}</Text>
              {isZh
                ? "这里的删除只会把素材从当前卖家素材库解绑，不会直接物理删除对象存储里的原始文件。"
                : "Delete only removes the seller library binding, not the original object from storage."}
            </Paragraph>
          </Space>
        ) : null}
      </Drawer>
    </section>
  );
}
