package docs

import (
	"github.com/swaggo/swag"
)

var SwaggerInfo = swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), &SwaggerInfo)
}
