package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type recommendationsConfigResource struct {
	client *client.Client
}

var (
	_ resource.Resource              = &recommendationsConfigResource{}
	_ resource.ResourceWithConfigure = &recommendationsConfigResource{}
)

func newRecommendationsConfigResource() resource.Resource {
	return &recommendationsConfigResource{}
}

func (r *recommendationsConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*resourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.client
}

func (r *recommendationsConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_recommendations_config", req.ProviderTypeName)
}

func (r *recommendationsConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"keep_labels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Description: "The array of labels to keep; labels not in this array will be aggregated.",
			},
			"auto_apply": schema.SingleNestedAttribute{
				Optional:    true,
				Description: privatePreviewWarning + "Configurations related to auto-applying recommendations.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: privatePreviewWarning + "Whether to automatically apply the generated recommendations in the default segment.",
					},
				},
			},
		},
	}
}

func (r *recommendationsConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.AggregationRecommendationConfigurationTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAggregationRecommendationsConfig(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update recommendations config", err.Error())
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	resp.Diagnostics.AddWarning(
		"Resource has not been created",
		"The recommendations config is a singleton that always exists for every tenant. Creating it adds it to Terraform state but nothing is actually created.",
	)
}

func (r *recommendationsConfigResource) Read(ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
	cfg, err := r.client.AggregationRecommendationsConfig()
	if err != nil {
		resp.Diagnostics.AddError("Unable to read recommendations config", err.Error())
		return
	}

	tf := cfg.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
}

func (r *recommendationsConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.AggregationRecommendationConfigurationTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAggregationRecommendationsConfig(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update recommendations config", err.Error())
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *recommendationsConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Resource has not been deleted",
		"The recommendations config is a singleton that always exists for every tenant. Deleting it removes it from Terraform state but does nothing to the underlying resource.",
	)
}
