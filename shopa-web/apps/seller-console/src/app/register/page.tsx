"use client";

import Link from "next/link";
import { useState } from "react";
import { Button, Card, Form, Input, Typography, notification } from "antd";
import { LockOutlined, MobileOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { registerByPassword, sendSmsCode } from "@/features/iam/api";

const { Title, Text } = Typography;

export default function RegisterPage() {
  const [sendingCode, setSendingCode] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [form] = Form.useForm<{ phone: string; smsCode: string; password: string; confirmPassword: string }>();

  const handleSendCode = async () => {
    const phone = form.getFieldValue("phone");
    if (!phone) {
      notification.warning({ message: "请先输入手机号" });
      return;
    }
    setSendingCode(true);
    try {
      const data = await sendSmsCode({ scene: 1, phone });
      const next = data?.resendAfterSeconds ?? 60;
      setCountdown(next);
      notification.success({ message: "验证码已发送" });
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
      notification.error({ message: "验证码发送失败，请稍后重试" });
    } finally {
      setSendingCode(false);
    }
  };

  const handleRegister = async (values: {
    phone: string;
    smsCode: string;
    password: string;
    confirmPassword: string;
  }) => {
    setSubmitting(true);
    try {
      await registerByPassword(values);
      notification.success({ message: "注册成功，请登录" });
      window.location.href = "/login";
    } catch (_err) {
      notification.error({ message: "注册失败，请检查信息后重试" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="tb-register-wrap">
      <Card className="tb-register-card" bordered={false}>
        <Title level={4} style={{ marginBottom: 8 }}>
          注册商家账号
        </Title>
        <Text type="secondary">完成注册后即可登录卖家中心。</Text>

        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 18 }}
          onFinish={handleRegister}
          autoComplete="off"
        >
          <Form.Item name="phone" label="手机号" rules={[{ required: true, message: "请输入手机号" }]}>
            <Input placeholder="请输入手机号" prefix={<MobileOutlined />} />
          </Form.Item>

          <Form.Item label="短信验证码" required>
            <div className="tb-register-code-row">
              <Form.Item
                name="smsCode"
                noStyle
                rules={[{ required: true, message: "请输入短信验证码" }]}
              >
                <Input placeholder="请输入验证码" prefix={<SafetyCertificateOutlined />} />
              </Form.Item>
              <Button onClick={handleSendCode} disabled={countdown > 0} loading={sendingCode}>
                {countdown > 0 ? `${countdown}s后重发` : "发送验证码"}
              </Button>
            </div>
          </Form.Item>

          <Form.Item
            name="password"
            label="登录密码"
            rules={[{ required: true, message: "请输入登录密码" }, { min: 8, message: "密码至少 8 位" }]}
          >
            <Input.Password placeholder="请输入登录密码" prefix={<LockOutlined />} />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认密码"
            dependencies={["password"]}
            rules={[
              { required: true, message: "请再次输入密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的密码不一致"));
                }
              })
            ]}
          >
            <Input.Password placeholder="请再次输入密码" prefix={<LockOutlined />} />
          </Form.Item>

          <Button type="primary" htmlType="submit" block loading={submitting}>
            立即注册
          </Button>
        </Form>

        <div className="tb-register-hint">
          <Text type="secondary">已有账号？</Text>
          <Link href="/login">去登录</Link>
        </div>
      </Card>
    </section>
  );
}
