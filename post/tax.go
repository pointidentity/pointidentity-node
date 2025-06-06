package posthandler

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	pointidante "github.com/pointidentity/pointidentity-node/ante"
	didtypes "github.com/pointidentity/pointidentity-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

// TaxDecorator will handle tax for all taxable messages
type TaxDecorator struct {
	accountKeeper   ante.AccountKeeper
	bankKeeper      BankKeeper
	feegrantKeeper  ante.FeegrantKeeper
	didKeeper       pointidante.DidKeeper
	resourceKeeper  pointidante.ResourceKeeper
	feemarketKeeper FeeMarketKeeper
}

// NewTaxDecorator returns a new taxDecorator
func NewTaxDecorator(ak ante.AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, dk pointidante.DidKeeper, rk pointidante.ResourceKeeper, fmk FeeMarketKeeper) TaxDecorator {
	return TaxDecorator{
		accountKeeper:   ak,
		bankKeeper:      bk,
		feegrantKeeper:  fk,
		didKeeper:       dk,
		resourceKeeper:  rk,
		feemarketKeeper: fmk,
	}
}

// AnteHandle handles tax for all taxable messages
func (td TaxDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (sdk.Context, error) {
	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, err
	}

	nativeDenom := params.FeeDenom
	// must implement FeeTx
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", tx)
	}
	// if simulate, perform no-op
	if simulate {
		return next(ctx, tx, simulate, success)
	}
	// get metrics for tax
	rewards, burn, taxable, err := td.isTaxable(ctx, feeTx)
	if err != nil {
		return ctx, err
	}
	if taxable {
		err := td.handleTaxableTransaction(ctx, feeTx, simulate, rewards, burn, tx)
		if err != nil {
			return ctx, err
		}
		return next(ctx, tx, simulate, success)
	}
	params, err = td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market params")
	}
	// return if disabled
	if !params.Enabled {
		return next(ctx, tx, simulate, success)
	}

	enabledHeight, err := td.feemarketKeeper.GetEnabledHeight(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market enabled height")
	}

	// if the current height is that which enabled the feemarket or lower, skip deduction
	if ctx.BlockHeight() <= enabledHeight {
		return next(ctx, tx, simulate, success)
	}

	// update fee market state
	state, err := td.feemarketKeeper.GetState(ctx)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get fee market state")
	}

	onlyNativeDenom := true
	for _, fee := range feeTx.GetFee() {
		if fee.Denom != nativeDenom {
			// If any other token besides the native denom is present, set the flag to false
			onlyNativeDenom = false
			break
		}
	}

	var feeCoins sdk.Coins

	// if IBC Denom fetch the balance of did module.
	if onlyNativeDenom {
		feeCoins = feeTx.GetFee()
	} else {
		addr := td.accountKeeper.GetModuleAddress(feemarkettypes.FeeCollectorName)
		feeBal := td.bankKeeper.GetBalance(ctx, addr, nativeDenom)
		feeCoins = sdk.NewCoins(feeBal)
	}

	gas := ctx.GasMeter().GasConsumed() // use context gas consumed

	if len(feeCoins) == 0 && !simulate {
		return ctx, errorsmod.Wrapf(feemarkettypes.ErrNoFeeCoins, "got length %d", len(feeCoins))
	}
	if len(feeCoins) > 1 {
		return ctx, errorsmod.Wrapf(feemarkettypes.ErrTooManyFeeCoins, "got length %d", len(feeCoins))
	}

	var feeCoin sdk.Coin
	if simulate && len(feeCoins) == 0 {
		// if simulating and user did not provider a fee - create a dummy value for them
		feeCoin = sdk.NewCoin(params.FeeDenom, math.OneInt())
	} else {
		feeCoin = feeCoins[0]
	}

	feeGas := int64(feeTx.GetGas())

	var (
		tip     = sdk.NewCoin(feeCoin.Denom, math.ZeroInt())
		payCoin = feeCoin
	)

	minGasPrice, err := td.feemarketKeeper.GetMinGasPrice(ctx, feeCoin.GetDenom())
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to get min gas price for denom %s", feeCoins[0].GetDenom())
	}

	ctx.Logger().Info("fee deduct post handle",
		"min gas prices", minGasPrice,
		"gas consumed", gas,
	)

	if !simulate {
		payCoin, tip, err = pointidante.CheckTxFee(ctx, minGasPrice, feeCoin, feeGas, false)
		if err != nil {
			return ctx, err
		}
	}

	ctx.Logger().Info("fee deduct post handle",
		"fee", payCoin,
		"tip", tip,
	)
	if err := td.PayOutFeeAndTip(ctx, payCoin, tip); err != nil {
		return ctx, err
	}
	err = state.Update(gas, params)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to update fee market state")
	}

	err = td.feemarketKeeper.SetState(ctx, state)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "unable to set fee market state")
	}

	return next(ctx, tx, simulate, success)
}

