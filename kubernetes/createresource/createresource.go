package createresource

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/yaml"
)

func CreateResourceWithJson(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jsondata, err := request.RequireString("jsondata")
	if err != nil {
		output := fmt.Sprintf("Provide jsonData to create resource")
		return mcp.NewToolResultText(string(output)), nil
	}
	_, dynamicClient, discoveryClient, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var obj unstructured.Unstructured
	if err = yaml.Unmarshal([]byte(jsondata), &obj); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in unmarshal the json/yaml data: %v", err)), nil
	}

	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting API group resources: %v", err)), nil
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in REST mapping for %s: %v", gvk.Kind, err)), nil
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == "root" {
		dr = dynamicClient.Resource(mapping.Resource)
	} else {
		ns := obj.GetNamespace()
		if ns == "" {
			ns = "default"
		}
		dr = dynamicClient.Resource(mapping.Resource).Namespace(ns)
	}
	_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating resource with jsondata: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully created resource with jsondata")), nil
}
