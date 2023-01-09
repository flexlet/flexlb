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

// ListHandlerFunc turns a function with the right signature into a list handler
type ListHandlerFunc func(ListParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListHandlerFunc) Handle(params ListParams) middleware.Responder {
	return fn(params)
}

// ListHandler interface for that can handle valid list params
type ListHandler interface {
	Handle(ListParams) middleware.Responder
}

// NewList creates a new http.Handler for the list operation
func NewList(ctx *middleware.Context, handler ListHandler) *List {
	return &List{Context: ctx, Handler: handler}
}

/* List swagger:route GET /instances Instance list

List Instances

*/
type List struct {
	Context *middleware.Context
	Handler ListHandler
}

func (o *List) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewListParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// ListBadRequestBody list bad request body
//
// swagger:model ListBadRequestBody
type ListBadRequestBody struct {

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this list bad request body
func (o *ListBadRequestBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListBadRequestBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("listBadRequest"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this list bad request body based on context it is used
func (o *ListBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ListBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListBadRequestBody) UnmarshalBinary(b []byte) error {
	var res ListBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ListForbiddenBody list forbidden body
//
// swagger:model ListForbiddenBody
type ListForbiddenBody struct {

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this list forbidden body
func (o *ListForbiddenBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListForbiddenBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("listForbidden"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this list forbidden body based on context it is used
func (o *ListForbiddenBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ListForbiddenBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListForbiddenBody) UnmarshalBinary(b []byte) error {
	var res ListForbiddenBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ListNotFoundBody list not found body
//
// swagger:model ListNotFoundBody
type ListNotFoundBody struct {

	// error
	// Required: true
	Error *string `json:"error"`

	// status
	// Required: true
	Status *string `json:"status"`
}

// Validate validates this list not found body
func (o *ListNotFoundBody) Validate(formats strfmt.Registry) error {
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

func (o *ListNotFoundBody) validateError(formats strfmt.Registry) error {

	if err := validate.Required("listNotFound"+"."+"error", "body", o.Error); err != nil {
		return err
	}

	return nil
}

func (o *ListNotFoundBody) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("listNotFound"+"."+"status", "body", o.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this list not found body based on context it is used
func (o *ListNotFoundBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ListNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListNotFoundBody) UnmarshalBinary(b []byte) error {
	var res ListNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
