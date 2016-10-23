package graphqlendpoint_test

import (
	"github.com/microbusinesses/Micro-Businesses-Core/system"
	"github.com/microbusinesses/TenantService/business/domain"
)

func createApplicationInfo() domain.Application {
	randomValue, _ := system.RandomUUID()
	return domain.Application{Name: randomValue.String()}
}
