package resolver

import (
	"github.com/keithpaterson/go-tools/env"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment Resolver", func() {
	type expectation struct {
		value string
		ok    bool
	}
	var (
		root     *rootResolver
		resolver *envResolver
	)
	BeforeEach(func() {
		root = NewResolver(&ResolverConfig{})
		resolver = NewEnvResolver(root)

		// must register myself with root to support recursive lookups
		root.WithResolver("env", resolver)
	})

	DescribeTable("resolve",
		func(startEnv env.Setup, tokens []string, expect []expectation) {
			// Arrange
			origEnv := startEnv.Apply()
			defer origEnv.Apply()

			// Act & Assert
			for index, token := range tokens {
				actual, ok := resolver.resolve("env", token)
				Expect(actual).To(Equal(expect[index].value))
				Expect(ok).To(Equal(expect[index].ok))
			}
		},
		Entry("no env, one token, no change", env.New(), []string{"input"}, []expectation{{"input", false}}),
		Entry("env, one token, replaced", env.New().Set("input", "test"), []string{"input"}, []expectation{{"test", true}}),
		Entry("env with token, one token, should resolve token from env", env.New().Set("input", "${env:foo}").Set("foo", "bar"), []string{"input"}, []expectation{{"bar", true}}),
		Entry("multiple tokens with expectations", env.New().Set("input", "test").Set("foo", "bar"),
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
