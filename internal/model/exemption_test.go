package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExemptionMatchTypeRoundTrip(t *testing.T) {
	for _, matchType := range []string{"exact", "prefix", "suffix"} {
		t.Run(matchType, func(t *testing.T) {
			input := Exemption{ID: "existing-exemption", Metric: "rpc_", MatchType: matchType, KeepLabels: []string{"method", "status"}}
			wire, err := json.Marshal(input)
			require.NoError(t, err)
			var decoded Exemption
			require.NoError(t, json.Unmarshal(wire, &decoded))
			state := decoded.ToTF()
			require.Equal(t, matchType, state.MatchType.ValueString())
			restored := state.ToAPIReq()
			require.Equal(t, input.ID, restored.ID)
			require.Equal(t, input.Metric, restored.Metric)
			require.Equal(t, input.MatchType, restored.MatchType)
			require.Equal(t, input.KeepLabels, restored.KeepLabels)
		})
	}
}

func TestLegacyExemptionDefaultsToExact(t *testing.T) {
	var input Exemption
	require.NoError(t, json.Unmarshal([]byte(`{"id":"legacy","metric":"rpc_requests_total"}`), &input))
	require.Equal(t, "exact", input.ToTF().MatchType.ValueString())
	require.Equal(t, "exact", input.ToTF().ToAPIReq().MatchType)
}
