package application

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"deploy-wizard/gen/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var templates = map[string][]string{
	"app":           {"service-account.yaml"},
	"services":      {"service.yaml"},
	"deployment":    {"deployment.yaml"},
	"kustomization": {"kustomization.yaml"},
}

var errTemplateUnreadableFormat = "the %q template must exist and be readable"

// Renderer is responsible for rendering manifests
type Renderer struct {
	templateDir string
}

// NewRenderer creates a new Renderer with the specified options
func NewRenderer(templateDir string) (*Renderer, error) {
	log.Infof("creating renderer with template directory %q", templateDir)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "template directory %q does not exist", templateDir)
	}

	// TODO: error if required templates do not exist?

	return &Renderer{templateDir}, nil
}

// BuildKustomization builds a kustomization file
func (r *Renderer) BuildKustomization(resources []string) (string, error) {
	data := struct {
		Resources []string
	}{Resources: resources}

	templateFile, err := templateFile(r.templateDir, templates["kustomization"][0])
	if err != nil {
		return "", errors.Wrapf(err, errTemplateUnreadableFormat)
	}

	result, err := renderTemplate(templateFile, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

// RenderManifests renders an application to individual Kubernetes manifest files
func (r *Renderer) RenderManifests(app *models.Application) (map[string]string, error) {
	manifests := map[string]string{}

	for _, tmpl := range templates["app"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return manifests, errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		log.Infof("rendering %q", templateFile)
		result, err := renderTemplate(templateFile, app)
		if err != nil {
			return manifests, err
		}
		manifests[tmpl] = result
	}

	serviceResults, err := r.renderServices(app)
	if err != nil {
		return manifests, err
	}
	for filename, content := range serviceResults {
		manifests[filename] = content
	}

	deploymentResults, err := r.renderDeployments(app)
	if err != nil {
		return manifests, err
	}
	for filename, content := range deploymentResults {
		manifests[filename] = content
	}

	var resources []string
	for filename := range manifests {
		resources = append(resources, filename)
	}

	log.Infof("manifest files: %v", resources)
	kustomizeFile, err := r.BuildKustomization(resources)
	if err != nil {
		return manifests, err
	}
	manifests[templates["kustomization"][0]] = kustomizeFile

	return manifests, nil
}

// RenderApplication renders an application to Kubernetes manifests
func (r *Renderer) RenderApplication(app *models.Application) (string, error) {
	var results []string
	manifests, err := r.RenderManifests(app)
	if err != nil {
		return "", err
	}

	// render in a specific order
	results = append(results, manifests["service-account.yaml"])
	for _, service := range app.Services {
		results = append(results, manifests[serviceName(service)])
		results = append(results, manifests[deploymentName(service)])
	}

	return strings.Join(results, "\n\n---\n"), nil
}

func (r *Renderer) renderServices(app *models.Application) (map[string]string, error) {
	manifests := map[string]string{}

	data := struct {
		App     *models.Application
		Service *models.Service
	}{App: app}

	for _, tmpl := range templates["services"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return manifests, errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		for _, service := range app.Services {
			log.Infof("rendering %q", templateFile)
			data.Service = service
			result, err := renderTemplate(templateFile, data)
			if err != nil {
				return manifests, err
			}
			manifests[serviceName(service)] = result
		}
	}

	return manifests, nil
}

func (r *Renderer) renderDeployments(app *models.Application) (map[string]string, error) {
	manifests := map[string]string{}

	data := struct {
		App        *models.Application
		Service    *models.Service
		Containers []*models.Container
	}{App: app}

	for _, tmpl := range templates["deployment"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return manifests, errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		for _, service := range app.Services {
			log.Infof("rendering %q", templateFile)
			data.Service = service
			data.Containers = service.Containers
			result, err := renderTemplate(templateFile, data)
			if err != nil {
				return manifests, err
			}
			manifests[deploymentName(service)] = result
		}
	}

	return manifests, nil
}

// RenderTemplate renders the specified template with the Application model
func renderTemplate(name string, obj interface{}) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read template %q", name)
	}

	t, err := template.New("resources").Parse(string(data))
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse template %q", name)
	}

	var rendered bytes.Buffer
	err = t.Execute(&rendered, obj)
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute template %q", name)
	}

	return rendered.String(), nil
}

func templateFile(templateDir, name string) (string, error) {
	filename := path.Join(templateDir, name)

	_, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func serviceName(s *models.Service) string {
	return fmt.Sprintf("service-%s.yaml", s.Name)
}

func deploymentName(s *models.Service) string {
	return fmt.Sprintf("deployment-%s.yaml", s.Name)
}
