resource "grafana-adaptive-metrics_recommendations_config" "singleton" {
  keep_labels = ["namespace"]
}
