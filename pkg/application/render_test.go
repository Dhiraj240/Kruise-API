package application_test

import (
	"strings"
	"testing"

	"deploy-wizard/gen/models"
	"deploy-wizard/pkg/application"

	"github.com/andreyvit/diff"
	"github.com/go-openapi/strfmt"
)

var (
	validApplication = &models.Application{
		Metadata: &models.Metadata{
			Name:      "app1",
			Namespace: "tenant1",
			Labels: &models.Labels{
				Version: "v1",
				Team:    "tenant1",
				Env:     "Dev",
				Region:  "STL",
			},
		},
		Spec: &models.Spec{
			Destination: &models.Destination{
				URL:            strfmt.URI("https://fusion.mastercard.int/stash/scm/ce/fake-repo.git/"),
				Path:           "/",
				TargetRevision: "HEAD",
			},
			ConfigMaps: []*models.ConfigMap{
				{
					Name: "config",
					Data: "debug: true",
				},
			},
			PersistentVolumes: []*models.PersistentVolume{
				{
					Name:             "data",
					Capacity:         20,
					AccessMode:       models.PersistentVolumeAccessModeReadWriteOnce,
					StorageClassName: "SSD",
				},
			},
			Components: []*models.Component{
				{
					Service: &models.Service{
						Name: "app1",
						Type: "ClusterIP",
						Ports: []*models.ServicePort{
							{
								Name: "http",
								Port: 8080,
							},
							{
								Name:       "metrics",
								Port:       8081,
								TargetPort: 8090,
								Protocol:   "TCP",
							},
						},
					},
					Containers: []*models.Container{
						{
							Name:            "app1",
							Image:           "nginx",
							ImagePullPolicy: "IfNotPresent",
							ImageTag:        "alpine",
							PortNames:       []string{"http", "metrics"},
							Volumes: []*models.VolumeMount{
								{
									MountPath: "/config",
									Name:      "config",
									ReadOnly:  true,
									Type:      "ConfigMap",
								},
								{
									MountPath: "/data",
									Name:      "data",
									ReadOnly:  false,
									Type:      "PersistentVolume",
								},
							},
						},
					},
					Ingresses: []*models.Ingress{
						{
							Host: "app1.mc.int",
							Paths: []*models.IngressPath{
								{
									Path:     "/",
									PortName: "http",
								},
							},
						},
					},
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
			- name: config
				configMap:
					name: config
			- name: data
				persistentVolumeClaim:
					claimName: data
			containers:
			- name: app1
				image: nginx:alpine
				imagePullPolicy: IfNotPresent
				volumeMounts:
				- mountPath: /etc/ssl/certs
					name: ca-bundles
					readOnly: true
				- mountPath: /config
					name: config
					readOnly: true
				- mountPath: /data
					name: data
					readOnly: false
				ports:
				- name: http
					containerPort: 8080
					protocol: TCP
				- name: metrics
					containerPort: 8090
					protocol: TCP


---
apiVersion: v1
kind: ConfigMap
metadata:
	labels:
		app: app1
		release: v1
	name: config
data:
	data: |
		debug: true


---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
	labels:
		app: app1
		release: v1
	name: data
spec:
	accessModes:
	- ReadWriteOnce
	resources:
		requests:
			storage: 20Gi
	storageClassName: SSD
`
)

func TestRenderManifests(t *testing.T) {
	app := application.ApplyDefaults(validApplication)

	renderer, err := application.NewRenderer("../../_templates")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	results, err := renderer.RenderManifests(app)

	if err != nil {
		t.Error(err)
	}

	if _, ok := results["service-account.yaml"]; !ok {
		t.Error("service-account.yaml not found")
		t.Log(results)
		t.FailNow()
	}

	if _, ok := results["service-app1.yaml"]; !ok {
		t.Error("service-app1.yaml not found")
		t.Log(results)
		t.FailNow()
	}

	if _, ok := results["deployment-app1.yaml"]; !ok {
		t.Error("deployment-app1.yaml not found")
		t.Log(results)
		t.FailNow()
	}
}

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
