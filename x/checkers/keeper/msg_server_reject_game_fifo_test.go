package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRejectMiddleGameHasSavedFifo(t *testing.T) {
    msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
    ctx := sdk.UnwrapSDKContext(context)
    msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: bob,
        Black:   carol,
        Red:     alice,
    })
    msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: carol,
        Black:   alice,
        Red:     bob,
    })
    msgServer.RejectGame(context, &types.MsgRejectGame{
        Creator:   carol,
        GameIndex: "1",
    })
    systemInfo, found := keeper.GetSystemInfo(ctx)
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId:        4,
        FifoHeadIndex: "2",
        FifoTailIndex: "3",
    }, systemInfo)
    game1, found := keeper.GetStoredGame(ctx, "1")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "0",
        Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "b",
        Black:       bob,
        Red:         carol,
        MoveCount:   uint64(0),
        BeforeIndex: "-1",
        AfterIndex:  "2",
    }, game1)
    game3, found := keeper.GetStoredGame(ctx, "3")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "2",
        Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn:        "b",
        Black:       alice,
        Red:         bob,
        MoveCount:   uint64(0),
        BeforeIndex: "0",
        AfterIndex:  "-1",
    }, game3)
}
