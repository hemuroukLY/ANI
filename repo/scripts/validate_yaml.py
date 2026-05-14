#!/usr/bin/env python3
"""Offline YAML syntax checks for generated deploy configuration."""

from __future__ import annotations

import pathlib
import sys

import yaml


def main() -> int:
    roots = [pathlib.Path(arg) for arg in sys.argv[1:]]
    if not roots:
        roots = [pathlib.Path(".")]

    files: list[pathlib.Path] = []
    for root in roots:
        if root.is_file() and root.suffix in (".yaml", ".yml"):
            files.append(root)
        elif root.is_dir():
            files.extend(sorted(root.rglob("*.yaml")))
            files.extend(sorted(root.rglob("*.yml")))

    if not files:
        print("no YAML files found", file=sys.stderr)
        return 1

    errors: list[str] = []
    count = 0
    for path in sorted(set(files)):
        try:
            with path.open("r", encoding="utf-8") as handle:
                docs = list(yaml.safe_load_all(handle))
        except yaml.YAMLError as err:
            errors.append(f"{path}: {err}")
            continue
        if not docs:
            errors.append(f"{path}: empty YAML file")
            continue
        count += 1

    if errors:
        for err in errors:
            print(err, file=sys.stderr)
        return 1

    print(f"validated {count} YAML files")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
