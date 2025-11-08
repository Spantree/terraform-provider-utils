# Terraform Spantree Utils Provider - Documentation

Complete documentation for the Terraform provider that offers conflict-free template rendering.

## Quick Links

- **[Getting Started](getting-started.md)** - Tutorial, examples, and best practices
- **[Reference](reference.md)** - Complete API documentation
- **[Contributing](contributing.md)** - Development setup and release process

## Why This Provider?

When deploying infrastructure with Terraform, you often need to inject Terraform-managed values (like image tags, namespaces, or configuration) into templates that contain their own templating syntax.

### The Problem

Standard Terraform templating uses `${}` syntax, which conflicts with:

- **Argo Workflows**: `{{inputs.parameters.*}}` and `{{workflow.*}}` for runtime parameters
- **Shell Scripts**: `${VAR}` for variable expansion and `$(command)` for command substitution
- **Bash**: `[[`, `]]`, `<<`, `>>`, `&&`, `||`, `((`, `))` for various operations
- **Windows Batch**: `%%VAR%%` for variable expansion
- **Helm Charts**: `{{ .Values.* }}` for templating

### The Solution

This provider uses `@@VAR@@` syntax which doesn't conflict with any common templating system, while staying lightweight and focused on just template rendering.

**Example:**

```hcl
data "utils_render_template" "workflow" {
  template = file("workflow.yaml")

  values = {
    NAMESPACE = "production"
    IMAGE_TAG = "v1.2.3"
  }
}
```

In your `workflow.yaml`:

```yaml
metadata:
  namespace: @@NAMESPACE@@  # Terraform replaces this
spec:
  container:
    image: myapp:@@IMAGE_TAG@@  # Terraform replaces this
    args:
      - echo "{{inputs.parameters.input}}"  # Argo evaluates this at runtime
      - echo "$HOME"  # Shell expands this at runtime
```

After rendering, the `@@` placeholders are replaced with Terraform values, while `{{}}` and `${}` remain intact for their respective systems.

## Key Features

✅ **Conflict-free syntax** - `@@VAR@@` doesn't interfere with other templating systems
✅ **Validation** - All placeholders must have values or an error is returned
✅ **Simple** - Pure data source, no infrastructure created
✅ **Fast** - Template rendering happens locally during plan/apply

---

**[← Back to Project README](../README.md)**
