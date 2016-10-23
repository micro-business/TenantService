package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

func getUpdateApplicationQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Updates existing application",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"applicationID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"application": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputApplicationType),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)
			applicationIDArg, _ := resolveParams.Args["applicationID"].(string)
			inputApplicationArgument, _ := resolveParams.Args["application"].(map[string]interface{})

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(tenantIDArg); err != nil {
				return false, err
			}

			var applicationID system.UUID

			if applicationID, err = system.ParseUUID(applicationIDArg); err != nil {
				return false, err
			}

			application := resolveApplicationFromInputApplicationArgument(inputApplicationArgument)

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.UpdateApplication(tenantID, applicationID, application)

			if err != nil {
				return false, err
			}

			return true, nil
		},
	}
}
