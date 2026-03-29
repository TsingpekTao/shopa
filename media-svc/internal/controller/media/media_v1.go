package media

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/media-svc/api/media/v1"
	pb "github.com/TsingpekTao/shopa/media-svc/api/v1"
)

func (c *ControllerV1) InitUpload(ctx context.Context, req *httpv1.InitUploadReq) (res *httpv1.InitUploadRes, err error) {
	out, err := c.svc.InitUpload(withRequestMetadata(ctx), &pb.InitUploadReq{
		SceneCode:      req.SceneCode,
		BizType:        req.BizType,
		BizNo:          req.BizNo,
		FileName:       req.FileName,
		MimeType:       req.MimeType,
		SizeBytes:      req.SizeBytes,
		ChecksumSha256: req.ChecksumSha256,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.InitUploadRes{
		Asset:        out.GetAsset(),
		UploadTicket: out.GetUploadTicket(),
	}, nil
}

func (c *ControllerV1) CompleteUpload(ctx context.Context, req *httpv1.CompleteUploadReq) (res *httpv1.CompleteUploadRes, err error) {
	out, err := c.svc.CompleteUpload(withRequestMetadata(ctx), &pb.CompleteUploadReq{
		AssetId:        req.AssetId,
		Etag:           req.Etag,
		SizeBytes:      req.SizeBytes,
		MimeType:       req.MimeType,
		ChecksumSha256: req.ChecksumSha256,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.CompleteUploadRes{Asset: out.GetAsset()}, nil
}

func (c *ControllerV1) IssueReadUrl(ctx context.Context, req *httpv1.IssueReadUrlReq) (res *httpv1.IssueReadUrlRes, err error) {
	out, err := c.svc.IssueReadUrl(withRequestMetadata(ctx), &pb.IssueReadUrlReq{
		AssetId:    req.AssetId,
		TtlSeconds: req.TtlSeconds,
	})
	if err != nil {
		return nil, err
	}
	expiredAt := ""
	if out.GetExpiredAt() != nil {
		expiredAt = out.GetExpiredAt().AsTime().Format("2006-01-02T15:04:05Z07:00")
	}

	return &httpv1.IssueReadUrlRes{
		AssetId:   out.GetAssetId(),
		Url:       out.GetUrl(),
		IsPublic:  out.GetIsPublic(),
		ExpiredAt: expiredAt,
	}, nil
}

func (c *ControllerV1) GetAssetProcessStatus(ctx context.Context, req *httpv1.GetAssetProcessStatusReq) (res *httpv1.GetAssetProcessStatusRes, err error) {
	out, err := c.svc.GetAssetProcessStatus(withRequestMetadata(ctx), &pb.GetAssetProcessStatusReq{
		AssetId: req.AssetId,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.GetAssetProcessStatusRes{
		AssetId:             out.GetAssetId(),
		ProcessStatus:       out.GetProcessStatus(),
		ProcessProgress:     out.GetProcessProgress(),
		ProcessErrorCode:    out.GetProcessErrorCode(),
		ProcessErrorMessage: out.GetProcessErrorMessage(),
		DerivedAssets:       out.GetDerivedAssets(),
	}, nil
}

func (c *ControllerV1) BatchGetAssetProcessStatus(ctx context.Context, req *httpv1.BatchGetAssetProcessStatusReq) (res *httpv1.BatchGetAssetProcessStatusRes, err error) {
	out, err := c.svc.BatchGetAssetProcessStatus(withRequestMetadata(ctx), &pb.BatchGetAssetProcessStatusReq{
		AssetIds: req.AssetIds,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.BatchGetAssetProcessStatusRes{Statuses: out.GetStatuses()}, nil
}

func (c *ControllerV1) BatchBindAssetsToBiz(ctx context.Context, req *httpv1.BatchBindAssetsToBizReq) (res *httpv1.BatchBindAssetsToBizRes, err error) {
	out, err := c.svc.BatchBindAssetsToBiz(withRequestMetadata(ctx), &pb.BatchBindAssetsToBizReq{
		SceneCode: req.SceneCode,
		BizType:   req.BizType,
		BizNo:     req.BizNo,
		Items:     req.Items,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.BatchBindAssetsToBizRes{Bindings: out.GetBindings()}, nil
}

func (c *ControllerV1) ReplaceBizAssetBindings(ctx context.Context, req *httpv1.ReplaceBizAssetBindingsReq) (res *httpv1.ReplaceBizAssetBindingsRes, err error) {
	out, err := c.svc.ReplaceBizAssetBindings(withRequestMetadata(ctx), &pb.ReplaceBizAssetBindingsReq{
		SceneCode:    req.SceneCode,
		BizType:      req.BizType,
		BizNo:        req.BizNo,
		BindingField: req.BindingField,
		Items:        req.Items,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.ReplaceBizAssetBindingsRes{Bindings: out.GetBindings()}, nil
}

func (c *ControllerV1) BatchUnbindAssetsFromBiz(ctx context.Context, req *httpv1.BatchUnbindAssetsFromBizReq) (res *httpv1.BatchUnbindAssetsFromBizRes, err error) {
	out, err := c.svc.BatchUnbindAssetsFromBiz(withRequestMetadata(ctx), &pb.BatchUnbindAssetsFromBizReq{
		SceneCode:    req.SceneCode,
		BizType:      req.BizType,
		BizNo:        req.BizNo,
		BindingField: req.BindingField,
		AssetIds:     req.AssetIds,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.BatchUnbindAssetsFromBizRes{AffectedRows: out.GetAffectedRows()}, nil
}

func (c *ControllerV1) GetBizAssets(ctx context.Context, req *httpv1.GetBizAssetsReq) (res *httpv1.GetBizAssetsRes, err error) {
	out, err := c.svc.GetBizAssets(withRequestMetadata(ctx), &pb.GetBizAssetsReq{
		SceneCode:       req.SceneCode,
		BizType:         req.BizType,
		BizNo:           req.BizNo,
		BindingField:    req.BindingField,
		IncludeInactive: req.IncludeInactive,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.GetBizAssetsRes{Assets: out.GetAssets()}, nil
}

func (c *ControllerV1) CreateDerivedAsset(ctx context.Context, req *httpv1.CreateDerivedAssetReq) (res *httpv1.CreateDerivedAssetRes, err error) {
	out, err := c.svc.CreateDerivedAsset(withRequestMetadata(ctx), &pb.CreateDerivedAssetReq{
		ParentAssetId:   req.ParentAssetId,
		SceneCode:       req.SceneCode,
		DerivedKind:     req.DerivedKind,
		FileName:        req.FileName,
		MimeType:        req.MimeType,
		SizeBytes:       req.SizeBytes,
		ChecksumSha256:  req.ChecksumSha256,
		Etag:            req.Etag,
		StorageProvider: req.StorageProvider,
		Bucket:          req.Bucket,
		ObjectKey:       req.ObjectKey,
		PublicUrl:       req.PublicUrl,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.CreateDerivedAssetRes{Asset: out.GetAsset()}, nil
}

func (c *ControllerV1) UpdateAssetProcessStatus(ctx context.Context, req *httpv1.UpdateAssetProcessStatusReq) (res *httpv1.UpdateAssetProcessStatusRes, err error) {
	out, err := c.svc.UpdateAssetProcessStatus(withRequestMetadata(ctx), &pb.UpdateAssetProcessStatusReq{
		AssetId:             req.AssetId,
		ProcessStatus:       req.ProcessStatus,
		ProcessProgress:     req.ProcessProgress,
		ProcessErrorCode:    req.ProcessErrorCode,
		ProcessErrorMessage: req.ProcessErrorMessage,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.UpdateAssetProcessStatusRes{Updated: out.GetUpdated()}, nil
}

