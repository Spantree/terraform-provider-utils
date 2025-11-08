# How to Set Up Local Development

This guide shows you how to set up the Spantree Utils provider for local development and testing.

## Prerequisites

- Go 1.24+ installed
- Terraform 1.0+ installed
- Git

## Step 1: Clone the Repository

```bash
git clone https://github.com/spantree/terraform-provider-utils.git
cd terraform-provider-utils
```

## Step 2: Build the Provider

```bash
make build
```

This creates the `terraform-provider-utils` binary in the project root.

## Step 3: Configure Dev Override

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

## Step 4: Test the Setup

Create a test directory:

```bash
mkdir ~/test-provider
cd ~/test-provider
```

Create `main.tf`:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source = "spantree/utils"
    }
  }
}

provider "spantree_utils" {}

data "spantree_utils_render_template" "test" {
  template = "Hello @@NAME@@!"
  values = {
    NAME = "Developer"
  }
}

output "result" {
  value = data.spantree_utils_render_template.test.result
}
```

Run it:

```bash
terraform init
terraform plan
```

You should see a warning about the dev override and the rendered output.

## Step 5: Make Changes

Edit the provider code in your repository, then:

```bash
cd /path/to/terraform-provider-spantree
make build
```

The changes are immediately available in your test configuration (no need to run `terraform init` again).

## Running Tests

### Unit Tests

```bash
make test
```

### Acceptance Tests

```bash
make testacc
```

### Specific Test

```bash
go test -v ./internal/provider/ -run TestRenderTemplate
```

## Code Quality

### Format Code

```bash
make fmt
```

### Lint Code

```bash
make lint
```

## Troubleshooting

**Problem**: Changes not reflected  
**Solution**: Rebuild with `make build` after code changes

**Problem**: "Provider not found" error  
**Solution**: Check the path in `~/.terraformrc` is absolute and correct

**Problem**: Tests fail  
**Solution**: Ensure Go modules are up to date: `go mod tidy`

## Next Steps

- [How to: Contributing](contributing.md)
- [Reference: Project Structure](../reference/project-structure.md)
