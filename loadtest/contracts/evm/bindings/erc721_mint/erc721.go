// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc721

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

// Erc721MetaData contains all meta data concerning the Erc721 contract.
var Erc721MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextTokenId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001600755348015610014575f5ffd5b50604051612800380380612800833981810160405281019061003691906101ad565b8181815f90816100469190610433565b5080600190816100569190610433565b5050505050610502565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6100bf82610079565b810181811067ffffffffffffffff821117156100de576100dd610089565b5b80604052505050565b5f6100f0610060565b90506100fc82826100b6565b919050565b5f67ffffffffffffffff82111561011b5761011a610089565b5b61012482610079565b9050602081019050919050565b8281835e5f83830152505050565b5f61015161014c84610101565b6100e7565b90508281526020810184848401111561016d5761016c610075565b5b610178848285610131565b509392505050565b5f82601f83011261019457610193610071565b5b81516101a484826020860161013f565b91505092915050565b5f5f604083850312156101c3576101c2610069565b5b5f83015167ffffffffffffffff8111156101e0576101df61006d565b5b6101ec85828601610180565b925050602083015167ffffffffffffffff81111561020d5761020c61006d565b5b61021985828601610180565b9150509250929050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061027157607f821691505b6020821081036102845761028361022d565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026102e67fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826102ab565b6102f086836102ab565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61033461032f61032a84610308565b610311565b610308565b9050919050565b5f819050919050565b61034d8361031a565b6103616103598261033b565b8484546102b7565b825550505050565b5f5f905090565b610378610369565b610383818484610344565b505050565b5b818110156103a65761039b5f82610370565b600181019050610389565b5050565b601f8211156103eb576103bc8161028a565b6103c58461029c565b810160208510156103d4578190505b6103e86103e08561029c565b830182610388565b50505b505050565b5f82821c905092915050565b5f61040b5f19846008026103f0565b1980831691505092915050565b5f61042383836103fc565b9150826002028217905092915050565b61043c82610223565b67ffffffffffffffff81111561045557610454610089565b5b61045f825461025a565b61046a8282856103aa565b5f60209050601f83116001811461049b575f8415610489578287015190505b6104938582610418565b8655506104fa565b601f1984166104a98661028a565b5f5b828110156104d0578489015182556001820191506020850194506020810190506104ab565b868310156104ed57848901516104e9601f8916826103fc565b8355505b6001600288020188555050505b505050505050565b6122f18061050f5f395ff3fe608060405234801561000f575f5ffd5b50600436106100f3575f3560e01c806370a0823111610095578063b88d4fde11610064578063b88d4fde14610281578063c87b56dd1461029d578063d0def521146102cd578063e985e9c5146102e9576100f3565b806370a08231146101f957806375794a3c1461022957806395d89b4114610247578063a22cb46514610265576100f3565b8063095ea7b3116100d1578063095ea7b31461017557806323b872dd1461019157806342842e0e146101ad5780636352211e146101c9576100f3565b806301ffc9a7146100f757806306fdde0314610127578063081812fc14610145575b5f5ffd5b610111600480360381019061010c91906117c1565b610319565b60405161011e9190611806565b60405180910390f35b61012f610379565b60405161013c919061188f565b60405180910390f35b61015f600480360381019061015a91906118e2565b610408565b60405161016c919061194c565b60405180910390f35b61018f600480360381019061018a919061198f565b610423565b005b6101ab60048036038101906101a691906119cd565b610439565b005b6101c760048036038101906101c291906119cd565b610538565b005b6101e360048036038101906101de91906118e2565b610557565b6040516101f0919061194c565b60405180910390f35b610213600480360381019061020e9190611a1d565b610568565b6040516102209190611a57565b60405180910390f35b61023161061e565b60405161023e9190611a57565b60405180910390f35b61024f610624565b60405161025c919061188f565b60405180910390f35b61027f600480360381019061027a9190611a9a565b6106b4565b005b61029b60048036038101906102969190611c04565b6106ca565b005b6102b760048036038101906102b291906118e2565b6106ef565b6040516102c4919061188f565b60405180910390f35b6102e760048036038101906102e29190611d22565b6107fa565b005b61030360048036038101906102fe9190611d7c565b61082c565b6040516103109190611806565b60405180910390f35b5f634906490660e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806103725750610371826108ba565b5b9050919050565b60605f805461038790611de7565b80601f01602080910402602001604051908101604052809291908181526020018280546103b390611de7565b80156103fe5780601f106103d5576101008083540402835291602001916103fe565b820191905f5260205f20905b8154815290600101906020018083116103e157829003601f168201915b5050505050905090565b5f6104128261099b565b5061041c82610a21565b9050919050565b6104358282610430610a5a565b610a61565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036104a9575f6040517f64a0ae920000000000000000000000000000000000000000000000000000000081526004016104a0919061194c565b60405180910390fd5b5f6104bc83836104b7610a5a565b610a73565b90508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610532578382826040517f64283d7b00000000000000000000000000000000000000000000000000000000815260040161052993929190611e17565b60405180910390fd5b50505050565b61055283838360405180602001604052805f8152506106ca565b505050565b5f6105618261099b565b9050919050565b5f5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036105d9575f6040517f89c62b640000000000000000000000000000000000000000000000000000000081526004016105d0919061194c565b60405180910390fd5b60035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b60075481565b60606001805461063390611de7565b80601f016020809104026020016040519081016040528092919081815260200182805461065f90611de7565b80156106aa5780601f10610681576101008083540402835291602001916106aa565b820191905f5260205f20905b81548152906001019060200180831161068d57829003601f168201915b5050505050905090565b6106c66106bf610a5a565b8383610c7e565b5050565b6106d5848484610439565b6106e96106e0610a5a565b85858585610de7565b50505050565b60606106fa8261099b565b505f60065f8481526020019081526020015f20805461071890611de7565b80601f016020809104026020016040519081016040528092919081815260200182805461074490611de7565b801561078f5780601f106107665761010080835404028352916020019161078f565b820191905f5260205f20905b81548152906001019060200180831161077257829003601f168201915b505050505090505f61079f610f93565b90505f8151036107b35781925050506107f5565b5f825111156107e75780826040516020016107cf929190611e86565b604051602081830303815290604052925050506107f5565b6107f084610fa9565b925050505b919050565b5f60075f81548092919061080d90611ed6565b91905055905061081d838261100f565b6108278183611102565b505050565b5f60055f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f7f80ac58cd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061098457507f5b5e139f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b8061099457506109938261115c565b5b9050919050565b5f5f6109a6836111c5565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610a1857826040517f7e273289000000000000000000000000000000000000000000000000000000008152600401610a0f9190611a57565b60405180910390fd5b80915050919050565b5f60045f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b5f33905090565b610a6e83838360016111fe565b505050565b5f5f610a7e846111c5565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614610abf57610abe8184866113bd565b5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610b4a57610afe5f855f5f6111fe565b600160035f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825403925050819055505b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614610bc957600160035f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8460025f8681526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4809150509392505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610cee57816040517f5b08ba18000000000000000000000000000000000000000000000000000000008152600401610ce5919061194c565b60405180910390fd5b8060055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c3183604051610dda9190611806565b60405180910390a3505050565b5f8373ffffffffffffffffffffffffffffffffffffffff163b1115610f8c578273ffffffffffffffffffffffffffffffffffffffff1663150b7a02868685856040518563ffffffff1660e01b8152600401610e459493929190611f6f565b6020604051808303815f875af1925050508015610e8057506040513d601f19601f82011682018060405250810190610e7d9190611fcd565b60015b610f01573d805f8114610eae576040519150601f19603f3d011682016040523d82523d5f602084013e610eb3565b606091505b505f815103610ef957836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610ef0919061194c565b60405180910390fd5b805181602001fd5b63150b7a0260e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614610f8a57836040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401610f81919061194c565b60405180910390fd5b505b5050505050565b606060405180602001604052805f815250905090565b6060610fb48261099b565b505f610fbe610f93565b90505f815111610fdc5760405180602001604052805f815250611007565b80610fe684611480565b604051602001610ff7929190611e86565b6040516020818303038152906040525b915050919050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361107f575f6040517f64a0ae92000000000000000000000000000000000000000000000000000000008152600401611076919061194c565b60405180910390fd5b5f61108b83835f610a73565b90505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146110fd575f6040517f73c6ac6e0000000000000000000000000000000000000000000000000000000081526004016110f4919061194c565b60405180910390fd5b505050565b8060065f8481526020019081526020015f2090816111209190612198565b507ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7826040516111509190611a57565b60405180910390a15050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f60025f8381526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b808061123657505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b15611368575f6112458461099b565b90505f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141580156112af57508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614155b80156112c257506112c0818461082c565b155b1561130457826040517fa9fbf51f0000000000000000000000000000000000000000000000000000000081526004016112fb919061194c565b60405180910390fd5b811561136657838573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b8360045f8581526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b6113c883838361154a565b61147b575f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361143c57806040517f7e2732890000000000000000000000000000000000000000000000000000000081526004016114339190611a57565b60405180910390fd5b81816040517f177e802f000000000000000000000000000000000000000000000000000000008152600401611472929190612267565b60405180910390fd5b505050565b60605f600161148e8461160a565b0190505f8167ffffffffffffffff8111156114ac576114ab611ae0565b5b6040519080825280601f01601f1916602001820160405280156114de5781602001600182028036833780820191505090505b5090505f82602001820190505b60011561153f578080600190039150507f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a85816115345761153361228e565b5b0494505f85036114eb575b819350505050919050565b5f5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561160157508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614806115c257506115c1848461082c565b5b8061160057508273ffffffffffffffffffffffffffffffffffffffff166115e883610a21565b73ffffffffffffffffffffffffffffffffffffffff16145b5b90509392505050565b5f5f5f90507a184f03e93ff9f4daa797ed6e38ed64bf6a1f0100000000000000008310611666577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000838161165c5761165b61228e565b5b0492506040810190505b6d04ee2d6d415b85acef810000000083106116a3576d04ee2d6d415b85acef810000000083816116995761169861228e565b5b0492506020810190505b662386f26fc1000083106116d257662386f26fc1000083816116c8576116c761228e565b5b0492506010810190505b6305f5e10083106116fb576305f5e10083816116f1576116f061228e565b5b0492506008810190505b61271083106117205761271083816117165761171561228e565b5b0492506004810190505b6064831061174357606483816117395761173861228e565b5b0492506002810190505b600a8310611752576001810190505b80915050919050565b5f604051905090565b5f5ffd5b5f5ffd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6117a08161176c565b81146117aa575f5ffd5b50565b5f813590506117bb81611797565b92915050565b5f602082840312156117d6576117d5611764565b5b5f6117e3848285016117ad565b91505092915050565b5f8115159050919050565b611800816117ec565b82525050565b5f6020820190506118195f8301846117f7565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6118618261181f565b61186b8185611829565b935061187b818560208601611839565b61188481611847565b840191505092915050565b5f6020820190508181035f8301526118a78184611857565b905092915050565b5f819050919050565b6118c1816118af565b81146118cb575f5ffd5b50565b5f813590506118dc816118b8565b92915050565b5f602082840312156118f7576118f6611764565b5b5f611904848285016118ce565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6119368261190d565b9050919050565b6119468161192c565b82525050565b5f60208201905061195f5f83018461193d565b92915050565b61196e8161192c565b8114611978575f5ffd5b50565b5f8135905061198981611965565b92915050565b5f5f604083850312156119a5576119a4611764565b5b5f6119b28582860161197b565b92505060206119c3858286016118ce565b9150509250929050565b5f5f5f606084860312156119e4576119e3611764565b5b5f6119f18682870161197b565b9350506020611a028682870161197b565b9250506040611a13868287016118ce565b9150509250925092565b5f60208284031215611a3257611a31611764565b5b5f611a3f8482850161197b565b91505092915050565b611a51816118af565b82525050565b5f602082019050611a6a5f830184611a48565b92915050565b611a79816117ec565b8114611a83575f5ffd5b50565b5f81359050611a9481611a70565b92915050565b5f5f60408385031215611ab057611aaf611764565b5b5f611abd8582860161197b565b9250506020611ace85828601611a86565b9150509250929050565b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611b1682611847565b810181811067ffffffffffffffff82111715611b3557611b34611ae0565b5b80604052505050565b5f611b4761175b565b9050611b538282611b0d565b919050565b5f67ffffffffffffffff821115611b7257611b71611ae0565b5b611b7b82611847565b9050602081019050919050565b828183375f83830152505050565b5f611ba8611ba384611b58565b611b3e565b905082815260208101848484011115611bc457611bc3611adc565b5b611bcf848285611b88565b509392505050565b5f82601f830112611beb57611bea611ad8565b5b8135611bfb848260208601611b96565b91505092915050565b5f5f5f5f60808587031215611c1c57611c1b611764565b5b5f611c298782880161197b565b9450506020611c3a8782880161197b565b9350506040611c4b878288016118ce565b925050606085013567ffffffffffffffff811115611c6c57611c6b611768565b5b611c7887828801611bd7565b91505092959194509250565b5f67ffffffffffffffff821115611c9e57611c9d611ae0565b5b611ca782611847565b9050602081019050919050565b5f611cc6611cc184611c84565b611b3e565b905082815260208101848484011115611ce257611ce1611adc565b5b611ced848285611b88565b509392505050565b5f82601f830112611d0957611d08611ad8565b5b8135611d19848260208601611cb4565b91505092915050565b5f5f60408385031215611d3857611d37611764565b5b5f611d458582860161197b565b925050602083013567ffffffffffffffff811115611d6657611d65611768565b5b611d7285828601611cf5565b9150509250929050565b5f5f60408385031215611d9257611d91611764565b5b5f611d9f8582860161197b565b9250506020611db08582860161197b565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680611dfe57607f821691505b602082108103611e1157611e10611dba565b5b50919050565b5f606082019050611e2a5f83018661193d565b611e376020830185611a48565b611e44604083018461193d565b949350505050565b5f81905092915050565b5f611e608261181f565b611e6a8185611e4c565b9350611e7a818560208601611839565b80840191505092915050565b5f611e918285611e56565b9150611e9d8284611e56565b91508190509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611ee0826118af565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611f1257611f11611ea9565b5b600182019050919050565b5f81519050919050565b5f82825260208201905092915050565b5f611f4182611f1d565b611f4b8185611f27565b9350611f5b818560208601611839565b611f6481611847565b840191505092915050565b5f608082019050611f825f83018761193d565b611f8f602083018661193d565b611f9c6040830185611a48565b8181036060830152611fae8184611f37565b905095945050505050565b5f81519050611fc781611797565b92915050565b5f60208284031215611fe257611fe1611764565b5b5f611fef84828501611fb9565b91505092915050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026120547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82612019565b61205e8683612019565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61209961209461208f846118af565b612076565b6118af565b9050919050565b5f819050919050565b6120b28361207f565b6120c66120be826120a0565b848454612025565b825550505050565b5f5f905090565b6120dd6120ce565b6120e88184846120a9565b505050565b5b8181101561210b576121005f826120d5565b6001810190506120ee565b5050565b601f8211156121505761212181611ff8565b61212a8461200a565b81016020851015612139578190505b61214d6121458561200a565b8301826120ed565b50505b505050565b5f82821c905092915050565b5f6121705f1984600802612155565b1980831691505092915050565b5f6121888383612161565b9150826002028217905092915050565b6121a18261181f565b67ffffffffffffffff8111156121ba576121b9611ae0565b5b6121c48254611de7565b6121cf82828561210f565b5f60209050601f831160018114612200575f84156121ee578287015190505b6121f8858261217d565b86555061225f565b601f19841661220e86611ff8565b5f5b8281101561223557848901518255600182019150602085019450602081019050612210565b86831015612252578489015161224e601f891682612161565b8355505b6001600288020188555050505b505050505050565b5f60408201905061227a5f83018561193d565b6122876020830184611a48565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffdfea2646970667358221220b0541f627e3ea3d74b6f626f40b2e2b6e45c1ad3a55bf8da86299a2d24b897b364736f6c634300081e0033",
}

