package service

import (
	"fmt"

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

// CreateTenant  creates a new tenant.
// tenant: Mandatory. The reference to the new tenant information
// Returns either the unique identifier of the new tenant or error if something goes wrong.
func (tenantDataService TenantDataService) CreateTenant(tenant contract.Tenant) (system.UUID, error) {
	diagnostics.IsNotNil(tenantDataService.UUIDGeneratorService, "tenantDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")

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

// UpdateTenant updates an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// tenant: Mandatory. The reference to the updated tenant information.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) UpdateTenant(tenantID system.UUID, tenant contract.Tenant) error {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	return addNewTenant(tenantID, tenant, session)
}

// ReadTenant retrieves an existing tenant.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// Returns either the tenant information or error if something goes wrong.
func (tenantDataService TenantDataService) ReadTenant(tenantID system.UUID) (contract.Tenant, error) {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return contract.Tenant{}, err
	}

	defer session.Close()

	return readTenant(tenantID, session)

}

// DeleteTenant deletes an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) DeleteTenant(tenantID system.UUID) error {
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)

	return session.Query(
		"DELETE FROM tenant"+
			" WHERE"+
			" tenant_id = ?",
		mappedTenantID).
		Exec()
}

// CreateApplication creates new application for the provided tenant.
// tenantID: Mandatory. The unique identifier of the tenant to create the application for.
// application: Mandatory. The reference to the new application to create for the provided tenant
// Returns either the unique identifier of the new application or error if something goes wrong.
func (tenantDataService TenantDataService) CreateApplication(tenantID system.UUID, application contract.Application) (system.UUID, error) {
	diagnostics.IsNotNil(tenantDataService.UUIDGeneratorService, "tenantDataService.UUIDGeneratorService", "UUIDGeneratorService must be provided.")
	diagnostics.IsNotNil(tenantDataService.ClusterConfig, "tenantDataServic.ClusterConfig", "ClusterConfig must be provided.")

	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return system.EmptyUUID, err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return system.EmptyUUID, fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	applicationID, err := tenantDataService.UUIDGeneratorService.GenerateRandomUUID()

	if err != nil {
		return system.EmptyUUID, err
	}

	err = addNewApplication(tenantID, applicationID, application, session)

	return applicationID, err
}

// UpdateApplication updates an existing tenant application.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// applicationID: Mandatory: The unique identifier of the existing application.
// application: Mandatory. The reference to the updated application information.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) UpdateApplication(tenantID system.UUID, applicationID system.UUID, application contract.Application) error {
	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	if !doesApplicationExist(tenantID, applicationID, session) {
		return fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), applicationID.String())
	}

	return addNewApplication(tenantID, applicationID, application, session)
}

// ReadApplication retrieves an existing tenant information.
// tenantID: Mandatory: The unique identifier of the existing tenant.
// applicationID: Mandatory: The unique identifier of the existing application.
// Returns either the tenant application information or error if something goes wrong.
func (tenantDataService TenantDataService) ReadApplication(tenantID system.UUID, applicationID system.UUID) (contract.Application, error) {
	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return contract.Application{}, err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return contract.Application{}, fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	return readApplication(tenantID, applicationID, session)
}

// DeleteApplication deletes an existing tenant application information.
// tenantID: Mandatory: The unique identifier of the existing tenant to remove.
// applicationID: Mandatory: The unique identifier of the existing application.
// Returns error if something goes wrong.
func (tenantDataService TenantDataService) DeleteApplication(tenantID system.UUID, applicationID system.UUID) error {
	session, err := tenantDataService.ClusterConfig.CreateSession()

	if err != nil {
		return err
	}

	defer session.Close()

	if !doesTenantExist(tenantID, session) {
		return fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	if !doesApplicationExist(tenantID, applicationID, session) {
		return fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), applicationID.String())
	}

	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)
	mappedApplicationID := mapSystemUUIDToGocqlUUID(applicationID)

	return session.Query(
		"DELETE FROM application"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?",
		mappedTenantID,
		mappedApplicationID).
		Exec()
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

// readTenant takes the provided tenantID and tries to read the tenant information from database
func readTenant(tenantID system.UUID, session *gocql.Session) (contract.Tenant, error) {
	iter := session.Query(
		"SELECT secret_key"+
			" FROM tenant"+
			" WHERE"+
			" tenant_id = ?",
		tenantID.String()).Iter()

	tenant := contract.Tenant{}

	if !iter.Scan(&tenant.SecretKey) {
		return contract.Tenant{}, fmt.Errorf("Tenant not found. Tenant ID: %s", tenantID.String())
	}

	return tenant, nil

}

// doesTenantExist checks whether the provided tenant exists in database
func doesTenantExist(tenantID system.UUID, session *gocql.Session) bool {
	iter := session.Query(
		"SELECT secret_key"+
			" FROM tenant"+
			" WHERE"+
			" tenant_id = ?",
		tenantID.String()).Iter()

	var secretKey string

	return iter.Scan(&secretKey)
}

// addNewApplication adds new qpplication to tenant application table
func addNewApplication(tenantID, applicationID system.UUID, application contract.Application, session *gocql.Session) error {
	mappedTenantID := mapSystemUUIDToGocqlUUID(tenantID)
	mappedApplicationID := mapSystemUUIDToGocqlUUID(applicationID)

	return session.Query(
		"INSERT INTO application"+
			" (tenant_id, application_, name)"+
			" VALUES(?, ?, ?)",
		mappedTenantID,
		mappedApplicationID,
		application.Name).
		Exec()
}

// readApplication takes the provided tenantID and applicationID and tries to read the tenant application information from database
func readApplication(tenantID, applicationID system.UUID, session *gocql.Session) (contract.Application, error) {
	iter := session.Query(
		"SELECT name"+
			" FROM application"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?",
		tenantID.String(),
		applicationID.String()).Iter()

	application := contract.Application{}

	if !iter.Scan(&application.Name) {
		return contract.Application{}, fmt.Errorf("Tenant Application not found. Tenant ID: %s, Application ID: %s", tenantID.String(), applicationID.String())
	}

	return application, nil

}

// doesApplicationExist checks whether the provided tenant application exists in database
func doesApplicationExist(tenantID system.UUID, applicationID system.UUID, session *gocql.Session) bool {
	iter := session.Query(
		"SELECT name"+
			" FROM application"+
			" WHERE"+
			" tenant_id = ?"+
			" AND application_id = ?",
		tenantID.String(),
		applicationID.String()).Iter()

	var name string

	return iter.Scan(&name)
}
