package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecommendationsConfigResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_recommendations_config" "test" {
	keep_labels = ["foobar"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.0", "foobar"),
				),
			},
			// Update + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_recommendations_config" "test" {
	keep_labels = ["foobar", "foobaz"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.0", "foobar"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.1", "foobaz"),
				),
			},
			// Update enabling auto_apply + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_recommendations_config" "test" {
	keep_labels = ["foobar", "foobaz"]
	auto_apply = {
		enabled = true
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.#", "2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.0", "foobar"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "keep_labels.1", "foobaz"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_recommendations_config.test", "auto_apply.enabled", "true"),
				),
			},
			// Delete happens automatically.
		},
	})
}
