package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExemptionResource(t *testing.T) {
	CheckAccTestsEnabled(t)

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
				),
			},
			// Delete happens automatically.
		},
	})
}
