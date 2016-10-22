package graphqlendpoint

import (
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/contract"
	"github.com/microbusinesses/TenantService/business/domain"
	"golang.org/x/net/context"
)

const (
	secretKey = "SecretKey"
)

type tenant struct {
	SecretKey string `json:"SecretKey"`
}

var tenantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tenant",
		Fields: graphql.Fields{
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

var rootQueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"tenant": getTenant(),
		},
	},
)

var rootMutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTenant": getCreateTenant(),
			"updateTenant": getUpdateTenant(),
			"deleteTenant": getDeleteTenant(),
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

func getTenant() *graphql.Field {
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

func getCreateTenant() *graphql.Field {
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

func getUpdateTenant() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.ID,
		Description: "Updates existing tenant",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"tenant": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputTenantType),
			},
		},
		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			inputTenantArgument, _ := resolveParams.Args["tenant"].(map[string]interface{})
			id, _ := resolveParams.Args["id"].(string)

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(id); err != nil {
				return nil, err
			}

			tenant := resolveTenantFromInputTenantArgument(inputTenantArgument)

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.UpdateTenant(tenantID, tenant)

			if err != nil {
				return nil, err
			}

			return tenantID.String(), nil
		},
	}
}

func getDeleteTenant() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.ID,
		Description: "Deletes existing tenant",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
		},
		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			id, _ := resolveParams.Args["id"].(string)

			var tenantID system.UUID
			var err error

			if tenantID, err = system.ParseUUID(id); err != nil {
				return nil, err
			}

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			err = executionContext.tenantService.DeleteTenant(tenantID)

			if err != nil {
				return nil, err
			}

			return tenantID.String(), nil

		},
	}
}

func resolveTenantFromInputTenantArgument(inputTenantArgument map[string]interface{}) domain.Tenant {
	tenant := domain.Tenant{}

	secretKeyArg, secretKeyArgProvided := inputTenantArgument[secretKey].(string)

	if secretKeyArgProvided {
		tenant.SecretKey = secretKeyArg
	}

	return tenant
}
