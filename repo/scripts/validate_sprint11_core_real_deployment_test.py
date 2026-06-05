#!/usr/bin/env python3
"""Tests for Sprint 11 Core real deployment profile validator."""

from __future__ import annotations

import copy

from validate_sprint11_core_real_deployment import DEFAULT_PROFILE, load_profile, validate_profile


def expect_failure(profile: dict, expected: str) -> None:
    try:
        validate_profile(profile)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def test_default_profile_valid() -> None:
    validate_profile(load_profile(DEFAULT_PROFILE))


def test_rejects_mutation_allowed() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["mutation_allowed"] = True
    expect_failure(profile, "mutation_allowed")


def test_rejects_rook_present_before_approval() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["cluster_observed"]["rook_ceph_namespace_present"] = True
    expect_failure(profile, "Rook-Ceph")


def test_rejects_missing_storage_gate() -> None:
    profile = copy.deepcopy(load_profile(DEFAULT_PROFILE))
    profile["gates"] = [gate for gate in profile["gates"] if gate["id"] != "storage-disk-plan"]
    expect_failure(profile, "storage-disk-plan")


if __name__ == "__main__":
    test_default_profile_valid()
    test_rejects_mutation_allowed()
    test_rejects_rook_present_before_approval()
    test_rejects_missing_storage_gate()
    print("Sprint 11 real deployment profile validator tests passed")
