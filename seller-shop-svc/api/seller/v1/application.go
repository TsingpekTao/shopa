package v1

import (
	pb "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateApplicationDraftReq struct {
	g.Meta `path:"/v1/seller/applications/draft" method:"post" tags:"SellerShop" summary:"创建入驻申请草稿"`
	Entity *pb.EntityProfile `json:"entity" v:"required#entity is required"`
	Shop   *pb.ShopProfile   `json:"shop" v:"required#shop is required"`
}

type CreateApplicationDraftRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type UpdateApplicationDraftReq struct {
	g.Meta          `path:"/v1/seller/applications/{applicationNo}/draft" method:"patch" tags:"SellerShop" summary:"更新入驻申请草稿"`
	ApplicationNo   string            `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
	ExpectedVersion int32             `json:"expectedVersion" v:"required#expectedVersion is required"`
	Entity          *pb.EntityProfile `json:"entity"`
	Shop            *pb.ShopProfile   `json:"shop"`
	UpdateMask      []string          `json:"updateMask"`
}

type UpdateApplicationDraftRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type SubmitApplicationReq struct {
	g.Meta          `path:"/v1/seller/applications/{applicationNo}/submit" method:"post" tags:"SellerShop" summary:"提交入驻申请"`
	ApplicationNo   string `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
	ExpectedVersion int32  `json:"expectedVersion" v:"required#expectedVersion is required"`
}

type SubmitApplicationRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type ResubmitApplicationReq struct {
	g.Meta                `path:"/v1/seller/applications/{rejectedApplicationNo}/resubmit" method:"post" tags:"SellerShop" summary:"驳回后重新提交申请"`
	RejectedApplicationNo string            `json:"rejectedApplicationNo" in:"path" v:"required#rejectedApplicationNo is required"`
	ExpectedVersion       int32             `json:"expectedVersion" v:"required#expectedVersion is required"`
	Entity                *pb.EntityProfile `json:"entity"`
	Shop                  *pb.ShopProfile   `json:"shop"`
	UpdateMask            []string          `json:"updateMask"`
	SubmitImmediately     bool              `json:"submitImmediately"`
}

type ResubmitApplicationRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type GetMyApplicationReq struct {
	g.Meta        `path:"/v1/seller/applications/{applicationNo}" method:"get" tags:"SellerShop" summary:"获取我的申请详情"`
	ApplicationNo string `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
}

type GetMyApplicationRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type ListMyApplicationsReq struct {
	g.Meta   `path:"/v1/seller/applications" method:"get" tags:"SellerShop" summary:"分页查询我的申请列表"`
	Page     int32   `json:"page" in:"query"`
	PageSize int32   `json:"pageSize" in:"query"`
	Statuses []int32 `json:"statuses" in:"query"`
}

type ListMyApplicationsRes struct {
	Applications []*pb.SellerApplication `json:"applications"`
	Page         int32                   `json:"page"`
	PageSize     int32                   `json:"pageSize"`
	Total        int64                   `json:"total"`
}

type ListApplicationsReq struct {
	g.Meta   `path:"/v1/admin/seller/applications" method:"get" tags:"SellerShopAdmin" summary:"分页查询申请列表（管理端）"`
	Page     int32   `json:"page" in:"query"`
	PageSize int32   `json:"pageSize" in:"query"`
	Statuses []int32 `json:"statuses" in:"query"`
	Keyword  string  `json:"keyword" in:"query"`
}

type ListApplicationsRes struct {
	Applications []*pb.SellerApplication `json:"applications"`
	Page         int32                   `json:"page"`
	PageSize     int32                   `json:"pageSize"`
	Total        int64                   `json:"total"`
}

type GetApplicationDetailReq struct {
	g.Meta        `path:"/v1/admin/seller/applications/{applicationNo}" method:"get" tags:"SellerShopAdmin" summary:"获取申请详情（管理端）"`
	ApplicationNo string `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
}

type GetApplicationDetailRes struct {
	Application *pb.SellerApplication `json:"application"`
}

