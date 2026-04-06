package points

import (
	pointsapi "github.com/TsingpekTao/shopa/points-svc/api/points"
	"github.com/TsingpekTao/shopa/points-svc/internal/service"
)

type ControllerV1 struct {
	points service.IPoints
}

func NewV1() pointsapi.IPointsV1 {
	return &ControllerV1{points: service.Points()}
}
