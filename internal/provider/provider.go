package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type TfSplitPoliciesProvider struct {
	version string
}

const VERSION = "0.1"

func (p *TfSplitPoliciesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tf-split-policies"
	resp.Version = p.version
}

func (p *TfSplitPoliciesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *TfSplitPoliciesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return NewTfSplitPoliciesDataSource()
		},
	}
}

func (p *TfSplitPoliciesProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *TfSplitPoliciesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &TfSplitPoliciesProvider{
			version: VERSION,
		}
	}
}
