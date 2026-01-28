package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record" // NEW: For Events
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

type SemaReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder // NEW: For communicating with the user
}

func (r *SemaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 1. Detect Intent to Freeze
	if pod.Annotations["semamesh.io/action"] == "PAUSE" {
		r.Recorder.Event(&pod, "Normal", "Freezing", "Agent reasoning gate triggered. Checkpointing state...")
		return r.handleFreeze(ctx, &pod)
	}

	// 2. Handle Auto-Reject Timeout
	if pod.Annotations["semamesh.io/status"] == "FROZEN" {
		freezeTime, err := time.Parse(time.RFC3339, pod.Annotations["semamesh.io/frozen-at"])
		if err == nil && time.Since(freezeTime) > 30*time.Minute {
			r.Recorder.Event(&pod, "Warning", "AutoReject", "Human approval timeout reached. Terminating pod for safety.")
			return ctrl.Result{}, r.Delete(ctx, &pod)
		}
		// Re-check every minute
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

func (r *SemaReconciler) handleFreeze(ctx context.Context, pod *corev1.Pod) (ctrl.Result, error) {
	// PRO TIP: In production, we'd use a dedicated ServiceAccount token
	// and a secure HTTP client that trusts the Kubelet CA.
	nodeIP := pod.Status.HostIP
	containerName := pod.Spec.Containers[0].Name
	checkpointURL := fmt.Sprintf("https://%s:10250/checkpoint/%s/%s/%s",
                        nodeIP, pod.Namespace, pod.Name, containerName)

	// We use a simplified POST here, but in Step 7 we will configure the
	// TLS transport to bypass the 'Insecure' error.
	req, _ := http.NewRequest("POST", checkpointURL, nil)
	// req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)

	if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
		r.Recorder.Event(pod, "Warning", "FreezeFailed", "Kubelet checkpoint call failed.")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	// Update metadata to show status is now FROZEN
	pod.Annotations["semamesh.io/status"] = "FROZEN"
	pod.Annotations["semamesh.io/frozen-at"] = time.Now().Format(time.RFC3339)
	delete(pod.Annotations, "semamesh.io/action")

	if err := r.Update(ctx, pod); err != nil {
		return ctrl.Result{}, err
	}

	r.Recorder.Event(pod, "Normal", "Frozen", "State successfully preserved. Awaiting 'kubectl sema approve'.")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SemaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}). // Tell the manager to watch Pods
		Complete(r)
}