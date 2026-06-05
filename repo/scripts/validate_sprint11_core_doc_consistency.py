#!/usr/bin/env python3
"""Validate Sprint 11 Core documentation and gate consistency."""

from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
SPRINT11_MARKER = "Sprint 11 / Core Real Deployment Validation 正式部署完成"
BOUNDARY_MARKER = "真实服务器只读验证已完成；Rook-Ceph 正式部署已完成"
EXECUTION_ENV_MARKER = "Sprint 11 执行环境：正式部署执行环境"
REQUIRED_MAKE_TARGETS = (
    "validate-sprint11-storage-disk-plan",
    "validate-sprint11-core-real-deployment",
    "validate-sprint11-rook-ceph-formal-deployment",
    "validate-sprint11-rook-ceph-live-deployment-result",
    "validate-sprint11-rook-ceph-vm-storage-smoke",
    "validate-sprint11-rook-ceph-reboot-resilience",
    "validate-sprint11-safe-completion",
    "validate-sprint11-core-doc-consistency",
    "validate-sprint11-real-deployment",
)
REQUIRED_RECORDS = (
    "SPRINT11-KICKOFF-A",
    "CORE-STORAGE-DISK-RISK-A",
    "CORE-REAL-DEPLOY-A",
    "CORE-ROOK-CEPH-FORMAL-DEPLOYMENT-A",
    "CORE-ROOK-CEPH-LIVE-DEPLOYMENT-A",
    "CORE-ROOK-CEPH-VM-STORAGE-SMOKE-A",
    "CORE-ROOK-CEPH-REBOOT-RESILIENCE-A",
    "CORE-SAFE-COMPLETION-A",
    "CORE-REAL-DEPLOY-DOC-CONSISTENCY-A",
    "SPRINT11-SAFE-CLOSURE-A",
)


def validate_workspace(root: Path) -> None:
    docs_index = read(root / "ANI-DOCS-INDEX.md")
    ani_06 = read(root / "ANI-06-开发计划.md")
    current = read(root / "repo" / "CURRENT-SPRINT.md")
    makefile = read(root / "repo" / "Makefile")
    records = read(root / "repo" / "development-records" / "README.md")
    repo_readme = read(root / "repo" / "README.md")

    for label, content in (
        ("ANI-DOCS-INDEX.md", docs_index),
        ("ANI-06-开发计划.md", ani_06),
        ("repo/CURRENT-SPRINT.md", current),
        ("repo/development-records/README.md", records),
        ("repo/README.md", repo_readme),
    ):
        if SPRINT11_MARKER not in content:
            raise SystemExit(f"{label}: missing Sprint 11 started marker")
        if BOUNDARY_MARKER not in content:
            raise SystemExit(f"{label}: missing Sprint 11 real deployment safety boundary marker")
        if EXECUTION_ENV_MARKER not in content:
            raise SystemExit(f"{label}: missing Sprint 11 execution environment marker")
    for target in REQUIRED_MAKE_TARGETS:
        if f"{target}:" not in makefile:
            raise SystemExit(f"repo/Makefile: missing Sprint 11 target {target}")
    for record in REQUIRED_RECORDS:
        if record not in records:
            raise SystemExit(f"repo/development-records/README.md: missing Sprint 11 record {record}")


def read(path: Path) -> str:
    if not path.exists():
        raise SystemExit(f"required document does not exist: {path}")
    return path.read_text(encoding="utf-8")


def main() -> int:
    validate_workspace(ROOT)
    print("Core Sprint 11 documentation consistency valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
