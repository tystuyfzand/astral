package router

import "testing"

import (
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Route", func() {
	var (
		parent *Route
	)
	ginkgo.BeforeEach(func() {
		parent = New()
	})
	ginkgo.Context("Paths", func() {
		ginkgo.It("Should properly serialize path", func() {
			p := parent.On("test", nil).On("something", nil).On("deeper", nil).Path()

			Expect(p).To(Equal([]string{"test", "something", "deeper"}))
		})
	})
	ginkgo.Context("Validation", func() {
		ginkgo.Context("Choices", func() {
			var (
				r *Route
			)
			ginkgo.BeforeEach(func() {
				parent = parent.On("content", nil)

				r = parent.On("track <type> <name>", nil)
				r.Arguments["type"].Choices = []StringChoice{
					{Name: "Twitch", Value: "twitch"},
				}
			})
			ginkgo.It("Should fail validation when choice is not valid", func() {
				err := r.Validate(&Context{
					Arguments: map[string]interface{}{
						"type": "youtube",
						"name": "test",
					},
				})

				Expect(err).ToNot(BeNil())
			})
			ginkgo.It("Should not fail validation with choices", func() {
				err := r.Validate(&Context{
					Arguments: map[string]interface{}{
						"type": "twitch",
						"name": "test",
					},
				})

				Expect(err).To(BeNil())
			})
		})
	})
})

func TestRoute_Validate(t *testing.T) {

}
