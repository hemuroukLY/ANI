#!/usr/bin/env python3
"""Validate Sprint 13 S01-S04 B-track production-shaped evidence boundaries."""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
RECORD_ROOT = ROOT / "development-records"

SLICES = {
    "S01": {
        "evidence": RECORD_ROOT / "live-evidence/sprint13-netroute-kubeovn-live-evidence.json",
        "result": RECORD_ROOT / "sprint13-netroute-kubeovn-live-result.md",
        "required_missing": {
            "production_rbac_and_credential_management",
            "persistent_route_metadata_reconciliation",
        },
    },
    "S02": {
        "evidence": RECORD_ROOT / "live-evidence/sprint13-k8s-workloads-vcluster-live-evidence.json",
        "result": RECORD_ROOT / "sprint13-k8s-workloads-vcluster-live-result.md",
        "required_missing": {
            "production_per_cluster_metadata_target",
            "production_tls_and_token_management",
        },
    },
    "S03": {
        "evidence": RECORD_ROOT / "live-evidence/sprint13-storage-rook-ceph-live-evidence.json",
        "result": RECORD_ROOT / "sprint13-storage-rook-ceph-live-result.md",
        "required_missing": {
            "production_serviceaccount_rbac",
            "tenant_storage_lifecycle_and_backup_restore",
        },
    },
    "S04": {
        "evidence": RECORD_ROOT / "live-evidence/sprint13-gpu-inventory-dcgm-live-evidence.json",
        "result": RECORD_ROOT / "sprint13-gpu-inventory-dcgm-live-result.md",
        "required_missing": {
            "production_in_cluster_kubernetes_api",
            "production_dcgm_service_or_prometheus_query",
        },
    },
}

ALLOWED_PRODUCTION_STATUSES = {"pending", "passed"}
PRODUCTION_FORBIDDEN_TRANSPORTS = {"lab_proxy", "local_port_forward", "dev_gateway"}


def fail(message: str) -> None:
    raise SystemExit(f"sprint13 production-shape guard invalid: {message}")


def load_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        fail(f"missing evidence {path.relative_to(ROOT)}")
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(f"malformed evidence {path.relative_to(ROOT)}: {exc}")
    if not isinstance(payload, dict):
        fail(f"evidence {path.relative_to(ROOT)} must be a JSON object")
    return payload


def validate_evidence(slice_id: str, path: Path) -> None:
    payload = load_json(path)
    if payload.get("status") != "passed":
        fail(f"{slice_id} evidence status must remain passed for real-provider gate")

    shape = payload.get("production_shape")
    if not isinstance(shape, dict):
        fail(f"{slice_id} evidence must include production_shape")

    status = shape.get("status")
    if status not in ALLOWED_PRODUCTION_STATUSES:
        fail(f"{slice_id} production_shape.status must be pending or passed")

    transport = shape.get("transport_profile")
    if not isinstance(transport, str) or not transport.strip():
        fail(f"{slice_id} production_shape.transport_profile must be non-empty")

    missing_items = shape.get("missing_items")
    if not isinstance(missing_items, list):
        if status == "pending":
            fail(f"{slice_id} pending production_shape must list missing_items")
        fail(f"{slice_id} production_shape.missing_items must be a list")
    missing_set = {str(item).strip() for item in missing_items if str(item).strip()}

    if status == "pending":
        if not missing_set:
            fail(f"{slice_id} pending production_shape must list missing_items")
        required = SLICES.get(slice_id, {}).get("required_missing", set())
        absent = set(required) - missing_set
        if absent:
            fail(f"{slice_id} production_shape missing_items must include {', '.join(sorted(absent))}")
        return

    if transport in PRODUCTION_FORBIDDEN_TRANSPORTS:
        fail(f"{slice_id} production_shape passed cannot use {transport}")
    if missing_set:
        fail(f"{slice_id} production_shape passed must not list missing_items")


def validate_result_doc(slice_id: str, path: Path) -> None:
    if not path.exists():
        fail(f"missing result doc {path.relative_to(ROOT)}")
    content = path.read_text(encoding="utf-8")
    required_tokens = ["Production-shaped gate", "production_shape"]
    for token in required_tokens:
        if token not in content:
            fail(f"{slice_id} result doc must reference {token}")
    if "not production ready" not in content and "不代表 production ready" not in content:
        fail(f"{slice_id} result doc must state not production ready")


def validate_all() -> None:
    for slice_id, spec in SLICES.items():
        validate_evidence(slice_id, spec["evidence"])
        validate_result_doc(slice_id, spec["result"])


def main() -> int:
    validate_all()
    print("Sprint 13 S01-S04 production-shaped evidence boundaries valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
