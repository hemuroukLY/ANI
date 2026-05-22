#!/usr/bin/env python3
"""Validate Sprint 4 DOC-API-A static API docs."""

from __future__ import annotations

from pathlib import Path
from typing import Any
import re

import yaml


ROOT = Path(__file__).resolve().parents[1]
DOC_ROOT = ROOT.parent
DOCS_DIR = ROOT / "docs/api"
HTTP_METHODS = {"get", "post", "put", "patch", "delete"}
SPECS = {
    "core": ROOT / "api/openapi/v1.yaml",
    "services": ROOT / "api/openapi/services/v1.yaml",
}


def fail(message: str) -> None:
    raise SystemExit(f"api docs contract invalid: {message}")


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def load_spec(path: Path) -> dict[str, Any]:
    with path.open(encoding="utf-8") as handle:
        data = yaml.safe_load(handle)
    if not isinstance(data, dict):
        fail(f"{path.relative_to(ROOT)} must parse to an object")
    return data


def fallback_operation_id(method: str, path: str) -> str:
    tokens = re.sub(r"[^a-zA-Z0-9]+", " ", path).title().replace(" ", "")
    return method.lower() + tokens


def collect_operations(spec: dict[str, Any]) -> list[tuple[str, str, str]]:
    operations: list[tuple[str, str, str]] = []
    for path, path_item in sorted(spec.get("paths", {}).items()):
        if not isinstance(path_item, dict):
            continue
        for method, operation in sorted(path_item.items()):
            if method not in HTTP_METHODS or not isinstance(operation, dict):
                continue
            operations.append((method.upper(), path, operation.get("operationId") or fallback_operation_id(method, path)))
    return operations


def validate_files() -> None:
    for relative in ("index.html", "core.html", "services.html"):
        path = DOCS_DIR / relative
        if not path.exists():
            fail(f"missing docs/api/{relative}; run make gen-api-docs")
        content = read(path)
        if "<html" not in content or "</html>" not in content:
            fail(f"docs/api/{relative} must be a complete HTML document")


def validate_layer_docs() -> None:
    for layer, spec_path in SPECS.items():
        spec = load_spec(spec_path)
        content = read(DOCS_DIR / f"{layer}.html")
        for method, path, operation_id in collect_operations(spec):
            for token in (method, path, operation_id):
                if token not in content:
                    fail(f"docs/api/{layer}.html missing {token}")
        for schema_name in spec.get("components", {}).get("schemas", {}).keys():
            if schema_name not in content:
                fail(f"docs/api/{layer}.html missing schema {schema_name}")
    core_content = read(DOCS_DIR / "core.html")
    for forbidden in ("/models", "/inference-services", "/knowledge-bases"):
        if forbidden in core_content:
            fail(f"Core API docs must not include Services path {forbidden}")


def validate_index() -> None:
    content = read(DOCS_DIR / "index.html")
    required = ["Core API", "Services API", "core.html", "services.html", "api/openapi/v1.yaml", "api/openapi/services/v1.yaml"]
    for token in required:
        if token not in content:
            fail(f"docs/api/index.html missing {token}")


def validate_docs_references() -> None:
    required = {
        "ANI-DOCS-INDEX.md": read(DOC_ROOT / "ANI-DOCS-INDEX.md"),
        "ANI-06-开发计划.md": read(DOC_ROOT / "ANI-06-开发计划.md"),
        "CURRENT-SPRINT.md": read(ROOT / "CURRENT-SPRINT.md"),
        "development-records/README.md": read(ROOT / "development-records/README.md"),
    }
    for path, content in required.items():
        if "DOC-API-A" not in content:
            fail(f"{path} must reference DOC-API-A")
    if "validate-doc-api" not in required["CURRENT-SPRINT.md"]:
        fail("CURRENT-SPRINT.md must list validate-doc-api")


def main() -> None:
    validate_files()
    validate_layer_docs()
    validate_index()
    validate_docs_references()
    print("api docs contract valid")


if __name__ == "__main__":
    main()
