package application

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"deploy-wizard/gen/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

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

	return &Renderer{templateDir}, nil
}

// RenderApplication renders an application to Kubernetes manifests
func (r *Renderer) RenderApplication(app *models.Application) (string, error) {
	pattern := filepath.Join(r.templateDir, "*.y[a]ml")
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return "", errors.Wrapf(err, "failed listing templates in %q", r.templateDir)
	}

	var result string
	for _, tmpl := range templates {
		templateName := filepath.Base(tmpl)
		if filepath.Base(tmpl) == "service-account.yaml" {
			log.Infof("rendering %q", templateName)
			result, err = renderTemplate(tmpl, app)
			if err != nil {
				return "", err
			}
		}
	}

	return result, nil
}

// RenderTemplate renders the specified template with the Application model
func renderTemplate(name string, app *models.Application) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read template %q", name)
	}

	t := template.Must(template.New("resources").Parse(string(data)))
	var rendered bytes.Buffer
	err = t.Execute(&rendered, app)
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute template %q", name)
	}

	return rendered.String(), nil
}
