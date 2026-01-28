package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SemaPolicySpec defines the rules of the road for an Agent
type SemaPolicySpec struct {
    // Selector identifies which agents (pods) this policy applies to
    Selector metav1.LabelSelector `json:"selector"`

    // Rules is a list of behaviors to govern
    Rules []PolicyRule `json:"rules"`

    // OnViolation defines the global behavior if no rule matches but risk is detected
    // +kubebuilder:default="DENY"
    DefaultAction string `json:"defaultAction,omitempty"`
}

type PolicyRule struct {
    // Name of the rule (e.g., "protect-production-storage")
    Name string `json:"name"`

    // IntentMatches: Semantic keywords or regex for the LLM prompt/thought
    // +optional
    IntentMatches []string `json:"intentMatches,omitempty"`

    // ToolMatches: Specific Kubernetes or external API tools the agent tries to call
    // +optional
    ToolMatches []string `json:"toolMatches,omitempty"`

    // RiskLevel: Evaluated by the Semantic Waypoint (Low, Medium, High, Critical)
    // +kubebuilder:validation:Enum=Low;Medium;High;Critical
    RiskLevel string `json:"riskLevel"`

    // Action: What to do? (ALLOW, DENY, PAUSE)
    // PAUSE triggers the CRIU Freeze logic
    // +kubebuilder:validation:Enum=ALLOW;DENY;PAUSE
    Action string `json:"action"`

    // PauseSettings: Only used if Action is PAUSE
    // +optional
    PauseSettings *PauseSettings `json:"pauseSettings,omitempty"`
}

type PauseSettings struct {
    // Timeout: Duration before Auto-Reject (e.g., "15m", "1h")
    // +kubebuilder:default="30m"
    Timeout string `json:"timeout,omitempty"`

    // Notify: Where to send the 'kubectl sema' alert (e.g., Slack Webhook URL)
    // +optional
    Notify string `json:"notify,omitempty"`
}

// +kubebuilder:object:root=true
// SemaPolicy is the Schema for the semapolicies API
type SemaPolicy struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec SemaPolicySpec `json:"spec,omitempty"`
}