package integration_test

import (
	. "github.com/loveandpeople-DAG/goClient/api"
	_ "github.com/loveandpeople-DAG/goClient/api/integration/gocks"
	. "github.com/loveandpeople-DAG/goClient/api/integration/samples"
	"github.com/loveandpeople-DAG/goClient/checksum"
	. "github.com/loveandpeople-DAG/goClient/consts"
	. "github.com/loveandpeople-DAG/goClient/trinary"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("FindTransactionObjects()", func() {
	api, err := ComposeAPI(HTTPClientSettings{}, nil)
	if err != nil {
		panic(err)
	}

	Context("call", func() {
		It("resolves to correct response", func() {
			addrWithChecksum, err := checksum.AddChecksum(Bundle[0].Address, true, AddressChecksumTrytesSize)
			Expect(err).ToNot(HaveOccurred())
			txs, err := api.FindTransactionObjects(FindTransactionsQuery{
				Addresses: Hashes{addrWithChecksum},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(txs[0]).To(Equal(Bundle[0]))
		})
	})

	Context("invalid input", func() {
		It("returns an error for invalid hashes", func() {
			_, err = api.GetTransactionObjects("asdf")
			Expect(errors.Cause(err)).To(Equal(ErrInvalidTransactionHash))
		})
	})
})
