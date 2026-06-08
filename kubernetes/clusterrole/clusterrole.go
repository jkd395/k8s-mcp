package clusterrole

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type crData struct {
	Name  string  `json:"name,omitempty"`
	Rules []rules `json:"rules,omitempty"`
}

type rules struct {
	ApiGroups []string `json:"apiGroups,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
}

func ListCR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	crs, err := clientset.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing clusterrole: %v", err)), nil
	}
	outFmt := request.GetString("output", "")
	if outFmt != "" {
		result, err := output.Format(outFmt, crs.Items)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}
	var output []crData
	for _, cr := range crs.Items {
		output = append(output, crData{
			Name: cr.Name,
		})
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func GetCR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for clusterrole")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	cr, err := clientset.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting clusterrole named %s: %v", name, err)), nil
	}
	outFmt := request.GetString("output", "")
	if outFmt != "" {
		result, err := output.Format(outFmt, cr)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	var crRules []rules

	for _, rule := range cr.Rules {
		crRules = append(crRules, rules{
			ApiGroups: rule.APIGroups,
			Resources: rule.Resources,
			Verbs:     rule.Verbs,
		})
	}

	output := crData{
		Name:  cr.Name,
		Rules: crRules,
	}

	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteCR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for clusterrole")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting clusterrole named %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted clusterrole named %s", name)), nil
}