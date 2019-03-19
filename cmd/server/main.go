package main

import (
	"flag"

	"deploy-wizard/gen/models"
	"deploy-wizard/gen/restapi"
	"deploy-wizard/gen/restapi/operations"
	"deploy-wizard/gen/restapi/operations/deployments"
	"deploy-wizard/gen/restapi/operations/general"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	var portFlag = flag.Int("port", 9801, "Port to run this service on")

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewDeployWizardAPI(swaggerSpec)
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

			return deployments.NewCreateDeploymentCreated().WithPayload(&models.Application{
				Name:              params.Application.Name,
				Tenant:            params.Application.Tenant,
				TargetEnvironment: params.Application.TargetEnvironment,
			})
		})

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
