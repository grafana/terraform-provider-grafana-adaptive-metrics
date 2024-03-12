data "grafana-adaptive-metrics_recommendations" "all" {
}

output "recs" {
  value = data.grafana-adaptive-metrics_recommendations.all
}