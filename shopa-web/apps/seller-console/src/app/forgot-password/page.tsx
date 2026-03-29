"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Button, Card, Form, Input, Typography, notification } from "antd";
import { LockOutlined, MobileOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { resetPasswordBySms, sendSmsCode } from "@/features/iam/api";

const { Title, Text } = Typography;

type ForgotPasswordForm = {
  phone: string;
  smsCode: string;
  newPassword: string;
  confirmPassword: string;
};

export default function ForgotPasswordPage() {
  const router = useRouter();
  const [sendingCode, setSendingCode] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [form] = Form.useForm<ForgotPasswordForm>();
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, []);

  const startCountdown = (seconds: number) => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
    }
    setCountdown(seconds);
    timerRef.current = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          if (timerRef.current) {
            clearInterval(timerRef.current);
          }
          timerRef.current = null;
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  const handleSendCode = async () => {
    const phone = form.getFieldValue("phone");
    if (!phone) {
      notification.warning({ message: "请先输入手机号" });
      return;
    }
    setSendingCode(true);
    try {
      const data = await sendSmsCode({ scene: 4, phone: phone.trim() });
      startCountdown(data?.resendAfterSeconds ?? 60);
      notification.success({ message: "验证码已发送" });
    } catch (err) {
      // 降低账号枚举风险：对外统一提示，不回显具体失败细节。
      console.error("[forgot-password] send sms failed", err);
      notification.error({ message: "验证码发送失败，请稍后再试" });
    } finally {
      setSendingCode(false);
    }
  };

  const handleSubmit = async (values: ForgotPasswordForm) => {
    setSubmitting(true);
    try {
      await resetPasswordBySms({
        phone: values.phone.trim(),
        smsCode: values.smsCode.trim(),
        newPassword: values.newPassword
      });
      // 重置成功后清理本地旧令牌，避免后续页面携带失效 token。
      localStorage.removeItem("shopa_access_token");
      localStorage.removeItem("shopa_refresh_token");
      localStorage.removeItem("shopa-auth");
      notification.success({ message: "密码重置成功，请使用新密码登录" });
      router.replace("/login");
    } catch (err) {
      // 降低账号枚举风险：统一提示验证码/手机号/密码重置失败。
      console.error("[forgot-password] reset password failed", err);
      notification.error({ message: "重置失败，请检查验证码后重试" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="tb-register-wrap">
      <Card className="tb-register-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          找回密码
        </Title>
        <Text type="secondary">通过手机号验证码重置登录密码。</Text>

        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 18 }}
          onFinish={handleSubmit}
          autoComplete="off"
        >
          <Form.Item name="phone" label="手机号" rules={[{ required: true, message: "请输入手机号" }]}>
            <Input placeholder="请输入手机号" prefix={<MobileOutlined />} />
          </Form.Item>

          <Form.Item label="短信验证码" required>
            <div className="tb-register-code-row">
              <Form.Item name="smsCode" noStyle rules={[{ required: true, message: "请输入短信验证码" }]}>
                <Input placeholder="请输入验证码" prefix={<SafetyCertificateOutlined />} />
              </Form.Item>
              <Button onClick={handleSendCode} disabled={countdown > 0} loading={sendingCode}>
                {countdown > 0 ? `${countdown}s后重发` : "发送验证码"}
              </Button>
            </div>
          </Form.Item>

          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: "请输入新密码" },
              { min: 8, message: "密码至少 8 位" },
              {
                pattern: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[^\w\s]).{8,20}$/,
                message: "密码需包含字母、数字和特殊字符"
              }
            ]}
          >
            <Input.Password placeholder="请输入新密码" prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认新密码"
            dependencies={["newPassword"]}
            rules={[
              { required: true, message: "请再次输入新密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("newPassword") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的密码不一致"));
                }
              })
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" prefix={<LockOutlined />} />
          </Form.Item>

          <Button type="primary" htmlType="submit" block loading={submitting}>
            确认重置
          </Button>
        </Form>

        <div className="tb-register-hint">
          <Text type="secondary">想起密码了？</Text>
          <Link href="/login">返回登录</Link>
        </div>
      </Card>
    </section>
  );
}
