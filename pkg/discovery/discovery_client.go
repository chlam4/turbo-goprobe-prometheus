package discovery

import (
	"context"
	"fmt"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/conf"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/registration"
	"github.com/golang/glog"
	prometheusHttpClient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
	"time"
)

// Discovery Client for the Prometheus Probe
// Implements the TurboDiscoveryClient interface
type PrometheusDiscoveryClient struct {
	TargetConf    *conf.PrometheusTargetConf
	PrometheusApi prometheus.API
}

func NewDiscoveryClient(confFile string) (*PrometheusDiscoveryClient, error) {
	//
	// Parse conf file to create clientConf
	//
	targetConf, err := conf.NewPrometheusTargetConf(confFile)
	if err != nil {
		return nil, err
	}
	glog.Infof("Target Conf %v\n", targetConf)
	//
	// Create a Prometheus client
	//
	promConfig := prometheusHttpClient.Config{Address: targetConf.Address}
	promHttpClient, err := prometheusHttpClient.NewClient(promConfig)
	if err != nil {
		return nil, err
	}

	return &PrometheusDiscoveryClient{
		TargetConf:    targetConf,
		PrometheusApi: prometheus.NewAPI(promHttpClient),
	}, nil
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (discClient *PrometheusDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	// Convert all parameters in clientConf to AccountValue list
	targetConf := discClient.TargetConf

	targetId := registration.TargetIdField
	targetIdVal := &proto.AccountValue{
		Key:         &targetId,
		StringValue: &targetConf.Address,
	}

	accountValues := []*proto.AccountValue{
		targetIdVal,
	}

	targetInfo := probe.NewTurboTargetInfoBuilder(registration.ProbeCategory, registration.TargetType,
		registration.TargetIdField, accountValues).Create()

	return targetInfo
}

// Validate the Target
func (discClient *PrometheusDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.Infof("BEGIN Validation for PrometheusDiscoveryClient %s\n", accountValues)

	validationResponse := &proto.ValidationResponse{}

	glog.Infof("Validation response %s\n", validationResponse)
	return validationResponse, nil
}

// Discover the Target Topology
func (discClient *PrometheusDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.Infof("========= Discovering Prometheus ============= %s\n", accountValues)
	var entities []*proto.EntityDTO

	for query, action := range Blueprint {
		value, err := discClient.PrometheusApi.Query(context.Background(), string(query), time.Now())
		if err != nil {
			glog.Errorf("Error while discovering Prometheus target %s: %s\n", discClient.TargetConf.Address, err)
			// If there is error during discovery, return an ErrorDTO.
			severity := proto.ErrorDTO_CRITICAL
			description := fmt.Sprintf("%v", err)
			errorDTO := &proto.ErrorDTO{
				Severity:    &severity,
				Description: &description,
			}
			discoveryResponse := &proto.DiscoveryResponse{
				ErrorDTO: []*proto.ErrorDTO{errorDTO},
			}
			return discoveryResponse, nil
		}

		propertyNamespace := "DEFAULT"
		propertyName := supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS
		for _, metric := range value.(model.Vector) {
			//
			// Extract entity id and node IP from the metric
			//
			instanceName := string(metric.Metric["instance"])
			appData, err := action.GetAppData(instanceName)
			if err != nil {
				glog.Errorf("Error extracting id and node IP info from metric %v: %s", metric, err)
				continue
			}
			//
			// Construct the commodity
			//
			commodity, err := builder.NewCommodityDTOBuilder(action.CommodityType).
				Capacity(action.Capacity).Used(float64(metric.Value)).Create()
			if err != nil {
				glog.Errorf("Error building a commodity: %s", err)
				continue
			}

			dto, err := builder.NewEntityDTOBuilder(proto.EntityDTO_APPLICATION, appData.id).
				DisplayName(appData.id).
				SellsCommodity(commodity).
				ApplicationData(&proto.EntityDTO_ApplicationData{
					Type:      &appData.appType,
					IpAddress: &appData.nodeIp,
				}).
				WithProperty(&proto.EntityDTO_EntityProperty{
					Namespace: &propertyNamespace,
					Name:      &propertyName,
					Value:     &appData.nodeIp,
				}).Create()

			if err != nil {
				glog.Errorf("Error building EntityDTO from metric %v: %s", metric, err)
				continue
			}
			entities = append(entities, dto)
		}
	}

	discoveryResponse := &proto.DiscoveryResponse{
		EntityDTO: entities,
	}
	glog.Infof("Prometheus discovery response %s\n", discoveryResponse)

	return discoveryResponse, nil
}
