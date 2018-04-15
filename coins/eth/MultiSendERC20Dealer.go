// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// MultiSendERC20DealerABI is the input ABI used to generate the binding from.
const MultiSendERC20DealerABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenDropLog\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_tokenAddr\",\"type\":\"address\"},{\"name\":\"_tokenSupplier\",\"type\":\"address\"},{\"name\":\"_dealer\",\"type\":\"address\"},{\"name\":\"_price\",\"type\":\"uint256\"},{\"name\":\"recipients\",\"type\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"multiSend\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]"

const MultiSendERC20DealerBin = "0x606060405260008054600160a060020a033316600160a060020a0319909116179055610386806100306000396000f3006060604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416638da5cb5b81146100585780639efb5a5514610087578063f2fde38b146100c4575b005b341561006357600080fd5b61006b6100e3565b604051600160a060020a03909116815260200160405180910390f35b61005660048035600160a060020a03908116916024803583169260443516916064359160843580820192908101359160a4359081019101356100f2565b34156100cf57600080fd5b610056600160a060020a03600435166102bf565b600054600160a060020a031681565b60008054819033600160a060020a0390811691161461011057600080fd5b84831461011c57600080fd5b603285111561012a57600080fd5b600160a060020a03881687156108fc0288604051600060405180830381858888f19350505050151561015b57600080fd5b5088905060005b848110156102b357600160a060020a0382166323b872dd8a88888581811061018657fe5b90506020020135600160a060020a031687878681811015156101a457fe5b905060200201356000604051602001526040517c010000000000000000000000000000000000000000000000000000000063ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401602060405180830381600087803b151561021c57600080fd5b6102c65a03f1151561022d57600080fd5b5050506040518051507f873074ca4c984cb6f1fb306295e51b68a065cf15a5fc80cea5d4655ab32300ab905086868381811061026557fe5b90506020020135600160a060020a0316858584818110151561028357fe5b90506020020135604051600160a060020a03909216825260208201526040908101905180910390a1600101610162565b50505050505050505050565b60005433600160a060020a039081169116146102da57600080fd5b600160a060020a03811615156102ef57600080fd5b600054600160a060020a0380831691167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03929092169190911790555600a165627a7a72305820ad439e4bf020a6046efeac395e5c539cf8281d84e21f46f57fc34c8686c635370029"

func DeployMultiSendERC20Dealer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultiSendERC20Dealer, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiSendERC20DealerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MultiSendERC20DealerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiSendERC20Dealer{MultiSendERC20DealerCaller: MultiSendERC20DealerCaller{contract: contract}, MultiSendERC20DealerTransactor: MultiSendERC20DealerTransactor{contract: contract}, MultiSendERC20DealerFilterer: MultiSendERC20DealerFilterer{contract: contract}}, nil
}

// MultiSendERC20Dealer is an auto generated Go binding around an Ethereum contract.
type MultiSendERC20Dealer struct {
	MultiSendERC20DealerCaller     // Read-only binding to the contract
	MultiSendERC20DealerTransactor // Write-only binding to the contract
	MultiSendERC20DealerFilterer   // Log filterer for contract events
}

// MultiSendERC20DealerCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiSendERC20DealerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendERC20DealerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiSendERC20DealerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendERC20DealerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiSendERC20DealerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendERC20DealerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiSendERC20DealerSession struct {
	Contract     *MultiSendERC20Dealer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MultiSendERC20DealerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiSendERC20DealerCallerSession struct {
	Contract *MultiSendERC20DealerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// MultiSendERC20DealerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiSendERC20DealerTransactorSession struct {
	Contract     *MultiSendERC20DealerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// MultiSendERC20DealerRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiSendERC20DealerRaw struct {
	Contract *MultiSendERC20Dealer // Generic contract binding to access the raw methods on
}

// MultiSendERC20DealerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiSendERC20DealerCallerRaw struct {
	Contract *MultiSendERC20DealerCaller // Generic read-only contract binding to access the raw methods on
}

