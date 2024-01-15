package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/koga45/golang-adot-test/bff/pkg"
)

type Service struct {
	Client pkg.ServiceClient
}

func NewService(ctx context.Context) (*Service, error) {
	u, err := url.Parse("http://backend:9090")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse url: %v", err))
	}
	conn, err := grpc.DialContext(ctx, u.Host, append([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	})...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create context: %v", err))
	}

	return &Service{
		Client: pkg.NewServiceClient(conn),
	}, nil
}

func (s *Service) Get(ctx context.Context, input string) (*string, error) {
	tracer := otel.Tracer("bff")
	ctx, span := tracer.Start(
		ctx,
		"bff.service.Get",
		trace.WithAttributes(attribute.String("input", input)))
	defer span.End()

	response, err := s.Client.Get(ctx, &pkg.GetRequest{
		Input: input,
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to grpc request: %v", err))
	}

	return &response.Output, nil
}
