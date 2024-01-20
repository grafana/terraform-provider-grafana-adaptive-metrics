package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/model"
)

type ruleResource struct {
	rules *AggregationRules
}

var (
	_ resource.Resource                = &ruleResource{}
	_ resource.ResourceWithConfigure   = &ruleResource{}
	_ resource.ResourceWithImportState = &ruleResource{}
)

func newRuleResource() resource.Resource {
	return &ruleResource{}
}

func (r *ruleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*AggregationRules)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.rules = data
}

func (r *ruleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_rule", req.ProviderTypeName)
}

func (r *ruleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metric": schema.StringAttribute{
				Required: true,
			},
			"match_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},

			"drop": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaultBoolFalse{},
			},
			"keep_labels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     defaultEmptyList{},
			},
			"drop_labels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     defaultEmptyList{},
			},

			"aggregations": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     defaultEmptyList{},
			},

			"aggregation_interval": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"aggregation_delay": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},

			"ingest": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  defaultBoolFalse{},
			},
		},
	}
}

func (r *ruleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.RuleTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.Create(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create aggregation rule", err.Error())
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ruleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.RuleTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.rules.Read(state.Metric.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read aggregation rule", err.Error())
	}

	tf := rule.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
}

func (r *ruleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.RuleTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.Update(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update aggregation rule", err.Error())
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ruleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.RuleTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.Delete(state.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete aggregation rule", err.Error())
	}
}

func (r *ruleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("metric"), req, resp)
}
