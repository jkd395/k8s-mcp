package job

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type jobData struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Completions *int32            `json:"completions,omitempty"`
	Succeeded   int32             `json:"succeeded,omitempty"`
	Failed      int32             `json:"failed,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

func ListJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	var output []jobData
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			jobs, err := clientset.BatchV1().Jobs(namespace.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing job in %s: %v", namespace.Name, err)), nil
			}
			for _, job := range jobs.Items {
				output = append(output, jobData{
					Name:        job.Name,
					Namespace:   job.Namespace,
					Completions: job.Spec.Completions,
					Succeeded:   job.Status.Succeeded,
					Failed:      job.Status.Failed,
					Labels:      job.Labels,
				})
			}
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	} else {
		jobs, err := clientset.BatchV1().Jobs(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing job in %s namespace: %v", ns, err)), nil
		}
		for _, job := range jobs.Items {
			output = append(output, jobData{
				Name:        job.Name,
				Namespace:   job.Namespace,
				Completions: job.Spec.Completions,
				Succeeded:   job.Status.Succeeded,
				Failed:      job.Status.Failed,
				Labels:      job.Labels,
			})
		}
		mcpOutput, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(string(mcpOutput)), nil
	}
}

func GetJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	job, err := clientset.BatchV1().Jobs(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting job %s/%s: %v", ns, name, err)), nil
	}
	output := jobData{
		Name:        job.Name,
		Namespace:   job.Namespace,
		Completions: job.Spec.Completions,
		Succeeded:   job.Status.Succeeded,
		Failed:      job.Status.Failed,
		Labels:      job.Labels,
	}
	mcpOutput, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func DeleteJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}
	err = clientset.BatchV1().Jobs(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in deleting job %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Job %s/%s is deleted", ns, name)
	return mcp.NewToolResultText(string(output)), nil
}

func CreateJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		output := fmt.Sprintf("Provide namespace for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		output := fmt.Sprintf("Provide name for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	labels := request.GetString("label", "")
	containerName, err := request.RequireString("containerName")
	if err != nil {
		output := fmt.Sprintf("Provide containerName for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	containerImage, err := request.RequireString("containerImage")
	if err != nil {
		output := fmt.Sprintf("Provide containerImage for job")
		return mcp.NewToolResultText(string(output)), nil
	}
	command := request.GetString("command", "")
	backoffLimit := request.GetInt("backoffLimit", 6)
	completions := request.GetInt("completions", 1)
	parallelism := request.GetInt("parallelism", 1)

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in intialize client: %v", err)), nil
	}

	lab := make(map[string]string)
	if labels != "" {
		jobLabel := strings.Split(labels, ",")
		for _, label := range jobLabel {
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

	completionsInt32 := int32(completions)
	parallelismInt32 := int32(parallelism)
	backoffLimitInt32 := int32(backoffLimit)

	containers := []v1.Container{
		{
			Name:  strings.TrimSpace(containerName),
			Image: strings.TrimSpace(containerImage),
		},
	}

	if command != "" {
		containers[0].Command = strings.Split(command, " ")
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels:    lab,
		},
		Spec: batchv1.JobSpec{
			Completions:  &completionsInt32,
			Parallelism:  &parallelismInt32,
			BackoffLimit: &backoffLimitInt32,
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
	}

	createJob, err := clientset.BatchV1().Jobs(ns).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating job %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully job %s/%s is created", createJob.Namespace, createJob.Name)
	return mcp.NewToolResultText(string(output)), nil
}
