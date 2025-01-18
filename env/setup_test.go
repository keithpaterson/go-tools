package env

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VarSetup", func() {
	type expectations struct {
		afterApply  Setup
		afterRevert Setup
	}

	DescribeTable("Apply and Revert",
		func(startEnv Setup, input Setup, expect expectations) {
			// Arrange
			resetTest := startEnv.Apply()
			defer resetTest.Apply()

			// Act
			afterApply := input.Apply()
			afterRevert := afterApply.Apply()

			// Assert
			Expect(afterApply).To(Equal(expect.afterApply))
			Expect(afterRevert).To(Equal(expect.afterRevert))
		},
		Entry("empty env",
			New(),
			New(),
			expectations{New(), New()}),
		Entry("add env",
			New(),
			New().Set("flim", "flam"),
			expectations{New().Unset("flim"), New().Set("flim", "flam")}),
		Entry("remove env",
			New().Set("flim", "flam"),
			New().Unset("flim"),
			expectations{New().Set("flim", "flam"), New().Unset("flim")}),
		Entry("replace env",
			New().Set("flim", "flam"),
			New().Set("flim", "pool"),
			expectations{New().Set("flim", "flam"), New().Set("flim", "pool")}),
		Entry("add, replace and unset",
			New().Set("flim", "flam").Set("junk", "trunk"),
			New().Set("foo", "bar").Set("junk", "pile").Unset("flim"),
			expectations{
				New().Unset("foo").Set("junk", "trunk").Set("flim", "flam"),
				New().Set("foo", "bar").Set("junk", "pile").Unset("flim"),
			}),
		Entry("add with integer",
			New(),
			New().Set("testint", 100),
			expectations{
				New().Unset("testint"),
				New().Set("testint", "100"),
			}),
	)

	Describe("validate actual env changes", func() {
		// It's worth noting that this test confirms that we can both add and remove environment variables.
		It("should add, update, and unset properly", func() {
			// start with no variable set
			name := "__test_add_env_var__"
			_, ok := os.LookupEnv(name)
			Expect(ok).To(BeFalse())

			// add the variaable
			setup := New().Set(name, "added")
			remove := setup.Apply()
			value, ok := os.LookupEnv(name)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("added"))

			// update the variable
			update := New().Set(name, "updated")
			update.Apply() // don't need to capture the result here
			value, ok = os.LookupEnv(name)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("updated"))

			// remove the variable
			remove.Apply()
			_, ok = os.LookupEnv(name)
			Expect(ok).To(BeFalse())
		})
	})
})
