package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccRuleSetResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	metricName := fmt.Sprintf("test_tf_metric_%s", RandString(6))
	t.Cleanup(func() {
		aggRules := AggregationRulesForAccTest(t)
		_ = aggRules.UpdateRuleSet("", nil)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a new ruleset
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_ruleset" "test" {
	rules = [{
		metric = "%s"
		drop = true
	}]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_delay", ""),
				),
			},
			// Create + Read, no existing rule.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.UpdateRuleSet("", nil))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_ruleset" "test" {
	rules = [{
		metric = "%s"
		drop = true
	}]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_delay", ""),
				),
			},
			// Update + Read.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_ruleset" "test" {
	rules = [{
		metric = "%s"
		match_type = "prefix"
		drop_labels = [ "instance" ]
		aggregations = [ "sum" ]
	}]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_delay", ""),
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.UpdateRuleSet("", nil))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_ruleset" "test" {
	rules = [{
		metric = "%s"
		match_type = "prefix"
		drop_labels = [ "instance" ]
		aggregations = [ "sum" ]
	}]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_delay", ""),
				),
			},
			// Update of metric name, TF should replace the old aggregation rule with a new one in-line.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_ruleset" "test" {
	rules = [{
		metric = "%s"
		match_type = "prefix"
		drop_labels = [ "instance" ]
		aggregations = [ "sum" ]
	}]
}
`, metricName+"_replaced"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.metric", metricName+"_replaced"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_ruleset.test", "rules.0.aggregation_delay", ""),
				),
			},
			// Delete happens automatically.
		},
	})
}
