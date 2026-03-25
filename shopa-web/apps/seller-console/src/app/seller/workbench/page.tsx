import { WorkbenchPanel } from "@/features/seller-shop/WorkbenchPanel";

export default function WorkbenchPage() {
  return (
    <section style={{ display: "grid", gap: 16 }}>
      <header
        style={{
          padding: 20,
          borderRadius: 12,
          background: "linear-gradient(135deg,#0f2a44,#162a63)",
          color: "#fff"
        }}
      >
        <h2 style={{ margin: "0 0 8px" }}>Seller Workbench</h2>
        <p style={{ margin: 0, color: "#c7d8ff" }}>
          Everything you need to keep your shops, applications, products and inventory under control.
        </p>
      </header>
      <WorkbenchPanel />
    </section>
  );
}
