#!/usr/bin/env python3
"""Validate Sprint 8 CORE-OFFLINE-PACK-A offline package lock."""

from __future__ import annotations

import hashlib
import re
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_LOCK = ROOT / "deploy" / "offline" / "core-package-lock.yaml"
SHA256_RE = re.compile(r"^[0-9a-f]{64}$")
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")


def load_lock(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"offline package lock does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: offline package lock must be a mapping")
    data["_path"] = str(path)
    return data


def validate_lock(lock: dict[str, Any]) -> None:
    path = lock.get("_path", "<memory>")
    if lock.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    source_manifest = require_string(lock, "source_manifest", path)
    source_manifest_path = ROOT / source_manifest
    if not source_manifest_path.exists():
        raise SystemExit(f"{path}: source_manifest does not exist: {source_manifest}")
    source_manifest_sha = require_string(lock, "source_manifest_sha256", path)
    if not SHA256_RE.match(source_manifest_sha):
        raise SystemExit(f"{path}: source_manifest_sha256 must be a 64 character hex string")
    actual_source_manifest_sha = hashlib.sha256(source_manifest_path.read_bytes()).hexdigest()
    if source_manifest_sha != actual_source_manifest_sha:
        raise SystemExit(
            f"{path}: source_manifest_sha256 does not match {source_manifest}: "
            f"want {actual_source_manifest_sha}"
        )
    artifact = lock.get("artifact")
    if not isinstance(artifact, dict):
        raise SystemExit(f"{path}: artifact must be a mapping")
    artifact_name = require_string(artifact, "name", path)
    artifact_sha = artifact.get("sha256")
    if any(term in artifact_name.lower() for term in FORBIDDEN_TERMS):
        raise SystemExit(f"{path}: forbidden Services artifact in Core offline package lock")
    if not isinstance(artifact_sha, str) or not SHA256_RE.match(artifact_sha):
        raise SystemExit(f"{path}: artifact.sha256 must be a 64 character hex string")
    if len(set(artifact_sha)) == 1:
        raise SystemExit(f"{path}: artifact.sha256 must not be a placeholder checksum")
    if artifact_sha != source_manifest_sha:
        raise SystemExit(f"{path}: artifact.sha256 must match source_manifest_sha256 for the Sprint 9 contract package")
    verification = lock.get("verification")
    if not isinstance(verification, list) or not verification:
        raise SystemExit(f"{path}: verification must be a non-empty list")
    commands: set[str] = set()
    for step in verification:
        if not isinstance(step, dict):
            raise SystemExit(f"{path}: verification entries must be mappings")
        command = require_string(step, "command", path)
        if any(term in command.lower() for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services verification command in Core offline package lock")
        commands.add(command)
    if "make validate-core-offline" not in commands:
        raise SystemExit(f"{path}: verification must include make validate-core-offline")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    lock = load_lock(DEFAULT_LOCK)
    validate_lock(lock)
    print(f"Core offline package lock valid: {lock['artifact']['name']}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
