import { apiClient } from "@shopa/api-client";
import type { Address, Profile } from "./types";

type PbImageRef = { url?: string };
type PbUserProfile = {
  display_name: string;
  updated_at?: string;
  avatar?: PbImageRef;
  profile_version?: number;
  ext?: Record<string, string>;
};

type PbAddress = {
  address_id: number;
  label?: string;
  receiver_name?: string;
  receiver_phone?: string;
  street?: string;
  detail?: string;
  is_default?: boolean;
  updated_at?: string;
  address_version?: number;
};

type FetchProfileResponse = {
  profile: PbUserProfile;
  addresses?: PbAddress[];
  address_book_version?: number;
};

type FetchAddressesResponse = {
  addresses: PbAddress[];
  address_book_version?: number;
};

type CreateAddressResponse = {
  address: PbAddress;
  address_book_version?: number;
};

type UpdateAddressResponse = {
  address: PbAddress;
  address_book_version?: number;
};

let cachedProfile: Profile | null = null;
let cachedAddresses: Address[] = [];
let profileVersion: number | null = null;
let addressBookVersion: number | null = null;

const toProfile = (pb: PbUserProfile, bookVersion?: number): Profile => {
  const bio = pb.ext?.bio;
  profileVersion = pb.profile_version ?? profileVersion;
  addressBookVersion = bookVersion ?? addressBookVersion;
  const profile: Profile = {
    displayName: pb.display_name,
    email: pb.ext?.email,
    phone: pb.ext?.phone,
    bio,
    avatarUrl: pb.avatar?.url,
    updatedAt: pb.updated_at,
    version: pb.profile_version,
    addressBookVersion: bookVersion ?? addressBookVersion ?? undefined
  };
  cachedProfile = profile;
  return profile;
};

const toAddress = (pb: PbAddress, bookVersion?: number): Address => {
  const parts = [pb.street, pb.detail].filter(Boolean);
  return {
    id: String(pb.address_id),
    label: pb.label ?? "Address",
    contact: pb.receiver_name ?? "",
    phone: pb.receiver_phone ?? "",
    fullAddress: parts.join(" ").trim() || "Address",
    isDefault: Boolean(pb.is_default),
    updatedAt: pb.updated_at,
    version: pb.address_version,
    addressBookVersion: bookVersion ?? addressBookVersion ?? undefined
  };
};

const buildAddressCreate = (address: Partial<Address>) => ({
  label: address.label,
  receiver_name: address.contact,
  receiver_phone: address.phone,
  street: address.fullAddress,
  detail: "",
  postal_code: "",
  latitude: 0,
  longitude: 0
});

const buildAddressPatch = (
  address: Partial<Address>
): { patch: Record<string, unknown>; mask: string[] } => {
  const patch: Record<string, unknown> = {};
  const mask: string[] = [];
  if (address.label) {
    patch.label = address.label;
    mask.push("label");
  }
  if (address.contact) {
    patch.receiver_name = address.contact;
    mask.push("receiver_name");
  }
  if (address.phone) {
    patch.receiver_phone = address.phone;
    mask.push("receiver_phone");
  }
  if (address.fullAddress) {
    patch.street = address.fullAddress;
    mask.push("street");
  }
  if (mask.length === 0) {
    mask.push("street");
    patch.street = address.fullAddress;
  }
  return { patch, mask };
};

const ensureAddressBookVersion = (): number => addressBookVersion ?? 0;

export async function fetchProfile(): Promise<Profile> {
  try {
    const res = await apiClient.get<FetchProfileResponse>("/v1/me/profile?include_addresses=true");
    if (res.profile) {
      toProfile(res.profile, res.address_book_version);
    }
    if (Array.isArray(res.addresses)) {
      cachedAddresses = res.addresses.map((address) => toAddress(address, res.address_book_version));
    }
    return cachedProfile ?? toProfile(res.profile, res.address_book_version);
  } catch (error) {
    if (cachedProfile) {
      return cachedProfile;
    }
    throw error;
  }
}

export async function saveProfile(payload: Partial<Profile>): Promise<Profile> {
  if (!cachedProfile) {
    await fetchProfile();
  }
  const updateMask: string[] = [];
  const profileBody: Record<string, unknown> = {};
  if (payload.displayName) {
    profileBody.display_name = payload.displayName;
    updateMask.push("display_name");
  }
  const ext: Record<string, string> = {};
  if (payload.bio) {
    ext.bio = payload.bio;
  }
  if (payload.email) {
    ext.email = payload.email;
  }
  if (payload.phone) {
    ext.phone = payload.phone;
  }
  if (Object.keys(ext).length > 0) {
    profileBody.ext = ext;
    updateMask.push("ext");
  }
  if (updateMask.length === 0) {
    throw new Error("At least one profile field must be provided.");
  }
  const patch = {
    profile: profileBody,
    update_mask: updateMask,
    expected_profile_version: profileVersion ?? 0
  };
  const res = await apiClient.patch<{ profile: PbUserProfile }>("/v1/me/profile", patch);
  return toProfile(res.profile, addressBookVersion ?? undefined);
}

export async function fetchAddresses(): Promise<Address[]> {
  try {
    const res = await apiClient.get<FetchAddressesResponse>("/v1/me/addresses");
    addressBookVersion = res.address_book_version ?? addressBookVersion;
    cachedAddresses = res.addresses.map((address) => toAddress(address, res.address_book_version));
    return cachedAddresses;
  } catch (error) {
    if (cachedAddresses.length > 0) {
      return cachedAddresses;
    }
    throw error;
  }
}

export async function saveAddress(address: Partial<Address>): Promise<Address> {
  if (!address.fullAddress || !address.contact) {
    throw new Error("Contact and full address are required.");
  }
  if (address.id) {
    const { patch, mask } = buildAddressPatch(address);
    const res = await apiClient.patch<UpdateAddressResponse>(`/v1/me/addresses/${address.id}`, {
      address: patch,
      update_mask: mask,
      expected_address_version: address.version ?? ensureAddressBookVersion()
    });
    addressBookVersion = res.address_book_version ?? addressBookVersion;
    const updated = toAddress(res.address, res.address_book_version ?? addressBookVersion ?? undefined);
    cachedAddresses = cachedAddresses.map((item) => (item.id === updated.id ? updated : item));
    return updated;
  }
  const res = await apiClient.post<CreateAddressResponse>("/v1/me/addresses", {
    address: buildAddressCreate(address),
    set_as_default: Boolean(address.isDefault),
    expected_address_book_version: ensureAddressBookVersion()
  });
  addressBookVersion = res.address_book_version ?? addressBookVersion;
  const created = toAddress(res.address, res.address_book_version ?? addressBookVersion ?? undefined);
  cachedAddresses = [...cachedAddresses, created];
  return created;
}

export async function setDefaultAddress(addressId?: string): Promise<void> {
  const res = await apiClient.post<{ default_address_id?: number; address_book_version?: number }>("/v1/me/addresses/default", {
    address_id: addressId ? Number(addressId) : undefined,
    expected_address_book_version: ensureAddressBookVersion()
  });
  addressBookVersion = res.address_book_version ?? addressBookVersion;
}
