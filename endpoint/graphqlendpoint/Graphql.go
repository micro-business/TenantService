package graphqlendpoint

import (
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/TenantService/business/contract"
	"golang.org/x/net/context"
)

var rootQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"tenant": getTenantQuery(),
		},
	},
)

var rootMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTenant": getCreateTenantQuery(),
			"updateTenant": getUpdateTenantQuery(),
			"deleteTenant": getDeleteTenantQuery(),
		},
	},
)

var tenantSchema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: rootQueryType, Mutation: rootMutationType})

type executionContext struct {
	tenantService contract.TenantService
}

// ExecuteQuery executes the provided query and returns the result.
func ExecuteQuery(query string, tenantService contract.TenantService) (interface{}, error) {
	result := graphql.Do(
		graphql.Params{
			Schema:        tenantSchema,
			RequestString: query,
			Context:       context.WithValue(context.Background(), "ExecutionContext", executionContext{tenantService}),
		})

	if result.HasErrors() {
		errorMessages := []string{}

		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Error())
		}

		return nil, errors.New(strings.Join(errorMessages, "\n"))
	}

	return result, nil
}
