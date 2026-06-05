#!/usr/bin/env python3
"""Tests for Sprint 11 Rook-Ceph VM storage smoke result validator."""

from __future__ import annotations

import copy

from validate_sprint11_rook_ceph_vm_storage_smoke import (
    DEFAULT_RESULT,
    load_result,
    validate_result,
)


def expect_failure(result: dict, expected: str) -> None:
    try:
        validate_result(result)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def load_default() -> dict:
    return load_result(DEFAULT_RESULT)


def test_default_result_valid() -> None:
    validate_result(load_default())


def test_rejects_missing_vm_storage_smoke() -> None:
    result = copy.deepcopy(load_default())
    result["live_vm_storage_smoke_executed"] = False
    expect_failure(result, "live_vm_storage_smoke_executed")


def test_rejects_default_storageclass_switch() -> None:
    result = copy.deepcopy(load_default())
    result["storageclass_default_switched"] = True
    expect_failure(result, "storageclass_default_switched")


def test_rejects_missing_guest_probe() -> None:
    result = copy.deepcopy(load_default())
    result["vm_result"]["guest_probe_completed"] = False
    expect_failure(result, "guest_probe_completed")


def test_rejects_unclean_pv_cleanup() -> None:
    result = copy.deepcopy(load_default())
    result["cleanup_result"]["retained_test_pv_left"] = True
    expect_failure(result, "retained_test_pv_left")


if __name__ == "__main__":
    test_default_result_valid()
    test_rejects_missing_vm_storage_smoke()
    test_rejects_default_storageclass_switch()
    test_rejects_missing_guest_probe()
    test_rejects_unclean_pv_cleanup()
    print("Sprint 11 Rook-Ceph VM storage smoke result validator tests passed")
