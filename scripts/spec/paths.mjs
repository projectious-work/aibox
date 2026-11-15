// Locations of the v1 specification material (issue #504). Contract files
// (schemas, catalog/, fixtures/) stay under spec/v1/ until each moves into its
// feature specification; change only this file when they do.
import path from 'node:path';
import {fileURLToPath} from 'node:url';

export const repo = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
export const docs = path.join(repo, 'docs/product/v1');
export const ledger = path.join(repo, 'ledger/v0-parity');
export const examples = path.join(repo, 'examples');
export const scripts = path.join(repo, 'scripts/spec');
export const roadmap = path.join(repo, 'roadmap.yaml');
export const roadmapSchema = path.join(repo, 'roadmap.schema.json');
export const contracts = path.join(repo, 'spec/v1');
export const toolCatalog = path.join(repo, 'spec/tool-catalog/v1');
export const nativeStarters = path.join(repo, 'spec/native-starters/v1');
export const workspaceCustomization = path.join(repo, 'spec/workspace-customization/v1');
export const templates = path.join(repo, 'templates/src');
