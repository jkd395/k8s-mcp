package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type nodeData struct {
	Name              string            `json:"name,omitempty"`
	Status            string            `json:"status,omitempty"`
	KubernetesVersion string            `json:"kubernetesVersion,omitempty"`
	OS                string            `json:"os,omitempty"`
	KernelVersion     string            `json:"kernelVersion,omitempty"`
	Architecture      string            `json:"architecture,omitempty"`
	PodCIDR           string            `json:"podCIDR,omitempty"`
	CapacityCPU       string            `json:"capacityCPU,omitempty"`
	CapacityMemory    string            `json:"capacityMemory,omitempty"`
	CapacityPods      string            `json:"capacityPods,omitempty"`
	AllocatableCPU    string            `json:"allocatableCPU,omitempty"`
	AllocatableMemory string            `json:"allocatableMemory,omitempty"`
	AllocatablePods   string            `json:"allocatablePods,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Taints            []string          `json:"taints,omitempty"`
}

func ListNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	outFormat := request.GetString("output", "")
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing node: %v", err)), nil
	}
	if outFormat != "" {
		result, err := output.Format(outFormat, nodes.Items)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}
	var summary []nodeData
	for _, node := range nodes.Items {
		var nodeStatus string
		for _, v := range node.Status.Conditions {
			if v.Type == "Ready" {
				if v.Status == "True" {
					nodeStatus = "Ready"
				} else {
					nodeStatus = "NotReady"
				}
			}
		}
		var taints []string
		for _, t := range node.Spec.Taints {
			taints = append(taints, fmt.Sprintf("%s=%s:%s", t.Key, t.Value, t.Effect))
		}
		summary = append(summary, nodeData{
			Name:           node.Name,
			Status:         nodeStatus,
			CapacityCPU:    node.Status.Capacity.Cpu().String(),
			CapacityMemory: node.Status.Capacity.Memory().String(),
			CapacityPods:   node.Status.Capacity.Pods().String(),
			AllocatableCPU:    node.Status.Allocatable.Cpu().String(),
			AllocatableMemory: node.Status.Allocatable.Memory().String(),
			AllocatablePods:   node.Status.Allocatable.Pods().String(),
			Labels:         node.Labels,
			Taints:         taints,
		})
	}
	mcpOutput, err := json.MarshalIndent(summary, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func GetNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		msg := fmt.Sprintf("Provide name for node")
		return mcp.NewToolResultText(msg), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	outFormat := request.GetString("output", "")
	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting node: %v", err)), nil
	}
	if outFormat != "" {
		result, err := output.Format(outFormat, node)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}
	var nodeStatus string
	for _, v := range node.Status.Conditions {
		if v.Type == "Ready" {
			if v.Status == "True" {
				nodeStatus = "Ready"
			} else {
				nodeStatus = "NotReady"
			}
		}
	}
	var taints []string
	for _, t := range node.Spec.Taints {
		taints = append(taints, fmt.Sprintf("%s=%s:%s", t.Key, t.Value, t.Effect))
	}
	summary := nodeData{
		Name:              node.Name,
		Status:            nodeStatus,
		KubernetesVersion: node.Status.NodeInfo.KubeletVersion,
		OS:                node.Status.NodeInfo.OSImage,
		KernelVersion:     node.Status.NodeInfo.KernelVersion,
		Architecture:      node.Status.NodeInfo.Architecture,
		PodCIDR:           node.Spec.PodCIDR,
		CapacityCPU:       node.Status.Capacity.Cpu().String(),
		CapacityMemory:    node.Status.Capacity.Memory().String(),
		CapacityPods:      node.Status.Capacity.Pods().String(),
		AllocatableCPU:    node.Status.Allocatable.Cpu().String(),
		AllocatableMemory: node.Status.Allocatable.Memory().String(),
		AllocatablePods:   node.Status.Allocatable.Pods().String(),
		Labels:            node.Labels,
		Taints:            taints,
	}
	mcpOutput, err := json.MarshalIndent(summary, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for node")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.CoreV1().Nodes().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting node: %v", err)), nil
	}
	output := fmt.Sprintf("Node %s is deleted", name)
	return mcp.NewToolResultText(string(output)), nil
}

func UpdateNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for node")
		return mcp.NewToolResultText(string(output)), nil
	}
	labels, err := request.RequireString("label")
	if err != nil {
		output := fmt.Sprintf("Provide label for node")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting node: %v", err)), nil
	}
	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}
	label := strings.Split(labels, ",")
	for _, lab := range label {
		kv := strings.SplitN(lab, "=", 2)
		if len(kv) == 2 {
			node.Labels[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	updateNode, err := clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in updating node %s with label %s: %v", name, labels, err)), nil
	}
	output := fmt.Sprintf("Successfully node %s updated with label %s", updateNode.Name, labels)
	return mcp.NewToolResultText(string(output)), nil
}