package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TfSplitPoliciesDataSource struct{}

// enforce interface impl
var _ datasource.DataSource = &TfSplitPoliciesDataSource{}

func NewTfSplitPoliciesDataSource() datasource.DataSource {
	return &TfSplitPoliciesDataSource{}
}

type TfSplitPoliciesDataSourceModel struct {
	MaximumChunkSize types.Int64 `tfsdk:"maximum_chunk_size"`
	Hash     types.String   `tfsdk:"hash"`
	Policies []types.String `tfsdk:"policies"`
	Chunks   [][]types.String `tfsdk:"chunks"`
}

func (d *TfSplitPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// we name it like the provider as this is our only exported data source type
	resp.TypeName = "tf_split_policies"
}

func (d *TfSplitPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	var l = types.ListType{}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Splitting multiple JSON AWS policies into chunks of a given maximum size",
		Attributes: map[string]schema.Attribute{
			"hash": schema.StringAttribute{
				MarkdownDescription: "The hash of all the inputs, usefull for update triggering",
				Computed:            true,
			},
			"chunks": schema.ListAttribute{
				ElementType: l.WithElementType(l.WithElementType(types.StringType)),
				Computed:    true,
			},
			"policies": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"maximum_chunk_size": schema.Int64Attribute {
				Optional: true,
			},
		},
	}
}

func (d *TfSplitPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TfSplitPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var hash = hashInputs(&data)
	data.Hash = types.StringValue(hash)
	tflog.Trace(ctx, "set the hash value of the input policies")

	// Convert from TF wrapped strings
	policies := []string{}
	for _, policy := range data.Policies {
		policies = append(policies, policy.ValueString())
	}

	// get the maximum size or default
	var maximum_size int = 6144
	if !data.MaximumChunkSize.IsUnknown() {
		// FIXME: we are truncating the width of the integer here, probably okay as we shouldn't need 64bits anyway. TF only exposes that int type in its schema.
		maximum_size = int(data.MaximumChunkSize.ValueInt64())
	}

	// Split the policies into chunks of at most 6144 characters
	chunks, err := split_policies(policies, maximum_size)
	if err != "" {
		d := diag.NewErrorDiagnostic("Failed to chunk policies", err)
		resp.Diagnostics.Append(d)
		return
	}

	// Convert to TF wrapped strings
	tfchunks := [][]types.String{}
	for _, chunk := range chunks {
		c := []types.String{}
		for _, item := range chunk {
			c = append(c, types.StringValue(item))
		}

		tfchunks = append(tfchunks, c)
	}
	tflog.Trace(ctx, "split the policies into chunks")

	// Finally hand the data back to TF
	data.Chunks = tfchunks
	

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
