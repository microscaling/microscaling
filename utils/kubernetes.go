// Package utils contains common shared code.
package utils

import (
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

// NewKubeClientset creates a Kubernetes clientset for accessing the Kubernetes APIs.
// Uses the provided kube config file or when running as a pod uses the built in config.
func NewKubeClientset(kubeConfig string, namespace string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Errorf("Error creating Kubernetes config: %v", err)
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Errorf("Error getting Kubernetes config: %v", err)
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Error creating Kubernetes clientset: %v", err)
		return nil, err
	}

	return clientset, err
}
