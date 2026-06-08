package role

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	k8soutput "k8s-mcp/kubernetes/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type roleData struct {
	Name      string  `json:"name,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
	Rules     []rules ` json:"rules,omitempty"`
}

type rules struct {
	ApiGroups []string `json:"apiGroups,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
}

func ListRole(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []roleData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var rawItems []any
		for _, namespace := range namespaces.Items {
			roles, err := clientset.RbacV1().Roles(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing role in namespace %s: %v", namespace.Name, err)), nil
			}
			for _, role := range roles.Items {
				rawItems = append(rawItems, role)
				output = append(output, roleData{
					Name:      role.Name,
					Namespace: role.Namespace,
				})
			}
		}
		if outputParam != "" {
			raw, err := k8soutput.Format(outputParam, rawItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(raw), nil
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		roles, err := clientset.RbacV1().Roles(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing role in %s: %v", ns, err)), nil
		}
		if outputParam != "" {
			raw, err := k8soutput.Format(outputParam, roles.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(raw), nil
		}
		for _, role := range roles.Items {
			output = append(output, roleData{
				Name:      role.Name,
				Namespace: role.Namespace,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetRole(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for role")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for role")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	role, err := clientset.RbacV1().Roles(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting role in %s/%s: %v", ns, name, err)), nil
	}

	if outputParam != "" {
		raw, err := k8soutput.Format(outputParam, role)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(raw), nil
	}

	var roleRules []rules

	for _, rule := range role.Rules {
		roleRules = append(roleRules, rules{
			ApiGroups: rule.APIGroups,
			Resources: rule.Resources,
			Verbs:     rule.Verbs,
		})
	}

	output := roleData{
		Name:      role.Name,
		Namespace: role.Namespace,
		Rules:     roleRules,
	}

	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteRole(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for role")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for role")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.RbacV1().Roles(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting role %s/%s: %v", ns, name, err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted role %s/%s", ns, name)), nil
}