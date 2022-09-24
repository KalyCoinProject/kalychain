package command

import "github.com/KalyCoinProject/kalychain/server"

const (
	DefaultGenesisFileName = "genesis.json"
	DefaultChainName       = "kalychain"
	DefaultChainID         = 100
	DefaultPremineBalance  = "0x3635C9ADC5DEA00000" // 1000 ETH
	DefaultConsensus       = server.IBFTConsensus
	DefaultGenesisGasUsed  = 458752  // 0x70000
	DefaultGenesisGasLimit = 5242880 // 0x500000
)

const (
	JSONOutputFlag     = "json"
	GRPCAddressFlag    = "grpc-address"
	JSONRPCFlag        = "jsonrpc"
	GraphQLAddressFlag = "graphql-address"
	PprofFlag          = "pprof"
	PprofAddressFlag   = "pprof-address"
)

// Legacy flag that needs to be present to preserve backwards
// compatibility with running clients
const (
	GRPCAddressFlagLEGACY = "grpc"
)
