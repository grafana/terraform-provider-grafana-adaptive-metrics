package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
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

func (r *ruleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_rule", req.ProviderTypeName)
}

func (r *ruleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	ruleSchemaCopy := schema.Schema{
		Attributes: ruleAttributes(true),
	}
	// These fields are not part of the shared schema, but are used by the provider to manage the resource.
	ruleSchemaCopy.Attributes["auto_import"] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
		Description: "When set to true, the rule will be automatically imported if it is not already in Terraform state.",
	}
	ruleSchemaCopy.Attributes["segment"] = schema.StringAttribute{
		Optional:    true,
		Description: "The name of the segment to aggregate metrics for.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	resp.Schema = ruleSchemaCopy
}

func (r *ruleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.RuleTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.AutoImport.ValueBool() {
		_, err := r.rules.Read(plan.Segment.ValueString(), plan.Metric.ValueString())
		if err != nil {
			// There is no existing rule for this metric; create it.
			err := r.rules.Create(plan.Segment.ValueString(), plan.ToAPIReq())
			if err != nil {
				resp.Diagnostics.AddError("Unable to create aggregation rule", err.Error())
				return
			}
		} else {
			// There is an existing rule for this metric; update it.
			err := r.rules.Update(plan.Segment.ValueString(), plan.ToAPIReq())
			if err != nil {
				resp.Diagnostics.AddError("Unable to update aggregation rule", err.Error())
				return
			}

			resp.Diagnostics.AddWarning("Existing aggregation rule for metric found", "The existing rule has been updated and imported into Terraform state; no aggregation rule has been created.")
		}
	} else {
		err := r.rules.Create(plan.Segment.ValueString(), plan.ToAPIReq())
		if err != nil {
			resp.Diagnostics.AddError("Unable to create aggregation rule", err.Error())
			return
		}
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ruleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.RuleTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.rules.Read(state.Segment.ValueString(), state.Metric.ValueString())
	if err != nil {
		if client.IsErrNotFound(err) {
			resp.Diagnostics.AddWarning("Aggregation rule not found", err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read aggregation rule", err.Error())
		return
	}

	tf := rule.ToTF()

	// Segment tells us where to put the rule later, but isn't actually a part
	// of the rule, so we set it separately.
	tf.Segment = state.Segment

	// AutoImport is a meta field used by this Terraform provider; the API never returns
	// a value for it so we keep it updated separately.
	tf.AutoImport = state.AutoImport

	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
}

func (r *ruleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.RuleTF
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state model.RuleTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.Update(plan.Segment.ValueString(), plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update aggregation rule", err.Error())
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ruleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.RuleTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.rules.Delete(state.Segment.ValueString(), state.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete aggregation rule", err.Error())
	}
}

func (r *ruleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("metric"), req, resp)
}
