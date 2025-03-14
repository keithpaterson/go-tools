package resolver

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Property Resolver", func() {
	type expectation struct {
		value string
		ok    bool
	}
	var (
		cfg      *ResolverConfig
		root     *rootResolver
		resolver *propertiesResolver
	)
	BeforeEach(func() {
		cfg = &ResolverConfig{}
		root = NewResolver(cfg)
		resolver = NewPropertiesResolver(root.WithResolver("env", NewEnvResolver(root)))

		// must register myself with root to support recursive lookups
		root.WithResolver("prop", resolver)
	})

	DescribeTable("resolve",
		func(properties Properties, tokens []string, expect []expectation) {
			// Arrange
			cfg.Properties = properties

			// Act & Assert
			for index, token := range tokens {
				actual, ok := resolver.resolve("prop", token)
				Expect(actual).To(Equal(expect[index].value))
				Expect(ok).To(Equal(expect[index].ok))
			}
		},
		Entry("no props, one token, no change", Properties{}, []string{"input"}, []expectation{{"input", false}}),
		Entry("props, one token, replaced", Properties{"input": "test"}, []string{"input"}, []expectation{{"test", true}}),
		Entry("props with token, one token, should resolve token from env", Properties{"input": "${prop:foo}", "foo": "bar"}, []string{"input"}, []expectation{{"bar", true}}),
		Entry("props with token missing prefix, one token, should resolve token from env", Properties{"input": "${foo}", "foo": "bar"}, []string{"input"}, []expectation{{"bar", true}}),
		Entry("multiple tokens with expectations", Properties{"input": "test", "foo": "bar"},
			[]string{"input", "hello", "barcelona", "foo"},
			[]expectation{{"test", true}, {"hello", false}, {"barcelona", false}, {"bar", true}},
		),
	)

	It("will ignore requests for invalid token name", func() {
		// Arrange

		// Act
		actual, ok := resolver.resolve("crumb", "fling")

		// Assert
		Expect(ok).To(BeFalse())
		Expect(actual).To(Equal("fling"))
	})
})
