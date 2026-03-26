"use client";

import {Avatar, Badge, Card, Col, Divider, Row, Space, Table, Tag, Typography} from "antd";
import {useQuery} from "@tanstack/react-query";
import {DegradedBanner} from "@shopa/ui";
import {fetchSellerWorkbench} from "./api";

const metrics: Array<[keyof {total: number; onShelf: number; offShelf: number; reviewing: number; draft: number; rejected: number}, string]> = [
  ["total", "Total Products"],
  ["onShelf", "On Shelf"],
  ["offShelf", "Off Shelf"],
  ["reviewing", "Reviewing"],
  ["draft", "Drafts"],
  ["rejected", "Rejected"]
];

export function WorkbenchPanel() {
  const {data, isLoading} = useQuery({
    queryKey: ["seller-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30 * 1000
  });

  return (
    <Space direction="vertical" size="large" style={{width: "100%"}}>
      {data?.partial && <DegradedBanner fields={data.degradedFields ?? []} />}
      <Row gutter={[16, 16]}>
        {metrics.map(([key, label]) => (
          <Col key={key} span={4}>
            <Card
              loading={isLoading}
              bodyStyle={{background: "linear-gradient(135deg,#001529,#002f6c)", color: "#fff"}}
            >
              <Typography.Text style={{color: "#8fb5ff"}}>{label}</Typography.Text>
              <Typography.Title level={3} style={{color: "#fff", marginTop: 8}}>
                {data?.productSummary?.[key as keyof typeof data["productSummary"]] ?? 0}
              </Typography.Title>
            </Card>
          </Col>
        ))}
      </Row>
      <Divider />
      <Row gutter={16}>
        <Col span={12}>
          <Card title="My Shops" loading={isLoading} bodyStyle={{padding: 0}}>
            <Table
              dataSource={data?.shops ?? []}
              rowKey="shopNo"
              pagination={false}
              size="small"
              columns={[
                {
                  title: "Shop",
                  dataIndex: "shopDisplayName",
                  render: (_: any, record: any) => (
                    <Space>
                      <Avatar size="small" shape="square">
                        {record.shopDisplayName.slice(0, 1)}
                      </Avatar>
                      {record.shopDisplayName}
                    </Space>
                  )
                },
                {title: "Status", dataIndex: "shopStatusCode", render: (status: string) => <Badge status="processing" text={status} />},
                {
                  title: "Visible",
                  dataIndex: "buyerVisible",
                  render: (visible: boolean) => <Badge status={visible ? "success" : "default"} text={visible ? "Live" : "Hidden"} />
                }
              ]}
            />
            {data?.shopsTruncated && <Typography.Text type="secondary" style={{display: "block", padding: 12}}>Showing 10 most recent shops.</Typography.Text>}
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Recent Applications" loading={isLoading} bodyStyle={{padding: 0}}>
            <Table
              dataSource={data?.latestApplications ?? []}
              rowKey="applicationNo"
              pagination={false}
              size="small"
              columns={[
                {title: "Application", dataIndex: "applicationNo"},
                {title: "Status", dataIndex: "applicationStatusCode", render: (value: string) => <Tag color="purple">{value}</Tag>},
                {title: "Updated", dataIndex: "updatedAt"}
              ]}
            />
            {data?.applicationsTruncated && <Typography.Text type="secondary" style={{display: "block", padding: 12}}>Only latest five applications are shown.</Typography.Text>}
          </Card>
        </Col>
      </Row>
    </Space>
  );
}
