# Terraform Spantree Utils Provider Documentation

Welcome to the documentation for the Terraform Spantree Utils Provider. This provider offers utility functions including template rendering with custom placeholder syntax that doesn't conflict with other templating systems.

## Documentation Structure

This documentation follows the [Diataxis framework](https://diataxis.fr/), organizing content into four categories based on your needs:

### 📚 [Tutorials](tutorials/) - *Learning by doing*

Step-by-step lessons to help you get started and learn the fundamentals. Perfect for beginners.

- [Getting Started](tutorials/getting-started.md) - Your first steps with the provider
- [Quickstart](tutorials/quickstart.md) - Quick reference for local development
- [Basic Template Rendering](tutorials/basic-template-rendering.md) - Learn template syntax and rendering
- [Deploying an Argo Workflow](tutorials/argo-workflow-deployment.md) - Complete real-world example

### 🔧 [How-to Guides](how-to/) - *Solving specific problems*

Practical guides for accomplishing specific tasks. For users who know what they want to do.

- [Set Up Local Development](how-to/local-development.md) - Configure dev environment
- [Set Up Releases and CI/CD](how-to/setup-release.md) - One-time GitHub Actions and release setup
- [Set Up Documentation Generation](how-to/setup-documentation-generation.md) - Configure doc generation
- [Release a New Version](how-to/release-version.md) - Publishing releases

### 📖 [Reference](reference/) - *Technical information*

Technical descriptions and specifications. For looking up details.

- [Data Source: render_template](reference/data-source-render-template.md) - Complete API reference
- [Placeholder Syntax](reference/placeholder-syntax.md) - Syntax rules and validation
- [GitHub Actions Workflows](reference/github-actions-workflows.md) - Workflow specifications

### 💡 [Explanation](explanation/) - *Understanding concepts*

Conceptual information to deepen your understanding. For learning the "why" behind decisions.

- [Why This Provider?](explanation/why-this-provider.md) - Problem statement and motivation
- [Design Requirements](explanation/design-requirements.md) - Original design goals and requirements
- [Documentation System](explanation/documentation-system.md) - How the dual documentation system works
- [Release Process](explanation/release-process.md) - How releases work

## Quick Links

### New to the provider?

Start with [Getting Started Tutorial](tutorials/getting-started.md)

### Need to accomplish something specific?

Check the [How-to Guides](how-to/)

### Looking for API details?

See the [Reference Documentation](reference/)

### Want to understand the concepts?

Read the [Explanations](explanation/)

## Additional Resources

- [GitHub Repository](https://github.com/spantree/terraform-provider-utils)
- [Terraform Registry](https://registry.terraform.io/providers/spantree/utils)
- [Examples](../examples/) - Working code examples
- [CHANGELOG](../CHANGELOG.md) - Version history

## Contributing

Contributions are welcome! See:

- [Local Development Setup](how-to/local-development.md)
- [Release Process](how-to/release-version.md)

## Support

- **Issues**: [GitHub Issues](https://github.com/spantree/terraform-provider-spantree/issues)
- **Discussions**: [GitHub Discussions](https://github.com/spantree/terraform-provider-spantree/discussions)

## License

This provider is released under the MIT License. See [LICENSE](../LICENSE) for details.
