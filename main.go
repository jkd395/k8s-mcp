package main

import (
	"flag"
	"fmt"
	"k8s-mcp/kubernetes/clusterhealth"
	"k8s-mcp/kubernetes/clusterrole"
	"k8s-mcp/kubernetes/clusterrolebinding"
	"k8s-mcp/kubernetes/componentstatus"
	"k8s-mcp/kubernetes/configmap"
	"k8s-mcp/kubernetes/crd"
	"k8s-mcp/kubernetes/createresource"
	"k8s-mcp/kubernetes/cronjob"
	"k8s-mcp/kubernetes/custom"
	"k8s-mcp/kubernetes/daemonset"
	"k8s-mcp/kubernetes/deployment"
	"k8s-mcp/kubernetes/endpoint"
	"k8s-mcp/kubernetes/event"
	"k8s-mcp/kubernetes/hpa"
	"k8s-mcp/kubernetes/ingress"
	"k8s-mcp/kubernetes/job"
	"k8s-mcp/kubernetes/limitrange"
	"k8s-mcp/kubernetes/namespace"
	"k8s-mcp/kubernetes/networkpolicy"
	"k8s-mcp/kubernetes/node"
	"k8s-mcp/kubernetes/pod"
	"k8s-mcp/kubernetes/pv"
	"k8s-mcp/kubernetes/pvc"
	"k8s-mcp/kubernetes/resourcequota"
	"k8s-mcp/kubernetes/role"
	"k8s-mcp/kubernetes/rolebinding"
	"k8s-mcp/kubernetes/secret"
	"k8s-mcp/kubernetes/service"
	"k8s-mcp/kubernetes/serviceaccount"
	"k8s-mcp/kubernetes/statefulset"
	"k8s-mcp/kubernetes/storageclass"
	"k8s-mcp/kubernetes/top"
	"k8s-mcp/tools"
	"log"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

var (
	mode   string
	apiKey string
)

func init() {
	flag.StringVar(&mode, "mode", "http", "MCP server mode")
	flag.StringVar(&apiKey, "apiKey", "", "API key for HTTP authentication (optional)")
}

// authMiddleware wraps an http.Handler with API key authentication.
// If no apiKey is configured, it passes through all requests.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiKey != "" {
			// Check Authorization: Bearer <key> or X-API-Key header
			authHeader := r.Header.Get("Authorization")
			xApiKey := r.Header.Get("X-API-Key")

			var providedKey string
			if strings.HasPrefix(authHeader, "Bearer ") {
				providedKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else if xApiKey != "" {
				providedKey = xApiKey
			}

			if providedKey != apiKey {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	s := server.NewMCPServer(
		"Kubernetes MCP",
		"1.0.0",
	)

	flag.Parse()

	// Pod tools
	s.AddTool(tools.ListPod, pod.ListPod)
	s.AddTool(tools.GetPod, pod.GetPod)
	s.AddTool(tools.DeletePod, pod.DeletePod)
	s.AddTool(tools.UpdatePod, pod.UpdatePod)
	s.AddTool(tools.CreatePod, pod.CreatePod)
	s.AddTool(tools.PodLog, pod.PodLog)

	// Namespace tools
	s.AddTool(tools.ListNS, namespace.ListNS)
	s.AddTool(tools.GetNS, namespace.GetNS)
	s.AddTool(tools.DeleteNS, namespace.DeleteNS)
	s.AddTool(tools.UpdateNS, namespace.UpdateNS)
	s.AddTool(tools.CreateNS, namespace.CreateNS)

	// Node tools
	s.AddTool(tools.ListNode, node.ListNode)
	s.AddTool(tools.GetNode, node.GetNode)
	s.AddTool(tools.DeleteNode, node.DeleteNode)
	s.AddTool(tools.UpdateNode, node.UpdateNode)

	// Deployment tools
	s.AddTool(tools.ListDeployment, deployment.ListDeployment)
	s.AddTool(tools.GetDeployment, deployment.GetDeployment)
	s.AddTool(tools.DeleteDeployment, deployment.DeleteDeployment)
	s.AddTool(tools.CreateDeployment, deployment.CreateDeployment)
	s.AddTool(tools.UpdateDeployment, deployment.UpdateDeployment)

	// Daemonset tools
	s.AddTool(tools.ListDaemonset, daemonset.ListDaemonset)
	s.AddTool(tools.GetDaemonset, daemonset.GetDaemonset)
	s.AddTool(tools.DeleteDaemonset, daemonset.DeleteDaemonset)
	s.AddTool(tools.UpdateDaemonset, daemonset.UpdateDaemonset)
	s.AddTool(tools.CreateDaemonset, daemonset.CreateDaemonset)

	// Statefulset tools
	s.AddTool(tools.ListStatefulset, statefulset.ListStatefulset)
	s.AddTool(tools.GetStatefulset, statefulset.GetStatefulset)
	s.AddTool(tools.DeleteStatefulset, statefulset.DeleteStatefulset)
	s.AddTool(tools.UpdateStatefulset, statefulset.UpdateStatefulset)
	s.AddTool(tools.CreateStatefulset, statefulset.CreateStatefulset)

	// Service tools
	s.AddTool(tools.ListService, service.ListService)
	s.AddTool(tools.GetService, service.GetService)
	s.AddTool(tools.DeleteService, service.DeleteService)
	s.AddTool(tools.UpdateService, service.UpdateService)
	s.AddTool(tools.CreateService, service.CreateService)

	// Configmap tools
	s.AddTool(tools.ListConfigmap, configmap.ListConfigmap)
	s.AddTool(tools.GetConfigmap, configmap.GetConfigmap)
	s.AddTool(tools.DeleteConfigmap, configmap.DeleteConfigmap)
	s.AddTool(tools.CreateConfigmap, configmap.CreateConfigmap)

	// Secret tools
	s.AddTool(tools.ListSecret, secret.ListSecret)
	s.AddTool(tools.GetSecret, secret.GetSecret)
	s.AddTool(tools.DeleteSecret, secret.DeleteSecret)
	s.AddTool(tools.CreateSecret, secret.CreateSecret)

	// ServiceAccount tools
	s.AddTool(tools.ListSA, serviceaccount.ListSA)
	s.AddTool(tools.GetSA, serviceaccount.GetSA)
	s.AddTool(tools.DeleteSA, serviceaccount.DeleteSA)
	s.AddTool(tools.CreateSA, serviceaccount.CreateSA)

	// Role tools
	s.AddTool(tools.ListRole, role.ListRole)
	s.AddTool(tools.GetRole, role.GetRole)
	s.AddTool(tools.DeleteRole, role.DeleteRole)

	// RoleBinding tools
	s.AddTool(tools.ListRB, rolebinding.ListRB)
	s.AddTool(tools.GetRB, rolebinding.GetRB)
	s.AddTool(tools.DeleteRB, rolebinding.DeleteRB)

	// PVC tools
	s.AddTool(tools.ListPVC, pvc.ListPVC)
	s.AddTool(tools.GetPVC, pvc.GetPVC)
	s.AddTool(tools.DeletePVC, pvc.DeletePVC)
	s.AddTool(tools.UpdatePVC, pvc.UpdatePVC)
	s.AddTool(tools.CreatePVC, pvc.CreatePVC)

	// PV tools
	s.AddTool(tools.ListPV, pv.ListPV)
	s.AddTool(tools.GetPV, pv.GetPV)
	s.AddTool(tools.DeletePV, pv.DeletePV)

	// ClusterRole tools
	s.AddTool(tools.ListCR, clusterrole.ListCR)
	s.AddTool(tools.GetCR, clusterrole.GetCR)
	s.AddTool(tools.DeleteCR, clusterrole.DeleteCR)

	// ClusterRoleBinding tools
	s.AddTool(tools.ListCRB, clusterrolebinding.ListCRB)
	s.AddTool(tools.GetCRB, clusterrolebinding.GetCRB)
	s.AddTool(tools.DeleteCRB, clusterrolebinding.DeleteCRB)

	// StorageClass tools
	s.AddTool(tools.ListSC, storageclass.ListSC)
	s.AddTool(tools.GetSC, storageclass.GetSC)
	s.AddTool(tools.DeleteSC, storageclass.DeleteSC)

	// CRD tools
	s.AddTool(tools.ListCRD, crd.ListCRD)
	s.AddTool(tools.GetCRD, crd.GetCRD)
	s.AddTool(tools.DeleteCRD, crd.DeleteCRD)
	s.AddTool(tools.CreateCRDWithJson, crd.CreateCRDWithJson)

	// Create Resource with Json tool
	s.AddTool(tools.CreateResourceWithJSon, createresource.CreateResourceWithJson)

	// Custom tool
	s.AddTool(tools.Custom, custom.Custom)

	// Ingress tools
	s.AddTool(tools.ListIngress, ingress.ListIngress)
	s.AddTool(tools.GetIngress, ingress.GetIngress)
	s.AddTool(tools.DeleteIngress, ingress.DeleteIngress)
	s.AddTool(tools.CreateIngress, ingress.CreateIngress)

	// HPA tools
	s.AddTool(tools.ListHPA, hpa.ListHPA)
	s.AddTool(tools.GetHPA, hpa.GetHPA)
	s.AddTool(tools.DeleteHPA, hpa.DeleteHPA)
	s.AddTool(tools.CreateHPA, hpa.CreateHPA)

	// Job tools
	s.AddTool(tools.ListJob, job.ListJob)
	s.AddTool(tools.GetJob, job.GetJob)
	s.AddTool(tools.DeleteJob, job.DeleteJob)
	s.AddTool(tools.CreateJob, job.CreateJob)

	// CronJob tools
	s.AddTool(tools.ListCronJob, cronjob.ListCronJob)
	s.AddTool(tools.GetCronJob, cronjob.GetCronJob)
	s.AddTool(tools.DeleteCronJob, cronjob.DeleteCronJob)
	s.AddTool(tools.CreateCronJob, cronjob.CreateCronJob)

	// NetworkPolicy tools
	s.AddTool(tools.ListNetworkPolicy, networkpolicy.ListNetworkPolicy)
	s.AddTool(tools.GetNetworkPolicy, networkpolicy.GetNetworkPolicy)
	s.AddTool(tools.DeleteNetworkPolicy, networkpolicy.DeleteNetworkPolicy)
	s.AddTool(tools.CreateNetworkPolicy, networkpolicy.CreateNetworkPolicy)

	// Event tools
	s.AddTool(tools.ListEvent, event.ListEvent)
	s.AddTool(tools.GetEvent, event.GetEvent)

	// ResourceQuota tools
	s.AddTool(tools.ListResourceQuota, resourcequota.ListResourceQuota)
	s.AddTool(tools.GetResourceQuota, resourcequota.GetResourceQuota)

	// LimitRange tools
	s.AddTool(tools.ListLimitRange, limitrange.ListLimitRange)
	s.AddTool(tools.GetLimitRange, limitrange.GetLimitRange)

	// Endpoint tools
	s.AddTool(tools.ListEndpoint, endpoint.ListEndpoint)
	s.AddTool(tools.GetEndpoint, endpoint.GetEndpoint)

	// ComponentStatus tools
	s.AddTool(tools.ListComponentStatus, componentstatus.ListComponentStatus)
	s.AddTool(tools.GetComponentStatus, componentstatus.GetComponentStatus)

	// Cluster Health tools
	s.AddTool(tools.GetClusterHealth, clusterhealth.GetClusterHealth)
	s.AddTool(tools.ListNodeHealth, clusterhealth.ListNodeHealth)

	// Top tools (metrics)
	s.AddTool(tools.TopPod, top.TopPod)
	s.AddTool(tools.TopNode, top.TopNode)

	if mode == "http" {
		handler := server.NewStreamableHTTPServer(s, server.WithDisableStreaming(true))
		mux := http.NewServeMux()
		mux.HandleFunc("/mcp", handler.ServeHTTP)
		log.Println("Starting Kubernetes MCP Server")
		if apiKey != "" {
			log.Println("Authentication enabled: API key required via Authorization: Bearer <key> or X-API-Key header")
		} else {
			log.Println("Authentication disabled: server is open to all requests")
		}
		if err := http.ListenAndServe(":8080", authMiddleware(mux)); err != nil {
			fmt.Printf("Error starting http server: %v\n", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Error starting stdio server: %v\n", err)
		}
	}
}
