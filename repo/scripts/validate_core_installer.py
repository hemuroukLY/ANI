#!/usr/bin/env python3
"""Validate Sprint 7 Core installer profiles."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE_DIR = ROOT / "installer" / "ani-installer" / "profiles"
REQUIRED_MODES = {"baremetal", "vm", "existing-k8s"}
FORBIDDEN_COMPONENTS = {"model-service", "kb-service", "rag-engine", "inference-operator", "console", "boss"}
REQUIRED_PHASES = {"preflight", "render", "validate"}


def load_profiles(profile_dir: Path) -> list[dict[str, Any]]:
    if not profile_dir.exists():
        raise SystemExit(f"profile directory does not exist: {profile_dir}")
    profiles: list[dict[str, Any]] = []
    for path in sorted(profile_dir.glob("*.yaml")):
        data = yaml.safe_load(path.read_text(encoding="utf-8"))
        if not isinstance(data, dict):
            raise SystemExit(f"{path}: installer profile must be a mapping")
        data["_path"] = str(path)
        profiles.append(data)
    if not profiles:
        raise SystemExit(f"no installer profiles found in {profile_dir}")
    return profiles


def validate_profiles(profiles: list[dict[str, Any]]) -> None:
    modes = set()
    for profile in profiles:
        path = profile.get("_path", "<memory>")
        mode = require_string(profile, "mode", path)
        modes.add(mode)
        if profile.get("scope") != "core":
            raise SystemExit(f"{path}: scope must be core")
        components = require_string_list(profile, "core_components", path)
        for component in components:
            if component in FORBIDDEN_COMPONENTS:
                raise SystemExit(f"{path}: forbidden Services component in Core installer profile: {component}")
        if "ani-gateway" not in components:
            raise SystemExit(f"{path}: core_components must include ani-gateway")
        preflight = profile.get("preflight")
        if not isinstance(preflight, dict):
            raise SystemExit(f"{path}: preflight must be a mapping")
        required_commands = preflight.get("required_commands")
        if not isinstance(required_commands, list) or not all(isinstance(item, str) and item for item in required_commands):
            raise SystemExit(f"{path}: preflight.required_commands must be non-empty strings")
        plan = profile.get("plan")
        if not isinstance(plan, dict):
            raise SystemExit(f"{path}: plan must be a mapping")
        phases = plan.get("phases")
        if not isinstance(phases, list) or not all(isinstance(item, str) and item for item in phases):
            raise SystemExit(f"{path}: plan.phases must be non-empty strings")
        missing = REQUIRED_PHASES - set(phases)
        if missing:
            raise SystemExit(f"{path}: plan phases must include {', '.join(sorted(missing))}")
    if modes != REQUIRED_MODES:
        raise SystemExit(f"installer modes must be {sorted(REQUIRED_MODES)}, got {sorted(modes)}")


def require_string(profile: dict[str, Any], key: str, path: object) -> str:
    value = profile.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def require_string_list(profile: dict[str, Any], key: str, path: object) -> list[str]:
    value = profile.get(key)
    if not isinstance(value, list) or not all(isinstance(item, str) and item for item in value):
        raise SystemExit(f"{path}: {key} must be a list of non-empty strings")
    return list(value)


def main() -> int:
    profiles = load_profiles(DEFAULT_PROFILE_DIR)
    validate_profiles(profiles)
    print(f"core installer profiles valid: {len(profiles)} profiles")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
