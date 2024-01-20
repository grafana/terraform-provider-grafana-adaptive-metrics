// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/client"
)

// Ensure AdaptiveMetricsProvider satisfies various provider interfaces.
var _ provider.Provider = &AdaptiveMetricsProvider{}

// AdaptiveMetricsProvider defines the provider implementation.
type AdaptiveMetricsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AdaptiveMetricsProviderModel describes the provider data model.
type AdaptiveMetricsProviderModel struct {
	URL              types.String `tfsdk:"url"`
	APIKey           types.String `tfsdk:"api_key"`
	HTTPHeaders      types.Map    `tfsdk:"http_headers"`
	Retries          types.Int64  `tfsdk:"retries"`
	RetryStatusCodes types.Set    `tfsdk:"retry_status_codes"`
	RetryWait        types.Int64  `tfsdk:"retry_wait"`

	UserAgent types.String `json:"-" tfsdk:"-"`
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

func (p *AdaptiveMetricsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "adaptive-metrics"
	resp.Version = p.version
}

func (p *AdaptiveMetricsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Grafana Cloud's API URL. May alternatively be set via the `GRAFANA_CLOUD_API_URL` environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Access Policy Token (or API key) for Grafana Cloud. May alternatively be set via the `GRAFANA_CLOUD_API_KEY` environment variable.",
			},
			"http_headers": schema.MapAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "HTTP headers mapping keys to values used for accessing Grafana Cloud APIs. May alternatively be set via the `GRAFANA_CLOUD_HTTP_HEADERS` environment variable in JSON format.",
				ElementType:         types.StringType,
			},
			"retries": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The amount of retries to use for Grafana API and Grafana Cloud API calls. Defaults to 3. May alternatively be set via the `GRAFANA_CLOUD_RETRIES` environment variable.",
			},
			"retry_status_codes": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "The status codes to retry on for Grafana API and Grafana Cloud API calls. Use `x` as a digit wildcard. Defaults to 429 and 5xx. May alternatively be set via the `GRAFANA_CLOUD_RETRY_STATUS_CODES` environment variable.",
				ElementType:         types.StringType,
			},
			"retry_wait": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The amount of time in seconds to wait between retries for Grafana Cloud API calls. May alternatively be set via the `GRAFANA_CLOUD_RETRY_WAIT` environment variable.",
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

	apiURL := getStringOverriddenByEnvOrDefault(cfg.URL, "GRAFANA_CLOUD_API_URL", "")
	if apiURL == "" {
		resp.Diagnostics.AddError("Missing required attribute", "url")
		return
	}

	apiKey := getStringOverriddenByEnvOrDefault(cfg.APIKey, "GRAFANA_CLOUD_API_KEY", "")
	retries, err := getIntOverriddenByEnvOrDefault(cfg.Retries, "GRAFANA_CLOUD_RETRIES", 3)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse GRAFANA_CLOUD_RETRIES", err.Error())
		return
	}
	retryTimeout, err := getIntOverriddenByEnvOrDefault(cfg.RetryWait, "GRAFANA_CLOUD_RETRY_WAIT", 0)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse GRAFANA_CLOUD_RETRY_WAIT", err.Error())
		return
	}

	var retryStatusCodes []string
	if envRetryStatusCodes := os.Getenv("GRAFANA_CLOUD_RETRY_STATUS_CODES"); envRetryStatusCodes != "" {
		retryStatusCodes = strings.Split(envRetryStatusCodes, ",")
	} else if !cfg.RetryStatusCodes.IsNull() {
		for _, v := range cfg.RetryStatusCodes.Elements() {
			if vStr, ok := v.(types.String); ok {
				retryStatusCodes = append(retryStatusCodes, vStr.ValueString())
			} else {
				resp.Diagnostics.AddError("Non-string value in retry_status_codes", fmt.Sprintf("got %v", v))
			}
		}
	} else {
		retryStatusCodes = []string{
			"429",
			"5xx",
			"401", // In high-load scenarios, Grafana sometimes returns 401s.
		}
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
		APIKey:           apiKey,
		NumRetries:       retries,
		RetryTimeout:     time.Second * time.Duration(retryTimeout),
		RetryStatusCodes: retryStatusCodes,
		HTTPHeaders:      httpHeaders,
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
	resp.ResourceData = aggRules
}

func (p *AdaptiveMetricsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newRuleResource,
	}
}

func (p *AdaptiveMetricsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AdaptiveMetricsProvider{
			version: version,
		}
	}
}
