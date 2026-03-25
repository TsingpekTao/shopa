import { Alert } from "antd";

type Props = {
  fields: string[];
};

export function DegradedBanner({ fields }: Props) {
  if (!fields.length) {
    return null;
  }

  return (
    <Alert
      type="warning"
      showIcon
      role="status"
      message="Some side services are temporarily degraded"
      description={`Impacted fields: ${fields.join(", ")}. Core actions remain available.`}
      style={{ marginBottom: 16 }}
    />
  );
}
