package controller

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"

	networkmachineryv1alpha1 "github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
)

func (r *ReconcileNetworkPerformanceTest) reconcilePerfTest(networkPerformanceTest *networkmachineryv1alpha1.NetworkPerformanceTest) error {

	// no need to do anything as the test is already in succeeded state
	if networkPerformanceTest.Status.Phase == networkmachineryv1alpha1.NetworkPerformanceTestSucceeded {
		return nil
	}

	pod, err := apimachinery.GetPod(r.config, networkPerformanceTest.Name)

	if err != nil {
		if errors.IsNotFound(err) {
			// create pod
			image := "pruthi/private-workspace:k8s-netperf"
			args := []string{"--image=pruthi/private-workspace:nptests", fmt.Sprintf("--iterations=%v", networkPerformanceTest.Spec.Iterations)}
			err := apimachinery.CreatePod(r.config, networkPerformanceTest.Name, image, args, Name)
			return err
		}
		return err
	}

	// pod already exists - check its status
	switch pod.Status.Phase {

	case corev1.PodRunning:

		statusUpdated, err := r.updateStatus(networkPerformanceTest)
		if err != nil {
			return err
		}
		if statusUpdated {
			return r.deletePerfTest(networkPerformanceTest.Name)
		}
		return nil

	case corev1.PodSucceeded:

		return nil

	case corev1.PodFailed:

		// as pod's restart policy is set to "onFailure" - pod will be restarted after it reaches this phase
		// so no need to handle this, as we will keep on trying to take the test status to Succeeded.
		fmt.Println("Pod in failed state : Will automatically be restarted by K8s")
		return nil

		// default will handle PodPending and PodUnknown phases
	default:

		return nil

	}

}

func (r *ReconcileNetworkPerformanceTest) deletePerfTest(name string) error {

	err := apimachinery.DeletePod(r.config, name)

	// if resouce does not exist, return nil error
	if errors.IsNotFound(err) {
		return nil
	}
	return err

}

func (r *ReconcileNetworkPerformanceTest) updateStatus(networkPerformanceTest *networkmachineryv1alpha1.NetworkPerformanceTest) (bool, error) {

	var status = &networkmachineryv1alpha1.NetworkPerformanceTestStatus{}

	logs := apimachinery.GetLogs(r.config, networkPerformanceTest.Name)
	podLogs, err := logs.Stream()
	if err != nil {
		return false, err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return false, err
	}
	str := buf.String()
	begin := strings.Index(str, "MSS")
	end := strings.Index(str, "Test concluded")

	if begin == -1 || end == -1 {
		return false, nil
	}

	status.Phase = networkmachineryv1alpha1.NetworkPerformanceTestSucceeded

	str = str[begin : end-2]
	utils.ParseNetPerfOutput(str, &status.Output)

	err = apimachinery.TryUpdateStatus(r.ctx, retry.DefaultBackoff, r.client, networkPerformanceTest, func() error {
		networkPerformanceTest.Status = *status
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
