# Placeholder Syntax Reference

Complete reference for the `@@VAR@@` placeholder syntax.

## Pattern

Placeholders must match the regular expression:

```regex
@@[A-Za-z0-9_]+@@
```

## Rules

1. **Must start and end with `@@`**
2. **Name must contain only**:
   - Uppercase letters (A-Z)
   - Lowercase letters (a-z)
   - Numbers (0-9)
   - Underscores (_)
3. **Must contain at least one character** between the `@@` markers
4. **Case-sensitive**: `@@VAR@@` and `@@var@@` are different placeholders

## Valid Examples

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

## Invalid Examples

```text
@@my-var@@          # Hyphens not allowed
@@my.var@@          # Dots not allowed
@@my var@@          # Spaces not allowed
@@my/var@@          # Slashes not allowed
@@my@var@@          # @ only allowed as delimiters
@@@@                # Must have content between @@
@@123@@             # Must contain at least one letter
@@-var@@            # Cannot start with hyphen
@@var-@@            # Cannot end with hyphen
```

## Naming Conventions

### Recommended

Use UPPERCASE_WITH_UNDERSCORES for consistency:

```text
@@APP_NAME@@
@@DATABASE_HOST@@
@@MAX_RETRIES@@
```

### Also Valid

lowercase_with_underscores:

```text
@@app_name@@
@@database_host@@
@@max_retries@@
```

camelCase or PascalCase:

```text
@@appName@@
@@AppName@@
```

### Best Practice

Choose one convention and stick to it throughout your project.

## Usage in Templates

### Single Line

```text
server_host = @@HOST@@
```

### Multiple Placeholders Per Line

```text
connection_string = "postgres://@@USER@@:@@PASS@@@@@HOST@@:@@PORT@@/@@DB@@"
```

**Note**: When placeholders are adjacent, they must each have their own `@@` delimiters.

### Multi-line

```yaml
metadata:
  name: @@APP_NAME@@
  namespace: @@NAMESPACE@@
spec:
  replicas: @@REPLICAS@@
```

### Repeated Placeholders

The same placeholder can appear multiple times:

```text
Welcome to @@APP_NAME@@!
@@APP_NAME@@ is great.
Visit https://@@APP_NAME@@.com
```

All instances will be replaced with the same value.

## Escaping

There is no escaping mechanism. If you need literal `@@` in your output:

### Option 1: Use a Placeholder

```hcl
template = "Email: user@@AT@@domain.com"
values = {
  AT = "@"
}
```

### Option 2: Split the Template

```hcl
template = "Part 1"  # Doesn't contain @@
```

### Option 3: Post-process

Apply additional transformations after rendering if needed.

## Validation

The provider validates:

### 1. All Placeholders Have Values

```hcl
template = "Hello @@NAME@@, @@TITLE@@"
values = {
  NAME = "Alice"
  # Missing TITLE
}
# Error: missing values for placeholders: TITLE
```

### 2. Placeholder Syntax

Invalid placeholder syntax is detected:

```hcl
template = "Hello @@my-name@@"
# Error: Invalid placeholder syntax
```

## Edge Cases

### Empty Values

Empty strings are valid:

```hcl
template = "Prefix: @@VALUE@@"
values = {
  VALUE = ""
}
# Result: "Prefix: "
```

### Whitespace

Whitespace around values is preserved:

```hcl
values = {
  NAME = "  Alice  "
}
# The spaces are kept in the output
```

### Special Characters in Values

Values can contain any characters:

```hcl
values = {
  SPECIAL = "Hello @@ World $$ {{test}}"
}
# The value is inserted as-is
```

### Unicode

Placeholder names must be ASCII, but values can be Unicode:

```hcl
template = "Greeting: @@MSG@@"
values = {
  MSG = "こんにちは"  # Japanese
}
# Works fine
```

## Performance

- Placeholder replacement is O(n) where n is template size
- Large templates (>1MB) may impact plan performance
- Consider splitting very large templates

## See Also

- [Data Source: render_template](data-source-render-template.md)
- [Error Messages](error-messages.md)
- [Tutorial: Basic Template Rendering](../tutorials/basic-template-rendering.md)
