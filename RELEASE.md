# Release Process

This document describes the process for releasing new versions of the Terraform Spantree Utils Provider to the Terraform Registry.

## Overview

The release process is automated using GitHub Actions and GoReleaser. When you push a version tag, the workflow:

1. Builds binaries for multiple platforms (Linux, macOS, Windows, FreeBSD)
2. Creates checksums for all binaries
3. Signs the checksums with GPG
4. Creates a GitHub release with all artifacts
5. Notifies the Terraform Registry via webhook

## Prerequisites

Before you can release, you need to complete the following one-time setup:

### 1. Generate a GPG Key

You need a GPG key to sign releases. The Terraform Registry requires signed releases for security.

```bash
# Generate a new GPG key (if you don't have one)
gpg --full-generate-key

# Choose:
# - RSA and RSA
# - 4096 bits
# - Key does not expire (or set expiration as needed)
# - Real name: Your name or "Spantree Terraform Provider"
# - Email: Your email
# - Comment: Optional

# List your keys to get the fingerprint
gpg --list-secret-keys --keyid-format=long

# Example output:
# sec   rsa4096/ABCD1234EFGH5678 2024-01-01 [SC]
#       1234567890ABCDEF1234567890ABCDEF12345678
# uid                 [ultimate] Your Name <your.email@example.com>

# The long string (1234567890ABCDEF...) is your key fingerprint
```

### 2. Export Your GPG Key

```bash
# Export the private key (you'll need this for GitHub Secrets)
gpg --armor --export-secret-keys YOUR_KEY_FINGERPRINT

# This will output something like:
# -----BEGIN PGP PRIVATE KEY BLOCK-----
# ...
# -----END PGP PRIVATE KEY BLOCK-----

# Copy this entire output including the BEGIN and END lines
```

### 3. Add GPG Key to GitHub Secrets

Go to your GitHub repository settings:

1. Navigate to **Settings** → **Secrets and variables** → **Actions**
2. Add the following secrets:
   - `GPG_PRIVATE_KEY`: Paste the entire private key (including BEGIN/END lines)
   - `GPG_PASSPHRASE`: The passphrase you set when creating the key

### 4. Register Your Provider on Terraform Registry

1. Go to https://registry.terraform.io
2. Sign in with your GitHub account
3. Click **Publish** → **Provider**
4. Select your GitHub repository: `spantree/terraform-provider-spantree`
5. Follow the prompts to complete registration

### 5. Add Your Public GPG Key to Terraform Registry

1. Export your public key:
   ```bash
   gpg --armor --export YOUR_KEY_FINGERPRINT
   ```

2. In the Terraform Registry:
   - Go to your provider's settings
   - Navigate to **Signing Keys**
   - Click **Add a key**
   - Paste your public key
   - Save

The registry will now verify all your releases against this public key.

## Release Steps

### 1. Update the Changelog

Before releasing, update `CHANGELOG.md` with all changes in the new version:

```markdown
## [1.1.0] - 2024-01-15

### Added
- New feature X
- New data source Y

### Fixed
- Bug fix Z

### Changed
- Updated dependency A
```

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (1.0.0 → 2.0.0): Breaking changes
- **MINOR** (1.0.0 → 1.1.0): New features, backwards compatible
- **PATCH** (1.0.0 → 1.0.1): Bug fixes, backwards compatible

### 2. Commit and Push Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare for v1.1.0 release"
git push origin main
```

### 3. Create and Push a Version Tag

```bash
# Create an annotated tag
git tag -a v1.1.0 -m "Release v1.1.0"

