package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServerWithOneGameForRejectGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	return server, *k, context
}

func TestRejectGameWrongByCreator(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   alice,
		GameIndex: "1",
	})
	require.Nil(t, rejectGameResponse)
	require.Equal(t, alice+": message creator is not a player", err.Error())
}

func TestRejectGameByBlackNoMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgRejectGameResponse{}, *rejectGameResponse)
}

func TestRejectGameByBlackNoMoveRemovedGame(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
}

func TestRejectGameByBlackNoMoveEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-rejected",
		Attributes: []sdk.Attribute{
			{Key: "creator", Value: bob},
			{Key: "game-index", Value: "1"},
		},
	}, event)
}

func TestRejectGameByRedNoMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgRejectGameResponse{}, *rejectGameResponse)
}

func TestRejectGameByRedNoMoveRemovedGame(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
}

func TestRejectGameByRedNoMoveEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-rejected",
		Attributes: []sdk.Attribute{
			{Key: "creator", Value: carol},
			{Key: "game-index", Value: "1"},
		},
	}, event)
}
