package snapshot

import (
	"context"

	"github.com/KalyCoinProject/kalychain/command"
	"github.com/KalyCoinProject/kalychain/command/helper"
	ibftOp "github.com/KalyCoinProject/kalychain/consensus/ibft/proto"
)

const (
	numberFlag = "number"
)

var (
	params = &snapshotParams{}
)

type snapshotParams struct {
	blockNumber int

	snapshot *ibftOp.Snapshot
}

func (p *snapshotParams) initSnapshot(grpcAddress string) error {
	ibftClient, err := helper.GetIBFTOperatorClientConnection(grpcAddress)
	if err != nil {
		return err
	}

	snapshot, err := ibftClient.GetSnapshot(
		context.Background(),
		p.getSnapshotRequest(),
	)
	if err != nil {
		return err
	}

	p.snapshot = snapshot

	return nil
}

func (p *snapshotParams) getSnapshotRequest() *ibftOp.SnapshotReq {
	req := &ibftOp.SnapshotReq{
		Latest: true,
	}

	if p.blockNumber >= 0 {
		req.Latest = false
		req.Number = uint64(p.blockNumber)
	}

	return req
}

func (p *snapshotParams) getResult() command.CommandResult {
	return newIBFTSnapshotResult(p.snapshot)
}
