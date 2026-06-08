package ingress

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ingressData struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Hosts     []string          `json:"hosts,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func ListIngress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var ingressList []ingressData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing ingress in %s: %v", namespace.Name, err)), nil
			}
			if outFmt != "" {
				result, err := output.Format(outFmt, ingresses.Items)
				if err != nil {
					return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
				}
				return mcp.NewToolResultText(result), nil
			}
			for _, ingress := range ingresses.Items {
				var hosts []string
				for _, rule := range ingress.Spec.Rules {
					hosts = append(hosts, rule.Host)
				}
				ingressList = append(ingressList, ingressData{
					Name:      ingress.Name,
					Namespace: ingress.Namespace,
					Hosts:     hosts,
					Labels:    ingress.Labels,
				})
			}
		}
		mcpOutput, err := json.MarshalIndent(ingressList, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		ingresses, err := clientset.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing ingress in %s namespace: %v", ns, err)), nil
		}
		if outFmt != "" {
			result, err := output.Format(outFmt, ingresses.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, ingress := range ingresses.Items {
			var hosts []string
			for _, rule := range ingress.Spec.Rules {
				hosts = append(hosts, rule.Host)
			}
			ingressList = append(ingressList, ingressData{
				Name:      ingress.Name,
				Namespace: ingress.Namespace,
				Hosts:     hosts,
				Labels:    ingress.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(ingressList, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetIngress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	ingress, err := clientset.NetworkingV1().Ingresses(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting ingress %s/%s: %v", ns, name, err)), nil
	}
	outFmt := request.GetString("output", "")
	if outFmt != "" {
		result, err := output.Format(outFmt, ingress)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}
	var hosts []string
	for _, rule := range ingress.Spec.Rules {
		hosts = append(hosts, rule.Host)
	}
	res := ingressData{
		Name:      ingress.Name,
		Namespace: ingress.Namespace,
		Hosts:     hosts,
		Labels:    ingress.Labels,
	}
	mcpOutput, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteIngress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting ingress %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Ingress %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateIngress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for ingress")
		return mcp.NewToolResultText(string(output)), nil
	}
	host := request.GetString("host", "")
	serviceName := request.GetString("serviceName", "")
	servicePort := request.GetInt("servicePort", 80)
	labels := request.GetString("label", "")
	ingressClassName := request.GetString("ingressClassName", "nginx")

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	lab := make(map[string]string)
	if labels != "" {
		ingLabel := strings.Split(labels, ",")
		for _, label := range ingLabel {
			kv := strings.SplitN(label, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				lab[key] = value
			}
		}
	}

	if len(lab) == 0 {
		lab["app"] = name
	}

	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    lab,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Number: int32(servicePort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	createIngress, err := clientset.NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating ingress %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully ingress %s/%s is created", createIngress.Namespace, createIngress.Name)
	return mcp.NewToolResultText(string(output)), nil
}
