#!/usr/bin/env python3
"""Validate Sprint 11 Core safe completion profile."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PROFILE = ROOT / "deploy" / "real-k8s-lab" / "sprint11-core-safe-completion.yaml"
REQUIRED_PROFILES = {
    "deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml",
    "deploy/real-k8s-lab/sprint11-core-real-deployment.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml",
}
REQUIRED_GATES = {
    "sprint11-storage-disk-plan",
    "sprint11-core-real-deployment",
    "sprint11-rook-ceph-formal-deployment",
    "sprint11-doc-consistency",
    "sprint11-safe-completion",
    "sprint10-release-prep",
}
REQUIRED_PRINCIPLES = {
    "upstream_kubernetes_apis_first",
    "rook_ceph_raw_unmounted_osd_devices",
    "persistent_device_identity_over_kernel_sd_letter",
    "fail_closed_before_destructive_disk_operations",
    "reproducible_read_only_validation_before_mutation",
    "document_human_approval_before_state_change",
}
REQUIRED_APPROVALS = {
    "wipefs",
    "sgdisk",
    "mkfs",
    "mount",
    "fstab_update",
    "reboot_for_storage",
    "rook_ceph_cluster_install",
    "rook_ceph_osd_claim",
    "storageclass_default_switch",
    "existing_pv_migration",
}
REQUIRED_ALLOWED_OPERATION_CLASSES = {
    "local_contract_validation",
    "read_only_real_environment_inventory",
    "documentation_consistency_validation",
}
REQUIRED_FORBIDDEN_OPERATION_CLASSES = {
    "disk_partition_or_filesystem_mutation",
    "data_disk_mount_or_fstab_mutation",
    "rook_ceph_install_or_osd_claim",
    "storageclass_default_switch",
    "server_reboot",
}
FORBIDDEN_TRUE_FLAGS = {
    "real_storage_deployment_complete",
    "rook_ceph_installed",
    "rook_ceph_osd_claimed",
    "data_disk_mounting_complete",
    "actual_release",
    "release_candidate",
    "production_release",
    "mutation_allowed",
    "destructive_operations_allowed",
    "real_server_mutation_performed",
    "server_reboot_performed",
    "data_loss_risk_accepted",
}


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 safe completion profile does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 safe completion profile must be a mapping")
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
    if profile.get("safe_completion_complete") is not True:
        raise SystemExit(f"{path}: safe_completion_complete must be true")
    if profile.get("sprint_final_safe_complete") is not True:
        raise SystemExit(f"{path}: sprint_final_safe_complete must be true")
    if profile.get("real_deployment_validation_safe_complete") is not True:
        raise SystemExit(f"{path}: real_deployment_validation_safe_complete must be true")

    execution = require_mapping(profile, "execution_environment", path)
    if execution.get("entered") is not True:
        raise SystemExit(f"{path}: execution_environment.entered must be true")
    if execution.get("name") != "sprint11-core-real-deployment-validation":
        raise SystemExit(f"{path}: execution_environment.name must be sprint11-core-real-deployment-validation")
    if execution.get("scope") != "read_only_validation":
        raise SystemExit(f"{path}: execution_environment.scope must be read_only_validation")
    if execution.get("data_safety_prerequisite_enforced") is not True:
        raise SystemExit(f"{path}: execution_environment.data_safety_prerequisite_enforced must be true")
    for key in (
        "server_mutation_allowed",
        "storage_mutation_allowed",
        "production_storage_deployment_allowed",
    ):
        if execution.get(key) is not False:
            raise SystemExit(f"{path}: execution_environment.{key} must be false")
    allowed_operation_classes = set(require_list(execution, "allowed_operation_classes", path))
    missing_allowed = REQUIRED_ALLOWED_OPERATION_CLASSES - allowed_operation_classes
    if missing_allowed:
        raise SystemExit(
            f"{path}: missing allowed execution operation classes: {', '.join(sorted(missing_allowed))}"
        )
    forbidden_operation_classes = set(require_list(execution, "forbidden_operation_classes", path))
    missing_forbidden = REQUIRED_FORBIDDEN_OPERATION_CLASSES - forbidden_operation_classes
    if missing_forbidden:
        raise SystemExit(
            f"{path}: missing forbidden execution operation classes: {', '.join(sorted(missing_forbidden))}"
        )

    for flag in sorted(FORBIDDEN_TRUE_FLAGS):
        if profile.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false for safe completion")

    principles = set(require_list(profile, "open_source_principles", path))
    missing_principles = REQUIRED_PRINCIPLES - principles
    if missing_principles:
        raise SystemExit(f"{path}: missing open-source principles: {', '.join(sorted(missing_principles))}")

    evidence = require_mapping(profile, "safe_completion_evidence", path)
    required_true_evidence = (
        "read_only_inventory_performed",
        "data_disks_unmounted",
        "system_disks_mounted",
        "stable_device_ids_recorded",
        "no_fstab_change",
        "no_partition_change",
        "no_filesystem_creation",
        "no_rook_ceph_change",
    )
    for key in required_true_evidence:
        if evidence.get(key) is not True:
            raise SystemExit(f"{path}: safe_completion_evidence.{key} must be true")
    if evidence.get("k8s_nodes_ready") != 3:
        raise SystemExit(f"{path}: safe_completion_evidence.k8s_nodes_ready must be 3")
    if evidence.get("kubevirt_phase") != "Deployed":
        raise SystemExit(f"{path}: safe_completion_evidence.kubevirt_phase must be Deployed")
    if evidence.get("storageclass_count") != 0:
        raise SystemExit(f"{path}: safe completion must preserve missing StorageClass as an explicit gap")

    profiles = set(require_list(profile, "required_profiles", path))
    missing_profiles = REQUIRED_PROFILES - profiles
    if missing_profiles:
        raise SystemExit(f"{path}: missing required profiles: {', '.join(sorted(missing_profiles))}")
    for rel_path in profiles:
        if not (ROOT / rel_path).exists():
            raise SystemExit(f"{path}: required profile does not exist: {rel_path}")

    gates = require_list(profile, "completion_gates", path)
    seen_gates: set[str] = set()
    for gate in gates:
        if not isinstance(gate, dict):
            raise SystemExit(f"{path}: completion gate entries must be mappings")
        gate_id = require_string(gate, "id", path)
        require_string(gate, "category", path)
        require_string(gate, "command", path)
        if gate_id in seen_gates:
            raise SystemExit(f"{path}: duplicate completion gate id: {gate_id}")
        seen_gates.add(gate_id)
    missing_gates = REQUIRED_GATES - seen_gates
    if missing_gates:
        raise SystemExit(f"{path}: missing completion gates: {', '.join(sorted(missing_gates))}")

    approvals = set(require_list(profile, "manual_approval_required_before", path))
    missing_approvals = REQUIRED_APPROVALS - approvals
    if missing_approvals:
        raise SystemExit(f"{path}: missing manual approval actions: {', '.join(sorted(missing_approvals))}")


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
    print(f"Sprint 11 Core safe completion profile valid: {len(profile['completion_gates'])} gates")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
