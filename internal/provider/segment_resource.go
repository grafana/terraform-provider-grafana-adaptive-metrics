package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type segmentResource struct {
	client *client.Client
}

func (e *segmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var config model.SegmentTF
	var state model.SegmentTF
	var plan model.SegmentTF

	// Config = what user wrote in .tf
	req.Config.Get(ctx, &config)
	// State = what Terraform knew before
	req.State.Get(ctx, &state)
	// Plan = what Terraform has planned so far
	req.Plan.Get(ctx, &plan)

	// if policy_id was removed on already created resource it's not a default policy, override a plan to set policy_id to default policy ID.
	if !state.ID.IsNull() && config.PolicyID.IsNull() && state.PolicyID.String() != state.ID.String() {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("The segment %q will be reset to default policy.", config.Name.ValueString()),
			"User removed policy_id field, setting policy_id to default policy ID.",
		)
		plan.PolicyID = state.ID
	}

	resp.Plan.Set(ctx, &plan)
}

var (
	_ resource.Resource                = &segmentResource{}
	_ resource.ResourceWithConfigure   = &segmentResource{}
	_ resource.ResourceWithImportState = &segmentResource{}
	_ resource.ResourceWithModifyPlan  = &segmentResource{}
)

func newSegmentResource() resource.Resource {
	return &segmentResource{}
}

func (e *segmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (e *segmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_segment", req.ProviderTypeName)
}

func (e *segmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A ULID that uniquely identifies the segment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the segment.",
			},
			"selector": schema.StringAttribute{
				Required:    true,
				Description: "The selector that defines the segment.",
			},
			"fallback_to_default": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to fallback to the default segment if the selector does not match any segments.",
			},
			"auto_apply": schema.SingleNestedAttribute{
				Optional:    true,
				Description: privatePreviewWarning + "Configurations related to auto-applying recommendations.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: privatePreviewWarning + "Whether to automatically apply the generated recommendations in this segment.",
					},
				},
			},
			"policy_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: privatePreviewWarning + "ID of the policy applied to this segment.",
			},
		},
	}
}

func (e *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.SegmentTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment, err := e.client.CreateSegment(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create segment", err.Error())
		return
	}

	state := segment.ToTF()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (e *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SegmentTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment, err := e.client.ReadSegment(state.ID.ValueString())
	if err != nil {
		if client.IsErrNotFound(err) {
			resp.Diagnostics.AddWarning("Segment not found", err.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read segment", err.Error())
		return
	}

	state = segment.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (e *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.SegmentTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state model.SegmentTF
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	segment := plan.ToAPIReq()
	segment.ID = state.ID.ValueString()

	err := e.client.UpdateSegment(segment)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update segment", err.Error())
		return
	}

	segment, err = e.client.ReadSegment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read segment after updating", err.Error())
		return
	}

	state = segment.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (e *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.SegmentTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := e.client.DeleteSegment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete segment", err.Error())
	}
}

func (e *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
