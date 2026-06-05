#!/usr/bin/env python3
"""Validate Sprint 11 Rook-Ceph backed KubeVirt VM storage smoke result."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_RESULT = ROOT / "deploy" / "real-k8s-lab" / "sprint11-rook-ceph-vm-storage-smoke-result.yaml"
REQUIRED_SOURCE_PROFILES = {
    "deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-live-deployment-result.yaml",
}


def load_result(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 Rook-Ceph VM storage smoke result does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 Rook-Ceph VM storage smoke result must be a mapping")
    data["_path"] = str(path)
    return data


def validate_result(result: dict[str, Any]) -> None:
    path = result.get("_path", "<memory>")
    expected_top = {
        "scope": "core",
        "version": "sprint11",
        "phase": "core-real-deployment-vm-storage-smoke-result",
        "batch": "CORE-ROOK-CEPH-VM-STORAGE-SMOKE-A",
        "deployment_target": "kubevirt-vm-with-rook-ceph-rbd-block-pvc",
    }
    for key, value in expected_top.items():
        if result.get(key) != value:
            raise SystemExit(f"{path}: {key} must be {value}")

    for flag in (
        "manual_approval_received",
        "data_safety_prerequisite_enforced",
        "live_vm_storage_smoke_executed",
        "real_server_mutation_performed",
    ):
        if result.get(flag) is not True:
            raise SystemExit(f"{path}: {flag} must be true")
    for flag in (
        "destructive_operations_directly_invoked_by_agent",
        "server_reboot_performed",
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

    validate_storageclass(require_mapping(result, "storageclass_result", path), path)
    validate_pvc(require_mapping(result, "pvc_result", path), path)
    validate_vm(require_mapping(result, "vm_result", path), path)
    validate_cleanup(require_mapping(result, "cleanup_result", path), path)
    validate_cluster(require_mapping(result, "cluster_after_cleanup", path), path)
    validate_safety(require_mapping(result, "safety_result", path), path)


def validate_storageclass(storageclass: dict[str, Any], path: object) -> None:
    expected = {
        "formal_storageclass": "ani-rbd-ssd",
        "formal_storageclass_reclaimPolicy": "Retain",
        "formal_storageclass_default": False,
        "temporary_storageclass": "ani-rbd-ssd-vm-smoke-delete",
        "temporary_storageclass_reclaimPolicy": "Delete",
        "provisioner": "rook-ceph.rbd.csi.ceph.com",
        "pool": "ceph-rbd-ssd",
        "volumeBindingMode": "WaitForFirstConsumer",
        "allowVolumeExpansion": True,
    }
    for key, value in expected.items():
        if storageclass.get(key) != value:
            raise SystemExit(f"{path}: storageclass_result.{key} must be {value}")


def validate_pvc(pvc: dict[str, Any], path: object) -> None:
    expected = {
        "namespace": "ani-sprint11-vm-storage-smoke",
        "pvc_name": "ani-vm-rbd-data",
        "volumeMode": "Block",
        "requested_storage": "1Gi",
        "pvc_bound": True,
        "pv_reclaimPolicy": "Delete",
        "csi_driver": "rook-ceph.rbd.csi.ceph.com",
    }
    for key, value in expected.items():
        if pvc.get(key) != value:
            raise SystemExit(f"{path}: pvc_result.{key} must be {value}")
    pv_name = pvc.get("pv_name_observed")
    if not isinstance(pv_name, str) or not pv_name.startswith("pvc-"):
        raise SystemExit(f"{path}: pvc_result.pv_name_observed must record the temporary PV name")


def validate_vm(vm: dict[str, Any], path: object) -> None:
    expected = {
        "vm_name": "ani-vm-rbd-smoke",
        "kubevirt_vm_ready": True,
        "kubevirt_vm_printableStatus": "Running",
        "kubevirt_vmi_phase": "Running",
        "virt_launcher_pod_running": True,
        "rbd_block_device_visible_in_guest": True,
        "guest_block_write_attempted": True,
        "guest_probe_completed": True,
    }
    for key, value in expected.items():
        if vm.get(key) != value:
            raise SystemExit(f"{path}: vm_result.{key} must be {value}")
    if vm.get("guest_block_device_observed") not in {"/dev/vdb", "/dev/vdc", "/dev/vdd"}:
        raise SystemExit(f"{path}: vm_result.guest_block_device_observed must be an expected virtio block device")
    node = vm.get("node_observed")
    if not isinstance(node, str) or not node.strip():
        raise SystemExit(f"{path}: vm_result.node_observed must be set")


def validate_cleanup(cleanup: dict[str, Any], path: object) -> None:
    for flag in (
        "temporary_vm_deleted",
        "temporary_vmi_deleted",
        "temporary_pvc_deleted",
        "temporary_pv_deleted",
        "temporary_storageclass_deleted",
        "temporary_namespace_delete_requested",
    ):
        if cleanup.get(flag) is not True:
            raise SystemExit(f"{path}: cleanup_result.{flag} must be true")
    for flag in ("retained_test_pv_left", "retained_test_storageclass_left"):
        if cleanup.get(flag) is not False:
            raise SystemExit(f"{path}: cleanup_result.{flag} must be false")


def validate_cluster(cluster: dict[str, Any], path: object) -> None:
    expected = {
        "cephcluster_phase": "Ready",
        "ceph_health": "HEALTH_OK",
        "cephblockpool_phase": "Ready",
    }
    for key, value in expected.items():
        if cluster.get(key) != value:
            raise SystemExit(f"{path}: cluster_after_cleanup.{key} must be {value}")


def validate_safety(safety: dict[str, Any], path: object) -> None:
    for flag in (
        "root_filesystem_preserved",
        "fstab_unchanged_by_sprint11",
        "no_server_reboot",
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
    print("Sprint 11 Rook-Ceph VM storage smoke result valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
