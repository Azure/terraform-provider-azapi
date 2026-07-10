# Updating Bicep Types

The AzAPI provider relies on generated Bicep type definitions under `internal/azure/generated` to determine import API versions and to perform schema validation. These files should be synchronized periodically with the [azure-rest-api-specs](https://github.com/Azure/azure-rest-api-specs).

## Update workflow

The repository includes a GitHub Actions workflow at `.github/workflows/bicep-types-update.yaml` that automates this process.

The workflow is configured to run on a weekly schedule, and pushes the generated changes to a temporary branch. A PR should then be created from that branch into `main` for review and merging.

## Manual trigger

In urgent situation where the most recent run is not sufficient, the workflow can be triggered manually by creating a branch named `manual-push/bicep-types-update` and pushing an empty commit:

```bash
git commit --allow-empty -m "Trigger bicep types update workflow"
```
