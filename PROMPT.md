# Prompt: Terraform Template Provider

Create a Terraform provider that enables template rendering with custom placeholder syntax. The provider should be designed to work alongside existing templating systems without conflicts.

## Problem Statement

There's a need to render templates in Terraform where some values come from Terraform configuration, but the templates themselves may contain syntax from other templating systems that must remain untouched. For example:

- **Argo WorkflowTemplates** use `{{}}` syntax for their own parameter substitution and workflow expressions
- **Shell scripts** use `${}` and `$()` for parameter expansion and command substitution  
- **Bash** uses `[[` and `]]` for conditional expressions, `<<` and `>>` for heredoc/redirection, `<>` for redirection, `&&` and `||` for logical operators, and `((` and `))` for arithmetic evaluation

When rendering a Kubernetes manifest or Argo WorkflowTemplate, you might need Terraform to inject values like image tags, dates, or namespaces, while the template also contains Argo's native `{{inputs.parameters.*}}` syntax that should be preserved for runtime evaluation.

## Requirements

The provider needs to:

- Render templates by replacing placeholders with values from Terraform configuration
- Use a placeholder syntax that doesn't conflict with common templating systems or shell scripting
- Preserve all non-placeholder content exactly as provided, including syntax from other systems
- Validate that all placeholders have corresponding values before rendering
- Provide clear error messages when placeholders are missing, preventing silent failures
- Function as a data source (read-only) since template rendering doesn't create infrastructure

## Desired User Experience

Users should be able to:

- Write templates with placeholders that get replaced by Terraform values
- Include those templates alongside Argo WorkflowTemplates, shell scripts, or other templating systems in the same file
- Receive immediate feedback when required placeholders are missing
- Use the rendered output seamlessly with other Terraform resources (e.g., passing to `kubernetes_manifest` or `yamldecode`)

## Technical Context

- The provider should use modern Terraform Plugin Framework patterns
- Follow contemporary Go development practices
- The placeholder syntax should be distinctive enough to avoid accidental matches in typical template content
- Error handling should be robust and user-friendly
- The implementation should be focused and minimal, doing one thing well

## Provider Naming

- The provider should be named `template`
- The data source should be named `render`
- This results in the Terraform invocation pattern: `data "template_render" "example" { ... }`

## Example Scenario

Consider an Argo WorkflowTemplate YAML file where:
- The namespace, image tag, and date ranges need to come from Terraform variables
- The workflow steps use Argo's `{{inputs.parameters.*}}` syntax for runtime parameter passing
- Both the Terraform placeholders and Argo syntax coexist in the same file
- After Terraform renders the template, the result is a valid WorkflowTemplate that can be deployed to Kubernetes

The provider should make this workflow natural and error-free.

Use the mcp sequential-thinking, the mcp exa and the mcp ref
