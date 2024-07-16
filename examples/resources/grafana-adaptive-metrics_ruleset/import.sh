
# Import the ruleset from the default segment
terraform import grafana-adaptive-metrics_ruleset.rules default

# Import the ruleset from a custom segment
terraform import grafana-adaptive-metrics_ruleset.rules $CUSTOM_SEGMENT_ID