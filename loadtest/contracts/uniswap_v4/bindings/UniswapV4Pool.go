// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswap

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// UniswapV4PoolMetaData contains all meta data concerning the UniswapV4Pool contract.
var UniswapV4PoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token1\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reserve0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserve1\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// UniswapV4PoolABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV4PoolMetaData.ABI instead.
var UniswapV4PoolABI = UniswapV4PoolMetaData.ABI

// UniswapV4Pool is an auto generated Go binding around an Ethereum contract.
type UniswapV4Pool struct {
	UniswapV4PoolCaller     // Read-only binding to the contract
	UniswapV4PoolTransactor // Write-only binding to the contract
	UniswapV4PoolFilterer   // Log filterer for contract events
}

// UniswapV4PoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV4PoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV4PoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV4PoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV4PoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV4PoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV4PoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV4PoolSession struct {
	Contract     *UniswapV4Pool    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV4PoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV4PoolCallerSession struct {
	Contract *UniswapV4PoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// UniswapV4PoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV4PoolTransactorSession struct {
	Contract     *UniswapV4PoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// UniswapV4PoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV4PoolRaw struct {
	Contract *UniswapV4Pool // Generic contract binding to access the raw methods on
}

// UniswapV4PoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV4PoolCallerRaw struct {
	Contract *UniswapV4PoolCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV4PoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV4PoolTransactorRaw struct {
	Contract *UniswapV4PoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV4Pool creates a new instance of UniswapV4Pool, bound to a specific deployed contract.
func NewUniswapV4Pool(address common.Address, backend bind.ContractBackend) (*UniswapV4Pool, error) {
	contract, err := bindUniswapV4Pool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV4Pool{UniswapV4PoolCaller: UniswapV4PoolCaller{contract: contract}, UniswapV4PoolTransactor: UniswapV4PoolTransactor{contract: contract}, UniswapV4PoolFilterer: UniswapV4PoolFilterer{contract: contract}}, nil
}

// NewUniswapV4PoolCaller creates a new read-only instance of UniswapV4Pool, bound to a specific deployed contract.
func NewUniswapV4PoolCaller(address common.Address, caller bind.ContractCaller) (*UniswapV4PoolCaller, error) {
	contract, err := bindUniswapV4Pool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolCaller{contract: contract}, nil
}

// NewUniswapV4PoolTransactor creates a new write-only instance of UniswapV4Pool, bound to a specific deployed contract.
func NewUniswapV4PoolTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV4PoolTransactor, error) {
	contract, err := bindUniswapV4Pool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolTransactor{contract: contract}, nil
}

// NewUniswapV4PoolFilterer creates a new log filterer instance of UniswapV4Pool, bound to a specific deployed contract.
func NewUniswapV4PoolFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV4PoolFilterer, error) {
	contract, err := bindUniswapV4Pool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolFilterer{contract: contract}, nil
}

// bindUniswapV4Pool binds a generic wrapper to an already deployed contract.
func bindUniswapV4Pool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV4PoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV4Pool *UniswapV4PoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV4Pool.Contract.UniswapV4PoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV4Pool *UniswapV4PoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.UniswapV4PoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV4Pool *UniswapV4PoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.UniswapV4PoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV4Pool *UniswapV4PoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV4Pool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV4Pool *UniswapV4PoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV4Pool *UniswapV4PoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.contract.Transact(opts, method, params...)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1)
func (_UniswapV4Pool *UniswapV4PoolCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0 *big.Int
		Reserve1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1)
func (_UniswapV4Pool *UniswapV4PoolSession) GetReserves() (struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}, error) {
	return _UniswapV4Pool.Contract.GetReserves(&_UniswapV4Pool.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) GetReserves() (struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}, error) {
	return _UniswapV4Pool.Contract.GetReserves(&_UniswapV4Pool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolSession) Owner() (common.Address, error) {
	return _UniswapV4Pool.Contract.Owner(&_UniswapV4Pool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) Owner() (common.Address, error) {
	return _UniswapV4Pool.Contract.Owner(&_UniswapV4Pool.CallOpts)
}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolCaller) Reserve0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "reserve0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolSession) Reserve0() (*big.Int, error) {
	return _UniswapV4Pool.Contract.Reserve0(&_UniswapV4Pool.CallOpts)
}

