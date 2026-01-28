package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SemaTokenQuotaSpec defines the desired state of Token Quotas
type SemaTokenQuotaSpec struct {
    // HardLimit is the maximum tokens allowed per cycle
    // +kubebuilder:validation:Minimum=1
    HardLimit int64 `json:"hardLimit"`

    // SoftLimit triggers an alert/warning but doesn't freeze the agent
    // +optional
    SoftLimit int64 `json:"softLimit,omitempty"`

    // ModelMatch allows limiting specific models (e.g., "gpt-4*", "claude-3-5*")
    // +kubebuilder:default="*"
    ModelMatch string `json:"modelMatch,omitempty"`

    // ResetInterval defines how often the quota refills (Daily, Weekly, Monthly)
    // +kubebuilder:validation:Enum=Daily;Weekly;Monthly;Never
    ResetInterval string `json:"resetInterval"`
}

// SemaTokenQuotaStatus defines the observed state
type SemaTokenQuotaStatus struct {
    // TokensConsumed is the current meter reading
    TokensConsumed int64 `json:"tokensConsumed"`

    // LastResetTime tracks when the cycle was last cleared
    LastResetTime metav1.Time `json:"lastResetTime"`

    // Phase indicates if the quota is "Active", "Warning", or "Exhausted"
    Phase string `json:"phase"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// SemaTokenQuota is the Schema for the sematokenquotas API
type SemaTokenQuota struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   SemaTokenQuotaSpec   `json:"spec,omitempty"`
    Status SemaTokenQuotaStatus `json:"status,omitempty"`
}