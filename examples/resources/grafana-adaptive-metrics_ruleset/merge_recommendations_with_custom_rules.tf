# This example demonstrates how to merge recommendations from the recommendations data source
# with custom rules, giving precedence to custom rules when there are conflicts.
#
# This is the recommended approach when you want to:
# - Apply most recommendations automatically
# - Override specific recommendations with custom rules
# - Add additional custom rules not in recommendations

# Fetch the latest recommendations
data "grafana-adaptive-metrics_recommendations" "default" {
  # Optionally filter by action
  # action = ["add", "update"]
}

# Define your custom rules as local values for easier management. These could also be defined in a file.
locals {
  custom_rules = [
    {
      metric       = "http_requests_total"
      drop_labels  = ["instance", "pod"]
      aggregations = ["sum:counter"]
    },
    {
      metric       = "cpu_usage_seconds_total"
      drop_labels  = ["container"]
      aggregations = ["sum:counter", "avg:gauge"]
    }
  ]

  # Get the metric names from custom rules for filtering
  custom_metrics = toset([for rule in local.custom_rules : rule.metric])

  # Filter out recommendations that conflict with custom rules
  # This gives precedence to your custom rules over recommendations
  filtered_recommendations = [
    for rec in data.grafana-adaptive-metrics_recommendations.default.recommendations :
    rec if !contains(local.custom_metrics, rec.metric)
  ]

  # Merge custom rules with filtered recommendations
  # Custom rules come first to ensure they take precedence
  merged_rules = concat(local.custom_rules, local.filtered_recommendations)
}

# Apply the merged ruleset
resource "grafana-adaptive-metrics_ruleset" "default" {
  rules = local.merged_rules
}

# Alternatively, if you want to completely exclude certain recommendations
# (e.g., drop rules you don't want to apply), you can add additional filtering:
locals {
  # Metrics to exclude from recommendations entirely
  excluded_metrics = toset([
    "critical_metric_to_preserve",
    "another_important_metric"
  ])

  # Filter out both custom rules and excluded metrics from recommendations
  advanced_filtered_recommendations = [
    for rec in data.grafana-adaptive-metrics_recommendations.default.recommendations :
    rec if !contains(local.custom_metrics, rec.metric) && !contains(local.excluded_metrics, rec.metric)
  ]

  advanced_merged_rules = concat(local.custom_rules, local.advanced_filtered_recommendations)
}

# Example with exclusions
resource "grafana-adaptive-metrics_ruleset" "with_exclusions" {
  rules = local.advanced_merged_rules
}

