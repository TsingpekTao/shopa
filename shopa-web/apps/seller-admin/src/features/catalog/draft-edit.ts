import { SpuStatusCode } from "./types";

export function canResumeSellerProductEditing(status: SpuStatusCode): boolean {
  return status === "SPU_STATUS_DRAFT" || status === "SPU_STATUS_REJECTED";
}

export function buildSellerProductEditPath(spuNo: string): string {
  const normalizedSpuNo = spuNo.trim();
  if (!normalizedSpuNo) {
    return "/seller/publish";
  }
  return `/seller/publish?spuNo=${encodeURIComponent(normalizedSpuNo)}`;
}
