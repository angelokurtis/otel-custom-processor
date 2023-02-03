package customprocessor

import (
	"context"
	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"log"
)

const (
	// The value of "type" key in configuration.
	typeStr = "custom"
	// The stability level of the processor.
	stability = component.StabilityLevelDevelopment
)

func NewFactory() processor.Factory {
	if err := view.Register(NewMetricViews()...); err != nil {
		log.Fatal(err)
	}
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
	tp := NewTracesProcessor(set.Logger)
	return processorhelper.NewTracesProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		tp.Count,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}))
}
