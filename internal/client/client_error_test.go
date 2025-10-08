package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractErrorMessage(t *testing.T) {
	client := &Client{}

	t.Run("JSON error response", func(t *testing.T) {
		body := []byte(`{"error": "error updating policy with ID \"01K6AY1A5X6C40K6PJE2YNB6FN\""}`)
		errorMsg := client.extractErrorMessage(body, 500)
		require.Equal(t, "API error (status 500): error updating policy with ID \"01K6AY1A5X6C40K6PJE2YNB6FN\"", errorMsg)
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		body := []byte(`invalid json`)
		errorMsg := client.extractErrorMessage(body, 400)
		require.Equal(t, "status: 400, body: invalid json", errorMsg)
	})

	t.Run("Empty error field", func(t *testing.T) {
		body := []byte(`{"error": ""}`)
		errorMsg := client.extractErrorMessage(body, 500)
		require.Equal(t, "status: 500, body: {\"error\": \"\"}", errorMsg)
	})

	t.Run("Missing error field", func(t *testing.T) {
		body := []byte(`{"message": "some other error"}`)
		errorMsg := client.extractErrorMessage(body, 500)
		require.Equal(t, "status: 500, body: {\"message\": \"some other error\"}", errorMsg)
	})
}
