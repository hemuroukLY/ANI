#!/usr/bin/env python3
"""Validate Sprint 9 CORE-RC-GATE-A release-candidate readiness gate."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "release" / "sprint9-core-rc.yaml"
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")
REQUIRED_GATES = {
    "core-api-compatibility",
    "architecture-boundary",
    "doc-entrypoints",
    "sdk-beta",
    "sdk-mock-smoke",
    "sprint8-core-release",
    "core-release-evidence",
    "core-offline-pack",
    "core-cli-version",
    "sprint9-doc-consistency",
}


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 9 RC profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 9 RC profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if profile.get("version") != "sprint9":
        raise SystemExit(f"{path}: version must be sprint9")
    if profile.get("rc_readiness") is not True:
        raise SystemExit(f"{path}: rc_readiness must be true")
    if profile.get("release_candidate") is not False:
        raise SystemExit(f"{path}: release_candidate must be false until actual RC cut")
    if profile.get("production_release") is not False:
        raise SystemExit(f"{path}: production_release must be false")

    gates = profile.get("gates")
    if not isinstance(gates, list) or not gates:
        raise SystemExit(f"{path}: gates must be a non-empty list")
    seen: set[str] = set()
    for gate in gates:
        if not isinstance(gate, dict):
            raise SystemExit(f"{path}: gate entries must be mappings")
        gate_id = require_string(gate, "id", path)
        category = require_string(gate, "category", path)
        command = require_string(gate, "command", path)
        combined = f"{gate_id} {category} {command}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services gate in Sprint 9 Core RC profile: {gate_id}")
        if "m1-real-lab-" in combined:
            raise SystemExit(f"{path}: Sprint 9 must not add REAL-K8S-LAB guard gates")
        if gate_id in seen:
            raise SystemExit(f"{path}: duplicate gate id: {gate_id}")
        seen.add(gate_id)
    missing = REQUIRED_GATES - seen
    if missing:
        raise SystemExit(f"{path}: missing Sprint 9 RC gates: {', '.join(sorted(missing))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Sprint 9 Core RC readiness profile valid: {len(profile['gates'])} gates")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
