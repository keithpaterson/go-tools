package resolver

import (
	"github.com/keithpaterson/go-tools/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolver", func() {
	var (
		cfg      *ResolverConfig
		origEnv  env.Setup
		resolver *rootResolver
	)
	BeforeEach(func() {
		cfg = &ResolverConfig{}
		cfg.Properties = Properties{"cat": "mousse", "dog": "${cat}", "pig": "manbear${prop:cat}", "cow": "moo${env:foo}"}
		origEnv = env.New().
			Set("foo", "bar").Set("big", "${env:foo}").Set("small", "${prop:pig}").Set("can", "${cow}").
			Unset("thing").Unset("wing").Unset("ding").
			Apply()

		resolver = NewResolver(cfg).WithStandardResolvers()
	})
	AfterEach(func() {
		origEnv.Apply()
	})

	DescribeTable("Resolve",
		func(input string, expect string) {
			// Act
			actual := resolver.Resolve(input)

			// Assert
			Expect(actual).To(Equal(expect))
		},
		Entry("with no token matches, replaces nothing", "test: ${env:thing}, ${bar}", "test: ${env:thing}, ${bar}"),
		Entry("with invalid token names, replaces nothing", "test: ${crumb:size}", "test: ${crumb:size}"),
		Entry("with date, time token matches, computes values", "test: ${date:2011-05-09} ${time:11:22:44}", "test: 2011-05-09 11:22:44"),
		Entry("with prop matches, computes values", "test: ${prop:cat} ${dog} ${cow}", "test: mousse mousse moobar"),
		Entry("with env matches, computes values", "test: ${env:foo} ${env:big} ${env:small} ${env:can}", "test: bar bar manbearmousse moobar"),
	)

	It("will resolve tokens dynamically", func() {
		// Act
		actual := resolver.Resolve("${env:small}")
		Expect(actual).To(Equal("manbearmousse"))

		cfg.Properties["pig"] = "tandem"
		actual = resolver.Resolve("${env:small}")
		Expect(actual).To(Equal("tandem"))
	})
})
