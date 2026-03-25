package me

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
)

type IMeV1 interface {
	GetOverview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error)
}
