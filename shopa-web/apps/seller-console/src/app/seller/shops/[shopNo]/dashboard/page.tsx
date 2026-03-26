import {ShopDashboardPanel} from "@/features/seller-shop/ShopDashboardPanel";

export function generateStaticParams() {
  return [{shopNo: "demo-shop"}];
}

export default function ShopDashboardPage({params}: {params: {shopNo: string}}) {
  return (
    <section style={{ display: "grid", gap: 16 }}>
      <header
        style={{
          padding: 20,
          borderRadius: 12,
          background: "linear-gradient(135deg,#0b1f33,#112849)",
          color: "#fff"
        }}
      >
        <h2 style={{ margin: "0 0 8px" }}>Shop {params.shopNo} Dashboard</h2>
        <p style={{ margin: 0, color: "#d0daff" }}>
          Live view sourced from seller-shop and inventory aggregates. Partial data is shown when any
          side service degrades.
        </p>
      </header>
      <ShopDashboardPanel shopNo={params.shopNo} />
    </section>
  );
}