// Reserve0 is a free data retrieval call binding the contract method 0x443cb4bc.
//
// Solidity: function reserve0() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) Reserve0() (*big.Int, error) {
	return _UniswapV4Pool.Contract.Reserve0(&_UniswapV4Pool.CallOpts)
}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolCaller) Reserve1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "reserve1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolSession) Reserve1() (*big.Int, error) {
	return _UniswapV4Pool.Contract.Reserve1(&_UniswapV4Pool.CallOpts)
}

// Reserve1 is a free data retrieval call binding the contract method 0x5a76f25e.
//
// Solidity: function reserve1() view returns(uint256)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) Reserve1() (*big.Int, error) {
	return _UniswapV4Pool.Contract.Reserve1(&_UniswapV4Pool.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolSession) Token0() (common.Address, error) {
	return _UniswapV4Pool.Contract.Token0(&_UniswapV4Pool.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) Token0() (common.Address, error) {
	return _UniswapV4Pool.Contract.Token0(&_UniswapV4Pool.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolSession) Token1() (common.Address, error) {
	return _UniswapV4Pool.Contract.Token1(&_UniswapV4Pool.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) Token1() (common.Address, error) {
	return _UniswapV4Pool.Contract.Token1(&_UniswapV4Pool.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x749388c4.
//
// Solidity: function burn(uint256 amount0, uint256 amount1, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) Burn(opts *bind.TransactOpts, amount0 *big.Int, amount1 *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "burn", amount0, amount1, to)
}

// Burn is a paid mutator transaction binding the contract method 0x749388c4.
//
// Solidity: function burn(uint256 amount0, uint256 amount1, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolSession) Burn(amount0 *big.Int, amount1 *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Burn(&_UniswapV4Pool.TransactOpts, amount0, amount1, to)
}

// Burn is a paid mutator transaction binding the contract method 0x749388c4.
//
// Solidity: function burn(uint256 amount0, uint256 amount1, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Burn(amount0 *big.Int, amount1 *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Burn(&_UniswapV4Pool.TransactOpts, amount0, amount1, to)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 amount0, uint256 amount1) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) Mint(opts *bind.TransactOpts, amount0 *big.Int, amount1 *big.Int) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "mint", amount0, amount1)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 amount0, uint256 amount1) returns()
func (_UniswapV4Pool *UniswapV4PoolSession) Mint(amount0 *big.Int, amount1 *big.Int) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Mint(&_UniswapV4Pool.TransactOpts, amount0, amount1)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 amount0, uint256 amount1) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Mint(amount0 *big.Int, amount1 *big.Int) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Mint(&_UniswapV4Pool.TransactOpts, amount0, amount1)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV4Pool *UniswapV4PoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.RenounceOwnership(&_UniswapV4Pool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.RenounceOwnership(&_UniswapV4Pool.TransactOpts)
}

// Swap is a paid mutator transaction binding the contract method 0x562e19df.
//
// Solidity: function swap(uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) Swap(opts *bind.TransactOpts, amount0In *big.Int, amount1In *big.Int, amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "swap", amount0In, amount1In, amount0Out, amount1Out, to)
}

// Swap is a paid mutator transaction binding the contract method 0x562e19df.
//
// Solidity: function swap(uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolSession) Swap(amount0In *big.Int, amount1In *big.Int, amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Swap(&_UniswapV4Pool.TransactOpts, amount0In, amount1In, amount0Out, amount1Out, to)
}

