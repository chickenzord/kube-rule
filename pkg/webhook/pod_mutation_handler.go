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
	admissiontypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type podMutationHandler struct {
	client  client.Client
	decoder admissiontypes.Decoder
}

var _ admission.Handler = &podMutationHandler{} // Implements admission.Handler.

// podMutationHandler try to mutate every incoming pods based on rules
func (a *podMutationHandler) Handle(ctx context.Context, req admissiontypes.Request) admissiontypes.Response {
	// Decode request and make a clone to mutate
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	clone := pod.DeepCopy()

	log.Info("handling pod",
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
	log.Info("applying mutations to pod",
		"pod", &pod,
		"rule", rule,
	)
	mutations := rule.Spec.Mutations

	// merge with existing annotations
	if mutations.Annotations == nil {
		mutations.Annotations = map[string]string{}
	}
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	for key, val := range mutations.Annotations {
		pod.Annotations[key] = val
	}

	// apply affinity only if not already exists
	if pod.Spec.Affinity != nil && mutations.Affinity != nil {
		pod.Spec.Affinity = mutations.Affinity
	}

	// apply nodeSelector only if not already exists
	if pod.Spec.NodeSelector == nil {
		pod.Spec.NodeSelector = map[string]string{}
	}
	if len(pod.Spec.NodeSelector) == 0 {
		for key, val := range mutations.NodeSelector {
			pod.Spec.NodeSelector[key] = val
		}
	}

	// append to existing tolerations
	if pod.Spec.Tolerations == nil {
		pod.Spec.Tolerations = []corev1.Toleration{}
	}
	if mutations.Tolerations != nil {
		for _, toleration := range mutations.Tolerations {
			pod.Spec.Tolerations = append(pod.Spec.Tolerations, toleration)
		}
	}

	// append imagePullSecrets
	if pod.Spec.ImagePullSecrets == nil {
		pod.Spec.ImagePullSecrets = []corev1.LocalObjectReference{}
	}
	if mutations.ImagePullSecrets != nil {
		for _, secret := range mutations.ImagePullSecrets {
			found := false
			for _, existing := range pod.Spec.ImagePullSecrets {
				if secret.Name == existing.Name {
					found = true
					break
				}
			}

			if !found {
				pod.Spec.ImagePullSecrets = append(pod.Spec.ImagePullSecrets, secret)
			}
		}
	}

	// TODO: add more mutations here

	return nil
}
