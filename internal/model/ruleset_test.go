package model

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlignUpstreamWithState(t *testing.T) {
	t.Run("it re-orders exact matches", func(t *testing.T) {

		state := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
			{Metric: "d", MatchType: "exact"},
		}

		config := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "d", MatchType: "exact"},
		}

		output := AlignUpstreamWithState(state, config)
		requireConsistentState(t, state, config, output)
	})

	t.Run("it keeps new items at the end", func(t *testing.T) {
		state := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
			{Metric: "d", MatchType: "exact"},
		}

		config := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "d", MatchType: "exact"},
			{Metric: "e", MatchType: "exact"},
			{Metric: "f", MatchType: "exact"},
		}

		output := AlignUpstreamWithState(state, config)
		requireConsistentState(t, state, config, output)
	})

	t.Run("it skips missing items", func(t *testing.T) {
		state := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
			{Metric: "d", MatchType: "exact"},
		}

		config := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "b", MatchType: "exact"},
		}

		output := AlignUpstreamWithState(state, config)
		requireConsistentState(t, state, config, output)
	})

	t.Run("it preserves order for non-exact matches", func(t *testing.T) {
		state := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "b", MatchType: "prefix"},
			{Metric: "d", MatchType: "prefix"},
		}

		config := AggregationRuleSet{
			{Metric: "c", MatchType: "exact"},
			{Metric: "a", MatchType: "exact"},
			{Metric: "b", MatchType: "prefix"},
			{Metric: "d", MatchType: "prefix"},
		}

		output := AlignUpstreamWithState(state, config)
		requireConsistentState(t, state, config, output)
	})

	t.Run("it gives up if non-exact matches aren't ordered the same", func(t *testing.T) {
		state := AggregationRuleSet{
			{Metric: "a", MatchType: "exact"},
			{Metric: "c", MatchType: "exact"},
			{Metric: "b", MatchType: "prefix"},
			{Metric: "d", MatchType: "prefix"},
		}

		config := AggregationRuleSet{
			{Metric: "c", MatchType: "exact"},
			{Metric: "a", MatchType: "exact"},
			{Metric: "d", MatchType: "prefix"},
			{Metric: "b", MatchType: "prefix"},
		}

		output := AlignUpstreamWithState(state, config)
		requireConsistentState(t, state, config, output)
	})
}

func requireConsistentState(t *testing.T, state, config, output AggregationRuleSet) {
	nonExactMatchesInState := filteredRuleNames(state, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})
	nonExactMatchesInConfig := filteredRuleNames(config, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})
	nonExactMatchesInOutput := filteredRuleNames(output, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})

	// Only if the state and config have matching orders of non-exact matches
	// can we assert that the output matches the state.
	if reflect.DeepEqual(nonExactMatchesInState, nonExactMatchesInConfig) {
		require.Equal(t, nonExactMatchesInState, nonExactMatchesInOutput, "order of non-exact rules should be preserved")
		requireStateOrderPreserved(t, state, config, output)
		requireConfigOrderPreserved(t, config, output)
	} else {
		require.Equal(t, config, output, "output should match config")
	}
}

func requireStateOrderPreserved(t *testing.T, state, config, output AggregationRuleSet) {
	exactMatchesInConfig := filteredRuleNames(config, func(r AggregationRule) bool {
		return r.IsExactMatch()
	})
	exactMatchesInState := filteredRuleNames(state, func(r AggregationRule) bool {
		return r.IsExactMatch() && slices.Contains(exactMatchesInConfig, r.Metric)
	})
	exactMatchesInOutput := filteredRuleNames(output, func(r AggregationRule) bool {
		return r.IsExactMatch() && slices.Contains(exactMatchesInState, r.Metric)
	})
	require.Equal(t, exactMatchesInState, exactMatchesInOutput, "order of exact rules should be preserved")

	nonExactMatchesInState := filteredRuleNames(state, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})
	nonExactMatchesInConfig := filteredRuleNames(output, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})
	nonExactMatchesInOutput := filteredRuleNames(output, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})

	// Only if the state and config have matching orders of non-exact matches
	// can we assert that the output matches the state.
	if reflect.DeepEqual(nonExactMatchesInState, nonExactMatchesInConfig) {
		require.Equal(t, nonExactMatchesInState, nonExactMatchesInOutput, "order of non-exact rules should be preserved")
	}
}

func requireConfigOrderPreserved(t *testing.T, config, output AggregationRuleSet) {
	require.ElementsMatch(t, config, output, "output should match config")

	// Collect all non-exact metrics in the config and check that they are in
	// the same order in the output. If this property is violated, then the
	// semantics of the ruleset are changed.
	nonExactMatchesInConfig := filteredRuleNames(output, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})
	nonExactMatchesInOutput := filteredRuleNames(output, func(r AggregationRule) bool {
		return !r.IsExactMatch()
	})

	require.Equal(t, nonExactMatchesInConfig, nonExactMatchesInOutput, "order of non-exact rules should be preserved")
}

func filteredRuleNames(rules AggregationRuleSet, f func(r AggregationRule) bool) []string {
	output := []string{}
	for _, rule := range rules {
		if f(rule) {
			output = append(output, rule.Metric)
		}
	}
	return output
}
