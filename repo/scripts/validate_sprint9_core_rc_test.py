#!/usr/bin/env python3
"""Tests for Sprint 9 CORE-RC-GATE-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_sprint9_core_rc as sprint9_rc


class Sprint9CoreRCValidationTest(unittest.TestCase):
    def test_default_profile_covers_required_core_rc_gates(self) -> None:
        profile = sprint9_rc.load_profile(sprint9_rc.DEFAULT_PROFILE)

        sprint9_rc.validate_profile(profile)

        self.assertEqual("core", profile["scope"])
        self.assertTrue(profile["rc_readiness"])
        self.assertFalse(profile["release_candidate"])
        gate_ids = {gate["id"] for gate in profile["gates"]}
        self.assertIn("sprint8-core-release", gate_ids)
        self.assertIn("core-release-evidence", gate_ids)
        self.assertIn("sprint9-doc-consistency", gate_ids)

    def test_validation_rejects_services_gate(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "sprint9-core-rc.yaml"
            profile.write_text(
                """
scope: core
version: sprint9
rc_readiness: true
release_candidate: false
production_release: false
gates:
  - id: model-service-rc
    category: services
    command: make test-model-service
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                sprint9_rc.validate_profile(sprint9_rc.load_profile(profile))

        self.assertIn("forbidden Services gate", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
