#!/usr/bin/env python3
"""Validate Sprint 7 Core regression gate profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "real-k8s-lab" / "sprint7-core-regression.yaml"


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 7 regression profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: regression profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    targets = profile.get("targets")
    if not isinstance(targets, list) or not targets:
        raise SystemExit(f"{path}: targets must be a non-empty list")
    seen = set()
    for target in targets:
        if not isinstance(target, dict):
            raise SystemExit(f"{path}: targets entries must be mappings")
        target_id = require_string(target, "id", path)
        command = require_string(target, "command", path)
        category = str(target.get("category", ""))
        if target_id in seen:
            raise SystemExit(f"{path}: duplicate target id: {target_id}")
        seen.add(target_id)
        if "M1-REAL-LAB" in target_id or "M1-REAL-LAB" in command or "M1-REAL-LAB" in category:
            raise SystemExit(f"{path}: Sprint 7 regression must not add M1-REAL-LAB guard targets")
    required = {
        "core-installer-contract",
        "core-offline-contract",
        "core-cli-contract",
        "real-k8s-profile-contract",
        "sdk-mock-smoke",
    }
    missing = required - seen
    if missing:
        raise SystemExit(f"{path}: missing regression targets: {', '.join(sorted(missing))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Sprint 7 Core regression profile valid: {len(profile['targets'])} targets")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
