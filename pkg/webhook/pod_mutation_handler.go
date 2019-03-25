package webhook

import (
	"context"
	"net/http"
	"sort"

	kuberule "github.com/chickenzord/kube-rule/pkg/apis/kuberule/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type podMutationHandler struct {
	client  client.Client
	decoder types.Decoder
}

var _ admission.Handler = &podMutationHandler{} // Implements admission.Handler.

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

	// Get matching rules sorted by ApplyOrder
	podRuleList := &kuberule.PodRuleList{}
	listOptions := client.InNamespace(req.AdmissionRequest.Namespace)
	if err := a.client.List(ctx, listOptions, podRuleList); err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	sort.Slice(podRuleList.Items, func(i, j int) bool {
		return podRuleList.Items[i].Spec.ApplyOrder < podRuleList.Items[j].Spec.ApplyOrder
	})
	for _, rule := range podRuleList.Items {
		// check matching pods, skip if doesn't match
		podSelector := labels.Set(rule.Spec.Selector.MatchLabels).AsSelector()
		if !podSelector.Matches(labels.Set(pod.Labels)) {
			continue
		}

		// apply mutations
		err = a.mutatePodsFn(ctx, clone, rule)
		if err != nil {
			return admission.ErrorResponse(http.StatusInternalServerError, err)
		}
	}

	// create patches
	return admission.PatchResponse(pod, clone)
}

// mutatePodsFn mutates the given pod
func (a *podMutationHandler) mutatePodsFn(ctx context.Context, pod *corev1.Pod, rule kuberule.PodRule) error {
	// for convenience
	mutations := rule.Spec.Mutations

	// apply annotations
	log.Info("applying mutations to pod",
		"mutations.annotations", mutations.Annotations,
		"mutations.nodeSelector", mutations.NodeSelector,
		"pod.name", pod.Name,
		"pod.generateName", pod.GenerateName,
		"pod.annotations", pod.Annotations,
		"rule.name", rule.ObjectMeta.Name,
	)

	// apply annotations by merging with existing
	if mutations.Annotations == nil {
		mutations.Annotations = map[string]string{}
	}
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	for key, val := range mutations.Annotations {
		pod.Annotations[key] = val
	}

	// apply nodeSelector only if not already set
	if pod.Spec.NodeSelector == nil {
		pod.Spec.NodeSelector = map[string]string{}
	}
	if len(pod.Spec.NodeSelector) == 0 {
		for key, val := range mutations.NodeSelector {
			pod.Spec.NodeSelector[key] = val
		}
	}

	// TODO: add more mutations here

	return nil
}
