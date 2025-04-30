package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type ruleSetResource struct {
	rules *AggregationRules
}

var (
	_ resource.Resource                = &ruleSetResource{}
	_ resource.ResourceWithConfigure   = &ruleSetResource{}
	_ resource.ResourceWithImportState = &ruleSetResource{}
)

func newRuleSetResource() resource.Resource {
	return &ruleSetResource{}
}

func (r *ruleSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.rules = data.aggRules
}

func (r *ruleSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_ruleset", req.ProviderTypeName)
}

func (r *ruleSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"segment": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the segment to aggregate metrics for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: ruleAttributes(false),
				},
			},
		},
	}
}

func (r *ruleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.RuleSetTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This object is a singleton per segment, so we don't need to check if it already exists.
	err := r.rules.UpdateRuleSet(plan.Segment.ValueString(), plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update aggregation rule set", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ruleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.RuleSetTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.rules.ReadRuleSet(state.Segment.ValueString())
	if err != nil {
		if client.IsErrNotFound(err) {
			resp.Diagnostics.AddWarning("Ruleset not found", err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read ruleset", err.Error())
		return
	}

	// Prevent unnecessary drift due to reordering
	rules = model.AlignUpstreamWithState(state.ToAPIReq(), rules)

	tf := rules.ToTF(state.Segment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
}

func (r *ruleSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.RuleSetTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state model.RuleSetTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.UpdateRuleSet(plan.Segment.ValueString(), plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update aggregation rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ruleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.RuleSetTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.UpdateRuleSet(state.Segment.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete aggregation rule", err.Error())
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *ruleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The default segment is a special case where we don't have an ID to import.
	if req.ID == "default" {
		resp.State.SetAttribute(ctx, path.Root("segment"), types.StringNull())
	} else {
		resource.ImportStatePassthroughID(ctx, path.Root("segment"), req, resp)
	}
}
