package endpoint

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

type endpointData struct {
	Name      string        `json:"name,omitempty"`
	Namespace string        `json:"namespace,omitempty"`
	Addresses []addressData `json:"addresses,omitempty"`
}

type addressData struct {
	IP       string     `json:"ip,omitempty"`
	NodeName string     `json:"nodeName,omitempty"`
	Ports    []portData `json:"ports,omitempty"`
}

type portData struct {
	Name string `json:"name,omitempty"`
	Port int32  `json:"port,omitempty"`
}

func ListEndpoint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []endpointData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var allItems []v1.Endpoints
		for _, namespace := range namespaces.Items {
			eps, err := clientset.CoreV1().Endpoints(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing endpoints in %s: %v", namespace.Name, err)), nil
			}
			allItems = append(allItems, eps.Items...)
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, allItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, ep := range allItems {
			output = append(output, toEndpointData(ep))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		eps, err := clientset.CoreV1().Endpoints(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing endpoints in %s: %v", ns, err)), nil
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, eps.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, ep := range eps.Items {
			output = append(output, toEndpointData(ep))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetEndpoint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for endpoint")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for endpoint")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	ep, err := clientset.CoreV1().Endpoints(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting endpoint %s/%s: %v", ns, name, err)), nil
	}

	if outputFmt != "" {
		result, err := outpkg.Format(outputFmt, ep)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	output := toEndpointData(*ep)
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func toEndpointData(ep v1.Endpoints) endpointData {
	var addrs []addressData
	for _, subset := range ep.Subsets {
		for _, addr := range subset.Addresses {
			var pts []portData
			for _, p := range subset.Ports {
				pts = append(pts, portData{Name: p.Name, Port: p.Port})
			}
			nodeName := ""
			if addr.NodeName != nil {
				nodeName = *addr.NodeName
			}
			addrs = append(addrs, addressData{
				IP:       addr.IP,
				NodeName: nodeName,
				Ports:    pts,
			})
		}
	}
	return endpointData{
		Name:      ep.Name,
		Namespace: ep.Namespace,
		Addresses: addrs,
	}
}
