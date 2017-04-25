package registration

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	builder "github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	respTimeType  = proto.CommodityDTO_RESPONSE_TIME
	respTimeTemplateComm  *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &respTimeType}
)

type SupplyChainFactory struct{}

func (this *SupplyChainFactory) CreateSupplyChain() ([]*proto.TemplateDTO, error) {

	appNode, _ := builder.NewSupplyChainNodeBuilder(proto.EntityDTO_APPLICATION).Buys(respTimeTemplateComm).Create()

	return builder.NewSupplyChainBuilder().
		Top(appNode).
		Create()
}

