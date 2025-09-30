package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPolicyToTF(t *testing.T) {
	usageSources := []UsageSource{UsageSourceDashboard, UsageSourceRules}
	unusedMetricsAction := UnusedMetricsActionDropCustomLabels
	customLabels := []string{"cluster", "namespace"}
	minQueryUsages := int32(5)

	policy := Policy{
		ID:                                   "01JQVN6036Z18P6Z958JNNTXRP",
		Name:                                 "test-policy",
		UsageSources:                         &usageSources,
		UnusedMetricsAction:                  &unusedMetricsAction,
		UnusedMetricsActionActOnCustomLabels: &customLabels,
		MinQueryUsages:                       &minQueryUsages,
	}

	tf := policy.ToTF()

	require.Equal(t, "01JQVN6036Z18P6Z958JNNTXRP", tf.ID.ValueString())
	require.Equal(t, "test-policy", tf.Name.ValueString())
	require.Equal(t, 2, len(tf.UsageSources.Elements()))
	require.Equal(t, "drop_custom_labels", tf.UnusedMetricsAction.ValueString())
	require.Equal(t, 2, len(tf.UnusedMetricsActionActOnCustomLabels.Elements()))
	require.Equal(t, int64(5), tf.MinQueryUsages.ValueInt64())
}

func TestPolicyTFToAPIReq(t *testing.T) {
	usageSources, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("dashboard"),
		types.StringValue("rules"),
	})
	unusedMetricsAction := types.StringValue("drop_custom_labels")
	customLabels, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("cluster"),
		types.StringValue("namespace"),
	})
	minQueryUsages := types.Int64Value(5)

	policyTF := PolicyTF{
		ID:                                   types.StringValue("01JQVN6036Z18P6Z958JNNTXRP"),
		Name:                                 types.StringValue("test-policy"),
		UsageSources:                         usageSources,
		UnusedMetricsAction:                  unusedMetricsAction,
		UnusedMetricsActionActOnCustomLabels: customLabels,
		MinQueryUsages:                       minQueryUsages,
	}

	apiReq := policyTF.ToAPIReq()

	require.Equal(t, "01JQVN6036Z18P6Z958JNNTXRP", apiReq.ID)
	require.Equal(t, "test-policy", apiReq.Name)
	require.NotNil(t, apiReq.UsageSources)
	require.Equal(t, 2, len(*apiReq.UsageSources))
	require.Equal(t, UsageSourceDashboard, (*apiReq.UsageSources)[0])
	require.Equal(t, UsageSourceRules, (*apiReq.UsageSources)[1])
	require.NotNil(t, apiReq.UnusedMetricsAction)
	require.Equal(t, UnusedMetricsActionDropCustomLabels, *apiReq.UnusedMetricsAction)
	require.NotNil(t, apiReq.UnusedMetricsActionActOnCustomLabels)
	require.Equal(t, 2, len(*apiReq.UnusedMetricsActionActOnCustomLabels))
	require.Equal(t, "cluster", (*apiReq.UnusedMetricsActionActOnCustomLabels)[0])
	require.Equal(t, "namespace", (*apiReq.UnusedMetricsActionActOnCustomLabels)[1])
	require.NotNil(t, apiReq.MinQueryUsages)
	require.Equal(t, int32(5), *apiReq.MinQueryUsages)
}
