package bindings

import "github.com/eni-chain/go-eni/x/tokenfactory/types"

type EniTokenFactoryQuery struct {
	// queries the tokenfactory authority metadata
	DenomAuthorityMetadata *types.QueryDenomAuthorityMetadataRequest `json:"denom_authority_metadata,omitempty"`
	// queries the tokenfactory denoms from a creator
	DenomsFromCreator *types.QueryDenomsFromCreatorRequest `json:"denoms_from_creator,omitempty"`
}
