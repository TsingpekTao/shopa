import assert from "node:assert/strict";
import {
  buildSellerProductEditPath,
  canResumeSellerProductEditing
} from "./draft-edit.js";

function runCanResumeSellerProductEditingTest() {
  assert.equal(canResumeSellerProductEditing("SPU_STATUS_DRAFT"), true);
  assert.equal(canResumeSellerProductEditing("SPU_STATUS_REJECTED"), true);
  assert.equal(canResumeSellerProductEditing("SPU_STATUS_REVIEWING"), false);
  assert.equal(canResumeSellerProductEditing("SPU_STATUS_ON_SHELF"), false);
}

function runBuildSellerProductEditPathTest() {
  assert.equal(
    buildSellerProductEditPath(" SPU/20260408  "),
    "/seller/publish?spuNo=SPU%2F20260408"
  );
  assert.equal(buildSellerProductEditPath("   "), "/seller/publish");
}

runCanResumeSellerProductEditingTest();
runBuildSellerProductEditPathTest();
console.log("seller product draft edit tests passed");
