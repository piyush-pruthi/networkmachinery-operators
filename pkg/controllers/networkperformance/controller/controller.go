package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) *ReconcileNetworkPerformanceTest {
	return &ReconcileNetworkPerformanceTest{
		logger:   log.Log.WithName("networkperformance-test-controller"),
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetEventRecorderFor(Name)}
}

// DefaultPredicates returns the default predicates for an infrastructure reconciler.
func DefaultPredicates() []predicate.Predicate {
	return []predicate.Predicate{GenerationChangedPredicate()}
}

// Add creates a new NetworkPerformance Controller and adds it to the Manager
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr), DefaultPredicates())
}

func add(mgr manager.Manager, r reconcile.Reconciler, predicates []predicate.Predicate) error {
	ctrl, err := controller.New(Name, mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: 5})
	if err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.NetworkPerformanceTest{}}, &handler.EnqueueRequestForObject{}, predicates...); err != nil {
		return err
	}

	return nil
}
