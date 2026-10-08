resource "azapi_data_plane_resource" "evaluation_run" {
  type = "Microsoft.Foundry/evaluation/runs@2025-05-01"

  parent_id = "foundry-account-id.services.ai.azure.com/api/projects/foundry-project-id"

  evaluation_id = "eval_id"

  body = {
    name = "evaluation-run-name"

    data_source = {
      type = "azure_ai_target_completions"

      source = {
        type = "file_id"
        id   = "azureai://accounts/foundry-account-id/projects/foundry-project-id/data/dataset/versions/1"
      }

      input_messages = {
        type = "template"

        template = [{
          role    = "user"
          content = "{{item.input}}"
        }]
      }

      target = {
        type = "azure_ai_agent"
        name = "agent-name"
      }
    }

    evaluation_level = "turn"
  }
}