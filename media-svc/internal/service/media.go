// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/media-svc/api/v1"
)

type (
	IMedia interface {
		InitUpload(ctx context.Context, req *v1.InitUploadReq) (*v1.InitUploadRes, error)
		CompleteUpload(ctx context.Context, req *v1.CompleteUploadReq) (*v1.CompleteUploadRes, error)
		IssueReadUrl(ctx context.Context, req *v1.IssueReadUrlReq) (*v1.IssueReadUrlRes, error)
		GetAssetProcessStatus(ctx context.Context, req *v1.GetAssetProcessStatusReq) (*v1.GetAssetProcessStatusRes, error)
		BatchGetAssetProcessStatus(ctx context.Context, req *v1.BatchGetAssetProcessStatusReq) (*v1.BatchGetAssetProcessStatusRes, error)
		BatchBindAssetsToBiz(ctx context.Context, req *v1.BatchBindAssetsToBizReq) (*v1.BatchBindAssetsToBizRes, error)
		ReplaceBizAssetBindings(ctx context.Context, req *v1.ReplaceBizAssetBindingsReq) (*v1.ReplaceBizAssetBindingsRes, error)
		BatchUnbindAssetsFromBiz(ctx context.Context, req *v1.BatchUnbindAssetsFromBizReq) (*v1.BatchUnbindAssetsFromBizRes, error)
		GetBizAssets(ctx context.Context, req *v1.GetBizAssetsReq) (*v1.GetBizAssetsRes, error)
		CreateDerivedAsset(ctx context.Context, req *v1.CreateDerivedAssetReq) (*v1.CreateDerivedAssetRes, error)
		UpdateAssetProcessStatus(ctx context.Context, req *v1.UpdateAssetProcessStatusReq) (*v1.UpdateAssetProcessStatusRes, error)
	}
)

var (
	localMedia IMedia
)

func Media() IMedia {
	if localMedia == nil {
		panic("implement not found for interface IMedia, forgot register?")
	}
	return localMedia
}

func RegisterMedia(i IMedia) {
	localMedia = i
}
