package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesslimited/Micro-Business-Core/system"
)

func getUpdateTenantQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Updates existing tenant",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"tenant": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputTenantType),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)
			inputTenantArgument, _ := resolveParams.Args["tenant"].(map[string]interface{})

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(tenantIDArg); err != nil {
				return false, err
			}

			tenant := resolveTenantFromInputTenantArgument(inputTenantArgument)

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.UpdateTenant(tenantID, tenant)

			if err != nil {
				return false, err
			}

			return true, nil
		},
	}
}
