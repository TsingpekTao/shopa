import assert from "node:assert/strict";
import { buildSellerSalesChartPoints, buildSellerSalesLinePath, buildSellerSalesPieSlices } from "./sales-chart.js";
import { SellerSalesBucket } from "./types.js";

const series: SellerSalesBucket[] = [
  {
    bucketKey: "2026-04-01",
    bucketLabel: "04-01",
    gmv: 28000,
    paidOrderCount: 12,
    paidBuyerCount: 9,
    refundAmount: 1200,
    refundRate: 0.0428
  },
  {
    bucketKey: "2026-04-02",
    bucketLabel: "04-02",
    gmv: 46000,
    paidOrderCount: 18,
    paidBuyerCount: 15,
    refundAmount: 800,
    refundRate: 0.0173
  },
  {
    bucketKey: "2026-04-03",
    bucketLabel: "04-03",
    gmv: 12000,
    paidOrderCount: 5,
    paidBuyerCount: 5,
    refundAmount: 0,
    refundRate: 0
  },
  {
    bucketKey: "2026-04-04",
    bucketLabel: "04-04",
    gmv: 6000,
    paidOrderCount: 3,
    paidBuyerCount: 3,
    refundAmount: 200,
    refundRate: 0.0333
  },
  {
    bucketKey: "2026-04-05",
    bucketLabel: "04-05",
    gmv: 9000,
    paidOrderCount: 4,
    paidBuyerCount: 4,
    refundAmount: 0,
    refundRate: 0
  },
  {
    bucketKey: "2026-04-06",
    bucketLabel: "04-06",
    gmv: 3000,
    paidOrderCount: 2,
    paidBuyerCount: 2,
    refundAmount: 0,
    refundRate: 0
  },
  {
    bucketKey: "2026-04-07",
    bucketLabel: "04-07",
    gmv: 15000,
    paidOrderCount: 6,
    paidBuyerCount: 6,
    refundAmount: 500,
    refundRate: 0.0333
  }
];

function runBuildSellerSalesChartPointsTest() {
  const points = buildSellerSalesChartPoints(series, "GMV");
  assert.equal(points.length, series.length);
  assert.equal(points[1].value, 46000);
  assert.equal(points[1].ratio, 1);
  assert.equal(points[3].ratio > 0, true);
}

function runBuildSellerSalesLinePathTest() {
  const points = buildSellerSalesChartPoints(series, "ORDERS");
  const path = buildSellerSalesLinePath(points, 520, 220);
  assert.match(path, /^M\s/);
  assert.match(path, /L/);
}

function runBuildSellerSalesPieSlicesTest() {
  const slices = buildSellerSalesPieSlices(series, "GMV");
  assert.equal(slices.length, 6);
  assert.equal(slices[0].value >= slices[1].value, true);
  assert.equal(slices.some((slice) => slice.label === "其他"), true);
  assert.equal(
    Math.round(slices.reduce((sum, slice) => sum + slice.ratio, 0) * 1000) / 1000,
    1
  );
}

runBuildSellerSalesChartPointsTest();
runBuildSellerSalesLinePathTest();
runBuildSellerSalesPieSlicesTest();
console.log("seller sales chart tests passed");
