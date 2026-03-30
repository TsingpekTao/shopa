package review

import (
	reviewapi "github.com/TsingpekTao/shopa/review-svc/api/review"
	"github.com/TsingpekTao/shopa/review-svc/internal/service"
)

type ControllerV1 struct {
	review service.IReview
}

func NewV1() reviewapi.IReviewV1 {
	return &ControllerV1{review: service.Review()}
}
