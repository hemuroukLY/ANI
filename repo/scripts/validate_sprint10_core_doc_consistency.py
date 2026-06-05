#!/usr/bin/env python3
"""Validate Sprint 10 Core documentation and gate consistency."""

from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
SPRINT10_MARKER = "Sprint 10 Core-only 代码开发已完成"
BOUNDARY_MARKER = "不是实际 v1.0.0 发布"
STALE_CURRENT_MARKERS = (
    "当前重心：Sprint 9 Core-only 代码开发已完成",
    "下一步准备 Sprint 10",
    "下一步：准备 Sprint 10",
)
REQUIRED_MAKE_TARGETS = (
    "validate-core-artifact-manifest",
    "validate-core-version-policy",
    "validate-sprint10-final-readiness",
    "validate-sprint10-core-doc-consistency",
    "validate-sprint10-release-prep",
)
REQUIRED_RECORDS = (
    "CORE-ARTIFACT-MANIFEST-A",
    "CORE-VERSION-POLICY-A",
    "CORE-FINAL-READINESS-A",
    "CORE-CLI-RELEASE-METADATA-A",
    "CORE-FINAL-DOC-CONSISTENCY-A",
    "SPRINT10-CLOSURE-A",
)


def validate_workspace(root: Path) -> None:
    docs_index = read(root / "ANI-DOCS-INDEX.md")
    ani_06 = read(root / "ANI-06-开发计划.md")
    current = read(root / "repo" / "CURRENT-SPRINT.md")
    makefile = read(root / "repo" / "Makefile")
    records = read(root / "repo" / "development-records" / "README.md")

    for marker in STALE_CURRENT_MARKERS:
        if marker in current:
            raise SystemExit(f"repo/CURRENT-SPRINT.md: stale Sprint 9 current marker: {marker}")
    for label, content in (
        ("ANI-DOCS-INDEX.md", docs_index),
        ("ANI-06-开发计划.md", ani_06),
        ("repo/CURRENT-SPRINT.md", current),
        ("repo/development-records/README.md", records),
    ):
        if SPRINT10_MARKER not in content:
            raise SystemExit(f"{label}: missing Sprint 10 completed marker")
        if BOUNDARY_MARKER not in content:
            raise SystemExit(f"{label}: missing Sprint 10 release boundary marker")
    for target in REQUIRED_MAKE_TARGETS:
        if f"{target}:" not in makefile:
            raise SystemExit(f"repo/Makefile: missing Sprint 10 target {target}")
    for record in REQUIRED_RECORDS:
        if record not in records:
            raise SystemExit(f"repo/development-records/README.md: missing Sprint 10 record {record}")


def read(path: Path) -> str:
    if not path.exists():
        raise SystemExit(f"required document does not exist: {path}")
    return path.read_text(encoding="utf-8")


def main() -> int:
    validate_workspace(ROOT)
    print("Core Sprint 10 documentation consistency valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
