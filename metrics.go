package customprocessor

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/obsreport"
)

var (
	serviceOwnerKey, _         = tag.NewKey("service_owner")
	telemetrySDKLanguageKey, _ = tag.NewKey("telemetry_sdk_language")
	telemetrySDKVersionKey, _  = tag.NewKey("telemetry_sdk_version")
	spansProcessedCount        = stats.Int64("spans_processed_total", "Counts the processing of new spans", stats.UnitDimensionless)
)

func NewMetricViews() []*view.View {
	return []*view.View{
		{
			Name:        obsreport.BuildProcessorCustomMetricName(typeStr, spansProcessedCount.Name()),
			Measure:     spansProcessedCount,
			Description: spansProcessedCount.Description(),
			TagKeys:     []tag.Key{serviceOwnerKey, telemetrySDKLanguageKey, telemetrySDKVersionKey},
			Aggregation: view.Sum(),
		},
	}
}
