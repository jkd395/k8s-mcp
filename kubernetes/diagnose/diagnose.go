package diagnose

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"k8s-mcp/kubernetes/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func DescribePod(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for pod"), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for pod"), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing client: %v", err)), nil
	}

	pod, err := clientset.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting pod %s/%s: %v", ns, name, err)), nil
	}

	var lines []string

	// Basic info
	lines = append(lines, fmt.Sprintf("Name:         %s", pod.Name))
	lines = append(lines, fmt.Sprintf("Namespace:    %s", pod.Namespace))
	lines = append(lines, fmt.Sprintf("Node:         %s", pod.Spec.NodeName))
	lines = append(lines, fmt.Sprintf("Start Time:   %s", pod.CreationTimestamp.Format(time.RFC3339)))
	lines = append(lines, fmt.Sprintf("Status:       %s", string(pod.Status.Phase)))
	lines = append(lines, fmt.Sprintf("IP:           %s", pod.Status.PodIP))
	if pod.Status.QOSClass != "" {
		lines = append(lines, fmt.Sprintf("QoS Class:    %s", pod.Status.QOSClass))
	}
	lines = append(lines, "")

	// Labels
	if len(pod.Labels) > 0 {
		lines = append(lines, "Labels:")
		for k, v := range pod.Labels {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Pod Conditions
	if len(pod.Status.Conditions) > 0 {
		lines = append(lines, "Conditions:")
		for _, c := range pod.Status.Conditions {
			status := "False"
			if c.Status == v1.ConditionTrue {
				status = "True"
			}
			lines = append(lines, fmt.Sprintf("  %-30s %s", string(c.Type), status))
			if c.Reason != "" {
				lines = append(lines, fmt.Sprintf("    Reason:  %s", c.Reason))
			}
			if c.Message != "" {
				lines = append(lines, fmt.Sprintf("    Message: %s", c.Message))
			}
		}
		lines = append(lines, "")
	}

	// Containers
	lines = append(lines, "Containers:")
	for i, c := range pod.Spec.Containers {
		lines = append(lines, fmt.Sprintf("  %s:", c.Name))
		lines = append(lines, fmt.Sprintf("    Image:         %s", c.Image))
		if len(c.Command) > 0 {
			lines = append(lines, fmt.Sprintf("    Command:       %s", strings.Join(c.Command, " ")))
		}
		if len(c.Args) > 0 {
			lines = append(lines, fmt.Sprintf("    Args:          %s", strings.Join(c.Args, " ")))
		}
		if len(c.Ports) > 0 {
			for _, p := range c.Ports {
				lines = append(lines, fmt.Sprintf("    Port:          %s:%d/%s", p.Name, p.ContainerPort, p.Protocol))
			}
		}
		// Resource requests/limits
		if c.Resources.Requests != nil || c.Resources.Limits != nil {
			reqCPU := c.Resources.Requests.Cpu().String()
			reqMem := c.Resources.Requests.Memory().String()
			limCPU := c.Resources.Limits.Cpu().String()
			limMem := c.Resources.Limits.Memory().String()
			if reqCPU != "0" || reqMem != "0" {
				lines = append(lines, fmt.Sprintf("    Requests:     cpu=%s, memory=%s", reqCPU, reqMem))
			}
			if limCPU != "0" || limMem != "0" {
				lines = append(lines, fmt.Sprintf("    Limits:       cpu=%s, memory=%s", limCPU, limMem))
			}
		}

		// Container status
		if i < len(pod.Status.ContainerStatuses) {
			cs := pod.Status.ContainerStatuses[i]
			lines = append(lines, fmt.Sprintf("    Ready:         %v", cs.Ready))
			lines = append(lines, fmt.Sprintf("    Restart Count: %d", cs.RestartCount))
			if cs.State.Waiting != nil {
				lines = append(lines, fmt.Sprintf("    State:         Waiting (%s: %s)", cs.State.Waiting.Reason, cs.State.Waiting.Message))
			}
			if cs.State.Running != nil {
				lines = append(lines, fmt.Sprintf("    State:         Running (started: %s)", cs.State.Running.StartedAt.Format(time.RFC3339)))
			}
			if cs.State.Terminated != nil {
				lines = append(lines, fmt.Sprintf("    State:         Terminated (reason: %s, exit: %d)", cs.State.Terminated.Reason, cs.State.Terminated.ExitCode))
				if cs.State.Terminated.Message != "" {
					lines = append(lines, fmt.Sprintf("    Message:       %s", cs.State.Terminated.Message))
				}
			}
			if cs.LastTerminationState.Terminated != nil {
				lt := cs.LastTerminationState.Terminated
				lines = append(lines, fmt.Sprintf("    Last State:    Terminated (reason: %s, exit: %d)", lt.Reason, lt.ExitCode))
			}
		}
		lines = append(lines, "")
	}

	// Init Containers
	if len(pod.Spec.InitContainers) > 0 {
		lines = append(lines, "Init Containers:")
		for i, c := range pod.Spec.InitContainers {
			lines = append(lines, fmt.Sprintf("  %s:", c.Name))
			lines = append(lines, fmt.Sprintf("    Image:         %s", c.Image))
			if i < len(pod.Status.InitContainerStatuses) {
				cs := pod.Status.InitContainerStatuses[i]
				lines = append(lines, fmt.Sprintf("    Restart Count: %d", cs.RestartCount))
				if cs.State.Waiting != nil {
					lines = append(lines, fmt.Sprintf("    State:         Waiting (%s: %s)", cs.State.Waiting.Reason, cs.State.Waiting.Message))
				}
				if cs.State.Running != nil {
					lines = append(lines, fmt.Sprintf("    State:         Running"))
				}
				if cs.State.Terminated != nil {
					lines = append(lines, fmt.Sprintf("    State:         Terminated (reason: %s, exit: %d)", cs.State.Terminated.Reason, cs.State.Terminated.ExitCode))
				}
			}
			lines = append(lines, "")
		}
	}

	// Events
	events, err := clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err == nil && len(events.Items) > 0 {
		lines = append(lines, fmt.Sprintf("Events:"))
		sort.Slice(events.Items, func(i, j int) bool {
			return events.Items[i].LastTimestamp.After(events.Items[j].LastTimestamp.Time)
		})
		for _, e := range events.Items {
			lines = append(lines, fmt.Sprintf("  %s  %s  %s  %s  (%d times)",
				e.LastTimestamp.Format("2006-01-02 15:04:05"),
				e.Type,
				e.Reason,
				e.Message,
				e.Count))
		}
		lines = append(lines, "")
	}

	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func DescribeNode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for node"), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing client: %v", err)), nil
	}

	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting node %s: %v", name, err)), nil
	}

	var lines []string

	// Basic info
	lines = append(lines, fmt.Sprintf("Name:             %s", node.Name))
	lines = append(lines, fmt.Sprintf("Creation:         %s", node.CreationTimestamp.Format(time.RFC3339)))
	lines = append(lines, fmt.Sprintf("Kubelet:          %s", node.Status.NodeInfo.KubeletVersion))
	lines = append(lines, fmt.Sprintf("OS Image:         %s", node.Status.NodeInfo.OSImage))
	lines = append(lines, fmt.Sprintf("Kernel:           %s", node.Status.NodeInfo.KernelVersion))
	lines = append(lines, fmt.Sprintf("Architecture:     %s", node.Status.NodeInfo.Architecture))
	lines = append(lines, fmt.Sprintf("PodCIDR:          %s", node.Spec.PodCIDR))
	lines = append(lines, fmt.Sprintf("ProviderID:       %s", node.Spec.ProviderID))
	lines = append(lines, "")

	// Labels
	if len(node.Labels) > 0 {
		lines = append(lines, "Labels:")
		for k, v := range node.Labels {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Annotations
	if len(node.Annotations) > 0 {
		lines = append(lines, "Annotations:")
		keys := make([]string, 0, len(node.Annotations))
		for k := range node.Annotations {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := node.Annotations[k]
			if len(v) > 100 {
				v = v[:100] + "..."
			}
			lines = append(lines, fmt.Sprintf("  %s: %s", k, v))
		}
		lines = append(lines, "")
	}

	// Taints
	if len(node.Spec.Taints) > 0 {
		lines = append(lines, "Taints:")
		for _, t := range node.Spec.Taints {
			lines = append(lines, fmt.Sprintf("  %s=%s:%s", t.Key, t.Value, t.Effect))
		}
		lines = append(lines, "")
	}

	// Conditions
	if len(node.Status.Conditions) > 0 {
		lines = append(lines, "Conditions:")
		for _, c := range node.Status.Conditions {
			status := "False"
			if c.Status == v1.ConditionTrue {
				status = "True"
			}
			lines = append(lines, fmt.Sprintf("  %-30s %s", string(c.Type), status))
			lines = append(lines, fmt.Sprintf("    Reason:  %s", c.Reason))
			lines = append(lines, fmt.Sprintf("    Message: %s", c.Message))
			lines = append(lines, fmt.Sprintf("    Last:    %s", c.LastHeartbeatTime.Format(time.RFC3339)))
		}
		lines = append(lines, "")
	}

	// Capacity
	lines = append(lines, "Capacity:")
	lines = append(lines, fmt.Sprintf("  cpu:                 %s", node.Status.Capacity.Cpu().String()))
	lines = append(lines, fmt.Sprintf("  memory:              %s", node.Status.Capacity.Memory().String()))
	lines = append(lines, fmt.Sprintf("  pods:                %s", node.Status.Capacity.Pods().String()))
	if ephemeral, ok := node.Status.Capacity["ephemeral-storage"]; ok {
		lines = append(lines, fmt.Sprintf("  ephemeral-storage:   %s", ephemeral.String()))
	}
	lines = append(lines, "")

	// Allocatable
	lines = append(lines, "Allocatable:")
	lines = append(lines, fmt.Sprintf("  cpu:                 %s", node.Status.Allocatable.Cpu().String()))
	lines = append(lines, fmt.Sprintf("  memory:              %s", node.Status.Allocatable.Memory().String()))
	lines = append(lines, fmt.Sprintf("  pods:                %s", node.Status.Allocatable.Pods().String()))
	if ephemeral, ok := node.Status.Allocatable["ephemeral-storage"]; ok {
		lines = append(lines, fmt.Sprintf("  ephemeral-storage:   %s", ephemeral.String()))
	}
	lines = append(lines, "")

	// Pods running on this node
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", name),
	})
	if err == nil {
		lines = append(lines, fmt.Sprintf("Pods: (%d total)", len(pods.Items)))
		type podOnNode struct {
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
			Status    string `json:"status"`
		}
		var podList []podOnNode
		for _, p := range pods.Items {
			podList = append(podList, podOnNode{
				Namespace: p.Namespace,
				Name:      p.Name,
				Status:    string(p.Status.Phase),
			})
		}
		sort.Slice(podList, func(i, j int) bool {
			if podList[i].Namespace != podList[j].Namespace {
				return podList[i].Namespace < podList[j].Namespace
			}
			return podList[i].Name < podList[j].Name
		})
		for _, p := range podList {
			lines = append(lines, fmt.Sprintf("  %s/%s - %s", p.Namespace, p.Name, p.Status))
		}
	}

	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func ListNodePods(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeName, err := request.RequireString("node")
	if err != nil {
		return mcp.NewToolResultText("Provide node name"), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing client: %v", err)), nil
	}

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error listing pods on node %s: %v", nodeName, err)), nil
	}

	type podInfo struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
		Status    string `json:"status"`
		Node      string `json:"node"`
	}

	var out []podInfo
	for _, p := range pods.Items {
		out = append(out, podInfo{
			Namespace: p.Namespace,
			Name:      p.Name,
			Status:    string(p.Status.Phase),
			Node:      p.Spec.NodeName,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Namespace != out[j].Namespace {
			return out[i].Namespace < out[j].Namespace
		}
		return out[i].Name < out[j].Name
	})

	b, _ := json.MarshalIndent(out, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func DescribeService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for service"), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for service"), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing client: %v", err)), nil
	}

	svc, err := clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting service %s/%s: %v", ns, name, err)), nil
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("Name:              %s", svc.Name))
	lines = append(lines, fmt.Sprintf("Namespace:         %s", svc.Namespace))
	lines = append(lines, fmt.Sprintf("Type:              %s", svc.Spec.Type))
	lines = append(lines, fmt.Sprintf("ClusterIP:         %s", svc.Spec.ClusterIP))
	if svc.Spec.ClusterIPs != nil {
		lines = append(lines, fmt.Sprintf("ClusterIPs:        %s", strings.Join(svc.Spec.ClusterIPs, ", ")))
	}
	if svc.Spec.ExternalIPs != nil {
		lines = append(lines, fmt.Sprintf("ExternalIPs:       %s", strings.Join(svc.Spec.ExternalIPs, ", ")))
	}
	if svc.Spec.ExternalName != "" {
		lines = append(lines, fmt.Sprintf("ExternalName:      %s", svc.Spec.ExternalName))
	}
	if svc.Spec.LoadBalancerIP != "" {
		lines = append(lines, fmt.Sprintf("LoadBalancerIP:    %s", svc.Spec.LoadBalancerIP))
	}
	if svc.Spec.SessionAffinity != "" {
		lines = append(lines, fmt.Sprintf("Session Affinity:  %s", svc.Spec.SessionAffinity))
	}
	if svc.Spec.ExternalTrafficPolicy != "" {
		lines = append(lines, fmt.Sprintf("External Traffic:  %s", svc.Spec.ExternalTrafficPolicy))
	}
	lines = append(lines, "")

	// Labels
	if len(svc.Labels) > 0 {
		lines = append(lines, "Labels:")
		for k, v := range svc.Labels {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Selector
	if len(svc.Spec.Selector) > 0 {
		lines = append(lines, "Selector:")
		for k, v := range svc.Spec.Selector {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Ports
	if len(svc.Spec.Ports) > 0 {
		lines = append(lines, "Ports:")
		for _, p := range svc.Spec.Ports {
			lines = append(lines, fmt.Sprintf("  %s  %d/%s -> %d", p.Name, p.Port, p.Protocol, p.TargetPort.IntVal))
		}
		lines = append(lines, "")
	}

	// Endpoints
	eps, err := clientset.CoreV1().Endpoints(ns).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		lines = append(lines, "Endpoints:")
		for _, subset := range eps.Subsets {
			var addrs []string
			for _, a := range subset.Addresses {
				if a.NodeName != nil {
					addrs = append(addrs, fmt.Sprintf("%s (node: %s)", a.IP, *a.NodeName))
				} else {
					addrs = append(addrs, a.IP)
				}
			}
			for _, p := range subset.Ports {
				lines = append(lines, fmt.Sprintf("  %s:%d -> %d pods", p.Name, p.Port, len(subset.Addresses)))
			}
			for _, addr := range addrs {
				lines = append(lines, fmt.Sprintf("    - %s", addr))
			}
		}
		if len(eps.Subsets) == 0 {
			lines = append(lines, "  <none>")
		}
		lines = append(lines, "")
	}

	// Events
	events, err := clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Service", name),
	})
	if err == nil && len(events.Items) > 0 {
		lines = append(lines, "Events:")
		sort.Slice(events.Items, func(i, j int) bool {
			return events.Items[i].LastTimestamp.After(events.Items[j].LastTimestamp.Time)
		})
		for _, e := range events.Items {
			lines = append(lines, fmt.Sprintf("  %s  %s  %s  %s",
				e.LastTimestamp.Format("2006-01-02 15:04:05"),
				e.Type,
				e.Reason,
				e.Message))
		}
	}

	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func DescribeDeployment(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for deployment"), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for deployment"), nil
	}

	clientset, _, _, _, _, err := client.InitializeClients()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing client: %v", err)), nil
	}

	deploy, err := clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting deployment %s/%s: %v", ns, name, err)), nil
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("Name:               %s", deploy.Name))
	lines = append(lines, fmt.Sprintf("Namespace:          %s", deploy.Namespace))
	lines = append(lines, fmt.Sprintf("Creation:           %s", deploy.CreationTimestamp.Format(time.RFC3339)))
	lines = append(lines, fmt.Sprintf("Strategy:           %s", string(deploy.Spec.Strategy.Type)))
	if deploy.Spec.Strategy.RollingUpdate != nil {
		if deploy.Spec.Strategy.RollingUpdate.MaxSurge != nil {
			lines = append(lines, fmt.Sprintf("Max Surge:          %s", deploy.Spec.Strategy.RollingUpdate.MaxSurge.String()))
		}
		if deploy.Spec.Strategy.RollingUpdate.MaxUnavailable != nil {
			lines = append(lines, fmt.Sprintf("Max Unavailable:    %s", deploy.Spec.Strategy.RollingUpdate.MaxUnavailable.String()))
		}
	}
	lines = append(lines, fmt.Sprintf("Replicas:           %d desired | %d updated | %d total | %d available | %d unavailable",
		*deploy.Spec.Replicas,
		deploy.Status.UpdatedReplicas,
		deploy.Status.Replicas,
		deploy.Status.AvailableReplicas,
		deploy.Status.UnavailableReplicas))
	if deploy.Spec.RevisionHistoryLimit != nil {
		lines = append(lines, fmt.Sprintf("Revision History:   %d (limit)", *deploy.Spec.RevisionHistoryLimit))
	}
	if deploy.Spec.MinReadySeconds > 0 {
		lines = append(lines, fmt.Sprintf("Min Ready:          %ds", deploy.Spec.MinReadySeconds))
	}
	if deploy.Status.ReadyReplicas != deploy.Status.Replicas {
		lines = append(lines, fmt.Sprintf("Rollout Status:     INCOMPLETE (ready: %d/%d)", deploy.Status.ReadyReplicas, deploy.Status.Replicas))
	} else {
		lines = append(lines, "Rollout Status:     COMPLETE")
	}
	lines = append(lines, "")

	// Labels
	if len(deploy.Labels) > 0 {
		lines = append(lines, "Labels:")
		for k, v := range deploy.Labels {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Selector
	if deploy.Spec.Selector != nil && len(deploy.Spec.Selector.MatchLabels) > 0 {
		lines = append(lines, "Selector:")
		for k, v := range deploy.Spec.Selector.MatchLabels {
			lines = append(lines, fmt.Sprintf("  %s=%s", k, v))
		}
		lines = append(lines, "")
	}

	// Conditions
	if len(deploy.Status.Conditions) > 0 {
		lines = append(lines, "Conditions:")
		for _, c := range deploy.Status.Conditions {
			status := "False"
			if c.Status == v1.ConditionTrue {
				status = "True"
			}
			lines = append(lines, fmt.Sprintf("  %-25s %s", string(c.Type), status))
			if c.Reason != "" {
				lines = append(lines, fmt.Sprintf("    Reason:     %s", c.Reason))
			}
			if c.Message != "" {
				lines = append(lines, fmt.Sprintf("    Message:    %s", c.Message))
			}
			lines = append(lines, fmt.Sprintf("    Last:       %s", c.LastUpdateTime.Format(time.RFC3339)))
		}
		lines = append(lines, "")
	}

	// Containers
	lines = append(lines, "Containers:")
	for _, c := range deploy.Spec.Template.Spec.Containers {
		lines = append(lines, fmt.Sprintf("  %s:", c.Name))
		lines = append(lines, fmt.Sprintf("    Image:     %s", c.Image))
		if c.Resources.Requests != nil || c.Resources.Limits != nil {
			reqCPU := c.Resources.Requests.Cpu().String()
			reqMem := c.Resources.Requests.Memory().String()
			limCPU := c.Resources.Limits.Cpu().String()
			limMem := c.Resources.Limits.Memory().String()
			if reqCPU != "0" || reqMem != "0" {
				lines = append(lines, fmt.Sprintf("    Requests:  cpu=%s, memory=%s", reqCPU, reqMem))
			}
			if limCPU != "0" || limMem != "0" {
				lines = append(lines, fmt.Sprintf("    Limits:    cpu=%s, memory=%s", limCPU, limMem))
			}
		}
	}
	lines = append(lines, "")

	// Pod status summary
	pods, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deploy.Spec.Selector),
	})
	if err == nil {
		phaseCount := map[string]int{}
		for _, p := range pods.Items {
			phaseCount[string(p.Status.Phase)]++
		}
		lines = append(lines, fmt.Sprintf("Pods: (%d total)", len(pods.Items)))
		for _, phase := range []string{"Running", "Pending", "Succeeded", "Failed", "Unknown"} {
			if c, ok := phaseCount[phase]; ok && c > 0 {
				lines = append(lines, fmt.Sprintf("  %-10s: %d", phase, c))
			}
		}
		lines = append(lines, "")
	}

	// Events
	events, err := clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Deployment", name),
	})
	if err == nil && len(events.Items) > 0 {
		lines = append(lines, "Events:")
		sort.Slice(events.Items, func(i, j int) bool {
			return events.Items[i].LastTimestamp.After(events.Items[j].LastTimestamp.Time)
		})
		for _, e := range events.Items {
			lines = append(lines, fmt.Sprintf("  %s  %s  %s  %s",
				e.LastTimestamp.Format("2006-01-02 15:04:05"),
				e.Type,
				e.Reason,
				e.Message))
		}
	}

	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func CheckAPIServerHealth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	restConfig, err := client.GetRestConfig()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting rest config: %v", err)), nil
	}

	transport, err := rest.TransportFor(restConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error creating transport: %v", err)), nil
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	host := restConfig.Host
	var lines []string
	lines = append(lines, fmt.Sprintf("API Server: %s", host))
	lines = append(lines, "")

	endpoints := []string{"/healthz?verbose", "/livez?verbose", "/readyz?verbose"}
	for _, ep := range endpoints {
		url := host + ep
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			lines = append(lines, fmt.Sprintf("[ERROR] %s: %v", ep, err))
			continue
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			lines = append(lines, fmt.Sprintf("[UNREACHABLE] %s: %v", ep, err))
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		status := "HEALTHY"
		if resp.StatusCode != 200 {
			status = fmt.Sprintf("UNHEALTHY (HTTP %d)", resp.StatusCode)
		}
		lines = append(lines, fmt.Sprintf(">>> %s", ep))
		lines = append(lines, fmt.Sprintf("  Status: %s", status))

		for _, line := range strings.Split(string(body), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if strings.Contains(line, "] ok") || strings.Contains(line, "healthy") {
				lines = append(lines, fmt.Sprintf("  %s", line))
				continue
			}
			lines = append(lines, fmt.Sprintf("  %s", line))
		}
		lines = append(lines, "")
	}

	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}

