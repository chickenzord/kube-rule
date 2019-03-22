package webhook

import (
	"context"
	"net/http"

	kuberule "github.com/chickenzord/kube-rule/pkg/apis/kuberule/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type podRuleMutationHandler struct {
	client  client.Client
	decoder types.Decoder
}

var defaultPodRuleSpecApplyOrder int32 = 1

// podRuleMutationHandler Implements admission.Handler.
var _ admission.Handler = &podRuleMutationHandler{}

// podRuleMutationHandler adds an annotation to every incoming pods.
func (a *podRuleMutationHandler) Handle(ctx context.Context, req types.Request) types.Response {
	podRule := &kuberule.PodRule{}

	err := a.decoder.Decode(req, podRule)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	copy := podRule.DeepCopy()
	err = a.mutatePodRuleFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}

	// admission.PatchResponse generates a Response containing patches.
	return admission.PatchResponse(podRule, copy)
}

// mutatePodRuleFn add an annotation to the given pod rule
func (a *podRuleMutationHandler) mutatePodRuleFn(ctx context.Context, podRule *kuberule.PodRule) error {
	return nil
}
