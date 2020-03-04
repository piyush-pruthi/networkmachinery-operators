package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
)

const (
	LogKey          = "NetworkPerformanceTest"
	FinalizerName   = "networkmachinery.io/networkperformance"
	reconcilePeriod = 5 * time.Second
)

type ReconcileNetworkPerformanceTest struct {
	config   *rest.Config
	logger   logr.Logger
	client   client.Client
	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// InjectConfig implements inject.Config.
func (r *ReconcileNetworkPerformanceTest) InjectConfig(config *rest.Config) error {
	r.config = config
	return nil
}

func (r *ReconcileNetworkPerformanceTest) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkPerformanceTest) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkPerformanceTest) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	networkPerformanceTest := &v1alpha1.NetworkPerformanceTest{}

	if err := r.client.Get(r.ctx, request.NamespacedName, networkPerformanceTest); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return apimachinery.ReconcileErr(err)
	}

	if networkPerformanceTest.DeletionTimestamp != nil {
		return r.delete(networkPerformanceTest)
	}

	r.logger.Info("Reconciling Network Performance Test", "Name", networkPerformanceTest.Name)
	return r.reconcile(networkPerformanceTest)
}

func (r *ReconcileNetworkPerformanceTest) reconcile(networkPerformanceTest *v1alpha1.NetworkPerformanceTest) (reconcile.Result, error) {

	var err error
	// Finalizer is needed to make sure to clean up installed debugging tools after reconciliation
	if err = apimachinery.EnsureFinalizer(r.ctx, r.client, FinalizerName, networkPerformanceTest); err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err = r.reconcilePerfTest(networkPerformanceTest)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{
		RequeueAfter: reconcilePeriod,
	}, nil
}

func (r *ReconcileNetworkPerformanceTest) delete(networkPerformanceTest *v1alpha1.NetworkPerformanceTest) (reconcile.Result, error) {

	hasFinalizer, err := apimachinery.HasFinalizer(networkPerformanceTest, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return apimachinery.ReconcileErr(err)
	}

	if !hasFinalizer {
		r.logger.Info("Deleting NetworkPerformanceTest causes a no-op as there is no finalizer.", LogKey, networkPerformanceTest.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Initiating NetworkPerformanceTest deletion ", LogKey, networkPerformanceTest.Name)
	r.recorder.Event(networkPerformanceTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting network performance test")

	// delete perf job
	err = r.deletePerfTest(networkPerformanceTest.Name)
	if err != nil {
		r.logger.Error(err, "Error deleting cascaded resource", LogKey, networkPerformanceTest.Name)
		return apimachinery.ReconcileErr(err)
	}

	// delete finalizer
	if err = apimachinery.DeleteFinalizer(r.ctx, r.client, FinalizerName, networkPerformanceTest); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkPerformance resource", LogKey, networkPerformanceTest.Name)
		return apimachinery.ReconcileErr(err)
	}

	r.logger.Info("Deletion successful ", LogKey, networkPerformanceTest.Name)
	r.recorder.Event(networkPerformanceTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Network Performance Test Deleted!")

	return reconcile.Result{}, nil
}
