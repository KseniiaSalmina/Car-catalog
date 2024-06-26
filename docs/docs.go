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
        "/cars": {
            "get": {
                "description": "get cars by filters. Filters are accepted in json format as a car structure",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Get cars",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit records by page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "=",
                            "\u003e=",
                            "\u003c="
                        ],
                        "type": "string",
                        "description": "can be =, \u003e=, \u003c=. By default =",
                        "name": "yearFilterMode",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "registration number of the car",
                        "name": "regNum",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "mark of the car",
                        "name": "mark",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "model of the car",
                        "name": "model",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "year of the car",
                        "name": "year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name of owner of the car",
                        "name": "ownerName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "surname of owner of the car",
                        "name": "ownerSurname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "to search for car owners without a patronymic, the patronymic field must contain the string ",
                        "name": "ownerPatronymic",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ASC",
                            "DESC"
                        ],
                        "type": "string",
                        "description": "relating to the car year. Can be ASC or DESC, by default DESC",
                        "name": "orderByMode",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CarsPage"
                        }
                    },
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "get info about cars from outsider service and put it to the database",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Post cars",
                "parameters": [
                    {
                        "description": "array of regNums",
                        "name": "regNums",
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
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cars/{regNum}": {
            "delete": {
                "description": "delete car by its registration number",
                "tags": [
                    "cars"
                ],
                "summary": "Delete car",
                "parameters": [
                    {
                        "type": "string",
                        "description": "registration number of the car",
                        "name": "regNum",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "update car info by accepted json",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Patch car",
                "parameters": [
                    {
                        "type": "string",
                        "description": "registration number of the car",
                        "name": "regNum",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "car with updated info",
                        "name": "car",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Car"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Car": {
            "type": "object",
            "properties": {
                "mark": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/models.Person"
                },
                "regNum": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "models.CarsPage": {
            "type": "object",
            "properties": {
                "cars": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Car"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "page_number": {
                    "type": "integer"
                },
                "pages_amount": {
                    "type": "integer"
                }
            }
        },
        "models.Person": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8088",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Car catalog",
	Description:      "microservice for storing cars info",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
