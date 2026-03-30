"use client";

import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { deleteMyAddress, listMyAddresses, setMyDefaultAddress } from "@/features/address/api";
import { UserAddress } from "@/features/address/types";

function buildAddressLine(address: Pick<UserAddress, "provinceName" | "cityName" | "districtName" | "street" | "detail">): string {
  return [address.provinceName, address.cityName, address.districtName, address.street, address.detail].filter((item) => item).join(" ");
}

function toErrorText(error: unknown, zhText: string, enText: string, isZh: boolean): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return isZh ? zhText : enText;
}

export function AddressBookManager({ embedded = false }: { embedded?: boolean }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const addressQuery = useQuery({
    queryKey: ["my-addresses"],
    queryFn: () => listMyAddresses(false),
    staleTime: 20_000,
    refetchOnWindowFocus: false
  });

  const addresses = addressQuery.data?.addresses ?? [];
  const addressBookVersion = addressQuery.data?.addressBookVersion ?? 0;

  const deleteMutation = useMutation({
    mutationFn: async (address: UserAddress) => deleteMyAddress(address.addressId, address.addressVersion),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["my-addresses"] });
      messageApi.success(isZh ? "地址已删除" : "Address deleted");
    },
    onError: (error) => {
      messageApi.error(toErrorText(error, "删除失败，请稍后重试", "Failed to delete address", isZh));
    }
  });

  const setDefaultMutation = useMutation({
    mutationFn: async (addressId: number) => setMyDefaultAddress(addressId, addressBookVersion),
    onSuccess: async (result) => {
      const defaultId = result.defaultAddressId;
      if (defaultId > 0) {
        window.localStorage.setItem("shopa_mall_default_address_id", String(defaultId));
      }
      await queryClient.invalidateQueries({ queryKey: ["my-addresses"] });
      messageApi.success(isZh ? "默认地址已更新" : "Default address updated");
    },
    onError: (error) => {
      messageApi.error(toErrorText(error, "设置默认地址失败", "Failed to set default address", isZh));
    }
  });

  return (
    <section className={embedded ? "tb-address-manager tb-address-manager-embedded" : "tb-address-manager"}>
      {contextHolder}
      <div className="tb-address-manager-header">
        <h2>{isZh ? "收货地址管理" : "Address Book"}</h2>
        <Link className="tb-address-new-link" href="/me/address/new">
          {isZh ? "新增地址" : "Add Address"}
        </Link>
      </div>

      {addresses.length === 0 && !addressQuery.isLoading ? (
        <div className="tb-address-empty">
          <p>{isZh ? "你还没有收货地址" : "No shipping address yet."}</p>
          <Link href="/me/address/new">{isZh ? "去新增地址" : "Create one now"}</Link>
        </div>
      ) : null}

      <div className="tb-address-grid">
        {addresses.map((address) => {
          const line = buildAddressLine(address);
          return (
            <article key={address.addressId} className={`tb-address-card ${address.isDefault ? "is-default" : ""}`}>
              <div className="tb-address-card-top">
                <strong>{address.receiverName || "--"}</strong>
                <span>{address.receiverPhone || "--"}</span>
                {address.isDefault ? <em>{isZh ? "默认" : "Default"}</em> : null}
              </div>
              <p className="tb-address-card-line">{line || "--"}</p>
              <p className="tb-address-card-label">
                {isZh ? "标签：" : "Label: "}
                {address.label || (isZh ? "未设置" : "N/A")}
              </p>
              <div className="tb-address-actions">
                <Link href={`/me/address/${address.addressId}/edit`}>{isZh ? "编辑" : "Edit"}</Link>
                <button
                  type="button"
                  onClick={() => setDefaultMutation.mutate(address.addressId)}
                  disabled={address.isDefault || setDefaultMutation.isLoading}
                >
                  {isZh ? "设为默认" : "Set Default"}
                </button>
                <button
                  type="button"
                  className="danger"
                  onClick={() => deleteMutation.mutate(address)}
                  disabled={deleteMutation.isLoading}
                >
                  {isZh ? "删除" : "Delete"}
                </button>
              </div>
            </article>
          );
        })}
      </div>

      {addressQuery.isLoading && <p>{isZh ? "正在加载地址..." : "Loading addresses..."}</p>}
      {addressQuery.isError && <p>{isZh ? "地址加载失败，请稍后重试。" : "Failed to load addresses."}</p>}
    </section>
  );
}
