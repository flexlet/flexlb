// Code generated by go-swagger; DO NOT EDIT.

package instance

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// StopHandlerFunc turns a function with the right signature into a stop handler
type StopHandlerFunc func(StopParams) middleware.Responder

// Handle executing the request and returning a response
func (fn StopHandlerFunc) Handle(params StopParams) middleware.Responder {
	return fn(params)
}

// StopHandler interface for that can handle valid stop params
type StopHandler interface {
	Handle(StopParams) middleware.Responder
}

// NewStop creates a new http.Handler for the stop operation
func NewStop(ctx *middleware.Context, handler StopHandler) *Stop {
	return &Stop{Context: ctx, Handler: handler}
}

/* Stop swagger:route POST /instances/{name}/stop Instance stop

Stop Instance

*/
type Stop struct {
	Context *middleware.Context
	Handler StopHandler
}

func (o *Stop) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewStopParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// StopBadRequestBody stop bad request body
//
// swagger:model StopBadRequestBody
type StopBadRequestBody struct {

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this stop bad request body
func (o *StopBadRequestBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *StopBadRequestBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("stopBadRequest"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this stop bad request body based on context it is used
func (o *StopBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *StopBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *StopBadRequestBody) UnmarshalBinary(b []byte) error {
	var res StopBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// StopForbiddenBody stop forbidden body
//
// swagger:model StopForbiddenBody
type StopForbiddenBody struct {

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this stop forbidden body
func (o *StopForbiddenBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *StopForbiddenBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("stopForbidden"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this stop forbidden body based on context it is used
func (o *StopForbiddenBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *StopForbiddenBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *StopForbiddenBody) UnmarshalBinary(b []byte) error {
	var res StopForbiddenBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// StopNotFoundBody stop not found body
//
// swagger:model StopNotFoundBody
type StopNotFoundBody struct {

	// error
	// Required: true
	Error *string `json:"error"`

	// status
	// Required: true
	Status *string `json:"status"`
}

// Validate validates this stop not found body
func (o *StopNotFoundBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateError(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *StopNotFoundBody) validateError(formats strfmt.Registry) error {

	if err := validate.Required("stopNotFound"+"."+"error", "body", o.Error); err != nil {
		return err
	}

	return nil
}

func (o *StopNotFoundBody) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("stopNotFound"+"."+"status", "body", o.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this stop not found body based on context it is used
func (o *StopNotFoundBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *StopNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *StopNotFoundBody) UnmarshalBinary(b []byte) error {
	var res StopNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
