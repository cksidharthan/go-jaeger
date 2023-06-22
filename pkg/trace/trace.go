package trace

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	Context context.Context
	Span    trace.Span
}

// NewTraceWithContext returns a new instance of a Trace object to be used for tracing.
func NewTraceWithContext(c context.Context) *Trace {
	span := trace.SpanFromContext(c)
	baseTrace := &Trace{
		Context: c,
		Span:    span,
	}

	return baseTrace
}

// GetTraceID returns the trace ID of the trace.
func (t *Trace) GetTraceID() string {
	return t.Span.SpanContext().TraceID().String()
}

// StartNewSpan starts a new span with the given span name.
func (t *Trace) StartNewSpan() *Trace {
	//nolint:dogsled // declaration has 3 blank identifiers
	curFunc, _, _, _ := runtime.Caller(1)
	funcPath := runtime.FuncForPC(curFunc).Name()
	funcName := funcPath[strings.LastIndex(funcPath, "/")+1:]

	return t.StartNewSpanf(funcName)
}

func (t *Trace) StartNewSpanWithName(spanName string) *Trace {
	return t.StartNewSpanf(spanName)
}

// StartNewSpanf starts a new span with the given span name and arguments.
func (t *Trace) StartNewSpanf(format string, args ...any) *Trace {
	// Use the global trace provider to instantiate a span anywhere
	tr := otel.GetTracerProvider().Tracer("")
	childContext, childSpan := tr.Start(t.Context, fmt.Sprintf(format, args...))

	newBaseTrace := &Trace{
		Context: childContext,
		Span:    childSpan,
	}

	return newBaseTrace
}

// SetHTTPHeaders sets the HTTP headers for the trace.
func (t *Trace) SetHTTPHeaders(request *http.Request) {
	spanContext := t.Span.SpanContext()

	request.Header.Set("X-B3-Sampled", "1")
	request.Header.Set("X-B3-Spanid", spanContext.SpanID().String())
	request.Header.Set("X-B3-Traceid", spanContext.TraceID().String())
	request.Header.Set("X-B3-Parentspanid", spanContext.SpanID().String())

	carrier := propagation.HeaderCarrier(request.Header)
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(t.Context, carrier)
}

// AddEvent adds an event to the trace.
func (t *Trace) AddEvent(event string, attributes map[string]string) {
	options := make([]attribute.KeyValue, 0)

	for key, value := range attributes {
		options = append(options, attribute.String(key, value))
	}

	t.Span.AddEvent(event, trace.WithAttributes(options...))
}

// RecordError records an error to the trace.
func (t *Trace) RecordError(err error, attributes map[string]string) {
	options := make([]attribute.KeyValue, 0)

	for key, value := range attributes {
		options = append(options, attribute.String(key, value))
	}

	t.Span.RecordError(err, trace.WithAttributes(options...))
}

// Failedf marks the trace as failed.
func (t *Trace) Failedf(format string, args ...any) {
	t.Span.SetStatus(codes.Error, fmt.Sprintf(format, args...))
}

// FailIf marks the trace as failed if the error is not nil.
func (t *Trace) FailIf(err error) error {
	if err != nil {
		t.Span.SetStatus(codes.Error, err.Error())
	}

	return err
}

// SetTag adds a tag to the trace.
func (t *Trace) SetTag(key, value string) {
	t.Span.SetAttributes(attribute.String(key, value))
}

// Close ends the trace.
func (t *Trace) Close() {
	if t.Span != nil {
		t.Span.End()
	}
}
