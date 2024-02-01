package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AggregationRule struct {
	Metric    string `json:"metric"`
	MatchType string `json:"match_type,omitempty"`

	Drop       bool     `json:"drop,omitempty"`
	KeepLabels []string `json:"keep_labels,omitempty"`
	DropLabels []string `json:"drop_labels,omitempty"`

	Aggregations []string `json:"aggregations,omitempty"`

	AggregationInterval string `json:"aggregation_interval,omitempty"`
	AggregationDelay    string `json:"aggregation_delay,omitempty"`

	Ingest bool `json:"ingest,omitempty"`
}

func (r AggregationRule) ToTF() RuleTF {
	return RuleTF{
		Metric:    types.StringValue(r.Metric),
		MatchType: types.StringValue(r.MatchType),

		Drop:       types.BoolValue(r.Drop),
		KeepLabels: toTypesStringSlice(r.KeepLabels),
		DropLabels: toTypesStringSlice(r.DropLabels),

		Aggregations: toTypesStringSlice(r.Aggregations),

		AggregationInterval: types.StringValue(r.AggregationInterval),
		AggregationDelay:    types.StringValue(r.AggregationDelay),

		Ingest: types.BoolValue(r.Ingest),
	}
}

func toTypesStringSlice(in []string) []types.String {
	out := make([]types.String, len(in))
	for i, s := range in {
		out[i] = types.StringValue(s)
	}
	return out
}

func toStringSlice(in []types.String) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = s.ValueString()
	}
	return out
}

type RuleTF struct {
	Metric    types.String `tfsdk:"metric"`
	MatchType types.String `tfsdk:"match_type"`

	Drop       types.Bool     `tfsdk:"drop"`
	KeepLabels []types.String `tfsdk:"keep_labels"`
	DropLabels []types.String `tfsdk:"drop_labels"`

	Aggregations []types.String `tfsdk:"aggregations"`

	AggregationInterval types.String `tfsdk:"aggregation_interval"`
	AggregationDelay    types.String `tfsdk:"aggregation_delay"`

	Ingest types.Bool `tfsdk:"ingest"`

	LastUpdated types.String `tfsdk:"-"`
}

func (r RuleTF) ToAPIReq() AggregationRule {
	return AggregationRule{
		Metric:    r.Metric.ValueString(),
		MatchType: r.MatchType.ValueString(),

		Drop:       r.Drop.ValueBool(),
		KeepLabels: toStringSlice(r.KeepLabels),
		DropLabels: toStringSlice(r.DropLabels),

		Aggregations: toStringSlice(r.Aggregations),

		AggregationInterval: r.AggregationInterval.ValueString(),
		AggregationDelay:    r.AggregationDelay.ValueString(),

		Ingest: r.Ingest.ValueBool(),
	}
}
