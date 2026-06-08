package job

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
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
	Finalizers  []string          `json:"finalizers,omitempty"`
}
func ListJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	var jobs []batchv1.Job
	if ns == "" {
		namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing namespace: %v", err)), nil
		}
		for _, namespace := range namespaces.Items {
			items, err := clientset.BatchV1().Jobs(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error in listing job in %s: %v", namespace.Name, err)), nil
			}
			jobs = append(jobs, items.Items...)
		}
	} else {
		items, err := clientset.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in listing job in %s namespace: %v", ns, err)), nil
		}
		jobs = items.Items
	}

	if outFmt != "" {
		out, err := output.Format(outFmt, jobs)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(out), nil
	}

	var out []jobData
	for _, job := range jobs {
		out = append(out, jobData{
			Name:        job.Name,
			Namespace:   job.Namespace,
			Completions: job.Spec.Completions,
			Succeeded:   job.Status.Succeeded,
			Failed:      job.Status.Failed,
			Labels:      job.Labels,
			Finalizers:  job.Finalizers,
		})
	}
	mcpOutput, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
	}
	return mcp.NewToolResultText(string(mcpOutput)), nil
}

func GetJob(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Provide namespace for job")), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Provide name for job")), nil
	}
	outFmt := request.GetString("output", "")
	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	job, err := clientset.BatchV1().Jobs(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in getting job %s/%s: %v", ns, name, err)), nil
	}

	if outFmt != "" {
		out, err := output.Format(outFmt, job)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error in marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(out), nil
	}

	raw := jobData{
		Name:        job.Name,
		Namespace:   job.Namespace,
		Completions: job.Spec.Completions,
		Succeeded:   job.Status.Succeeded,
		Failed:      job.Status.Failed,
		Labels:      job.Labels,
		Finalizers:  job.Finalizers,
	}
	mcpOutput, err := json.MarshalIndent(raw, "", " ")
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
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
	}
	err = clientset.BatchV1().Jobs(ns).Delete(ctx, name, metav1.DeleteOptions{})
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
		return mcp.NewToolResultText(fmt.Sprintf("Error in initialize client: %v", err)), nil
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

	createJob, err := clientset.BatchV1().Jobs(ns).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error in creating job %s/%s: %v", ns, name, err)), nil
	}
	output := fmt.Sprintf("Successfully job %s/%s is created", createJob.Namespace, createJob.Name)
	return mcp.NewToolResultText(string(output)), nil
}
