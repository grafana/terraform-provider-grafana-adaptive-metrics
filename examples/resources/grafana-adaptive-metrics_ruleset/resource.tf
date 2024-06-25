
# Apply an inline ruleset
resource "grafana-adaptive-metrics_ruleset" "default" {
  rules = [
    {
      "metric" : "cpu_usage_seconds_total",
      "drop_labels" : [
        "instance"
      ],
      "aggregations" : [
        "sum:counter"
      ],
    }
  ]
}

# Apply a ruleset from a file
resource "grafana-adaptive-metrics_ruleset" "default" {
  rules = jsondecode(file("${path.module}/rules.json"))
}

# Apply the latest recommendations on each apply
data "grafana-adaptive-metrics_recommendations" "default" {
  verbose = true
}

resource "grafana-adaptive-metrics_ruleset" "default" {
  # stable_sort_rules ensures that the rules are always applied in the same order, regardless of the ordering of the recommendations
  rules = provider::grafana-adaptive-metrics::stable_sort_rules(data.grafana-adaptive-metrics_recommendations.default.recommendations)
}
