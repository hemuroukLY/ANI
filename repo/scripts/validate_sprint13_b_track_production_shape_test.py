#!/usr/bin/env python3
"""Tests for Sprint 13 S01-S04 production-shaped evidence guard."""

from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path

import validate_sprint13_b_track_production_shape as guard


class Sprint13ProductionShapeGuardTest(unittest.TestCase):
    def test_repository_records_are_explicit_about_production_shape(self) -> None:
        guard.validate_all()

    def test_evidence_requires_production_shape_block(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            evidence = Path(tmp) / "evidence.json"
            evidence.write_text(json.dumps({"status": "passed"}) + "\n", encoding="utf-8")

            with self.assertRaises(SystemExit) as raised:
                guard.validate_evidence("S01", evidence)

        self.assertIn("must include production_shape", str(raised.exception))

    def test_pending_production_shape_requires_missing_items(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            evidence = Path(tmp) / "evidence.json"
            evidence.write_text(json.dumps({
                "status": "passed",
                "production_shape": {"status": "pending", "transport_profile": "lab_proxy"},
            }) + "\n", encoding="utf-8")

            with self.assertRaises(SystemExit) as raised:
                guard.validate_evidence("S01", evidence)

        self.assertIn("pending production_shape must list missing_items", str(raised.exception))

    def test_production_passed_rejects_lab_transport(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            evidence = Path(tmp) / "evidence.json"
            evidence.write_text(json.dumps({
                "status": "passed",
                "production_shape": {
                    "status": "passed",
                    "transport_profile": "lab_kubeconfig_and_dev_gateway",
                    "missing_items": [],
                },
            }) + "\n", encoding="utf-8")

            with self.assertRaises(SystemExit) as raised:
                guard.validate_evidence("S01", evidence)

        self.assertIn("production_shape passed cannot use lab_kubeconfig_and_dev_gateway", str(raised.exception))

    def test_production_passed_requires_proof_items(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            evidence = Path(tmp) / "evidence.json"
            evidence.write_text(json.dumps({
                "status": "passed",
                "production_shape": {
                    "status": "passed",
                    "transport_profile": "in_cluster_serviceaccount",
                    "missing_items": [],
                },
            }) + "\n", encoding="utf-8")

            with self.assertRaises(SystemExit) as raised:
                guard.validate_evidence("S01", evidence)

        self.assertIn("production_shape passed requires proof_items", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
