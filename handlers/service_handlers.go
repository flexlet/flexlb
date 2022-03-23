package handlers

import (
	"flexlb/config"
	"flexlb/models"
	"flexlb/restapi/operations/service"

	"github.com/go-openapi/runtime/middleware"
)

type ReadyzHandlerImpl struct {
}

func (h *ReadyzHandlerImpl) Handle(params service.ReadyzParams) middleware.Responder {
	var readyStatus models.ReadyStatus
	readyStatus.Status = config.ReadyStatus
	return service.NewReadyzOK().WithPayload(&readyStatus)
}
