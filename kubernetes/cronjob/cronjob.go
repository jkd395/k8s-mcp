package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	outpkg "k8s-mcp/kubernetes/output"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type cronjobData struct {
	Name             string            `json:"name,omitempty"`
	Namespace        string            `json:"namespace,omitempty"`
	Schedule         string            `json:"schedule,omitempty"`
	Suspend          *bool             `json:"suspend,omitempty"`
	LastScheduleTime *string           `json:"lastScheduleTime,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
}

func ListCronJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var output []cronjobData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		var allItems []batchv1.CronJob
		for _, namespace := range namespaces.Items {
			cronjobs, err := clientset.BatchV1().CronJobs(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing cronjob in %s: %v", namespace.Name, err)), nil
			}
			allItems = append(allItems, cronjobs.Items...)
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, allItems)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, cj := range allItems {
			var lastSchedule *string
			if cj.Status.LastScheduleTime != nil {
				t := cj.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
				lastSchedule = &t
			}
			output = append(output, cronjobData{
				Name:             cj.Name,
				Namespace:        cj.Namespace,
				Schedule:         cj.Spec.Schedule,
				Suspend:          cj.Spec.Suspend,
				LastScheduleTime: lastSchedule,
				Labels:           cj.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		cronjobs, err := clientset.BatchV1().CronJobs(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing cronjob in %s namespace: %v", ns, err)), nil
		}
		if outputFmt != "" {
			result, err := outpkg.Format(outputFmt, cronjobs.Items)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
			}
			return mcp.NewToolResultText(result), nil
		}
		for _, cj := range cronjobs.Items {
			var lastSchedule *string
			if cj.Status.LastScheduleTime != nil {
				t := cj.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
				lastSchedule = &t
			}
			output = append(output, cronjobData{
				Name:             cj.Name,
				Namespace:        cj.Namespace,
				Schedule:         cj.Spec.Schedule,
				Suspend:          cj.Spec.Suspend,
				LastScheduleTime: lastSchedule,
				Labels:           cj.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetCronJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	outputFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	cj, err := clientset.BatchV1().CronJobs(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting cronjob %s/%s: %v", ns, name, err)), nil
	}

	if outputFmt != "" {
		result, err := outpkg.Format(outputFmt, cj)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error formatting output: %v", err)), nil
		}
		return mcp.NewToolResultText(result), nil
	}

	var lastSchedule *string
	if cj.Status.LastScheduleTime != nil {
		t := cj.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		lastSchedule = &t
	}
	output := cronjobData{
		Name:             cj.Name,
		Namespace:        cj.Namespace,
		Schedule:         cj.Spec.Schedule,
		Suspend:          cj.Spec.Suspend,
		LastScheduleTime: lastSchedule,
		Labels:           cj.Labels,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteCronJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.BatchV1().CronJobs(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting cronjob %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("CronJob %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateCronJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	schedule, err := request.RequireString("schedule")
	if err != nil {
		output := fmt.Sprintf("Provide schedule for cronjob (e.g., '*/5 * * * *')")
		return mcp.NewToolResultText(string(output)), nil
	}
	labels := request.GetString("label", "")
	containerName, err := request.RequireString("containerName")
	if err != nil {
		output := fmt.Sprintf("Provide containerName for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	containerImage, err := request.RequireString("containerImage")
	if err != nil {
		output := fmt.Sprintf("Provide containerImage for cronjob")
		return mcp.NewToolResultText(string(output)), nil
	}
	command := request.GetString("command", "")

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}

	lab := make(map[string]string)
	if labels != "" {
		cjLabel := strings.Split(labels, ",")
		for _, label := range cjLabel {
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

	containers := []v1.Container{
		{
			Name:  strings.TrimSpace(containerName),
			Image: strings.TrimSpace(containerImage),
		},
	}

	if command != "" {
		containers[0].Command = strings.Split(command, " ")
	}

	cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    lab,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lab,
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: lab,
						},
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							Containers:    containers,
						},
					},
				},
			},
		},
	}

	createCronJob, err := clientset.BatchV1().CronJobs(ns).Create(ctx, cronjob, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating cronjob %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully cronjob %s/%s is created", createCronJob.Namespace, createCronJob.Name)
	return mcp.NewToolResultText(string(output)), nil
}
