// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// ReadyStatus Ready status
//
// swagger:model ReadyStatus
type ReadyStatus map[string]string

// Validate validates this ready status
func (m ReadyStatus) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this ready status based on context it is used
func (m ReadyStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