func CheckAPIServerMetrics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	restConfig, err := client.GetRestConfig()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting rest config: %v", err)), nil
	}

	transport, err := rest.TransportFor(restConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error creating transport: %v", err)), nil
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", restConfig.Host+"/metrics", nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error creating request: %v", err)), nil
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error fetching /metrics: %v", err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return mcp.NewToolResultText(fmt.Sprintf("Metrics endpoint returned HTTP %d (needs --authorization-always-allow-paths=/metrics or RBAC)",
			resp.StatusCode)), nil
	}

	scanner := bufio.NewScanner(resp.Body)

	inflight := map[string]int64{}
	reqTotal := map[string]int64{}
	reqErr := map[string]int64{}
	durationBuckets := map[string]map[string]int64{}
	reqDuration := map[string]float64{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		if strings.Contains(line, "apiserver_current_inflight_requests{") {
			if v := parseLabelValue(line, "kind"); v != "" {
				if val := extractValue(line); val >= 0 {
					inflight[v] = int64(val)
				}
			}
		}

		if strings.Contains(line, "apiserver_request_total{") {
			code := parseLabelValue(line, "code")
			verb := parseLabelValue(line, "verb")
			if code != "" && verb != "" {
				key := verb + "/" + code
				val := extractValue(line)
				reqTotal[key] += int64(val)
				if strings.HasPrefix(code, "5") || strings.HasPrefix(code, "4") {
					reqErr[key] += int64(val)
				}
			}
		}

		if strings.Contains(line, "apiserver_request_duration_seconds_bucket{") {
			verb := parseLabelValue(line, "verb")
			le := parseLabelValue(line, "le")
			if verb != "" && le != "" {
				if _, ok := durationBuckets[verb]; !ok {
					durationBuckets[verb] = map[string]int64{}
				}
				durationBuckets[verb][le] += int64(extractValue(line))
			}
		}

		if strings.Contains(line, "apiserver_request_duration_seconds_sum{") {
			verb := parseLabelValue(line, "verb")
			if verb != "" {
				reqDuration[verb] += extractValue(line)
			}
		}
	}

	var out []string
	out = append(out, "=== API Server Metrics ===")
	out = append(out, "")

	// Inflight requests
	out = append(out, "--- Current Inflight Requests ---")
	for _, kind := range []string{"mutating", "readOnly"} {
		if v, ok := inflight[kind]; ok {
			label := "readOnly"
			if kind == "mutating" {
				label = "mutating"
			}
			out = append(out, fmt.Sprintf("  %-10s: %d", label, v))
		}
	}
	out = append(out, "")

	// Request rate by verb/status
	out = append(out, "--- Request Counts ---")
	verbs := map[string]int64{}
	totalReq := int64(0)
	totalErr := int64(0)
	for key, count := range reqTotal {
		parts := strings.SplitN(key, "/", 2)
		verb := parts[0]
		verbs[verb] += count
		totalReq += count
	}
	for _, count := range reqErr {
		totalErr += count
	}
	out = append(out, fmt.Sprintf("  Total requests:   %d", totalReq))
	out = append(out, fmt.Sprintf("  Error responses:  %d (4xx+5xx)", totalErr))
	if totalReq > 0 {
		out = append(out, fmt.Sprintf("  Error rate:       %.2f%%", float64(totalErr)/float64(totalReq)*100))
	}
	out = append(out, "")
	out = append(out, "  By Verb:")
	for _, v := range []string{"GET", "LIST", "WATCH", "POST", "PUT", "PATCH", "DELETE"} {
		if c, ok := verbs[v]; ok {
			out = append(out, fmt.Sprintf("    %-8s: %d", v, c))
		}
	}
	out = append(out, "")

	// Latency
	out = append(out, "--- Request Latency (seconds) ---")
	if len(durationBuckets) == 0 {
		out = append(out, "  (no latency histogram data)")
	} else {
		for _, verb := range []string{"GET", "LIST", "WATCH", "POST", "PUT", "PATCH", "DELETE"} {
			b, ok := durationBuckets[verb]
			if !ok {
				continue
			}
			total := b["+Inf"]
			if total == 0 {
				continue
			}
			out = append(out, fmt.Sprintf("  %s:", verb))
			out = append(out, fmt.Sprintf("    p50: %.4fs", estimateQuantile(0.50, b)))
			out = append(out, fmt.Sprintf("    p90: %.4fs", estimateQuantile(0.90, b)))
			out = append(out, fmt.Sprintf("    p99: %.4fs", estimateQuantile(0.99, b)))
		}
	}

	out = append(out, "")

	// Top error endpoints (parse apiserver_request_total with resource/verb)
	out = append(out, "--- Top Error Endpoints (by code) ---")
	type errEntry struct {
		code  string
		verb  string
		count int64
	}
	var errList []errEntry
	for key, count := range reqErr {
		parts := strings.SplitN(key, "/", 2)
		errList = append(errList, errEntry{code: parts[1], verb: parts[0], count: count})
	}
	sort.Slice(errList, func(i, j int) bool {
		return errList[i].count > errList[j].count
	})
	for _, e := range errList {
		out = append(out, fmt.Sprintf("  %s %s: %d", e.verb, e.code, e.count))
	}

	return mcp.NewToolResultText(strings.Join(out, "\n")), nil
}

