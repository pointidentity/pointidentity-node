package keeper

import (
	didkeeper "github.com/pointidentity/pointidentity-node/x/did/keeper"
	"github.com/pointidentity/pointidentity-node/x/resource/types"
)

type msgServer struct {
	Keeper
	didKeeper didkeeper.Keeper
}

// NewMsgServerImpl returns an implementation of the x/auth MsgServer interface.
func NewMsgServerImpl(keeper Keeper, didKeeper didkeeper.Keeper) types.MsgServer {
	return &msgServer{
		Keeper:    keeper,
		didKeeper: didKeeper,
	}
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper, pidKeeper didkeeper.Keeper) types.MsgServer {
	return &msgServer{
		Keeper:    keeper,
		didKeeper: pidKeeper,
	}
}

var _ types.MsgServer = msgServer{}
