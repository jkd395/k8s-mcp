package event

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

type eventData struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Type      string `json:"type,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Message   string `json:"message,omitempty"`
	Source    string `json:"source,omitempty"`
	FirstTime string `json:"firstTime,omitempty"`
	LastTime  string `json:"lastTime,omitempty"`
	Count     int32  `json:"count,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Involved  string `json:"involved,omitempty"`
}

func ListEvent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []eventData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var allItems []v1.Event
		for _, namespace := range namespaces.Items {
			events, err := clientset.CoreV1().Events(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing events in %s: %v", namespace.Name, err)), nil
			}
			allItems = append(allItems, events.Items...)
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, allItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, e := range allItems {
			output = append(output, toEventData(e))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		events, err := clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing events in %s: %v", ns, err)), nil
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, events.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, e := range events.Items {
			output = append(output, toEventData(e))
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetEvent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for event")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for event")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	e, err := clientset.CoreV1().Events(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting event %s/%s: %v", ns, name, err)), nil
	}

	if outputFmt != "" {
		result, err := outpkg.Format(outputFmt, e)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	output := toEventData(*e)
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func toEventData(e v1.Event) eventData {
	source := e.Source.Component
	if e.Source.Host != "" {
		source = source + "/" + e.Source.Host
	}
	return eventData{
		Name:      e.Name,
		Namespace: e.Namespace,
		Type:      e.Type,
		Reason:    e.Reason,
		Message:   e.Message,
		Source:    source,
		FirstTime: e.FirstTimestamp.Format("2006-01-02 15:04:05"),
		LastTime:  e.LastTimestamp.Format("2006-01-02 15:04:05"),
		Count:     e.Count,
		Kind:      e.InvolvedObject.Kind,
		Involved:  e.InvolvedObject.Name,
	}
}
