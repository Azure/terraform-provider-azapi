package clients

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	moduleName    = "resource"
	moduleVersion = "v0.1.0"
)

type NewResourceClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

func NewNewResourceClient(subscriptionID string, credential azcore.TokenCredential, opt *arm.ClientOptions) *NewResourceClient {
	client := &NewResourceClient{
		subscriptionID: subscriptionID,
		host:           string(opt.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, opt),
	}
	return client
}
