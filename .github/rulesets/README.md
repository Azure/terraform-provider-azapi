# GitHub Rulesets

This directory contains GitHub repository ruleset definitions in JSONC format (JSON with comments).

> **Note:** Changes to the ruleset JSON files in this repository are **not** applied automatically. You must manually import or update them via the GitHub web UI.

## Importing a Ruleset

1. Go to **Settings** → **Rules** → **Rulesets** in the repository.
2. Click **New ruleset** → **Import a ruleset**.
3. Strip comments from the JSONC file first (the importer requires plain JSON):

   ```bash
   npx strip-json-comments-cli .github/rulesets/restrict-acctests-branches.jsonc > /tmp/ruleset.json
   ```

4. Upload `/tmp/ruleset.json` in the import dialog.
5. Review the settings and click **Create**.

## Updating an Existing Ruleset

1. Go to **Settings** → **Rules** → **Rulesets** and open the ruleset to edit.
2. Make changes directly in the UI, **or** use the **Export ruleset** button to download, compare with the file in this repo, and delete and re-import if needed.

## Prerequisites

- Write access to the repository (Owner or Admin role)
- [Node.js](https://nodejs.org/) installed (for `npx strip-json-comments-cli`)
