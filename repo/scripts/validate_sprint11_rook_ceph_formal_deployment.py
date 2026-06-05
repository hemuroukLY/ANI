#!/usr/bin/env python3
"""Validate Sprint 11 Rook-Ceph formal deployment code bundle."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_BUNDLE = ROOT / "deploy" / "real-k8s-lab" / "sprint11-rook-ceph-formal-deployment.yaml"
STORAGE_PLAN = ROOT / "deploy" / "real-k8s-lab" / "sprint11-storage-disk-plan.yaml"
REQUIRED_MANIFEST_KINDS = {"Namespace", "CephCluster", "CephBlockPool", "StorageClass"}
REQUIRED_FORBIDDEN_ACTIONS = {
    "kubectl_apply",
    "wipefs",
    "sgdisk",
    "mkfs",
    "mount",
    "fstab_update",
    "rook_ceph_cluster_install",
    "rook_ceph_osd_claim",
    "storageclass_default_switch",
    "server_reboot",
}


def load_yaml(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"required Sprint 11 file does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: expected YAML mapping")
    data["_path"] = str(path)
    return data


def validate_bundle(bundle: dict[str, Any], storage_plan: dict[str, Any]) -> None:
    path = bundle.get("_path", "<memory>")
    if bundle.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if bundle.get("version") != "sprint11":
        raise SystemExit(f"{path}: version must be sprint11")
    if bundle.get("phase") != "core-formal-deployment-code":
        raise SystemExit(f"{path}: phase must be core-formal-deployment-code")
    if bundle.get("batch") != "CORE-ROOK-CEPH-FORMAL-DEPLOYMENT-A":
        raise SystemExit(f"{path}: unexpected batch id")
    if bundle.get("formal_deployment_code_complete") is not True:
        raise SystemExit(f"{path}: formal_deployment_code_complete must be true")
    if bundle.get("rook_operator_version") != "v1.20.0":
        raise SystemExit(f"{path}: rook_operator_version must be v1.20.0")

    for flag in (
        "live_deployment_executed",
        "real_server_mutation_performed",
        "mutation_allowed_without_approval",
        "production_storage_deployment_complete",
        "rook_ceph_cluster_installed",
        "rook_ceph_osd_claimed",
        "storageclass_default_switched",
    ):
        if bundle.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false before approved live deployment")
    if bundle.get("data_safety_prerequisite_enforced") is not True:
        raise SystemExit(f"{path}: data_safety_prerequisite_enforced must be true")

    validate_source_profiles(bundle, path)
    validate_policy(bundle, path)
    candidate_devices = expected_candidate_devices(storage_plan)
    manifests = require_list(bundle, "manifests", path)
    by_kind = validate_manifest_set(manifests, path)
    validate_ceph_cluster(by_kind["CephCluster"], candidate_devices, path)
    validate_block_pool(by_kind["CephBlockPool"], path)
    validate_storage_class(by_kind["StorageClass"], path)


def validate_source_profiles(bundle: dict[str, Any], path: object) -> None:
    profiles = set(require_list(bundle, "source_profiles", path))
    for rel_path in (
        "deploy/real-k8s-lab/sprint11-storage-disk-plan.yaml",
        "deploy/real-k8s-lab/sprint11-core-safe-completion.yaml",
    ):
        if rel_path not in profiles:
            raise SystemExit(f"{path}: missing source profile {rel_path}")
        if not (ROOT / rel_path).exists():
            raise SystemExit(f"{path}: source profile does not exist: {rel_path}")


def validate_policy(bundle: dict[str, Any], path: object) -> None:
    device_policy = require_mapping(bundle, "device_selection_policy", path)
    if device_policy.get("selector") != "explicit_persistent_by_id":
        raise SystemExit(f"{path}: device selector must be explicit_persistent_by_id")
    for key in (
        "require_all_devices_from_storage_disk_plan",
        "reject_kernel_sd_names",
        "reject_rotational_hdd_initial_pool",
        "reject_mounted_or_formatted_devices",
        "reject_unapproved_wipe",
    ):
        if device_policy.get(key) is not True:
            raise SystemExit(f"{path}: device_selection_policy.{key} must be true")

    approval = require_mapping(bundle, "approval_gate", path)
    if approval.get("required") is not True:
        raise SystemExit(f"{path}: approval_gate.required must be true")
    require_string(approval, "approval_token_name", path)
    confirmations = set(require_list(approval, "approval_must_confirm", path))
    for required in (
        "exact_by_id_devices_match_current_read_only_inventory",
        "every_selected_device_is_unmounted",
        "every_selected_device_has_no_filesystem_partition_lvm_or_mountpoint",
        "rollback_and_recovery_window_is_approved",
    ):
        if required not in confirmations:
            raise SystemExit(f"{path}: approval gate missing confirmation {required}")

    forbidden = set(require_list(bundle, "forbidden_without_approval", path))
    missing_forbidden = REQUIRED_FORBIDDEN_ACTIONS - forbidden
    if missing_forbidden:
        raise SystemExit(f"{path}: missing forbidden actions: {', '.join(sorted(missing_forbidden))}")

    prerequisites = require_mapping(bundle, "operator_prerequisites", path)
    for key in ("rook_operator_required_before_apply", "rook_crds_required_before_apply"):
        if prerequisites.get(key) is not True:
            raise SystemExit(f"{path}: operator_prerequisites.{key} must be true")


def expected_candidate_devices(storage_plan: dict[str, Any]) -> dict[str, set[str]]:
    result: dict[str, set[str]] = {}
    for node in require_list(storage_plan, "nodes", storage_plan.get("_path", "<storage-plan>")):
        k8s_name = require_string(node, "k8s_node_observed", storage_plan.get("_path", "<storage-plan>"))
        devices: set[str] = set()
        for disk in require_list(node, "data_disks", storage_plan.get("_path", "<storage-plan>")):
            if disk.get("candidate_for_rook_osd") is True:
                if disk.get("rotational") is True:
                    raise SystemExit(f"{storage_plan.get('_path')}: rotational disk cannot be an SSD OSD candidate")
                if disk.get("mounted") is not False or disk.get("filesystem_observed") != "none":
                    raise SystemExit(f"{storage_plan.get('_path')}: OSD candidates must be unmounted and unformatted")
                devices.add(require_by_id(disk, storage_plan.get("_path", "<storage-plan>")))
        if devices:
            result[k8s_name] = devices
    if len(result) != 3:
        raise SystemExit(f"{storage_plan.get('_path')}: expected SSD OSD candidates on all three nodes")
    return result


def validate_manifest_set(manifests: list[Any], path: object) -> dict[str, dict[str, Any]]:
    by_kind: dict[str, dict[str, Any]] = {}
    for manifest in manifests:
        if not isinstance(manifest, dict):
            raise SystemExit(f"{path}: manifest entries must be mappings")
        kind = require_string(manifest, "kind", path)
        if kind in by_kind:
            raise SystemExit(f"{path}: duplicate manifest kind {kind}")
        by_kind[kind] = manifest
        metadata = require_mapping(manifest, "metadata", path)
        require_string(metadata, "name", path)
    missing = REQUIRED_MANIFEST_KINDS - set(by_kind)
    if missing:
        raise SystemExit(f"{path}: missing manifests: {', '.join(sorted(missing))}")
    return by_kind


def validate_ceph_cluster(manifest: dict[str, Any], expected_devices: dict[str, set[str]], path: object) -> None:
    if manifest.get("apiVersion") != "ceph.rook.io/v1":
        raise SystemExit(f"{path}: CephCluster apiVersion must be ceph.rook.io/v1")
    spec = require_mapping(manifest, "spec", path)
    storage = require_mapping(spec, "storage", path)
    if storage.get("useAllNodes") is not False or storage.get("useAllDevices") is not False:
        raise SystemExit(f"{path}: CephCluster must use explicit nodes and devices only")
    if require_mapping(spec, "mon", path).get("count") != 3:
        raise SystemExit(f"{path}: CephCluster mon.count must be 3")
    require_string(require_mapping(spec, "cephVersion", path), "image", path)
    validate_control_plane_toleration(spec, path)

    nodes = require_list(storage, "nodes", path)
    actual_devices: dict[str, set[str]] = {}
    for node in nodes:
        node_name = require_string(node, "name", path)
        devices: set[str] = set()
        for device in require_list(node, "devices", path):
            device_name = require_by_id(device, path)
            if "/dev/sd" in device_name or not device_name.startswith("/dev/disk/by-id/"):
                raise SystemExit(f"{path}: CephCluster devices must use persistent /dev/disk/by-id names")
            devices.add(device_name)
        actual_devices[node_name] = devices
    if actual_devices != expected_devices:
        raise SystemExit(f"{path}: CephCluster devices must exactly match storage disk plan SSD candidates")


def validate_control_plane_toleration(spec: dict[str, Any], path: object) -> None:
    placement = require_mapping(spec, "placement", path)
    all_placement = require_mapping(placement, "all", path)
    tolerations = require_list(all_placement, "tolerations", path)
    for toleration in tolerations:
        if not isinstance(toleration, dict):
            raise SystemExit(f"{path}: CephCluster tolerations must be mappings")
        if (
            toleration.get("key") == "node-role.kubernetes.io/control-plane"
            and toleration.get("operator") == "Exists"
            and toleration.get("effect") == "NoSchedule"
        ):
            return
    raise SystemExit(f"{path}: CephCluster must tolerate the control-plane NoSchedule taint")


def validate_block_pool(manifest: dict[str, Any], path: object) -> None:
    spec = require_mapping(manifest, "spec", path)
    if spec.get("failureDomain") != "host":
        raise SystemExit(f"{path}: CephBlockPool failureDomain must be host")
    replicated = require_mapping(spec, "replicated", path)
    if replicated.get("size") != 3:
        raise SystemExit(f"{path}: CephBlockPool replicated.size must be 3")
    if replicated.get("requireSafeReplicaSize") is not True:
        raise SystemExit(f"{path}: CephBlockPool must require safe replica size")


def validate_storage_class(manifest: dict[str, Any], path: object) -> None:
    metadata = require_mapping(manifest, "metadata", path)
    annotations = require_mapping(metadata, "annotations", path)
    if annotations.get("storageclass.kubernetes.io/is-default-class") != "false":
        raise SystemExit(f"{path}: Sprint 11 StorageClass must not become default automatically")
    if manifest.get("provisioner") != "rook-ceph.rbd.csi.ceph.com":
        raise SystemExit(f"{path}: StorageClass provisioner must be rook-ceph.rbd.csi.ceph.com")
    if manifest.get("reclaimPolicy") != "Retain":
        raise SystemExit(f"{path}: StorageClass reclaimPolicy must be Retain for data safety")
    if manifest.get("volumeBindingMode") != "WaitForFirstConsumer":
        raise SystemExit(f"{path}: StorageClass volumeBindingMode must be WaitForFirstConsumer")
    if manifest.get("allowVolumeExpansion") is not True:
        raise SystemExit(f"{path}: StorageClass allowVolumeExpansion must be true")
    parameters = require_mapping(manifest, "parameters", path)
    if parameters.get("clusterID") != "rook-ceph" or parameters.get("pool") != "ceph-rbd-ssd":
        raise SystemExit(f"{path}: StorageClass must target rook-ceph/ceph-rbd-ssd")


def require_by_id(mapping: dict[str, Any], path: object) -> str:
    value = require_string(mapping, "stable_by_id" if "stable_by_id" in mapping else "name", path)
    if not value.startswith("/dev/disk/by-id/"):
        raise SystemExit(f"{path}: device must use /dev/disk/by-id")
    return value


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
    bundle = load_yaml(DEFAULT_BUNDLE)
    storage_plan = load_yaml(STORAGE_PLAN)
    validate_bundle(bundle, storage_plan)
    print("Sprint 11 Rook-Ceph formal deployment bundle valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
