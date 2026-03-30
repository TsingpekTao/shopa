package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/notification-svc/api/v1"
	"github.com/TsingpekTao/shopa/notification-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 承载 notification 的三类 gRPC 服务。
type Controller struct {
	v1.UnimplementedNotificationServiceServer
	v1.UnimplementedUserPreferenceServiceServer
	v1.UnimplementedInternalNotificationServiceServer
}

// Register 注册 notification gRPC 服务。
func Register(s *grpcx.GrpcServer) {
	// 注册通知发送服务。
	v1.RegisterNotificationServiceServer(s.Server, &Controller{})
	// 注册用户偏好服务。
	v1.RegisterUserPreferenceServiceServer(s.Server, &Controller{})
	// 注册内部通知服务。
	v1.RegisterInternalNotificationServiceServer(s.Server, &Controller{})
}

// SendTemplateMessage 调用 service 发送模板消息。
func (*Controller) SendTemplateMessage(ctx context.Context, req *v1.SendTemplateMessageReq) (*v1.SendTemplateMessageRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().SendTemplateMessage(ctx, req)
}

// BatchSendTemplateMessage 调用 service 批量发送模板消息。
func (*Controller) BatchSendTemplateMessage(ctx context.Context, req *v1.BatchSendTemplateMessageReq) (*v1.BatchSendTemplateMessageRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().BatchSendTemplateMessage(ctx, req)
}

// GetDeliveryStatus 调用 service 查询投递状态。
func (*Controller) GetDeliveryStatus(ctx context.Context, req *v1.GetDeliveryStatusReq) (*v1.GetDeliveryStatusRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().GetDeliveryStatus(ctx, req)
}

// GetMyNotificationPreference 调用 service 查询当前用户偏好。
func (*Controller) GetMyNotificationPreference(ctx context.Context, req *v1.GetMyNotificationPreferenceReq) (*v1.GetMyNotificationPreferenceRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().GetMyNotificationPreference(ctx, req)
}

// UpdateMyNotificationPreference 调用 service 更新当前用户偏好。
func (*Controller) UpdateMyNotificationPreference(ctx context.Context, req *v1.UpdateMyNotificationPreferenceReq) (*v1.UpdateMyNotificationPreferenceRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().UpdateMyNotificationPreference(ctx, req)
}

// RetryFailedMessage 调用 service 执行失败消息重试。
func (*Controller) RetryFailedMessage(ctx context.Context, req *v1.RetryFailedMessageReq) (*v1.RetryFailedMessageRes, error) {
	// 委派给 service 层执行业务逻辑。
	return service.Notification().RetryFailedMessage(ctx, req)
}
