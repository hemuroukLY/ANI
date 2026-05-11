#!/usr/bin/env python3
"""Offline checks for M1-INSTANCE-K provider adapter contracts."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "KubernetesProviderAdapter",
    "KubernetesProviderClient",
    "WorkloadProviderDryRun",
    "WorkloadProviderApply",
    "WorkloadProviderStatusReader",
    "Kubernetes Deployment",
    "Kubernetes Job",
    "KubeVirt VirtualMachine",
    "dryRun=All",
    "apply is disabled by default",
    "audit id and permission proof",
    "resource refs",
    "Business services must use ANI ports",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-k")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated provider adapter contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
