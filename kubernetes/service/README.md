# Service 操作

### 创建

创建 Service 可用字段：
- Namespace：必填
- Name：必填
- SelectorLabel：必填
- TargetPort：必填（多个目标端口用逗号分隔，例如：8080,9090）
- ServicePort：必填（多个端口用逗号分隔，例如：http:8080,metrics:9090）
- ServiceType：可选

### 列表

指定命名空间列出 Service：
- Namespace：必填

列出所有命名空间的 Service 无需参数。

### 查询

查询指定命名空间的 Service：
- Namespace：必填
- Name：必填

### 删除

删除指定命名空间的 Service：
- Namespace：必填
- Name：必填

### 更新

更新 Service（支持选择器标签或服务类型，二选一）：
- Namespace：必填
- Name：必填
- SelectorLabel：可选
- Service Type：可选