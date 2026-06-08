package configmap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	k8soutput "k8s-mcp/kubernetes/output"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type cmData struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

func ListConfigmap(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []cmData
	var rawItems []v1.ConfigMap
	if ns == ""{
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			configmaps, err := clientset.CoreV1().ConfigMaps(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing configmaps in %s: %v", namespace.Name, err)), nil
			}
			rawItems = append(rawItems, configmaps.Items...)
			for _, configmap := range configmaps.Items {
				output = append(output, cmData{
					Name:      configmap.Name,
					Namespace: configmap.Namespace,
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
		configmaps, err := clientset.CoreV1().ConfigMaps(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing configmaps in %s namespace: %v", ns, err)), nil
		}
		if outputParam != "" {
			raw, err := k8soutput.Format(outputParam, configmaps.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(raw), nil
		}
		for _, configmap := range configmaps.Items {
			output = append(output, cmData{
				Name:      configmap.Name,
				Namespace: configmap.Namespace,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetConfigmap(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for configmap")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for configmap")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputParam := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	configmap, err := clientset.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting configmap in %s/%s: %v", ns, name, err)), nil
	}

	if outputParam != "" {
		raw, err := k8soutput.Format(outputParam, configmap)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(raw), nil
	}

	output := cmData{
		Name:      configmap.Name,
		Namespace: configmap.Namespace,
		Data:      configmap.Data,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteConfigmap(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for configmap")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for configmap")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.CoreV1().ConfigMaps(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting configmaps in %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Configmap %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateConfigmap(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for configmap creation")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for configmap creation")
		return mcp.NewToolResultText(string(output)), nil
	}
	data, err := request.RequireString("data")
	if err != nil {
		output := fmt.Sprintf("Provide datas for configmap creation like password=kubernetes123,username=kubernetes")
		return mcp.NewToolResultText(string(output)), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	configmapData := make(map[string]string)

	cmData := strings.Split(data, ",")
	for _, datas := range cmData {
		kv := strings.SplitN(datas, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			configmapData[key] = value
		}
	}

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: configmapData,
	}
	createConfigmap, err := clientset.CoreV1().ConfigMaps(ns).Create(ctx, configmap, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating configmap in %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully configmap %s/%s is created", createConfigmap.Namespace, createConfigmap.Name)
	return mcp.NewToolResultText(string(output)), nil
}