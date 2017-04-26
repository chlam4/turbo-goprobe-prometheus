package discovery

import (
	"fmt"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/conf"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/registration"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

// Discovery Client for the Prometheus Probe
// Implements the TurboDiscoveryClient interface
type PrometheusDiscoveryClient struct {
	ClientConf *conf.PrometheusTargetConf
}

func NewDiscoveryClient(confFile string) (*PrometheusDiscoveryClient, error) {
	// Parse conf file to create clientConf
	clientConf, _ := conf.NewPrometheusTargetConf(confFile)
	glog.Infof("Target Conf %v\n", clientConf)
	client := &PrometheusDiscoveryClient{
		ClientConf: clientConf,
	}

	return client, nil
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (handler *PrometheusDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	// Convert all parameters in clientConf to AccountValue list
	clientConf := handler.ClientConf

	targetId := registration.TargetIdField
	targetIdVal := &proto.AccountValue{
		Key:         &targetId,
		StringValue: &clientConf.Address,
	}

	accountValues := []*proto.AccountValue{
		targetIdVal,
	}

	targetInfo := probe.NewTurboTargetInfoBuilder(registration.ProbeCategory, registration.TargetType,
		registration.TargetIdField, accountValues).Create()

	return targetInfo
}

// Validate the Target
func (handler *PrometheusDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.Infof("BEGIN Validation for PrometheusDiscoveryClient %s\n", accountValues)
	// TODO: connect to the client and get validation response
	validationResponse := &proto.ValidationResponse{}

	glog.Infof("Validation response %s\n", validationResponse)
	return validationResponse, nil
}

// Discover the Target Topology
func (handler *PrometheusDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.Infof("========= Discovering Prometheus ============= %s\n", accountValues)

	respTimeCommodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).Create()
	propertyNamespace := "DEFAULT"
	propertyName := supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS
	ipAddress := "10.10.174.90"
	appDto, err := builder.NewEntityDTOBuilder(proto.EntityDTO_APPLICATION, "https://10.10.172.235:9400/com.vmturbo.UI/UIMain.html").
		DisplayName("OpsMgr-10.10.172.235").
		SellsCommodity(respTimeCommodity).
		WithProperty(&proto.EntityDTO_EntityProperty{
			Namespace: &propertyNamespace,
			Name:      &propertyName,
			Value:     &ipAddress,
		}).Create()

	var discoveryResponse *proto.DiscoveryResponse
	if err != nil {
		// If there is error during discovery, return an ErrorDTO.
		severity := proto.ErrorDTO_CRITICAL
		description := fmt.Sprintf("%v", err)
		errorDTO := &proto.ErrorDTO{
			Severity:    &severity,
			Description: &description,
		}
		discoveryResponse = &proto.DiscoveryResponse{
			ErrorDTO: []*proto.ErrorDTO{errorDTO},
		}
	} else {
		// No error. Return the result entityDTOs.
		discoveryResponse = &proto.DiscoveryResponse{
			EntityDTO: []*proto.EntityDTO{appDto},
		}
	}
	glog.Infof("Prometheus discovery response %s\n", discoveryResponse)
	return discoveryResponse, nil
}
