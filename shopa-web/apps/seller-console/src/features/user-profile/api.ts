import { apiClient } from "@shopa/api-client";
import { Address, Profile } from "./types";

const fallbackProfile: Profile = {
  displayName: "Shopa Seller",
  email: "seller@example.com",
  phone: "13800138000",
  bio: "Focus on high quality products and stable service.",
  updatedAt: undefined
};

const fallbackAddresses: Address[] = [
  {
    id: "local-1",
    contact: "Default Receiver",
    phone: "13800138000",
    label: "Home",
    fullAddress: "No.1 Zhongguancun Street, Beijing",
    isDefault: true,
    updatedAt: undefined
  },
  {
    id: "local-2",
    contact: "Warehouse",
    phone: "01088880000",
    label: "Warehouse",
    fullAddress: "Airport Highway Side, Shunyi, Beijing",
    isDefault: false,
    updatedAt: undefined
  }
];

export async function fetchProfile(): Promise<Profile> {
  try {
    return await apiClient.get<Profile>("/v1/me/profile");
  } catch (error) {
    console.warn("profile fallback", error);
    return fallbackProfile;
  }
}

export async function saveProfile(payload: Partial<Profile>): Promise<Profile> {
  try {
    return await apiClient.put<Profile>("/v1/me/profile", payload);
  } catch (error) {
    console.warn("profile save fallback", error);
    return { ...fallbackProfile, ...payload };
  }
}

export async function fetchAddresses(): Promise<Address[]> {
  try {
    return await apiClient.get<Address[]>("/v1/me/addresses");
  } catch (error) {
    console.warn("address fallback", error);
    return fallbackAddresses;
  }
}

export async function saveAddress(address: Partial<Address>): Promise<Address> {
  if (!address.fullAddress || !address.contact) {
    throw new Error("Contact and full address are required.");
  }

  if (address.id) {
    await apiClient.put(`/v1/me/addresses/${address.id}`, address);
    return {
      id: address.id,
      label: address.label ?? "Address",
      contact: address.contact,
      phone: address.phone ?? "",
      fullAddress: address.fullAddress,
      isDefault: Boolean(address.isDefault),
      updatedAt: new Date().toISOString()
    };
  }

  await apiClient.post("/v1/me/addresses", address);
  return {
    id: `local-${Date.now()}`,
    label: address.label ?? "Address",
    contact: address.contact,
    phone: address.phone ?? "",
    fullAddress: address.fullAddress,
    isDefault: false,
    updatedAt: new Date().toISOString()
  };
}

export async function setDefaultAddress(addressId: string): Promise<void> {
  try {
    await apiClient.put(`/v1/me/addresses/${addressId}/default`, {});
  } catch (error) {
    console.warn("set default fallback", error);
  }
}
