#!/usr/bin/env python3
"""Validate Sprint 7 Core offline package manifest."""

from __future__ import annotations

import re
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MANIFEST = ROOT / "deploy" / "offline" / "core-package.yaml"
FORBIDDEN_IMAGES = {"model-service", "kb-service", "rag-engine", "inference-operator", "console", "boss"}
DIGEST_RE = re.compile(r"@sha256:[0-9a-f]{64}$")


def load_manifest(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"offline manifest does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: offline manifest must be a mapping")
    data["_path"] = str(path)
    return data


def validate_manifest(manifest: dict[str, Any]) -> None:
    path = manifest.get("_path", "<memory>")
    if manifest.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    images = require_mapping_list(manifest, "images", path)
    image_names = set()
    for image in images:
        name = require_string(image, "name", path)
        ref = require_string(image, "image", path)
        image_names.add(name)
        if name in FORBIDDEN_IMAGES:
            raise SystemExit(f"{path}: forbidden Services image in Core offline package: {name}")
        if not DIGEST_RE.search(ref):
            raise SystemExit(f"{path}: image {name} must be pinned by digest")
    if "ani-gateway" not in image_names:
        raise SystemExit(f"{path}: images must include ani-gateway")
    for key in ("charts", "scripts"):
        entries = require_mapping_list(manifest, key, path)
        for entry in entries:
            rel_path = require_string(entry, "path", path)
            if not (ROOT / rel_path).exists():
                raise SystemExit(f"{path}: {key} path does not exist: {rel_path}")


def require_mapping_list(manifest: dict[str, Any], key: str, path: object) -> list[dict[str, Any]]:
    value = manifest.get(key)
    if not isinstance(value, list) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty list")
    for item in value:
        if not isinstance(item, dict):
            raise SystemExit(f"{path}: {key} entries must be mappings")
    return list(value)


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    manifest = load_manifest(DEFAULT_MANIFEST)
    validate_manifest(manifest)
    print(f"core offline manifest valid: {len(manifest['images'])} images")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
