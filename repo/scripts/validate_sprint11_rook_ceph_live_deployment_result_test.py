#!/usr/bin/env python3
"""Tests for Sprint 11 Rook-Ceph live deployment result validator."""

from __future__ import annotations

import copy

from validate_sprint11_rook_ceph_live_deployment_result import (
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


def test_rejects_missing_live_execution() -> None:
    result = copy.deepcopy(load_default())
    result["live_deployment_executed"] = False
    expect_failure(result, "live_deployment_executed")


def test_rejects_default_storageclass_switch() -> None:
    result = copy.deepcopy(load_default())
    result["storageclass_result"]["default_storageclass"] = True
    expect_failure(result, "default_storageclass")


def test_rejects_wrong_osd_count() -> None:
    result = copy.deepcopy(load_default())
    result["cluster_result"]["osd_count_running"] = 6
    expect_failure(result, "osd_count_running")


def test_rejects_unclean_smoke_test() -> None:
    result = copy.deepcopy(load_default())
    result["smoke_test"]["retained_test_pv_left"] = True
    expect_failure(result, "retained_test_pv_left")


if __name__ == "__main__":
    test_default_result_valid()
    test_rejects_missing_live_execution()
    test_rejects_default_storageclass_switch()
    test_rejects_wrong_osd_count()
    test_rejects_unclean_smoke_test()
    print("Sprint 11 Rook-Ceph live deployment result validator tests passed")
