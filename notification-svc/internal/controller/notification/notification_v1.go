package notification

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/notification-svc/api/notification/v1"
	pb "github.com/TsingpekTao/shopa/notification-svc/api/v1"
)

func (c *ControllerV1) SendTemplateMessage(ctx context.Context, req *httpv1.SendTemplateMessageReq) (*httpv1.SendTemplateMessageRes, error) {
	return c.notification.SendTemplateMessage(ctx, &req.SendTemplateMessageReq)
}

func (c *ControllerV1) BatchSendTemplateMessage(ctx context.Context, req *httpv1.BatchSendTemplateMessageReq) (*httpv1.BatchSendTemplateMessageRes, error) {
	return c.notification.BatchSendTemplateMessage(ctx, &req.BatchSendTemplateMessageReq)
}

func (c *ControllerV1) GetDeliveryStatus(ctx context.Context, req *httpv1.GetDeliveryStatusReq) (*httpv1.GetDeliveryStatusRes, error) {
	return c.notification.GetDeliveryStatus(ctx, &pb.GetDeliveryStatusReq{
		NotificationNo: req.NotificationNo,
	})
}

func (c *ControllerV1) GetMyNotificationPreference(ctx context.Context, req *httpv1.GetMyNotificationPreferenceReq) (*httpv1.GetMyNotificationPreferenceRes, error) {
	return c.notification.GetMyNotificationPreference(ctx, &pb.GetMyNotificationPreferenceReq{})
}

func (c *ControllerV1) UpdateMyNotificationPreference(ctx context.Context, req *httpv1.UpdateMyNotificationPreferenceReq) (*httpv1.UpdateMyNotificationPreferenceRes, error) {
	return c.notification.UpdateMyNotificationPreference(ctx, &req.UpdateMyNotificationPreferenceReq)
}

func (c *ControllerV1) RetryFailedMessage(ctx context.Context, req *httpv1.RetryFailedMessageReq) (*httpv1.RetryFailedMessageRes, error) {
	return c.notification.RetryFailedMessage(ctx, &pb.RetryFailedMessageReq{
		NotificationNo: req.NotificationNo,
		ReasonCode:     req.ReasonCode,
		IdempotencyKey: req.IdempotencyKey,
	})
}
