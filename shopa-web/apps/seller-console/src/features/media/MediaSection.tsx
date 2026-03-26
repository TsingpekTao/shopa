'use client';

import { Card, Divider, Typography } from 'antd';
import { MediaLibrary } from './MediaLibrary';

const { Text } = Typography;

export const MediaSection = () => (
  <Card
    title="Media + Proofs"
    bordered={false}
    style={{
      background: 'var(--media-card-bg, #f9fbff)',
      borderRadius: 16,
      boxShadow: '0 15px 35px rgba(15, 23, 42, 0.08)',
    }}
  >
    <Text type="secondary">
      Upload stories, documents, or certs that the catalog and inventory teams will consume. Rendered once
      and persisted across the publish workflow.
    </Text>
    <Divider />
    <MediaLibrary />
  </Card>
);

export default MediaSection;
