package astral

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
	ginkgo.Context("Groups", func() {
		ginkgo.It("Should copy middleware from parent to group", func() {
			parent.Group(func(r *Route) {
				r.Use(func(fn Handler) Handler {
					return func(ctx *Context) {
						fn(ctx)
					}
				})

				r.On("sub", func(ctx *Context) {
					// Nothing
				})
			})

			test := parent.Find("sub")

			Expect(test).ToNot(BeNil())
			Expect(len(test.middleware)).To(Equal(1))
		})
		ginkgo.It("Should assign parent when using group", func() {
			sub := parent.On("sub", nil)

			sub.Group(func(r *Route) {
				r.On("test", func(ctx *Context) {
					// Nothing
				})
			})

			test := parent.Find("sub", "test")

			Expect(test).ToNot(BeNil())
			Expect(test.Path()).To(Equal([]string{"sub", "test"}))
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
