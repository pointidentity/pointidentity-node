package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/pointidentity/pointidentity-node/x/did/types"
)

var _ = Describe("Service tests", func() {
	type TestCaseServiceStruct struct {
		service           *Service
		baseDid           string
		allowedNamespaces []string
		isValid           bool
		errorMsg          string
	}

	DescribeTable("Service validation tests", func(testCase TestCaseServiceStruct) {
		err := testCase.service.Validate(testCase.baseDid, testCase.allowedNamespaces)

		if testCase.isValid {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testCase.errorMsg))
		}
	},

		Entry(
			"Valid service entry",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:pointid:aABCDEFG123456789abcd#service1",
					ServiceType:     "DIDCommMessaging",
					ServiceEndpoint: []string{"endpoint1", "endpoint2"},
				},
				baseDid:           "did:pointid:aABCDEFG123456789abcd",
				allowedNamespaces: []string{""},
				isValid:           true,
				errorMsg:          "",
			}),

		Entry(
			"Namespace is not allowed",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:pointid:zABCDEFG123456789abcd#service1",
					ServiceType:     "DIDCommMessaging",
					ServiceEndpoint: []string{"endpoint"},
				},
				allowedNamespaces: []string{"mainnet"},
				isValid:           false,
				errorMsg:          "id: did namespace must be one of: mainnet.",
			}),

		Entry(
			"Base DID is not the same as in ID",
			TestCaseServiceStruct{
				service: &Service{
					Id:              "did:pointid:zABCDEFG123456789abcd#service1",
					ServiceType:     "DIDCommMessaging",
					ServiceEndpoint: []string{"endpoint"},
				},
				baseDid:  "did:pointid:zABCDEFG987654321abcd",
				isValid:  false,
				errorMsg: "id: must have prefix: did:pointid:zABCDEFG987654321abcd",
			}),
	)
})
