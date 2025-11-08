# Documentation Setup Guide

This project uses **two documentation systems** working together:

## 1. User Documentation (`user-docs/`) ✍️

**Manually maintained** - Organized using [Diataxis framework](https://diataxis.fr/)

- **Tutorials**: Step-by-step learning
- **How-to Guides**: Task-focused solutions
- **Reference**: Technical specifications
- **Explanation**: Conceptual understanding

**Location**: `user-docs/`  
**Edit**: Manually - commit changes  
**Purpose**: Comprehensive learning and reference

## 2. Terraform Registry Documentation (`docs/`) 🤖

**Auto-generated** - For Terraform Registry publication

- **index.md**: Provider overview
- **data-sources/**: Data source API docs
- **resources/**: Resource API docs (future)
- **functions/**: Function docs (future)

**Location**: `docs/` (gitignored)  
**Edit**: Never - auto-generated from code  
**Purpose**: Terraform Registry integration

## Setup Steps

### Step 1: Add Dependencies

```bash
# Add terraform-plugin-docs to go.mod
go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Tidy dependencies
go mod tidy
```

### Step 2: Generate Registry Documentation

```bash
# Generate docs
make docs

# Or manually
go generate ./...
```

This creates the `docs/` directory with:

- `docs/index.md` - Provider overview
- `docs/data-sources/render_template.md` - Data source documentation

### Step 3: Review Generated Docs

```bash
# View generated files
ls -la docs/
cat docs/index.md
cat docs/data-sources/render_template.md
```

### Step 4: Commit Generated Docs

**Important**: The `docs/` directory should be committed to the repository (even though it's auto-generated) because:

1. The Terraform Registry reads docs from Git tags/releases
2. Users browsing GitHub can see the registry-formatted docs
3. The CI/CD validates docs are up to date

```bash
git add docs/
git commit -m "docs: generate Terraform Registry documentation"
```

## Workflow

### When You Change Code

1. Update schema descriptions in `internal/provider/*.go`
2. Run `make docs` to regenerate
3. Commit both code and `docs/` changes

### When You Update Examples

1. Edit files in `examples/`
2. Run `make docs` to regenerate (includes examples in generated docs)
3. Commit both `examples/` and `docs/` changes

### When You Update User Documentation

1. Edit files in `user-docs/`
2. No regeneration needed
3. Commit `user-docs/` changes

## CI/CD Integration

The GitHub Actions `generate` job ensures docs stay in sync:

```yaml
generate:
  - run: go generate ./...
  - run: git diff --exit-code || exit 1  # Fails if docs are out of date
```

If you forget to run `make docs` after code changes, CI will fail.

## Customizing Registry Documentation

Edit templates in `templates/`:

- `templates/index.md.tmpl` - Provider index page template
- `templates/data-sources/render_template.md.tmpl` - Data source page template

After editing templates:

```bash
make docs  # Regenerate with new templates
git add docs/ templates/
git commit -m "docs: update registry documentation templates"
```

## File Structure

```text
terraform-provider-utils/
├── user-docs/               # ✍️ Manual - Diataxis framework
│   ├── index.md             # Documentation hub
│   ├── README.md            # Documentation overview
│   ├── tutorials/           # Learning-oriented
│   ├── how-to/              # Problem-oriented
│   ├── reference/           # Information-oriented
│   └── explanation/         # Understanding-oriented
│
├── docs/                    # 🤖 Auto-generated - Terraform Registry
│   ├── index.md             # Generated from templates/index.md.tmpl
│   └── data-sources/        # Generated from schema + templates
│       └── render_template.md
│
├── templates/               # 🎨 Templates for doc generation
│   ├── index.md.tmpl
│   └── data-sources/
│       └── render_template.md.tmpl
│
├── examples/                # 💻 Examples (used by both systems)
│   ├── basic/
│   └── argo-workflow/
│
└── tools/                   # 🔧 Build tools
    └── tools.go             # Dependency tracking
```

## Troubleshooting

### "Unexpected difference in directories after code generation"

You forgot to run `make docs` after changing code. Run it and commit:

```bash
make docs
git add docs/
git commit -m "docs: regenerate registry documentation"
```

### "terraform-plugin-docs not found"

Install the tool:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

### Templates Not Being Used

Ensure templates are in the correct location:

- `templates/index.md.tmpl`
- `templates/data-sources/<resource_name>.md.tmpl`

### Docs Not Showing on Registry

1. Ensure `docs/` is committed to the repository
2. Create a release/tag
3. The registry pulls docs from the Git tag

## Best Practices

✅ **DO**:

- Keep `user-docs/` updated with tutorials and explanations
- Run `make docs` after every schema change
- Commit the generated `docs/` directory
- Add good descriptions to schema attributes
- Create rich examples in `examples/`

❌ **DON'T**:

- Manually edit files in `docs/` (they're auto-generated)
- Ignore the generate CI check
- Put tutorials in `docs/` (use `user-docs/`)
- Delete the `docs/` directory from git

## References

- [Terraform Registry Provider Docs Format](https://developer.hashicorp.com/terraform/registry/providers/docs)
- [terraform-plugin-docs GitHub](https://github.com/hashicorp/terraform-plugin-docs)
- [Diataxis Framework](https://diataxis.fr/)
- [Main Documentation](DOCUMENTATION.md)