// PayOutFeeAndTip deducts the provided fee and tip from the fee payer.
// If the tx uses a feegranter, the fee granter address will pay the fee instead of the tx signer.
func (td TaxDecorator) PayOutFeeAndTip(ctx sdk.Context, fee, tip sdk.Coin) error {
	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("error getting feemarket params: %v", err)
	}

	var events sdk.Events
	// deduct the fees and tip
	if !fee.IsNil() {
		err := DeductCoins(td.bankKeeper, ctx, sdk.NewCoins(fee), params.DistributeFees)
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			feemarkettypes.EventTypeFeePay,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		))
	}

	proposer := sdk.AccAddress(ctx.BlockHeader().ProposerAddress)
	if !tip.IsNil() {
		err := SendTip(td.bankKeeper, ctx, proposer, sdk.NewCoins(tip))
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			feemarkettypes.EventTypeTipPay,
			sdk.NewAttribute(feemarkettypes.AttributeKeyTip, tip.String()),
			sdk.NewAttribute(feemarkettypes.AttributeKeyTipPayee, proposer.String()),
		))
	}

	ctx.EventManager().EmitEvents(events)
	return nil
}

// DeductCoins deducts coins from the given account.
// Coins can be sent to the default fee collector (
// causes coins to be distributed to stakers) or kept in the fee collector account (soft burn).
func DeductCoins(bankKeeper BankKeeper, ctx sdk.Context, coins sdk.Coins, distributeFees bool) error {
	if distributeFees {
		err := bankKeeper.SendCoinsFromModuleToModule(ctx, feemarkettypes.FeeCollectorName, types.FeeCollectorName, coins)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendTip sends a tip to the current block proposer.
func SendTip(bankKeeper BankKeeper, ctx sdk.Context, proposer sdk.AccAddress, coins sdk.Coins) error {
	err := bankKeeper.SendCoinsFromModuleToAccount(ctx, feemarkettypes.FeeCollectorName, proposer, coins)
	if err != nil {
		return err
	}

	return nil
}

// isTaxable returns true if the message is taxable and returns
func (td TaxDecorator) isTaxable(ctx sdk.Context, sdkTx sdk.Tx) (rewards sdk.Coins, burn sdk.Coins, taxable bool, err error) {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdk.Coins{}, sdk.Coins{}, false, errorsmod.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", sdkTx)
	}
	// run lite validation
	taxable = pointidante.IsTaxableTxLite(feeTx)
	if taxable {
		// run full validation
		_, rewards, burn = pointidante.IsTaxableTx(ctx, td.didKeeper, td.resourceKeeper, feeTx)
		return rewards, burn, taxable, nil
	}

	return rewards, burn, taxable, err
}

// getFeePayer returns the fee payer and checks if a fee grant exists
func (td TaxDecorator) getFeePayer(ctx sdk.Context, feeTx sdk.FeeTx, tax sdk.Coins, msgs []sdk.Msg) (types.AccountI, error) {
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFrom := feePayer
	if feeGranter != nil {
		// check if fee grant is supported
		if td.feegrantKeeper == nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			// check if fee grant exists
			err := td.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, tax, msgs)
			if err != nil {
				return nil, errorsmod.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}
		deductFrom = feeGranter
	}

	deductFromAcc := td.accountKeeper.GetAccount(ctx, deductFrom)
	if deductFromAcc == nil {
		return nil, sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFrom)
	}

	return deductFromAcc, nil
}

