package notification

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/notification-svc/api/notification/v1"
)

type INotificationV1 interface {
	SendTemplateMessage(ctx context.Context, req *v1.SendTemplateMessageReq) (res *v1.SendTemplateMessageRes, err error)
	BatchSendTemplateMessage(ctx context.Context, req *v1.BatchSendTemplateMessageReq) (res *v1.BatchSendTemplateMessageRes, err error)
	GetDeliveryStatus(ctx context.Context, req *v1.GetDeliveryStatusReq) (res *v1.GetDeliveryStatusRes, err error)
	GetMyNotificationPreference(ctx context.Context, req *v1.GetMyNotificationPreferenceReq) (res *v1.GetMyNotificationPreferenceRes, err error)
	UpdateMyNotificationPreference(ctx context.Context, req *v1.UpdateMyNotificationPreferenceReq) (res *v1.UpdateMyNotificationPreferenceRes, err error)
	RetryFailedMessage(ctx context.Context, req *v1.RetryFailedMessageReq) (res *v1.RetryFailedMessageRes, err error)
}
