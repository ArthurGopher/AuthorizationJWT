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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/address/geocode": {
            "post": {
                "description": "Возвращает адреса по переданным широте и долготе через DaData API.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Геокодирование"
                ],
                "summary": "Обратное геокодирование",
                "parameters": [
                    {
                        "description": "Широта и долгота",
                        "name": "coordinates",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.GeocodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.GeocodeResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка вызова API",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "description": "Ищет адреса по переданному параметру query через DaData API.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Адреса"
                ],
                "summary": "Поиск адресов",
                "parameters": [
                    {
                        "description": "Запрос с адресом",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SearchRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.SearchResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка вызова API",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Выполняет вход пользователя и возвращает JWT-токен, если учетные данные верны",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "Вход в систему пользователя",
                "parameters": [
                    {
                        "description": "Учетные данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь не существует или пароль не совпадает",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Регистрирует нового пользователя, сохраняя в памяти его имя пользователя и хэшированный пароль",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "Регистрация нового пользователя",
                "parameters": [
                    {
                        "description": "Учетные данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь успешно зарегистрирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Address": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                }
            }
        },
        "main.GeocodeRequest": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "string"
                },
                "lng": {
                    "type": "string"
                }
            }
        },
        "main.GeocodeResponse": {
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Address"
                    }
                }
            }
        },
        "main.SearchRequest": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string"
                }
            }
        },
        "main.SearchResponse": {
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Address"
                    }
                }
            }
        },
        "main.User": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
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
