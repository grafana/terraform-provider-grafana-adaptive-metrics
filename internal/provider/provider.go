// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
)

const privatePreviewWarning = "WARNING: contact Grafana Cloud support before use. This feature is in private preview and may change without notice, including in ways that may break your configuration. "

// Ensure AdaptiveMetricsProvider satisfies various provider interfaces.
var _ provider.Provider = &AdaptiveMetricsProvider{}

// AdaptiveMetricsProvider defines the provider implementation.
type AdaptiveMetricsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// commit is set to the provider commit on release, "unknown" when the
	// provider is built and ran locally or when running acceptance
	// testing.
	commit string
}

// AdaptiveMetricsProviderModel describes the provider data model.
type AdaptiveMetricsProviderModel struct {
	URL         types.String `tfsdk:"url"`
	APIKey      types.String `tfsdk:"api_key"`
	HTTPHeaders types.Map    `tfsdk:"http_headers"`
	Retries     types.Int64  `tfsdk:"retries"`
	Debug       types.Bool   `tfsdk:"debug"`
}

func getStringOverriddenByEnvOrDefault(s types.String, envKey string, valDefault string) string {
	val, ok := os.LookupEnv(envKey)
	if ok {
		return val
	}

	if !s.IsNull() {
		return s.ValueString()
	}

	return valDefault
}

func getIntOverriddenByEnvOrDefault(s types.Int64, envKey string, valDefault int) (int, error) {
	val, ok := os.LookupEnv(envKey)
	if ok {
		return strconv.Atoi(val)
	}

	if !s.IsNull() {
		return int(s.ValueInt64()), nil
	}

	return valDefault, nil
}

func getBooleanOverriddenByEnvOrDefault(s types.Bool, envKey string, valDefault bool) (bool, error) {
	val, ok := os.LookupEnv(envKey)
	if ok {
		return strconv.ParseBool(val)
	}

	if !s.IsNull() {
		return s.ValueBool(), nil
	}

	return valDefault, nil
}

func (p *AdaptiveMetricsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "grafana-adaptive-metrics"
	resp.Version = p.version
}

func (p *AdaptiveMetricsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Grafana Cloud's API URL. May alternatively be set via the `GRAFANA_AM_API_URL` environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Tenant ID and Access Policy Token (or API key) for Grafana Cloud in the format '<tenant-id>:<token-or-api-key>'. May alternatively be set via the `GRAFANA_AM_API_KEY` environment variable.",
			},
			"http_headers": schema.MapAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "HTTP headers mapping keys to values used for accessing Grafana Cloud APIs. May alternatively be set via the `GRAFANA_AM_HTTP_HEADERS` environment variable in JSON format.",
				ElementType:         types.StringType,
			},
			"retries": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The amount of retries to use for Grafana API and Grafana Cloud API calls. Defaults to 3. May alternatively be set via the `GRAFANA_AM_RETRIES` environment variable.",
			},
			"debug": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable debug logging. Defaults to false.",
			},
		},
	}
}

func (p *AdaptiveMetricsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg AdaptiveMetricsProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := getStringOverriddenByEnvOrDefault(cfg.URL, "GRAFANA_AM_API_URL", "")
	if apiURL == "" {
		resp.Diagnostics.AddError("Missing required attribute 'url'", "This may alternatively be set via the `GRAFANA_AM_API_URL` environment variable.")
		return
	}

	apiKey := getStringOverriddenByEnvOrDefault(cfg.APIKey, "GRAFANA_AM_API_KEY", "")
	debug, err := getBooleanOverriddenByEnvOrDefault(cfg.Debug, "GRAFANA_AM_DEBUG", false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse GRAFANA_AM_DEBUG", err.Error())
		return
	}
	retries, err := getIntOverriddenByEnvOrDefault(cfg.Retries, "GRAFANA_AM_RETRIES", 3)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse GRAFANA_AM_RETRIES", err.Error())
		return
	}
	httpClient := cleanhttp.DefaultClient()
	if retries > 0 {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = retries
		httpClient = retryClient.StandardClient()
	}

	httpHeaders := make(map[string]string)
	if envHeaders := os.Getenv("GRAFANA_HTTP_HEADERS"); envHeaders != "" {
		err = json.Unmarshal([]byte(envHeaders), &httpHeaders)
		if err != nil {
			resp.Diagnostics.AddError("Failed to parse GRAFANA_HTTP_HEADERS", err.Error())
			return
		}
	} else if !cfg.HTTPHeaders.IsNull() {
		for k, v := range cfg.HTTPHeaders.Elements() {
			if vStr, ok := v.(types.String); ok {
				httpHeaders[k] = vStr.ValueString()
			} else {
				resp.Diagnostics.AddError("Non-string value in http_headers", fmt.Sprintf("got %v for key %s", v, k))
			}
		}
	}

	c, err := client.New(apiURL, &client.Config{
		APIKey:      apiKey,
		HTTPHeaders: httpHeaders,
		Debug:       debug,
		HttpClient:  httpClient,
		UserAgent:   fmt.Sprintf("Terraform/%s grafana-adaptive-metrics-provider/%s (commit:%s)", req.TerraformVersion, p.version, p.commit),
	})
	if err != nil {
		resp.Diagnostics.AddError("Could not instantiate the API client.", err.Error())
		return
	}

	aggRules := NewAggregationRules(c)
	if err = aggRules.Init(); err != nil {
		resp.Diagnostics.AddError("Could not initialize internal state.", err.Error())
		return
	}

	resp.DataSourceData = c // TODO
	resp.ResourceData = &resourceData{
		aggRules: aggRules,
		client:   c,
	}
}

func (p *AdaptiveMetricsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newRuleSetResource,
		newRuleResource,
		newExemptionResource,
		newRecommendationsConfigResource,
		newSegmentResource,
		newPolicyResource,
	}
}

func (p *AdaptiveMetricsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newRecommendationDatasource,
	}
}

func New(version string, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &AdaptiveMetricsProvider{
			version: version,
			commit:  commit,
		}
	}
}

type resourceData struct {
	aggRules *AggregationRules
	client   *client.Client
}
