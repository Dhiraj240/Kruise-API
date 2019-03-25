package application

import (
	"bytes"
	"deploy-wizard/gen/models"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// ValidateApplication returns of map with key = field and value = error
func ValidateApplication(appdata interface{}) map[string]string {
	errors := map[string]string{}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(appdata)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf("invalid json payload")
		errors[""] = "invalid json payload"
		return errors

	}

	var app *models.Application
	err = json.Unmarshal(b.Bytes(), &app)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf("not an application object")
		errors[""] = "not an application object"
		return errors
	}

	if app.Name == "" {
		errors["name"] = "name is required"
	}

	if app.Release == "" {
		errors["release"] = "release is required"
	}

	if app.Environment == "" {
		errors["environment"] = "environment is required"
	}

	if app.Tenant == "" {
		errors["tenant"] = "tenant is required"
	}

	if app.Namespace == "" {
		errors["namespace"] = "namespace is required"
	}

	if app.Path == "" {
		errors["path"] = "path is required"
	}

	if app.Region == "" {
		errors["region"] = "region is required"
	}

	if app.RepoURL == "" {
		errors["repoURL"] = "repoURL is required"
	}

	if app.TargetRevision == "" {
		errors["targetRevision"] = "targetRevision is required"
	}

	return errors
}
