# Deploying an Argo Workflow Tutorial

This tutorial demonstrates a real-world use case: deploying an Argo WorkflowTemplate to Kubernetes using Terraform, where you need both Terraform-managed values and Argo's runtime parameters to coexist.

**Time to complete**: 20 minutes  
**Prerequisites**:

- Complete [Getting Started](getting-started.md) tutorial
- Kubernetes cluster with Argo Workflows installed
- `kubectl` configured

## The Challenge

When deploying Argo WorkflowTemplates with Terraform, you face a problem:

- **Terraform needs to inject** values like namespace, image tags, and configuration
- **Argo needs to preserve** its `{{inputs.parameters.*}}` and `{{workflow.*}}` syntax for runtime
- Standard Terraform templating conflicts with Argo's `{{}}` syntax

The Spantree Utils provider solves this with its distinctive `@@VAR@@` syntax.

## What You'll Learn

- How to template Argo WorkflowTemplates
- How to preserve Argo's runtime syntax
- How to deploy the result to Kubernetes
- How shell script syntax coexists with both

## Step 1: Set Up the Project

Create a new directory:

```bash
mkdir argo-workflow-demo
cd argo-workflow-demo
```

## Step 2: Create the Workflow Template

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

**Notice**:

- `@@NAMESPACE@@`, `@@IMAGE_REPOSITORY@@`, `@@IMAGE_TAG@@`, `@@START_DATE@@`, `@@END_DATE@@` - Terraform will replace these
- `{{workflow.parameters.*}}`, `{{inputs.parameters.*}}` - Argo will evaluate these at runtime
- `$START_DATE`, `$END_DATE`, `[[ ]]`, `$(( ))` - Shell syntax is preserved

## Step 3: Create the Terraform Configuration

Create `main.tf`:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source  = "spantree/utils"
      version = "~> 1.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "spantree_utils" {}

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
  default     = "2024-01-01"
}

variable "end_date" {
  description = "Default end date"
  type        = string
  default     = "2024-12-31"
}

# Render the template
data "spantree_utils_render_template" "workflow" {
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
  manifest = yamldecode(data.spantree_utils_render_template.workflow.result)
}

output "rendered_workflow" {
  description = "The rendered workflow"
  value       = data.spantree_utils_render_template.workflow.result
}
```

## Step 4: Preview the Rendered Template

Initialize and plan:

```bash
terraform init
terraform plan
```

Look at the output - you'll see the rendered YAML with:

- `@@NAMESPACE@@` replaced with "argo"
- `@@IMAGE_TAG@@` replaced with "3.11-slim"
- `{{workflow.parameters.start-date}}` **preserved** for Argo

## Step 5: Deploy to Kubernetes

Apply the configuration:

```bash
terraform apply
```

Type `yes` to confirm.

## Step 6: Verify the Deployment

Check that the WorkflowTemplate was created:

```bash
kubectl get workflowtemplate -n argo
```

You should see:

```text
NAME                         AGE
data-processing-workflow     1m
```

View the template:

```bash
kubectl get workflowtemplate data-processing-workflow -n argo -o yaml
```

Notice:

- The namespace is "argo" (from Terraform)
- The image is "python:3.11-slim" (from Terraform)
- The `{{inputs.parameters.*}}` syntax is intact (for Argo)

## Step 7: Run the Workflow

Submit the workflow:

```bash
argo submit --from workflowtemplate/data-processing-workflow -n argo
```

Watch it run:

```bash
argo watch @latest -n argo
```

## Step 8: Update with Different Values

Let's update the image tag. Modify your `terraform.tfvars` or use command line:

```bash
terraform apply -var="image_tag=3.12-slim" -var="start_date=2024-06-01"
```

Terraform will update the WorkflowTemplate with the new values.

## Understanding the Separation of Concerns

This example demonstrates three layers of templating:

### Layer 1: Terraform Values (@@...@@)

**When**: During `terraform plan/apply`  
**What**: Infrastructure configuration

```yaml
namespace: @@NAMESPACE@@          # → argo
image: @@IMAGE_REPOSITORY@@:@@IMAGE_TAG@@  # → python:3.11-slim
```

### Layer 2: Argo Parameters ({{...}})

**When**: During workflow execution  
**What**: Runtime parameters

```yaml
value: "{{workflow.parameters.start-date}}"  # Evaluated by Argo
value: "{{inputs.parameters.end-date}}"      # Evaluated by Argo
```

### Layer 3: Shell Syntax ($..., [[...]], etc.)

**When**: During container execution  
**What**: Shell operations

```bash
START_DATE="{{inputs.parameters.start-date}}"  # Argo fills this
echo "Date: $START_DATE"                        # Shell expands this
if [[ -n "$START_DATE" ]]; then                # Shell evaluates this
```

## Real-World Patterns

### Pattern 1: Environment-Specific Configuration

```hcl
variable "environment" {
  type = string
}

