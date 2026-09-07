package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

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

			autoApply := schemaResponse.Schema.Attributes["auto_apply"].(resourceschema.SingleNestedAttribute)
			gate := autoApply.Attributes["gate"].(resourceschema.SingleNestedAttribute)
			policy := gate.Attributes["policy"].(resourceschema.StringAttribute)
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
