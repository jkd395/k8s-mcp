package componentstatus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type csData struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

func ListComponentStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	css, err := clientset.CoreV1().ComponentStatuses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing componentstatus: %v", err)), nil
	}
	var output []csData
	for _, cs := range css.Items {
		status := "Unknown"
		ctype := ""
		for _, c := range cs.Conditions {
			status = string(c.Status)
			ctype = string(c.Type)
		}
		output = append(output, csData{
			Name:   cs.Name,
			Status: status,
			Type:   ctype,
		})
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func GetComponentStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for componentstatus")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	cs, err := clientset.CoreV1().ComponentStatuses().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting componentstatus %s: %v", name, err)), nil
	}
	status := "Unknown"
	ctype := ""
	for _, c := range cs.Conditions {
		status = string(c.Status)
		ctype = string(c.Type)
	}
	out := csData{
		Name:   cs.Name,
		Status: status,
		Type:   ctype,
	}
	mcpOutput, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}
