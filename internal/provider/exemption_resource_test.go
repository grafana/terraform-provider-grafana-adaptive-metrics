package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccExemptionResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	var exemptionID string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_exemption" "test" {
	metric = "test_tf_metric"
	keep_labels = ["namespace"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "metric", "test_tf_metric"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.0", "namespace"),
				),
			},
			// ImportState.
			{
				ResourceName:      "grafana-adaptive-metrics_exemption.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the
				// aggregations API, therefore there is no value for it during
				// import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_exemption" "test" {
	metric = "test_tf_metric"
	keep_labels = ["namespace", "cluster"]
	reason = "testing"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "metric", "test_tf_metric"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.0", "namespace"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.1", "cluster"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "reason", "testing"),
					func(s *terraform.State) error {
						// Capture the exemption ID for later use.
						exemptionID = s.RootModule().Resources["grafana-adaptive-metrics_exemption.test"].Primary.ID
						return nil
					},
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					client := ClientForAccTest(t)
					require.NoError(t, client.DeleteExemption("", exemptionID))
				},
				Config: providerConfig + `
resource "grafana-adaptive-metrics_exemption" "test" {
	metric = "test_tf_metric"
	keep_labels = ["namespace", "cluster"]
	reason = "testing"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "metric", "test_tf_metric"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.0", "namespace"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "keep_labels.1", "cluster"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "reason", "testing"),
				),
			},
			// Update + Read, setting disable_recommendations=true
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_exemption" "test" {
	metric = "test_tf_metric"
	disable_recommendations = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "metric", "test_tf_metric"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_exemption.test", "disable_recommendations", "true"),
				),
			},
			// Delete happens automatically.
		},
	})
}
