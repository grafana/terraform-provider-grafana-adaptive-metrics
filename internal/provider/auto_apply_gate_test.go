package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestAutoApply_DefaultDisabled(t *testing.T) {
	resources := map[string]resource.Resource{
		"segment":                newSegmentResource(),
		"recommendations config": newRecommendationsConfigResource(),
	}

	for name, resourceUnderTest := range resources {
		t.Run(name, func(t *testing.T) {
			var schemaResponse resource.SchemaResponse
			resourceUnderTest.Schema(t.Context(), resource.SchemaRequest{}, &schemaResponse)
			require.False(t, schemaResponse.Diagnostics.HasError())

			autoApply, ok := schemaResponse.Schema.Attributes["auto_apply"].(resourceschema.SingleNestedAttribute)
			require.True(t, ok)
			require.NotNil(t, autoApply.Default)
			require.True(t, autoApply.Computed)
			enabled, ok := autoApply.Attributes["enabled"].(resourceschema.BoolAttribute)
			require.True(t, ok)
			require.True(t, enabled.Computed)

			var response defaults.ObjectResponse
			autoApply.Default.DefaultObject(t.Context(), defaults.ObjectRequest{}, &response)
			require.False(t, response.Diagnostics.HasError())
			require.Equal(t, types.BoolValue(false), response.PlanValue.Attributes()["enabled"])
			require.True(t, response.PlanValue.Attributes()["gate"].IsNull())
		})
	}
}

func TestAutoApplyGate_Validation(t *testing.T) {
	resources := map[string]resource.Resource{
		"segment":                newSegmentResource(),
		"recommendations config": newRecommendationsConfigResource(),
	}

	for name, resourceUnderTest := range resources {
		t.Run(name, func(t *testing.T) {
			var schemaResponse resource.SchemaResponse
			resourceUnderTest.Schema(t.Context(), resource.SchemaRequest{}, &schemaResponse)
			require.False(t, schemaResponse.Diagnostics.HasError())

			autoApply, ok := schemaResponse.Schema.Attributes["auto_apply"].(resourceschema.SingleNestedAttribute)
			require.True(t, ok)
			gate, ok := autoApply.Attributes["gate"].(resourceschema.SingleNestedAttribute)
			require.True(t, ok)
			policy, ok := gate.Attributes["policy"].(resourceschema.StringAttribute)
			require.True(t, ok)
			require.Len(t, policy.Validators, 1)

			for _, value := range []string{"unbounded", "no-increase"} {
				var response validator.StringResponse
				policy.Validators[0].ValidateString(t.Context(), validator.StringRequest{ConfigValue: types.StringValue(value)}, &response)
				require.False(t, response.Diagnostics.HasError())
			}

			var response validator.StringResponse
			policy.Validators[0].ValidateString(t.Context(), validator.StringRequest{ConfigValue: types.StringValue("invalid")}, &response)
			require.True(t, response.Diagnostics.HasError())
		})
	}
}
