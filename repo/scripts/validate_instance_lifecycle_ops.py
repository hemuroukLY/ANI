#!/usr/bin/env python3
"""Offline checks for M1-INSTANCE-M lifecycle and visual ops contracts."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "WorkloadInstanceService",
    "Start",
    "Stop",
    "Restart",
    "Resize",
    "Delete",
    "WorkloadInstanceStore",
    "WorkloadInstanceOps",
    "LocalInstanceOpsGuard",
    "WorkloadInstanceService.Ops",
    "logs",
    "events",
    "metrics",
    "terminal",
    "exec",
    "ops are disabled by default",
    "terminal and exec are container-only",
    "business services must not call Kubernetes exec/logs/events APIs directly",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-m")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated instance lifecycle and ops contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
