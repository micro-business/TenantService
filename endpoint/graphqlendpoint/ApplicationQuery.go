package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
)

func getApplicationQuery() *graphql.Field {
	return &graphql.Field{
		Type:        applicationType,
		Description: "Returns an existing application",
		Args: graphql.FieldConfigArgument{
			"tenantID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"applicationID": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},

		Resolve: func(resolveParams graphql.ResolveParams) (interface{}, error) {
			tenantIDArg, _ := resolveParams.Args["tenantID"].(string)
			applicationIDArg, _ := resolveParams.Args["applicationID"].(string)

			tenantID, err := system.ParseUUID(tenantIDArg)

			if err != nil {
				return nil, err
			}

			applicationID, err := system.ParseUUID(applicationIDArg)

			if err != nil {
				return nil, err
			}

			var returnedApplication domain.Application

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			if returnedApplication, err = executionContext.tenantService.ReadApplication(tenantID, applicationID); err != nil {
				return nil, err
			}

			return application{ID: applicationID.String(), Name: returnedApplication.Name}, nil
		},
	}
}

func getApplicationsQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(applicationType),
		Description: "Returns all registered applications for the provided tenant",
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

			var returnedApplications map[system.UUID]domain.Application

			executionContext := resolveParams.Context.Value("ExecutionContext").(executionContext)

			if returnedApplications, err = executionContext.tenantService.ReadAllApplications(tenantID); err != nil {
				return nil, err
			}

			applications := make([]application, 0, len(returnedApplications))

			for applicationID, app := range returnedApplications {
				applications = append(applications, application{ID: applicationID.String(), Name: app.Name})
			}

			return applications, nil
		},
	}
}
