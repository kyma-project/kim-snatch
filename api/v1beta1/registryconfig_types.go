/*
Copyright 2024.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RegistryCacheConfigSpec defines the desired state of RegistryCacheConfig.
type RegistryCacheConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// RegistryCaches stores configuration for registry caches.
	RegistryCaches RegistryCache `json:"caches,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RegistryCacheConfig is the Schema for the registrycacheconfigs API.
type RegistryCacheConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistryCacheConfigSpec   `json:"spec,omitempty"`
	Status RegistryCacheConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistryCacheConfigList contains a list of CustomConfig.
type RegistryCacheConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RegistryCacheConfig `json:"items"`
}

type State string

const (
	ReadyState   State = "Ready"
	ErrorState   State = "Error"
	PendingState State = "Pending"
)

type ConditionReason string

const (
	ConditionReasonRegistryCacheConfigured     ConditionReason = "RegistryCacheConfigured"
	ConditionReasonFailedToGetSecret           ConditionReason = "FailedToGetCredentialsSecret"
	ConditionReasonSecretHasIncorrectStructure ConditionReason = "SecretHasIncorrectStructure"
	ConditionReasonFailedToResolveRegistryURL  ConditionReason = "FailedToResolveRegistryURL"
)

type RegistryCacheConfigStatus struct {
	// State signifies current state of Runtime
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Pending;Ready;Terminating;Failed
	State State `json:"state,omitempty"`

	// List of status conditions to indicate the status of a ServiceInstance.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ProvisioningCompleted indicates if the initial provisioning of the cluster is completed
	ProvisioningCompleted bool `json:"provisioningCompleted,omitempty"`
}

func init() {
	SchemeBuilder.Register(&RegistryCacheConfig{}, &RegistryCacheConfigList{})
}
