#!/usr/bin/env python3
"""Tests for Sprint 9 CORE-RELEASE-EVIDENCE-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_release_evidence as evidence


class CoreReleaseEvidenceValidationTest(unittest.TestCase):
    def test_default_manifest_covers_core_rc_readiness_evidence(self) -> None:
        manifest = evidence.load_manifest(evidence.DEFAULT_MANIFEST)

        evidence.validate_manifest(manifest)

        self.assertEqual("core", manifest["scope"])
        self.assertTrue(manifest["rc_readiness"])
        self.assertFalse(manifest["release_candidate"])
        self.assertFalse(manifest["production_release"])
        evidence_ids = {entry["id"] for entry in manifest["evidence"]}
        self.assertIn("sprint8-core-release", evidence_ids)
        self.assertIn("core-cli-version", evidence_ids)

    def test_validation_rejects_secret_bearing_evidence(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            manifest = Path(tmpdir) / "core-release-evidence.yaml"
            manifest.write_text(
                """
scope: core
version: sprint9
target_release: v1.0.0
rc_readiness: true
release_candidate: false
production_release: false
evidence:
  - id: core-api-compatibility
    command: make validate-core-api-compatibility
    artifact: repo/api/core-v1-compatibility-baseline.yaml
  - id: architecture-boundary
    command: make validate-architecture
    artifact: repo/architecture/component-import-allowlist.yaml
  - id: doc-entrypoints
    command: make validate-doc-entrypoints
    artifact: ANI-DOCS-INDEX.md
  - id: sdk-beta
    command: make validate-sdk-beta
    artifact: repo/sdks/core/sdk-metadata.json
  - id: sdk-mock-smoke
    command: make validate-sdk-mock-smoke
    artifact: repo/sdks/core/
  - id: sprint8-core-release
    command: make validate-sprint8-core-release
    artifact: repo/deploy/release/core-hardening.yaml
  - id: core-offline-pack
    command: make validate-core-offline-pack
    artifact: repo/deploy/offline/core-package-lock.yaml
  - id: core-cli-version
    command: make validate-core-cli
    artifact: local-secret-token.txt
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                evidence.validate_manifest(evidence.load_manifest(manifest))

        self.assertIn("forbidden secret-bearing evidence", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
