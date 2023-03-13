package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TfSplitPoliciesDataSource{}

// NewTfSplitPoliciesDataSource creates a new tf-splitpolicies data source
func NewTfSplitPoliciesDataSource() datasource.DataSource {
	return &TfSplitPoliciesDataSource{}
}

// TfSplitPoliciesDataSource defines the data source implementation.
type TfSplitPoliciesDataSource struct {
}

// TfSplitPoliciesDataSourceModel describes the data source data model.
type TfSplitPoliciesDataSourceModel struct {
	ID               types.String   `tfsdk:"id"`
	MaximumChunkSize types.Int64    `tfsdk:"maximum_chunk_size"`
	Hash             types.String   `tfsdk:"hash"`
	Policies         []types.String `tfsdk:"policies"`
	Chunks           types.Map      `tfsdk:"chunks"`
}

// Metadata provides the name of the data source
func (*TfSplitPoliciesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

var chunksType = basetypes.ListType{ElemType: basetypes.StringType{}}

// Schema returns the schema for this data source
func (*TfSplitPoliciesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Splitting multiple JSON AWS policies into chunks of a given maximum size",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"hash": schema.StringAttribute{
				MarkdownDescription: "The hash of all the inputs, usefull for update triggering",
				Computed:            true,
			},
			"chunks": schema.MapAttribute{
				ElementType: chunksType,
				Computed:    true,
			},
			"policies": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"maximum_chunk_size": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

// Configure configures the data source
func (*TfSplitPoliciesDataSource) Configure(
	_ context.Context,
	_ datasource.ConfigureRequest,
	_ *datasource.ConfigureResponse) {
}

// Read is called when the provider reads from the data source
func (*TfSplitPoliciesDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse) {
	var data TfSplitPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var hash, err = hashInputs(&data)
	if err != nil {
		newDiagnostic := diag.NewErrorDiagnostic("Failed to chunk policies", err.Error())
		resp.Diagnostics.Append(newDiagnostic)
		return
	}
	data.Hash = types.StringValue(hash)
	tflog.Trace(ctx, "set the hash value of the input policies")

	// Convert from TF wrapped strings
	policies := []string{}
	for _, policy := range data.Policies {
		policies = append(policies, policy.ValueString())
	}

	// get the maximum size or default
	var maximumSize = 6144
	if !data.MaximumChunkSize.IsNull() {
		maximumSize = int(data.MaximumChunkSize.ValueInt64())
	}

	// Split the policies into chunks of at most maximumSize characters
	chunks, err := splitPolicies(policies, maximumSize)
	if err != nil {
		newDiagnostic := diag.NewErrorDiagnostic("Failed to chunk policies", err.Error())
		resp.Diagnostics.Append(newDiagnostic)
		return
	}

	// Convert to TF wrapped strings
	tfchunks := make(map[string]attr.Value)
	for i, chunk := range chunks {
		v, newDiagnostic := types.ListValueFrom(ctx, types.StringType, chunk)
		resp.Diagnostics.Append(newDiagnostic...)
		if resp.Diagnostics.HasError() {
			return
		}

		tfchunks[fmt.Sprintf("%d", i)] = v
	}

	tflog.Trace(ctx, "split the policies into chunks")

	// Finally hand the data back to TF
	mv, newDiagnostic := basetypes.NewMapValue(chunksType, tfchunks)
	resp.Diagnostics.Append(newDiagnostic...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Chunks = mv
	data.ID = types.StringValue("some-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
