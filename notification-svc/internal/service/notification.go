// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/notification-svc/api/v1"
)

type (
	INotification interface {
		// SendTemplateMessage 发送单条模板消息。
		SendTemplateMessage(ctx context.Context, req *v1.SendTemplateMessageReq) (*v1.SendTemplateMessageRes, error)
		// BatchSendTemplateMessage 批量发送模板消息。
		BatchSendTemplateMessage(ctx context.Context, req *v1.BatchSendTemplateMessageReq) (*v1.BatchSendTemplateMessageRes, error)
		// GetDeliveryStatus 查询单条通知投递状态。
		GetDeliveryStatus(ctx context.Context, req *v1.GetDeliveryStatusReq) (*v1.GetDeliveryStatusRes, error)
		// GetMyNotificationPreference 查询当前用户通知偏好。
		GetMyNotificationPreference(ctx context.Context, req *v1.GetMyNotificationPreferenceReq) (*v1.GetMyNotificationPreferenceRes, error)
		// UpdateMyNotificationPreference 更新当前用户通知偏好。
		UpdateMyNotificationPreference(ctx context.Context, req *v1.UpdateMyNotificationPreferenceReq) (*v1.UpdateMyNotificationPreferenceRes, error)
		// RetryFailedMessage 重试失败消息。
		RetryFailedMessage(ctx context.Context, req *v1.RetryFailedMessageReq) (*v1.RetryFailedMessageRes, error)
	}
)

var (
	localNotification INotification
)

func Notification() INotification {
	if localNotification == nil {
		panic("implement not found for interface INotification, forgot register?")
	}
	return localNotification
}

func RegisterNotification(i INotification) {
	localNotification = i
}
