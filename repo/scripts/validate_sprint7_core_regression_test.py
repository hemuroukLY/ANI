#!/usr/bin/env python3
"""Tests for Sprint 7 Core regression gate composition."""

from __future__ import annotations

import unittest

import validate_sprint7_core_regression as regression


class Sprint7CoreRegressionValidationTest(unittest.TestCase):
    def test_regression_targets_include_sprint7_and_historical_gates(self) -> None:
        profile = regression.load_profile(regression.DEFAULT_PROFILE)

        regression.validate_profile(profile)

        target_ids = {target["id"] for target in profile["targets"]}
        self.assertIn("core-installer-contract", target_ids)
        self.assertIn("core-offline-contract", target_ids)
        self.assertIn("core-cli-contract", target_ids)
        self.assertIn("real-k8s-profile-contract", target_ids)
        self.assertIn("sdk-mock-smoke", target_ids)

    def test_validation_rejects_new_real_lab_guard_target(self) -> None:
        profile = {
            "scope": "core",
            "targets": [
                {
                    "id": "bad-guard",
                    "command": "make validate-m1-real-lab-ky",
                    "category": "M1-REAL-LAB-KY",
                }
            ],
        }

        with self.assertRaises(SystemExit) as raised:
            regression.validate_profile(profile)

        self.assertIn("must not add M1-REAL-LAB", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
