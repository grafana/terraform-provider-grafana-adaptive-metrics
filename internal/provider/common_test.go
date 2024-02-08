package provider

import (
	"os"
	"strconv"
	"testing"
)

func CheckAccTestsEnabled(t *testing.T) {
	t.Helper()

	if enabled, _ := strconv.ParseBool(os.Getenv("TF_ACC")); enabled {
		for _, env := range []string{"GRAFANA_CLOUD_API_URL", "GRAFANA_CLOUD_API_KEY"} {
			if _, ok := os.LookupEnv(env); !ok {
				t.Fatalf("Missing required env var: %s", env)
			}
		}
		return
	}

	t.Skip("Set TF_ACC=true to enable acceptance tests.")
}
