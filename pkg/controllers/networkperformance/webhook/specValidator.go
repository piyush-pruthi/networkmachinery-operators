package webhook

import (
	"context"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

type SpecValidator struct {
	client  client.Client
	decoder *admission.Decoder
}

func (v *SpecValidator) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

func (v *SpecValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func (v *SpecValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	networkPerformanceTest := &v1alpha1.NetworkPerformanceTest{}

	err := v.decoder.Decode(req, networkPerformanceTest)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	allowed, reason, err := v.validateSpec(networkPerformanceTest)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.ValidationResponse(allowed, reason)
}

func (v *SpecValidator) validateSpec(npt *v1alpha1.NetworkPerformanceTest) (bool, string, error) {

	// can validate more parameters of spec here as the list grows

	if npt.Spec.Iterations != 1 {
		return false, "spec.iterations must be equal to 1 in current release", nil
	}

	return true, "", nil

}
