package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
	"github.com/stretchr/testify/require"
)

func TestExemptionImportKeepsSegmentAndMatchType(t *testing.T) {
	for _, segment := range []string{"", "segment-ethereum"} {
		t.Run("segment="+segment, func(t *testing.T) {
			reads := 0
			updates := 0
			stored := model.Exemption{ID: "exemption-1", Metric: "rpc_", MatchType: "prefix", KeepLabels: []string{"method", "status"}}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/v1/recommendations/exemptions/exemption-1", r.URL.Path)
				require.Equal(t, segment, r.URL.Query().Get("segment"))
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					reads++
					require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"result": stored}))
				case http.MethodPut:
					updates++
					require.NoError(t, json.NewDecoder(r.Body).Decode(&stored))
					require.Equal(t, "exemption-1", stored.ID)
					require.Equal(t, "prefix", stored.MatchType)
					require.Equal(t, []string{"method", "status"}, stored.KeepLabels)
					require.Equal(t, "updated reason", stored.Reason)
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			defer server.Close()
			api, err := client.New(server.URL, &client.Config{})
			require.NoError(t, err)
			r := &exemptionResource{client: api}
			var schema resource.SchemaResponse
			r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			require.False(t, state.Set(t.Context(), &model.ExemptionTF{}).HasError())
			id := "exemption-1"
			if segment != "" {
				id = segment + "/" + id
			}
			imported := resource.ImportStateResponse{State: state}
			r.ImportState(t.Context(), resource.ImportStateRequest{ID: id}, &imported)
			require.False(t, imported.Diagnostics.HasError(), imported.Diagnostics)
			read := resource.ReadResponse{State: imported.State}
			r.Read(t.Context(), resource.ReadRequest{State: imported.State}, &read)
			require.False(t, read.Diagnostics.HasError(), read.Diagnostics)
			var actual model.ExemptionTF
			require.False(t, read.State.Get(t.Context(), &actual).HasError())
			require.Equal(t, "exemption-1", actual.ID.ValueString())
			require.Equal(t, segment, actual.Segment.ValueString())
			require.Equal(t, "prefix", actual.MatchType.ValueString())
			require.Equal(t, []types.String{types.StringValue("method"), types.StringValue("status")}, actual.KeepLabels)
			require.Equal(t, 1, reads)
			actual.Reason = types.StringValue("updated reason")
			planState := tfsdk.State{Schema: schema.Schema}
			require.False(t, planState.Set(t.Context(), &actual).HasError())
			update := resource.UpdateResponse{State: read.State}
			r.Update(t.Context(), resource.UpdateRequest{
				State: read.State,
				Plan:  tfsdk.Plan{Schema: schema.Schema, Raw: planState.Raw},
			}, &update)
			require.False(t, update.Diagnostics.HasError(), update.Diagnostics)
			require.False(t, update.State.Get(t.Context(), &actual).HasError())
			require.Equal(t, "prefix", actual.MatchType.ValueString())
			require.Equal(t, segment, actual.Segment.ValueString())
			require.Equal(t, "updated reason", actual.Reason.ValueString())
			require.Equal(t, 1, updates)
			require.Equal(t, 2, reads)
		})
	}
}

func TestExemptionImportRejectsMalformedIdentity(t *testing.T) {
	for _, id := range []string{"", "/id", "segment/", "segment/id/extra"} {
		t.Run(id, func(t *testing.T) {
			r := &exemptionResource{}
			var schema resource.SchemaResponse
			r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			require.False(t, state.Set(t.Context(), &model.ExemptionTF{}).HasError())
			imported := resource.ImportStateResponse{State: state}
			r.ImportState(t.Context(), resource.ImportStateRequest{ID: id}, &imported)
			require.True(t, imported.Diagnostics.HasError())
			var actual types.String
			require.False(t, imported.State.GetAttribute(t.Context(), path.Root("id"), &actual).HasError())
			require.True(t, actual.IsNull())
		})
	}
}

func TestExemptionMatchTypeValidation(t *testing.T) {
	r := &exemptionResource{}
	var schema resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
	matchType := schema.Schema.Attributes["match_type"].(resourceschema.StringAttribute)
	require.True(t, matchType.Optional && matchType.Computed)
	require.Len(t, matchType.Validators, 1)
	for _, value := range []types.String{types.StringValue("exact"), types.StringValue("prefix"), types.StringValue("suffix"), types.StringNull(), types.StringUnknown()} {
		var response validator.StringResponse
		matchType.Validators[0].ValidateString(t.Context(), validator.StringRequest{ConfigValue: value}, &response)
		require.False(t, response.Diagnostics.HasError())
	}
	for _, value := range []string{"", "regex", "PREFIX"} {
		var response validator.StringResponse
		matchType.Validators[0].ValidateString(t.Context(), validator.StringRequest{ConfigValue: types.StringValue(value)}, &response)
		require.True(t, response.Diagnostics.HasError())
	}
}
