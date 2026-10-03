# 更新 K8s 测试环境（10.10.1.66）的 auth-service 和 gateway 操作步骤

> 目标：把本地最新代码构建成镜像，部署到 K8s 测试环境 10.10.1.66，并配置 pilot 鉴权环境变量。
> 涉及三台机器：本机（Windows）+ 构建机（192.168.18.35）+ K8s 集群（10.10.1.66）。

---

## 机器信息

| 角色 | 地址 | 账号 | 用途 |
|---|---|---|---|
| 本机 | Windows | - | 修改代码、上传代码 |
| 构建机 | 192.168.18.35 | root / user@dev2025 | Docker 构建镜像 |
| K8s 集群 | 10.10.1.66 | root / User@dev123 | 部署 Pod、NodePort 30080 |

---

## 一、修改本机 Dockerfile（添加 GOPROXY 解决网络问题）

### 1.1 修改 `services/auth-service/Dockerfile`

添加 `GOPROXY` 和 `GOWORK=off`，去掉 `COPY go.work`：

```dockerfile
FROM golang:1.25.13-alpine AS build

ENV GOPROXY=https://goproxy.cn,direct
ENV GOWORK=off

WORKDIR /src
COPY pkg ./pkg
COPY services/auth-service ./services/auth-service

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/auth-service && \
    CGO_ENABLED=0 go build -tags stdjson -ldflags "-s -w" -o /out/auth-service .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -H -u 65532 ani

COPY --from=build /out/auth-service /usr/local/bin/auth-service

USER 65532:65532
EXPOSE 9101 9201

ENTRYPOINT ["/usr/local/bin/auth-service"]
```

### 1.2 修改 `services/ani-gateway/Dockerfile`

添加 `GOPROXY`，把 tools stage 改为 `COPY` 本地预下载工具（避免 build 时 curl 下载超时）：

```dockerfile
FROM golang:1.25.13-alpine AS build

ENV GOPROXY=https://goproxy.cn,direct
ENV GOWORK=off

WORKDIR /src
COPY pkg ./pkg
COPY services/ani-gateway ./services/ani-gateway

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/ani-gateway && \
    CGO_ENABLED=0 go build -tags stdjson -ldflags "-s -w" -o /out/ani-gateway .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -H -u 65532 ani

COPY --from=build /out/ani-gateway /usr/local/bin/ani-gateway
COPY ani-tools/helm /usr/local/bin/helm
COPY ani-tools/vcluster /usr/local/bin/vcluster
COPY ani-tools/kubectl /usr/local/bin/kubectl

USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/ani-gateway"]
```

---

## 二、在构建机（192.168.18.35）上预下载工具

ani-gateway 镜像需要 helm、vcluster、kubectl 三个二进制，预下载到 `/root/ani-build/ani-tools/`：

```bash
ssh root@192.168.18.35   # 密码: user@dev2025

mkdir -p /root/ani-build/ani-tools && cd /root/ani-build/ani-tools

# 下载 helm（3 秒）
curl -fsSL --connect-timeout 10 --max-time 120 \
  "https://get.helm.sh/helm-v3.19.0-linux-amd64.tar.gz" -o helm.tar.gz
tar -xzf helm.tar.gz && cp linux-amd64/helm . && chmod +x helm

# 下载 vcluster（18 秒）
curl -fsSL --connect-timeout 10 --max-time 180 \
  "https://gh-proxy.com/https://github.com/loft-sh/vcluster/releases/download/v0.34.1/vcluster-linux-amd64" -o vcluster
chmod +x vcluster

# 下载 kubectl（11 秒）
curl -fsSL --connect-timeout 10 --max-time 120 \
  "https://dl.k8s.io/release/v1.36.1/bin/linux/amd64/kubectl" -o kubectl
chmod +x kubectl

# 验证
ls -la /root/ani-build/ani-tools/
# 预期: helm(60MB) + vcluster(96MB) + kubectl(59MB) = 215MB
```

