package discovery

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"net/url"
	"strings"
)

const (
	AppType_Web   = "web"
	AppType_MySQL = "mysql"
)

type Query string

type AppData struct {
	appType string
	id      string
	nodeIp  string
}

type Action struct {
	CommodityType proto.CommodityDTO_CommodityType
	Capacity      float64
	GetAppData    func(instanceName string) (AppData, error)
}

var Blueprint = map[Query]Action{
	"(navigation_timing_response_end_seconds-navigation_timing_request_start_seconds)*1000": {
		CommodityType: proto.CommodityDTO_RESPONSE_TIME,
		Capacity:      20,
		GetAppData: func(instanceName string) (AppData, error) {
			entityUrl, err := url.Parse(instanceName)
			if err != nil {
				glog.Errorf("Error instance field %v is not a valid URL: %s", instanceName, err)
				return AppData{}, err
			}
			appData := AppData{
				appType: AppType_Web,
				id:      AppType_Web + "_" + entityUrl.Hostname() + ":" + entityUrl.Port(),
				nodeIp:  entityUrl.Hostname(), // TODO perform DNS lookup here
			}
			return appData, nil
		},
	},
	"rate(mysql_global_status_queries[5m])": {
		CommodityType: proto.CommodityDTO_TRANSACTION,
		Capacity:      20,
		GetAppData: func(instanceName string) (AppData, error) {
			nodeIp := strings.Split(instanceName, ":")[0]
			return AppData{
				appType: AppType_MySQL,
				id:      AppType_MySQL + "_" + nodeIp,
				nodeIp:  nodeIp,
			}, nil
		},
	},
}
