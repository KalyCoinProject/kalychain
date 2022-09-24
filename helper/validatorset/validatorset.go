package validatorset

import (
	"errors"
	"math/big"

	"github.com/KalyCoinProject/kalychain/chain"
	"github.com/KalyCoinProject/kalychain/helper/common"
	"github.com/KalyCoinProject/kalychain/helper/hex"
	"github.com/KalyCoinProject/kalychain/helper/keccak"
	"github.com/KalyCoinProject/kalychain/types"
)

// getAddressMapping returns the key for the SC storage mapping (address => something)
//
// More information:
// https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html
func getAddressMapping(address types.Address, slot int64) []byte {
	bigSlot := big.NewInt(slot)

	finalSlice := append(
		common.PadLeftOrTrim(address.Bytes(), 32),
		common.PadLeftOrTrim(bigSlot.Bytes(), 32)...,
	)
	keccakValue := keccak.Keccak256(nil, finalSlice)

	return keccakValue
}

// getIndexWithOffset is a helper method for adding an offset to the already found keccak hash
func getIndexWithOffset(keccakHash []byte, offset int64) []byte {
	bigOffset := big.NewInt(offset)
	bigKeccak := big.NewInt(0).SetBytes(keccakHash)

	bigKeccak.Add(bigKeccak, bigOffset)

	return bigKeccak.Bytes()
}

// getStorageIndexes is a helper function for getting the correct indexes
// of the storage slots which need to be modified during bootstrap.
//
// It is SC dependant, and based on the SC located at:
// https://github.com/KalyCoinProject/kalychain-contracts
func getStorageIndexes(address types.Address, index int64) *StorageIndexes {
	storageIndexes := StorageIndexes{}

	// Get the indexes for the mappings
	// The index for the mapping is retrieved with:
	// keccak(address . slot)
	// . stands for concatenation (basically appending the bytes)
	storageIndexes.AddressToIsValidatorIndex = getAddressMapping(address, addressToIsValidatorSlot)
	storageIndexes.AddressToStakedAmountIndex = getAddressMapping(address, addressToStakedAmountSlot)
	storageIndexes.AddressToValidatorIndexIndex = getAddressMapping(address, addressToValidatorIndexSlot)

	// Get the indexes for _status, _owner, _validators, _stakedAmount, etc
	// Index for regular types is calculated as just the regular slot
	storageIndexes.StatusIndex = big.NewInt(statusSlot).Bytes()
	storageIndexes.OwnerIndex = big.NewInt(ownerSlot).Bytes()
	storageIndexes.ThresholdIndex = big.NewInt(thresholdSlot).Bytes()
	storageIndexes.MinimumIndex = big.NewInt(minimumSlot).Bytes()
	storageIndexes.StakedAmountIndex = big.NewInt(stakedAmountSlot).Bytes()

	// Index for array types is calculated as keccak(slot) + index
	// The slot for the dynamic arrays that's put in the keccak needs to be in hex form (padded 64 chars)
	storageIndexes.ValidatorsIndex = getIndexWithOffset(
		keccak.Keccak256(nil, common.PadLeftOrTrim(big.NewInt(validatorsSlot).Bytes(), 32)),
		index,
	)

	// For any dynamic array in Solidity, the size of the actual array should be
	// located on slot x
	storageIndexes.ValidatorsArraySizeIndex = []byte{byte(validatorsSlot)}

	return &storageIndexes
}

// PredeployParams contains the values used to predeploy the PoS staking contract
type PredeployParams struct {
	Owner      types.Address
	Validators []types.Address
}

// StorageIndexes is a wrapper for different storage indexes that
// need to be modified
type StorageIndexes struct {
	StatusIndex                  []byte // uint256
	OwnerIndex                   []byte // address
	ThresholdIndex               []byte // uint256
	MinimumIndex                 []byte // uint256
	ValidatorsIndex              []byte // []address
	ValidatorsArraySizeIndex     []byte // []address size
	AddressToIsValidatorIndex    []byte // mapping(address => bool)
	AddressToStakedAmountIndex   []byte // mapping(address => uint256)
	AddressToValidatorIndexIndex []byte // mapping(address => uint256)
	StakedAmountIndex            []byte // uint256
}

// Slot definitions for SC storage
const (
	statusSlot = int64(iota) // Slot 0
	ownerSlot
	thresholdSlot
	minimumSlot
	validatorsSlot
	addressToIsValidatorSlot
	addressToStakedAmountSlot
	addressToValidatorIndexSlot
	stakedAmountSlot
)

