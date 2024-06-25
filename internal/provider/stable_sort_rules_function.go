package provider

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &StableSortRulesFunction{}

type StableSortRulesFunction struct{}

func NewStableSortRulesFunction() function.Function {
	return &StableSortRulesFunction{}
}

func (f *StableSortRulesFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "stable_sort_rules"
}

func (f *StableSortRulesFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	ruleAttrsCopy := ruleAttributes(false)
	ruleAttrTypes := make(map[string]attr.Type, len(ruleAttrsCopy))

	for name, attr := range ruleAttrsCopy {
		ruleAttrTypes[name] = attr.GetType()
	}

	resp.Definition = function.Definition{
		Summary:     "Sort rule lists in a stable order",
		Description: "Sort a list of rules in a stable order to prevent unnecessary updates in a ruleset resource. Rules are sorted by match type, then metric name.",

		Parameters: []function.Parameter{
			function.ListParameter{
				Name:        "rules",
				Description: "The list of rules to sort.",
				ElementType: basetypes.ObjectType{
					AttrTypes: ruleAttrTypes,
				},
			},
		},
		Return: function.ListReturn{
			ElementType: basetypes.ObjectType{
				AttrTypes: ruleAttrTypes,
			},
		},
	}
}

func (f *StableSortRulesFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input []model.RuleSetRuleTF

	// Read Terraform argument data into the variable
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &input))

	sort.SliceStable(input, func(i, j int) bool {
		if input[i].MatchType.ValueString() != input[j].MatchType.ValueString() {
			return input[i].MatchType.ValueString() < input[j].MatchType.ValueString()
		}
		return input[i].Metric.ValueString() < input[j].Metric.ValueString()
	})

	// Set the result to the same data
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, input))
}
