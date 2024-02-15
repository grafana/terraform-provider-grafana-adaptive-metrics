resource "grafana-adaptive-metrics_rule" "agent_request_duration_seconds_sum" {
  metric       = "agent_request_duration_seconds_sum"
  drop_labels  = ["namespace", "pod"]
  aggregations = ["sum:counter"]
}

resource "grafana-adaptive-metrics_rule" "prometheus_request_duration_seconds_sum" {
  metric       = "prometheus_request_duration_seconds_sum"
  drop_labels  = ["container", "instance", "ws"]
  aggregations = ["sum:counter"]
}
