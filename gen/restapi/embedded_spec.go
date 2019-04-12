// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Kruise deployment wizard",
    "title": "Kruise deployment wizard",
    "version": "0.0.1"
  },
  "paths": {
    "/apps": {
      "post": {
        "description": "Generates a new Kruise application",
        "produces": [
          "text/plain"
        ],
        "tags": [
          "apps"
        ],
        "operationId": "createApp",
        "parameters": [
          {
            "description": "The application to create",
            "name": "application",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/application"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "created",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "$ref": "#/responses/BadRequest"
          },
          "default": {
            "$ref": "#/responses/InternalServerError"
          }
        }
      }
    },
    "/health": {
      "get": {
        "description": "Get the current health of the API",
        "tags": [
          "general"
        ],
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "Get current health response",
            "schema": {
              "$ref": "#/definitions/healthStatus"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/validates/application": {
      "post": {
        "description": "validates the application details",
        "tags": [
          "validations"
        ],
        "operationId": "validateApplication",
        "parameters": [
          {
            "description": "the application object",
            "name": "application",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "validated",
            "schema": {
              "$ref": "#/definitions/validationResponse"
            }
          },
          "400": {
            "$ref": "#/responses/BadRequest"
          },
          "default": {
            "$ref": "#/responses/InternalServerError"
          }
        }
      }
    }
  },
  "definitions": {
    "application": {
      "type": "object",
      "required": [
        "name",
        "release",
        "tenant",
        "environment",
        "region",
        "namespace",
        "repoURL"
      ],
      "properties": {
        "environment": {
          "description": "The environment to deploy to",
          "type": "string",
          "default": "Dev",
          "minLength": 1,
          "enum": [
            "Dev",
            "Stage",
            "Prod"
          ],
          "x-nullable": false
        },
        "ingress": {
          "$ref": "#/definitions/ingress"
        },
        "name": {
          "description": "The name of the application",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "namespace": {
          "description": "The namespace to deploy to",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "path": {
          "description": "The relative path to the manifests in the git repo",
          "type": "string",
          "format": "filepath",
          "default": "/",
          "minLength": 1,
          "x-nullable": false
        },
        "region": {
          "description": "The region to deploy to",
          "type": "string",
          "default": "STL",
          "minLength": 1,
          "enum": [
            "STL",
            "KCI",
            "BEL"
          ],
          "x-nullable": false
        },
        "release": {
          "description": "The version or release name of the application",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "repoURL": {
          "description": "The git repo URL",
          "type": "string",
          "format": "uri",
          "minLength": 1,
          "x-nullable": false
        },
        "services": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/service"
          }
        },
        "targetRevision": {
          "description": "Defines the commit, tag, or branch in which to sync the application to.",
          "type": "string",
          "default": "HEAD",
          "minLength": 1,
          "x-nullable": false
        },
        "tenant": {
          "description": "The name of the tenant",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "container": {
      "type": "object",
      "properties": {
        "command": {
          "description": "The command to run for the docker image's entrypoint.",
          "type": "string",
          "x-nullable": true
        },
        "image": {
          "description": "The docker image name for the container",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "imagePullPolicy": {
          "description": "Image pull policy. One of Always or IfNotPresent.",
          "type": "string",
          "default": "IfNotPresent",
          "minLength": 1,
          "enum": [
            "Always",
            "IfNotPresent"
          ],
          "x-nullable": false
        },
        "imageTag": {
          "description": "The docker image tag for the container",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "name": {
          "description": "The name of this container within the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "healthStatus": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string"
        }
      }
    },
    "ingress": {
      "type": "object",
      "required": [
        "name",
        "rules"
      ],
      "properties": {
        "name": {
          "description": "The name of the ingress",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "rules": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ingressRule"
          }
        }
      }
    },
    "ingressRule": {
      "type": "object",
      "required": [
        "host",
        "serviceName",
        "servicePort"
      ],
      "properties": {
        "host": {
          "description": "Host is the fully qualified domain name of a network host, as defined by RFC 3986",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "path": {
          "description": "Path is an extended POSIX regex as defined by IEEE Std 1003.1, (i.e this follows the egrep/unix syntax, not the perl syntax) matched against the path of an incoming request",
          "type": "string"
        },
        "serviceName": {
          "description": "Specifies the name of the referenced service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "servicePort": {
          "description": "Specifies the port of the referenced service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "service": {
      "type": "object",
      "required": [
        "name",
        "type",
        "ports"
      ],
      "properties": {
        "containers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/container"
          }
        },
        "name": {
          "description": "The name of the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "ports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/servicePort"
          }
        },
        "tier": {
          "description": "The tier for the service",
          "type": "string",
          "enum": [
            "Frontend",
            "API",
            "Backend",
            "Cache"
          ],
          "x-nullable": true
        },
        "type": {
          "description": "The service type",
          "type": "string",
          "enum": [
            "ClusterIP",
            "ExternalName",
            "LoadBalancer"
          ],
          "x-nullable": false
        }
      }
    },
    "servicePort": {
      "type": "object",
      "required": [
        "name",
        "port"
      ],
      "properties": {
        "name": {
          "description": "The name of this port within the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "port": {
          "description": "The port that will be exposed by this service",
          "type": "integer",
          "x-nullable": false
        },
        "protocol": {
          "description": "The IP protocol for this port. Supports \"TCP\" and \"UDP\". Default is TCP",
          "type": "string",
          "default": "TCP",
          "minLength": 1,
          "enum": [
            "TCP",
            "UDP"
          ],
          "x-nullable": false
        },
        "targetPort": {
          "description": "Number or name of the port to access on the pods targeted by the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "validationResponse": {
      "type": "object",
      "properties": {
        "errors": {
          "type": "object",
          "additionalProperties": {
            "type": "object"
          }
        }
      }
    }
  },
  "responses": {
    "BadRequest": {
      "description": "Bad request",
      "schema": {
        "type": "string"
      }
    },
    "InternalServerError": {
      "description": "Internal server error",
      "schema": {
        "type": "string"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Kruise deployment wizard",
    "title": "Kruise deployment wizard",
    "version": "0.0.1"
  },
  "paths": {
    "/apps": {
      "post": {
        "description": "Generates a new Kruise application",
        "produces": [
          "text/plain"
        ],
        "tags": [
          "apps"
        ],
        "operationId": "createApp",
        "parameters": [
          {
            "description": "The application to create",
            "name": "application",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/application"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "created",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Internal server error",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/health": {
      "get": {
        "description": "Get the current health of the API",
        "tags": [
          "general"
        ],
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "Get current health response",
            "schema": {
              "$ref": "#/definitions/healthStatus"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/validates/application": {
      "post": {
        "description": "validates the application details",
        "tags": [
          "validations"
        ],
        "operationId": "validateApplication",
        "parameters": [
          {
            "description": "the application object",
            "name": "application",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "validated",
            "schema": {
              "$ref": "#/definitions/validationResponse"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Internal server error",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "application": {
      "type": "object",
      "required": [
        "name",
        "release",
        "tenant",
        "environment",
        "region",
        "namespace",
        "repoURL"
      ],
      "properties": {
        "environment": {
          "description": "The environment to deploy to",
          "type": "string",
          "default": "Dev",
          "minLength": 1,
          "enum": [
            "Dev",
            "Stage",
            "Prod"
          ],
          "x-nullable": false
        },
        "ingress": {
          "$ref": "#/definitions/ingress"
        },
        "name": {
          "description": "The name of the application",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "namespace": {
          "description": "The namespace to deploy to",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "path": {
          "description": "The relative path to the manifests in the git repo",
          "type": "string",
          "format": "filepath",
          "default": "/",
          "minLength": 1,
          "x-nullable": false
        },
        "region": {
          "description": "The region to deploy to",
          "type": "string",
          "default": "STL",
          "minLength": 1,
          "enum": [
            "STL",
            "KCI",
            "BEL"
          ],
          "x-nullable": false
        },
        "release": {
          "description": "The version or release name of the application",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "repoURL": {
          "description": "The git repo URL",
          "type": "string",
          "format": "uri",
          "minLength": 1,
          "x-nullable": false
        },
        "services": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/service"
          }
        },
        "targetRevision": {
          "description": "Defines the commit, tag, or branch in which to sync the application to.",
          "type": "string",
          "default": "HEAD",
          "minLength": 1,
          "x-nullable": false
        },
        "tenant": {
          "description": "The name of the tenant",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "container": {
      "type": "object",
      "properties": {
        "command": {
          "description": "The command to run for the docker image's entrypoint.",
          "type": "string",
          "x-nullable": true
        },
        "image": {
          "description": "The docker image name for the container",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "imagePullPolicy": {
          "description": "Image pull policy. One of Always or IfNotPresent.",
          "type": "string",
          "default": "IfNotPresent",
          "minLength": 1,
          "enum": [
            "Always",
            "IfNotPresent"
          ],
          "x-nullable": false
        },
        "imageTag": {
          "description": "The docker image tag for the container",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "name": {
          "description": "The name of this container within the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "healthStatus": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string"
        }
      }
    },
    "ingress": {
      "type": "object",
      "required": [
        "name",
        "rules"
      ],
      "properties": {
        "name": {
          "description": "The name of the ingress",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "rules": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ingressRule"
          }
        }
      }
    },
    "ingressRule": {
      "type": "object",
      "required": [
        "host",
        "serviceName",
        "servicePort"
      ],
      "properties": {
        "host": {
          "description": "Host is the fully qualified domain name of a network host, as defined by RFC 3986",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "path": {
          "description": "Path is an extended POSIX regex as defined by IEEE Std 1003.1, (i.e this follows the egrep/unix syntax, not the perl syntax) matched against the path of an incoming request",
          "type": "string"
        },
        "serviceName": {
          "description": "Specifies the name of the referenced service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "servicePort": {
          "description": "Specifies the port of the referenced service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "service": {
      "type": "object",
      "required": [
        "name",
        "type",
        "ports"
      ],
      "properties": {
        "containers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/container"
          }
        },
        "name": {
          "description": "The name of the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "ports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/servicePort"
          }
        },
        "tier": {
          "description": "The tier for the service",
          "type": "string",
          "enum": [
            "Frontend",
            "API",
            "Backend",
            "Cache"
          ],
          "x-nullable": true
        },
        "type": {
          "description": "The service type",
          "type": "string",
          "enum": [
            "ClusterIP",
            "ExternalName",
            "LoadBalancer"
          ],
          "x-nullable": false
        }
      }
    },
    "servicePort": {
      "type": "object",
      "required": [
        "name",
        "port"
      ],
      "properties": {
        "name": {
          "description": "The name of this port within the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        },
        "port": {
          "description": "The port that will be exposed by this service",
          "type": "integer",
          "x-nullable": false
        },
        "protocol": {
          "description": "The IP protocol for this port. Supports \"TCP\" and \"UDP\". Default is TCP",
          "type": "string",
          "default": "TCP",
          "minLength": 1,
          "enum": [
            "TCP",
            "UDP"
          ],
          "x-nullable": false
        },
        "targetPort": {
          "description": "Number or name of the port to access on the pods targeted by the service",
          "type": "string",
          "minLength": 1,
          "x-nullable": false
        }
      }
    },
    "validationResponse": {
      "type": "object",
      "properties": {
        "errors": {
          "type": "object",
          "additionalProperties": {
            "type": "object"
          }
        }
      }
    }
  },
  "responses": {
    "BadRequest": {
      "description": "Bad request",
      "schema": {
        "type": "string"
      }
    },
    "InternalServerError": {
      "description": "Internal server error",
      "schema": {
        "type": "string"
      }
    }
  }
}`))
}
