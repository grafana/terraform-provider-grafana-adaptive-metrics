package provider

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func CheckAccTestsEnabled(t *testing.T) {
	t.Helper()

	if enabled, _ := strconv.ParseBool(os.Getenv("TF_ACC")); enabled {
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