func parseLabelValue(line, label string) string {
	braceIdx := strings.Index(line, "{")
	closeIdx := strings.Index(line, "}")
	if braceIdx < 0 || closeIdx < 0 {
		return ""
	}
	labels := strings.Split(line[braceIdx+1:closeIdx], ",")
	for _, l := range labels {
		if strings.HasPrefix(l, label+"=") {
			return strings.Trim(strings.TrimPrefix(l, label+"="), "\" ")
		}
	}
	return ""
}

func extractValue(line string) float64 {
	closeIdx := strings.LastIndex(line, "}")
	if closeIdx < 0 {
		return 0
	}
	part := strings.TrimSpace(line[closeIdx+1:])
	v, err := strconv.ParseFloat(part, 64)
	if err != nil {
		return 0
	}
	return v
}

func estimateQuantile(q float64, buckets map[string]int64) float64 {
	total := buckets["+Inf"]
	if total == 0 {
		return 0
	}
	target := float64(total) * q
	cumul := int64(0)

	boundaries := []struct {
		le  string
		val float64
	}{
		{"0.001", 0.001}, {"0.0025", 0.0025}, {"0.005", 0.005},
		{"0.01", 0.01}, {"0.025", 0.025}, {"0.05", 0.05},
		{"0.1", 0.1}, {"0.25", 0.25}, {"0.5", 0.5},
		{"1", 1}, {"2.5", 2.5}, {"5", 5}, {"10", 10},
		{"+Inf", 1e18},
	}
	for _, b := range boundaries {
		if v, ok := buckets[b.le]; ok {
			cumul += v
		}
		if float64(cumul) >= target {
			return b.val
		}
	}
	return 0
}
