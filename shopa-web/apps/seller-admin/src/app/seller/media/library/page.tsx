"use client";

import { useMemo, useState } from "react";
import { Button, Card, Drawer, Segmented, Space, Tag, Typography } from "antd";
import { EyeOutlined, PlusOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";

const { Title, Paragraph, Text } = Typography;

type AssetRow = {
  id: string;
  name: string;
  type: "IMAGE" | "VIDEO";
  scene: "PRODUCT" | "CERT";
  size: string;
  uploadedAt: string;
};

const assets: AssetRow[] = [
  { id: "A9001", name: "shoes-main-1.jpg", type: "IMAGE", scene: "PRODUCT", size: "1.8MB", uploadedAt: "2026-03-28 09:20" },
  { id: "A9002", name: "shop-license.jpg", type: "IMAGE", scene: "CERT", size: "2.4MB", uploadedAt: "2026-03-27 19:12" },
  { id: "A9003", name: "lookbook-video.mp4", type: "VIDEO", scene: "PRODUCT", size: "18.6MB", uploadedAt: "2026-03-27 15:43" },
  { id: "A9004", name: "id-front.jpg", type: "IMAGE", scene: "CERT", size: "1.1MB", uploadedAt: "2026-03-27 11:08" }
];

export default function MediaLibraryPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [sceneFilter, setSceneFilter] = useState<"ALL" | "PRODUCT" | "CERT">("ALL");
  const [selected, setSelected] = useState<AssetRow | null>(null);

  const rows = useMemo(() => {
    if (sceneFilter === "ALL") return assets;
    return assets.filter((item) => item.scene === sceneFilter);
  }, [sceneFilter]);

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "素材库" : "Media Library"}</Title>
          <Text type="secondary">{isZh ? "集中管理商品素材和证件素材。" : "Manage product assets and certificate assets."}</Text>
        </div>
        <Space wrap>
          <Button type="primary" icon={<PlusOutlined />}>
            {isZh ? "上传素材" : "Upload Asset"}
          </Button>
        </Space>
      </header>

      <Card
        extra={
          <Segmented
            value={sceneFilter}
            onChange={(val) => setSceneFilter(val as any)}
            options={[
              { label: isZh ? "全部" : "All", value: "ALL" },
              { label: isZh ? "商品素材" : "Product", value: "PRODUCT" },
              { label: isZh ? "证件素材" : "Certificate", value: "CERT" }
            ]}
          />
        }
      >
        <div className="seller-media-grid">
          {rows.map((item) => (
            <article key={item.id} className="seller-media-card">
              <div className="seller-media-preview" />
              <h4>{item.name}</h4>
              <p>{item.id}</p>
              <div className="seller-media-meta">
                <Tag>{item.type}</Tag>
                <Tag color={item.scene === "CERT" ? "gold" : "blue"}>{item.scene}</Tag>
                <Tag>{item.size}</Tag>
              </div>
              <Button icon={<EyeOutlined />} onClick={() => setSelected(item)}>
                {isZh ? "查看详情" : "View"}
              </Button>
            </article>
          ))}
        </div>
      </Card>

      <Drawer
        open={Boolean(selected)}
        width={420}
        onClose={() => setSelected(null)}
        title={isZh ? "素材详情" : "Asset Detail"}
      >
        {selected ? (
          <Space direction="vertical" size={12} style={{ width: "100%" }}>
            <div className="seller-media-preview seller-media-preview-large" />
            <Paragraph>
              <Text strong>{isZh ? "文件名：" : "File: "}</Text>
              {selected.name}
            </Paragraph>
            <Paragraph>
              <Text strong>ID: </Text>
              {selected.id}
            </Paragraph>
            <Paragraph>
              <Text strong>{isZh ? "上传时间：" : "Uploaded: "}</Text>
              {selected.uploadedAt}
            </Paragraph>
            <Paragraph>
              <Text strong>{isZh ? "场景：" : "Scene: "}</Text>
              {selected.scene}
            </Paragraph>
            <Space>
              <Button type="primary">{isZh ? "复制素材ID" : "Copy Asset ID"}</Button>
              <Button danger>{isZh ? "移除引用" : "Remove Binding"}</Button>
            </Space>
          </Space>
        ) : null}
      </Drawer>
    </section>
  );
}
