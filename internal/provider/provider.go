package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure TfSplitPoliciesProvider satisfies various provider interfaces.
var _ provider.Provider = &TfSplitPoliciesProvider{}

// TfSplitPoliciesProvider defines the provider implementation.
type TfSplitPoliciesProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TfSplitPoliciesProviderModel describes the provider data model.
type TfSplitPoliciesProviderModel struct {
}

// Metadata returns the metadata for the provider
func (p *TfSplitPoliciesProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "splitpolicies"
	resp.Version = p.version
}

// Schema returns the schema for the provider
func (*TfSplitPoliciesProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	_ *provider.SchemaResponse) {
}

// Configure configures the provider
func (*TfSplitPoliciesProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse) {
	var data TfSplitPoliciesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Resources returns a slice of functions to instantiate each Resource implementation.
func (*TfSplitPoliciesProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns a slice of functions to instantiate each DataSource implementation.
func (*TfSplitPoliciesProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTfSplitPoliciesDataSource,
	}
}

// New creates a new tf-splitpolicies provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TfSplitPoliciesProvider{
			version: version,
		}
	}
}
