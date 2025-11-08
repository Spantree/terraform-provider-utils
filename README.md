# Terraform Spantree Utils Provider

A Terraform provider offering utility functions including template rendering with custom placeholder syntax that doesn't conflict with other templating systems.

## Why This Provider?

When working with Terraform to deploy resources like Argo Workflows, Kubernetes manifests, or shell scripts, you often need to inject Terraform-managed values (like namespaces, image tags, or dates) into templates that also contain their own templating syntax.

**The problem:**

- Standard Terraform templating conflicts with Argo `{{}}`, shell `${}`, and bash operators
- Other providers like `kbst/kustomization` or Helm are often overkill for simple template substitution
- Managing Kustomize overlays or Helm charts adds unnecessary complexity when you just need variable injection

**Common conflicts:**

- **Argo Workflows**: `{{inputs.parameters.*}}` syntax
- **Shell scripts**: `${VAR}`, `$(command)` syntax
- **Bash**: `[[ ]]`, `<<`, `>>`, `<>`, `&&`, `||`, `(( ))` operators
- **Windows batch files**: `%%VAR%%` syntax

**This provider solves it** by using a distinctive placeholder syntax (`@@VAR@@`) that doesn't conflict with anything, while staying lightweight and focused on just template rendering.

## Features

- **Conflict-free syntax**: `@@VAR@@` placeholders don't interfere with other templating systems
- **Validation**: All placeholders must have corresponding values or an error is returned
- **Simple**: Pure data source, no infrastructure created
- **Fast**: Template rendering happens locally during plan/apply

## 📚 Documentation

**[📖 View Complete Documentation](user-docs/)** - Comprehensive guides for users and contributors

- **[Getting Started](user-docs/getting-started.md)** - Tutorial, examples, and best practices
- **[Reference](user-docs/reference.md)** - Complete API documentation
- **[Contributing](user-docs/contributing.md)** - Development setup and releases

## Installation

### Terraform 1.0+

Add to your `terraform` block:

```hcl
terraform {
  required_providers {
    utils = {
      source  = "spantree/utils"
      version = "~> 0.2"
    }
  }
}

provider "utils" {}
```

### Local Development

For local development and testing:

1. Build the provider:

   ```bash
   make build
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
data "utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@!"

  values = {
    NAME  = "Alice"
    PLACE = "Wonderland"
  }
}

output "greeting" {
  value = data.utils_render_template.greeting.result
  # Output: "Hello Alice, welcome to Wonderland!"
}
```

### Configuration File Example

```hcl
data "utils_render_template" "config" {
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
data "utils_render_template" "script" {
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

data "utils_render_template" "workflow" {
  template = file("${path.module}/workflow.yaml")

  values = {
    NAMESPACE = var.namespace
    IMAGE_TAG = var.image_tag
  }
}

resource "kubernetes_manifest" "workflow" {
  manifest = yamldecode(data.utils_render_template.workflow.result)
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

## Quick Reference

### Placeholder Syntax

Use `@@VARIABLE_NAME@@` format - alphanumeric and underscores only.

**Learn more**: [Reference Documentation](user-docs/reference.md)

### Key Features

✅ Conflict-free with Argo `{{}}`, shell `${}`, bash operators  
✅ Validates all placeholders have values  
✅ Clear error messages  
✅ No infrastructure created

**Learn more**: [Getting Started](user-docs/getting-started.md)

## Development

See [Contributing Guide](user-docs/contributing.md) for detailed instructions.

**Quick start:**

```bash
# Build
make build

# Test
make test

# Run examples
cd examples/basic && terraform init && terraform plan
```

## Contributing

Contributions welcome! See the [Contributing Guide](user-docs/contributing.md).

## Resources

- 📚 **[Documentation](user-docs/)** - Complete user guides and API reference
- 💻 **[Examples](examples/)** - Working code examples
- 🐛 **[Issues](https://github.com/spantree/terraform-provider-utils/issues)** - Report bugs or request features
- 📦 **[Terraform Registry](https://registry.terraform.io/providers/spantree/utils)** - Official provider listing

## License

MIT License - see [LICENSE](LICENSE) for details.
