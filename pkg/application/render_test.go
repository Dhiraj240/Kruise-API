package application_test

import (
	"strings"
	"testing"

	"deploy-wizard/gen/models"
	"deploy-wizard/pkg/application"
	"github.com/andreyvit/diff"
	"github.com/go-openapi/strfmt"
)

func defaultTier() *string {
	t := "Frontend"
	return &t
}

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
				Tier: defaultTier(),
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
				Containers: []*models.Container{
					{
						Name:            "app1",
						Image:           "nginx",
						ImagePullPolicy: "IfNotPresent",
						ImageTag:        "alpine",
					},
				},
			},
		},
		Ingress: &models.Ingress{
			Name: "app1-ingress",
			Rules: []*models.IngressRule{
				{
					Host:        "app1.mc.int",
					ServiceName: "app1",
					ServicePort: "http",
				},
			},
		},
	}

	expected = `apiVersion: v1
kind: ServiceAccount
metadata:
	labels:
		app: app1
		release: v1
	name: app1


---
apiVersion: v1
kind: Service
metadata:
	labels:
		component: app1
		app: app1
		release: v1
		tier: Frontend
	name: app1
spec:
	ports:
	- name: http
		port: 8080
		protocol: TCP
	- name: metrics
		port: 8081
		protocol: TCP
		targetPort: 8090
	selector:
		app: app1
		component: app1
	type: ClusterIP


---
apiVersion: apps/v1
kind: Deployment
metadata:
	name: app1
	labels:
		app: app1
		component: app1
		release: v1
spec:
	replicas: 1
	selector:
		matchLabels:
			app: app1
			component: app1
	strategy:
		type: RollingUpdate
	template:
		metadata:
			labels:
				app: app1
				component: app1
				release: v1
		spec:
			affinity:
				podAntiAffinity:
					preferredDuringSchedulingIgnoredDuringExecution:
					- podAffinityTerm:
							labelSelector:
								matchLabels:
									app: app1
									component: app1
									release: v1
							topologyKey: kubernetes.io/hostname
						weight: 100
			volumes:
			- name: ca-bundles
				configMap:
					name: ca-bundles
			containers:
			- name: app1
				image: nginx:alpine
				imagePullPolicy: IfNotPresent
				volumeMounts:
				- mountPath: "/etc/ssl/certs"
					name: ca-bundles
					readOnly: true
				ports:
				- name: http
					containerPort: 8080
					protocol: TCP
				- name: metrics
					containerPort: 8090
					protocol: TCP
`
)

func TestRenderApplication(t *testing.T) {
	app := application.ApplyDefaults(validApplication)

	renderer, err := application.NewRenderer("../../_templates")
	if err != nil {
		t.Error(err)
	}

	result, err := renderer.RenderApplication(app)

	if err != nil {
		t.Error(err)
	}

	expected = strings.Replace(expected, "\t", "  ", -1)
	if result != expected {
		t.Errorf("Result not as expected:\n%v", diff.LineDiff(result, expected))
	}
}
