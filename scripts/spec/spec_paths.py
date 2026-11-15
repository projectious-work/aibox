"""Locations of the v1 specification material (issue #504).

Contract files (schemas, catalog/, fixtures/) stay under spec/v1/ until each
moves into its feature specification; change only this file when they do.
"""
from pathlib import Path

REPO = Path(__file__).resolve().parents[2]
DOCS = REPO / "docs/product/v1"
LEDGER = REPO / "ledger/v0-parity"
EXAMPLES = REPO / "examples"
SCRIPTS = REPO / "scripts/spec"
ROADMAP = REPO / "roadmap.yaml"
ROADMAP_SCHEMA = REPO / "roadmap.schema.json"
CONTRACTS = REPO / "spec/v1"
NATIVE_STARTERS = REPO / "spec/native-starters/v1"
TEMPLATES = REPO / "templates/src"
