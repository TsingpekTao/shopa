"use client";

import { ReactNode, createContext, useContext, useMemo } from "react";

type AccessContextValue = {
  permissions: string[];
};

const AccessContext = createContext<AccessContextValue>({ permissions: [] });

export function PermissionProvider({ permissions, children }: { permissions: string[]; children: ReactNode }) {
  const value = useMemo<AccessContextValue>(() => ({ permissions: permissions ?? [] }), [permissions]);
  return <AccessContext.Provider value={value}>{children}</AccessContext.Provider>;
}

export function usePermissions() {
  return useContext(AccessContext).permissions;
}

export function hasPermission(permissions: string[], required: string): boolean {
  const normalizedRequired = (required || "").trim();
  if (!normalizedRequired) {
    return true;
  }
  for (const item of permissions || []) {
    const permission = (item || "").trim();
    if (!permission) {
      continue;
    }
    if (permission === normalizedRequired || permission === "*:*:*") {
      return true;
    }
    if (!permission.includes("*")) {
      continue;
    }
    const source = permission.split(":");
    const target = normalizedRequired.split(":");
    if (source.length !== target.length) {
      continue;
    }
    let matched = true;
    for (let i = 0; i < source.length; i += 1) {
      if (source[i] !== "*" && source[i] !== target[i]) {
        matched = false;
        break;
      }
    }
    if (matched) {
      return true;
    }
  }
  return false;
}

export function AccessControl({ require, fallback = null, children }: { require: string; fallback?: ReactNode; children: ReactNode }) {
  const permissions = usePermissions();
  if (!hasPermission(permissions, require)) {
    return <>{fallback}</>;
  }
  return <>{children}</>;
}
