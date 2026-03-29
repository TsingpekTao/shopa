package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/media-svc/api/v1"
	"github.com/TsingpekTao/shopa/media-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 鐎圭偟骞囨禍鍡楊嚠婢舵牔绗岄崘鍛村劥閻ㄥ嫬鐛熸担?RPC 閺堝秴濮熼妴?
type Controller struct {
	v1.UnimplementedMediaServiceServer
	v1.UnimplementedMediaInternalServiceServer
}

// Register 鐏忓棗鐛熸担?RPC 閺堝秴濮熷▔銊ュ斀閸?gRPC 閺堝秴濮熼崳銊ｂ偓?
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterMediaServiceServer(s.Server, ctrl)
	v1.RegisterMediaInternalServiceServer(s.Server, ctrl)
}

func (*Controller) InitUpload(ctx context.Context, req *v1.InitUploadReq) (res *v1.InitUploadRes, err error) {
	return service.Media().InitUpload(ctx, req)
}

func (*Controller) CompleteUpload(ctx context.Context, req *v1.CompleteUploadReq) (res *v1.CompleteUploadRes, err error) {
	return service.Media().CompleteUpload(ctx, req)
}

func (*Controller) IssueReadUrl(ctx context.Context, req *v1.IssueReadUrlReq) (res *v1.IssueReadUrlRes, err error) {
	return service.Media().IssueReadUrl(ctx, req)
}

func (*Controller) GetAssetProcessStatus(ctx context.Context, req *v1.GetAssetProcessStatusReq) (res *v1.GetAssetProcessStatusRes, err error) {
	return service.Media().GetAssetProcessStatus(ctx, req)
}

func (*Controller) BatchGetAssetProcessStatus(ctx context.Context, req *v1.BatchGetAssetProcessStatusReq) (res *v1.BatchGetAssetProcessStatusRes, err error) {
	return service.Media().BatchGetAssetProcessStatus(ctx, req)
}

func (*Controller) BatchBindAssetsToBiz(ctx context.Context, req *v1.BatchBindAssetsToBizReq) (res *v1.BatchBindAssetsToBizRes, err error) {
	return service.Media().BatchBindAssetsToBiz(ctx, req)
}

func (*Controller) ReplaceBizAssetBindings(ctx context.Context, req *v1.ReplaceBizAssetBindingsReq) (res *v1.ReplaceBizAssetBindingsRes, err error) {
	return service.Media().ReplaceBizAssetBindings(ctx, req)
}

func (*Controller) BatchUnbindAssetsFromBiz(ctx context.Context, req *v1.BatchUnbindAssetsFromBizReq) (res *v1.BatchUnbindAssetsFromBizRes, err error) {
	return service.Media().BatchUnbindAssetsFromBiz(ctx, req)
}

func (*Controller) GetBizAssets(ctx context.Context, req *v1.GetBizAssetsReq) (res *v1.GetBizAssetsRes, err error) {
	return service.Media().GetBizAssets(ctx, req)
}

func (*Controller) CreateDerivedAsset(ctx context.Context, req *v1.CreateDerivedAssetReq) (res *v1.CreateDerivedAssetRes, err error) {
	return service.Media().CreateDerivedAsset(ctx, req)
}

func (*Controller) UpdateAssetProcessStatus(ctx context.Context, req *v1.UpdateAssetProcessStatusReq) (res *v1.UpdateAssetProcessStatusRes, err error) {
	return service.Media().UpdateAssetProcessStatus(ctx, req)
}

