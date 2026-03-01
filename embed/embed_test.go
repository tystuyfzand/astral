package embed

import (
	_ "embed"

	"github.com/auroradevllc/astral/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type simpleTestStruct struct {
	Text string `hcl:"text"`
}

type complexBlock struct {
	Label string `hcl:",label"`
}

type complexInStruct struct {
	Text  string       `hcl:"text"`
	Block complexBlock `hcl:"block,block"`
}

type complexOutBlock struct {
	Label string
}

type complexOutStruct struct {
	Text  string
	Block complexOutBlock
}

var _ = Describe("Embeds", func() {
	Context("Simple Decoding", func() {
		const simpleType = "simple"
		BeforeEach(func() {
			RegisterDecoder(simpleType, NewSimpleDecoder[simpleTestStruct]())
		})
		It("Should have registered a simple decoder", func() {
			Expect(HasDecoder(simpleType)).To(BeTrue())
		})
		It("Should decode a simple struct", func() {
			s, err := Render[simpleTestStruct](simpleType, astral.Embed{
				Templates: astral.TemplateMap(map[astral.ClientType]string{
					simpleType: `
text = "Test"
`,
				}),
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(s).ToNot(BeNil())
			Expect(s.Text).To(Equal("Test"))
		})
	})
	Context("Complex Decoding", func() {
		const complexType = "complex"
		BeforeEach(func() {
			RegisterDecoder(complexType, NewComplexDecoder[complexInStruct, complexOutStruct]())
		})
		It("Should have registered a complex decoder", func() {
			Expect(HasDecoder(complexType)).To(BeTrue())
		})
		It("Should decode a complex struct", func() {
			s, err := Render[complexOutStruct](complexType, astral.Embed{
				Templates: astral.TemplateMap(map[astral.ClientType]string{
					complexType: `
text = "Test"

block "Test" {
}
`,
				}),
			})

			Expect(err).NotTo(HaveOccurred())

			// Shouldn't be nil, but check anyway
			Expect(s).ToNot(BeNil())

			// This will be a runtime error if the type isn't complexOutStruct, but good to verify
			Expect(s).To(BeAssignableToTypeOf(&complexOutStruct{}))

			// Validate fields decoded as expected
			Expect(s.Text).To(Equal("Test"))
		})
	})
})