type ApproveApplicationReq struct {
	g.Meta          `path:"/v1/admin/seller/applications/{applicationNo}/approve" method:"post" tags:"SellerShopAdmin" summary:"审核通过申请"`
	ApplicationNo   string `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
	ExpectedVersion int32  `json:"expectedVersion" v:"required#expectedVersion is required"`
	ReviewComment   string `json:"reviewComment"`
}

type ApproveApplicationRes struct {
	ApplicationNo string `json:"applicationNo"`
	ShopNo        string `json:"shopNo"`
	ShopStatus    int32  `json:"shopStatus"`
}

type RejectApplicationReq struct {
	g.Meta           `path:"/v1/admin/seller/applications/{applicationNo}/reject" method:"post" tags:"SellerShopAdmin" summary:"审核驳回申请"`
	ApplicationNo    string `json:"applicationNo" in:"path" v:"required#applicationNo is required"`
	ExpectedVersion  int32  `json:"expectedVersion" v:"required#expectedVersion is required"`
	RejectReasonCode string `json:"rejectReasonCode" v:"required#rejectReasonCode is required"`
	RejectComment    string `json:"rejectComment"`
}

type RejectApplicationRes struct {
	ApplicationNo string `json:"applicationNo"`
	Status        int32  `json:"status"`
}

type FreezeShopReq struct {
	g.Meta          `path:"/v1/admin/seller/shops/{shopNo}/freeze" method:"post" tags:"SellerShopAdmin" summary:"冻结店铺"`
	ShopNo          string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	ExpectedVersion int32  `json:"expectedVersion" v:"required#expectedVersion is required"`
	ReasonCode      string `json:"reasonCode" v:"required#reasonCode is required"`
	Reason          string `json:"reason"`
}

type FreezeShopRes struct {
	ShopNo     string `json:"shopNo"`
	ShopStatus int32  `json:"shopStatus"`
}

type CloseShopReq struct {
	g.Meta          `path:"/v1/admin/seller/shops/{shopNo}/close" method:"post" tags:"SellerShopAdmin" summary:"关闭店铺"`
	ShopNo          string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	ExpectedVersion int32  `json:"expectedVersion" v:"required#expectedVersion is required"`
	ReasonCode      string `json:"reasonCode" v:"required#reasonCode is required"`
	Reason          string `json:"reason"`
}

type CloseShopRes struct {
	ShopNo     string `json:"shopNo"`
	ShopStatus int32  `json:"shopStatus"`
}

type GetShopByNoReq struct {
	g.Meta `path:"/v1/internal/shops/{shopNo}" method:"get" tags:"SellerShopInternal" summary:"按 shopNo 查询店铺"`
	ShopNo string `json:"shopNo" in:"path" v:"required#shopNo is required"`
}

type GetShopByNoRes struct {
	Shop *pb.Shop `json:"shop"`
}

type BatchGetShopsByNoReq struct {
	g.Meta  `path:"/v1/internal/shops/batch-get" method:"post" tags:"SellerShopInternal" summary:"批量按 shopNo 查询店铺"`
	ShopNos []string `json:"shopNos" v:"required#shopNos is required"`
}

type BatchGetShopsByNoRes struct {
	Shops []*pb.Shop `json:"shops"`
}

type ListShopsByOwnerUserIdReq struct {
	g.Meta      `path:"/v1/internal/shops/by-owner/{ownerUserId}" method:"get" tags:"SellerShopInternal" summary:"按 ownerUserId 查询店铺列表"`
	OwnerUserId uint64 `json:"ownerUserId" in:"path" v:"required#ownerUserId is required"`
}

type ListShopsByOwnerUserIdRes struct {
	Shops []*pb.ShopSummary `json:"shops"`
}

type IsUserShopOwnerReq struct {
	g.Meta `path:"/v1/internal/shops/{shopNo}/owners/{userId}/check" method:"get" tags:"SellerShopInternal" summary:"校验用户是否为店铺所有者"`
	ShopNo string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	UserId uint64 `json:"userId" in:"path" v:"required#userId is required"`
}

type IsUserShopOwnerRes struct {
	IsOwner bool `json:"isOwner"`
}
