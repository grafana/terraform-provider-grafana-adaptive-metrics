resource "grafana-adaptive-metrics_exemption" "ex1" {
  metric      = "prometheus_request_duration_seconds_sum"
  keep_labels = ["namespace", "cluster"]
}