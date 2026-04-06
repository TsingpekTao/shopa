import Link from "next/link";
import { ServiceContent } from "../_info-content";

export default function ServicePage() {
  return (
    <>
      <section className="tb-assistant-entry">
        <div>
          <span className="tb-chip">Shopa AI</span>
          <h1>智能客服与导购助手</h1>
          <p>从商品推荐、订单进度到售后规则，统一在一个入口里问清楚。</p>
        </div>
        <Link href="/service/assistant" className="tb-btn-main">
          立即进入
        </Link>
      </section>
      <ServiceContent />
    </>
  );
}
