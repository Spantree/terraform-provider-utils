# Basic Template Rendering Tutorial

In this tutorial, you'll learn the fundamentals of template rendering with the Spantree Utils provider. We'll explore different use cases and patterns.

**Time to complete**: 15 minutes  
**Prerequisites**: Complete the [Getting Started](getting-started.md) tutorial

## What You'll Learn

- How to render multi-line templates
- How to use templates from files
- How to reuse the same placeholder multiple times
- How to work with different data types
- Best practices for template organization

## Step 1: Multi-line Templates

Create a new directory and `main.tf`:

```bash
mkdir template-tutorial
cd template-tutorial
```

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

data "spantree_utils_render_template" "config" {
  template = <<-EOT
    server {
      host = "@@HOST@@"
      port = @@PORT@@
      environment = "@@ENV@@"
      debug = @@DEBUG@@
    }
  EOT
  
  values = {
    HOST  = "localhost"
    PORT  = "8080"
    ENV   = "production"
    DEBUG = "false"
  }
}

output "config" {
  value = data.spantree_utils_render_template.config.result
}
```

Run it:

```bash
terraform init
terraform apply
```

You'll see:

```text
config = <<EOT
server {
  host = "localhost"
  port = 8080
  environment = "production"
  debug = false
}
EOT
```

**Key points**:

- Use `<<-EOT ... EOT` for multi-line strings
- Indentation is preserved
- All placeholders are replaced, regardless of line

## Step 2: Templates from Files

Create a file named `config.tpl`:

```text
# Application Configuration
# Generated on @@DATE@@

[server]
host = @@HOST@@
port = @@PORT@@

[database]
connection_string = "postgres://@@DB_USER@@:@@DB_PASS@@@@@DB_HOST@@:@@DB_PORT@@/@@DB_NAME@@"

[features]
enable_cache = @@CACHE_ENABLED@@
cache_ttl = @@CACHE_TTL@@
```

Now create `main.tf`:

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

data "spantree_utils_render_template" "app_config" {
  template = file("${path.module}/config.tpl")
  
  values = {
    DATE          = "2024-01-15"
    HOST          = "0.0.0.0"
    PORT          = "3000"
    DB_USER       = "appuser"
    DB_PASS       = "secret123"
    DB_HOST       = "db.example.com"
    DB_PORT       = "5432"
    DB_NAME       = "myapp"
    CACHE_ENABLED = "true"
    CACHE_TTL     = "3600"
  }
}

output "config_file" {
  value = data.spantree_utils_render_template.app_config.result
}
```

Run it:

```bash
terraform init
terraform apply
```

**Best practice**: Keep complex templates in separate files for better organization and syntax highlighting.

## Step 3: Reusing Placeholders

Placeholders can appear multiple times in a template:

```hcl
data "spantree_utils_render_template" "repeated" {
  template = <<-EOT
    Welcome to @@APP_NAME@@!
    
    @@APP_NAME@@ is a powerful tool for managing your infrastructure.
    
    To get started with @@APP_NAME@@, run:
      @@APP_NAME@@ init
      @@APP_NAME@@ plan
      @@APP_NAME@@ apply
    
    Visit https://docs.example.com/@@APP_NAME@@ for more information.
  EOT
  
  values = {
    APP_NAME = "TerraformUtils"
  }
}

output "docs" {
  value = data.spantree_utils_render_template.repeated.result
}
```

All instances of `@@APP_NAME@@` will be replaced with "TerraformUtils".

## Step 4: Working with Variables

You can use Terraform variables to make templates dynamic:

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

data "spantree_utils_render_template" "deployment" {
  template = <<-EOT
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: myapp-@@ENV@@
      labels:
        app: myapp
        version: "@@VERSION@@"
        environment: @@ENV@@
    spec:
      replicas: @@REPLICAS@@
      selector:
        matchLabels:
          app: myapp
      template:
        metadata:
          labels:
            app: myapp
            version: "@@VERSION@@"
        spec:
          containers:
          - name: myapp
            image: myapp:@@VERSION@@
            env:
            - name: ENVIRONMENT
              value: "@@ENV@@"
  EOT
  
  values = {
    ENV      = var.environment
    VERSION  = var.app_version
    REPLICAS = tostring(var.replicas)
  }
}

