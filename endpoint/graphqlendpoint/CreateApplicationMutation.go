package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/micro-business/Micro-Business-Core/system"
)

func getCreateApplicationQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.ID,
		Description: "Creates new application",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"application": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputApplicationType),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)
			inputApplicationArgument, _ := resolveParams.Args["application"].(map[string]interface{})

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(tenantIDArg); err != nil {
				return nil, err
			}

			application := resolveApplicationFromInputApplicationArgument(inputApplicationArgument)

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			applicationID, err := executionContext.tenantService.CreateApplication(tenantID, application)

			if err != nil {
				return nil, err
			}

			return applicationID.String(), nil
		},
	}
}
