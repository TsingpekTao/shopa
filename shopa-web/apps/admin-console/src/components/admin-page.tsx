"use client";

import { ReactNode } from "react";
import { Space, Typography } from "antd";

const { Title, Text } = Typography;

export function AdminPage({
  title,
  subtitle,
  extra,
  children
}: {
  title: string;
  subtitle?: string;
  extra?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="admin-page-section">
      <header className="admin-page-head">
        <div className="admin-page-head-main">
          <Title level={4} className="admin-page-title">
            {title}
          </Title>
          {subtitle ? (
            <Text type="secondary" className="admin-page-subtitle">
              {subtitle}
            </Text>
          ) : null}
        </div>
        {extra ? <Space>{extra}</Space> : null}
      </header>
      <div className="admin-page-body">{children}</div>
    </section>
  );
}
