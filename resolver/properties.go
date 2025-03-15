package resolver

import (
	"fmt"
)

type propertiesResolver struct {
	ResolverImpl
}

func NewPropertiesResolver() *propertiesResolver {
	return &propertiesResolver{}
}

func (r *propertiesResolver) Resolve(name string, token string) (string, bool) {
	if name != "prop" {
		return token, false
	}

	value, ok := r.root.config.Properties[token]
	if !ok {
		return token, false
	}
	strVal := fmt.Sprintf("%v", value)

	return r.resolveValue(strVal), true
}
