package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

var _ validator.String = gatePolicyValidator{}

type gatePolicyValidator struct{}

func (gatePolicyValidator) Description(context.Context) string {
	return fmt.Sprintf("value must be %q or %q", model.GatePolicyUnbounded, model.GatePolicyNoIncrease)
}

func (v gatePolicyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (gatePolicyValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == model.GatePolicyUnbounded || value == model.GatePolicyNoIncrease {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid gate policy",
		fmt.Sprintf("Expected %q or %q, got %q.", model.GatePolicyUnbounded, model.GatePolicyNoIncrease, value),
	)
}
