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

// Create creates a new tenant.
// tenant: Mandatory. The reference to the new tenant information
// Returns either the unique identifier of the new tenant or error if something goes wrong.
func (tenantService TenantService) Create(tenant domain.Tenant) (system.UUID, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")

	validateTenant(tenant)

	return tenantService.TenantDataService.CreateTenant(mapToDataTenant(tenant))
}

// Update updates an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// tenant: Mandatory. The reference to the updated tenant information.
// Returns error if something goes wrong.
func (tenantService TenantService) Update(tenantID system.UUID, tenant domain.Tenant) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	validateTenant(tenant)

	return tenantService.TenantDataService.UpdateTenant(tenantID, mapToDataTenant(tenant))
}

// Read retrieves an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// Returns either the tenant information or error if something goes wrong.
func (tenantService TenantService) Read(tenantID system.UUID) (domain.Tenant, error) {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	tenant, err := tenantService.TenantDataService.ReadTenant(tenantID)

	if err != nil {
		return domain.Tenant{}, err
	}

	return mapFromDataTenant(tenant), nil
}

// Delete deletes an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// Returns error if something goes wrong.
func (tenantService TenantService) Delete(tenantID system.UUID) error {
	diagnostics.IsNotNil(tenantService.TenantDataService, "tenantService.TenantDataServic", "TenantDataService must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	return tenantService.TenantDataService.DeleteTenant(tenantID)
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
