#!/usr/bin/env python3
"""Tests for Sprint 11 Rook-Ceph formal deployment validator."""

from __future__ import annotations

import copy

from validate_sprint11_rook_ceph_formal_deployment import (
    DEFAULT_BUNDLE,
    STORAGE_PLAN,
    load_yaml,
    validate_bundle,
)


def expect_failure(bundle: dict, storage_plan: dict, expected: str) -> None:
    try:
        validate_bundle(bundle, storage_plan)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def load_defaults() -> tuple[dict, dict]:
    return load_yaml(DEFAULT_BUNDLE), load_yaml(STORAGE_PLAN)


def find_manifest(bundle: dict, kind: str) -> dict:
    for manifest in bundle["manifests"]:
        if manifest["kind"] == kind:
            return manifest
    raise AssertionError(f"missing manifest kind {kind}")


def test_default_bundle_valid() -> None:
    validate_bundle(*load_defaults())


def test_rejects_live_deployment_claim() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    bundle["live_deployment_executed"] = True
    expect_failure(bundle, storage_plan, "live_deployment_executed")


def test_rejects_operator_version_drift() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    bundle["rook_operator_version"] = "v1.19.0"
    expect_failure(bundle, storage_plan, "rook_operator_version")


def test_rejects_missing_control_plane_toleration() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    cluster = find_manifest(bundle, "CephCluster")
    cluster["spec"].pop("placement")
    expect_failure(bundle, storage_plan, "placement")


def test_rejects_kernel_sd_device_name() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    cluster = find_manifest(bundle, "CephCluster")
    cluster["spec"]["storage"]["nodes"][0]["devices"][0]["name"] = "/dev/sda"
    expect_failure(bundle, storage_plan, "/dev/disk/by-id")


def test_rejects_extra_hdd_or_unplanned_device() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    cluster = find_manifest(bundle, "CephCluster")
    cluster["spec"]["storage"]["nodes"][2]["devices"].append(
        {"name": "/dev/disk/by-id/scsi-SSEAGATE_ST4000NM0125_ZC119FL80000R7271R0L"}
    )
    expect_failure(bundle, storage_plan, "exactly match")


def test_rejects_default_storageclass_switch() -> None:
    bundle, storage_plan = load_defaults()
    bundle = copy.deepcopy(bundle)
    storage_class = find_manifest(bundle, "StorageClass")
    storage_class["metadata"]["annotations"]["storageclass.kubernetes.io/is-default-class"] = "true"
    expect_failure(bundle, storage_plan, "must not become default")


if __name__ == "__main__":
    test_default_bundle_valid()
    test_rejects_live_deployment_claim()
    test_rejects_operator_version_drift()
    test_rejects_missing_control_plane_toleration()
    test_rejects_kernel_sd_device_name()
    test_rejects_extra_hdd_or_unplanned_device()
    test_rejects_default_storageclass_switch()
    print("Sprint 11 Rook-Ceph formal deployment validator tests passed")
