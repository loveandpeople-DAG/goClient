package integration_test

import (
	. "github.com/loveandpeople-DAG/goClient/api"
	. "github.com/loveandpeople-DAG/goClient/api/integration/samples"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsAddressUsed()", func() {

	api, err := ComposeAPI(HTTPClientSettings{}, nil)
	if err != nil {
		panic(err)
	}

	It("returns true for spent address", func() {
		used, err := api.IsAddressUsed(SampleAddressesWithChecksum[0])
		Expect(err).ToNot(HaveOccurred())
		Expect(used).To(BeTrue())
	})

	It("returns true for address with transactions", func() {
		used, err := api.IsAddressUsed(SampleAddressesWithChecksum[1])
		Expect(err).ToNot(HaveOccurred())
		Expect(used).To(BeTrue())
	})

	It("returns false for unused address", func() {
		used, err := api.IsAddressUsed(SampleAddressesWithChecksum[2])
		Expect(err).ToNot(HaveOccurred())
		Expect(used).To(BeFalse())
	})

})
