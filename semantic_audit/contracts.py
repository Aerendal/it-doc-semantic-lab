"""
Class-role registry loader.
P2 component: loads document_classes.yaml and role_contracts.yaml.
Status: stub — to be implemented in exp_001.
"""
from pathlib import Path
from typing import Any, Dict
import yaml


def load_class_registry(schema_dir: Path) -> Dict[str, Any]:
    """Load document class registry from schemas/document_classes.yaml."""
    p = schema_dir / "document_classes.yaml"
    return yaml.safe_load(p.read_text(encoding="utf-8"))


def load_role_contracts(schema_dir: Path) -> Dict[str, Any]:
    """Load role contracts from schemas/role_contracts.yaml."""
    p = schema_dir / "role_contracts.yaml"
    return yaml.safe_load(p.read_text(encoding="utf-8"))
