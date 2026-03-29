"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import { Button, Card, Col, Collapse, Row, Space, Tabs, Tag, Typography } from "antd";
import {
  CheckCircleOutlined,
  FieldTimeOutlined,
  PlayCircleOutlined,
  QrcodeOutlined,
  RocketOutlined,
  SafetyCertificateOutlined,
  TeamOutlined
} from "@ant-design/icons";
import { useI18n } from "@shopa/ui";

const { Title, Paragraph, Text } = Typography;

type ProcessStep = {
  time: string;
  titleZh: string;
  titleEn: string;
  descZh: string;
  descEn: string;
};

type ProcessGroup = {
  key: string;
  badgeZh: string;
  badgeEn: string;
  titleZh: string;
  titleEn: string;
  descZh: string;
  descEn: string;
  steps: ProcessStep[];
};

type Industry = {
  nameZh: string;
  nameEn: string;
  targetZh: string;
  targetEn: string;
  benefitsZh: string[];
  benefitsEn: string[];
};

const processGroups: ProcessGroup[] = [
  {
    key: "personal",
    badgeZh: "个人身份",
    badgeEn: "Personal",
    titleZh: "适合无营业执照的个人卖家",
    titleEn: "For Individual Sellers",
    descZh: "用个人身份快速开店，流程更轻量。",
    descEn: "A lightweight path for individual sellers.",
    steps: [
      {
        time: "约1分钟",
        titleZh: "选择店铺类型",
        titleEn: "Choose Store Type",
        descZh: "选择个人或企业，系统会按账号认证信息校验。",
        descEn: "Choose personal or business type based on account verification."
      },
      {
        time: "实时审核",
        titleZh: "填写开店信息",
        titleEn: "Fill Store Information",
        descZh: "填写店铺名称、主体信息并上传证件材料。",
        descEn: "Complete shop profile and upload required documents."
      },
      {
        time: "实时审核",
        titleZh: "实人扫脸",
        titleEn: "Face Verification",
        descZh: "开店主体本人进行人脸识别验证。",
        descEn: "Complete face verification by the account owner."
      },
      {
        time: "完成后",
        titleZh: "发布商品",
        titleEn: "Publish Products",
        descZh: "发布商品后即可进入正式经营阶段。",
        descEn: "Publish products and start your operation."
      }
    ]
  },
  {
    key: "business",
    badgeZh: "企业身份",
    badgeEn: "Business",
    titleZh: "适合有营业执照的企业商家",
    titleEn: "For Business Sellers",
    descZh: "适合公司/个体工商户，支持更完整经营能力。",
    descEn: "For registered companies and licensed entities.",
    steps: [
      {
        time: "约1分钟",
        titleZh: "选择店铺类型",
        titleEn: "Choose Store Type",
        descZh: "根据淘宝账号绑定的认证类型自动匹配。",
        descEn: "Type is matched with account verification profile."
      },
      {
        time: "实时审核",
        titleZh: "填写主体信息",
        titleEn: "Entity Information",
        descZh: "上传法人证件、营业执照并完善经营信息。",
        descEn: "Upload legal entity documents and business profile."
      },
      {
        time: "实时审核",
        titleZh: "经营人验证",
        titleEn: "Operator Verification",
        descZh: "支持法人或实际经营人完成人脸验证。",
        descEn: "Face verification by legal person or operator."
      },
      {
        time: "完成后",
        titleZh: "发布商品",
        titleEn: "Publish Products",
        descZh: "通过后发布商品并领取新商权益。",
        descEn: "After approval, publish and claim seller incentives."
      }
    ]
  }
];

const advantageTabs = [
  {
    key: "search",
    titleZh: "淘宝搜索",
    titleEn: "Taobao Search",
    pointsZh: ["搜索冷启流量扶持", "关键词投放建议", "商品标题优化模板"],
    pointsEn: ["Cold-start traffic support", "Keyword recommendations", "Title optimization templates"]
  },
  {
    key: "home",
    titleZh: "首页推荐",
    titleEn: "Homepage Feed",
    pointsZh: ["平台内容推荐加权", "新商曝光加速", "活动场域优先入选"],
    pointsEn: ["Recommendation boost", "Exposure acceleration", "Campaign placement priority"]
  },
  {
    key: "discover",
    titleZh: "逛逛种草",
    titleEn: "Discover",
    pointsZh: ["图文短视频种草", "达人联动选品", "内容转化链路打通"],
    pointsEn: ["Short-form content", "Influencer collaboration", "Content-to-conversion flow"]
  },
  {
    key: "campaign",
    titleZh: "促销活动",
    titleEn: "Campaigns",
    pointsZh: ["大促活动报名", "新店专属激励", "经营数据复盘"],
    pointsEn: ["Campaign enrollment", "New-store incentives", "Operation analytics"]
  }
];

