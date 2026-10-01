#!/usr/bin/env python3
"""Check the frozen v0 census against the V1-05 native selection review."""
import json
from collections import Counter
from pathlib import Path

root = Path(__file__).resolve().parents[1]
ledger = json.loads((root / "ledger/addons.json").read_text())
table = json.loads((root / "ledger/v1-05-tool-selection.json").read_text())
audit = json.loads((root / "evidence/V1-05/source-audit.json").read_text())
oci = json.loads((root / "evidence/V1-05/oci-evidence.json").read_text())
fixture = json.loads((root / "fixtures/v1-05-native-features.json").read_text())

expected = {f"{recipe['id']}/{tool['name']}": (recipe, tool)
            for recipe in ledger for tool in recipe['tools']}
rows = table['rows']
actual = {row['id']: row for row in rows}
assert len(expected) == len(rows) == len(actual) == 98
assert set(actual) == set(expected)
assert table['baselineRevision'] == audit['baseline_revision']
assert {row['id'] for row in audit['rows']} == set(expected)
audited = {row['id']: row for row in audit['rows']}

for key, row in actual.items():
    source, _ = expected[key]
    assert row['v0-source'] == source['source']
    assert row['v0-source-sha256'] == source['sha256']
    assert row['classification'] in ('qualified_upstream', 'gap', 'uncertain')
    assert row['classification'] == audited[key]['status']
    if row['classification'] == 'qualified_upstream':
        assert row['v1-ref'].endswith('@' + row['digest'])
        assert row['v1-option'] and row['install-owner'] and row['license']
        assert row['platforms'] == ['linux/amd64', 'linux/arm64']
        name = row['v1-ref'].split('/features/')[1].split('@')[0]
        artifact = oci[name]
        assert row['digest'] == artifact['manifest_digest']
        assert row['install-owner'] == artifact['publisher']
        assert set(row['v1-option']) <= set(artifact['metadata']['options'])
    else:
        assert row['v1-ref'] is None and row['digest'] is None
        assert row['decision'] == 'owner-review-pending'
        assert row['proposed-route'] and row['proposed-route']['proposal']

assert len(fixture['selected']) == 1
for ref, options in fixture['selected'].items():
    matches = [row for row in rows if row['v1-ref'] == ref and row['v1-option'] == options]
    assert len(matches) == 1, f'example Feature is not in the reviewed selection table: {ref}'

counts = Counter(row['classification'] for row in rows)
assert dict(counts) == audit['counts']
print(f"V1-05 catalog: {len(rows)} rows, {counts['qualified_upstream']} source-reviewed upstream Features, "
      f"{counts['gap'] + counts['uncertain']} owner-review rows")
