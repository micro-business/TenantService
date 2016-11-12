package graphqlendpoint_test

import (
	"github.com/microbusinesslimited/Micro-Business-Core/system"
	"github.com/microbusinesslimited/TenantService/business/domain"
)

func createApplicationInfo() domain.Application {
	randomValue, _ := system.RandomUUID()
	return domain.Application{Name: randomValue.String()}
}
