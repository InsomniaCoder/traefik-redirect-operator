/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TraefikType type of traffic to expose in Traefik's ingress
type TraefikType string

const (
	// Public Any traffic from public space
	Public TraefikType = "public"
	// External Whitelisted public traffic
	External TraefikType = "external"
	// Internal Traffic running internally include intranet and internal api
	Internal TraefikType = "internal"
	// Private Cluster-only api
	Private TraefikType = "private"
)

// TODO: add marker for regex validation of host and redirect to url
// TraefikRedirectSpec defines the desired state of TraefikRedirect
type TraefikRedirectSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// TraefikType traffic type: public, external, internal, private
	TraefikType TraefikType `json:"traefik-type"`
	// TraefikHost host name to expose in Traefik's ingress
	// Valid values are:
	// - "public": Any traffic from public space
	// - "external":  Whitelisted public traffic
	// - "internal": Traffic running internally include intranet and internal api
	// - "private": Cluster-only api
	TraefikHost string `json:"traefik-host"`
	// RedirectTo url, or ip of ExternalName
	RedirectTo string `json:"redirect-to"`
	// Port port of the External name's server
	Port int `json:"port"`
}

// TraefikRedirectStatus defines the observed state of TraefikRedirect
type TraefikRedirectStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TraefikRedirect is the Schema for the traefikredirects API
type TraefikRedirect struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraefikRedirectSpec   `json:"spec,omitempty"`
	Status TraefikRedirectStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TraefikRedirectList contains a list of TraefikRedirect
type TraefikRedirectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TraefikRedirect `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TraefikRedirect{}, &TraefikRedirectList{})
}
