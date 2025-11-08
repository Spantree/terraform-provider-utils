# Reference Documentation

Complete API reference for the Terraform Spantree Utils Provider.

## Data Source: utils_render_template

Renders a template by replacing placeholders with provided values.

### Arguments

#### Required

- **`template`** (String) - The template string containing placeholders in `@@VAR@@` format. Can be inline or loaded from a file using `file()`.

- **`values`** (Map of String) - Map of placeholder names to their replacement values. All placeholders in the template must have a corresponding value.

#### Computed

- **`result`** (String) - The rendered template with all placeholders replaced.

### Basic Example

```hcl
data "utils_render_template" "config" {
  template = "Server: @@HOST@@:@@PORT@@"
  
  values = {
    HOST = "localhost"
    PORT = "8080"
  }
}

output "result" {
  value = data.utils_render_template.config.result
  # Output: "Server: localhost:8080"
}
```

### Examples

#### From File

```hcl
data "utils_render_template" "from_file" {
  template = file("${path.module}/template.tpl")
  
  values = {
    VAR1 = "value1"
    VAR2 = "value2"
  }
}
```

#### With Variables

```hcl
variable "environment" {
  type = string
}

data "utils_render_template" "config" {
  template = "ENV=@@ENVIRONMENT@@"
  
  values = {
    ENVIRONMENT = var.environment
  }
}
```

#### Multi-line Template

```hcl
data "utils_render_template" "script" {
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

#### With Conditional Values

```hcl
variable "enable_debug" {
  type = bool
}

data "utils_render_template" "config" {
  template = "DEBUG=@@DEBUG@@"
  
  values = {
    DEBUG = var.enable_debug ? "true" : "false"
  }
}
```

#### Reused Placeholders

The same placeholder can appear multiple times:

```hcl
data "utils_render_template" "repeated" {
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

### Type Conversion

All values must be strings. Convert other types:

```hcl
values = {
  PORT     = tostring(var.port)           # number to string
  ENABLED  = var.enabled ? "true" : "false"  # bool to string
  CONFIG   = jsonencode(var.config_map)   # map/list to JSON string
}
```

---

## Placeholder Syntax

Complete reference for the `@@VAR@@` placeholder syntax.

### Pattern

Placeholders must match the regular expression:

```regex
@@[A-Za-z0-9_]+@@
```

### Rules

1. **Must start and end with `@@`**
2. **Name must contain only**:
   - Uppercase letters (A-Z)
   - Lowercase letters (a-z)
   - Numbers (0-9)
   - Underscores (_)
3. **Must contain at least one character** between the `@@` markers
4. **Case-sensitive**: `@@VAR@@` and `@@var@@` are different placeholders

### Valid Examples

```text
@@NAME@@
@@IMAGE_TAG@@
@@my_var_123@@
@@NAMESPACE@@
@@DB_HOST@@
@@MAX_CONNECTIONS@@
@@env@@
@@x1@@
```

### Invalid Examples

```text
@@my-var@@          # Hyphens not allowed
@@my.var@@          # Dots not allowed
@@my var@@          # Spaces not allowed
@@my/var@@          # Slashes not allowed
@@my@var@@          # @ only allowed as delimiters
@@@@                # Must have content between @@
@@123@@             # Must contain at least one letter or underscore
@@-var@@            # Cannot start with hyphen
@@var-@@            # Cannot end with hyphen
```

### Naming Conventions

#### Recommended (UPPERCASE_WITH_UNDERSCORES)

```text
@@APP_NAME@@
@@DATABASE_HOST@@
@@MAX_RETRIES@@
```

#### Also Valid (lowercase_with_underscores)

```text
@@app_name@@
@@database_host@@
@@max_retries@@
```

#### Also Valid (camelCase or PascalCase)

```text
@@appName@@
@@AppName@@
```

**Best Practice**: Choose one convention and stick to it throughout your project.

### Usage Patterns

#### Single Line

```text
server_host = @@HOST@@
```

#### Multiple Placeholders Per Line

```text
connection_string = "postgres://@@USER@@:@@PASS@@@@@HOST@@:@@PORT@@/@@DB@@"
```

#### Multi-line

```yaml
metadata:
  name: @@APP_NAME@@
  namespace: @@NAMESPACE@@
spec:
  replicas: @@REPLICAS@@
```

#### Repeated Placeholders

The same placeholder can appear multiple times - all instances will be replaced with the same value:

```text
Welcome to @@APP_NAME@@!
@@APP_NAME@@ is great.
Visit https://@@APP_NAME@@.com
```

### Validation

The provider validates:

#### 1. All Placeholders Have Values

```hcl
template = "Hello @@NAME@@, @@TITLE@@"
values = {
  NAME = "Alice"
  # Missing TITLE
}
# Error: missing values for placeholders: TITLE
```

#### 2. Placeholder Syntax

Invalid placeholder syntax is detected:

```hcl
template = "Hello @@my-name@@"
# Error: Invalid placeholder syntax
```

### Edge Cases

#### Empty Values

Empty strings are valid:

```hcl
template = "Prefix: @@VALUE@@"
values = {
  VALUE = ""
}
# Result: "Prefix: "
```

#### Whitespace

Whitespace around values is preserved:

```hcl
values = {
  NAME = "  Alice  "
}
# The spaces are kept in the output
```

#### Special Characters in Values

Values can contain any characters:

```hcl
values = {
  SPECIAL = "Hello @@ World $$ {{test}}"
}
# The value is inserted as-is
```

#### Unicode

Placeholder names must be ASCII, but values can be Unicode:

```hcl
template = "Greeting: @@MSG@@"
values = {
  MSG = "こんにちは"  # Japanese
}
# Works fine
```

---

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

---

## Performance

- Template rendering happens during the plan phase
- Rendering is O(n) where n is template size
- Large templates (>1MB) may impact plan performance
- Consider splitting very large templates into multiple data sources

---

**[← Back to Documentation](README.md)**
