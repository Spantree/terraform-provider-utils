package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &renderTemplateDataSource{}
)

// NewRenderTemplateDataSource is a helper function to simplify the provider implementation.
func NewRenderTemplateDataSource() datasource.DataSource {
	return &renderTemplateDataSource{}
}

// renderTemplateDataSource is the data source implementation.
type renderTemplateDataSource struct{}

// renderTemplateDataSourceModel maps the data source schema data.
type renderTemplateDataSourceModel struct {
	Template types.String `tfsdk:"template"`
	Values   types.Map    `tfsdk:"values"`
	Result   types.String `tfsdk:"result"`
}

// Metadata returns the data source type name.
func (d *renderTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_render_template"
}

// Schema defines the schema for the data source.
func (d *renderTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Renders a template by replacing placeholders with values. " +
			"Uses @@VAR@@ syntax for placeholders to avoid conflicts with other templating systems. " +
			"All placeholders must have corresponding values or an error will be returned.",
		MarkdownDescription: "Renders a template by replacing placeholders with values.\n\n" +
			"**Placeholder Syntax:** `@@VAR@@`\n\n" +
			"This syntax is designed to not conflict with:\n" +
			"- Argo Workflows: `{{inputs.parameters.*}}`\n" +
			"- Shell variables: `${VAR}`, `$(command)`\n" +
			"- Bash conditionals: `[[ ]]`\n" +
			"- Heredoc/redirection: `<<`, `>>`, `<>`\n" +
			"- Logical operators: `&&`, `||`\n" +
			"- Arithmetic: `(( ))`\n" +
			"- Windows batch files: `%%VAR%%`\n\n" +
			"All placeholders must have corresponding values or an error will be returned.",
		Attributes: map[string]schema.Attribute{
			"template": schema.StringAttribute{
				Description: "The template string containing placeholders in @@VAR@@ format. " +
					"Placeholder names must be alphanumeric or underscore.",
				MarkdownDescription: "The template string containing placeholders in `@@VAR@@` format. " +
					"Placeholder names must be alphanumeric or underscore.",
				Required: true,
			},
			"values": schema.MapAttribute{
				Description: "Map of placeholder names to their replacement values. " +
					"All placeholders in the template must have a corresponding value.",
				MarkdownDescription: "Map of placeholder names to their replacement values. " +
					"All placeholders in the template must have a corresponding value.",
				ElementType: types.StringType,
				Required:    true,
			},
			"result": schema.StringAttribute{
				Description:         "The rendered template with all placeholders replaced by their values.",
				MarkdownDescription: "The rendered template with all placeholders replaced by their values.",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *renderTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data renderTemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template := data.Template.ValueString()

	// Convert types.Map to map[string]string
	values := make(map[string]string)
	diags := data.Values.ElementsAs(ctx, &values, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Rendering template", map[string]any{
		"template_length": len(template),
		"values_count":    len(values),
	})

	// Render the template
	result, err := renderTemplate(template, values)
	if err != nil {
		resp.Diagnostics.AddError(
			"Template Rendering Failed",
			fmt.Sprintf("Failed to render template: %s", err.Error()),
		)
		return
	}

	data.Result = types.StringValue(result)

	tflog.Trace(ctx, "Successfully rendered template")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// renderTemplate performs the actual template rendering.
// It finds all @@VAR@@ placeholders and replaces them with values from the map.
// Returns an error if any placeholder doesn't have a corresponding value.
func renderTemplate(template string, values map[string]string) (string, error) {
	// Regular expression to match @@VAR@@ where VAR is alphanumeric or underscore
	placeholderRegex := regexp.MustCompile(`@@([A-Za-z0-9_]+)@@`)

	// Find all unique placeholders in the template
	matches := placeholderRegex.FindAllStringSubmatch(template, -1)
	uniquePlaceholders := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			uniquePlaceholders[match[1]] = true
		}
	}

	// Check that all placeholders have values
	var missingPlaceholders []string
	for placeholder := range uniquePlaceholders {
		if _, exists := values[placeholder]; !exists {
			missingPlaceholders = append(missingPlaceholders, placeholder)
		}
	}

	if len(missingPlaceholders) > 0 {
		return "", fmt.Errorf("missing values for placeholders: %s", strings.Join(missingPlaceholders, ", "))
	}

	// Replace all placeholders with their values
	result := placeholderRegex.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the placeholder name (without @@)
		placeholder := match[2 : len(match)-2]
		if value, exists := values[placeholder]; exists {
			return value
		}
		// This shouldn't happen due to the check above, but handle it safely
		return match
	})

	return result, nil
}
