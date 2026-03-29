"use client";

import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { SellerApplicationDetail } from "@/features/seller-shop/types";

type OnboardingState = {
  applicationNo: string;
  expectedVersion: number;
  statusCode: string;
  entityName: string;
  contactName: string;
  contactPhone: string;
  contactEmail: string;
  shopName: string;
  shopDisplayName: string;
  servicePhone: string;
  description: string;
  idCardFrontAssetId: string;
  idCardBackAssetId: string;
  businessLicenseAssetId: string;
  setField: <K extends keyof OnboardingState>(key: K, value: OnboardingState[K]) => void;
  applyServerDraft: (draft: SellerApplicationDetail) => void;
  clearDraft: () => void;
};

const initialState = {
  applicationNo: "",
  expectedVersion: 0,
  statusCode: "",
  entityName: "",
  contactName: "",
  contactPhone: "",
  contactEmail: "",
  shopName: "",
  shopDisplayName: "",
  servicePhone: "",
  description: "",
  idCardFrontAssetId: "",
  idCardBackAssetId: "",
  businessLicenseAssetId: ""
};

const storage = typeof window !== "undefined" ? createJSONStorage(() => localStorage) : undefined;

export const useOnboardingStore = create<OnboardingState>()(
  persist(
    (set) => ({
      ...initialState,
      setField: (key, value) => set((state) => ({ ...state, [key]: value })),
      applyServerDraft: (draft) =>
        set((state) => ({
          ...state,
          applicationNo: draft.applicationNo || state.applicationNo,
          expectedVersion: draft.version || state.expectedVersion,
          statusCode: draft.applicationStatusCode || state.statusCode,
          entityName: draft.entity?.entityName || state.entityName,
          contactName: draft.entity?.contactName || state.contactName,
          contactPhone: draft.entity?.contactPhone || state.contactPhone,
          contactEmail: draft.entity?.contactEmail || state.contactEmail,
          shopName: draft.shop?.shopName || state.shopName,
          shopDisplayName: draft.shop?.shopDisplayName || state.shopDisplayName,
          servicePhone: draft.shop?.servicePhone || state.servicePhone,
          description: draft.shop?.description || state.description,
          idCardFrontAssetId: draft.entity?.ext?.id_card_front_asset_id || state.idCardFrontAssetId,
          idCardBackAssetId: draft.entity?.ext?.id_card_back_asset_id || state.idCardBackAssetId,
          businessLicenseAssetId: draft.entity?.ext?.business_license_asset_id || state.businessLicenseAssetId
        })),
      clearDraft: () => set({ ...initialState })
    }),
    {
      name: "shopa_seller_onboarding_draft_v1",
      storage
    }
  )
);
