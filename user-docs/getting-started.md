# Getting Started

Complete tutorial for using the Terraform Spantree Utils Provider.

## Quick Start (5 minutes)

### 1. Install the Provider

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

### 2. Create Your First Template

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
}
```

### 3. Test It

```bash
terraform init
terraform plan
```

You should see:

```text
Changes to Outputs:
  + greeting = "Hello Alice, welcome to Wonderland!"
```

**That's it!** The `@@NAME@@` and `@@PLACE@@` placeholders were replaced with your values.

---

## Tutorial (15 minutes)

### Step 1: Multi-line Templates

Use heredoc syntax for complex templates:

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

output "config" {
  value = data.utils_render_template.config.result
}
```

### Step 2: Templates from Files

Create `config.tpl`:

```text
# Application Configuration
# Generated on @@DATE@@

[server]
host = @@HOST@@
port = @@PORT@@

[database]
connection_string = "postgres://@@DB_USER@@:@@DB_PASS@@@@@DB_HOST@@:@@DB_PORT@@/@@DB_NAME@@"
```

Then reference it:

```hcl
data "utils_render_template" "app_config" {
  template = file("${path.module}/config.tpl")

  values = {
    DATE    = "2025-01-15"
    HOST    = "0.0.0.0"
    PORT    = "3000"
    DB_USER = "appuser"
    DB_PASS = "secret123"
    DB_HOST = "db.example.com"
    DB_PORT = "5432"
    DB_NAME = "myapp"
  }
}
```

### Step 3: Using with Terraform Variables

```hcl
variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "development"
}

variable "app_version" {
  description = "Application version"
  type        = string
  default     = "1.0.0"
}

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = 3
}

data "utils_render_template" "deployment" {
  template = <<-EOT
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: myapp-@@ENV@@
      labels:
        version: "@@VERSION@@"
    spec:
      replicas: @@REPLICAS@@
      template:
        spec:
          containers:
          - name: myapp
            image: myapp:@@VERSION@@
  EOT

  values = {
    ENV      = var.environment
    VERSION  = var.app_version
    REPLICAS = tostring(var.replicas)  # Convert number to string
  }
}
```

**Important**: All values must be strings. Use `tostring()` for numbers and `jsonencode()` for complex types.

Test with different values:

```bash
terraform apply -var="environment=production" -var="app_version=2.1.0" -var="replicas=5"
```

### Step 4: Error Handling

The provider validates that all placeholders have values:

```hcl
data "utils_render_template" "greeting" {
  template = "Hello @@NAME@@, today is @@DAY@@."

  values = {
    NAME = "Bob"
    # Oops! Missing DAY
  }
}
```

You'll get a clear error:

```text
Error: Template Rendering Failed

  with data.utils_render_template.greeting,
  on main.tf line 11, in data "utils_render_template" "greeting":

Failed to render template: missing values for placeholders: DAY
```

This prevents accidentally deploying templates with unreplaced placeholders.

---

## Real-World Example: Argo Workflows

The primary use case - deploying Argo WorkflowTemplates where you need both Terraform values and Argo's runtime parameters.

### The Challenge

When deploying Argo WorkflowTemplates with Terraform:

- **Terraform needs to inject** values like namespace, image tags, and configuration
- **Argo needs to preserve** its `{{inputs.parameters.*}}` and `{{workflow.*}}` syntax for runtime
- Standard Terraform templating conflicts with Argo's `{{}}` syntax

The Spantree Utils provider solves this with its distinctive `@@VAR@@` syntax.

### Complete Example

Create `workflow-template.yaml`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: data-processing-workflow
  namespace: @@NAMESPACE@@
  labels:
    app: data-processor
    managed-by: terraform
spec:
  entrypoint: main

  arguments:
    parameters:
      - name: start-date
        value: "@@START_DATE@@"
      - name: end-date
        value: "@@END_DATE@@"

  templates:
    - name: main
      steps:
        - - name: process-data
            template: processor
            arguments:
              parameters:
                - name: start-date
                  value: "{{workflow.parameters.start-date}}"
                - name: end-date
                  value: "{{workflow.parameters.end-date}}"

    - name: processor
      inputs:
        parameters:
          - name: start-date
          - name: end-date
      container:
        image: @@IMAGE_REPOSITORY@@:@@IMAGE_TAG@@
        command: [sh, -c]
        args:
          - |
            #!/bin/bash
            set -e

            # These are Argo parameters ({{...}})
            START_DATE="{{inputs.parameters.start-date}}"
            END_DATE="{{inputs.parameters.end-date}}"

            echo "Processing from $START_DATE to $END_DATE"

            # Bash conditionals work fine
            if [[ -n "$START_DATE" ]]; then
              echo "Valid date range"
            fi

            # Shell arithmetic works fine
            DAYS=$(($(date -d "$END_DATE" +%s) - $(date -d "$START_DATE" +%s)))
            echo "Processing $DAYS seconds of data"
