package types_test

import (
	. "github.com/pointidentity/pointidentity-node/x/did/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create DID Payload Validation tests", func() {
	type TestCaseUUIDDidStruct struct {
		inputID    string
		expectedID string
	}

	DescribeTable("UUID validation tests", func(testCase TestCaseUUIDDidStruct) {
		inputMsg := MsgCreateDidDocPayload{
			Id:             testCase.inputID,
			Authentication: []string{testCase.inputID + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:                     testCase.inputID + "#key1",
					VerificationMethodType: Ed25519VerificationKey2020Type,
					Controller:             testCase.inputID,
				},
			},
		}
		expectedMsg := MsgCreateDidDocPayload{
			Id:             testCase.expectedID,
			Authentication: []string{testCase.expectedID + "#key1"},
			VerificationMethod: []*VerificationMethod{
				{
					Id:                     testCase.expectedID + "#key1",
					VerificationMethodType: Ed25519VerificationKey2020Type,
					Controller:             testCase.expectedID,
				},
			},
		}
		inputMsg.Normalize()
		Expect(inputMsg).To(Equal(expectedMsg))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDDidStruct{
				inputID:    "did:pointid:testnet:zABCDEFG123456789abcd",
				expectedID: "did:pointid:testnet:zABCDEFG123456789abcd",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:pointid:testnet:BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedID: "did:pointid:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:pointid:testnet:babbba14-f294-458a-9b9c-474d188680fd",
				expectedID: "did:pointid:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDDidStruct{
				inputID:    "did:pointid:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedID: "did:pointid:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
})
