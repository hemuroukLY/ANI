#!/usr/bin/env python3
"""Offline checks for M1-INSTANCE-I instance orchestration contracts."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "WorkloadProviderStatusReader",
    "LocalProviderStatusReader",
    "WorkloadProviderObservation",
    "provider apply must be applied",
    "KubeVirt",
    "Business services must not import provider SDKs",
    "WorkloadInstanceOrchestrator",
    "LocalInstanceOrchestrator",
    "WorkloadPlanAuditStore",
    "WorkloadProviderDryRun",
    "WorkloadProviderApply",
    "WorkloadStatusReconciler",
    "permission proof",
    "final workload status",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-i")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated instance orchestration contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
