package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run "operator-sdk generate k8s && operator-sdk generate openapi" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

// CustomerImageBuildSpec defines the desired state of CustomerImageBuild
// +k8s:openapi-gen=true
type CustomerImageBuildSpec struct {
	RepoURL string `json:"repoURL"`
	GitRef string `json:"gitRef"`
	ImageRepo string `json:"imageRepo"`
	ImageTag string `json:"imageTag"`

	Retries int32 `json:"retries,omitempty"` // +optional
}

// CustomerImageBuildStatus defines the observed state of CustomerImageBuild
// +k8s:openapi-gen=true
type CustomerImageBuildStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomerImageBuild is the Schema for the customerimagebuilds API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CustomerImageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomerImageBuildSpec   `json:"spec,omitempty"`
	Status CustomerImageBuildStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomerImageBuildList contains a list of CustomerImageBuild
type CustomerImageBuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomerImageBuild `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomerImageBuild{}, &CustomerImageBuildList{})
}
