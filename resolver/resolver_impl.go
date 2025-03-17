package resolver

// ResolverImpl is the base implementation for a resolver.
//
// Expected use:
//
//	type MyResolver struct {
//	  ResolverImpl
//	  ...
//	}
type ResolverImpl struct {
	root *rootResolver
}

// ResolveValue is Used by Resolver implementations to evaluate and possibly
// resolve the value string if it is a new token.  This allows token-chaining.
// e.g.
//
//	prop["foo"] = "my-${env:bar}"
//	env["bar"] = "bar"
//	token = "${prop:foo}"
//	value := r.NewResolver().WithStandardResolvers().Resolve(token)
//
//	1. prop resolver resolves "prop:foo" as "my-${env:bar}"
//	2. prop resolver calls `r.ResolveValue("my-${env:bar}")`
//	3. env resolver resolves "env:bar" as "bar"
//
//	final resolved value = "my-bar"
func (ri *ResolverImpl) ResolveValue(value string) string {
	return ri.root.Resolve(value)
}

func (ri *ResolverImpl) setRoot(root *rootResolver) {
	ri.root = root
}
