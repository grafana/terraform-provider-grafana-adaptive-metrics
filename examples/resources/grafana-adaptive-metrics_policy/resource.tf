resource "grafana-adaptive-metrics_policy" "example" {
  name = "example-policy"

  usage_sources = ["dashboard", "rules", "queries"]

  unused_metrics_action = "drop_custom_labels"

  unused_metrics_action_act_on_custom_labels = [
    "cluster",
    "namespace",
    "environment"
  ]

  min_query_usages = 5
}
