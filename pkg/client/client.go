package client

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Client is the implementation of tracing client for Jaeger.
type Client struct {
	Provider     *sdktrace.TracerProvider
	Exporter     *jaeger.Exporter
	Logger       *logrus.Logger
	CollectorURL string
	ServiceName  string
	Environment  string
}

// Opts is the configuration for Client.
type Opts struct {
	Logger       *logrus.Logger
	CollectorURL string
	ServiceName  string
	Environment  string
}

// New creates a new tracing client for Jaeger with the given configuration.
func New(opts *Opts) (*Client, error) {
	client := &Client{
		CollectorURL: opts.CollectorURL,
		ServiceName:  opts.ServiceName,
		Environment:  opts.Environment,
		Logger:       opts.Logger,
	}

	return client.Connect()
}

// Connect connects to the tracing server and returns a shutdown function and error.
func (jc *Client) Connect() (*Client, error) {
	var err error

	jc.Exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jc.CollectorURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	jc.Provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jc.Exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(jc.ServiceName),
			attribute.String("environment", jc.Environment),
		)),
	)

	otel.SetTracerProvider(jc.Provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	jc.Logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))

	return jc, nil
}

// Disconnect disconnects from the tracing server and returns an error.
func (jc *Client) Disconnect(ctx context.Context) error {
	otel.SetTracerProvider(nil)

	err := jc.Provider.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	err = jc.Exporter.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown exporter: %w", err)
	}

	return nil
}
