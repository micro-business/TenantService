package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

func getDeleteTenantQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Deletes existing tenant",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(tenantIDArg); err != nil {
				return false, err
			}

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.DeleteTenant(tenantID)

			if err != nil {
				return false, err
			}

			return true, nil
		},
	}
}
