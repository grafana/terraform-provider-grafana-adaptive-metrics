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

type policyResource struct {
	client *client.Client
}

var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

func newPolicyResource() resource.Resource {
	return &policyResource{}
}

func (e *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	e.client = data.client
}

func (e *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_policy", req.ProviderTypeName)
}

func (e *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A policy defines the configuration for Adaptive Metrics recommendations generation. Policies can be assigned to segments to customize recommendation behavior.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A ULID that uniquely identifies the policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the policy.",
			},
			"usage_sources": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The sources of usage data to consider when generating recommendations. Valid values are: dashboard, rules, queries.",
			},
			"unused_metrics_action": schema.StringAttribute{
				Optional:    true,
				Description: "The action to take for metrics with no recorded usage. Valid values are: drop_all_labels, drop_custom_labels, keep_custom_labels, best_guess.",
			},
			"unused_metrics_action_act_on_custom_labels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The list of custom labels to act on when unused_metrics_action is drop_custom_labels or keep_custom_labels.",
			},
			"min_query_usages": schema.Int64Attribute{
				Optional:    true,
				Description: "The minimum number of query usages required for a metric to be considered for recommendations.",
			},
		},
	}
}

func (e *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.PolicyTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := e.client.CreatePolicy(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create policy: "+err.Error(), err.Error())
		return
	}

	state := policy.ToTF()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (e *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.PolicyTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := e.client.ReadPolicy(state.ID.ValueString())
	if err != nil {
		if client.IsErrNotFound(err) {
			resp.Diagnostics.AddWarning("Policy not found", err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read policy: "+err.Error(), err.Error())
		return
	}

	state = policy.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (e *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.PolicyTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state model.PolicyTF
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := plan.ToAPIReq()
	policy.ID = state.ID.ValueString()

	err := e.client.UpdatePolicy(policy)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update policy: "+err.Error(), err.Error())
		return
	}

	policy, err = e.client.ReadPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read policy after updating: "+err.Error(), err.Error())
		return
	}

	state = policy.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (e *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.PolicyTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := e.client.DeletePolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete policy: "+err.Error(), err.Error())
	}
}

func (e *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
