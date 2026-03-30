"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { useMutation, useQuery } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { createMyAddress, listMyAddresses, setMyDefaultAddress, updateMyAddress } from "@/features/address/api";
import { AddressFormInput } from "@/features/address/types";
import {
  findCityByCode,
  findCityCodeByName,
  findDistrictByCode,
  findDistrictCodeByName,
  findProvinceByCode,
  findProvinceCodeByName,
  getCityOptions,
  getDistrictOptions,
  getProvinceOptions
} from "@/features/address/region-data";

type AddressEditorMode = "create" | "edit";

const EMPTY_FORM: AddressFormInput = {
  label: "",
  receiverName: "",
  receiverPhone: "",
  provinceCode: "",
  provinceName: "",
  cityCode: "",
  cityName: "",
  districtCode: "",
  districtName: "",
  street: "",
  detail: "",
  postalCode: "",
  setAsDefault: true
};

function isValidPhone(value: string): boolean {
  return /^1\d{10}$/.test(value.trim());
}

function toErrorText(error: unknown, zhText: string, enText: string, isZh: boolean): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return isZh ? zhText : enText;
}

function ensureOption(options: { code: string; name: string }[], code: string, name: string): { code: string; name: string }[] {
  if (!code || !name) {
    return options;
  }
  if (options.some((item) => item.code === code)) {
    return options;
  }
  return [{ code, name }, ...options];
}

