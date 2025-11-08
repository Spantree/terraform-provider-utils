# Terraform Spantree Utils Provider

A Terraform provider offering utility functions including template rendering with custom placeholder syntax that doesn't conflict with other templating systems.

## Why This Provider?

When working with Terraform to deploy resources like Argo Workflows, Kubernetes manifests, or shell scripts, you often need to inject Terraform-managed values (like namespaces, image tags, or dates) into templates that also contain their own templating syntax. Standard Terraform templating can conflict with:

- **Argo Workflows**: `{{inputs.parameters.*}}` syntax
- **Shell scripts**: `${VAR}`, `$(command)` syntax
- **Bash**: `[[ ]]`, `<<`, `>>`, `<>`, `&&`, `||`, `(( ))` operators
- **Windows batch files**: `%%VAR%%` syntax

This provider solves that problem by using a distinctive placeholder syntax: `@@VAR@@`

## Features

- **Conflict-free syntax**: `@@VAR@@` placeholders don't interfere with other templating systems
- **Validation**: All placeholders must have corresponding values or an error is returned
- **Simple**: Pure data source, no infrastructure created
- **Fast**: Template rendering happens locally during plan/apply

## Installation

### Terraform 0.13+

Add to your `terraform` block:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source  = "spantree/utils"
      version = "~> 1.0"
    }
  }
}

provider "spantree_utils" {}
```

### Local Development

For local development and testing:

1. Build the provider:

   ```bash
   go build -o terraform-provider-template
   ```

2. Create a local provider override in `~/.terraformrc`:

   ```hcl
   provider_installation {
     dev_overrides {
       "spantree/utils" = "/path/to/provider/directory"
     }
     direct {}
   }
   ```

## Usage

### Basic Example

```hcl
data "spantree_utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@!"

  values = {
    NAME  = "Alice"
    PLACE = "Wonderland"
  }
}

output "greeting" {
  value = data.spantree_utils_render_template.greeting.result
  # Output: "Hello Alice, welcome to Wonderland!"
}
```

### Configuration File Example

```hcl
data "spantree_utils_render_template" "config" {
  template = <<-EOT
    server {
      host = "@@HOST@@"
      port = @@PORT@@
      environment = "@@ENV@@"
    }
  EOT

  values = {
    HOST = "localhost"
    PORT = "8080"
    ENV  = "production"
  }
}
```

### Shell Script with Preserved Syntax

```hcl
data "spantree_utils_render_template" "script" {
  template = <<-EOT
    #!/bin/bash
    
    # Terraform-injected values
    NAMESPACE=@@NAMESPACE@@
    IMAGE_TAG=@@IMAGE_TAG@@
    
    # Shell variables remain as-is
    echo "Namespace: $${NAMESPACE}"
    echo "Tag: $(echo $IMAGE_TAG)"
    
    # Bash conditionals work fine
    if [[ -n "$NAMESPACE" ]]; then
      echo "Valid"
    fi
  EOT

  values = {
    NAMESPACE = "production"
    IMAGE_TAG = "v1.2.3"
  }
}
```

### Argo Workflow Template Example

This is the primary use case - deploying Argo WorkflowTemplates where you need both Terraform values and Argo's runtime parameters:

```hcl
variable "namespace" {
  default = "argo"
}

variable "image_tag" {
  default = "v1.0.0"
}

data "spantree_utils_render_template" "workflow" {
  template = file("${path.module}/workflow.yaml")

  values = {
    NAMESPACE = var.namespace
    IMAGE_TAG = var.image_tag
  }
}

resource "kubernetes_manifest" "workflow" {
  manifest = yamldecode(data.spantree_utils_render_template.workflow.result)
}
```

Where `workflow.yaml` contains:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: my-workflow
  namespace: @@NAMESPACE@@
spec:
  templates:
    - name: process
      inputs:
        parameters:
          - name: input
      container:
        image: myapp:@@IMAGE_TAG@@
        command: [sh, -c]
        args:
          # Argo's {{}} syntax is preserved
          - echo "Processing {{inputs.parameters.input}}"
```

After rendering, the `@@NAMESPACE@@` and `@@IMAGE_TAG@@` are replaced with Terraform values, while `{{inputs.parameters.input}}` remains intact for Argo to evaluate at runtime.

## Data Source: `spantree_utils_render_template`

### Arguments

- `template` (String, Required) - The template string containing placeholders in `@@VAR@@` format. Placeholder names must be alphanumeric or underscore.
- `values` (Map of String, Required) - Map of placeholder names to their replacement values. All placeholders in the template must have a corresponding value.

### Attributes

- `result` (String) - The rendered template with all placeholders replaced by their values.

## Placeholder Syntax

Placeholders use the format: `@@VARIABLE_NAME@@`

**Rules:**

- Placeholder names must be alphanumeric or underscore: `[A-Za-z0-9_]+`
- All placeholders in the template must have a corresponding value in the `values` map
- Extra values in the map (not used in the template) are ignored
- The same placeholder can appear multiple times in the template

**Valid placeholders:**

- `@@NAME@@`
- `@@IMAGE_TAG@@`
- `@@my_var_123@@`

**Invalid placeholders:**

- `@@my-var@@` (hyphens not allowed)
- `@@my.var@@` (dots not allowed)
- `@@my var@@` (spaces not allowed)

## Error Handling

The provider will return an error if:

1. **Missing placeholder values**: Any placeholder in the template doesn't have a corresponding value

   ```
   Error: Template Rendering Failed
   Failed to render template: missing values for placeholders: NAMESPACE, IMAGE_TAG
   ```

2. **Invalid configuration**: Required arguments are missing or have wrong types

## Comparison with Other Solutions

### vs. Terraform's built-in `templatefile()`

Terraform's `templatefile()` function uses `${}` syntax which conflicts with shell scripts and can be confusing when mixed with Terraform's own interpolation syntax.

**This provider:**

- Uses distinctive `%%VAR%%` syntax
- Explicit validation of all placeholders
- Clear error messages
- Works as a data source with proper state management

### vs. Terraform's deprecated `template` provider

The old `template` provider used `${}` syntax and has been deprecated. This provider:

- Uses modern Plugin Framework
- Uses conflict-free syntax
- Provides better error messages
- Actively maintained

### vs. External template processors

Using external tools (like `envsubst`, `sed`, or custom scripts) requires:

- Additional dependencies
- Complex null_resource or local-exec provisioners
- Less integration with Terraform's plan/apply workflow

**This provider:**

- Pure Terraform solution
- No external dependencies
- Works in Terraform plan phase
- Proper state management

## Development

### Requirements

- Go 1.23+
- Terraform 1.0+

### Building

```bash
go build -o terraform-provider-utils
```

### Testing

Run unit tests:

```bash
go test -v ./internal/provider/
```

Run acceptance tests:

```bash
TF_ACC=1 go test -v ./internal/provider/
```

### Running Examples

```bash
cd examples/basic
terraform init
terraform plan
terraform apply
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

For information about releasing new versions, see [RELEASE.md](RELEASE.md).

## License

This provider is released under the MIT License. See LICENSE for details.

## Support

For issues, questions, or contributions, please visit:
<https://github.com/spantree/terraform-provider-utils>

## Examples

More examples can be found in the `examples/` directory:

- `examples/basic/` - Basic usage examples
- `examples/argo-workflow/` - Complete Argo Workflow deployment example

## Acknowledgments

This provider was built to solve real-world problems when deploying Argo Workflows and other Kubernetes resources with Terraform, where multiple templating systems need to coexist peacefully.
