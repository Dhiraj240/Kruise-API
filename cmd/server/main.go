package main

import (
	"flag"

	"deploy-wizard/gen/models"
	"deploy-wizard/gen/restapi"
	"deploy-wizard/gen/restapi/operations"
	"deploy-wizard/gen/restapi/operations/deployments"
	"deploy-wizard/gen/restapi/operations/general"
	"deploy-wizard/pkg/application"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

const (
	defaultApplicationTargetRevision = "HEAD"
	defaultApplicationPath           = "/"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var portFlag = flag.Int("port", 9801, "Port to run this service on")

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

	api.GeneralGetHealthHandler = general.GetHealthHandlerFunc(
		func(params general.GetHealthParams) middleware.Responder {
			return general.NewGetHealthOK().WithPayload(&models.HealthStatus{"OK"})
		})

	api.DeploymentsCreateDeploymentHandler = deployments.CreateDeploymentHandlerFunc(
		func(params deployments.CreateDeploymentParams) middleware.Responder {
			if params.Application == nil {
				return deployments.NewCreateDeploymentBadRequest().WithPayload("application is required")
			}

			return deployments.NewCreateDeploymentCreated().WithPayload(application.ApplyDefaults(params.Application))
		})

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
