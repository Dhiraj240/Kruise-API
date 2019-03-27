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
func ValidateApplication(appdata interface{}) map[string]interface{} {
	errors := map[string]interface{}{}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(appdata)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgInvalidJSON)
		errors[""] = errMsgInvalidJSON
		return errors

	}

	var app *models.Application
	err = json.Unmarshal(b.Bytes(), &app)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgNotAnApplication)
		errors[""] = errMsgNotAnApplication
		return errors
	}

	if app.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if app.Release == "" {
		errors["release"] = newRequiredValidationError("release")
	}

	if app.Environment == "" {
		errors["environment"] = newRequiredValidationError("environment")
	}

	if app.Tenant == "" {
		errors["tenant"] = newRequiredValidationError("tenant")
	}

	if app.Namespace == "" {
		errors["namespace"] = newRequiredValidationError("namespace")
	}

	if app.Path == "" {
		errors["path"] = newRequiredValidationError("path")
	}

	if app.Region == "" {
		errors["region"] = newRequiredValidationError("region")
	}

	if app.RepoURL == "" {
		errors["repoURL"] = newRequiredValidationError("repoURL")
	}

	if app.TargetRevision == "" {
		errors["targetRevision"] = newRequiredValidationError("targetRevision")
	}

	return errors
}

func newRequiredValidationError(field string) string {
	return fmt.Sprintf("%s is required", field)
}
