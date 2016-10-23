package graphqlendpoint

import "github.com/graphql-go/graphql"

func getCreateTenantQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.ID,
		Description: "Creates new tenant",
		Args: graphql.FieldConfigArgument{
			"tenant": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputTenantType),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			inputTenantArgument, _ := resolveParams.Args["tenant"].(map[string]interface{})

			tenant := resolveTenantFromInputTenantArgument(inputTenantArgument)

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			tenantID, err := executionContext.tenantService.CreateTenant(tenant)

			if err != nil {
				return nil, err
			}

			return tenantID.String(), nil
		},
	}
}
