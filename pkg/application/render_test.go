package application_test

import (
	"strings"
	"testing"

	"deploy-wizard/gen/models"
	"deploy-wizard/pkg/application"
	"github.com/go-openapi/strfmt"
)

var (
	name            = "app1"
	tenant          = "tenant1"
	environment     = "Dev"
	region          = "STL"
	namespace       = "tenant1"
	repoURL         = strfmt.URI("https://fusion.mastercard.int/stash/scm/ce/fake-repo.git/")
	path            = "/"
	targetRevision  = "HEAD"
	tier            = "Frontend"
	httpPorts       = []int64{8080, 8000}
	metricsPorts    = []int64{8081, 8001}
	httpPortName    = "http"
	metricsPortName = "metrics"

	validApplication = &models.Application{
		Name:           &name,
		Tenant:         &tenant,
		Environment:    &environment,
		Region:         &region,
		Namespace:      &namespace,
		RepoURL:        &repoURL,
		Path:           path,
		TargetRevision: targetRevision,
		Services: []*models.Service{
			{
				Name: &name,
				Tier: &tier,
				Ports: []*models.ServicePort{
					{
						Name: &httpPortName,
						Port: &httpPorts[0],
					},
					{
						Name: &metricsPortName,
						Port: &metricsPorts[0],
					},
				},
			},
		},
		Ingresses: []*models.Ingress{},
	}

	expected = `apiVersion: v1
kind: ServiceAccount
metadata:
	labels:
		app: app1
		release: v1
	name: app1
`
)

func TestRenderApplication(t *testing.T) {
	renderer, err := application.NewRenderer("../../_templates")
	if err != nil {
		t.Error(err)
	}

	result, err := renderer.RenderApplication(validApplication)

	if err != nil {
		t.Error(err)
	}

	expected = strings.Replace(expected, "\t", "  ", -1)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
