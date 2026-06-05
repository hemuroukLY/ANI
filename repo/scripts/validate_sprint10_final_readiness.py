#!/usr/bin/env python3
"""Validate Sprint 10 CORE-FINAL-READINESS-A release-prep readiness profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "release" / "sprint10-core-final-readiness.yaml"
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")
REQUIRED_GATES = {
    "core-api-compatibility",
    "architecture-boundary",
    "doc-entrypoints",
    "sdk-beta",
    "sdk-mock-smoke",
    "sprint9-rc",
    "core-artifact-manifest",
    "core-version-policy",
    "core-cli-build",
    "sprint10-doc-consistency",
}


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 10 final readiness profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 10 final readiness profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if profile.get("version") != "sprint10":
        raise SystemExit(f"{path}: version must be sprint10")
    if profile.get("target_release") != "v1.0.0":
        raise SystemExit(f"{path}: target_release must be v1.0.0")
    if profile.get("release_preparation_complete") is not True:
        raise SystemExit(f"{path}: release_preparation_complete must be true")
    for flag in ("actual_release", "release_candidate", "production_release"):
        if profile.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false before formal release approval")

    gates = profile.get("gates")
    if not isinstance(gates, list) or not gates:
        raise SystemExit(f"{path}: gates must be a non-empty list")
    seen: set[str] = set()
    for gate in gates:
        if not isinstance(gate, dict):
            raise SystemExit(f"{path}: gate entries must be mappings")
        gate_id = require_string(gate, "id", path)
        command = require_string(gate, "command", path)
        category = require_string(gate, "category", path)
        combined = f"{gate_id} {command} {category}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services gate in Sprint 10 Core readiness profile: {gate_id}")
        if "m1-real-lab-" in combined:
            raise SystemExit(f"{path}: Sprint 10 must not add REAL-K8S-LAB guard gates")
        if gate_id in seen:
            raise SystemExit(f"{path}: duplicate gate id: {gate_id}")
        seen.add(gate_id)
    missing = REQUIRED_GATES - seen
    if missing:
        raise SystemExit(f"{path}: missing Sprint 10 readiness gates: {', '.join(sorted(missing))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Sprint 10 Core final readiness profile valid: {len(profile['gates'])} gates")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
