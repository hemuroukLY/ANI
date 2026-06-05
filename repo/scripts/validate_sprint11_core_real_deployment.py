#!/usr/bin/env python3
"""Validate Sprint 11 Core real deployment validation profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "real-k8s-lab" / "sprint11-core-real-deployment.yaml"
FORBIDDEN_TERMS = ("model-service", "kb-service", "/ai/", "inference", "boss")
REQUIRED_GATES = {
    "sprint10-release-prep",
    "real-k8s-profile",
    "k8s-read-only-inventory",
    "kubevirt-vm-readiness",
    "storageclass-readiness",
    "storage-disk-plan",
    "rook-ceph-formal-deployment",
    "sprint11-doc-consistency",
    "sprint11-safe-completion",
}
REQUIRED_BOUNDARIES = {
    "sprint11_may_run_read_only_inventory",
    "sprint11_may_validate_existing_core_contracts",
    "sprint11_must_not_partition_format_mount_or_edit_fstab",
    "sprint11_must_not_install_rook_ceph_without_disk_plan_approval",
    "sprint11_must_not_claim_production_ready_or_actual_release",
    "sprint11_safe_completion_requires_no_server_mutation",
}


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 real deployment profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 real deployment profile must be a mapping")
    data["_path"] = str(path)
    return data


def validate_profile(profile: dict[str, Any]) -> None:
    path = profile.get("_path", "<memory>")
    if profile.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if profile.get("version") != "sprint11":
        raise SystemExit(f"{path}: version must be sprint11")
    if profile.get("phase") != "core-real-deployment-validation":
        raise SystemExit(f"{path}: phase must be core-real-deployment-validation")
    if profile.get("real_deployment_validation_started") is not True:
        raise SystemExit(f"{path}: real_deployment_validation_started must be true")
    if profile.get("real_deployment_validation_complete") is not True:
        raise SystemExit(f"{path}: real_deployment_validation_complete must be true after safe completion")
    if profile.get("formal_deployment_code_complete") is not True:
        raise SystemExit(f"{path}: formal_deployment_code_complete must be true")
    for flag in (
        "mutation_allowed",
        "destructive_operations_allowed",
        "real_server_mutation_performed",
        "actual_release",
        "release_candidate",
        "production_release",
    ):
        if profile.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false")

    observed = require_mapping(profile, "cluster_observed", path)
    if observed.get("nodes_ready") != 3:
        raise SystemExit(f"{path}: cluster_observed.nodes_ready must be 3")
    if observed.get("kubevirt_phase") != "Deployed":
        raise SystemExit(f"{path}: cluster_observed.kubevirt_phase must be Deployed")
    if observed.get("storageclass_count") != 0:
        raise SystemExit(f"{path}: Sprint 11 kickoff must record the missing StorageClass state")
    if observed.get("rook_ceph_namespace_present") is not False:
        raise SystemExit(f"{path}: Rook-Ceph must not be marked present before install approval")

    boundaries = set(require_list(profile, "validation_boundaries", path))
    missing_boundaries = REQUIRED_BOUNDARIES - boundaries
    if missing_boundaries:
        raise SystemExit(f"{path}: missing validation boundaries: {', '.join(sorted(missing_boundaries))}")

    gates = require_list(profile, "gates", path)
    seen: set[str] = set()
    for gate in gates:
        if not isinstance(gate, dict):
            raise SystemExit(f"{path}: gate entries must be mappings")
        gate_id = require_string(gate, "id", path)
        category = require_string(gate, "category", path)
        command = require_string(gate, "command", path)
        combined = f"{gate_id} {category} {command}".lower()
        if any(term in combined for term in FORBIDDEN_TERMS):
            raise SystemExit(f"{path}: forbidden Services gate in Sprint 11 Core profile: {gate_id}")
        if "m1-real-lab-" in combined:
            raise SystemExit(f"{path}: Sprint 11 must not add REAL-K8S-LAB guard gates")
        if gate_id in seen:
            raise SystemExit(f"{path}: duplicate gate id: {gate_id}")
        seen.add(gate_id)
    missing_gates = REQUIRED_GATES - seen
    if missing_gates:
        raise SystemExit(f"{path}: missing Sprint 11 gates: {', '.join(sorted(missing_gates))}")

    approvals = set(require_list(profile, "next_manual_approval_required_for", path))
    for action in ("wipefs", "sgdisk", "mkfs", "mount", "fstab_update", "rook_ceph_osd_claim"):
        if action not in approvals:
            raise SystemExit(f"{path}: {action} must require manual approval")


def require_mapping(mapping: dict[str, Any], key: str, path: object) -> dict[str, Any]:
    value = mapping.get(key)
    if not isinstance(value, dict):
        raise SystemExit(f"{path}: {key} must be a mapping")
    return value


def require_list(mapping: dict[str, Any], key: str, path: object) -> list[Any]:
    value = mapping.get(key)
    if not isinstance(value, list) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty list")
    return value


def require_string(mapping: dict[str, Any], key: str, path: object) -> str:
    value = mapping.get(key)
    if not isinstance(value, str) or not value:
        raise SystemExit(f"{path}: {key} must be a non-empty string")
    return value


def main() -> int:
    profile = load_profile(DEFAULT_PROFILE)
    validate_profile(profile)
    print(f"Sprint 11 Core real deployment validation profile valid: {len(profile['gates'])} gates")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
