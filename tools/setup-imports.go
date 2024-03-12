package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

func main() {
	rulesFile := flag.String("rules-file", "", "Filepath to an existing rules file.")
	flag.Parse()

	if *rulesFile == "" {
		fmt.Println("missing required argument \"--rules-file\"")
		os.Exit(1)
	}

	rulesJson, err := os.ReadFile(*rulesFile)
	if err != nil {
		fmt.Printf("could not read file %s: %v\n", *rulesFile, err)
		os.Exit(1)
	}

	var rules []model.AggregationRule
	if err = json.Unmarshal(rulesJson, &rules); err != nil {
		fmt.Printf("could not unmarshal json: %v\n", err)
		os.Exit(1)
	}

	for _, rule := range rules {
		cleanMetricName := strings.ReplaceAll(rule.Metric, ":", "_")

		fmt.Println("import {")
		fmt.Printf("  to = grafana-adaptive-metrics_rule.%s\n", cleanMetricName)
		fmt.Printf("  id = \"%s\"\n", rule.Metric)
		fmt.Println("}")
		fmt.Println()
		fmt.Printf("resource \"grafana-adaptive-metrics_rule\" \"%s\" {\n", cleanMetricName)
		fmt.Printf("  metric               = \"%s\"\n", rule.Metric)
		if rule.MatchType != "" {
			fmt.Printf("  match_type           = \"%s\"\n", rule.MatchType)
		}
		if rule.Drop {
			fmt.Printf("  drop                 = %t\n", rule.Drop)
		}
		if len(rule.KeepLabels) > 0 {
			fmt.Printf("  keep_labels          = %s\n", strList(rule.KeepLabels))
		}
		if len(rule.DropLabels) > 0 {
			fmt.Printf("  drop_labels          = %s\n", strList(rule.DropLabels))
		}
		if len(rule.Aggregations) > 0 {
			fmt.Printf("  aggregations         = %s\n", strList(rule.Aggregations))
		}
		if rule.AggregationInterval != "" {
			fmt.Printf("  aggregation_interval = \"%s\"\n", rule.AggregationInterval)
		}
		if rule.AggregationDelay != "" {
			fmt.Printf("  aggregation_delay    = \"%s\"\n", rule.AggregationDelay)
		}
		fmt.Println("}")
		fmt.Println()
	}
}

func strList(in []string) string {
	for i, s := range in {
		in[i] = fmt.Sprintf("\"%s\"", s)
	}
	return fmt.Sprintf("[%s]", strings.Join(in, ", "))
}
