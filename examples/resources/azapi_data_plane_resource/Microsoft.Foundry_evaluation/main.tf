resource "azapi_data_plane_resource" "evaluation" {
  type = "Microsoft.Foundry/evaluation/versions@2025-05-01"

  parent_id = "foundry-account-id.services.ai.azure.com/api/projects/foundry-project-id"

  body = {
    name = "evaluation-name"

    data_source_config = {
      type = "custom"

      item_schema = {
        type = "object"

        properties = {
          input           = { type = "string" }
          expected_output = { type = "string" }
        }

        required = []
      }

      include_sample_schema = false
    }

    testing_criteria = [{
      type           = "azure_ai_evaluator"
      name           = "TaskCompletion"
      evaluator_name = "builtin.task_completion"

      data_mapping = {
        query            = "{{item.input}}"
        response         = "{{sample.output_text}}"
        ground_truth     = "{{item.expected_output}}"
        tool_calls       = "{{sample.tool_calls}}"
        tool_definitions = "{{sample.tool_definitions}}"
      }

      initialization_parameters = {
        deployment_name = "gpt5.4-nano"
      }
    }]
  }
}