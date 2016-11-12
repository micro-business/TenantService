package graphqlendpoint

import (
	"github.com/graphql-go/graphql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/business/domain"
)

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
