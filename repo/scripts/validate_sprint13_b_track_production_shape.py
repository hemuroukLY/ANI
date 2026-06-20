#!/usr/bin/env python3
"""Validate Sprint 13 S01-S04 B-track production-shaped evidence boundaries."""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
RECORD_ROOT = ROOT / "development-records"
PRODUCTION_PROFILE = ROOT / "deploy/real-k8s-lab/sprint13-production-shaped-gateway-profile.yaml"
PRODUCTION_RBAC = ROOT / "deploy/real-k8s-lab/sprint13-production-shaped-gateway-rbac.yaml"

SLICES = {
    "S01": {
        "evidence": RECORD_ROOT / "live-evidence/sprint13-netroute-kubeovn-live-evidence.json",
        "result": RECORD_ROOT / "sprint13-netroute-kubeovn-live-result.md",
        "required_missing": {
            "production_rbac_and_credential_management",
            "persistent_route_metadata_reconciliation",
        },
        "required_proof": {
            "production_gateway",
            "in_cluster_serviceaccount_rbac",
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
        "required_proof": {
            "production_gateway",
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
        "required_proof": {
            "production_gateway",
            "in_cluster_serviceaccount_rbac",
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
        "required_proof": {
            "production_gateway",
            "in_cluster_kubernetes_api",
            "production_dcgm_service_or_prometheus_query",
        },
    },
}

ALLOWED_PRODUCTION_STATUSES = {"pending", "passed"}
PRODUCTION_FORBIDDEN_TRANSPORT_TOKENS = {"lab", "local", "port_forward", "port-forward", "dev_gateway", "dev-gateway", "kubectl_proxy", "kubectl-proxy"}
REQUIRED_RBAC_KINDS = {"ServiceAccount", "ClusterRole", "ClusterRoleBinding"}
REQUIRED_RBAC_RESOURCES = {"nodes", "services", "persistentvolumeclaims", "networkpolicies", "vpcs", "subnets", "volumesnapshots"}
REQUIRED_STANDARD_SLICES = {"S01", "S02", "S03", "S04", "S05", "S06", "S07"}


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

    if any(token in transport for token in PRODUCTION_FORBIDDEN_TRANSPORT_TOKENS):
        fail(f"{slice_id} production_shape passed cannot use {transport}")
    if missing_set:
        fail(f"{slice_id} production_shape passed must not list missing_items")
    proof_items = shape.get("proof_items")
    if not isinstance(proof_items, list):
        fail(f"{slice_id} production_shape passed requires proof_items")
    proof_set = {str(item).strip() for item in proof_items if str(item).strip()}
    if not proof_set:
        fail(f"{slice_id} production_shape passed requires proof_items")
    required_proof = SLICES.get(slice_id, {}).get("required_proof", set())
    absent_proof = set(required_proof) - proof_set
    if absent_proof:
        fail(f"{slice_id} production_shape proof_items must include {', '.join(sorted(absent_proof))}")


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


def validate_production_profile() -> None:
    if not PRODUCTION_PROFILE.exists():
        fail(f"missing production profile {PRODUCTION_PROFILE.relative_to(ROOT)}")
    if not PRODUCTION_RBAC.exists():
        fail(f"missing production RBAC {PRODUCTION_RBAC.relative_to(ROOT)}")
    try:
        profile = yaml.safe_load(PRODUCTION_PROFILE.read_text(encoding="utf-8"))
    except yaml.YAMLError as exc:
        fail(f"malformed production profile {PRODUCTION_PROFILE.relative_to(ROOT)}: {exc}")
    if not isinstance(profile, dict):
        fail("production profile must be a YAML object")
    if profile.get("profile") != "SPRINT13-B-TRACK-PRODUCTION-SHAPED-GATEWAY":
        fail("production profile id must be SPRINT13-B-TRACK-PRODUCTION-SHAPED-GATEWAY")
    gateway = profile.get("gateway")
    if not isinstance(gateway, dict):
        fail("production profile must include gateway block")
    if gateway.get("deployment_mode") != "in_cluster":
        fail("production profile gateway deployment_mode must be in_cluster")
    kube_client = gateway.get("kubernetes_client")
    if not isinstance(kube_client, dict):
        fail("production profile must include gateway.kubernetes_client")
    expected_sources = {
        "host_source": "in_cluster_service",
        "token_source": "service_account_projected_token",
        "ca_source": "service_account_ca_bundle",
    }
    for field, expected in expected_sources.items():
        if kube_client.get(field) != expected:
            fail(f"production profile Kubernetes client {field} must be {expected}")
    proof_items = profile.get("slice_proof_items")
    if not isinstance(proof_items, dict):
        fail("production profile must include slice_proof_items")
    absent_slices = REQUIRED_STANDARD_SLICES - set(proof_items)
    if absent_slices:
        fail(f"production profile slice_proof_items missing {', '.join(sorted(absent_slices))}")
    for slice_id in sorted(REQUIRED_STANDARD_SLICES):
        items = proof_items.get(slice_id)
        if not isinstance(items, list) or "production_gateway" not in items:
            fail(f"production profile {slice_id} proof items must include production_gateway")

    try:
        docs = list(yaml.safe_load_all(PRODUCTION_RBAC.read_text(encoding="utf-8")))
    except yaml.YAMLError as exc:
        fail(f"malformed production RBAC {PRODUCTION_RBAC.relative_to(ROOT)}: {exc}")
    docs = [doc for doc in docs if isinstance(doc, dict)]
    kinds = {str(doc.get("kind", "")) for doc in docs}
    missing_kinds = REQUIRED_RBAC_KINDS - kinds
    if missing_kinds:
        fail(f"production RBAC missing {', '.join(sorted(missing_kinds))}")
    cluster_roles = [doc for doc in docs if doc.get("kind") == "ClusterRole"]
    if len(cluster_roles) != 1:
        fail("production RBAC must include exactly one ClusterRole")
    resources = set()
    for rule in cluster_roles[0].get("rules", []):
        if isinstance(rule, dict) and isinstance(rule.get("resources"), list):
            resources.update(str(resource) for resource in rule["resources"])
    missing_resources = REQUIRED_RBAC_RESOURCES - resources
    if missing_resources:
        fail(f"production RBAC missing resources {', '.join(sorted(missing_resources))}")
    if "*" in resources:
        fail("production RBAC must not grant wildcard resources")


def validate_all() -> None:
    validate_production_profile()
    for slice_id, spec in SLICES.items():
        validate_evidence(slice_id, spec["evidence"])
        validate_result_doc(slice_id, spec["result"])


def main() -> int:
    validate_all()
    print("Sprint 13 S01-S04 production-shaped evidence boundaries valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
