package registration

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	vCpuType     = proto.CommodityDTO_VCPU
	vMemType     = proto.CommodityDTO_VMEM
	respTimeType = proto.CommodityDTO_RESPONSE_TIME
	transactionType = proto.CommodityDTO_TRANSACTION

	vCpuTemplateComm     *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vCpuType}
	vMemTemplateComm     *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vMemType}
	respTimeTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &respTimeType}
	transactionTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &transactionType}
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

	return supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_MACHINE).
		Sells(vCpuTemplateComm).
		Sells(vMemTemplateComm).
		Create()
}

func (factory *SupplyChainFactory) buildAppSupplyBuilder() (*proto.TemplateDTO, error) {
	//create
	appSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_APPLICATION)
		//Sells(respTimeTemplateComm).Sells(transactionTemplateComm)
		//Provider(proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING).
		//Buys(vCpuTemplateComm).
		//Buys(vMemTemplateComm)

	// Link from App to VM
	vmAppExtLinkBuilder := supplychain.NewExternalEntityLinkBuilder()
	vmAppExtLinkBuilder.Link(proto.EntityDTO_APPLICATION, proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING).
		Commodity(vCpuType, false).Commodity(vMemType, false).
		ProbeEntityPropertyDef(supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS, "IP Address where the app is running").
		ExternalEntityPropertyDef(supplychain.VM_IP)
	vmAppExternalLink, err := vmAppExtLinkBuilder.Build()
	if err != nil {
		return nil, err
	}

	return appSupplyChainNodeBuilder.ConnectsTo(vmAppExternalLink).Create()
}
