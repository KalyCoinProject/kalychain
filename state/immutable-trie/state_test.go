package itrie

import (
	"testing"

	"github.com/KalyCoinProject/kalychain/state"
)

func TestState(t *testing.T) {
	state.TestState(t, buildPreState)
}

func buildPreState(pre state.PreStates) (state.State, state.Snapshot) {
	storage := NewMemoryStorage()
	st := NewState(storage)
	snap := st.NewSnapshot()

	return st, snap
}
