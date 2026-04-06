'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { Alert, Button, Card, Empty, List, Progress, Space, Tag, Typography, notification } from 'antd';
import { CloudUploadOutlined, ReloadOutlined } from '@ant-design/icons';
import { fetchMediaLibrary, OSS_CORS_BLOCKED_ERROR, uploadSellerMediaAsset } from './api';
import { MediaAsset } from './types';

const { Text, Paragraph } = Typography;

const statusMeta: Record<MediaAsset['status'], { label: string; color: string }> = {
  idle: { label: 'Pending', color: 'default' },
  uploading: { label: 'Uploading', color: 'processing' },
  processing: { label: 'Processing', color: 'orange' },
  done: { label: 'Ready', color: 'success' },
  failed: { label: 'Failed', color: 'error' },
};

type Props = {
  bizNo?: string;
  bizType?: string;
  bindingField?: string;
  title?: string;
  description?: string;
  onAssetsChange?: (assets: MediaAsset[]) => void;
};

export const MediaLibrary = ({
  bizNo = 'seller-console-library',
  bizType,
  bindingField,
  title = '商品素材库',
  description = '上传商品主图、细节图和 SKU 图。这里只展示真正已经绑定到当前发布会话的素材。',
  onAssetsChange,
}: Props) => {
  const [assets, setAssets] = useState<MediaAsset[]>([]);
  const [loading, setLoading] = useState(false);
  const [uploadingNames, setUploadingNames] = useState<string[]>([]);
  const [errorMessage, setErrorMessage] = useState('');
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const hasProcessingAsset = useMemo(
    () => assets.some((asset) => asset.status === 'uploading' || asset.status === 'processing'),
    [assets],
  );

  const loadAssets = async () => {
    setLoading(true);
    try {
      const library = await fetchMediaLibrary({ bizNo, bizType, bindingField });
      setAssets(library);
      onAssetsChange?.(library);
      setErrorMessage('');
    } catch (error) {
      const message = error instanceof Error ? error.message : '素材库加载失败';
      setErrorMessage(message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadAssets();
  }, [bizNo, bizType, bindingField]);

  useEffect(() => {
    if (!hasProcessingAsset) {
      return;
    }
    const timer = window.setInterval(() => {
      void loadAssets();
    }, 4000);
    return () => window.clearInterval(timer);
  }, [hasProcessingAsset, bizNo, bizType, bindingField]);

  const handlePickFiles = () => {
    fileInputRef.current?.click();
  };

  const handleFilesSelected = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files ?? []);
    event.target.value = '';
    if (!files.length) {
      return;
    }

    setUploadingNames(files.map((file) => file.name));
    setErrorMessage('');

    let currentAssets = [...assets];

    for (const file of files) {
      try {
        await uploadSellerMediaAsset({
          file,
          bizNo,
          bizType,
          bindingField,
          existingAssetIds: currentAssets.map((asset) => asset.assetId),
        });
        currentAssets = await fetchMediaLibrary({ bizNo, bizType, bindingField });
        setAssets(currentAssets);
        onAssetsChange?.(currentAssets);
      } catch (error) {
        const message =
          error instanceof Error && error.message === OSS_CORS_BLOCKED_ERROR
            ? 'OSS 跨域未放行，浏览器已被拦截，请先补齐 bucket CORS 配置。'
            : error instanceof Error
              ? error.message
              : '素材上传失败';
        setErrorMessage(message);
        notification.error({
          message: '上传失败',
          description: message,
        });
        break;
      }
    }

    setUploadingNames([]);
  };

  return (
    <Card
      bordered={false}
      title={title}
      extra={
        <Space>
          <Button icon={<ReloadOutlined />} onClick={() => void loadAssets()} loading={loading}>
            刷新
          </Button>
          <Button type="primary" icon={<CloudUploadOutlined />} onClick={handlePickFiles} loading={uploadingNames.length > 0}>
            上传素材
          </Button>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            hidden
            onChange={(event) => void handleFilesSelected(event)}
          />
        </Space>
      }
    >
      <Paragraph type="secondary" style={{ marginBottom: 16 }}>
        {description}
      </Paragraph>
      {errorMessage ? (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          message="素材库需要处理"
          description={errorMessage}
        />
      ) : null}
      {uploadingNames.length ? (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          message="正在上传"
          description={`当前文件: ${uploadingNames.join('、')}`}
        />
      ) : null}
      <List
        grid={{ gutter: 16, xs: 1, sm: 2, xl: 3 }}
        loading={loading}
        dataSource={assets}
        locale={{
          emptyText: (
            <Empty
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description="还没有上传任何商品素材，先把主图和细节图放进来。"
            />
          ),
        }}
        renderItem={(asset) => (
          <List.Item>
            <Card hoverable bodyStyle={{ padding: 14 }}>
              <Space direction="vertical" size={12} style={{ width: '100%' }}>
                <div
                  style={{
                    width: '100%',
                    aspectRatio: '1 / 1',
                    borderRadius: 16,
                    overflow: 'hidden',
                    background: 'linear-gradient(135deg, #fff1e6, #ffe1c2)',
                    border: '1px solid #ffd8bf',
                  }}
                >
                  <img
                    src={asset.publicUrl || asset.thumbnail}
                    alt={asset.name}
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                  />
                </div>
                <div>
                  <Text strong>{asset.name}</Text>
                  <div style={{ marginTop: 8 }}>
                    <Tag color={statusMeta[asset.status].color}>{statusMeta[asset.status].label}</Tag>
                    <Text type="secondary">Asset #{asset.assetId}</Text>
                  </div>
                </div>
                <Progress
                  percent={asset.progress}
                  status={asset.status === 'failed' ? 'exception' : asset.status === 'done' ? 'success' : 'active'}
                />
                <Text type="secondary">
                  {asset.uploadedAt ? `最近处理: ${new Date(asset.uploadedAt).toLocaleString()}` : '等待处理时间'}
                </Text>
              </Space>
            </Card>
          </List.Item>
        )}
      />
    </Card>
  );
};

export default MediaLibrary;

