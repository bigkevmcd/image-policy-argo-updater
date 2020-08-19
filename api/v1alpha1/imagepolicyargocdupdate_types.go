/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImagePolicyArgoCDUpdateSpec defines the desired state of ImagePolicyArgoCDUpdate
type ImagePolicyArgoCDUpdateSpec struct {
	ApplicationRef corev1.LocalObjectReference  `json:"applicationRef"`
	ImagePolicyRef *corev1.LocalObjectReference `json:"imagePolicyRef"`
}

// ImagePolicyArgoCDUpdateStatus defines the observed state of ImagePolicyArgoCDUpdate
type ImagePolicyArgoCDUpdateStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ImagePolicyArgoCDUpdate is the Schema for the imagepolicyargocdupdates API
type ImagePolicyArgoCDUpdate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImagePolicyArgoCDUpdateSpec   `json:"spec,omitempty"`
	Status ImagePolicyArgoCDUpdateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImagePolicyArgoCDUpdateList contains a list of ImagePolicyArgoCDUpdate
type ImagePolicyArgoCDUpdateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImagePolicyArgoCDUpdate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImagePolicyArgoCDUpdate{}, &ImagePolicyArgoCDUpdateList{})
}
