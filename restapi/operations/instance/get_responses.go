// Code generated by go-swagger; DO NOT EDIT.

package instance

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"flexlb/models"
)

// GetOKCode is the HTTP code returned for type GetOK
const GetOKCode int = 200

/*GetOK Get Instance succeeded

swagger:response getOK
*/
type GetOK struct {

	/*
	  In: Body
	*/
	Payload *models.Instance `json:"body,omitempty"`
}

// NewGetOK creates GetOK with default headers values
func NewGetOK() *GetOK {

	return &GetOK{}
}

// WithPayload adds the payload to the get o k response
func (o *GetOK) WithPayload(payload *models.Instance) *GetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get o k response
func (o *GetOK) SetPayload(payload *models.Instance) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetBadRequestCode is the HTTP code returned for type GetBadRequest
const GetBadRequestCode int = 400

/*GetBadRequest get bad request

swagger:response getBadRequest
*/
type GetBadRequest struct {

	/*
	  In: Body
	*/
	Payload *GetBadRequestBody `json:"body,omitempty"`
}

// NewGetBadRequest creates GetBadRequest with default headers values
func NewGetBadRequest() *GetBadRequest {

	return &GetBadRequest{}
}

// WithPayload adds the payload to the get bad request response
func (o *GetBadRequest) WithPayload(payload *GetBadRequestBody) *GetBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get bad request response
func (o *GetBadRequest) SetPayload(payload *GetBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUnauthorizedCode is the HTTP code returned for type GetUnauthorized
const GetUnauthorizedCode int = 401

/*GetUnauthorized get unauthorized

swagger:response getUnauthorized
*/
type GetUnauthorized struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetUnauthorized creates GetUnauthorized with default headers values
func NewGetUnauthorized() *GetUnauthorized {

	return &GetUnauthorized{}
}

// WithPayload adds the payload to the get unauthorized response
func (o *GetUnauthorized) WithPayload(payload interface{}) *GetUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get unauthorized response
func (o *GetUnauthorized) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetForbiddenCode is the HTTP code returned for type GetForbidden
const GetForbiddenCode int = 403

/*GetForbidden get forbidden

swagger:response getForbidden
*/
type GetForbidden struct {

	/*
	  In: Body
	*/
	Payload *GetForbiddenBody `json:"body,omitempty"`
}

// NewGetForbidden creates GetForbidden with default headers values
func NewGetForbidden() *GetForbidden {

	return &GetForbidden{}
}

// WithPayload adds the payload to the get forbidden response
func (o *GetForbidden) WithPayload(payload *GetForbiddenBody) *GetForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get forbidden response
func (o *GetForbidden) SetPayload(payload *GetForbiddenBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNotFoundCode is the HTTP code returned for type GetNotFound
const GetNotFoundCode int = 404

/*GetNotFound get not found

swagger:response getNotFound
*/
type GetNotFound struct {

	/*
	  In: Body
	*/
	Payload *GetNotFoundBody `json:"body,omitempty"`
}

// NewGetNotFound creates GetNotFound with default headers values
func NewGetNotFound() *GetNotFound {

	return &GetNotFound{}
}

// WithPayload adds the payload to the get not found response
func (o *GetNotFound) WithPayload(payload *GetNotFoundBody) *GetNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get not found response
func (o *GetNotFound) SetPayload(payload *GetNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetInternalServerErrorCode is the HTTP code returned for type GetInternalServerError
const GetInternalServerErrorCode int = 500

/*GetInternalServerError get internal server error

swagger:response getInternalServerError
*/
type GetInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetInternalServerError creates GetInternalServerError with default headers values
func NewGetInternalServerError() *GetInternalServerError {

	return &GetInternalServerError{}
}

// WithPayload adds the payload to the get internal server error response
func (o *GetInternalServerError) WithPayload(payload interface{}) *GetInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get internal server error response
func (o *GetInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}