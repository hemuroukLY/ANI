#!/usr/bin/env python3
"""Offline Kubernetes YAML sanity checks for generated manifests."""

from __future__ import annotations

import pathlib
import sys

import yaml


REQUIRED_FIELDS = ("apiVersion", "kind", "metadata")


def main() -> int:
    root = pathlib.Path(sys.argv[1] if len(sys.argv) > 1 else ".")
    files = sorted(root.glob("*.yaml")) + sorted(root.glob("*.yml"))
    if not files:
        print(f"no YAML files found under {root}", file=sys.stderr)
        return 1

    errors: list[str] = []
    count = 0
    for path in files:
        with path.open("r", encoding="utf-8") as handle:
            docs = list(yaml.safe_load_all(handle))
        for index, doc in enumerate(docs, start=1):
            if doc is None:
                continue
            count += 1
            if not isinstance(doc, dict):
                errors.append(f"{path}:{index}: document must be a mapping")
                continue
            for field in REQUIRED_FIELDS:
                if field not in doc:
                    errors.append(f"{path}:{index}: missing {field}")
            metadata = doc.get("metadata")
            if isinstance(metadata, dict) and not metadata.get("name"):
                errors.append(f"{path}:{index}: metadata.name is required")

    if count == 0:
        errors.append(f"{root}: no Kubernetes documents found")

    if errors:
        for err in errors:
            print(err, file=sys.stderr)
        return 1

    print(f"validated {count} Kubernetes YAML documents under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
