package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	appparams "github.com/eni-chain/go-eni/app/params"
	"github.com/eni-chain/go-eni/x/tokenfactory/types"
)

func TestDecomposeDenoms(t *testing.T) {
	appparams.SetAddressPrefixes()
	for _, tc := range []struct {
		desc  string
		denom string
		valid bool
	}{
		{
			desc:  "empty is invalid",
			denom: "",
			valid: false,
		},
		{
			desc:  "normal",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/bitcoin",
			valid: true,
		},
		{
			desc:  "multiple slashes in subdenom",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/bitcoin/1",
			valid: true,
		},
		{
			desc:  "no subdenom",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/",
			valid: true,
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/bitcoin",
			valid: false,
		},
		{
			desc:  "subdenom of only slashes",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/////",
			valid: true,
		},
		{
			desc:  "too long name",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			valid: false,
		},
		{
			desc:  "too long creator name",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3fasdfasdfasdfasdfasdfasdfadfasdfasdfasdfasdfasdfas/bitcoin",
			valid: false,
		},
		{
			desc:  "empty subdenom",
			denom: "factory/eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f/",
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, _, err := types.DeconstructDenom(tc.denom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGetTokenDenom(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		creator  string
		subdenom string
		valid    bool
	}{
		{
			desc:     "normal",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "bitcoin",
			valid:    true,
		},
		{
			desc:     "multiple slashes in subdenom",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "bitcoin/1",
			valid:    true,
		},
		{
			desc:     "no subdenom",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "",
			valid:    true,
		},
		{
			desc:     "subdenom of only slashes",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "/////",
			valid:    true,
		},
		{
			desc:     "too long name",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			valid:    false,
		},
		{
			desc:     "subdenom is exactly max length",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "bitcoinfsadfsdfeadfsafwefsefsefsdfsdafasefsf",
			valid:    true,
		},
		{
			desc:     "creator is exactly max length",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3fhjkljkljkljkljkljkljkljkljkljkljk",
			subdenom: "bitcoin",
			valid:    true,
		},
		{
			desc:     "empty subdenom",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "",
			valid:    true,
		},
		{
			desc:     "non standard UTF-8",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "\u2603",
			valid:    false,
		},
		{
			desc:     "non standard ASCII",
			creator:  "eni1y3pxq5dp900czh0mkudhjdqjq5m8cpmmujwv3f",
			subdenom: "\n\t",
			valid:    false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := types.GetTokenDenom(tc.creator, tc.subdenom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
