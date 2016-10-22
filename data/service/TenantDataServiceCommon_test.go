// +build integration

package service_test

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
	"github.com/microbusinesses/TenantService/data/service"
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

func createTenant(keyspace string) (system.UUID, error) {
	randomValue, _ := system.RandomUUID()
	validTenant := contract.Tenant{SecretKey: randomValue.String()}

	return createService().CreateTenant(validTenant)
}

func createApplication(keyspace string) (system.UUID, system.UUID, error) {
	tenantID, err := createTenant(keyspace)

	if err != nil {
		return system.EmptyUUID, system.EmptyUUID, err
	}

	randomValue, _ := system.RandomUUID()
	validApplication := contract.Application{Name: randomValue.String()}

	applicationID, err := createService().CreateApplication(tenantID, validApplication)

	if err != nil {
		return system.EmptyUUID, system.EmptyUUID, err
	}

	return tenantID, applicationID, nil
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
