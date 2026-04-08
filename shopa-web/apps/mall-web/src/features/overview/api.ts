import { apiClient } from "@/lib/api-client";
import { MallOverview } from "./types";

type RawOverviewRole = {
  roleCode?: string;
  role_code?: string;
  scopeTypeCode?: string;
  scope_type_code?: string;
  scopeNo?: string;
  scope_no?: string;
};

type RawOverview = {
  userId?: number | string;
  user_id?: number | string;
  accountStatusCode?: string;
  account_status_code?: string;
  roles?: RawOverviewRole[];
  displayName?: string;
  display_name?: string;
  avatarUrl?: string;
  avatar_url?: string;
  points?: number | string;
  partial?: boolean;
  degradedFields?: string[];
  degraded_fields?: string[];
};

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value);
  }
  return "";
}

export async function getMyOverview(): Promise<MallOverview> {
  const response = await apiClient.get<RawOverview>("/v1/me/overview");
  return {
    userId: toNumber(response?.userId ?? response?.user_id),
    accountStatusCode: toString(response?.accountStatusCode ?? response?.account_status_code),
    roles: Array.isArray(response?.roles)
      ? response.roles.map((role) => ({
          roleCode: toString(role.roleCode ?? role.role_code),
          scopeTypeCode: toString(role.scopeTypeCode ?? role.scope_type_code),
          scopeNo: toString(role.scopeNo ?? role.scope_no)
        }))
      : [],
    displayName: toString(response?.displayName ?? response?.display_name),
    avatarUrl: toString(response?.avatarUrl ?? response?.avatar_url),
    points: toNumber(response?.points),
    partial: Boolean(response?.partial),
    degradedFields: Array.isArray(response?.degradedFields ?? response?.degraded_fields)
      ? (response?.degradedFields ?? response?.degraded_fields ?? []).map((item) => toString(item)).filter(Boolean)
      : []
  };
}
