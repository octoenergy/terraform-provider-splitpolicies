package provider

import (
	"context"
	"fmt"
	"net/http"

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

func NewTfSplitPoliciesDataSource() datasource.DataSource {
	return &TfSplitPoliciesDataSource{}
}

// TfSplitPoliciesDataSource defines the data source implementation.
type TfSplitPoliciesDataSource struct {
	client *http.Client
}

// TfSplitPoliciesDataSourceModel describes the data source data model.
type TfSplitPoliciesDataSourceModel struct {
	Id               types.String   `tfsdk:"id"`
	MaximumChunkSize types.Int64    `tfsdk:"maximum_chunk_size"`
	Hash             types.String   `tfsdk:"hash"`
	Policies         []types.String `tfsdk:"policies"`
	Chunks           types.Map      `tfsdk:"chunks"`
}

func (d *TfSplitPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

var chunksType = basetypes.ListType{ElemType: basetypes.StringType{}}

func (d *TfSplitPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *TfSplitPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
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
	var maximumSize = 6144
	if !data.MaximumChunkSize.IsNull() {
		maximumSize = int(data.MaximumChunkSize.ValueInt64())
	}

	// Split the policies into chunks of at most 6144 characters
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
	data.Id = types.StringValue("some-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
