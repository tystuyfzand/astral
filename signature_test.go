package astral

import (
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Signatures", func() {
	var (
		r *Route
	)
	ginkgo.BeforeEach(func() {
		r = New()
	})
	ginkgo.Context("Arguments", func() {
		ginkgo.Context("Strings", func() {
			ginkgo.It("Should properly parse string arguments", func() {
				parseSignature(r, "test <arg>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeBasic))
			})
			ginkgo.It("Should properly parse optional string arguments", func() {
				parseSignature(r, "test [arg]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeBasic))
			})
		})
		ginkgo.Context("Emojis", func() {
			ginkgo.It("Should properly parse emoji arguments", func() {
				parseSignature(r, "test <:emoji>")

				arg := r.Arguments["emoji"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeEmoji))
			})
			ginkgo.It("Should properly parse optional emoji arguments", func() {
				parseSignature(r, "test [:emoji]")

				arg := r.Arguments["emoji"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeEmoji))
			})
		})
		ginkgo.Context("Users", func() {
			ginkgo.It("Should properly parse required mention arguments", func() {
				parseSignature(r, "test <@mention>")

				arg := r.Arguments["mention"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeUserMention))
			})
			ginkgo.It("Should properly parse optional mention arguments", func() {
				parseSignature(r, "test [@mention]")

				arg := r.Arguments["mention"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeUserMention))
			})
		})
		ginkgo.Context("Channels", func() {
			ginkgo.It("Should properly parse channel arguments", func() {
				parseSignature(r, "test <#channel>")

				arg := r.Arguments["channel"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeChannelMention))
			})
			ginkgo.It("Should properly parse optional channel arguments", func() {
				parseSignature(r, "test [#channel]")

				arg := r.Arguments["channel"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeChannelMention))
			})
		})
		ginkgo.Context("Ints", func() {
			ginkgo.It("Should properly parse int arguments", func() {
				parseSignature(r, "test <arg int>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeInt))
			})
			ginkgo.It("Should properly parse optional int arguments", func() {
				parseSignature(r, "test [arg int]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeInt))
			})
			ginkgo.It("Should properly parse int arguments with minimums", func() {
				parseSignature(r, "test <arg int min:2>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeInt))
				Expect(arg.Min).To(BeEquivalentTo(2))
			})
			ginkgo.It("Should properly parse optional int arguments with minimums", func() {
				parseSignature(r, "test [arg int min:3 max:100]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeInt))
				Expect(arg.Min).To(BeEquivalentTo(3))
				Expect(arg.Max).To(BeEquivalentTo(100))
			})
		})
		ginkgo.Context("Floats", func() {
			ginkgo.It("Should properly parse float arguments", func() {
				parseSignature(r, "test <arg float>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeFloat))
			})
			ginkgo.It("Should properly parse optional float arguments", func() {
				parseSignature(r, "test [arg float]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeFloat))
			})
			ginkgo.It("Should properly parse float arguments with minimums", func() {
				parseSignature(r, "test <arg float min:2>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeFloat))
				Expect(arg.Min).To(BeEquivalentTo(2))
			})
			ginkgo.It("Should properly parse optional float arguments with minimums", func() {
				parseSignature(r, "test [arg float min:3]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeFloat))
				Expect(arg.Min).To(BeEquivalentTo(3))
			})
		})
		ginkgo.Context("Booleans", func() {
			ginkgo.It("Should properly parse bool arguments", func() {
				parseSignature(r, "test <arg bool>")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeTrue())
				Expect(arg.Type).To(Equal(ArgumentTypeBool))
			})
			ginkgo.It("Should properly parse optional bool arguments", func() {
				parseSignature(r, "test [arg bool]")

				arg := r.Arguments["arg"]

				Expect(arg).ToNot(BeNil())
				Expect(arg.Required).To(BeFalse())
				Expect(arg.Type).To(Equal(ArgumentTypeBool))
			})
		})
	})
})
