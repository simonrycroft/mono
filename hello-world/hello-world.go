package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // Alias this import to 'sdktrace'
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

var tracer trace.Tracer

func main() {
	// Setup Prometheus metrics
	requests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hello_world_requests_total",
		Help: "Total number of hello-world requests.",
	})
	prometheus.MustRegister(requests)

	// Initialize the OpenTelemetry setup
	ctx := context.Background()

	// Create an exporter that will send trace data to the OpenTelemetry Collector
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Create a TracerProvider and configure it to use the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("hello-world-service"),
		)),
	)

	// Register the TracerProvider as the global tracer provider
	otel.SetTracerProvider(tp)
	defer func() { _ = tp.Shutdown(ctx) }()

	// Create a tracer
	tracer = otel.Tracer("hello-world-tracer")

	// HTTP handler for /hello endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Start a span for the request
		_, span := tracer.Start(r.Context(), "HelloHandler")
		defer span.End()

		// Adding an example attribute to the span
		span.SetAttributes(attribute.String("http.method", r.Method))

		// Increase the Prometheus counter
		requests.Inc()

		fmt.Fprintln(w, "Hello, World!")
	})

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
