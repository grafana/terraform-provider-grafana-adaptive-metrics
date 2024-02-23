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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read non-verbose.
			{
				Config: providerConfig + `
data "grafana-adaptive-metrics_recommendations" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "metric", "am_terraform_provider_acceptance_test_metric"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop", "false"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "keep_labels.#", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.#", "4"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.0", "this"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.1", "metric"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.2", "doesnt"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.3", "exist"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregations.#", "1"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregations.0", "count"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregation_interval", ""),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregation_delay", ""),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "ingest", "false"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "recommended_action", ""),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_rules", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_queries", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_dashboards", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "kept_labels.#", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "total_series_before_aggregation", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "total_series_after_aggregation", "0"),
				),
			},
			// Read verbose.
			{
				Config: providerConfig + `
data "grafana-adaptive-metrics_recommendations" "test" {
	verbose = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "metric", "am_terraform_provider_acceptance_test_metric"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop", "false"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "keep_labels.#", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.#", "4"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.0", "this"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.1", "metric"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.2", "doesnt"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "drop_labels.3", "exist"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregations.#", "1"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "aggregations.0", "count"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "recommended_action", "keep"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_rules", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_queries", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "usages_in_dashboards", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "kept_labels.#", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "total_series_before_aggregation", "0"),
					checkMetricRecommendationAttr("data.grafana-adaptive-metrics_recommendations.test", "am_terraform_provider_acceptance_test_metric", "total_series_after_aggregation", "0"),
				),
			},
		},
	})
}

var metricPathRegex = regexp.MustCompile(`recommendations\.\d+\.metric`)

// checkMetricRecommendationAttr finds the recommendation for a metric and
// checks the value of an attribute. Recommendations do not have a predictable
// order, so we need to find the recommendation for the metric first.
func checkMetricRecommendationAttr(name, metric, attr, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		root := s.RootModule()
		r, ok := root.Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s", name)
		}

		if r.Primary == nil {
			return fmt.Errorf("Primary instance not found: %s", name)
		}

		var prefix string
		primary := r.Primary
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
