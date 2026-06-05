#!/usr/bin/env python3
"""Tests for Sprint 10 CORE-FINAL-DOC-CONSISTENCY-A validation."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import validate_sprint10_core_doc_consistency as doc_consistency


class Sprint10CoreDocConsistencyValidationTest(unittest.TestCase):
    def test_default_docs_and_makefile_are_sprint10_aligned(self) -> None:
        doc_consistency.validate_workspace(doc_consistency.ROOT)

    def test_validation_rejects_stale_sprint9_current_marker(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            root = Path(tmpdir)
            (root / "repo" / "development-records").mkdir(parents=True)
            (root / "ANI-DOCS-INDEX.md").write_text(
                "Sprint 10 Core-only 代码开发已完成\n不是实际 v1.0.0 发布\n", encoding="utf-8"
            )
            (root / "ANI-06-开发计划.md").write_text(
                "Sprint 10 Core-only 代码开发已完成\n不是实际 v1.0.0 发布\n", encoding="utf-8"
            )
            (root / "repo" / "CURRENT-SPRINT.md").write_text(
                "当前重心：Sprint 9 Core-only 代码开发已完成\n不是实际 v1.0.0 发布\n", encoding="utf-8"
            )
            (root / "repo" / "Makefile").write_text("validate-sprint10-release-prep:\n", encoding="utf-8")
            (root / "repo" / "development-records" / "README.md").write_text(
                "Sprint 10 Core-only 代码开发已完成\n不是实际 v1.0.0 发布\nSPRINT10-CLOSURE-A\n",
                encoding="utf-8",
            )

            with self.assertRaises(SystemExit) as raised:
                doc_consistency.validate_workspace(root)

        self.assertIn("stale Sprint 9 current marker", str(raised.exception))


if __name__ == "__main__":
    unittest.main()
