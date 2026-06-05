#!/usr/bin/env python3
"""Tests for Sprint 11 Core safe completion validator."""

from __future__ import annotations

import copy

from validate_sprint11_safe_completion import DEFAULT_PROFILE, load_profile, validate_profile


def expect_failure(profile: dict, expected: str) -> None:
    try:
        validate_profile(profile)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def test_default_profile_valid() -> None:
    validate_profile(load_profile(DEFAULT_PROFILE))


def test_rejects_mutation() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["mutation_allowed"] = True
    expect_failure(profile, "mutation_allowed")


def test_rejects_missing_final_safe_completion_marker() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["sprint_final_safe_complete"] = False
    expect_failure(profile, "sprint_final_safe_complete")


def test_rejects_mutating_execution_environment() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["execution_environment"]["storage_mutation_allowed"] = True
    expect_failure(profile, "execution_environment.storage_mutation_allowed")


def test_rejects_missing_forbidden_execution_class() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["execution_environment"]["forbidden_operation_classes"].remove("server_reboot")
    expect_failure(profile, "server_reboot")


def test_rejects_rook_install_claim() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["rook_ceph_installed"] = True
    expect_failure(profile, "rook_ceph_installed")


def test_rejects_missing_open_source_principle() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["open_source_principles"].remove("persistent_device_identity_over_kernel_sd_letter")
    expect_failure(profile, "persistent_device_identity_over_kernel_sd_letter")


def test_rejects_missing_manual_approval() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["manual_approval_required_before"].remove("reboot_for_storage")
    expect_failure(profile, "reboot_for_storage")


if __name__ == "__main__":
    test_default_profile_valid()
    test_rejects_mutation()
    test_rejects_missing_final_safe_completion_marker()
    test_rejects_mutating_execution_environment()
    test_rejects_missing_forbidden_execution_class()
    test_rejects_rook_install_claim()
    test_rejects_missing_open_source_principle()
    test_rejects_missing_manual_approval()
    print("Sprint 11 safe completion validator tests passed")