output "deployment_yaml" {
  value = data.spantree_utils_render_template.deployment.result
}
```

**Important**: All values must be strings. Use `tostring()` for numbers and `jsonencode()` for complex types.

Test with different variables:

```bash
terraform apply -var="environment=production" -var="app_version=2.1.0" -var="replicas=5"
```

## Step 5: Conditional Values

Use Terraform's conditional expressions:

```hcl
variable "enable_debug" {
  type    = bool
  default = false
}

variable "environment" {
  type    = string
  default = "development"
}

data "spantree_utils_render_template" "app_config" {
  template = <<-EOT
    DEBUG_MODE=@@DEBUG@@
    LOG_LEVEL=@@LOG_LEVEL@@
    ENVIRONMENT=@@ENV@@
  EOT
  
  values = {
    DEBUG     = var.enable_debug ? "true" : "false"
    LOG_LEVEL = var.enable_debug ? "debug" : "info"
    ENV       = var.environment
  }
}
```

## Step 6: Using with Other Resources

The rendered template can be used with other Terraform resources:

```hcl
data "spantree_utils_render_template" "script" {
  template = file("${path.module}/setup.sh.tpl")
  
  values = {
    APP_NAME    = "myapp"
    INSTALL_DIR = "/opt/myapp"
    VERSION     = "1.0.0"
  }
}

resource "local_file" "setup_script" {
  content  = data.spantree_utils_render_template.script.result
  filename = "${path.module}/setup.sh"
  
  file_permission = "0755"
}
```

Or with Kubernetes:

```hcl
data "spantree_utils_render_template" "k8s_manifest" {
  template = file("${path.module}/deployment.yaml.tpl")
  
  values = {
    NAMESPACE = "production"
    IMAGE_TAG = "v1.2.3"
  }
}

resource "kubernetes_manifest" "deployment" {
  manifest = yamldecode(data.spantree_utils_render_template.k8s_manifest.result)
}
```

## Best Practices

### ✅ DO

1. **Use descriptive placeholder names**:

   ```hcl
   # Good
   @@DATABASE_HOST@@
   @@API_KEY@@
   
   # Avoid
   @@H@@
   @@K@@
   ```

2. **Keep templates in separate files** for complex content:

   ```hcl
   template = file("${path.module}/templates/config.tpl")
   ```

3. **Use consistent naming conventions**:

   ```hcl
   # All uppercase with underscores
   @@APP_NAME@@
   @@DB_HOST@@
   @@MAX_CONNECTIONS@@
   ```

4. **Document your placeholders**:

   ```hcl
   # config.tpl
   # Placeholders:
   # - @@APP_NAME@@: Application name
   # - @@VERSION@@: Semantic version (e.g., 1.2.3)
   # - @@ENV@@: Environment (dev/staging/prod)
   ```

### ❌ DON'T

1. **Don't use hyphens or special characters** in placeholder names:

   ```hcl
   # Invalid
   @@app-name@@
   @@db.host@@
   ```

2. **Don't leave unused placeholders** in templates:

   ```hcl
   # Bad - @@UNUSED@@ will cause an error
   template = "Hello @@NAME@@, @@UNUSED@@"
   values = {
     NAME = "Alice"
   }
   ```

3. **Don't forget to convert non-string values**:

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

## What You've Learned

✅ How to use multi-line templates with heredoc syntax  
✅ How to load templates from external files  
✅ How to reuse placeholders multiple times  
✅ How to integrate with Terraform variables  
✅ How to use conditional values  
✅ How to pass rendered templates to other resources  
✅ Best practices for naming and organization

## Next Steps

- **Real-world example**: Try the [Argo Workflow Deployment](argo-workflow-deployment.md) tutorial
- **Advanced usage**: See [How-to: Handle Multiple Placeholders](../how-to/multiple-placeholders.md)
- **Syntax details**: Read [Reference: Placeholder Syntax](../reference/placeholder-syntax.md)

## Clean Up

```bash
cd ..
rm -rf template-tutorial
```
