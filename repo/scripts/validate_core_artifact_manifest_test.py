#!/usr/bin/env python3
"""Tests for Sprint 10 CORE-ARTIFACT-MANIFEST-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_artifact_manifest as artifact_manifest


class CoreArtifactManifestValidationTest(unittest.TestCase):
    def test_default_manifest_covers_core_release_artifacts(self) -> None:
        manifest = artifact_manifest.load_manifest(artifact_manifest.DEFAULT_MANIFEST)

        artifact_manifest.validate_manifest(manifest)

        self.assertEqual("core", manifest["scope"])
        self.assertFalse(manifest["actual_release"])
        artifact_ids = {artifact["id"] for artifact in manifest["artifacts"]}
        self.assertIn("core-openapi", artifact_ids)
        self.assertIn("core-cli-source", artifact_ids)

    def test_validation_rejects_checksum_mismatch(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            manifest = Path(tmpdir) / "core-artifacts.yaml"
            wrong_digest = "0" * 64
            manifest.write_text(
                f"""
scope: core
version: sprint10
target_release: v1.0.0
actual_release: false
artifacts:
  - id: core-openapi
    path: api/openapi/v1.yaml
    sha256: "{wrong_digest}"
  - id: core-sdk-metadata
    path: sdks/core/sdk-metadata.json
    sha256: "{wrong_digest}"
  - id: core-cli-source
    path: cli/ani/main.go
    sha256: "{wrong_digest}"
  - id: core-offline-lock
    path: deploy/offline/core-package-lock.yaml
    sha256: "{wrong_digest}"
  - id: core-release-evidence
    path: deploy/release/core-release-evidence.yaml
    sha256: "{wrong_digest}"
verification:
  - command: make validate-core-api-compatibility
  - command: make validate-sdk-beta
  - command: make validate-core-cli
  - command: make validate-core-offline-pack
  - command: make validate-core-release-evidence
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                artifact_manifest.validate_manifest(artifact_manifest.load_manifest(manifest))

        self.assertIn("sha256 does not match", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
