package builder

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type ReplacementEntityMetaDataBuilder struct {
	metaData *proto.EntityDTO_ReplacementEntityMetaData
}

func NewReplacementEntityMetaDataBuilder() *ReplacementEntityMetaDataBuilder {
	replacementEntityMetaData := &proto.EntityDTO_ReplacementEntityMetaData{
		IdentifyingProp:  []string{},
		BuyingCommTypes:  []proto.CommodityDTO_CommodityType{},
		SellingCommTypes: []proto.CommodityDTO_CommodityType{},
	}
	return &ReplacementEntityMetaDataBuilder{
		metaData: replacementEntityMetaData,
	}
}

func (builder *ReplacementEntityMetaDataBuilder) Build() *proto.EntityDTO_ReplacementEntityMetaData {
	return builder.metaData
}

// Specifies the name of the property whose value will be used to find the server entity
// for which builder entity is a proxy. The value for the property must be set while building the
// entity.
// Specific properties are pre-defined for some entity types. See the constants defined in
// supply_chain_constants for the names of the specific properties.
func (builder *ReplacementEntityMetaDataBuilder) Matching(property string) *ReplacementEntityMetaDataBuilder {
	builder.metaData.IdentifyingProp = append(builder.metaData.GetIdentifyingProp(), property)
	return builder
}

// Set the commodity type whose metric values will be transferred to the entity
// builder DTO will be replaced by.
func (builder *ReplacementEntityMetaDataBuilder) PatchBuying(commType proto.CommodityDTO_CommodityType) *ReplacementEntityMetaDataBuilder {
	builder.metaData.BuyingCommTypes = append(builder.metaData.GetBuyingCommTypes(), commType)
	return builder
}

// Set the commodity type whose metric values will be transferred to the entity
//  builder DTO will be replaced by.
func (builder *ReplacementEntityMetaDataBuilder) PatchSelling(commType proto.CommodityDTO_CommodityType) *ReplacementEntityMetaDataBuilder {
	builder.metaData.SellingCommTypes = append(builder.metaData.GetSellingCommTypes(), commType)
	return builder
}
