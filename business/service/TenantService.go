package service

import (
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
	"github.com/microbusinesses/TenantService/data/contract"
)

// TenantService provides access to add new tenant and update/retrieve/remove an existing tenant.
type TenantService struct {
	TenantDataService contract.TenantDataService
}

// CreateTenant creates a new tenant.
// tenant: Mandatory. The reference to the new tenant information
// Returns either the unique identifier of the new tenant or error if something goes wrong.
func (tenantService TenantService) CreateTenant(tenant domain.Tenant) (system.UUID, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")

	validateTenant(tenant)

	return tenantService.TenantDataService.CreateTenant(mapToDataTenant(tenant))
}

// UpdateTenant updates an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// tenant: Mandatory. The reference to the updated tenant information.
// Returns error if something goes wrong.
func (tenantService TenantService) UpdateTenant(tenantID system.UUID, tenant domain.Tenant) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	validateTenant(tenant)

	return tenantService.TenantDataService.UpdateTenant(tenantID, mapToDataTenant(tenant))
}

// ReadTenant retrieves an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// Returns either the tenant information or error if something goes wrong.
func (tenantService TenantService) ReadTenant(tenantID system.UUID) (domain.Tenant, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	tenant, err := tenantService.TenantDataService.ReadTenant(tenantID)

	if err != nil {
		return domain.Tenant{}, err
	}

	return mapFromDataTenant(tenant), nil
}

// DeleteTenant deletes an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// Returns error if something goes wrong.
func (tenantService TenantService) DeleteTenant(tenantID system.UUID) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	return tenantService.TenantDataService.DeleteTenant(tenantID)
}

// CreateApplication creates new application for the provided tenant.
// tenantID: Mandatory. The unique identifier of the tenant to create the application for.
// application: Mandatory. The reference to the new application to create for the provided tenant
// Returns either the unique identifier of the new application or error if something goes wrong.
func (tenantService TenantService) CreateApplication(tenantID system.UUID, application domain.Application) (system.UUID, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	validateApplication(application)

	return tenantService.TenantDataService.CreateApplication(tenantID, mapToDataApplication(application))
}

// Update updates an existing tenant application.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// applicationID: Mandatory: The unique identifier of the existing application.
// application: Mandatory. The reference to the updated application information.
// Returns error if something goes wrong.
func (tenantService TenantService) UpdateApplication(tenantID system.UUID, applicationID system.UUID, application domain.Application) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")

	validateApplication(application)

	return tenantService.TenantDataService.UpdateApplication(tenantID, applicationID, mapToDataApplication(application))
}

// Read retrieves an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// applicationID: Mandatory: The unique identifier of the existing application.
// Returns either the tenant application information or error if something goes wrong.
func (tenantService TenantService) ReadApplication(tenantID system.UUID, applicationID system.UUID) (domain.Application, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")

	application, err := tenantService.TenantDataService.ReadApplication(tenantID, applicationID)

	if err != nil {
		return domain.Application{}, err
	}

	return mapFromDataApplication(application), nil
}

// Delete deletes an existing tenant application information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// applicationID: Mandatory: The unique identifier of the existing application.
// Returns error if something goes wrong.
func (tenantService TenantService) DeleteApplication(tenantID system.UUID, applicationID system.UUID) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")
	diagnostics.IsNotNilOrEmpty(applicationID, "applicationID", "applicationID must be provided.")

	return tenantService.TenantDataService.DeleteApplication(tenantID, applicationID)
}

// validateTenant validates the tenant domain object and make sure the data is consistent and valid.
func validateTenant(tenant domain.Tenant) {
	diagnostics.IsNotNilOrEmptyOrWhitespace(tenant.SecretKey, "tenant.SecretKey", "SecretKey must be provided.")
}

// mapToDataTenant Maps the domain tenant object to the tenant object used in data layer.
// tenant: Mandatory. The tenant domain object
// Returns the converted tenant object used in data layer
func mapToDataTenant(tenant domain.Tenant) contract.Tenant {
	return contract.Tenant{SecretKey: tenant.SecretKey}
}

// mapFromDataTenant Maps the tenant object used in data layer to the tenant domain object.
// tenant: Mandatory. The tenant object used in data layer
// Returns the converted tenant domain object
func mapFromDataTenant(tenant contract.Tenant) domain.Tenant {
	return domain.Tenant{SecretKey: tenant.SecretKey}
}

// validateApplication validates the tenant application domain object and make sure the data is consistent and valid.
func validateApplication(application domain.Application) {
	diagnostics.IsNotNilOrEmptyOrWhitespace(application.Name, "application.Name", "Name must be provided.")
}

// mapToDataApplication Maps the domain tenant application object to the tenant application object used in data layer.
// application: Mandatory. The tenant application domain object
// Returns the converted tenant application object used in data layer
func mapToDataApplication(application domain.Application) contract.Application {
	return contract.Application{Name: application.Name}
}

// mapFromDataApplication Maps the tenant application object used in data layer to the tenant application domain object.
// application: Mandatory. The tenant application object used in data layer
// Returns the converted tenant application domain object
func mapFromDataApplication(application contract.Application) domain.Application {
	return domain.Application{Name: application.Name}
}
