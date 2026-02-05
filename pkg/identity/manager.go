package identity

import (
	"fmt"
	"log"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	// metav1 removed
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	// context removed
)

// PodMetadata holds the identity info we care about
type PodMetadata struct {
	Namespace      string
	PodName        string
	ServiceAccount string
}

// Manager is a thread-safe store for IP -> Identity lookups
type Manager struct {
	ipMap   map[string]PodMetadata
	mutex   sync.RWMutex
	devMode bool
}

// NewManager creates the store
func NewManager(devMode bool) *Manager {
	return &Manager{
		ipMap:   make(map[string]PodMetadata),
		devMode: devMode,
	}
}

// StartWatcher connects to K8s and listens for Pod IP changes
func (m *Manager) StartWatcher(kubeconfigPath string) error {
	if m.devMode {
		log.Println("⚠️ Identity Manager: Running in Dev Mode (No K8s connection)")
		return nil
	}

	// 1. Build K8s Config (In-Cluster or Kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	// 2. Create the Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %v", err)
	}

	// 3. Create an Informer Factory (Resync every 10 minutes)
	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	podInformer := factory.Core().V1().Pods().Informer()

	// 4. Register Event Handlers (Add, Update, Delete)
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			m.handlePodUpdate(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			m.handlePodUpdate(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			m.handlePodDelete(obj)
		},
	})

	// 5. Start the Watcher in the background
	stopper := make(chan struct{})
	log.Println("⚡ Connected to Kubernetes API. Watching Pods...")
	go podInformer.Run(stopper)

	// Wait for cache sync so we don't miss existing pods
	if !cache.WaitForCacheSync(stopper, podInformer.HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	return nil
}

// handlePodUpdate extracts IP and Identity from a Pod object
func (m *Manager) handlePodUpdate(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	// Only care if the Pod has an IP
	if pod.Status.PodIP == "" {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ipMap[pod.Status.PodIP] = PodMetadata{
		Namespace:      pod.Namespace,
		PodName:        pod.Name,
		ServiceAccount: pod.Spec.ServiceAccountName,
	}
}

// handlePodDelete removes the IP from the map
func (m *Manager) handlePodDelete(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.ipMap, pod.Status.PodIP)
}

// GetIdentity looks up the metadata for a given IP
func (m *Manager) GetIdentity(ip string) (PodMetadata, bool) {
	// 1. Dev Mode Bypass (Localhost always returns fake identity)
	if m.devMode && (ip == "127.0.0.1" || ip == "::1") {
		return PodMetadata{
			Namespace:      "dev-workspace",
			PodName:        "curl-terminal-agent",
			ServiceAccount: "admin-user",
		}, true
	}

	// 2. Real Lookup
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	meta, exists := m.ipMap[ip]
	return meta, exists
}