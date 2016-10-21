// Package contract defines the tenant service contract.
package contract

import (
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
)

// TenantService contract, it can add new tenant and update/retrieve/remove an existing tenant.
type TenantService interface {
	// CreateTenant creates a new tenant.
	// tenant: Mandatory. The reference to the new tenant information
	// Returns either the unique identifier of the new tenant or error if something goes wrong.
	CreateTenant(tenant domain.Tenant) (system.UUID, error)

	// UpdateTenant updates an existing tenant.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// tenant: Mandatory. The reference to the updated tenant information.
	// Returns error if something goes wrong.
	UpdateTenant(tenantID system.UUID, tenant domain.Tenant) error

	// ReadTenant retrieves an existing tenant.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// Returns either the tenant information or error if something goes wrong.
	ReadTenant(tenantID system.UUID) (domain.Tenant, error)

	// DeleteTenant deletes an existing tenant information.
	// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
	// Returns error if something goes wrong.
	DeleteTenant(tenantID system.UUID) error

	// CreateApplication creates new application for the provided tenant.
	// tenantID: Mandatory. The unique identifier of the tenant to create the application for.
	// application: Mandatory. The reference to the new application to create for the provided tenant
	// Returns either the unique identifier of the new application or error if something goes wrong.
	CreateApplication(tenantID system.UUID, application domain.Application) (system.UUID, error)

	// UpdateApplication updates an existing tenant application.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// applicationID: Mandatory: The unique identifier of the existing application.
	// application: Mandatory. The reference to the updated application information.
	// Returns error if something goes wrong.
	UpdateApplication(tenantID system.UUID, applicationID system.UUID, application domain.Application) error

	// ReadApplication retrieves an existing tenant information.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// applicationID: Mandatory: The unique identifier of the existing application.
	// Returns either the tenant application information or error if something goes wrong.
	ReadApplication(tenantID system.UUID, applicationID system.UUID) (domain.Application, error)

	// DeleteApplication deletes an existing tenant application information.
	// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
	// applicationID: Mandatory: The unique identifier of the existing application.
	// Returns error if something goes wrong.
	DeleteApplication(tenantID system.UUID, applicationID system.UUID) error
}
