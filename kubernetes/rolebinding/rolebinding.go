package rolebinding

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	k8soutput "k8s-mcp/kubernetes/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type rbData struct {
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	RoleRef   roleRef    `json:"roleRef,omitempty"`
	Subjects  []subjects `json:"subjects,omitempty"`
}

type roleRef struct {
	ApiGroup string `json:"apiGroup,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Name     string `json:"name,omitempty"`
}

type subjects struct {
	ApiGroup  string `json:"apiGroup,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func ListRB(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []rbData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var rawItems []any
		for _, namespace := range namespaces.Items {
			rbs, err := clientset.RbacV1().RoleBindings(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing rolebinding in namespace %s: %v", namespace.Name, err)), nil
			}
			for _, rb := range rbs.Items {
				rawItems = append(rawItems, rb)
				output = append(output, rbData{
					Name:      rb.Name,
					Namespace: rb.Namespace,
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
		rbs, err := clientset.RbacV1().RoleBindings(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing rolebinding in %s: %v", ns, err)), nil
		}
		if outputParam != "" {
			raw, err := k8soutput.Format(outputParam, rbs.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(raw), nil
		}
		for _, rb := range rbs.Items {
			output = append(output, rbData{
				Name:      rb.Name,
				Namespace: rb.Namespace,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetRB(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for rolebinding")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for rolebinding")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	rb, err := clientset.RbacV1().RoleBindings(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting rolebinding in %s/%s: %v", ns, name, err)), nil
	}

	if outputParam != "" {
		raw, err := k8soutput.Format(outputParam, rb)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(raw), nil
	}

	var saDetails []subjects

	for _, rolebind := range rb.Subjects {
		saDetails = append(saDetails, subjects{
			ApiGroup:  rolebind.APIGroup,
			Kind:      rolebind.Kind,
			Name:      rolebind.Name,
			Namespace: rolebind.Namespace,
		})
	}

	var rRef roleRef
	rRef = roleRef{
		ApiGroup: rb.RoleRef.APIGroup,
		Kind:     rb.RoleRef.Kind,
		Name:     rb.RoleRef.Name,
	}

	output := rbData{
		Name:      rb.Name,
		Namespace: rb.Namespace,
		RoleRef:   rRef,
		Subjects:  saDetails,
	}

	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteRB(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for rolebinding")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for rolebinding")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.RbacV1().RoleBindings(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting rolebinding %s/%s: %v", ns, name, err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted rolebinding %s/%s", ns, name)), nil
}