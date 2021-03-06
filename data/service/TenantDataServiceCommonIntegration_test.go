// +build integration

package service_test

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/micro-business/Micro-Business-Core/system"
	"github.com/micro-business/TenantService/data/contract"
	"github.com/micro-business/TenantService/data/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const databasePreparationMaxTimeout = time.Minute

var keyspace string

var _ = BeforeSuite(func() {
	keyspace = createRandomKeyspace()
	createTenantKeyspaceAndAllRequiredTables(keyspace)
})

var _ = AfterSuite(func() {
	dropKeyspace(keyspace)
})

func getClusterConfig() *gocql.ClusterConfig {
	cassandraIPAddress := os.Getenv("CASSANDRA_ADDRESS")

	if len(cassandraIPAddress) == 0 {
		cassandraIPAddress = "127.0.0.1"
	}

	config := gocql.NewCluster(cassandraIPAddress)

	cassandraProtocolVersion := os.Getenv("CASSANDRA_PROTOCOL_VERSION")

	if len(cassandraProtocolVersion) != 0 {
		if protocolVersion, err := strconv.Atoi(cassandraProtocolVersion); err == nil {
			config.ProtoVersion = protocolVersion
		}
	}

	config.Consistency = gocql.Quorum

	return config
}

func createRandomKeyspace() string {
	keyspaceRandomValue, _ := system.RandomUUID()

	return strings.ToLower("a" + strings.Replace(keyspaceRandomValue.String(), "-", "", -1))
}

func createTenantKeyspaceAndAllRequiredTables(keyspace string) {
	config := getClusterConfig()
	config.Timeout = databasePreparationMaxTimeout
	session, err := config.CreateSession()

	Expect(err).To(BeNil())

	defer session.Close()

	Expect(session.Query(
		"CREATE KEYSPACE " +
			keyspace +
			" with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").
		Exec()).To(BeNil())

	Expect(session.Query(
		"CREATE TABLE " +
			keyspace +
			".tenant(tenant_id UUID, secret_key text," +
			" PRIMARY KEY(tenant_id));").
		Exec()).To(BeNil())

	Expect(session.Query(
		"CREATE TABLE " +
			keyspace +
			".application(tenant_id UUID, application_id UUID, name text," +
			" PRIMARY KEY(tenant_id, application_id));").
		Exec()).To(BeNil())
}

func createTenant(keyspace string) (system.UUID, contract.Tenant, error) {
	tenant := createTenantInfo()
	tenantID, err := createService().CreateTenant(tenant)

	if err != nil {
		return system.EmptyUUID, contract.Tenant{}, nil
	}

	return tenantID, tenant, nil
}

func createApplication(keyspace string) (system.UUID, contract.Tenant, system.UUID, contract.Application, error) {
	tenantID, tenant, err := createTenant(keyspace)

	if err != nil {
		return system.EmptyUUID, contract.Tenant{}, system.EmptyUUID, contract.Application{}, err
	}

	application := createApplicationInfo()
	applicationID, err := createService().CreateApplication(tenantID, application)

	if err != nil {
		return system.EmptyUUID, contract.Tenant{}, system.EmptyUUID, contract.Application{}, err
	}

	return tenantID, tenant, applicationID, application, nil
}

func createApplications(keyspace string) (system.UUID, contract.Tenant, map[system.UUID]contract.Application, error) {
	tenantID, tenant, err := createTenant(keyspace)

	if err != nil {
		return system.EmptyUUID, contract.Tenant{}, nil, err
	}

	applications := make(map[system.UUID]contract.Application)

	for idx := 0; idx < rand.Intn(5)+1; idx++ {
		application := createApplicationInfo()
		applicationID, err := createService().CreateApplication(tenantID, application)

		if err != nil {
			return system.EmptyUUID, contract.Tenant{}, nil, err
		}

		applications[applicationID] = application
	}

	return tenantID, tenant, applications, nil
}

func dropKeyspace(keyspace string) {
	config := getClusterConfig()
	config.Timeout = databasePreparationMaxTimeout
	session, err := config.CreateSession()

	Expect(err).To(BeNil())

	defer session.Close()

	err = session.Query("DROP KEYSPACE " + keyspace + " ;").Exec()

	Expect(err).To(BeNil())
}

func mapSystemUUIDToGocqlUUID(uuid system.UUID) gocql.UUID {
	mappedUUID, _ := gocql.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

func mapGocqlUUIDToSystemUUID(uuid gocql.UUID) system.UUID {
	mappedUUID, _ := system.UUIDFromBytes(uuid.Bytes())

	return mappedUUID
}

func createService() contract.TenantDataService {
	clusterConfig := getClusterConfig()
	clusterConfig.Keyspace = keyspace

	return &service.TenantDataService{UUIDGeneratorService: system.UUIDGeneratorServiceImpl{}, ClusterConfig: clusterConfig}
}
