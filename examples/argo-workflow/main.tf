terraform {
  required_providers {
    spantree_utils = {
      source = "spantree/utils"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "spantree_utils" {}

provider "kubernetes" {
  # Configure your Kubernetes provider here
  # config_path = "~/.kube/config"
}

# Variables that would typically come from Terraform configuration
variable "namespace" {
  description = "Kubernetes namespace for the workflow"
  type        = string
  default     = "argo"
}

variable "image_repository" {
  description = "Container image repository"
  type        = string
  default     = "myregistry.io/myapp"
}

variable "image_tag" {
  description = "Container image tag"
  type        = string
  default     = "v1.0.0"
}

variable "start_date" {
  description = "Start date for data processing"
  type        = string
  default     = "2024-01-01"
}

variable "end_date" {
  description = "End date for data processing"
  type        = string
  default     = "2024-12-31"
}

# Render the Argo WorkflowTemplate with Terraform values
# Note: Argo's {{inputs.parameters.*}} syntax is preserved
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

# Deploy the rendered workflow to Kubernetes
resource "kubernetes_manifest" "workflow_template" {
  manifest = yamldecode(data.spantree_utils_render_template.workflow.result)
}

output "rendered_workflow" {
  description = "The rendered workflow template"
  value       = data.spantree_utils_render_template.workflow.result
}