const (
	DefaultStakedBalance    = "0x84595161401484A000000" // 10_000_000 DC
	DefaultStatusNotEntered = 1                         // ReentrancyGuard status contant
	//nolint:lll
	StakingSCBytecode = "0x608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063ca1e781911610071578063ca1e781914610281578063d0a5e6ce1461029f578063f2fde38b146102bb578063f90ecacc146102d7578063facd743b1461030757610116565b80638da5cb5b1461020f578063960bfe041461022d5780639cbfc76514610249578063a694fc3a1461026557610116565b8063373d6132116100e9578063373d61321461018f57806342cde4e8146101ad5780634d238c8e146101cb57806352d6804d146101e7578063715018a61461020557610116565b80630c340a241461011b5780632367f6b5146101395780632def6620146101695780633209e9e614610173575b600080fd5b610123610337565b6040516101309190611b59565b60405180910390f35b610153600480360381019061014e919061189f565b610361565b6040516101609190611d71565b60405180910390f35b6101716103aa565b005b61018d600480360381019061018891906118f9565b6106f1565b005b6101976107d8565b6040516101a49190611d71565b60405180910390f35b6101b56107e2565b6040516101c29190611d71565b60405180910390f35b6101e560048036038101906101e0919061189f565b6107ec565b005b6101ef610aee565b6040516101fc9190611d71565b60405180910390f35b61020d610af8565b005b610217610b94565b6040516102249190611b59565b60405180910390f35b610247600480360381019061024291906118f9565b610bbe565b005b610263600480360381019061025e919061189f565b610ca5565b005b61027f600480360381019061027a91906118f9565b610e70565b005b610289611120565b6040516102969190611bd4565b60405180910390f35b6102b960048036038101906102b4919061189f565b6111ae565b005b6102d560048036038101906102d0919061189f565b61131b565b005b6102f160048036038101906102ec91906118f9565b611427565b6040516102fe9190611b59565b60405180910390f35b610321600480360381019061031c919061189f565b611466565b60405161032e9190611bf6565b60405180910390f35b6000600960009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b3373ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff1614610418576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040f90611cd1565b60405180910390fd5b6002600054141561045e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045590611d51565b60405180910390fd5b60026000819055506000600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600081116104ed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104e490611c51565b60405180910390fd5b6000600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610547816008546114bc90919063ffffffff16565b600881905550600960009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33836040518363ffffffff1660e01b81526004016105aa929190611bab565b602060405180830381600087803b1580156105c457600080fd5b505af11580156105d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105fc91906118cc565b50600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156106a25760035460048054905011610698576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068f90611cf1565b60405180910390fd5b6106a1336114d2565b5b803373ffffffffffffffffffffffffffffffffffffffff167f0f5bb82176feb1b5e747e28471aa92156a04d9f3ab9f45f28e2d704232b93f7560405160405180910390a3506001600081905550565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610781576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161077890611c91565b60405180910390fd5b600060035490508160038190555081813373ffffffffffffffffffffffffffffffffffffffff167f6eb5ec46450e0c6e94bb67a32e6bca9ec9ff819009505cbc6b886caf512d37bc60405160405180910390a45050565b6000600854905090565b6000600254905090565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461087c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087390611c91565b60405180910390fd5b600254600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610900576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108f790611d31565b60405180910390fd5b600560008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161561098d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161098490611d11565b60405180910390fd5b600480549050600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506001600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506004819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8064a302796c89446a96d63470b5b036212da26bd2debe5bec73e0170a9a5e8360405160405180910390a350565b6000600354905090565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610b88576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b7f90611c91565b60405180910390fd5b610b926000611784565b565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610c4e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c4590611c91565b60405180910390fd5b600060025490508160028190555081813373ffffffffffffffffffffffffffffffffffffffff167fed4e7b6d1951b75b13e101295f8473d6492319d89608bbfbfdbc643d96246f7d60405160405180910390a45050565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610d35576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d2c90611c91565b60405180910390fd5b60035460048054905011610d7e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d7590611cf1565b60405180910390fd5b600560008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610e0a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e0190611c71565b60405180910390fd5b610e13816114d2565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f4edebfc5ffaa4271f94ab363e643701124f2b4381b7a4f614dbdf75f166dc0cb60405160405180910390a350565b3373ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff1614610ede576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ed590611cd1565b60405180910390fd5b60026000541415610f24576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f1b90611d51565b60405180910390fd5b600260008190555060008111610f6f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f6690611c31565b60405180910390fd5b600960009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b8152600401610fce93929190611b74565b602060405180830381600087803b158015610fe857600080fd5b505af1158015610ffc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061102091906118cc565b5061107381600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461184a90919063ffffffff16565b600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506110cb8160085461184a90919063ffffffff16565b600881905550803373ffffffffffffffffffffffffffffffffffffffff167f9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d60405160405180910390a3600160008190555050565b606060048054806020026020016040519081016040528092919081815260200182805480156111a457602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001906001019080831161115a575b5050505050905090565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461123e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161123590611c91565b60405180910390fd5b6000600960009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600960006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f48a312d70029e6dd97980e9e051e1ff0b8b8be967450af46ce6dc5fa9830428f60405160405180910390a45050565b3373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146113ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113a290611c91565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561141b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161141290611c11565b60405180910390fd5b61142481611784565b50565b6004818154811061143757600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b600081836114ca9190611e2c565b905092915050565b600480549050600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410611558576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161154f90611cb1565b60405180910390fd5b6000600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060006115b760016004805490506114bc90919063ffffffff16565b90508082146116a6576000600482815481106115d6576115d5611f06565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050806004848154811061161857611617611f06565b5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555082600760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff0219169055600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600480548061174a57611749611ed7565b5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690559055505050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600081836118589190611dd6565b905092915050565b60008135905061186f8161216f565b92915050565b60008151905061188481612186565b92915050565b6000813590506118998161219d565b92915050565b6000602082840312156118b5576118b4611f35565b5b60006118c384828501611860565b91505092915050565b6000602082840312156118e2576118e1611f35565b5b60006118f084828501611875565b91505092915050565b60006020828403121561190f5761190e611f35565b5b600061191d8482850161188a565b91505092915050565b6000611932838361193e565b60208301905092915050565b61194781611e60565b82525050565b61195681611e60565b82525050565b600061196782611d9c565b6119718185611db4565b935061197c83611d8c565b8060005b838110156119ad5781516119948882611926565b975061199f83611da7565b925050600181019050611980565b5085935050505092915050565b6119c381611e72565b82525050565b60006119d6602683611dc5565b91506119e182611f3a565b604082019050919050565b60006119f9600e83611dc5565b9150611a0482611f89565b602082019050919050565b6000611a1c601d83611dc5565b9150611a2782611fb2565b602082019050919050565b6000611a3f601983611dc5565b9150611a4a82611fdb565b602082019050919050565b6000611a62601c83611dc5565b9150611a6d82612004565b602082019050919050565b6000611a85601283611dc5565b9150611a908261202d565b602082019050919050565b6000611aa8601a83611dc5565b9150611ab382612056565b602082019050919050565b6000611acb602583611dc5565b9150611ad68261207f565b604082019050919050565b6000611aee602583611dc5565b9150611af9826120ce565b604082019050919050565b6000611b11601d83611dc5565b9150611b1c8261211d565b602082019050919050565b6000611b34601f83611dc5565b9150611b3f82612146565b602082019050919050565b611b5381611e9e565b82525050565b6000602082019050611b6e600083018461194d565b92915050565b6000606082019050611b89600083018661194d565b611b96602083018561194d565b611ba36040830184611b4a565b949350505050565b6000604082019050611bc0600083018561194d565b611bcd6020830184611b4a565b9392505050565b60006020820190508181036000830152611bee818461195c565b905092915050565b6000602082019050611c0b60008301846119ba565b92915050565b60006020820190508181036000830152611c2a816119c9565b9050919050565b60006020820190508181036000830152611c4a816119ec565b9050919050565b60006020820190508181036000830152611c6a81611a0f565b9050919050565b60006020820190508181036000830152611c8a81611a32565b9050919050565b60006020820190508181036000830152611caa81611a55565b9050919050565b60006020820190508181036000830152611cca81611a78565b9050919050565b60006020820190508181036000830152611cea81611a9b565b9050919050565b60006020820190508181036000830152611d0a81611abe565b9050919050565b60006020820190508181036000830152611d2a81611ae1565b9050919050565b60006020820190508181036000830152611d4a81611b04565b9050919050565b60006020820190508181036000830152611d6a81611b27565b9050919050565b6000602082019050611d866000830184611b4a565b92915050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b6000611de182611e9e565b9150611dec83611e9e565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611e2157611e20611ea8565b5b828201905092915050565b6000611e3782611e9e565b9150611e4283611e9e565b925082821015611e5557611e54611ea8565b5b828203905092915050565b6000611e6b82611e7e565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600080fd5b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b7f496e76616c696420616d6f756e74000000000000000000000000000000000000600082015250565b7f4f6e6c79207374616b65722063616e2063616c6c2066756e6374696f6e000000600082015250565b7f4163636f756e74206d7573742062652076616c696461746f7200000000000000600082015250565b7f4f6e6c79206f776e65722063616e2063616c6c2066756e6374696f6e00000000600082015250565b7f696e646578206f7574206f662072616e67650000000000000000000000000000600082015250565b7f4f6e6c7920454f412063616e2063616c6c2066756e6374696f6e000000000000600082015250565b7f56616c696461746f72732063616e2774206265206c657373207468616e206d6960008201527f6e696d756d000000000000000000000000000000000000000000000000000000602082015250565b7f4163636f756e742063616e6e6f7420616c726561647920626520612076616c6960008201527f6461746f72000000000000000000000000000000000000000000000000000000602082015250565b7f4163636f756e74206d757374206265207374616b656420656e6f756768000000600082015250565b7f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00600082015250565b61217881611e60565b811461218357600080fd5b50565b61218f81611e72565b811461219a57600080fd5b50565b6121a681611e9e565b81146121b157600080fd5b5056fea26469706673582212207b4f1a27bec5e8a044f17b8b4069e66cb9e4d2ad6ef3d1654e32ea80cff42c5d64736f6c63430008060033"
)

