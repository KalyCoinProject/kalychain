package protocol

import (
	"math/big"

	"github.com/KalyCoinProject/kalychain/blockchain"
	"github.com/KalyCoinProject/kalychain/types"
)

// Blockchain is the interface required by the syncer to connect to the blockchain
type blockchainShim interface {
	SubscribeEvents() blockchain.Subscription
	Header() *types.Header
	CurrentTD() *big.Int

	GetTD(hash types.Hash) (*big.Int, bool)
	GetReceiptsByHash(types.Hash) ([]*types.Receipt, error)
	GetBodyByHash(types.Hash) (*types.Body, bool)
	GetHeaderByHash(types.Hash) (*types.Header, bool)
	GetHeaderByNumber(n uint64) (*types.Header, bool)

	// advance chain methods
	WriteBlock(block *types.Block) error
	VerifyFinalizedBlock(block *types.Block) error
	CalculateGasLimit(number uint64) (uint64, error)
}
