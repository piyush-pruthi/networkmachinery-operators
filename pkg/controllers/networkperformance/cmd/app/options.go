package app

import (
	"time"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkperformance/controller"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
)

type NetworkPerformanceTestCmdOpts struct {
	disableWebhookConfigInstaller bool
	leaderElectionRetryPeriod     *time.Duration
	*genericclioptions.ConfigFlags
	controllers.LeaderElectionOptions
}

func (npt *NetworkPerformanceTestCmdOpts) InjectLeaderElectionOpts(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.LeaderElectionID = npt.LeaderElectionID
	mgrOpts.LeaderElectionNamespace = npt.LeaderElectionNamespace
	mgrOpts.LeaderElection = npt.LeaderElection
	return mgrOpts
}

func (npt *NetworkPerformanceTestCmdOpts) InjectRetryOptions(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.RetryPeriod = npt.leaderElectionRetryPeriod
	return mgrOpts
}

func (npt *NetworkPerformanceTestCmdOpts) InitConfig() *rest.Config {
	config, err := npt.ToRESTConfig()

	if err != nil {
		utils.LogErrAndExit(err, "Error Getting Rest Api Configuration")
	}
	config.UserAgent = controller.Name
	return config
}

func (npt *NetworkPerformanceTestCmdOpts) AddWebHookFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&npt.disableWebhookConfigInstaller, "disable-webhook-config-installer", true,
		"disable the installer in the webhook server, so it won't install webhook configuration resources during bootstrapping")
}

func (npt *NetworkPerformanceTestCmdOpts) AddAllFlags(flags *pflag.FlagSet) {
	npt.AddWebHookFlags(flags)
	npt.AddFlags(flags)
}
