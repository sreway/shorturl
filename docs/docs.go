// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "a.y.oleynik@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "post": {
                "description": "add short URL",
                "produces": [
                    "text/plain"
                ],
                "summary": "add short URL",
                "operationId": "addURL",
                "parameters": [
                    {
                        "description": "long URL to shorten",
                        "name": "longURL",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        },
        "/api/shorten": {
            "post": {
                "description": "create short URL",
                "produces": [
                    "application/json"
                ],
                "summary": "create short URL",
                "operationId": "shortURL",
                "parameters": [
                    {
                        "description": "long URL to shorten",
                        "name": "longURL",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.shortURLRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/http.shortURLResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        },
        "/api/shorten/batch": {
            "post": {
                "description": "create of several short URLs",
                "produces": [
                    "application/json"
                ],
                "summary": "create of several short URLs",
                "operationId": "batchURL",
                "parameters": [
                    {
                        "description": "several long URLs to shorten",
                        "name": "bathURL",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/http.batchURLRequest"
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/http.batchURLResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        },
        "/api/shorten/user/urls": {
            "get": {
                "description": "get short URLs for user ID",
                "produces": [
                    "application/json"
                ],
                "summary": "get short URLs for user ID",
                "operationId": "userURL",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/http.userURLResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "remove multiple short URLs",
                "produces": [
                    "application/json"
                ],
                "summary": "remove multiple short URLs",
                "operationId": "deleteURL",
                "parameters": [
                    {
                        "description": "short URL ids to delete",
                        "name": "ids",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "410": {
                        "description": "Gone",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "health check shortener storage",
                "summary": "health check shortener storage",
                "operationId": "ping",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        },
        "/{id}": {
            "post": {
                "description": "get short URL",
                "produces": [
                    "text/plain"
                ],
                "summary": "get short URL",
                "operationId": "getURL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "short URL id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/http.errResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.batchURLRequest": {
            "type": "object",
            "properties": {
                "correlation_id": {
                    "type": "string"
                },
                "original_url": {
                    "type": "string"
                }
            }
        },
        "http.batchURLResponse": {
            "type": "object",
            "properties": {
                "correlation_id": {
                    "type": "string"
                },
                "short_url": {
                    "type": "string"
                }
            }
        },
        "http.errResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "http.shortURLRequest": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "http.shortURLResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "string"
                }
            }
        },
        "http.userURLResponse": {
            "type": "object",
            "properties": {
                "original_url": {
                    "type": "string"
                },
                "short_url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
