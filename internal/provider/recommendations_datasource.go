package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type recommendationDatasource struct {
	client *client.Client
}

var (
	_ datasource.DataSource              = &recommendationDatasource{}
	_ datasource.DataSourceWithConfigure = &recommendationDatasource{}
)

func newRecommendationDatasource() datasource.DataSource {
	return &recommendationDatasource{}
}

func (r *recommendationDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected datasource configure type",
			fmt.Sprintf("Got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data
}

func (r *recommendationDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_recommendations", req.ProviderTypeName)
}

func (r *recommendationDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"segment": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the segment to get recommendations for.",
			},
			"verbose": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, the response will include additional information about the recommendation, such as the number of rules, queries, and dashboards that use the metric.",
			},
			"action": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Limit the types of recommended actions to list. Valid recommended actions are 'add', 'remove', 'keep', and 'update'. Defaults to listing all actions.",
			},
			"recommendations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"metric": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the metric to be aggregated.",
						},
						"match_type": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies how the metric field matches to incoming metric names. Can be 'prefix', 'suffix', or 'exact', defaults to 'exact'.",
						},

						"drop": schema.BoolAttribute{
							Computed:    true,
							Description: "Set to true to skip both ingestion and aggregation and drop the metric entirely.",
						},
						"keep_labels": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The array of labels to keep; labels not in this array will be aggregated.",
						},
						"drop_labels": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The array of labels that will be aggregated.",
						},

						"aggregations": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The array of aggregation types to calculate for this metric.",
						},

						"aggregation_interval": schema.StringAttribute{
							Computed:    true,
							Description: "The interval at which to generate the aggregated series.",
						},
						"aggregation_delay": schema.StringAttribute{
							Computed:    true,
							Description: "The delay until aggregation is performed.",
						},

						"recommended_action": schema.StringAttribute{
							Computed:    true,
							Description: "The recommended action for the aggregation rule.",
						},

						"usages_in_rules": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of rules that use this metric.",
						},

						"usages_in_queries": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of queries that use this metric..",
						},

						"usages_in_dashboards": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of dashboards that use this metric.",
						},

						"kept_labels": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The array of labels that will be kept.",
						},

						"total_series_after_aggregation": schema.Int64Attribute{
							Computed:    true,
							Description: "The total number of series after aggregation.",
						},

						"total_series_before_aggregation": schema.Int64Attribute{
							Computed:    true,
							Description: "The total number of series before aggregation.",
						},
					},
				},
			},
		},
	}
}

func (r *recommendationDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.AggregationRecommendationListTF
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	recs, err := r.client.AggregationRecommendations(state.Segment.ValueString(), state.IsVerbose(), state.GetActionIn())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read aggregation rule", err.Error())
		return
	}

	for _, ar := range recs {
		state.Recommendations = append(state.Recommendations, ar.ToTF())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
