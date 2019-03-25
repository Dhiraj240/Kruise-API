package main

import (
	"flag"
	"net/http"
	"strconv"

	"deploy-wizard/gen/models"
	"deploy-wizard/gen/restapi"
	"deploy-wizard/gen/restapi/operations"
	"deploy-wizard/gen/restapi/operations/apps"
	"deploy-wizard/gen/restapi/operations/general"
	"deploy-wizard/gen/restapi/operations/validations"
	"deploy-wizard/pkg/application"
	"deploy-wizard/pkg/metrics"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	defaultApplicationTargetRevision = "HEAD"
	defaultApplicationPath           = "/"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var portFlag = flag.Int("port", 9801, "Port to run this service on")
	var metricsPortFlag = flag.Int("metrics-port", 9802, "Metrics port")

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	api := operations.NewDeployWizardAPI(swaggerSpec)
	api.Logger = log.Infof

	server := restapi.NewServer(api)
	defer func() {
		_ = server.Shutdown()
	}()

	// parse flags
	flag.Parse()
	// set the port this service will be run on
	server.Port = *portFlag

	// TODO: flag template directory
	renderer, err := application.NewRenderer("./_templates")
	if err != nil {
		log.Fatal(err)
	}

	api.GeneralGetHealthHandler = general.GetHealthHandlerFunc(
		func(params general.GetHealthParams) middleware.Responder {
			return general.NewGetHealthOK().WithPayload(&models.HealthStatus{"OK"})
		})

	api.ValidationsValidateApplicationHandler = validations.ValidateApplicationHandlerFunc(
		func(params validations.ValidateApplicationParams) middleware.Responder {
			response := &models.ValidationResponse{}

			validationErrors := application.ValidateApplication(params.Application)
			for name, error := range validationErrors {
				response.Errors = append(response.Errors, &models.ValidationError{
					Name:  name,
					Error: error,
				})
			}

			return validations.NewValidateApplicationOK().WithPayload(response)
		})

	api.AppsCreateAppHandler = apps.CreateAppHandlerFunc(
		func(params apps.CreateAppParams) middleware.Responder {
			if params.Application == nil {
				return apps.NewCreateAppBadRequest().WithPayload("application is required")
			}

			app := application.ApplyDefaults(params.Application)

			rendered, err := renderer.RenderApplication(app)
			if err != nil {
				return apps.NewCreateAppDefault(500).WithPayload(err.Error())
			}

			metrics.AppsRenderedCount.WithLabelValues(app.Name).Inc()
			return apps.NewCreateAppCreated().WithPayload(rendered)
		})

	server.ConfigureAPI()

	go func() {
		http.Handle("/metrics", prometheus.Handler())
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*metricsPortFlag), nil))
	}()

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
