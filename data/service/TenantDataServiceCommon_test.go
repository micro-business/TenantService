package service_test

import (
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/data/contract"
)

func createTenantInfo() contract.Tenant {
	randomValue, _ := system.RandomUUID()
	return contract.Tenant{SecretKey: randomValue.String()}
}

func createApplicationInfo() contract.Application {
	randomValue, _ := system.RandomUUID()
	return contract.Application{Name: randomValue.String()}
}
