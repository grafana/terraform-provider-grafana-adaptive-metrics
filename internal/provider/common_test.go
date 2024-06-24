package provider

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
)

func CheckAccTestsEnabled(t *testing.T) {
	t.Helper()

	if enabled, _ := strconv.ParseBool(os.Getenv(resource.EnvTfAcc)); enabled {
		for _, env := range []string{"GRAFANA_AM_API_URL", "GRAFANA_AM_API_KEY"} {
			if _, ok := os.LookupEnv(env); !ok {
				t.Fatalf("Missing required env var: %s", env)
			}
		}
		return
	}

	t.Skip("Set TF_ACC=true to enable acceptance tests.")
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ClientForAccTest(t *testing.T) *client.Client {
	t.Helper()

	apiURL := os.Getenv("GRAFANA_AM_API_URL")
	apiKey := os.Getenv("GRAFANA_AM_API_KEY")

	c, err := client.New(apiURL, &client.Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	return c
}

func AggregationRulesForAccTest(t *testing.T) *AggregationRules {
	t.Helper()

	c := ClientForAccTest(t)

	aggRules := NewAggregationRules(c)
	require.NoError(t, aggRules.Init())

	return aggRules
}
