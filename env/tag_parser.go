package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	// Declares the prefix to use when parsing environment variables from the tags.
	//
	// You can change this in your main() fuction to whatever prefix you want to use.
	EnvTagPrefix = "ENV_"
)

var (
	ErrEnvParseFailure = errors.New("failed to parse env tags")
)

const (
	tagName     = "envp"
	propEnv     = "env"
	propDefault = "default"

	errMsgFmt = "%w: %s"
	errErrFmt = "%w: %w"
)

type tagProperties struct {
	envSuffix    string // suffix for the environment variable to test
	defaultValue string // default value as string
}

// Parser for converting 'envp' tags in structs and assigning values to fields, using a 'name' to compose
// the environment variable name to look up.
//
// Example:
//
// Given an environment with
//
//	["ENV_FOO_HOST_NAME": "monty", "ENV_FOO_PORT": nil]
//
// (where 'nil' means it is unset):
//
//	type MyFoo struct {
//	    Host string `envp:"host_name"`
//	    Port int    `envp:"port,default=1234"`
//	}
//
//	var foo MyFoo
//	ResolveWithName("foo", &foo)
//	fmt.Printf("%+v\n", foo)
//
// Will output:
//
//	{Host: "monty" Port: 1234}
//
// Note that this function prefers to use a 'name'-specific environment variable, and will fall-back to
// a non-named variable when the named one cannot be found.  For example, given this environment:
//
//	["ENV_HOST_NAME": "winston", "ENV_FOO_HOST_NAME": nil]
//
// the output would instead be:
//
//	{Host: "winston" Port: 1234}
func ResolveEnvWithName(name string, data interface{}) error {
	if data == nil {
		// nothing to do
		return nil
	}

	if reflect.TypeOf(data).Kind() != reflect.Pointer {
		// we can't usefully modify struct passed by value
		// we can only modify pointers to structs, so any non-pointer here is of no value
		return nil
	}

	parser := envpTagParser{name: name}
	return parser.resolve(reflect.ValueOf(data).Elem())
}

// allows parsing 'envp' tags without requiring a 'name' (so it would only look up 'base' environment values)
func ResolveEnv(data interface{}) error {
	return ResolveEnvWithName("", data)
}

type envpTagParser struct {
	name string
}

func (p *envpTagParser) resolve(value reflect.Value) error {
	for index := 0; index < value.Type().NumField(); index++ {
		field := value.Field(index)
		if !field.IsZero() || !field.CanSet() {
			// don't override an existing nonzero value
			// skip fields that cannot be set
			continue
		}
		fieldType := value.Type().Field(index)
		newValue := p.resolveFieldValue(fieldType)
		if err := p.setFieldValue(value.Field(index), newValue); err != nil {
			return err
		}
	}
	return nil
}

func (p *envpTagParser) resolveFieldValue(field reflect.StructField) string {
	tag := field.Tag.Get(tagName)
	properties := p.getTagProperties(tag)
	value := p.getEnvWithDefault(properties.envSuffix, properties.defaultValue)
	return value
}

func (p *envpTagParser) setFieldValue(field reflect.Value, value string) error {
	var err error

	switch field.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		var intValue int64
		if intValue, err = strconv.ParseInt(value, 0, 64); err != nil {
			return fmt.Errorf(errErrFmt, ErrEnvParseFailure, err)
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var uintValue uint64
		if uintValue, err = strconv.ParseUint(value, 0, 64); err != nil {
			return fmt.Errorf(errErrFmt, ErrEnvParseFailure, err)
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		var f64Value float64
		if f64Value, err = strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf(errErrFmt, ErrEnvParseFailure, err)
		}
		field.SetFloat(f64Value)
	case reflect.Bool:
		var boolValue bool
		if boolValue, err = strconv.ParseBool(value); err != nil {
			return fmt.Errorf(errErrFmt, ErrEnvParseFailure, err)
		}
		field.SetBool(boolValue)
	case reflect.String:
		field.SetString(value)
	case reflect.Struct:
		// recursive
		p.resolve(field)
	case reflect.Pointer, reflect.Interface:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		// recursive
		p.resolve(field.Elem())
	default:
		return fmt.Errorf("%w: unsupported field type '%s'", ErrEnvParseFailure, field.Kind().String())
	}

	return nil
}

func (p *envpTagParser) getTagProperties(tag string) tagProperties {
	properties := tagProperties{}
	params := strings.Split(tag, ",")
	for _, p := range params {
		trimmedParam := strings.TrimSpace(p)
		keyValue := strings.SplitN(trimmedParam, "=", 2)
		switch keyValue[0] {
		case trimmedParam: // specified without "="; assumes the 'env' prefix
			properties.envSuffix = trimmedParam
		case propEnv:
			properties.envSuffix = keyValue[1]
		case propDefault:
			properties.defaultValue = keyValue[1]
		}
	}
	return properties
}

func (p *envpTagParser) getEnvWithDefault(key string, defaultValue string) string {
	// e.g. ("bags", "bag_size") => "ENV_BAGS_BAG_SIZE" / "ENV_BAG_SIZE"
	// e.g. ("bytes", "foo_bar") => "ENV_BYTES_FOO_BAR" / "ENV_FOO_BAR"
	envNames := []string{
		// order matters here
		fmt.Sprintf("%s%s_%s", EnvTagPrefix, strings.ToUpper(p.name), strings.ToUpper(key)),
		fmt.Sprintf("%s%s", EnvTagPrefix, strings.ToUpper(key)),
	}

	result := ""
	for _, envName := range envNames {
		if value, found := os.LookupEnv(envName); found {
			result = value
			break
		}
	}
	if result == "" {
		result = defaultValue
	}
	return result
}
