package client

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

var (
	// minifiedJson is the json equivalent for rulesPayload and recsPayload below.
	minifiedJson        = []byte(`[{"metric":"kube_persistentvolumeclaim_created","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"]},{"metric":"kube_persistentvolumeclaim_resource_requests_storage_bytes","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"]}]`)
	minifiedVerboseJson = []byte(`[{"metric":"kube_persistentvolumeclaim_created","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"],"recommended_action":"keep"},{"metric":"kube_persistentvolumeclaim_resource_requests_storage_bytes","drop_labels":["persistentvolumeclaim"],"aggregations":["count","sum"],"recommended_action":"update"}]`)

	segmentedMinified = []byte(`[{"segment":{"name":"segmentname","id":"segmentid"},"rules":` + string(minifiedJson) + `}]`)

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
	verboseRecsPayload = []model.AggregationRecommendation{
		{
			AggregationRule: model.AggregationRule{
				Metric:       "kube_persistentvolumeclaim_created",
				DropLabels:   []string{"persistentvolumeclaim"},
				Aggregations: []string{"count", "sum"},
			},
			RecommendedAction: "keep",
		},
		{
			AggregationRule: model.AggregationRule{
				Metric:       "kube_persistentvolumeclaim_resource_requests_storage_bytes",
				DropLabels:   []string{"persistentvolumeclaim"},
				Aggregations: []string{"count", "sum"},
			},
			RecommendedAction: "update",
		},
	}
)