---

## 三、上传代码到构建机

把本机 `c:\ProgramProject\ChangQinYun\kuberai\ANI\repo` 下的必要目录上传到 `192.168.18.35:/root/ani-build/`：

需要上传的内容：
- `repo/pkg/` → `/root/ani-build/pkg/`
- `repo/services/auth-service/` → `/root/ani-build/services/auth-service/`
- `repo/services/ani-gateway/` → `/root/ani-build/services/ani-gateway/`

用 scp 或 SFTP 上传（如果之前已上传过代码，只需上传修改后的两个 Dockerfile 覆盖）：

```bash
# 在本机 PowerShell 执行
scp -r c:\ProgramProject\ChangQinYun\kuberai\ANI\repo\pkg root@192.168.18.35:/root/ani-build/
scp -r c:\ProgramProject\ChangQinYun\kuberai\ANI\repo\services\auth-service root@192.168.18.35:/root/ani-build/services/
scp -r c:\ProgramProject\ChangQinYun\kuberai\ANI\repo\services\ani-gateway root@192.168.18.35:/root/ani-build/services/
```

---

## 四、在构建机上构建镜像

```bash
ssh root@192.168.18.35   # 密码: user@dev2025
cd /root/ani-build

# 构建 auth-service 镜像
docker build -f services/auth-service/Dockerfile \
  -t docker.changqingyun.cn/ani/ani-auth-service:dev-20260826 .

# 构建 ani-gateway 镜像
docker build -f services/ani-gateway/Dockerfile \
  -t docker.changqingyun.cn/ani/ani-gateway:dev-20260826 .

# 查看构建结果
docker images | grep -E "ani-auth-service|ani-gateway"
# 预期:
#   ani-auth-service:dev-20260826   34.1MB
#   ani-gateway:dev-20260826         262MB
```

---

## 五、推送镜像到镜像仓库

```bash
# 在构建机 192.168.18.35 上执行
# 如果未登录，先登录（账号密码问仓库管理员）
# docker login docker.changqingyun.cn

# push auth-service
docker push docker.changqingyun.cn/ani/ani-auth-service:dev-20260826

# push ani-gateway
docker push docker.changqingyun.cn/ani/ani-gateway:dev-20260826
```

---

## 六、在 K8s 集群（10.10.1.66）上更新 Deployment

```bash
ssh root@10.10.1.66   # 密码: User@dev123

# 6.1 更新 auth-service 镜像
kubectl set image deployment/ani-auth-service \
  auth-service=docker.changqingyun.cn/ani/ani-auth-service:dev-20260826 \
  -n ani-system

# 6.2 更新 ani-gateway 镜像
kubectl set image deployment/ani-gateway \
  ani-gateway=docker.changqingyun.cn/ani/ani-gateway:dev-20260826 \
  -n ani-system

# 6.3 添加 pilot 鉴权环境变量（pilot 新链路验证用）
kubectl set env deployment/ani-gateway -n ani-system \
  GATEWAY_AUTHZ_POLICY_MODE=pilot \
  GATEWAY_AUTHZ_PILOT_OPERATIONS=listQuotaMeta

# 6.4 ★ 关键：把 ANI_AUTH_MODE 从 dev 改成 auth_service（否则 pilot 不生效）
kubectl set env deployment/ani-gateway -n ani-system \
  ANI_AUTH_MODE=auth_service

# 6.5 等待两个 deployment rollout 完成
kubectl rollout status deployment/ani-auth-service -n ani-system --timeout=60s
kubectl rollout status deployment/ani-gateway -n ani-system --timeout=60s
# 预期: deployment "xxx" successfully rolled out
```

---

## 七、验证部署成功

