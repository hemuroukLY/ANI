#!/usr/bin/env python3
"""Validate Sprint 9 CORE-RELEASE-EVIDENCE-A release evidence manifest."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MANIFEST = ROOT / "deploy" / "release" / "core-release-evidence.yaml"
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")
FORBIDDEN_SECRET_TERMS = ("password", "secret", "token", "credential", "private-key")
REQUIRED_EVIDENCE = {
    "core-api-compatibility",
    "architecture-boundary",
    "doc-entrypoints",
    "sdk-beta",
    "sdk-mock-smoke",
    "sprint8-core-release",
    "core-offline-pack",
    "core-cli-version",
}


def load_manifest(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"release evidence manifest does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: release evidence manifest must be a mapping")
    data["_path"] = str(path)
    return data


def validate_manifest(manifest: dict[str, Any]) -> None:
    path = manifest.get("_path", "<memory>")
    if manifest.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if manifest.get("version") != "sprint9":
        raise SystemExit(f"{path}: version must be sprint9")
    if manifest.get("target_release") != "v1.0.0":
        raise SystemExit(f"{path}: target_release must be v1.0.0")
    if manifest.get("rc_readiness") is not True:
        raise SystemExit(f"{path}: rc_readiness must be true")
    if manifest.get("release_candidate") is not False:
        raise SystemExit(f"{path}: release_candidate must be false until actual RC cut")
    if manifest.get("production_release") is not False:
        raise SystemExit(f"{path}: production_release must be false")

    entries = manifest.get("evidence")
    if not isinstance(entries, list) or not entries:
        raise SystemExit(f"{path}: evidence must be a non-empty list")
    seen: set[str] = set()
    for entry in entries:
        if not isinstance(entry, dict):
            raise SystemExit(f"{path}: evidence entries must be mappings")
        evidence_id = require_string(entry, "id", path)
        command = require_string(entry, "command", path)
        artifact = require_string(entry, "artifact", path)
        combined = f"{evidence_id} {command} {artifact}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services evidence in Core release manifest: {evidence_id}")
        if any(term in combined for term in FORBIDDEN_SECRET_TERMS):
            raise SystemExit(f"{path}: forbidden secret-bearing evidence in release manifest: {evidence_id}")
        if evidence_id in seen:
            raise SystemExit(f"{path}: duplicate evidence id: {evidence_id}")
        seen.add(evidence_id)
    missing = REQUIRED_EVIDENCE - seen
    if missing:
        raise SystemExit(f"{path}: missing release evidence: {', '.join(sorted(missing))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    manifest = load_manifest(DEFAULT_MANIFEST)
    validate_manifest(manifest)
    print(f"Core release evidence manifest valid: {len(manifest['evidence'])} entries")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
