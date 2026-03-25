package media

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/media-svc/api/media/v1"
)

// IMediaV1 定义 media-svc 的 HTTP API。
type IMediaV1 interface {
	// 外部媒体接口
	InitUpload(ctx context.Context, req *v1.InitUploadReq) (res *v1.InitUploadRes, err error)
	CompleteUpload(ctx context.Context, req *v1.CompleteUploadReq) (res *v1.CompleteUploadRes, err error)
	IssueReadUrl(ctx context.Context, req *v1.IssueReadUrlReq) (res *v1.IssueReadUrlRes, err error)
	GetAssetProcessStatus(ctx context.Context, req *v1.GetAssetProcessStatusReq) (res *v1.GetAssetProcessStatusRes, err error)
	BatchGetAssetProcessStatus(ctx context.Context, req *v1.BatchGetAssetProcessStatusReq) (res *v1.BatchGetAssetProcessStatusRes, err error)
	BatchBindAssetsToBiz(ctx context.Context, req *v1.BatchBindAssetsToBizReq) (res *v1.BatchBindAssetsToBizRes, err error)
	ReplaceBizAssetBindings(ctx context.Context, req *v1.ReplaceBizAssetBindingsReq) (res *v1.ReplaceBizAssetBindingsRes, err error)
	BatchUnbindAssetsFromBiz(ctx context.Context, req *v1.BatchUnbindAssetsFromBizReq) (res *v1.BatchUnbindAssetsFromBizRes, err error)
	GetBizAssets(ctx context.Context, req *v1.GetBizAssetsReq) (res *v1.GetBizAssetsRes, err error)

	// 内部媒体接口
	CreateDerivedAsset(ctx context.Context, req *v1.CreateDerivedAssetReq) (res *v1.CreateDerivedAssetRes, err error)
	UpdateAssetProcessStatus(ctx context.Context, req *v1.UpdateAssetProcessStatusReq) (res *v1.UpdateAssetProcessStatusRes, err error)
}
