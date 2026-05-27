# StatefulSet 操作

### 创建

创建 StatefulSet 可用字段（当前仅支持单容器）：
- Namespace：必填
- Name：必填
- ContainerNames：可选
- ContainerImages：必填
- ContainerPorts：可选
- StorageValue：必填（例如：1Gi）
- MountPath：必填
- PVCName：可选
- ServiceType：可选
- ServicePort：可选
- Label：可选
- Replica：可选

### 列表

指定命名空间列出 StatefulSet：
- Namespace：必填
- Label：可选

列出所有命名空间的 StatefulSet：
- Label：可选

### 查询

查询指定命名空间的 StatefulSet：
- Namespace：必填
- Name：必填

### 删除

删除指定命名空间的 StatefulSet：
- Namespace：必填
- Name：必填

### 更新

更新 StatefulSet（支持标签、注解、副本数或镜像）：
- Namespace：必填
- Name：必填
- Label：可选
- Annotation：可选
- Replica：可选
- ContainerName：可选
- Image：可选（多容器时需指定容器名称）