# Push the tag to GitHub
git push origin v1.1.0
```

**Important**: The tag must:
- Start with `v` (e.g., `v1.1.0`, not `1.1.0`)
- Follow semantic versioning (MAJOR.MINOR.PATCH)
- Be an annotated tag (use `-a` flag)

### 4. Monitor the Release Workflow

1. Go to your GitHub repository
2. Click on **Actions**
3. You should see a "Release" workflow running
4. Wait for it to complete (usually 2-5 minutes)

The workflow will:
- Build binaries for all platforms
- Generate SHA256 checksums
- Sign the checksums with your GPG key
- Create a GitHub release (draft mode by default)
- Upload all artifacts

### 5. Review and Publish the Release

1. Go to **Releases** in your GitHub repository
2. You should see a draft release for your version
3. Review the release notes and artifacts:
   - Binaries for all platforms
   - `terraform-provider-spantree_X.Y.Z_SHA256SUMS`
   - `terraform-provider-spantree_X.Y.Z_SHA256SUMS.sig` (GPG signature)
   - `terraform-provider-spantree_X.Y.Z_manifest.json`
4. Edit the release notes if needed (add highlights, breaking changes, etc.)
5. Click **Publish release**

### 6. Verify Terraform Registry Update

After publishing the release on GitHub:

1. Wait 1-2 minutes for the webhook to trigger
2. Go to https://registry.terraform.io/providers/spantree/utils
3. Verify that your new version appears in the version list
4. Click on the version to see the documentation

If the version doesn't appear:
- Check the webhook deliveries in GitHub Settings → Webhooks
- Verify your GPG signature is valid
- Check that all required files are present in the release

### 7. Test the New Release

Create a test Terraform configuration:

```hcl
terraform {
  required_providers {
    spantree_utils = {
      source  = "spantree/utils"
      version = "1.1.0"  # Your new version
    }
  }
}

provider "spantree_utils" {}

data "spantree_utils_render_template" "test" {
  template = "Hello @@NAME@@!"
  values = {
    NAME = "World"
  }
}

output "result" {
  value = data.spantree_utils_render_template.test.result
}
```

Run:
```bash
terraform init
terraform plan
terraform apply
```

Verify that Terraform downloads and uses your new version.

## Pre-releases

To create a pre-release (beta, alpha, rc):

```bash
# Tag with a pre-release suffix
git tag -a v1.1.0-beta.1 -m "Release v1.1.0-beta.1"
git push origin v1.1.0-beta.1
```

When publishing on GitHub, check the **"This is a pre-release"** box. The Terraform Registry will index it but won't show it as the latest stable version.

## Troubleshooting

### Workflow Failed: "No secret key"

The GPG key wasn't imported correctly. Verify:
- `GPG_PRIVATE_KEY` secret contains the entire private key block
- `GPG_PASSPHRASE` is correct
- The key hasn't expired

### Release Not Appearing on Terraform Registry

Check:
1. **Webhook**: Go to GitHub Settings → Webhooks, verify `registry.terraform.io` webhook exists
2. **Signature**: Verify the GPG signature is valid:
   ```bash
   gpg --verify terraform-provider-spantree_X.Y.Z_SHA256SUMS.sig terraform-provider-spantree_X.Y.Z_SHA256SUMS
   ```
3. **Public Key**: Ensure your public key is registered in Terraform Registry
4. **Release Status**: Ensure the release is published (not draft)

### Wrong Version Number

If you tagged the wrong version:

```bash
# Delete the tag locally
git tag -d v1.1.0

# Delete the tag on GitHub
git push origin :refs/tags/v1.1.0

# Delete the draft release on GitHub (if created)
# Then create the correct tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

**Never** reuse a version number. If v1.1.0 was already published, create v1.1.1 instead.

### Build Failed for Specific Platform

Check `.goreleaser.yml` configuration. You can test locally:

```bash
# Install goreleaser
brew install goreleaser

# Test the build (without releasing)
goreleaser build --snapshot --clean

# Check the dist/ directory for built binaries
ls -la dist/
```

## Best Practices

1. **Always test before releasing**: Run full test suite with `make test` and `make testacc`
2. **Update documentation**: Ensure README and examples reflect new features
3. **Follow semantic versioning**: Be clear about breaking changes
4. **Write good release notes**: Help users understand what changed
5. **Test the release**: Always verify the new version works from the registry
6. **Keep GPG key secure**: Store it safely, back it up, set an expiration date
7. **Monitor the workflow**: Watch the GitHub Actions logs for any issues

## Rollback

If you need to rollback a release:

1. **Don't delete the release** - this breaks checksums for anyone who downloaded it
2. Instead, release a new version with the fix or revert
3. Mark the bad release with a note: "⚠️ This release has issues, please use vX.Y.Z instead"
4. Update documentation to point to the correct version

## Security

- **Never commit GPG keys** to the repository
- **Use GitHub Secrets** for sensitive data
- **Rotate keys periodically** (update both GitHub Secrets and Terraform Registry)
- **Enable 2FA** on your GitHub and Terraform Registry accounts
- **Review release artifacts** before publishing

## Resources

- [Terraform Provider Publishing Guide](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GPG Quick Start](https://www.gnupg.org/gph/en/manual/c14.html)

