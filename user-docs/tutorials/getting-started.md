# Getting Started Tutorial

Welcome! This tutorial will guide you through your first steps with the Terraform Spantree Utils Provider. By the end, you'll understand how to install the provider and render your first template.

**Time to complete**: 10 minutes  
**Prerequisites**: Terraform 1.0+ installed

## What You'll Learn

- How to install the provider
- How to configure it in your Terraform code
- How to render a simple template
- How to use the rendered output

## Step 1: Create a New Directory

First, let's create a clean workspace for this tutorial:

```bash
mkdir terraform-utils-tutorial
cd terraform-utils-tutorial
```

## Step 2: Create Your First Configuration

Create a file named `main.tf` with the following content:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source  = "spantree/utils"
      version = "~> 1.0"
    }
  }
}

provider "spantree_utils" {}
```

This tells Terraform to use the Spantree Utils provider.

## Step 3: Add a Simple Template

Now let's add a template data source. Add this to your `main.tf`:

```hcl
data "spantree_utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@!"
  
  values = {
    NAME  = "Alice"
    PLACE = "Wonderland"
  }
}

output "greeting" {
  value = data.spantree_utils_render_template.greeting.result
}
```

**What's happening here?**

- We're using the `spantree_utils_render_template` data source
- The `template` contains placeholders: `@@NAME@@` and `@@PLACE@@`
- The `values` map provides the replacement values
- The `output` will show us the rendered result

## Step 4: Initialize Terraform

Run the initialization command:

```bash
terraform init
```

You should see output like:

```text
Initializing the backend...
Initializing provider plugins...
- Finding spantree/utils versions matching "~> 1.0"...
- Installing spantree/utils v1.0.0...
- Installed spantree/utils v1.0.0

Terraform has been successfully initialized!
```

## Step 5: Preview the Changes

Run a plan to see what Terraform will do:

```bash
terraform plan
```

You should see:

```text
Changes to Outputs:
  + greeting = "Hello Alice, welcome to Wonderland!"
```

Notice how `@@NAME@@` was replaced with "Alice" and `@@PLACE@@` with "Wonderland"!

## Step 6: Apply the Configuration

Apply the configuration to see the final output:

```bash
terraform apply
```

Type `yes` when prompted. You'll see:

```text
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

greeting = "Hello Alice, welcome to Wonderland!"
```

**Important**: Notice that no resources were created. The template rendering happens during the plan phase and doesn't create any infrastructure.

## Step 7: Experiment with Different Values

Let's try changing the values. Update your `main.tf`:

```hcl
data "spantree_utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@!"
  
  values = {
    NAME  = "Bob"
    PLACE = "Terraform Land"
  }
}
```

Run `terraform plan` again:

```bash
terraform plan
```

You'll see the output has changed:

```text
Changes to Outputs:
  ~ greeting = "Hello Alice, welcome to Wonderland!" -> "Hello Bob, welcome to Terraform Land!"
```

## Step 8: Try a Missing Placeholder

Let's see what happens when we forget to provide a value. Update your template:

```hcl
data "spantree_utils_render_template" "greeting" {
  template = "Hello @@NAME@@, welcome to @@PLACE@@! Today is @@DAY@@."
  
  values = {
    NAME  = "Bob"
    PLACE = "Terraform Land"
    # Oops! We forgot to add DAY
  }
}
```

Run `terraform plan`:

```bash
terraform plan
```

You'll get a clear error:

```text
Error: Template Rendering Failed

  with data.spantree_utils_render_template.greeting,
  on main.tf line 11, in data "spantree_utils_render_template" "greeting":
  11: data "spantree_utils_render_template" "greeting" {

Failed to render template: missing values for placeholders: DAY
```

This validation ensures you never accidentally deploy a template with unreplaced placeholders!

## What You've Learned

✅ How to configure the Spantree Utils provider  
✅ How to use the `@@PLACEHOLDER@@` syntax  
✅ How to provide values for placeholders  
✅ How to access the rendered result  
✅ How the provider validates all placeholders have values

## Next Steps

Now that you understand the basics, you can:

1. **Learn more about templates**: Try the [Basic Template Rendering](basic-template-rendering.md) tutorial
2. **See a real-world example**: Check out the [Argo Workflow Deployment](argo-workflow-deployment.md) tutorial
3. **Explore how-to guides**: Learn specific tasks in the [How-to Guides](../how-to/) section

## Clean Up

To remove the tutorial files:

```bash
cd ..
rm -rf terraform-utils-tutorial
```

## Troubleshooting

**Problem**: `terraform init` fails with "provider not found"  
**Solution**: Check your internet connection and verify the provider name is spelled correctly: `spantree/utils`

**Problem**: Template doesn't render as expected  
**Solution**: Ensure placeholder names in the template exactly match the keys in your `values` map (case-sensitive)

**Problem**: Getting "missing values" error  
**Solution**: Every `@@PLACEHOLDER@@` in your template must have a corresponding entry in the `values` map

## Get Help

- [Reference: render_template data source](../reference/data-source-render-template.md)
- [Explanation: Why this provider?](../explanation/why-this-provider.md)
- [GitHub Issues](https://github.com/spantree/terraform-provider-spantree/issues)
