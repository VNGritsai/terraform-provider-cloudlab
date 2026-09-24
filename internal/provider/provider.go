package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/VNGritsai/terraform-provider-cloudlab/internal/client"
)

var _ provider.Provider = &cloudlabProvider{}

type cloudlabProvider struct {
	version string
}

type cloudlabProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *cloudlabProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudlab"
	resp.Version = p.version
}

func (p *cloudlabProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Провайдер для управления ресурсами cloudlab.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:    true,
				Description: "Базовый URL API cloudlab, например http://localhost:8080.",
			},
		},
	}
}

func (p *cloudlabProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data cloudlabProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.IsNull() || data.Endpoint.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Не задан endpoint",
			"Укажите атрибут endpoint в блоке provider \"cloudlab\".",
		)
		return
	}

	c := client.NewClient(data.Endpoint.ValueString())

	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *cloudlabProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInstanceResource,
	}
}

func (p *cloudlabProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cloudlabProvider{
			version: version,
		}
	}
}
