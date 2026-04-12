import { HelpContent, ServiceContent } from "../_info-content";
import { AssistantConsole } from "@/features/agent/assistant-console";

type HelpPageProps = {
  searchParams?: {
    tab?: string;
  };
};

export default function HelpPage({ searchParams }: HelpPageProps) {
  const tab = searchParams?.tab ?? "";

  if (tab === "official") {
    return <AssistantConsole variant="official" backHref="/help" />;
  }

  if (tab === "merchant") {
    return <ServiceContent />;
  }

  return <HelpContent />;
}
