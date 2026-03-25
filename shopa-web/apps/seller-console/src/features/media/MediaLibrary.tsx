'use client';

import { Button, Card, Divider, List, Progress, Space, Tag, Typography } from 'antd';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { fetchMediaLibrary, pollUploadStatus, startUpload } from './api';
import { MediaAsset, UploadJob } from './types';

const { Text } = Typography;

const statusMeta: Record<MediaAsset['status'], { label: string; color: string }> = {
  idle: { label: 'Pending', color: 'default' },
  uploading: { label: 'Uploading', color: 'processing' },
  processing: { label: 'Processing', color: 'orange' },
  done: { label: 'Complete', color: 'success' },
  failed: { label: 'Failed', color: 'error' },
};

export const MediaLibrary = () => {
  const [assets, setAssets] = useState<MediaAsset[]>([]);
  const [job, setJob] = useState<UploadJob | null>(null);
  const [statusHint, setStatusHint] = useState('');

  useEffect(() => {
    fetchMediaLibrary().then(setAssets);
  }, []);

  useEffect(() => {
    if (!job || job.status === 'done') {
      return;
    }
    const handle = setInterval(async () => {
      const current = await pollUploadStatus(job.uploadId);
      setJob(current);
      setStatusHint(`Upload ${current.uploadId}: ${current.status} (${current.progress}%)`);
    }, 3000);
    return () => clearInterval(handle);
  }, [job]);

  const uploading = useMemo(
    () => assets.some((asset) => asset.status === 'uploading' || asset.status === 'processing'),
    [assets],
  );

  const handleStartUpload = useCallback(async () => {
    const newJob = await startUpload();
    setJob(newJob);
    setStatusHint(`Upload ${newJob.uploadId} queued. ${uploading ? 'Waiting for previous upload to finish.' : ''}`);
  }, [uploading]);

  return (
    <Card
      title="Media Assets"
      bordered={false}
      extra={
        <Button type="primary" onClick={handleStartUpload} aria-label="Start new media upload">
          Upload Asset
        </Button>
      }
    >
      <Text role="status" aria-live="polite" type={job?.status === 'failed' ? 'danger' : 'secondary'}>
        {statusHint || 'All assets are synchronized.'}
      </Text>
      <Divider />
      <List
        grid={{ column: 2, gutter: 16 }}
        dataSource={assets}
        locale={{ emptyText: 'No media has been uploaded yet.' }}
        renderItem={(asset) => (
          <List.Item>
            <Card hoverable bodyStyle={{ padding: 16 }}>
              <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <Space align="center" size="middle">
                  <img
                    src={asset.thumbnail}
                    alt={`${asset.name} thumbnail`}
                    width={64}
                    height={64}
                    style={{ objectFit: 'cover', borderRadius: 8 }}
                  />
                  <div>
                    <Text strong>{asset.name}</Text>
                    <div>
                      <Tag color={statusMeta[asset.status].color}>{statusMeta[asset.status].label}</Tag>
                    </div>
                  </div>
                </Space>
                <Progress
                  percent={asset.progress}
                  status={asset.status === 'failed' ? 'exception' : 'active'}
                  aria-label={`Asset ${asset.name} upload progress ${asset.progress}%`}
                />
                <Space size="small">
                  <Text type="secondary">Last updated:</Text>
                  <Text type="secondary" strong>
                    {asset.uploadedAt ? new Date(asset.uploadedAt).toLocaleString() : 'Not yet uploaded'}
                  </Text>
                </Space>
              </Space>
            </Card>
          </List.Item>
        )}
      />
    </Card>
  );
};

