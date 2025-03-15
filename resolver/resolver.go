package resolver

import (
	"regexp"
	"strings"
)

const (
	// DefaultResolver determines which resolver is used when the 'type:' prefix is missing.
	// default is "prop"
	DefaultResolver = "prop"
)

var (
	regToken = regexp.MustCompile(`(?U)\${(.*)}`) // matches a "${type:value}" token, captures "type:token"
)

// A Resolver is used to convert property tokens in the form '${type:value}' into actual data.
type Resolver interface {
	Resolve(name string, token string) (string, bool)

	setRoot(root *rootResolver)
}

// base implementation for a resolver
type ResolverImpl struct {
	root *rootResolver
}

type rootResolver struct {
	config    *ResolverConfig
	resolvers resolversMap
}

func NewResolver(cfg *ResolverConfig) *rootResolver {
	return &rootResolver{config: cfg, resolvers: make(resolversMap)}
}

func (r *rootResolver) WithStandardResolvers() *rootResolver {
	return r.
		WithResolver("env", NewEnvResolver()).
		WithResolver("prop", NewPropertiesResolver()).
		WithDateTimeResolvers()
}

func (r *rootResolver) WithResolver(name string, resolver Resolver) *rootResolver {
	resolver.setRoot(r)
	r.resolvers[name] = resolver
	return r
}

func (r *rootResolver) WithDateTimeResolvers() *rootResolver {
	dtr := NewDateTimeResolver()
	return r.WithResolver("date", dtr).
		WithResolver("time", dtr).
		WithResolver("datetime", dtr).
		WithResolver("epoch", dtr)
}

type resolversMap map[string]Resolver

func (ri *ResolverImpl) resolveValue(value string) string {
	return ri.root.Resolve(value)
}

func (ri *ResolverImpl) setRoot(root *rootResolver) {
	ri.root = root
}

//
// rootResolver
//

func (r *rootResolver) Resolve(input string) string {
	tokens := regToken.FindAllString(input, -1)
	if tokens == nil {
		return input
	}

	replacements := r.resolveTokenValues(tokens)
	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(input)
}

func (r *rootResolver) resolveTokenValues(tokens []string) []string {
	result := make([]string, len(tokens)*2)
	for index, token := range tokens {
		matches := regToken.FindStringSubmatch(token)
		if len(matches) < 2 {
			continue // shouldn't happen but this appears to be an invalid token
		}
		if resolved, ok := r.resolveToken(matches[1]); ok {
			offset := index * 2
			result[offset] = token
			result[offset+1] = resolved
		}
	}
	return result
}

func (r *rootResolver) resolveToken(token string) (string, bool) {
	// expect "name:value"
	// e.g. "prop:foo", "env:MY_ENV_VAR", "datetime:now.(RSS3339)", "time:now.(TimeOnly) + 30s"
	var name, value string
	var ok bool
	if name, value, ok = strings.Cut(token, ":"); !ok {
		value = name
		name = DefaultResolver
	}
	name = strings.ToLower(name)

	resolver, ok := r.resolvers[name]
	if !ok {
		// log a warning?
		return token, false
	}
	return resolver.Resolve(name, value)
}
