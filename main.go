package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"text/template"
)

func main() {
	log.Println("starting deploy-wizard")

	tmpl, err := loadTemplate("ingress-fanout.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"applicationName": "molly",
		"services": []map[string]string{
			{
				"hostFQDN":    "molly-data.mc.int",
				"serviceName": "molly-data",
				"servicePort": "8080",
				"servicePath": "/",
			},
			{
				"hostFQDN":    "molly-storage.mc.int",
				"serviceName": "molly-storage",
				"servicePort": "9801",
			},
		},
	}
	rendered, err := renderTemplate(tmpl, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(rendered)
}

func renderTemplate(tmpl []byte, data map[string]interface{}) (string, error) {
	t := template.Must(template.New("resources").Parse(string(tmpl)))
	var rendered bytes.Buffer
	err := t.Execute(&rendered, data)
	if err != nil {
		return "", err
	}

	return rendered.String(), nil
}

func loadTemplate(name string) ([]byte, error) {
	templateFile := path.Join("./_templates", name)
	return ioutil.ReadFile(templateFile)
}
