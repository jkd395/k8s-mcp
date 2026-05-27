package top

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type topPodData struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	CPU       string `json:"cpu,omitempty"`
	Memory    string `json:"memory,omitempty"`
}

type topNodeData struct {
	Name   string `json:"name,omitempty"`
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

func TopPod(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	clientset, _, _, _, mc, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	var pms []metricsv1beta1.PodMetrics
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			pmList, err := mc.MetricsV1beta1().PodMetricses(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				continue
			}
			pms = append(pms, pmList.Items...)
		}
	} else {
		pmList, err := mc.MetricsV1beta1().PodMetricses(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing pod metrics in %s: %v", ns, err)), nil
		}
		pms = pmList.Items
	}
	var output []topPodData
	for _, pm := range pms {
		cpu, mem := sumContainerMetrics(pm.Containers)
		output = append(output, topPodData{
			Name:      pm.Name,
			Namespace: pm.Namespace,
			CPU:       cpu,
			Memory:    mem,
		})
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func TopNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, _, _, _, mc, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	nodesMetrics, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing node metrics: %v", err)), nil
	}
	var output []topNodeData
	for _, nm := range nodesMetrics.Items {
		output = append(output, topNodeData{
			Name:   nm.Name,
			CPU:    nm.Usage.Cpu().String(),
			Memory: nm.Usage.Memory().String(),
		})
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func sumContainerMetrics(containers []metricsv1beta1.ContainerMetrics) (string, string) {
	var totalCPU, totalMemory int64
	for _, c := range containers {
		if cpu, ok := c.Usage["cpu"]; ok {
			totalCPU += cpu.MilliValue()
		}
		if mem, ok := c.Usage["memory"]; ok {
			totalMemory += mem.Value()
		}
	}
	cpuStr := fmt.Sprintf("%dm", totalCPU)
	memStr := fmt.Sprintf("%dMi", totalMemory/(1024*1024))
	if totalMemory < 1024*1024 {
		memStr = fmt.Sprintf("%dKi", totalMemory/1024)
	}
	return cpuStr, memStr
}
