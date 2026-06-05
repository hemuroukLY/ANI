#!/usr/bin/env python3
"""Tests for Sprint 7 CORE-OFFLINE-A manifest validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_core_offline as offline


class CoreOfflineValidationTest(unittest.TestCase):
    def test_default_manifest_contains_only_core_artifacts(self) -> None:
        manifest = offline.load_manifest(offline.DEFAULT_MANIFEST)

        offline.validate_manifest(manifest)

        image_names = {image["name"] for image in manifest["images"]}
        self.assertIn("ani-gateway", image_names)
        self.assertIn("auth-service", image_names)
        self.assertNotIn("model-service", image_names)
        self.assertNotIn("kb-service", image_names)
        self.assertEqual("core", manifest["scope"])

    def test_validation_rejects_unpinned_image(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            manifest = Path(tmpdir) / "offline.yaml"
            manifest.write_text(
                """
scope: core
images:
  - name: ani-gateway
    image: harbor.ani.internal/ani/ani-gateway:latest
charts:
  - name: ani-platform
    path: deploy/helm/ani-platform
scripts:
  - path: scripts/validate_core_installer.py
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                offline.validate_manifest(offline.load_manifest(manifest))

        self.assertIn("must be pinned by digest", str(raised.exception))

    def test_validation_rejects_services_image(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            manifest = Path(tmpdir) / "offline.yaml"
            manifest.write_text(
                """
scope: core
images:
  - name: model-service
    image: harbor.ani.internal/ani/model-service@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
charts:
  - name: ani-platform
    path: deploy/helm/ani-platform
scripts:
  - path: scripts/validate_core_installer.py
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                offline.validate_manifest(offline.load_manifest(manifest))

        self.assertIn("forbidden Services image", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