const industries: Industry[] = [
  {
    nameZh: "女装",
    nameEn: "Women Fashion",
    targetZh: "直播型商家 / 红人型商家 / 官方半托管",
    targetEn: "Live sellers / Influencer sellers / Managed model",
    benefitsZh: ["冷启期流量扶持", "定向商业化激励", "运费险专属优惠"],
    benefitsEn: ["Cold-start traffic", "Commercial incentives", "Shipping insurance discount"]
  },
  {
    nameZh: "美妆",
    nameEn: "Beauty",
    targetZh: "品牌渠道商家 / 人设店铺 / 直播型商家",
    targetEn: "Brand channels / Persona stores / Live sellers",
    benefitsZh: ["0元入驻", "新商免费流量包", "广告流量券"],
    benefitsEn: ["Zero-deposit onboarding", "Free traffic package", "Ad coupons"]
  },
  {
    nameZh: "个护家清",
    nameEn: "Personal & Home Care",
    targetZh: "源头工厂 / 优质组货商 / 产业带商家",
    targetEn: "Source factories / Aggregators / Industrial-belt sellers",
    benefitsZh: ["运营培训课程", "数据选品支持", "官方活动资源"],
    benefitsEn: ["Operation training", "Data-driven selection", "Official campaign resources"]
  },
  {
    nameZh: "运动户外",
    nameEn: "Sports & Outdoor",
    targetZh: "奥莱渠道商 / 原创品牌 / 产业带商家",
    targetEn: "Outlet channels / Original brands / Industrial-belt sellers",
    benefitsZh: ["0元开店", "直播激励", "免费课程"],
    benefitsEn: ["Zero-barrier onboarding", "Live incentives", "Free courses"]
  }
];

const stories = [
  {
    shop: "吉小七希乐正品小店",
    tags: "单月破百万 · 年轻人群定位",
    desc: "通过活动节奏+付费拉新，3个月完成冷启动并实现规模增长。"
  },
  {
    shop: "观匠设计家具店",
    tags: "00后创业 · 3个月破百万",
    desc: "工厂直供+内容种草，持续优化选品和店铺转化结构。"
  },
  {
    shop: "西大瓜植物商店",
    tags: "直播带货 · 年成交破百万",
    desc: "围绕直播场景搭建产品矩阵，稳定提升复购率。"
  }
];

const supportItems = [
  {
    icon: <RocketOutlined />,
    titleZh: "新商成长计划",
    titleEn: "Seller Growth Program",
    descZh: "入驻前90天提供经营路径、选品与投放建议。",
    descEn: "90-day playbook for operation, products and traffic."
  },
  {
    icon: <SafetyCertificateOutlined />,
    titleZh: "交易与风控保障",
    titleEn: "Transaction Protection",
    descZh: "完善的交易链路和风控策略，保障店铺稳定经营。",
    descEn: "Stable transaction and risk-control infrastructure."
  },
  {
    icon: <TeamOutlined />,
    titleZh: "客服与培训体系",
    titleEn: "Support & Training",
    descZh: "覆盖入驻、发货、售后全流程问题处理。",
    descEn: "End-to-end support from onboarding to after-sales."
  }
];

const faqPanels = [
  {
    key: "1",
    zh: "没有营业执照可以开店吗？",
    en: "Can I open a store without a business license?",
    ansZh: "可以，支持个人身份开店；若为企业经营，建议使用企业身份入驻。",
    ansEn: "Yes. Personal onboarding is supported; business entities should use business onboarding."
  },
  {
    key: "2",
    zh: "审核一般多久有结果？",
    en: "How long does review take?",
    ansZh: "多数资料可实时校验，复杂场景通常在24小时内反馈。",
    ansEn: "Most checks are real-time; complex cases typically return within 24 hours."
  },
  {
    key: "3",
    zh: "审核通过后如何进入后台？",
    en: "How do I access the seller admin after approval?",
    ansZh: "审核通过后系统会自动跳转到商家后台工作台。",
    ansEn: "After approval, the system automatically routes you to seller workbench."
  }
];

