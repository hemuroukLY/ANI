#!/usr/bin/env python3
"""Validate Sprint 11 Rook-Ceph live deployment result."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_RESULT = ROOT / "deploy" / "real-k8s-lab" / "sprint11-rook-ceph-live-deployment-result.yaml"
REQUIRED_SOURCE_PROFILES = {
    "deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml",
    "deploy/real-k8s-lab/sprint11-rook-ceph-formal-deployment.yaml",
}


def load_result(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 Rook-Ceph live result does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 Rook-Ceph live result must be a mapping")
    data["_path"] = str(path)
    return data


def validate_result(result: dict[str, Any]) -> None:
    path = result.get("_path", "<memory>")
    if result.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if result.get("version") != "sprint11":
        raise SystemExit(f"{path}: version must be sprint11")
    if result.get("phase") != "core-real-deployment-live-result":
        raise SystemExit(f"{path}: phase must be core-real-deployment-live-result")
    if result.get("batch") != "CORE-ROOK-CEPH-LIVE-DEPLOYMENT-A":
        raise SystemExit(f"{path}: unexpected batch id")

    for flag in (
        "manual_approval_received",
        "data_safety_prerequisite_enforced",
        "live_deployment_executed",
        "real_server_mutation_performed",
    ):
        if result.get(flag) is not True:
            raise SystemExit(f"{path}: {flag} must be true")
    for flag in (
        "destructive_operations_directly_invoked_by_agent",
        "server_reboot_performed",
        "system_disk_touched",
        "data_disk_mount_or_fstab_mutation",
        "storageclass_default_switched",
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

    validate_operator(require_mapping(result, "operator_installation", path), path)
    validate_cluster(require_mapping(result, "cluster_result", path), path)
    validate_storageclass(require_mapping(result, "storageclass_result", path), path)
    validate_smoke(require_mapping(result, "smoke_test", path), path)
    validate_safety(require_mapping(result, "safety_result", path), path)


def validate_operator(operator: dict[str, Any], path: object) -> None:
    if operator.get("rook_operator_version") != "v1.20.0":
        raise SystemExit(f"{path}: rook_operator_version must be v1.20.0")
    if operator.get("ceph_image") != "quay.io/ceph/ceph:v19.2.3":
        raise SystemExit(f"{path}: ceph_image must be quay.io/ceph/ceph:v19.2.3")
    if operator.get("csi_addons_version") != "v0.14.0":
        raise SystemExit(f"{path}: csi_addons_version must be v0.14.0")
    for flag in (
        "rook_crds_installed",
        "rook_common_installed",
        "csi_operator_installed",
        "rook_operator_available",
        "csi_addons_crds_installed",
        "csi_addons_sidecars_running_after_crd_remediation",
    ):
        if operator.get(flag) is not True:
            raise SystemExit(f"{path}: operator_installation.{flag} must be true")


def validate_cluster(cluster: dict[str, Any], path: object) -> None:
    expected = {
        "kubernetes_nodes_ready": 3,
        "cephcluster_phase": "Ready",
        "ceph_health": "HEALTH_OK",
        "cephblockpool_phase": "Ready",
        "mon_count_running": 3,
        "mgr_count_running": 1,
        "osd_count_running": 5,
        "osd_prepare_jobs_completed": 3,
        "ani3_ssd_candidate_count": 2,
    }
    for key, value in expected.items():
        if cluster.get(key) != value:
            raise SystemExit(f"{path}: cluster_result.{key} must be {value}")
    if cluster.get("cephcluster_name") != "rook-ceph" or cluster.get("cephcluster_namespace") != "rook-ceph":
        raise SystemExit(f"{path}: cluster_result must target rook-ceph/rook-ceph")
    if cluster.get("cephblockpool_name") != "ceph-rbd-ssd":
        raise SystemExit(f"{path}: cluster_result.cephblockpool_name must be ceph-rbd-ssd")
    topology = require_mapping(cluster, "osd_topology", path)
    if topology != {"kubercloud": 1, "dev-phys-02": 2, "dev-phys-03": 2}:
        raise SystemExit(f"{path}: cluster_result.osd_topology must match the 5 approved SSD OSDs")
    if cluster.get("ani3_rotational_hdd_excluded") is not True:
        raise SystemExit(f"{path}: ANI3 rotational HDD must be excluded from VM primary pool")


def validate_storageclass(storageclass: dict[str, Any], path: object) -> None:
    expected = {
        "name": "ani-rbd-ssd",
        "provisioner": "rook-ceph.rbd.csi.ceph.com",
        "reclaimPolicy": "Retain",
        "volumeBindingMode": "WaitForFirstConsumer",
        "allowVolumeExpansion": True,
        "default_storageclass": False,
    }
    for key, value in expected.items():
        if storageclass.get(key) != value:
            raise SystemExit(f"{path}: storageclass_result.{key} must be {value}")


def validate_smoke(smoke: dict[str, Any], path: object) -> None:
    expected = {
        "executed": True,
        "temporary_storageclass": "ani-rbd-ssd-smoke-delete",
        "temporary_storageclass_reclaimPolicy": "Delete",
        "pvc_bound": True,
        "pod_ready": True,
        "mounted_volume_write_read_marker": "sprint11-rbd-smoke",
        "temporary_resources_deleted": True,
        "temporary_pv_deleted": True,
        "retained_test_pv_left": False,
    }
    for key, value in expected.items():
        if smoke.get(key) != value:
            raise SystemExit(f"{path}: smoke_test.{key} must be {value}")


def validate_safety(safety: dict[str, Any], path: object) -> None:
    for flag in (
        "root_filesystem_preserved",
        "fstab_unchanged_by_sprint11",
        "no_server_reboot",
        "no_manual_mount_created",
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
    print("Sprint 11 Rook-Ceph live deployment result valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
