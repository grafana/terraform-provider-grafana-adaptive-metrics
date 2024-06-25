package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ruleAttributes(replaceOnChange bool) map[string]schema.Attribute {
	// When used in the rule resource, we must destroy / create on metric change.
	// When used in the ruleset resource, we can just update the ruleset.
	var modifiers []planmodifier.String
	if replaceOnChange {
		modifiers = append(modifiers, stringplanmodifier.RequiresReplace())
	}

	return map[string]schema.Attribute{
		"metric": schema.StringAttribute{
			Required:      true,
			Description:   "The name of the metric to be aggregated.",
			PlanModifiers: modifiers,
		},
		"match_type": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "Specifies how the metric field matches to incoming metric names. Can be 'prefix', 'suffix', or 'exact', defaults to 'exact'.",
		},

		"drop": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
			Description: "Set to true to skip both ingestion and aggregation and drop the metric entirely.",
		},
		"keep_labels": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			Description: "The array of labels to keep; labels not in this array will be aggregated.",
		},
		"drop_labels": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			Description: "The array of labels that will be aggregated.",
		},

		"aggregations": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			Description: "The array of aggregation types to calculate for this metric.",
		},

		"aggregation_interval": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "The interval at which to generate the aggregated series.",
		},
		"aggregation_delay": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "The delay until aggregation is performed.",
		},
	}
}
