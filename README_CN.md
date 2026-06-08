# K8S MCP SERVER

k8s-mcp 是一个基于 Golang 的模型上下文协议（MCP）服务器，它将 Kubernetes 资源暴露为结构化的 MCP 工具，使 AI 代理能够安全地与 Kubernetes 集群进行交互。

## Features

### 支持的 Kubernetes 资源和操作（100+ tools）：

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
- Helm: List/Get/Install/Upgrade/Uninstall/Rollback releases, Get history/manifest/notes/values, Manage repositories
- Cluster Diagnosis: describe-pod, describe-node, describe-deployment, describe-service, list-node-pods, check-apiserver-health, check-apiserver-metrics

Create any kubernetes resource by passing json data.

所有 `list-*` 和 `get-*` 工具都支持可选的 `output` 参数（`output=json` 或 `output=yaml`），用于返回完整的 K8s 对象而非摘要格式。

### Custom Resource:

除了上述预定义的 Kubernetes 资源工具外，本 MCP 服务器还提供了一个通用的自定义工具，允许用户与任何未显式支持的 Kubernetes 自定义资源进行交互。

这在以下场景中非常有用：

- 你想处理自定义资源（CRDs）。
- 你想访问新的 Kubernetes 资源而无需更新 MCP 服务器。

MCP 服务器会将提供的参数转发给 gRPC 后端服务器，该服务器会动态解析资源并执行请求的操作。

自定义工具支持的操作：Create, Get, List and Delete。

注意：每个资源的参数详情可在 kubernetes 目录下相应的 `README.md` 文件中查看。

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

通过使用 `--customURL` 标志启用自定义工具：

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

### InCluster 模式（Kubernetes Pod 内运行）：

```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "k8s-mcp",
      "args": ["-mode", "stdio"]
    }
  }
}
```
不传 `-kubeconfigPath` 时自动使用 Pod 的 ServiceAccount 通过 `rest.InClusterConfig()` 连接集群。需确保 RBAC 权限已通过 ClusterRole/ServiceAccount 授予。

### HTTP 模式（远程 / Cursor）：

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
| `-mode` | `http` | 服务器模式：`http` 或 `stdio` |
| `-kubeconfigPath` | `""` | kubeconfig 文件路径。如果为空，则使用 InCluster 配置 |
| `-apiKey` | `""` | HTTP Bearer / X-API-Key 认证的 API 密钥 |
| `-customURL` | `""` | 自定义资源工具的 gRPC 服务器 URL |

## Security

- 访问权限完全由 kubeconfig 中定义的 RBAC 权限控制。
- 仅执行 kubeconfig 允许的操作。
- HTTP 模式支持可选的 API 密钥认证。
