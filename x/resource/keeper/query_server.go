package keeper

import (
	didkeeper "github.com/pointidentity/pointidentity-node/x/did/keeper"
	"github.com/pointidentity/pointidentity-node/x/resource/types"
)

type queryServer struct {
	Keeper
	didKeeper didkeeper.Keeper
}

// NewQueryServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewQueryServer(keeper Keeper, pidKeeper didkeeper.Keeper) types.QueryServer {
	return &queryServer{
		Keeper:    keeper,
		didKeeper: pidKeeper,
	}
}

var _ types.QueryServer = queryServer{}
