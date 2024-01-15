package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/koga45/golang-adot-test/bff/controller"
	"github.com/koga45/golang-adot-test/bff/service"
)

func main() {

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
			semconv.ServiceName("bff"),
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

	s, err := service.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}
	c := &controller.Controller{
		Service: s,
	}

	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()

	filter := func(req *http.Request) bool {
		return req.URL.Path != "/" && !strings.Contains(req.URL.Path, "/health")
	}
	handler.Use(otelgin.Middleware("bff-http-server", otelgin.WithFilter(filter)))

	test := handler.Group("/test")
	test.GET("/:input/", c.Get)

	health := handler.Group("/health")
	health.GET("/", func(gc *gin.Context) {
		gc.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	handler.NoRoute(func(gc *gin.Context) {
		gc.JSON(http.StatusNotFound, gin.H{"message": "Page not Found"})
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 20 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Print("Server start")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-quit

	log.Print("Server shutdown")
}
