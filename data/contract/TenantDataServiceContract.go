// Package contract defines the tenant data service contract.
package contract

import "github.com/microbusinesses/Micro-Businesses-Core/system"

type Tenant struct {
	SecretKey string
}

// TenantDataService service can add new tenant and update/retrieve/remove existing tenant.
type TenantDataService interface {
	// Create creates a new tenant.
	// tenant: Mandatory. The reference to the new tenant information
	// Returns either the unique identifier of the new tenant or error if something goes wrong.
	Create(tenant Tenant) (system.UUID, error)

	// Update updates an existing tenant.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// tenant: Mandatory. The reference to the updated tenant information.
	// Returns error if something goes wrong.
	Update(tenantID system.UUID, tenant Tenant) error

	// Read retrieves an existing tenant.
	// tenantID: Mandatory: The unique identifier of the existing tenant.
	// Returns either the tenant information or error if something goes wrong.
	Read(tenantID system.UUID) (Tenant, error)

	// Deletes deletes an existing tenant information.
	// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
	// Returns error if something goes wrong.
	Delete(tenantID system.UUID) error
}
