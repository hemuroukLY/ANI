# 2026-05-12 Demo Handoff

## Status

Paused for the day.

The staged Demo console is ready for local presentation use. It is a
presentation layer and frontend experience reference, not a replacement for the
formal ANI product roadmap.

## Demo Boundary

- Real M1-backed flows:
  - VM/container/GPU container create
  - VM/container/GPU container list/get
  - VM/container/GPU container lifecycle operations
  - VM console session page
  - Demo VM shell execution in a scoped backend workspace
- Mock-only presentation flows:
  - dashboard
  - model management
  - knowledge base
  - usage report
  - settings
  - image/network/storage/GPU pool/console session/task/audit resource-domain
    pages

Mock pages exist to make the demo clickable and to validate console information
architecture early. They must not be treated as production backend completion.

## Start Tomorrow

Terminal 1:

```bash
cd ~/ANI/repo
ANI_AUTH_MODE=dev ./bin/ani-gateway
```

Terminal 2:

```bash
cd ~/ANI/repo/frontends/console
npm run dev
```

Open:

```text
http://127.0.0.1:5173/
```

VM console opens from the instance workspace and uses:

```text
http://127.0.0.1:5173/demo-console?instance_id=<id>&protocol=vnc
```

## Presentation Checklist

1. Open the console at `http://127.0.0.1:5173/`.
2. Confirm the Evergreen/KuberCloud brand icon is visible in the dark top bar.
3. Navigate the left service menu:
   - 运营总览
   - 实例工作台
   - 模型管理
   - 知识库
   - 用量报表
   - 平台设置
4. In 实例工作台, create a VM.
5. Run VM lifecycle actions:
   - 启动
   - 停止
   - 重启
   - 变配
   - 删除
6. Open VM 控制台 in a new page.
7. Run a simple read-only shell command in the Demo console, for example:

```bash
pwd
```

## Validation Commands

Run before a formal demo:

```bash
cd ~/ANI/repo
make validate-demo-instances
git diff --check
```

Run frontend build check:

```bash
cd ~/ANI/repo/frontends/console
npm run build
```

Optional local HTTP probes:

```bash
curl -s -i http://127.0.0.1:5173/
curl -s -i http://127.0.0.1:8080/health
curl -s -i http://127.0.0.1:5173/assets/brand/d-logo.png
```

If port `8080` is occupied by a stale gateway process:

```bash
pkill -f ani-gateway
```

Then restart Terminal 1.

## Guardrails For Future Development

- The Demo must not change the formal M1/M2/M3 development order.
- Mock data must stay in frontend/demo presentation surfaces unless promoted by
  a formal implementation slice.
- Real instance behavior must continue to go through `WorkloadInstanceService`
  and `WorkloadInstanceOps`.
- Provider details must stay behind M1 ports/adapters.
- Future production VM console integration should replace the Demo shell
  backend with noVNC/WebSocket or cloud-provider console proxy without changing
  the business-layer contract.

## Next Development Options

Recommended order after resuming:

1. `M1-DEMO-SMOKE-A` if the next milestone requires a real presentation
   environment that creates resources in Kubernetes/KubeVirt/GPU cluster.
2. `M3-MODEL-A` if the M1/M2/GPU/Runtime/Instance baseline is accepted and the
   product work should move to model metadata and object-storage boundaries.
