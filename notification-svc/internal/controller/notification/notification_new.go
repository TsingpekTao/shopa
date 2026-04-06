package notification

import (
	notificationapi "github.com/TsingpekTao/shopa/notification-svc/api/notification"
	"github.com/TsingpekTao/shopa/notification-svc/internal/service"
)

type ControllerV1 struct {
	notification service.INotification
}

func NewV1() notificationapi.INotificationV1 {
	return &ControllerV1{notification: service.Notification()}
}
