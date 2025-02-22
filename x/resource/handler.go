package resource

import (
	"fmt"

	didkeeper "github.com/pointidentity/pointidentity-node/x/did/keeper"

	errorsmod "cosmossdk.io/errors"
	"github.com/pointidentity/pointidentity-node/x/resource/keeper"
	"github.com/pointidentity/pointidentity-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper, pidKeeper didkeeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServer(k, pidKeeper)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateResource:
			res, err := msgServer.CreateResource(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
