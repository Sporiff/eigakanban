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
            "name": "Ciarán Ainsworth",
            "url": "https://codeberg.org/sporiff/eigakanban/issues",
            "email": "cda@sporiff.dev"
        },
        "license": {
            "name": "AGPL3 or Later",
            "url": "https://codeberg.org/sporiff/eigakanban/src/branch/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Log in to user account using email or username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Log in",
                "parameters": [
                    {
                        "description": "Login details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginUser.LoginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful login",
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginUser.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Missing mandatory fields",
                        "schema": {
                            "$ref": "#/definitions/handlers.MissingFieldResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginUser.NoUserFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Log out of the app",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Log out",
                "responses": {
                    "200": {
                        "description": "Logout successful",
                        "schema": {
                            "$ref": "#/definitions/handlers.LogoutUser.LogoutSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Missing refresh token",
                        "schema": {
                            "$ref": "#/definitions/handlers.LogoutUser.RefreshTokenMissingResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user account",
                "parameters": [
                    {
                        "description": "User details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterUser.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered successfully",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Missing mandatory fields",
                        "schema": {
                            "$ref": "#/definitions/handlers.MissingFieldResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/boards": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all boards in a paginated list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "boards"
                ],
                "summary": "Get all boards",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.GetAllBoards.PaginatedBoardsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Add a new board for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "boards"
                ],
                "summary": "Add a new board",
                "parameters": [
                    {
                        "description": "Board details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AddBoard.AddBoardRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Board added successfully",
                        "schema": {
                            "$ref": "#/definitions/handlers.BoardsResponse"
                        }
                    },
                    "400": {
                        "description": "Missing mandatory fields",
                        "schema": {
                            "$ref": "#/definitions/handlers.MissingFieldResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/boards/{uuid}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a board by UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "boards"
                ],
                "summary": "Get board by UUID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Board UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.BoardsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete a board by UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "boards"
                ],
                "summary": "Delete board",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Board UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Board deleted successfully",
                        "schema": {
                            "$ref": "#/definitions/handlers.DeleteBoard.BoardDeletedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user/{uuid}/boards": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all boards for a user in a paginated list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "boards"
                ],
                "summary": "Get all boards for a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.GetAllBoards.PaginatedBoardsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all users in a paginated list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.GetAllUsers.PaginatedUsersResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{uuid}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a user by UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user by UUID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete a user by UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted successfully",
                        "schema": {
                            "$ref": "#/definitions/handlers.DeleteUser.UserDeletedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update user details by UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User details to update",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.UpdateUser.UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.AddBoard.AddBoardRequest": {
            "type": "object",
            "required": [
                "name",
                "user_uuid"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "A short description"
                },
                "name": {
                    "type": "string",
                    "example": "My Queue"
                },
                "user_uuid": {
                    "type": "string",
                    "example": "00ca71c5-7c8a-4470-ab47-f962d33c1303"
                }
            }
        },
        "handlers.BoardsResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "A short description"
                },
                "name": {
                    "type": "string",
                    "example": "My queue"
                },
                "uuid": {
                    "type": "string",
                    "example": "00000000-0000-0000-0000-000000000000"
                }
            }
        },
        "handlers.DeleteBoard.BoardDeletedResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string",
                    "example": "Board deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"
                }
            }
        },
        "handlers.DeleteUser.UserDeletedResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string",
                    "example": "user deleted: 77b62cff-0020-43d9-a90c-5d35bff89f7a"
                }
            }
        },
        "handlers.GetAllBoards.PaginatedBoardsResponse": {
            "type": "object",
            "properties": {
                "boards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.BoardsResponse"
                    }
                },
                "pagination": {
                    "$ref": "#/definitions/types.Pagination"
                }
            }
        },
        "handlers.GetAllUsers.PaginatedUsersResponse": {
            "type": "object",
            "properties": {
                "pagination": {
                    "$ref": "#/definitions/types.Pagination"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.UserResponse"
                    }
                }
            }
        },
        "handlers.LoginUser.LoginUserRequest": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "username": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "handlers.LoginUser.NoUserFoundResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "User not found"
                }
            }
        },
        "handlers.LoginUser.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "jwt-token-string"
                }
            }
        },
        "handlers.LogoutUser.AlreadyLoggedOutResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Already logged out"
                }
            }
        },
        "handlers.LogoutUser.LogoutSuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Logged out successfully"
                }
            }
        },
        "handlers.LogoutUser.RefreshTokenMissingResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Refresh token is required"
                }
            }
        },
        "handlers.MissingFieldResponse": {
            "description": "an example of a missing field response",
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "properties": {
                        "username": {
                            "type": "string",
                            "example": "This field is required"
                        }
                    }
                }
            }
        },
        "handlers.RegisterUser.RegisterUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "username": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "handlers.UpdateBoard.UpdateBoardRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "A short description"
                },
                "name": {
                    "type": "string",
                    "example": "My Board"
                }
            }
        },
        "handlers.UpdateUser.UpdateUserRequest": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string",
                    "example": "This is a bio"
                },
                "full_name": {
                    "type": "string",
                    "example": "Tim Test"
                },
                "username": {
                    "type": "string",
                    "example": "new_username"
                }
            }
        },
        "handlers.UserResponse": {
            "description": "JSON representation of a user in the system",
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string",
                    "example": "This is a bio"
                },
                "full_name": {
                    "type": "string",
                    "example": "Tim Test"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                },
                "uuid": {
                    "type": "string",
                    "example": "77b62cff-0020-43d9-a90c-5d35bff89f7a"
                }
            }
        },
        "types.ErrorResponse": {
            "description": "an unknown error",
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "internal server error"
                }
            }
        },
        "types.Pagination": {
            "description": "pagination information",
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer",
                    "example": 1
                },
                "page_size": {
                    "type": "integer",
                    "example": 50
                },
                "total": {
                    "type": "integer",
                    "example": 2
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "eigakanban API",
	Description:      "The REST API for the eigakanban server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
