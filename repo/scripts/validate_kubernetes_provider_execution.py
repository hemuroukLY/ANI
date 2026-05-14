#!/usr/bin/env python3
"""Offline checks for the M1 Kubernetes provider execution profile."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "M1-INSTANCE-N",
    "Kubernetes provider execution profile",
    "KubernetesProviderAdapter",
    "KubernetesProviderClient",
    "KubernetesProviderClient.ServerSideDryRun",
    "dryRun=All",
    "KubernetesProviderClient.Apply",
    "KubernetesProviderClient.Observe",
    "WorkloadInstanceOrchestrator",
    "WorkloadPlanAuditStore",
    "WorkloadAdmission",
    "WorkloadProviderDryRun",
    "WorkloadProviderApply",
    "WorkloadProviderStatusReader",
    "WorkloadStatusReconciler",
    "WorkloadInstanceStore",
    "apply disabled by default",
    "audit id",
    "permission proof",
    "resource refs",
    "adapter-owned package",
    "Business services must not import client-go",
    "KubeVirt client",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-n")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated Kubernetes provider execution profile under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
