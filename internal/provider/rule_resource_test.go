package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRuleResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	metricName := fmt.Sprintf("test_tf_metric_%s", RandString(6))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "match_type", ""),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "drop", "true"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "drop_labels.#", "0"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregations.#", "0"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregation_delay", ""),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "ingest", "false"),
				),
			},
			// ImportState.
			{
				ResourceName:                         "adaptive-metrics_rule.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        metricName,
				ImportStateVerifyIdentifierAttribute: "metric",
				// The last_updated attribute does not exist in the
				// aggregations API, therefore there is no value for it during
				// import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update + Read.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
	ingest = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "aggregation_delay", ""),
					resource.TestCheckResourceAttr("adaptive-metrics_rule.test", "ingest", "true"),
				),
			},
			// Delete happens automatically.
		},
	})
}