```bash
# 7.1 确认 Pod 状态（两个都应该是 1/1 Running）
kubectl get pods -n ani-system -l 'app.kubernetes.io/name in (ani-auth-service,ani-gateway)'
# 预期:
#   NAME                               READY   STATUS    RESTARTS   AGE
#   ani-auth-service-xxx-xxx           1/1     Running   0          xm
#   ani-gateway-xxx-xxx                1/1     Running   0          xm

# 7.2 确认镜像版本
kubectl get deploy -n ani-system ani-gateway ani-auth-service \
  -o custom-columns=NAME:.metadata.name,IMAGE:.spec.template.spec.containers[0].image
# 预期:
#   NAME               IMAGE
#   ani-gateway        docker.changqingyun.cn/ani/ani-gateway:dev-20260826
#   ani-auth-service   docker.changqingyun.cn/ani/ani-auth-service:dev-20260826

# 7.3 确认环境变量（pilot 三件套）
kubectl get deploy ani-gateway -n ani-system -o jsonpath='{range .spec.template.spec.containers[0].env[*]}{.name}={.value}{"\n"}{end}' | grep -iE 'AUTH_MODE|GATEWAY_AUTHZ'
# 预期:
#   ANI_AUTH_MODE=auth_service
#   GATEWAY_AUTHZ_POLICY_MODE=pilot
#   GATEWAY_AUTHZ_PILOT_OPERATIONS=listQuotaMeta

# 7.4 健康检查（NodePort 30080）
curl -s http://10.10.1.66:30080/healthz
# 预期: 200
```

---

## 八、功能验证（pilot 新链路测试）

```bash
BASE="http://10.10.1.66:30080"

# 8.1 平台登录取 token
TOK=$(curl -s -X POST $BASE/api/v1/auth/platform/password/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"Correct@123"}' | \
  grep -oP '"access_token":"\K[^"]+')
echo "token length: ${#TOK}"   # 预期: 822

# 8.2 调用 quota-meta（pilot V2 新链路，应 200）
curl -s -o /dev/null -w "%{http_code}" -X GET $BASE/api/v1/admin/quota-meta \
  -H "Authorization: Bearer $TOK"
# 预期: 200

# 8.3 无凭证调 quota-meta（V2 应拒绝 401）
curl -s -o /dev/null -w "%{http_code}" -X GET $BASE/api/v1/admin/quota-meta
# 预期: 401

# 8.4 legacy 旧链路无效 token（应 401）
curl -s -o /dev/null -w "%{http_code}" -X GET $BASE/api/v1/svc/models \
  -H "Authorization: Bearer invalid.token.here"
# 预期: 401

# 8.5 public 路由直通（应 200）
curl -s -o /dev/null -w "%{http_code}" -X GET $BASE/api/v1/branding
# 预期: 200
```

---

## 九、恢复环境（可选）

如果不需要 pilot 鉴权模式，改回 dev 旁路认证：

```bash
kubectl set env deployment/ani-gateway -n ani-system ANI_AUTH_MODE=dev
kubectl rollout status deployment/ani-gateway -n ani-system --timeout=60s
```

---

## 附录：常见问题

| 问题 | 原因 | 解决 |
|---|---|---|
| docker build 时 Go 模块下载超时 | `proxy.golang.org` 被墙 | Dockerfile 添加 `ENV GOPROXY=https://goproxy.cn,direct` |
| docker build 时 go.work 引用模块缺失 | workspace 模式找不到其他模块 | Dockerfile 添加 `ENV GOWORK=off` |
| ani-gateway build 时 curl 下载 helm/vcluster/kubectl 超时 | GitHub/get.helm.sh 下载慢 | 预下载工具到 `ani-tools/`，Dockerfile 改用 `COPY` |
| kubectl set image 后 Pod 还是旧版本 | hostPath 挂载二进制覆盖了镜像 | 确认 Pod 实际配置（不是 deployment 的 last-applied） |
| pilot 模式不生效 | `ANI_AUTH_MODE=dev` 旁路了认证 | 改成 `ANI_AUTH_MODE=auth_service` |