func (td TaxDecorator) validateTax(tax sdk.Coins, simulate bool) error {
	// no-op if simulate
	if simulate {
		return nil
	}
	// check if denom is accepted
	if !tax.DenomsSubsetOf(sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1)))) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: %s", tax)
	}
	// check if tax is positive
	if !tax.IsAllPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid tax: %s", tax)
	}
	return nil
}

// deductTaxFromFeePayer deducts fees from the account
func (td TaxDecorator) deductTaxFromFeePayer(ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	// ensure module account has been set
	if addr := td.accountKeeper.GetModuleAddress(didtypes.ModuleName); addr == nil {
		return fmt.Errorf("pointidentity fee collector module account (%s) has not been set", didtypes.ModuleName)
	}
	// deduct fees to did module account
	err := td.bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), didtypes.ModuleName, fees)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "failed to deduct fees from %s: %s", acc.GetAddress(), err)
	}
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeTx,
				sdk.NewAttribute(sdk.AttributeKeyFee, fees.String()),
				sdk.NewAttribute(sdk.AttributeKeyFeePayer, acc.GetAddress().String()),
			),
		},
	)
	return nil
}

// distributeRewards distributes rewards to the fee collector
func (td TaxDecorator) distributeRewards(ctx sdk.Context, rewards sdk.Coins) error {
	// move rewards to fee collector
	err := td.bankKeeper.SendCoinsFromModuleToModule(ctx, didtypes.ModuleName, types.FeeCollectorName, rewards)
	if err != nil {
		return err
	}
	return nil
}

// burnFees burns fees from the module account
func (td TaxDecorator) burnFees(ctx sdk.Context, fees sdk.Coins) error {
	// burn fees
	err := td.bankKeeper.BurnCoins(ctx, didtypes.ModuleName, fees)
	if err != nil {
		return err
	}
	return nil
}

func (td *TaxDecorator) handleTaxableTransaction(
	ctx sdk.Context,
	feeTx sdk.FeeTx,
	simulate bool,
	rewards, burn sdk.Coins,
	tx sdk.Tx,
) error {
	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	nativeDenom := params.FeeDenom

	// Check if fees contain only the native denom, if it is IBC denom fee abs logic applies.
	// The fees will be deducted from the fee payer and the tokens already sent to did module account.
	onlyNativeDenom := true
	for _, fee := range feeTx.GetFee() {
		if fee.Denom != nativeDenom {
			// If any other token besides the native denom is present, set the flag to false
			onlyNativeDenom = false
			break
		}
	}

	// if the fees are only in native denom then fee-abs logic won't be applied and deduct the fees from fee payer.
	if onlyNativeDenom {
		// Validate the tax
		if err := td.validateTax(feeTx.GetFee(), simulate); err != nil {
			return err
		}

		// Get fee payer and check if fee grant exists
		tax := rewards.Add(burn...)
		feePayer, err := td.getFeePayer(ctx, feeTx, tax, tx.GetMsgs())
		if err != nil {
			return err
		}
		// Deduct tax from fee payer
		if err := td.deductTaxFromFeePayer(ctx, feePayer, tax); err != nil {
			return err
		}
	}

	// Distribute rewards to fee collector
	if err := td.distributeRewards(ctx, rewards); err != nil {
		return err
	}

	// Burn the tax portion
	if err := td.burnFees(ctx, burn); err != nil {
		return err
	}

	return nil
}
