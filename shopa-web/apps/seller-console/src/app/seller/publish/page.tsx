"use client";

import { useEffect } from "react";
import { Button, Card, Col, Row, Tag, Typography } from "antd";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { DegradedBanner } from "@shopa/ui";
import { CatalogSection } from "@/features/catalog/CatalogSection";
import { MediaSection } from "@/features/media/MediaSection";
import { InventorySection } from "@/features/inventory/InventorySection";

const { Title, Paragraph } = Typography;

type PublishState = {
  draftId?: string;
  currentStep: "media" | "catalog" | "inventory" | "review";
  lastSavedAt?: string;
  setStep: (step: PublishState["currentStep"]) => void;
  touch: () => void;
};

const usePublishStore = create<PublishState>()(
  persist(
    (set) => ({
      currentStep: "media",
      setStep: (step) => set({ currentStep: step, lastSavedAt: new Date().toISOString() }),
      touch: () => set({ lastSavedAt: new Date().toISOString() })
    }),
    {
      name: "shopa-publish-draft"
    }
  )
);

export default function PublishPage() {
  const { currentStep, setStep, lastSavedAt, touch } = usePublishStore();

  useEffect(() => {
    if (!lastSavedAt) {
      setStep("media");
    }
  }, [lastSavedAt, setStep]);

  return (
    <div>
      <Title level={3}>Publish Product</Title>
      <Paragraph>
        Single journey for seller product publish: upload media, edit catalog draft, adjust initial
        inventory, and submit for review.
      </Paragraph>
      <DegradedBanner fields={["inventory_risk"]} />
      <Row gutter={[24, 24]}>
        <Col span={16}>
          <Card title="1. Media Assets">
            <MediaSection />
          </Card>
          <Card title="2. Catalog Draft" style={{ marginTop: 24 }}>
            <CatalogSection onSaved={() => setStep("inventory")} />
          </Card>
          <Card title="3. Initial Inventory" style={{ marginTop: 24 }}>
            <InventorySection compact />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="Publish State">
            <Tag color="blue">{currentStep}</Tag>
            <Paragraph type="secondary" style={{ marginTop: 12 }}>
              Last saved: {lastSavedAt ?? "not yet"}
            </Paragraph>
            <Button block type="primary" onClick={touch}>
              Save Draft Checkpoint
            </Button>
            <Button block style={{ marginTop: 12 }}>
              Submit for Review
            </Button>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
