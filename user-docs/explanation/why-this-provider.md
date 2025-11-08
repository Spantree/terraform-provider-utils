# Why This Provider?

This document explains the problem this provider solves and why it exists.

## The Problem

When deploying infrastructure with Terraform, you often need to inject Terraform-managed values (like image tags, namespaces, or configuration) into templates that contain their own templating syntax.

### Common Conflicts

Standard Terraform templating uses `${}` syntax, which conflicts with:

- **Argo Workflows**: Uses `{{inputs.parameters.*}}` and `{{workflow.*}}` for runtime parameters
- **Shell Scripts**: Uses `${VAR}` for variable expansion and `$(command)` for command substitution
- **Bash**: Uses `[[`, `]]`, `<<`, `>>`, `&&`, `||`, `((`, `))` for various operations
- **Windows Batch**: Uses `%%VAR%%` for variable expansion
- **Jinja2/Ansible**: Uses `{{ var }}` for templating
- **Helm Charts**: Uses `{{ .Values.* }}` for templating

### Real-World Example

Consider deploying an Argo WorkflowTemplate:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: my-workflow
  namespace: ???  # Needs to come from Terraform
spec:
  templates:
    - name: process
      inputs:
        parameters:
          - name: input
      container:
        image: myapp:???  # Needs to come from Terraform
        args:
          - echo "Processing {{inputs.parameters.input}}"  # Must stay for Argo
```

You need:

1. Terraform to inject the namespace and image tag
2. Argo's `{{...}}` syntax to remain untouched for runtime evaluation

## Why Not Use Existing Solutions?

### Terraform's `templatefile()` Function

**Problem**: Uses `${}` syntax which conflicts with shell scripts and is confusing when mixed with Terraform's own interpolation.

```hcl
# Terraform's templatefile()
template = templatefile("script.sh", {
  var = "value"
})

# In script.sh:
echo "${var}"  # Terraform replaces this
echo "$OTHER"  # Shell tries to expand this, but Terraform might too
```

### Terraform's Deprecated `template` Provider

**Problems**:

- Deprecated and no longer maintained
- Used `${}` syntax with same conflicts
- Complex escaping required

### External Tools (envsubst, sed, etc.)

**Problems**:

- Requires external dependencies
- Needs `null_resource` or `local-exec` provisioners
- Doesn't integrate with Terraform's plan/apply workflow
- No validation of placeholders
- Complex error handling

### Manual String Replacement

**Problems**:

- Error-prone
- No validation
- Hard to maintain
- Doesn't scale

## The Solution: Distinctive Syntax

This provider uses `@@VAR@@` syntax which:

1. **Doesn't conflict** with any common templating system
2. **Is visually distinctive** - easy to spot in templates
3. **Is simple to validate** - regex pattern matching
4. **Provides clear errors** - missing placeholders are caught immediately

### How It Works

```hcl
data "spantree_render_template" "workflow" {
  template = file("workflow.yaml")
  
  values = {
    NAMESPACE = "production"
    IMAGE_TAG = "v1.2.3"
  }
}
```

In `workflow.yaml`:

```yaml
metadata:
  namespace: @@NAMESPACE@@  # Terraform replaces this
spec:
  container:
    image: myapp:@@IMAGE_TAG@@  # Terraform replaces this
    args:
      - echo "{{inputs.parameters.input}}"  # Argo evaluates this
      - echo "$HOME"  # Shell expands this
      - if [[ -n "$VAR" ]]; then  # Bash evaluates this
```

After rendering:

```yaml
metadata:
  namespace: production  # ✓ Replaced by Terraform
spec:
  container:
    image: myapp:v1.2.3  # ✓ Replaced by Terraform
    args:
      - echo "{{inputs.parameters.input}}"  # ✓ Preserved for Argo
      - echo "$HOME"  # ✓ Preserved for shell
      - if [[ -n "$VAR" ]]; then  # ✓ Preserved for bash
```

## Key Benefits

### 1. Conflict-Free

The `@@VAR@@` syntax doesn't appear in:

- Programming languages
- Shell scripts
- Other templating systems
- Configuration files
- YAML/JSON/TOML

### 2. Validation

All placeholders must have values:

```hcl
template = "Hello @@NAME@@, @@MISSING@@"
values = {
  NAME = "Alice"
}
# Error: missing values for placeholders: MISSING
```

This prevents accidentally deploying templates with unreplaced placeholders.

### 3. Simple and Focused

- One data source
- Clear purpose
- No complex configuration
- Works in plan phase
- No infrastructure created

### 4. Terraform Native

- Pure Terraform solution
- No external dependencies
- Integrates with plan/apply workflow
- Proper state management
- Works with all Terraform features

## Use Cases

### Argo Workflows

Deploy WorkflowTemplates where Terraform manages infrastructure config and Argo manages runtime parameters.

### Kubernetes Manifests

Inject environment-specific values while preserving Helm-like templating or init container scripts.

### Shell Scripts

Generate scripts with Terraform values while preserving shell syntax for runtime execution.

### Configuration Files

Create config files with Terraform-managed values that also contain application-specific templating.

### Multi-Stage Templating

Handle scenarios where multiple systems need to template the same file at different stages.

## Design Philosophy

1. **Do one thing well**: Template rendering with conflict-free syntax
2. **Fail fast**: Validate all placeholders have values
3. **Be explicit**: Clear syntax, clear errors
4. **Stay simple**: No complex features, just reliable templating
5. **Integrate naturally**: Work with Terraform's existing patterns

## Comparison Summary

| Solution | Syntax | Conflicts | Validation | Terraform Native |
|----------|--------|-----------|------------|------------------|
| This Provider | `@@VAR@@` | None | Yes | Yes |
| templatefile() | `${}` | Many | No | Yes |
| template provider | `${}` | Many | Limited | Deprecated |
| envsubst | `${}` | Many | No | No |
| sed/awk | Custom | Varies | No | No |

## When to Use This Provider

✅ **Use when**:

- Deploying Argo Workflows with Terraform
- Managing Kubernetes manifests with embedded scripts
- Generating shell scripts with Terraform values
- Working with multiple templating systems
- Need validation of placeholder replacement

❌ **Don't use when**:

- Simple string interpolation (use Terraform's native interpolation)
- No templating conflicts exist
- Only using Terraform's own syntax

## See Also

- [Design Decisions](design-decisions.md)
- [Placeholder Syntax Choice](placeholder-syntax.md)
- [Comparison with Alternatives](comparison.md)
- [Tutorial: Argo Workflow Deployment](../tutorials/argo-workflow-deployment.md)
