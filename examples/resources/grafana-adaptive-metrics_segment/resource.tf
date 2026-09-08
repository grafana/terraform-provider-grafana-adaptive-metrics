resource "grafana-adaptive-metrics_segment" "s1" {
  name                = "mimir team"
  selector            = "{namespace=\"mimir\"}"
  fallback_to_default = true

  auto_apply = {
    enabled = true
    gate = {
      policy = "no-increase"
    }
  }
}
