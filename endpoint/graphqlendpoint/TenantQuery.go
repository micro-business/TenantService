package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
)

func getTenantQuery() *graphql.Field {
	return &graphql.Field{
		Type:        tenantType,
		Description: "Returns an existing tenant",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)
			id, _ := resolveParams.Args["id"].(string)

			tenantID, err := system.ParseUUID(id)

			if err != nil {
				return nil, err
			}

			var returnedTenant domain.Tenant

			if returnedTenant, err = executionContext.tenantService.ReadTenant(tenantID); err != nil {
				return nil, err
			}

			return tenant{SecretKey: returnedTenant.SecretKey}, nil
		},
	}
}
