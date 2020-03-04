package main

import (
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkperformance/cmd/app"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
)

func main() {
	log.SetLogger(log.ZapLogger(false))
	cmd := app.NewNetworkPerformanceTestCmd(utils.SetupSignalHandlerContext())

	if err := cmd.Execute(); err != nil {
		utils.LogErrAndExit(err, "error executing main controller command")
	}
}
