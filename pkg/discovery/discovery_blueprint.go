package discovery

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"

type AppType string

const (
	AppType_Web   AppType = "web"
	AppType_MySQL AppType = "mysql"
)

type Query string

type Action struct {
	CommodityType proto.CommodityDTO_CommodityType
	Capacity      float64
	AppType       AppType
}

var Blueprint = map[Query]Action{
	"(navigation_timing_response_end_seconds-navigation_timing_request_start_seconds)*1000": Action{
		CommodityType: proto.CommodityDTO_RESPONSE_TIME,
		Capacity:      10,
		AppType:       AppType_Web,
	},
	"rate(mysql_global_status_queries[5m])": Action{
		CommodityType: proto.CommodityDTO_TRANSACTION,
		Capacity:      100,
		AppType:       AppType_MySQL,
	},
}
