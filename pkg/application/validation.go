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

// TODO: validate integers are not 0 (port, targetPort, capacity)

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

	if app.Metadata == nil {
		errors[""] = newRequiredValidationError("metadata")
		return errors
	}

	if app.Spec == nil {
		errors[""] = newRequiredValidationError("spec")
		return errors
	}

	mdErrors := ValidateMetadata(app.Metadata)
	if len(mdErrors) > 0 {
		errors["metadata"] = mdErrors
	}

	specErrors := ValidateSpec(app.Spec)
	if len(specErrors) > 0 {
		errors["spec"] = specErrors
	}

	return errors
}

// ValidateMetadata returns of map with key = field and value = error
func ValidateMetadata(md *models.Metadata) map[string]interface{} {
	errors := map[string]interface{}{}
	if md.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if md.Namespace == "" {
		errors["namespace"] = newRequiredValidationError("namespace")
	}

	lblErrors := ValidateLabels(md.Labels)
	if len(lblErrors) > 0 {
		errors["labels"] = lblErrors
	}

	return errors
}

// ValidateLabels returns of map with key = field and value = error
func ValidateLabels(labels *models.Labels) map[string]interface{} {
	errors := map[string]interface{}{}
	if labels.Env == "" {
		errors["environment"] = newRequiredValidationError("environment")
	}

	if labels.Team == "" {
		errors["tenant"] = newRequiredValidationError("tenant")
	}

	if labels.Version == "" {
		errors["version"] = newRequiredValidationError("version")
	}

	if labels.Region == "" {
		errors["region"] = newRequiredValidationError("region")
	}

	return errors
}

// ValidateDestination returns of map with key = field and value = error
func ValidateDestination(dest *models.Destination) map[string]interface{} {
	errors := map[string]interface{}{}

	if dest.URL == "" {
		errors["url"] = newRequiredValidationError("url")
	}

	if dest.Path == "" {
		errors["path"] = newRequiredValidationError("path")
	}

	if dest.TargetRevision == "" {
		errors["targetRevision"] = newRequiredValidationError("targetRevision")
	}

	return errors
}

// ValidateSpec returns of map with key = field and value = error
func ValidateSpec(spec *models.Spec) map[string]interface{} {
	errors := map[string]interface{}{}
	if spec.Destination == nil {
		errors["destination"] = newRequiredValidationError("destination")
		return errors
	}
	if verrs := ValidateDestination(spec.Destination); len(verrs) > 0 {
		errors["destination"] = verrs
	}

	if len(spec.Components) == 0 {
		errors["components"] = newRequiredValidationError("components")
		return errors
	}

	if verrs := ValidateComponents(spec.Components); len(verrs) > 0 {
		errors["components"] = verrs
	}

	if verrs := ValidateConfigMaps(spec.ConfigMaps); len(verrs) > 0 {
		errors["configMaps"] = verrs
	}
	if verrs := ValidatePersistentVolumes(spec.PersistentVolumes); len(verrs) > 0 {
		errors["persistentVolumes"] = verrs
	}

	return errors
}

// ValidateConfigMaps returns of map with key = field and value = error
func ValidateConfigMaps(configMaps []*models.ConfigMap) map[string]interface{} {
	errors := map[string]interface{}{}
	for i, comp := range configMaps {
		errs := ValidateConfigMap(comp)
		idx := strconv.Itoa(i)
		if len(errs) > 0 {
			errors[idx] = errs
		}
	}
	return errors
}

// ValidateConfigMap returns of map with key = field and value = error
func ValidateConfigMap(configMap *models.ConfigMap) map[string]interface{} {
	errors := map[string]interface{}{}

	if configMap.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if configMap.Data == "" {
		errors["data"] = newRequiredValidationError("data")
	}

	return errors
}

// ValidatePersistentVolumes returns of map with key = field and value = error
func ValidatePersistentVolumes(persistentVolumes []*models.PersistentVolume) map[string]interface{} {
	errors := map[string]interface{}{}
	for i, comp := range persistentVolumes {
		errs := ValidatePersistentVolume(comp)
		idx := strconv.Itoa(i)
		if len(errs) > 0 {
			errors[idx] = errs
		}
	}
	return errors
}

// ValidatePersistentVolume returns of map with key = field and value = error
func ValidatePersistentVolume(persistentVolume *models.PersistentVolume) map[string]interface{} {
	errors := map[string]interface{}{}

	if persistentVolume.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}

	if persistentVolume.AccessMode == "" {
		errors["accessMode"] = newRequiredValidationError("accessMode")
	}

	if persistentVolume.Capacity <= 0 {
		errors["capacity"] = "capacity must be greater than 0"
	}

	if persistentVolume.StorageClassName == "" {
		errors["storageClassName"] = newRequiredValidationError("storageClassName")
	}

	return errors
}

