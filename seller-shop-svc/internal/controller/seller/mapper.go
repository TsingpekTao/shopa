package seller

import (
	"context"

	sellerv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/seller/v1"
	pb "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// withRequestMetadata 灏?HTTP 璇锋眰澶撮€忎紶涓?gRPC metadata锛屽鐢?logic 灞傚凡鏈夐壌鏉?瀹¤璇诲彇閫昏緫銆?
func withRequestMetadata(ctx context.Context) context.Context {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ctx
	}

	md := metadata.Pairs(
		"x-user-id", req.GetHeader("X-User-Id"),
		"x-operator-user-id", req.GetHeader("X-Operator-User-Id"),
		"x-request-id", req.GetHeader("X-Request-Id"),
		"x-idempotency-key", req.GetHeader("X-Idempotency-Key"),
		"x-forwarded-for", req.GetHeader("X-Forwarded-For"),
		"x-client-ip", req.GetClientIp(),
		"x-user-agent", req.GetHeader("User-Agent"),
		"authorization", req.GetHeader("Authorization"),
		"x-access-token", req.GetHeader("X-Access-Token"),
	)
	return metadata.NewIncomingContext(ctx, md)
}

func toPBApplicationStatuses(in []int32) []pb.ApplicationStatus {
	if len(in) == 0 {
		return nil
	}
	out := make([]pb.ApplicationStatus, 0, len(in))
	for _, v := range in {
		out = append(out, pb.ApplicationStatus(v))
	}
	return out
}

func toFieldMask(in []string) *fieldmaskpb.FieldMask {
	if len(in) == 0 {
		return nil
	}
	return &fieldmaskpb.FieldMask{Paths: in}
}

func toHTTPListMyApplicationsRes(in *pb.ListMyApplicationsRes) *sellerv1.ListMyApplicationsRes {
	if in == nil {
		return &sellerv1.ListMyApplicationsRes{}
	}
	return &sellerv1.ListMyApplicationsRes{
		Applications: in.GetApplications(),
		Page:         in.GetPage(),
		PageSize:     in.GetPageSize(),
		Total:        in.GetTotal(),
	}
}

func toHTTPListApplicationsRes(in *pb.ListApplicationsRes) *sellerv1.ListApplicationsRes {
	if in == nil {
		return &sellerv1.ListApplicationsRes{}
	}
	return &sellerv1.ListApplicationsRes{
		Applications: in.GetApplications(),
		Page:         in.GetPage(),
		PageSize:     in.GetPageSize(),
		Total:        in.GetTotal(),
	}
}

