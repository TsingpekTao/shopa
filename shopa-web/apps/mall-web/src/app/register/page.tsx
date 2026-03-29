"use client";

import Link from "next/link";
import { useState } from "react";
import { Button, Card, Form, Input, Typography, notification } from "antd";
import { LockOutlined, MobileOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { registerByPassword, sendSmsCode } from "@/features/iam/api";
import { useI18n } from "@shopa/ui";

const { Title, Text } = Typography;

type FormValue = {
  phone: string;
  smsCode: string;
  password: string;
  confirmPassword: string;
};

export default function RegisterPage() {
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
      const data = await sendSmsCode({ scene: 1, phone });
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

  const handleRegister = async (values: FormValue) => {
    setSubmitting(true);
    try {
      await registerByPassword(values);
      notification.success({ message: isZh ? "注册成功" : "Register success" });
      window.location.href = "/login";
    } catch (_err) {
      notification.error({ message: isZh ? "注册失败" : "Register failed" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="tb-register-wrap">
      <Card className="tb-register-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          {isZh ? "买家注册" : "Buyer Register"}
        </Title>
        <Text type="secondary">{isZh ? "创建买家账号，开始逛商城。" : "Create a customer account to start shopping."}</Text>

        <Form form={form} layout="vertical" style={{ marginTop: 18 }} onFinish={handleRegister} autoComplete="off">
          <Form.Item name="phone" label={isZh ? "手机号" : "Phone"} rules={[{ required: true, message: isZh ? "请输入手机号" : "Please input phone" }]}>
            <Input placeholder={isZh ? "手机号" : "Phone"} prefix={<MobileOutlined />} />
          </Form.Item>

          <Form.Item label={isZh ? "短信验证码" : "SMS Code"} required>
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
            name="password"
            label={isZh ? "密码" : "Password"}
            rules={[
              { required: true, message: isZh ? "请输入密码" : "Please input password" },
              { min: 8, message: isZh ? "至少 8 位字符" : "At least 8 chars" }
            ]}
          >
            <Input.Password placeholder={isZh ? "密码" : "Password"} prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label={isZh ? "确认密码" : "Confirm Password"}
            dependencies={["password"]}
            rules={[
              { required: true, message: isZh ? "请确认密码" : "Please confirm password" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error(isZh ? "两次密码不一致" : "Passwords do not match"));
                }
              })
            ]}
          >
            <Input.Password placeholder={isZh ? "确认密码" : "Confirm password"} prefix={<LockOutlined />} />
          </Form.Item>

          <Button type="primary" htmlType="submit" block loading={submitting}>
            {isZh ? "注册" : "Register"}
          </Button>
        </Form>

        <div className="tb-register-hint">
          <Text type="secondary">{isZh ? "已有账号？" : "Already have an account?"}</Text>
          <Link href="/login">{isZh ? "去登录" : "Login"}</Link>
        </div>
      </Card>
    </section>
  );
}
