#!/usr/bin/env python3
"""Validate Sprint 5 REAL-K8S-LAB-A real provider lab gate."""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Any

import yaml


ROOT = Path(__file__).resolve().parents[1]
DOC_ROOT = ROOT.parent
PROFILE = ROOT / "deploy/real-k8s-lab/profile.yaml"
REQUIRED_COMPONENTS = {"kubernetes", "kube_ovn", "kubevirt", "vcluster"}
REQUIRED_DOC_TOKENS = [
    "REAL-K8S-LAB-A",
    "validate-real-k8s-profile",
    "Kube-OVN",
    "KubeVirt",
    "vCluster",
    "local profile",
]


def fail(message: str) -> None:
    raise SystemExit(f"real k8s profile invalid: {message}")


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def load_profile(path: Path) -> dict[str, Any]:
    if not path.exists():
        fail(f"missing {path.relative_to(ROOT)}")
    with path.open(encoding="utf-8") as handle:
        data = yaml.safe_load(handle)
    if not isinstance(data, dict):
        fail(f"{path.relative_to(ROOT)} must be a YAML object")
    return data


def validate_contract(profile: dict[str, Any]) -> None:
    if profile.get("profile") != "REAL-K8S-LAB-A":
        fail("profile must be REAL-K8S-LAB-A")
    if profile.get("status") not in {"contract", "live", "production_like"}:
        fail("status must be contract, live or production_like")
    if int(profile.get("minimum_nodes", 0)) < 3:
        fail("minimum_nodes must be at least 3")
    components = profile.get("components")
    if not isinstance(components, dict):
        fail("components must be an object")
    missing = REQUIRED_COMPONENTS - set(components)
    if missing:
        fail(f"missing required components: {', '.join(sorted(missing))}")
    for name, component in components.items():
        if not isinstance(component, dict):
            fail(f"{name} component must be an object")
        if "purpose" not in component:
            fail(f"{name} must document purpose")
        checks = component.get("live_checks")
        if not isinstance(checks, list) or not checks:
            fail(f"{name} must define live_checks")
        for check in checks:
            if not isinstance(check, dict):
                fail(f"{name} live check must be an object")
            for field in ("id", "command", "pass_condition"):
                if not check.get(field):
                    fail(f"{name} live check missing {field}")


def validate_docs() -> None:
    docs = {
        "CLAUDE.md": DOC_ROOT / "CLAUDE.md",
        "ANI-DOCS-INDEX.md": DOC_ROOT / "ANI-DOCS-INDEX.md",
        "ANI-06-开发计划.md": DOC_ROOT / "ANI-06-开发计划.md",
        "CURRENT-SPRINT.md": ROOT / "CURRENT-SPRINT.md",
        "development-records/README.md": ROOT / "development-records/README.md",
    }
    for label, path in docs.items():
        content = read(path)
        for token in REQUIRED_DOC_TOKENS:
            if token not in content:
                fail(f"{label} must reference {token}")


def kubectl(args: list[str], kubeconfig: str | None) -> subprocess.CompletedProcess[str]:
    command = ["kubectl", *args]
    env = os.environ.copy()
    if kubeconfig:
        env["KUBECONFIG"] = kubeconfig
    return subprocess.run(command, env=env, text=True, capture_output=True, check=False)


def run_json_check(command: str, kubeconfig: str | None) -> Any:
    if not command.startswith("kubectl "):
        fail(f"live command must start with kubectl: {command}")
    result = kubectl(command.split()[1:], kubeconfig)
    if result.returncode != 0:
        fail(f"{command} failed: {result.stderr.strip() or result.stdout.strip()}")
    try:
        return json.loads(result.stdout)
    except json.JSONDecodeError as err:
        fail(f"{command} did not return JSON: {err}")


def condition_passed(condition: str, command: str, kubeconfig: str | None, minimum_nodes: int) -> bool:
    if condition == "stdout_yes":
        result = kubectl(command.split()[1:], kubeconfig)
        return result.returncode == 0 and result.stdout.strip().lower() == "yes"

    data = run_json_check(command, kubeconfig)
    if condition == "at_least_minimum_nodes_ready":
        items = data.get("items", [])
        ready = 0
        for node in items:
            conditions = node.get("status", {}).get("conditions", [])
            if any(item.get("type") == "Ready" and item.get("status") == "True" for item in conditions):
                ready += 1
        return ready >= minimum_nodes
    if condition == "at_least_one_storageclass":
        return len(data.get("items", [])) >= 1
    if condition == "crd_exists":
        return bool(data.get("metadata", {}).get("name"))
    if condition == "at_least_one_item":
        return len(data.get("items", [])) >= 1
    fail(f"unsupported pass_condition {condition}")


def validate_live(profile: dict[str, Any], kubeconfig: str | None) -> None:
    if shutil.which("kubectl") is None:
        fail("kubectl is required for --live")
    minimum_nodes = int(profile.get("minimum_nodes", 3))
    for component_name, component in profile["components"].items():
        if not component.get("required", False):
            continue
        for check in component.get("live_checks", []):
            if not condition_passed(check["pass_condition"], check["command"], kubeconfig, minimum_nodes):
                fail(f"{component_name}/{check['id']} did not satisfy {check['pass_condition']}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--profile", default=str(PROFILE), help="real lab profile YAML")
    parser.add_argument("--live", action="store_true", help="run live kubectl checks")
    parser.add_argument("--kubeconfig", default=os.getenv("KUBECONFIG"), help="kubeconfig for live checks")
    args = parser.parse_args()

    profile = load_profile(Path(args.profile))
    validate_contract(profile)
    validate_docs()
    if args.live:
        validate_live(profile, args.kubeconfig)
        print("REAL-K8S-LAB-A live checks valid")
    else:
        print("REAL-K8S-LAB-A contract valid; use --live with KUBECONFIG to verify a real lab")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