// Swap is a paid mutator transaction binding the contract method 0x562e19df.
//
// Solidity: function swap(uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Swap(amount0In *big.Int, amount1In *big.Int, amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Swap(&_UniswapV4Pool.TransactOpts, amount0In, amount1In, amount0Out, amount1Out, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV4Pool *UniswapV4PoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.TransferOwnership(&_UniswapV4Pool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.TransferOwnership(&_UniswapV4Pool.TransactOpts, newOwner)
}

// UniswapV4PoolBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the UniswapV4Pool contract.
type UniswapV4PoolBurnIterator struct {
	Event *UniswapV4PoolBurn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UniswapV4PoolBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV4PoolBurn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UniswapV4PoolBurn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UniswapV4PoolBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV4PoolBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV4PoolBurn represents a Burn event raised by the UniswapV4Pool contract.
type UniswapV4PoolBurn struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) FilterBurn(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*UniswapV4PoolBurnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.FilterLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolBurnIterator{contract: _UniswapV4Pool.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *UniswapV4PoolBurn, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.WatchLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV4PoolBurn)
				if err := _UniswapV4Pool.contract.UnpackLog(event, "Burn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBurn is a log parse operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) ParseBurn(log types.Log) (*UniswapV4PoolBurn, error) {
	event := new(UniswapV4PoolBurn)
	if err := _UniswapV4Pool.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV4PoolMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the UniswapV4Pool contract.
type UniswapV4PoolMintIterator struct {
	Event *UniswapV4PoolMint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UniswapV4PoolMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV4PoolMint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UniswapV4PoolMint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UniswapV4PoolMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV4PoolMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV4PoolMint represents a Mint event raised by the UniswapV4Pool contract.
type UniswapV4PoolMint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolFilterer) FilterMint(opts *bind.FilterOpts, sender []common.Address) (*UniswapV4PoolMintIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.FilterLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolMintIterator{contract: _UniswapV4Pool.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *UniswapV4PoolMint, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.WatchLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV4PoolMint)
				if err := _UniswapV4Pool.contract.UnpackLog(event, "Mint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolFilterer) ParseMint(log types.Log) (*UniswapV4PoolMint, error) {
	event := new(UniswapV4PoolMint)
	if err := _UniswapV4Pool.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV4PoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UniswapV4Pool contract.
type UniswapV4PoolOwnershipTransferredIterator struct {
	Event *UniswapV4PoolOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UniswapV4PoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV4PoolOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UniswapV4PoolOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UniswapV4PoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV4PoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV4PoolOwnershipTransferred represents a OwnershipTransferred event raised by the UniswapV4Pool contract.
type UniswapV4PoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UniswapV4Pool *UniswapV4PoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UniswapV4PoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolOwnershipTransferredIterator{contract: _UniswapV4Pool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UniswapV4Pool *UniswapV4PoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UniswapV4PoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV4PoolOwnershipTransferred)
				if err := _UniswapV4Pool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UniswapV4Pool *UniswapV4PoolFilterer) ParseOwnershipTransferred(log types.Log) (*UniswapV4PoolOwnershipTransferred, error) {
	event := new(UniswapV4PoolOwnershipTransferred)
	if err := _UniswapV4Pool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV4PoolSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the UniswapV4Pool contract.
type UniswapV4PoolSwapIterator struct {
	Event *UniswapV4PoolSwap // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UniswapV4PoolSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV4PoolSwap)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UniswapV4PoolSwap)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UniswapV4PoolSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV4PoolSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV4PoolSwap represents a Swap event raised by the UniswapV4Pool contract.
type UniswapV4PoolSwap struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*UniswapV4PoolSwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV4PoolSwapIterator{contract: _UniswapV4Pool.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *UniswapV4PoolSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV4Pool.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV4PoolSwap)
				if err := _UniswapV4Pool.contract.UnpackLog(event, "Swap", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwap is a log parse operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV4Pool *UniswapV4PoolFilterer) ParseSwap(log types.Log) (*UniswapV4PoolSwap, error) {
	event := new(UniswapV4PoolSwap)
	if err := _UniswapV4Pool.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