// MultiSendERC20DealerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiSendERC20DealerTransactorRaw struct {
	Contract *MultiSendERC20DealerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiSendERC20Dealer creates a new instance of MultiSendERC20Dealer, bound to a specific deployed contract.
func NewMultiSendERC20Dealer(address common.Address, backend bind.ContractBackend) (*MultiSendERC20Dealer, error) {
	contract, err := bindMultiSendERC20Dealer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20Dealer{MultiSendERC20DealerCaller: MultiSendERC20DealerCaller{contract: contract}, MultiSendERC20DealerTransactor: MultiSendERC20DealerTransactor{contract: contract}, MultiSendERC20DealerFilterer: MultiSendERC20DealerFilterer{contract: contract}}, nil
}

// NewMultiSendERC20DealerCaller creates a new read-only instance of MultiSendERC20Dealer, bound to a specific deployed contract.
func NewMultiSendERC20DealerCaller(address common.Address, caller bind.ContractCaller) (*MultiSendERC20DealerCaller, error) {
	contract, err := bindMultiSendERC20Dealer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20DealerCaller{contract: contract}, nil
}

// NewMultiSendERC20DealerTransactor creates a new write-only instance of MultiSendERC20Dealer, bound to a specific deployed contract.
func NewMultiSendERC20DealerTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiSendERC20DealerTransactor, error) {
	contract, err := bindMultiSendERC20Dealer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20DealerTransactor{contract: contract}, nil
}

// NewMultiSendERC20DealerFilterer creates a new log filterer instance of MultiSendERC20Dealer, bound to a specific deployed contract.
func NewMultiSendERC20DealerFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiSendERC20DealerFilterer, error) {
	contract, err := bindMultiSendERC20Dealer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20DealerFilterer{contract: contract}, nil
}

// bindMultiSendERC20Dealer binds a generic wrapper to an already deployed contract.
func bindMultiSendERC20Dealer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiSendERC20DealerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSendERC20Dealer *MultiSendERC20DealerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiSendERC20Dealer.Contract.MultiSendERC20DealerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSendERC20Dealer *MultiSendERC20DealerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.MultiSendERC20DealerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSendERC20Dealer *MultiSendERC20DealerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.MultiSendERC20DealerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSendERC20Dealer *MultiSendERC20DealerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiSendERC20Dealer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultiSendERC20Dealer *MultiSendERC20DealerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MultiSendERC20Dealer.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultiSendERC20Dealer *MultiSendERC20DealerSession) Owner() (common.Address, error) {
	return _MultiSendERC20Dealer.Contract.Owner(&_MultiSendERC20Dealer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultiSendERC20Dealer *MultiSendERC20DealerCallerSession) Owner() (common.Address, error) {
	return _MultiSendERC20Dealer.Contract.Owner(&_MultiSendERC20Dealer.CallOpts)
}

// MultiSend is a paid mutator transaction binding the contract method 0x9efb5a55.
//
// Solidity: function multiSend(_tokenAddr address, _tokenSupplier address, _dealer address, _price uint256, recipients address[], amounts uint256[]) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactor) MultiSend(opts *bind.TransactOpts, _tokenAddr common.Address, _tokenSupplier common.Address, _dealer common.Address, _price *big.Int, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.contract.Transact(opts, "multiSend", _tokenAddr, _tokenSupplier, _dealer, _price, recipients, amounts)
}

// MultiSend is a paid mutator transaction binding the contract method 0x9efb5a55.
//
// Solidity: function multiSend(_tokenAddr address, _tokenSupplier address, _dealer address, _price uint256, recipients address[], amounts uint256[]) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerSession) MultiSend(_tokenAddr common.Address, _tokenSupplier common.Address, _dealer common.Address, _price *big.Int, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.MultiSend(&_MultiSendERC20Dealer.TransactOpts, _tokenAddr, _tokenSupplier, _dealer, _price, recipients, amounts)
}

