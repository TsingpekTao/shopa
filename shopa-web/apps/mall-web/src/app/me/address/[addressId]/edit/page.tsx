"use client";

import { useMemo } from "react";
import { useParams } from "next/navigation";
import { AddressEditor } from "@/components/address-editor";

function parseAddressId(raw: string | string[] | undefined): number {
  const value = Array.isArray(raw) ? raw[0] : raw;
  const id = Number(value ?? "");
  if (!Number.isFinite(id) || id <= 0) {
    return 0;
  }
  return Math.trunc(id);
}

export default function EditAddressPage() {
  const params = useParams<{ addressId: string }>();
  const addressId = useMemo(() => parseAddressId(params?.addressId), [params?.addressId]);

  return <AddressEditor mode="edit" addressId={addressId} />;
}

