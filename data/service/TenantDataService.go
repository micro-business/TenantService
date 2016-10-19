package service

import (
	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
)

// TenantDataService provides access to add new tenant and update/retrieve/remove an existing tenant.
type TenantDataService struct {
	UUIDGeneratorService system.UUIDGeneratorService
	ClusterConfig        *gocql.ClusterConfig
}

// Create creates a new tenant.
// tenant: Mandatory. The reference to the new tenant information
// Returns either the unique identifier of the new tenant or error if something goes wrong.
func (tenantDataService TenantDataService) Create(tenant contract.Tenant) (system.UUID, error) {
	diagnostics.IsNotNil(tenantDataService.UUIDGeneratorService, "tenantDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmptyOrWhitespace(tenant.SecretKey, "tenant.SecretKey", "SecretKey must be provided.")

	tenantID, err := tenantDataService.UUIDGeneratorService.GenerateRandomUUID()

	if err != nil {
		return system.EmptyUUID, err
	}

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return system.EmptyUUID, err
	}

	defer session.Close()

	err = addNewTenant(tenantID, tenant, session)

	return tenantID, err
}

// Update updates an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// tenant: Mandatory. The reference to the updated tenant information.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) Update(tenantID system.UUID, tenant contract.Tenant) error {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmptyOrWhitespace(tenant.SecretKey, "tenant.SecretKey", "SecretKey must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return addNewTenant(tenantID, tenant, session)
}

// Read retrieves an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// Returns either the tenant information or error if something goes wrong.
func (tenantDataService TenantDataService) Read(tenantID system.UUID) (contract.Tenant, error) {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	return contract.Tenant{}, nil
}

// Delete deletes an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) Delete(tenantID system.UUID) error {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")
	diagnostics.IsNotNilOrEmpty(tenantID, "tenantID", "tenantID must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	return removeExistingTenant(tenantID, session)
}

// mapSystemUUIDToGocqlUUID maps the system type UUID to gocql UUID type
func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

// addNewTenant adds new tenant to tenant table
func addNewTenant(tenantID system.UUID, tenant contract.Tenant, session *gocql.Session) error {

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)

	return session.Query(
		"INSERT INTO tenant"+
			" (tenant_id, secret_key)"+
			" VALUES(?, ?)",
		mappedTenantID,
		tenant.SecretKey).
		Exec()

}

// removeExistingTenant adds new tenant to tenant table
func removeExistingTenant(tenantID system.UUID, session *gocql.Session) error {

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)

	return session.Query(
		"DELETE FROM tenant"+
			" WHERE"+
			" tenant_id = ?",
		mappedTenantID).
		Exec()
}