export default function SellerOpenShopPortalPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [activeProcess, setActiveProcess] = useState("personal");

  const currentProcess = useMemo(() => processGroups.find((item) => item.key === activeProcess) ?? processGroups[0], [activeProcess]);

  return (
    <div className="tb-open-portal">
      <section className="tb-open-topbar">
        <div className="tb-open-container tb-open-topbar-inner">
          <span>{isZh ? "欢迎来到 SHOPA 开店" : "Welcome to SHOPA Open-Shop"}</span>
          <Space size={16}>
            <a href="#process">{isZh ? "入驻流程" : "Process"}</a>
            <a href="#industry">{isZh ? "招商类目" : "Industries"}</a>
            <a href="#faq">FAQ</a>
          </Space>
        </div>
      </section>

      <section className="tb-open-nav">
        <div className="tb-open-container tb-open-nav-inner">
          <div className="tb-open-logo">淘宝开店 · SHOPA</div>
          <Space size={20} className="tb-open-menu">
            <a href="#advantage">{isZh ? "平台优势" : "Advantages"}</a>
            <a href="#stories">{isZh ? "商家案例" : "Stories"}</a>
            <a href="#support">{isZh ? "服务保障" : "Support"}</a>
          </Space>
          <Link href="/login">
            <Button type="primary" size="large">
              {isZh ? "我要开店" : "Open Shop"}
            </Button>
          </Link>
        </div>
      </section>

      <section className="tb-open-hero-v2">
        <div className="tb-open-container tb-open-hero-grid">
          <div className="tb-open-hero-main">
            <Tag color="gold">{isZh ? "淘宝风格开店门户" : "Taobao-Style Portal"}</Tag>
            <Title level={1}>{isZh ? "一站式开店，经营更简单" : "One-Stop Onboarding, Easy Operations"}</Title>
            <Paragraph>
              {isZh
                ? "先浏览入驻规则和平台能力，再点击“我要开店”进入登录/注册。登录后先提交证件资料，审核通过即可进入商家后台管理。"
                : "Review onboarding info first, then click Open Shop to login/register. After login, upload documents; once approved, enter seller admin."}
            </Paragraph>
            <Space wrap>
              <Link href="/login">
                <Button type="primary" size="large">
                  {isZh ? "我要开店" : "Open My Shop"}
                </Button>
              </Link>
              <Link href="/register">
                <Button size="large">{isZh ? "新商注册" : "Seller Register"}</Button>
              </Link>
            </Space>
            <div className="tb-open-note">
              <CheckCircleOutlined />
              <span>{isZh ? "审核通过后自动进入商家后台工作台" : "Auto route to seller workbench after approval"}</span>
            </div>
          </div>

          <Card className="tb-open-hero-side" bordered={false}>
            <div className="tb-open-side-header">
              <PlayCircleOutlined />
              <strong>{isZh ? "开店数据看板" : "Onboarding KPI"}</strong>
            </div>
            <div className="tb-open-kpi-grid">
              <div>
                <strong>4</strong>
                <span>{isZh ? "步完成开店" : "4 steps"}</span>
              </div>
              <div>
                <strong>24h</strong>
                <span>{isZh ? "审核反馈" : "Review SLA"}</span>
              </div>
              <div>
                <strong>100%</strong>
                <span>{isZh ? "在线化流程" : "Online flow"}</span>
              </div>
              <div>
                <strong>7x12</strong>
                <span>{isZh ? "商家服务" : "Seller support"}</span>
              </div>
            </div>
          </Card>
        </div>
      </section>

      <section className="tb-open-subnav-v2">
        <div className="tb-open-container">
          <Space wrap size={12}>
            <Tag color="orange">{isZh ? "新商权益" : "New Seller Benefits"}</Tag>
            <Tag color="orange">{isZh ? "免费课程" : "Free Courses"}</Tag>
            <Tag color="orange">{isZh ? "平台扶持" : "Platform Support"}</Tag>
            <Tag color="orange">{isZh ? "直播激励" : "Live Incentives"}</Tag>
          </Space>
        </div>
      </section>

      <section id="process" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>{isZh ? "入驻流程" : "Onboarding Process"}</Title>
            <Text type="secondary">{isZh ? "个人/企业双路径，按你的经营类型选择。" : "Personal and business paths available."}</Text>
          </header>

          <Tabs
            activeKey={activeProcess}
            onChange={setActiveProcess}
            items={processGroups.map((group) => ({
              key: group.key,
              label: isZh ? group.badgeZh : group.badgeEn
            }))}
          />

          <Card className="tb-process-card" bordered={false}>
            <div className="tb-process-title-row">
              <Tag color="orange">{isZh ? currentProcess.badgeZh : currentProcess.badgeEn}</Tag>
              <h4>{isZh ? currentProcess.titleZh : currentProcess.titleEn}</h4>
              <p>{isZh ? currentProcess.descZh : currentProcess.descEn}</p>
            </div>
            <div className="tb-process-steps">
              {currentProcess.steps.map((step, index) => (
                <article key={`${currentProcess.key}-${step.titleZh}`}>
                  <div className="tb-process-step-index">{String(index + 1).padStart(2, "0")}</div>
                  <div className="tb-process-step-time">
                    <FieldTimeOutlined /> {step.time}
                  </div>
                  <h5>{isZh ? step.titleZh : step.titleEn}</h5>
                  <p>{isZh ? step.descZh : step.descEn}</p>
                </article>
              ))}
            </div>
          </Card>
        </div>
      </section>

      <section id="advantage" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>{isZh ? "平台优势" : "Platform Advantages"}</Title>
          </header>

          <Tabs
            items={advantageTabs.map((tab) => ({
              key: tab.key,
              label: isZh ? tab.titleZh : tab.titleEn,
              children: (
                <div className="tb-advantage-panel">
                  <div className="tb-advantage-preview" />
                  <ul>
                    {(isZh ? tab.pointsZh : tab.pointsEn).map((point) => (
                      <li key={point}>{point}</li>
                    ))}
                  </ul>
                </div>
              )
            }))}
          />
        </div>
      </section>

      <section id="industry" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>{isZh ? "热门招商行业" : "Hot Recruiting Industries"}</Title>
            <Text type="secondary">{isZh ? "查询入驻材料及费用，匹配你的商家类型。" : "Find requirements and costs by category."}</Text>
          </header>
          <Row gutter={[12, 12]}>
            {industries.map((item) => (
              <Col xs={24} md={12} key={item.nameZh}>
                <Card className="tb-industry-card" bordered={false}>
                  <h4>{isZh ? item.nameZh : item.nameEn}</h4>
                  <p>{isZh ? item.targetZh : item.targetEn}</p>
                  <div className="tb-industry-benefits">
                    {(isZh ? item.benefitsZh : item.benefitsEn).map((benefit) => (
                      <Tag key={benefit}>{benefit}</Tag>
                    ))}
                  </div>
                  <Link href="/login" className="tb-industry-link">
                    {isZh ? "查询入驻材料及费用 >" : "Check requirements and fees >"}
                  </Link>
                </Card>
              </Col>
            ))}
          </Row>
        </div>
      </section>

      <section id="stories" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>{isZh ? "商家案例" : "Merchant Stories"}</Title>
          </header>
          <div className="tb-story-grid">
            {stories.map((item) => (
              <Card key={item.shop} className="tb-story-card" bordered={false}>
                <h4>{item.shop}</h4>
                <Tag color="gold">{item.tags}</Tag>
                <p>{item.desc}</p>
              </Card>
            ))}
          </div>
        </div>
      </section>

      <section id="support" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>{isZh ? "平台支持" : "Platform Support"}</Title>
          </header>
          <Row gutter={[12, 12]}>
            {supportItems.map((item) => (
              <Col xs={24} md={8} key={item.titleZh}>
                <Card bordered={false} className="tb-support-card-v2">
                  <div className="tb-support-icon-v2">{item.icon}</div>
                  <h4>{isZh ? item.titleZh : item.titleEn}</h4>
                  <p>{isZh ? item.descZh : item.descEn}</p>
                </Card>
              </Col>
            ))}
          </Row>
        </div>
      </section>

      <section id="faq" className="tb-open-section-v2">
        <div className="tb-open-container">
          <header className="tb-open-head-v2">
            <Title level={3}>FAQ</Title>
          </header>
          <Collapse
            items={faqPanels.map((item) => ({
              key: item.key,
              label: isZh ? item.zh : item.en,
              children: isZh ? item.ansZh : item.ansEn
            }))}
          />
        </div>
      </section>

      <section className="tb-open-section-v2">
        <div className="tb-open-container tb-qr-panel">
          <div>
            <Title level={4}>{isZh ? "商家工作台" : "Seller Workbench"}</Title>
            <Paragraph>{isZh ? "登录后先上传入驻证件，审核通过即可进入后台管理商品、库存、店铺。" : "Upload onboarding documents after login, then enter admin after approval."}</Paragraph>
            <Space wrap>
              <Link href="/login">
                <Button type="primary" size="large">
                  {isZh ? "我要开店" : "Open Shop"}
                </Button>
              </Link>
              <Link href="/register">
                <Button size="large">{isZh ? "注册账号" : "Create Account"}</Button>
              </Link>
            </Space>
          </div>
          <div className="tb-qr-box-v2" aria-label="seller-workbench-qr">
            <QrcodeOutlined />
            <span>{isZh ? "淘宝/钉钉扫码" : "Scan with app"}</span>
          </div>
        </div>
      </section>
    </div>
  );
}
