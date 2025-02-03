package bindings

import "github.com/eni-chain/go-eni/x/epoch/types"

type EniEpochQuery struct {
	// queries the current Epoch
	Epoch *types.QueryEpochRequest `json:"epoch,omitempty"`
}
