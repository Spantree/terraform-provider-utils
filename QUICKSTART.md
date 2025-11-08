# Quick Start Guide

## Installation for Local Development

1. **Build the provider:**

   ```bash
   make build
   ```

2. **Create a dev override** in `~/.terraformrc`:

   ```hcl
   provider_installation {
     dev_overrides {
       "spantree/utils" = "/Users/feniix/src/spantree/spiderrock/sr-quant-iac/provider/terraform-provider-utils"
     }
     direct {}
   }
   ```

3. **Test with the basic example:**

   ```bash
   cd examples/basic
   terraform init
   terraform plan
   ```

## Quick Test

Create a `test.tf` file:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source = "spantree/utils"
    }
  }
}

provider "spantree_utils" {}

data "spantree_utils_render_template" "hello" {
  template = "Hello @@NAME@@!"
  values = {
    NAME = "World"
  }
}

output "result" {
  value = data.spantree_utils_render_template.hello.result
}
```

Run:

```bash
terraform init
terraform plan
```

Expected output:

```text
Changes to Outputs:
  + result = "Hello World!"
```

## Running Tests

```bash
# Unit tests
make test

# Build
make build

# Format code
make fmt
```

## Key Features

✅ **Conflict-free syntax**: Uses `@@VAR@@` instead of `${}`, `{{}}`, or `%%`  
✅ **Validation**: Errors if placeholders are missing values  
✅ **Preserves other syntax**: Argo `{{}}`, shell `${}`, bash `[[]]`, Windows `%%` all work  
✅ **Simple**: Just a data source, no infrastructure created  

## Next Steps

- See `examples/basic/` for simple examples
- See `examples/argo-workflow/` for a complete Argo Workflow example
- Read `README.md` for full documentation
