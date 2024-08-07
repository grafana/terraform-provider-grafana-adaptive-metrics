---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "grafana-adaptive-metrics Provider"
subcategory: ""
description: |-
  
---

# grafana-adaptive-metrics Provider



## Example Usage

```terraform
provider "grafana-adaptive-metrics" {
  url     = "https://my-prometheus-url.net"
  api_key = "my-tenant-id:my-api-key"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) Tenant ID and Access Policy Token (or API key) for Grafana Cloud in the format '<tenant-id>:<token-or-api-key>'. May alternatively be set via the `GRAFANA_AM_API_KEY` environment variable.
- `debug` (Boolean) Whether to enable debug logging. Defaults to false.
- `http_headers` (Map of String, Sensitive) HTTP headers mapping keys to values used for accessing Grafana Cloud APIs. May alternatively be set via the `GRAFANA_AM_HTTP_HEADERS` environment variable in JSON format.
- `retries` (Number) The amount of retries to use for Grafana API and Grafana Cloud API calls. Defaults to 3. May alternatively be set via the `GRAFANA_AM_RETRIES` environment variable.
- `url` (String) Grafana Cloud's API URL. May alternatively be set via the `GRAFANA_AM_API_URL` environment variable.
