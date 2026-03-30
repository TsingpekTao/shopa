import { apiClient } from "@/lib/api-client";
import { AddressBook, AddressFormInput, UserAddress } from "./types";

type RawUserAddress = {
  address_id?: number | string;
  addressId?: number | string;
  address_version?: number | string;
  addressVersion?: number | string;
  label?: string;
  receiver_name?: string;
  receiverName?: string;
  receiver_phone?: string;
  receiverPhone?: string;
  country_code?: string;
  countryCode?: string;
  province_code?: string;
  provinceCode?: string;
  province_name?: string;
  provinceName?: string;
  city_code?: string;
  cityCode?: string;
  city_name?: string;
  cityName?: string;
  district_code?: string;
  districtCode?: string;
  district_name?: string;
  districtName?: string;
  street?: string;
  detail?: string;
  postal_code?: string;
  postalCode?: string;
  is_default?: boolean;
  isDefault?: boolean;
  status?: number | string;
};

type RawListMyAddressesRes = {
  addresses?: RawUserAddress[];
  address_book_version?: number | string;
  addressBookVersion?: number | string;
};

type RawCreateMyAddressRes = {
  address?: RawUserAddress;
  address_book_version?: number | string;
  addressBookVersion?: number | string;
};

type RawUpdateMyAddressRes = RawCreateMyAddressRes;
type RawDeleteMyAddressRes = {
  address_book_version?: number | string;
  addressBookVersion?: number | string;
};
type RawSetDefaultAddressRes = {
  default_address_id?: number | string;
  defaultAddressId?: number | string;
  address_book_version?: number | string;
  addressBookVersion?: number | string;
};

function toNumber(value: unknown, fallback = 0): number {
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function normalizeAddress(raw: RawUserAddress): UserAddress {
  return {
    addressId: toNumber(raw.address_id ?? raw.addressId, 0),
    addressVersion: toNumber(raw.address_version ?? raw.addressVersion, 0),
    label: toString(raw.label),
    receiverName: toString(raw.receiver_name ?? raw.receiverName),
    receiverPhone: toString(raw.receiver_phone ?? raw.receiverPhone),
    countryCode: toString(raw.country_code ?? raw.countryCode),
    provinceCode: toString(raw.province_code ?? raw.provinceCode),
    provinceName: toString(raw.province_name ?? raw.provinceName),
    cityCode: toString(raw.city_code ?? raw.cityCode),
    cityName: toString(raw.city_name ?? raw.cityName),
    districtCode: toString(raw.district_code ?? raw.districtCode),
    districtName: toString(raw.district_name ?? raw.districtName),
    street: toString(raw.street),
    detail: toString(raw.detail),
    postalCode: toString(raw.postal_code ?? raw.postalCode),
    isDefault: Boolean(raw.is_default ?? raw.isDefault),
    status: toNumber(raw.status, 0)
  };
}

function buildAddressCreatePayload(input: AddressFormInput) {
  return {
    label: input.label.trim(),
    receiver_name: input.receiverName.trim(),
    receiver_phone: input.receiverPhone.trim(),
    country_code: "CN",
    province_code: input.provinceCode.trim(),
    province_name: input.provinceName.trim(),
    city_code: input.cityCode.trim(),
    city_name: input.cityName.trim(),
    district_code: input.districtCode.trim(),
    district_name: input.districtName.trim(),
    street: input.street.trim(),
    detail: input.detail.trim(),
    postal_code: input.postalCode.trim()
  };
}

function buildAddressPatchPayload(input: AddressFormInput) {
  return {
    label: input.label.trim(),
    receiver_name: input.receiverName.trim(),
    receiver_phone: input.receiverPhone.trim(),
    country_code: "CN",
    province_code: input.provinceCode.trim(),
    province_name: input.provinceName.trim(),
    city_code: input.cityCode.trim(),
    city_name: input.cityName.trim(),
    district_code: input.districtCode.trim(),
    district_name: input.districtName.trim(),
    street: input.street.trim(),
    detail: input.detail.trim(),
    postal_code: input.postalCode.trim()
  };
}

function addressPatchMask(): string[] {
  return [
    "label",
    "receiver_name",
    "receiver_phone",
    "country_code",
    "province_code",
    "province_name",
    "city_code",
    "city_name",
    "district_code",
    "district_name",
    "street",
    "detail",
    "postal_code"
  ];
}

export async function listMyAddresses(includeDeleted = false): Promise<AddressBook> {
  const response = await apiClient.get<RawListMyAddressesRes>("/v1/me/addresses", {
    params: { include_deleted: includeDeleted }
  });
  return {
    addresses: Array.isArray(response?.addresses) ? response.addresses.map((item) => normalizeAddress(item)) : [],
    addressBookVersion: toNumber(response?.address_book_version ?? response?.addressBookVersion, 0)
  };
}

export async function createMyAddress(input: AddressFormInput, addressBookVersion: number): Promise<{ address: UserAddress; addressBookVersion: number }> {
  const response = await apiClient.post<RawCreateMyAddressRes>("/v1/me/addresses", {
    address: buildAddressCreatePayload(input),
    set_as_default: Boolean(input.setAsDefault),
    expected_address_book_version: addressBookVersion
  });
  return {
    address: normalizeAddress(response?.address ?? {}),
    addressBookVersion: toNumber(response?.address_book_version ?? response?.addressBookVersion, 0)
  };
}

export async function updateMyAddress(addressId: number, input: AddressFormInput, expectedAddressVersion: number): Promise<{ address: UserAddress; addressBookVersion: number }> {
  const mask = addressPatchMask();
  const response = await apiClient.patch<RawUpdateMyAddressRes>(`/v1/me/addresses/${addressId}`, {
    address: buildAddressPatchPayload(input),
    update_mask: mask,
    expected_address_version: expectedAddressVersion
  });
  return {
    address: normalizeAddress(response?.address ?? {}),
    addressBookVersion: toNumber(response?.address_book_version ?? response?.addressBookVersion, 0)
  };
}

export async function deleteMyAddress(addressId: number, expectedAddressVersion: number): Promise<{ addressBookVersion: number }> {
  const response = await apiClient.delete<RawDeleteMyAddressRes>(`/v1/me/addresses/${addressId}`, {
    params: {
      expected_address_version: expectedAddressVersion
    },
    data: {
      expected_address_version: expectedAddressVersion
    }
  });
  return {
    addressBookVersion: toNumber(response?.address_book_version ?? response?.addressBookVersion, 0)
  };
}

export async function setMyDefaultAddress(addressId: number, expectedAddressBookVersion: number): Promise<{ defaultAddressId: number; addressBookVersion: number }> {
  const response = await apiClient.post<RawSetDefaultAddressRes>("/v1/me/addresses/default", {
    address_id: addressId,
    expected_address_book_version: expectedAddressBookVersion
  });
  return {
    defaultAddressId: toNumber(response?.default_address_id ?? response?.defaultAddressId, 0),
    addressBookVersion: toNumber(response?.address_book_version ?? response?.addressBookVersion, 0)
  };
}
