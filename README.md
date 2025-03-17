# Personal Go tools

A few things I seem to keep repeating.

Dive into the packages for more detail; this is just a simple summary.

## package env

### Setup

Simplifies configuring (and un-configuring) your environment variables.
This is quite useful for testing, particularly with web services, so that you can
easily and reliably adjust the runtime environment during test without having
to deal with complicated scripting.

Example:

```
// build your environment  (add `$FOO` with value "foo", and remove `$BAR`)
var (
  myEnv := env.New().Set("FOO", "foo").Unset("BAR")
)

func myFunc() {
  // apply the environment and revert on exit
  origEnv := myEnv.Apply()
  defer origEnv.Apply()

  ... // your code goes here

  // when the function exits, the environment reverts to its original state
}
```

# package resolver

Provides a customizable text tokenizer.

Tokens are specified with the pattern `${type:value}` and are resolved via the
provided resolvers.

Default resolver types include:

- `date`, `time`, `datetime`, `epoch` for date formatting and manipulation
  (e.g. you can specify a date like "now + 1D" to get "now" plus a day)
- `env` looks up environment variables
- `prop` looks up values contained in a properties collection

You can implement custom resolvers using the `Resolver` interface and
the `ResolverImpl` base struct,