locals {
  config = {
    dev = {
      namespace = "argo-dev"
      image_tag = "latest"
      replicas  = "1"
    }
    prod = {
      namespace = "argo-prod"
      image_tag = "v1.2.3"
      replicas  = "3"
    }
  }
  
  env_config = local.config[var.environment]
}

data "spantree_utils_render_template" "workflow" {
  template = file("${path.module}/workflow-template.yaml")
  
  values = {
    NAMESPACE = local.env_config.namespace
    IMAGE_TAG = local.env_config.image_tag
    REPLICAS  = local.env_config.replicas
  }
}
```

### Pattern 2: Dynamic Workflow Generation

```hcl
variable "workflows" {
  type = map(object({
    image_tag  = string
    start_date = string
    end_date   = string
  }))
}

resource "kubernetes_manifest" "workflows" {
  for_each = var.workflows
  
  manifest = yamldecode(
    data.spantree_utils_render_template.workflow[each.key].result
  )
}

data "spantree_utils_render_template" "workflow" {
  for_each = var.workflows
  
  template = file("${path.module}/workflow-template.yaml")
  
  values = {
    WORKFLOW_NAME = each.key
    IMAGE_TAG     = each.value.image_tag
    START_DATE    = each.value.start_date
    END_DATE      = each.value.end_date
  }
}
```

## What You've Learned

✅ How to template Argo WorkflowTemplates with Terraform  
✅ How `@@VAR@@` syntax coexists with Argo's `{{}}` syntax  
✅ How shell syntax is preserved through both layers  
✅ How to deploy the result to Kubernetes  
✅ Real-world patterns for environment configuration

## Next Steps

- **Explore examples**: Check the [examples/argo-workflow](../../examples/argo-workflow/) directory
- **Learn more**: Read [Explanation: Why This Provider?](../explanation/why-this-provider.md)
- **Advanced usage**: See [How-to: Work with Kubernetes Manifests](../how-to/kubernetes-manifests.md)

## Clean Up

To remove the workflow:

```bash
terraform destroy
```

To clean up the directory:

```bash
cd ..
rm -rf argo-workflow-demo
```

## Troubleshooting

**Problem**: Argo syntax is being replaced  
**Solution**: Ensure you're using `@@VAR@@` for Terraform values, not `{{VAR}}`

**Problem**: Workflow fails at runtime  
**Solution**: Check that Argo's `{{...}}` syntax wasn't accidentally replaced. View the deployed YAML with `kubectl get workflowtemplate -o yaml`

**Problem**: Shell commands fail  
**Solution**: Verify shell syntax like `$VAR`, `[[ ]]`, `$(( ))` is preserved in the rendered output

## Additional Resources

- [Argo Workflows Documentation](https://argoproj.github.io/argo-workflows/)
- [Reference: Placeholder Syntax](../reference/placeholder-syntax.md)
- [Explanation: Design Decisions](../explanation/design-decisions.md)
