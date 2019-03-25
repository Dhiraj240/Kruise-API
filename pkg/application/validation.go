package application

import (
	"bytes"
	"deploy-wizard/gen/models"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	errMsgInvalidJSON      = "invalid json payload"
	errMsgNotAnApplication = "not an application object"
)

// ValidateApplication returns of map with key = field and value = error
func ValidateApplication(appdata interface{}) []*models.ValidationError {
	// errors := map[string]string{}
	var errors []*models.ValidationError

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(appdata)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgInvalidJSON)
		errors = append(errors, newValidationError("", errMsgInvalidJSON))
		return errors

	}

	var app *models.Application
	err = json.Unmarshal(b.Bytes(), &app)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgNotAnApplication)
		errors = append(errors, newValidationError("", errMsgNotAnApplication))
		return errors
	}

	if app.Name == "" {
		errors = append(errors, newRequiredValidationError("name"))
	}

	if app.Release == "" {
		errors = append(errors, newRequiredValidationError("release"))
	}

	if app.Environment == "" {
		errors = append(errors, newRequiredValidationError("environment"))
	}

	if app.Tenant == "" {
		errors = append(errors, newRequiredValidationError("tenant"))
	}

	if app.Namespace == "" {
		errors = append(errors, newRequiredValidationError("namespace"))
	}

	if app.Path == "" {
		errors = append(errors, newRequiredValidationError("path"))
	}

	if app.Region == "" {
		errors = append(errors, newRequiredValidationError("region"))
	}

	if app.RepoURL == "" {
		errors = append(errors, newRequiredValidationError("repoURL"))
	}

	if app.TargetRevision == "" {
		errors = append(errors, newRequiredValidationError("targetRevision"))
	}

	return errors
}

func newValidationError(name, error string) *models.ValidationError {
	return &models.ValidationError{
		Name:  name,
		Error: error,
	}
}

func newRequiredValidationError(name string) *models.ValidationError {
	return newValidationError(name, fmt.Sprintf("%s is required", name))
}
