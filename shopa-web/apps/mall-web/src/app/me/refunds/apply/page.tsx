"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { formatCnyFromCents } from "@/lib/price";
import { applyRefundBatch, previewRefund } from "@/features/refund/api";
import { getRefundReasonLabel } from "@/features/refund/helpers";

const REASON_OPTIONS = [
  { value: "changed_mind", zh: "不想要了", en: "Changed my mind" },
  { value: "found_a_better_price", zh: "发现更优惠的价格", en: "Found a better price" },
  { value: "shipping_delay", zh: "发货时间太久", en: "Shipping delayed" }
] as const;

export default function RefundApplyPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const searchParams = useSearchParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();

  const orderNo = (searchParams.get("orderNo") ?? "").trim();
  const subOrderNo = (searchParams.get("subOrderNo") ?? "").trim();
  const itemNo = (searchParams.get("itemNo") ?? "").trim();

  const [reasonCode, setReasonCode] = useState<(typeof REASON_OPTIONS)[number]["value"]>("changed_mind");
  const [buyerRemark, setBuyerRemark] = useState("");

  const previewQuery = useQuery({
    queryKey: ["refund-preview", orderNo, subOrderNo, itemNo],
    queryFn: () => previewRefund({ orderNo, subOrderNo, itemNo }),
    enabled: Boolean(orderNo && subOrderNo),
    staleTime: 30_000,
    refetchOnWindowFocus: false
  });

  const snapshot = previewQuery.data?.snapshot;
  const targets = useMemo(() => {
    if (!orderNo || !subOrderNo) {
      return [];
    }
    const selectedItemNos = snapshot?.items.map((item) => item.itemNo).filter(Boolean) ?? [];
    return [
      {
        orderNo,
        subOrderNo,
        itemNo: itemNo || snapshot?.itemNo || "",
        selectedItemNos
      }
    ];
  }, [orderNo, subOrderNo, itemNo, snapshot]);

  const reasonLabelZh = useMemo(() => getRefundReasonLabel(reasonCode, true), [reasonCode]);

  const applyMutation = useMutation({
    mutationFn: () =>
      applyRefundBatch({
        targets,
        reasonCode,
        reasonDesc: reasonLabelZh,
        buyerRemark: buyerRemark.trim()
      }),
    onSuccess: async (detail) => {
      if (!detail?.batch.refundBatchNo) {
        messageApi.error(isZh ? "退款申请失败，请重试" : "Unable to create refund");
        return;
      }
      await queryClient.invalidateQueries({ queryKey: ["refund-batches"] });
      messageApi.success(isZh ? "退款申请已提交" : "Refund request submitted");
      router.push(`/me/refunds/${encodeURIComponent(detail.batch.refundBatchNo)}`);
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "提交失败" : "Submit failed";
      messageApi.error(errorMessage);
    }
  });

  const previewEmpty = !orderNo || !subOrderNo;

  return (
    <section className="tb-refund-page">
      {contextHolder}

      <header className="tb-refund-apply-header tb-refund-panel">
        <div>
          <p className="tb-refund-subhead">{isZh ? "申请退款" : "Request refund"}</p>
          <h1>{isZh ? "预发货退款" : "Pre-shipment refund"}</h1>
          <p>{isZh ? "先确认退款金额，再把申请提交给商家审核。" : "Confirm the refundable amount before sending the request to the seller."}</p>
        </div>
        <Link href="/me/orders" className="tb-order-btn tb-order-btn-secondary">
          {isZh ? "查看订单" : "View orders"}
        </Link>
      </header>

      {previewEmpty ? (
        <section className="tb-order-empty-card">
          <strong>{isZh ? "请从订单页发起退款" : "Start from the order page"}</strong>
          <p>{isZh ? "在待发货订单里点击“申请退款”即可进入这里。" : "Open a pending shipment order and click 'Request Refund'."}</p>
        </section>
      ) : previewQuery.isLoading ? (
        <p className="tb-chat-empty">{isZh ? "正在获取退款预览..." : "Loading preview..."}</p>
      ) : !snapshot ? (
        <p className="tb-chat-empty">{isZh ? "暂时无法获取退款预览，请稍后重试。" : "Unable to load the refund preview right now."}</p>
      ) : (
        <>
          <section className="tb-refund-panel tb-refund-apply-preview">
            <div className="tb-refund-apply-summary">
              <div>
                <p>{isZh ? "可退金额" : "Refundable"}</p>
                <strong className="is-price">{`CNY ${formatCnyFromCents(snapshot.refundableAmount)}`}</strong>
              </div>
              <div>
                <p>{isZh ? "原始订单" : "Order"}</p>
                <strong>{snapshot.orderNo}</strong>
                <small>{snapshot.subOrderNo}</small>
              </div>
              <div>
                <p>{isZh ? "支付状态" : "Payment"}</p>
                <strong>{snapshot.paymentStatus || (isZh ? "未知" : "Unknown")}</strong>
              </div>
            </div>

            <div className="tb-refund-apply-items">
              {(snapshot.items ?? []).map((item) => (
                <article key={item.itemNo} className="tb-refund-apply-item">
                  <div>
                    <strong>{item.spuTitle || item.skuName || item.skuNo}</strong>
                    <p>{`${item.skuName || item.skuNo} x${item.qty}`}</p>
                  </div>
                  <span className="is-price">{`CNY ${formatCnyFromCents(item.salePrice)}`}</span>
                </article>
              ))}
            </div>
          </section>

          <section className="tb-refund-panel tb-refund-apply-form">
            <label className="tb-refund-form-label">
              {isZh ? "退款原因" : "Refund reason"}
              <select value={reasonCode} onChange={(event) => setReasonCode(event.target.value as (typeof REASON_OPTIONS)[number]["value"])}>
                {REASON_OPTIONS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {isZh ? option.zh : option.en}
                  </option>
                ))}
              </select>
            </label>

            <label className="tb-refund-form-label">
              {isZh ? "补充说明" : "Additional note"}
              <textarea
                value={buyerRemark}
                onChange={(event) => setBuyerRemark(event.target.value)}
                placeholder={isZh ? "可以补充你的退款说明，方便商家处理" : "Add extra details to help the seller review the refund"}
                rows={4}
              />
            </label>

            <div className="tb-refund-apply-actions">
              <button
                type="button"
                className="tb-order-btn tb-order-btn-primary"
                disabled={applyMutation.isPending || !snapshot}
                onClick={() => applyMutation.mutate()}
              >
                {applyMutation.isPending ? (isZh ? "提交中..." : "Submitting...") : isZh ? "提交退款申请" : "Submit refund"}
              </button>
              <Link href="/me/refunds" className="tb-order-btn tb-order-btn-secondary">
                {isZh ? "返回退款列表" : "Back to refunds"}
              </Link>
            </div>
          </section>
        </>
      )}
    </section>
  );
}