```

Create `main.tf`:

```hcl
terraform {
  required_providers {
    utils = {
      source  = "spantree/utils"
      version = "~> 0.2"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "utils" {}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
  default     = "argo"
}

variable "image_repository" {
  description = "Container image repository"
  type        = string
  default     = "python"
}

variable "image_tag" {
  description = "Container image tag"
  type        = string
  default     = "3.11-slim"
}

variable "start_date" {
  description = "Default start date"
  type        = string
  default     = "2025-01-01"
}

variable "end_date" {
  description = "Default end date"
  type        = string
  default     = "2025-12-31"
}

# Render the template
data "utils_render_template" "workflow" {
  template = file("${path.module}/workflow-template.yaml")

  values = {
    NAMESPACE        = var.namespace
    IMAGE_REPOSITORY = var.image_repository
    IMAGE_TAG        = var.image_tag
    START_DATE       = var.start_date
    END_DATE         = var.end_date
  }
}

# Deploy to Kubernetes
resource "kubernetes_manifest" "workflow_template" {
  manifest = yamldecode(data.utils_render_template.workflow.result)
}

output "rendered_workflow" {
  description = "The rendered workflow"
  value       = data.utils_render_template.workflow.result
}
```

### Deploy It

```bash
terraform init
terraform plan
terraform apply
```

Verify the deployment:

```bash
kubectl get workflowtemplate -n argo
kubectl get workflowtemplate data-processing-workflow -n argo -o yaml
```

Notice:

- The namespace is "argo" (from Terraform)
- The image is "python:3.11-slim" (from Terraform)
- The `{{inputs.parameters.*}}` syntax is intact (for Argo)

### Understanding the Layers

This example demonstrates three layers of templating:

**Layer 1: Terraform Values (`@@...@@`)** - During `terraform plan/apply`

```yaml
namespace: @@NAMESPACE@@          # → argo
image: @@IMAGE_REPOSITORY@@:@@IMAGE_TAG@@  # → python:3.11-slim
```

**Layer 2: Argo Parameters (`{{...}}`)** - During workflow execution

```yaml
value: "{{workflow.parameters.start-date}}"  # Evaluated by Argo
value: "{{inputs.parameters.end-date}}"      # Evaluated by Argo
```

**Layer 3: Shell Syntax (`$...`, `[[...]]`, etc.)** - During container execution

```bash
START_DATE="{{inputs.parameters.start-date}}"  # Argo fills this
echo "Date: $START_DATE"                        # Shell expands this
if [[ -n "$START_DATE" ]]; then                # Shell evaluates this
```

---

## Best Practices

### ✅ DO

#### 1. Use descriptive placeholder names

```hcl
# Good
@@DATABASE_HOST@@
@@API_KEY@@

# Avoid
@@H@@
@@K@@
```

#### 2. Keep templates in separate files for complex content

```hcl
template = file("${path.module}/templates/config.tpl")
```

#### 3. Use consistent naming conventions

```hcl
# All uppercase with underscores
@@APP_NAME@@
@@DB_HOST@@
@@MAX_CONNECTIONS@@
```

#### 4. Convert non-string values

```hcl
values = {
  PORT     = tostring(var.port)           # number to string
  ENABLED  = var.enabled ? "true" : "false"  # bool to string
  CONFIG   = jsonencode(var.config_map)   # map/list to JSON string
}
```

### ❌ DON'T

#### 1. Don't use hyphens or special characters in placeholder names

```hcl
# Invalid
@@app-name@@
@@db.host@@
```

#### 2. Don't leave unused placeholders in templates

```hcl
# Bad - @@UNUSED@@ will cause an error
template = "Hello @@NAME@@, @@UNUSED@@"
values = {
  NAME = "Alice"
}
```

#### 3. Don't use raw numbers

```hcl
# Bad
values = {
  PORT = 8080  # Error: must be string
}

# Good
values = {
  PORT = "8080"  # or tostring(var.port)
}
```

---

## Troubleshooting

**Problem**: `terraform init` fails with "provider not found"
**Solution**: Check your internet connection and verify the provider name is spelled correctly: `spantree/utils`

**Problem**: Template doesn't render as expected
**Solution**: Ensure placeholder names in the template exactly match the keys in your `values` map (case-sensitive)

**Problem**: Getting "missing values" error
**Solution**: Every `@@PLACEHOLDER@@` in your template must have a corresponding entry in the `values` map

**Problem**: Argo syntax is being replaced
**Solution**: Ensure you're using `@@VAR@@` for Terraform values, not `{{VAR}}`

**Problem**: Shell commands fail
**Solution**: Verify shell syntax like `$VAR`, `[[ ]]`, `$(( ))` is preserved in the rendered output

---

## Next Steps

- **API Details**: See [Reference](reference.md) for complete data source documentation
- **Contributing**: See [Contributing](contributing.md) to help develop the provider
- **More Examples**: Check the [examples/](../examples/) directory in the repository

---

**[← Back to Documentation](README.md)**
