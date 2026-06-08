package resourcequota

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	outpkg "k8s-mcp/kubernetes/output"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type rqData struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Hard      map[string]string `json:"hard,omitempty"`
	Used      map[string]string `json:"used,omitempty"`
}

func ListResourceQuota(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []rqData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var allItems []v1.ResourceQuota
		for _, namespace := range namespaces.Items {
			rqs, err := clientset.CoreV1().ResourceQuotas(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing resourcequota in %s: %v", namespace.Name, err)), nil
			}
			allItems = append(allItems, rqs.Items...)
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, allItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, rq := range allItems {
			output = append(output, toRQData(rq))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		rqs, err := clientset.CoreV1().ResourceQuotas(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing resourcequota in %s: %v", ns, err)), nil
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, rqs.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, rq := range rqs.Items {
			output = append(output, toRQData(rq))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetResourceQuota(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for resourcequota")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for resourcequota")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	rq, err := clientset.CoreV1().ResourceQuotas(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting resourcequota %s/%s: %v", ns, name, err)), nil
	}

	if outputFmt != "" {
		result, err := outpkg.Format(outputFmt, rq)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	output := toRQData(*rq)
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func toRQData(rq v1.ResourceQuota) rqData {
	hard := make(map[string]string)
	for k, v := range rq.Status.Hard {
		hard[string(k)] = v.String()
	}
	used := make(map[string]string)
	for k, v := range rq.Status.Used {
		used[string(k)] = v.String()
	}
	return rqData{
		Name:      rq.Name,
		Namespace: rq.Namespace,
		Hard:      hard,
		Used:      used,
	}
}
