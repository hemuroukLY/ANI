#!/usr/bin/env python3
"""Tests for Sprint 10 CORE-VERSION-POLICY-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_version_policy as version_policy


class CoreVersionPolicyValidationTest(unittest.TestCase):
    def test_default_policy_keeps_release_flags_false(self) -> None:
        policy = version_policy.load_policy(version_policy.DEFAULT_POLICY)

        version_policy.validate_policy(policy)

        self.assertEqual("core", policy["scope"])
        self.assertFalse(policy["actual_release"])
        self.assertFalse(policy["production_release"])

    def test_validation_rejects_actual_release_true(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            root = Path(tmpdir)
            release_dir = root / "repo" / "deploy" / "release"
            release_dir.mkdir(parents=True)
            (root / "repo").mkdir(exist_ok=True)
            (root / "ANI-DOCS-INDEX.md").write_text("不是实际 v1.0.0 发布\n", encoding="utf-8")
            (root / "ANI-06-开发计划.md").write_text("不是实际 v1.0.0 发布\n", encoding="utf-8")
            (root / "repo" / "CURRENT-SPRINT.md").write_text("不是实际 v1.0.0 发布\n", encoding="utf-8")
            policy = {
                "scope": "core",
                "version": "sprint10",
                "target_release": "v1.0.0",
                "actual_release": False,
                "release_candidate": False,
                "production_release": False,
                "rc_cut": False,
            }
            (release_dir / "bad.yaml").write_text("scope: core\nactual_release: true\n", encoding="utf-8")

            with self.assertRaises(SystemExit) as raised:
                version_policy.validate_policy(policy, root=root)

        self.assertIn("actual_release must not be true", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
