#!/usr/bin/env python3
"""Tests for Sprint 7 CORE-INSTALLER-A contract validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_installer as installer


class CoreInstallerValidationTest(unittest.TestCase):
    def test_default_profiles_cover_core_install_modes(self) -> None:
        profiles = installer.load_profiles(installer.DEFAULT_PROFILE_DIR)

        installer.validate_profiles(profiles)

        self.assertEqual({"baremetal", "vm", "existing-k8s"}, {profile["mode"] for profile in profiles})
        for profile in profiles:
            self.assertEqual("core", profile["scope"])
            self.assertIn("ani-gateway", profile["core_components"])
            self.assertNotIn("model-service", profile["core_components"])
            self.assertNotIn("kb-service", profile["core_components"])

    def test_validation_rejects_services_component(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "bad.yaml"
            profile.write_text(
                """
profile: bad
mode: baremetal
scope: core
core_components:
  - ani-gateway
  - model-service
preflight:
  required_commands: [kubectl]
plan:
  phases: [preflight, render, validate]
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                installer.validate_profiles(installer.load_profiles(Path(tmpdir)))

        self.assertIn("forbidden Services component", str(raised.exception))

    def test_validation_rejects_missing_preflight_phase(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "bad.yaml"
            profile.write_text(
                """
profile: bad
mode: vm
scope: core
core_components:
  - ani-gateway
preflight:
  required_commands: [kubectl]
plan:
  phases: [render, validate]
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                installer.validate_profiles(installer.load_profiles(Path(tmpdir)))

        self.assertIn("plan phases must include preflight", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
