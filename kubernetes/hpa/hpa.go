package hpa

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"

	"github.com/mark3labs/mcp-go/mcp"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type hpaData struct {
	Name                  string `json:"name,omitempty"`
	Namespace             string `json:"namespace,omitempty"`
	MinReplicas           int32  `json:"minReplicas,omitempty"`
	MaxReplicas           int32  `json:"maxReplicas,omitempty"`
	CurrentReplicas       int32  `json:"currentReplicas,omitempty"`
	DesiredReplicas       int32  `json:"desiredReplicas,omitempty"`
	TargetCPUUtilization  *int32 `json:"targetCPUUtilization,omitempty"`
	CurrentCPUUtilization *int32 `json:"currentCPUUtilization,omitempty"`
	ScaleTargetRefKind    string `json:"scaleTargetRefKind,omitempty"`
	ScaleTargetRefName    string `json:"scaleTargetRefName,omitempty"`
}

func ListHPA(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	var output []hpaData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing hpa in %s: %v", namespace.Name, err)), nil
			}
			for _, h := range hpas.Items {
				minReplicas := int32(1)
				if h.Spec.MinReplicas != nil {
					minReplicas = *h.Spec.MinReplicas
				}
				output = append(output, hpaData{
					Name:               h.Name,
					Namespace:          h.Namespace,
					MinReplicas:        minReplicas,
					MaxReplicas:        h.Spec.MaxReplicas,
					CurrentReplicas:    h.Status.CurrentReplicas,
					DesiredReplicas:    h.Status.DesiredReplicas,
					ScaleTargetRefKind: h.Spec.ScaleTargetRef.Kind,
					ScaleTargetRefName: h.Spec.ScaleTargetRef.Name,
				})
			}
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing hpa in %s namespace: %v", ns, err)), nil
		}
		for _, h := range hpas.Items {
			minReplicas := int32(1)
			if h.Spec.MinReplicas != nil {
				minReplicas = *h.Spec.MinReplicas
			}
			output = append(output, hpaData{
				Name:               h.Name,
				Namespace:          h.Namespace,
				MinReplicas:        minReplicas,
				MaxReplicas:        h.Spec.MaxReplicas,
				CurrentReplicas:    h.Status.CurrentReplicas,
				DesiredReplicas:    h.Status.DesiredReplicas,
				ScaleTargetRefKind: h.Spec.ScaleTargetRef.Kind,
				ScaleTargetRefName: h.Spec.ScaleTargetRef.Name,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetHPA(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	h, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting hpa %s/%s: %v", ns, name, err)), nil
	}

	var targetCPU *int32
	if h.Spec.Metrics != nil {
		for _, metric := range h.Spec.Metrics {
			if metric.Resource != nil && metric.Resource.Name == "cpu" {
				if metric.Resource.Target.AverageUtilization != nil {
					targetCPU = metric.Resource.Target.AverageUtilization
				}
			}
		}
	}

	var currentCPU *int32
	if h.Status.CurrentMetrics != nil {
		for _, metric := range h.Status.CurrentMetrics {
			if metric.Resource != nil && metric.Resource.Name == "cpu" {
				if metric.Resource.Current.AverageUtilization != nil {
					currentCPU = metric.Resource.Current.AverageUtilization
				}
			}
		}
	}

	minReplicas := int32(1)
	if h.Spec.MinReplicas != nil {
		minReplicas = *h.Spec.MinReplicas
	}
	output := hpaData{
		Name:                  h.Name,
		Namespace:             h.Namespace,
		MinReplicas:           minReplicas,
		MaxReplicas:           h.Spec.MaxReplicas,
		CurrentReplicas:       h.Status.CurrentReplicas,
		DesiredReplicas:       h.Status.DesiredReplicas,
		TargetCPUUtilization:  targetCPU,
		CurrentCPUUtilization: currentCPU,
		ScaleTargetRefKind:    h.Spec.ScaleTargetRef.Kind,
		ScaleTargetRefName:    h.Spec.ScaleTargetRef.Name,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteHPA(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	err = clientset.AutoscalingV2().HorizontalPodAutoscalers(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting hpa %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("HPA %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateHPA(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	targetKind := request.GetString("targetKind", "Deployment")
	targetName, err := request.RequireString("targetName")
	if err != nil {
		output := fmt.Sprintf("Provide targetName for hpa (the workload name to scale)")
		return mcp.NewToolResultText(string(output)), nil
	}
	minReplicas := request.GetInt("minReplicas", 1)
	maxReplicas, err := request.RequireInt("maxReplicas")
	if err != nil {
		output := fmt.Sprintf("Provide maxReplicas for hpa")
		return mcp.NewToolResultText(string(output)), nil
	}
	cpuTarget := request.GetInt("cpuTarget", 50)

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}

	minReplicasInt32 := int32(minReplicas)
	maxReplicasInt32 := int32(maxReplicas)
	cpuTargetInt32 := int32(cpuTarget)

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind:       targetKind,
				Name:       targetName,
				APIVersion: "apps/v1",
			},
			MinReplicas: &minReplicasInt32,
			MaxReplicas: maxReplicasInt32,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2.MetricTarget{
							Type:               autoscalingv2.UtilizationMetricType,
							AverageUtilization: &cpuTargetInt32,
						},
					},
				},
			},
		},
	}

	createHPA, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(ns).Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating hpa %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully hpa %s/%s is created", createHPA.Namespace, createHPA.Name)
	return mcp.NewToolResultText(string(output)), nil
}