// Erc721ABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc721MetaData.ABI instead.
var Erc721ABI = Erc721MetaData.ABI

// Erc721Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Erc721MetaData.Bin instead.
var Erc721Bin = Erc721MetaData.Bin

// DeployErc721 deploys a new Ethereum contract, binding an instance of Erc721 to it.
func DeployErc721(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *Erc721, error) {
	parsed, err := Erc721MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Erc721Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// Erc721 is an auto generated Go binding around an Ethereum contract.
type Erc721 struct {
	Erc721Caller     // Read-only binding to the contract
	Erc721Transactor // Write-only binding to the contract
	Erc721Filterer   // Log filterer for contract events
}

// Erc721Caller is an auto generated read-only Go binding around an Ethereum contract.
type Erc721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc721Session struct {
	Contract     *Erc721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc721CallerSession struct {
	Contract *Erc721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Erc721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc721TransactorSession struct {
	Contract     *Erc721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721Raw is an auto generated low-level Go binding around an Ethereum contract.
type Erc721Raw struct {
	Contract *Erc721 // Generic contract binding to access the raw methods on
}

// Erc721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc721CallerRaw struct {
	Contract *Erc721Caller // Generic read-only contract binding to access the raw methods on
}

// Erc721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc721TransactorRaw struct {
	Contract *Erc721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewErc721 creates a new instance of Erc721, bound to a specific deployed contract.
func NewErc721(address common.Address, backend bind.ContractBackend) (*Erc721, error) {
	contract, err := bindErc721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// NewErc721Caller creates a new read-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Caller(address common.Address, caller bind.ContractCaller) (*Erc721Caller, error) {
	contract, err := bindErc721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Caller{contract: contract}, nil
}

// NewErc721Transactor creates a new write-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Transactor(address common.Address, transactor bind.ContractTransactor) (*Erc721Transactor, error) {
	contract, err := bindErc721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Transactor{contract: contract}, nil
}

// NewErc721Filterer creates a new log filterer instance of Erc721, bound to a specific deployed contract.
func NewErc721Filterer(address common.Address, filterer bind.ContractFilterer) (*Erc721Filterer, error) {
	contract, err := bindErc721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc721Filterer{contract: contract}, nil
}

// bindErc721 binds a generic wrapper to an already deployed contract.
func bindErc721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Erc721MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.Erc721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721.Contract.IsApprovedForAll(&_Erc721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721.Contract.IsApprovedForAll(&_Erc721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721Session) Name() (string, error) {
	return _Erc721.Contract.Name(&_Erc721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721CallerSession) Name() (string, error) {
	return _Erc721.Contract.Name(&_Erc721.CallOpts)
}

// NextTokenId is a free data retrieval call binding the contract method 0x75794a3c.
//
// Solidity: function nextTokenId() view returns(uint256)
func (_Erc721 *Erc721Caller) NextTokenId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "nextTokenId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextTokenId is a free data retrieval call binding the contract method 0x75794a3c.
//
// Solidity: function nextTokenId() view returns(uint256)
func (_Erc721 *Erc721Session) NextTokenId() (*big.Int, error) {
	return _Erc721.Contract.NextTokenId(&_Erc721.CallOpts)
}

// NextTokenId is a free data retrieval call binding the contract method 0x75794a3c.
//
// Solidity: function nextTokenId() view returns(uint256)
func (_Erc721 *Erc721CallerSession) NextTokenId() (*big.Int, error) {
	return _Erc721.Contract.NextTokenId(&_Erc721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721.Contract.SupportsInterface(&_Erc721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721.Contract.SupportsInterface(&_Erc721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721Session) Symbol() (string, error) {
	return _Erc721.Contract.Symbol(&_Erc721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721CallerSession) Symbol() (string, error) {
	return _Erc721.Contract.Symbol(&_Erc721.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721Session) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721.Contract.TokenURI(&_Erc721.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721.Contract.TokenURI(&_Erc721.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) returns()
func (_Erc721 *Erc721Transactor) Mint(opts *bind.TransactOpts, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "mint", to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) returns()
func (_Erc721 *Erc721Session) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Erc721.Contract.Mint(&_Erc721.TransactOpts, to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) returns()
func (_Erc721 *Erc721TransactorSession) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Erc721.Contract.Mint(&_Erc721.TransactOpts, to, tokenURI)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom0(&_Erc721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom0(&_Erc721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.Contract.SetApprovalForAll(&_Erc721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.Contract.SetApprovalForAll(&_Erc721.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// Erc721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc721 contract.
type Erc721ApprovalIterator struct {
	Event *Erc721Approval // Event containing the contract specifics and raw log

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
func (it *Erc721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Approval)
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
		it.Event = new(Erc721Approval)
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
func (it *Erc721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Approval represents a Approval event raised by the Erc721 contract.
type Erc721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*Erc721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721ApprovalIterator{contract: _Erc721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Approval)
				if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) ParseApproval(log types.Log) (*Erc721Approval, error) {
	event := new(Erc721Approval)
	if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Erc721 contract.
type Erc721ApprovalForAllIterator struct {
	Event *Erc721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *Erc721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721ApprovalForAll)
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
		it.Event = new(Erc721ApprovalForAll)
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
func (it *Erc721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721ApprovalForAll represents a ApprovalForAll event raised by the Erc721 contract.
type Erc721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*Erc721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Erc721ApprovalForAllIterator{contract: _Erc721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Erc721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721ApprovalForAll)
				if err := _Erc721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) ParseApprovalForAll(log types.Log) (*Erc721ApprovalForAll, error) {
	event := new(Erc721ApprovalForAll)
	if err := _Erc721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721BatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the Erc721 contract.
type Erc721BatchMetadataUpdateIterator struct {
	Event *Erc721BatchMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *Erc721BatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721BatchMetadataUpdate)
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
		it.Event = new(Erc721BatchMetadataUpdate)
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
func (it *Erc721BatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721BatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721BatchMetadataUpdate represents a BatchMetadataUpdate event raised by the Erc721 contract.
type Erc721BatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_Erc721 *Erc721Filterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*Erc721BatchMetadataUpdateIterator, error) {

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &Erc721BatchMetadataUpdateIterator{contract: _Erc721.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_Erc721 *Erc721Filterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *Erc721BatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721BatchMetadataUpdate)
				if err := _Erc721.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
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

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_Erc721 *Erc721Filterer) ParseBatchMetadataUpdate(log types.Log) (*Erc721BatchMetadataUpdate, error) {
	event := new(Erc721BatchMetadataUpdate)
	if err := _Erc721.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721MetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the Erc721 contract.
type Erc721MetadataUpdateIterator struct {
	Event *Erc721MetadataUpdate // Event containing the contract specifics and raw log

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
func (it *Erc721MetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721MetadataUpdate)
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
		it.Event = new(Erc721MetadataUpdate)
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
func (it *Erc721MetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721MetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721MetadataUpdate represents a MetadataUpdate event raised by the Erc721 contract.
type Erc721MetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_Erc721 *Erc721Filterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*Erc721MetadataUpdateIterator, error) {

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &Erc721MetadataUpdateIterator{contract: _Erc721.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_Erc721 *Erc721Filterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *Erc721MetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721MetadataUpdate)
				if err := _Erc721.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_Erc721 *Erc721Filterer) ParseMetadataUpdate(log types.Log) (*Erc721MetadataUpdate, error) {
	event := new(Erc721MetadataUpdate)
	if err := _Erc721.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc721 contract.
type Erc721TransferIterator struct {
	Event *Erc721Transfer // Event containing the contract specifics and raw log

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
func (it *Erc721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Transfer)
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
		it.Event = new(Erc721Transfer)
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
func (it *Erc721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Transfer represents a Transfer event raised by the Erc721 contract.
type Erc721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*Erc721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721TransferIterator{contract: _Erc721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Transfer)
				if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) ParseTransfer(log types.Log) (*Erc721Transfer, error) {
	event := new(Erc721Transfer)
	if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
