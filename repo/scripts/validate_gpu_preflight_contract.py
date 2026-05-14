#!/usr/bin/env python3
"""Offline contract checks for M1-INFRA-F GPU preflight manifests."""

from __future__ import annotations

import pathlib
import sys
from typing import Any

import yaml


ROOT = pathlib.Path("deploy/manifests/m1-infra-f")
REQUIRED_CONFIG_FRAGMENTS = (
    "ani.kubercloud.io/gpu-node=true",
    "gpu-vendor",
    "gpu-model",
    "queues.scheduling.volcano.sh",
    "ANI_GPU_REQUIRE_HAMI_ALLOCATABLE",
    "ANI_GPU_REQUIRE_DCGM_SERVICE",
)


def load_docs(path: pathlib.Path) -> list[dict[str, Any]]:
    with path.open("r", encoding="utf-8") as handle:
        docs = [doc for doc in yaml.safe_load_all(handle) if doc is not None]
    return [doc for doc in docs if isinstance(doc, dict)]


def metadata_name(doc: dict[str, Any]) -> str:
    metadata = doc.get("metadata")
    if not isinstance(metadata, dict):
        return ""
    name = metadata.get("name")
    return name if isinstance(name, str) else ""


def main() -> int:
    root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else ROOT
    errors: list[str] = []

    required_files = (
        "00-gpu-e2e-preflight-rbac.yaml",
        "10-gpu-e2e-preflight-config.yaml",
        "20-gpu-e2e-preflight-job.yaml",
        "30-gpu-smoke-workload-template.yaml",
    )
    for filename in required_files:
        if not (root / filename).exists():
            errors.append(f"{root / filename}: missing required file")

    docs: list[dict[str, Any]] = []
    for path in sorted(root.glob("*.yaml")):
        docs.extend(load_docs(path))

    kinds_by_name = {(doc.get("kind"), metadata_name(doc)) for doc in docs}
    for expected in (
        ("ServiceAccount", "ani-gpu-e2e-preflight"),
        ("ClusterRole", "ani-gpu-e2e-preflight"),
        ("ClusterRoleBinding", "ani-gpu-e2e-preflight"),
        ("ConfigMap", "ani-gpu-e2e-preflight"),
        ("Job", "ani-gpu-e2e-preflight"),
        ("ConfigMap", "ani-gpu-smoke-workload-template"),
    ):
        if expected not in kinds_by_name:
            errors.append(f"{root}: missing {expected[0]} {expected[1]}")

    configmaps = [
        doc
        for doc in docs
        if doc.get("kind") == "ConfigMap" and metadata_name(doc) == "ani-gpu-e2e-preflight"
    ]
    if not configmaps:
        errors.append(f"{root}: missing ani-gpu-e2e-preflight ConfigMap")
    else:
        data = configmaps[0].get("data")
        script = data.get("preflight.sh") if isinstance(data, dict) else None
        if not isinstance(script, str):
            errors.append(f"{root}: missing preflight.sh")
        if isinstance(data, dict):
            combined_data = "\n".join(str(value) for value in data.values())
            for fragment in REQUIRED_CONFIG_FRAGMENTS:
                if fragment not in combined_data:
                    errors.append(f"{root}: preflight ConfigMap missing {fragment}")

    jobs = [
        doc
        for doc in docs
        if doc.get("kind") == "Job" and metadata_name(doc) == "ani-gpu-e2e-preflight"
    ]
    if jobs:
        spec = jobs[0].get("spec", {})
        template = spec.get("template", {}) if isinstance(spec, dict) else {}
        pod_spec = template.get("spec", {}) if isinstance(template, dict) else {}
        service_account = pod_spec.get("serviceAccountName") if isinstance(pod_spec, dict) else None
        if service_account != "ani-gpu-e2e-preflight":
            errors.append(f"{root}: Job must use ani-gpu-e2e-preflight ServiceAccount")

    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    print(f"validated GPU preflight contract under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
