package v1

import (
	pb "github.com/TsingpekTao/shopa/media-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type InitUploadReq struct {
	g.Meta         `path:"/v1/media/upload/init" method:"post" tags:"Media" summary:"初始化上传，返回上传凭证"`
	SceneCode      string  `json:"sceneCode" v:"required#sceneCode is required"`
	BizType        *string `json:"bizType"`
	BizNo          *string `json:"bizNo"`
	FileName       string  `json:"fileName" v:"required#fileName is required"`
	MimeType       string  `json:"mimeType" v:"required#mimeType is required"`
	SizeBytes      uint64  `json:"sizeBytes" v:"required#sizeBytes is required"`
	ChecksumSha256 string  `json:"checksumSha256" v:"required#checksumSha256 is required"`
}

type InitUploadRes struct {
	Asset        *pb.Asset        `json:"asset"`
	UploadTicket *pb.UploadTicket `json:"uploadTicket"`
}

type CompleteUploadReq struct {
	g.Meta         `path:"/v1/media/upload/complete" method:"post" tags:"Media" summary:"确认上传完成"`
	AssetId        uint64 `json:"assetId" v:"required#assetId is required"`
	Etag           string `json:"etag"`
	SizeBytes      uint64 `json:"sizeBytes" v:"required#sizeBytes is required"`
	MimeType       string `json:"mimeType" v:"required#mimeType is required"`
	ChecksumSha256 string `json:"checksumSha256" v:"required#checksumSha256 is required"`
}

type CompleteUploadRes struct {
	Asset *pb.Asset `json:"asset"`
}

type IssueReadUrlReq struct {
	g.Meta     `path:"/v1/media/assets/{assetId}/read-url" method:"get" tags:"Media" summary:"签发读取 URL"`
	AssetId    uint64 `json:"assetId" in:"path" v:"required#assetId is required"`
	TtlSeconds uint32 `json:"ttlSeconds" in:"query"`
}

type IssueReadUrlRes struct {
	AssetId   uint64 `json:"assetId"`
	Url       string `json:"url"`
	IsPublic  bool   `json:"isPublic"`
	ExpiredAt string `json:"expiredAt"`
}

type GetAssetProcessStatusReq struct {
	g.Meta  `path:"/v1/media/assets/{assetId}/process-status" method:"get" tags:"Media" summary:"查询单个资产处理状态"`
	AssetId uint64 `json:"assetId" in:"path" v:"required#assetId is required"`
}

type GetAssetProcessStatusRes struct {
	AssetId             uint64                    `json:"assetId"`
	ProcessStatus       pb.ProcessStatus          `json:"processStatus"`
	ProcessProgress     uint32                    `json:"processProgress"`
	ProcessErrorCode    string                    `json:"processErrorCode"`
	ProcessErrorMessage string                    `json:"processErrorMessage"`
	DerivedAssets       []*pb.DerivedAssetSummary `json:"derivedAssets"`
}

type BatchGetAssetProcessStatusReq struct {
	g.Meta   `path:"/v1/media/assets/process-status/batch-get" method:"post" tags:"Media" summary:"批量查询资产处理状态"`
	AssetIds []uint64 `json:"assetIds" v:"required#assetIds is required"`
}

type BatchGetAssetProcessStatusRes struct {
	Statuses []*pb.GetAssetProcessStatusRes `json:"statuses"`
}

type BatchBindAssetsToBizReq struct {
	g.Meta    `path:"/v1/media/bindings/batch-bind" method:"post" tags:"MediaBinding" summary:"批量绑定资产到业务对象"`
	SceneCode string            `json:"sceneCode" v:"required#sceneCode is required"`
	BizType   string            `json:"bizType" v:"required#bizType is required"`
	BizNo     string            `json:"bizNo" v:"required#bizNo is required"`
	Items     []*pb.BindingItem `json:"items" v:"required#items is required"`
}

type BatchBindAssetsToBizRes struct {
	Bindings []*pb.BizAssetBinding `json:"bindings"`
}

