package astral

import (
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Validator", func() {
	ginkgo.Context("Integers", func() {
		ginkgo.It("Should error on a minimum value that is out of bounds", func() {
			err := validateInt(nil, &Argument{
				Min: int64(1),
			}, int64(0))

			Expect(err).ToNot(BeNil())
		})
		ginkgo.It("Should error on a maximum value that is out of bounds", func() {
			err := validateInt(nil, &Argument{
				Max: int64(1),
			}, int64(2))

			Expect(err).ToNot(BeNil())
		})
	})
	ginkgo.Context("Floats", func() {
		ginkgo.It("Should error on a minimum value that is out of bounds", func() {
			err := validateFloat(nil, &Argument{
				Min: float64(1),
			}, float64(0))

			Expect(err).ToNot(BeNil())
		})
		ginkgo.It("Should error on a maximum value that is out of bounds", func() {
			err := validateFloat(nil, &Argument{
				Max: float64(1),
			}, float64(2))

			Expect(err).ToNot(BeNil())
		})
	})
	ginkgo.Context("Strings", func() {
		ginkgo.It("Should error on a minimum value that is out of bounds", func() {
			err := validateString(nil, &Argument{
				Min: int64(1),
			}, "")

			Expect(err).ToNot(BeNil())
		})
		ginkgo.It("Should error on a maximum value that is out of bounds", func() {
			err := validateString(nil, &Argument{
				Max: int64(1),
			}, "ab")

			Expect(err).ToNot(BeNil())
		})
	})
})
