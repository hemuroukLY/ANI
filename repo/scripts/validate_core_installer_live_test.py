#!/usr/bin/env python3
"""Tests for Sprint 8 CORE-INSTALLER-LIVE-A readiness validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_installer_live as installer_live


class CoreInstallerLiveValidationTest(unittest.TestCase):
    def test_default_profile_requires_all_install_modes(self) -> None:
        profile = installer_live.load_profile(installer_live.DEFAULT_PROFILE)

        installer_live.validate_profile(profile)

        self.assertEqual("core", profile["scope"])
        self.assertEqual({"baremetal", "vm", "existing-k8s"}, {target["mode"] for target in profile["targets"]})

    def test_validation_rejects_production_ready_claim(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "installer-live.yaml"
            profile.write_text(
                """
scope: core
evidence_level: contract-local
production_ready: true
targets:
  - id: baremetal
    mode: baremetal
    profile: installer/ani-installer/profiles/baremetal.yaml
    command: make validate-core-installer
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                installer_live.validate_profile(installer_live.load_profile(profile))

        self.assertIn("must not claim production ready", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
