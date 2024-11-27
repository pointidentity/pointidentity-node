package utils_test

import (
	. "github.com/pointidentity/pointidentity-node/x/did/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DID checks", func() {
	type TestCaseUUIDIdStruct struct {
		inputDID    string
		expectedDID string
	}

	DescribeTable("Check is valid ID (for example did:pointiD:testnet:zABCDEFG123456789abcd and ID is zABCDEFG123456789abcd)",

		func(expected bool, did string) {
			Expect(IsValidID(did)).To(Equal(expected))
		},
		Entry("Valid indy-style identifier 22 chars", true, "FgvC1XcJMRcdRz243A38s5"),
		Entry("Valid indy-style identifier 21 chars", true, "esUzsmZQKCjHHhdkQ5vJA"),
		Entry("Valid UUID", true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"),
		Entry("Not valid indy-style identifier length", false, "esUzsmZQKCjHHhdkQ5vJ"),
		Entry("Not valid, not base58 symbols", false, "12345678abcdIlO0"),
		Entry("Not valid, length", false, "sdf"),
		Entry("Not valid, length and format", false, "sdf:sdf"),
		Entry("Not valid, length", false, "12345"),
	)

	DescribeTable("DID validation",

		func(expected bool, did string, method string, allowedNamespaces []string) {
			Expect(IsValidDID(did, method, allowedNamespaces)).To(Equal(expected))
		},

		Entry("Valid: Inputs: Method and namespace are set", true, "did:pointid:testnet:zABCDEFG123456789abcd", "cheqd", []string{"testnet"}),
		Entry("Valid: Inputs: Method and namespaces are set", true, "did:pointid:testnet:zABCDEFG123456789abcd", "cheqd", []string{"testnet", "mainnet"}),
		Entry("Valid: Inputs: Method not set", true, "did:pointid:testnet:zABCDEFG123456789abcd", "", []string{"testnet"}),
		Entry("Valid: Inputs: Method and namespaces are empty", true, "did:pointid:testnet:zABCDEFG123456789abcd", "", []string{}),
		Entry("Valid: Namespace is absent in DID", true, "did:pointid:zABCDEFG123456789abcd", "", []string{}),
		// Generic method validation
		Entry("Valid: Inputs: Method is not set and passed for NOTpointid", true, "did:NOTpointid:zABCDEFG123456789abcd", "", []string{}),

		Entry("Valid: Inputs: Order of namespaces changed", true, "did:pointid:testnet:zABCDEFG123456789abcd", "pointid", []string{"mainnet", "testnet"}),
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		Entry("Not valid: DID is not started from 'did'", false, "did123:pointid:::zABCDEFG123456789abcd", "pointid", []string{"testnet"}),
		Entry("Not valid: empty namespace", false, "did:pointid::zABCDEFG123456789abcd", "pointid", []string{"testnet"}),
		Entry("Not valid: a lot of ':'", false, "did:pointid:::zABCDEFG123456789abcd", "pointid", []string{"testnet"}),
		Entry("Not valid: several DIDs in one string", false, "did:pointid:testnet:zABCDEFG123456789abcd:did:pointid:testnet:zABCDEFG123456789abcd", "pointid", []string{"testnet"}),
		// Wrong method
		Entry("Not valid: method in DID is not the same as from Inputs", false, "did:NOTpointid:testnet:zABCDEFG123456789abcd", "pointid", []string{"mainnet", "testnet"}),
		Entry("Not valid: method in Inputs is not the same as from DID", false, "did:pointid:testnet:zABCDEFG123456789abcd", "NOTpointid", []string{"mainnet", "testnet"}),
		// Wrong namespace (^[a-zA-Z0-9]*)
		Entry("Not valid: / is not allowed for namespace", false, "did:pointid:testnet/:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: _ is not allowed for namespace", false, "did:pointid:testnet_:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: % is not allowed for namespace", false, "did:pointid:testnet%:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: * is not allowed for namespace", false, "did:pointid:testnet*:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: & is not allowed for namespace", false, "did:pointid:testnet&:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: @ is not allowed for namespace", false, "did:pointid:testnet@:zABCDEFG123456789abcd", "pointid", []string{}),
		Entry("Not valid: namespace from Inputs is not the same as from DID", false, "did:pointid:testnet:zABCDEFG123456789abcd", "pointid", []string{"not_testnet"}),
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		Entry("Not valid: O - is not allowed for UniqueID", false, "did:pointid:testnet:123456789abcdefO", "pointid", []string{}),
		Entry("Not valid: I - is not allowed for UniqueID", false, "did:pointid:testnet:123456789abcdefI", "pointid", []string{}),
		Entry("Not valid: l - is not allowed for UniqueID", false, "did:pointid:testnet:123456789abcdefl", "pointid", []string{}),
		Entry("Not valid: 0 - is not allowed for UniqueID", false, "did:pointid:testnet:123456789abcdef0", "pointid", []string{}),
		// Length checks (should be exactly 16 bytes)
		Entry("Not valid: UniqueID less then 16 bytes", false, "did:pointid:testnet:123", "pointid", []string{}),
	)

	// The next test will check more functionality aspects. This one is for testing corner cases
	Describe("Check splitting functionality", func() {
		Context("Valid DID", func() {
			It("should return expected method, namespace and DID", func() {
				method, namespace, id := MustSplitDID("did:pointid:mainnet:qqqqqqqqqqqqqqqq")
				Expect(method).To(Equal("cheqd"))
				Expect(namespace).To(Equal("mainnet"))
				Expect(id).To(Equal("qqqqqqqqqqqqqqqq"))
			})
		})

		Context("Not valid DID string at all", func() {
			It("should panic", func() {
				panicDID := "Not 	"
				Expect(func() {
					MustSplitDID(panicDID)
				}).To(Panic())
			})
		})
	})

	DescribeTable("Check DID splitting and joining",

		func(did string) {
			method, namespace, id := MustSplitDID(did)
			Expect(did).To(Equal(JoinDID(method, namespace, id)))
		},
		Entry("Full DID", "did:cheqd:testnet:zABCDEFG123456789abcd"),
		Entry("Without namespace", "did:cheqd:zABCDEFG123456789abcd"),
		Entry("Not cheqd method", "did:NOTcheqd:zABCDEFG123456789abcd"),
	)

	DescribeTable("UUID test cases", func(testCase TestCaseUUIDIdStruct) {
		result := NormalizeDID(testCase.inputDID)
		Expect(result).To(Equal(testCase.expectedDID))
	},

		Entry(
			"base58 identifier - not changed",
			TestCaseUUIDIdStruct{
				inputDID:    "did:cheqd:testnet:zABCDEFG123456789abcd",
				expectedDID: "did:cheqd:testnet:zABCDEFG123456789abcd",
			}),

		Entry(
			"Mixed case UUID",
			TestCaseUUIDIdStruct{
				inputDID:    "did:cheqd:testnet:BAbbba14-f294-458a-9b9c-474d188680fd",
				expectedDID: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Low case UUID",
			TestCaseUUIDIdStruct{
				inputDID:    "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
				expectedDID: "did:cheqd:testnet:babbba14-f294-458a-9b9c-474d188680fd",
			}),

		Entry(
			"Upper case UUID",
			TestCaseUUIDIdStruct{
				inputDID:    "did:cheqd:testnet:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
				expectedDID: "did:cheqd:testnet:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			}),
	)
	// ToDo: tests for list of DIDs
})
