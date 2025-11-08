# Data Source: spantree_utils_render_template

Renders a template by replacing placeholders with provided values.

## Example Usage

```hcl
data "spantree_utils_render_template" "config" {
  template = "Server: @@HOST@@:@@PORT@@"
  
  values = {
    HOST = "localhost"
    PORT = "8080"
  }
}

output "result" {
  value = data.spantree_utils_render_template.config.result
}
```

## Argument Reference

### Required Arguments

- **`template`** (String) - The template string containing placeholders in `@@VAR@@` format. Can be inline or loaded from a file using `file()`.

- **`values`** (Map of String) - Map of placeholder names to their replacement values. All placeholders in the template must have a corresponding value.

## Attribute Reference

- **`result`** (String) - The rendered template with all placeholders replaced.

- **`id`** (String) - Internal identifier (hash of template and values).

## Placeholder Syntax

Placeholders must follow the pattern: `@@[A-Za-z0-9_]+@@`

### Valid Examples

```text
@@NAME@@
@@IMAGE_TAG@@
@@my_var_123@@
@@NAMESPACE@@
```

### Invalid Examples

```text
@@my-var@@      # Hyphens not allowed
@@my.var@@      # Dots not allowed
@@my var@@      # Spaces not allowed
@@123@@         # Must contain at least one letter
```

## Validation

The data source validates:

1. **All placeholders have values**: If a placeholder exists in the template but not in `values`, an error is returned.

2. **Placeholder syntax**: Placeholders must match the valid pattern.

## Examples

### From File

```hcl
data "spantree_utils_render_template" "from_file" {
  template = file("${path.module}/template.tpl")
  
  values = {
    VAR1 = "value1"
    VAR2 = "value2"
  }
}
```

### With Variables

```hcl
variable "environment" {
  type = string
}

data "spantree_utils_render_template" "config" {
  template = "ENV=@@ENVIRONMENT@@"
  
  values = {
    ENVIRONMENT = var.environment
  }
}
```

### Multi-line Template

```hcl
data "spantree_utils_render_template" "script" {
  template = <<-EOT
    #!/bin/bash
    export APP_NAME="@@APP_NAME@@"
    export VERSION="@@VERSION@@"
    
    echo "Starting $APP_NAME v$VERSION"
  EOT
  
  values = {
    APP_NAME = "myapp"
    VERSION  = "1.0.0"
  }
}
```

### With Conditional Values

```hcl
variable "enable_debug" {
  type = bool
}

data "spantree_utils_render_template" "config" {
  template = "DEBUG=@@DEBUG@@"
  
  values = {
    DEBUG = var.enable_debug ? "true" : "false"
  }
}
```

### Reused Placeholders

```hcl
data "spantree_utils_render_template" "repeated" {
  template = <<-EOT
    Welcome to @@APP@@!
    @@APP@@ is great.
    Use @@APP@@ today!
  EOT
  
  values = {
    APP = "TerraformUtils"
  }
}

# Result:
# Welcome to TerraformUtils!
# TerraformUtils is great.
# Use TerraformUtils today!
```

## Error Messages

### Missing Placeholder Values

```text
Error: Template Rendering Failed

Failed to render template: missing values for placeholders: VAR1, VAR2
```

**Solution**: Add the missing placeholders to your `values` map.

### Invalid Placeholder Syntax

```text
Error: Template Rendering Failed

Invalid placeholder syntax: @@my-var@@
```

**Solution**: Use only alphanumeric characters and underscores in placeholder names.

## Type Conversion

All values must be strings. Convert other types:

```hcl
values = {
  PORT     = tostring(var.port)           # number to string
  ENABLED  = var.enabled ? "true" : "false"  # bool to string
  CONFIG   = jsonencode(var.config_map)   # map/list to JSON string
}
```

## Performance Considerations

- Template rendering happens during the plan phase
- Large templates (>1MB) may impact plan performance
- Consider splitting very large templates into multiple data sources

## See Also

- [Placeholder Syntax Reference](placeholder-syntax.md)
- [Error Messages Reference](error-messages.md)
- [Tutorial: Basic Template Rendering](../tutorials/basic-template-rendering.md)
