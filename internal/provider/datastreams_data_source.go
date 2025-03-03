package provider

import (
	"context"
	"fmt"
	"terraform-provider-google-analyticsadmin/internal/provider/datasource_datastreams"

	admin "cloud.google.com/go/analytics/admin/apiv1alpha"
	"cloud.google.com/go/analytics/admin/apiv1alpha/adminpb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/iterator"
)

var _ datasource.DataSource = (*datastreamsDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*datastreamsDataSource)(nil)

func NewDatastreamsDataSource() datasource.DataSource {
	return &datastreamsDataSource{}
}

type datastreamsDataSource struct {
	client *admin.AnalyticsAdminClient
}

func (d *datastreamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datastreams"
}

func (d *datastreamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*admin.AnalyticsAdminClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *admin.AnalyticsAdminClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *datastreamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_datastreams.DatastreamsDataSourceSchema(ctx)
}

func (d *datastreamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_datastreams.DatastreamsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	it := d.client.ListDataStreams(ctx, &adminpb.ListDataStreamsRequest{
		Parent: fmt.Sprintf("properties/%s", data.PropertyId.ValueString()),
	})

	var streams []datasource_datastreams.DataStreamsValue
	for {
		apiresp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating DataStreams",
				"Could not read DataStreams, unexpected error: "+err.Error(),
			)
			return
		}

		streams = append(streams, datasource_datastreams.DataStreamsValue{
			Name:            types.StringValue(apiresp.Name),
			DisplayName:     types.StringValue(apiresp.DisplayName),
			DataStreamsType: types.StringValue(apiresp.Type.String()),
		})
	}

	list, diags := types.ListValueFrom(ctx, datasource_datastreams.DataStreamsValue{}.Type(ctx), streams)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DataStreams = list

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
