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
	destination := app.Spec.Destination
	if destination.TargetRevision == "" {
		destination.TargetRevision = defaultApplicationTargetRevision
	}
	if destination.Path == "" {
		destination.Path = defaultApplicationPath
	}

	for _, component := range app.Spec.Components {
		applyServiceDefaults(component.Service)
	}

	return app
}

func applyServiceDefaults(service *models.Service) {
	if service.Type == "" {
		service.Type = "ClusterIP"
	}
	for _, port := range service.Ports {
		if port.Protocol == "" {
			port.Protocol = "TCP"
		}
	}
}
