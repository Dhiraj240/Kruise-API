package application

import (
	"bytes"
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
	"app":        {"service-account.yaml"},
	"services":   {"service.yaml"},
	"deployment": {"deployment.yaml"},
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

// RenderApplication renders an application to Kubernetes manifests
func (r *Renderer) RenderApplication(app *models.Application) (string, error) {
	var results []string
	for _, tmpl := range templates["app"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return "", errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		log.Infof("rendering %q", templateFile)
		result, err := renderTemplate(templateFile, app)
		if err != nil {
			return "", err
		}
		results = append(results, result)
	}

	serviceResults, err := r.renderServices(app)
	if err != nil {
		return "", err
	}
	results = append(results, serviceResults...)

	deploymentResults, err := r.renderDeployments(app)
	if err != nil {
		return "", err
	}
	results = append(results, deploymentResults...)

	return strings.Join(results, "\n\n---\n"), nil
}

func (r *Renderer) renderServices(app *models.Application) ([]string, error) {
	var results []string

	data := struct {
		App     *models.Application
		Service *models.Service
	}{App: app}

	for _, tmpl := range templates["services"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return []string{}, errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		for _, service := range app.Services {
			log.Infof("rendering %q", templateFile)
			data.Service = service
			result, err := renderTemplate(templateFile, data)
			if err != nil {
				return []string{}, err
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func (r *Renderer) renderDeployments(app *models.Application) ([]string, error) {
	var results []string

	data := struct {
		App       *models.Application
		Service   *models.Service
		Container *models.Container
	}{App: app}

	for _, tmpl := range templates["deployment"] {
		templateFile, err := templateFile(r.templateDir, tmpl)
		if err != nil {
			return []string{}, errors.Wrapf(err, errTemplateUnreadableFormat)
		}

		for _, service := range app.Services {
			log.Infof("rendering %q", templateFile)
			data.Service = service
			for _, container := range service.Containers {
				data.Container = container
				result, err := renderTemplate(templateFile, data)
				if err != nil {
					return []string{}, err
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
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
