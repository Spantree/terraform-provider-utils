terraform {
  required_providers {
    spantree_utils = {
      source = "spantree/utils"
    }
  }
}

provider "spantree_utils" {}

# Basic example: Simple string interpolation
data "spantree_utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@!"

  values = {
    NAME  = "Alice"
    PLACE = "Wonderland"
  }
}

output "greeting" {
  value = data.spantree_utils_render_template.greeting.result
}

# Example: Configuration file with multiple variables
data "spantree_utils_render_template" "config" {
  template = <<-EOT
    server {
      host = "@@HOST@@"
      port = @@PORT@@
      environment = "@@ENV@@"
      
      database {
        connection_string = "@@DB_CONN@@"
      }
    }
  EOT

  values = {
    HOST    = "localhost"
    PORT    = "8080"
    ENV     = "production"
    DB_CONN = "postgresql://user:pass@db:5432/mydb"
  }
}

output "config" {
  value = data.spantree_utils_render_template.config.result
}

# Example: Script with shell variables preserved
data "spantree_utils_render_template" "script" {
  template = <<-EOT
    #!/bin/bash
    
    # Terraform-injected values
    NAMESPACE=@@NAMESPACE@@
    IMAGE_TAG=@@IMAGE_TAG@@
    
    # Shell variables that remain as-is
    echo "Namespace: $${NAMESPACE}"
    echo "Image tag: $(echo $IMAGE_TAG)"
    
    # Bash conditionals work fine
    if [[ -n "$NAMESPACE" ]]; then
      echo "Namespace is set"
    fi
    
    # Arithmetic expressions preserved
    count=$((1 + 2))
    echo "Count: $count"
  EOT

  values = {
    NAMESPACE = "production"
    IMAGE_TAG = "v1.2.3"
  }
}

output "script" {
  value = data.spantree_utils_render_template.script.result
}

