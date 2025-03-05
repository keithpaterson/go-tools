package env

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env Tag Parser", func() {
	var (
		noEnv      = New().Unset("ENV_VALUE").Unset("ENV_TEST_VALUE")
		baseEnv    = New().Set("ENV_VALUE", 1).Unset("ENV_TEST_VALUE")
		testEnv    = New().Set("ENV_VALUE", 1).Set("ENV_TEST_VALUE", 100)
		badBaseEnv = New().Set("ENV_VALUE", "BAD").Unset("ENV_TEST_VALUE")
		badTestEnv = New().Set("ENV_VALUE", 1).Set("ENV_TEST_VALUE", "EVIL")
	)

	Context("ResolveEnv", func() {
		DescribeTable("ensure ResolveEnv call ResolveEnvWithName with empty name",
			func(testEnv Setup, input int, expectedValue int, expectedErr error) {
				// Arrange
				type TestStruct struct {
					Value int `envp:"env=value,default=10"`
				}
				origEnv := testEnv.Apply()
				defer origEnv.Apply()

				// Act
				s := TestStruct{Value: input}
				err := ResolveEnv(&s)

				// Assert
				Expect(s.Value).To(Equal(expectedValue))
				if expectedErr != nil {
					Expect(err).To(MatchError(expectedErr))
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
			},
			Entry("no env returns default", noEnv, 0, 10, nil),
			Entry("base env returns base env value", baseEnv, 0, 1, nil),
			Entry("base + test env returns base env value", testEnv, 0, 1, nil),
			Entry("base + test env returns preset value", testEnv, 1000, 1000, nil),
			Entry("invalid base env value returns error", badBaseEnv, 0, 0, ErrEnvParseFailure),
			Entry("invalid test env value returns base env value", badTestEnv, 0, 1, nil),
		)
	})

	Context("ResolveEnvWithName", func() {
		It("will do nothing if input is nil", func() {
			// Arrange
			type TestStruct struct {
				Value int `envp:"env=value,default=10"`
			}
			origEnv := testEnv.Apply()
			defer origEnv.Apply()

			// Act
			var s TestStruct
			err := ResolveEnvWithName("test", nil)

			// Assert
			Expect(s.Value).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("will do nothing if input is not a pointer", func() {
			// Arrange
			type TestStruct struct {
				Value int `envp:"env=foo_value,default=10"`
			}
			origEnv := testEnv.Apply()
			defer origEnv.Apply()

			// Act
			var s TestStruct
			err := ResolveEnvWithName("test", s)

			// Assert
			Expect(s.Value).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("will report an error if struct has unsupported types", func() {
			// Arrange
			type TestStruct struct {
				Value complex64 `envp:"env=value,default=10"`
			}
			origEnv := testEnv.Apply()
			defer origEnv.Apply()

			// Act
			var s TestStruct
			err := ResolveEnvWithName("test", &s)

			// Assert
			Expect(s.Value).To(BeZero())
			Expect(err).To(MatchError(ErrEnvParseFailure))
		})

		It("will ignore unsettable fields", func() {
			// Arrange
			type TestStruct struct {
				Value int `envp:"env=value,default=10"`
				foo   int
			}
			origEnv := testEnv.Apply()
			defer origEnv.Apply()

			// Act
			s := TestStruct{}
			err := ResolveEnvWithName("test", &s)

			// Assert
			Expect(s.Value).To(Equal(100))
			Expect(s.foo).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("will assume 'res' for missing key", func() {
			// Arrange
			type TestStruct struct {
				Value int `envp:"value,default=10"`
				foo   int
			}
			origEnv := testEnv.Apply()
			defer origEnv.Apply()

			// Act
			s := TestStruct{}
			err := ResolveEnvWithName("test", &s)

			// Assert
			Expect(s.Value).To(Equal(100))
			Expect(s.foo).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("will use the tag's 'env' value to lookup the environment", func() {
			// Arrange
			type TestStruct struct {
				Value int `envp:"env=foo_value,default=10"`
			}
			origEnv := New().Set("ENV_FOO_VALUE", 1).Set("ENV_TEST_FOO_VALUE", 100).Apply()
			defer origEnv.Apply()

			// Act
			s := TestStruct{}
			err := ResolveEnvWithName("test", &s)

			// Assert
			Expect(s.Value).To(Equal(100))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("signed integers", func() {
			DescribeTable("will convert signed int",
				func(testEnv Setup, input int, expectedValue int, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value int `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					s := TestStruct{Value: input}
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, 0, 10, nil),
				Entry("base env returns base env value", baseEnv, 0, 1, nil),
				Entry("base + test env returns test env value", testEnv, 0, 100, nil),
				Entry("base + test env returns preset value", testEnv, 1000, 1000, nil),
				Entry("invalid base env value returns error", badBaseEnv, 0, 0, ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, 0, 0, ErrEnvParseFailure),
			)

			DescribeTable("will convert signed int16",
				func(testEnv Setup, expectedValue int16, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value int16 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, int16(10), nil),
				Entry("base env returns base env value", baseEnv, int16(1), nil),
				Entry("base + test env returns test env value", testEnv, int16(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, int16(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, int16(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert signed int32",
				func(testEnv Setup, expectedValue int32, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value int32 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, int32(10), nil),
				Entry("base env returns base env value", baseEnv, int32(1), nil),
				Entry("base + test env returns test env value", testEnv, int32(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, int32(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, int32(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert signed int64",
				func(testEnv Setup, expectedValue int64, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value int64 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, int64(10), nil),
				Entry("base env returns base env value", baseEnv, int64(1), nil),
				Entry("base + test env returns test env value", testEnv, int64(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, int64(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, int64(0), ErrEnvParseFailure),
			)
		})

		Context("unsigned integers", func() {
			DescribeTable("will convert unsigned int",
				func(testEnv Setup, expectedValue uint, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value uint `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, uint(10), nil),
				Entry("base env returns base env value", baseEnv, uint(1), nil),
				Entry("base + test env returns test env value", testEnv, uint(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, uint(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, uint(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert unsigned int16",
				func(testEnv Setup, expectedValue uint16, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value uint16 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, uint16(10), nil),
				Entry("base env returns base env value", baseEnv, uint16(1), nil),
				Entry("base + test env returns test env value", testEnv, uint16(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, uint16(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, uint16(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert unsigned int32",
				func(testEnv Setup, expectedValue uint32, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value uint32 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, uint32(10), nil),
				Entry("base env returns base env value", baseEnv, uint32(1), nil),
				Entry("base + test env returns test env value", testEnv, uint32(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, uint32(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, uint32(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert unsigned int64",
				func(testEnv Setup, expectedValue uint64, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value uint64 `envp:"env=value,default=10"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, uint64(10), nil),
				Entry("base env returns base env value", baseEnv, uint64(1), nil),
				Entry("base + test env returns test env value", testEnv, uint64(100), nil),
				Entry("invalid base env value returns error", badBaseEnv, uint64(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, uint64(0), ErrEnvParseFailure),
			)
		})

		Context("floating point", func() {
			var (
				baseEnv    = New().Set("ENV_VALUE", 1.1).Unset("ENV_TEST_VALUE")
				testEnv    = New().Set("ENV_VALUE", 1.1).Set("ENV_TEST_VALUE", 100.1)
				badBaseEnv = New().Set("ENV_VALUE", "BAD").Unset("ENV_TEST_VALUE")
				badTestEnv = New().Set("ENV_VALUE", 1.1).Set("ENV_TEST_VALUE", "EVIL")
			)
			DescribeTable("will convert float32",
				func(testEnv Setup, expectedValue float32, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value float32 `envp:"env=value,default=10.1"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(BeNumerically("==", expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, float32(10.1), nil),
				Entry("base env returns base env value", baseEnv, float32(1.1), nil),
				Entry("base + test env returns test env value", testEnv, float32(100.1), nil),
				Entry("invalid base env value returns error", badBaseEnv, float32(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, float32(0), ErrEnvParseFailure),
			)

			DescribeTable("will convert float64",
				func(testEnv Setup, expectedValue float64, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value float64 `envp:"env=value,default=10.1"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(BeNumerically("==", expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, float64(10.1), nil),
				Entry("base env returns base env value", baseEnv, float64(1.1), nil),
				Entry("base + test env returns test env value", testEnv, float64(100.1), nil),
				Entry("invalid base env value returns error", badBaseEnv, float64(0), ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, float64(0), ErrEnvParseFailure),
			)
		})

		Context("boolean", func() {
			var (
				baseEnv    = New().Set("ENV_VALUE", true).Unset("ENV_TEST_VALUE")
				testEnv    = New().Set("ENV_VALUE", true).Set("ENV_TEST_VALUE", false)
				badBaseEnv = New().Set("ENV_VALUE", "BAD").Unset("ENV_TEST_VALUE")
				badTestEnv = New().Set("ENV_VALUE", true).Set("ENV_TEST_VALUE", "EVIL")
			)
			DescribeTable("will convert bool",
				func(testEnv Setup, expectedValue bool, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value bool `envp:"env=value,default=false"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, false, nil),
				Entry("base env returns base env value", baseEnv, true, nil),
				Entry("base + test env test env value", testEnv, false, nil),
				Entry("invalid base env value returns error", badBaseEnv, false, ErrEnvParseFailure),
				Entry("invalid test env value returns error", badTestEnv, false, ErrEnvParseFailure),
			)
		})

		Context("strings", func() {
			var (
				baseEnv = New().Set("ENV_VALUE", "foo").Unset("ENV_TEST_VALUE")
				testEnv = New().Set("ENV_VALUE", "foo").Set("ENV_TEST_VALUE", "bar")
			)
			DescribeTable("will convert string",
				func(testEnv Setup, expectedValue string, expectedErr error) {
					// Arrange
					type TestStruct struct {
						Value string `envp:"env=value,default=eric"`
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, "eric", nil),
				Entry("base env returns base env value", baseEnv, "foo", nil),
				Entry("base + test env test env value", testEnv, "bar", nil),
			)
		})

		Context("nested structs", func() {
			var (
				noEnv   = New().Unset("ENV_VALUE").Unset("ENV_TEST_VALUE").Unset("ENV_FAB").Unset("ENV_TEST_FAB")
				baseEnv = New().Set("ENV_VALUE", 10).Unset("ENV_TEST_VALUE").Set("ENV_FAB", "tastic").Unset("ENV_TEST_FAB")
				testEnv = New().Set("ENV_VALUE", 10).Set("ENV_TEST_VALUE", 100).Unset("ENV_FAB").Set("ENV_TEST_FAB", "ulous")
			)
			DescribeTable("will convert nested structs",
				func(testEnv Setup, expectedValue int, expectedFab string, expectedErr error) {
					// Arrange
					type InnerStruct struct {
						Fab string `envp:"env=fab,default=four"`
					}
					type TestStruct struct {
						Value int `envp:"env=value,default=1"`
						Inner InnerStruct
					}
					origEnv := testEnv.Apply()
					defer origEnv.Apply()

					// Act
					var s TestStruct
					err := ResolveEnvWithName("test", &s)

					// Assert
					Expect(s.Value).To(Equal(expectedValue))
					Expect(s.Inner.Fab).To(Equal(expectedFab))
					if expectedErr != nil {
						Expect(err).To(MatchError(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
				},
				Entry("no env returns default", noEnv, 1, "four", nil),
				Entry("base env returns base env value", baseEnv, 10, "tastic", nil),
				Entry("base + test env test env value", testEnv, 100, "ulous", nil),
			)

			It("will follow pointers in structs", func() {
				// Arrange
				type InnerStruct struct {
					Fab string `envp:"env=fab,default=four"`
				}
				type TestStruct struct {
					Value int `envp:"env=value,default=1"`
					Inner *InnerStruct
				}
				origEnv := testEnv.Apply()
				defer origEnv.Apply()

				// Act
				s := TestStruct{}
				err := ResolveEnvWithName("test", &s)

				// Assert
				Expect(s.Value).To(Equal(100))
				Expect(s.Inner.Fab).To(Equal("ulous"))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
