# Secret 操作

### 创建

创建 Secret 可用字段：
- Namespace：必填
- Name：必填
- Data：必填（多个数据用逗号分隔，例如：password=Passw0rd@123,username=admin）

### 列表

指定命名空间列出 Secret：
- Namespace：必填

列出所有命名空间的 Secret 无需参数。

### 查询

查询指定命名空间的 Secret：
- Namespace：必填
- Name：必填

### 删除

删除指定命名空间的 Secret：
- Namespace：必填
- Name：必填