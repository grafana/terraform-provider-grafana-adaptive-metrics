
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
}

resource "grafana-adaptive-metrics_ruleset" "default" {
  rules = data.grafana-adaptive-metrics_recommendations.default.recommendations
}
