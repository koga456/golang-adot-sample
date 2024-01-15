package controller

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/koga45/golang-adot-test/backend/pkg"
	"github.com/koga45/golang-adot-test/backend/service"
)

type TestController struct {
	Service *service.Service
}

func (c *TestController) Get(ctx context.Context, message *pkg.GetRequest) (*pkg.GetResponse, error) {
	tracer := otel.Tracer("backend")
	ctx, span := tracer.Start(
		ctx,
		"backend.controller.Get",
		trace.WithAttributes(attribute.String("message", fmt.Sprintf("%#v", message))))
	defer span.End()

	output, _ := c.Service.Get(ctx, message.Input)

	return &pkg.GetResponse{
		Output: *output,
	}, nil
}
