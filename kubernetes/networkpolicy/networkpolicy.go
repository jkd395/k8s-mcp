package networkpolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	outpkg "k8s-mcp/kubernetes/output"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type networkPolicyData struct {
	Name        string                    `json:"name,omitempty"`
	Namespace   string                    `json:"namespace,omitempty"`
	PodSelector map[string]string         `json:"podSelector,omitempty"`
	PolicyTypes []networkingv1.PolicyType `json:"policyTypes,omitempty"`
	Labels      map[string]string         `json:"labels,omitempty"`
}

func ListNetworkPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []networkPolicyData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var allItems []networkingv1.NetworkPolicy
		for _, namespace := range namespaces.Items {
			policies, err := clientset.NetworkingV1().NetworkPolicies(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing networkpolicy in %s: %v", namespace.Name, err)), nil
			}
			allItems = append(allItems, policies.Items...)
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, allItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, np := range allItems {
			output = append(output, networkPolicyData{
				Name:        np.Name,
				Namespace:   np.Namespace,
				PodSelector: np.Spec.PodSelector.MatchLabels,
				PolicyTypes: np.Spec.PolicyTypes,
				Labels:      np.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		policies, err := clientset.NetworkingV1().NetworkPolicies(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing networkpolicy in %s namespace: %v", ns, err)), nil
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, policies.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, np := range policies.Items {
			output = append(output, networkPolicyData{
				Name:        np.Name,
				Namespace:   np.Namespace,
				PodSelector: np.Spec.PodSelector.MatchLabels,
				PolicyTypes: np.Spec.PolicyTypes,
				Labels:      np.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetNetworkPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	np, err := clientset.NetworkingV1().NetworkPolicies(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting networkpolicy %s/%s: %v", ns, name, err)), nil
	}

	if outputFmt != "" {
		result, err := outpkg.Format(outputFmt, np)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	output := networkPolicyData{
		Name:        np.Name,
		Namespace:   np.Namespace,
		PodSelector: np.Spec.PodSelector.MatchLabels,
		PolicyTypes: np.Spec.PolicyTypes,
		Labels:      np.Labels,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteNetworkPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.NetworkingV1().NetworkPolicies(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting networkpolicy %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("NetworkPolicy %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateNetworkPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for networkpolicy")
		return mcp.NewToolResultText(string(output)), nil
	}
	podSelector := request.GetString("podSelector", "app=myapp")
	policyTypes := request.GetString("policyTypes", "Ingress")

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	// Parse pod selector labels
	podSel := make(map[string]string)
	ps := strings.Split(podSelector, ",")
	for _, p := range ps {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			podSel[key] = value
		}
	}

	// Parse policy types
	var types []networkingv1.PolicyType
	pt := strings.Split(policyTypes, ",")
	for _, t := range pt {
		t = strings.TrimSpace(t)
		if t == "Ingress" || t == "Egress" {
			types = append(types, networkingv1.PolicyType(t))
		}
	}

	if len(types) == 0 {
		types = []networkingv1.PolicyType{networkingv1.PolicyTypeIngress}
	}

	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: podSel,
			},
			PolicyTypes: types,
		},
	}

	createNP, err := clientset.NetworkingV1().NetworkPolicies(ns).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating networkpolicy %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully networkpolicy %s/%s is created", createNP.Namespace, createNP.Name)
	return mcp.NewToolResultText(string(output)), nil
}
