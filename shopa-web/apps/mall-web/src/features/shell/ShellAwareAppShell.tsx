"use client";

import { ReactNode } from "react";
import { AppShell } from "@shopa/ui";
import { useShellStatus } from "./status";

export function ShellAwareAppShell({ children }: { children: ReactNode }) {
  const { cartTypeCount, unreadMessages } = useShellStatus();
  return (
    <AppShell
      mode="mall"
      shellState={{
        cartTypeCount,
        unreadMessages
      }}
    >
      {children}
    </AppShell>
  );
}
