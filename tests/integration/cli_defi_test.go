//go:build integration

package integration

import (
	cli "github.com/pointidentity/pointidentity-node/tests/integration/cli"
	helpers "github.com/pointidentity/pointidentity-node/tests/integration/helpers"

	"github.com/pointidentity/pointidentity-node/tests/integration/testdata"
	didtypes "github.com/pointidentity/pointidentity-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Burn coins from relevant message signer", func() {
	It("should burn the coins from the given address", func() {
		// define the coins to burn, in which case 1,000,000 upoint or 0.01 POINT
		burnCoins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000))

		// get the balance of the account before burning the coins
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// generate fixed fees, in which case 500,000,000 upoint or 0.5 POINT
		fees := helpers.GenerateFees("500000000upoint")

		// burn the coins
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_1, burnCoins.String(), fees)

		// assert no error
		Expect(err).To(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))

		// get the balance of the account after burning the coins
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_1_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// calculate the difference between the balance before and after burning the coins
		diff := balanceBefore.Sub(balanceAfter)

		// assert the difference is equal to the coins burnt
		total := burnCoins.Add(sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(500_000_000)))

		// assert the difference is equal to the coins burnt
		Expect(diff).To(Equal(total))
	})

	It("shouldn't burn as their are insufficient funds in the sender", func() {
		// define the coins to burn, in which case 1,000,000 upoint or 0.01 POINT
		coins := sdk.NewCoin(didtypes.BaseMinimalDenom, sdk.NewInt(1_000_000))

		// get the balance of the account before burning the coins
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// generate fixed fees, in which case 500,000,000 upoint or 0.5 POINT
		fees := helpers.GenerateFees("500000000upoint")

		// burn the coins
		res, err := cli.BurnMsg(testdata.BASE_ACCOUNT_3, coins.String(), fees)

		// assert no error
		Expect(err).NotTo(BeNil())

		// assert the response code is 0
		Expect(res.Code).To(BeEquivalentTo(0))

		// get the balance of the account after burning the coins
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_3_ADDR, didtypes.BaseMinimalDenom)

		// assert no error
		Expect(err).To(BeNil())

		// assert the balance before and after burning the coins are equal
		Expect(balanceBefore).To(Equal(balanceAfter))
	})
})
