# How to Set Up Releases and GitHub Actions

This guide walks you through the one-time setup for automated releases and CI/CD. For detailed information, see [Release Process](../explanation/release-process.md).

## One-Time Setup Checklist

### ✅ Step 1: Generate GPG Key

```bash
# Generate key
gpg --full-generate-key

# Choose RSA and RSA, 4096 bits
# Set name and email
# Set a strong passphrase

# Get your key fingerprint
gpg --list-secret-keys --keyid-format=long
```

### ✅ Step 2: Export GPG Keys

```bash
# Export private key (for GitHub Secrets)
gpg --armor --export-secret-keys YOUR_FINGERPRINT > private-key.asc

# Export public key (for Terraform Registry)
gpg --armor --export YOUR_FINGERPRINT > public-key.asc
```

### ✅ Step 3: Add GitHub Secrets

Go to: `https://github.com/spantree/terraform-provider-spantree/settings/secrets/actions`

Add two secrets:

- **GPG_PRIVATE_KEY**: Contents of `private-key.asc` (entire file including BEGIN/END)
- **GPG_PASSPHRASE**: The passphrase you set when creating the key

### ✅ Step 4: Register Provider on Terraform Registry

1. Go to <https://registry.terraform.io>
2. Sign in with GitHub
3. Click **Publish** → **Provider**
4. Select repository: `spantree/terraform-provider-spantree`
5. Complete the registration

### ✅ Step 5: Add Public Key to Terraform Registry

1. In Terraform Registry, go to your provider settings
2. Navigate to **Signing Keys**
3. Click **Add a key**
4. Paste contents of `public-key.asc`
5. Save

### ✅ Step 6: Verify GitHub Actions are Enabled

1. Go to `https://github.com/spantree/terraform-provider-utils/actions`
2. Ensure Actions are enabled
3. You should see two workflows:
   - **Tests** (runs on PRs and pushes)
   - **Release** (runs on version tags)

For workflow details, see [GitHub Actions Workflows Reference](../reference/github-actions-workflows.md).

## Testing the Setup

### Test Without Publishing

```bash
# Install goreleaser locally
brew install goreleaser  # macOS
# or
go install github.com/goreleaser/goreleaser@latest

# Test build locally (won't create a release)
goreleaser build --snapshot --clean

# Check the dist/ directory
ls -la dist/
```

### Test the Full Release Process

1. Update `CHANGELOG.md` with test changes
2. Commit and push to main
3. Create a test tag:

   ```bash
   git tag -a v0.0.1-test -m "Test release"
   git push origin v0.0.1-test
   ```

4. Go to GitHub Actions and watch the workflow
5. Check the Releases page for a draft release
6. If successful, you can delete the test release and tag

## Quick Release Steps

Once setup is complete, releasing is simple:

```bash
# 1. Update CHANGELOG.md
vim CHANGELOG.md

# 2. Commit changes
git add CHANGELOG.md
git commit -m "chore: prepare for v1.0.0 release"
git push origin main

# 3. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. Wait for GitHub Actions to complete
# 5. Review and publish the draft release on GitHub
# 6. Verify on Terraform Registry
```

## Troubleshooting

### "No secret key" error in GitHub Actions

- Verify `GPG_PRIVATE_KEY` secret is set correctly
- Ensure it includes the full BEGIN/END blocks
- Check that `GPG_PASSPHRASE` is correct

### Release not appearing on Terraform Registry

- Ensure the release is **published** (not draft)
- Check webhook in GitHub Settings → Webhooks
- Verify public key is registered in Terraform Registry
- Wait 2-3 minutes for webhook to process

### Build fails for a platform

- Check `.goreleaser.yml` configuration
- Test locally with `goreleaser build --snapshot --clean`
- Review the GitHub Actions logs

## Security Notes

- ⚠️ **Never commit GPG keys** to the repository
- ⚠️ Delete `private-key.asc` and `public-key.asc` after setup
- ⚠️ Keep your GPG passphrase secure
- ⚠️ Enable 2FA on GitHub and Terraform Registry

## See Also

- [Release Process Explanation](../explanation/release-process.md) - How releases work
- [How to Release a Version](release-version.md) - Actual release steps
- [GitHub Actions Workflows Reference](../reference/github-actions-workflows.md) - Workflow specs
- [Terraform Provider Publishing](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [GoReleaser Documentation](https://goreleaser.com/)
