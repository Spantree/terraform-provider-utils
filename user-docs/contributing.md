# Contributing Guide

Guide for contributing to the Terraform Spantree Utils Provider.

## Development Setup

### Prerequisites

- Go 1.24+ installed
- Terraform 1.0+ installed
- Git
- [pre-commit](https://pre-commit.com/) (for automatic code quality checks and documentation generation)

### Installing Pre-commit Hooks (Recommended)

We use pre-commit hooks to automatically generate documentation and enforce code quality standards.

**Install pre-commit:**

```bash
# macOS
brew install pre-commit

# Linux
pip install pre-commit

# Windows (with Python)
pip install pre-commit
```

**Setup hooks (one-time, after cloning):**

```bash
cd terraform-provider-utils
pre-commit install
```

**What the hooks do:**

- ✅ Remove trailing whitespace
- ✅ Fix end-of-file newlines
- ✅ Validate YAML files (excluding template files)
- ✅ Check for accidentally committed large files
- ✅ **Auto-generate Terraform Registry documentation** from provider code

> **Note**: Documentation in `docs/` is automatically generated from your provider code, templates, and examples before each commit. You don't need to manually run documentation commands!

### Quick Start

#### 1. Clone the repository

```bash
git clone https://github.com/spantree/terraform-provider-utils.git
cd terraform-provider-utils
```

#### 2. Build the provider

```bash
make build
```

#### 3. Configure dev override

Create or edit `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "spantree/utils" = "/absolute/path/to/terraform-provider-utils"
  }
  direct {}
}
```

**Important**: Use the absolute path to your cloned repository.

#### 4. Test it

Create a test directory:

```bash
mkdir ~/test-provider
cd ~/test-provider
```

Create `main.tf`:

```hcl
terraform {
  required_providers {
    utils = {
      source = "spantree/utils"
    }
  }
}

provider "utils" {}

data "utils_render_template" "test" {
  template = "Hello @@NAME@@!"
  values = {
    NAME = "Developer"
  }
}

output "result" {
  value = data.utils_render_template.test.result
}
```

Run it:

```bash
terraform init
terraform plan
```

You should see a warning about the dev override and the rendered output.

### Making Changes

1. Edit the provider code in your repository
2. Rebuild: `make build`
3. Changes are immediately available (no need to run `terraform init` again)

### Running Tests

```bash
# Unit tests
make test

# Acceptance tests
make testacc

# Specific test
go test -v ./internal/provider/ -run TestRenderTemplate

# Format code
make fmt

# Lint code
make lint
```

---

## Documentation

This project maintains **two documentation systems**:

### 1. User Documentation (`user-docs/`)

**Manual** - Comprehensive guides for users and contributors.

```text
user-docs/
├── README.md            # Navigation hub (auto-displayed on GitHub)
├── getting-started.md   # Tutorial and examples
├── reference.md         # API documentation
└── contributing.md      # This file
```

**How to update**: Just edit the markdown files directly.

### 2. Terraform Registry Documentation (`docs/`)

**Auto-generated** - API docs for the Terraform Registry.

```text
docs/
├── index.md
└── data-sources/
    └── render_template.md
```

**How it's generated (automatically!)**:

- ✅ **Pre-commit hook** (recommended): Runs before each `git commit`
- ✅ **GitHub Actions**: Runs as backup when code is pushed to `main`

You **don't need to manually generate docs** - they're created automatically from:

- Provider schema (in `internal/provider/*.go`)
- Templates (in `templates/*.md.tmpl`)
- Examples (in `examples/`)

**Manual generation** (if needed):

```bash
# Using make
make docs

# Or directly
go generate ./...
```

> **Important**: Never manually edit files in `docs/` - your changes will be overwritten!

This extracts documentation from:

- Schema descriptions in `internal/provider/*.go`
- Examples in `examples/`
- Templates in `templates/`

**Important**: The `docs/` directory is auto-generated but committed to the repository because:

- Terraform Registry reads docs from Git tags/releases
- CI validates docs are up-to-date

**Workflow**:

1. Update schema descriptions in `internal/provider/*.go`
2. Run `make docs` to regenerate
3. Commit both code and `docs/` changes

If you forget to run `make docs`, CI will fail.

---

## Releasing (for Maintainers)

### One-Time Setup

#### 1. Generate GPG Key

```bash
gpg --full-generate-key

# Choose RSA and RSA, 4096 bits
# Get your key fingerprint
gpg --list-secret-keys --keyid-format=long
```

#### 2. Add GitHub Secrets

Go to `https://github.com/spantree/terraform-provider-utils/settings/secrets/actions`

Add two secrets:

- **GPG_PRIVATE_KEY**: Contents of `gpg --armor --export-secret-keys YOUR_FINGERPRINT`
- **GPG_PASSPHRASE**: The passphrase you set when creating the key

#### 3. Register Provider on Terraform Registry

1. Go to <https://registry.terraform.io>
2. Sign in with GitHub
3. Click **Publish** → **Provider**
4. Select repository: `spantree/terraform-provider-utils`
5. Add your GPG public key in provider settings

### Release Steps

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes

**Steps:**

**1. Update `CHANGELOG.md`:**

```markdown
## [0.3.0] - 2025-01-15

### Added
- New feature X

### Fixed
- Bug Y
```

**2. Commit and push:**

```bash
git add CHANGELOG.md
git commit -m "chore: prepare for v0.3.0 release"
git push origin main
```

**3. Create and push tag:**

```bash
git tag -a v0.3.0 -m "Release v0.3.0"
git push origin v0.3.0
```

**4. Monitor GitHub Actions** at `https://github.com/spantree/terraform-provider-utils/actions`

**5. Go to Releases**, review the draft release, then publish

**6. Verify on Terraform Registry:** `https://registry.terraform.io/providers/spantree/utils` (may take 1-2 minutes)

### Testing Releases Locally

```bash
# Install goreleaser
brew install goreleaser

# Test build (creates binaries in dist/ without releasing)
goreleaser build --snapshot --clean

# Check the output
ls -la dist/
```

### Troubleshooting Releases

#### Workflow failed: "No secret key"

- Verify `GPG_PRIVATE_KEY` secret is set correctly
- Ensure it includes the full BEGIN/END blocks
- Check that `GPG_PASSPHRASE` is correct

#### Release not on Terraform Registry

- Ensure the release is **published** (not draft)
- Verify public key is registered in Terraform Registry
- Wait 2-3 minutes for webhook to process

---

## Project Structure

```text
terraform-provider-utils/
├── internal/provider/      # Provider implementation
│   ├── provider.go         # Provider definition (TypeName = "utils")
│   ├── render_template_data_source.go
│   └── *_test.go
├── examples/               # Working examples
│   ├── basic/
│   └── argo-workflow/
├── docs/                   # Auto-generated Terraform Registry docs
├── user-docs/              # Manual user documentation
├── templates/              # Templates for doc generation
├── main.go                 # Provider entry point
├── go.mod                  # Go dependencies
├── Makefile                # Build tasks
└── .github/workflows/      # CI/CD
```

---

## Development Troubleshooting

### Changes not reflected

Rebuild with `make build` after code changes

### "Provider not found" error

Check the path in `~/.terraformrc` is absolute and correct

### Tests fail

Ensure Go modules are up to date: `go mod tidy`

### "Unexpected difference in directories after code generation"

Run `make docs` and commit the generated files

---

## Pull Request Guidelines

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make test` and `make fmt`
5. Run `make docs` if you changed schemas
6. Push and create a PR
7. Test workflow runs automatically
8. Address any failures
9. Get review and merge

---

**[← Back to Documentation](README.md)**
