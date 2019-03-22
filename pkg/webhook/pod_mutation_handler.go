package webhook

import (
	"context"
	"net/http"
	"sort"

	kuberule "github.com/chickenzord/kube-rule/pkg/apis/kuberule/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type podMutationHandler struct {
	client  client.Client
	decoder types.Decoder
}

var _ admission.Handler = &podMutationHandler{} // Implements admission.Handler.
var log = logf.Log.WithName("webhook.pods.mutation")

// podMutationHandler try to mutate every incoming pods based on rules
func (a *podMutationHandler) Handle(ctx context.Context, req types.Request) types.Response {
	// Decode request and make a clone to mutate
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	clone := pod.DeepCopy()

	log.Info("receiving pod to handle",
		"request.namespace", req.AdmissionRequest.Namespace,
		"request.operation", req.AdmissionRequest.Operation,
		"pod.name", pod.Name,
		"pod.generateName", pod.GenerateName,
	)

	// Getting matching rules
	podRuleList := &kuberule.PodRuleList{}
	listOptions := client.InNamespace(req.AdmissionRequest.Namespace)
	if err := a.client.List(ctx, listOptions, podRuleList); err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	if podRuleList.Items == nil {
		return admission.PatchResponse(pod, clone)
	}
	sort.Slice(podRuleList.Items, func(i, j int) bool {
		return podRuleList.Items[i].Spec.ApplyOrder < podRuleList.Items[j].Spec.ApplyOrder
	})

	for _, rule := range podRuleList.Items {
		// check matching pods, skip if doesn't match
		rulePodSelector := labels.Set(rule.Spec.Selector.MatchLabels).AsSelector()
		if !rulePodSelector.Matches(labels.Set(pod.Labels)) {
			continue
		}

		// apply mutations
		err = a.mutatePodsFn(ctx, clone, rule)
		if err != nil {
			return admission.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	// admission.PatchResponse generates a Response containing patches.
	return admission.PatchResponse(pod, clone)
}

// mutatePodsFn mutates the given pod
func (a *podMutationHandler) mutatePodsFn(ctx context.Context, pod *corev1.Pod, rule kuberule.PodRule) error {
	if rule.Spec.Mutations == nil {
		return nil
	}

	// apply annotations
	if rule.Spec.Mutations.Annotations != nil {
		log.Info("applying annotations to pod",
			"pod.name", pod.Name,
			"pod.generateName", pod.GenerateName,
			"rule.name", rule.ObjectMeta.Name,
		)
		if pod.Annotations == nil {
			pod.Annotations = map[string]string{}
		}

		for key, val := range rule.Spec.Mutations.Annotations {
			pod.Annotations[key] = val
		}
	}

	// TODO: add more mutations here

	return nil
}
