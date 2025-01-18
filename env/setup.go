package env

import (
	"fmt"
	"os"
)

// Provides a means for configuring, applying, and reverting a specific set of environment variables
//
// For example:
//
// If you want to ensure that your environment has the following values:
//
// ["XXX": "triple-x", "YYY" = "100", "XXX" = nil]
// (where 'nil' means it is unset)
//
//	myEnv := env.New().Set("XXX", "triple-x").Set("YYY", 100).Unset("ZZZ")
//	oldEnv := myEnv.Apply()
//	defer oldEnv.Apply() // reverts the changes made by myEnv.Apply()
type Setup []applicator

// Returns an empty environment Setup.
//
// Use Set/Unset methods to configure your setup.
//
// This is useful under test to ensure that your test environment conforms to a known setup,
// but there are other useful applications for this such as sanitizing your runtime environment
// in main().
func New() Setup {
	return Setup{}
}

// Adds an entry that will create or update an environment variable, ensuring that its value
// equals the value specified here
func (a Setup) Set(key string, value interface{}) Setup {
	return append(a, &addOrUpdateEnv{key: key, value: fmt.Sprint(value)})
}

// Adds an entry that will remove an environment variable (if it exists).
func (a Setup) Unset(key string) Setup {
	return append(a, &removeEnv{key: key})
}

// Applies the environment (adding, updating, and removing variables) and returns a new
// Setup that will revert your environment back to its starting state.
func (a Setup) Apply() Setup {
	inverters := Setup{}

	for _, applicator := range a {
		if i := applicator.apply(); i != nil {
			inverters = append(inverters, i)
		}
	}
	return inverters
}

type applicator interface {
	apply() applicator
}

type addOrUpdateEnv struct {
	key   string
	value string
}

type removeEnv struct {
	key string
}

func (a *addOrUpdateEnv) apply() applicator {
	var inverter applicator
	if value, ok := os.LookupEnv(a.key); ok {
		inverter = &addOrUpdateEnv{key: a.key, value: value}
	} else {
		inverter = &removeEnv{key: a.key}
	}
	os.Setenv(a.key, a.value)
	return inverter
}

func (a *removeEnv) apply() applicator {
	var inverter applicator = nil
	if value, ok := os.LookupEnv(a.key); ok {
		inverter = &addOrUpdateEnv{key: a.key, value: value}
	}
	os.Unsetenv(a.key)
	return inverter
}
