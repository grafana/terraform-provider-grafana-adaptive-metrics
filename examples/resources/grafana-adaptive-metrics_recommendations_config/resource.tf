resource "grafana-adaptive-metrics_resource_config" "singleton" {
  keep_labels = ["namespace"]
}