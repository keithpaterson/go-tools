package env

import (
	"fmt"
	"os"
)

type Setup []applicator

func New() Setup {
	return Setup{}
}

func (a Setup) Set(key string, value interface{}) Setup {
	return append(a, &addOrUpdateEnv{key: key, value: fmt.Sprint(value)})
}

func (a Setup) Unset(key string) Setup {
	return append(a, &removeEnv{key: key})
}

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
