package crd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type crdData struct {
	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`
}

func ListCRD(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, _, _, apiClient, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	crds, err := apiClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in listing crds: %v", err)), nil
	}
	var output []string
	for _, crd := range crds.Items {
		output = append(output, crd.Name)
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func GetCRD(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for crd")
		return mcp.NewToolResultText(string(output)), nil
	}
	_, _, _, apiClient, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	crds, err := apiClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting crd: %v", err)), nil
	}
	var output crdData
	output = crdData{
		Name: crds.Name,
		Kind: crds.Spec.Names.Kind,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteCRD(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for crd")
		return mcp.NewToolResultText(string(output)), nil
	}
	_, _, _, apiClient, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = apiClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting crds: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted crd %s", name)), nil
}

func CreateCRDWithJson(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jsondata, err := request.RequireString("jsondata")
	if err != nil {
		output := fmt.Sprintf("Provide jsonData for crd")
		return mcp.NewToolResultText(string(output)), nil
	}
	_, _, _, apiClient, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	var crd apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal([]byte(jsondata), &crd); err != nil {
		return nil, err
	}

	_, err = apiClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, &crd, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating crd: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully created crd")), nil
}
