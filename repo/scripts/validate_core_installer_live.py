#!/usr/bin/env python3
"""Validate Sprint 8 CORE-INSTALLER-LIVE-A installer live readiness profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "real-k8s-lab" / "sprint8-core-installer-live.yaml"
REQUIRED_MODES = {"baremetal", "vm", "existing-k8s"}
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"installer live profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: installer live profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if profile.get("production_ready") is True:
        raise SystemExit(f"{path}: installer live profile must not claim production ready")
    if profile.get("evidence_level") not in {"contract-local", "live-ready"}:
        raise SystemExit(f"{path}: evidence_level must be contract-local or live-ready")
    targets = profile.get("targets")
    if not isinstance(targets, list) or not targets:
        raise SystemExit(f"{path}: targets must be a non-empty list")
    modes: set[str] = set()
    for target in targets:
        if not isinstance(target, dict):
            raise SystemExit(f"{path}: target entries must be mappings")
        target_id = require_string(target, "id", path)
        mode = require_string(target, "mode", path)
        profile_path = require_string(target, "profile", path)
        command = require_string(target, "command", path)
        combined = f"{target_id} {mode} {profile_path} {command}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services target in Core installer live profile: {target_id}")
        if not (ROOT / profile_path).exists():
            raise SystemExit(f"{path}: target profile path does not exist: {profile_path}")
        modes.add(mode)
    if modes != REQUIRED_MODES:
        raise SystemExit(f"{path}: installer live targets must cover {sorted(REQUIRED_MODES)}, got {sorted(modes)}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Core installer live readiness profile valid: {len(profile['targets'])} targets")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
