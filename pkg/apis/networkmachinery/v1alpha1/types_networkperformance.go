package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkPerformanceTestPhase string

const (
	NetworkPerformanceTestPending   NetworkPerformanceTestPhase = "Pending"
	NetworkPerformanceTestFailed    NetworkPerformanceTestPhase = "Failed"
	NetworkPerformanceTestSucceeded NetworkPerformanceTestPhase = "Succeeded"
	NetworkPerformanceTestUnknown   NetworkPerformanceTestPhase = "Unknown"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkPerformanceTest represents a network performance test
type NetworkPerformanceTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkPerformanceTestSpec   `json:"spec,omitempty"`
	Status NetworkPerformanceTestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkPerformanceTestList is a list of network performance tests
type NetworkPerformanceTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkPerformanceTest `json:"items,omitempty"`
}

type NetworkPerformanceTestSpec struct {
	Type       string `json:"type"`
	Iterations int    `json:"iterations"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkPerformanceTestStatus contains status of test carried out by netperf
type NetworkPerformanceTestStatus struct {
	metav1.TypeMeta `json:",inline"`

	Phase  NetworkPerformanceTestPhase  `json:"phase"`
	Output NetworkPerformanceTestOutput `json:"output,omitempty"`
}

type NetworkPerformanceTestOutput struct {
	Bandwidth map[string]map[string]string `json:"bandwidth"`
}
