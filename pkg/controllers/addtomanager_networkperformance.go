package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkperformance/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
