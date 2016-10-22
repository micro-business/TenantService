package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

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