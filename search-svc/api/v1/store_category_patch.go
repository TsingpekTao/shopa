package searchv1

import timestamppb "google.golang.org/protobuf/types/known/timestamppb"

type PatchSpuStoreCategoryReq struct {
	SpuNo             string                 `json:"spu_no,omitempty"`
	StoreCategoryId   uint64                 `json:"store_category_id,omitempty"`
	StoreCategoryL1   uint64                 `json:"store_category_l1,omitempty"`
	StoreCategoryL2   uint64                 `json:"store_category_l2,omitempty"`
	StoreCategoryPath []uint64               `json:"store_category_path,omitempty"`
	UpdatedAt         *timestamppb.Timestamp `json:"updated_at,omitempty"`
}

type PatchSpuStoreCategoryRes struct {
	Ok bool `json:"ok,omitempty"`
}
