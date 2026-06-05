#!/usr/bin/env python3
"""Tests for Sprint 8 CORE-HARDEN-A release hardening validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_release_hardening as hardening


class CoreReleaseHardeningValidationTest(unittest.TestCase):
    def test_default_profile_covers_required_core_gates(self) -> None:
        profile = hardening.load_profile(hardening.DEFAULT_PROFILE)

        hardening.validate_profile(profile)

        gate_ids = {gate["id"] for gate in profile["gates"]}
        self.assertIn("core-api-compatibility", gate_ids)
        self.assertIn("architecture-boundary", gate_ids)
        self.assertIn("doc-entrypoints", gate_ids)
        self.assertIn("sdk-mock-smoke", gate_ids)
        self.assertEqual("core", profile["scope"])

    def test_validation_rejects_services_gate(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "hardening.yaml"
            profile.write_text(
                """
scope: core
gates:
  - id: model-service-hardening
    category: services
    command: make test-model-service
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                hardening.validate_profile(hardening.load_profile(profile))

        self.assertIn("forbidden Services gate", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
