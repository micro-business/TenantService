package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/business/domain"
)

func getTenantQuery() *graphql.Field {
	return &graphql.Field{
		Type:        tenantType,
		Description: "Returns an existing tenant",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)

			tenantID, err := system.ParseUUID(tenantIDArg)

			if err != nil {
				return nil, err
			}

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			var returnedTenant domain.Tenant

			if returnedTenant, err = executionContext.tenantService.ReadTenant(tenantID); err != nil {
				return nil, err
			}

			return tenant{ID: tenantID.String(), SecretKey: returnedTenant.SecretKey}, nil
		},
	}
}
