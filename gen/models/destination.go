// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Destination destination
// swagger:model destination
type Destination struct {

	// The relative path to the manifests in the git repo
	// Min Length: 1
	Path string `json:"path,omitempty"`

	// Defines the commit, tag, or branch in which to sync the application to.
	// Min Length: 1
	TargetRevision string `json:"targetRevision,omitempty"`

	// The git repo URL
	// Required: true
	// Min Length: 1
	// Format: uri
	URL strfmt.URI `json:"url"`
}

// Validate validates this destination
func (m *Destination) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePath(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTargetRevision(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURL(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Destination) validatePath(formats strfmt.Registry) error {

	if swag.IsZero(m.Path) { // not required
		return nil
	}

	if err := validate.MinLength("path", "body", string(m.Path), 1); err != nil {
		return err
	}

	return nil
}

func (m *Destination) validateTargetRevision(formats strfmt.Registry) error {

	if swag.IsZero(m.TargetRevision) { // not required
		return nil
	}

	if err := validate.MinLength("targetRevision", "body", string(m.TargetRevision), 1); err != nil {
		return err
	}

	return nil
}

func (m *Destination) validateURL(formats strfmt.Registry) error {

	if err := validate.Required("url", "body", strfmt.URI(m.URL)); err != nil {
		return err
	}

	if err := validate.MinLength("url", "body", string(m.URL), 1); err != nil {
		return err
	}

	if err := validate.FormatOf("url", "body", "uri", m.URL.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Destination) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Destination) UnmarshalBinary(b []byte) error {
	var res Destination
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}