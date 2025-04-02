package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRecommendationDatasource(t *testing.T) {
	CheckAccTestsEnabled(t)

	testAttr := func(attr, value string) resource.TestCheckFunc {
		return checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", attr, value)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read non-verbose.
			{
				Config: providerConfig + `
data "grafana-adaptive-metrics_recommendations" "test" {
  segment = "01JQVN6036Z18P6Z958JNNTXRP"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAttr("metric", "am_terraform_provider_acceptance_test_metric"),
					testAttr("drop", "false"),
					testAttr("keep_labels.#", "0"),
					testAttr("drop_labels.#", "4"),
					testAttr("drop_labels.0", "this"),
					testAttr("drop_labels.1", "metric"),
					testAttr("drop_labels.2", "doesnt"),
					testAttr("drop_labels.3", "exist"),
					testAttr("aggregations.#", "1"),
					testAttr("aggregations.0", "count"),
					testAttr("aggregation_interval", ""),
					testAttr("aggregation_delay", ""),
					testAttr("recommended_action", ""),
					testAttr("usages_in_rules", "0"),
					testAttr("usages_in_queries", "0"),
					testAttr("usages_in_dashboards", "0"),
					testAttr("kept_labels.#", "0"),
					testAttr("total_series_before_aggregation", "0"),
					testAttr("total_series_after_aggregation", "0"),
				),
			},
			// Read verbose.
			{
				Config: providerConfig + `
data "grafana-adaptive-metrics_recommendations" "test" {
	segment = "01JQVN6036Z18P6Z958JNNTXRP"
	verbose = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAttr("metric", "am_terraform_provider_acceptance_test_metric"),
					testAttr("drop", "false"),
					testAttr("keep_labels.#", "0"),
					testAttr("drop_labels.#", "4"),
					testAttr("drop_labels.0", "this"),
					testAttr("drop_labels.1", "metric"),
					testAttr("drop_labels.2", "doesnt"),
					testAttr("drop_labels.3", "exist"),
					testAttr("aggregations.#", "1"),
					testAttr("aggregations.0", "count"),
					testAttr("recommended_action", "keep"),
					testAttr("usages_in_rules", "0"),
					testAttr("usages_in_queries", "0"),
					testAttr("usages_in_dashboards", "0"),
					testAttr("kept_labels.#", "0"),
					testAttr("total_series_before_aggregation", "0"),
					testAttr("total_series_after_aggregation", "0"),
				),
			},
		},
	})
}

var metricPathRegex = regexp.MustCompile(`recommendations\.\d+\.metric`)

func findPrimaryInstance(s *terraform.State, name string) (*terraform.InstanceState, error) {
	root := s.RootModule()
	r, ok := root.Resources[name]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", name)
	}

	if r.Primary == nil {
		return nil, fmt.Errorf("primary instance not found: %s", name)
	}

	return r.Primary, nil
}

// checkMetricRecommendationAttr finds the recommendation for a metric and
// checks the value of an attribute. Recommendations do not have a predictable
// order, so we need to find the recommendation for the metric first.
func checkMetricRecommendationAttr(name, metric, attr, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		primary, err := findPrimaryInstance(s, name)
		if err != nil {
			return err
		}

		var prefix string
		for k, v := range primary.Attributes {
			if metricPathRegex.MatchString(k) && v == metric {
				prefix = strings.TrimSuffix(k, ".metric")
				break
			}
		}

		attrPath := prefix + "." + attr
		v, ok := primary.Attributes[attrPath]
		if !ok {
			attrsWithPrefix := make([]string, 0, len(primary.Attributes))
			for k := range primary.Attributes {
				if strings.HasPrefix(k, prefix+".") {
					attrsWithPrefix = append(attrsWithPrefix, k)
				}
			}

			return fmt.Errorf("Attribute not found: %s, available attributes with prefix %s: %v", attrPath, prefix, attrsWithPrefix)
		}

		if v != value {
			return fmt.Errorf("Expected %s to be %s, got %s", attrPath, value, v)
		}

		return nil
	}
}
