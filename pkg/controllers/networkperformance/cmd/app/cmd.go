package app

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkperformance/controller"
	networkmachineryhandlers "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkperformance/webhook"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
)

var log = logf.Log.WithName(controller.Name)

const (
	specValidationServerPath = "/validate-spec-v1alpha1-networkperformancetest"

	webhookServerPort = 9876
)

func NewNetworkPerformanceTestCmd(ctx context.Context) *cobra.Command {
	var (
		entryLog                      = log.WithName("networkperformance-test-cmd")
		retryDuration                 = 100 * time.Millisecond
		networkPerformanceTestCmdOpts = NetworkPerformanceTestCmdOpts{
			ConfigFlags: genericclioptions.NewConfigFlags(true),
			LeaderElectionOptions: controllers.LeaderElectionOptions{
				LeaderElection:          true,
				LeaderElectionNamespace: "default",
				LeaderElectionID:        utils.LeaderElectionNameID(controller.Name),
			},
			leaderElectionRetryPeriod: &retryDuration,
		}
	)

	cmd := &cobra.Command{
		Use: "networkperformance-test-controller",
		Run: func(cmd *cobra.Command, args []string) {
			mgrOptions := &manager.Options{
				Port: webhookServerPort,
			}
			mgr, err := manager.New(networkPerformanceTestCmdOpts.InitConfig(), *networkPerformanceTestCmdOpts.InjectRetryOptions(networkPerformanceTestCmdOpts.InjectLeaderElectionOpts(mgrOptions)))
			if err != nil {
				utils.LogErrAndExit(err, "Could not instantiate manager")
			}
			if err := install.AddToScheme(mgr.GetScheme()); err != nil {
				utils.LogErrAndExit(err, "Could not update manager scheme")
			}

			entryLog.Info("Setting up webhook server")
			admissionServer := mgr.GetWebhookServer()

			entryLog.Info("registering webhooks to the webhook server")
			admissionServer.Register(specValidationServerPath, &webhook.Admission{Handler: &networkmachineryhandlers.SpecValidator{}})

			if err := controller.Add(mgr); err != nil {
				utils.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				utils.LogErrAndExit(err, "Error running manager")
			}
		},
	}
	networkPerformanceTestCmdOpts.AddAllFlags(cmd.Flags())
	return cmd
}
