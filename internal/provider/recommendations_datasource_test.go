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

	testAttr := func(resource, attr, value string) resource.TestCheckFunc {
		return checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations."+resource, "am_terraform_provider_acceptance_test_metric", attr, value)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// This test requires a rule and exemption are present in order to guarantee some recommendations are returned.
			// We do this in a separate segment to avoid conflicts with the ruleSetResourceTest.
			{
				Config: providerConfig + `
resource "grafana-adaptive-metrics_exemption" "test" {
	metric = "am_terraform_provider_acceptance_test_metric"
	disable_recommendations = true
	segment = "01JQVN6036Z18P6Z958JNNTXRP"
}

resource "grafana-adaptive-metrics_rule" "test" {
	metric = "am_terraform_provider_acceptance_test_metric"
	drop_labels = ["this", "metric", "doesnt", "exist"]
	aggregations = ["count"]
	segment = "01JQVN6036Z18P6Z958JNNTXRP"
}

data "grafana-adaptive-metrics_recommendations" "non_verbose" {
	segment = "01JQVN6036Z18P6Z958JNNTXRP"
}

data "grafana-adaptive-metrics_recommendations" "verbose" {
	segment = "01JQVN6036Z18P6Z958JNNTXRP"
	verbose = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAttr("non_verbose", "metric", "am_terraform_provider_acceptance_test_metric"),
					testAttr("non_verbose", "drop", "false"),
					testAttr("non_verbose", "keep_labels.#", "0"),
					testAttr("non_verbose", "drop_labels.#", "4"),
					testAttr("non_verbose", "drop_labels.0", "this"),
					testAttr("non_verbose", "drop_labels.1", "metric"),
					testAttr("non_verbose", "drop_labels.2", "doesnt"),
					testAttr("non_verbose", "drop_labels.3", "exist"),
					testAttr("non_verbose", "aggregations.#", "1"),
					testAttr("non_verbose", "aggregations.0", "count"),
					testAttr("non_verbose", "aggregation_interval", ""),
					testAttr("non_verbose", "aggregation_delay", ""),
					testAttr("non_verbose", "recommended_action", ""),
					testAttr("non_verbose", "usages_in_rules", "0"),
					testAttr("non_verbose", "usages_in_queries", "0"),
					testAttr("non_verbose", "usages_in_dashboards", "0"),
					testAttr("non_verbose", "kept_labels.#", "0"),
					testAttr("non_verbose", "total_series_before_aggregation", "0"),
					testAttr("non_verbose", "total_series_after_aggregation", "0"),

					testAttr("verbose", "metric", "am_terraform_provider_acceptance_test_metric"),
					testAttr("verbose", "drop", "false"),
					testAttr("verbose", "keep_labels.#", "0"),
					testAttr("verbose", "drop_labels.#", "4"),
					testAttr("verbose", "drop_labels.0", "this"),
					testAttr("verbose", "drop_labels.1", "metric"),
					testAttr("verbose", "drop_labels.2", "doesnt"),
					testAttr("verbose", "drop_labels.3", "exist"),
					testAttr("verbose", "aggregations.#", "1"),
					testAttr("verbose", "aggregations.0", "count"),
					testAttr("verbose", "recommended_action", "keep"),
					testAttr("verbose", "usages_in_rules", "0"),
					testAttr("verbose", "usages_in_queries", "0"),
					testAttr("verbose", "usages_in_dashboards", "0"),
					testAttr("verbose", "kept_labels.#", "0"),
					testAttr("verbose", "total_series_before_aggregation", "0"),
					testAttr("verbose", "total_series_after_aggregation", "0"),
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
