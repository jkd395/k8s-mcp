# Kubernetes MCP Server - Tool Usage Guide

This project is a Kubernetes management service based on [MCP (Model Context Protocol)](https://github.com/mark3labs/mcp-go), providing **100+ tools** for operating Kubernetes cluster resources. Supports HTTP and stdio runtime modes.

## Architecture Overview

- [`main.go`](main.go:1) — Entry point, registers all tools to the MCP Server
- [`tools/tools.go`](tools/tools.go:1) — Defines schemas for all tools (name, description, parameters)
- [`kubernetes/client/client.go`](kubernetes/client/client.go:1) — Initializes the Kubernetes client (supports InCluster and kubeconfig)
- `kubernetes/<resource>/` — Specific implementations for each resource
- [`proto/`](proto/custom_tool.proto:1) — gRPC custom tool protocol

---

## Common Parameters

All `list-*` and `get-*` tools support the following optional parameters:

| Parameter | Type   | Required | Description                                                       |
|-----------|--------|----------|-------------------------------------------------------------------|
| `output`  | string | no       | Output format: empty=flattened summary, `json`=full K8s object JSON, `yaml`=full K8s object YAML |

---

## 1. Pod Tools

### 1. `list-pod`
List Pods, supports filtering by namespace and label.

| Parameter   | Type   | Required | Description                                              |
|-------------|--------|----------|----------------------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all namespaces            |
| `label`     | string | no       | Label selector, e.g. `app=nginx,env=prod`                |

**Output**: JSON array containing `name`, `namespace`, `status`, `labels`

### 2. `get-pod`
Get detailed information for a specific Pod.

| Parameter   | Type   | Required | Description                |
|-------------|--------|----------|----------------------------|
| `namespace` | string | yes      | Pod namespace              |
| `name`      | string | yes      | Pod name                   |

**Output**: JSON containing `name`, `namespace`, `status`, `labels`, `containerNames`

### 3. `delete-pod`
Delete a specific Pod.

| Parameter   | Type   | Required | Description     |
|-------------|--------|----------|-----------------|
| `namespace` | string | yes      | Pod namespace   |
| `name`      | string | yes      | Pod name        |

### 4. `update-pod`
Update Pod labels (appended to existing labels, not overwritten).

| Parameter   | Type   | Required | Description                                   |
|-------------|--------|----------|-----------------------------------------------|
| `namespace` | string | yes      | Pod namespace                                 |
| `name`      | string | yes      | Pod name                                      |
| `label`     | string | yes      | Labels, format `key=value,key2=value2`        |

### 5. `create-pod`
Create a Pod.

| Parameter        | Type   | Required | Description                                                       |
|------------------|--------|----------|-------------------------------------------------------------------|
| `namespace`      | string | yes      | Target namespace                                                  |
| `name`           | string | yes      | Pod name                                                          |
| `label`          | string | no       | Labels, format `key=value,key2=value2`                            |
| `containerNames` | string | yes      | Container names, comma-separated                                  |
| `containerImages`| string | yes      | Container images, comma-separated                                 |
| `containerPorts` | string | no       | Ports, format `name:port,name2:port2`, default `http:8080`        |

> Note: The number of `containerNames` and `containerImages` must match. Ports support `|` to separate multiple port definitions.

### 6. `pod-log`
Get Pod logs.

| Parameter       | Type   | Required | Description                                      |
|-----------------|--------|----------|--------------------------------------------------|
| `namespace`     | string | yes      | Pod namespace                                    |
| `name`          | string | yes      | Pod name                                         |
| `containerName` | string | no       | Container name (optional for single-container Pod)|
| `tailLine`      | number | no       | Number of log lines to return, default 100       |

---

## 2. Namespace Tools

### 7. `list-ns`
List all namespaces.

**No parameters**. Output JSON array containing `name`, `status`.

### 8. `get-ns`
Get details of a specific namespace.

| Parameter | Type   | Required | Description       |
|-----------|--------|----------|-------------------|
| `name`    | string | yes      | Namespace name    |

### 9. `delete-ns`
Delete a namespace.

| Parameter | Type   | Required | Description    |
|-----------|--------|----------|----------------|
| `name`    | string | yes      | Namespace name |

### 10. `update-ns`
Update namespace labels or annotations (appended to existing values, not overwritten).

| Parameter    | Type   | Required | Description                         |
|--------------|--------|----------|-------------------------------------|
| `name`       | string | yes      | Namespace name                      |
| `label`      | string | no       | Labels, format `key=value`          |
| `annotation` | string | no       | Annotations, format `key=value`     |

> `label` and `annotation` are mutually exclusive, cannot be used together.

### 11. `create-ns`
Create a namespace.

| Parameter | Type   | Required | Description                  |
|-----------|--------|----------|------------------------------|
| `name`    | string | yes      | Namespace name               |
| `label`   | string | no       | Labels, format `key=value`   |

---

## 3. Deployment Tools

### 12. `list-deployment`
List Deployments.

| Parameter   | Type   | Required | Description              |
|-------------|--------|----------|--------------------------|
| `namespace` | string | no       | Specify namespace        |
| `label`     | string | no       | Label selector           |

**Output**: JSON array containing `name`, `namespace`, `availableInstance` (Ready/Total), `labels`

### 13. `get-deployment`
Get Deployment details.

| Parameter   | Type   | Required | Description          |
|-------------|--------|----------|----------------------|
| `namespace` | string | yes      | Namespace            |
| `name`      | string | yes      | Deployment name      |

**Output**: JSON containing `containerName`, `containerImage`, etc.

### 14. `delete-deployment`
Delete a Deployment.

| Parameter   | Type   | Required | Description      |
|-------------|--------|----------|------------------|
| `namespace` | string | yes      | Namespace        |
| `name`      | string | yes      | Deployment name  |

### 15. `update-deployment`
Update Deployment labels, annotations, replicas, or image.

| Parameter       | Type   | Required | Description                                                    |
|-----------------|--------|----------|----------------------------------------------------------------|
| `namespace`     | string | yes      | Namespace                                                      |
| `name`          | string | yes      | Deployment name                                                |
| `label`         | string | no       | Labels (appended, not overwritten)                             |
| `annotation`    | string | no       | Annotations (appended, not overwritten)                        |
| `replica`       | number | no       | Replica count                                                  |
| `containerName` | string | no       | Container name (required when updating image for multi-container)|
| `image`         | string | no       | New image                                                      |

> Update priority: `label` > `annotation` > `image` > `replica`, only one type can be updated at a time.

### 16. `create-deployment`
Create a Deployment.

| Parameter         | Type   | Required | Description                                |
|-------------------|--------|----------|--------------------------------------------|
| `namespace`       | string | yes      | Namespace                                  |
| `name`            | string | yes      | Deployment name                            |
| `label`           | string | no       | Labels, default `app=<name>`               |
| `replica`         | number | no       | Replica count, default 1                   |
| `containerNames`  | string | yes      | Container names, comma-separated           |
| `containerImages` | string | yes      | Images, comma-separated                    |
| `containerPorts`  | string | no       | Ports, default `http:8080`                 |

---

## 4. Service Tools

### 17. `list-service`
List Services.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

**Output**: JSON containing `name`, `namespace`, `type`

### 18. `get-service`
Get Service details.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | yes      | Namespace         |
| `name`      | string | yes      | Service name      |

**Output**: JSON containing `type`, `internalIP`, `externalIP` (for LoadBalancer type), `selectorLabel`

### 19. `delete-service`
Delete a Service.

| Parameter   | Type   | Required | Description   |
|-------------|--------|----------|---------------|
| `namespace` | string | yes      | Namespace     |
| `name`      | string | yes      | Service name  |

### 20. `update-service`
Update Service selector labels or type.

| Parameter        | Type   | Required | Description                                              |
|------------------|--------|----------|----------------------------------------------------------|
| `namespace`      | string | yes      | Namespace                                                |
| `name`           | string | yes      | Service name                                             |
| `selectorLabel`  | string | no       | Selector labels                                          |
| `svctype`        | string | no       | Service type, e.g. `ClusterIP`, `NodePort`, `LoadBalancer`|

### 21. `create-service`
Create a Service.

| Parameter       | Type   | Required | Description                                                                |
|-----------------|--------|----------|----------------------------------------------------------------------------|
| `namespace`     | string | yes      | Namespace                                                                  |
| `name`          | string | yes      | Service name                                                               |
| `selectorLabel` | string | yes      | Selector labels                                                            |
| `svcPort`       | string | yes      | Service port, supports `name:port` or `port` format, multi-port comma-separated |
| `targetPort`    | string | yes      | Target port, multi-port comma-separated                                    |
| `svcType`       | string | no       | Service type, default `ClusterIP`                                          |

> Supports multi-port, `svcPort` and `targetPort` are comma-separated and must have matching counts.

---

## 5. StatefulSet Tools

### 22. `list-statefulset`
List StatefulSets.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `namespace` | string | no       | Specify namespace  |
| `label`     | string | no       | Label selector     |

### 23. `get-statefulset`
Get StatefulSet details.

| Parameter   | Type   | Required | Description         |
|-------------|--------|----------|---------------------|
| `namespace` | string | yes      | Namespace           |
| `name`      | string | yes      | StatefulSet name    |

### 24. `delete-statefulset`
Delete a StatefulSet.

| Parameter   | Type   | Required | Description      |
|-------------|--------|----------|------------------|
| `namespace` | string | yes      | Namespace        |
| `name`      | string | yes      | StatefulSet name |

### 25. `update-statefulset`
Update StatefulSet (same update logic as Deployment).

| Parameter       | Type   | Required | Description                                      |
|-----------------|--------|----------|--------------------------------------------------|
| `namespace`     | string | yes      | Namespace                                        |
| `name`          | string | yes      | StatefulSet name                                 |
| `label`         | string | no       | Labels (appended, not overwritten)               |
| `annotation`    | string | no       | Annotations (appended, not overwritten)           |
| `replica`       | number | no       | Replica count                                    |
| `containerName` | string | no       | Container name                                   |
| `image`         | string | no       | New image                                        |

### 26. `create-statefulset`
Create a StatefulSet (**also creates associated Service and PVC**).

| Parameter         | Type   | Required | Description                                |
|-------------------|--------|----------|--------------------------------------------|
| `namespace`       | string | yes      | Namespace                                  |
| `name`            | string | yes      | StatefulSet name                           |
| `label`           | string | no       | Labels, default `app=<name>`               |
| `containerNames`  | string | no       | Container names, defaults to `name`        |
| `containerImages` | string | yes      | Container images                           |
| `containerPorts`  | number | no       | Container ports, default 8080              |
| `storageValue`    | string | yes      | PVC size, e.g. `1Gi`                       |
| `mountPath`       | string | yes      | Mount path                                 |
| `pvcName`         | string | no       | PVC name, defaults to `name`               |
| `svcType`         | string | no       | Service type, default `ClusterIP`          |
| `svcPort`         | number | no       | Service port, default 8080                 |
| `replica`         | number | no       | Replica count, default 1                   |

---

## 6. DaemonSet Tools

### 27. `list-daemonset`
List DaemonSets.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |
| `label`     | string | no       | Label selector    |

### 28. `get-daemonset`
Get DaemonSet details.

| Parameter   | Type   | Required | Description      |
|-------------|--------|----------|------------------|
| `namespace` | string | yes      | Namespace        |
| `name`      | string | yes      | DaemonSet name   |

### 29. `delete-daemonset`
Delete a DaemonSet.

| Parameter   | Type   | Required | Description   |
|-------------|--------|----------|---------------|
| `namespace` | string | yes      | Namespace     |
| `name`      | string | yes      | DaemonSet name|

### 30. `update-daemonset`
Update DaemonSet (supports labels, annotations, image, **does NOT support replicas**).

| Parameter       | Type   | Required | Description                                |
|-----------------|--------|----------|--------------------------------------------|
| `namespace`     | string | yes      | Namespace                                  |
| `name`          | string | yes      | DaemonSet name                             |
| `label`         | string | no       | Labels (appended, not overwritten)         |
| `annotation`    | string | no       | Annotations (appended, not overwritten)    |
| `containerName` | string | no       | Container name                             |
| `image`         | string | no       | New image                                  |

### 31. `create-daemonset`
Create a DaemonSet.

| Parameter         | Type   | Required | Description                              |
|-------------------|--------|----------|------------------------------------------|
| `namespace`       | string | yes      | Namespace                                |
| `name`            | string | yes      | DaemonSet name                           |
| `label`           | string | no       | Labels, default `app=<name>`             |
| `containerNames`  | string | yes      | Container names, comma-separated         |
| `containerImages` | string | yes      | Images, comma-separated                  |
| `containerPorts`  | string | no       | Ports, default `http:8080`               |

---

## 7. ConfigMap Tools

### 32. `list-configmap`
List ConfigMaps.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

### 33. `get-configmap`
Get ConfigMap details (including data).

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | yes      | Namespace         |
| `name`      | string | yes      | ConfigMap name    |

### 34. `delete-configmap`
Delete a ConfigMap.

| Parameter   | Type   | Required | Description    |
|-------------|--------|----------|----------------|
| `namespace` | string | yes      | Namespace      |
| `name`      | string | yes      | ConfigMap name |

### 35. `create-configmap`
Create a ConfigMap.

| Parameter   | Type   | Required | Description                                  |
|-------------|--------|----------|----------------------------------------------|
| `namespace` | string | yes      | Namespace                                    |
| `name`      | string | yes      | ConfigMap name                               |
| `data`      | string | yes      | Data, format `key=value,key2=value2`         |

---

## 8. Secret Tools

### 36. `list-secret`
List Secrets.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

### 37. `get-secret`
Get Secret details. **Note**: Returned data is base64-encoded raw values, containing passwords, keys, and other sensitive information.

| Parameter   | Type   | Required | Description    |
|-------------|--------|----------|----------------|
| `namespace` | string | yes      | Namespace      |
| `name`      | string | yes      | Secret name    |

### 38. `delete-secret`
Delete a Secret.

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `namespace` | string | yes      | Namespace   |
| `name`      | string | yes      | Secret name |

### 39. `create-secret`
Create an Opaque type Secret.

| Parameter   | Type   | Required | Description                                  |
|-------------|--------|----------|----------------------------------------------|
| `namespace` | string | yes      | Namespace                                    |
| `name`      | string | yes      | Secret name                                  |
| `data`      | string | yes      | Data, format `key=value,key2=value2`         |

---

## 9. Node Tools

### 40. `list-node`
List all nodes.

**No parameters**. Output JSON array containing `name`, `status` (Ready/NotReady), `capacityCPU`, `capacityMemory`, `capacityPods`, `allocatableCPU`, `allocatableMemory`, `allocatablePods`, `labels`, `taints`.

### 41. `get-node`
Get node details.

| Parameter | Type   | Required | Description  |
|-----------|--------|----------|--------------|
| `name`    | string | yes      | Node name    |

**Output**: JSON containing `status`, `kubernetesVersion`, `os`, `kernelVersion`, `architecture`, `podCIDR`, `capacityCPU`, `capacityMemory`, `capacityPods`, `allocatableCPU`, `allocatableMemory`, `allocatablePods`, `labels`, `taints`.

### 42. `delete-node`
Delete a node.

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | yes      | Node name   |

### 43. `update-node`
Update node labels (appended to existing labels, not overwritten).

| Parameter | Type   | Required | Description                         |
|-----------|--------|----------|-------------------------------------|
| `name`    | string | yes      | Node name                           |
| `label`   | string | yes      | Labels, format `key=value`          |

---

## 10. ServiceAccount Tools

### 44. `list-serviceAccount`
List ServiceAccounts.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |
| `label`     | string | no       | Label selector    |

### 45. `get-serviceAccount`
Get ServiceAccount details.

| Parameter   | Type   | Required | Description           |
|-------------|--------|----------|-----------------------|
| `namespace` | string | yes      | Namespace             |
| `name`      | string | yes      | ServiceAccount name   |

### 46. `delete-serviceAccount`
Delete a ServiceAccount.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `namespace` | string | yes      | Namespace          |
| `name`      | string | yes      | ServiceAccount name|

### 47. `create-serviceAccount`
Create a ServiceAccount.

| Parameter   | Type   | Required | Description               |
|-------------|--------|----------|---------------------------|
| `namespace` | string | yes      | Namespace                 |
| `name`      | string | yes      | ServiceAccount name       |
| `label`     | string | no       | Labels                    |

---

## 11. PVC Tools

### 48. `list-pvc`
List PVCs.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

**Output**: JSON array containing `name`, `namespace`, `capacity`, `status`

### 49. `get-pvc`
Get PVC details.

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `namespace` | string | yes      | Namespace   |
| `name`      | string | yes      | PVC name    |

**Output**: JSON containing `capacity`, `accessMode`, `storageClass`, `volume`, `status`

### 50. `delete-pvc`
Delete a PVC.

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `namespace` | string | yes      | Namespace   |
| `name`      | string | yes      | PVC name    |

### 51. `update-pvc`
Update PVC size.

| Parameter   | Type   | Required | Description                  |
|-------------|--------|----------|------------------------------|
| `namespace` | string | yes      | Namespace                    |
| `name`      | string | yes      | PVC name                     |
| `size`      | string | yes      | New size, e.g. `2Gi`         |

### 52. `create-pvc`
Create a PVC.

| Parameter      | Type   | Required | Description                                                    |
|----------------|--------|----------|----------------------------------------------------------------|
| `namespace`    | string | yes      | Namespace                                                      |
| `name`         | string | yes      | PVC name                                                       |
| `size`         | string | yes      | Size, e.g. `1Gi`                                               |
| `storageClass` | string | no       | StorageClass name (empty uses default StorageClass)            |
| `accessMode`   | string | no       | Access mode, default `ReadWriteOnce`, supports comma-separated multiple |

---

## 12. PV Tools

### 53. `list-pv`
List all PVs.

**No parameters**.

### 54. `get-pv`
Get PV details.

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | yes      | PV name     |

### 55. `delete-pv`
Delete a PV.

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | yes      | PV name     |

---

## 13. Role Tools

### 56. `list-role`
List Roles.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

### 57. `get-role`
Get Role details (including rules).

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `namespace` | string | yes      | Namespace   |
| `name`      | string | yes      | Role name   |

### 58. `delete-role`
Delete a Role.

| Parameter   | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `namespace` | string | yes      | Namespace   |
| `name`      | string | yes      | Role name   |

---

## 14. RoleBinding Tools

### 59. `list-rolebinding`
List RoleBindings.

| Parameter   | Type   | Required | Description       |
|-------------|--------|----------|-------------------|
| `namespace` | string | no       | Specify namespace |

### 60. `get-rolebinding`
Get RoleBinding details (including RoleRef and Subjects).

| Parameter   | Type   | Required | Description          |
|-------------|--------|----------|----------------------|
| `namespace` | string | yes      | Namespace            |
| `name`      | string | yes      | RoleBinding name     |

### 61. `delete-rolebinding`
Delete a RoleBinding.

| Parameter   | Type   | Required | Description      |
|-------------|--------|----------|------------------|
| `namespace` | string | yes      | Namespace        |
| `name`      | string | yes      | RoleBinding name |

---

## 15. ClusterRole Tools

### 62. `list-clusterrole`
List all ClusterRoles.

**No parameters**.

### 63. `get-clusterrole`
Get ClusterRole details (including rules).

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| `name`    | string | yes      | ClusterRole name    |

### 64. `delete-clusterrole`
Delete a ClusterRole.

| Parameter | Type   | Required | Description     |
|-----------|--------|----------|-----------------|
| `name`    | string | yes      | ClusterRole name|

---

## 16. ClusterRoleBinding Tools

### 65. `list-clusterrolebinding`
List all ClusterRoleBindings.

**No parameters**.

### 66. `get-clusterrolebinding`
Get ClusterRoleBinding details.

| Parameter | Type   | Required | Description              |
|-----------|--------|----------|--------------------------|
| `name`    | string | yes      | ClusterRoleBinding name  |

### 67. `delete-clusterrolebinding`
Delete a ClusterRoleBinding.

| Parameter | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| `name`    | string | yes      | ClusterRoleBinding name |

---

## 17. StorageClass Tools

### 68. `list-storageClass`
List all StorageClasses.

**No parameters**.

### 69. `get-storageClass`
Get StorageClass details.

| Parameter | Type   | Required | Description           |
|-----------|--------|----------|-----------------------|
| `name`    | string | yes      | StorageClass name     |

**Output**: JSON containing `name`, `provisioner`, `reclaimPolicy`

### 70. `delete-storageClass`
Delete a StorageClass.

| Parameter | Type   | Required | Description       |
|-----------|--------|----------|-------------------|
| `name`    | string | yes      | StorageClass name |

---

## 18. CRD Tools

### 71. `list-crd`
List all CRDs.

**No parameters**.

### 72. `get-crd`
Get CRD details.

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | yes      | CRD name    |

### 73. `delete-crd`
Delete a CRD.

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | yes      | CRD name    |

### 74. `create-crd-with-json`
Create a CRD via JSON/YAML data.

| Parameter  | Type   | Required | Description                        |
|------------|--------|----------|------------------------------------|
| `jsondata` | string | yes      | JSON or YAML definition of the CRD |

---

## 19. Generic Resource Creation Tool

### 75. `create-resource-with-json`
Create any Kubernetes resource via JSON/YAML data (using Dynamic Client).

| Parameter  | Type   | Required | Description                              |
|------------|--------|----------|------------------------------------------|
| `jsondata` | string | yes      | Full JSON or YAML definition of resource |

> Automatically identifies resource type (GVR), supports cluster-scoped and namespace-scoped resources, default namespace is `default`.

---

## 20. Custom Tools (gRPC)

### 76. `custom`
Call an external service via gRPC to operate custom resources.

| Parameter   | Type   | Required | Description                     |
|-------------|--------|----------|---------------------------------|
| `kind`      | string | yes      | Custom resource type            |
| `method`    | string | yes      | Operation method                |
| `name`      | string | no       | Resource name                   |
| `namespace` | string | no       | Namespace                       |
| `jsondata`  | string | no       | JSON data                       |

> Requires the `-customURL` parameter at startup to specify the gRPC service address. Proto definition at [`proto/custom_tool.proto`](proto/custom_tool.proto:1).

---

## 21. Event Tools (Inspection)

### 77. `list-event`
List cluster events, supports filtering by namespace.

| Parameter   | Type   | Required | Description                                      |
|-------------|--------|----------|--------------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all               |

**Output**: JSON array containing `name`, `namespace`, `type` (Normal/Warning), `reason`, `message`, `source`, `firstTime`, `lastTime`, `count`, `kind`, `involved`

### 78. `get-event`
Get detailed information for a specific event.

| Parameter   | Type   | Required | Description           |
|-------------|--------|----------|-----------------------|
| `namespace` | string | yes      | Event namespace       |
| `name`      | string | yes      | Event name            |

---

## 22. ResourceQuota Tools (Inspection)

### 79. `list-resourcequota`
List ResourceQuotas, supports filtering by namespace.

| Parameter   | Type   | Required | Description                                     |
|-------------|--------|----------|-------------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all              |

**Output**: JSON array containing `name`, `namespace`, `hard` (quota limits), `used` (current usage)

### 80. `get-resourcequota`
Get ResourceQuota details.

| Parameter   | Type   | Required | Description            |
|-------------|--------|----------|------------------------|
| `namespace` | string | yes      | Namespace              |
| `name`      | string | yes      | ResourceQuota name     |

---

## 23. LimitRange Tools (Inspection)

### 81. `list-limitrange`
List LimitRanges, supports filtering by namespace.

| Parameter   | Type   | Required | Description                                 |
|-------------|--------|----------|---------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all          |

**Output**: JSON array containing `name`, `namespace`, `limits` (including max/min/default CPU/Memory)

### 82. `get-limitrange`
Get LimitRange details.

| Parameter   | Type   | Required | Description         |
|-------------|--------|----------|---------------------|
| `namespace` | string | yes      | Namespace           |
| `name`      | string | yes      | LimitRange name     |

---

## 24. Endpoint Tools (Inspection)

### 83. `list-endpoint`
List Endpoints, supports filtering by namespace.

| Parameter   | Type   | Required | Description                                  |
|-------------|--------|----------|----------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all           |

**Output**: JSON array containing `name`, `namespace`, `addresses` (including IP, NodeName, Ports)

### 84. `get-endpoint`
Get Endpoint details (including IP and port information).

| Parameter   | Type   | Required | Description     |
|-------------|--------|----------|-----------------|
| `namespace` | string | yes      | Namespace       |
| `name`      | string | yes      | Endpoint name   |

---

## 25. ComponentStatus Tools (Inspection)

### 85. `list-componentstatus`
List health status of all Kubernetes control plane components.

**No parameters**. Output JSON array containing `name` (e.g. etcd-0, kube-apiserver), `status` (True/False), `type`

### 86. `get-componentstatus`
Get health status of a specific component.

| Parameter | Type   | Required | Description                                                    |
|-----------|--------|----------|----------------------------------------------------------------|
| `name`    | string | yes      | Component name (e.g. `etcd-0`, `kube-apiserver`, `kube-scheduler`) |

---

## 26. Top Tools (Monitoring, requires metrics-server)

### 87. `top-pod`
Display CPU and memory usage of Pods (requires metrics-server deployed in the cluster).

| Parameter   | Type   | Required | Description                                      |
|-------------|--------|----------|--------------------------------------------------|
| `namespace` | string | no       | Specify namespace, empty lists all               |

**Output**: JSON array containing `name`, `namespace`, `cpu` (e.g. `100m`), `memory` (e.g. `128Mi`)

### 88. `top-node`
Display CPU and memory usage of Nodes (requires metrics-server deployed in the cluster).

**No parameters**. Output JSON array containing `name`, `cpu`, `memory`

---

## 27. Cluster Health Tools (Monitoring)

### 89. `cluster-health`
Display overall cluster health overview: node summary (Total/Ready/NotReady) + control plane component Pod status.

**No parameters**. Output JSON containing:
- `nodes.total`, `nodes.ready`, `nodes.notReady`
- `controlPlanePods[].name`, `controlPlanePods[].namespace`, `controlPlanePods[].component`, `controlPlanePods[].status`, `controlPlanePods[].restarts`, `controlPlanePods[].node`

### 90. `node-health`
Display detailed health status of all nodes.

**No parameters**. Output JSON array, each node containing `name`, `status` (Ready/NotReady), `kubelet`, `cpu`, `memory`, `pods`, `labels`

---

## 28. Deep Diagnosis Tools (Cluster Fault Diagnosis)

### 91. `describe-pod`
Deep Pod inspection, equivalent to `kubectl describe pod`. Outputs container-level status, Conditions, QoS, resource requests/limits, associated Events.

| Parameter   | Type   | Required | Description          |
|-------------|--------|----------|----------------------|
| `namespace` | string | yes      | Pod namespace        |
| `name`      | string | yes      | Pod name             |

**Output**: Plain text format, containing:
- **Basic Info**: Name, Namespace, Node, Start Time, Status, IP, QoS Class
- **Labels**
- **Conditions**: PodScheduled, Initialized, ContainersReady, Ready status
- **Containers**: Image, Command, Ports, Resource Requests/Limits, Ready, Restart Count, State (Waiting/Running/Terminated with reason and exit code), Last State for each container
- **Init Containers**: Status details
- **Events**: Reverse chronological, showing LastTimestamp, Type, Reason, Message, Count

### 92. `describe-node`
Deep node inspection, equivalent to `kubectl describe node`.

| Parameter | Type   | Required | Description  |
|-----------|--------|----------|--------------|
| `name`    | string | yes      | Node name    |

**Output**: Plain text format, containing:
- **Basic Info**: Creation, Kubelet, OS Image, Kernel, Architecture, PodCIDR, ProviderID
- **Labels**, **Annotations**, **Taints**
- **Conditions**: Ready/DiskPressure/MemoryPressure/PIDPressure/NetworkUnavailable status, Reason, Message, Last Heartbeat
- **Capacity**: cpu, memory, pods, ephemeral-storage
- **Allocatable**: Allocatable amounts for the same dimensions
- **Pods**: List of all Pods running on the node (ns/name/status)

### 93. `list-node-pods`
List all Pods running on a specific node.

| Parameter | Type   | Required | Description  |
|-----------|--------|----------|--------------|
| `node`    | string | yes      | Node name    |

**Output**: JSON array containing `namespace`, `name`, `status`, `node`

### 94. `describe-service`
Deep Service inspection, equivalent to `kubectl describe service`.

| Parameter   | Type   | Required | Description     |
|-------------|--------|----------|-----------------|
| `namespace` | string | yes      | Namespace       |
| `name`      | string | yes      | Service name    |

**Output**: Plain text format, containing:
- **Basic Info**: Type, ClusterIP, ExternalIPs, ExternalName, LoadBalancerIP, Session Affinity
- **Labels**, **Selector**
- **Ports**: name, port/protocol → targetPort
- **Endpoints**: addresses for each subset (including NodeName)
- **Events**

### 95. `describe-deployment`
Deep inspect a deployment, similar to `kubectl describe deployment`.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | Yes | Namespace of the deployment |
| `name` | string | Yes | Name of the deployment |

**Output**: Plain text including:
- **Basic info**: Strategy (RollingUpdate MaxSurge/MaxUnavailable), Replicas (desired/updated/total/available/unavailable), Revision History Limit, Min Ready Seconds, Rollout Status
- **Labels**, **Selector**
- **Conditions**: Available/Progressing/ReplicaFailure status, Reason, Message, LastUpdateTime
- **Containers**: Image, Resource Requests/Limits
- **Pods**: Count by Phase (Running/Pending/Failed etc.)
- **Events**

### 111. `check-apiserver-health`
Probe API Server health endpoints. **Does not rely on Pod labels or ComponentStatus API, works for binary-deployed apiserver as well.**

**No parameters**. Directly calls three apiserver endpoints:
- `/healthz?verbose` — Detailed health check results
- `/livez?verbose` — Liveness probe
- `/readyz?verbose` — Readiness probe

**Output**: Plain text, HTTP status code for each endpoint and all non-`ok` check items.

### 111. `check-apiserver-metrics`
Fetch API Server `/metrics` and analyze key performance indicators. **Works for binary-deployed apiserver as well.**

**No parameters**. Requires apiserver to have `--authorization-always-allow-paths=/metrics` configured or RBAC permission granted.

**Output**: Contains:
- **Current Inflight Requests**: Number of mutating and readOnly requests currently being processed
- **Request Counts**: Total requests, errors, and error rate by Verb (GET/LIST/WATCH/POST/PUT/PATCH/DELETE)
- **Request Latency**: Estimated p50/p90/p99 latency by Verb
- **Top Error Endpoints**: 4xx/5xx errors sorted by code

---

## 29. Helm Tools

### 111. `helm-list-releases`
List Helm Releases.

| Parameter       | Type    | Required | Description                                      |
|-----------------|---------|----------|--------------------------------------------------|
| `namespace`     | string  | no       | Filter by namespace                              |
| `allNamespaces` | boolean | no       | List releases in all namespaces                  |
| `output`        | string  | no       | Output format: `json` or `yaml` for full objects |

### 111. `helm-get-release`
Get Release details.

| Parameter   | Type   | Required | Description               |
|-------------|--------|----------|---------------------------|
| `name`      | string | yes      | Release name              |
| `namespace` | string | yes      | Release namespace         |
| `output`    | string | no       | Output format             |

### 111. `helm-get-values`
Get Release values.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `name`      | string | yes      | Release name       |
| `namespace` | string | yes      | Release namespace  |

### 111. `helm-install`
Install a Chart.

| Parameter   | Type   | Required | Description                                                    |
|-------------|--------|----------|----------------------------------------------------------------|
| `name`      | string | yes      | Release name                                                   |
| `namespace` | string | yes      | Target namespace                                               |
| `chart`     | string | yes      | Chart reference, e.g. `stable/nginx-ingress`, or local path    |
| `version`   | string | no       | Chart version                                                  |
| `values`    | string | no       | Inline YAML/JSON format values                                 |
| `set`       | string | no       | Command-line set values, format `key1=val1,key2=val2`          |

> `values` and `set` can be used together, `set` has higher priority.

### 111. `helm-upgrade`
Upgrade a Release.

| Parameter     | Type    | Required | Description                                     |
|---------------|---------|----------|-------------------------------------------------|
| `name`        | string  | yes      | Release name                                    |
| `namespace`   | string  | yes      | Release namespace                               |
| `chart`       | string  | yes      | New Chart reference                             |
| `version`     | string  | no       | Chart version                                   |
| `values`      | string  | no       | Inline YAML/JSON format values                  |
| `set`         | string  | no       | Command-line set values                         |
| `reuseValues` | boolean | no       | Whether to reuse previous values, default `true`|

### 111. `helm-uninstall`
Uninstall a Release.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `name`      | string | yes      | Release name       |
| `namespace` | string | yes      | Release namespace  |

### 111. `helm-rollback`
Rollback a Release to a specific revision.

| Parameter   | Type   | Required | Description                                    |
|-------------|--------|----------|------------------------------------------------|
| `name`      | string | yes      | Release name                                   |
| `namespace` | string | yes      | Release namespace                              |
| `revision`  | number | no       | Target revision number, defaults to previous   |

### 111. `helm-history`
Get Release revision history.

| Parameter   | Type   | Required | Description                |
|-------------|--------|----------|----------------------------|
| `name`      | string | yes      | Release name               |
| `namespace` | string | yes      | Release namespace          |
| `max`       | number | no       | Maximum number of revisions|

**Output**: JSON array, each revision containing `revision`, `status`, `chart`, `appVersion`, `description`, `updated`

### 111. `helm-get-manifest`
Get rendered manifest of a Release.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `name`      | string | yes      | Release name       |
| `namespace` | string | yes      | Release namespace  |

### 111. `helm-get-notes`
Get NOTES.txt content of a Release.

| Parameter   | Type   | Required | Description        |
|-------------|--------|----------|--------------------|
| `name`      | string | yes      | Release name       |
| `namespace` | string | yes      | Release namespace  |

### 111. `helm-list-repos`
List added Helm repositories.

**No parameters**.

### 111. `helm-add-repo`
Add a Helm repository and download its index.

| Parameter | Type   | Required | Description      |
|-----------|--------|----------|------------------|
| `name`    | string | yes      | Repository name  |
| `url`     | string | yes      | Repository URL   |

### 111. `helm-remove-repo`
Remove a Helm repository.

| Parameter | Type   | Required | Description      |
|-----------|--------|----------|------------------|
| `name`    | string | yes      | Repository name  |

### 111. `helm-update-repos`
Update Helm repository index files.

| Parameter | Type   | Required | Description                                 |
|-----------|--------|----------|---------------------------------------------|
| `name`    | string | no       | Repository name, empty updates all repos    |

### 111. `helm-search-repo`
Search Charts across all added repositories.

| Parameter | Type   | Required | Description                                                     |
|-----------|--------|----------|-----------------------------------------------------------------|
| `keyword` | string | yes      | Search keyword (matches chart name, description, appVersion)   |

---

## Startup

```bash
# HTTP mode (default, listens on :8080)
go run main.go

# Specify kubeconfig
go run main.go -kubeconfigPath /path/to/kubeconfig

# stdio mode
go run main.go -mode stdio

# Specify gRPC custom service address
go run main.go -customURL localhost:50051

# Enable API key authentication
go run main.go -apiKey my-secret-key
```

## Command-Line Flags

| Flag              | Default  | Description                                              |
|-------------------|----------|----------------------------------------------------------|
| `-mode`           | `http`   | Runtime mode: `http` or `stdio`                          |
| `-kubeconfigPath` | `""`     | kubeconfig file path, empty uses InCluster mode          |
| `-apiKey`         | `""`     | HTTP API key authentication (optional)                   |
| `-customURL`      | `""`     | gRPC custom resource service address                     |

## Client Initialization Logic

[`kubernetes/client/client.go`](kubernetes/client/client.go:19) [`InitializeClients()`](kubernetes/client/client.go:19) returns 6 clients:

1. `*kubernetes.Clientset` — Standard Kubernetes client
2. `dynamic.Interface` — Dynamic client (for generic resource creation)
3. `discovery.DiscoveryInterface` — API discovery client
4. `*apiextensionsclient.Clientset` — CRD client
5. `*metricsv.Interface` — Metrics client (for Top tools)

Connection priority: `-kubeconfigPath` flag > InCluster mode (running inside a Pod).

## Tool Registration Mapping Table

[`main.go`](main.go:38) registers all tools. Below is the mapping between tool names and implementation functions:

| Tool Name                  | Schema Variable              | Implementation Function                     |
|----------------------------|------------------------------|---------------------------------------------|
| `list-pod`                 | `tools.ListPod`              | `pod.ListPod`                               |
| `get-pod`                  | `tools.GetPod`               | `pod.GetPod`                                |
| `delete-pod`               | `tools.DeletePod`            | `pod.DeletePod`                             |
| `update-pod`               | `tools.UpdatePod`            | `pod.UpdatePod`                             |
| `create-pod`               | `tools.CreatePod`            | `pod.CreatePod`                             |
| `pod-log`                  | `tools.PodLog`               | `pod.PodLog`                                |
| `list-ns`                  | `tools.ListNS`               | `namespace.ListNS`                          |
| `get-ns`                   | `tools.GetNS`                | `namespace.GetNS`                           |
| `delete-ns`                | `tools.DeleteNS`             | `namespace.DeleteNS`                        |
| `update-ns`                | `tools.UpdateNS`             | `namespace.UpdateNS`                        |
| `create-ns`                | `tools.CreateNS`             | `namespace.CreateNS`                        |
| `list-node`                | `tools.ListNode`             | `node.ListNode`                             |
| `get-node`                 | `tools.GetNode`              | `node.GetNode`                              |
| `delete-node`              | `tools.DeleteNode`           | `node.DeleteNode`                           |
| `update-node`              | `tools.UpdateNode`           | `node.UpdateNode`                           |
| `list-deployment`          | `tools.ListDeployment`       | `deployment.ListDeployment`                 |
| `get-deployment`           | `tools.GetDeployment`        | `deployment.GetDeployment`                  |
| `delete-deployment`        | `tools.DeleteDeployment`     | `deployment.DeleteDeployment`               |
| `create-deployment`        | `tools.CreateDeployment`     | `deployment.CreateDeployment`               |
| `update-deployment`        | `tools.UpdateDeployment`     | `deployment.UpdateDeployment`               |
| `list-daemonset`           | `tools.ListDaemonset`        | `daemonset.ListDaemonset`                   |
| `get-daemonset`            | `tools.GetDaemonset`         | `daemonset.GetDaemonset`                    |
| `delete-daemonset`         | `tools.DeleteDaemonset`      | `daemonset.DeleteDaemonset`                 |
| `update-daemonset`         | `tools.UpdateDaemonset`      | `daemonset.UpdateDaemonset`                 |
| `create-daemonset`         | `tools.CreateDaemonset`      | `daemonset.CreateDaemonset`                 |
| `list-statefulset`         | `tools.ListStatefulset`      | `statefulset.ListStatefulset`               |
| `get-statefulset`          | `tools.GetStatefulset`       | `statefulset.GetStatefulset`                |
| `delete-statefulset`       | `tools.DeleteStatefulset`    | `statefulset.DeleteStatefulset`             |
| `update-statefulset`       | `tools.UpdateStatefulset`    | `statefulset.UpdateStatefulset`             |
| `create-statefulset`       | `tools.CreateStatefulset`    | `statefulset.CreateStatefulset`             |
| `list-service`             | `tools.ListService`          | `service.ListService`                       |
| `get-service`              | `tools.GetService`           | `service.GetService`                        |
| `delete-service`           | `tools.DeleteService`        | `service.DeleteService`                     |
| `update-service`           | `tools.UpdateService`        | `service.UpdateService`                     |
| `create-service`           | `tools.CreateService`        | `service.CreateService`                     |
| `list-configmap`           | `tools.ListConfigmap`        | `configmap.ListConfigmap`                   |
| `get-configmap`            | `tools.GetConfigmap`         | `configmap.GetConfigmap`                    |
| `delete-configmap`         | `tools.DeleteConfigmap`      | `configmap.DeleteConfigmap`                 |
| `create-configmap`         | `tools.CreateConfigmap`      | `configmap.CreateConfigmap`                 |
| `list-secret`              | `tools.ListSecret`           | `secret.ListSecret`                         |
| `get-secret`               | `tools.GetSecret`            | `secret.GetSecret`                          |
| `delete-secret`            | `tools.DeleteSecret`         | `secret.DeleteSecret`                       |
| `create-secret`            | `tools.CreateSecret`         | `secret.CreateSecret`                       |
| `list-serviceAccount`      | `tools.ListSA`               | `serviceaccount.ListSA`                     |
| `get-serviceAccount`       | `tools.GetSA`                | `serviceaccount.GetSA`                      |
| `delete-serviceAccount`    | `tools.DeleteSA`             | `serviceaccount.DeleteSA`                   |
| `create-serviceAccount`    | `tools.CreateSA`             | `serviceaccount.CreateSA`                   |
| `list-role`                | `tools.ListRole`             | `role.ListRole`                             |
| `get-role`                 | `tools.GetRole`              | `role.GetRole`                              |
| `delete-role`              | `tools.DeleteRole`           | `role.DeleteRole`                           |
| `list-rolebinding`         | `tools.ListRB`               | `rolebinding.ListRB`                        |
| `get-rolebinding`          | `tools.GetRB`                | `rolebinding.GetRB`                         |
| `delete-rolebinding`       | `tools.DeleteRB`             | `rolebinding.DeleteRB`                      |
| `list-pvc`                 | `tools.ListPVC`              | `pvc.ListPVC`                               |
| `get-pvc`                  | `tools.GetPVC`               | `pvc.GetPVC`                                |
| `delete-pvc`               | `tools.DeletePVC`            | `pvc.DeletePVC`                             |
| `update-pvc`               | `tools.UpdatePVC`            | `pvc.UpdatePVC`                             |
| `create-pvc`               | `tools.CreatePVC`            | `pvc.CreatePVC`                             |
| `list-pv`                  | `tools.ListPV`               | `pv.ListPV`                                 |
| `get-pv`                   | `tools.GetPV`                | `pv.GetPV`                                  |
| `delete-pv`                | `tools.DeletePV`             | `pv.DeletePV`                               |
| `list-clusterrole`         | `tools.ListCR`               | `clusterrole.ListCR`                        |
| `get-clusterrole`          | `tools.GetCR`                | `clusterrole.GetCR`                         |
| `delete-clusterrole`       | `tools.DeleteCR`             | `clusterrole.DeleteCR`                      |
| `list-clusterrolebinding`  | `tools.ListCRB`              | `clusterrolebinding.ListCRB`                |
| `get-clusterrolebinding`   | `tools.GetCRB`               | `clusterrolebinding.GetCRB`                 |
| `delete-clusterrolebinding`| `tools.DeleteCRB`            | `clusterrolebinding.DeleteCRB`              |
| `list-storageClass`        | `tools.ListSC`               | `storageclass.ListSC`                       |
| `get-storageClass`         | `tools.GetSC`                | `storageclass.GetSC`                        |
| `delete-storageClass`      | `tools.DeleteSC`             | `storageclass.DeleteSC`                     |
| `list-crd`                 | `tools.ListCRD`              | `crd.ListCRD`                               |
| `get-crd`                  | `tools.GetCRD`               | `crd.GetCRD`                                |
| `delete-crd`               | `tools.DeleteCRD`            | `crd.DeleteCRD`                             |
| `create-crd-with-json`     | `tools.CreateCRDWithJson`    | `crd.CreateCRDWithJson`                     |
| `create-resource-with-json`| `tools.CreateResourceWithJSon`| `createresource.CreateResourceWithJson`    |
| `custom`                   | `tools.Custom`               | `custom.Custom`                             |
| `list-event`               | `tools.ListEvent`            | `event.ListEvent`                           |
| `get-event`                | `tools.GetEvent`             | `event.GetEvent`                            |
| `list-resourcequota`       | `tools.ListResourceQuota`    | `resourcequota.ListResourceQuota`           |
| `get-resourcequota`        | `tools.GetResourceQuota`     | `resourcequota.GetResourceQuota`            |
| `list-limitrange`          | `tools.ListLimitRange`       | `limitrange.ListLimitRange`                 |
| `get-limitrange`           | `tools.GetLimitRange`        | `limitrange.GetLimitRange`                  |
| `list-endpoint`            | `tools.ListEndpoint`         | `endpoint.ListEndpoint`                     |
| `get-endpoint`             | `tools.GetEndpoint`          | `endpoint.GetEndpoint`                      |
| `list-componentstatus`     | `tools.ListComponentStatus`  | `componentstatus.ListComponentStatus`       |
| `get-componentstatus`      | `tools.GetComponentStatus`   | `componentstatus.GetComponentStatus`        |
| `top-pod`                  | `tools.TopPod`               | `top.TopPod`                                |
| `top-node`                 | `tools.TopNode`              | `top.TopNode`                               |
| `cluster-health`           | `tools.GetClusterHealth`     | `clusterhealth.GetClusterHealth`            |
| `node-health`              | `tools.ListNodeHealth`       | `clusterhealth.ListNodeHealth`              |
| `describe-pod`             | `tools.DescribePod`          | `diagnose.DescribePod`                      |
| `describe-node`            | `tools.DescribeNode`         | `diagnose.DescribeNode`                     |
| `list-node-pods`           | `tools.ListNodePods`         | `diagnose.ListNodePods`                     |
| `describe-service`         | `tools.DescribeService`      | `diagnose.DescribeService`                  |
| `describe-deployment`      | `tools.DescribeDeployment`   | `diagnose.DescribeDeployment`               |
| `check-apiserver-health`   | `tools.CheckAPIServerHealth` | `diagnose.CheckAPIServerHealth`             |
| `check-apiserver-metrics`  | `tools.CheckAPIServerMetrics`| `diagnose.CheckAPIServerMetrics`            |
| `helm-list-releases`       | `tools.ListHelmReleases`     | `helm.ListHelmReleases`                     |
| `helm-get-release`         | `tools.GetHelmRelease`       | `helm.GetHelmRelease`                       |
| `helm-get-values`          | `tools.GetHelmReleaseValues` | `helm.GetHelmReleaseValues`                 |
| `helm-install`             | `tools.InstallHelmRelease`   | `helm.InstallHelmRelease`                   |
| `helm-upgrade`             | `tools.UpgradeHelmRelease`   | `helm.UpgradeHelmRelease`                   |
| `helm-uninstall`           | `tools.UninstallHelmRelease` | `helm.UninstallHelmRelease`                 |
| `helm-rollback`            | `tools.RollbackHelmRelease`  | `helm.RollbackHelmRelease`                  |
| `helm-history`             | `tools.GetHelmReleaseHistory`| `helm.GetHelmReleaseHistory`                |
| `helm-get-manifest`        | `tools.GetHelmReleaseManifest`| `helm.GetHelmReleaseManifest`              |
| `helm-get-notes`           | `tools.GetHelmReleaseNotes`  | `helm.GetHelmReleaseNotes`                  |
| `helm-list-repos`          | `tools.ListHelmRepos`        | `helm.ListHelmRepos`                        |
| `helm-add-repo`            | `tools.AddHelmRepo`          | `helm.AddHelmRepo`                          |
| `helm-remove-repo`         | `tools.RemoveHelmRepo`       | `helm.RemoveHelmRepo`                       |
| `helm-update-repos`        | `tools.UpdateHelmRepos`      | `helm.UpdateHelmRepos`                      |
| `helm-search-repo`         | `tools.SearchHelmRepo`       | `helm.SearchHelmRepo`                       |

---
