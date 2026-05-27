# K8S MCP SERVER

k8s-mcp is a Golang based Model Context Protocol (MCP) server that expose Kubernetes resources as structured MCP tools, enable AI agents to safely interact with Kubernetes cluster.

## Features

### Supported Kubernetes resources and operations (90 tools):

- Pod: Create, Get, List, Update, Delete and Log.
- Deployment: Create, Get, List, Update and Delete.
- Daemonset: Create, Get, List, Update and Delete.
- Statefulset: Create, Get, List, Update and Delete.
- Namespace: Create, Get, List, Update and Delete.
- Service: Create, Get, List, Update and Delete.
- Configmap: Create, Get, List and Delete.
- Secret: Create, Get, List and Delete.
- Node: List, Get, Update and Delete.
- ServiceAccount: Create, Get, List and Delete.
- PVC: Create, Get, List, Update and Delete.
- PV: List, Get and Delete.
- Role: Get, List and Delete.
- RoleBinding: Get, List and Delete.
- ClusterRole: Get, List and Delete.
- ClusterRoleBinding: Get, List and Delete.
- Storageclass: Get, List and Delete.
- CRD: Create, Get, List and Delete.
- Event: List and Get.
- ResourceQuota: List and Get.
- LimitRange: List and Get.
- Endpoint: List and Get.
- ComponentStatus: List and Get.
- ClusterHealth: Get overall cluster health and node health.
- Top Pod/Node: Show resource usage (requires metrics-server).
- Ingress: Create, Get, List and Delete.
- HPA: Create, Get, List and Delete.
- Job: Create, Get, List and Delete.
- CronJob: Create, Get, List and Delete.
- NetworkPolicy: Create, Get, List and Delete.

Create any kubernetes resource by passing json data.

### Custom Resource:

In addition to predefined kubernetes resource tools listed above, this MCP Server also provides a generic custom tool that allows user to interact with any kubernetes custom resource that is not explicitly supported.

This is useful when:

- You want to work with custom resource(CRDs).
- You want to access new kubernetes resources without updating MCP Server.

The MCP Server forwards the provided parameters to the gRPC backend server, which dynamically resolves the resource and perform the requested action.

Custom Tool Supported Operations: Create, Get, List and Delete.

NOTE: Parameter details for each resource are available in the respective `README.md` file under the kubernetes directory.

## Prerequisites

- Go 1.25+
- Access to kubernetes cluster
- Kubeconfig file
- An application with MCP supported
- (Optional) metrics-server for Top Pod/Node tools

## Installation

```
go install k8s-mcp@latest
```

## Running MCP Server

### stdio mode (Claude Desktop):

```json
{
    "mcpServers": {
        "Kubernetes": {
            "command": "k8s-mcp",
            "args": ["--kubeconfigPath=<Path to kubeconfig file>", "--mode=stdio"]
        }
    }
}
```

Enable the custom tool by using the `--customURL` flag:

```json
{
    "mcpServers": {
        "Kubernetes": {
            "command": "k8s-mcp",
            "args": ["--kubeconfigPath=<Path to kubeconfig file>","--customURL=<grpc server url>", "--mode=stdio"]
        }
    }
}
```

### HTTP mode (remote / Cursor):

```json
{
  "mcpServers": {
    "kubernetes": {
      "url": "http://<server>:8080/mcp"
    }
  }
}
```

## Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-mode` | `http` | Server mode: `http` or `stdio` |
| `-kubeconfigPath` | `""` | Path to kubeconfig file. If empty, uses InCluster config |
| `-apiKey` | `""` | API key for HTTP Bearer / X-API-Key authentication |
| `-customURL` | `""` | gRPC server URL for custom resource tool |

## Security

- Access is fully controlled by the RBAC permissions defined in the kubeconfig.
- Only operations allowed by the kubeconfig are executed.
- Optional API key authentication for HTTP mode.
