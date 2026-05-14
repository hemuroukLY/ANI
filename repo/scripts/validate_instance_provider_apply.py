#!/usr/bin/env python3
"""Offline checks for M1-INSTANCE-G provider apply gate contracts."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "WorkloadProviderApply",
    "LocalProviderApply",
    "disabled by default",
    "fail closed",
    "audit id",
    "permission proof",
    "admission must be allowed",
    "provider dry-run must be accepted",
    "create only",
    "server-side dry-run",
    "user permission checks",
    "tenant scoping",
    "resourceRefs",
)


def load_text(root: pathlib.Path) -> str:
    chunks: list[str] = []
    for path in sorted(root.glob("*.yaml")):
        with path.open("r", encoding="utf-8") as handle:
            for doc in yaml.safe_load_all(handle):
                if isinstance(doc, dict):
                    chunks.append(str(doc))
    return "\n".join(chunks)


def main() -> int:
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-g")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated instance provider apply gate contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
