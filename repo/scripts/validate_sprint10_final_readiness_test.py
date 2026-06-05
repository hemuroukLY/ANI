#!/usr/bin/env python3
"""Tests for Sprint 10 CORE-FINAL-READINESS-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_sprint10_final_readiness as final_readiness


class Sprint10FinalReadinessValidationTest(unittest.TestCase):
    def test_default_profile_covers_required_release_prep_gates(self) -> None:
        profile = final_readiness.load_profile(final_readiness.DEFAULT_PROFILE)

        final_readiness.validate_profile(profile)

        self.assertEqual("core", profile["scope"])
        self.assertTrue(profile["release_preparation_complete"])
        self.assertFalse(profile["actual_release"])
        gate_ids = {gate["id"] for gate in profile["gates"]}
        self.assertIn("sprint9-rc", gate_ids)
        self.assertIn("core-artifact-manifest", gate_ids)

    def test_validation_rejects_actual_release_true(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            profile = Path(tmpdir) / "sprint10-core-final-readiness.yaml"
            profile.write_text(
                """
scope: core
version: sprint10
target_release: v1.0.0
release_preparation_complete: true
actual_release: true
release_candidate: false
production_release: false
gates: []
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                final_readiness.validate_profile(final_readiness.load_profile(profile))

        self.assertIn("actual_release must be false", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
