import { Card } from "antd";

type Props = {
  title: string;
  description: string;
};

export function FeatureCard({ title, description }: Props) {
  return (
    <Card hoverable size="small" className="feature-card" role="article" aria-label={title}>
      <Card.Meta title={title} description={description} />
    </Card>
  );
}
