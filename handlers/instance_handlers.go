package handlers

import (
	"flexlb/config"
	"flexlb/restapi/operations/instance"

	"github.com/go-openapi/runtime/middleware"
)

type InstanceCreateHandlerImpl struct {
}

func (h *InstanceCreateHandlerImpl) Handle(params instance.CreateParams) middleware.Responder {
	inst, err := config.CreateInstance(params.Config)
	if err != nil {
		var errMsg instance.CreateBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewCreateBadRequest().WithPayload(&errMsg)
	}
	return instance.NewCreateOK().WithPayload(inst)
}

type InstanceListHandlerImpl struct {
}

func (h *InstanceListHandlerImpl) Handle(params instance.ListParams) middleware.Responder {
	insts := config.ListInstances(params.Name)
	return instance.NewListOK().WithPayload(insts)
}

type InstanceGetHandlerImpl struct {
}

func (h *InstanceGetHandlerImpl) Handle(params instance.GetParams) middleware.Responder {
	inst, err := config.GetInstance(params.Name)
	if err != nil {
		var errMsg instance.GetBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewGetBadRequest().WithPayload(&errMsg)
	}
	return instance.NewGetOK().WithPayload(inst)
}

type InstanceModifyHandlerImpl struct {
}

func (h *InstanceModifyHandlerImpl) Handle(params instance.ModifyParams) middleware.Responder {
	inst, err := config.ModifyInstance(params.Name, params.Config)
	if err != nil {
		var errMsg instance.ModifyBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewModifyBadRequest().WithPayload(&errMsg)
	}
	return instance.NewModifyOK().WithPayload(inst)
}

type InstanceDeleteHandlerImpl struct {
}

func (h *InstanceDeleteHandlerImpl) Handle(params instance.DeleteParams) middleware.Responder {
	err := config.DeleteInstance(params.Name)
	if err != nil {
		var errMsg instance.DeleteBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewDeleteBadRequest().WithPayload(&errMsg)
	}
	return instance.NewDeleteOK()
}

type InstanceStopHandlerImpl struct {
}

func (h *InstanceStopHandlerImpl) Handle(params instance.StopParams) middleware.Responder {
	inst, err := config.StopInstance(params.Name)
	if err != nil {
		var errMsg instance.StopBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewStopBadRequest().WithPayload(&errMsg)
	}
	return instance.NewStopOK().WithPayload(inst)
}

type InstanceStartHandlerImpl struct {
}

func (h *InstanceStartHandlerImpl) Handle(params instance.StartParams) middleware.Responder {
	inst, err := config.StartInstance(params.Name)
	if err != nil {
		var errMsg instance.StartBadRequestBody
		msg := err.Error()
		errMsg.Message = &msg
		return instance.NewStartBadRequest().WithPayload(&errMsg)
	}
	return instance.NewStartOK().WithPayload(inst)
}
