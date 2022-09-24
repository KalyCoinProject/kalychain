package server

import (
	"github.com/KalyCoinProject/kalychain/consensus"
	consensusDev "github.com/KalyCoinProject/kalychain/consensus/dev"
	consensusDummy "github.com/KalyCoinProject/kalychain/consensus/dummy"
	consensusIBFT "github.com/KalyCoinProject/kalychain/consensus/ibft"
	"github.com/KalyCoinProject/kalychain/secrets"
	"github.com/KalyCoinProject/kalychain/secrets/awsssm"
	"github.com/KalyCoinProject/kalychain/secrets/hashicorpvault"
	"github.com/KalyCoinProject/kalychain/secrets/local"
)

type ConsensusType string

const (
	DevConsensus   ConsensusType = "dev"
	IBFTConsensus  ConsensusType = "ibft"
	DummyConsensus ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:   consensusDev.Factory,
	IBFTConsensus:  consensusIBFT.Factory,
	DummyConsensus: consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
