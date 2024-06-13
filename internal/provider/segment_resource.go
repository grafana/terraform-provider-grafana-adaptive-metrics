package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type segmentResource struct {
	client *client.Client
}

var (
	_ resource.Resource                = &segmentResource{}
	_ resource.ResourceWithConfigure   = &segmentResource{}
	_ resource.ResourceWithImportState = &segmentResource{}
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
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A UILD that uniquely identifies the segment.",
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
				Default:     defaultBoolTrue{},
				Description: "Whether to fallback to the default segment if the selector does not match any segments.",
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

	s, err := e.client.CreateSegment(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create segment", err.Error())
		return
	}

	state := s.ToTF()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (e *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.SegmentTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ex, err := e.client.ReadSegment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read segment", err.Error())
		return
	}

	tf := ex.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
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

	ex := plan.ToAPIReq()
	ex.ID = state.ID.ValueString()

	err := e.client.UpdateSegment(ex)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update segment", err.Error())
		return
	}

	ex, err = e.client.ReadSegment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read segment after updating", err.Error())
		return
	}

	state = ex.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (e *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// var state model.SegmentTF
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// err := e.client.DeleteSegment(state.Selector.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Unable to delete segment", err.Error())
	// }
}

func (e *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
