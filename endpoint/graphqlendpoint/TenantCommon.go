package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/TenantService/business/domain"
)

const (
	tenantID  = "ID"
	secretKey = "SecretKey"
)

type tenant struct {
	ID        string `json:"ID"`
	SecretKey string `json:"SecretKey"`
}

var tenantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tenant",
		Fields: graphql.Fields{
			tenantID:  &graphql.Field{Type: graphql.String},
			secretKey: &graphql.Field{Type: graphql.String},
		},
	},
)

var inputTenantType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Tenant",
		Fields: graphql.InputObjectConfigFieldMap{
			secretKey: &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	},
)

func resolveTenantFromInputTenantArgument(inputTenantArgument map[string]interface{}) domain.Tenant {
	tenant := domain.Tenant{}

	secretKeyArg, secretKeyArgProvided := inputTenantArgument[secretKey].(string)

	if secretKeyArgProvided {
		tenant.SecretKey = secretKeyArg
	}

	return tenant
}
