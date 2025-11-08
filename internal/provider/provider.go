package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &spantreeUtilsProvider{}
)

// spantreeUtilsProvider is the provider implementation.
type spantreeUtilsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &spantreeUtilsProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *spantreeUtilsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utils"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *spantreeUtilsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Terraform provider offering utility functions including template rendering with custom placeholder syntax (@@VAR@@) " +
			"that doesn't conflict with other templating systems like Argo Workflows ({{}}), " +
			"shell scripts (${}, $()), bash syntax, or Windows batch files (%%).",
	}
}

// Configure prepares the provider for data sources and resources.
func (p *spantreeUtilsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No provider-level configuration needed for this simple provider
}

// DataSources defines the data sources implemented in the provider.
func (p *spantreeUtilsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRenderTemplateDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *spantreeUtilsProvider) Resources(_ context.Context) []func() resource.Resource {
	// This provider only implements data sources, no resources
	return []func() resource.Resource{}
}
