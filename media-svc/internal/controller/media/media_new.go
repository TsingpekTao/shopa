package media

import (
	mediaapi "github.com/TsingpekTao/shopa/media-svc/api/media"
	"github.com/TsingpekTao/shopa/media-svc/internal/service"
)

type ControllerV1 struct {
	svc service.IMedia
}

// NewV1 创建 media HTTP 控制器。
func NewV1() mediaapi.IMediaV1 {
	return &ControllerV1{svc: service.Media()}
}
