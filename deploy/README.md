# 部署说明

## 文件清单

| 文件 | 用途 |
|------|------|
| `namespace.yaml` | 命名空间 `k8s-mcp` |
| `serviceaccount.yaml` | ServiceAccount `k8s-mcp` |
| `rbac.yaml` | ClusterRole + ClusterRoleBinding，授权 MCP Server 全部资源操作 |
| `deployment.yaml` | Deployment，使用 InCluster 模式连接集群 |
| `service.yaml` | ClusterIP Service，暴露端口 8080 |
| `ingress.yaml` | Ingress，通过域名对外暴露 `/mcp` 端点 |

## 部署步骤

### 1. 构建镜像

```bash
# 在项目根目录执行
docker build -t k8s-mcp:latest .
# 或推送到仓库
# docker build -t registry.example.com/k8s-mcp:latest .
# docker push registry.example.com/k8s-mcp:latest
```

### 2. 部署到集群

```bash
kubectl apply -f namespace.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

### 3. 验证

```bash
kubectl -n k8s-mcp get pods
kubectl -n k8s-mcp logs deploy/k8s-mcp
```

### 4. 修改 Ingress 域名

编辑 `ingress.yaml`，将 `mcp.example.com` 替换为实际域名，然后重新应用。

### 5. 配置 MCP 客户端

```json
{
  "mcpServers": {
    "kubernetes": {
      "url": "http://mcp.example.com/mcp"
    }
  }
}
```

## 安全建议

- 生产环境建议在 Ingress 上配置 TLS
- 如需 API 密钥认证，在 `deployment.yaml` 的 `args` 中添加 `-apiKey=<密钥>`，客户端配置 `Authorization: Bearer <密钥>` 头
- 如需缩小权限范围，编辑 `rbac.yaml` 中的 ClusterRole rules
