package clusterhealth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type healthReport struct {
	Nodes          nodeSummary     `json:"nodes,omitempty"`
	ControlPlanePods []cpPodStatus `json:"controlPlanePods,omitempty"`
}

type nodeSummary struct {
	Total   int `json:"total"`
	Ready   int `json:"ready"`
	NotReady int `json:"notReady"`
}

type cpPodStatus struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Component string `json:"component,omitempty"`
	Status    string `json:"status,omitempty"`
	Restarts  int32  `json:"restarts,omitempty"`
	Node      string `json:"node,omitempty"`
}

type nodeHealthData struct {
	Name      string            `json:"name,omitempty"`
	Status    string            `json:"status,omitempty"`
	Kubelet   string            `json:"kubelet,omitempty"`
	CPU       string            `json:"cpu,omitempty"`
	Memory    string            `json:"memory,omitempty"`
	Pods      string            `json:"pods,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func GetClusterHealth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}

	report := healthReport{}

	// 1. Node health summary
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing nodes: %v", err)), nil
	}
	report.Nodes.Total = len(nodes.Items)
	for _, n := range nodes.Items {
		ready := false
		for _, c := range n.Status.Conditions {
			if c.Type == v1.NodeReady && c.Status == v1.ConditionTrue {
				ready = true
				break
			}
		}
		if ready {
			report.Nodes.Ready++
		} else {
			report.Nodes.NotReady++
		}
	}

	// 2. Control plane component pods in kube-system
	componentLabels := []struct{
		label     string
		component string
	}{
		{"component=etcd", "etcd"},
		{"component=kube-apiserver", "kube-apiserver"},
		{"component=kube-scheduler", "kube-scheduler"},
		{"component=kube-controller-manager", "kube-controller-manager"},
		{"k8s-app=kube-dns", "coredns"},
		{"k8s-app=coredns", "coredns"},
	}

	seen := make(map[string]bool)
	for _, cl := range componentLabels {
		pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{
			LabelSelector: cl.label,
		})
		if err != nil {
			continue
		}
		for _, p := range pods.Items {
			if seen[p.Name] {
				continue
			}
			seen[p.Name] = true
			status := podStatusSummary(p)
			var restarts int32
			for _, cs := range p.Status.ContainerStatuses {
				restarts += cs.RestartCount
			}
			report.ControlPlanePods = append(report.ControlPlanePods, cpPodStatus{
				Name:      p.Name,
				Namespace: p.Namespace,
				Component: cl.component,
				Status:    status,
				Restarts:  restarts,
				Node:      p.Spec.NodeName,
			})
		}
	}

	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(output)), nil
}

func ListNodeHealth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing nodes: %v", err)), nil
	}
	var output []nodeHealthData
	for _, n := range nodes.Items {
		ready := false
		for _, c := range n.Status.Conditions {
			if c.Type == v1.NodeReady && c.Status == v1.ConditionTrue {
				ready = true
				break
			}
		}
		status := "NotReady"
		if ready {
			status = "Ready"
		}
		output = append(output, nodeHealthData{
			Name:      n.Name,
			Status:    status,
			Kubelet:   n.Status.NodeInfo.KubeletVersion,
			CPU:       n.Status.Allocatable.Cpu().String(),
			Memory:    n.Status.Allocatable.Memory().String(),
			Pods:      n.Status.Allocatable.Pods().String(),
			Labels:    n.Labels,
		})
	}
	mcpOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func podStatusSummary(p v1.Pod) string {
	for _, cs := range p.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			return fmt.Sprintf("Waiting:%s", cs.State.Waiting.Reason)
		}
		if cs.State.Terminated != nil && cs.State.Terminated.Reason != "" {
			return fmt.Sprintf("Terminated:%s", cs.State.Terminated.Reason)
		}
	}
	if p.Status.Phase == v1.PodRunning {
		return "Running"
	}
	return string(p.Status.Phase)
}

// Keep componentstatus fallback for older clusters
