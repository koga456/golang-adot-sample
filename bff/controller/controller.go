package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/koga45/golang-adot-test/bff/service"
)

type Controller struct {
	Service *service.Service
}

func (c *Controller) Get(gc *gin.Context) {

	input := gc.Params.ByName("input")

	tracer := otel.Tracer("bff")
	ctx, span := tracer.Start(
		gc.Request.Context(),
		"bff.controller.Get",
		trace.WithAttributes(attribute.String("input", input)))
	defer span.End()

	output, err := c.Service.Get(ctx, input)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, Response{
		Output: *output,
	})
}

type Response struct {
	Output string `json:"output"`
}
