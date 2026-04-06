package v1

import (
	pb "github.com/TsingpekTao/shopa/notification-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type SendTemplateMessageReq struct {
	g.Meta `path:"/v1/notifications/template/send" method:"post" tags:"Notification" summary:"Send template message"`
	pb.SendTemplateMessageReq
}

type SendTemplateMessageRes = pb.SendTemplateMessageRes

type BatchSendTemplateMessageReq struct {
	g.Meta `path:"/v1/notifications/template/batch-send" method:"post" tags:"Notification" summary:"Batch send template messages"`
	pb.BatchSendTemplateMessageReq
}

type BatchSendTemplateMessageRes = pb.BatchSendTemplateMessageRes

type GetDeliveryStatusReq struct {
	g.Meta         `path:"/v1/notifications/{notificationNo}/delivery-status" method:"get" tags:"Notification" summary:"Get notification delivery status"`
	NotificationNo string `json:"notificationNo" in:"path" v:"required#notificationNo is required"`
}

type GetDeliveryStatusRes = pb.GetDeliveryStatusRes

type GetMyNotificationPreferenceReq struct {
	g.Meta `path:"/v1/me/notifications/preference" method:"get" tags:"Notification" summary:"Get my notification preference"`
}

type GetMyNotificationPreferenceRes = pb.GetMyNotificationPreferenceRes

type UpdateMyNotificationPreferenceReq struct {
	g.Meta `path:"/v1/me/notifications/preference" method:"patch" tags:"Notification" summary:"Update my notification preference"`
	pb.UpdateMyNotificationPreferenceReq
}

type UpdateMyNotificationPreferenceRes = pb.UpdateMyNotificationPreferenceRes

type RetryFailedMessageReq struct {
	g.Meta         `path:"/v1/internal/notifications/{notificationNo}/retry" method:"post" tags:"NotificationInternal" summary:"Retry failed notification"`
	NotificationNo string `json:"notificationNo" in:"path" v:"required#notificationNo is required"`
	ReasonCode     string `json:"reasonCode"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type RetryFailedMessageRes = pb.RetryFailedMessageRes
