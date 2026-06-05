#!/usr/bin/env python3
"""Tests for Sprint 11 Core storage disk risk plan validator."""

from __future__ import annotations

import copy

from validate_sprint11_storage_disk_plan import DEFAULT_PLAN, load_plan, validate_plan


def expect_failure(plan: dict, expected: str) -> None:
    try:
        validate_plan(plan)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def test_default_plan_valid() -> None:
    validate_plan(load_plan(DEFAULT_PLAN))


def test_rejects_destructive_mode() -> None:
    plan = copy.deepcopy(load_plan(DEFAULT_PLAN))
    plan["destructive_operations_allowed"] = True
    expect_failure(plan, "destructive_operations_allowed")


def test_rejects_sd_letter_identity() -> None:
    plan = copy.deepcopy(load_plan(DEFAULT_PLAN))
    plan["nodes"][0]["data_disks"][0]["stable_by_id"] = "/dev/sda"
    expect_failure(plan, "/dev/disk/by-id")


def test_rejects_hdd_initial_osd_candidate() -> None:
    plan = copy.deepcopy(load_plan(DEFAULT_PLAN))
    plan["nodes"][2]["data_disks"][2]["candidate_for_rook_osd"] = True
    expect_failure(plan, "rotational disks")


if __name__ == "__main__":
    test_default_plan_valid()
    test_rejects_destructive_mode()
    test_rejects_sd_letter_identity()
    test_rejects_hdd_initial_osd_candidate()
    print("Sprint 11 storage disk plan validator tests passed")
