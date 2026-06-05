#!/usr/bin/env python3
"""Validate Sprint 11 Core storage disk risk plan."""

from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_PLAN = ROOT / "deploy" / "real-k8s-lab" / "sprint11-storage-disk-plan.yaml"
EXPECTED_NODES = {"ani1", "ani2", "ani3"}
FORBIDDEN_SELECTORS = {"/dev/sdX", "kernel_name_order"}
REQUIRED_RISK_CONTROLS = {
    "never_change_boot_disk_order_for_alignment",
    "never_reference_osd_disks_by_sd_letter_in_automation",
    "fail_closed_if_filesystem_partition_lvm_or_mountpoint_is_detected",
    "require_human_approval_before_wipefs_sgdisk_mkfs_mount_or_fstab_change",
}


def load_plan(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"Sprint 11 storage disk plan does not exist: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise SystemExit(f"{path}: Sprint 11 storage disk plan must be a mapping")
    data["_path"] = str(path)
    return data


def validate_plan(plan: dict[str, Any]) -> None:
    path = plan.get("_path", "<memory>")
    if plan.get("scope") != "core":
        raise SystemExit(f"{path}: scope must be core")
    if plan.get("version") != "sprint11":
        raise SystemExit(f"{path}: version must be sprint11")
    for flag in ("mutation_allowed", "destructive_operations_allowed", "real_server_mutation_performed"):
        if plan.get(flag) is not False:
            raise SystemExit(f"{path}: {flag} must be false during Sprint 11 storage preflight")
    if plan.get("disk_order_policy") != "do_not_align_sd_letters":
        raise SystemExit(f"{path}: disk_order_policy must reject sd-letter alignment")
    if plan.get("storage_device_identity_policy") != "persistent_device_id_required":
        raise SystemExit(f"{path}: storage_device_identity_policy must require persistent device ids")

    rook_policy = require_mapping(plan, "rook_ceph_policy", path)
    if rook_policy.get("pre_mount_osd_devices") is not False:
        raise SystemExit(f"{path}: Rook-Ceph OSD devices must not be pre-mounted")
    if rook_policy.get("require_explicit_wipe_approval") is not True:
        raise SystemExit(f"{path}: Rook-Ceph plan must require explicit wipe approval")
    forbidden = set(require_list(rook_policy, "forbidden_device_selectors", path))
    if not FORBIDDEN_SELECTORS.issubset(forbidden):
        raise SystemExit(f"{path}: Rook-Ceph plan must forbid sd-letter and kernel-order selectors")
    allowed = set(require_list(rook_policy, "allowed_device_selectors", path))
    if "/dev/disk/by-id" not in allowed:
        raise SystemExit(f"{path}: Rook-Ceph plan must allow /dev/disk/by-id selectors")

    controls = set(require_list(plan, "risk_controls", path))
    missing_controls = REQUIRED_RISK_CONTROLS - controls
    if missing_controls:
        raise SystemExit(f"{path}: missing risk controls: {', '.join(sorted(missing_controls))}")

    nodes = require_list(plan, "nodes", path)
    seen_nodes: set[str] = set()
    data_disk_count = 0
    ssd_osd_candidates = 0
    for node in nodes:
        if not isinstance(node, dict):
            raise SystemExit(f"{path}: node entries must be mappings")
        node_id = require_string(node, "node_id", path)
        if node_id in seen_nodes:
            raise SystemExit(f"{path}: duplicate node_id: {node_id}")
        seen_nodes.add(node_id)
        system_disk = require_mapping(node, "system_disk", path)
        if system_disk.get("root_mounted") is not True:
            raise SystemExit(f"{path}: {node_id} system disk must have root_mounted=true")
        validate_stable_by_id(system_disk, f"{path}: {node_id} system_disk")
        data_disks = require_list(node, "data_disks", path)
        if not data_disks:
            raise SystemExit(f"{path}: {node_id} must define at least one data disk")
        for disk in data_disks:
            if not isinstance(disk, dict):
                raise SystemExit(f"{path}: {node_id} data disk entries must be mappings")
            data_disk_count += 1
            validate_stable_by_id(disk, f"{path}: {node_id} data disk")
            if disk.get("mounted") is not False:
                raise SystemExit(f"{path}: {node_id} data disk {disk.get('slot')} must be unmounted")
            if disk.get("filesystem_observed") != "none":
                raise SystemExit(f"{path}: {node_id} data disk {disk.get('slot')} must have no filesystem")
            if disk.get("candidate_for_rook_osd") is True and disk.get("rotational") is False:
                ssd_osd_candidates += 1
            if disk.get("rotational") is True and disk.get("candidate_for_rook_osd") is True:
                raise SystemExit(f"{path}: rotational disks must not join the initial VM SSD pool")
    missing_nodes = EXPECTED_NODES - seen_nodes
    if missing_nodes:
        raise SystemExit(f"{path}: missing physical nodes: {', '.join(sorted(missing_nodes))}")
    if data_disk_count < 6:
        raise SystemExit(f"{path}: expected at least six data disks across the lab")
    if ssd_osd_candidates < 5:
        raise SystemExit(f"{path}: expected at least five SSD OSD candidates for VM-first validation")


def validate_stable_by_id(mapping: dict[str, Any], label: str) -> None:
    stable_by_id = require_string(mapping, "stable_by_id", label)
    if not stable_by_id.startswith("/dev/disk/by-id/"):
        raise SystemExit(f"{label}: stable_by_id must use /dev/disk/by-id")
    kernel_name = require_string(mapping, "kernel_name_observed", label)
    if stable_by_id.endswith(kernel_name) or stable_by_id == f"/dev/{kernel_name}":
        raise SystemExit(f"{label}: stable_by_id must not be a kernel sdX name")


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
    plan = load_plan(DEFAULT_PLAN)
    validate_plan(plan)
    print("Sprint 11 Core storage disk plan valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
