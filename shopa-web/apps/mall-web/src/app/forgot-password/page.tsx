"use client";

import Link from "next/link";
import { useState } from "react";
import { Button, Card, Form, Input, Typography, notification } from "antd";
import { LockOutlined, MobileOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { resetPasswordBySms, sendSmsCode } from "@/features/iam/api";
import { useI18n } from "@shopa/ui";

const { Title, Text } = Typography;

type FormValue = {
  phone: string;
  smsCode: string;
  newPassword: string;
  confirmPassword: string;
};

export default function ForgotPasswordPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [sendingCode, setSendingCode] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [form] = Form.useForm<FormValue>();

  const handleSendCode = async () => {
    const phone = form.getFieldValue("phone");
    if (!phone) {
      notification.warning({ message: isZh ? "请先输入手机号" : "Please input phone first" });
      return;
    }
    setSendingCode(true);
    try {
      const data = await sendSmsCode({ scene: 4, phone });
      const next = data?.resendAfterSeconds ?? 60;
      setCountdown(next);
      notification.success({ message: isZh ? "验证码已发送" : "Verification code sent" });
      const timer = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            clearInterval(timer);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    } catch (_err) {
      notification.error({ message: isZh ? "验证码发送失败" : "Failed to send verification code" });
    } finally {
      setSendingCode(false);
    }
  };

  const handleReset = async (values: FormValue) => {
    if (values.newPassword !== values.confirmPassword) {
      notification.error({ message: isZh ? "两次密码不一致" : "Password confirmation does not match" });
      return;
    }

    setSubmitting(true);
    try {
      await resetPasswordBySms({
        phone: values.phone,
        smsCode: values.smsCode,
        newPassword: values.newPassword
      });
      localStorage.removeItem("shopa_mall_access_token");
      localStorage.removeItem("shopa_mall_refresh_token");
      localStorage.removeItem("shopa-mall-auth");
      notification.success({ message: isZh ? "密码重置成功" : "Password reset success" });
      window.location.href = "/login";
    } catch (_err) {
      notification.error({ message: isZh ? "密码重置失败" : "Password reset failed" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="tb-register-wrap">
      <Card className="tb-register-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          {isZh ? "找回密码" : "Reset Password"}
        </Title>
        <Text type="secondary">{isZh ? "通过短信验证码重置登录密码。" : "Use phone verification code to reset your password."}</Text>

        <Form form={form} layout="vertical" style={{ marginTop: 18 }} onFinish={handleReset} autoComplete="off">
          <Form.Item name="phone" label={isZh ? "手机号" : "Phone"} rules={[{ required: true, message: isZh ? "请输入手机号" : "Please input phone" }]}>
            <Input placeholder={isZh ? "手机号" : "Phone"} prefix={<MobileOutlined />} />
          </Form.Item>

          <Form.Item label={isZh ? "验证码" : "Verification Code"} required>
            <div className="tb-register-code-row">
              <Form.Item name="smsCode" noStyle rules={[{ required: true, message: isZh ? "请输入验证码" : "Please input code" }]}>
                <Input placeholder={isZh ? "短信验证码" : "SMS code"} prefix={<SafetyCertificateOutlined />} />
              </Form.Item>
              <Button onClick={handleSendCode} disabled={countdown > 0} loading={sendingCode}>
                {countdown > 0 ? `${countdown}s` : isZh ? "发送验证码" : "Send Code"}
              </Button>
            </div>
          </Form.Item>

          <Form.Item
            name="newPassword"
            label={isZh ? "新密码" : "New Password"}
            rules={[
              { required: true, message: isZh ? "请输入新密码" : "Please input new password" },
              { min: 8, message: isZh ? "至少 8 位字符" : "At least 8 chars" }
            ]}
          >
            <Input.Password placeholder={isZh ? "新密码" : "New password"} prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label={isZh ? "确认密码" : "Confirm Password"}
            rules={[{ required: true, message: isZh ? "请确认密码" : "Please confirm password" }]}
          >
            <Input.Password placeholder={isZh ? "确认密码" : "Confirm password"} prefix={<LockOutlined />} />
          </Form.Item>

          <Button type="primary" htmlType="submit" block loading={submitting}>
            {isZh ? "重置密码" : "Reset Password"}
          </Button>
        </Form>

        <div className="tb-register-hint">
          <Text type="secondary">{isZh ? "返回" : "Back to"}</Text>
          <Link href="/login">{isZh ? "登录" : "Login"}</Link>
        </div>
      </Card>
    </section>
  );
}
