package handlers

import (
	"github.com/flexlet/flexlb/pkg/config"
	"github.com/flexlet/flexlb/restapi/operations/service"

	"github.com/go-openapi/runtime/middleware"
)

type ReadyzHandlerImpl struct {
}

func (h *ReadyzHandlerImpl) Handle(params service.ReadyzParams) middleware.Responder {
	return service.NewReadyzOK().WithPayload(config.LB.Status)
}
