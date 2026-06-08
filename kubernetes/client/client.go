package client

import (
	"flag"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var kubeconfigPath string

func init() {
	flag.StringVar(&kubeconfigPath, "kubeconfigPath", "", "Path to kubeconfig file")
}

func GetRestConfig() (*rest.Config, error) {
	if kubeconfigPath == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func InitializeClients() (*kubernetes.Clientset, dynamic.Interface, discovery.DiscoveryInterface, *apiextensionsclient.Clientset, metricsv.Interface, error) {

	config, err := GetRestConfig()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	discoverClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	apiClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return clientset, dynamicClient, discoverClient, apiClient, metricsClient, nil
}