export function AddressEditor({ mode, addressId }: { mode: AddressEditorMode; addressId?: number }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const router = useRouter();
  const [form, setForm] = useState<AddressFormInput>(EMPTY_FORM);

  const addressQuery = useQuery({
    queryKey: ["my-addresses", "editor"],
    queryFn: () => listMyAddresses(false),
    staleTime: 15_000,
    refetchOnWindowFocus: false
  });

  const addresses = addressQuery.data?.addresses ?? [];
  const addressBookVersion = addressQuery.data?.addressBookVersion ?? 0;

  const editingAddress = useMemo(() => {
    if (mode !== "edit" || !addressId) {
      return null;
    }
    return addresses.find((item) => item.addressId === addressId) ?? null;
  }, [addresses, addressId, mode]);

  useEffect(() => {
    if (mode !== "edit") {
      setForm((prev) => ({ ...prev, setAsDefault: true }));
      return;
    }
    if (!editingAddress) {
      return;
    }

    const provinceCode = editingAddress.provinceCode || findProvinceCodeByName(editingAddress.provinceName);
    const cityCode = editingAddress.cityCode || findCityCodeByName(provinceCode, editingAddress.cityName);
    const districtCode = editingAddress.districtCode || findDistrictCodeByName(provinceCode, cityCode, editingAddress.districtName);

    setForm({
      label: editingAddress.label || "",
      receiverName: editingAddress.receiverName || "",
      receiverPhone: editingAddress.receiverPhone || "",
      provinceCode,
      provinceName: editingAddress.provinceName || "",
      cityCode,
      cityName: editingAddress.cityName || "",
      districtCode,
      districtName: editingAddress.districtName || "",
      street: editingAddress.street || "",
      detail: editingAddress.detail || "",
      postalCode: editingAddress.postalCode || "",
      setAsDefault: editingAddress.isDefault
    });
  }, [editingAddress, mode]);

  const provinceOptions = useMemo(() => getProvinceOptions(), []);
  const cityOptions = useMemo(
    () => ensureOption(getCityOptions(form.provinceCode), form.cityCode, form.cityName),
    [form.cityCode, form.cityName, form.provinceCode]
  );
  const districtOptions = useMemo(
    () => ensureOption(getDistrictOptions(form.provinceCode, form.cityCode), form.districtCode, form.districtName),
    [form.cityCode, form.districtCode, form.districtName, form.provinceCode]
  );

  const saveMutation = useMutation({
    mutationFn: async () => {
      if (!form.receiverName.trim()) {
        throw new Error(isZh ? "收件人不能为空" : "Receiver name is required");
      }
      if (!isValidPhone(form.receiverPhone)) {
        throw new Error(isZh ? "手机号格式不正确" : "Invalid phone number");
      }
      if (!form.provinceCode || !form.cityCode || !form.districtCode) {
        throw new Error(isZh ? "请选择省、市、区" : "Please select province, city and district");
      }
      if (!form.detail.trim()) {
        throw new Error(isZh ? "详细地址不能为空" : "Detail address is required");
      }

      if (mode === "edit") {
        if (!editingAddress) {
          throw new Error(isZh ? "地址不存在或已被删除" : "Address not found");
        }
        const result = await updateMyAddress(editingAddress.addressId, form, editingAddress.addressVersion);
        if (form.setAsDefault && !editingAddress.isDefault) {
          await setMyDefaultAddress(result.address.addressId, result.addressBookVersion || addressBookVersion);
        }
        return result.address.addressId;
      }

      const result = await createMyAddress(form, addressBookVersion);
      return result.address.addressId;
    },
    onSuccess: (savedAddressId) => {
      if (savedAddressId > 0) {
        window.localStorage.setItem("shopa_mall_default_address_id", String(savedAddressId));
      }
      messageApi.success(mode === "edit" ? (isZh ? "地址已更新" : "Address updated") : isZh ? "地址已创建" : "Address created");
      router.push("/me/address");
      router.refresh();
    },
    onError: (error) => {
      messageApi.error(toErrorText(error, "保存失败，请稍后重试", "Failed to save address", isZh));
    }
  });

  return (
    <section className="tb-address-editor">
      {contextHolder}
      <header className="tb-address-editor-header">
        <h1>{mode === "edit" ? (isZh ? "编辑收货地址" : "Edit Address") : isZh ? "新增收货地址" : "Create Address"}</h1>
        <Link href="/me/address">{isZh ? "返回地址列表" : "Back to addresses"}</Link>
      </header>

      <form
        className="tb-address-form"
        onSubmit={(event) => {
          event.preventDefault();
          saveMutation.mutate();
        }}
      >
        <label>
          <span>{isZh ? "收件人" : "Receiver"}</span>
          <input value={form.receiverName} onChange={(e) => setForm((prev) => ({ ...prev, receiverName: e.target.value }))} />
        </label>
        <label>
          <span>{isZh ? "手机号" : "Phone"}</span>
          <input value={form.receiverPhone} onChange={(e) => setForm((prev) => ({ ...prev, receiverPhone: e.target.value }))} />
        </label>

        <label>
          <span>{isZh ? "省份" : "Province"}</span>
          <select
            value={form.provinceCode}
            onChange={(e) => {
              const code = e.target.value;
              const province = findProvinceByCode(code);
              setForm((prev) => ({
                ...prev,
                provinceCode: code,
                provinceName: province?.name ?? "",
                cityCode: "",
                cityName: "",
                districtCode: "",
                districtName: ""
              }));
            }}
          >
            <option value="">{isZh ? "请选择省份" : "Select province"}</option>
            {provinceOptions.map((item) => (
              <option key={item.code} value={item.code}>
                {item.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          <span>{isZh ? "城市" : "City"}</span>
          <select
            value={form.cityCode}
            onChange={(e) => {
              const code = e.target.value;
              const city = findCityByCode(form.provinceCode, code);
              setForm((prev) => ({
                ...prev,
                cityCode: code,
                cityName: city?.name ?? "",
                districtCode: "",
                districtName: ""
              }));
            }}
            disabled={!form.provinceCode}
          >
            <option value="">{isZh ? "请选择城市" : "Select city"}</option>
            {cityOptions.map((item) => (
              <option key={item.code} value={item.code}>
                {item.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          <span>{isZh ? "区县" : "District"}</span>
          <select
            value={form.districtCode}
            onChange={(e) => {
              const code = e.target.value;
              const district = findDistrictByCode(form.provinceCode, form.cityCode, code);
              setForm((prev) => ({
                ...prev,
                districtCode: code,
                districtName: district?.name ?? ""
              }));
            }}
            disabled={!form.cityCode}
          >
            <option value="">{isZh ? "请选择区县" : "Select district"}</option>
            {districtOptions.map((item) => (
              <option key={item.code} value={item.code}>
                {item.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          <span>{isZh ? "街道" : "Street"}</span>
          <input value={form.street} onChange={(e) => setForm((prev) => ({ ...prev, street: e.target.value }))} />
        </label>

        <label className="span-2">
          <span>{isZh ? "详细地址" : "Detail Address"}</span>
          <input value={form.detail} onChange={(e) => setForm((prev) => ({ ...prev, detail: e.target.value }))} />
        </label>

        <label>
          <span>{isZh ? "邮编" : "Postal Code"}</span>
          <input value={form.postalCode} onChange={(e) => setForm((prev) => ({ ...prev, postalCode: e.target.value }))} />
        </label>

        <label>
          <span>{isZh ? "标签" : "Label"}</span>
          <input value={form.label} onChange={(e) => setForm((prev) => ({ ...prev, label: e.target.value }))} />
        </label>

        <label className="tb-address-checkbox span-2">
          <input
            type="checkbox"
            checked={Boolean(form.setAsDefault)}
            onChange={(e) => setForm((prev) => ({ ...prev, setAsDefault: e.target.checked }))}
          />
          <span>{isZh ? "保存后设为默认地址" : "Set as default after saving"}</span>
        </label>

        <div className="tb-address-form-actions span-2">
          <Link className="tb-address-cancel-link" href="/me/address">
            {isZh ? "取消" : "Cancel"}
          </Link>
          <button type="submit" disabled={saveMutation.isLoading || (mode === "edit" && !!addressId && !editingAddress && !addressQuery.isLoading)}>
            {saveMutation.isLoading ? (isZh ? "保存中..." : "Saving...") : isZh ? "保存地址" : "Save Address"}
          </button>
        </div>
      </form>

      {mode === "edit" && !editingAddress && !addressQuery.isLoading ? (
        <p className="tb-address-editor-tip">{isZh ? "未找到该地址，可能已被删除。" : "Address not found, it may have been deleted."}</p>
      ) : null}
    </section>
  );
}
