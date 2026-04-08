import { getSellerSalesMetricValue, SellerSalesMetricType } from "./sales-helpers";
import { SellerSalesBucket } from "./types";

export type SellerSalesChartMode = "BAR" | "LINE" | "PIE";

export interface SellerSalesChartPoint {
  key: string;
  label: string;
  value: number;
  ratio: number;
}

export interface SellerSalesLinePoint extends SellerSalesChartPoint {
  x: number;
  y: number;
  showLabel: boolean;
}

export interface SellerSalesPieSlice {
  key: string;
  label: string;
  value: number;
  ratio: number;
  startAngle: number;
  endAngle: number;
}

const MAX_PIE_SLICES = 6;
const DEFAULT_PADDING_X = 22;
const DEFAULT_PADDING_Y = 18;

export function buildSellerSalesChartPoints(
  series: SellerSalesBucket[],
  metric: SellerSalesMetricType
): SellerSalesChartPoint[] {
  const values = series.map((bucket) => Math.max(0, getSellerSalesMetricValue(bucket, metric)));
  const maxValue = Math.max(...values, 0);
  return series.map((bucket, index) => ({
    key: bucket.bucketKey,
    label: bucket.bucketLabel,
    value: values[index],
    ratio: maxValue > 0 ? values[index] / maxValue : 0
  }));
}

export function shouldShowSellerSalesPointLabel(index: number, total: number): boolean {
  if (total <= 8) {
    return true;
  }
  if (index === 0 || index === total - 1) {
    return true;
  }
  const step = total <= 14 ? 2 : total <= 21 ? 3 : 4;
  return index % step === 0;
}

export function buildSellerSalesLinePoints(
  points: SellerSalesChartPoint[],
  width: number,
  height: number,
  paddingX = DEFAULT_PADDING_X,
  paddingY = DEFAULT_PADDING_Y
): SellerSalesLinePoint[] {
  if (!points.length || width <= 0 || height <= 0) {
    return [];
  }
  const usableWidth = Math.max(width - paddingX * 2, 1);
  const usableHeight = Math.max(height - paddingY * 2, 1);
  const denominator = Math.max(points.length - 1, 1);

  return points.map((point, index) => ({
    ...point,
    x: paddingX + (usableWidth / denominator) * index,
    y: height - paddingY - point.ratio * usableHeight,
    showLabel: shouldShowSellerSalesPointLabel(index, points.length)
  }));
}

export function buildSellerSalesLinePath(points: SellerSalesChartPoint[], width: number, height: number): string {
  const linePoints = buildSellerSalesLinePoints(points, width, height);
  if (!linePoints.length) {
    return "";
  }
  return linePoints
    .map((point, index) => `${index === 0 ? "M" : "L"} ${point.x.toFixed(2)} ${point.y.toFixed(2)}`)
    .join(" ");
}

export function buildSellerSalesLineAreaPath(points: SellerSalesChartPoint[], width: number, height: number): string {
  const linePoints = buildSellerSalesLinePoints(points, width, height);
  if (!linePoints.length) {
    return "";
  }
  const first = linePoints[0];
  const last = linePoints[linePoints.length - 1];
  const linePath = buildSellerSalesLinePath(points, width, height);
  return `${linePath} L ${last.x.toFixed(2)} ${(height - DEFAULT_PADDING_Y).toFixed(2)} L ${first.x.toFixed(2)} ${(height - DEFAULT_PADDING_Y).toFixed(2)} Z`;
}

export function buildSellerSalesPieSlices(
  series: SellerSalesBucket[],
  metric: SellerSalesMetricType
): SellerSalesPieSlice[] {
  const points = buildSellerSalesChartPoints(series, metric)
    .filter((point) => point.value > 0)
    .sort((left, right) => right.value - left.value);

  if (!points.length) {
    return [];
  }

  const sliceSource =
    points.length <= MAX_PIE_SLICES
      ? points
      : [
          ...points.slice(0, MAX_PIE_SLICES - 1),
          {
            key: "others",
            label: "其他",
            value: points.slice(MAX_PIE_SLICES - 1).reduce((sum, point) => sum + point.value, 0),
            ratio: 0
          }
        ];

  const total = sliceSource.reduce((sum, point) => sum + point.value, 0);
  let startAngle = -Math.PI / 2;

  return sliceSource.map((point) => {
    const ratio = total > 0 ? point.value / total : 0;
    const endAngle = startAngle + ratio * Math.PI * 2;
    const slice: SellerSalesPieSlice = {
      key: point.key,
      label: point.label,
      value: point.value,
      ratio,
      startAngle,
      endAngle
    };
    startAngle = endAngle;
    return slice;
  });
}
