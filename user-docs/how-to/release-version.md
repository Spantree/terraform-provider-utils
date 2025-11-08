# How to Release a New Version

This guide walks you through releasing a new version of the provider to the Terraform Registry.

## Prerequisites

- Maintainer access to the repository
- GPG key configured (see [Setup Release Guide](../../SETUP_RELEASE.md))
- All tests passing on `main` branch

## Quick Steps

### 1. Update the Changelog

Edit `CHANGELOG.md`:

```markdown
## [1.1.0] - 2024-01-15

### Added
- New feature X

### Fixed
- Bug Y

### Changed
- Updated dependency Z
```

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes

### 2. Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare for v1.1.0 release"
git push origin main
```

### 3. Create and Push Tag

```bash
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0
```

### 4. Monitor GitHub Actions

Go to the [Actions tab](https://github.com/spantree/terraform-provider-spantree/actions) and watch the release workflow.

### 5. Review and Publish

1. Go to [Releases](https://github.com/spantree/terraform-provider-spantree/releases)
2. Review the draft release
3. Edit release notes if needed
4. Click **Publish release**

### 6. Verify on Registry

Wait 1-2 minutes, then check:
<https://registry.terraform.io/providers/spantree/utils>

## Troubleshooting

**Workflow failed**: Check the [Release Documentation](../explanation/release-process.md)

**Not on registry**: Verify webhook and GPG signature

## See Also

- [Complete Release Documentation](../explanation/release-process.md)
- [Setup Release Guide](setup-release.md)
- [GitHub Actions Workflows Reference](../reference/github-actions-workflows.md)
