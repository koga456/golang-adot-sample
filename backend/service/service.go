package service

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Service struct{}

func (s *Service) Get(ctx context.Context, input string) (*string, error) {
	tracer := otel.Tracer("backend")
	ctx, span := tracer.Start(
		ctx,
		"backend.controller.Get",
		trace.WithAttributes(attribute.String("input", input)))
	defer span.End()

	output := "hoge"
	return &output, nil
}
