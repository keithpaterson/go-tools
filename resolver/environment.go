package resolver

import (
	"os"
)

type envResolver struct {
	ResolverImpl
}

func NewEnvResolver() *envResolver {
	return &envResolver{}
}

func (r *envResolver) Resolve(name string, token string) (string, bool) {
	if name != "env" {
		return token, false
	}

	envValue, ok := os.LookupEnv(token)
	if !ok {
		return token, false
	}

	return r.ResolveValue(envValue), true
}
