package provider

import (
	"context"

	admin "cloud.google.com/go/analytics/admin/apiv1alpha"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/api/option"
)

var _ provider.Provider = (*analyticsadminProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &analyticsadminProvider{
			version: version,
		}
	}
}

type analyticsadminProvider struct {
	version string
}

type analyticsadminProviderModel struct {
	Scopes types.List `tfsdk:"scopes"`
}

func (p *analyticsadminProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with AnalyticsAdmin.",
		Attributes: map[string]schema.Attribute{
			"scopes": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The Google analytics admin scopes. Check <https://developers.google.com/identity/protocols/oauth2/scopes#analytics>",
				Required:    true,
			},
		},
	}
}

func (p *analyticsadminProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configure AnalyticsAdmin client")

	// Retrieve provider data from configuration
	var config analyticsadminProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get oauth scopes", map[string]interface{}{"scopes": config.Scopes.String()})

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Scopes.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("scopes"),
			"Unknown AnalyticsAdmin API scopes",
			"The provider cannot create the AnalyticsAdmin API client as there is an unknown configuration value for the Analytics API scopes. "+
				"Target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "scopes", config.Scopes.String())
	tflog.Debug(ctx, "Creating HashiCups client")

	// Create a new HashiCups client using the configuration values
	var scopes []string
	config.Scopes.ElementsAs(ctx, &scopes, false)

	client, err := admin.NewAnalyticsAdminClient(ctx, option.WithScopes(scopes...))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create AnalyticsAdmin Client",
			"An unexpected error occurred when creating the AnalyticsAdmin API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"AnalyticsAdmin Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured AnalyticsAdmin client", map[string]any{"success": true})
}

func (p *analyticsadminProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "analyticsadmin"
}

func (p *analyticsadminProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *analyticsadminProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBigquerylinkResource,
	}
}
