package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccSegmentResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	t.Cleanup(func() {
		c := ClientForAccTest(t)
		segments, err := c.ListSegments()
		require.NoError(t, err)

		for _, s := range segments {
			err = c.DeleteSegment(s.ID)
			require.NoError(t, err)
		}
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_segment" "test" {
	name = "test segment"
	selector = "{namespace=\"test\"}"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "name", "test segment"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "selector", "{namespace=\"test\"}"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "fallback_to_default", "true"),
				),
			},
			// ImportState.
			{
				ResourceName:      "grafana-adaptive-metrics_segment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update + Read.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_segment" "test" {
	name = "test segment 2"
	selector = "{namespace=\"test\"}"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "name", "test segment 2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "selector", "{namespace=\"test\"}"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "fallback_to_default", "true"),
				),
			},
			// Update + Read, setting fallback_to_default=true
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_segment" "test" {
	name = "test segment 2"
	selector = "{namespace=\"test\"}"
	fallback_to_default = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "name", "test segment 2"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "selector", "{namespace=\"test\"}"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_segment.test", "fallback_to_default", "false"),
				),
			},
			// Delete happens automatically.
		},
	})
}
