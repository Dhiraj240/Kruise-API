package application

import (
	"bytes"
	"deploy-wizard/gen/models"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	errMsgInvalidJSON      = "invalid json payload"
	errMsgNotAnApplication = "not an application object"
)

var (
	regexDNSName = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)
)

// ValidateApplication returns of map with key = field and value = error
func ValidateApplication(appdata interface{}) map[string]interface{} {
	errors := map[string]interface{}{}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(appdata)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgInvalidJSON)
		errors[""] = errMsgInvalidJSON
		return errors

	}

	var app *models.Application
	err = json.Unmarshal(b.Bytes(), &app)
	if err != nil {
		log.WithField("f", "application.ValidateApplication").WithError(err).Warnf(errMsgNotAnApplication)
		errors[""] = errMsgNotAnApplication
		return errors
	}

	if app.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if app.Release == "" {
		errors["release"] = newRequiredValidationError("release")
	}

	if app.Environment == "" {
		errors["environment"] = newRequiredValidationError("environment")
	}

	if app.Tenant == "" {
		errors["tenant"] = newRequiredValidationError("tenant")
	}

	if app.Namespace == "" {
		errors["namespace"] = newRequiredValidationError("namespace")
	}

	if app.Path == "" {
		errors["path"] = newRequiredValidationError("path")
	}

	if app.Region == "" {
		errors["region"] = newRequiredValidationError("region")
	}

	if app.RepoURL == "" {
		errors["repoURL"] = newRequiredValidationError("repoURL")
	}

	if app.TargetRevision == "" {
		errors["targetRevision"] = newRequiredValidationError("targetRevision")
	}

	if len(app.Services) > 0 {
		servicesErrors := ValidateServices(app.Services)
		if len(servicesErrors) > 0 {
			errors["services"] = servicesErrors
		}
	}

	if app.Ingress != nil {
		ingressErrors := ValidateIngress(app.Ingress, app.Services)
		if len(ingressErrors) > 0 {
			errors["ingress"] = ingressErrors
		}
	}

	return errors
}

// ValidateServices returns of map with key = field and value = error
func ValidateServices(services []*models.Service) map[string]interface{} {
	errors := map[string]interface{}{}

	for i, svc := range services {
		svcErrors := ValidateService(svc)
		idx := strconv.Itoa(i)

		if len(svcErrors) > 0 {
			errors[idx] = svcErrors
		}
	}

	return errors
}

// ValidateService returns of map with key = field and value = error
func ValidateService(svc *models.Service) map[string]interface{} {
	errors := map[string]interface{}{}

	if svc.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if svc.Type == "" {
		errors["type"] = newRequiredValidationError("type")
	}

	if len(svc.Ports) == 0 {
		errors["ports"] = newRequiredValidationError("ports")
		return errors
	}

	portsErrors := ValidateServicePorts(svc.Ports)
	if len(portsErrors) > 0 {
		errors["ports"] = portsErrors
	}

	containerErrors := ValidateContainers(svc.Containers)
	if len(containerErrors) > 0 {
		errors["containers"] = containerErrors
	}

	return errors
}

// ValidateContainers returns of map with key = field and value = error
func ValidateContainers(containers []*models.Container) map[string]interface{} {
	errors := map[string]interface{}{}

	for i, container := range containers {
		containerErrors := ValidateContainer(container)

		if len(containerErrors) > 0 {
			errors[strconv.Itoa(i)] = containerErrors
		}
	}

	return errors
}

// ValidateContainer returns of map with key = field and value = error
func ValidateContainer(container *models.Container) map[string]interface{} {
	errors := map[string]interface{}{}

	if container.Name == "" {
		errors["name"] = newRequiredValidationError("container")
	}

	if container.Image == "" {
		errors["image"] = newRequiredValidationError("image")
	}

	return errors
}

// ValidateServicePorts returns of map with key = field and value = error
func ValidateServicePorts(ports []*models.ServicePort) map[string]interface{} {
	errors := map[string]interface{}{}

	for i, port := range ports {
		portErrors := ValidateServicePort(port)

		if len(portErrors) > 0 {
			errors[strconv.Itoa(i)] = portErrors
		}
	}

	return errors
}

// ValidateServicePort returns of map with key = field and value = error
func ValidateServicePort(port *models.ServicePort) map[string]interface{} {
	errors := map[string]interface{}{}

	if port.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if port.Port == 0 {
		errors["port"] = newRequiredValidationError("port")
	}

	return errors
}

// ValidateIngress returns of map with key = field and value = error
func ValidateIngress(ingress *models.Ingress, services []*models.Service) map[string]interface{} {
	errors := map[string]interface{}{}

	if ingress.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	rulesErrors := ValidateIngressRules(ingress.Rules, services)
	if len(rulesErrors) > 0 {
		errors["rules"] = rulesErrors
	}

	return errors
}

// ValidateIngressRules returns of map with key = field and value = error
func ValidateIngressRules(ingressRules []*models.IngressRule, services []*models.Service) map[string]interface{} {
	errors := map[string]interface{}{}

	for i, ingressRule := range ingressRules {
		ruleErrors := ValidateIngressRule(ingressRule, services)
		idx := strconv.Itoa(i)

		if len(ruleErrors) > 0 {
			errors[idx] = ruleErrors
		}
	}

	return errors
}

// ValidateIngressRule returns of map with key = field and value = error
func ValidateIngressRule(ingressRule *models.IngressRule, services []*models.Service) map[string]interface{} {
	errors := map[string]interface{}{}

	if ingressRule.Host == "" {
		errors["host"] = newRequiredValidationError("host")
	} else if !isValidDNSName(ingressRule.Host) {
		errors["host"] = fmt.Sprintf("%q must be a valid host name", ingressRule.Host)
	}

	var backendService *models.Service

	if ingressRule.ServiceName == "" {
		errors["serviceName"] = newRequiredValidationError("serviceName")
	} else {
		for _, service := range services {
			if service.Name == ingressRule.ServiceName {
				backendService = service
				break
			}
		}
		if backendService == nil {
			errors["serviceName"] = fmt.Sprintf("%q does not match an exisiting service", ingressRule.ServiceName)
		}
	}

	if ingressRule.ServicePort == "" {
		errors["servicePort"] = newRequiredValidationError("servicePort")
	} else {
		if backendService != nil {
			var validPort bool
			for _, port := range backendService.Ports {
				if port.Name == ingressRule.ServicePort {
					validPort = true
					break
				}
			}

			if !validPort {
				errors["servicePort"] = fmt.Sprintf("%q does not match a service port for %q", ingressRule.ServicePort, backendService.Name)
			}
		}
	}

	return errors
}

func newRequiredValidationError(field string) string {
	return fmt.Sprintf("%q is a required field", field)
}

func isValidDNSName(host string) bool {
	if host == "" || len(strings.Replace(host, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !isIP(host) && regexDNSName.MatchString(host)
}

func isIP(host string) bool {
	return net.ParseIP(host) != nil
}