// ValidateComponent returns of map with key = field and value = error
func ValidateComponent(component *models.Component) map[string]interface{} {
	errors := map[string]interface{}{}

	if verrs := ValidateService(component.Service); len(verrs) > 0 {
		errors["service"] = verrs
	}

	if verrs := ValidateIngresses(component.Ingresses); len(verrs) > 0 {
		errors["ingresses"] = verrs
	}

	if verrs := ValidateContainers(component.Containers); len(verrs) > 0 {
		errors["containers"] = verrs
	}

	return errors
}

// ValidateComponents returns of map with key = field and value = error
func ValidateComponents(components []*models.Component) map[string]interface{} {
	errors := map[string]interface{}{}
	for i, comp := range components {
		errs := ValidateComponent(comp)
		idx := strconv.Itoa(i)
		if len(errs) > 0 {
			errors[idx] = errs
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
		errors["name"] = newRequiredValidationError("name")
	}

	if container.Image == "" {
		errors["image"] = newRequiredValidationError("image")
	}

	if container.ImageTag == "" {
		errors["imageTag"] = newRequiredValidationError("imageTag")
	}

	if len(container.PortNames) == 0 {
		errors["portNames"] = newRequiredValidationError("portNames")
	}

	if len(container.Volumes) > 0 {
		if verrs := ValidateVolumeMounts(container.Volumes); len(verrs) > 0 {
			errors["volumes"] = verrs
		}
	}

	return errors
}

// ValidateVolumeMounts returns of map with key = field and value = error
func ValidateVolumeMounts(mounts []*models.VolumeMount) map[string]interface{} {
	errors := map[string]interface{}{}
	for i, vm := range mounts {
		if verrs := ValidateVolumeMount(vm); len(verrs) > 0 {
			errors[strconv.Itoa(i)] = verrs
		}
	}
	return errors
}

// ValidateVolumeMount returns of map with key = field and value = error
func ValidateVolumeMount(mount *models.VolumeMount) map[string]interface{} {
	errors := map[string]interface{}{}
	if mount.Name == "" {
		errors["name"] = newRequiredValidationError("name")
	}
	if mount.Type == "" {
		errors["type"] = newRequiredValidationError("type")
	}
	if mount.MountPath == "" {
		errors["mountPath"] = newRequiredValidationError("mountPath")
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
		errors["port"] = fmt.Sprintf("%d is not a valid port number", port.Port)
	}

	return errors
}

// ValidateIngress returns of map with key = field and value = error
func ValidateIngresses(ingresses []*models.Ingress) map[string]interface{} {
	errors := map[string]interface{}{}
	for i, ingress := range ingresses {
		if verrs := ValidateIngress(ingress); len(verrs) > 0 {
			errors[strconv.Itoa(i)] = verrs
		}
	}
	return errors
}

// ValidateIngress returns of map with key = field and value = error
func ValidateIngress(ingress *models.Ingress) map[string]interface{} {
	errors := map[string]interface{}{}

	if ingress.Host == "" {
		errors["host"] = newRequiredValidationError("host")
	}

	if !isValidDNSName(ingress.Host) {
		errors["host"] = fmt.Sprintf("%q must be a valid host name", ingress.Host)
	}

	if verrs := ValidateIngressPaths(ingress.Paths); len(verrs) > 0 {
		errors["paths"] = verrs
	}

	return errors
}

// ValidateIngressPaths returns of map with key = field and value = error
func ValidateIngressPaths(ingressPaths []*models.IngressPath) map[string]interface{} {
	errors := map[string]interface{}{}

	for i, ingressPath := range ingressPaths {
		pathErrors := ValidateIngressPath(ingressPath)
		idx := strconv.Itoa(i)

		if len(pathErrors) > 0 {
			errors[idx] = pathErrors
		}
	}

	return errors
}

// ValidateIngressPath returns of map with key = field and value = error
func ValidateIngressPath(ingressPath *models.IngressPath) map[string]interface{} {
	errors := map[string]interface{}{}

	if ingressPath.Path == "" {
		errors["path"] = newRequiredValidationError("path")
	}

	if ingressPath.PortName == "" {
		errors["portName"] = newRequiredValidationError("portName")
	}

	// TODO: portName matches one of service port names

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
