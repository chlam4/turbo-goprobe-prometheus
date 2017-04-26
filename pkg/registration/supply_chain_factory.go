package registration

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	respTimeType                                     = proto.CommodityDTO_RESPONSE_TIME

	respTimeTemplateComm    *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &respTimeType}
)

type SupplyChainFactory struct{}

func (factory *SupplyChainFactory) CreateSupplyChain() ([]*proto.TemplateDTO, error) {

	vmNode, err := factory.buildNodeSupplyBuilder()
	if err != nil {
		return nil, err
	}

	appNode, err := factory.buildAppSupplyBuilder()
	if err != nil {
		return nil, err
	}

	return supplychain.NewSupplyChainBuilder().Top(appNode).Entity(vmNode).Create()
}

func (factory *SupplyChainFactory) buildNodeSupplyBuilder() (*proto.TemplateDTO, error) {

	return supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_MACHINE).Create()
}

func (factory *SupplyChainFactory) buildAppSupplyBuilder() (*proto.TemplateDTO, error) {
	appSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_APPLICATION).
		Sells(respTimeTemplateComm).
		Provider(proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING)

	// Link from App to VM
	vmAppExtLinkBuilder := supplychain.NewExternalEntityLinkBuilder()
	vmAppExtLinkBuilder.Link(proto.EntityDTO_APPLICATION, proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING).
		ProbeEntityPropertyDef(supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS, "IP Address where the app is running").
		ExternalEntityPropertyDef(supplychain.VM_IP)
	vmAppExternalLink, err := vmAppExtLinkBuilder.Build()
	if err != nil {
		return nil, err
	}

	return appSupplyChainNodeBuilder.ConnectsTo(vmAppExternalLink).Create()
}
