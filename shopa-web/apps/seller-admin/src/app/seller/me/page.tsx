"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Alert, Button, Card, Col, Input, Modal, Row, Space, Tag, Typography, notification } from "antd";
import { ExclamationCircleOutlined, LogoutOutlined, SafetyOutlined, StopOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { useAuth } from "@/features/iam/useAuth";
import { useAuthStore } from "@/features/iam/store";

const { Title, Text, Paragraph } = Typography;

export default function MyAccountPage() {
  const router = useRouter();
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const { logout } = useAuth();
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const [deactivateOpen, setDeactivateOpen] = useState(false);
  const [confirmText, setConfirmText] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleLogout = async () => {
    setSubmitting(true);
    try {
      await logout();
      notification.success({ message: isZh ? "已退出登录" : "Signed out" });
      router.replace("/login");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDeactivate = async () => {
    if (confirmText !== "注销") {
      notification.warning({ message: isZh ? "请输入“注销”以确认" : "Please type the confirmation text" });
      return;
    }

    setSubmitting(true);
    try {
      // 当前版本先完成前端侧流程，后续可接入后端真实注销接口。
      await logout();
      notification.success({
        message: isZh ? "注销申请已提交（演示流程）" : "Deactivation request submitted (demo flow)",
        description: isZh ? "已为你执行安全退出。" : "You have been signed out safely."
      });
      router.replace("/login");
    } finally {
      setSubmitting(false);
      setDeactivateOpen(false);
      setConfirmText("");
    }
  };

  return (
    <section className="seller-page">
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "我的" : "My Account"}</Title>
          <Text type="secondary">{isZh ? "管理账号安全、登录状态与风险操作。" : "Manage security, session status and risk actions."}</Text>
        </div>
      </header>

      <Alert
        type="info"
        showIcon
        message={isZh ? "账号安全提醒" : "Security Notice"}
        description={
          isZh
            ? "建议定期更改密码、开启登录保护，并避免在公共设备保存登录状态。"
            : "Change password regularly, enable login protection, and avoid saving sessions on public devices."
        }
      />

      <Row gutter={[14, 14]}>
        <Col xs={24} xl={14}>
          <Card title={isZh ? "账号概览" : "Account Overview"}>
            <div className="seller-health-list">
              <div>
                <span>{isZh ? "当前状态" : "Current Status"}</span>
                <Tag color="success">{isZh ? "已登录" : "Signed In"}</Tag>
              </div>
              <div>
                <span>{isZh ? "Token SID" : "Token SID"}</span>
                <Text code>{tokenPair?.sid || "-"}</Text>
              </div>
              <div>
                <span>{isZh ? "Access 令牌" : "Access Token"}</span>
                <Tag color={tokenPair?.accessToken ? "processing" : "default"}>{tokenPair?.accessToken ? "ACTIVE" : "N/A"}</Tag>
              </div>
              <div>
                <span>{isZh ? "Refresh 令牌" : "Refresh Token"}</span>
                <Tag color={tokenPair?.refreshToken ? "processing" : "default"}>{tokenPair?.refreshToken ? "ACTIVE" : "N/A"}</Tag>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={24} xl={10}>
          <Card title={isZh ? "安全中心" : "Security Center"} className="seller-side-card">
            <Space direction="vertical" size={10} style={{ width: "100%" }}>
              <Link href="/forgot-password">
                <Button block icon={<SafetyOutlined />}>
                  {isZh ? "修改登录密码" : "Change Password"}
                </Button>
              </Link>
              <Button block icon={<LogoutOutlined />} onClick={handleLogout} loading={submitting}>
                {isZh ? "退出登录" : "Sign Out"}
              </Button>
              <Button block danger icon={<StopOutlined />} onClick={() => setDeactivateOpen(true)}>
                {isZh ? "注销账号" : "Deactivate Account"}
              </Button>
            </Space>
            <Paragraph type="secondary" style={{ marginTop: 12, marginBottom: 0 }}>
              {isZh ? "注销后将无法继续使用当前账号登录商家后台。" : "After deactivation, this account cannot access seller admin."}
            </Paragraph>
          </Card>
        </Col>
      </Row>

      <Card title={isZh ? "常用入口" : "Quick Access"}>
        <Space wrap>
          <Link href="/seller/workbench">
            <Button>{isZh ? "返回工作台" : "Back to Workbench"}</Button>
          </Link>
          <Link href="/seller/sales">
            <Button>{isZh ? "查看销售数据" : "View Sales Analytics"}</Button>
          </Link>
          <Link href="/seller/profile">
            <Button>{isZh ? "编辑店主资料" : "Edit Profile"}</Button>
          </Link>
        </Space>
      </Card>

      <Modal
        title={isZh ? "确认注销账号" : "Confirm Account Deactivation"}
        open={deactivateOpen}
        onCancel={() => {
          setDeactivateOpen(false);
          setConfirmText("");
        }}
        onOk={handleDeactivate}
        okButtonProps={{ danger: true, loading: submitting }}
        okText={isZh ? "确认注销" : "Confirm"}
        cancelText={isZh ? "取消" : "Cancel"}
      >
        <Space direction="vertical" style={{ width: "100%" }}>
          <Alert
            type="warning"
            showIcon
            icon={<ExclamationCircleOutlined />}
            message={isZh ? "此操作不可逆" : "This action is irreversible"}
            description={isZh ? "请输入“注销”以确认操作。" : "Type the confirmation text to continue."}
          />
          <Input
            value={confirmText}
            onChange={(e) => setConfirmText(e.target.value)}
            placeholder={isZh ? "请输入：注销" : "Type confirmation text"}
          />
        </Space>
      </Modal>
    </section>
  );
}
