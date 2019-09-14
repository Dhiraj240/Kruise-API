package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"deploy-wizard/gen/models"
	"deploy-wizard/gen/restapi"
	"deploy-wizard/gen/restapi/operations"
	"deploy-wizard/gen/restapi/operations/apps"
	"deploy-wizard/gen/restapi/operations/general"
	"deploy-wizard/gen/restapi/operations/validations"
	"deploy-wizard/pkg/application"
	"deploy-wizard/pkg/git"
	"deploy-wizard/pkg/metrics"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	envUsernameVar      = "KRUISE_STASH_USERNAME"
	envPasswordVar      = "KRUISE_STASH_PASSWORD"
	codeRenderError     = 101
	codeRepoCloneError  = 301
	codeRepoCommitError = 302
	codeRepoPushError   = 303
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var (
		stashUserFile         string
		stashUser             string
		stashPasswordFile     string
		stashPassword         string
		gitInsecureSkipVerify bool
	)

	var portFlag = flag.Int("port", 9801, "Port to run this service on")
	var metricsPortFlag = flag.Int("metrics-port", 9802, "Metrics port")
	flag.BoolVar(&gitInsecureSkipVerify, "git-insecure-skip-verify", false, "If true, will ignore TLS verification errors (insecure)")
	flag.StringVar(&stashUserFile, "username-file", "", "Path to a file that contains the stash username")
	flag.StringVar(&stashPasswordFile, "password-file", "", "Path to a file that contains the stash password")

	// parse flags
	flag.Parse()

	if gitInsecureSkipVerify {
		log.Warn("Ignoring TLS verification errors")
	}

	if stashUserFile != "" {
		stashUser = loadFromFile(stashUserFile)
	}

	if stashUser == "" {
		stashUser = os.Getenv(envUsernameVar)
		if stashUser == "" {
			log.Fatalf("set a valid username for stash in a an environment variable called %s", envUsernameVar)
		}
	}

	if stashPasswordFile != "" {
		stashPassword = loadFromFile(stashPasswordFile)
	}

	if stashPassword == "" {
		stashPassword = os.Getenv(envPasswordVar)
		if stashPassword == "" {
			log.Fatalf("set a valid password for stash in a an environment variable called %s", envPasswordVar)
		}
	}

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
				return apps.NewPreviewAppBadRequest().WithPayload("application is required")
			}

			app := application.ApplyDefaults(params.Application)

			rendered, err := renderer.RenderApplication(app)
			if err != nil {
				errResp := &models.Error{Code: codeRenderError, Message: err.Error()}
				return apps.NewPreviewAppDefault(500).WithPayload(errResp)
			}

			metrics.AppsRenderedCount.WithLabelValues(app.Metadata.Name).Inc()
			return apps.NewPreviewAppCreated().WithPayload(rendered)
		})

	api.AppsReleaseAppHandler = apps.ReleaseAppHandlerFunc(
		func(params apps.ReleaseAppParams) middleware.Responder {
			if params.Application == nil {
				return apps.NewReleaseAppBadRequest().WithPayload(&models.ValidationResponse{
					Errors: map[string]interface{}{
						"application": fmt.Sprintf("%q is required", "application"),
					},
				})
			}

			app := application.ApplyDefaults(params.Application)
			validationErrors := application.ValidateApplication(app)

			if len(validationErrors) > 0 {
				return apps.NewReleaseAppBadRequest().WithPayload(&models.ValidationResponse{Errors: validationErrors})
			}

			rendered, err := renderer.RenderManifests(app)
			if err != nil {
				errResp := &models.Error{Code: codeRenderError, Message: err.Error()}
				return apps.NewReleaseAppDefault(500).WithPayload(errResp)
			}

			repo := git.NewRepo(
				app.Spec.Destination.URL.String(),
				app.Spec.Destination.Path,
				app.Spec.Destination.TargetRevision,
				&git.RepoCreds{
					Username: stashUser,
					Password: stashPassword,
				}, gitInsecureSkipVerify)

			err = repo.Clone()
			if err != nil {
				errResp := &models.Error{Code: codeRepoCloneError, Message: err.Error()}
				return apps.NewReleaseAppDefault(500).WithPayload(errResp)
			}

			for filename, content := range rendered {
				log.Infof("adding file %q (%d bytes)", filename, len(content))
				repo.AddFile(filename, content)
			}

			err = repo.Commit(
				fmt.Sprintf("kruise release for %s:%s",
					app.Metadata.Name,
					app.Metadata.Labels.Version,
				))
			if err != nil {
				errResp := &models.Error{Code: codeRepoCommitError, Message: err.Error()}
				return apps.NewReleaseAppDefault(500).WithPayload(errResp)
			}

			err = repo.Push()
			if err != nil {
				errResp := &models.Error{Code: codeRepoPushError, Message: err.Error()}
				return apps.NewReleaseAppDefault(500).WithPayload(errResp)
			}

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

func loadFromFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Warnf("could not open file %q: %s", filename, err.Error())
		return ""
	}
	defer func() { _ = f.Close() }()

	s := bufio.NewScanner(f)
	s.Scan()
	return s.Text()
}
