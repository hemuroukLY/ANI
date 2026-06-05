#!/usr/bin/env python3
"""Tests for Sprint 11 Rook-Ceph reboot resilience result validator."""

from __future__ import annotations

import copy

from validate_sprint11_rook_ceph_reboot_resilience import (
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


def test_rejects_missing_reboot_execution() -> None:
    result = copy.deepcopy(load_default())
    result["live_reboot_resilience_executed"] = False
    expect_failure(result, "live_reboot_resilience_executed")


def test_rejects_simultaneous_reboot() -> None:
    result = copy.deepcopy(load_default())
    result["simultaneous_reboot_performed"] = True
    expect_failure(result, "simultaneous_reboot_performed")


def test_rejects_broken_node_order() -> None:
    result = copy.deepcopy(load_default())
    result["node_reboot_results"] = list(reversed(result["node_reboot_results"]))
    expect_failure(result, "worker-first control-plane-last")


def test_rejects_missing_worker_vm_recovery() -> None:
    result = copy.deepcopy(load_default())
    result["node_reboot_results"][0]["vm_ready_after_reboot"] = False
    expect_failure(result, "vm_ready_after_reboot")


def test_rejects_unclean_cleanup() -> None:
    result = copy.deepcopy(load_default())
    result["cleanup_result"]["retained_test_pv_left"] = True
    expect_failure(result, "retained_test_pv_left")


def test_rejects_missing_time_sync_remediation() -> None:
    result = copy.deepcopy(load_default())
    result["time_sync_remediation"]["remediation_performed"] = False
    expect_failure(result, "time_sync_remediation.remediation_performed")


if __name__ == "__main__":
    test_default_result_valid()
    test_rejects_missing_reboot_execution()
    test_rejects_simultaneous_reboot()
    test_rejects_broken_node_order()
    test_rejects_missing_worker_vm_recovery()
    test_rejects_unclean_cleanup()
    test_rejects_missing_time_sync_remediation()
    print("Sprint 11 Rook-Ceph reboot resilience result validator tests passed")
