package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccPolicyResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	var policyID string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_policy" "test" {
	name = "test policy"
	usage_sources = ["dashboard", "rules"]
	unused_metrics_action = "drop_custom_labels"
	unused_metrics_action_act_on_custom_labels = ["cluster", "namespace"]
	min_query_usages = 5
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "name", "test policy"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.0", "dashboard"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.1", "rules"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action", "drop_custom_labels"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action_act_on_custom_labels.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action_act_on_custom_labels.0", "cluster"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action_act_on_custom_labels.1", "namespace"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "min_query_usages", "5"),
					resource.TestCheckResourceAttrSet("grafana-adaptive-metrics_policy.test", "id"),
				),
			},
			// ImportState.
			{
				ResourceName:      "grafana-adaptive-metrics_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_policy" "test" {
	name = "test policy updated"
	usage_sources = ["dashboard", "rules", "queries"]
	unused_metrics_action = "best_guess"
	min_query_usages = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "name", "test policy updated"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.#", "3"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.0", "dashboard"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.1", "rules"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.2", "queries"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action", "best_guess"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action_act_on_custom_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "min_query_usages", "10"),
					func(s *terraform.State) error {
						// Capture the policy ID for later use.
						policyID = s.RootModule().Resources["grafana-adaptive-metrics_policy.test"].Primary.ID
						return nil
					},
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					client := ClientForAccTest(t)
					require.NoError(t, client.DeletePolicy(policyID))
				},
				Config: providerConfig + `
resource "grafana-adaptive-metrics_policy" "test" {
	name = "test policy updated"
	usage_sources = ["dashboard", "rules", "queries"]
	unused_metrics_action = "best_guess"
	min_query_usages = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "name", "test policy updated"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources.#", "3"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action", "best_guess"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "min_query_usages", "10"),
				),
			},
			// Update to minimal policy (only name required).
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_policy" "test" {
	name = "minimal policy"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_policy.test", "name", "minimal policy"),
					resource.TestCheckNoResourceAttr("grafana-adaptive-metrics_policy.test", "usage_sources"),
					resource.TestCheckNoResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action_act_on_custom_labels"),
					resource.TestCheckNoResourceAttr("grafana-adaptive-metrics_policy.test", "min_query_usages"),
					resource.TestCheckNoResourceAttr("grafana-adaptive-metrics_policy.test", "unused_metrics_action"),
				),
			},
			// Delete happens automatically.
		},
	})
}
