#!/usr/bin/env python3
"""Validate Sprint 10 CORE-ARTIFACT-MANIFEST-A Core artifact manifest."""

from __future__ import annotations

import hashlib
import re
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MANIFEST = ROOT / "deploy" / "release" / "core-artifacts.yaml"
SHA256_RE = re.compile(r"^[0-9a-f]{64}$")
FORBIDDEN_TERMS = ("model-service", "kb-service", "rag", "inference", "console", "boss")
REQUIRED_ARTIFACTS = {
    "core-openapi",
    "core-sdk-metadata",
    "core-cli-source",
    "core-offline-lock",
    "core-release-evidence",
}
REQUIRED_COMMANDS = {
    "make validate-core-api-compatibility",
    "make validate-sdk-beta",
    "make validate-core-cli",
    "make validate-core-offline-pack",
    "make validate-core-release-evidence",
}


def load_manifest(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Core artifact manifest does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Core artifact manifest must be a mapping")
    data["_path"] = str(path)
    return data


def validate_manifest(manifest: dict[str, Any]) -> None:
    path = manifest.get("_path", "<memory>")
    if manifest.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if manifest.get("version") != "sprint10":
        raise SystemExit(f"{path}: version must be sprint10")
    if manifest.get("target_release") != "v1.0.0":
        raise SystemExit(f"{path}: target_release must be v1.0.0")
    if manifest.get("actual_release") is not False:
        raise SystemExit(f"{path}: actual_release must be false until formal release approval")

    artifacts = manifest.get("artifacts")
    if not isinstance(artifacts, list) or not artifacts:
        raise SystemExit(f"{path}: artifacts must be a non-empty list")
    seen: set[str] = set()
    for artifact in artifacts:
        if not isinstance(artifact, dict):
            raise SystemExit(f"{path}: artifact entries must be mappings")
        artifact_id = require_string(artifact, "id", path)
        rel_path = require_string(artifact, "path", path)
        digest = require_string(artifact, "sha256", path)
        combined = f"{artifact_id} {rel_path}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services artifact in Core manifest: {artifact_id}")
        if artifact_id in seen:
            raise SystemExit(f"{path}: duplicate artifact id: {artifact_id}")
        seen.add(artifact_id)
        if not SHA256_RE.match(digest):
            raise SystemExit(f"{path}: artifact {artifact_id} sha256 must be a 64 character hex string")
        artifact_path = ROOT / rel_path
        if not artifact_path.exists() or not artifact_path.is_file():
            raise SystemExit(f"{path}: artifact path does not exist: {rel_path}")
        actual = hashlib.sha256(artifact_path.read_bytes()).hexdigest()
        if digest != actual:
            raise SystemExit(f"{path}: artifact {artifact_id} sha256 does not match {rel_path}: want {actual}")

    missing = REQUIRED_ARTIFACTS - seen
    if missing:
        raise SystemExit(f"{path}: missing Core artifacts: {', '.join(sorted(missing))}")

    verification = manifest.get("verification")
    if not isinstance(verification, list) or not verification:
        raise SystemExit(f"{path}: verification must be a non-empty list")
    commands = {require_string(step, "command", path) for step in verification if isinstance(step, dict)}
    missing_commands = REQUIRED_COMMANDS - commands
    if missing_commands:
        raise SystemExit(f"{path}: missing verification commands: {', '.join(sorted(missing_commands))}")


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    manifest = load_manifest(DEFAULT_MANIFEST)
    validate_manifest(manifest)
    print(f"Core artifact manifest valid: {len(manifest['artifacts'])} artifacts")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
