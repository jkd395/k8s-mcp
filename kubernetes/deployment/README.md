# Deployment 操作

### 创建

创建 Deployment 可用字段：
- Namespace：必填
- Name：必填
- ContainerNames：必填（多个容器用逗号分隔，例如：nginx,apache）
- ContainerImages：必填（多个镜像用逗号分隔，例如：nginx:latest,apache2@latest）
- ContainerPorts：必填（多个容器端口用逗号分隔，单个容器多个端口用 `|` 分隔，例如：http:80|https:443,http:80）
- Label：可选
- Replica：可选

### 列表

指定命名空间列出 Deployment：
- Namespace：必填
- Label：可选

列出所有命名空间的 Deployment：
- Label：可选

### 查询

查询指定命名空间的 Deployment：
- Namespace：必填
- Name：必填

### 删除

删除指定命名空间的 Deployment：
- Namespace：必填
- Name：必填

### 更新

更新 Deployment（支持标签、注解、副本数或镜像，一次只能更新一种）：
- Namespace：必填
- Name：必填
- Label：可选
- Annotation：可选
- Replica：可选
- ContainerName：可选
- Image：可选（多容器时需指定容器名称）