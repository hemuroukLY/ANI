#!/usr/bin/env python3
"""Validate Sprint 8 CORE-HARDEN-A release hardening profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "release" / "core-hardening.yaml"
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")
REQUIRED_GATES = {
    "core-api-compatibility",
    "architecture-boundary",
    "doc-entrypoints",
    "sdk-beta",
    "sdk-mock-smoke",
    "core-cli",
}


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"release hardening profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: hardening profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
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
            raise SystemExit(f"{path}: forbidden Services gate in Core hardening profile: {gate_id}")
        if gate_id in seen:
            raise SystemExit(f"{path}: duplicate gate id: {gate_id}")
        seen.add(gate_id)
    missing = REQUIRED_GATES - seen
    if missing:
        raise SystemExit(f"{path}: missing hardening gates: {', '.join(sorted(missing))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Core release hardening profile valid: {len(profile['gates'])} gates")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
