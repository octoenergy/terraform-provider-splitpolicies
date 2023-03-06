package provider

import (
	"context"
	"net/http"

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

func (p *TfSplitPoliciesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tf-split-policies"
	resp.Version = p.version
}

func (p *TfSplitPoliciesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
}

func (p *TfSplitPoliciesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TfSplitPoliciesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TfSplitPoliciesProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *TfSplitPoliciesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTfSplitPoliciesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TfSplitPoliciesProvider{
			version: version,
		}
	}
}
