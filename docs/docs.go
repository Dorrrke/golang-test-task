// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/test_task/api/user": {
            "post": {
                "description": "Add user in data base",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Add user",
                "operationId": "add-user",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "write data in db error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/test_task/api/user/{id}": {
            "get": {
                "description": "Get user from db",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user",
                "operationId": "get-user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error encoding to json",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/test_task/api/users": {
            "get": {
                "description": "Get all user from db",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all user",
                "operationId": "all-users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "204": {
                        "description": "no users in db",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error encoding to json",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.User": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "occupation": {
                    "type": "string"
                },
                "salary": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/test_task/api",
	Schemes:          []string{},
	Title:            "Blueprint Swagger API",
	Description:      "Swagger API for Golang Project Blueprint.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}