func TestClientAuths(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	apiHeader := make(http.Header)
	apiHeader.Add("Authorization", "Bearer apikey")

	s.addExpected("GET", "/aggregations/recommendations/config",
		withReqHeader(apiHeader),
		withRespBody([]byte(`{}`)),
	)

	scopeHeader := make(http.Header)
	scopeHeader.Add("x-scope-orgid", "9960")

	s.addExpected("GET", "/aggregations/recommendations/config",
		withReqHeader(scopeHeader),
		withRespBody([]byte(`{}`)),
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

	s.addExpected("GET", "/aggregations/recommendations",
		withRespBody(minifiedJson),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendations("", false, nil)
	require.NoError(t, err)

	require.Equal(t, recsPayload, actual)
}

func TestAggregationSegmentedRecommendations(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("GET", "/aggregations/recommendations",
		withRespBody(minifiedVerboseJson),
		withParams(url.Values{"verbose": []string{"true"}, "segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendations("segment-id", true, nil)
	require.NoError(t, err)

	require.Equal(t, verboseRecsPayload, actual)
}

func TestAggregationVerboseRecommendations(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("GET", "/aggregations/recommendations",
		withRespBody(minifiedVerboseJson),
		withParams(url.Values{"verbose": []string{"true"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendations("", true, nil)
	require.NoError(t, err)

	require.Equal(t, verboseRecsPayload, actual)
}

func TestAggregationRecommendationsWithAction(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("GET", "/aggregations/recommendations",
		withRespBody(minifiedJson),
		withParams(url.Values{"action": []string{"add", "update"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.AggregationRecommendations("", false, []string{"add", "update"})
	require.NoError(t, err)

	require.Equal(t, recsPayload, actual)
}

func TestUpdateAggregationRecommendationsConfig(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("POST", "/aggregations/recommendations/config",
		withReqBody([]byte(`{"keep_labels":["namespace"]}`)),
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

	s.addExpected("GET", "/aggregations/recommendations/config",
		withRespBody([]byte(`{"keep_labels":["namespace"]}`)),
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

	s.addExpected("GET", "/aggregations/segmented_rules",
		withRespBody(segmentedMinified),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actualRules, err := c.SegmentedAggregationRules()
	require.NoError(t, err)

	require.Equal(t, rulesPayload, actualRules[0].Rules)
}

func TestReadAggregationRuleSet(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	respHeader := make(http.Header)
	respHeader.Set("ETag", etag)

	s.addExpected("GET", "/aggregations/rules",
		withRespHeader(respHeader),
		withRespBody([]byte(`[{"metric":"test_metric","drop":true}]`)),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, newEtag, err := c.ReadAggregationRuleSet("segment-id")
	require.NoError(t, err)

	require.Equal(t, etag, newEtag)
	require.Equal(t, []model.AggregationRule{{Metric: "test_metric", Drop: true}}, actual)
}

func TestUpdateAggregationRuleSet(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected("POST", "/aggregations/rules",
		withReqHeader(expectedHeader),
		withReqBody([]byte(`[{"metric":"test_metric","drop":true}]`)),
		withRespHeader(respHeader),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.UpdateAggregationRuleSet("segment-id", []model.AggregationRule{{Metric: "test_metric", Drop: true}}, etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestUpdateAggregationRuleSetWithNilRules(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	expectedHeader := make(http.Header)
	expectedHeader.Set("If-Match", etag)

	respHeader := make(http.Header)
	respHeader.Set("ETag", "\"updated-fake-etag\"")

	s.addExpected("POST", "/aggregations/rules",
		withReqHeader(expectedHeader),
		withReqBody([]byte(`[]`)),
		withRespHeader(respHeader),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.UpdateAggregationRuleSet("segment-id", nil, etag)
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

	s.addExpected("POST", "/aggregations/rule/test_metric",
		withReqHeader(expectedHeader),
		withReqBody([]byte(`{"metric":"test_metric","drop":true}`)),
		withRespHeader(respHeader),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.CreateAggregationRule("segment-id", model.AggregationRule{Metric: "test_metric", Drop: true}, etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestReadAggregationRule(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	const etag = "\"fake-etag\""
	respHeader := make(http.Header)
	respHeader.Set("ETag", etag)

	s.addExpected("GET", "/aggregations/rule/test_metric",
		withRespHeader(respHeader),
		withRespBody([]byte(`{"metric":"test_metric","drop":true}`)),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, newEtag, err := c.ReadAggregationRule("segment-id", "test_metric")
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

	s.addExpected("PUT", "/aggregations/rule/test_metric",
		withReqHeader(expectedHeader),
		withReqBody([]byte(`{"metric":"test_metric","drop":true}`)),
		withRespHeader(respHeader),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.UpdateAggregationRule("segment-id", model.AggregationRule{Metric: "test_metric", Drop: true}, etag)
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

	s.addExpected("DELETE", "/aggregations/rule/test_metric",
		withReqHeader(expectedHeader),
		withRespHeader(respHeader),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	newEtag, err := c.DeleteAggregationRule("segment-id", "test_metric", etag)
	require.NoError(t, err)

	require.Equal(t, "\"updated-fake-etag\"", newEtag)
}

func TestCreateExemption(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	reqBody := []byte(`{"id":"","metric":"test_metric","keep_labels":["foobar"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`)
	respBody := []byte(`{"result":{"id":"generated-ulid","metric":"test_metric","keep_labels":["foobar"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`)

	s.addExpected("POST", "/v1/recommendations/exemptions",
		withReqBody(reqBody),
		withRespBody(respBody),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.CreateExemption("segment-id", model.Exemption{
		Metric:     "test_metric",
		KeepLabels: []string{"foobar"},
	})
	require.NoError(t, err)

	expected := model.Exemption{
		ID:         "generated-ulid",
		Metric:     "test_metric",
		KeepLabels: []string{"foobar"},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	require.Equal(t, expected, actual)
}

func TestReadExemption(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	respBody := []byte(`{"result":{"id":"generated-ulid","metric":"test_metric","keep_labels":["foobar"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`)
	expected := model.Exemption{
		ID:         "generated-ulid",
		Metric:     "test_metric",
		KeepLabels: []string{"foobar"},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	s.addExpected("GET", "/v1/recommendations/exemptions/generated-ulid",
		withRespBody(respBody),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.ReadExemption("segment-id", "generated-ulid")
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestUpdateExemption(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	reqBody := []byte(`{"id":"generated-ulid","metric":"test_metric","keep_labels":["foobar"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`)

	s.addExpected("PUT", "/v1/recommendations/exemptions/generated-ulid",
		withReqBody(reqBody),
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	err = c.UpdateExemption("segment-id", model.Exemption{
		ID:         "generated-ulid",
		Metric:     "test_metric",
		KeepLabels: []string{"foobar"},
	})
	require.NoError(t, err)
}

func TestDeleteExemption(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("DELETE", "/v1/recommendations/exemptions/generated-ulid",
		withParams(url.Values{"segment": []string{"segment-id"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	err = c.DeleteExemption("segment-id", "generated-ulid")
	require.NoError(t, err)
}

func TestCreateSegment(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	reqBody := []byte(`{"id":"","name":"segment name","selector":"{foo=\"bar\"}","fallback_to_default":true}`)
	respBody := []byte(`{"name":"segment name","selector":"{foo=\"bar\"}","fallback_to_default":true,"id":"generated-ulid"}`)

	s.addExpected("POST", "/aggregations/rules/segments",
		withReqBody(reqBody),
		withRespBody(respBody),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.CreateSegment(model.Segment{
		Name:              "segment name",
		Selector:          "{foo=\"bar\"}",
		FallbackToDefault: true,
	})
	require.NoError(t, err)

	expected := model.Segment{
		ID:                "generated-ulid",
		Name:              "segment name",
		Selector:          "{foo=\"bar\"}",
		FallbackToDefault: true,
	}

	require.Equal(t, expected, actual)
}

func TestReadSegment(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	respBody := []byte(`[{"name":"segment name","selector":"{foo=\"bar\"}","fallback_to_default":true,"id":"generated-ulid"}]`)
	expected := model.Segment{
		ID:                "generated-ulid",
		Name:              "segment name",
		Selector:          "{foo=\"bar\"}",
		FallbackToDefault: true,
	}

	s.addExpected("GET", "/aggregations/rules/segments",
		withRespBody(respBody),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	actual, err := c.ReadSegment("generated-ulid")
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestUpdateSegment(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	reqBody := []byte(`{"id":"generated-ulid","name":"segment name","selector":"{foo=\"bar\"}","fallback_to_default":true}`)

	s.addExpected("PUT", "/aggregations/rules/segments",
		withReqBody(reqBody),
		withParams(url.Values{"segment": []string{"generated-ulid"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	err = c.UpdateSegment(model.Segment{
		ID:                "generated-ulid",
		Name:              "segment name",
		Selector:          "{foo=\"bar\"}",
		FallbackToDefault: true,
	})
	require.NoError(t, err)
}

func TestDeleteSegment(t *testing.T) {
	s := newMockServer(t)
	defer s.close()

	s.addExpected("DELETE", "/aggregations/rules/segments",
		withParams(url.Values{"segment": []string{"generated-ulid"}}),
	)

	c, err := New(s.server.URL, &Config{})
	require.NoError(t, err)

	err = c.DeleteSegment("generated-ulid")
	require.NoError(t, err)
}
