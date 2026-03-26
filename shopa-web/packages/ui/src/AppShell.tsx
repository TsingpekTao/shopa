"use client";

import { ReactNode } from "react";
import { ConfigProvider, Layout } from "antd";
import "./styles.css";

type Props = {
  children: ReactNode;
};

export function AppShell({ children }: Props) {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#ff6a1a",
          borderRadius: 12,
          colorText: "#2f1b09",
          colorBgContainer: "#ffffff"
        }
      }}
    >
      <a className="skip-link" href="#main-content">
        Skip to main content
      </a>
      <Layout className="app-shell">
        <Layout.Header className="app-header" role="banner">
          <div className="logo">Shopa Seller Console</div>
          <div className="header-badge" aria-label="Seller workspace status">
            Seller Workspace
          </div>
        </Layout.Header>
        <Layout.Content id="main-content" className="app-content" role="main" tabIndex={-1}>
          {children}
        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
}
