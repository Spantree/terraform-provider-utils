# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-11-07

### Added

- Initial release of terraform-provider-template
- `template_render` data source for rendering templates with custom placeholder syntax
- Support for `@@VAR@@` placeholder syntax that doesn't conflict with:
  - Argo Workflows: `{{inputs.parameters.*}}`
  - Shell scripts: `${VAR}`, `$(command)`
  - Bash: `[[ ]]`, `<<`, `>>`, `<>`, `&&`, `||`, `(( ))`
  - Windows batch files: `%%VAR%%`
- Validation that all placeholders have corresponding values
- Comprehensive test suite with 11 test cases
- Examples for basic usage and Argo Workflow deployment
- Complete documentation (README, QUICKSTART, LICENSE)
- MIT License for maximum compatibility and ease of use

### Design Decisions

- **Placeholder Syntax**: Chose `@@VAR@@` over `%%VAR%%` to avoid conflicts with Windows batch file syntax
- **Data Source Only**: Implemented as a pure data source (no resources) since template rendering doesn't create infrastructure
- **Strict Validation**: All placeholders must have values to prevent silent failures
- **Modern Framework**: Built with Terraform Plugin Framework (not legacy SDK)
