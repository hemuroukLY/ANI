#!/usr/bin/env python3
"""Validate Sprint 4 MOCK-A Core mock server coverage."""

from __future__ import annotations

from pathlib import Path
from typing import Any
import json

import yaml

from serve_core_mock import HTTP_METHODS, build_routes, fallback_operation_id, load_spec, mock_value, server_base_path


ROOT = Path(__file__).resolve().parents[1]
DOC_ROOT = ROOT.parent


def fail(message: str) -> None:
    raise SystemExit(f"mock server contract invalid: {message}")


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def spec_operations(spec: dict[str, Any]) -> list[tuple[str, str, str, dict[str, Any]]]:
    operations: list[tuple[str, str, str, dict[str, Any]]] = []
    for path, path_item in sorted(spec.get("paths", {}).items()):
        if not isinstance(path_item, dict):
            continue
        for method, operation in sorted(path_item.items()):
            if method not in HTTP_METHODS or not isinstance(operation, dict):
                continue
            operation_id = operation.get("operationId") or fallback_operation_id(method, path)
            operations.append((method.upper(), path, operation_id, operation))
    return operations


def validate_route_coverage(spec: dict[str, Any]) -> None:
    operations = spec_operations(spec)
    routes = build_routes(spec)
    route_keys = {(route.method, route.path_template, route.operation_id) for route in routes}
    operation_keys = {(method, path, operation_id) for method, path, operation_id, _ in operations}
    missing = sorted(operation_keys - route_keys)
    extra = sorted(route_keys - operation_keys)
    if missing:
        fail(f"mock routes missing operations: {missing}")
    if extra:
        fail(f"mock routes contain unknown operations: {extra}")
    base_path = server_base_path(spec)
    if base_path != "/api/v1":
        fail(f"mock server base path must be /api/v1, got {base_path}")


def validate_mock_bodies(spec: dict[str, Any]) -> None:
    routes = build_routes(spec)
    if len(routes) < 10:
        fail("mock server must cover Core API operations, got too few routes")
    for route in routes:
        if route.status_code == 204:
            continue
        try:
            json.dumps(route.body, ensure_ascii=False)
        except TypeError as err:
            fail(f"{route.operation_id} mock body is not JSON serializable: {err}")
        if route.body is None:
            fail(f"{route.operation_id} mock body must not be empty for {route.status_code}")


def validate_error_mock(spec: dict[str, Any]) -> None:
    error_schema = spec.get("components", {}).get("schemas", {}).get("ErrorResponse", {})
    body = mock_value(spec, error_schema, "mockError")
    for field in ("code", "message", "request_id"):
        if field not in body:
            fail(f"ErrorResponse mock missing {field}")


def validate_docs() -> None:
    required = {
        "ANI-DOCS-INDEX.md": read(DOC_ROOT / "ANI-DOCS-INDEX.md"),
        "ANI-06-开发计划.md": read(DOC_ROOT / "ANI-06-开发计划.md"),
        "CURRENT-SPRINT.md": read(ROOT / "CURRENT-SPRINT.md"),
        "development-records/README.md": read(ROOT / "development-records/README.md"),
    }
    for path, content in required.items():
        if "MOCK-A" not in content:
            fail(f"{path} must reference MOCK-A")
    if "validate-mock-a" not in required["CURRENT-SPRINT.md"]:
        fail("CURRENT-SPRINT.md must list validate-mock-a")


def main() -> None:
    spec = load_spec(ROOT / "api/openapi/v1.yaml")
    validate_route_coverage(spec)
    validate_mock_bodies(spec)
    validate_error_mock(spec)
    validate_docs()
    print("mock server contract valid")


if __name__ == "__main__":
    main()
