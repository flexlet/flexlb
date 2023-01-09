// Code generated by go-swagger; DO NOT EDIT.

package instance

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewStopParams creates a new StopParams object
//
// There are no default values defined in the spec.
func NewStopParams() StopParams {

	return StopParams{}
}

// StopParams contains all the bound params for the stop operation
// typically these are obtained from a http.Request
//
// swagger:parameters stop
type StopParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Instance name
	  Required: true
	  Pattern: ^[A-Za-z0-9\-_.]{1,32}$
	  In: path
	*/
	Name string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewStopParams() beforehand.
func (o *StopParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindName binds and validates parameter Name from path.
func (o *StopParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.Name = raw

	if err := o.validateName(formats); err != nil {
		return err
	}

	return nil
}

// validateName carries on validations for parameter Name
func (o *StopParams) validateName(formats strfmt.Registry) error {

	if err := validate.Pattern("name", "path", o.Name, `^[A-Za-z0-9\-_.]{1,32}$`); err != nil {
		return err
	}

	return nil
}
