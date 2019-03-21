package application_test

import (
	"strings"
	"testing"

	"deploy-wizard/gen/models"
	"deploy-wizard/pkg/application"
	"github.com/go-openapi/strfmt"
)

var (
	validApplication = &models.Application{
		Name:           "app1",
		Release:        "v1",
		Tenant:         "tenant1",
		Environment:    "Dev",
		Region:         "STL",
		Namespace:      "tenant1",
		RepoURL:        strfmt.URI("https://fusion.mastercard.int/stash/scm/ce/fake-repo.git/"),
		Path:           "/",
		TargetRevision: "HEAD",
		Services: []*models.Service{
			{
				Name: "app1",
				Tier: "Frontend",
				Ports: []*models.ServicePort{
					{
						Name: "http",
						Port: 8080,
					},
					{
						Name:       "metrics",
						Port:       8081,
						TargetPort: "8090",
						Protocol:   "TCP",
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
