[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)

# nodeport-control-policy

这是一个 Kubewarden 策略，用于控制 Kubernetes 集群中 NodePort 类型服务的创建。通过此策略，集群管理员可以限制用户创建 NodePort 类型的服务，从而加强集群的网络安全控制。

## 简介

NodePort 服务类型允许从集群外部直接访问集群内的服务，这可能会带来潜在的安全风险。本策略允许集群管理员：
- 完全禁止创建 NodePort 类型的服务
- 允许创建 NodePort 类型的服务
- 对其他类型的服务（如 ClusterIP、LoadBalancer）不做限制

## 配置

策略的配置非常简单，只需要一个布尔值参数：

```json
{
  "disable_nodeport": true
}
```

配置参数说明：
- `disable_nodeport`: 布尔值
  - `true`: 禁止创建 NodePort 类型的服务
  - `false`: 允许创建 NodePort 类型的服务
  - 默认值：`false`（如果未设置）

## 使用示例

1. 禁止创建 NodePort 服务：
```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: disable-nodeport
spec:
  module: registry://ghcr.io/kubewarden/policies/nodeport-control-policy:v0.1.0
  settings:
    disable_nodeport: true
  rules:
    - apiGroups: [""]
      apiVersions: ["v1"]
      resources: ["services"]
      operations: ["CREATE"]
```

2. 允许创建 NodePort 服务：
```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: allow-nodeport
spec:
  module: registry://ghcr.io/kubewarden/policies/nodeport-control-policy:v0.1.0
  settings:
    disable_nodeport: false
  rules:
    - apiGroups: [""]
      apiVersions: ["v1"]
      resources: ["services"]
      operations: ["CREATE"]
```

## 代码组织

- `settings.go`: 包含策略配置的解析和验证逻辑
- `validate.go`: 包含服务类型验证的核心逻辑
- `main.go`: 包含策略入口点的注册代码

## 测试

### 单元测试

运行单元测试：
```console
make test
```

### 端到端测试

运行端到端测试：
```console
make e2e-tests
```

测试场景包括：
1. 当策略配置为禁用 NodePort 时，拒绝创建 NodePort 类型的服务
2. 当策略配置为允许 NodePort 时，允许创建 NodePort 类型的服务
3. 无论策略如何配置，始终允许创建 ClusterIP 类型的服务
4. 使用默认配置时的行为测试

## 实现细节

本策略使用 Go 语言编写，并使用 TinyGo 编译为 WebAssembly。主要使用了以下库：
- [Kubewarden's Go SDK](https://github.com/kubewarden/policy-sdk-go)
- [Kubernetes Go types](https://github.com/kubewarden/k8s-objects)

## 安全考虑

NodePort 服务类型会在所有节点上开放端口，可能带来以下安全风险：
1. 增加集群的攻击面
2. 可能暴露内部服务
3. 可能违反网络安全策略

因此，在生产环境中建议谨慎使用 NodePort 服务，可以考虑使用此策略来控制其使用。

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个策略。

## 许可证

Apache-2.0
