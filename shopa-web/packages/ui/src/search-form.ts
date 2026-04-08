export type GlobalSearchTab = "item" | "shop";

export function normalizeGlobalSearchTab(value: string | null | undefined): GlobalSearchTab {
  return value === "shop" ? "shop" : "item";
}

export function buildGlobalSearchHref(query: string, tab: GlobalSearchTab): string {
  const params = new URLSearchParams();
  const normalizedQuery = query.trim();
  if (normalizedQuery) {
    params.set("q", normalizedQuery);
  }
  params.set("tab", normalizeGlobalSearchTab(tab));
  return `/search?${params.toString()}`;
}
