package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

func TestRuleResourceReadMatchType(t *testing.T) {
	tests := []struct {
		name      string
		prior     string
		apiFields string
		want      string
	}{
		{name: "omitted default keeps explicit exact", prior: "exact", want: "exact"},
		{name: "empty default keeps explicit exact", prior: "exact", apiFields: `,"match_type":""`, want: "exact"},
		{name: "omitted default keeps implicit exact"},
		{name: "omitted default exposes drift from prefix", prior: "prefix"},
		{name: "prefix exposes drift from exact", prior: "exact", apiFields: `,"match_type":"prefix"`, want: "prefix"},
		{name: "suffix exposes drift from exact", prior: "exact", apiFields: `,"match_type":"suffix"`, want: "suffix"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				if req.URL.Path == "/aggregations/segmented_rules" {
					_, err := fmt.Fprint(w, "[]")
					assert.NoError(t, err)
					return
				}
				assert.Equal(t, "/aggregations/rule/test_metric", req.URL.Path)
				w.Header().Set("ETag", `"test-etag"`)
				_, err := fmt.Fprintf(w, `{"metric":"test_metric"%s}`, tc.apiFields)
				assert.NoError(t, err)
			}))
			defer server.Close()

			api, err := client.New(server.URL, &client.Config{HttpClient: server.Client()})
			require.NoError(t, err)
			rules := NewAggregationRules(api)
			require.NoError(t, rules.Init())
			rule := &ruleResource{rules: rules}
			ctx := t.Context()
			var schemaResp frameworkresource.SchemaResponse
			rule.Schema(ctx, frameworkresource.SchemaRequest{}, &schemaResp)
			state := tfsdk.State{Schema: schemaResp.Schema}
			prior := model.AggregationRule{Metric: "test_metric", MatchType: tc.prior}.ToTF()
			require.False(t, state.Set(ctx, prior).HasError())

			resp := frameworkresource.ReadResponse{State: state}
			rule.Read(ctx, frameworkresource.ReadRequest{State: state}, &resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			var refreshed model.RuleTF
			require.False(t, resp.State.Get(ctx, &refreshed).HasError())
			require.Equal(t, types.StringValue(tc.want), refreshed.MatchType)
		})
	}
}

func TestAccRuleResource(t *testing.T) {
	CheckAccTestsEnabled(t)

	metricName := fmt.Sprintf("test_tf_metric_%s", RandString(6))
	t.Cleanup(func() {
		aggRules := AggregationRulesForAccTest(t)
		_ = aggRules.Delete("", model.AggregationRule{Metric: metricName})
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create + Read an existing rule w/ auto_import=false (results in an error).
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Create("", model.AggregationRule{Metric: metricName, DropLabels: []string{"foobar"}, Aggregations: []string{"sum"}}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
}
`, metricName),
				ExpectError: regexp.MustCompile("Unable to create aggregation rule"),
			},
			// Create + Read an existing rule w/ auto_import=true (results in an update).
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
	auto_import = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "auto_import", "true"),
				),
			},
			// Create + Read, no existing rule.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Delete("", model.AggregationRule{Metric: metricName}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	drop = true
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "true"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// ImportState.
			{
				ResourceName:                         "grafana-adaptive-metrics_rule.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        metricName,
				ImportStateVerifyIdentifierAttribute: "metric",
				// The last_updated and auto_import attributes do not exist in the
				// aggregations API, therefore there is no value for it during
				// import.
				ImportStateVerifyIgnore: []string{"last_updated", "auto_import"},
			},
			// Update + Read.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// External delete of resource, TF should recreate it.
			{
				PreConfig: func() {
					aggRules := AggregationRulesForAccTest(t)
					require.NoError(t, aggRules.Delete("", model.AggregationRule{Metric: metricName}))
				},
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// Update of metric name, TF should replace the old aggregation rule with a new one in-line.
			{
				Config: providerConfig + fmt.Sprintf(`
resource "grafana-adaptive-metrics_rule" "test" {
	metric = "%s"
	match_type = "prefix"
	drop_labels = [ "instance" ]
	aggregations = [ "sum" ]
}
`, metricName+"_replaced"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "metric", metricName+"_replaced"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "match_type", "prefix"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop", "false"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "keep_labels.#", "0"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "drop_labels.0", "instance"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.#", "1"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregations.0", "sum"),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_interval", ""),
					resource.TestCheckResourceAttr("grafana-adaptive-metrics_rule.test", "aggregation_delay", ""),
				),
			},
			// Delete happens automatically.
		},
	})
}
