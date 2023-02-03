package customprocessor

import (
	"context"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type TracesProcessor struct {
	log *zap.Logger
}

func NewTracesProcessor(logger *zap.Logger) *TracesProcessor {
	return &TracesProcessor{log: logger}
}

func (tp *TracesProcessor) Count(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		rs := traces.ResourceSpans().At(i)
		attrs := rs.Resource().Attributes()

		tags := make([]tag.Mutator, 0, 0)
		if val, ok := attrs.Get("service.name"); ok {
			tags = append(tags, tag.Upsert(serviceOwnerKey, val.Str()))
		} else {
			tags = append(tags, tag.Upsert(serviceOwnerKey, "other-services"))
		}
		if val, ok := attrs.Get("telemetry.sdk.language"); ok {
			tags = append(tags, tag.Upsert(telemetrySDKLanguageKey, val.Str()))
		}
		if val, ok := attrs.Get("telemetry.sdk.version"); ok {
			tags = append(tags, tag.Upsert(telemetrySDKVersionKey, val.Str()))
		}

		spanCount := 0
		ilss := rs.ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			spanCount += ilss.At(j).Spans().Len()
		}

		tp.log.Debug("running otelcustomprocessor developed by Kurtis", zap.Int("span_count", spanCount))
		if err := stats.RecordWithTags(ctx, tags, spansProcessedCount.M(int64(spanCount))); err != nil {
			tp.log.Warn("Error registering metrics", zap.Any("error", err))
		}
	}
	return traces, nil
}
