package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
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

func getDeleteApplicationQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Deletes existing application",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"applicationID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)
			applicationIDArg, _ := resolveParams.Args["applicationID"].(string)

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(tenantIDArg); err != nil {
				return false, err
			}

			var applicationID system.UUID

			if applicationID, err = system.ParseUUID(applicationIDArg); err != nil {
				return false, err
			}

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.DeleteApplication(tenantID, applicationID)

			if err != nil {
				return false, err
			}

			return true, nil
		},
	}
}
