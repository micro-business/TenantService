package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesslimited/TenantService/business/domain"
)

const (
	applicationID = "ID"
	name          = "Name"
)

type application struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

var applicationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Application",
		Fields: graphql.Fields{
			applicationID: &graphql.Field{Type: graphql.String},
			name:          &graphql.Field{Type: graphql.String},
		},
	},
)

var inputApplicationType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Application",
		Fields: graphql.InputObjectConfigFieldMap{
			name: &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	},
)

func resolveApplicationFromInputApplicationArgument(inputApplicationArgument map[string]interface{}) domain.Application {
	application := domain.Application{}

	nameKeyArg, nameArgProvided := inputApplicationArgument[name].(string)

	if nameArgProvided {
		application.Name = nameKeyArg
	}

	return application
}
