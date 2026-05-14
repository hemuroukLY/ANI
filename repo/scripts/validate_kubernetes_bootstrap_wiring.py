#!/usr/bin/env python3
"""Offline checks for Kubernetes provider bootstrap wiring."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_KEYWORDS = (
    "M1-INSTANCE-P",
    "Kubernetes provider bootstrap wiring",
    "WORKLOAD_PROVIDER",
    "kubernetes_rest",
    "WORKLOAD_PROVIDER_APPLY_ENABLED",
    "KUBERNETES_API_HOST",
    "KUBERNETES_BEARER_TOKEN",
    "KUBERNETES_PROVIDER_FIELD_MANAGER",
    "KubernetesRESTClient",
    "KubernetesProviderAdapter",
    "apply disabled by default",
    "local provider remains default",
    "WorkloadProviderDryRun",
    "WorkloadProviderApply",
    "WorkloadProviderStatusReader",
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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-p")
    text = load_text(root)
    errors = [f"{root}: missing {keyword}" for keyword in REQUIRED_KEYWORDS if keyword not in text]
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated Kubernetes bootstrap wiring profile under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
