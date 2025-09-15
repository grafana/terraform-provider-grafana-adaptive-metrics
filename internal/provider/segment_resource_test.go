package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccSegmentResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	t.Cleanup(func() {
		c := ClientForAccTest(t)
		segments, err := c.ListSegments()
		require.NoError(t, err)

		for _, s := range segments {
			if s.ID == "01JQVN6036Z18P6Z958JNNTXRP" {
				// Recommendations test segment, do not delete.
				continue
			}
			err = c.DeleteSegment(s.ID)
			require.NoError(t, err)
		}
	})

	var segmentID string
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
					resource.TestCheckResourceAttrSet("grafana-adaptive-metrics_segment.test", "policy_id"),
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
					func(s *terraform.State) error {
						// Capture the exemption ID for later use.
						segmentID = s.RootModule().Resources["grafana-adaptive-metrics_segment.test"].Primary.ID
						return nil
					},
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					client := ClientForAccTest(t)
					require.NoError(t, client.DeleteSegment(segmentID))
				},
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
