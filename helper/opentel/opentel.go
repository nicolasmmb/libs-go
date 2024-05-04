package opentel

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel/trace"
)

var (
	tracer *trace.Tracer
	once   sync.Once
)

type TraceConfig struct {
	TracerName            string
	Endpoint              string
	SchemaURL             string
	ServiceName           string
	ServiceVersion        string
	ServiceNamespace      string
	DeploymentEnvironment string
}

func (cfg *TraceConfig) Validate() error {
	if cfg.TracerName == "" {
		return errors.New("--> TracerName cannot be empty")
	}

	if cfg.Endpoint == "" {
		return errors.New("--> Endpoint cannot be empty")
	}

	if cfg.SchemaURL == "" {
		return errors.New("--> SchemaURL cannot be empty")
	}

	if cfg.ServiceName == "" {
		return errors.New("--> ServiceName cannot be empty")
	}

	if cfg.ServiceVersion == "" {
		return errors.New("--> ServiceVersion cannot be empty")
	}

	if cfg.ServiceNamespace == "" {
		return errors.New("--> ServiceNamespace cannot be empty")
	}

	if cfg.DeploymentEnvironment == "" {
		return errors.New("--> DeploymentEnvironment cannot be empty")
	}

	return nil
}

func NewTraceConfig(tracerName, endpoint, schemaURL, serviceName, serviceVersion, serviceNamespace, deploymentEnvironment string) TraceConfig {
	cfg := TraceConfig{
		TracerName:            tracerName,
		Endpoint:              endpoint,
		SchemaURL:             schemaURL,
		ServiceName:           serviceName,
		ServiceVersion:        serviceVersion,
		ServiceNamespace:      serviceNamespace,
		DeploymentEnvironment: deploymentEnvironment,
	}
	if err := cfg.Validate(); err != nil {
		log.Panicln(err)
	}
	return cfg
}

func InitTracer(cfg TraceConfig) error {

	// Create a new OTLP trace exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(cfg.Endpoint),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		otlptracehttp.WithTimeout(30*time.Second),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 5 * time.Second,
			MaxInterval:     30 * time.Second,
			MaxElapsedTime:  5 * time.Minute,
		}),
	)

	// Handle any errors
	if err != nil {
		return err
	}

	// Close the exporter on shutdown
	defer func() {
		if err := exporter.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Create and configure a new resource
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			cfg.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.ServiceNamespace(cfg.SchemaURL),
			semconv.DeploymentEnvironment(cfg.DeploymentEnvironment),
		),
	)

	// Handle any errors
	if err != nil {
		return err
	}

	// Create a new OTLP trace provider
	tracerProvider := oteltrace.NewTracerProvider(
		oteltrace.WithSampler(oteltrace.AlwaysSample()),
		oteltrace.WithBatcher(exporter),
		oteltrace.WithResource(resource),
	)

	// Set the global trace provider
	once.Do(func() {
		otel.SetTracerProvider(tracerProvider)
		trace := tracerProvider.Tracer("")
		tracer = &trace
		log.Println("--> OpenTelemetry Tracer initialized")
	})

	return nil
}

func GetTracer() *trace.Tracer {
	if tracer != nil {
		return tracer
	} else {
		log.Println("OpenTelemetry Tracer not initialized")
		return nil
	}
}
