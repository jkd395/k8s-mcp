package limitrange

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type limitData struct {
	Name       string              `json:"name,omitempty"`
	Namespace  string              `json:"namespace,omitempty"`
	Limits     []limitItemData     `json:"limits,omitempty"`
}

type limitItemData struct {
	Type           string `json:"type,omitempty"`
	MaxCPU         string `json:"maxCPU,omitempty"`
	MaxMemory      string `json:"maxMemory,omitempty"`
	MinCPU         string `json:"minCPU,omitempty"`
	MinMemory      string `json:"minMemory,omitempty"`
	DefaultCPU     string `json:"defaultCPU,omitempty"`
	DefaultMemory  string `json:"defaultMemory,omitempty"`
	DefaultRequestCPU    string `json:"defaultRequestCPU,omitempty"`
	DefaultRequestMemory string `json:"defaultRequestMemory,omitempty"`
	MaxLimitRequestRatioCPU    string `json:"maxLimitRequestRatioCPU,omitempty"`
	MaxLimitRequestRatioMemory string `json:"maxLimitRequestRatioMemory,omitempty"`
}

func ListLimitRange(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	var output []limitData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			lrs, err := clientset.CoreV1().LimitRanges(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing limitrange in %s: %v", namespace.Name, err)), nil
			}
			for _, lr := range lrs.Items {
				output = append(output, toLimitData(lr))
			}
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		lrs, err := clientset.CoreV1().LimitRanges(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing limitrange in %s: %v", ns, err)), nil
		}
		for _, lr := range lrs.Items {
			output = append(output, toLimitData(lr))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetLimitRange(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for limitrange")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for limitrange")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	lr, err := clientset.CoreV1().LimitRanges(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting limitrange %s/%s: %v", ns, name, err)), nil
	}
	output := toLimitData(*lr)
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func toLimitData(lr v1.LimitRange) limitData {
	var items []limitItemData
	for _, l := range lr.Spec.Limits {
		item := limitItemData{Type: string(l.Type)}
		if l.Max != nil {
			if cpu, ok := l.Max[v1.ResourceCPU]; ok {
				item.MaxCPU = cpu.String()
			}
			if mem, ok := l.Max[v1.ResourceMemory]; ok {
				item.MaxMemory = mem.String()
			}
		}
		if l.Min != nil {
			if cpu, ok := l.Min[v1.ResourceCPU]; ok {
				item.MinCPU = cpu.String()
			}
			if mem, ok := l.Min[v1.ResourceMemory]; ok {
				item.MinMemory = mem.String()
			}
		}
		if l.Default != nil {
			if cpu, ok := l.Default[v1.ResourceCPU]; ok {
				item.DefaultCPU = cpu.String()
			}
			if mem, ok := l.Default[v1.ResourceMemory]; ok {
				item.DefaultMemory = mem.String()
			}
		}
		if l.DefaultRequest != nil {
			if cpu, ok := l.DefaultRequest[v1.ResourceCPU]; ok {
				item.DefaultRequestCPU = cpu.String()
			}
			if mem, ok := l.DefaultRequest[v1.ResourceMemory]; ok {
				item.DefaultRequestMemory = mem.String()
			}
		}
		if l.MaxLimitRequestRatio != nil {
			if cpu, ok := l.MaxLimitRequestRatio[v1.ResourceCPU]; ok {
				item.MaxLimitRequestRatioCPU = cpu.String()
			}
			if mem, ok := l.MaxLimitRequestRatio[v1.ResourceMemory]; ok {
				item.MaxLimitRequestRatioMemory = mem.String()
			}
		}
		items = append(items, item)
	}
	return limitData{
		Name:      lr.Name,
		Namespace: lr.Namespace,
		Limits:    items,
	}
}
