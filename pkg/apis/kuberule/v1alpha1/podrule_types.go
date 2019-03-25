package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodRuleMutations defines mutations to be applied on the selected pods
type PodRuleMutations struct {

	// Annotations to be merged with selected pods' existing annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// NodeSelector to be added to selected pods according to nodeSelectorStrategy
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// PodRuleSpec defines the desired state of PodRule
type PodRuleSpec struct {
	// Arbitrary number to define ordering of multiple rules matching same pods.
	// Higher number will be applied later, but might override mutations of smaller number.
	ApplyOrder int32 `json:"applyOrder"`

	// Label selector for pods
	Selector metav1.LabelSelector `json:"selector"`

	// Mutations to be done on the selected pods
	Mutations PodRuleMutations `json:"mutations,omitempty"`
}

// PodRuleStatus defines the observed state of PodRule
type PodRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodRule is the Schema for the podrules API
// +k8s:openapi-gen=true
type PodRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodRuleSpec   `json:"spec,omitempty"`
	Status PodRuleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodRuleList contains a list of PodRule
type PodRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodRule{}, &PodRuleList{})
}
