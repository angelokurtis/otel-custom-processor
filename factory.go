package customprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/zap"
)

const (
	// The value of "type" key in configuration.
	typeStr = "custom"
	// The stability level of the processor.
	stability = component.StabilityLevelDevelopment
)

var processorCapabilities = consumer.Capabilities{MutatesData: false}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesProcessor(ctx context.Context, set processor.CreateSettings, cfg component.Config, nextConsumer consumer.Traces) (processor.Traces, error) {
	return processorhelper.NewTracesProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		func(ctx context.Context, traces ptrace.Traces) (ptrace.Traces, error) {
			set.Logger.Warn("running otelcustomprocessor developed by Kurtis", zap.Int("span_count", traces.SpanCount()))
			return traces, nil
		},
		processorhelper.WithCapabilities(processorCapabilities))
}
