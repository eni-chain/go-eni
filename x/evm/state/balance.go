package state

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"
)

func (s *DBImpl) SubBalance(evmAddr common.Address, amt *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	if amt.Sign() == 0 {
		return uint256.Int{}
	}
	if amt.Sign() < 0 {
		s.AddBalance(evmAddr, new(uint256.Int).Neg(amt), reason)
		return uint256.Int{}
	}
	ctx := s.ctx
	bigAmt := amt.ToBig()

	// this avoids emitting cosmos events for ephemeral bookkeeping transfers like send_native
	if s.eventsSuppressed {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
	}

	ueni, wei := SplitUeniWeiAmount(bigAmt)
	addr := s.getEniAddress(evmAddr)
	err := s.k.BankKeeper().SubUnlockedCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), ueni)), true)
	if err != nil {
		s.err = err
		return uint256.Int{}
	}
	err = s.k.BankKeeper().SubWei(ctx, addr, wei)
	if err != nil {
		s.err = err
		return uint256.Int{}
	}

	if s.logger != nil && s.logger.OnBalanceChange != nil {
		// We could modify AddWei instead so it returns us the old/new balance directly.
		newBalance := s.GetBalance(evmAddr).ToBig()
		oldBalance := new(big.Int).Add(newBalance, bigAmt)

		s.logger.OnBalanceChange(evmAddr, oldBalance, newBalance, reason)
	}

	s.tempStateCurrent.surplus = s.tempStateCurrent.surplus.Add(sdk.NewIntFromBigInt(bigAmt))
	return uint256.Int{}
}

func (s *DBImpl) AddBalance(evmAddr common.Address, amt *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	if amt.Sign() == 0 {
		return uint256.Int{}
	}
	if amt.Sign() < 0 {
		s.SubBalance(evmAddr, new(uint256.Int).Neg(amt), reason)
		return uint256.Int{}
	}

	ctx := s.ctx
	bigAmt := amt.ToBig()
	// this avoids emitting cosmos events for ephemeral bookkeeping transfers like send_native
	if s.eventsSuppressed {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
	}

	ueni, wei := SplitUeniWeiAmount(bigAmt)
	addr := s.getEniAddress(evmAddr)
	err := s.k.BankKeeper().AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), ueni)), true)
	if err != nil {
		s.err = err
		return uint256.Int{}
	}
	err = s.k.BankKeeper().AddWei(ctx, addr, wei)
	if err != nil {
		s.err = err
		return uint256.Int{}
	}

	if s.logger != nil && s.logger.OnBalanceChange != nil {
		// We could modify AddWei instead so it returns us the old/new balance directly.
		newBalance := s.GetBalance(evmAddr).ToBig()
		oldBalance := new(big.Int).Sub(newBalance, bigAmt)

		s.logger.OnBalanceChange(evmAddr, oldBalance, newBalance, reason)
	}

	s.tempStateCurrent.surplus = s.tempStateCurrent.surplus.Sub(sdk.NewIntFromBigInt(bigAmt))
	return uint256.Int{}
}

func (s *DBImpl) GetBalance(evmAddr common.Address) *uint256.Int {
	bigBalance := s.getBalance(evmAddr)
	balance, _ := uint256.FromBig(bigBalance)
	return balance
}
func (s *DBImpl) getBalance(evmAddr common.Address) *big.Int {
	eniAddr := s.getEniAddress(evmAddr)
	return s.k.GetBalance(s.ctx, eniAddr)
}

// should only be called during simulation
func (s *DBImpl) SetBalance(evmAddr common.Address, amt *uint256.Int, reason tracing.BalanceChangeReason) {
	if !s.simulation {
		panic("should never call SetBalance in a non-simulation setting")
	}
	eniAddr := s.getEniAddress(evmAddr)
	moduleAddr := s.k.AccountKeeper().GetModuleAddress(types.ModuleName)
	s.send(eniAddr, moduleAddr, s.getBalance(evmAddr))
	if s.err != nil {
		panic(s.err)
	}
	a := amt.ToBig()
	ueni, _ := SplitUeniWeiAmount(a)
	coinsAmt := sdk.NewCoins(sdk.NewCoin(s.k.GetBaseDenom(s.ctx), ueni.Add(sdk.OneInt())))
	if err := s.k.BankKeeper().MintCoins(s.ctx, types.ModuleName, coinsAmt); err != nil {
		panic(err)
	}
	s.send(moduleAddr, eniAddr, a)
	if s.err != nil {
		panic(s.err)
	}
}

func (s *DBImpl) getEniAddress(evmAddr common.Address) sdk.AccAddress {
	if s.coinbaseEvmAddress.Cmp(evmAddr) == 0 {
		return s.coinbaseAddress
	}
	return s.k.GetEniAddressOrDefault(s.ctx, evmAddr)
}

func (s *DBImpl) send(from sdk.AccAddress, to sdk.AccAddress, amt *big.Int) {
	ueni, wei := SplitUeniWeiAmount(amt)
	err := s.k.BankKeeper().SendCoinsAndWei(s.ctx, from, to, ueni, wei)
	if err != nil {
		s.err = err
	}
}
