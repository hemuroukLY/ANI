#!/usr/bin/env python3
"""Offline checks for M1-RUNTIME-A workload runtime contracts."""

from __future__ import annotations

import pathlib
import sys
from typing import Any

import yaml


REQUIRED_KINDS = (
    "vm",
    "container",
    "gpu_container",
    "inference",
    "notebook",
    "agent_sandbox",
    "batch_job",
)
REQUIRED_STATES = ("pending", "provisioning", "running", "stopped", "failed", "deleting")


def load_docs(root: pathlib.Path) -> list[dict[str, Any]]:
    docs: list[dict[str, Any]] = []
    for path in sorted(root.glob("*.yaml")):
        with path.open("r", encoding="utf-8") as handle:
            docs.extend(doc for doc in yaml.safe_load_all(handle) if isinstance(doc, dict))
    return docs


def configmap_data(docs: list[dict[str, Any]], name: str) -> dict[str, str]:
    for doc in docs:
        if doc.get("kind") != "ConfigMap":
            continue
        metadata = doc.get("metadata")
        if isinstance(metadata, dict) and metadata.get("name") == name:
            data = doc.get("data")
            return data if isinstance(data, dict) else {}
    return {}


def main() -> int:
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-runtime-a")
    docs = load_docs(root)
    errors: list[str] = []

    runtime = configmap_data(docs, "ani-workload-runtime-contract")
    runtime_class = configmap_data(docs, "ani-runtime-class-contract")
    instance_schema = configmap_data(docs, "ani-instance-schema")

    if not runtime:
        errors.append("missing ani-workload-runtime-contract ConfigMap")
    if not runtime_class:
        errors.append("missing ani-runtime-class-contract ConfigMap")
    if not instance_schema:
        errors.append("missing ani-instance-schema ConfigMap")

    combined_runtime = "\n".join(str(value) for value in runtime.values())
    for kind in REQUIRED_KINDS:
        if kind not in combined_runtime:
            errors.append(f"runtime contract missing workload kind {kind}")

    combined_runtime_class = "\n".join(str(value) for value in runtime_class.values())
    for state in REQUIRED_STATES:
        if state not in combined_runtime_class:
            errors.append(f"runtime contract missing lifecycle state {state}")

    combined_instance = "\n".join(str(value) for value in instance_schema.values())
    for field in ("tenantID", "instanceID", "providerID", "gpu_container", "tenantIsolated"):
        if field not in combined_instance:
            errors.append(f"instance schema missing {field}")

    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated workload runtime contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
