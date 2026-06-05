#!/usr/bin/env python3
"""Tests for Sprint 11 Core documentation consistency validator."""

from __future__ import annotations

from pathlib import Path
from tempfile import TemporaryDirectory

from validate_sprint11_core_doc_consistency import (
    BOUNDARY_MARKER,
    EXECUTION_ENV_MARKER,
    REQUIRED_MAKE_TARGETS,
    REQUIRED_RECORDS,
    SPRINT11_MARKER,
    validate_workspace,
)


def write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def build_workspace(root: Path, omit_marker: bool = False) -> None:
    marker = "" if omit_marker else f"{SPRINT11_MARKER}\n{BOUNDARY_MARKER}\n{EXECUTION_ENV_MARKER}\n"
    for path in (
        root / "ANI-DOCS-INDEX.md",
        root / "ANI-06-开发计划.md",
        root / "repo" / "CURRENT-SPRINT.md",
        root / "repo" / "README.md",
    ):
        write(path, marker)
    make_targets = "\n".join(f"{target}:\n\t@true" for target in REQUIRED_MAKE_TARGETS)
    write(root / "repo" / "Makefile", make_targets)
    records = marker + "\n".join(REQUIRED_RECORDS)
    write(root / "repo" / "development-records" / "README.md", records)


def expect_failure(root: Path, expected: str) -> None:
    try:
        validate_workspace(root)
    except SystemExit as exc:
        assert expected in str(exc), f"expected {expected!r}, got {exc!s}"
        return
    raise AssertionError(f"expected validation failure containing {expected!r}")


def test_valid_workspace() -> None:
    with TemporaryDirectory() as tmp:
        root = Path(tmp)
        build_workspace(root)
        validate_workspace(root)


def test_missing_marker_fails() -> None:
    with TemporaryDirectory() as tmp:
        root = Path(tmp)
        build_workspace(root, omit_marker=True)
        expect_failure(root, "missing Sprint 11 started marker")


if __name__ == "__main__":
    test_valid_workspace()
    test_missing_marker_fails()
    print("Sprint 11 documentation consistency validator tests passed")
