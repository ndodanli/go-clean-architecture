// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "email": "ndodanli14@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/auth/user": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show an account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK. On success.",
                        "schema": {
                            "$ref": "#/definitions/res.SwaggerSuccessRes-httpctrl_GetUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request. On any validation error.",
                        "schema": {
                            "$ref": "#/definitions/res.SwaggerValidationErrRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized.",
                        "schema": {
                            "$ref": "#/definitions/res.SwaggerUnauthorizedErrRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error.",
                        "schema": {
                            "$ref": "#/definitions/res.SwaggerInternalErrRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "httpctrl.GetUserResponse": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "res.SwaggerInternalErrRes": {
            "type": "object",
            "properties": {
                "M": {
                    "type": "string",
                    "example": "Internal Server Error"
                },
                "S": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "res.SwaggerSuccessRes-httpctrl_GetUserResponse": {
            "type": "object",
            "properties": {
                "D": {
                    "$ref": "#/definitions/httpctrl.GetUserResponse"
                },
                "M": {
                    "type": "string",
                    "example": "XXX Created/Updated/Deleted Successfully"
                },
                "S": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "res.SwaggerUnauthorizedErrRes": {
            "type": "object",
            "properties": {
                "M": {
                    "type": "string",
                    "example": "Unauthorized"
                },
                "S": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "res.SwaggerValidationErrRes": {
            "type": "object",
            "properties": {
                "S": {
                    "type": "boolean",
                    "example": false
                },
                "V": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/res.ValidationError"
                    }
                }
            }
        },
        "res.ValidationError": {
            "type": "object",
            "properties": {
                "E": {
                    "type": "string",
                    "example": "Age must be greater than 0"
                },
                "F": {
                    "type": "string",
                    "example": "Age"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:5005",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Auth API",
	Description:      "This is a server for authentication and authorization",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
