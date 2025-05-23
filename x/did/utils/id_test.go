package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/pointidentity/pointidentity-node/x/did/utils"
)

var _ = Describe("Identifier Validation tests", func() {
	DescribeTable("Id validation tests", func(isValid bool, id string) {
		isDid := IsValidID(id)

		if isValid {
			Expect(isDid).To(BeTrue())
		} else {
			Expect(isDid).To(BeFalse())
		}
	},

		Entry("Base58 string, 16 bytes", true, "zABCDEFG123456789abcd"),
		Entry("UUID string", true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"),
		Entry("Too short", false, "sdf"),
		Entry("Unexpected :", false, "sdf:sdf"),
		Entry("Too short", false, "12345"),
	)
})
