package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/model"
)

var (
	// minifiedJson is the json equivalent for rulesPayload and recsPayload below.
	minifiedJson = []byte(`[{"metric":"kube_persistentvolumeclaim_created","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"]},{"metric":"kube_persistentvolumeclaim_resource_requests_storage_bytes","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"]}]`)

	rulesPayload = []model.AggregationRule{
		{
			Metric:       "kube_persistentvolumeclaim_created",
			DropLabels:   []string{"persistentvolumeclaim"},
			Aggregations: []string{"count", "sum"},
		},
		{
			Metric:       "kube_persistentvolumeclaim_resource_requests_storage_bytes",
			DropLabels:   []string{"persistentvolumeclaim"},
			Aggregations: []string{"count", "sum"},
		},
	}
	recsPayload = []model.AggregationRecommendation{
		{
			AggregationRule: model.AggregationRule{
				Metric:       "kube_persistentvolumeclaim_created",
				DropLabels:   []string{"persistentvolumeclaim"},
				Aggregations: []string{"count", "sum"},
			},
		},
		{
			AggregationRule: model.AggregationRule{
				Metric:       "kube_persistentvolumeclaim_resource_requests_storage_bytes",
				DropLabels:   []string{"persistentvolumeclaim"},
				Aggregations: []string{"count", "sum"},
			},
		},
	}
)

func TestClientAuths(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	apiHeader := make(http.Header)
	apiHeader.Add("Authorization", "Bearer apikey")

	s.addExpected(
		"GET", "/aggregations/recommendations/config", apiHeader, nil,
		nil, []byte(`{}`),
	)

	scopeHeader := make(http.Header)
	scopeHeader.Add("x-scope-orgid", "9960")

	s.addExpected(
		"GET", "/aggregations/recommendations/config", scopeHeader, nil,
		nil, []byte(`{}`),
	)

	cAPI, err := New(s.server.URL, &Config{APIKey: "apikey"})
	require.NoError(t, err)

	_, err = cAPI.AggregationRecommendationsConfig()
	require.NoError(t, err)

	cScope, err := New(s.server.URL, &Config{HTTPHeaders: map[string]string{"x-scope-orgid": "9960"}})
	require.NoError(t, err)

	_, err = cScope.AggregationRecommendationsConfig()
	require.NoError(t, err)
}

func TestAggregationRecommendations(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected(
		"GET", "/aggregations/recommendations", nil, nil,
		nil, minifiedJson,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendations()
	require.NoError(t, err)

	require.Equal(t, recsPayload, actual)
}

func TestUpdateAggregationRecommendationsConfig(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected(
		"POST", "/aggregations/recommendations/config", nil, []byte(`{"keep_labels":["namespace"]}`),
		nil, nil,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	require.NoError(t, c.UpdateAggregationRecommendationsConfig(model.AggregationRecommendationConfiguration{
		KeepLabels: []string{"namespace"},
	}))
}

func TestAggregationRecommendationsConfig(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected(
		"GET", "/aggregations/recommendations/config", nil, nil,
		nil, []byte(`{"keep_labels":["namespace"]}`),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendationsConfig()
	require.NoError(t, err)

	require.Equal(t, model.AggregationRecommendationConfiguration{
		KeepLabels: []string{"namespace"},
	}, actual)
}

func TestAggregationRules(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	header := make(http.Header)
	header.Set("Etag", etag)

	s.addExpected(
		"GET", "/aggregations/rules", nil, nil,
		header, minifiedJson,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actualRules, actualEtag, err := c.AggregationRules()
	require.NoError(t, err)

	require.Equal(t, etag, actualEtag)

	require.Equal(t, rulesPayload, actualRules)
}

func TestUpdateAggregationRules(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected(
		"POST", "/aggregations/rules", expectedHeader, minifiedJson,
		respHeader, nil,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.UpdateAggregationRules(rulesPayload, etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestCreateAggregationRule(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected(
		"POST", "/aggregations/rule/test_metric", expectedHeader, []byte(`{"metric":"test_metric","drop":true}`),
		respHeader, nil,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.CreateAggregationRule(model.AggregationRule{Metric: "test_metric", Drop: true}, etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestReadAggregationRule(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	respHeader := make(http.Header)
	respHeader.Set("ETag", etag)

	s.addExpected(
		"GET", "/aggregations/rule/test_metric", nil, nil,
		respHeader, []byte(`{"metric":"test_metric","drop":true}`),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, newEtag, err := c.ReadAggregationRule("test_metric")
	require.NoError(t, err)

	require.Equal(t, etag, newEtag)
	require.Equal(t, model.AggregationRule{Metric: "test_metric", Drop: true}, actual)
}

func TestUpdateAggregationRule(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected(
		"PUT", "/aggregations/rule/test_metric", expectedHeader, []byte(`{"metric":"test_metric","drop":true}`),
		respHeader, nil,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.UpdateAggregationRule(model.AggregationRule{Metric: "test_metric", Drop: true}, etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestDeleteAggregationRule(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected(
		"DELETE", "/aggregations/rule/test_metric", expectedHeader, nil,
		respHeader, nil,
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.DeleteAggregationRule("test_metric", etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}
