package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/koga45/golang-adot-test/backend/controller"
	"github.com/koga45/golang-adot-test/backend/pkg"
	"github.com/koga45/golang-adot-test/backend/service"
)

func main() {
	port := 9090
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("adot-collector:4317"))
	if err != nil {
		log.Fatalf("failed to creating exporter: %v", err)
	}

	rsc, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithProcessPID(),
		resource.WithProcessOwner(),
		resource.WithProcessRuntimeDescription(),
		resource.WithOSDescription(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName("backend"),
			attribute.String("hoge.fuga", "piyo"),
		),
	)
	if err != nil {
		log.Fatalf("failed to creating resouce: %v", err)
	}

	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithResource(rsc),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithIDGenerator(idg),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithInterceptorFilter(filters.None(filters.HealthCheck()))),
		)),
	)

	pkg.RegisterServiceServer(server, &controller.TestController{
		Service: &service.Service{},
	})

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Print("Server start")

	go func() {
		if err := server.Serve(listenPort); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-quit

	log.Print("Server shutdown")
}
