// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// UniswapV4PoolABI is the input ABI used to generate the binding from.
const UniswapV4PoolABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token1\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reserve0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reserve1\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_blockTimestampLast\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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
	parsed, err := abi.JSON(strings.NewReader(UniswapV4PoolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV4Pool *UniswapV4PoolCaller) GetReserves(opts *bind.CallOpts) (*big.Int, *big.Int, uint32, error) {
	var out []interface{}
	err := _UniswapV4Pool.contract.Call(opts, &out, "getReserves")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(uint32), err
	}

	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int),
		*abi.ConvertType(out[1], new(*big.Int)).(**big.Int),
		*abi.ConvertType(out[2], new(uint32)).(*uint32),
		nil

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV4Pool *UniswapV4PoolSession) GetReserves() (*big.Int, *big.Int, uint32, error) {
	return _UniswapV4Pool.Contract.GetReserves(&_UniswapV4Pool.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 _reserve0, uint256 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV4Pool *UniswapV4PoolCallerSession) GetReserves() (*big.Int, *big.Int, uint32, error) {
	return _UniswapV4Pool.Contract.GetReserves(&_UniswapV4Pool.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolTransactor) Burn(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "burn", to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolSession) Burn(to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Burn(&_UniswapV4Pool.TransactOpts, to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Burn(to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Burn(&_UniswapV4Pool.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV4Pool *UniswapV4PoolTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV4Pool *UniswapV4PoolSession) Mint(to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Mint(&_UniswapV4Pool.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Mint(&_UniswapV4Pool.TransactOpts, to)
}

// Swap is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactor) Swap(opts *bind.TransactOpts, amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.contract.Transact(opts, "swap", amount0Out, amount1Out, to)
}

// Swap is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Swap(&_UniswapV4Pool.TransactOpts, amount0Out, amount1Out, to)
}

// Swap is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to) returns()
func (_UniswapV4Pool *UniswapV4PoolTransactorSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV4Pool.Contract.Swap(&_UniswapV4Pool.TransactOpts, amount0Out, amount1Out, to)
}

// UniswapV4PoolBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the UniswapV4Pool contract.
type UniswapV4PoolBurnIterator struct {
	Event *UniswapV4PoolBurn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for events, 'quit' is on the first error
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

// FilterBurn is a free log retrieval operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
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

// WatchBurn is a free log subscription operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
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

// ParseBurn is a log parse operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
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
	sub  ethereum.Subscription // Subscription for events, 'quit' is on the first error
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

// UniswapV4PoolSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the UniswapV4Pool contract.
type UniswapV4PoolSwapIterator struct {
	Event *UniswapV4PoolSwap // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for events, 'quit' is on the first error
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
	Sender    common.Address
	Amount0In *big.Int
	Amount1In *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To        common.Address
	Raw       types.Log // Blockchain specific contextual infos
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