// MultiSend is a paid mutator transaction binding the contract method 0x9efb5a55.
//
// Solidity: function multiSend(_tokenAddr address, _tokenSupplier address, _dealer address, _price uint256, recipients address[], amounts uint256[]) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactorSession) MultiSend(_tokenAddr common.Address, _tokenSupplier common.Address, _dealer common.Address, _price *big.Int, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.MultiSend(&_MultiSendERC20Dealer.TransactOpts, _tokenAddr, _tokenSupplier, _dealer, _price, recipients, amounts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.TransferOwnership(&_MultiSendERC20Dealer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultiSendERC20Dealer *MultiSendERC20DealerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultiSendERC20Dealer.Contract.TransferOwnership(&_MultiSendERC20Dealer.TransactOpts, newOwner)
}

// MultiSendERC20DealerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MultiSendERC20Dealer contract.
type MultiSendERC20DealerOwnershipTransferredIterator struct {
	Event *MultiSendERC20DealerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MultiSendERC20DealerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiSendERC20DealerOwnershipTransferred)
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
		it.Event = new(MultiSendERC20DealerOwnershipTransferred)
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
func (it *MultiSendERC20DealerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiSendERC20DealerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiSendERC20DealerOwnershipTransferred represents a OwnershipTransferred event raised by the MultiSendERC20Dealer contract.
type MultiSendERC20DealerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_MultiSendERC20Dealer *MultiSendERC20DealerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultiSendERC20DealerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiSendERC20Dealer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20DealerOwnershipTransferredIterator{contract: _MultiSendERC20Dealer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_MultiSendERC20Dealer *MultiSendERC20DealerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultiSendERC20DealerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultiSendERC20Dealer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiSendERC20DealerOwnershipTransferred)
				if err := _MultiSendERC20Dealer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// MultiSendERC20DealerTokenDropLogIterator is returned from FilterTokenDropLog and is used to iterate over the raw logs and unpacked data for TokenDropLog events raised by the MultiSendERC20Dealer contract.
type MultiSendERC20DealerTokenDropLogIterator struct {
	Event *MultiSendERC20DealerTokenDropLog // Event containing the contract specifics and raw log

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
func (it *MultiSendERC20DealerTokenDropLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiSendERC20DealerTokenDropLog)
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
		it.Event = new(MultiSendERC20DealerTokenDropLog)
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
func (it *MultiSendERC20DealerTokenDropLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiSendERC20DealerTokenDropLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiSendERC20DealerTokenDropLog represents a TokenDropLog event raised by the MultiSendERC20Dealer contract.
type MultiSendERC20DealerTokenDropLog struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenDropLog is a free log retrieval operation binding the contract event 0x873074ca4c984cb6f1fb306295e51b68a065cf15a5fc80cea5d4655ab32300ab.
//
// Solidity: event TokenDropLog(receiver address, amount uint256)
func (_MultiSendERC20Dealer *MultiSendERC20DealerFilterer) FilterTokenDropLog(opts *bind.FilterOpts) (*MultiSendERC20DealerTokenDropLogIterator, error) {

	logs, sub, err := _MultiSendERC20Dealer.contract.FilterLogs(opts, "TokenDropLog")
	if err != nil {
		return nil, err
	}
	return &MultiSendERC20DealerTokenDropLogIterator{contract: _MultiSendERC20Dealer.contract, event: "TokenDropLog", logs: logs, sub: sub}, nil
}

// WatchTokenDropLog is a free log subscription operation binding the contract event 0x873074ca4c984cb6f1fb306295e51b68a065cf15a5fc80cea5d4655ab32300ab.
//
// Solidity: event TokenDropLog(receiver address, amount uint256)
func (_MultiSendERC20Dealer *MultiSendERC20DealerFilterer) WatchTokenDropLog(opts *bind.WatchOpts, sink chan<- *MultiSendERC20DealerTokenDropLog) (event.Subscription, error) {

	logs, sub, err := _MultiSendERC20Dealer.contract.WatchLogs(opts, "TokenDropLog")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiSendERC20DealerTokenDropLog)
				if err := _MultiSendERC20Dealer.contract.UnpackLog(event, "TokenDropLog", log); err != nil {
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
