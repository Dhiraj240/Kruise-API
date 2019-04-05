package application

import (
	"deploy-wizard/gen/models"
)

const (
	defaultApplicationTargetRevision = "HEAD"
	defaultApplicationPath           = "/"
)

// ApplyDefaults applies defaults to the Application model
func ApplyDefaults(app *models.Application) *models.Application {
	if app.TargetRevision == "" {
		app.TargetRevision = defaultApplicationTargetRevision
	}
	if app.Path == "" {
		app.Path = defaultApplicationPath
	}

	for _, service := range app.Services {
		applyServiceDefaults(service)
	}

	return app
}

func applyServiceDefaults(service *models.Service) {
	for _, port := range service.Ports {
		if port.Protocol == "" {
			port.Protocol = "TCP"
		}
	}
}
