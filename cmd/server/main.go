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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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
			return general.NewGetHealthOK().WithPayload(&models.HealthStatus{Status: "OK"})
		})

	api.ValidationsValidateApplicationHandler = validations.ValidateApplicationHandlerFunc(
		func(params validations.ValidateApplicationParams) middleware.Responder {
			validationErrors := application.ValidateApplication(params.Application)
			return validations.NewValidateApplicationOK().WithPayload(&models.ValidationResponse{Errors: validationErrors})
		})

	api.AppsPreviewAppHandler = apps.PreviewAppHandlerFunc(
		func(params apps.PreviewAppParams) middleware.Responder {
			if params.Application == nil {
				return apps.NewPreviewAppBadRequest().
					WithPayload("application is required")
			}

			app := application.ApplyDefaults(params.Application)

			rendered, err := renderer.RenderApplication(app)
			if err != nil {
				return apps.NewPreviewAppDefault(500).WithPayload(err.Error())
			}

			metrics.AppsRenderedCount.WithLabelValues(app.Name).Inc()
			return apps.NewPreviewAppCreated().WithPayload(rendered)
		})

	api.AppsReleaseAppHandler = apps.ReleaseAppHandlerFunc(
		func(params apps.ReleaseAppParams) middleware.Responder {
			validationErrors := application.ValidateApplication(params.Application)
			return apps.NewReleaseAppCreated().WithPayload(&models.ValidationResponse{Errors: validationErrors})
		})

	server.ConfigureAPI()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*metricsPortFlag), nil))
	}()

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
