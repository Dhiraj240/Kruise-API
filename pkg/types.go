package types

// Application configures an application deployment
type Application struct {
	Name           string    `json:"name"`
	Tenant         string    `json:"tenant"`
	Environment    string    `json:"environment"`
	Region         string    `json:"region"`
	Namespace      string    `json:"namespace"`
	RepoURL        string    `json:"repoUrl"`
	Path           string    `json:"path"`
	TargetRevision string    `json:"targetRevision"`
	Services       []Service `json:"services"`
}

// Service configures a service resource
type Service struct {
	Name  string        `json:"name"`
	Tier  string        `json:"tier"`
	Ports []ServicePort `json:"ports"`
}

// ServicePort configures a service port
type ServicePort struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	TargetPort string `json:"targetPort"`
}

// Ingress configures an ingress resource
type Ingress struct {
	Name  string        `json:"name"`
	Rules []IngressRule `json:"rules"`
}

// IngressRule configures an ingress rule
type IngressRule struct {
	Host        string `json:"host"`
	Path        string `json:"path"`
	ServiceName string `json:"serviceName"`
	ServicePort string `json:"servicePort"`
}
