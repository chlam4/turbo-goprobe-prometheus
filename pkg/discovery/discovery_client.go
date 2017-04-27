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
	"net/url"
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
		value, err := discClient.PrometheusApi.Query(context.Background(), query, time.Now())
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
			respTimeCommodity, _ := builder.NewCommodityDTOBuilder(action.CommodityType).
				Capacity(action.Capacity).Used(float64(metric.Value)).Create()
			//vcpuCommodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).Used(3.5).Create()
			//vmemCommodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).Used(6.5).Create()

			entityName := string(metric.Metric["instance"])
			entityUrl, err := url.Parse(entityName)
			if err != nil {
				glog.Errorf("Error extracting the instance URL from metric %v: %s", metric, err)
				continue
			}
			if (action.AppType != AppType_Web) {
				entityName = action.AppType + "_" + entityUrl.Hostname()
			}
			ipAddress := entityUrl.Hostname() // TODO perform DNS lookup here
			appDto, err := builder.NewEntityDTOBuilder(proto.EntityDTO_APPLICATION, entityName).
				DisplayName(entityName).
				SellsCommodity(respTimeCommodity).
				ApplicationData(&proto.EntityDTO_ApplicationData{
					Type:      &action.AppType,
					IpAddress: &ipAddress,
				}).
				WithProperty(&proto.EntityDTO_EntityProperty{
					Namespace: &propertyNamespace,
					Name:      &propertyName,
					Value:     &ipAddress,
				}).Create()

			if err != nil {
				glog.Errorf("Error building EntityDTO from metric %v: %s", metric, err)
				continue
			}
			entities = append(entities, appDto)
		}
	}

	discoveryResponse := &proto.DiscoveryResponse{
		EntityDTO: entities,
	}
	glog.Infof("Prometheus discovery response %s\n", discoveryResponse)

	return discoveryResponse, nil
}
