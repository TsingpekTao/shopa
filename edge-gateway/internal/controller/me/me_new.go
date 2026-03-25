package me

import (
	"github.com/TsingpekTao/shopa/edge-gateway/api/me"
)

type ControllerV1 struct{}

func NewV1() me.IMeV1 {
	return &ControllerV1{}
}
