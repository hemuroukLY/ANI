#!/usr/bin/env python3
"""Offline checks for M1-INSTANCE-A instance contracts."""

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
REQUIRED_PLANES = (
    "tenant_vpc",
    "foundation_mesh",
    "storage",
    "management",
    "public_ingress",
)
REQUIRED_STORAGE = (
    "root_disk",
    "data_disk",
    "shared_pvc",
    "object_fuse",
    "ephemeral",
)


def load_docs(root: pathlib.Path) -> list[dict[str, Any]]:
    docs: list[dict[str, Any]] = []
    for path in sorted(root.glob("*.yaml")):
        with path.open("r", encoding="utf-8") as handle:
            docs.extend(doc for doc in yaml.safe_load_all(handle) if isinstance(doc, dict))
    return docs


def data_for(docs: list[dict[str, Any]], name: str) -> dict[str, str]:
    for doc in docs:
        if doc.get("kind") != "ConfigMap":
            continue
        metadata = doc.get("metadata")
        if isinstance(metadata, dict) and metadata.get("name") == name:
            data = doc.get("data")
            return data if isinstance(data, dict) else {}
    return {}


def main() -> int:
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-instance-a")
    docs = load_docs(root)
    errors: list[str] = []

    object_contract = data_for(docs, "ani-instance-object-contract")
    network_plan = data_for(docs, "ani-instance-network-plan")
    storage_plan = data_for(docs, "ani-instance-storage-plan")
    examples = data_for(docs, "ani-instance-examples")

    for name, data in (
        ("ani-instance-object-contract", object_contract),
        ("ani-instance-network-plan", network_plan),
        ("ani-instance-storage-plan", storage_plan),
        ("ani-instance-examples", examples),
    ):
        if not data:
            errors.append(f"missing {name} ConfigMap")

    combined_object = "\n".join(str(value) for value in object_contract.values())
    for kind in REQUIRED_KINDS:
        if kind not in combined_object:
            errors.append(f"instance object contract missing kind {kind}")

    combined_network = "\n".join(str(value) for value in network_plan.values())
    for plane in REQUIRED_PLANES:
        if plane not in combined_network:
            errors.append(f"instance network plan missing plane {plane}")

    combined_storage = "\n".join(str(value) for value in storage_plan.values())
    for storage in REQUIRED_STORAGE:
        if storage not in combined_storage:
            errors.append(f"instance storage plan missing attachment kind {storage}")

    combined_examples = "\n".join([*examples.keys(), *[str(value) for value in examples.values()]])
    for keyword in ("vm-instance", "container-instance", "gpu-container-instance", "tenant_vpc", "foundation_mesh"):
        if keyword not in combined_examples:
            errors.append(f"instance examples missing {keyword}")

    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated instance contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
