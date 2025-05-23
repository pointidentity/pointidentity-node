//go:build integration

package integration

import (
	"crypto/ed25519"

	"github.com/pointidentity/pointidentity-node/tests/integration/cli"
	"github.com/pointidentity/pointidentity-node/tests/integration/helpers"
	"github.com/pointidentity/pointidentity-node/tests/integration/network"
	"github.com/pointidentity/pointidentity-node/tests/integration/testdata"
	didcli "github.com/pointidentity/pointidentity-node/x/did/client/cli"
	testsetup "github.com/pointidentity/pointidentity-node/x/did/tests/setup"
	"github.com/pointidentity/pointidentity-node/x/did/types"
	"github.com/google/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// cases:
//   1. fixed fee - invalid fee denom create
//   2. fixed fee - invalid fee denom update
//   3. fixed fee - invalid fee denom deactivate
//   4. fixed fee - invalid fee amount create
//   5. fixed fee - invalid fee amount update
//   6. fixed fee - invalid fee amount deactivate
//   10. fixed fee - charge only tax if fee is more than tax create
//   11. fixed fee - charge only tax if fee is more than tax update
//   12. fixed fee - charge only tax if fee is more than tax deactivate
//   13. fixed fee - insufficient funds create
//   14. fixed fee - insufficient funds update
//   15. fixed fee - insufficient funds deactivate

var _ = Describe("pointidentity cli - negative diddoc pricing", func() {
	var tmpDir string
	var feeParams types.FeeParams
	var payload didcli.DIDDocument
	var signInputs []didcli.SignInput

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		// Query fee params
		res, err := cli.QueryParams(types.ModuleName, string(types.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = helpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())

		// Create a new DID Doc
		did := "did:pointid:" + network.DidNamespace + ":" + uuid.NewString()
		keyId := did + "#key1"

		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		publicKeyMultibase := testsetup.GenerateEd25519VerificationKey2020VerificationMaterial(publicKey)

		payload = didcli.DIDDocument{
			ID: did,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 keyId,
					"type":               "Ed25519VerificationKey2020",
					"controller":         did,
					"publicKeyMultibase": publicKeyMultibase,
				},
			},
			Authentication: []string{keyId},
		}

		signInputs = []didcli.SignInput{
			{
				VerificationMethodID: keyId,
				PrivKey:              privateKey,
			},
		}
	})

	It("should not succeed in create diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting create diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.CreateDid.Amount.Int64()))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in update diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.GetUpdateDid().Amount.Int64()))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, invalid denom", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with invalid denom")
		invalidTax := sdk.NewCoin("invalid", sdk.NewInt(feeParams.GetDeactivateDid().Amount.Int64()))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(invalidTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(10))
	})

	It("should not fail in create diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting create diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.CreateDid.Denom, sdk.NewInt(feeParams.CreateDid.Amount.Int64()-1))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not fail in update diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.UpdateDid.Denom, sdk.NewInt(feeParams.UpdateDid.Amount.Int64()-1))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not fail in deactivate diddoc message - case: fixed fee, lower amount than required", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with lower amount than required")
		lowerTax := sdk.NewCoin(feeParams.DeactivateDid.Denom, sdk.NewInt(feeParams.DeactivateDid.Amount.Int64()-1))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(lowerTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should not charge more than tax for create diddoc message - case: fixed fee", func() {
		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the create diddoc message with double the tax")
		tax := feeParams.CreateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax for update diddoc message - case: fixed fee", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               "Ed25519VerificationKey2020",
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the update diddoc message with double the tax")
		tax := feeParams.UpdateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not charge more than tax for deactivate diddoc message - case: fixed fee", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id:        payload.ID,
			VersionId: uuid.NewString(),
		}

		By("querying the fee payer account balance before the transaction")
		balanceBefore, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("submitting the deactivate diddoc message with double the tax")
		tax := feeParams.DeactivateDid
		doubleTax := sdk.NewCoin(types.BaseMinimalDenom, tax.Amount.Mul(sdk.NewInt(2)))
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(doubleTax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("querying the fee payer account balance after the transaction")
		balanceAfter, err := cli.QueryBalance(testdata.BASE_ACCOUNT_5_ADDR, types.BaseMinimalDenom)
		Expect(err).To(BeNil())

		By("checking that the fee payer account balance has been decreased by the tax")
		diff := balanceBefore.Amount.Sub(balanceAfter.Amount)
		Expect(diff).To(Equal(tax.Amount))
	})

	It("should not succeed in create diddoc create message - case: fixed fee, insufficient funds", func() {
		By("submitting create diddoc message with insufficient funds")
		tax := feeParams.CreateDid
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in update diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_5, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the update diddoc message")
		payload2 := didcli.DIDDocument{
			ID: payload.ID,
			VerificationMethod: []didcli.VerificationMethod{
				map[string]any{
					"id":                 payload.VerificationMethod[0]["id"],
					"type":               payload.VerificationMethod[0]["type"],
					"controller":         payload.VerificationMethod[0]["controller"],
					"publicKeyMultibase": payload.VerificationMethod[0]["publicKeyMultibase"],
				},
			},
			Authentication:  payload.Authentication,
			AssertionMethod: []string{payload.VerificationMethod[0]["id"].(string)}, // <-- changed
		}

		By("submitting update diddoc message with insufficient funds")
		tax := feeParams.UpdateDid
		res, err = cli.UpdateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})

	It("should not succeed in deactivate diddoc message - case: fixed fee, insufficient funds", func() {
		By("submitting the create diddoc message")
		res, err := cli.CreateDidDoc(tmpDir, payload, signInputs, "", testdata.BASE_ACCOUNT_4, helpers.GenerateFees(feeParams.CreateDid.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		By("preparing the deactivate diddoc message")
		payload2 := types.MsgDeactivateDidDocPayload{
			Id: payload.ID,
		}

		By("submitting deactivate diddoc message with insufficient funds")
		tax := feeParams.DeactivateDid
		res, err = cli.DeactivateDidDoc(tmpDir, payload2, signInputs, "", testdata.BASE_ACCOUNT_6, helpers.GenerateFees(tax.String()))
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(5))
	})
})
