package handlers

import (
	"fmt"

	"github.com/flexlet/flexlb/pkg/config"
	"github.com/flexlet/flexlb/restapi/operations/instance"

	"github.com/go-openapi/runtime/middleware"
)

type InstanceCreateHandlerImpl struct {
}

func (h *InstanceCreateHandlerImpl) Handle(params instance.CreateParams) middleware.Responder {
	// cluster instancce blackout
	if params.Config.Name == config.LB.Cluster.Name {
		msg := fmt.Sprintf("cluster instance '%s' cannot be created manually", params.Config.Name)
		return instance.NewCreateBadRequest().WithPayload(&instance.CreateBadRequestBody{Message: &msg})
	}

	// create instance
	inst, err := config.CreateInstance(params.Config)
	if err != nil {
		msg := err.Error()
		return instance.NewCreateBadRequest().WithPayload(&instance.CreateBadRequestBody{Message: &msg})
	}

	// notify other nodes
	config.GossipCreateInstance(inst)

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
		msg := err.Error()
		return instance.NewGetBadRequest().WithPayload(&instance.GetBadRequestBody{Message: &msg})
	}
	return instance.NewGetOK().WithPayload(inst)
}

type InstanceModifyHandlerImpl struct {
}

func (h *InstanceModifyHandlerImpl) Handle(params instance.ModifyParams) middleware.Responder {
	// cluster instancce blackout
	if params.Config.Name == config.LB.Cluster.Name {
		msg := fmt.Sprintf("cluster instance '%s' cannot be modified manually", params.Config.Name)
		return instance.NewModifyBadRequest().WithPayload(&instance.ModifyBadRequestBody{Message: &msg})
	}

	// modify instance
	inst, err := config.ModifyInstance(params.Config)
	if err != nil {
		msg := err.Error()
		return instance.NewModifyBadRequest().WithPayload(&instance.ModifyBadRequestBody{Message: &msg})
	}

	// notify other nodes
	config.GossipModifyInstance(inst)

	return instance.NewModifyOK().WithPayload(inst)
}

type InstanceDeleteHandlerImpl struct {
}

func (h *InstanceDeleteHandlerImpl) Handle(params instance.DeleteParams) middleware.Responder {
	// cluster instancce blackout
	if params.Name == config.LB.Cluster.Name {
		msg := fmt.Sprintf("cluster instance '%s' cannot be deleted manually", params.Name)
		return instance.NewDeleteBadRequest().WithPayload(&instance.DeleteBadRequestBody{Message: &msg})
	}

	// delete instance
	err := config.DeleteInstance(params.Name)
	if err != nil {
		msg := err.Error()
		return instance.NewDeleteBadRequest().WithPayload(&instance.DeleteBadRequestBody{Message: &msg})
	}

	// notify other nodes
	config.GossipDeleteInstance(params.Name)

	return instance.NewDeleteOK()
}

type InstanceStopHandlerImpl struct {
}

func (h *InstanceStopHandlerImpl) Handle(params instance.StopParams) middleware.Responder {
	// cluster instancce blackout
	if params.Name == config.LB.Cluster.Name {
		msg := fmt.Sprintf("cluster instance '%s' cannot be stoped manually", params.Name)
		return instance.NewStopBadRequest().WithPayload(&instance.StopBadRequestBody{Message: &msg})
	}

	// stop instance
	inst, err := config.StopInstance(params.Name)
	if err != nil {
		msg := err.Error()
		return instance.NewStopBadRequest().WithPayload(&instance.StopBadRequestBody{Message: &msg})
	}

	// notify other nodes
	config.GossipStopInstance(params.Name)

	return instance.NewStopOK().WithPayload(inst)
}

type InstanceStartHandlerImpl struct {
}

func (h *InstanceStartHandlerImpl) Handle(params instance.StartParams) middleware.Responder {
	// cluster instancce blackout
	if params.Name == config.LB.Cluster.Name {
		msg := fmt.Sprintf("cluster instance '%s' cannot be started manually", params.Name)
		return instance.NewStartBadRequest().WithPayload(&instance.StartBadRequestBody{Message: &msg})
	}

	// start instance
	inst, err := config.StartInstance(params.Name)
	if err != nil {
		msg := err.Error()
		return instance.NewStartBadRequest().WithPayload(&instance.StartBadRequestBody{Message: &msg})
	}

	// notify other nodes
	config.GossipStartInstance(params.Name)

	return instance.NewStartOK().WithPayload(inst)
}
