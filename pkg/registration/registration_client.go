package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdField string = "targetIdentifier"
	ProbeCategory string = "CloudNative"
	TargetType string = "Prometheus"
)

// Registration Client for the Prometheus probe
// Implements the TurboRegistrationClient interface
type PrometheusRegistrationClient struct {
}

func (myProbe *PrometheusRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	glog.Infoln("Building a supply chain ..........")

	// 2. Build supply chain.
	supplyChainFactory := &SupplyChainFactory{}
	templateDtos, err := supplyChainFactory.CreateSupplyChain()
	if err != nil {
		glog.Infoln("[ExampleProbe] Error creating Supply chain for the example probe")
		return nil
	}
	glog.Infoln("[ExampleProbe] Supply chain for the example probe is created.")
	return templateDtos
}

func (registrationClient *PrometheusRegistrationClient) GetIdentifyingFields() string {
	return TargetIdField
}

func (myProbe *PrometheusRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {

	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdField, "URL",
		"URL of the Prometheus target", ".*", true, false).Create()

	return []*proto.AccountDefEntry{
		targetIDAcctDefEntry,
	}
}
