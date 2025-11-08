# Documentation System

This project maintains **two separate documentation systems** to serve different audiences:

## 1. User Documentation (`user-docs/`)

**Purpose**: Comprehensive learning and reference material for users and contributors

**Structure**: Organized using the [Diataxis framework](https://diataxis.fr/)

```text
user-docs/
├── tutorials/          # Learning-oriented (step-by-step lessons)
├── how-to/            # Problem-oriented (task guides)
├── reference/         # Information-oriented (technical specs)
└── explanation/       # Understanding-oriented (concepts)
```

**Audience**:

- Users learning the provider
- Contributors understanding the codebase
- Developers wanting detailed explanations

**Access**: [Browse user-docs/](user-docs/)

## 2. Terraform Registry Documentation (`docs/`)

**Purpose**: Auto-generated API documentation for the Terraform Registry

**Structure**: Following [Terraform Registry requirements](https://developer.hashicorp.com/terraform/registry/providers/docs)

```text
docs/                  # Auto-generated - DO NOT EDIT
├── index.md
├── data-sources/
│   └── render_template.md
└── (future: resources/, functions/, guides/)
```

**Audience**:

- Users browsing the Terraform Registry
- Terraform CLI (`terraform init` command)

**Access**: Published automatically to [registry.terraform.io/providers/spantree/utils](https://registry.terraform.io/providers/spantree/utils)

## Generating Registry Documentation

The `docs/` directory is **auto-generated** using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs).

### Prerequisites

Install the tool:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

### Generate

```bash
# Using make
make docs

# Or directly
go generate ./...
```

This will:

1. Format examples with `terraform fmt`
2. Extract schema descriptions from the provider code
3. Generate documentation files in `docs/`
4. Use templates from `templates/` for customization

### Customizing Templates

Template files in `templates/` control how the generated documentation looks:

- `templates/index.md.tmpl` - Provider index page
- `templates/data-sources/render_template.md.tmpl` - Data source documentation

The templates use Go template syntax and have access to:

- `{{ .SchemaMarkdown }}` - Auto-generated schema table
- `{{ .Description }}` - Description from schema
- `{{tffile "path/to/example.tf"}}` - Include example files

## Workflow

### Adding Examples

1. Add or update examples in `examples/`
2. Run `make docs` to regenerate
3. Commit both `examples/` and `docs/` changes

### Updating Schema Descriptions

1. Edit descriptions in `internal/provider/*.go` files
2. Run `make docs` to regenerate
3. Commit the code and `docs/` changes

### Updating User Documentation

1. Edit files in `user-docs/`
2. Follow Diataxis principles:
   - **Tutorials**: Teaching through doing
   - **How-to**: Solving specific problems
   - **Reference**: Technical information
   - **Explanation**: Conceptual understanding
3. Update `user-docs/index.md` if adding new pages
4. Commit `user-docs/` changes

## CI/CD Integration

The GitHub Actions workflow verifies generated docs are up to date:

```yaml
generate:
  - run: go generate ./...
  - name: git diff
    run: git diff --exit-code || exit 1
```

If the `docs/` directory is out of sync with the code, the CI fails.

## Why Two Documentation Systems?

### User Documentation (user-docs/)

- **Rich learning paths**: Tutorials guide users step-by-step
- **Task-focused**: How-to guides for specific problems
- **Deep explanations**: Understanding the "why" behind decisions
- **Better organization**: Diataxis framework improves discoverability

### Registry Documentation (docs/)

- **Registry requirement**: Terraform Registry expects specific format
- **Schema-based**: Auto-generated from provider code
- **Always accurate**: Reflects actual provider behavior
- **Version-specific**: Each release has its own docs

## Best Practices

### DO

✅ Keep `user-docs/` manually updated with tutorials and explanations  
✅ Let `docs/` be auto-generated from code  
✅ Add good descriptions to schema in provider code  
✅ Create rich examples in `examples/`  
✅ Run `make docs` before committing  
✅ Commit both systems to the repository

### DON'T

❌ Manually edit files in `docs/` (they'll be overwritten)  
❌ Forget to run `make docs` after schema changes  
❌ Duplicate content between the two systems  
❌ Put Diataxis docs in `docs/` (use `user-docs/`)

## Directory Structure

```text
terraform-provider-utils/
├── user-docs/               # User-facing documentation (Diataxis)
│   ├── index.md
│   ├── tutorials/
│   ├── how-to/
│   ├── reference/
│   └── explanation/
├── docs/                    # Auto-generated (Terraform Registry)
│   ├── index.md
│   └── data-sources/
├── templates/               # Templates for doc generation
│   ├── index.md.tmpl
│   └── data-sources/
├── examples/                # Working examples (used by both)
│   ├── basic/
│   └── argo-workflow/
└── internal/provider/       # Provider code with schema descriptions
```

## References

- [Terraform Registry Provider Documentation](https://developer.hashicorp.com/terraform/registry/providers/docs)
- [terraform-plugin-docs GitHub](https://github.com/hashicorp/terraform-plugin-docs)
- [Diataxis Framework](https://diataxis.fr/)
- [User Documentation](user-docs/)
