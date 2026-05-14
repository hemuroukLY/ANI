#!/usr/bin/env python3
"""Offline checks for M1-GPU-A heterogeneous GPU contracts."""

from __future__ import annotations

import pathlib
import sys
from typing import Any

import yaml


REQUIRED_VENDORS = ("nvidia", "huawei", "hygon")
REQUIRED_LABELS = (
    "ani.kubercloud.io/gpu-node",
    "ani.kubercloud.io/gpu-vendor",
    "ani.kubercloud.io/gpu-model",
    "ani.kubercloud.io/gpu-pool",
    "ani.kubercloud.io/gpu-kernel-compatible",
)
REQUIRED_FIELDS = (
    "resourceName",
    "runtimeClassName",
    "schedulerName",
    "queueName",
)


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
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("deploy/manifests/m1-gpu-a")
    docs = load_docs(root)
    errors: list[str] = []

    contract = configmap_data(docs, "ani-heterogeneous-gpu-contract")
    node_schema = configmap_data(docs, "ani-gpu-node-class-schema")
    compatibility = configmap_data(docs, "ani-gpu-runtime-compatibility-matrix")

    if not contract:
        errors.append("missing ani-heterogeneous-gpu-contract ConfigMap")
    if not node_schema:
        errors.append("missing ani-gpu-node-class-schema ConfigMap")
    if not compatibility:
        errors.append("missing ani-gpu-runtime-compatibility-matrix ConfigMap")

    combined_contract = "\n".join(str(value) for value in contract.values())
    for vendor in REQUIRED_VENDORS:
        if vendor not in combined_contract:
            errors.append(f"GPU contract missing vendor {vendor}")
    for label in REQUIRED_LABELS:
        if label not in combined_contract:
            errors.append(f"GPU contract missing label {label}")

    combined_schema = "\n".join(str(value) for value in node_schema.values())
    for field in REQUIRED_FIELDS:
        if field not in combined_schema:
            errors.append(f"GPU node schema missing scheduling field {field}")

    combined_compatibility = "\n".join(str(value) for value in compatibility.values())
    for keyword in ("kernel", "driver", "device plugin", "runtime", "isolate"):
        if keyword not in combined_compatibility:
            errors.append(f"GPU compatibility policy missing {keyword}")

    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated heterogeneous GPU contracts under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
