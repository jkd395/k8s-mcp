# Pod 操作

### 创建

创建 Pod 可用字段：
- Namespace：必填
- Name：必填
- ContainerNames：必填（多个容器用逗号分隔，例如：nginx,apache）
- ContainerImages：必填（多个镜像用逗号分隔，例如：nginx:latest,apache2@latest）
- ContainerPorts：必填（多个容器端口用逗号分隔，单个容器多个端口用 `|` 分隔，例如：http:80|https:443,http:80）
- Label：可选

### 列表

指定命名空间列出 Pod：
- Namespace：必填
- Label：可选

列出所有命名空间的 Pod：
- Label：可选

### 查询

查询指定命名空间的 Pod：
- Namespace：必填
- Name：必填

### 删除

删除指定命名空间的 Pod：
- Namespace：必填
- Name：必填

### 更新

更新 Pod（仅支持标签）：
- Namespace：必填
- Name：必填
- label：必填

### 日志

Pod 日志参数：
- Namespace：必填
- Name：必填
- ContainerName：必填
- Tailline：可选
