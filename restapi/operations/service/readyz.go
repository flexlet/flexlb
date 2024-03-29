// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReadyzHandlerFunc turns a function with the right signature into a readyz handler
type ReadyzHandlerFunc func(ReadyzParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ReadyzHandlerFunc) Handle(params ReadyzParams) middleware.Responder {
	return fn(params)
}

// ReadyzHandler interface for that can handle valid readyz params
type ReadyzHandler interface {
	Handle(ReadyzParams) middleware.Responder
}

// NewReadyz creates a new http.Handler for the readyz operation
func NewReadyz(ctx *middleware.Context, handler ReadyzHandler) *Readyz {
	return &Readyz{Context: ctx, Handler: handler}
}

/* Readyz swagger:route GET /readyz Service readyz

Ready status

*/
type Readyz struct {
	Context *middleware.Context
	Handler ReadyzHandler
}

func (o *Readyz) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReadyzParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