// PredeploySC is a helper method for setting up the ValidatorSet smart contract account,
// using the passed in validators as pre-staked validators
func PredeploySC(params PredeployParams) (*chain.GenesisAccount, error) {
	// Set the code for the staking smart contract
	// Code retrieved from https://github.com/KalyCoinProject/kalychain-contracts
	scHex, _ := hex.DecodeHex(StakingSCBytecode)
	stakingAccount := &chain.GenesisAccount{
		Code: scHex,
	}

	if params.Owner == types.ZeroAddress {
		return nil, errors.New("contract owner should not be empty")
	}

	// Generate the empty account storage map
	storageMap := make(map[types.Hash]types.Hash)
	bigOne := big.NewInt(1)
	bigTrueValue := big.NewInt(1)
	stakedAmount := big.NewInt(0)
	notEnteredStatus := big.NewInt(DefaultStatusNotEntered)

	for indx, validator := range params.Validators {
		// Get the storage indexes
		storageIndexes := getStorageIndexes(validator, int64(indx))

		// Set the value for the owner
		storageMap[types.BytesToHash(storageIndexes.OwnerIndex)] =
			types.BytesToHash(params.Owner.Bytes())

		// Set the value for the owner
		storageMap[types.BytesToHash(storageIndexes.MinimumIndex)] =
			types.BytesToHash(bigOne.Bytes())

		// Set the value for the validators array
		storageMap[types.BytesToHash(storageIndexes.ValidatorsIndex)] =
			types.BytesToHash(
				validator.Bytes(),
			)

		// Set the value for the address -> validator array index mapping
		storageMap[types.BytesToHash(storageIndexes.AddressToIsValidatorIndex)] =
			types.BytesToHash(bigTrueValue.Bytes())

		// Set the value for the address -> validator index mapping
		storageMap[types.BytesToHash(storageIndexes.AddressToValidatorIndexIndex)] =
			types.StringToHash(hex.EncodeUint64(uint64(indx)))

		// Set the value for the total staked amount
		storageMap[types.BytesToHash(storageIndexes.StakedAmountIndex)] =
			types.BytesToHash(stakedAmount.Bytes())

		// Set the value for the size of the validators array
		storageMap[types.BytesToHash(storageIndexes.ValidatorsArraySizeIndex)] =
			types.StringToHash(hex.EncodeUint64(uint64(indx + 1)))

		// Set the default status
		storageMap[types.BytesToHash(storageIndexes.StatusIndex)] =
			types.BytesToHash(notEnteredStatus.Bytes())
	}

	// Save the storage map
	stakingAccount.Storage = storageMap

	// Set the Staking SC balance to numValidators * defaultStakedBalance
	stakingAccount.Balance = stakedAmount

	return stakingAccount, nil
}
