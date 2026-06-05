#!/usr/bin/env python3
"""Validate Sprint 11 Rook-Ceph reboot resilience result."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_RESULT = ROOT / "deploy" / "real-k8s-lab" / "sprint11-rook-ceph-reboot-resilience-result.yaml"
REQUIRED_SOURCE_PROFILES = {
    "deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-vm-storage-smoke-result.yaml",
}
EXPECTED_REBOOT_ORDER = ("dev-phys-02", "dev-phys-03", "kubercloud")


def load_result(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 reboot resilience result does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 reboot resilience result must be a mapping")
    data["_path"] = str(path)
    return data


def validate_result(result: dict[str, Any]) -> None:
    path = result.get("_path", "<memory>")
    expected_top = {
        "scope": "core",
        "version": "sprint11",
        "phase": "core-real-deployment-reboot-resilience-result",
        "batch": "CORE-ROOK-CEPH-REBOOT-RESILIENCE-A",
        "deployment_target": "rook-ceph-kubevirt-rbd-reboot-resilience",
    }
    for key, value in expected_top.items():
        if result.get(key) != value:
            raise SystemExit(f"{path}: {key} must be {value}")

    for flag in (
        "manual_approval_received",
        "data_safety_prerequisite_enforced",
        "live_reboot_resilience_executed",
        "real_server_mutation_performed",
        "server_reboot_performed",
    ):
        if result.get(flag) is not True:
            raise SystemExit(f"{path}: {flag} must be true")
    for flag in (
        "simultaneous_reboot_performed",
        "destructive_disk_operations_directly_invoked_by_agent",
        "system_disk_touched",
        "manual_mount_performed",
        "fstab_mutation_performed",
        "data_disk_mount_or_fstab_mutation",
        "storageclass_default_switched",
        "existing_pvc_migration_performed",
        "actual_release",
        "production_release",
    ):
        if result.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false")

    profiles = set(require_list(result, "source_profiles", path))
    missing_profiles = REQUIRED_SOURCE_PROFILES - profiles
    if missing_profiles:
        raise SystemExit(f"{path}: missing source profiles: {', '.join(sorted(missing_profiles))}")
    for rel_path in profiles:
        if not (ROOT / rel_path).exists():
            raise SystemExit(f"{path}: source profile does not exist: {rel_path}")

    validate_preflight(require_mapping(result, "preflight", path), path)
    validate_execution_policy(require_mapping(result, "execution_policy", path), path)
    validate_node_results(require_list(result, "node_reboot_results", path), path)
    validate_post_reboot_cluster(require_mapping(result, "post_reboot_cluster", path), path)
    validate_time_sync_remediation(require_mapping(result, "time_sync_remediation", path), path)
    validate_cleanup(require_mapping(result, "cleanup_result", path), path)
    validate_safety(require_mapping(result, "safety_result", path), path)


def validate_preflight(preflight: dict[str, Any], path: object) -> None:
    expected = {
        "kubernetes_nodes_ready": 3,
        "cephcluster_phase": "Ready",
        "ceph_health": "HEALTH_OK",
        "cephblockpool_phase": "Ready",
        "kubevirt_phase": "Deployed",
        "formal_storageclass": "ani-rbd-ssd",
        "formal_storageclass_reclaimPolicy": "Retain",
        "formal_storageclass_default": False,
    }
    for key, value in expected.items():
        if preflight.get(key) != value:
            raise SystemExit(f"{path}: preflight.{key} must be {value}")


def validate_execution_policy(policy: dict[str, Any], path: object) -> None:
    for flag in (
        "worker_nodes_first",
        "control_plane_last",
        "one_node_at_a_time",
        "stop_on_unhealthy_gate",
        "no_drain_performed",
        "no_mount_or_fstab_change",
    ):
        if policy.get(flag) is not True:
            raise SystemExit(f"{path}: execution_policy.{flag} must be true")


def validate_node_results(nodes: list[Any], path: object) -> None:
    if [node.get("node") for node in nodes if isinstance(node, dict)] != list(EXPECTED_REBOOT_ORDER):
        raise SystemExit(f"{path}: node_reboot_results must preserve worker-first control-plane-last order")
    for node in nodes:
        if not isinstance(node, dict):
            raise SystemExit(f"{path}: node_reboot_results entries must be mappings")
        name = node.get("node")
        if node.get("reboot_executed") is not True:
            raise SystemExit(f"{path}: node {name} reboot_executed must be true")
        if node.get("reboot_event_observed") is not True:
            raise SystemExit(f"{path}: node {name} reboot_event_observed must be true")
        if node.get("node_ready_after_reboot") is not True:
            raise SystemExit(f"{path}: node {name} node_ready_after_reboot must be true")
        if node.get("ceph_health_after_reboot") != "HEALTH_OK":
            raise SystemExit(f"{path}: node {name} ceph_health_after_reboot must be HEALTH_OK")
        if name in {"dev-phys-02", "dev-phys-03"}:
            validate_worker_node(node, path)
        elif name == "kubercloud":
            validate_control_plane_node(node, path)
        else:
            raise SystemExit(f"{path}: unexpected reboot node {name}")


def validate_worker_node(node: dict[str, Any], path: object) -> None:
    name = node.get("node")
    if node.get("role") != "worker":
        raise SystemExit(f"{path}: node {name} role must be worker")
    if node.get("osd_pods_on_node_recovered") != 2:
        raise SystemExit(f"{path}: node {name} must recover two OSD pods")
    for flag in (
        "vm_pvc_baseline_before_reboot",
        "vm_ready_after_reboot",
        "pvc_bound_after_reboot",
        "guest_rbd_block_device_visible_after_reboot",
        "guest_block_write_attempted_after_reboot",
    ):
        if node.get(flag) is not True:
            raise SystemExit(f"{path}: node {name} {flag} must be true")


def validate_control_plane_node(node: dict[str, Any], path: object) -> None:
    if node.get("role") != "control-plane":
        raise SystemExit(f"{path}: control-plane node role must be control-plane")
    expected = {
        "api_readyz_after_reboot": True,
        "mon_pod_on_node_recovered": True,
        "mgr_pod_on_node_recovered": True,
        "worker_vm_pvc_observable_after_control_plane_reboot": True,
        "worker_vm_ready_after_control_plane_reboot": True,
        "worker_pvc_bound_after_control_plane_reboot": True,
        "control_plane_recovery_required_extended_wait": True,
    }
    for key, value in expected.items():
        if node.get(key) != value:
            raise SystemExit(f"{path}: control-plane {key} must be {value}")
    if node.get("osd_pods_on_node_recovered") != 1:
        raise SystemExit(f"{path}: control-plane must recover one OSD pod")


def validate_post_reboot_cluster(cluster: dict[str, Any], path: object) -> None:
    expected = {
        "kubernetes_nodes_ready": 3,
        "cephcluster_phase": "Ready",
        "ceph_health": "HEALTH_OK",
        "clock_skew_warning_observed_after_reboot": True,
        "clock_skew_warning_cleared": True,
        "cephblockpool_phase": "Ready",
        "kubevirt_phase": "Deployed",
        "osd_count_running": 5,
        "formal_storageclass_reclaimPolicy": "Retain",
        "formal_storageclass_default": False,
    }
    for key, value in expected.items():
        if cluster.get(key) != value:
            raise SystemExit(f"{path}: post_reboot_cluster.{key} must be {value}")
    if require_mapping(cluster, "osd_topology", path) != {"kubercloud": 1, "dev-phys-02": 2, "dev-phys-03": 2}:
        raise SystemExit(f"{path}: post_reboot_cluster.osd_topology must match approved SSD OSD layout")
    if require_mapping(cluster, "osd_restart_observed_after_reboot", path) != {
        "kubercloud": 1,
        "dev-phys-02": 2,
        "dev-phys-03": 2,
    }:
        raise SystemExit(f"{path}: post_reboot_cluster.osd_restart_observed_after_reboot must reflect each node reboot")


def validate_time_sync_remediation(remediation: dict[str, Any], path: object) -> None:
    expected = {
        "remediation_performed": True,
        "reason": "ceph_mon_clock_skew_after_reboot",
        "systemd_timesyncd_restarted": True,
        "ntp_renegotiated_on_high_offset_node": True,
        "final_ceph_health": "HEALTH_OK",
    }
    for key, value in expected.items():
        if remediation.get(key) != value:
            raise SystemExit(f"{path}: time_sync_remediation.{key} must be {value}")


def validate_cleanup(cleanup: dict[str, Any], path: object) -> None:
    for flag in (
        "temporary_worker_vm_namespaces_deleted",
        "temporary_control_plane_observation_namespace_deleted",
        "temporary_pvc_deleted",
        "temporary_pv_deleted",
        "temporary_storageclass_deleted",
    ):
        if cleanup.get(flag) is not True:
            raise SystemExit(f"{path}: cleanup_result.{flag} must be true")
    for flag in ("retained_test_pv_left", "retained_test_storageclass_left"):
        if cleanup.get(flag) is not False:
            raise SystemExit(f"{path}: cleanup_result.{flag} must be false")


def validate_safety(safety: dict[str, Any], path: object) -> None:
    for flag in (
        "root_filesystem_preserved",
        "fstab_unchanged_by_sprint11",
        "no_manual_mount_created",
        "no_default_storageclass_switch",
        "no_existing_pvc_migration",
        "hdd_not_in_vm_primary_pool",
        "storage_automation_uses_persistent_by_id",
    ):
        if safety.get(flag) is not True:
            raise SystemExit(f"{path}: safety_result.{flag} must be true")


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


def main() -> int:
    result = load_result(DEFAULT_RESULT)
    validate_result(result)
    print("Sprint 11 Rook-Ceph reboot resilience result valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
