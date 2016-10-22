package endpoint

import (
	"log"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/TenantService/business/contract"
	"github.com/microbusinesses/TenantService/config"
	"github.com/microbusinesses/TenantService/endpoint/transport"
	"golang.org/x/net/context"
)

// Endpoint implements method to start the service. The structure contains all the dependencies required by the Endpoint service.
type Endpoint struct {
	ConfigurationReader config.ConfigurationReader
	TenantService       contract.TenantService
}

// StartServer creates all the endpoints and starts the server.
func (endpoint Endpoint) StartServer() {
	diagnostics.IsNotNil(endpoint.TenantService, "endpoint.TenantService", "TenantService must be provided.")
	diagnostics.IsNotNil(endpoint.ConfigurationReader, "endpoint.ConfigurationReader", "ConfigurationReader must be provided.")

	http.Handle("/Api", httptransport.NewServer(
		context.Background(),
		createAPIEndpoint(endpoint.TenantService),
		transport.DecodeAPIRequest,
		transport.EncodeAPIResponse))

	if listeningPort, err := endpoint.ConfigurationReader.GetListeningPort(); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(listeningPort), nil))
	}
}
