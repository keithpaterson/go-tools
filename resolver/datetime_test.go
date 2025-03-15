package resolver

import (
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	// ALWAYS use EDT so that no matter the time zone the "local" tests work
	testLocal = time.Date(2025, 06, 20, 13, 49, 38, 0, time.FixedZone("Test", -4*60*60))

	useLocal = false
)

func testNow() time.Time {
	// always return local time because UTC is a separate call on that
	return testLocal
}

func toEpochStr(date time.Time, mod int64) string {
	return strconv.FormatInt(date.Unix()+mod, 10)
}

var _ = Describe("Date/Time Resolver", func() {
	type expectation struct {
		value string
		ok    bool
	}
	var (
		//root     *rootResolver
		resolver *dateTimeResolver
	)
	BeforeEach(func() {
		//root = NewResolver(&ResolverConfig{})
		resolver = NewDateTimeResolver()
		//resolver.root = root
		resolver.nowFn = testNow
	})

	Describe("Local Time", func() {
		BeforeEach(func() {
			useLocal = true
		})

		DescribeTable("time",
			func(token string, expect expectation) {
				// Arrange
				useLocal = true

				// Act
				actual, ok := resolver.Resolve("time", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now returns time", "now", expectation{"13:49:38", true}),
			Entry("now plus seconds returns correct time", "now + 10s", expectation{"13:49:48", true}),
			Entry("now minus seconds returns correct time", "now - 10s", expectation{"13:49:28", true}),
			Entry("now plus minutes returns correct time", "now + 10m", expectation{"13:59:38", true}),
			Entry("now plus hours returns correct time", "now + 10h", expectation{"23:49:38", true}),
			Entry("now with 24h format returns time", "now.(HHmmss)", expectation{"134938", true}),
			Entry("now with 12h format returns time", "now.(hh.mm.ss-AM)", expectation{"01.49.38-PM", true}),
			Entry("specific time format returns time", "15:23:11.(TimeOnly) - 4m", expectation{"15:19:11", true}),
			Entry("invalid time modifier returns input", "now +", expectation{"now +", false}),
			Entry("invalid time modifier delta returns input", "now + 6", expectation{"now + 6", false}),
			Entry("invalid time modifier delta scope returns input", "now + 6X", expectation{"now + 6X", false}),
		)

		DescribeTable("date",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("date", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now returns date", "now", expectation{"2025-06-20", true}),
			Entry("now plus days returns correct date", "now + 10D", expectation{"2025-06-30", true}),
			Entry("now plus months returns correct date", "now + 5M", expectation{"2025-11-20", true}),
			Entry("now minus months returns correct date", "now - 5M", expectation{"2025-01-20", true}),
			Entry("now plus weeks returns correct date", "now + 2W", expectation{"2025-07-04", true}),
			Entry("now plus years returns correct date", "now + 10Y", expectation{"2035-06-20", true}),
			Entry("now with format returns date", "now.(YYYYMMDD)", expectation{"20250620", true}),
			Entry("now with another format returns date", "now.(YY/MM/DD)", expectation{"25/06/20", true}),
			Entry("specific date returns date", "2012/09/13.(YYYY/MM/DD) - 10D", expectation{"2012/09/03", true}),
		)

		DescribeTable("datetime",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("datetime", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now returns datetime", "now", expectation{"2025-06-20T13:49:38-04:00", true}),
			Entry("now plus days returns correct datetime", "now + 10D", expectation{"2025-06-30T13:49:38-04:00", true}),
			Entry("now plus months returns correct datetime", "now + 5M", expectation{"2025-11-20T13:49:38-04:00", true}),
			Entry("now plus years returns correct datetime", "now + 10Y", expectation{"2035-06-20T13:49:38-04:00", true}),
			Entry("now minus years returns correct datetime", "now - 10Y", expectation{"2015-06-20T13:49:38-04:00", true}),
			Entry("now with 24h format returns datetime", "now.(YYYYMMDDTHHmmss)", expectation{"20250620T134938", true}),
			Entry("now with 12h format returns datetime", "now.(YY/MM/DD hh:mm:ss AM)", expectation{"25/06/20 01:49:38 PM", true}),
			Entry("specific datetime returns datetime", "2011/03/23 09:10:11 AM.(YYYY/MM/DD hh:mm:ss AM) + 5Y", expectation{"2016/03/23 09:10:11 AM", true}),
		)

		DescribeTable("epoch",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("epoch", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now returns epoch", "now", expectation{toEpochStr(testLocal, 0), true}),
			Entry("now plus seconds returns correct epoch", "now + 30s", expectation{toEpochStr(testLocal, 30), true}),
			Entry("now minus seconds returns correct epoch", "now - 30s", expectation{toEpochStr(testLocal, -30), true}),
			Entry("specific epoch minus seconds returns expected epoch", "1750427679 - 100s", expectation{"1750427579", true}),
		)
	})

	Describe("UTC Time", func() {
		BeforeEach(func() {
			useLocal = false
		})

		DescribeTable("time",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("time", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now returns time", "now.utc", expectation{"17:49:38", true}),
			Entry("now plus seconds returns correct time", "now.utc + 10s", expectation{"17:49:48", true}),
			Entry("now minus seconds returns correct time", "now.utc - 10s", expectation{"17:49:28", true}),
			Entry("now plus minutes returns correct time", "now.utc + 10m", expectation{"17:59:38", true}),
			Entry("now plus hours returns correct time", "now.utc + 10h", expectation{"03:49:38", true}),
			Entry("now with 24h format returns time", "now.utc.(HHmmss)", expectation{"174938", true}),
			Entry("now with 12h format returns time", "now.utc.(hh.mm.ss-AM)", expectation{"05.49.38-PM", true}),
			Entry("specific time format returns time", "15:23:11.(TimeOnly) - 4m", expectation{"15:19:11", true}),
			Entry("invalid time modifier returns input", "now.utc +", expectation{"now.utc +", false}),
			Entry("invalid time modifier delta returns input", "now.utc + 6", expectation{"now.utc + 6", false}),
			Entry("invalid time modifier delta scope returns input", "now.utc + 6X", expectation{"now.utc + 6X", false}),
		)

		DescribeTable("date",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("date", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now.utc returns date", "now.utc", expectation{"2025-06-20", true}),
			Entry("now.utc plus days returns correct date", "now.utc + 10D", expectation{"2025-06-30", true}),
			Entry("now.utc plus months returns correct date", "now.utc + 5M", expectation{"2025-11-20", true}),
			Entry("now.utc minus months returns correct date", "now.utc - 5M", expectation{"2025-01-20", true}),
			Entry("now.utc plus weeks returns correct date", "now.utc + 2W", expectation{"2025-07-04", true}),
			Entry("now.utc plus years returns correct date", "now.utc + 10Y", expectation{"2035-06-20", true}),
			Entry("now.utc with format returns date", "now.utc.(YYYYMMDD)", expectation{"20250620", true}),
			Entry("now.utc with another format returns date", "now.utc.(YY/MM/DD)", expectation{"25/06/20", true}),
			Entry("specific date returns date", "2012/09/13.(YYYY/MM/DD) - 10D", expectation{"2012/09/03", true}),
		)

		DescribeTable("datetime",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("datetime", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now.utc returns datetime", "now.utc", expectation{"2025-06-20T17:49:38Z", true}),
			Entry("now.utc plus days returns correct datetime", "now.utc + 10D", expectation{"2025-06-30T17:49:38Z", true}),
			Entry("now.utc plus months returns correct datetime", "now.utc + 5M", expectation{"2025-11-20T17:49:38Z", true}),
			Entry("now.utc plus years returns correct datetime", "now.utc + 10Y", expectation{"2035-06-20T17:49:38Z", true}),
			Entry("now.utc minus years returns correct datetime", "now.utc - 10Y", expectation{"2015-06-20T17:49:38Z", true}),
			Entry("now.utc with 24h format returns datetime", "now.utc.(YYYYMMDDTHHmmss)", expectation{"20250620T174938", true}),
			Entry("now.utc with 12h format returns datetime", "now.utc.(YY/MM/DD hh:mm:ss AM)", expectation{"25/06/20 05:49:38 PM", true}),
			Entry("specific datetime returns datetime", "2011/03/23 09:10:11 AM.(YYYY/MM/DD hh:mm:ss AM) + 5Y", expectation{"2016/03/23 09:10:11 AM", true}),
		)

		DescribeTable("epoch",
			func(token string, expect expectation) {
				// Act
				actual, ok := resolver.Resolve("epoch", token)

				// Assert
				Expect(actual).To(Equal(expect.value))
				Expect(ok).To(Equal(expect.ok))
			},
			Entry("now.utc returns epoch", "now.utc", expectation{toEpochStr(testLocal, 0), true}),
			Entry("now.utc plus seconds returns correct epoch", "now.utc + 30s", expectation{toEpochStr(testLocal, 30), true}),
			Entry("now.utc minus seconds returns correct epoch", "now.utc - 30s", expectation{toEpochStr(testLocal, -30), true}),
			Entry("specific epoch minus seconds returns expected epoch", "1750427679 - 100s", expectation{"1750427579", true}),
		)
	})

	It("will ignore requests for invalid token name", func() {
		// Act
		actual, ok := resolver.Resolve("crumb", "now + 6D")

		// Assert
		Expect(ok).To(BeFalse())
		Expect(actual).To(Equal("now + 6D"))
	})
})
