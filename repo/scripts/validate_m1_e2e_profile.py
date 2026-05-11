#!/usr/bin/env python3
"""Offline checks for the M1 end-to-end integration profile."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "M1 end-to-end integration profile",
    "VM",
    "container",
    "GPU container",
    "WorkloadInstanceService.Create",
    "WorkloadInstanceOrchestrator",
    "WorkloadRuntime",
    "WorkloadRenderer",
    "WorkloadAdmission",
    "WorkloadPlanAuditStore",
    "WorkloadProviderDryRun",
    "WorkloadProviderApply",
    "WorkloadProviderStatusReader",
    "WorkloadStatusReconciler",
    "WorkloadInstanceStore",
    "WorkloadInstanceService.Start",
    "WorkloadInstanceService.Stop",
    "WorkloadInstanceService.Restart",
    "WorkloadInstanceService.Resize",
    "WorkloadInstanceService.Delete",
    "WorkloadInstanceService.Get",
    "WorkloadInstanceService.List",
    "WorkloadInstanceService.Ops",
    "WorkloadInstanceOps",
    "GPUInventory",
    "KubeVirt VirtualMachine",
    "Kubernetes Deployment",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-e2e")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated M1 e2e profile under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
