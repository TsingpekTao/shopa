import { AddressPanel } from "@/features/user-profile/AddressPanel";
import { ProfilePanel } from "@/features/user-profile/ProfilePanel";

export default function ProfilePage() {
  return (
    <section style={{ display: "grid", gap: 24, gridTemplateColumns: "repeat(auto-fit, minmax(360px, 1fr))" }}>
      <ProfilePanel />
      <AddressPanel />
    </section>
  );
}
