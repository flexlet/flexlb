// Code generated by go-swagger; DO NOT EDIT.

package instance

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/flexlet/flexlb/models"
)

// StartOKCode is the HTTP code returned for type StartOK
const StartOKCode int = 200

/*StartOK Start instance succeeded

swagger:response startOK
*/
type StartOK struct {

	/*
	  In: Body
	*/
	Payload *models.Instance `json:"body,omitempty"`
}

// NewStartOK creates StartOK with default headers values
func NewStartOK() *StartOK {

	return &StartOK{}
}

// WithPayload adds the payload to the start o k response
func (o *StartOK) WithPayload(payload *models.Instance) *StartOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start o k response
func (o *StartOK) SetPayload(payload *models.Instance) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// StartBadRequestCode is the HTTP code returned for type StartBadRequest
const StartBadRequestCode int = 400

/*StartBadRequest start bad request

swagger:response startBadRequest
*/
type StartBadRequest struct {

	/*
	  In: Body
	*/
	Payload *StartBadRequestBody `json:"body,omitempty"`
}

// NewStartBadRequest creates StartBadRequest with default headers values
func NewStartBadRequest() *StartBadRequest {

	return &StartBadRequest{}
}

// WithPayload adds the payload to the start bad request response
func (o *StartBadRequest) WithPayload(payload *StartBadRequestBody) *StartBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start bad request response
func (o *StartBadRequest) SetPayload(payload *StartBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// StartUnauthorizedCode is the HTTP code returned for type StartUnauthorized
const StartUnauthorizedCode int = 401

/*StartUnauthorized start unauthorized

swagger:response startUnauthorized
*/
type StartUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewStartUnauthorized creates StartUnauthorized with default headers values
func NewStartUnauthorized() *StartUnauthorized {

	return &StartUnauthorized{}
}

// WithPayload adds the payload to the start unauthorized response
func (o *StartUnauthorized) WithPayload(payload interface{}) *StartUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start unauthorized response
func (o *StartUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// StartForbiddenCode is the HTTP code returned for type StartForbidden
const StartForbiddenCode int = 403

/*StartForbidden start forbidden

swagger:response startForbidden
*/
type StartForbidden struct {

	/*
	  In: Body
	*/
	Payload *StartForbiddenBody `json:"body,omitempty"`
}

// NewStartForbidden creates StartForbidden with default headers values
func NewStartForbidden() *StartForbidden {

	return &StartForbidden{}
}

// WithPayload adds the payload to the start forbidden response
func (o *StartForbidden) WithPayload(payload *StartForbiddenBody) *StartForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start forbidden response
func (o *StartForbidden) SetPayload(payload *StartForbiddenBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// StartNotFoundCode is the HTTP code returned for type StartNotFound
const StartNotFoundCode int = 404

/*StartNotFound start not found

swagger:response startNotFound
*/
type StartNotFound struct {

	/*
	  In: Body
	*/
	Payload *StartNotFoundBody `json:"body,omitempty"`
}

// NewStartNotFound creates StartNotFound with default headers values
func NewStartNotFound() *StartNotFound {

	return &StartNotFound{}
}

// WithPayload adds the payload to the start not found response
func (o *StartNotFound) WithPayload(payload *StartNotFoundBody) *StartNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start not found response
func (o *StartNotFound) SetPayload(payload *StartNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// StartInternalServerErrorCode is the HTTP code returned for type StartInternalServerError
const StartInternalServerErrorCode int = 500

/*StartInternalServerError start internal server error

swagger:response startInternalServerError
*/
type StartInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewStartInternalServerError creates StartInternalServerError with default headers values
func NewStartInternalServerError() *StartInternalServerError {

	return &StartInternalServerError{}
}

// WithPayload adds the payload to the start internal server error response
func (o *StartInternalServerError) WithPayload(payload interface{}) *StartInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the start internal server error response
func (o *StartInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StartInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
