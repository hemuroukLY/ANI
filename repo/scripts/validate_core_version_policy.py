#!/usr/bin/env python3
"""Validate Sprint 10 CORE-VERSION-POLICY-A version policy guard."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[2]
REPO_ROOT = ROOT / "repo"
DEFAULT_POLICY = REPO_ROOT / "deploy" / "release" / "core-version-policy.yaml"
FORBIDDEN_TRUE_FLAGS = (
    "actual_release",
    "release_candidate",
    "production_release",
    "rc_cut",
)
REQUIRED_DOC_MARKER = "不是实际 v1.0.0 发布"


def load_policy(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Core version policy does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Core version policy must be a mapping")
    data["_path"] = str(path)
    return data


def validate_policy(policy: dict[str, Any], root: Path = ROOT) -> None:
    path = policy.get("_path", "<memory>")
    if policy.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if policy.get("version") != "sprint10":
        raise SystemExit(f"{path}: version must be sprint10")
    if policy.get("target_release") != "v1.0.0":
        raise SystemExit(f"{path}: target_release must be v1.0.0")
    for flag in FORBIDDEN_TRUE_FLAGS:
        if policy.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false before formal release approval")

    release_dir = root / "repo" / "deploy" / "release"
    for release_file in release_dir.glob("*.yaml"):
        data = yaml.safe_load(release_file.read_text(encoding="utf-8"))
        if not isinstance(data, dict):
            continue
        for flag in FORBIDDEN_TRUE_FLAGS:
            if data.get(flag) is True:
                raise SystemExit(f"{release_file}: {flag} must not be true before formal release approval")
        if data.get("released_version") == "v1.0.0":
            raise SystemExit(f"{release_file}: released_version must not be v1.0.0 before formal release approval")

    for doc in (
        root / "ANI-DOCS-INDEX.md",
        root / "ANI-06-开发计划.md",
        root / "repo" / "CURRENT-SPRINT.md",
    ):
        content = read(doc)
        if REQUIRED_DOC_MARKER not in content:
            raise SystemExit(f"{doc}: missing Sprint 10 actual release boundary marker")


def read(path: Path) -> str:
    if not path.exists():
        raise SystemExit(f"required document does not exist: {path}")
    return path.read_text(encoding="utf-8")


def main() -> int:
    policy = load_policy(DEFAULT_POLICY)
    validate_policy(policy)
    print("Core version policy valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
