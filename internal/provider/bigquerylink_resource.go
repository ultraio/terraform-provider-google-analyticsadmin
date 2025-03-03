package provider

import (
	"context"
	"fmt"
	"terraform-provider-google-analyticsadmin/internal/provider/resource_bigquerylink"
	"time"

	admin "cloud.google.com/go/analytics/admin/apiv1alpha"
	"cloud.google.com/go/analytics/admin/apiv1alpha/adminpb"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var _ resource.Resource = (*bigquerylinkResource)(nil)
var _ resource.ResourceWithConfigure = (*bigquerylinkResource)(nil)

func NewBigquerylinkResource() resource.Resource {
	return &bigquerylinkResource{}
}

type bigquerylinkResource struct {
	client *admin.AnalyticsAdminClient
}

func (r *bigquerylinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bigquerylink"
}

func (r *bigquerylinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_bigquerylink.BigquerylinkResourceSchema(ctx)
}

func (r *bigquerylinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *bigquerylinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_bigquerylink.BigquerylinkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	var link adminpb.BigQueryLink
	resp.Diagnostics.Append(marshal(ctx, &data, &link)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apireq := &adminpb.CreateBigQueryLinkRequest{
		Parent:       fmt.Sprintf("properties/%s", data.PropertyId.ValueString()),
		BigqueryLink: &link,
	}

	apiresp, err := r.client.CreateBigQueryLink(ctx, apireq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating BigQueryLink",
			"Could not create BigQueryLink, unexpected error: "+err.Error(),
		)
		return
	}

	setComputedField(&data, apiresp)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bigquerylinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_bigquerylink.BigquerylinkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	apireq := &adminpb.GetBigQueryLinkRequest{
		Name: data.Name.ValueString(),
	}

	apiresp, err := r.client.GetBigQueryLink(ctx, apireq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating BigQueryLink",
			"Could not create BigQueryLink, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(unmarshal(ctx, &data, apiresp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bigquerylinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_bigquerylink.BigquerylinkModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	var link adminpb.BigQueryLink
	resp.Diagnostics.Append(marshal(ctx, &data, &link)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resource_bigquerylink.BigquerylinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	link.Name = state.Name.ValueString()

	apireq := &adminpb.UpdateBigQueryLinkRequest{
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"daily_export_enabled"},
		},
		BigqueryLink: &link,
	}

	apiresp, err := r.client.UpdateBigQueryLink(ctx, apireq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating BigQueryLink",
			"Could not create BigQueryLink, unexpected error: "+err.Error(),
		)
		return
	}

	setComputedField(&data, apiresp)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bigquerylinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_bigquerylink.BigquerylinkModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	apireq := &adminpb.DeleteBigQueryLinkRequest{
		Name: data.Name.ValueString(),
	}

	err := r.client.DeleteBigQueryLink(ctx, apireq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating BigQueryLink",
			"Could not create BigQueryLink, unexpected error: "+err.Error(),
		)
		return
	}
}

func setComputedField(data *resource_bigquerylink.BigquerylinkModel, link *adminpb.BigQueryLink) {
	data.Name = types.StringValue(link.Name)
	data.CreateTime = types.StringValue(link.CreateTime.AsTime().Format(time.RFC3339))
}

func marshal(ctx context.Context, data *resource_bigquerylink.BigquerylinkModel, link *adminpb.BigQueryLink) diag.Diagnostics {
	// Create API call logic
	var streams []string
	diags := data.ExportStreams.ElementsAs(ctx, &streams, false)
	if diags.HasError() {
		return diags
	}

	var events []string
	diags = data.ExcludeEvents.ElementsAs(ctx, &events, false)
	if diags.HasError() {
		return diags
	}

	*link = adminpb.BigQueryLink{
		Project:                 data.Project.ValueString(),
		DailyExportEnabled:      data.DailyExport.ValueBool(),
		StreamingExportEnabled:  data.StreamingExport.ValueBool(),
		FreshDailyExportEnabled: data.FreshDailyExport.ValueBool(),
		IncludeAdvertisingId:    data.IncludeAdvertisingId.ValueBool(),
		ExportStreams:           streams,
		ExcludedEvents:          events,
		DatasetLocation:         data.DatasetLocation.ValueString(),
	}

	return diag.Diagnostics{}
}

func unmarshal(ctx context.Context, data *resource_bigquerylink.BigquerylinkModel, link *adminpb.BigQueryLink) diag.Diagnostics {
	events, diags := types.ListValueFrom(ctx, types.StringType, link.ExcludedEvents)
	if diags.HasError() {
		return diags
	}

	streams, diags := types.ListValueFrom(ctx, types.StringType, link.ExportStreams)
	if diags.HasError() {
		return diags
	}

	data.DailyExport = types.BoolValue(link.DailyExportEnabled)
	data.ExcludeEvents = events
	data.ExportStreams = streams
	data.FreshDailyExport = types.BoolValue(link.FreshDailyExportEnabled)
	data.IncludeAdvertisingId = types.BoolValue(link.IncludeAdvertisingId)
	data.StreamingExport = types.BoolValue(link.StreamingExportEnabled)
	return diag.Diagnostics{}
}
