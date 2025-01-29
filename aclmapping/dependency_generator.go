package aclmapping

import (
	aclkeeper "github.com/cosmos/cosmos-sdk/x/accesscontrol/keeper"
	aclbankmapping "github.com/eni-chain/go-eni/aclmapping/bank"
	aclevmmapping "github.com/eni-chain/go-eni/aclmapping/evm"
	acloraclemapping "github.com/eni-chain/go-eni/aclmapping/oracle"
	acltokenfactorymapping "github.com/eni-chain/go-eni/aclmapping/tokenfactory"
	aclwasmmapping "github.com/eni-chain/go-eni/aclmapping/wasm"
	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
)

type CustomDependencyGenerator struct{}

func NewCustomDependencyGenerator() CustomDependencyGenerator {
	return CustomDependencyGenerator{}
}

func (customDepGen CustomDependencyGenerator) GetCustomDependencyGenerators(evmKeeper evmkeeper.Keeper) aclkeeper.DependencyGeneratorMap {
	dependencyGeneratorMap := make(aclkeeper.DependencyGeneratorMap)
	wasmDependencyGenerators := aclwasmmapping.NewWasmDependencyGenerator()

	dependencyGeneratorMap = dependencyGeneratorMap.Merge(aclbankmapping.GetBankDepedencyGenerator())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(acltokenfactorymapping.GetTokenFactoryDependencyGenerators())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(wasmDependencyGenerators.GetWasmDependencyGenerators())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(acloraclemapping.GetOracleDependencyGenerator())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(aclevmmapping.GetEVMDependencyGenerators(evmKeeper))

	return dependencyGeneratorMap
}
