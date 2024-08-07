package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ClientConfigDataSourceModel struct {
	ID             types.String   `tfsdk:"id"`
	TenantID       types.String   `tfsdk:"tenant_id"`
	SubscriptionID types.String   `tfsdk:"subscription_id"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

type ClientConfigDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &ClientConfigDataSource{}
var _ datasource.DataSourceWithConfigure = &ClientConfigDataSource{}

func (r *ClientConfigDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ClientConfigDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_client_config"
}

func (r *ClientConfigDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},

			"tenant_id": schema.StringAttribute{
				Computed: true,
			},

			"subscription_id": schema.StringAttribute{
				Computed: true,
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Read: true,
			}),
		},
	}
}

func (r *ClientConfigDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ClientConfigDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	subscriptionId := r.ProviderData.Account.GetSubscriptionId()
	tenantId := r.ProviderData.Account.GetTenantId()

	model.ID = types.StringValue(fmt.Sprintf("clientConfigs/subscriptionId=%s;tenantId=%s", subscriptionId, tenantId))
	model.SubscriptionID = types.StringValue(subscriptionId)
	model.TenantID = types.StringValue(tenantId)
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
