package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// We expect the GRAFANA_AM_API_URL and GRAFANA_AM_API_URL env vars to be populated.
	providerConfig = `provider "grafana-adaptive-metrics" {}`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"grafana-adaptive-metrics": providerserver.NewProtocol6WithError(New("test", "unknown")()),
	}
)
