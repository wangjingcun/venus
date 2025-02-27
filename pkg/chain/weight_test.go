package chain

import (
	"bytes"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	tf "github.com/filecoin-project/venus/pkg/testhelpers/testflags"
	"github.com/filecoin-project/venus/venus-shared/types"
)

func TestWeightTieBreaker(t *testing.T) {
	tf.UnitTest(t)
	smallerVrf1 := []byte{0} // blake2b starts with byte 3
	largerVrf1 := []byte{2}  // blake2b starts with byte 187

	smallerVrfSum1 := blake2b.Sum256(smallerVrf1)
	largerVrfSum1 := blake2b.Sum256(largerVrf1)

	// this checks that smallerVrf is in fact smaller than LargerVrf
	require.True(t, bytes.Compare(smallerVrfSum1[:], largerVrfSum1[:]) < 0)

	smallerVrf2 := []byte{3} // blake2b starts with byte 232
	largerVrf2 := []byte{1}  // blake2b starts with byte 238

	smallerVrfSum2 := blake2b.Sum256(smallerVrf2)
	largerVrfSum2 := blake2b.Sum256(largerVrf2)

	// this checks that smallerVrf is in fact smaller than LargerVrf
	require.True(t, bytes.Compare(smallerVrfSum2[:], largerVrfSum2[:]) < 0)

	ts1 := mockTipSet(t, largerVrf1, smallerVrf2)
	ts2 := mockTipSet(t, smallerVrf1, largerVrf2)

	// ts1's first block has a larger VRF than ts2's, so it should lose
	require.False(t, breakWeightTie(ts1, ts2))
}

func mockTipSet(t *testing.T, vrf1, vrf2 []byte) *types.TipSet {
	minerAct, err := address.NewIDAddress(0)
	require.NoError(t, err)
	c, err := cid.Decode("QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH")
	require.NoError(t, err)
	blks := []*types.BlockHeader{
		{
			Miner:                 minerAct,
			Height:                abi.ChainEpoch(1),
			ParentStateRoot:       c,
			ParentMessageReceipts: c,
			Messages:              c,
			Ticket:                &types.Ticket{VRFProof: vrf1},
		},
		{
			Miner:                 minerAct,
			Height:                abi.ChainEpoch(1),
			ParentStateRoot:       c,
			ParentMessageReceipts: c,
			Messages:              c,
			Ticket:                &types.Ticket{VRFProof: vrf2},
		},
	}
	ts, err := types.NewTipSet(blks)
	require.NoError(t, err)
	return ts
}
