# SPRINT5-KICKOFF-A · Sprint 5 启动与执行入口切换

- 完成日期：2026-05-20
- 批次：SPRINT5-KICKOFF-A
- 所属 Sprint：Sprint 5（K8s 集群 + 控制器 + 加解密）

## 目标

1. 在 Sprint 4 已提交后的前提下，把执行入口从 Sprint 4 收尾切换为 Sprint 5。
2. 保持 `ANI-DOCS-INDEX.md`、`repo/CURRENT-SPRINT.md`、`ANI-06-开发计划.md` 三份入口文档状态一致。
3. 明确 Sprint 5 的首个开发切片与验收命令，确保后续批次可持续按“开发-测试-文档”闭环推进。

## 本批次完成内容

1. 将 `repo/CURRENT-SPRINT.md` 切换为 Sprint 5 执行入口，给出 `M1-K8S-A`、`M1-ENCRYPT-A` 两个 P0 批次最小切片和验收命令。
2. 更新 `ANI-DOCS-INDEX.md` 的“当前结论”和“推荐阅读路径”，把“下一步入口”显式指向 Sprint 5。
3. 更新 `ANI-06-开发计划.md` Section 零，将“当前阶段”切换为 Sprint 5 执行中并保持 Sprint 4 已完成的历史事实。
4. 更新 `development-records/README.md` 增加 `SPRINT5-KICKOFF-A` 归档索引条目。

## 验证命令

```bash
python scripts/validate_yaml.py api/openapi/v1.yaml api/openapi/services/v1.yaml
make test
make validate-architecture
git diff --check
```

## 备注

本批次只进行执行入口和文档一致性切换，不提前实现 Sprint 5 业务能力代码；业务能力按 `M1-K8S-A` 与 `M1-ENCRYPT-A` 分批次开发并逐批验证。
