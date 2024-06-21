package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

func TestAccRuleResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	metricName := fmt.Sprintf("test_tf_metric_%s", RandString(6))
	t.Cleanup(func() {
		aggRules := AggregationRulesForAccTest(t)
		_ = aggRules.Delete("", model.AggregationRule{Metric: metricName})
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read an existing rule w/ auto_import=false (results in an error).
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Create("", model.AggregationRule{Metric: metricName, DropLabels: []string{"foobar"}, Aggregations: []string{"sum"}}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
}
`, metricName),
				ExpectError: regexp.MustCompile("Unable to create aggregation rule"),
			},
			// Create + Read an existing rule w/ auto_import=true (results in an update).
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
	auto_import = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "auto_import", "true"),
				),
			},
			// Create + Read, no existing rule.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Delete("", model.AggregationRule{Metric: metricName}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// ImportState.
			{
				ResourceName:                         "grafana-adaptive-metrics_rule.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        metricName,
				ImportStateVerifyIdentifierAttribute: "metric",
				// The last_updated and auto_import attributes do not exist in the
				// aggregations API, therefore there is no value for it during
				// import.
				ImportStateVerifyIgnore: []string{"last_updated", "auto_import"},
			},
			// Update + Read.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Delete("", model.AggregationRule{Metric: metricName}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// Update of metric name, TF should replace the old aggregation rule with a new one in-line.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName+"_replaced"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName+"_replaced"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// Delete happens automatically.
		},
	})
}
