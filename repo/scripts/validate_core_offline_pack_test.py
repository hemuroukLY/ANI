#!/usr/bin/env python3
"""Tests for Sprint 8 CORE-OFFLINE-PACK-A package contract validation."""

from __future__ import annotations

import tempfile
import unittest
from hashlib import sha256
from pathlib import Path

import validate_core_offline_pack as offline_pack


class CoreOfflinePackValidationTest(unittest.TestCase):
    def test_default_lock_references_core_manifest_and_verification(self) -> None:
        lock = offline_pack.load_lock(offline_pack.DEFAULT_LOCK)

        offline_pack.validate_lock(lock)

        self.assertEqual("core", lock["scope"])
        self.assertEqual("deploy/offline/core-package.yaml", lock["source_manifest"])
        manifest = offline_pack.ROOT / lock["source_manifest"]
        actual_digest = sha256(manifest.read_bytes()).hexdigest()
        self.assertEqual(actual_digest, lock["source_manifest_sha256"])
        self.assertEqual(actual_digest, lock["artifact"]["sha256"])
        self.assertIn("make validate-core-offline", [step["command"] for step in lock["verification"]])

    def test_validation_rejects_missing_signature_checksum(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            lock = Path(tmpdir) / "core-package-lock.yaml"
            manifest = offline_pack.ROOT / "deploy/offline/core-package.yaml"
            actual_digest = sha256(manifest.read_bytes()).hexdigest()
            lock.write_text(
                f"""
scope: core
source_manifest: deploy/offline/core-package.yaml
source_manifest_sha256: "{actual_digest}"
artifact:
  name: ani-core-offline-bundle.tar.zst
verification:
  - id: offline-contract
    command: make validate-core-offline
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                offline_pack.validate_lock(offline_pack.load_lock(lock))

        self.assertIn("artifact.sha256 must be a 64 character hex string", str(raised.exception))

    def test_validation_rejects_source_manifest_checksum_mismatch(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            lock = Path(tmpdir) / "core-package-lock.yaml"
            wrong_digest = "0" * 64
            lock.write_text(
                f"""
scope: core
source_manifest: deploy/offline/core-package.yaml
source_manifest_sha256: "{wrong_digest}"
artifact:
  name: ani-core-offline-bundle.tar.zst
  sha256: "{wrong_digest}"
verification:
  - id: offline-contract
    command: make validate-core-offline
""".lstrip(),
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                offline_pack.validate_lock(offline_pack.load_lock(lock))

        self.assertIn("source_manifest_sha256 does not match", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
