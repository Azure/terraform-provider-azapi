# AzAPI provider agents

The `.x` directory contains the agent definitions, configuration and tools for triaging and analysing issues in `terraform-provider-azapi`. Keeping these files here lets maintainers review and update the agents alongside the provider code.

## How it works

When enabled, the X Engineering Agent loads this package from an approved repository commit and automatically discovers its agents and tools. Each run uses a fixed version of the package so subsequent edits do not change work already in progress.

## Updating the package

- `definitions/`: agent instructions and responsibilities.
- `tools/`: tools used by each agent role.
- `x.yml`: repository settings and workflow routing.

With the X Engineering Agent installed, validate the package from the repository root:

```sh
python -m x_engineering_agent.repository_packages --root .x --repository Azure/terraform-provider-azapi --revision "$(git rev-parse HEAD)"
```
