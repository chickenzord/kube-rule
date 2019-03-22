package webhook

import (
	"context"
	"fmt"
	"net/http"

	kuberule "github.com/chickenzord/kube-rule/pkg/apis/kuberule/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type podRuleValidationHandler struct {
	client  client.Client
	decoder types.Decoder
}

// podRuleValidationHandler Implements admission.Handler.
var _ admission.Handler = &podRuleValidationHandler{}

// podRuleValidationHandler handle pod rules validation
func (a *podRuleValidationHandler) Handle(ctx context.Context, req types.Request) types.Response {
	podRule := &kuberule.PodRule{}

	err := a.decoder.Decode(req, podRule)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	if err := a.validatePodRuleFn(ctx, podRule); err != nil {
		return admission.ValidationResponse(false, err.Error())
	}

	return admission.ValidationResponse(true, "OK")
}

// validatePodRuleFn validates the given pod rule
func (a *podRuleValidationHandler) validatePodRuleFn(ctx context.Context, podRule *kuberule.PodRule) error {
	if podRule.Spec.ApplyOrder < 0 {
		return fmt.Errorf("podrule.spec.applyOrder must be >= 0")
	}

	return nil
}