type ReplaceBizAssetBindingsReq struct {
	g.Meta       `path:"/v1/media/bindings/replace" method:"post" tags:"MediaBinding" summary:"按业务坑位整体替换绑定"`
	SceneCode    string                  `json:"sceneCode" v:"required#sceneCode is required"`
	BizType      string                  `json:"bizType" v:"required#bizType is required"`
	BizNo        string                  `json:"bizNo" v:"required#bizNo is required"`
	BindingField string                  `json:"bindingField" v:"required#bindingField is required"`
	Items        []*pb.SimpleBindingItem `json:"items"`
}

type ReplaceBizAssetBindingsRes struct {
	Bindings []*pb.BizAssetBinding `json:"bindings"`
}

type BatchUnbindAssetsFromBizReq struct {
	g.Meta       `path:"/v1/media/bindings/batch-unbind" method:"post" tags:"MediaBinding" summary:"批量解绑业务资产"`
	SceneCode    string   `json:"sceneCode" v:"required#sceneCode is required"`
	BizType      string   `json:"bizType" v:"required#bizType is required"`
	BizNo        string   `json:"bizNo" v:"required#bizNo is required"`
	BindingField string   `json:"bindingField" v:"required#bindingField is required"`
	AssetIds     []uint64 `json:"assetIds" v:"required#assetIds is required"`
}

type BatchUnbindAssetsFromBizRes struct {
	AffectedRows uint64 `json:"affectedRows"`
}

type GetBizAssetsReq struct {
	g.Meta          `path:"/v1/media/biz-assets" method:"get" tags:"MediaBinding" summary:"查询业务对象绑定资产"`
	SceneCode       string `json:"sceneCode" in:"query" v:"required#sceneCode is required"`
	BizType         string `json:"bizType" in:"query" v:"required#bizType is required"`
	BizNo           string `json:"bizNo" in:"query" v:"required#bizNo is required"`
	BindingField    string `json:"bindingField" in:"query"`
	IncludeInactive bool   `json:"includeInactive" in:"query"`
}

type GetBizAssetsRes struct {
	Assets []*pb.BizAsset `json:"assets"`
}

type CreateDerivedAssetReq struct {
	g.Meta          `path:"/v1/internal/media/derived-assets" method:"post" tags:"MediaInternal" summary:"创建衍生资产（内部）"`
	ParentAssetId   uint64         `json:"parentAssetId" v:"required#parentAssetId is required"`
	SceneCode       string         `json:"sceneCode" v:"required#sceneCode is required"`
	DerivedKind     pb.DerivedKind `json:"derivedKind" v:"required#derivedKind is required"`
	FileName        string         `json:"fileName" v:"required#fileName is required"`
	MimeType        string         `json:"mimeType" v:"required#mimeType is required"`
	SizeBytes       uint64         `json:"sizeBytes" v:"required#sizeBytes is required"`
	ChecksumSha256  string         `json:"checksumSha256" v:"required#checksumSha256 is required"`
	Etag            string         `json:"etag"`
	StorageProvider string         `json:"storageProvider" v:"required#storageProvider is required"`
	Bucket          string         `json:"bucket" v:"required#bucket is required"`
	ObjectKey       string         `json:"objectKey" v:"required#objectKey is required"`
	PublicUrl       string         `json:"publicUrl"`
}

type CreateDerivedAssetRes struct {
	Asset *pb.Asset `json:"asset"`
}

type UpdateAssetProcessStatusReq struct {
	g.Meta              `path:"/v1/internal/media/assets/process-status" method:"post" tags:"MediaInternal" summary:"更新资产处理状态（内部）"`
	AssetId             uint64           `json:"assetId" v:"required#assetId is required"`
	ProcessStatus       pb.ProcessStatus `json:"processStatus" v:"required#processStatus is required"`
	ProcessProgress     uint32           `json:"processProgress"`
	ProcessErrorCode    string           `json:"processErrorCode"`
	ProcessErrorMessage string           `json:"processErrorMessage"`
}

type UpdateAssetProcessStatusRes struct {
	Updated bool `json:"updated"`
}
