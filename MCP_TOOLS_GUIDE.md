# Kubernetes MCP Server - 工具使用手册

该项目是一个基于 [MCP (Model Context Protocol)](https://github.com/mark3labs/mcp-go) 的 Kubernetes 管理服务，提供 **100+ 个工具** 用于操作 Kubernetes 集群资源。支持 HTTP 和 stdio 两种运行模式。

## 架构概览

- [`main.go`](main.go:1) — 入口，注册所有工具到 MCP Server
- [`tools/tools.go`](tools/tools.go:1) — 定义所有工具的 Schema（名称、描述、参数）
- [`kubernetes/client/client.go`](kubernetes/client/client.go:1) — 初始化 Kubernetes 客户端（支持 InCluster 和 kubeconfig）
- `kubernetes/<resource>/` — 各资源的具体实现
- [`proto/`](proto/custom_tool.proto:1) — gRPC 自定义工具协议

---

## 通用参数

所有 `list-*` 和 `get-*` 工具均支持以下可选参数：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `output` | string | 否 | 输出格式：空=平铺摘要，`json`=完整 K8s 对象 JSON，`yaml`=完整 K8s 对象 YAML |

---

## 一、Pod 工具

### 1. `list-pod`
列出 Pod，支持按命名空间和标签过滤。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有命名空间 |
| `label` | string | 否 | 标签选择器，如 `app=nginx,env=prod` |

**输出**：JSON 数组，包含 `name`、`namespace`、`status`、`labels`

### 2. `get-pod`
获取指定 Pod 的详细信息。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | Pod 所在命名空间 |
| `name` | string | 是 | Pod 名称 |

**输出**：JSON，包含 `name`、`namespace`、`status`、`labels`、`containerNames`

### 3. `delete-pod`
删除指定 Pod。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | Pod 所在命名空间 |
| `name` | string | 是 | Pod 名称 |

### 4. `update-pod`
更新 Pod 的标签（追加到现有标签，不覆盖）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | Pod 所在命名空间 |
| `name` | string | 是 | Pod 名称 |
| `label` | string | 是 | 标签，格式 `key=value,key2=value2` |

### 5. `create-pod`
创建 Pod。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 目标命名空间 |
| `name` | string | 是 | Pod 名称 |
| `label` | string | 否 | 标签，格式 `key=value,key2=value2` |
| `containerNames` | string | 是 | 容器名，多个用逗号分隔 |
| `containerImages` | string | 是 | 容器镜像，多个用逗号分隔 |
| `containerPorts` | string | 否 | 端口，格式 `name:port,name2:port2`，默认 `http:8080` |

> 注意：`containerNames` 和 `containerImages` 数量必须一致。端口支持 `|` 分隔多个端口定义。

### 6. `pod-log`
获取 Pod 日志。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | Pod 所在命名空间 |
| `name` | string | 是 | Pod 名称 |
| `containerName` | string | 否 | 容器名称（单容器 Pod 可不填） |
| `tailLine` | number | 否 | 返回日志行数，默认 100 |

---

## 二、Namespace 工具

### 7. `list-ns`
列出所有命名空间。

**无参数**。输出 JSON 数组，包含 `name`、`status`。

### 8. `get-ns`
获取指定命名空间详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 命名空间名称 |

### 9. `delete-ns`
删除命名空间。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 命名空间名称 |

### 10. `update-ns`
更新命名空间的标签或注解（追加到现有值，不覆盖）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 命名空间名称 |
| `label` | string | 否 | 标签，格式 `key=value` |
| `annotation` | string | 否 | 注解，格式 `key=value` |

> `label` 和 `annotation` 二选一，不能同时使用。

### 11. `create-ns`
创建命名空间。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 命名空间名称 |
| `label` | string | 否 | 标签，格式 `key=value` |

---

## 三、Deployment 工具

### 12. `list-deployment`
列出 Deployment。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |
| `label` | string | 否 | 标签选择器 |

**输出**：JSON 数组，包含 `name`、`namespace`、`availableInstance`（Ready/Total）、`labels`

### 13. `get-deployment`
获取 Deployment 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Deployment 名称 |

**输出**：JSON，包含 `containerName`、`containerImage` 等

### 14. `delete-deployment`
删除 Deployment。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Deployment 名称 |

### 15. `update-deployment`
更新 Deployment 的标签、注解、副本数或镜像。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Deployment 名称 |
| `label` | string | 否 | 标签（追加，不覆盖） |
| `annotation` | string | 否 | 注解（追加，不覆盖） |
| `replica` | number | 否 | 副本数 |
| `containerName` | string | 否 | 容器名（多容器时更新镜像必填） |
| `image` | string | 否 | 新镜像 |

> 更新操作按优先级：`label` > `annotation` > `image` > `replica`，一次只能更新一种。

### 16. `create-deployment`
创建 Deployment。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Deployment 名称 |
| `label` | string | 否 | 标签，默认 `app=<name>` |
| `replica` | number | 否 | 副本数，默认 1 |
| `containerNames` | string | 是 | 容器名，逗号分隔 |
| `containerImages` | string | 是 | 镜像，逗号分隔 |
| `containerPorts` | string | 否 | 端口，默认 `http:8080` |

---

## 四、Service 工具

### 17. `list-service`
列出 Service。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

**输出**：JSON，包含 `name`、`namespace`、`type`

### 18. `get-service`
获取 Service 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Service 名称 |

**输出**：JSON，包含 `type`、`internalIP`、`externalIP`（LoadBalancer 类型）、`selectorLabel`

### 19. `delete-service`
删除 Service。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Service 名称 |

### 20. `update-service`
更新 Service 的选择器标签或类型。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Service 名称 |
| `selectorLabel` | string | 否 | 选择器标签 |
| `svctype` | string | 否 | 服务类型，如 `ClusterIP`、`NodePort`、`LoadBalancer` |

### 21. `create-service`
创建 Service。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Service 名称 |
| `selectorLabel` | string | 是 | 选择器标签 |
| `svcPort` | string | 是 | 服务端口，支持 `name:port` 或 `port` 格式，多端口用逗号分隔 |
| `targetPort` | string | 是 | 目标端口，多端口用逗号分隔 |
| `svcType` | string | 否 | 服务类型，默认 `ClusterIP` |

> 支持多端口，`svcPort` 和 `targetPort` 用逗号分隔，数量必须一致。

---

## 五、StatefulSet 工具

### 22. `list-statefulset`
列出 StatefulSet。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |
| `label` | string | 否 | 标签选择器 |

### 23. `get-statefulset`
获取 StatefulSet 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | StatefulSet 名称 |

### 24. `delete-statefulset`
删除 StatefulSet。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | StatefulSet 名称 |

### 25. `update-statefulset`
更新 StatefulSet（与 Deployment 更新逻辑相同）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | StatefulSet 名称 |
| `label` | string | 否 | 标签（追加，不覆盖） |
| `annotation` | string | 否 | 注解（追加，不覆盖） |
| `replica` | number | 否 | 副本数 |
| `containerName` | string | 否 | 容器名 |
| `image` | string | 否 | 新镜像 |

### 26. `create-statefulset`
创建 StatefulSet（**同时创建关联的 Service 和 PVC**）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | StatefulSet 名称 |
| `label` | string | 否 | 标签，默认 `app=<name>` |
| `containerNames` | string | 否 | 容器名，默认等于 `name` |
| `containerImages` | string | 是 | 容器镜像 |
| `containerPorts` | number | 否 | 容器端口，默认 8080 |
| `storageValue` | string | 是 | PVC 大小，如 `1Gi` |
| `mountPath` | string | 是 | 挂载路径 |
| `pvcName` | string | 否 | PVC 名称，默认等于 `name` |
| `svcType` | string | 否 | Service 类型，默认 `ClusterIP` |
| `svcPort` | number | 否 | Service 端口，默认 8080 |
| `replica` | number | 否 | 副本数，默认 1 |

---

## 六、DaemonSet 工具

### 27. `list-daemonset`
列出 DaemonSet。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |
| `label` | string | 否 | 标签选择器 |

### 28. `get-daemonset`
获取 DaemonSet 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | DaemonSet 名称 |

### 29. `delete-daemonset`
删除 DaemonSet。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | DaemonSet 名称 |

### 30. `update-daemonset`
更新 DaemonSet（支持标签、注解、镜像，**不支持副本数**）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | DaemonSet 名称 |
| `label` | string | 否 | 标签（追加，不覆盖） |
| `annotation` | string | 否 | 注解（追加，不覆盖） |
| `containerName` | string | 否 | 容器名 |
| `image` | string | 否 | 新镜像 |

### 31. `create-daemonset`
创建 DaemonSet。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | DaemonSet 名称 |
| `label` | string | 否 | 标签，默认 `app=<name>` |
| `containerNames` | string | 是 | 容器名，逗号分隔 |
| `containerImages` | string | 是 | 镜像，逗号分隔 |
| `containerPorts` | string | 否 | 端口，默认 `http:8080` |

---

## 七、ConfigMap 工具

### 32. `list-configmap`
列出 ConfigMap。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

### 33. `get-configmap`
获取 ConfigMap 详情（含数据）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ConfigMap 名称 |

### 34. `delete-configmap`
删除 ConfigMap。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ConfigMap 名称 |

### 35. `create-configmap`
创建 ConfigMap。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ConfigMap 名称 |
| `data` | string | 是 | 数据，格式 `key=value,key2=value2` |

---

## 八、Secret 工具

### 36. `list-secret`
列出 Secret。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

### 37. `get-secret`
获取 Secret 详情。**注意**：返回的数据为 base64 编码的原始值，包含密码、密钥等敏感信息。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Secret 名称 |

### 38. `delete-secret`
删除 Secret。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Secret 名称 |

### 39. `create-secret`
创建 Opaque 类型 Secret。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Secret 名称 |
| `data` | string | 是 | 数据，格式 `key=value,key2=value2` |

---

## 九、Node 工具

### 40. `list-node`
列出所有节点。

**无参数**。输出 JSON 数组，包含 `name`、`status`（Ready/NotReady）、`capacityCPU`、`capacityMemory`、`capacityPods`、`allocatableCPU`、`allocatableMemory`、`allocatablePods`、`labels`、`taints`。

### 41. `get-node`
获取节点详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 节点名称 |

**输出**：JSON，包含 `status`、`kubernetesVersion`、`os`、`kernelVersion`、`architecture`、`podCIDR`、`capacityCPU`、`capacityMemory`、`capacityPods`、`allocatableCPU`、`allocatableMemory`、`allocatablePods`、`labels`、`taints`。

### 42. `delete-node`
删除节点。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 节点名称 |

### 43. `update-node`
更新节点标签（追加到现有标签，不覆盖）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 节点名称 |
| `label` | string | 是 | 标签，格式 `key=value` |

---

## 十、ServiceAccount 工具

### 44. `list-serviceAccount`
列出 ServiceAccount。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |
| `label` | string | 否 | 标签选择器 |

### 45. `get-serviceAccount`
获取 ServiceAccount 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ServiceAccount 名称 |

### 46. `delete-serviceAccount`
删除 ServiceAccount。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ServiceAccount 名称 |

### 47. `create-serviceAccount`
创建 ServiceAccount。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ServiceAccount 名称 |
| `label` | string | 否 | 标签 |

---

## 十一、PVC 工具

### 48. `list-pvc`
列出 PVC。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

**输出**：JSON 数组，包含 `name`、`namespace`、`capacity`、`status`

### 49. `get-pvc`
获取 PVC 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | PVC 名称 |

**输出**：JSON，包含 `capacity`、`accessMode`、`storageClass`、`volume`、`status`

### 50. `delete-pvc`
删除 PVC。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | PVC 名称 |

### 51. `update-pvc`
更新 PVC 大小。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | PVC 名称 |
| `size` | string | 是 | 新大小，如 `2Gi` |

### 52. `create-pvc`
创建 PVC。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | PVC 名称 |
| `size` | string | 是 | 大小，如 `1Gi` |
| `storageClass` | string | 否 | StorageClass 名称（为空则使用默认 StorageClass） |
| `accessMode` | string | 否 | 访问模式，默认 `ReadWriteOnce`，支持逗号分隔多个 |

---

## 十二、PV 工具

### 53. `list-pv`
列出所有 PV。

**无参数**。

### 54. `get-pv`
获取 PV 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | PV 名称 |

### 55. `delete-pv`
删除 PV。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | PV 名称 |

---

## 十三、Role 工具

### 56. `list-role`
列出 Role。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

### 57. `get-role`
获取 Role 详情（含规则）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Role 名称 |

### 58. `delete-role`
删除 Role。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Role 名称 |

---

## 十四、RoleBinding 工具

### 59. `list-rolebinding`
列出 RoleBinding。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间 |

### 60. `get-rolebinding`
获取 RoleBinding 详情（含 RoleRef 和 Subjects）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | RoleBinding 名称 |

### 61. `delete-rolebinding`
删除 RoleBinding。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | RoleBinding 名称 |

---

## 十五、ClusterRole 工具

### 62. `list-clusterrole`
列出所有 ClusterRole。

**无参数**。

### 63. `get-clusterrole`
获取 ClusterRole 详情（含规则）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | ClusterRole 名称 |

### 64. `delete-clusterrole`
删除 ClusterRole。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | ClusterRole 名称 |

---

## 十六、ClusterRoleBinding 工具

### 65. `list-clusterrolebinding`
列出所有 ClusterRoleBinding。

**无参数**。

### 66. `get-clusterrolebinding`
获取 ClusterRoleBinding 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | ClusterRoleBinding 名称 |

### 67. `delete-clusterrolebinding`
删除 ClusterRoleBinding。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | ClusterRoleBinding 名称 |

---

## 十七、StorageClass 工具

### 68. `list-storageClass`
列出所有 StorageClass。

**无参数**。

### 69. `get-storageClass`
获取 StorageClass 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | StorageClass 名称 |

**输出**：JSON，包含 `name`、`provisioner`、`reclaimPolicy`

### 70. `delete-storageClass`
删除 StorageClass。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | StorageClass 名称 |

---

## 十八、CRD 工具

### 71. `list-crd`
列出所有 CRD。

**无参数**。

### 72. `get-crd`
获取 CRD 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | CRD 名称 |

### 73. `delete-crd`
删除 CRD。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | CRD 名称 |

### 74. `create-crd-with-json`
通过 JSON/YAML 数据创建 CRD。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `jsondata` | string | 是 | CRD 的 JSON 或 YAML 定义 |

---

## 十九、通用资源创建工具

### 75. `create-resource-with-json`
通过 JSON/YAML 数据创建任意 Kubernetes 资源（使用 Dynamic Client）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `jsondata` | string | 是 | 资源的完整 JSON 或 YAML 定义 |

> 自动识别资源类型（GVR），支持集群作用域和命名空间作用域资源，默认命名空间为 `default`。

---

## 二十、自定义工具（gRPC）

### 76. `custom`
通过 gRPC 调用外部服务操作自定义资源。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `kind` | string | 是 | 自定义资源类型 |
| `method` | string | 是 | 操作方法 |
| `name` | string | 否 | 资源名称 |
| `namespace` | string | 否 | 命名空间 |
| `jsondata` | string | 否 | JSON 数据 |

> 需要启动时通过 `-customURL` 参数指定 gRPC 服务地址。Proto 定义见 [`proto/custom_tool.proto`](proto/custom_tool.proto:1)。

---

## 二十一、Event 工具（巡检）

### 77. `list-event`
列出集群事件，支持按命名空间过滤。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有 |

**输出**：JSON 数组，包含 `name`、`namespace`、`type`（Normal/Warning）、`reason`、`message`、`source`、`firstTime`、`lastTime`、`count`、`kind`、`involved`

### 78. `get-event`
获取指定事件的详细信息。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 事件所在命名空间 |
| `name` | string | 是 | 事件名称 |

---

## 二十二、ResourceQuota 工具（巡检）

### 79. `list-resourcequota`
列出 ResourceQuota，支持按命名空间过滤。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有 |

**输出**：JSON 数组，包含 `name`、`namespace`、`hard`（配额上限）、`used`（当前使用量）

### 80. `get-resourcequota`
获取 ResourceQuota 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | ResourceQuota 名称 |

---

## 二十三、LimitRange 工具（巡检）

### 81. `list-limitrange`
列出 LimitRange，支持按命名空间过滤。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有 |

**输出**：JSON 数组，包含 `name`、`namespace`、`limits`（含 max/min/default CPU/Memory）

### 82. `get-limitrange`
获取 LimitRange 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | LimitRange 名称 |

---

## 二十四、Endpoint 工具（巡检）

### 83. `list-endpoint`
列出 Endpoints，支持按命名空间过滤。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有 |

**输出**：JSON 数组，包含 `name`、`namespace`、`addresses`（含 IP、NodeName、Ports）

### 84. `get-endpoint`
获取 Endpoint 详情（含 IP 和端口信息）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Endpoint 名称 |

---

## 二十五、ComponentStatus 工具（巡检）

### 85. `list-componentstatus`
列出所有 Kubernetes 控制面组件健康状态。

**无参数**。输出 JSON 数组，包含 `name`（如 etcd-0、kube-apiserver）、`status`（True/False）、`type`

### 86. `get-componentstatus`
获取指定组件的健康状态。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 组件名称（如 `etcd-0`、`kube-apiserver`、`kube-scheduler`） |

---

## 二十六、Top 工具（监控，需 metrics-server）

### 87. `top-pod`
显示 Pod 的 CPU 和内存使用量（需要集群中部署了 metrics-server）。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 指定命名空间，为空则列出所有 |

**输出**：JSON 数组，包含 `name`、`namespace`、`cpu`（如 `100m`）、`memory`（如 `128Mi`）

### 88. `top-node`
显示 Node 的 CPU 和内存使用量（需要集群中部署了 metrics-server）。

**无参数**。输出 JSON 数组，包含 `name`、`cpu`、`memory`

---

## 二十七、Cluster Health 工具（监控）

### 89. `cluster-health`
显示集群整体健康概览：节点汇总（Total/Ready/NotReady）+ 控制面组件 Pod 状态。

**无参数**。输出 JSON，包含：
- `nodes.total`、`nodes.ready`、`nodes.notReady`
- `controlPlanePods[].name`、`controlPlanePods[].namespace`、`controlPlanePods[].component`、`controlPlanePods[].status`、`controlPlanePods[].restarts`、`controlPlanePods[].node`

### 90. `node-health`
显示所有节点的详细健康状态。

**无参数**。输出 JSON 数组，每个节点包含 `name`、`status`（Ready/NotReady）、`kubelet`、`cpu`、`memory`、`pods`、`labels`

---

## 二十八、Deep Diagnosis 工具（集群故障诊断）

### 91. `describe-pod`
深度检查 Pod，对标 `kubectl describe pod`。输出容器级状态、Conditions、QoS、资源请求/限制、关联 Events。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | Pod 所在命名空间 |
| `name` | string | 是 | Pod 名称 |

**输出**：纯文本格式，包含：
- **基本信息**：Name、Namespace、Node、Start Time、Status、IP、QoS Class
- **Labels**
- **Conditions**：PodScheduled、Initialized、ContainersReady、Ready 状态
- **Containers**：每个容器的 Image、Command、Ports、Resource Requests/Limits、Ready、Restart Count、State（Waiting/Running/Terminated 含原因和退出码）、Last State
- **Init Containers**：状态详情
- **Events**：按时间倒序，显示 LastTimestamp、Type、Reason、Message、Count

### 92. `describe-node`
深度检查节点，对标 `kubectl describe node`。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 节点名称 |

**输出**：纯文本格式，包含：
- **基本信息**：Creation、Kubelet、OS Image、Kernel、Architecture、PodCIDR、ProviderID
- **Labels**、**Annotations**、**Taints**
- **Conditions**：Ready/DiskPressure/MemoryPressure/PIDPressure/NetworkUnavailable 状态、Reason、Message、Last Heartbeat
- **Capacity**：cpu、memory、pods、ephemeral-storage
- **Allocatable**：同维度可分配量
- **Pods**：节点上运行的所有 Pod 列表（ns/name/status）

### 93. `list-node-pods`
列出指定节点上运行的所有 Pod。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node` | string | 是 | 节点名称 |

**输出**：JSON 数组，包含 `namespace`、`name`、`status`、`node`

### 94. `describe-service`
深度检查 Service，对标 `kubectl describe service`。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Service 名称 |

**输出**：纯文本格式，包含：
- **基本信息**：Type、ClusterIP、ExternalIPs、ExternalName、LoadBalancerIP、Session Affinity
- **Labels**、**Selector**
- **Ports**：name、port/protocol → targetPort
- **Endpoints**：每个 subset 的地址（含 NodeName）
- **Events**

### 95. `describe-deployment`
深度检查 Deployment，对标 `kubectl describe deployment`。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 是 | 命名空间 |
| `name` | string | 是 | Deployment 名称 |

**输出**：纯文本格式，包含：
- **基本信息**：Strategy（RollingUpdate 的 MaxSurge/MaxUnavailable）、Replicas（desired/updated/total/available/unavailable）、Revision History Limit、Min Ready Seconds、Rollout Status
- **Labels**、**Selector**
- **Conditions**：Available/Progressing/ReplicaFailure 状态、Reason、Message、LastUpdateTime
- **Containers**：Image、Resource Requests/Limits
- **Pods**：按 Phase 统计（Running/Pending/Failed 等）
- **Events**

### 111. `check-apiserver-health`
探测 API Server 的健康端点。**不依赖 Pod 标签或 ComponentStatus API，二进制部署的 apiserver 同样可用。**

**无参数**。直接调用 apiserver 的三个端点：
- `/healthz?verbose` — 详细健康检查结果
- `/livez?verbose` — 存活探活
- `/readyz?verbose` — 就绪探活

**输出**：纯文本，每个端点的 HTTP 状态码和所有非 `ok` 的检查项。

### 111. `check-apiserver-metrics`
拉取 API Server 的 `/metrics` 并分析关键性能指标。**二进制部署的 apiserver 同样可用。**

**无参数**。需要 apiserver 配置 `--authorization-always-allow-paths=/metrics` 或在 RBAC 中放权。

**输出**：包含：
- **Current Inflight Requests**：mutating 和 readOnly 当前处理中的请求数
- **Request Counts**：按 Verb（GET/LIST/WATCH/POST/PUT/PATCH/DELETE）统计请求总量、错误量、错误率
- **Request Latency**：按 Verb 估算 p50/p90/p99 延迟
- **Top Error Endpoints**：4xx/5xx 错误按 code 排序

---

## 二十九、Helm 工具

### 111. `helm-list-releases`
列出 Helm Releases。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `namespace` | string | 否 | 过滤命名空间 |
| `allNamespaces` | boolean | 否 | 列出所有命名空间的 release |
| `output` | string | 否 | 输出格式：`json` 或 `yaml` 获取完整对象 |

### 111. `helm-get-release`
获取 Release 详情。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |
| `output` | string | 否 | 输出格式 |

### 111. `helm-get-values`
获取 Release 的 Values。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |

### 111. `helm-install`
安装 Chart。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | 目标命名空间 |
| `chart` | string | 是 | Chart 引用，如 `stable/nginx-ingress`，或本地路径 |
| `version` | string | 否 | Chart 版本 |
| `values` | string | 否 | 内联 YAML/JSON 格式的 values |
| `set` | string | 否 | 命令行设置 values，格式 `key1=val1,key2=val2` |

> `values` 和 `set` 可以同时使用，`set` 优先级更高。

### 111. `helm-upgrade`
升级 Release。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |
| `chart` | string | 是 | 新 Chart 引用 |
| `version` | string | 否 | Chart 版本 |
| `values` | string | 否 | 内联 YAML/JSON 格式的 values |
| `set` | string | 否 | 命令行设置 values |
| `reuseValues` | boolean | 否 | 是否复用之前的 values，默认 `true` |

### 111. `helm-uninstall`
卸载 Release。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |

### 111. `helm-rollback`
回滚 Release 到指定版本。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |
| `revision` | number | 否 | 目标版本号，默认回滚到上一个版本 |

### 111. `helm-history`
获取 Release 的修订历史。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |
| `max` | number | 否 | 最大显示版本数 |

**输出**：JSON 数组，每个版本包含 `revision`、`status`、`chart`、`appVersion`、`description`、`updated`

### 111. `helm-get-manifest`
获取 Release 的渲染后 Manifest。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |

### 111. `helm-get-notes`
获取 Release 的 NOTES.txt 内容。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Release 名称 |
| `namespace` | string | 是 | Release 所在命名空间 |

### 111. `helm-list-repos`
列出已添加的 Helm 仓库。

**无参数**。

### 111. `helm-add-repo`
添加 Helm 仓库并下载索引。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 仓库名称 |
| `url` | string | 是 | 仓库 URL |

### 111. `helm-remove-repo`
移除 Helm 仓库。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 仓库名称 |

### 111. `helm-update-repos`
更新 Helm 仓库索引文件。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 否 | 仓库名称，为空则更新所有仓库 |

### 111. `helm-search-repo`
在所有已添加的仓库中搜索 Chart。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 搜索关键词（匹配 chart 名称、描述、appVersion） |

---

## 启动方式

```bash
# HTTP 模式（默认，监听 :8080）
go run main.go

# 指定 kubeconfig
go run main.go -kubeconfigPath /path/to/kubeconfig

# stdio 模式
go run main.go -mode stdio

# 指定 gRPC 自定义服务地址
go run main.go -customURL localhost:50051

# 启用 API 密钥认证
go run main.go -apiKey my-secret-key
```

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-mode` | `http` | 运行模式：`http` 或 `stdio` |
| `-kubeconfigPath` | `""` | kubeconfig 文件路径，为空则使用 InCluster 模式 |
| `-apiKey` | `""` | HTTP API 密钥认证（可选） |
| `-customURL` | `""` | gRPC 自定义资源服务地址 |

## 客户端初始化逻辑

[`kubernetes/client/client.go`](kubernetes/client/client.go:19) 中的 [`InitializeClients()`](kubernetes/client/client.go:19) 返回 6 个客户端：

1. `*kubernetes.Clientset` — 标准 Kubernetes 客户端
2. `dynamic.Interface` — 动态客户端（用于通用资源创建）
3. `discovery.DiscoveryInterface` — API 发现客户端
4. `*apiextensionsclient.Clientset` — CRD 客户端
5. `*metricsv.Interface` — Metrics 客户端（用于 Top 工具）

连接优先级：`-kubeconfigPath` 参数 > InCluster 模式（Pod 内运行）。

## 工具注册映射表

[`main.go`](main.go:38) 中注册了所有工具，以下是工具名与实现函数的对应关系：

| 工具名 | Schema 变量 | 实现函数 |
|--------|------------|----------|
| `list-pod` | `tools.ListPod` | `pod.ListPod` |
| `get-pod` | `tools.GetPod` | `pod.GetPod` |
| `delete-pod` | `tools.DeletePod` | `pod.DeletePod` |
| `update-pod` | `tools.UpdatePod` | `pod.UpdatePod` |
| `create-pod` | `tools.CreatePod` | `pod.CreatePod` |
| `pod-log` | `tools.PodLog` | `pod.PodLog` |
| `list-ns` | `tools.ListNS` | `namespace.ListNS` |
| `get-ns` | `tools.GetNS` | `namespace.GetNS` |
| `delete-ns` | `tools.DeleteNS` | `namespace.DeleteNS` |
| `update-ns` | `tools.UpdateNS` | `namespace.UpdateNS` |
| `create-ns` | `tools.CreateNS` | `namespace.CreateNS` |
| `list-node` | `tools.ListNode` | `node.ListNode` |
| `get-node` | `tools.GetNode` | `node.GetNode` |
| `delete-node` | `tools.DeleteNode` | `node.DeleteNode` |
| `update-node` | `tools.UpdateNode` | `node.UpdateNode` |
| `list-deployment` | `tools.ListDeployment` | `deployment.ListDeployment` |
| `get-deployment` | `tools.GetDeployment` | `deployment.GetDeployment` |
| `delete-deployment` | `tools.DeleteDeployment` | `deployment.DeleteDeployment` |
| `create-deployment` | `tools.CreateDeployment` | `deployment.CreateDeployment` |
| `update-deployment` | `tools.UpdateDeployment` | `deployment.UpdateDeployment` |
| `list-daemonset` | `tools.ListDaemonset` | `daemonset.ListDaemonset` |
| `get-daemonset` | `tools.GetDaemonset` | `daemonset.GetDaemonset` |
| `delete-daemonset` | `tools.DeleteDaemonset` | `daemonset.DeleteDaemonset` |
| `update-daemonset` | `tools.UpdateDaemonset` | `daemonset.UpdateDaemonset` |
| `create-daemonset` | `tools.CreateDaemonset` | `daemonset.CreateDaemonset` |
| `list-statefulset` | `tools.ListStatefulset` | `statefulset.ListStatefulset` |
| `get-statefulset` | `tools.GetStatefulset` | `statefulset.GetStatefulset` |
| `delete-statefulset` | `tools.DeleteStatefulset` | `statefulset.DeleteStatefulset` |
| `update-statefulset` | `tools.UpdateStatefulset` | `statefulset.UpdateStatefulset` |
| `create-statefulset` | `tools.CreateStatefulset` | `statefulset.CreateStatefulset` |
| `list-service` | `tools.ListService` | `service.ListService` |
| `get-service` | `tools.GetService` | `service.GetService` |
| `delete-service` | `tools.DeleteService` | `service.DeleteService` |
| `update-service` | `tools.UpdateService` | `service.UpdateService` |
| `create-service` | `tools.CreateService` | `service.CreateService` |
| `list-configmap` | `tools.ListConfigmap` | `configmap.ListConfigmap` |
| `get-configmap` | `tools.GetConfigmap` | `configmap.GetConfigmap` |
| `delete-configmap` | `tools.DeleteConfigmap` | `configmap.DeleteConfigmap` |
| `create-configmap` | `tools.CreateConfigmap` | `configmap.CreateConfigmap` |
| `list-secret` | `tools.ListSecret` | `secret.ListSecret` |
| `get-secret` | `tools.GetSecret` | `secret.GetSecret` |
| `delete-secret` | `tools.DeleteSecret` | `secret.DeleteSecret` |
| `create-secret` | `tools.CreateSecret` | `secret.CreateSecret` |
| `list-serviceAccount` | `tools.ListSA` | `serviceaccount.ListSA` |
| `get-serviceAccount` | `tools.GetSA` | `serviceaccount.GetSA` |
| `delete-serviceAccount` | `tools.DeleteSA` | `serviceaccount.DeleteSA` |
| `create-serviceAccount` | `tools.CreateSA` | `serviceaccount.CreateSA` |
| `list-role` | `tools.ListRole` | `role.ListRole` |
| `get-role` | `tools.GetRole` | `role.GetRole` |
| `delete-role` | `tools.DeleteRole` | `role.DeleteRole` |
| `list-rolebinding` | `tools.ListRB` | `rolebinding.ListRB` |
| `get-rolebinding` | `tools.GetRB` | `rolebinding.GetRB` |
| `delete-rolebinding` | `tools.DeleteRB` | `rolebinding.DeleteRB` |
| `list-pvc` | `tools.ListPVC` | `pvc.ListPVC` |
| `get-pvc` | `tools.GetPVC` | `pvc.GetPVC` |
| `delete-pvc` | `tools.DeletePVC` | `pvc.DeletePVC` |
| `update-pvc` | `tools.UpdatePVC` | `pvc.UpdatePVC` |
| `create-pvc` | `tools.CreatePVC` | `pvc.CreatePVC` |
| `list-pv` | `tools.ListPV` | `pv.ListPV` |
| `get-pv` | `tools.GetPV` | `pv.GetPV` |
| `delete-pv` | `tools.DeletePV` | `pv.DeletePV` |
| `list-clusterrole` | `tools.ListCR` | `clusterrole.ListCR` |
| `get-clusterrole` | `tools.GetCR` | `clusterrole.GetCR` |
| `delete-clusterrole` | `tools.DeleteCR` | `clusterrole.DeleteCR` |
| `list-clusterrolebinding` | `tools.ListCRB` | `clusterrolebinding.ListCRB` |
| `get-clusterrolebinding` | `tools.GetCRB` | `clusterrolebinding.GetCRB` |
| `delete-clusterrolebinding` | `tools.DeleteCRB` | `clusterrolebinding.DeleteCRB` |
| `list-storageClass` | `tools.ListSC` | `storageclass.ListSC` |
| `get-storageClass` | `tools.GetSC` | `storageclass.GetSC` |
| `delete-storageClass` | `tools.DeleteSC` | `storageclass.DeleteSC` |
| `list-crd` | `tools.ListCRD` | `crd.ListCRD` |
| `get-crd` | `tools.GetCRD` | `crd.GetCRD` |
| `delete-crd` | `tools.DeleteCRD` | `crd.DeleteCRD` |
| `create-crd-with-json` | `tools.CreateCRDWithJson` | `crd.CreateCRDWithJson` |
| `create-resource-with-json` | `tools.CreateResourceWithJSon` | `createresource.CreateResourceWithJson` |
| `custom` | `tools.Custom` | `custom.Custom` |
| `list-event` | `tools.ListEvent` | `event.ListEvent` |
| `get-event` | `tools.GetEvent` | `event.GetEvent` |
| `list-resourcequota` | `tools.ListResourceQuota` | `resourcequota.ListResourceQuota` |
| `get-resourcequota` | `tools.GetResourceQuota` | `resourcequota.GetResourceQuota` |
| `list-limitrange` | `tools.ListLimitRange` | `limitrange.ListLimitRange` |
| `get-limitrange` | `tools.GetLimitRange` | `limitrange.GetLimitRange` |
| `list-endpoint` | `tools.ListEndpoint` | `endpoint.ListEndpoint` |
| `get-endpoint` | `tools.GetEndpoint` | `endpoint.GetEndpoint` |
| `list-componentstatus` | `tools.ListComponentStatus` | `componentstatus.ListComponentStatus` |
| `get-componentstatus` | `tools.GetComponentStatus` | `componentstatus.GetComponentStatus` |
| `top-pod` | `tools.TopPod` | `top.TopPod` |
| `top-node` | `tools.TopNode` | `top.TopNode` |
| `cluster-health` | `tools.GetClusterHealth` | `clusterhealth.GetClusterHealth` |
| `node-health` | `tools.ListNodeHealth` | `clusterhealth.ListNodeHealth` |
| `describe-pod` | `tools.DescribePod` | `diagnose.DescribePod` |
| `describe-node` | `tools.DescribeNode` | `diagnose.DescribeNode` |
| `list-node-pods` | `tools.ListNodePods` | `diagnose.ListNodePods` |
| `describe-service` | `tools.DescribeService` | `diagnose.DescribeService` |
| `describe-deployment` | `tools.DescribeDeployment` | `diagnose.DescribeDeployment` |
| `check-apiserver-health` | `tools.CheckAPIServerHealth` | `diagnose.CheckAPIServerHealth` |
| `check-apiserver-metrics` | `tools.CheckAPIServerMetrics` | `diagnose.CheckAPIServerMetrics` |
| `helm-list-releases` | `tools.ListHelmReleases` | `helm.ListHelmReleases` |
| `helm-get-release` | `tools.GetHelmRelease` | `helm.GetHelmRelease` |
| `helm-get-values` | `tools.GetHelmReleaseValues` | `helm.GetHelmReleaseValues` |
| `helm-install` | `tools.InstallHelmRelease` | `helm.InstallHelmRelease` |
| `helm-upgrade` | `tools.UpgradeHelmRelease` | `helm.UpgradeHelmRelease` |
| `helm-uninstall` | `tools.UninstallHelmRelease` | `helm.UninstallHelmRelease` |
| `helm-rollback` | `tools.RollbackHelmRelease` | `helm.RollbackHelmRelease` |
| `helm-history` | `tools.GetHelmReleaseHistory` | `helm.GetHelmReleaseHistory` |
| `helm-get-manifest` | `tools.GetHelmReleaseManifest` | `helm.GetHelmReleaseManifest` |
| `helm-get-notes` | `tools.GetHelmReleaseNotes` | `helm.GetHelmReleaseNotes` |
| `helm-list-repos` | `tools.ListHelmRepos` | `helm.ListHelmRepos` |
| `helm-add-repo` | `tools.AddHelmRepo` | `helm.AddHelmRepo` |
| `helm-remove-repo` | `tools.RemoveHelmRepo` | `helm.RemoveHelmRepo` |
| `helm-update-repos` | `tools.UpdateHelmRepos` | `helm.UpdateHelmRepos` |
| `helm-search-repo` | `tools.SearchHelmRepo` | `helm.SearchHelmRepo` |

---
