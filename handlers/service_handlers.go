package handlers

import (
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/restapi/operations/service"

	"github.com/go-openapi/runtime/middleware"
)

type ReadyzHandlerImpl struct {
}

func (h *ReadyzHandlerImpl) Handle(params service.ReadyzParams) middleware.Responder {
	return service.NewReadyzOK().WithPayload(config.LB.Status)
}
